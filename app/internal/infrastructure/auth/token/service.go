package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/storage"
)

type Service struct {
	storage    storage.StorageAdapter
	generators map[TokenGeneratorAlias]TokenGenerator // key: alias, value: generator function
}

func NewService(storage storage.StorageAdapter) *Service {
	return &Service{
		storage:    storage,
		generators: make(map[TokenGeneratorAlias]TokenGenerator),
	}
}

// RegisterGenerator registra uma função geradora de token para um alias
func (s *Service) RegisterGenerator(alias TokenGeneratorAlias, generator TokenGenerator) {
	s.generators[alias] = generator
}

// getTokenKey gera a chave para armazenamento baseado no escopo e alias
func (s *Service) getTokenKey(scope TokenScope, alias TokenGeneratorAlias) string {
	return fmt.Sprintf("token_%s_%s", string(scope), alias)
}

// Get busca um token válido para o escopo e alias
func (s *Service) Get(scope TokenScope, alias TokenGeneratorAlias) (string, error) {
	key := s.getTokenKey(scope, alias)

	// 1. Buscar token existente
	if s.storage.Exists(key) {
		tokenData, err := s.getStoredToken(key)
		if err == nil {
			// 2. Validar se ainda é válido
			if s.isValid(tokenData) {
				return tokenData.Token, nil
			}
		}
	}

	// 3. Se não existe ou inválido, criar novo
	return s.create(scope, alias)
}

// getStoredToken recupera e deserializa token do storage
func (s *Service) getStoredToken(key string) (*TokenData, error) {
	data, err := s.storage.Get(key)
	if err != nil {
		return nil, err
	}

	var tokenData TokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &tokenData, nil
}

// isValid verifica se o token não expirou
func (s *Service) isValid(tokenData *TokenData) bool {
	if tokenData == nil || tokenData.Token == "" {
		return false
	}
	return time.Now().Before(tokenData.ExpiresAt)
}

// create gera um novo token usando o generator registrado
func (s *Service) create(scope TokenScope, alias TokenGeneratorAlias) (string, error) {
	generator, exists := s.generators[alias]
	if !exists {
		return "", fmt.Errorf("no generator registered for alias: %s", alias)
	}

	// Gerar token via função registrada
	tokenResponse, err := generator(scope)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Salvar no storage
	tokenData := TokenData{
		Token:     tokenResponse.AccessToken,
		ExpiresAt: time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
		Scope:     scope,
		Alias:     alias,
	}

	if err := s.save(tokenData); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	return tokenData.Token, nil
}

// save serializa e armazena o token
func (s *Service) save(tokenData TokenData) error {
	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	key := s.getTokenKey(tokenData.Scope, tokenData.Alias)
	return s.storage.Set(key, string(data))
}

// Delete remove um token específico
func (s *Service) Delete(scope TokenScope, alias TokenGeneratorAlias) error {
	key := s.getTokenKey(scope, alias)
	return s.storage.Delete(key)
}

// Exists verifica se existe token para o escopo e alias
func (s *Service) Exists(scope TokenScope, alias TokenGeneratorAlias) bool {
	key := s.getTokenKey(scope, alias)
	return s.storage.Exists(key)
}
