package relationships

import "github.com/gin-gonic/gin"

func RelationshipRoutes(rg *gin.RouterGroup, h *RelationshipHandler, authMiddleware gin.HandlerFunc) {
	relationshipGroup := rg.Group("/relationship")

	{
		relationshipGroup.Use(authMiddleware)
		relationshipGroup.POST("/", h.CreateRelationship)
	}
}
