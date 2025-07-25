package storage

import (
	"errors"
	"sync"
)

type MemoryAdapter struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryAdapter() *MemoryAdapter {
	return &MemoryAdapter{
		data: make(map[string]string),
	}
}

func (m *MemoryAdapter) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *MemoryAdapter) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if value, exists := m.data[key]; exists {
		return value, nil
	}
	return "", errors.New("key not found")
}

func (m *MemoryAdapter) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryAdapter) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.data[key]
	return exists
}