package contacts

import "github.com/gin-gonic/gin"

// ContactRoutes registers all contact-related HTTP endpoints.
//
// Route map:
//
//	POST   /contacts/relationships/:id/contacts     — log a new contact interaction
//	GET    /contacts/relationships/:id/contacts     — list all contacts for a relationship
//	GET    /contacts/:id                            — get a single contact log
//	PATCH  /contacts/:id                            — partially update a contact log
//	DELETE /contacts/:id                            — delete a contact log
//
//	GET    /contacts/drift                          — list relationships with drift data (per user)
//	GET    /contacts/drift/reminders                — list overdue drift reminders (global)
//	PATCH  /contacts/drift/reminders/:id/sent       — mark a drift reminder as sent
//
//	GET    /contacts/search?q=<query>               — full-text search over relationships
func ContactRoutes(rg *gin.RouterGroup, h *ContactHandler, authMiddleware gin.HandlerFunc) {
	g := rg.Group("/contacts")
	g.Use(authMiddleware)

	// --- Contact log endpoints (scoped to a relationship) ---
	g.POST("/relationships/:id/contacts", h.LogContact)
	g.GET("/relationships/:id/contacts", h.GetContactsByRelationship)

	// --- Contact log single-entry endpoints ---
	// Note: these must be declared BEFORE the drift group to avoid route conflicts.
	g.GET("/:id", h.GetContact)
	g.PATCH("/:id", h.UpdateContact)
	g.DELETE("/:id", h.DeleteContact)

	// --- Drift / reminder endpoints ---
	g.GET("/drift", h.ListRelationshipsForDrift)
	g.GET("/drift/reminders", h.ListDriftReminders)
	g.PATCH("/drift/reminders/:id/sent", h.MarkReminderSent)

	// --- Search ---
	g.GET("/search", h.SearchRelationships)
}
