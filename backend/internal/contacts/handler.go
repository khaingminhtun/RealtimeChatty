package contacts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/contextutils"
)

type ContactHandler struct {
	svc ContactService
}

func NewContactHandler(svc ContactService) *ContactHandler {
	return &ContactHandler{svc: svc}
}

// ---------------------------------------------------------------------------
// POST /api/v1/contacts/relationships/:id/contacts
// Log a new contact interaction for a relationship.
// ---------------------------------------------------------------------------

func (h *ContactHandler) LogContact(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	relationshipID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	var body struct {
		Channel     string  `json:"channel"`
		Note        *string `json:"note"`
		ContactedAt *string `json:"contacted_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if body.Channel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel is required"})
		return
	}

	result, err := h.svc.LogContact(c.Request.Context(), LogContactDTO{
		RelationshipID: relationshipID,
		UserID:         userID,
		Channel:        body.Channel,
		Note:           body.Note,
		ContactedAt:    body.ContactedAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log contact"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// ---------------------------------------------------------------------------
// GET /api/v1/contacts/relationships/:id/contacts
// List all contact logs for a relationship (newest first).
// ---------------------------------------------------------------------------

func (h *ContactHandler) GetContactsByRelationship(c *gin.Context) {
	_, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	relationshipID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	results, err := h.svc.GetContactsByRelationship(c.Request.Context(), relationshipID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ---------------------------------------------------------------------------
// GET /api/v1/contacts/:id
// Fetch a single contact log entry (ownership-checked via user_id).
// ---------------------------------------------------------------------------

func (h *ContactHandler) GetContact(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	contactID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact ID"})
		return
	}

	result, err := h.svc.GetContact(c.Request.Context(), contactID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contact log not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ---------------------------------------------------------------------------
// PATCH /api/v1/contacts/:id
// Partially update a contact log entry (ownership-checked via user_id).
// ---------------------------------------------------------------------------

func (h *ContactHandler) UpdateContact(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	contactID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact ID"})
		return
	}

	var body struct {
		Channel     *string `json:"channel"`
		Note        *string `json:"note"`
		ContactedAt *string `json:"contacted_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.svc.UpdateContact(c.Request.Context(), UpdateContactDTO{
		ID:          contactID,
		UserID:      userID,
		Channel:     body.Channel,
		Note:        body.Note,
		ContactedAt: body.ContactedAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update contact"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ---------------------------------------------------------------------------
// DELETE /api/v1/contacts/:id
// Delete a contact log entry (ownership-checked via user_id).
// ---------------------------------------------------------------------------

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	contactID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact ID"})
		return
	}

	if err := h.svc.DeleteContact(c.Request.Context(), contactID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete contact log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "contact log deleted successfully"})
}

// ---------------------------------------------------------------------------
// GET /api/v1/contacts/drift
// List all relationships for the authenticated user with drift scheduling data.
// ---------------------------------------------------------------------------

func (h *ContactHandler) ListRelationshipsForDrift(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	results, err := h.svc.ListRelationshipsForDrift(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch drift data"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ---------------------------------------------------------------------------
// GET /api/v1/contacts/drift/reminders
// List all overdue drift reminders (relationships past their next_contact_at).
// ---------------------------------------------------------------------------

func (h *ContactHandler) ListDriftReminders(c *gin.Context) {
	_, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	results, err := h.svc.ListDriftReminders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch drift reminders"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ---------------------------------------------------------------------------
// PATCH /api/v1/contacts/drift/reminders/:id/sent
// Mark a drift reminder as sent for a specific relationship.
// ---------------------------------------------------------------------------

func (h *ContactHandler) MarkReminderSent(c *gin.Context) {
	_, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	relationshipID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid relationship ID"})
		return
	}

	if err := h.svc.MarkReminderSent(c.Request.Context(), relationshipID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark reminder as sent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reminder marked as sent"})
}

// ---------------------------------------------------------------------------
// GET /api/v1/contacts/search?q=<query>
// Full-text search over the authenticated user's relationships.
// ---------------------------------------------------------------------------

func (h *ContactHandler) SearchRelationships(c *gin.Context) {
	userID, exists := contextutils.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results, err := h.svc.SearchRelationships(c.Request.Context(), userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
