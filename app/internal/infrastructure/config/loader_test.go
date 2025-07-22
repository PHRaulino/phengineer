package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSaveSettingsToFile testa o salvamento de configurações em arquivo
func TestSaveSettingsToFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "settings-save-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Teste 1: Salvar configurações válidas
	t.Run("Save valid settings", func(t *testing.T) {
		settings := GetDefaultSettings(".phengineer")
		settings.Project.Type = "test-app"
		settings.Project.Language.Name = "javascript"

		settingsPath := filepath.Join(tempDir, "save-test.yml")
		err := SaveSettingsToFile(settings, settingsPath)
		if err != nil {
			t.Fatalf("Failed to save settings: %v", err)
		}

		// Verifica se arquivo foi criado
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			t.Error("Settings file was not created")
		}

		// Verifica se pode carregar novamente
		loadedSettings, err := LoadSettingsFromFile(settingsPath)
		if err != nil {
			t.Fatalf("Failed to load saved settings: %v", err)
		}

		if loadedSettings.Project.Type != "test-app" {
			t.Errorf("Expected project type 'test-app', got '%s'", loadedSettings.Project.Type)
		}
		if loadedSettings.Project.Language.Name != "javascript" {
			t.Errorf("Expected language 'javascript', got '%s'", loadedSettings.Project.Language.Name)
		}
	})

	// Teste 2: Criar diretório se não existir
	t.Run("Create directory if not exists", func(t *testing.T) {
		nestedDir := filepath.Join(tempDir, "nested", "deep", "config")
		settingsPath := filepath.Join(nestedDir, "settings.yml")

		settings := GetDefaultSettings(".phengineer")
		err := SaveSettingsToFile(settings, settingsPath)
		if err != nil {
			t.Fatalf("Failed to save settings with nested directory: %v", err)
		}

		// Verifica se diretório foi criado
		if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
			t.Error("Nested directory was not created")
		}

		// Verifica se arquivo foi criado
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			t.Error("Settings file was not created in nested directory")
		}
	})

	// Teste 3: Configurações inválidas
	t.Run("Save invalid settings", func(t *testing.T) {
		invalidSettings := &Settings{
			Project: Project{
				Type: "", // inválido
				Language: Language{
					Name:    "go",
					Version: "1.21",
				},
			},
		}

		settingsPath := filepath.Join(tempDir, "invalid-save.yml")
		err := SaveSettingsToFile(invalidSettings, settingsPath)
		if err == nil {
			t.Error("Expected error when saving invalid settings")
		}
		if !strings.Contains(err.Error(), "settings validation failed") {
			t.Errorf("Expected 'validation failed' error, got: %v", err)
		}
	})

	// Teste 4: Verificar conteúdo YAML gerado
	t.Run("Verify YAML content", func(t *testing.T) {
		settings := &Settings{
			Project: Project{
				Type: "cli",
				Language: Language{
					Name:    "rust",
					Version: "1.70",
				},
			},
			Analysis: Analysis{
				AnalysisFilesPath: "src/**/*.rs",
				IgnoreFilesPath:   "target/**",
				FileLimits: Limits{
					MaxFileSize: "20MB",
					MaxFiles:    2000,
				},
			},
		}

		settingsPath := filepath.Join(tempDir, "yaml-content.yml")
		err := SaveSettingsToFile(settings, settingsPath)
		if err != nil {
			t.Fatalf("Failed to save settings: %v", err)
		}

		// Lê o arquivo e verifica conteúdo
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			t.Fatalf("Failed to read saved file: %v", err)
		}

		content := string(data)
		expectedStrings := []string{
			"type: cli",
			"name: rust",
			"version: \"1.70\"",
			"files_include_path: src/**/*.rs",
			"max_file_size: 20MB",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(content, expected) {
				t.Errorf("Expected YAML to contain '%s', got:\n%s", expected, content)
			}
		}
	})
}

// TestLoadOrCreateSettings testa carregamento ou criação de configurações
func TestLoadOrCreateSettingsLoader(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "settings-load-create-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	settingsPath := filepath.Join(tempDir, "settings.yml")

	// Teste 1: Arquivo não existe - deve criar padrão
	t.Run("Create default when file not exists", func(t *testing.T) {
		settings, err := LoadOrCreateSettings(settingsPath, ".phengineer")
		if err != nil {
			t.Fatalf("Failed to load or create settings: %v", err)
		}

		// Verifica se retornou configurações padrão
		defaults := GetDefaultSettings(".phengineer")
		if settings.Project.Type != defaults.Project.Type {
			t.Errorf("Expected default project type '%s', got '%s'", defaults.Project.Type, settings.Project.Type)
		}

		// Verifica se arquivo foi criado
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			t.Error("Settings file should have been created")
		}
	})

	// Teste 3: Falha ao criar arquivo padrão
	t.Run("Fail to create default file", func(t *testing.T) {
		// Tenta criar em diretório que não pode ser criado (permissão)
		invalidPath := "/root/cannot-create/settings.yml"
		_, err := LoadOrCreateSettings(invalidPath, ".phengineer")
		if err == nil {
			t.Error("Expected error when cannot create default settings file")
		}
	})
}

// BenchmarkSaveSettingsToFile testa performance do salvamento
func BenchmarkSaveSettingsToFile(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "benchmark-save-*")
	defer os.RemoveAll(tempDir)

	settings := GetDefaultSettings(".phengineer")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		settingsPath := filepath.Join(tempDir, "bench-settings.yml")
		SaveSettingsToFile(settings, settingsPath)
	}
}
