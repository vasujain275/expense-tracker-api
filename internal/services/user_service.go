package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"github.com/vasujain275/expense-tracker-api/internal/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with validation
func (s *userService) CreateUser(email, name, currency string) (*models.User, error) {
	// Validate inputs
	if err := s.validateUserInput(email, name, currency); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create user
	user := &models.User{
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Name:     strings.TrimSpace(name),
		Currency: strings.ToUpper(strings.TrimSpace(currency)),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	return s.userRepo.GetByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	return s.userRepo.GetByEmail(strings.ToLower(strings.TrimSpace(email)))
}

// UpdateUser updates user information
func (s *userService) UpdateUser(id uuid.UUID, name, currency string) (*models.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate inputs
	if name != "" {
		if len(strings.TrimSpace(name)) < 2 {
			return nil, errors.New("name must be at least 2 characters long")
		}
		user.Name = strings.TrimSpace(name)
	}

	if currency != "" {
		if len(strings.TrimSpace(currency)) != 3 {
			return nil, errors.New("currency must be a 3-letter code (e.g., USD, EUR)")
		}
		user.Currency = strings.ToUpper(strings.TrimSpace(currency))
	}

	// Update user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	exists, err := s.userRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

// validateUserInput validates user input fields
func (s *userService) validateUserInput(email, name, currency string) error {
	// Validate email
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	if !s.isValidEmail(email) {
		return errors.New("invalid email format")
	}

	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	// Validate currency
	currency = strings.TrimSpace(currency)
	if currency == "" {
		return errors.New("currency is required")
	}
	if len(currency) != 3 {
		return errors.New("currency must be a 3-letter code (e.g., USD, EUR)")
	}

	return nil
}

// isValidEmail performs basic email validation
func (s *userService) isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
