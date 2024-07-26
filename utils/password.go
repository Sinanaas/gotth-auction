package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(passwornd string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwornd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
} 