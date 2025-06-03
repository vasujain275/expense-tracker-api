package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vasujain275/expense-tracker-api/internal/config"
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

	router := gin.Default()

	dsn := cfg.PostgresConnectionDsn()

}
