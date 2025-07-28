package forms

import (
	"strings"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	BaseField
	textInput textinput.Model
	multiline bool
}

func NewInput() *Input {
	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 256

	return &Input{
		BaseField: BaseField{
			theme: styles.DefaultTheme,
			width: 40,
		},
		textInput: ti,
	}
}

// Builder pattern methods
func (i *Input) WithPlaceholder(placeholder string) *Input {
	i.placeholder = placeholder
	i.textInput.Placeholder = placeholder
	return i
}

func (i *Input) WithValue(value string) *Input {
	i.value = value
	i.textInput.SetValue(value)
	return i
}

func (i *Input) WithHelp(help string) *Input {
	i.help = help
	return i
}

func (i *Input) WithValidation(fn ValidationFunc) *Input {
	i.validation = fn
	return i
}

func (i *Input) WithCharLimit(limit int) *Input {
	i.textInput.CharLimit = limit
	return i
}

func (i *Input) Required() *Input {
	i.required = true
	return i
}

func (i *Input) Init() tea.Cmd {
	return textinput.Blink
}

func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if i.focused {
		i.textInput, cmd = i.textInput.Update(msg)
		i.value = i.textInput.Value()

		// Validar em tempo real se tiver erro
		if i.error != "" {
			i.Validate()
		}
	}

	return i, cmd
}

func (i *Input) View() string {
	var b strings.Builder

	// Container style
	containerStyle := lipgloss.NewStyle().Width(i.width)

	// Input style
	inputStyle := lipgloss.NewStyle().
		Width(i.width).
		Border(lipgloss.NormalBorder()).
		BorderForeground(i.theme.Border).
		Padding(0, 1)

	if i.focused {
		inputStyle = inputStyle.BorderForeground(i.theme.Primary)
	}

	if i.error != "" {
		inputStyle = inputStyle.BorderForeground(i.theme.Error)
	}

	// Renderizar input
	i.textInput.Width = i.width - 4 // Compensar padding e border
	inputView := inputStyle.Render(i.textInput.View())
	b.WriteString(inputView)

	// Help text
	if i.help != "" && i.focused && i.error == "" {
		helpStyle := lipgloss.NewStyle().
			Foreground(i.theme.Muted).
			Italic(true).
			MarginTop(1)
		b.WriteString("\n" + helpStyle.Render(i.help))
	}

	// Error message
	if i.error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(i.theme.Error).
			MarginTop(1)
		b.WriteString("\n" + errorStyle.Render("⚠ "+i.error))
	}

	return containerStyle.Render(b.String())
}

func (i *Input) Focus() tea.Cmd {
	i.BaseField.Focus()
	i.textInput.Focus()
	return textinput.Blink
}

func (i *Input) Blur() {
	i.BaseField.Blur()
	i.textInput.Blur()
}

func (i *Input) SetTheme(theme *styles.Theme) {
	i.BaseField.SetTheme(theme)

	i.textInput.PromptStyle = lipgloss.NewStyle().Foreground(theme.Primary)
	i.textInput.TextStyle = lipgloss.NewStyle().Foreground(theme.Foreground)
	i.textInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(theme.Muted)
	i.textInput.Cursor.Style = lipgloss.NewStyle().Foreground(theme.Primary)
}

func (i *Input) Reset() {
	i.BaseField.Reset()
	i.textInput.SetValue("")
}
