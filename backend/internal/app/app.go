package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/khaingminhtun/realtimechatty/internal/auth"
	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/router"
	"github.com/khaingminhtun/realtimechatty/internal/user"
)

type App struct {
	router *gin.Engine
}

func NewApp(pool *pgxpool.Pool) *App {

	queries := db.New(pool)

	userRepo := user.NewUserRepository(queries)
	authRepo := auth.NewAuthRepository(queries)

	authSvc := auth.NewAuthService(userRepo, authRepo)
	authHandler := auth.NewAuthHandler(authSvc)

	// IMPORTANT: must be used immediately
	r := router.SetupRouter(
		router.Handlers{
			AuthHandler: authHandler,
		},
	)

	return &App{
		router: r,
	}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
