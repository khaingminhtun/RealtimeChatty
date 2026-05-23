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
func (h *AuthHandler) VerifyOTPHandler(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and verification OTP token are required"})
		return
	}

	// Extract device fingerprints using Gin context abstractions
	metadata := extractClientMetadata(c)

	response, err := h.authService.VerifyOTP(c.Request.Context(), req.Email, req.OTP, metadata)
	if err != nil {
		// 401 Unauthorized for expired or wrong tokens
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		response.RefreshTokenExpiry,
		"/",
		"",
		true,
		true,
	)

	c.SetCookie(
		"session_token",
		response.SessionToken,
		response.SessionTokenExpiry,
		"/",
		"",
		true,
		true,
	)

	// 3. Construct the clean JSON response payload (Omitting raw long-lived tokens)
	jsonPayload := VerifyOTPUserResponse{
		UserID:      response.UserID,
		Username:    response.Username,
		Email:       response.Email,
		DisplayName: response.DisplayName,
		AccessToken: response.AccessToken,
	}

	// 4. Return operational access token and user profile via JSON
	c.JSON(http.StatusOK, jsonPayload)
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


