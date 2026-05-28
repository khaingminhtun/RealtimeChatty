package user

import (
	"github.com/gin-gonic/gin"
)

// UserRoutes registers protected user profile management endpoints
func UserRoutes(rg *gin.RouterGroup, h *UserHandler, authMiddleware gin.HandlerFunc) {
	userGroup := rg.Group("/users")
	{
		// Inject the JWT protection middleware explicitly across this group
		userGroup.Use(authMiddleware)

		userGroup.GET("/me", h.GetMe)
		userGroup.PATCH("/me", h.UpdateMe)
		userGroup.DELETE("/me", h.DeleteMe)
	}
}
