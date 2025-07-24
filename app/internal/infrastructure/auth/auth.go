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

// auth/client.go
package auth

import (
	"errors"
	"fmt"
)

const (
	ServiceName    = "stackspot-cli"
	ClientIDKey    = "client_id"
	ClientSecretKey = "client_secret"
	TokenKey       = "access_token"
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
	"errors"
	"fmt"
	"time"
)

type TokenService struct {
	storage StorageAdapter
}

func NewTokenService(storage StorageAdapter) *TokenService {
	return &TokenService{
		storage: storage,
	}
}

// SaveToken salva o token de acesso
func (t *TokenService) SaveToken(token string) error {
	if token == "" {
		return errors.New("token cannot be empty")
	}

	return t.storage.Set(ServiceName, TokenKey, token)
}

// GetToken recupera o token de acesso
func (t *TokenService) GetToken() (string, error) {
	token, err := t.storage.Get(ServiceName, TokenKey)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	return token, nil
}

// HasToken verifica se o token existe
func (t *TokenService) HasToken() bool {
	return t.storage.Exists(ServiceName, TokenKey)
}

// DeleteToken remove o token
func (t *TokenService) DeleteToken() error {
	return t.storage.Delete(ServiceName, TokenKey)
}

// ValidateToken verifica se o token é válido (placeholder)
// Você deve implementar a lógica específica da sua API
func (t *TokenService) ValidateToken(token string) bool {
	// TODO: Implementar validação real do token
	// Por exemplo: fazer request para endpoint de validação
	// Por enquanto, verifica apenas se não está vazio
	return token != ""
}

// auth/service.go
package auth

import (
	"fmt"
)

type AuthService struct {
	clientService *ClientService
	tokenService  *TokenService
}

func NewAuthService(storage StorageAdapter) *AuthService {
	return &AuthService{
		clientService: NewClientService(storage),
		tokenService:  NewTokenService(storage),
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

// GetValidToken retorna um token válido (gera novo se necessário)
func (a *AuthService) GetValidToken() (string, error) {
	// Verifica se existe token
	if !a.tokenService.HasToken() {
		return a.generateNewToken()
	}

	// Recupera token existente
	token, err := a.tokenService.GetToken()
	if err != nil {
		return a.generateNewToken()
	}

	// Valida token
	if !a.tokenService.ValidateToken(token) {
		return a.generateNewToken()
	}

	return token, nil
}

// generateNewToken gera um novo token usando as credenciais
func (a *AuthService) generateNewToken() (string, error) {
	clientID, clientSecret, err := a.clientService.GetCredentials()
	if err != nil {
		return "", fmt.Errorf("credentials not found: %w", err)
	}

	// TODO: Implementar sua função de geração de token aqui
	token, err := a.callTokenAPI(clientID, clientSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Salva o novo token
	if err := a.tokenService.SaveToken(token); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	return token, nil
}

// callTokenAPI chama sua API para gerar token (placeholder)
func (a *AuthService) callTokenAPI(clientID, clientSecret string) (string, error) {
	// TODO: Implementar chamada real para sua API
	// Por enquanto retorna um token fictício
	return fmt.Sprintf("token_%d", time.Now().Unix()), nil
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
	storage := auth.NewKeyringAdapter()
	authService := auth.NewAuthService(storage)

	// Salvar credenciais
	if err := authService.SetupCredentials(clientID, clientSecret); err != nil {
		return fmt.Errorf("failed to setup credentials: %w", err)
	}

	fmt.Println("✅ Credentials saved successfully!")
	return nil
}

// main.go ou onde você inicializa
package main

import (
	"fmt"
	"log"
	"your-project/auth"
)

func main() {
	// Inicializar serviço de auth
	storage := auth.NewKeyringAdapter()
	authService := auth.NewAuthService(storage)

	// Verificar se está configurado
	if !authService.IsSetup() {
		fmt.Println("❌ Authentication not configured.")
		fmt.Println("Run: stackspot-cli auth user --id CLIENT_ID --secret CLIENT_SECRET")
		return
	}

	// Obter token válido
	token, err := authService.GetValidToken()
	if err != nil {
		log.Fatalf("Failed to get valid token: %v", err)
	}

	fmt.Printf("✅ Token obtained: %s\n", token[:20]+"...")
	
	// Continuar com o resto da aplicação...
}