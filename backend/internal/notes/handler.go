package notes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/contextutils"
)

type NoteHandler struct {
	svc NoteService
}

func NewNoteHandler(svc NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

// GET /api/v1/relationships/:id/notes
func (h *NoteHandler) GetNote(c *gin.Context) {
	relID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.svc.GetNote(c.Request.Context(), relID, userID)
	if err != nil {
		// It's normal to not have a note, return empty status or a 404 depending on your frontend preference
		c.JSON(http.StatusNotFound, gin.H{"error": "no private note found for this relationship"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PUT /api/v1/relationships/:id/notes
func (h *NoteHandler) UpsertNote(c *gin.Context) {
	relID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content string is required"})
		return
	}

	result, err := h.svc.UpsertNote(c.Request.Context(), relID, userID, body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save note"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DELETE /api/v1/relationships/:id/notes
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	relID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.svc.DeleteNote(c.Request.Context(), relID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "private note cleared successfully"})
}
