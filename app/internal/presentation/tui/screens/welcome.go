package screens

import (
	"strings"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/messages"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/models"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeScreen struct {
	models.BaseModel
	options []string
	cursor  int
}

func NewWelcomeScreen() *WelcomeScreen {
	return &WelcomeScreen{
		BaseModel: models.BaseModel{
			Theme: styles.DefaultTheme,
		},
		options: []string{
			"🔐 Configurar Autenticação",
			"⚙️  Configurações",
			"📋 Ver Status",
			"❌ Sair",
		},
		cursor: 0,
	}
}

func (s *WelcomeScreen) Init() tea.Cmd {
	return nil
}

func (s *WelcomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit

		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}

		case "down", "j":
			if s.cursor < len(s.options)-1 {
				s.cursor++
			}

		case "enter", " ":
			switch s.cursor {
			case 0:
				return s, func() tea.Msg {
					return messages.ChangeScreenMsg{Screen: NewAuthSetupScreen()}
				}
			case 1:
				// TODO: Implementar tela de configurações
			case 2:
				// TODO: Implementar tela de status
			case 3:
				return s, tea.Quit
			}
		}
	}

	return s, nil
}

func (s *WelcomeScreen) View() string {
	var doc strings.Builder

	// Logo/Title
	titleStyle := s.Theme.GetStyles().Title.
		Width(s.Width).
		Align(lipgloss.Center).
		MarginBottom(2)

	doc.WriteString(titleStyle.Render("🔧 PHRaulino Engineer CLI"))
	doc.WriteString("\n")

	// Subtitle
	subtitleStyle := s.Theme.GetStyles().Subtitle.
		Width(s.Width).
		Align(lipgloss.Center).
		MarginBottom(3)

	doc.WriteString(subtitleStyle.Render("Gerenciador de configurações e autenticação"))
	doc.WriteString("\n")

	// Menu options
	for i, option := range s.options {
		cursor := " "
		if s.cursor == i {
			cursor = ">"
		}

		optionStyle := lipgloss.NewStyle().
			Foreground(s.Theme.Foreground).
			MarginLeft(2)

		if s.cursor == i {
			optionStyle = optionStyle.
				Background(s.Theme.Primary).
				Foreground(s.Theme.Background).
				Bold(true).
				Padding(0, 1)
		}

		line := cursor + " " + option
		if s.cursor == i {
			line = optionStyle.Render(option)
		} else {
			line = optionStyle.Render(line)
		}

		doc.WriteString(line)
		doc.WriteString("\n")
	}

	// Help text
	doc.WriteString("\n")
	helpStyle := s.Theme.GetStyles().Info.
		Width(s.Width).
		Align(lipgloss.Center).
		MarginTop(2)

	doc.WriteString(helpStyle.Render("Use ↑/↓ para navegar • Enter para selecionar • q/Ctrl+C para sair"))

	// Container
	containerStyle := lipgloss.NewStyle().
		Width(s.Width).
		Height(s.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Theme.Border).
		Padding(2, 4)

	return containerStyle.Render(doc.String())
}

func (s *WelcomeScreen) SetSize(width, height int) {
	s.BaseModel.SetSize(width, height)
}

func (s *WelcomeScreen) SetTheme(theme *styles.Theme) {
	s.BaseModel.Theme = theme
}

func (s *WelcomeScreen) GetTitle() string {
	return "Welcome"
}

func (s *WelcomeScreen) HandleError(err error) tea.Cmd {
	s.BaseModel.Error = err
	return nil
}