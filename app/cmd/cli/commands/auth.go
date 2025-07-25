package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
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
