// internal/api/routes.go
package api

import (
	"finbro-backend-go/internal/api/handlers"
	"finbro-backend-go/internal/api/middleware"
	"finbro-backend-go/internal/config"
	"finbro-backend-go/internal/db"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	database *db.DB,
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	accountHandler *handlers.AccountHandler,
	transactionHandler *handlers.TransactionHandler,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Finbro API"})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	jwtSecret := cfg.JWT.Secret

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no middleware)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken) // Ensure RefreshToken uses user_id from context
			// --- Add Google OAuth routes ---
			authGroup.GET("/google", authHandler.GoogleLogin)
			authGroup.GET("/google/callback", authHandler.GoogleCallback)
		}

		// Protected routes
		protected := v1.Group("/")
		// --- FIXED: Pass the correct JWT secret ---
		protected.Use(middleware.AuthRequired(jwtSecret)) // Use the extracted secret
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
