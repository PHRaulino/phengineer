package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewRequirementsValidator testa a criação de uma nova instância
func TestNewRequirementsValidator(t *testing.T) {
	configFolder := "test-config"
	validator := NewRequirementsValidator(configFolder)

	if validator == nil {
		t.Fatal("NewRequirementsValidator returned nil")
	}

	if validator.ConfigFolderName != configFolder {
		t.Errorf("Expected ConfigFolderName to be '%s', got '%s'", configFolder, validator.ConfigFolderName)
	}
}

// TestValidationError_Error testa a implementação do erro customizado
func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Requirement: "Test Requirement",
		Message:     "Test error message",
	}

	expected := "Requirement 'Test Requirement' failed: Test error message"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestValidateGitInstalled testa a validação de instalação do Git
func TestValidateGitInstalled(t *testing.T) {
	validator := NewRequirementsValidator("config")

	err := validator.validateGitInstalled()

	// Verifica se o git está realmente instalado
	if _, lookupErr := exec.LookPath("git"); lookupErr != nil {
		// Git não está instalado, então esperamos um erro
		if err == nil {
			t.Error("Expected validation to fail when Git is not installed")
		}

		if validationErr, ok := err.(ValidationError); ok {
			if validationErr.Requirement != "Git Installation" {
				t.Errorf("Expected requirement 'Git Installation', got '%s'", validationErr.Requirement)
			}
		} else {
			t.Error("Expected ValidationError type")
		}
	} else {
		// Git está instalado, então esperamos sucesso
		if err != nil {
			t.Errorf("Expected validation to pass when Git is installed, got error: %v", err)
		}
	}
}

// TestValidateInGitRepository testa a validação de repositório Git
func TestValidateInGitRepository(t *testing.T) {
	validator := NewRequirementsValidator("config")

	// Verifica se estamos em um repositório Git
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, cmdErr := cmd.CombinedOutput()

	err := validator.validateInGitRepository()

	if cmdErr != nil || strings.TrimSpace(string(output)) != "true" {
		// Não estamos em um repositório Git, então esperamos um erro
		if err == nil {
			t.Error("Expected validation to fail when not in Git repository")
		}

		if validationErr, ok := err.(ValidationError); ok {
			if validationErr.Requirement != "Git Repository" {
				t.Errorf("Expected requirement 'Git Repository', got '%s'", validationErr.Requirement)
			}
		} else {
			t.Error("Expected ValidationError type")
		}
	} else {
		// Estamos em um repositório Git, então esperamos sucesso
		if err != nil {
			t.Errorf("Expected validation to pass when in Git repository, got error: %v", err)
		}
	}
}

// TestValidateGitRemoteExists testa a validação de remotes Git
func TestValidateGitRemoteExists(t *testing.T) {
	validator := NewRequirementsValidator("config")

	// Verifica se há remotes configurados
	cmd := exec.Command("git", "remote")
	output, cmdErr := cmd.CombinedOutput()

	err := validator.validateGitRemoteExists()

	if cmdErr != nil || strings.TrimSpace(string(output)) == "" {
		// Não há remotes configurados, então esperamos um erro
		if err == nil {
			t.Error("Expected validation to fail when no Git remotes exist")
		}

		if validationErr, ok := err.(ValidationError); ok {
			if validationErr.Requirement != "Git Remote" {
				t.Errorf("Expected requirement 'Git Remote', got '%s'", validationErr.Requirement)
			}
		} else {
			t.Error("Expected ValidationError type")
		}
	} else {
		// Há remotes configurados, então esperamos sucesso
		if err != nil {
			t.Errorf("Expected validation to pass when Git remotes exist, got error: %v", err)
		}
	}
}

// TestValidateConfigFolder testa a validação da pasta de configuração
func TestValidateConfigFolder(t *testing.T) {
	validator := NewRequirementsValidator("config")

	// Tenta obter a raiz do Git
	gitRoot, gitErr := validator.getGitRoot()

	err := validator.validateConfigFolder()

	if gitErr != nil {
		// Não conseguimos obter a raiz do Git, então esperamos um erro
		if err == nil {
			t.Error("Expected validation to fail when Git root cannot be determined")
		}
		return
	}

	// Verifica se a pasta de config existe
	configPath := filepath.Join(gitRoot, validator.ConfigFolderName)
	_, statErr := os.Stat(configPath)

	if os.IsNotExist(statErr) {
		// Pasta não existe, então esperamos um erro
		if err == nil {
			t.Error("Expected validation to fail when config folder does not exist")
		}

		if validationErr, ok := err.(ValidationError); ok {
			if validationErr.Requirement != "Config Folder" {
				t.Errorf("Expected requirement 'Config Folder', got '%s'", validationErr.Requirement)
			}
		} else {
			t.Error("Expected ValidationError type")
		}
	} else {
		// Pasta existe, então esperamos sucesso
		if err != nil {
			t.Errorf("Expected validation to pass when config folder exists, got error: %v", err)
		}
	}
}

// TestGetGitRoot testa a obtenção da raiz do repositório Git
func TestGetGitRoot(t *testing.T) {
	validator := NewRequirementsValidator("config")

	gitRoot, err := validator.GetGitRoot()

	// Verifica se estamos em um repositório Git
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, cmdErr := cmd.CombinedOutput()

	if cmdErr != nil {
		// Não estamos em um repositório Git, então esperamos um erro
		if err == nil {
			t.Error("Expected GetGitRoot to fail when not in Git repository")
		}
	} else {
		// Estamos em um repositório Git, então esperamos sucesso
		if err != nil {
			t.Errorf("Expected GetGitRoot to succeed when in Git repository, got error: %v", err)
		}

		expected := strings.TrimSpace(string(output))
		if gitRoot != expected {
			t.Errorf("Expected Git root '%s', got '%s'", expected, gitRoot)
		}
	}
}

