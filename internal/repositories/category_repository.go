package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create creates a new category
func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetByID retrieves a category by ID
func (r *categoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// GetAll retrieves all categories
func (r *categoryRepository) GetAll() ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

// GetByType retrieves categories by type
func (r *categoryRepository) GetByType(categoryType models.CategoryType) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Where("type = ?", categoryType).Order("name ASC").Find(&categories).Error
	return categories, err
}

// Update updates a category
func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category by ID
func (r *categoryRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Category{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// Exists checks if a category exists by ID
func (r *categoryRepository) Exists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
