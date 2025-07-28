package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/storage"
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/token"
)

const (
	StackSpotTokenURL = "https://idm.stackspot.com/realms/stackspot/protocol/openid-connect/token"
)

type StackSpotProvider struct {
	storage storage.StorageAdapter
}

func NewStackSpotProvider(storage storage.StorageAdapter) *StackSpotProvider {
	return &StackSpotProvider{
		storage: storage,
	}
}

type StackSpotTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope,omitempty"`
}

type StackSpotTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func (p *StackSpotProvider) GetToken(scope token.TokenScope) (token.TokenResponse, error) {
	// Buscar credenciais armazenadas
	clientID, err := p.storage.Get("stackspot_client_id")
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("client_id não encontrado: %w", err)
	}

	clientSecret, err := p.storage.Get("stackspot_client_secret")
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("client_secret não encontrado: %w", err)
	}

	// Preparar dados para requisição
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	
	// Mapear escopo interno para escopo StackSpot
	if stackSpotScope := p.mapScope(scope); stackSpotScope != "" {
		data.Set("scope", stackSpotScope)
	}

	// Fazer requisição HTTP
	req, err := http.NewRequest("POST", StackSpotTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("erro na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return token.TokenResponse{}, fmt.Errorf("erro na autenticação StackSpot: status %d", resp.StatusCode)
	}

	// Decodificar resposta
	var tokenResp StackSpotTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return token.TokenResponse{}, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return token.TokenResponse{
		AccessToken: tokenResp.AccessToken,
		ExpiresIn:   tokenResp.ExpiresIn,
	}, nil
}

func (p *StackSpotProvider) SaveCredentials(clientID, clientSecret string) error {
	if err := p.storage.Set("stackspot_client_id", clientID); err != nil {
		return fmt.Errorf("erro ao salvar client_id: %w", err)
	}

	if err := p.storage.Set("stackspot_client_secret", clientSecret); err != nil {
		return fmt.Errorf("erro ao salvar client_secret: %w", err)
	}

	return nil
}

func (p *StackSpotProvider) mapScope(scope token.TokenScope) string {
	switch scope {
	case token.ScopeExecution:
		return "execution"
	case token.ScopeCreation:
		return "creation"
	case token.ScopeRead:
		return "read"
	case token.ScopeWrite:
		return "write"
	default:
		return ""
	}
}