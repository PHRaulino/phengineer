package auth

// StorageAdapter interface para diferentes tipos de armazenamento
type StorageAdapter interface {
	Set(service, key, value string) error
	Get(service, key string) (string, error)
	Delete(service, key string) error
	Exists(service, key string) bool
}
