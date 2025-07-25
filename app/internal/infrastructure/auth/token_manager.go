package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// TokenScope define os escopos disponíveis
type TokenScope string

const (
	ScopeExecution TokenScope = "execution"
	ScopeCreation  TokenScope = "creation"
	ScopeRead      TokenScope = "read"
	ScopeWrite     TokenScope = "write"
)

// TokenData representa os dados do token
type TokenData struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Scope     TokenScope `json:"scope"`
}

// TokenResponse representa a resposta da API de token
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // segundos
}

type TokenService struct {
	storage StorageAdapter
}

func NewTokenService(storage StorageAdapter) *TokenService {
	return &TokenService{
		storage: storage,
	}
}

// getTokenKey gera a chave para o token baseado no escopo
func (t *TokenService) getTokenKey(scope TokenScope) string {
	return fmt.Sprintf("token_%s", string(scope))
}

// SaveToken salva o token com escopo e expiração
func (t *TokenService) SaveToken(scope TokenScope, tokenResponse TokenResponse) error {
	if tokenResponse.AccessToken == "" {
		return errors.New("token cannot be empty")
	}

	tokenData := TokenData{
		Token:     tokenResponse.AccessToken,
		ExpiresAt: time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
		Scope:     scope,
	}

	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	key := t.getTokenKey(scope)
	return t.storage.Set(ServiceName, key, string(data))
}

// GetToken recupera o token para um escopo específico
func (t *TokenService) GetToken(scope TokenScope) (*TokenData, error) {
	key := t.getTokenKey(scope)
	data, err := t.storage.Get(ServiceName, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get token for scope %s: %w", scope, err)
	}

	var tokenData TokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &tokenData, nil
}

// HasToken verifica se o token existe para um escopo
func (t *TokenService) HasToken(scope TokenScope) bool {
	key := t.getTokenKey(scope)
	return t.storage.Exists(ServiceName, key)
}

// DeleteToken remove o token para um escopo
func (t *TokenService) DeleteToken(scope TokenScope) error {
	key := t.getTokenKey(scope)
	return t.storage.Delete(ServiceName, key)
}

// IsTokenValid verifica se o token é válido (não expirado)
func (t *TokenService) IsTokenValid(tokenData *TokenData) bool {
	if tokenData == nil || tokenData.Token == "" {
		return false
	}
	return time.Now().Before(tokenData.ExpiresAt)
}

// GetValidToken retorna um token válido para o escopo (gera novo se necessário)
func (t *TokenService) GetValidToken(scope TokenScope, generateTokenFunc func(TokenScope) (TokenResponse, error)) (string, error) {
	// Verifica se existe token
	if t.HasToken(scope) {
		tokenData, err := t.GetToken(scope)
		if err == nil && t.IsTokenValid(tokenData) {
			return tokenData.Token, nil
		}
	}

	// Gera novo token
	tokenResponse, err := generateTokenFunc(scope)
	if err != nil {
		return "", fmt.Errorf("failed to generate token for scope %s: %w", scope, err)
	}

	// Salva o novo token
	if err := t.SaveToken(scope, tokenResponse); err != nil {
		return "", fmt.Errorf("failed to save token for scope %s: %w", scope, err)
	}

	return tokenResponse.AccessToken, nil
}
