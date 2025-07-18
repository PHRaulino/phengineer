package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
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
		settings := GetDefaultSettings()
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

		settings := GetDefaultSettings()
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
				FilesIncludePath: "src/**/*.rs",
				FilesExcludePath: "target/**",
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
		settings, err := LoadOrCreateSettings(settingsPath)
		if err != nil {
			t.Fatalf("Failed to load or create settings: %v", err)
		}

		// Verifica se retornou configurações padrão
		defaults := GetDefaultSettings()
		if settings.Project.Type != defaults.Project.Type {
			t.Errorf("Expected default project type '%s', got '%s'", defaults.Project.Type, settings.Project.Type)
		}

		// Verifica se arquivo foi criado
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			t.Error("Settings file should have been created")
		}
	})

	// Teste 2: Arquivo existe - deve carregar existente
	t.Run("Load existing file", func(t *testing.T) {
		// Modifica o arquivo existente
		existingSettings, _ := LoadSettingsFromFile(settingsPath)
		existingSettings.Project.Type = "modified"
		SaveSettingsToFile(existingSettings, settingsPath)

		// Carrega novamente
		loadedSettings, err := LoadOrCreateSettings(settingsPath)
		if err != nil {
			t.Fatalf("Failed to load existing settings: %v", err)
		}

		if loadedSettings.Project.Type != "modified" {
			t.Errorf("Expected to load modified settings with type 'modified', got '%s'", loadedSettings.Project.Type)
		}
	})

	// Teste 3: Falha ao criar arquivo padrão
	t.Run("Fail to create default file", func(t *testing.T) {
		// Tenta criar em diretório que não pode ser criado (permissão)
		invalidPath := "/root/cannot-create/settings.yml"
		_, err := LoadOrCreateSettings(invalidPath)
		if err == nil {
			t.Error("Expected error when cannot create default settings file")
		}
	})
}

// TestGetSettingsTemplate testa o template de configurações
func TestGetSettingsTemplate(t *testing.T) {
	template := GetSettingsTemplate()

	// Teste 1: Template não deve estar vazio
	if template == "" {
		t.Error("Settings template should not be empty")
	}

	// Teste 2: Deve conter seções principais
	expectedSections := []string{
		"project:",
		"analysis:",
		"type:",
		"language:",
		"files_include_path:",
		"files_exclude_path:",
		"file_limits:",
		"max_file_size:",
		"max_files:",
	}

	for _, section := range expectedSections {
		if !strings.Contains(template, section) {
			t.Errorf("Template should contain section '%s'", section)
		}
	}

	// Teste 3: Deve ser YAML válido
	var testParse interface{}
	err := yaml.Unmarshal([]byte(template), &testParse)
	if err != nil {
		t.Errorf("Template should be valid YAML: %v", err)
	}

	// Teste 4: Deve conter comentários explicativos
	expectedComments := []string{
		"# Project Analysis Configuration",
		"# application, library, cli, etc.",
		"# Glob pattern for files to include",
		"# Maximum file size to analyze",
	}

	for _, comment := range expectedComments {
		if !strings.Contains(template, comment) {
			t.Errorf("Template should contain comment '%s'", comment)
		}
	}
}

// TestWriteSettingsTemplate testa a escrita do template
func TestWriteSettingsTemplate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "template-write-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Teste 1: Escrever template em arquivo
	t.Run("Write template to file", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "template.yml")
		err := WriteSettingsTemplate(templatePath)
		if err != nil {
			t.Fatalf("Failed to write settings template: %v", err)
		}

		// Verifica se arquivo foi criado
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			t.Error("Template file was not created")
		}

		// Verifica se conteúdo está correto
		data, err := os.ReadFile(templatePath)
		if err != nil {
			t.Fatalf("Failed to read template file: %v", err)
		}

		content := string(data)
		expectedTemplate := GetSettingsTemplate()
		if content != expectedTemplate {
			t.Error("Written template content does not match expected template")
		}
	})

	// Teste 2: Criar diretório se não existir
	t.Run("Create directory for template", func(t *testing.T) {
		nestedPath := filepath.Join(tempDir, "deep", "nested", "template.yml")
		err := WriteSettingsTemplate(nestedPath)
		if err != nil {
			t.Fatalf("Failed to write template in nested directory: %v", err)
		}

		// Verifica se diretório foi criado
		if _, err := os.Stat(filepath.Dir(nestedPath)); os.IsNotExist(err) {
			t.Error("Nested directory was not created for template")
		}

		// Verifica se arquivo foi criado
		if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
			t.Error("Template file was not created in nested directory")
		}
	})

	// Teste 3: Template deve ser carregável como configurações válidas
	t.Run("Template should be loadable as valid settings", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "loadable-template.yml")
		err := WriteSettingsTemplate(templatePath)
		if err != nil {
			t.Fatalf("Failed to write template: %v", err)
		}

		// Tenta carregar o template como configurações
		settings, err := LoadSettingsFromFile(templatePath)
		if err != nil {
			t.Fatalf("Template should be loadable as valid settings: %v", err)
		}

		// Verifica se tem valores padrão esperados
		if settings.Project.Type == "" {
			t.Error("Template should have project type")
		}
		if settings.Project.Language.Name == "" {
			t.Error("Template should have language name")
		}
		if settings.Analysis.FilesIncludePath == "" {
			t.Error("Template should have files include path")
		}
	})
}

// Benchmarks

// BenchmarkLoadSettingsFromFile testa performance do carregamento
func BenchmarkLoadSettingsFromFile(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "benchmark-load-*")
	defer os.RemoveAll(tempDir)

	settingsPath := filepath.Join(tempDir, "settings.yml")
	WriteSettingsTemplate(settingsPath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadSettingsFromFile(settingsPath)
	}
}

// BenchmarkSaveSettingsToFile testa performance do salvamento
func BenchmarkSaveSettingsToFile(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "benchmark-save-*")
	defer os.RemoveAll(tempDir)

	settings := GetDefaultSettings()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		settingsPath := filepath.Join(tempDir, "bench-settings.yml")
		SaveSettingsToFile(settings, settingsPath)
	}
}

// BenchmarkGetSettingsTemplate testa performance da geração do template
func BenchmarkGetSettingsTemplate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSettingsTemplate()
	}
}
