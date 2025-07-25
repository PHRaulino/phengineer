package storage

import (
	"github.com/zalando/go-keyring"
)

const ServiceName = "phengineer"

type KeyringAdapter struct{}

func NewKeyringAdapter() *KeyringAdapter {
	return &KeyringAdapter{}
}

func (k *KeyringAdapter) Set(key, value string) error {
	return keyring.Set(ServiceName, key, value)
}

func (k *KeyringAdapter) Get(key string) (string, error) {
	return keyring.Get(ServiceName, key)
}

func (k *KeyringAdapter) Delete(key string) error {
	return keyring.Delete(ServiceName, key)
}

func (k *KeyringAdapter) Exists(key string) bool {
	_, err := keyring.Get(ServiceName, key)
	return err == nil
}