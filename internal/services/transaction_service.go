package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"github.com/vasujain275/expense-tracker-api/internal/repositories"
)

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	accountRepo     repositories.AccountRepository
	categoryRepo    repositories.CategoryRepository
	userRepo        repositories.UserRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	accountRepo repositories.AccountRepository,
	categoryRepo repositories.CategoryRepository,
	userRepo repositories.UserRepository,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		userRepo:        userRepo,
	}
}

// CreateTransaction creates a new transaction and updates account balance
func (s *transactionService) CreateTransaction(req TransactionCreateRequest) (*models.Transaction, error) {
	// Validate request
	if err := s.validateTransactionCreateRequest(req); err != nil {
		return nil, err
	}

	// Verify user exists
	exists, err := s.userRepo.Exists(req.UserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Verify account exists and belongs to user
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, errors.New("account not found")
	}
	if account.UserID != req.UserID {
		return nil, errors.New("account does not belong to user")
	}

	// Verify category exists
	exists, err = s.categoryRepo.Exists(req.CategoryID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("category not found")
	}

	// Create transaction
	transaction := &models.Transaction{
		UserID:      req.UserID,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: strings.TrimSpace(req.Description),
		Date:        req.Date,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	// Update account balance
	newBalance := account.Balance.Add(req.Amount)
	if err := s.accountRepo.UpdateBalance(req.AccountID, newBalance); err != nil {
		// TODO: Consider implementing transaction rollback here
		return nil, errors.New("failed to update account balance")
	}

	// Get transaction with related data
	return s.transactionRepo.GetByID(transaction.ID)
}

// GetTransactionByID retrieves a transaction by ID
func (s *transactionService) GetTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid transaction ID")
	}

	return s.transactionRepo.GetByID(id)
}

// GetTransactions retrieves transactions with filtering and pagination
func (s *transactionService) GetTransactions(req TransactionListRequest) ([]*models.Transaction, int64, error) {
	// Validate request
	if err := s.validateTransactionListRequest(req); err != nil {
		return nil, 0, err
	}

	// Verify user exists
	exists, err := s.userRepo.Exists(req.UserID)
	if err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, errors.New("user not found")
	}

	// Create filter
	filter := repositories.TransactionFilter{
		UserID:     req.UserID,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		MinAmount:  req.MinAmount,
		MaxAmount:  req.MaxAmount,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	// Get transactions and count
	transactions, err := s.transactionRepo.GetByFilter(filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.transactionRepo.Count(filter)
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

// UpdateTransaction updates a transaction and adjusts account balances
func (s *transactionService) UpdateTransaction(id uuid.UUID, req TransactionUpdateRequest) (*models.Transaction, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid transaction ID")
	}

	// Get existing transaction
	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Store old values for balance adjustment
	oldAccountID := transaction.AccountID
	oldAmount := transaction.Amount

	// Validate and update fields
	if req.AccountID != nil {
		// Verify new account exists and belongs to user
		account, err := s.accountRepo.GetByID(*req.AccountID)
		if err != nil {
			return nil, errors.New("account not found")
		}
		if account.UserID != transaction.UserID {
			return nil, errors.New("account does not belong to user")
		}
		transaction.AccountID = *req.AccountID
	}

	if req.CategoryID != nil {
		// Verify category exists
		exists, err := s.categoryRepo.Exists(*req.CategoryID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("category not found")
		}
		transaction.CategoryID = *req.CategoryID
	}

	if req.Amount != nil {
		if req.Amount.IsZero() {
			return nil, errors.New("transaction amount cannot be zero")
		}
		transaction.Amount = *req.Amount
	}

	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		if desc == "" {
			return nil, errors.New("transaction description cannot be empty")
		}
		transaction.Description = desc
	}

	if req.Date != nil {
		transaction.Date = *req.Date
	}

	// Update transaction
	if err := s.transactionRepo.Update(transaction); err != nil {
		return nil, err
	}

	// Update account balances if account or amount changed
	if req.AccountID != nil || req.Amount != nil {
		// Revert old transaction from old account
		if oldAccountID != uuid.Nil {
			oldAccount, err := s.accountRepo.GetByID(oldAccountID)
			if err == nil {
				revertedBalance := oldAccount.Balance.Sub(oldAmount)
				s.accountRepo.UpdateBalance(oldAccountID, revertedBalance)
			}
		}

		// Apply new transaction to new account
		newAccount, err := s.accountRepo.GetByID(transaction.AccountID)
		if err == nil {
			newBalance := newAccount.Balance.Add(transaction.Amount)
			s.accountRepo.UpdateBalance(transaction.AccountID, newBalance)
		}
	}

	// Get updated transaction with related data
	return s.transactionRepo.GetByID(transaction.ID)
}

// DeleteTransaction deletes a transaction and adjusts account balance
func (s *transactionService) DeleteTransaction(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid transaction ID")
	}

	// Get transaction to adjust account balance
	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete transaction
	if err := s.transactionRepo.Delete(id); err != nil {
		return err
	}

	// Adjust account balance (reverse the transaction)
	account, err := s.accountRepo.GetByID(transaction.AccountID)
	if err == nil {
		newBalance := account.Balance.Sub(transaction.Amount)
		s.accountRepo.UpdateBalance(transaction.AccountID, newBalance)
	}

	return nil
}

// GetTransactionSummary gets spending summary by category
func (s *transactionService) GetTransactionSummary(userID uuid.UUID, startDate, endDate *time.Time) ([]*repositories.TransactionSummary, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	// Verify user exists
	exists, err := s.userRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	return s.transactionRepo.GetSummaryByCategory(userID, startDate, endDate)
}

// GetMonthlyTotal gets total transactions for a specific month
func (s *transactionService) GetMonthlyTotal(userID uuid.UUID, year int, month int) (decimal.Decimal, error) {
	if userID == uuid.Nil {
		return decimal.Zero, errors.New("invalid user ID")
	}

	// Verify user exists
	exists, err := s.userRepo.Exists(userID)
	if err != nil {
		return decimal.Zero, err
	}
	if !exists {
		return decimal.Zero, errors.New("user not found")
	}

	loc := time.UTC
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	endDate := startDate.AddDate(0, 1, -1)

	return s.transactionRepo.GetTotalByDateRange(userID, startDate, endDate)
}

// validateTransactionCreateRequest validates the transaction creation request
func (s *transactionService) validateTransactionCreateRequest(req TransactionCreateRequest) error {
	if req.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if req.AccountID == uuid.Nil {
		return errors.New("account ID is required")
	}
	if req.CategoryID == uuid.Nil {
		return errors.New("category ID is required")
	}
	if req.Amount.IsZero() {
		return errors.New("transaction amount cannot be zero")
	}
	if req.Description == "" {
		return errors.New("transaction description cannot be empty")
	}
	if req.Date.IsZero() {
		return errors.New("transaction date is required")
	}
	return nil
}

// validateTransactionListRequest validates the transaction list request
func (s *transactionService) validateTransactionListRequest(req TransactionListRequest) error {
	// Add validation logic if needed
	return nil
}
