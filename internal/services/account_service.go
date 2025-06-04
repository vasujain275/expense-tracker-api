package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"github.com/vasujain275/expense-tracker-api/internal/repositories"
)

type accountService struct {
	accountRepo repositories.AccountRepository
	userRepo    repositories.UserRepository
}

// NewAccountService creates a new account service
func NewAccountService(accountRepo repositories.AccountRepository, userRepo repositories.UserRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
	}
}

// CreateAccount creates a new account for a user
func (s *accountService) CreateAccount(userID uuid.UUID, name string, accountType models.AccountType, initialBalance decimal.Decimal) (*models.Account, error) {
	// Validate inputs
	if err := s.validateAccountInput(userID, name, accountType); err != nil {
		return nil, err
	}

	// Check if user exists
	exists, err := s.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Create account
	account := &models.Account{
		UserID:   userID,
		Name:     strings.TrimSpace(name),
		Type:     accountType,
		Balance:  initialBalance,
		IsActive: true,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByID retrieves an account by ID
func (s *accountService) GetAccountByID(id uuid.UUID) (*models.Account, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid account ID")
	}

	return s.accountRepo.GetByID(id)
}

// GetUserAccounts retrieves all accounts for a user
func (s *accountService) GetUserAccounts(userID uuid.UUID, activeOnly bool) ([]*models.Account, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if user exists
	exists, err := s.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	if activeOnly {
		return s.accountRepo.GetActiveByUserID(userID)
	}

	return s.accountRepo.GetByUserID(userID)
}

// UpdateAccount updates an account
func (s *accountService) UpdateAccount(id uuid.UUID, name string, accountType models.AccountType, isActive bool) (*models.Account, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid account ID")
	}

	// Get existing account
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate and update fields
	if name != "" {
		name = strings.TrimSpace(name)
		if len(name) < 2 {
			return nil, errors.New("account name must be at least 2 characters long")
		}
		account.Name = name
	}

	if accountType != "" {
		if !s.isValidAccountType(accountType) {
			return nil, errors.New("invalid account type")
		}
		account.Type = accountType
	}

	account.IsActive = isActive

	// Update account
	if err := s.accountRepo.Update(account); err != nil {
		return nil, err
	}

	return account, nil
}

// DeleteAccount deletes an account
func (s *accountService) DeleteAccount(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid account ID")
	}

	// Check if account exists
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if account has zero balance before deletion
	if !account.Balance.IsZero() {
		return errors.New("cannot delete account with non-zero balance")
	}

	return s.accountRepo.Delete(id)
}

// GetAccountBalance retrieves the current balance of an account
func (s *accountService) GetAccountBalance(id uuid.UUID) (decimal.Decimal, error) {
	if id == uuid.Nil {
		return decimal.Zero, errors.New("invalid account ID")
	}

	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return decimal.Zero, err
	}

	return account.Balance, nil
}

// validateAccountInput validates account input fields
func (s *accountService) validateAccountInput(userID uuid.UUID, name string, accountType models.AccountType) error {
	// Validate user ID
	if userID == uuid.Nil {
		return errors.New("user ID is required")
	}

	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("account name is required")
	}
	if len(name) < 2 {
		return errors.New("account name must be at least 2 characters long")
	}

	// Validate account type
	if !s.isValidAccountType(accountType) {
		return errors.New("invalid account type")
	}

	return nil
}

// isValidAccountType checks if the account type is valid
func (s *accountService) isValidAccountType(accountType models.AccountType) bool {
	switch accountType {
	case models.AccountTypeBank, models.AccountTypeCash, models.AccountTypeCreditCard:
		return true
	default:
		return false
	}
}
