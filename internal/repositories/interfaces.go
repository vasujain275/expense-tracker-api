package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
)

// UserRepository interface defines methods for user data access
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
}

// AccountRepository interface defines methods for account data access
type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id uuid.UUID) (*models.Account, error)
	GetByUserID(userID uuid.UUID) ([]*models.Account, error)
	Update(account *models.Account) error
	Delete(id uuid.UUID) error
	UpdateBalance(id uuid.UUID, balance decimal.Decimal) error
	GetActiveByUserID(userID uuid.UUID) ([]*models.Account, error)
}

// CategoryRepository interface defines methods for category data access
type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uuid.UUID) (*models.Category, error)
	GetAll() ([]*models.Category, error)
	GetByType(categoryType models.CategoryType) ([]*models.Category, error)
	Update(category *models.Category) error
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
}

// TransactionFilter holds filter parameters for transaction queries
type TransactionFilter struct {
	UserID     uuid.UUID
	AccountID  *uuid.UUID
	CategoryID *uuid.UUID
	StartDate  *time.Time
	EndDate    *time.Time
	MinAmount  *decimal.Decimal
	MaxAmount  *decimal.Decimal
	Limit      int
	Offset     int
}

// TransactionSummary represents spending summary by category
type TransactionSummary struct {
	CategoryID   uuid.UUID       `json:"category_id"`
	CategoryName string          `json:"category_name"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	Count        int64           `json:"count"`
}

// TransactionRepository interface defines methods for transaction data access
type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetByFilter(filter TransactionFilter) ([]*models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id uuid.UUID) error
	GetSummaryByCategory(userID uuid.UUID, startDate, endDate *time.Time) ([]*TransactionSummary, error)
	GetTotalByDateRange(userID uuid.UUID, startDate, endDate time.Time) (decimal.Decimal, error)
	Count(filter TransactionFilter) (int64, error)
}
