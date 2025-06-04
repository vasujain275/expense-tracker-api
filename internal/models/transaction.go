package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	AccountID   uuid.UUID       `json:"account_id" gorm:"type:uuid;not null"`
	CategoryID  uuid.UUID       `json:"category_id" gorm:"type:uuid;not null"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:decimal(15,2);not null"`
	Description string          `json:"description" gorm:"not null"`
	Date        time.Time       `json:"date" gorm:"not null;index"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`

	// Relationships
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Account  Account  `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// IsIncome returns true if the transaction amount is positive (income)
func (t *Transaction) IsIncome() bool {
	return t.Amount.GreaterThan(decimal.Zero)
}

// IsExpense returns true if the transaction amount is negative (expense)
func (t *Transaction) IsExpense() bool {
	return t.Amount.LessThan(decimal.Zero)
}

// AbsAmount returns the absolute value of the transaction amount
func (t *Transaction) AbsAmount() decimal.Decimal {
	return t.Amount.Abs()
}
