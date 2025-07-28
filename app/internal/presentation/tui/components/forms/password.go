package forms

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Password struct {
	Input
	showPassword bool
	strength     PasswordStrength
}

type PasswordStrength int

const (
	PasswordWeak PasswordStrength = iota
	PasswordFair
	PasswordGood
	PasswordStrong
)

func NewPassword() *Password {
	p := &Password{
		Input: *NewInput(),
	}

	// Configurar como senha
	p.textInput.EchoMode = textinput.EchoPassword
	p.textInput.EchoCharacter = '•'

	return p
}

func (p *Password) WithShowToggle() *Password {
	p.showPassword = false
	return p
}

func (p *Password) WithStrengthIndicator() *Password {
	p.strength = PasswordWeak
	return p
}

func (p *Password) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if p.focused {
		// Toggle show/hide com Ctrl+S
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+s" {
			p.toggleShow()
		}
	}

	newInput, cmd := p.Input.Update(msg)
	p.Input = *newInput.(*Input)

	// Calcular força da senha
	if p.value != "" {
		p.calculateStrength()
	}

	return p, cmd
}

func (p *Password) View() string {
	var b strings.Builder

	// View base do input
	b.WriteString(p.Input.View())

	// Indicador de força
	if p.strength >= 0 && p.value != "" && p.focused {
		b.WriteString("\n" + p.renderStrength())
	}

	// Ajuda para toggle
	if p.showPassword && p.focused {
		helpStyle := lipgloss.NewStyle().
			Foreground(p.theme.Muted).
			Italic(true)
		b.WriteString("\n" + helpStyle.Render("Ctrl+S para ocultar senha"))
	}

	return b.String()
}

func (p *Password) toggleShow() {
	p.showPassword = !p.showPassword
	if p.showPassword {
		p.textInput.EchoMode = textinput.EchoNormal
	} else {
		p.textInput.EchoMode = textinput.EchoPassword
	}
}

func (p *Password) calculateStrength() {
	length := len(p.value)
	score := 0

	// Comprimento
	switch {
	case length >= 12:
		score += 3
	case length >= 8:
		score += 2
	case length >= 6:
		score += 1
	}

	// Complexidade
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range p.value {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score += 2
	}

	// Determinar força
	switch {
	case score >= 8:
		p.strength = PasswordStrong
	case score >= 6:
		p.strength = PasswordGood
	case score >= 4:
		p.strength = PasswordFair
	default:
		p.strength = PasswordWeak
	}
}

func (p *Password) renderStrength() string {
	var label, bar string
	var color lipgloss.Color

	barWidth := 20

	switch p.strength {
	case PasswordWeak:
		label = "Fraca"
		color = p.theme.Error
		bar = strings.Repeat("█", barWidth/4) + strings.Repeat("░", 3*barWidth/4)

	case PasswordFair:
		label = "Razoável"
		color = p.theme.Warning
		bar = strings.Repeat("█", barWidth/2) + strings.Repeat("░", barWidth/2)

	case PasswordGood:
		label = "Boa"
		color = lipgloss.Color("#FFA500") // Orange
		bar = strings.Repeat("█", 3*barWidth/4) + strings.Repeat("░", barWidth/4)

	case PasswordStrong:
		label = "Forte"
		color = p.theme.Success
		bar = strings.Repeat("█", barWidth)
	}

	style := lipgloss.NewStyle().Foreground(color)

	return style.Render(fmt.Sprintf("Força: %s %s", bar, label))
}
