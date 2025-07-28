package messages

import "github.com/PHRaulino/phengineer/internal/presentation/tui/models"

// Messages para navegação entre telas
type ChangeScreenMsg struct {
	Screen models.Screen
}

type PopScreenMsg struct{}