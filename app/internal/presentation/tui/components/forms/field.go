package forms

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
)

// Field é a interface que todos os campos devem implementar
type Field interface {
	tea.Model

	// Getters
	Value() string
	IsValid() bool
	Error() string
	IsFocused() bool
	IsRequired() bool

	// Setters
	SetValue(string)
	Focus() tea.Cmd
	Blur()
	SetTheme(*styles.Theme)
	SetWidth(int)

	// Validation
	Validate() error
	SetValidation(ValidationFunc)
	SetRequired(bool)

	// Helpers
	Reset()
	GetLabel() string
	SetLabel(string)
}

// BaseField implementa funcionalidades comuns
type BaseField struct {
	label       string
	value       string
	placeholder string
	help        string
	required    bool
	focused     bool
	error       string
	width       int
	theme       *styles.Theme
	validation  ValidationFunc
}

func (f *BaseField) GetLabel() string {
	return f.label
}

func (f *BaseField) SetLabel(label string) {
	f.label = label
}

func (f *BaseField) Value() string {
	return f.value
}

func (f *BaseField) SetValue(value string) {
	f.value = value
	// Limpar erro ao mudar valor
	f.error = ""
}

func (f *BaseField) IsValid() bool {
	return f.error == ""
}

func (f *BaseField) Error() string {
	return f.error
}

func (f *BaseField) IsFocused() bool {
	return f.focused
}

func (f *BaseField) IsRequired() bool {
	return f.required
}

func (f *BaseField) SetRequired(required bool) {
	f.required = required
}

func (f *BaseField) Focus() tea.Cmd {
	f.focused = true
	return nil
}

func (f *BaseField) Blur() {
	f.focused = false
	f.Validate()
}

func (f *BaseField) SetTheme(theme *styles.Theme) {
	f.theme = theme
}

func (f *BaseField) SetWidth(width int) {
	f.width = width
}

func (f *BaseField) SetValidation(fn ValidationFunc) {
	f.validation = fn
}

func (f *BaseField) Validate() error {
	// Validação de campo obrigatório
	if f.required && f.value == "" {
		f.error = "Este campo é obrigatório"
		return ErrFieldRequired
	}

	// Validação customizada
	if f.validation != nil {
		if err := f.validation(f.value); err != nil {
			f.error = err.Error()
			return err
		}
	}

	f.error = ""
	return nil
}

func (f *BaseField) Reset() {
	f.value = ""
	f.error = ""
	f.focused = false
}
