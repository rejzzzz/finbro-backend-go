package main

import (
	"log"

	"finbro-backend-go/internal/api"
	"finbro-backend-go/internal/api/handlers"
	"finbro-backend-go/internal/config"
	"finbro-backend-go/internal/db"
	"finbro-backend-go/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database, err := db.Initialize(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)

	}

	userService := services.NewUserService(database)

	authHandler := handlers.NewAuthHandler(database, cfg, userService)

	userHandler := handlers.NewUserHandler(database)
	accountHandler := handlers.NewAccountHandler(database)
	transactionHandler := handlers.NewTransactionHandler(database)

	router := api.SetupRouter(
		database,
		cfg,
		authHandler,
		userHandler,
		accountHandler,
		transactionHandler,
	)

	address := cfg.Server.Address
	if address == "" {
		address = ":8081"
	}

	log.Printf("Server starting on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
