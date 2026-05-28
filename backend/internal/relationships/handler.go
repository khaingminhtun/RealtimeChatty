package relationships

import (
	"net/http"
	"strconv"
	"time"

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

func (h *RelationshipHandler) GetRelationshipByID(c *gin.Context) {
	// Gin extracts named path variables using c.Param("id")
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.svc.GetRelationship(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "relationship profile not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *RelationshipHandler) ListRelationships(c *gin.Context) {
	// Gin extracts query parameters (?type=friend) using c.Query()
	filterType := c.Query("type")

	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	results, err := h.svc.ListRelationships(c.Request.Context(), userID, filterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query relationships"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// PATCH /api/v1/relationships/:id
func (h *RelationshipHandler) UpdateRelationship(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Pointers help identify fields explicitly passed in PATCH payload
	var body struct {
		Name      *string    `json:"name"`
		Type      *string    `json:"type"`
		HowWeMet  *string    `json:"how_we_met"`
		Birthday  *time.Time `json:"birthday"`
		Location  *string    `json:"location"`
		AvatarURL *string    `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload formatting"})
		return
	}

	result, err := h.svc.UpdateRelationship(c.Request.Context(), UpdateRelationshipDTO{
		ID:        id,
		OwnerID:   userID,
		Name:      body.Name,
		Type:      body.Type,
		HowWeMet:  body.HowWeMet,
		Birthday:  body.Birthday,
		Location:  body.Location,
		AvatarURL: body.AvatarURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DELETE /api/v1/relationships/:id
func (h *RelationshipHandler) DeleteRelationship(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.svc.DeleteRelationship(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete target profile"})
		return
	}

	// Your PostgreSQL schema uses ON DELETE CASCADE, meaning all linked rows
	// in private_notes and contacts are safely removed automatically!
	c.JSON(http.StatusOK, gin.H{"message": "relationship and all cascade linked data deleted successfully"})
}
