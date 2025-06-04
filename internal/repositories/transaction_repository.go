package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create creates a new transaction
func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// GetByID retrieves a transaction by ID with related data
func (r *transactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Account").Preload("Category").First(&transaction, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

// GetByFilter retrieves transactions based on filter criteria
func (r *transactionRepository) GetByFilter(filter TransactionFilter) ([]*models.Transaction, error) {
	query := r.db.Preload("Account").Preload("Category").Where("user_id = ?", filter.UserID)

	// Apply filters
	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.StartDate != nil {
		query = query.Where("date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", *filter.EndDate)
	}
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var transactions []*models.Transaction
	err := query.Order("date DESC, created_at DESC").Find(&transactions).Error
	return transactions, err
}

// Update updates a transaction
func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

// Delete deletes a transaction by ID
func (r *transactionRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Transaction{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

// GetSummaryByCategory gets spending summary grouped by category
func (r *transactionRepository) GetSummaryByCategory(userID uuid.UUID, startDate, endDate *time.Time) ([]*TransactionSummary, error) {
	query := r.db.Table("transactions").
		Select("category_id, categories.name as category_name, SUM(ABS(amount)) as total_amount, COUNT(*) as count").
		Joins("JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.user_id = ?", userID).
		Group("category_id, categories.name")

	if startDate != nil {
		query = query.Where("transactions.date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("transactions.date <= ?", *endDate)
	}

	var summaries []*TransactionSummary
	err := query.Order("total_amount DESC").Find(&summaries).Error
	return summaries, err
}

// GetTotalByDateRange gets total transaction amount for a date range
func (r *transactionRepository) GetTotalByDateRange(userID uuid.UUID, startDate, endDate time.Time) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// Count counts transactions based on filter criteria
func (r *transactionRepository) Count(filter TransactionFilter) (int64, error) {
	query := r.db.Model(&models.Transaction{}).Where("user_id = ?", filter.UserID)

	// Apply same filters as GetByFilter
	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.StartDate != nil {
		query = query.Where("date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", *filter.EndDate)
	}
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}
