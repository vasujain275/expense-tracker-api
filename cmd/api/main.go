package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vasujain275/expense-tracker-api/docs" // This is required for Swagger
	"github.com/vasujain275/expense-tracker-api/internal/config"
	"github.com/vasujain275/expense-tracker-api/internal/database"
	"github.com/vasujain275/expense-tracker-api/internal/handlers"
	"github.com/vasujain275/expense-tracker-api/internal/repositories"
	"github.com/vasujain275/expense-tracker-api/internal/services"
)

// @title           Expense Tracker API
// @version         1.0
// @description     A RESTful API for tracking personal expenses
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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

	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo, userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo, accountRepo, categoryRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Initialize router
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// User routes
		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)

		// Account routes
		v1.POST("/accounts", accountHandler.CreateAccount)
		v1.GET("/accounts", accountHandler.GetUserAccounts)
		v1.GET("/accounts/:id", accountHandler.GetAccount)
		v1.PUT("/accounts/:id", accountHandler.UpdateAccount)
		v1.DELETE("/accounts/:id", accountHandler.DeleteAccount)
		v1.GET("/accounts/:id/balance", accountHandler.GetAccountBalance)

		// Category routes
		v1.POST("/categories", categoryHandler.CreateCategory)
		v1.GET("/categories", categoryHandler.GetAllCategories)
		v1.GET("/categories/:id", categoryHandler.GetCategory)
		v1.PUT("/categories/:id", categoryHandler.UpdateCategory)
		v1.DELETE("/categories/:id", categoryHandler.DeleteCategory)
		v1.GET("/categories/type/:type", categoryHandler.GetCategoriesByType)

		// Transaction routes
		v1.POST("/transactions", transactionHandler.CreateTransaction)
		v1.GET("/transactions", transactionHandler.GetTransactions)
		v1.GET("/transactions/:id", transactionHandler.GetTransaction)
		v1.PUT("/transactions/:id", transactionHandler.UpdateTransaction)
		v1.DELETE("/transactions/:id", transactionHandler.DeleteTransaction)
		v1.GET("/transactions/summary", transactionHandler.GetTransactionSummary)
		v1.GET("/transactions/monthly-total", transactionHandler.GetMonthlyTotal)
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
