package app

import (
	"e-commerce-api/internal/config"
	"e-commerce-api/internal/database"
	"e-commerce-api/internal/modules/auth"

	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
)

type App struct {
	Engine *gin.Engine
}

func New(cfg *config.Env) (*App, error) {
	db := database.InitDB(cfg.DB_DSN)
	database.Migrate(db)

	mv, err := auth.NewJWTMiddleware(cfg.JWT_SECRET)
	if err != nil {
		return nil, err
	}

	argon := argon2.DefaultConfig()

	// Init engine
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.SetTrustedProxies([]string{"192.168.0.0/24"})

	// AUTH
	authService := auth.NewAuthService(db, mv, &argon)
	authHandler := auth.NewAuthHandler(authService)

	r := engine.Group("/api/v1")

	auth.RegisterAuthRouters(r, authHandler)

	return &App{
		Engine: engine,
	}, nil
}
