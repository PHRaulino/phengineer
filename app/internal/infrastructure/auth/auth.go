package auth

import (
	"fmt"
	"time"

	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/token"
	"github.com/spf13/viper"
)

func SetupGenerators() {
	tokenService := token.GetService()

	tokenService.RegisterGenerator(token.TokenGenHC, func(scope token.TokenScope) (token.TokenResponse, error) {
		// Chama function AWS e hashicorp

		return token.TokenResponse{
			AccessToken: fmt.Sprintf("system_token_%s_%d", scope, time.Now().Unix()),
			ExpiresIn:   viper.GetInt("auth.token_expires"),
		}, nil
	})

	tokenService.RegisterGenerator(token.TokenGenSTK, func(scope token.TokenScope) (token.TokenResponse, error) {
		authMode := viper.GetString("auth.mode")

		switch authMode {
		case "stackspot_user":
			fmt.Println("client_id_user")
		case "stackspot_service":
			fmt.Println("client_id_service")
		default:
			fmt.Println("client_id_user")
		}

		// Chama a function do stackspot para gerar o token

		return token.TokenResponse{
			AccessToken: fmt.Sprintf("user_token_%s_%d", scope, time.Now().Unix()),
			ExpiresIn:   viper.GetInt("auth.token_expires"),
		}, nil
	})
}
