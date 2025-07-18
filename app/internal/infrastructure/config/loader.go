package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadSettingsFromFile carrega as configurações de um arquivo YAML
func LoadSettingsFromFile(filePath string) (*Settings, error) {
	// Verifica se o arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("settings file not found: %s", filePath)
	}

	// Lê o arquivo
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	// Parse do YAML
	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings YAML: %w", err)
	}

	// Mescla com valores padrão
	settings.MergeWithDefaults()

	// Valida as configurações
	if err := settings.Validate(); err != nil {
		return nil, fmt.Errorf("settings validation failed: %w", err)
	}

	return &settings, nil
}

// SaveSettingsToFile salva as configurações em um arquivo YAML
func SaveSettingsToFile(settings *Settings, filePath string) error {
	// Valida as configurações antes de salvar
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("settings validation failed: %w", err)
	}

	// Cria o diretório se não existir
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Converte para YAML
	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings to YAML: %w", err)
	}

	// Escreve o arquivo
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// LoadOrCreateSettings carrega as configurações ou cria um arquivo padrão se não existir
func LoadOrCreateSettings(filePath string) (*Settings, error) {
	// Tenta carregar o arquivo existente
	if _, err := os.Stat(filePath); err == nil {
		return LoadSettingsFromFile(filePath)
	}

	// Arquivo não existe, cria um padrão
	fmt.Printf("Settings file not found, creating default: %s\n", filePath)

	defaultSettings := GetDefaultSettings()
	if err := SaveSettingsToFile(defaultSettings, filePath); err != nil {
		return nil, fmt.Errorf("failed to create default settings: %w", err)
	}

	return defaultSettings, nil
}

// GetSettingsTemplate retorna um template do arquivo settings.yml
func GetSettingsTemplate() string {
	return `# Project Analysis Configuration

project:
  type: "application"  # application, library, cli, etc.
  language:
    name: "go"
    version: "1.21"

analysis:
  files_include_path: "**/*"  # Glob pattern for files to include
  files_exclude_path: "node_modules/**,vendor/**,.git/**,*.log"  # Glob pattern for files to exclude
  file_limits:
    max_file_size: "10MB"  # Maximum file size to analyze
    max_files: 1000      # Maximum number of files to analyze
`
}

// WriteSettingsTemplate escreve um template de exemplo
func WriteSettingsTemplate(filePath string) error {
	template := GetSettingsTemplate()

	// Cria o diretório se não existir
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Escreve o template
	if err := os.WriteFile(filePath, []byte(template), 0o644); err != nil {
		return fmt.Errorf("failed to write settings template: %w", err)
	}

	return nil
}
