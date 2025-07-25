package token

import (
	"sync"

	"github.com/PHRaulino/phengineer/internal/infrastructure/auth/storage"
	"github.com/spf13/viper"
)

var (
	globalService *Service
	once          sync.Once
)

// GetService retorna a instância global do token service
func GetService() *Service {
	once.Do(func() {
		var adapter storage.StorageAdapter

		authMode := viper.GetString("auth.mode")
		switch authMode {
		case "stackspot_user":
			adapter = storage.NewKeyringAdapter()
		case "stackspot_service":
			adapter = storage.NewMemoryAdapter()
		default:
			// Default para user
			adapter = storage.NewKeyringAdapter()
		}

		globalService = NewService(adapter)
	})
	return globalService
}
