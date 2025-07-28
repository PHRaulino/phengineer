package cli

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui"
	"github.com/spf13/cobra"
)

var authSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure authentication interactively",
	Long: `Configure authentication for different providers:
	
• StackSpot User (Client ID/Secret)
• StackSpot Service (via HashiCorp Vault + AWS)
• GitHub (Personal Access Token)

This command opens an interactive TUI for credential configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := tui.NewApp()
		return app.Start()
	},
}

// GetAuthSetupCmd returns the auth setup command for external use
func GetAuthSetupCmd() *cobra.Command {
	return authSetupCmd
}
