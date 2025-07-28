package display

import (
	"fmt"
	"strings"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

type ProgressBar struct {
	Width       int
	Current     int
	Total       int
	Theme       *styles.Theme
	ShowPercent bool
	ShowNumbers bool
	Label       string
}

func (p *ProgressBar) View() string {
	if p.Total == 0 {
		return ""
	}

	percent := float64(p.Current) / float64(p.Total)
	filled := int(percent * float64(p.Width))

	// Estilos
	filledStyle := lipgloss.NewStyle().
		Foreground(p.Theme.Primary).
		Background(p.Theme.Primary)

	emptyStyle := lipgloss.NewStyle().
		Foreground(p.Theme.Muted)

	// Barra
	bar := strings.Repeat("█", filled) +
		strings.Repeat("░", p.Width-filled)

	// Renderizar
	var parts []string

	if p.Label != "" {
		parts = append(parts, p.Label)
	}

	parts = append(parts,
		filledStyle.Render(bar[:filled])+
			emptyStyle.Render(bar[filled:]))

	if p.ShowPercent {
		parts = append(parts, fmt.Sprintf("%.0f%%", percent*100))
	}

	if p.ShowNumbers {
		parts = append(parts, fmt.Sprintf("%d/%d", p.Current, p.Total))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}
