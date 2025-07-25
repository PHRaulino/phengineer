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