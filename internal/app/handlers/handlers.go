package handlers

import (
	"github.com/aslon1213/comnet_task/internal/app/handlers/userHandlers"
)

// userhandlers "github.com/aslon1213/comnet_task/handlers/user"

type Handlers struct {
	UserHandlers *userHandlers.UserHandlers
}

func New(userHandlers *userHandlers.UserHandlers) *Handlers {
	return &Handlers{
		UserHandlers: userHandlers,
	}
}
