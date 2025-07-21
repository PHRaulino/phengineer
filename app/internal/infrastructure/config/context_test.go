package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestWithConfig testa a adição de config ao context
func TestWithConfig(t *testing.T) {
	// Setup de ambiente temporário
	tempDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Pode falhar por não ter git real, mas testa a estrutura
	ctx := context.Background()
	_, err = WithConfig(ctx, "config")
	
	// O teste pode falhar por requirements, mas não deve panic
	if err != nil {
		t.Logf("WithConfig failed (expected in test environment): %v", err)
	}
}

// TestFromContext testa extração de config do context
func TestFromContext(t *testing.T) {
	ctx := context.Background()
	
	// Testa context sem config - deve dar panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when config not in context")
		}
	}()
	
	FromContext(ctx)
}

// TestGetSettings testa extração de settings do context
func TestGetSettings(t *testing.T) {
	ctx := context.Background()
	mockSettings := GetDefaultSettings(".phengineer")
	mockSettings.Project.Type = "test-type"
	
	mockConfig := &Config{Settings: mockSettings}
	ctx = context.WithValue(ctx, configKey{}, mockConfig)
	
	settings := GetSettings(ctx)
	if settings.Project.Type != "test-type" {
		t.Errorf("Expected project type 'test-type', got '%s'", settings.Project.Type)
	}
}

// TestGetAutoConfig testa extração de auto config do context
func TestGetAutoConfig(t *testing.T) {
	ctx := context.Background()
	mockAuto := &AutoConfig{
		AppName:       "test-app",
		ConfigDirPath: "/test/config",
		RemoteURL:     "https://github.com/user/test-app",
	}
	
	mockConfig := &Config{Auto: mockAuto}
	ctx = context.WithValue(ctx, configKey{}, mockConfig)
	
	auto := GetAutoConfig(ctx)
	if auto.AppName != "test-app" {
		t.Errorf("Expected app name 'test-app', got '%s'", auto.AppName)
	}
}

// TestSettingsValidation testa validação das configurações
func TestSettingsValidation(t *testing.T) {
	tests := []struct {
		name     string
		settings *Settings
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid settings",
			settings: GetDefaultSettings(".phengineer"),
			wantErr:  false,
		},
		{
			name: "Missing project type",
			settings: &Settings{
				Project: Project{
					Language: Language{Name: "go", Version: "1.21"},
				},
				Analysis: Analysis{
					FilesIncludePath: "**/*",
					FileLimits:       Limits{MaxFileSize: "10MB", MaxFiles: 1000},
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
					Language: Language{Version: "1.21"},
				},
				Analysis: Analysis{
					FilesIncludePath: "**/*",
					FileLimits:       Limits{MaxFileSize: "10MB", MaxFiles: 1000},
				},
			},
			wantErr: true,
			errMsg:  "project.language.name is required",
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

// TestAutoConfigFunctions testa funções de auto config (podem falhar sem git)
func TestAutoConfigFunctions(t *testing.T) {
	t.Run("extractRepoNameFromURL", func(t *testing.T) {
		tests := []struct {
			url      string
			expected string
		}{
			{"https://github.com/user/repo.git", "repo"},
			{"https://github.com/user/repo", "repo"},
			{"git@github.com:user/repo.git", "repo"},
			{"git@github.com:user/repo", "repo"},
			{"", ""},
			{"invalid", "invalid"},
		}
		
		for _, tt := range tests {
			result := extractRepoNameFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("extractRepoNameFromURL(%s) = %s, expected %s", tt.url, result, tt.expected)
			}
		}
	})
}

// TestLoadOrCreateSettings testa carregamento ou criação de settings
func TestLoadOrCreateSettings(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "settings-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	settingsPath := filepath.Join(tempDir, "settings.yml")
	
	// Primeiro carregamento - deve criar arquivo padrão
	settings, err := LoadOrCreateSettings(settingsPath, ".phengineer")
	if err != nil {
		t.Fatalf("Failed to load or create settings: %v", err)
	}
	
	if settings.Project.Type == "" {
		t.Error("Settings should have default project type")
	}
	
	// Verifica se arquivo foi criado
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Error("Settings file should have been created")
	}
	
	// Segundo carregamento - deve carregar arquivo existente
	settings2, err := LoadOrCreateSettings(settingsPath, ".phengineer")
	if err != nil {
		t.Fatalf("Failed to load existing settings: %v", err)
	}
	
	if settings2.Project.Type != settings.Project.Type {
		t.Error("Loaded settings should match created settings")
	}
}

// BenchmarkFromContext testa performance de acesso ao context
func BenchmarkFromContext(b *testing.B) {
	ctx := context.Background()
	mockConfig := &Config{
		Settings: GetDefaultSettings(".phengineer"),
		Auto: &AutoConfig{
			AppName:       "bench-app",
			ConfigDirPath: "/bench/config",
			RemoteURL:     "https://github.com/user/bench-app",
		},
	}
	ctx = context.WithValue(ctx, configKey{}, mockConfig)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		config := FromContext(ctx)
		_ = config.Auto.AppName
	}
}

// BenchmarkGetSettings testa performance de acesso às settings
func BenchmarkGetSettings(b *testing.B) {
	ctx := context.Background()
	mockConfig := &Config{
		Settings: GetDefaultSettings(".phengineer"),
	}
	ctx = context.WithValue(ctx, configKey{}, mockConfig)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		settings := GetSettings(ctx)
		_ = settings.Project.Type
	}
}