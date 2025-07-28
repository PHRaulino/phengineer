package models

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/components/common"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// BaseModel contém funcionalidades comuns a todos os models
type BaseModel struct {
	Width      int
	Height     int
	Theme      *styles.Theme
	Header     *common.Header
	Footer     *common.Footer
	Error      error
	Loading    bool
	LoadingMsg string
}

// Screen interface que todos os models devem implementar
type Screen interface {
	tea.Model
	// Métodos adicionais úteis
	SetSize(width, height int)
	SetTheme(theme *styles.Theme)
	GetTitle() string
	HandleError(error) tea.Cmd
}

// Resize atualiza as dimensões
func (m *BaseModel) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	if m.Header != nil {
		m.Header.SetWidth(width)
	}
	if m.Footer != nil {
		m.Footer.SetWidth(width)
	}
}

// Métodos comuns...
