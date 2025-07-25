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