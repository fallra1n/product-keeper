package shared

// Crypto
type Crypto interface {
	HashPassword(password string) (string, error)
}
