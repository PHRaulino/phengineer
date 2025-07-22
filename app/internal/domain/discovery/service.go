package discovery

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PHRaulino/phengineer/internal/infrastructure/config"
)

// Service descoberta de arquivos
type Service struct{}

// NewService cria uma nova instância do service
func NewService() *Service {
	return &Service{}
}

// DiscoverFiles descobre arquivos baseado nas configurações
func (s *Service) DiscoverFiles(ctx context.Context) (*DiscoveryResult, error) {
	cfg := config.FromContext(ctx)
	analysisFilesPatternsPath := path.Join(cfg.Auto.ConfigDirPath, cfg.Settings.Analysis.AnalysisFilesPath)

	// Carrega apenas os patterns de análise (include)
	patterns, err := s.loadAnalysisPatterns(analysisFilesPatternsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load analysis patterns: %w", err)
	}

	// Converte tamanho máximo para bytes
	maxSizeBytes, err := s.parseFileSize(cfg.Settings.Analysis.FileLimits.MaxFileSize)
	if err != nil {
		return nil, fmt.Errorf("invalid max file size: %w", err)
	}

	// Descobre arquivos
	result, err := s.walkAndFilter(cfg.Auto.RootAppPath, patterns, maxSizeBytes, cfg.Settings.Analysis.FileLimits.MaxFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to discover files: %w", err)
	}

	result.MaxSizeBytes = maxSizeBytes
	result.Timestamp = time.Now().Unix()
	result.GitCommit = s.getCurrentGitCommit(cfg.Auto.RootAppPath)

	return result, nil
}

// DiscoverFilesWithLock descobre arquivos e salva/compara com lock
func (s *Service) DiscoverFilesWithLock(ctx context.Context) (*DiscoveryResult, *ChangedFilesResult, error) {
	cfg := config.FromContext(ctx)
	lockPath := filepath.Join(cfg.Auto.ConfigDirPath, "discovery-lock.json")

	// Descobre arquivos atuais
	currentResult, err := s.DiscoverFiles(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Carrega lock anterior (se existir)
	previousLock, err := s.loadDiscoveryLock(lockPath)
	if err != nil {
		// Se não conseguir carregar, cria novo lock
		err = s.saveDiscoveryLock(lockPath, currentResult)
		return currentResult, &ChangedFilesResult{
			NewFiles:   currentResult.Files,
			HasChanges: true,
		}, err
	}

	// Compara mudanças
	changedFiles := s.compareWithLock(currentResult, previousLock.LastDiscovery)

	// Salva novo lock
	err = s.saveDiscoveryLock(lockPath, currentResult)
	if err != nil {
		return currentResult, changedFiles, fmt.Errorf("failed to save lock: %w", err)
	}

	return currentResult, changedFiles, nil
}

// loadAnalysisPatterns carrega patterns tipados do arquivo de análise
func (s *Service) loadAnalysisPatterns(analysisPath string) ([]TypedPattern, error) {
	if analysisPath == "" {
		return []TypedPattern{}, nil
	}

	if _, err := os.Stat(analysisPath); os.IsNotExist(err) {
		return []TypedPattern{}, nil
	}

	file, err := os.Open(analysisPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open analysis file: %w", err)
	}
	defer file.Close()

	var patterns []TypedPattern
	currentType := PatternTypeCustom // Default
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignora linhas vazias
		if line == "" {
			continue
		}

		// Verifica se é uma seção tipada: # [CODE], # [CONFIG], etc.
		if strings.HasPrefix(line, "# [") && strings.Contains(line, "]") {
			sectionType := s.parseSectionType(line)
			if sectionType != "" {
				currentType = PatternType(sectionType)
				continue
			}
		}

		// Ignora outros comentários
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Verifica se é um pattern de negação (!)
		isNegated := strings.HasPrefix(line, "!")
		pattern := line
		if isNegated {
			pattern = strings.TrimPrefix(line, "!")
		}

		patterns = append(patterns, TypedPattern{
			Pattern:   pattern,
			Type:      currentType,
			IsNegated: isNegated,
		})
	}

	return patterns, scanner.Err()
}

