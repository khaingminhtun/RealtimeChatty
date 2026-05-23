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

// TokenPairResponse collects the generated credentials to pass down to handlers
type TokenPairResponse struct {
	AccessToken       string    
	AccessTokenExpiry time.Time 

	RawRefreshToken    string     // Never expose in JSON response body
	HashedRefreshToken string    
	RefreshTokenExpiry time.Time 

	RawSessionToken    string  
	HashedSessionToken string    
	SessionTokenExpiry time.Time 
}

// TokenManager orchestrates secure token operations
type TokenManager struct {
	secretKey           []byte
	accessTokenDuration time.Duration
	refreshTokenDuration time.Duration
	sessionTokenDuration time.Duration
}

// NewTokenManager initializes the JWT manager utilizing runtime environment variables and explicit durations
func NewTokenManager() *TokenManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-to-a-secure-32-byte-random-string-in-production"
	}

	// Parse custom durations from environment, falling back to industry defaults if missing/malformed
	accessDur := parseDurationEnv("ACCESS_TOKEN_EXPIRY", 15*time.Minute)
	refreshDur := parseDurationEnv("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour)
	sessionDur := parseDurationEnv("SESSION_TOKEN_EXPIRY", 30*24*time.Hour)

	return &TokenManager{
		secretKey:            []byte(secret),
		accessTokenDuration:  accessDur,
		refreshTokenDuration: refreshDur,
		sessionTokenDuration: sessionDur,
	}
}

// GenerateAuthTokens creates an atomic set of fresh credentials (Access, Refresh, and Session) for an identity
func (m *TokenManager) GenerateAuthTokens(userID int64) (*TokenPairResponse, error) {
	now := time.Now()

	// 1. Generate Access Token (JWT)
	accessExpiry := now.Add(m.accessTokenDuration)
	accessToken, err := m.generateAccessToken(userID, accessExpiry)
	if err != nil {
		return nil, err
	}

	// 2. Generate Refresh Token (Opaque)
	rawRefresh, hashedRefresh, err := m.generateOpaqueToken()
	if err != nil {
		return nil, err
	}

	// 3. Generate Session Token (Opaque)
	rawSession, hashedSession, err := m.generateOpaqueToken()
	if err != nil {
		return nil, err
	}

	return &TokenPairResponse{
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessExpiry,
		RawRefreshToken:    rawRefresh,
		HashedRefreshToken: hashedRefresh,
		RefreshTokenExpiry: now.Add(m.refreshTokenDuration),
		RawSessionToken:    rawSession,
		HashedSessionToken: hashedSession,
		SessionTokenExpiry: now.Add(m.sessionTokenDuration),
	}, nil
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

// HashOpaqueToken provides a public helper to hash raw cookies on incoming requests (e.g. for lookups)
func (m *TokenManager) HashOpaqueToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}

// Internal token generator builders
func (m *TokenManager) generateAccessToken(userID int64, expiresAt time.Time) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *TokenManager) generateOpaqueToken() (raw string, hashed string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(bytes)
	hashed = m.HashOpaqueToken(raw)
	return raw, hashed, nil
}

// Helper function to safely parse durations from system envs
func parseDurationEnv(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback // Keeps service running safely using the standard fallback if configuration is mistyped
	}
	return d
}