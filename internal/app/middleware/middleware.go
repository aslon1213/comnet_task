package middlewares

import authmiddleware "github.com/aslon1213/comnet_task/internal/app/middleware/AuthMiddleware"

type Middlewares struct {
	Auth *authmiddleware.AuthMiddleware
}

func New(Auth *authmiddleware.AuthMiddleware) *Middlewares {
	return &Middlewares{
		Auth: Auth,
	}
}
