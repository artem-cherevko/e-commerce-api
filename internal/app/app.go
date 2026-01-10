package app

import (
	"e-commerce-api/internal/config"
	"e-commerce-api/internal/database"
	"e-commerce-api/internal/modules/auth"
	"e-commerce-api/internal/modules/cart"
	"e-commerce-api/internal/modules/payments"
	"e-commerce-api/internal/modules/products"

	"github.com/gin-gonic/gin"
	"github.com/matthewhartstonge/argon2"
	"github.com/stripe/stripe-go/v84"
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
	stripe.Key = cfg.STRIPE_KEY

	// Init engine
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.SetTrustedProxies([]string{"192.168.0.0/24"})

	// AUTH
	authService := auth.NewAuthService(db, mv, &argon)
	authHandler := auth.NewAuthHandler(authService)

	// PRODUCTS
	productsService := products.NewProductsService(db)
	productsHandler := products.NewProductsHandler(productsService)

	// CART
	cartService := cart.NewCartService(db, productsService)
	cartHandler := cart.NewCartHandler(cartService)

	// PAYMENTS
	paymentsService := payments.NewPaymentsService(db)
	paymentsHandler := payments.NewPaymentsHandler(paymentsService, cfg.STRIPE_WEBHOOK_SECRET)

	r := engine.Group("/api/v1")

	auth.RegisterAuthRouters(r, authHandler)
	products.RegisterProductsRouters(r, productsHandler, mv)
	cart.RegisterCartRouters(r, cartHandler, mv)
	payments.NewPaymentsRouter(r, paymentsHandler, mv)

	return &App{
		Engine: engine,
	}, nil
}
