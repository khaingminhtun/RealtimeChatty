package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(bytes), err
}

func CheckPassword(
	password string,
	hash string,
) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

func HashString(val string) string {
	hash := sha256.Sum256([]byte(val))
	return hex.EncodeToString(hash[:])
}
