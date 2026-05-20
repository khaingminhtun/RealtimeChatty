package router

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/internal/auth"
)

type Handlers struct {
	AuthHandler *auth.AuthHandler
}

func SetupRouter(h Handlers) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")

	auth.AuthRoutes(
		api,
		h.AuthHandler,
	)

	return r
}
