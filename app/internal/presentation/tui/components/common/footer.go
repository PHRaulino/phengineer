package common

import (
	"strings"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FooterItem struct {
	Key   string
	Label string
}

type Footer struct {
	items []FooterItem
	width int
	theme *styles.Theme
}

func NewFooter() *Footer {
	return &Footer{
		theme: styles.DefaultTheme,
		items: []FooterItem{
			{Key: "↑/↓", Label: "Navigate"},
			{Key: "Enter", Label: "Select"},
			{Key: "Esc", Label: "Back"},
			{Key: "q", Label: "Quit"},
		},
	}
}

func (f *Footer) SetItems(items []FooterItem) *Footer {
	f.items = items
	return f
}

func (f *Footer) SetWidth(width int) {
	f.width = width
}

func (f *Footer) SetTheme(theme *styles.Theme) {
	f.theme = theme
}

func (f *Footer) Init() tea.Cmd {
	return nil
}

func (f *Footer) View() string {
	containerStyle := lipgloss.NewStyle().
		Width(f.width).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(f.theme.Border).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().
		Foreground(f.theme.Primary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(f.theme.Muted)

	separatorStyle := lipgloss.NewStyle().
		Foreground(f.theme.Border).
		Padding(0, 1)

	// Construir items
	var parts []string
	for _, item := range f.items {
		part := keyStyle.Render(item.Key) + " " + labelStyle.Render(item.Label)
		parts = append(parts, part)
	}

	content := strings.Join(parts, separatorStyle.Render("•"))

	// Centralizar
	contentWidth := lipgloss.Width(content)
	if contentWidth < f.width-2 {
		padding := (f.width - 2 - contentWidth) / 2
		content = strings.Repeat(" ", padding) + content
	}

	return containerStyle.Render(content)
}
