package notes

import (
	"github.com/gin-gonic/gin"
)

func NotesRoute(rg *gin.RouterGroup, h *NoteHandler, authMiddleware gin.HandlerFunc) {
	notegroup := rg.Group("/note")
	{
		notegroup.Use(authMiddleware)
		notegroup.GET("/:id/notes", h.GetNote)       // GET /api/v1/relationships/:id/notes
		notegroup.PUT("/:id/notes", h.UpsertNote)    // PUT /api/v1/relationships/:id/notes
		notegroup.DELETE("/:id/notes", h.DeleteNote) // DELETE /api/v1/relationships/:id/notes
	}

}
