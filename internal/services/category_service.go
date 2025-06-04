package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"github.com/vasujain275/expense-tracker-api/internal/repositories"
)

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new category
func (s *categoryService) CreateCategory(name string, categoryType models.CategoryType, color string) (*models.Category, error) {
	// Validate inputs
	if err := s.validateCategoryInput(name, categoryType, color); err != nil {
		return nil, err
	}

	// Set default color if not provided
	if color == "" {
		color = "#007bff"
	}

	// Create category
	category := &models.Category{
		Name:  strings.TrimSpace(name),
		Type:  categoryType,
		Color: color,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category by ID
func (s *categoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid category ID")
	}

	return s.categoryRepo.GetByID(id)
}

// GetAllCategories retrieves all categories
func (s *categoryService) GetAllCategories() ([]*models.Category, error) {
	return s.categoryRepo.GetAll()
}

// GetCategoriesByType retrieves categories by type
func (s *categoryService) GetCategoriesByType(categoryType models.CategoryType) ([]*models.Category, error) {
	if !s.isValidCategoryType(categoryType) {
		return nil, errors.New("invalid category type")
	}

	return s.categoryRepo.GetByType(categoryType)
}

// UpdateCategory updates a category
func (s *categoryService) UpdateCategory(id uuid.UUID, name, color string) (*models.Category, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid category ID")
	}

	// Get existing category
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate and update fields
	if name != "" {
		name = strings.TrimSpace(name)
		if len(name) < 2 {
			return nil, errors.New("category name must be at least 2 characters long")
		}
		category.Name = name
	}

	if color != "" {
		if !s.isValidColor(color) {
			return nil, errors.New("invalid color format (expected hex color like #007bff)")
		}
		category.Color = color
	}

	// Update category
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid category ID")
	}

	// Check if category exists
	exists, err := s.categoryRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("category not found")
	}

	// TODO: Check if category is being used by any transactions
	// This would require a transaction repository dependency

	return s.categoryRepo.Delete(id)
}

// validateCategoryInput validates category input fields
func (s *categoryService) validateCategoryInput(name string, categoryType models.CategoryType, color string) error {
	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("category name is required")
	}
	if len(name) < 2 {
		return errors.New("category name must be at least 2 characters long")
	}

	// Validate category type
	if !s.isValidCategoryType(categoryType) {
		return errors.New("invalid category type")
	}

	// Validate color if provided
	if color != "" && !s.isValidColor(color) {
		return errors.New("invalid color format (expected hex color like #007bff)")
	}

	return nil
}

// isValidCategoryType checks if the category type is valid
func (s *categoryService) isValidCategoryType(categoryType models.CategoryType) bool {
	switch categoryType {
	case models.CategoryTypeIncome, models.CategoryTypeExpense:
		return true
	default:
		return false
	}
}

// isValidColor checks if the color is a valid hex color
func (s *categoryService) isValidColor(color string) bool {
	if len(color) != 7 {
		return false
	}
	if color[0] != '#' {
		return false
	}

	// Check if remaining characters are valid hex
	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}
