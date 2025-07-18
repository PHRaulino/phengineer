package config

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// collectAutoConfig coleta configurações automáticas do ambiente
func collectAutoConfig(configFolderName string) (*AutoConfig, error) {
	auto := &AutoConfig{}

	// Coleta nome do app (nome do repositório)
	appName, err := getRepositoryName()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository name: %w", err)
	}
	auto.AppName = appName

	// Coleta caminho da pasta de config
	configDirPath, err := getConfigDirPath(configFolderName)
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir path: %w", err)
	}
	auto.ConfigDirPath = configDirPath

	// Coleta URL do remote
	remoteURL, err := getRemoteURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote URL: %w", err)
	}
	auto.RemoteURL = remoteURL

	return auto, nil
}

// getRepositoryName obtém o nome do repositório atual
func getRepositoryName() (string, error) {
	// Primeiro tenta pegar do remote
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.CombinedOutput()
	if err == nil {
		remoteURL := strings.TrimSpace(string(output))
		if name := extractRepoNameFromURL(remoteURL); name != "" {
			return name, nil
		}
	}

	// Se não conseguir do remote, pega o nome da pasta
	cmd = exec.Command("git", "rev-parse", "--show-toplevel")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get git root: %v", err)
	}

	gitRoot := strings.TrimSpace(string(output))
	return filepath.Base(gitRoot), nil
}

// extractRepoNameFromURL extrai o nome do repositório de uma URL
func extractRepoNameFromURL(url string) string {
	// Remove .git no final se existir
	url = strings.TrimSuffix(url, ".git")

	// Divide por / e pega o último elemento
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

// getConfigDirPath obtém o caminho completo da pasta de configuração
func getConfigDirPath(configFolderName string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get git root: %v", err)
	}

	gitRoot := strings.TrimSpace(string(output))
	return filepath.Join(gitRoot, configFolderName), nil
}

// getRemoteURL obtém a URL do remote origin sem .git
func getRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %v", err)
	}

	remoteURL := strings.TrimSpace(string(output))

	// Remove .git no final se existir
	remoteURL = strings.TrimSuffix(remoteURL, ".git")

	return remoteURL, nil
}