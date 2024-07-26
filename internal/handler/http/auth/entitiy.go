package authhttphandler

type DefaultResponse struct {
	Message string `json:"message"`
}

type AuthRequest struct {
	Name     string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
