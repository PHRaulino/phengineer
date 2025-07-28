package common

import (
	"strings"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type KeyBinding struct {
	keys key.Binding
	help string
}

type Help struct {
	bindings   []KeyBinding
	expanded   bool
	theme      *styles.Theme
	showToggle bool
}

func NewHelp() *Help {
	return &Help{
		bindings:   make([]KeyBinding, 0),
		expanded:   false,
		theme:      styles.DefaultTheme,
		showToggle: true,
	}
}

func (h *Help) AddBinding(keys, help string) *Help {
	h.bindings = append(h.bindings, KeyBinding{
		keys: key.NewBinding(key.WithKeys(keys)),
		help: help,
	})
	return h
}

func (h *Help) AddBindings(bindings ...KeyBinding) *Help {
	h.bindings = append(h.bindings, bindings...)
	return h
}

func (h *Help) Toggle() {
	h.expanded = !h.expanded
}

func (h *Help) View() string {
	if len(h.bindings) == 0 {
		return ""
	}

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(h.theme.Border).
		Padding(0, 1).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(h.theme.Primary).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(h.theme.Muted)

	titleStyle := lipgloss.NewStyle().
		Foreground(h.theme.Foreground).
		Bold(true).
		MarginBottom(1)

	var content string

	if h.expanded {
		content = titleStyle.Render("🔑 Keyboard Shortcuts")

		maxKeyWidth := 0
		for _, b := range h.bindings {
			if w := lipgloss.Width(b.keys.Keys()[0]); w > maxKeyWidth {
				maxKeyWidth = w
			}
		}

		for _, b := range h.bindings {
			key := keyStyle.Width(maxKeyWidth + 2).Render(b.keys.Keys()[0])
			help := helpStyle.Render(b.help)
			content += "\n" + key + help
		}

		if h.showToggle {
			content += "\n\n" + helpStyle.Render("Press ? to hide")
		}
	} else {
		// Vista compacta
		var items []string
		for i, b := range h.bindings {
			if i >= 3 { // Mostrar solo los primeros 3
				items = append(items, "...")
				break
			}
			item := keyStyle.Render(b.keys.Keys()[0]) + " " + helpStyle.Render(b.help)
			items = append(items, item)
		}

		content = strings.Join(items, " • ")

		if h.showToggle {
			content += " • " + keyStyle.Render("?") + " " + helpStyle.Render("more")
		}
	}

	return containerStyle.Render(content)
}

// Atalhos comuns pré-definidos
var (
	NavigationBindings = []KeyBinding{
		{keys: key.NewBinding(key.WithKeys("up", "k")), help: "Move up"},
		{keys: key.NewBinding(key.WithKeys("down", "j")), help: "Move down"},
		{keys: key.NewBinding(key.WithKeys("left", "h")), help: "Move left"},
		{keys: key.NewBinding(key.WithKeys("right", "l")), help: "Move right"},
	}

	CommonBindings = []KeyBinding{
		{keys: key.NewBinding(key.WithKeys("enter")), help: "Select"},
		{keys: key.NewBinding(key.WithKeys("esc")), help: "Back"},
		{keys: key.NewBinding(key.WithKeys("q", "ctrl+c")), help: "Quit"},
		{keys: key.NewBinding(key.WithKeys("?")), help: "Toggle help"},
	}
)
