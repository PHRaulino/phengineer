package common

import (
	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

type Container struct {
	width   int
	height  int
	padding int
	border  bool
	title   string
	theme   *styles.Theme
}

func NewContainer() *Container {
	return &Container{
		padding: 1,
		border:  true,
		theme:   styles.DefaultTheme,
	}
}

func (c *Container) SetSize(width, height int) *Container {
	c.width = width
	c.height = height
	return c
}

func (c *Container) SetTitle(title string) *Container {
	c.title = title
	return c
}

func (c *Container) SetPadding(padding int) *Container {
	c.padding = padding
	return c
}

func (c *Container) SetBorder(border bool) *Container {
	c.border = border
	return c
}

func (c *Container) Render(content string) string {
	style := lipgloss.NewStyle().
		Width(c.width).
		Height(c.height).
		Padding(c.padding)

	if c.border {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(c.theme.Border)
	}

	if c.title != "" {
		// Title will be rendered separately above the container
		titleStyle := lipgloss.NewStyle().
			Foreground(c.theme.Primary).
			Bold(true).
			Width(c.width).
			Align(lipgloss.Center).
			MarginBottom(1)
		
		titleContent := titleStyle.Render(c.title)
		return titleContent + "\n" + style.Render(content)
	}

	return style.Render(content)
}

// CenteredContainer centraliza o conteúdo
func CenteredContainer(width, height int, content string) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// SplitContainer divide a tela em duas partes
func SplitContainer(width, height int, left, right string, ratio float64) string {
	leftWidth := int(float64(width) * ratio)
	rightWidth := width - leftWidth

	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(height).
		Border(lipgloss.NormalBorder(), false, true, false, false)

	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(height)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)
}
