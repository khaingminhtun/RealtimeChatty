package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/khaingminhtun/realtimechatty/internal/auth"
	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/mail"
	middleware "github.com/khaingminhtun/realtimechatty/internal/middlewares"
	"github.com/khaingminhtun/realtimechatty/internal/router"
	"github.com/khaingminhtun/realtimechatty/internal/user"
)

type App struct {
	router *gin.Engine
}

func NewApp(pool *pgxpool.Pool, redis *redis.Client) *App {

	queries := db.New(pool)

	userRepo := user.NewUserRepository(queries)
	authRepo := auth.NewAuthRepository(pool)

	tokenManager := auth.NewTokenManager()
	mailer := mail.NewSendGridMailer()

	authGuard := middleware.AuthMiddleware(tokenManager, "access_token")

	authSvc := auth.NewAuthService(userRepo, authRepo, *tokenManager, mailer, redis)
	userSvc := user.NewUserService(userRepo)

	authHandler := auth.NewAuthHandler(authSvc)
	userHandler := user.NewUserHandler(userSvc)

	// IMPORTANT: must be used immediately
	r := router.SetupRouter(
		router.Handlers{
			AuthHandler:    authHandler,
			UserHandler:    userHandler,
			AuthMiddleware: authGuard,
		},
	)

	return &App{
		router: r,
	}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
