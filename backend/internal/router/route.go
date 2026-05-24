package router

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/internal/auth"
	"github.com/khaingminhtun/realtimechatty/internal/user"
)

type Handlers struct {
	AuthHandler    *auth.AuthHandler
	UserHandler    *user.UserHandler
	AuthMiddleware gin.HandlerFunc
}

func SetupRouter(h Handlers) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")

	auth.AuthRoutes(
		api,
		h.AuthHandler,
	)

	user.UserRoutes(api, h.UserHandler, h.AuthMiddleware)

	return r
}
