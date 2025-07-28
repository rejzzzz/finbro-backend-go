// internal/api/routes.go
package api

import (
	"finbro-backend-go/internal/api/handlers"
	"finbro-backend-go/internal/api/middleware"
	"finbro-backend-go/internal/config"
	"finbro-backend-go/internal/db"

	"github.com/gin-gonic/gin"
)

func SetupRouter(database *db.DB, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(database, cfg)
	userHandler := handlers.NewUserHandler(database)
	accountHandler := handlers.NewAccountHandler(database)
	transactionHandler := handlers.NewTransactionHandler(database)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no middleware)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			// auth.GET("/google", authHandler.GoogleLogin)
			// auth.GET("/google/callback", authHandler.GoogleCallback)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.DELETE("/account", userHandler.DeleteAccount)
			}

			// Account routes
			accounts := protected.Group("/accounts")
			{
				accounts.GET("/", accountHandler.GetAccounts)
				accounts.POST("/", accountHandler.CreateAccount)
				accounts.GET("/:id", accountHandler.GetAccount)
				accounts.PUT("/:id", accountHandler.UpdateAccount)
				accounts.DELETE("/:id", accountHandler.DeleteAccount)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.GET("/", transactionHandler.GetTransactions)
				transactions.POST("/", transactionHandler.CreateTransaction)
				transactions.GET("/:id", transactionHandler.GetTransaction)
				transactions.PUT("/:id", transactionHandler.UpdateTransaction)
				transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			}
		}
	}

	return router
}
