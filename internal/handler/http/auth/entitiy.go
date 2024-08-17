package authhttphandler

// DefaultResponse ...
type DefaultResponse struct {
	Message string `json:"message"`
}

// AuthRequest ...
type AuthRequest struct {
	Name     string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse ...
type LoginResponse struct {
	Token string `json:"token"`
}
