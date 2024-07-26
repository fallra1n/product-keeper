package middleware

const (
	AuthHeader  = "Authorization"
	UserContext = "username"
)

type DefaultResponse struct {
	Message string `json:"message"`
}