// TestGetValidationSummary testa o resumo de validação
func TestGetValidationSummary(t *testing.T) {
	validator := NewRequirementsValidator("config")

	summary := validator.GetValidationSummary()

	// Verifica se todas as validações estão presentes
	expectedRequirements := []string{
		"Git Installation",
		"Git Repository",
		"Git Remote",
		"Config Folder",
	}

	for _, requirement := range expectedRequirements {
		if _, exists := summary[requirement]; !exists {
			t.Errorf("Expected requirement '%s' to be in summary", requirement)
		}
	}

	// Verifica se não há requirements extras
	if len(summary) != len(expectedRequirements) {
		t.Errorf("Expected %d requirements in summary, got %d", len(expectedRequirements), len(summary))
	}
}

// TestValidateAll testa a validação completa
func TestValidateAll(t *testing.T) {
	validator := NewRequirementsValidator("config")

	err := validator.ValidateAll()
	// O resultado depende do ambiente de teste
	// Vamos apenas verificar se o método executa sem panic
	if err != nil {
		// Se há erro, deve ser do tipo ValidationError
		if _, ok := err.(ValidationError); !ok {
			t.Errorf("Expected ValidationError type, got %T", err)
		}
	}
}

// TestIsValid testa o método de conveniência IsValid
func TestIsValid(t *testing.T) {
	validator := NewRequirementsValidator("config")

	isValid := validator.IsValid()
	err := validator.ValidateAll()

	// IsValid deve retornar true se ValidateAll não retornar erro
	if (err == nil) != isValid {
		t.Errorf("IsValid() returned %v, but ValidateAll() error was %v", isValid, err)
	}
}

// TestValidateWithDetails testa a validação com detalhes
func TestValidateWithDetails(t *testing.T) {
	validator := NewRequirementsValidator("config")

	results, err := validator.ValidateWithDetails()

	// Verifica se temos todos os requirements
	expectedRequirements := []string{
		"Git Installation",
		"Git Repository",
		"Git Remote",
		"Config Folder",
	}

	if len(results) != len(expectedRequirements) {
		t.Errorf("Expected %d requirement results, got %d", len(expectedRequirements), len(results))
	}

	// Verifica se todos os requirements estão presentes
	for _, expected := range expectedRequirements {
		found := false
		for _, result := range results {
			if result.Name == expected {
				found = true

				// Verifica se todos os campos estão preenchidos
				if result.Description == "" {
					t.Errorf("Expected description for requirement '%s'", expected)
				}

				// Se passou, não deve ter erro
				if result.Passed && result.ErrorMsg != "" {
					t.Errorf("Requirement '%s' passed but has error message: %s", expected, result.ErrorMsg)
				}

				// Se não passou, deve ter erro
				if !result.Passed && result.ErrorMsg == "" {
					t.Errorf("Requirement '%s' failed but has no error message", expected)
				}

				break
			}
		}

		if !found {
			t.Errorf("Expected requirement '%s' not found in results", expected)
		}
	}

	// Verifica consistência com ValidateAll
	validateAllErr := validator.ValidateAll()
	if (validateAllErr == nil) != (err == nil) {
		t.Errorf("ValidateWithDetails error consistency: ValidateAll=%v, ValidateWithDetails=%v", validateAllErr, err)
	}
}

// TestValidateAll_WithMockConfigFolder testa a validação com uma pasta de config personalizada
func TestValidateAll_WithMockConfigFolder(t *testing.T) {
	validator := NewRequirementsValidator("config")
	gitRoot, err := validator.GetGitRoot()
	if err != nil {
		t.Skip("Skipping test: not in a Git repository")
	}

	// Cria uma pasta de config temporária
	configDir := filepath.Join(gitRoot, "temp-test-config")
	err = os.Mkdir(configDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	defer os.RemoveAll(configDir)

	// Testa com a pasta de config existente
	testValidator := NewRequirementsValidator("temp-test-config")
	err = testValidator.validateConfigFolder()
	if err != nil {
		t.Errorf("Expected validation to pass with existing config folder, got error: %v", err)
	}
}

// TestRequirementStatus testa a estrutura RequirementStatus
func TestRequirementStatus(t *testing.T) {
	status := RequirementStatus{
		Name:        "Test Requirement",
		Passed:      true,
		ErrorMsg:    "",
		Description: "Test description",
	}

	if status.Name != "Test Requirement" {
		t.Errorf("Expected Name 'Test Requirement', got '%s'", status.Name)
	}

	if !status.Passed {
		t.Error("Expected Passed to be true")
	}

	if status.ErrorMsg != "" {
		t.Errorf("Expected empty ErrorMsg, got '%s'", status.ErrorMsg)
	}

	if status.Description != "Test description" {
		t.Errorf("Expected Description 'Test description', got '%s'", status.Description)
	}
}

// BenchmarkValidateAll testa a performance da validação completa
func BenchmarkValidateAll(b *testing.B) {
	validator := NewRequirementsValidator("config")

	for i := 0; i < b.N; i++ {
		validator.ValidateAll()
	}
}

// BenchmarkGetValidationSummary testa a performance do resumo de validação
func BenchmarkGetValidationSummary(b *testing.B) {
	validator := NewRequirementsValidator("config")

	for i := 0; i < b.N; i++ {
		validator.GetValidationSummary()
	}
}

// BenchmarkValidateWithDetails testa a performance da validação com detalhes
func BenchmarkValidateWithDetails(b *testing.B) {
	validator := NewRequirementsValidator("config")

	for i := 0; i < b.N; i++ {
		validator.ValidateWithDetails()
	}
}
