package middlewares

import authmiddleware "github.com/aslon1213/comnet_task/internal/app/middleware/AuthMiddleware"

type Middlewares struct {
	auth *authmiddleware.AuthMiddleware
}

func New() *Middlewares {
	return &Middlewares{}
}
