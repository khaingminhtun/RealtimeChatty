package relationships

import "github.com/gin-gonic/gin"

func RelationshipRoutes(rg *gin.RouterGroup, h *RelationshipHandler, authMiddleware gin.HandlerFunc) {
	relationshipGroup := rg.Group("/relationship")

	{
		relationshipGroup.Use(authMiddleware)
		relationshipGroup.POST("/", h.CreateRelationship)
		relationshipGroup.GET("/:id", h.GetRelationshipByID)   // GET /api/v1/relationships/:id
		relationshipGroup.GET("/", h.ListRelationships)        // GET /api/v1/relationships?type=friend
		relationshipGroup.PATCH("/:id", h.UpdateRelationship)  // PATCH /api/v1/relationships/:id
		relationshipGroup.DELETE("/:id", h.DeleteRelationship) // DELETE /api/v1/relationships/:id

		relationshipGroup.PUT("/:id/tags", h.ReplaceTags)       // PUT /api/v1/relationships/:id/tags
		relationshipGroup.POST("/:id/tags", h.AppendTags)       // POST /api/v1/relationships/:id/tags
		relationshipGroup.DELETE("/:id/tags/:tag", h.RemoveTag) // DELETE /api/v1/relationships/:id/tags/:tag
	}
}
