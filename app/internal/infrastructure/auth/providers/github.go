package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/storage"
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/token"
)

type GitHubProvider struct {
	storage storage.StorageAdapter
}

func NewGitHubProvider(storage storage.StorageAdapter) *GitHubProvider {
	return &GitHubProvider{
		storage: storage,
	}
}

type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *GitHubProvider) GetToken(scope token.TokenScope) (token.TokenResponse, error) {
	// Buscar token GitHub armazenado
	githubToken, err := p.storage.Get("github_token")
	if err != nil {
		// Tentar variável de ambiente como fallback
		githubToken = os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			return token.TokenResponse{}, fmt.Errorf("github_token não encontrado: %w", err)
		}
	}

	// Validar token fazendo uma requisição de teste
	if err := p.validateToken(githubToken); err != nil {
		return token.TokenResponse{}, fmt.Errorf("token GitHub inválido: %w", err)
	}

	// GitHub tokens geralmente não expiram (Personal Access Tokens)
	// Retornar com expiração longa
	return token.TokenResponse{
		AccessToken: githubToken,
		ExpiresIn:   86400 * 365, // 1 ano
	}, nil
}

func (p *GitHubProvider) validateToken(token string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição de validação: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro na requisição de validação: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token inválido: status %d", resp.StatusCode)
	}

	return nil
}

func (p *GitHubProvider) GetUser() (*GitHubUser, error) {
	githubToken, err := p.storage.Get("github_token")
	if err != nil {
		githubToken = os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			return nil, fmt.Errorf("github_token não encontrado: %w", err)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao buscar usuário: status %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("erro ao decodificar usuário: %w", err)
	}

	return &user, nil
}

func (p *GitHubProvider) SaveToken(githubToken string) error {
	// Validar token antes de salvar
	if err := p.validateToken(githubToken); err != nil {
		return fmt.Errorf("token inválido: %w", err)
	}

	if err := p.storage.Set("github_token", githubToken); err != nil {
		return fmt.Errorf("erro ao salvar github_token: %w", err)
	}

	return nil
}

// ListRepositories lista repositórios do usuário
func (p *GitHubProvider) ListRepositories() ([]map[string]interface{}, error) {
	githubToken, err := p.storage.Get("github_token")
	if err != nil {
		githubToken = os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			return nil, fmt.Errorf("github_token não encontrado: %w", err)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos?sort=updated&per_page=100", nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao listar repositórios: status %d", resp.StatusCode)
	}

	var repos []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("erro ao decodificar repositórios: %w", err)
	}

	return repos, nil
}