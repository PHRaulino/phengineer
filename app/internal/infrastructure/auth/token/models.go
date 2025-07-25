package token

import "time"

type TokenGeneratorAlias string

const (
	TokenGenSTK TokenGeneratorAlias = "stackspot-api"
	TokenGenHC  TokenGeneratorAlias = "hashicorp-vault"
)

// TokenScope define os escopos disponíveis
type TokenScope string

const (
	ScopeExecution TokenScope = "execution"
	ScopeCreation  TokenScope = "creation"
	ScopeRead      TokenScope = "read"
	ScopeWrite     TokenScope = "write"
)

// TokenData representa os dados do token armazenados
type TokenData struct {
	Token     string              `json:"token"`
	ExpiresAt time.Time           `json:"expires_at"`
	Scope     TokenScope          `json:"scope"`
	Alias     TokenGeneratorAlias `json:"alias"`
}

// TokenResponse representa a resposta da API de token
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // segundos
}

// TokenGenerator função para gerar tokens
type TokenGenerator func(scope TokenScope) (TokenResponse, error)
