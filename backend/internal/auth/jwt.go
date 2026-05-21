package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims holds the core user data embedded within the Access Token payload
type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenManager orchestrates secure token operations
type TokenManager struct {
	secretKey []byte
}

// NewTokenManager initializes the JWT manager utilizing the runtime environment variables
func NewTokenManager() *TokenManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-to-a-secure-32-byte-random-string-in-production"
	}
	return &TokenManager{
		secretKey: []byte(secret),
	}
}

// GenerateAccessToken signs a short-lived JWT containing identity state
func (m *TokenManager) GenerateAccessToken(userID int64, duration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ValidateAccessToken parses and validates an incoming JWT token string
func (m *TokenManager) ValidateAccessToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signing algorithm")
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token payload claims")
}

// GenerateRefreshToken creates an opaque random token string for session extensions.
// It returns the raw token (to send to the client) and the securely hashed token (for DB storage).
func (m *TokenManager) GenerateRefreshToken() (rawToken string, hashedToken string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	rawToken = hex.EncodeToString(bytes)

	hash := sha256.Sum256([]byte(rawToken))
	hashedToken = hex.EncodeToString(hash[:])
	return rawToken, hashedToken, nil
}
