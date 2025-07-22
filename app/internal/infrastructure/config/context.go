package config

import (
	"context"

	"go.uber.org/zap"
)

// configKey é a chave para armazenar a config no context
type configKey struct{}

// WithConfig adiciona a configuração ao context com saída detalhada
func WithConfig(ctx context.Context, configFolderName string) (context.Context, error) {
	// Valida requirements com detalhes
	validator := NewRequirementsValidator(configFolderName)
	results, err := validator.ValidateWithDetails()
	// Mostra o status de cada requirement
	if err != nil {
		zap.L().Error("requirements validation failed", zap.Error(err))

		for _, result := range results {
			status := "OK"
			if !result.Passed {
				status = "ERROR"
			}
			zap.L().Info("requirements",
				zap.String("status", status),
				zap.String("name", result.Name),
				zap.String("description", result.Description),
			)
		}
		return nil, err
	}

	// Coleta configurações automáticas
	autoConfig, err := collectAutoConfig(configFolderName)
	if err != nil {
		zap.L().Error("failed to collect auto config", zap.Error(err))
		return nil, err
	}

	// Carrega configurações do settings.yml

	settings, err := LoadOrCreateSettings(autoConfig.ConfigDirPath, configFolderName)
	if err != nil {
		zap.L().Error("failed to load settings", zap.Error(err))
		return nil, err
	}

	// Cria a configuração completa
	config := &Config{
		Settings:   settings,
		Auto:       autoConfig,
		ConfigPath: autoConfig.ConfigDirPath,
	}

	return context.WithValue(ctx, configKey{}, config), nil
}

// FromContext extrai a configuração do context
func FromContext(ctx context.Context) *Config {
	if config, ok := ctx.Value(configKey{}).(*Config); ok {
		return config
	}
	zap.L().Panic("config not found in context - make sure to call config.WithConfig() first")
	return nil
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

	zap.L().Debug("config", zap.String("key", "app_name"), zap.String("value", config.Auto.AppName))
	zap.L().Debug("config", zap.String("key", "config_path"), zap.String("value", config.ConfigPath))
	zap.L().Debug("config", zap.String("key", "remote_url"), zap.String("value", config.Auto.RemoteURL))
	zap.L().Debug("config", zap.String("key", "root_app_path"), zap.String("value", config.Auto.RootAppPath))

	zap.L().Debug("config", zap.String("key", "project_type"), zap.String("value", config.Settings.Project.Type))
	zap.L().Debug("config", zap.String("key", "language"), zap.String("value", config.Settings.Project.Language.Name))
	zap.L().Debug("config", zap.String("key", "language_version"), zap.String("value", config.Settings.Project.Language.Version))

	zap.L().Debug("config", zap.String("key", "analyse_files"), zap.String("value", config.Settings.Analysis.AnalysisFilesPath))
	zap.L().Debug("config", zap.String("key", "max_file_size"), zap.String("value", config.Settings.Analysis.FileLimits.MaxFileSize))
	zap.L().Debug("config", zap.String("key", "max_files"), zap.Int64("value", config.Settings.Analysis.FileLimits.MaxFiles))
}
