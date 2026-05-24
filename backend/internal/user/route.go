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

		userGroup.GET("/meauth", h.GetMe)
		userGroup.PATCH("/meauth", h.UpdateMe)
		userGroup.DELETE("/meauth", h.DeleteMe)
	}
}
