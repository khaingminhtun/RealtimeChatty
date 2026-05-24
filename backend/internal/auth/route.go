package auth

import "github.com/gin-gonic/gin"

func AuthRoutes(rg *gin.RouterGroup, h *AuthHandler) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.RegisterHandler)
		authGroup.POST("/verify-otp", h.VerifyOTPHandler)
		authGroup.POST("/resend-otp", h.ResendOTPHandler)
		authGroup.POST("/login", h.LoginHandler)
		authGroup.POST("/forgot-password", h.ForgotPasswordHandler)
		authGroup.POST("/reset-password", h.ResetPasswordHandler)
	}
}
