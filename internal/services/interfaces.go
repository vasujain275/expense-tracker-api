package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"github.com/vasujain275/expense-tracker-api/internal/repositories"
)

// UserService interface defines business logic for user operations
type UserService interface {
	CreateUser(email, name, currency string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uuid.UUID, name, currency string) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

// AccountService interface defines business logic for account operations
type AccountService interface {
	CreateAccount(userID uuid.UUID, name string, accountType models.AccountType, initialBalance decimal.Decimal) (*models.Account, error)
	GetAccountByID(id uuid.UUID) (*models.Account, error)
	GetUserAccounts(userID uuid.UUID, activeOnly bool) ([]*models.Account, error)
	UpdateAccount(id uuid.UUID, name string, accountType models.AccountType, isActive bool) (*models.Account, error)
	DeleteAccount(id uuid.UUID) error
	GetAccountBalance(id uuid.UUID) (decimal.Decimal, error)
}

// CategoryService interface defines business logic for category operations
type CategoryService interface {
	CreateCategory(name string, categoryType models.CategoryType, color string) (*models.Category, error)
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	GetAllCategories() ([]*models.Category, error)
	GetCategoriesByType(categoryType models.CategoryType) ([]*models.Category, error)
	UpdateCategory(id uuid.UUID, name, color string) (*models.Category, error)
	DeleteCategory(id uuid.UUID) error
}

// TransactionCreateRequest represents a request to create a transaction
type TransactionCreateRequest struct {
	UserID      uuid.UUID       `json:"user_id"`
	AccountID   uuid.UUID       `json:"account_id"`
	CategoryID  uuid.UUID       `json:"category_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
}

// TransactionUpdateRequest represents a request to update a transaction
type TransactionUpdateRequest struct {
	AccountID   *uuid.UUID       `json:"account_id,omitempty"`
	CategoryID  *uuid.UUID       `json:"category_id,omitempty"`
	Amount      *decimal.Decimal `json:"amount,omitempty"`
	Description *string          `json:"description,omitempty"`
	Date        *time.Time       `json:"date,omitempty"`
}

// TransactionListRequest represents a request to list transactions
type TransactionListRequest struct {
	UserID     uuid.UUID        `json:"user_id"`
	AccountID  *uuid.UUID       `json:"account_id,omitempty"`
	CategoryID *uuid.UUID       `json:"category_id,omitempty"`
	StartDate  *time.Time       `json:"start_date,omitempty"`
	EndDate    *time.Time       `json:"end_date,omitempty"`
	MinAmount  *decimal.Decimal `json:"min_amount,omitempty"`
	MaxAmount  *decimal.Decimal `json:"max_amount,omitempty"`
	Limit      int              `json:"limit"`
	Offset     int              `json:"offset"`
}

// TransactionService interface defines business logic for transaction operations
type TransactionService interface {
	CreateTransaction(req TransactionCreateRequest) (*models.Transaction, error)
	GetTransactionByID(id uuid.UUID) (*models.Transaction, error)
	GetTransactions(req TransactionListRequest) ([]*models.Transaction, int64, error)
	UpdateTransaction(id uuid.UUID, req TransactionUpdateRequest) (*models.Transaction, error)
	DeleteTransaction(id uuid.UUID) error
	GetTransactionSummary(userID uuid.UUID, startDate, endDate *time.Time) ([]*repositories.TransactionSummary, error)
	GetMonthlyTotal(userID uuid.UUID, year int, month int) (decimal.Decimal, error)
}
