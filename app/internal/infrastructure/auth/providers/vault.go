package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/storage"
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/token"
)

type VaultProvider struct {
	storage storage.StorageAdapter
}

func NewVaultProvider(storage storage.StorageAdapter) *VaultProvider {
	return &VaultProvider{
		storage: storage,
	}
}

type VaultAWSAuthRequest struct {
	Role string `json:"role"`
}

type VaultAWSAuthResponse struct {
	Auth struct {
		ClientToken string `json:"client_token"`
		LeaseDuration int  `json:"lease_duration"`
		Renewable     bool `json:"renewable"`
	} `json:"auth"`
}

type VaultSecretResponse struct {
	Data struct {
		Data map[string]interface{} `json:"data"`
	} `json:"data"`
}

func (p *VaultProvider) GetToken(scope token.TokenScope) (token.TokenResponse, error) {
	// Buscar configurações do Vault
	vaultURL, err := p.storage.Get("vault_url")
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("vault_url não encontrado: %w", err)
	}

	vaultRole, err := p.storage.Get("vault_aws_role")
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("vault_aws_role não encontrado: %w", err)
	}

	// Autenticar com Vault usando AWS IAM
	vaultToken, err := p.authenticateWithAWS(vaultURL, vaultRole)
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("erro na autenticação AWS/Vault: %w", err)
	}

	// Buscar credenciais StackSpot no Vault
	stackSpotCreds, err := p.getStackSpotCredentials(vaultURL, vaultToken)
	if err != nil {
		return token.TokenResponse{}, fmt.Errorf("erro ao buscar credenciais StackSpot: %w", err)
	}

	// Usar provider StackSpot para gerar token
	stackSpotProvider := NewStackSpotProvider(p.storage)
	
	// Salvar temporariamente as credenciais
	if err := stackSpotProvider.SaveCredentials(
		stackSpotCreds["client_id"].(string),
		stackSpotCreds["client_secret"].(string),
	); err != nil {
		return token.TokenResponse{}, fmt.Errorf("erro ao salvar credenciais temporárias: %w", err)
	}

	return stackSpotProvider.GetToken(scope)
}

func (p *VaultProvider) authenticateWithAWS(vaultURL, role string) (string, error) {
	// Verificar se estamos em ambiente AWS (EC2, ECS, Lambda, etc.)
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = os.Getenv("AWS_DEFAULT_REGION")
	}
	if awsRegion == "" {
		awsRegion = "us-east-1" // fallback
	}

	// Para simplificar, vamos usar token IAM role-based
	// Em produção, seria necessário implementar assinatura AWS SigV4
	client := &http.Client{Timeout: 30 * time.Second}
	
	// Endpoint para auth AWS no Vault
	authURL := fmt.Sprintf("%s/v1/auth/aws/login", strings.TrimSuffix(vaultURL, "/"))
	
	authReq := VaultAWSAuthRequest{
		Role: role,
	}

	reqBody, _ := json.Marshal(authReq)
	req, err := http.NewRequest("POST", authURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição de auth: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na requisição de auth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("falha na autenticação AWS: status %d", resp.StatusCode)
	}

	var authResp VaultAWSAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta de auth: %w", err)
	}

	return authResp.Auth.ClientToken, nil
}

func (p *VaultProvider) getStackSpotCredentials(vaultURL, vaultToken string) (map[string]interface{}, error) {
	secretPath, err := p.storage.Get("vault_stackspot_path")
	if err != nil {
		// Usar path padrão se não configurado
		secretPath = "secret/data/stackspot"
	}

	client := &http.Client{Timeout: 30 * time.Second}
	
	secretURL := fmt.Sprintf("%s/v1/%s", strings.TrimSuffix(vaultURL, "/"), secretPath)
	
	req, err := http.NewRequest("GET", secretURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição de secret: %w", err)
	}

	req.Header.Set("X-Vault-Token", vaultToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição de secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao buscar secret: status %d", resp.StatusCode)
	}

	var secretResp VaultSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&secretResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar secret: %w", err)
	}

	return secretResp.Data.Data, nil
}

func (p *VaultProvider) SaveConfig(vaultURL, awsRole, stackspotPath string) error {
	if err := p.storage.Set("vault_url", vaultURL); err != nil {
		return fmt.Errorf("erro ao salvar vault_url: %w", err)
	}

	if err := p.storage.Set("vault_aws_role", awsRole); err != nil {
		return fmt.Errorf("erro ao salvar vault_aws_role: %w", err)
	}

	if stackspotPath != "" {
		if err := p.storage.Set("vault_stackspot_path", stackspotPath); err != nil {
			return fmt.Errorf("erro ao salvar vault_stackspot_path: %w", err)
		}
	}

	return nil
}