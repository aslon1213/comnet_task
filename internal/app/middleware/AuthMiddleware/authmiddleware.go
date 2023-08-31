package authmiddleware

type AuthMiddleware struct {
}

func New() *AuthMiddleware {
	return &AuthMiddleware{}
}
