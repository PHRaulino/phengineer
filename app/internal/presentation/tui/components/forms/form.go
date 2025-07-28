package forms

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SubmitMsg é enviada quando o formulário é submetido
type SubmitMsg struct {
	Values map[string]string
}

// SubmitHandler é chamado quando o formulário é válido
type SubmitHandler func(values map[string]string) tea.Cmd

// Form é o builder principal para formulários
type Form struct {
	title       string
	description string
	fields      []Field
	labels      []string
	focused     int
	theme       *styles.Theme
	width       int
	submitLabel string
	cancelLabel string
	onSubmit    SubmitHandler
	onCancel    func() tea.Cmd
	showButtons bool
	layout      string // "vertical" ou "horizontal"
}

// NewForm cria um novo formulário
func NewForm(title, description string) *Form {
	return &Form{
		title:       title,
		description: description,
		fields:      make([]Field, 0),
		labels:      make([]string, 0),
		theme:       styles.DefaultTheme,
		width:       60,
		submitLabel: "Confirmar",
		cancelLabel: "Cancelar",
		showButtons: true,
		layout:      "vertical",
	}
}

// Builder methods
func (f *Form) AddField(label string, field Field) *Form {
	f.fields = append(f.fields, field)
	f.labels = append(f.labels, label)
	field.SetLabel(label)
	field.SetTheme(f.theme)
	field.SetWidth(f.width - 4)
	return f
}

func (f *Form) OnSubmit(handler SubmitHandler) *Form {
	f.onSubmit = handler
	return f
}

func (f *Form) OnCancel(handler func() tea.Cmd) *Form {
	f.onCancel = handler
	return f
}

func (f *Form) SetTheme(theme *styles.Theme) *Form {
	f.theme = theme
	for _, field := range f.fields {
		field.SetTheme(theme)
	}
	return f
}

func (f *Form) SetWidth(width int) *Form {
	f.width = width
	for _, field := range f.fields {
		field.SetWidth(width - 4)
	}
	return f
}

func (f *Form) SetSubmitLabel(label string) *Form {
	f.submitLabel = label
	return f
}

func (f *Form) SetCancelLabel(label string) *Form {
	f.cancelLabel = label
	return f
}

func (f *Form) ShowButtons(show bool) *Form {
	f.showButtons = show
	return f
}

func (f *Form) SetLayout(layout string) *Form {
	f.layout = layout
	return f
}

// Focus management
func (f *Form) focusNext() {
	if f.focused < len(f.fields)-1 {
		f.fields[f.focused].Blur()
		f.focused++
		f.fields[f.focused].Focus()
	} else if f.showButtons && f.focused == len(f.fields)-1 {
		f.fields[f.focused].Blur()
		f.focused++ // Foco nos botões
	}
}

func (f *Form) focusPrevious() {
	if f.focused > 0 {
		if f.focused <= len(f.fields) {
			f.fields[f.focused-1].Blur()
		}
		f.focused--
		if f.focused < len(f.fields) {
			f.fields[f.focused].Focus()
		}
	}
}

// Validation
func (f *Form) validate() bool {
	isValid := true
	for _, field := range f.fields {
		if err := field.Validate(); err != nil {
			isValid = false
		}
	}
	return isValid
}

func (f *Form) getValues() map[string]string {
	values := make(map[string]string)
	for i, field := range f.fields {
		values[f.labels[i]] = field.Value()
	}
	return values
}

// Reset
func (f *Form) Reset() {
	for _, field := range f.fields {
		field.Reset()
	}
	f.focused = 0
	if len(f.fields) > 0 {
		f.fields[0].Focus()
	}
}

// Tea.Model implementation
func (f *Form) Init() tea.Cmd {
	if len(f.fields) > 0 {
		return f.fields[0].Focus()
	}
	return nil
}

