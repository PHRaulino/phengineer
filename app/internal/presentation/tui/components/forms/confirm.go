package forms

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Confirm struct {
	BaseField
	affirmative string
	negative    string
	confirmed   bool
}

func NewConfirm(question string) *Confirm {
	c := &Confirm{
		BaseField: BaseField{
			label: question,
			theme: styles.DefaultTheme,
			width: 40,
		},
		affirmative: "Sim",
		negative:    "Não",
		confirmed:   false,
	}

	c.value = c.negative
	return c
}

func (c *Confirm) WithOptions(yes, no string) *Confirm {
	c.affirmative = yes
	c.negative = no
	c.value = no
	return c
}

func (c *Confirm) WithDefault(confirmed bool) *Confirm {
	c.confirmed = confirmed
	if confirmed {
		c.value = c.affirmative
	} else {
		c.value = c.negative
	}
	return c
}

func (c *Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !c.focused {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "n":
			c.confirmed = false
			c.value = c.negative

		case "right", "l", "y", "s":
			c.confirmed = true
			c.value = c.affirmative

		case "tab", " ", "enter":
			c.confirmed = !c.confirmed
			if c.confirmed {
				c.value = c.affirmative
			} else {
				c.value = c.negative
			}
		}
	}

	return c, nil
}

func (c *Confirm) View() string {
	questionStyle := lipgloss.NewStyle().
		Foreground(c.theme.Foreground).
		Bold(true).
		MarginBottom(1)

	buttonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(c.theme.Border)

	selectedStyle := buttonStyle.Copy().
		Background(c.theme.Primary).
		Foreground(c.theme.Background).
		BorderForeground(c.theme.Primary).
		Bold(true)

	// Renderizar botões
	var noButton, yesButton string

	if !c.confirmed {
		noButton = selectedStyle.Render(c.negative)
		yesButton = buttonStyle.Render(c.affirmative)
	} else {
		noButton = buttonStyle.Render(c.negative)
		yesButton = selectedStyle.Render(c.affirmative)
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		noButton,
		"  ",
		yesButton,
	)

	// Help
	var help string
	if c.focused {
		helpStyle := lipgloss.NewStyle().
			Foreground(c.theme.Muted).
			Italic(true).
			MarginTop(1)
		help = helpStyle.Render("← → ou Y/N para selecionar • Enter para confirmar")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		questionStyle.Render(c.label),
		buttons,
		help,
	)
}

func (c *Confirm) IsConfirmed() bool {
	return c.confirmed
}

func (c *Confirm) Init() tea.Cmd {
	return nil
}

func (c *Confirm) Value() string {
	if c.confirmed {
		return "true"
	}
	return "false"
}
