package common

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

type Logo struct {
	theme *styles.Theme
}

func NewLogo() *Logo {
	return &Logo{
		theme: styles.DefaultTheme,
	}
}

func (l *Logo) View() string {
	// ASCII art do logo
	logoArt := `
███████╗████████╗ █████╗  ██████╗██╗  ██╗███████╗██████╗  ██████╗ ████████╗
██╔════╝╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗██╔═══██╗╚══██╔══╝
███████╗   ██║   ███████║██║     █████╔╝ ███████╗██████╔╝██║   ██║   ██║   
╚════██║   ██║   ██╔══██║██║     ██╔═██╗ ╚════██║██╔═══╝ ██║   ██║   ██║   
███████║   ██║   ██║  ██║╚██████╗██║  ██╗███████║██║     ╚██████╔╝   ██║   
╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝      ╚═════╝    ╚═╝   
`

	style := lipgloss.NewStyle().
		Foreground(l.theme.Primary).
		Bold(true).
		Align(lipgloss.Center)

	return style.Render(logoArt)
}

// SmallLogo para espaços menores
func (l *Logo) SmallView() string {
	smallLogo := "STACKSPOT CLI"

	style := lipgloss.NewStyle().
		Foreground(l.theme.Primary).
		Bold(true).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(l.theme.Primary).
		Padding(0, 2)

	return style.Render(smallLogo)
}
