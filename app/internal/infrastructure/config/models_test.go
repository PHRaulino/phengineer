package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestGetDefaultSettings testa as configurações padrão
func TestGetDefaultSettings(t *testing.T) {
	defaults := GetDefaultSettings()

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
	if defaults.Project.Language.Version != "1.21" {
		t.Errorf("Expected default language version '1.21', got '%s'", defaults.Project.Language.Version)
	}

	// Verifica valores da Analysis
	if defaults.Analysis.FilesIncludePath == "" {
		t.Error("Default files include path should not be empty")
	}
	if defaults.Analysis.FilesIncludePath != ".phengineer/.includeFiles" {
		t.Errorf("Expected default include path '.phengineer/.includeFiles', got '%s'", defaults.Analysis.FilesIncludePath)
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
			settings: GetDefaultSettings(),
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

// TestSettingsMergeWithDefaults testa a mesclagem com valores padrão
func TestSettingsMergeWithDefaults(t *testing.T) {
	tests := []struct {
		name     string
		initial  *Settings
		expected *Settings
	}{
		{
			name:     "Completely empty settings",
			initial:  &Settings{},
			expected: GetDefaultSettings(),
		},
		{
			name: "Partial project settings",
			initial: &Settings{
				Project: Project{
					Type: "custom-type",
					// Language será preenchido com defaults
				},
			},
			expected: &Settings{
				Project: Project{
					Type: "custom-type",
					Language: Language{
						Name:    "go",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: ".phengineer/.includeFiles",
					FilesExcludePath: ".phengineer/.ignoreFiles",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    1000,
					},
				},
			},
		},
		{
			name: "Custom language settings",
			initial: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "rust",
						Version: "1.70",
					},
				},
			},
			expected: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "rust",
						Version: "1.70",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: ".phengineer/.includeFiles",
					FilesExcludePath: ".phengineer/.ignoreFiles",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    1000,
					},
				},
			},
		},
		{
			name: "Custom analysis settings",
			initial: &Settings{
				Analysis: Analysis{
					FilesIncludePath: "custom/path/**",
					FilesExcludePath: "custom/ignore/**",
					FileLimits: Limits{
						MaxFileSize: "20MB",
						MaxFiles:    2000,
					},
				},
			},
			expected: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "go",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "custom/path/**",
					FilesExcludePath: "custom/ignore/**",
					FileLimits: Limits{
						MaxFileSize: "20MB",
						MaxFiles:    2000,
					},
				},
			},
		},
		{
			name: "Partial limits settings",
			initial: &Settings{
				Analysis: Analysis{
					FileLimits: Limits{
						MaxFileSize: "50MB",
						// MaxFiles será preenchido com default
					},
				},
			},
			expected: &Settings{
				Project: Project{
					Type: "application",
					Language: Language{
						Name:    "go",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: ".phengineer/.includeFiles",
					FilesExcludePath: ".phengineer/.ignoreFiles",
					FileLimits: Limits{
						MaxFileSize: "50MB",
						MaxFiles:    1000,
					},
				},
			},
		},
		{
			name: "Settings with some empty fields",
			initial: &Settings{
				Project: Project{
					Type: "cli",
					Language: Language{
						Name:    "javascript",
						Version: "", // será preenchido com default
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "src/**/*.js",
					FilesExcludePath: "", // será preenchido com default
					FileLimits: Limits{
						MaxFileSize: "", // será preenchido com default
						MaxFiles:    500,
					},
				},
			},
			expected: &Settings{
				Project: Project{
					Type: "cli",
					Language: Language{
						Name:    "javascript",
						Version: "1.21",
					},
				},
				Analysis: Analysis{
					FilesIncludePath: "src/**/*.js",
					FilesExcludePath: ".phengineer/.ignoreFiles",
					FileLimits: Limits{
						MaxFileSize: "10MB",
						MaxFiles:    500,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cria uma cópia para não modificar o original
			settings := *tt.initial
			settings.MergeWithDefaults()

			// Compara cada campo
			if settings.Project.Type != tt.expected.Project.Type {
				t.Errorf("Expected project type '%s', got '%s'", tt.expected.Project.Type, settings.Project.Type)
			}
			if settings.Project.Language.Name != tt.expected.Project.Language.Name {
				t.Errorf("Expected language name '%s', got '%s'", tt.expected.Project.Language.Name, settings.Project.Language.Name)
			}
			if settings.Project.Language.Version != tt.expected.Project.Language.Version {
				t.Errorf("Expected language version '%s', got '%s'", tt.expected.Project.Language.Version, settings.Project.Language.Version)
			}
			if settings.Analysis.FilesIncludePath != tt.expected.Analysis.FilesIncludePath {
				t.Errorf("Expected include path '%s', got '%s'", tt.expected.Analysis.FilesIncludePath, settings.Analysis.FilesIncludePath)
			}
			if settings.Analysis.FilesExcludePath != tt.expected.Analysis.FilesExcludePath {
				t.Errorf("Expected exclude path '%s', got '%s'", tt.expected.Analysis.FilesExcludePath, settings.Analysis.FilesExcludePath)
			}
			if settings.Analysis.FileLimits.MaxFileSize != tt.expected.Analysis.FileLimits.MaxFileSize {
				t.Errorf("Expected max file size '%s', got '%s'", tt.expected.Analysis.FileLimits.MaxFileSize, settings.Analysis.FileLimits.MaxFileSize)
			}
			if settings.Analysis.FileLimits.MaxFiles != tt.expected.Analysis.FileLimits.MaxFiles {
				t.Errorf("Expected max files %d, got %d", tt.expected.Analysis.FileLimits.MaxFiles, settings.Analysis.FileLimits.MaxFiles)
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
	settings := GetDefaultSettings()
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
	defaults := GetDefaultSettings()

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

// TestMergeWithDefaultsImmutability testa se MergeWithDefaults não modifica defaults
func TestMergeWithDefaultsImmutability(t *testing.T) {
	// Pega defaults originais
	originalDefaults := GetDefaultSettings()
	originalType := originalDefaults.Project.Type

	// Cria settings que vão usar merge
	settings := &Settings{
		Project: Project{
			Type: "custom",
		},
	}

	// Faz merge
	settings.MergeWithDefaults()

	// Verifica se defaults originais não foram modificados
	newDefaults := GetDefaultSettings()
	if newDefaults.Project.Type != originalType {
		t.Error("MergeWithDefaults should not modify the original defaults")
	}

	// Verifica se settings tem o valor customizado
	if settings.Project.Type != "custom" {
		t.Errorf("Expected custom type to be preserved, got '%s'", settings.Project.Type)
	}
}

// Benchmarks

// BenchmarkGetDefaultSettings testa performance da criação de defaults
func BenchmarkGetDefaultSettings(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDefaultSettings()
	}
}

// BenchmarkSettingsValidate testa performance da validação
func BenchmarkSettingsValidate(b *testing.B) {
	settings := GetDefaultSettings()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		settings.Validate()
	}
}

// BenchmarkMergeWithDefaults testa performance do merge
func BenchmarkMergeWithDefaults(b *testing.B) {
	_ = &Settings{
		Project: Project{
			Type: "custom",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Cria nova instância para cada iteração
		testSettings := &Settings{
			Project: Project{
				Type: "custom",
			},
		}
		testSettings.MergeWithDefaults()
	}
}

// BenchmarkYAMLMarshal testa performance da serialização YAML
func BenchmarkYAMLMarshal(b *testing.B) {
	settings := GetDefaultSettings()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		yaml.Marshal(settings)
	}
}
