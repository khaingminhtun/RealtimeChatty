package router

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/internal/auth"
	"github.com/khaingminhtun/realtimechatty/internal/contacts"
	"github.com/khaingminhtun/realtimechatty/internal/notes"
	"github.com/khaingminhtun/realtimechatty/internal/relationships"
	"github.com/khaingminhtun/realtimechatty/internal/user"
)

type Handlers struct {
	AuthHandler         *auth.AuthHandler
	UserHandler         *user.UserHandler
	RelationshipHandler *relationships.RelationshipHandler
	NoteHandler         *notes.NoteHandler
	ContactHandler      *contacts.ContactHandler
	AuthMiddleware      gin.HandlerFunc
}

func SetupRouter(h Handlers) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")

	auth.AuthRoutes(
		api,
		h.AuthHandler,
	)

	user.UserRoutes(api,
		h.UserHandler,
		h.AuthMiddleware,
	)

	relationships.RelationshipRoutes(
		api,
		h.RelationshipHandler,
		h.AuthMiddleware,
	)

	notes.NotesRoute(
		api,
		h.NoteHandler,
		h.AuthMiddleware,
	)

	contacts.ContactRoutes(
		api,
		h.ContactHandler,
		h.AuthMiddleware,
	)

	return r
}
