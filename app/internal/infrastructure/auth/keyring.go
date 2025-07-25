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