// parseSectionType extrai o tipo da seção do comentário
func (s *Service) parseSectionType(line string) string {
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")

	if start == -1 || end == -1 || end <= start {
		return ""
	}

	sectionName := strings.ToLower(strings.TrimSpace(line[start+1 : end]))

	switch sectionName {
	case "code", "source", "src":
		return string(PatternTypeSnippet)
	case "config", "configuration", "docs", "documentation", "infra", "infrastructure":
		return string(PatternTypeCustom)
	default:
		return string(PatternTypeCustom)
	}
}

// shouldIncludeFile verifica se arquivo deve ser incluído e retorna o tipo
func (s *Service) shouldIncludeFile(relativePath string, patterns []TypedPattern) (bool, PatternType) {
	if len(patterns) == 0 {
		return false, PatternTypeCustom // Se não tem patterns, não inclui nada
	}

	var matchedType PatternType = PatternTypeCustom
	isIncluded := false

	// Primeira passada: verifica includes
	for _, typedPattern := range patterns {
		if !typedPattern.IsNegated && s.matchPattern(relativePath, typedPattern.Pattern) {
			isIncluded = true
			matchedType = typedPattern.Type
		}
	}

	// Se não foi incluído, retorna false
	if !isIncluded {
		return false, PatternTypeCustom
	}

	// Segunda passada: verifica excludes (patterns com !)
	for _, typedPattern := range patterns {
		if typedPattern.IsNegated && s.matchPattern(relativePath, typedPattern.Pattern) {
			return false, PatternTypeCustom // Excluído explicitamente
		}
	}

	return true, matchedType
}

func (s *Service) matchPattern(path, pattern string) bool {
	// Match exato do padrão
	if matched, _ := filepath.Match(pattern, path); matched {
		return true
	}

	// Match do nome do arquivo apenas
	if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
		return true
	}

	// Para patterns como **/*.go, remove ** e testa a extensão
	if strings.HasPrefix(pattern, "**/") {
		simplePattern := strings.TrimPrefix(pattern, "**/")
		if matched, _ := filepath.Match(simplePattern, filepath.Base(path)); matched {
			return true
		}
	}

	return false
}

// walkAndFilter percorre diretórios e filtra arquivos
func (s *Service) walkAndFilter(rootPath string, patterns []TypedPattern, maxSizeBytes, maxFiles int64) (*DiscoveryResult, error) {
	result := &DiscoveryResult{
		Files:          make([]File, 0),
		OversizedFiles: make([]File, 0),
	}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignora erros e continua
		}

		// Só processa arquivos
		if info.IsDir() {
			return nil
		}

		result.TotalFound++

		// Calcula path relativo
		relativePath, err := filepath.Rel(rootPath, path)
		if err != nil {
			relativePath = path
		}

		// Verifica se deve incluir baseado nos patterns
		shouldInclude, patternType := s.shouldIncludeFile(relativePath, patterns)
		if !shouldInclude {
			return nil
		}

		result.TotalFiltered++

		// Pega commit hash do arquivo
		commitHash := s.getFileCommitHash(rootPath, relativePath)

		// Cria objeto File
		file := File{
			Name:        info.Name(),
			Path:        filepath.Dir(relativePath),
			Size:        info.Size(),
			Type:        s.getFileType(info.Name()),
			PatternType: patternType,
			CommitHash:  commitHash,
			ModTime:     info.ModTime().Unix(),
		}

		// Verifica tamanho
		if info.Size() > maxSizeBytes {
			result.OversizedFiles = append(result.OversizedFiles, file)
			return nil
		}

		// Verifica limite de arquivos
		if int64(len(result.Files)) >= maxFiles {
			return nil // Continua mas não adiciona mais
		}

		result.Files = append(result.Files, file)
		return nil
	})

	return result, err
}

