package tui

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/messages"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/models"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/screens"
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	currentScreen models.Screen
	screenStack   []models.Screen
	width         int
	height        int
	theme         *styles.Theme
}

func NewApp() *App {
	return &App{
		theme:       styles.DefaultTheme,
		screenStack: make([]models.Screen, 0),
	}
}

func (a *App) Start() error {
	// Tela inicial
	a.currentScreen = screens.NewWelcomeScreen()

	p := tea.NewProgram(a, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		a.currentScreen.Init(),
	)
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.currentScreen.SetSize(msg.Width, msg.Height)

	case messages.ChangeScreenMsg:
		a.pushScreen(msg.Screen)
		return a, a.currentScreen.Init()

	case messages.PopScreenMsg:
		if len(a.screenStack) > 0 {
			a.popScreen()
			return a, a.currentScreen.Init()
		}
	}

	// Delegar para a tela atual
	newScreen, cmd := a.currentScreen.Update(msg)
	a.currentScreen = newScreen.(models.Screen)

	return a, cmd
}

func (a *App) View() string {
	return a.currentScreen.View()
}

func (a *App) pushScreen(screen models.Screen) {
	a.screenStack = append(a.screenStack, a.currentScreen)
	a.currentScreen = screen
	a.currentScreen.SetSize(a.width, a.height)
	a.currentScreen.SetTheme(a.theme)
}

func (a *App) popScreen() {
	if len(a.screenStack) > 0 {
		a.currentScreen = a.screenStack[len(a.screenStack)-1]
		a.screenStack = a.screenStack[:len(a.screenStack)-1]
	}
}
