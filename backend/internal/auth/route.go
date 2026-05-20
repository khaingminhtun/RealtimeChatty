package auth

import "github.com/gin-gonic/gin"

func AuthRoutes(rg *gin.RouterGroup, h *AuthHandler) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
	}
}
