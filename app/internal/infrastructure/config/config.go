package config

import (
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() error {
	slug := "phengineer"

	repoName, _, rootPath := detectGitInfo()

	viper.SetConfigName(slug)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Defaults
	viper.SetDefault("project.type", "lambda")
	viper.SetDefault("project.repo_name", repoName)
	viper.SetDefault("project.root_path", rootPath)

	viper.SetDefault("project.language.name", "python")
	viper.SetDefault("project.language.version", "3.13")

	viper.SetDefault("analysis.file_limits.max_file_size", "10MB")
	viper.SetDefault("analysis.file_limits.max_files", 100)

	// Env vars
	viper.SetEnvPrefix(strings.ToUpper(slug))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return viper.ReadInConfig()
}
