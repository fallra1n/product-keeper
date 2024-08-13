package shared

// Crypto interface for working with cryptography
type Crypto interface {
	HashPassword(password string) (string, error)
}
