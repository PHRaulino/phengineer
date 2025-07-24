// auth/storage.go
package auth

import (
	"errors"
	"fmt"
)

// StorageAdapter interface para diferentes tipos de armazenamento
type StorageAdapter interface {
	Set(service, key, value string) error
	Get(service, key string) (string, error)
	Delete(service, key string) error
	Exists(service, key string) bool
}

// auth/keyring_adapter.go
package auth

import (
	"github.com/zalando/go-keyring"
)

type KeyringAdapter struct{}

func NewKeyringAdapter() *KeyringAdapter {
	return &KeyringAdapter{}
}

func (k *KeyringAdapter) Set(service, key, value string) error {
	return keyring.Set(service, key, value)
}

func (k *KeyringAdapter) Get(service, key string) (string, error) {
	return keyring.Get(service, key)
}

func (k *KeyringAdapter) Delete(service, key string) error {
	return keyring.Delete(service, key)
}

func (k *KeyringAdapter) Exists(service, key string) bool {
	_, err := keyring.Get(service, key)
	return err == nil
}

// auth/memory_adapter.go
package auth

import (
	"errors"
	"sync"
)

type MemoryAdapter struct {
	data map[string]map[string]string
	mu   sync.RWMutex
}

func NewMemoryAdapter() *MemoryAdapter {
	return &MemoryAdapter{
		data: make(map[string]map[string]string),
	}
}

func (m *MemoryAdapter) Set(service, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.data[service] == nil {
		m.data[service] = make(map[string]string)
	}
	m.data[service][key] = value
	return nil
}

func (m *MemoryAdapter) Get(service, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if serviceData, exists := m.data[service]; exists {
		if value, exists := serviceData[key]; exists {
			return value, nil
		}
	}
	return "", errors.New("key not found")
}

func (m *MemoryAdapter) Delete(service, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if serviceData, exists := m.data[service]; exists {
		delete(serviceData, key)
	}
	return nil
}

func (m *MemoryAdapter) Exists(service, key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if serviceData, exists := m.data[service]; exists {
		_, exists := serviceData[key]
		return exists
	}
	return false
}

// auth/client.go
package auth

import (
	"errors"
	"fmt"
)

const (
	ServiceName     = "stackspot-cli"
	ClientIDKey     = "client_id"
	ClientSecretKey = "client_secret"
)

type ClientService struct {
	storage StorageAdapter
}

func NewClientService(storage StorageAdapter) *ClientService {
	return &ClientService{
		storage: storage,
	}
}

// SetCredentials salva client_id e client_secret
func (c *ClientService) SetCredentials(clientID, clientSecret string) error {
	if clientID == "" || clientSecret == "" {
		return errors.New("client_id and client_secret cannot be empty")
	}

	if err := c.storage.Set(ServiceName, ClientIDKey, clientID); err != nil {
		return fmt.Errorf("failed to save client_id: %w", err)
	}

	if err := c.storage.Set(ServiceName, ClientSecretKey, clientSecret); err != nil {
		return fmt.Errorf("failed to save client_secret: %w", err)
	}

	return nil
}

// GetCredentials recupera client_id e client_secret
func (c *ClientService) GetCredentials() (clientID, clientSecret string, err error) {
	clientID, err = c.storage.Get(ServiceName, ClientIDKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to get client_id: %w", err)
	}

	clientSecret, err = c.storage.Get(ServiceName, ClientSecretKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to get client_secret: %w", err)
	}

	return clientID, clientSecret, nil
}

// HasCredentials verifica se as credenciais existem
func (c *ClientService) HasCredentials() bool {
	return c.storage.Exists(ServiceName, ClientIDKey) && 
		   c.storage.Exists(ServiceName, ClientSecretKey)
}

// DeleteCredentials remove as credenciais
func (c *ClientService) DeleteCredentials() error {
	if err := c.storage.Delete(ServiceName, ClientIDKey); err != nil {
		return fmt.Errorf("failed to delete client_id: %w", err)
	}

	if err := c.storage.Delete(ServiceName, ClientSecretKey); err != nil {
		return fmt.Errorf("failed to delete client_secret: %w", err)
	}

	return nil
}

