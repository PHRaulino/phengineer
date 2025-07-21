package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadSettingsFromFile carrega as configurações de um arquivo YAML
func LoadSettingsFromFile(filePath string) (*Settings, error) {
	// Verifica se o arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("settings file not found: %s", filePath)
	}

	// Lê o arquivo
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	// Parse do YAML
	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings YAML: %w", err)
	}

	// Valida as configurações
	if err := settings.Validate(); err != nil {
		return nil, fmt.Errorf("settings validation failed: %w", err)
	}

	return &settings, nil
}

// SaveSettingsToFile salva as configurações em um arquivo YAML
func SaveSettingsToFile(settings *Settings, filePath string) error {
	// Valida as configurações antes de salvar
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("settings validation failed: %w", err)
	}

	// Cria o diretório se não existir
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Converte para YAML
	data, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings to YAML: %w", err)
	}

	// Escreve o arquivo
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// LoadOrCreateSettings carrega as configurações ou cria um arquivo padrão se não existir
func LoadOrCreateSettings(configFolderPath, configFolderName string) (*Settings, error) {

	settingsPath := filepath.Join(configFolderPath, "settings.yml")
	// Tenta carregar o arquivo existente
	if _, err := os.Stat(settingsPath); err == nil {
		return LoadSettingsFromFile(settingsPath)
	}

	// Arquivo não existe, cria um padrão
	fmt.Printf("Settings file not found, creating default: %s\n", settingsPath)

	defaultSettings := GetDefaultSettings(configFolderName)
	if err := SaveSettingsToFile(defaultSettings, settingsPath); err != nil {
		return nil, fmt.Errorf("failed to create default settings: %w", err)
	}

	return defaultSettings, nil
}