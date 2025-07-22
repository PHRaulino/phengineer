package discovery

// PatternType define o tipo de processamento para um pattern
type PatternType string

const (
	PatternTypeSnippet PatternType = "snippet"
	PatternTypeCustom  PatternType = "custom"
)

// File representa um arquivo descoberto
type File struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Size        int64       `json:"size"`
	Type        string      `json:"type"`
	PatternType PatternType `json:"pattern_type"`
	CommitHash  string      `json:"commit_hash"` // Hash do último commit que alterou o arquivo
	ModTime     int64       `json:"mod_time"`    // Timestamp de modificação
}
type DicoveredFiles struct {
	Snippets []File
	Docs     []File
}

// TypedPattern representa um pattern com seu tipo
type TypedPattern struct {
	Pattern   string      `json:"pattern"`
	Type      PatternType `json:"type"`
	IsNegated bool        `json:"is_negated"` // Se começa com !
}

// DiscoveryResult contém o resultado da descoberta de arquivos
type DiscoveryResult struct {
	TotalFound     int64  `json:"total_found"`
	TotalFiltered  int64  `json:"total_filtered"`
	MaxSizeBytes   int64  `json:"max_size_bytes"`
	Timestamp      int64  `json:"timestamp"`  // Quando foi executado
	GitCommit      string `json:"git_commit"` // Commit atual do repositório
	Files          []File `json:"files"`
	OversizedFiles []File `json:"oversized_files"`
}

// DiscoveryLock representa o arquivo de lock
type DiscoveryLock struct {
	LastDiscovery *DiscoveryResult `json:"last_discovery"`
	CreatedAt     int64            `json:"created_at"`
	UpdatedAt     int64            `json:"updated_at"`
}

// ChangedFilesResult resultado da comparação com o lock
type ChangedFilesResult struct {
	ChangedFiles   []File `json:"changed_files"`   // Arquivos que mudaram
	NewFiles       []File `json:"new_files"`       // Arquivos novos
	DeletedFiles   []File `json:"deleted_files"`   // Arquivos deletados
	UnchangedFiles []File `json:"unchanged_files"` // Arquivos inalterados
	HasChanges     bool   `json:"has_changes"`     // Se houve mudanças
}