// auth/token.go
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

// auth/service.go
package auth

import (
	"fmt"
)

// TokenGenerator função para gerar tokens
type TokenGenerator func(clientID, clientSecret string, scope TokenScope) (TokenResponse, error)

type AuthService struct {
	clientService    *ClientService
	tokenService     *TokenService
	tokenGenerator   TokenGenerator
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

// auth/factory.go
package auth

import (
	"github.com/spf13/viper"
)

// StorageType define o tipo de storage
type StorageType string

const (
	StorageKeyring StorageType = "keyring"
	StorageMemory  StorageType = "memory"
)

// NewStorageAdapter cria o adapter baseado na configuração
func NewStorageAdapter() StorageAdapter {
	storageType := viper.GetString("auth.storage_type")
	
	switch StorageType(storageType) {
	case StorageMemory:
		return NewMemoryAdapter()
	default:
		return NewKeyringAdapter()
	}
}

// cmd/auth.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"your-project/auth"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User authentication setup",
	RunE:  runUserAuth,
}

func init() {
	userCmd.Flags().String("id", "", "Client ID")
	userCmd.Flags().String("secret", "", "Client Secret")
	userCmd.MarkFlagRequired("id")
	userCmd.MarkFlagRequired("secret")

	authCmd.AddCommand(userCmd)
	rootCmd.AddCommand(authCmd)
}

func runUserAuth(cmd *cobra.Command, args []string) error {
	clientID, _ := cmd.Flags().GetString("id")
	clientSecret, _ := cmd.Flags().GetString("secret")

	// Inicializar serviço de auth
	storage := auth.NewStorageAdapter()
	
	// Token generator placeholder
	tokenGenerator := func(clientID, clientSecret string, scope auth.TokenScope) (auth.TokenResponse, error) {
		// TODO: Implementar chamada real para API
		return auth.TokenResponse{
			AccessToken: fmt.Sprintf("token_%s_%d", scope, time.Now().Unix()),
			ExpiresIn:   3600, // 1 hora
		}, nil
	}
	
	authService := auth.NewAuthService(storage, tokenGenerator)

	// Salvar credenciais
	if err := authService.SetupCredentials(clientID, clientSecret); err != nil {
		return fmt.Errorf("failed to setup credentials: %w", err)
	}

	fmt.Println("✅ Credentials saved successfully!")
	return nil
}

// main.go - Exemplo de uso
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/spf13/viper"
	"your-project/auth"
)

// tokenGenerator implementa a geração real de tokens
func tokenGenerator(clientID, clientSecret string, scope auth.TokenScope) (auth.TokenResponse, error) {
	// TODO: Implementar chamada real para sua API
	// Por enquanto simula uma resposta
	return auth.TokenResponse{
		AccessToken: fmt.Sprintf("real_token_%s_%d", scope, time.Now().Unix()),
		ExpiresIn:   3600, // 1 hora
	}, nil
}

func main() {
	// Configurar storage type via config
	viper.SetDefault("auth.storage_type", "keyring") // ou "memory" para Lambda
	
	// Inicializar serviço de auth
	storage := auth.NewStorageAdapter()
	authService := auth.NewAuthService(storage, tokenGenerator)

	// Verificar se está configurado
	if !authService.IsSetup() {
		fmt.Println("❌ Authentication not configured.")
		fmt.Println("Run: stackspot-cli auth user --id CLIENT_ID --secret CLIENT_SECRET")
		return
	}

	// Obter tokens para diferentes escopos
	executionToken, err := authService.GetValidToken(auth.ScopeExecution)
	if err != nil {
		log.Fatalf("Failed to get execution token: %v", err)
	}

	creationToken, err := authService.GetValidToken(auth.ScopeCreation)
	if err != nil {
		log.Fatalf("Failed to get creation token: %v", err)
	}

	fmt.Printf("✅ Execution token: %s...\n", executionToken[:20])
	fmt.Printf("✅ Creation token: %s...\n", creationToken[:20])
	
	// Continuar com o resto da aplicação...
}