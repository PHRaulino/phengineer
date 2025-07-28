package styles

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	Border     lipgloss.Color
	Muted      lipgloss.Color
	Highlight  lipgloss.Color
}

var (
	// Tema padrão moderno
	DefaultTheme = &Theme{
		Name:       "Modern Dark",
		Primary:    lipgloss.Color("#7C3AED"), // Roxo vibrante
		Secondary:  lipgloss.Color("#10B981"), // Verde esmeralda
		Success:    lipgloss.Color("#10B981"),
		Warning:    lipgloss.Color("#F59E0B"), // Âmbar
		Error:      lipgloss.Color("#EF4444"), // Vermelho
		Background: lipgloss.Color("#111827"), // Gray-900
		Foreground: lipgloss.Color("#F9FAFB"), // Gray-50
		Border:     lipgloss.Color("#4B5563"), // Gray-600
		Muted:      lipgloss.Color("#6B7280"), // Gray-500
		Highlight:  lipgloss.Color("#8B5CF6"), // Violet-500
	}

	// Tema claro alternativo
	LightTheme = &Theme{
		Name:       "Modern Light",
		Primary:    lipgloss.Color("#7C3AED"),
		Secondary:  lipgloss.Color("#059669"),
		Success:    lipgloss.Color("#059669"),
		Warning:    lipgloss.Color("#D97706"),
		Error:      lipgloss.Color("#DC2626"),
		Background: lipgloss.Color("#FFFFFF"),
		Foreground: lipgloss.Color("#111827"),
		Border:     lipgloss.Color("#E5E7EB"),
		Muted:      lipgloss.Color("#9CA3AF"),
		Highlight:  lipgloss.Color("#A78BFA"),
	}
)

// Styles contém estilos pré-definidos
type Styles struct {
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Box      lipgloss.Style
	Success  lipgloss.Style
	Error    lipgloss.Style
	Warning  lipgloss.Style
	Info     lipgloss.Style
	Button   lipgloss.Style
	Input    lipgloss.Style
	Border   lipgloss.Style
}

// Estilos pré-definidos baseados no tema
func (t *Theme) GetStyles() *Styles {
	return &Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(t.Muted).
			Italic(true),

		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Padding(1, 2),

		Success: lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(t.Error).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(t.Warning).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(t.Secondary),

		Button: lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(t.Border),

		Input: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(t.Border).
			Padding(0, 1),

		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border),
	}
}
