package common

import (
	"fmt"
	"time"

	"github.com/PHRaulino/phengineer/internal/presentation/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NotificationType string

const (
	NotificationSuccess NotificationType = "success"
	NotificationError   NotificationType = "error"
	NotificationWarning NotificationType = "warning"
	NotificationInfo    NotificationType = "info"
)

type Notification struct {
	id        string
	message   string
	type_     NotificationType
	duration  time.Duration
	createdAt time.Time
	theme     *styles.Theme
}

func NewNotification(message string, notifType NotificationType) *Notification {
	return &Notification{
		id:        fmt.Sprintf("%d", time.Now().UnixNano()),
		message:   message,
		type_:     notifType,
		duration:  3 * time.Second,
		createdAt: time.Now(),
		theme:     styles.DefaultTheme,
	}
}

func (n *Notification) SetDuration(d time.Duration) *Notification {
	n.duration = d
	return n
}

func (n *Notification) View() string {
	var icon string
	var style lipgloss.Style

	baseStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		Bold(true)

	switch n.type_ {
	case NotificationSuccess:
		icon = "✅"
		style = baseStyle.
			BorderForeground(n.theme.Success).
			Foreground(n.theme.Success)

	case NotificationError:
		icon = "❌"
		style = baseStyle.
			BorderForeground(n.theme.Error).
			Foreground(n.theme.Error)

	case NotificationWarning:
		icon = "⚠️"
		style = baseStyle.
			BorderForeground(n.theme.Warning).
			Foreground(n.theme.Warning)

	case NotificationInfo:
		icon = "ℹ️"
		style = baseStyle.
			BorderForeground(n.theme.Primary).
			Foreground(n.theme.Primary)
	}

	content := fmt.Sprintf("%s %s", icon, n.message)
	return style.Render(content)
}

// NotificationManager gerencia múltiplas notificações
type NotificationManager struct {
	notifications []*Notification
	maxVisible    int
	position      string // "top-right", "top-left", "bottom-right", "bottom-left"
	theme         *styles.Theme
}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		notifications: make([]*Notification, 0),
		maxVisible:    3,
		position:      "top-right",
		theme:         styles.DefaultTheme,
	}
}

func (nm *NotificationManager) Add(notification *Notification) tea.Cmd {
	nm.notifications = append(nm.notifications, notification)

	// Remover após duração
	return tea.Tick(notification.duration, func(t time.Time) tea.Msg {
		return removeNotificationMsg{id: notification.id}
	})
}

func (nm *NotificationManager) Update(msg tea.Msg) (*NotificationManager, tea.Cmd) {
	switch msg := msg.(type) {
	case removeNotificationMsg:
		nm.removeNotification(msg.id)
	}
	return nm, nil
}

func (nm *NotificationManager) removeNotification(id string) {
	filtered := make([]*Notification, 0)
	for _, n := range nm.notifications {
		if n.id != id {
			filtered = append(filtered, n)
		}
	}
	nm.notifications = filtered
}

func (nm *NotificationManager) View(screenWidth, screenHeight int) string {
	if len(nm.notifications) == 0 {
		return ""
	}

	var views []string
	count := len(nm.notifications)
	if count > nm.maxVisible {
		count = nm.maxVisible
	}

	for i := 0; i < count; i++ {
		views = append(views, nm.notifications[i].View())
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, views...)

	// Posicionar baseado na configuração
	switch nm.position {
	case "top-right":
		return lipgloss.Place(
			screenWidth,
			screenHeight,
			lipgloss.Right,
			lipgloss.Top,
			joined,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
		)
	case "top-left":
		return lipgloss.Place(
			screenWidth,
			screenHeight,
			lipgloss.Left,
			lipgloss.Top,
			joined,
		)
		// ... outros casos
	}

	return joined
}

type removeNotificationMsg struct {
	id string
}