// parseFileSize converte string como "10MB" para bytes
func (s *Service) parseFileSize(sizeStr string) (int64, error) {
	if sizeStr == "" {
		return 10 << 20, nil // Default 10MB
	}

	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	var multiplier int64 = 1
	var numStr string

	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1 << 10
		numStr = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1 << 20
		numStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1 << 30
		numStr = strings.TrimSuffix(sizeStr, "GB")
	} else {
		numStr = sizeStr // Assume bytes
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}

	return int64(num * float64(multiplier)), nil
}

// getCurrentGitCommit pega o commit atual do repositório
func (s *Service) getCurrentGitCommit(rootPath string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = rootPath
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getFileCommitHash pega o hash do último commit que alterou o arquivo
func (s *Service) getFileCommitHash(rootPath, relativePath string) string {
	cmd := exec.Command("git", "log", "-1", "--format=%H", "--", relativePath)
	cmd.Dir = rootPath
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// loadDiscoveryLock carrega o arquivo de lock
func (s *Service) loadDiscoveryLock(lockPath string) (*DiscoveryLock, error) {
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("lock file not found")
	}

	data, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock DiscoveryLock
	err = json.Unmarshal(data, &lock)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lock, nil
}

// saveDiscoveryLock salva o arquivo de lock
func (s *Service) saveDiscoveryLock(lockPath string, result *DiscoveryResult) error {
	// Cria diretório se não existir
	dir := filepath.Dir(lockPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	now := time.Now().Unix()

	// Carrega lock existente para preservar CreatedAt
	existingLock, _ := s.loadDiscoveryLock(lockPath)
	createdAt := now
	if existingLock != nil {
		createdAt = existingLock.CreatedAt
	}

	lock := DiscoveryLock{
		LastDiscovery: result,
		CreatedAt:     createdAt,
		UpdatedAt:     now,
	}

	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock: %w", err)
	}

	err = os.WriteFile(lockPath, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}

// compareWithLock compara resultado atual com o lock anterior
func (s *Service) compareWithLock(current, previous *DiscoveryResult) *ChangedFilesResult {
	result := &ChangedFilesResult{
		ChangedFiles:   make([]File, 0),
		NewFiles:       make([]File, 0),
		DeletedFiles:   make([]File, 0),
		UnchangedFiles: make([]File, 0),
	}

	if previous == nil {
		result.NewFiles = current.Files
		result.HasChanges = len(current.Files) > 0
		return result
	}

	// Cria mapa dos arquivos anteriores para lookup rápido
	previousMap := make(map[string]File)
	for _, file := range previous.Files {
		key := filepath.Join(file.Path, file.Name)
		previousMap[key] = file
	}

	// Cria mapa dos arquivos atuais
	currentMap := make(map[string]File)
	for _, file := range current.Files {
		key := filepath.Join(file.Path, file.Name)
		currentMap[key] = file

		if prevFile, exists := previousMap[key]; exists {
			// Arquivo existe nos dois - verifica se mudou
			if prevFile.CommitHash != file.CommitHash || prevFile.ModTime != file.ModTime {
				result.ChangedFiles = append(result.ChangedFiles, file)
				result.HasChanges = true
			} else {
				result.UnchangedFiles = append(result.UnchangedFiles, file)
			}
		} else {
			// Arquivo novo
			result.NewFiles = append(result.NewFiles, file)
			result.HasChanges = true
		}
	}

	// Verifica arquivos deletados
	for _, prevFile := range previous.Files {
		key := filepath.Join(prevFile.Path, prevFile.Name)
		if _, exists := currentMap[key]; !exists {
			result.DeletedFiles = append(result.DeletedFiles, prevFile)
			result.HasChanges = true
		}
	}

	return result
}

func (s *Service) getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".go":
		return "go"
	case ".js", ".jsx":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".py":
		return "python"
	case ".java":
		return "java"
	case ".rs":
		return "rust"
	case ".php":
		return "php"
	case ".rb":
		return "ruby"
	case ".cs":
		return "csharp"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".xml":
		return "xml"
	case ".md":
		return "markdown"
	case ".txt":
		return "text"
	case ".toml":
		return "toml"
	case ".dockerfile":
		return "dockerfile"
	default:
		if strings.Contains(filename, "Dockerfile") {
			return "dockerfile"
		}
		return "other"
	}
}
