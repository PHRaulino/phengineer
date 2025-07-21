package config

import (
	"fmt"
	"path/filepath"
)

// Settings representa a estrutura do arquivo settings.yml
type Settings struct {
	Project  Project  `yaml:"project"`
	Analysis Analysis `yaml:"analysis"`
}

// Project representa as configurações do projeto
type Project struct {
	Type     string   `yaml:"type"`
	Language Language `yaml:"language"`
}

// Language representa as configurações da linguagem
type Language struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// Analysis representa as configurações de análise
type Analysis struct {
	FilesIncludePath string `yaml:"files_include_path"`
	FilesExcludePath string `yaml:"files_exclude_path"`
	FileLimits       Limits `yaml:"file_limits"`
}

// Limits representa os limites de arquivos
type Limits struct {
	MaxFileSize string `yaml:"max_file_size"`
	MaxFiles    int64  `yaml:"max_files"`
}

// AutoConfig representa as configurações automáticas coletadas do ambiente
type AutoConfig struct {
	AppName       string // Nome do repositório
	ConfigDirPath string // Caminho da pasta de configs
	RemoteURL     string // URL do remote sem .git
}

// Config representa a configuração completa da aplicação
type Config struct {
	Settings   *Settings
	Auto       *AutoConfig
	ConfigPath string
}

// GetDefaultSettings retorna as configurações padrão
func GetDefaultSettings(configFolderName string) *Settings {
	return &Settings{
		Project: Project{
			Type: "application",
			Language: Language{
				Name:    "go",
				Version: "1.24",
			},
		},
		Analysis: Analysis{
			FilesIncludePath: filepath.Join(configFolderName, ".analyzeFiles"),
			FilesExcludePath: filepath.Join(configFolderName, ".ignoreFiles"),
			FileLimits: Limits{
				MaxFileSize: "10MB",
				MaxFiles:    1000,
			},
		},
	}
}

// Validate valida as configurações
func (s *Settings) Validate() error {
	// Valida Project
	if s.Project.Type == "" {
		return fmt.Errorf("project.type is required")
	}

	if s.Project.Language.Name == "" {
		return fmt.Errorf("project.language.name is required")
	}

	if s.Project.Language.Version == "" {
		return fmt.Errorf("project.language.version is required")
	}

	// Valida Analysis
	if s.Analysis.FilesIncludePath == "" {
		return fmt.Errorf("analysis.files_include_path is required")
	}

	if s.Analysis.FileLimits.MaxFileSize == "" {
		return fmt.Errorf("analysis.file_limits.max_file_size is required")
	}

	if s.Analysis.FileLimits.MaxFiles == 0 {
		return fmt.Errorf("analysis.file_limits.max_files is required")
	}

	return nil
}