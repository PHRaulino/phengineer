package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestGetDefaultSettings testa as configurações padrão
func TestGetDefaultSettings(t *testing.T) {
	defaults := GetDefaultSettings(".phengineer")

	// Verifica se não é nil
	if defaults == nil {
		t.Fatal("GetDefaultSettings should not return nil")
	}

	// Verifica valores do Project
	if defaults.Project.Type == "" {
		t.Error("Default project type should not be empty")
	}
	if defaults.Project.Type != "application" {
		t.Errorf("Expected default project type 'application', got '%s'", defaults.Project.Type)
	}

	// Verifica valores da Language
	if defaults.Project.Language.Name == "" {
		t.Error("Default language name should not be empty")
	}
	if defaults.Project.Language.Name != "go" {
		t.Errorf("Expected default language 'go', got '%s'", defaults.Project.Language.Name)
	}
	if defaults.Project.Language.Version == "" {
		t.Error("Default language version should not be empty")
	}
	if defaults.Project.Language.Version != "1.24" {
		t.Errorf("Expected default language version '1.24', got '%s'", defaults.Project.Language.Version)
	}

	// Verifica valores da Analysis
	if defaults.Analysis.FilesIncludePath == "" {
		t.Error("Default files include path should not be empty")
	}
	if defaults.Analysis.FilesIncludePath != ".phengineer/.analyzeFiles" {
		t.Errorf("Expected default analyze path '.phengineer/.analyzeFiles', got '%s'", defaults.Analysis.FilesIncludePath)
	}

	if defaults.Analysis.FilesExcludePath == "" {
		t.Error("Default files exclude path should not be empty")
	}
	if defaults.Analysis.FilesExcludePath != ".phengineer/.ignoreFiles" {
		t.Errorf("Expected default exclude path '.phengineer/.ignoreFiles', got '%s'", defaults.Analysis.FilesExcludePath)
	}

	// Verifica valores dos Limits
	if defaults.Analysis.FileLimits.MaxFileSize == "" {
		t.Error("Default max file size should not be empty")
	}
	if defaults.Analysis.FileLimits.MaxFileSize != "10MB" {
		t.Errorf("Expected default max file size '10MB', got '%s'", defaults.Analysis.FileLimits.MaxFileSize)
	}

	if defaults.Analysis.FileLimits.MaxFiles == 0 {
		t.Error("Default max files should not be zero")
	}
	if defaults.Analysis.FileLimits.MaxFiles != 1000 {
		t.Errorf("Expected default max files 1000, got %d", defaults.Analysis.FileLimits.MaxFiles)
	}
}

// TestSettingsValidate testa a validação das configurações
func TestSettingsValidate(t *testing.T) {
	tests := []struct {
		name     string
		settings *Settings
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid default settings",
			settings: GetDefaultSettings(".phengineer"),
			wantErr:  false,
		},
		{
			name: "Valid custom settings",
			settings: &Settings{
				Project: Project{
					Type: "library",
					Language: Language{
						Name:    "python",
						Version: "3.11",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "src/**/*.py",
					FilesExcludePath: "__pycache__/**",
					FileLimits: Limits{
						MaxFileSize: "5MB",
						MaxFiles:    500,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing project type",
			settings: &Settings{
				Project: Project{
					Type: "",
					Language: Language{
						Name:    "go",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "**/*.go",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    1000,
					},
				},
			},
			wantErr: true,
			errMsg:  "project.type is required",
		},
		{
			name: "Missing language name",
			settings: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "**/*.go",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    1000,
					},
				},
			},
			wantErr: true,
			errMsg:  "project.language.name is required",
		},
		{
			name: "Missing language version",
			settings: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "go",
						Version: "",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "**/*.go",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    1000,
					},
				},
			},
			wantErr: true,
			errMsg:  "project.language.version is required",
		},
		{
			name: "Missing files include path",
			settings: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "go",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    1000,
					},
				},
			},
			wantErr: true,
			errMsg:  "analysis.files_include_path is required",
		},
		{
			name: "Missing max file size",
			settings: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "go",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "**/*.go",
					FileLimits: Limits{
						MaxFileSize: "",
						MaxFiles:    1000,
					},
				},
			},
			wantErr: true,
			errMsg:  "analysis.file_limits.max_file_size is required",
		},
		{
			name: "Zero max files",
			settings: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "go",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "**/*.go",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    0,
					},
				},
			},
			wantErr: true,
			errMsg:  "analysis.file_limits.max_files is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.Validate()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected validation error, got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, got: %v", err)
				}
			}
		})
	}
}


