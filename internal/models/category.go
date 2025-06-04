package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string       `json:"name" gorm:"not null"`
	Type      CategoryType `json:"type" gorm:"not null"`
	Color     string       `json:"color" gorm:"not null;default:'#007bff'"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	// Relationships
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:CategoryID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "categories"
}
