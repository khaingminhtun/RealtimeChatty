package relationships

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/contextutils"
)

type RelationshipHandler struct {
	svc RelationshipService
}

func NewHandler(svc RelationshipService) *RelationshipHandler {
	return &RelationshipHandler{svc: svc}
}

func (h *RelationshipHandler) CreateRelationship(c *gin.Context) {
	// 1. Get the authenticated user ID from context (set by your Auth middleware)
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Bind and validate the incoming JSON body
	var req CreateRelationshipDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
		return
	}

	// Manual field validation (Alternatively, use struct tags like `binding:"required"`)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	req.OwnerID = userID

	// 3. Call service and receive the clean DTO struct
	responseDTO, err := h.svc.AddRelationship(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create relationship"})
		return
	}

	// 4. Return the successful response
	c.JSON(http.StatusCreated, responseDTO)
}
