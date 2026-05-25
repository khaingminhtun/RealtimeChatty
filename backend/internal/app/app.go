package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/khaingminhtun/realtimechatty/internal/auth"
	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/mail"
	middleware "github.com/khaingminhtun/realtimechatty/internal/middlewares"
	"github.com/khaingminhtun/realtimechatty/internal/relationships"
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
	relationshiprepo := relationships.NewRelationShipRepository(queries)

	tokenManager := auth.NewTokenManager()
	mailer := mail.NewSendGridMailer()

	authGuard := middleware.AuthMiddleware(tokenManager, "access_token")

	authSvc := auth.NewAuthService(userRepo, authRepo, *tokenManager, mailer, redis)
	userSvc := user.NewUserService(userRepo)
	relationshipservice := relationships.NewRelationshipService(relationshiprepo)

	authHandler := auth.NewAuthHandler(authSvc)
	userHandler := user.NewUserHandler(userSvc)
	relationshipHandler := relationships.NewHandler(relationshipservice)

	// IMPORTANT: must be used immediately
	r := router.SetupRouter(
		router.Handlers{
			AuthHandler:         authHandler,
			UserHandler:         userHandler,
			RelationshipHandler: relationshipHandler,
			AuthMiddleware:      authGuard,
		},
	)

	return &App{
		router: r,
	}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
