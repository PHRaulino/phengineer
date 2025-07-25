package storage

// StorageAdapter interface para diferentes tipos de armazenamento
type StorageAdapter interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) bool
}
