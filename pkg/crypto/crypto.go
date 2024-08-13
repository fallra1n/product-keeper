package crypto

import "golang.org/x/crypto/bcrypt"

type Crypto struct{}

func NewCrypto() *Crypto {
	return &Crypto{}
}

// HashPassword get hashed password
func (c Crypto) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareHashAndPassword compare passwords
func (c Crypto) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
