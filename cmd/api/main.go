package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vasujain275/expense-tracker-api/internal/config"
	"github.com/vasujain275/expense-tracker-api/internal/database"
	"github.com/vasujain275/expense-tracker-api/internal/repositories"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := database.Connect(cfg.PostgresConnectionDsn())

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

}
