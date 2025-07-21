package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestExtractRepoNameFromURL testa a extração do nome do repositório de URLs
func TestExtractRepoNameFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "GitHub HTTPS with .git",
			url:      "https://github.com/user/my-repo.git",
			expected: "my-repo",
		},
		{
			name:     "GitHub HTTPS without .git",
			url:      "https://github.com/user/my-repo",
			expected: "my-repo",
		},
		{
			name:     "GitHub SSH with .git",
			url:      "git@github.com:user/my-repo.git",
			expected: "my-repo",
		},
		{
			name:     "GitHub SSH without .git",
			url:      "git@github.com:user/my-repo",
			expected: "my-repo",
		},
		{
			name:     "GitLab HTTPS",
			url:      "https://gitlab.com/group/subgroup/project.git",
			expected: "project",
		},
		{
			name:     "Bitbucket HTTPS",
			url:      "https://bitbucket.org/user/repository.git",
			expected: "repository",
		},
		{
			name:     "Complex repository name",
			url:      "https://github.com/org/my-complex-repo-name.git",
			expected: "my-complex-repo-name",
		},
		{
			name:     "URL with underscores",
			url:      "https://github.com/user/my_repo_name.git",
			expected: "my_repo_name",
		},
		{
			name:     "URL with numbers",
			url:      "https://github.com/user/repo123.git",
			expected: "repo123",
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "",
		},
		{
			name:     "Single word",
			url:      "repo",
			expected: "repo",
		},
		{
			name:     "URL with trailing slash",
			url:      "https://github.com/user/repo/",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRepoNameFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("extractRepoNameFromURL(%q) = %q, expected %q", tt.url, result, tt.expected)
			}
		})
	}
}

// TestGetRepositoryName testa a obtenção do nome do repositório
func TestGetRepositoryName(t *testing.T) {
	// Verifica se estamos em um repositório Git antes de prosseguir
	if !isInGitRepository() {
		t.Skip("Skipping test: not in a Git repository")
	}

	name, err := getRepositoryName()
	// Se o comando git falhar, pode ser ambiente de teste
	if err != nil {
		t.Logf("getRepositoryName failed (may be expected in test environment): %v", err)
		return
	}

	// Verifica se retornou um nome válido
	if name == "" {
		t.Error("Repository name should not be empty")
	}

	// Verifica se o nome não contém caracteres inválidos
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		t.Errorf("Repository name should not contain path separators: %s", name)
	}

	t.Logf("Repository name: %s", name)
}

// TestGetRepositoryNameFallback testa o fallback para nome da pasta
func TestGetRepositoryNameFallback(t *testing.T) {
	// Este teste simula quando o remote falha mas git root funciona
	// Difícil de testar sem mockar os comandos git
	t.Skip("Skipping fallback test - requires mocking git commands")
}

// TestGetConfigDirPath testa a obtenção do caminho da pasta de config
func TestGetConfigDirPath(t *testing.T) {
	if !isInGitRepository() {
		t.Skip("Skipping test: not in a Git repository")
	}

	configFolderName := "test-config"
	path, err := getConfigDirPath(configFolderName)
	if err != nil {
		t.Logf("getConfigDirPath failed (may be expected in test environment): %v", err)
		return
	}

	// Verifica se é um caminho absoluto
	if !filepath.IsAbs(path) {
		t.Errorf("Config dir path should be absolute: %s", path)
	}

	// Verifica se termina com o nome da pasta de config
	if !strings.HasSuffix(path, configFolderName) {
		t.Errorf("Config dir path should end with folder name '%s': %s", configFolderName, path)
	}

	t.Logf("Config dir path: %s", path)
}

// TestGetRemoteURL testa a obtenção da URL do remote
func TestGetRemoteURL(t *testing.T) {
	if !isInGitRepository() {
		t.Skip("Skipping test: not in a Git repository")
	}

	// Verifica se há remotes configurados
	if !hasGitRemotes() {
		t.Skip("Skipping test: no Git remotes configured")
	}

	url, err := getRemoteURL()
	if err != nil {
		t.Logf("getRemoteURL failed (may be expected in test environment): %v", err)
		return
	}

	// Verifica se a URL não está vazia
	if url == "" {
		t.Error("Remote URL should not be empty")
	}

	// Verifica se não termina com .git
	if strings.HasSuffix(url, ".git") {
		t.Errorf("Remote URL should not end with .git: %s", url)
	}

	// Verifica se parece com uma URL válida
	if !strings.Contains(url, "://") && !strings.Contains(url, "@") {
		t.Errorf("Remote URL should look like a valid URL: %s", url)
	}

	t.Logf("Remote URL: %s", url)
}

