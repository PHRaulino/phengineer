package common

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	spinner spinner.Model
	message string
	theme   *styles.Theme
}

func NewSpinner(message string) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.DefaultTheme.Primary)

	return &Spinner{
		spinner: s,
		message: message,
		theme:   styles.DefaultTheme,
	}
}

func (s *Spinner) SetMessage(message string) {
	s.message = message
}

func (s *Spinner) SetTheme(theme *styles.Theme) {
	s.theme = theme
	s.spinner.Style = lipgloss.NewStyle().Foreground(theme.Primary)
}

func (s *Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

func (s *Spinner) View() string {
	messageStyle := lipgloss.NewStyle().
		Foreground(s.theme.Foreground).
		MarginLeft(1)

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		s.spinner.View(),
		messageStyle.Render(s.message),
	)
}

// Spinner presets
var (
	SpinnerDots    = spinner.Dot
	SpinnerLine    = spinner.Line
	SpinnerMiniDot = spinner.MiniDot
	SpinnerJump    = spinner.Jump
	SpinnerPulse   = spinner.Pulse
	SpinnerPoints  = spinner.Points
	SpinnerGlobe   = spinner.Globe
	SpinnerMoon    = spinner.Moon
	SpinnerMonkey  = spinner.Monkey
)

// LoadingSpinner é um spinner pronto para uso
func LoadingSpinner(message string) *Spinner {
	return NewSpinner(message)
}

// ProcessingSpinner com estilo diferente
func ProcessingSpinner(message string) *Spinner {
	s := NewSpinner(message)
	s.spinner.Spinner = spinner.Points
	return s
}
