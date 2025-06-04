package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

// Create creates a new account
func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

// GetByID retrieves an account by ID with user information
func (r *accountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.Preload("User").First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

// GetByUserID retrieves all accounts for a user
func (r *accountRepository) GetByUserID(userID uuid.UUID) ([]*models.Account, error) {
	var accounts []*models.Account
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&accounts).Error
	return accounts, err
}

// GetActiveByUserID retrieves all active accounts for a user
func (r *accountRepository) GetActiveByUserID(userID uuid.UUID) ([]*models.Account, error) {
	var accounts []*models.Account
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").Find(&accounts).Error
	return accounts, err
}

// Update updates an account
func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

// UpdateBalance updates only the balance of an account
func (r *accountRepository) UpdateBalance(id uuid.UUID, balance decimal.Decimal) error {
	result := r.db.Model(&models.Account{}).Where("id = ?", id).Update("balance", balance)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}

// Delete deletes an account by ID
func (r *accountRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Account{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}
