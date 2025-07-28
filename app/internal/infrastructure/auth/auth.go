package auth

import (
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/providers"
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/storage"
	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/token"
	"github.com/spf13/viper"
)

func SetupGenerators() {
	tokenService := token.GetService()
	
	// Inicializar storage (keyring por padrão)
	authStorage := storage.NewKeyringAdapter()

	// Provider para HashiCorp Vault + AWS
	vaultProvider := providers.NewVaultProvider(authStorage)
	tokenService.RegisterGenerator(token.TokenGenHC, func(scope token.TokenScope) (token.TokenResponse, error) {
		return vaultProvider.GetToken(scope)
	})

	// Provider para StackSpot
	stackspotProvider := providers.NewStackSpotProvider(authStorage)
	tokenService.RegisterGenerator(token.TokenGenSTK, func(scope token.TokenScope) (token.TokenResponse, error) {
		authMode := viper.GetString("auth.mode")
		
		switch authMode {
		case "stackspot_service":
			// Usar Vault para buscar credenciais
			return vaultProvider.GetToken(scope)
		case "stackspot_user":
			fallthrough
		default:
			// Usar credenciais diretas do usuário
			return stackspotProvider.GetToken(scope)
		}
	})

	// Provider para GitHub
	githubProvider := providers.NewGitHubProvider(authStorage)
	tokenService.RegisterGenerator(token.TokenGenGH, func(scope token.TokenScope) (token.TokenResponse, error) {
		return githubProvider.GetToken(scope)
	})
}

// GetStackSpotProvider retorna uma instância do provider StackSpot
func GetStackSpotProvider() *providers.StackSpotProvider {
	authStorage := storage.NewKeyringAdapter()
	return providers.NewStackSpotProvider(authStorage)
}

// GetVaultProvider retorna uma instância do provider Vault
func GetVaultProvider() *providers.VaultProvider {
	authStorage := storage.NewKeyringAdapter()
	return providers.NewVaultProvider(authStorage)
}

// GetGitHubProvider retorna uma instância do provider GitHub
func GetGitHubProvider() *providers.GitHubProvider {
	authStorage := storage.NewKeyringAdapter()
	return providers.NewGitHubProvider(authStorage)
}