func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			f.focusNext()
			return f, nil

		case "shift+tab", "up":
			f.focusPrevious()
			return f, nil

		case "enter":
			// Se estiver nos botões
			if f.showButtons && f.focused == len(f.fields) {
				if f.validate() {
					if f.onSubmit != nil {
						return f, f.onSubmit(f.getValues())
					}
					return f, func() tea.Msg {
						return SubmitMsg{Values: f.getValues()}
					}
				}
			}

		case "esc":
			if f.onCancel != nil {
				return f, f.onCancel()
			}

		case "ctrl+c":
			return f, tea.Quit
		}
	}

	// Atualizar campo focado
	if f.focused < len(f.fields) {
		newField, cmd := f.fields[f.focused].Update(msg)
		f.fields[f.focused] = newField.(Field)
		return f, cmd
	}

	return f, nil
}

func (f *Form) View() string {
	// Container style
	containerStyle := lipgloss.NewStyle().
		Width(f.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(f.theme.Border).
		Padding(1, 2)

	// Title
	titleStyle := f.theme.GetStyles().Title
	titleView := titleStyle.Render(f.title)

	// Description
	var descView string
	if f.description != "" {
		descStyle := f.theme.GetStyles().Subtitle
		descView = descStyle.Render(f.description) + "\n"
	}

	// Fields
	var fieldsView string
	if f.layout == "horizontal" && len(f.fields) == 2 {
		// Layout horizontal para 2 campos
		fieldsView = f.renderHorizontalFields()
	} else {
		// Layout vertical padrão
		fieldsView = f.renderVerticalFields()
	}

	// Buttons
	var buttonsView string
	if f.showButtons {
		buttonsView = "\n" + f.renderButtons()
	}

	// Compose
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleView,
		descView,
		fieldsView,
		buttonsView,
	)

	return containerStyle.Render(content)
}

func (f *Form) renderVerticalFields() string {
	var fields []string

	labelStyle := lipgloss.NewStyle().
		Foreground(f.theme.Foreground).
		Bold(true).
		Width(f.width - 4).
		MarginBottom(1)

	for i, field := range f.fields {
		// Label
		label := labelStyle.Render(f.labels[i] + ":")

		// Field
		fieldView := field.View()

		// Container para o campo
		fieldContainer := lipgloss.NewStyle().
			MarginBottom(2).
			Render(label + "\n" + fieldView)

		fields = append(fields, fieldContainer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, fields...)
}

func (f *Form) renderHorizontalFields() string {
	if len(f.fields) != 2 {
		return f.renderVerticalFields()
	}

	halfWidth := (f.width - 6) / 2

	labelStyle := lipgloss.NewStyle().
		Foreground(f.theme.Foreground).
		Bold(true).
		MarginBottom(1)

	// Ajustar largura dos campos
	f.fields[0].SetWidth(halfWidth)
	f.fields[1].SetWidth(halfWidth)

	// Renderizar campos
	left := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render(f.labels[0]+":"),
		f.fields[0].View(),
	)

	right := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render(f.labels[1]+":"),
		f.fields[1].View(),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		"  ",
		right,
	)
}

func (f *Form) renderButtons() string {
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(f.theme.Border)

	activeButtonStyle := buttonStyle.Copy().
		Background(f.theme.Primary).
		Foreground(f.theme.Background).
		BorderForeground(f.theme.Primary).
		Bold(true)

	var cancelBtn, submitBtn string

	if f.focused == len(f.fields) {
		// Botões focados
		cancelBtn = buttonStyle.Render(f.cancelLabel)
		submitBtn = activeButtonStyle.Render(f.submitLabel)
	} else {
		cancelBtn = buttonStyle.Render(f.cancelLabel)
		submitBtn = buttonStyle.Render(f.submitLabel)
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		cancelBtn,
		"  ",
		submitBtn,
	)

	return lipgloss.Place(
		f.width-4,
		1,
		lipgloss.Right,
		lipgloss.Center,
		buttons,
	)
}
