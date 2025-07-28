package common

import (
	"strings"
	"time"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	title    string
	subtitle string
	showTime bool
	width    int
	theme    *styles.Theme
}

func NewHeader(title, subtitle string) *Header {
	return &Header{
		title:    title,
		subtitle: subtitle,
		showTime: true,
		theme:    styles.DefaultTheme,
	}
}

func (h *Header) SetWidth(width int) {
	h.width = width
}

func (h *Header) SetTheme(theme *styles.Theme) {
	h.theme = theme
}

func (h *Header) Init() tea.Cmd {
	if h.showTime {
		return tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return timeUpdateMsg(t)
		})
	}
	return nil
}

func (h *Header) Update(msg tea.Msg) (*Header, tea.Cmd) {
	switch msg.(type) {
	case timeUpdateMsg:
		return h, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return timeUpdateMsg(t)
		})
	}
	return h, nil
}

func (h *Header) View() string {
	// Estilos
	containerStyle := lipgloss.NewStyle().
		Width(h.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(h.theme.Border).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(h.theme.Primary)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(h.theme.Muted).
		Italic(true)

	timeStyle := lipgloss.NewStyle().
		Foreground(h.theme.Muted)

	// Conteúdo
	left := titleStyle.Render(h.title)
	if h.subtitle != "" {
		left += " " + subtitleStyle.Render(h.subtitle)
	}

	right := ""
	if h.showTime {
		right = timeStyle.Render(time.Now().Format("15:04:05"))
	}

	// Layout
	gap := h.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 0 {
		gap = 0
	}

	content := left + strings.Repeat(" ", gap) + right

	return containerStyle.Render(content)
}

type timeUpdateMsg time.Time
