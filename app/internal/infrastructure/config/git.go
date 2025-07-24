package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func detectGitInfo() (repoName, remoteName, rootPath string) {
	// Root path = diretório atual
	rootPath, _ = os.Getwd()

	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return filepath.Base(rootPath), "", rootPath
	}

	// Remote origin
	remotes, _ := repo.Remotes()
	if len(remotes) > 0 {
		url := remotes[0].Config().URLs[0]
		remoteName = extractRepoNameFromURL(url)
	}

	return remoteName, remoteName, rootPath
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
