package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterHandler handles POST /api/v1/auth/register
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	// Automatically parses JSON and validates required fields
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	// Pass Gin's native context lifecycle down to the service
	err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// 202 Status indicates account registration is initiated successfully but pending validation
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Account created successfully. Please verify your email via the OTP sent.",
	})
}

// VerifyOTPHandler handles POST /api/v1/auth/verify-otp
// VerifyOTPHandler handles POST /api/v1/auth/verify-otp
func (h *AuthHandler) VerifyOTPHandler(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and verification OTP token are required"})
		return
	}

	// No metadata needed here anymore since we aren't creating a session yet!
	err := h.authService.VerifyOTP(c.Request.Context(), req.Email, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Tell the frontend verification succeeded so it can redirect to login
	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully! You can now log in to your account.",
	})
}

func (h *AuthHandler) ResendOTPHandler(c *gin.Context) {
	var req ResendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid target email address is required"})
		return
	}

	err := h.authService.ResendOTP(c.Request.Context(), req.Email)
	if err != nil {
		// Return a generic bad request or unauthorized status depending on internal error match
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 200 OK because the action is fully completed successfully
	c.JSON(http.StatusOK, gin.H{
		"message": "A fresh verification code has been dispatched to your email address.",
	})
}

// LoginHandler handles POST /api/v1/auth/login
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid email address and password are required"})
		return
	}

	metadata := extractClientMetadata(c)

	// Call service layer (still returns the full token info internally)
	authData, err := h.authService.Login(c.Request.Context(), &req, metadata)
	if err != nil {
		if err.Error() == "please verify your email address before logging in" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",             // Cookie Name
		authData.RefreshToken,       // Token Value
		authData.RefreshTokenExpiry, // Expiry MaxAge in seconds
		"/",                         // Path scope
		"",                          // Domain (empty defaults to current host)
		false,                       // Secure flag (Set to TRUE in production for HTTPS!)
		true,                        // HttpOnly flag (CRITICAL: Prevents JS reading the token)
	)

	c.SetCookie(
		"session_token",             // Cookie Name
		authData.SessionToken,       // Token Value
		authData.SessionTokenExpiry, // Expiry MaxAge in seconds
		"/",                         // Path scope
		"",                          // Domain
		false,                       // Secure flag (Set to TRUE in production for HTTPS!)
		true,                        // HttpOnly flag (CRITICAL)
	)

	c.JSON(http.StatusOK, gin.H{
		"user_id":      authData.UserID,
		"username":     authData.Username,
		"email":        authData.Email,
		"display_name": authData.DisplayName,
		"access_token": authData.AccessToken, // Sent directly in JSON
	})
}

// ForgotPasswordHandler handles POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPasswordHandler(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid account email is required"})
		return
	}

	_ = h.authService.ForgotPassword(c.Request.Context(), req.Email)

	// Always return 200 OK with the exact same message to protect user privacy!
	c.JSON(http.StatusOK, gin.H{
		"message": "If that email matches an account, a secure recovery code has been dispatched.",
	})
}

// ResetPasswordHandler handles POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPasswordHandler(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully. You can now log in using your new credentials.",
	})
}

// Private helper to extract user device fingerprints using Gin abstractions
func extractClientMetadata(c *gin.Context) ClientMetadata {
	// Gin's ClientIP automatically handles tracking proxy chains (X-Forwarded-For, X-Real-IP)
	ipAddress := c.ClientIP()
	if strings.Contains(ipAddress, ":") {
		ipAddress = strings.Split(ipAddress, ":")[0]
	}

	return ClientMetadata{
		DeviceName: c.GetHeader("X-Device-Name"), // Custom app client identifying header
		DeviceType: c.GetHeader("X-Device-Type"), // e.g., iOS, Android, Desktop-Web
		IPAddress:  ipAddress,
		UserAgent:  c.Request.UserAgent(),
	}
}
