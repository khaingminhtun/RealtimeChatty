package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/contextutils"
)

type UserHandler struct {
	userInterface UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{userInterface: s}
}

// GET /api/v1/users/meauth
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userInterface.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PATCH /api/v1/users/meauth
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed JSON payload"})
		return
	}

	dto := UpdateProfileDTO{
		ID:          userID,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		Timezone:    req.Timezone,
		PushToken:   req.PushToken,
	}

	updatedUser, err := h.userInterface.UpdateProfile(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DELETE /api/v1/users/meauth
func (h *UserHandler) DeleteMe(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.userInterface.DeleteAccount(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	// 204 No Content has no body response payload
	c.Status(http.StatusNoContent)
}
