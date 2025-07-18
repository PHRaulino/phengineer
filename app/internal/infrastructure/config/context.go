package config

import (
	"context"
	"fmt"
	"path/filepath"
)

// configKey é a chave para armazenar a config no context
type configKey struct{}

// WithConfig adiciona a configuração ao context
func WithConfig(ctx context.Context, configFolderName string) (context.Context, error) {
	// Valida requirements primeiro
	validator := NewRequirementsValidator(configFolderName)
	if err := validator.ValidateAll(); err != nil {
		return nil, fmt.Errorf("requirements validation failed: %w", err)
	}

	// Coleta configurações automáticas
	autoConfig, err := collectAutoConfig(configFolderName)
	if err != nil {
		return nil, fmt.Errorf("failed to collect auto config: %w", err)
	}

	// Carrega configurações do settings.yml
	settingsPath := filepath.Join(autoConfig.ConfigDirPath, "settings.yml")
	settings, err := LoadOrCreateSettings(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	// Cria a configuração completa
	config := &Config{
		Settings:   settings,
		Auto:       autoConfig,
		ConfigPath: autoConfig.ConfigDirPath,
	}

	return context.WithValue(ctx, configKey{}, config), nil
}

// WithConfigVerbose adiciona a configuração ao context com saída detalhada
func WithConfigVerbose(ctx context.Context, configFolderName string) (context.Context, error) {
	fmt.Println("🔍 Validating requirements...")

	// Valida requirements com detalhes
	validator := NewRequirementsValidator(configFolderName)
	results, err := validator.ValidateWithDetails()

	// Mostra o status de cada requirement
	for _, result := range results {
		status := "✅"
		if !result.Passed {
			status = "❌"
		}
		fmt.Printf("  %s %s: %s\n", status, result.Name, result.Description)

		if !result.Passed {
			fmt.Printf("    Error: %s\n", result.ErrorMsg)
		}
	}

	if err != nil {
		fmt.Printf("\n❌ Requirements validation failed: %v\n", err)
		return nil, err
	}

	fmt.Println("\n📋 Collecting auto configurations...")

	// Coleta configurações automáticas
	autoConfig, err := collectAutoConfig(configFolderName)
	if err != nil {
		fmt.Printf("❌ Failed to collect auto config: %v\n", err)
		return nil, err
	}

	fmt.Printf("  📱 App Name: %s\n", autoConfig.AppName)
	fmt.Printf("  📁 Config Dir: %s\n", autoConfig.ConfigDirPath)
	fmt.Printf("  🌐 Remote URL: %s\n", autoConfig.RemoteURL)

	fmt.Println("\n📄 Loading settings.yml...")

	// Carrega configurações do settings.yml
	settingsPath := filepath.Join(autoConfig.ConfigDirPath, "settings.yml")
	settings, err := LoadOrCreateSettings(settingsPath)
	if err != nil {
		fmt.Printf("❌ Failed to load settings: %v\n", err)
		return nil, err
	}

	fmt.Printf("  ⚙️  Project Type: %s\n", settings.Project.Type)
	fmt.Printf("  💻 Language: %s %s\n", settings.Project.Language.Name, settings.Project.Language.Version)
	fmt.Printf("  📊 Include Path: %s\n", settings.Analysis.FilesIncludePath)
	fmt.Printf("  🚫 Exclude Path: %s\n", settings.Analysis.FilesExcludePath)

	// Cria a configuração completa
	config := &Config{
		Settings:   settings,
		Auto:       autoConfig,
		ConfigPath: autoConfig.ConfigDirPath,
	}

	fmt.Println("\n✅ Configuration loaded successfully!")

	return context.WithValue(ctx, configKey{}, config), nil
}

// FromContext extrai a configuração do context
func FromContext(ctx context.Context) *Config {
	if config, ok := ctx.Value(configKey{}).(*Config); ok {
		return config
	}
	panic("config not found in context - make sure to call config.WithConfig() first")
}

// MustFromContext extrai a configuração do context ou retorna erro
func MustFromContext(ctx context.Context) (*Config, error) {
	if config, ok := ctx.Value(configKey{}).(*Config); ok {
		return config, nil
	}
	return nil, fmt.Errorf("config not found in context")
}

// HasConfig verifica se o context possui configuração
func HasConfig(ctx context.Context) bool {
	_, ok := ctx.Value(configKey{}).(*Config)
	return ok
}

// GetSettings extrai apenas as configurações do settings.yml do context
func GetSettings(ctx context.Context) *Settings {
	return FromContext(ctx).Settings
}

// GetAutoConfig extrai apenas as configurações automáticas do context
func GetAutoConfig(ctx context.Context) *AutoConfig {
	return FromContext(ctx).Auto
}

// PrintDiagnostics imprime informações de diagnóstico da configuração
func PrintDiagnostics(ctx context.Context) {
	config := FromContext(ctx)

	fmt.Println("📋 Configuration Diagnostics")
	fmt.Printf("  📱 App Name: %s\n", config.Auto.AppName)
	fmt.Printf("  📁 Config Path: %s\n", config.ConfigPath)
	fmt.Printf("  🌐 Remote URL: %s\n", config.Auto.RemoteURL)

	fmt.Println("\n⚙️  Project Settings:")
	fmt.Printf("  Type: %s\n", config.Settings.Project.Type)
	fmt.Printf("  Language: %s %s\n", config.Settings.Project.Language.Name, config.Settings.Project.Language.Version)

	fmt.Println("\n📊 Analysis Settings:")
	fmt.Printf("  Include: %s\n", config.Settings.Analysis.FilesIncludePath)
	fmt.Printf("  Exclude: %s\n", config.Settings.Analysis.FilesExcludePath)
	fmt.Printf("  Max File Size: %s\n", config.Settings.Analysis.FileLimits.MaxFileSize)
	fmt.Printf("  Max Files: %d\n", config.Settings.Analysis.FileLimits.MaxFiles)
}
