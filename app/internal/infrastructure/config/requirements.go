package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RequirementsValidator é responsável por validar os pré-requisitos da CLI
type RequirementsValidator struct {
	ConfigFolderName string
}

// NewRequirementsValidator cria uma nova instância do validador
func NewRequirementsValidator(configFolderName string) *RequirementsValidator {
	return &RequirementsValidator{
		ConfigFolderName: configFolderName,
	}
}

// ValidationError representa um erro de validação de requirement
type ValidationError struct {
	Requirement string
	Message     string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Requirement '%s' failed: %s", e.Requirement, e.Message)
}

// RequirementStatus representa o status de um requirement
type RequirementStatus struct {
	Name        string
	Passed      bool
	ErrorMsg    string
	Description string
}

// ValidateAll executa todas as validações necessárias
func (rv *RequirementsValidator) ValidateAll() error {
	validators := []func() error{
		rv.validateGitInstalled,
		rv.validateInGitRepository,
		rv.validateGitRemoteExists,
		rv.validateConfigFolder,
	}

	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateWithDetails executa todas as validações e retorna detalhes completos
func (rv *RequirementsValidator) ValidateWithDetails() ([]RequirementStatus, error) {
	requirements := []struct {
		name        string
		description string
		validator   func() error
	}{
		{
			name:        "Git Installation",
			description: "Git must be installed and available in PATH",
			validator:   rv.validateGitInstalled,
		},
		{
			name:        "Git Repository",
			description: "Must be executed within a Git repository",
			validator:   rv.validateInGitRepository,
		},
		{
			name:        "Git Remote",
			description: "Repository must have at least one remote configured",
			validator:   rv.validateGitRemoteExists,
		},
		{
			name:        "Config Folder",
			description: fmt.Sprintf("Config folder '%s' must exist in repository root", rv.ConfigFolderName),
			validator:   rv.validateConfigFolder,
		},
	}

	var results []RequirementStatus
	var firstError error

	for _, req := range requirements {
		status := RequirementStatus{
			Name:        req.name,
			Description: req.description,
			Passed:      true,
		}

		if err := req.validator(); err != nil {
			status.Passed = false
			status.ErrorMsg = err.Error()

			if firstError == nil {
				firstError = err
			}
		}

		results = append(results, status)
	}

	return results, firstError
}

// validateGitInstalled verifica se o Git está instalado no sistema
func (rv *RequirementsValidator) validateGitInstalled() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return ValidationError{
			Requirement: "Git Installation",
			Message:     "Git is not installed or not available in PATH",
		}
	}
	return nil
}

// validateInGitRepository verifica se estamos dentro de um repositório Git
func (rv *RequirementsValidator) validateInGitRepository() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ValidationError{
			Requirement: "Git Repository",
			Message:     "Not inside a Git repository",
		}
	}

	if strings.TrimSpace(string(output)) != "true" {
		return ValidationError{
			Requirement: "Git Repository",
			Message:     "Not inside a Git working tree",
		}
	}

	return nil
}

// validateGitRemoteExists verifica se há pelo menos um remote configurado
func (rv *RequirementsValidator) validateGitRemoteExists() error {
	cmd := exec.Command("git", "remote")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ValidationError{
			Requirement: "Git Remote",
			Message:     "Failed to check Git remotes",
		}
	}

	remotes := strings.TrimSpace(string(output))
	if remotes == "" {
		return ValidationError{
			Requirement: "Git Remote",
			Message:     "No Git remotes configured",
		}
	}

	return nil
}

// validateConfigFolder verifica se a pasta de config existe na raiz do repositório
func (rv *RequirementsValidator) validateConfigFolder() error {
	gitRoot, err := rv.getGitRoot()
	if err != nil {
		return ValidationError{
			Requirement: "Config Folder",
			Message:     fmt.Sprintf("Failed to find Git root: %v", err),
		}
	}

	configPath := filepath.Join(gitRoot, rv.ConfigFolderName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return ValidationError{
			Requirement: "Config Folder",
			Message:     fmt.Sprintf("Config folder '%s' not found in Git root (%s)", rv.ConfigFolderName, gitRoot),
		}
	}

	return nil
}

// getGitRoot retorna o caminho para a raiz do repositório Git
func (rv *RequirementsValidator) getGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get Git root: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetGitRoot é a versão pública do método para obter a raiz do Git
func (rv *RequirementsValidator) GetGitRoot() (string, error) {
	return rv.getGitRoot()
}

// GetValidationSummary retorna um resumo do status de todos os requirements
func (rv *RequirementsValidator) GetValidationSummary() map[string]bool {
	summary := make(map[string]bool)

	validators := map[string]func() error{
		"Git Installation": rv.validateGitInstalled,
		"Git Repository":   rv.validateInGitRepository,
		"Git Remote":       rv.validateGitRemoteExists,
		"Config Folder":    rv.validateConfigFolder,
	}

	for name, validator := range validators {
		summary[name] = validator() == nil
	}

	return summary
}

// IsValid verifica se todos os requirements são válidos (método de conveniência)
func (rv *RequirementsValidator) IsValid() bool {
	return rv.ValidateAll() == nil
}
