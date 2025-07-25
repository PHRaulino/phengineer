package auth

import (
	"errors"
	"fmt"
)

// TokenGenerator função para gerar tokens
type TokenGenerator func(clientID, clientSecret string, scope TokenScope) (TokenResponse, error)

type AuthService struct {
	clientService  *ClientService
	tokenService   *TokenService
	tokenGenerator TokenGenerator
}

func NewAuthService(storage StorageAdapter, tokenGenerator TokenGenerator) *AuthService {
	return &AuthService{
		clientService:  NewClientService(storage),
		tokenService:   NewTokenService(storage),
		tokenGenerator: tokenGenerator,
	}
}

// SetupCredentials configura as credenciais do cliente
func (a *AuthService) SetupCredentials(clientID, clientSecret string) error {
	return a.clientService.SetCredentials(clientID, clientSecret)
}

// IsSetup verifica se as credenciais estão configuradas
func (a *AuthService) IsSetup() bool {
	return a.clientService.HasCredentials()
}

// GetValidToken retorna um token válido para o escopo especificado
func (a *AuthService) GetValidToken(scope TokenScope) (string, error) {
	if !a.IsSetup() {
		return "", errors.New("credentials not configured")
	}

	generateFunc := func(s TokenScope) (TokenResponse, error) {
		clientID, clientSecret, err := a.clientService.GetCredentials()
		if err != nil {
			return TokenResponse{}, fmt.Errorf("failed to get credentials: %w", err)
		}

		return a.tokenGenerator(clientID, clientSecret, s)
	}

	return a.tokenService.GetValidToken(scope, generateFunc)
}

// InvalidateToken remove um token específico (força nova geração)
func (a *AuthService) InvalidateToken(scope TokenScope) error {
	return a.tokenService.DeleteToken(scope)
}

// InvalidateAllTokens remove todos os tokens
func (a *AuthService) InvalidateAllTokens() error {
	scopes := []TokenScope{ScopeExecution, ScopeCreation, ScopeRead, ScopeWrite}

	for _, scope := range scopes {
		if err := a.tokenService.DeleteToken(scope); err != nil {
			return fmt.Errorf("failed to delete token for scope %s: %w", scope, err)
		}
	}
	return nil
}