// TestSettingsYAMLSerialization testa serialização/deserialização YAML
func TestSettingsYAMLSerialization(t *testing.T) {
	// Configurações originais
	original := &Settings{
		Project: Project{
			Type: "library",
			Language: Language{
				Name:    "rust",
				Version: "1.70",
			},
		},
		Analysis: Analysis{
			FilesIncludePath: "src/**/*.rs",
			FilesExcludePath: "target/**",
			FileLimits: Limits{
				MaxFileSize: "15MB",
				MaxFiles:    1500,
			},
		},
	}

	// Serializa para YAML
	yamlData, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal settings to YAML: %v", err)
	}

	// Verifica se YAML contém campos esperados
	yamlStr := string(yamlData)
	expectedFields := []string{
		"project:",
		"type: library",
		"language:",
		"name: rust",
		"version: \"1.70\"",
		"analysis:",
		"files_include_path: src/**/*.rs",
		"files_exclude_path: target/**",
		"file_limits:",
		"max_file_size: 15MB",
		"max_files: 1500",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("YAML should contain '%s', got:\n%s", field, yamlStr)
		}
	}

	// Deserializa de volta
	var deserialized Settings
	err = yaml.Unmarshal(yamlData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML to settings: %v", err)
	}

	// Compara valores
	if deserialized.Project.Type != original.Project.Type {
		t.Errorf("Expected project type '%s', got '%s'", original.Project.Type, deserialized.Project.Type)
	}
	if deserialized.Project.Language.Name != original.Project.Language.Name {
		t.Errorf("Expected language name '%s', got '%s'", original.Project.Language.Name, deserialized.Project.Language.Name)
	}
	if deserialized.Project.Language.Version != original.Project.Language.Version {
		t.Errorf("Expected language version '%s', got '%s'", original.Project.Language.Version, deserialized.Project.Language.Version)
	}
	if deserialized.Analysis.FilesIncludePath != original.Analysis.FilesIncludePath {
		t.Errorf("Expected include path '%s', got '%s'", original.Analysis.FilesIncludePath, deserialized.Analysis.FilesIncludePath)
	}
	if deserialized.Analysis.FileLimits.MaxFiles != original.Analysis.FileLimits.MaxFiles {
		t.Errorf("Expected max files %d, got %d", original.Analysis.FileLimits.MaxFiles, deserialized.Analysis.FileLimits.MaxFiles)
	}
}

// TestConfigStructs testa estruturas auxiliares
func TestConfigStructs(t *testing.T) {
	// Teste AutoConfig
	autoConfig := &AutoConfig{
		AppName:       "test-app",
		ConfigDirPath: "/path/to/config",
		RemoteURL:     "https://github.com/user/test-app",
	}

	if autoConfig.AppName != "test-app" {
		t.Errorf("Expected app name 'test-app', got '%s'", autoConfig.AppName)
	}
	if autoConfig.ConfigDirPath != "/path/to/config" {
		t.Errorf("Expected config dir '/path/to/config', got '%s'", autoConfig.ConfigDirPath)
	}
	if autoConfig.RemoteURL != "https://github.com/user/test-app" {
		t.Errorf("Expected remote URL 'https://github.com/user/test-app', got '%s'", autoConfig.RemoteURL)
	}

	// Teste Config
	settings := GetDefaultSettings(".phengineer")
	config := &Config{
		Settings:   settings,
		Auto:       autoConfig,
		ConfigPath: "/path/to/config",
	}

	if config.Settings != settings {
		t.Error("Config should contain the provided settings")
	}
	if config.Auto != autoConfig {
		t.Error("Config should contain the provided auto config")
	}
	if config.ConfigPath != "/path/to/config" {
		t.Errorf("Expected config path '/path/to/config', got '%s'", config.ConfigPath)
	}
}

// TestDefaultSettingsValidity testa se configurações padrão são válidas
func TestDefaultSettingsValidity(t *testing.T) {
	defaults := GetDefaultSettings(".phengineer")

	// Configurações padrão devem passar na validação
	err := defaults.Validate()
	if err != nil {
		t.Errorf("Default settings should be valid, got error: %v", err)
	}

	// Configurações padrão devem ser serializáveis
	_, err = yaml.Marshal(defaults)
	if err != nil {
		t.Errorf("Default settings should be serializable to YAML, got error: %v", err)
	}
}

// Benchmarks

// BenchmarkGetDefaultSettings testa performance da criação de defaults
func BenchmarkGetDefaultSettings(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDefaultSettings(".phengineer")
	}
}

// BenchmarkSettingsValidate testa performance da validação
func BenchmarkSettingsValidate(b *testing.B) {
	settings := GetDefaultSettings(".phengineer")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		settings.Validate()
	}
}

// BenchmarkYAMLMarshal testa performance da serialização YAML
func BenchmarkYAMLMarshal(b *testing.B) {
	settings := GetDefaultSettings(".phengineer")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		yaml.Marshal(settings)
	}
}
