package auth

import (
	"errors"
	"sync"
)

type MemoryAdapter struct {
	data map[string]map[string]string
	mu   sync.RWMutex
}

func NewMemoryAdapter() *MemoryAdapter {
	return &MemoryAdapter{
		data: make(map[string]map[string]string),
	}
}

func (m *MemoryAdapter) Set(service, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data[service] == nil {
		m.data[service] = make(map[string]string)
	}
	m.data[service][key] = value
	return nil
}

func (m *MemoryAdapter) Get(service, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if serviceData, exists := m.data[service]; exists {
		if value, exists := serviceData[key]; exists {
			return value, nil
		}
	}
	return "", errors.New("key not found")
}

func (m *MemoryAdapter) Delete(service, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if serviceData, exists := m.data[service]; exists {
		delete(serviceData, key)
	}
	return nil
}

func (m *MemoryAdapter) Exists(service, key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if serviceData, exists := m.data[service]; exists {
		_, exists := serviceData[key]
		return exists
	}
	return false
}
