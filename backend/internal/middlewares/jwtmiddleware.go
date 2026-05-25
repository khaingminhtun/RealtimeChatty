package middleware

import (
	"net/http"
	"strings"

	"github.com/khaingminhtun/realtimechatty/internal/auth" // Replace with your actual auth package path

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepts requests to validate the access token from an HTTP-only cookie
func AuthMiddleware(tm *auth.TokenManager, cookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Attempt to extract token from the specified HTTP-Only Cookie
		tokenString, err := c.Cookie(cookieName)

		// 2. Fallback: Check standard Authorization Header (Bearer Token)
		// This is extremely helpful if you want to test endpoints via Postman or support mobile clients later
		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// If no token was found anywhere, block the request immediately
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token missing or missing cookie"})
			c.Abort()
			return
		}

		// 3. Leverage your auth package TokenManager to validate the JWT
		claims, err := tm.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid, altered, or expired access token"})
			c.Abort()
			return
		}

		// 4. Inject the parsed UserID context securely into Gin's pipeline
		// Your existing handlers use c.Get("userID") to safely inspect this value
		c.Set("userID", claims.UserID)

		// 5. Transfer execution seamlessly to the next handler block
		c.Next()
	}
}