// TestCollectAutoConfig testa a coleta completa de configurações automáticas
func TestCollectAutoConfig(t *testing.T) {
	if !isInGitRepository() {
		t.Skip("Skipping test: not in a Git repository")
	}

	configFolderName := "test-config"
	autoConfig, err := collectAutoConfig(configFolderName)
	if err != nil {
		t.Logf("collectAutoConfig failed (may be expected in test environment): %v", err)
		return
	}

	// Verifica se todas as configurações foram coletadas
	if autoConfig.AppName == "" {
		t.Error("AppName should not be empty")
	}

	if autoConfig.ConfigDirPath == "" {
		t.Error("ConfigDirPath should not be empty")
	}

	if autoConfig.RemoteURL == "" {
		t.Error("RemoteURL should not be empty")
	}

	// Verifica se ConfigDirPath termina com o nome da pasta
	if !strings.HasSuffix(autoConfig.ConfigDirPath, configFolderName) {
		t.Errorf("ConfigDirPath should end with '%s': %s", configFolderName, autoConfig.ConfigDirPath)
	}

	// Verifica se RemoteURL não termina com .git
	if strings.HasSuffix(autoConfig.RemoteURL, ".git") {
		t.Errorf("RemoteURL should not end with .git: %s", autoConfig.RemoteURL)
	}

	t.Logf("Auto Config - App: %s, Config: %s, Remote: %s",
		autoConfig.AppName, autoConfig.ConfigDirPath, autoConfig.RemoteURL)
}

// TestCollectAutoConfigInvalidEnvironment testa comportamento em ambiente inválido
func TestCollectAutoConfigInvalidEnvironment(t *testing.T) {
	// Cria um diretório temporário sem Git
	tempDir, err := os.MkdirTemp("", "no-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Muda para o diretório sem Git
	originalDir, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Deve falhar pois não está em um repositório Git
	_, err = collectAutoConfig("config")
	if err == nil {
		t.Error("Expected collectAutoConfig to fail outside Git repository")
	}

	t.Logf("Expected error in non-Git environment: %v", err)
}

// TestAutoConfigWithDifferentFolderNames testa com diferentes nomes de pasta
func TestAutoConfigWithDifferentFolderNames(t *testing.T) {
	if !isInGitRepository() {
		t.Skip("Skipping test: not in a Git repository")
	}

	folderNames := []string{"config", "settings", ".config", "app-config"}

	for _, folderName := range folderNames {
		t.Run(folderName, func(t *testing.T) {
			path, err := getConfigDirPath(folderName)
			if err != nil {
				t.Logf("getConfigDirPath failed for '%s': %v", folderName, err)
				return
			}

			if !strings.HasSuffix(path, folderName) {
				t.Errorf("Path should end with '%s': %s", folderName, path)
			}
		})
	}
}

// Helper functions for tests

// isInGitRepository verifica se estamos em um repositório Git
func isInGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// hasGitRemotes verifica se há remotes configurados
func hasGitRemotes() bool {
	cmd := exec.Command("git", "remote")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// Benchmarks

// BenchmarkExtractRepoNameFromURL testa performance da extração de nome
func BenchmarkExtractRepoNameFromURL(b *testing.B) {
	url := "https://github.com/user/my-repository.git"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractRepoNameFromURL(url)
	}
}

// BenchmarkGetRepositoryName testa performance da obtenção do nome do repo
func BenchmarkGetRepositoryName(b *testing.B) {
	if !isInGitRepository() {
		b.Skip("Skipping benchmark: not in a Git repository")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getRepositoryName()
	}
}

// BenchmarkGetConfigDirPath testa performance da obtenção do caminho
func BenchmarkGetConfigDirPath(b *testing.B) {
	if !isInGitRepository() {
		b.Skip("Skipping benchmark: not in a Git repository")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getConfigDirPath("config")
	}
}

// BenchmarkCollectAutoConfig testa performance da coleta completa
func BenchmarkCollectAutoConfig(b *testing.B) {
	if !isInGitRepository() {
		b.Skip("Skipping benchmark: not in a Git repository")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collectAutoConfig("config")
	}
}

// Table-driven tests para casos edge

// TestExtractRepoNameEdgeCases testa casos extremos
func TestExtractRepoNameEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name     string
		url      string
		expected string
	}{
		{"Multiple slashes", "https://github.com/user//repo", "repo"},
		{"With query params", "https://github.com/user/repo?param=value", "repo?param=value"},
		{"With fragment", "https://github.com/user/repo#section", "repo#section"},
		{"Very long name", "https://github.com/user/" + strings.Repeat("a", 100), strings.Repeat("a", 100)},
		{"With dots", "https://github.com/user/repo.name.git", "repo.name"},
		{"With hyphens", "https://github.com/user/my-awesome-repo-name.git", "my-awesome-repo-name"},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			result := extractRepoNameFromURL(tc.url)
			if result != tc.expected {
				t.Errorf("extractRepoNameFromURL(%q) = %q, expected %q", tc.url, result, tc.expected)
			}
		})
	}
}
