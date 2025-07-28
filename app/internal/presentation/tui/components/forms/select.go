package forms

import (
	"strings"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Select struct {
	BaseField
	options      []string
	selected     int
	showOptions  bool
	maxVisible   int
	scrollOffset int
}

func NewSelect(options []string) *Select {
	return &Select{
		BaseField: BaseField{
			theme: styles.DefaultTheme,
			width: 40,
		},
		options:    options,
		selected:   0,
		maxVisible: 5,
	}
}

func (s *Select) WithDefault(index int) *Select {
	if index >= 0 && index < len(s.options) {
		s.selected = index
		s.value = s.options[index]
	}
	return s
}

func (s *Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			if s.showOptions {
				s.value = s.options[s.selected]
				s.showOptions = false
				s.Validate()
			} else {
				s.showOptions = true
			}

		case "up", "k":
			if s.showOptions && s.selected > 0 {
				s.selected--
				s.adjustScroll()
			}

		case "down", "j":
			if s.showOptions && s.selected < len(s.options)-1 {
				s.selected++
				s.adjustScroll()
			}

		case "esc":
			s.showOptions = false

		case "tab", "shift+tab":
			s.showOptions = false
		}
	}

	return s, nil
}

func (s *Select) View() string {
	var b strings.Builder

	// Container
	containerStyle := lipgloss.NewStyle().Width(s.width)

	// Select box style
	boxStyle := lipgloss.NewStyle().
		Width(s.width).
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.theme.Border).
		Padding(0, 1)

	if s.focused {
		boxStyle = boxStyle.BorderForeground(s.theme.Primary)
	}

	if s.error != "" {
		boxStyle = boxStyle.BorderForeground(s.theme.Error)
	}

	// Valor selecionado
	displayValue := s.value
	if displayValue == "" && len(s.options) > 0 {
		displayValue = "Selecione uma opção..."
	}

	// Ícone do dropdown
	icon := "▼"
	if s.showOptions {
		icon = "▲"
	}

	valueWidth := s.width - 4 - lipgloss.Width(icon)
	displayStyle := lipgloss.NewStyle().
		Width(valueWidth).
		Foreground(s.theme.Foreground)

	if s.value == "" {
		displayStyle = displayStyle.Foreground(s.theme.Muted)
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		displayStyle.Render(displayValue),
		lipgloss.NewStyle().Foreground(s.theme.Muted).Render(icon),
	)

	b.WriteString(boxStyle.Render(content))

	// Lista de opções
	if s.showOptions && s.focused {
		b.WriteString("\n" + s.renderOptions())
	}

	// Error
	if s.error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(s.theme.Error).
			MarginTop(1)
		b.WriteString("\n" + errorStyle.Render("⚠ "+s.error))
	}

	return containerStyle.Render(b.String())
}

func (s *Select) renderOptions() string {
	optionsStyle := lipgloss.NewStyle().
		Width(s.width).
		Border(lipgloss.NormalBorder()).
		BorderForeground(s.theme.Border).
		BorderTop(false).
		Padding(0, 1).
		MarginTop(-1)

	var visibleOptions []string
	start := s.scrollOffset
	end := start + s.maxVisible

	if end > len(s.options) {
		end = len(s.options)
	}

	for i := start; i < end; i++ {
		option := s.options[i]
		style := lipgloss.NewStyle().Width(s.width - 4)

		if i == s.selected {
			style = style.
				Background(s.theme.Primary).
				Foreground(s.theme.Background).
				Bold(true)
			option = "▸ " + option
		} else {
			option = "  " + option
		}

		visibleOptions = append(visibleOptions, style.Render(option))
	}

	// Indicadores de scroll
	if start > 0 {
		scrollUp := lipgloss.NewStyle().
			Foreground(s.theme.Muted).
			Width(s.width - 4).
			Align(lipgloss.Center).
			Render("▲ mais opções acima")
		visibleOptions = append([]string{scrollUp}, visibleOptions...)
	}

	if end < len(s.options) {
		scrollDown := lipgloss.NewStyle().
			Foreground(s.theme.Muted).
			Width(s.width - 4).
			Align(lipgloss.Center).
			Render("▼ mais opções abaixo")
		visibleOptions = append(visibleOptions, scrollDown)
	}

	return optionsStyle.Render(strings.Join(visibleOptions, "\n"))
}

func (s *Select) adjustScroll() {
	if s.selected < s.scrollOffset {
		s.scrollOffset = s.selected
	} else if s.selected >= s.scrollOffset+s.maxVisible {
		s.scrollOffset = s.selected - s.maxVisible + 1
	}
}

func (s *Select) Value() string {
	if s.selected >= 0 && s.selected < len(s.options) {
		return s.options[s.selected]
	}
	return ""
}

func (s *Select) Focus() tea.Cmd {
	s.BaseField.Focus()
	return nil
}

func (s *Select) Init() tea.Cmd {
	return nil
}

func (s *Select) Blur() {
	s.BaseField.Blur()
	s.showOptions = false
}
