package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
)

// UserHandler interface defines methods for user-related HTTP handlers
type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

// AccountHandler interface defines methods for account-related HTTP handlers
type AccountHandler interface {
	CreateAccount(c *gin.Context)
	GetAccount(c *gin.Context)
	GetUserAccounts(c *gin.Context)
	UpdateAccount(c *gin.Context)
	DeleteAccount(c *gin.Context)
	GetAccountBalance(c *gin.Context)
}

// CategoryHandler interface defines methods for category-related HTTP handlers
type CategoryHandler interface {
	CreateCategory(c *gin.Context)
	GetCategory(c *gin.Context)
	GetAllCategories(c *gin.Context)
	GetCategoriesByType(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

// TransactionHandler interface defines methods for transaction-related HTTP handlers
type TransactionHandler interface {
	CreateTransaction(c *gin.Context)
	GetTransaction(c *gin.Context)
	GetTransactions(c *gin.Context)
	UpdateTransaction(c *gin.Context)
	DeleteTransaction(c *gin.Context)
	GetTransactionSummary(c *gin.Context)
	GetMonthlyTotal(c *gin.Context)
}

// Request/Response structs for handlers

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR GBP JPY"`
}

// CreateAccountRequest represents a request to create an account
type CreateAccountRequest struct {
	Name           string             `json:"name" binding:"required"`
	Type           models.AccountType `json:"type" binding:"required,oneof=bank cash credit_card"`
	InitialBalance decimal.Decimal    `json:"initial_balance" binding:"required"`
}

// UpdateAccountRequest represents a request to update an account
type UpdateAccountRequest struct {
	Name     string             `json:"name,omitempty"`
	Type     models.AccountType `json:"type,omitempty" binding:"omitempty,oneof=bank cash credit_card"`
	IsActive bool               `json:"is_active,omitempty"`
}

// CreateCategoryRequest represents a request to create a category
type CreateCategoryRequest struct {
	Name  string              `json:"name" binding:"required"`
	Type  models.CategoryType `json:"type" binding:"required,oneof=income expense"`
	Color string              `json:"color" binding:"required,hexcolor"`
}

// UpdateCategoryRequest represents a request to update a category
type UpdateCategoryRequest struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty" binding:"omitempty,hexcolor"`
}

// CreateTransactionRequest represents a request to create a transaction
type CreateTransactionRequest struct {
	AccountID   uuid.UUID       `json:"account_id" binding:"required"`
	CategoryID  uuid.UUID       `json:"category_id" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Description string          `json:"description" binding:"required"`
	Date        string          `json:"date" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// UpdateTransactionRequest represents a request to update a transaction
type UpdateTransactionRequest struct {
	AccountID   *uuid.UUID       `json:"account_id,omitempty"`
	CategoryID  *uuid.UUID       `json:"category_id,omitempty"`
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	Description *string          `json:"description,omitempty"`
	Date        *string          `json:"date,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

// TransactionListRequest represents a request to list transactions
type TransactionListRequest struct {
	AccountID  *uuid.UUID       `json:"account_id,omitempty"`
	CategoryID *uuid.UUID       `json:"category_id,omitempty"`
	StartDate  *string          `json:"start_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	EndDate    *string          `json:"end_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	MinAmount  *decimal.Decimal `json:"min_amount,omitempty"`
	MaxAmount  *decimal.Decimal `json:"max_amount,omitempty"`
	Limit      int              `json:"limit" binding:"min=1,max=100"`
	Offset     int              `json:"offset" binding:"min=0"`
}
