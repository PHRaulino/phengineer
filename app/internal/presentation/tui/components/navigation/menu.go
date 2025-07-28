package navigation

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	Title       string
	Description string
	Icon        string
	Action      tea.Cmd
}

type Menu struct {
	items    []MenuItem
	cursor   int
	theme    *styles.Theme
	width    int
	showHelp bool
}

func (m *Menu) View() string {
	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border).
		Padding(1, 2).
		Width(m.width)

	var items []string
	for i, item := range m.items {
		cursor := "  "
		style := lipgloss.NewStyle()

		if i == m.cursor {
			cursor = "▸ "
			style = style.
				Foreground(m.theme.Primary).
				Bold(true)
		}

		icon := item.Icon
		if icon != "" {
			icon += " "
		}

		title := style.Render(icon + item.Title)
		desc := lipgloss.NewStyle().
			Foreground(m.theme.Muted).
			Render(item.Description)

		items = append(items, cursor+title+"\n    "+desc)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)

	if m.showHelp {
		help := lipgloss.NewStyle().
			Foreground(m.theme.Muted).
			MarginTop(1).
			Render("↑/↓ navegar • enter selecionar • q sair")
		content = lipgloss.JoinVertical(lipgloss.Left, content, help)
	}

	return menuStyle.Render(content)
}
