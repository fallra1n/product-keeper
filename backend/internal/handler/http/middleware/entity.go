package middleware

const (
	// AuthHeader ...
	AuthHeader = "Authorization"
	// UserContext ...
	UserContext = "username"
)

// DefaultResponse ...
type DefaultResponse struct {
	Message string `json:"message"`
}
