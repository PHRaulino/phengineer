package forms

import (
	"fmt"
	"strings"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextArea struct {
	BaseField
	textArea    textarea.Model
	rows        int
	maxLength   int
	showCounter bool
}

func NewTextArea() *TextArea {
	ta := textarea.New()
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.CharLimit = 1000

	return &TextArea{
		BaseField: BaseField{
			theme: styles.DefaultTheme,
			width: 60,
		},
		textArea:    ta,
		rows:        5,
		maxLength:   1000,
		showCounter: true,
	}
}

func (t *TextArea) WithRows(rows int) *TextArea {
	t.rows = rows
	t.textArea.SetHeight(rows)
	return t
}

func (t *TextArea) WithMaxLength(max int) *TextArea {
	t.maxLength = max
	t.textArea.CharLimit = max
	return t
}

func (t *TextArea) WithPlaceholder(placeholder string) *TextArea {
	t.placeholder = placeholder
	t.textArea.Placeholder = placeholder
	return t
}

func (t *TextArea) ShowLineNumbers(show bool) *TextArea {
	t.textArea.ShowLineNumbers = show
	return t
}

func (t *TextArea) ShowCounter(show bool) *TextArea {
	t.showCounter = show
	return t
}

func (t *TextArea) Init() tea.Cmd {
	return textarea.Blink
}

func (t *TextArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if t.focused {
		t.textArea, cmd = t.textArea.Update(msg)
		t.value = t.textArea.Value()

		// Validar em tempo real se tiver erro
		if t.error != "" {
			t.Validate()
		}
	}

	return t, cmd
}

func (t *TextArea) View() string {
	var b strings.Builder

	// Container
	containerStyle := lipgloss.NewStyle().Width(t.width)

	// TextArea style
	areaStyle := lipgloss.NewStyle().
		Width(t.width).
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.theme.Border)

	if t.focused {
		areaStyle = areaStyle.BorderForeground(t.theme.Primary)
	}

	if t.error != "" {
		areaStyle = areaStyle.BorderForeground(t.theme.Error)
	}

	// Configurar dimensões
	t.textArea.SetWidth(t.width - 2)
	t.textArea.SetHeight(t.rows)

	// Renderizar textarea
	areaView := areaStyle.Render(t.textArea.View())
	b.WriteString(areaView)

	// Counter e help
	if t.focused {
		var footer []string

		// Character counter
		if t.showCounter {
			count := len(t.value)
			counterStyle := lipgloss.NewStyle().
				Foreground(t.theme.Muted).
				Align(lipgloss.Right)

			if count > int(float64(t.maxLength)*0.9) {
				counterStyle = counterStyle.Foreground(t.theme.Warning)
			}

			counter := counterStyle.Render(fmt.Sprintf("%d/%d", count, t.maxLength))
			footer = append(footer, counter)
		}

		// Help text
		if t.help != "" && t.error == "" {
			helpStyle := lipgloss.NewStyle().
				Foreground(t.theme.Muted).
				Italic(true)
			footer = append(footer, helpStyle.Render(t.help))
		}

		if len(footer) > 0 {
			footerView := lipgloss.JoinHorizontal(
				lipgloss.Top,
				strings.Join(footer, " • "),
			)
			b.WriteString("\n" + footerView)
		}
	}

	// Error
	if t.error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(t.theme.Error).
			MarginTop(1)
		b.WriteString("\n" + errorStyle.Render("⚠ "+t.error))
	}

	return containerStyle.Render(b.String())
}

func (t *TextArea) Focus() tea.Cmd {
	t.BaseField.Focus()
	t.textArea.Focus()
	return textarea.Blink
}

func (t *TextArea) Blur() {
	t.BaseField.Blur()
	t.textArea.Blur()
}

func (t *TextArea) SetTheme(theme *styles.Theme) {
	t.BaseField.SetTheme(theme)

	t.textArea.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("237"))

	t.textArea.FocusedStyle.Prompt = lipgloss.NewStyle().
		Foreground(theme.Primary)

	t.textArea.BlurredStyle.Prompt = lipgloss.NewStyle().
		Foreground(theme.Muted)
}

func (t *TextArea) Reset() {
	t.BaseField.Reset()
	t.textArea.SetValue("")
}
