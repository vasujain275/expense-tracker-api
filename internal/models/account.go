package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeBank       AccountType = "bank"
	AccountTypeCash       AccountType = "cash"
	AccountTypeCreditCard AccountType = "credit_card"
)

type Account struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Name      string          `json:"name" gorm:"not null"`
	Type      AccountType     `json:"type" gorm:"not null"`
	Balance   decimal.Decimal `json:"balance" gorm:"type:decimal(15,2);not null;default:0"`
	IsActive  bool            `json:"is_active" gorm:"not null;default:true"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`

	// Relationships
	User         User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:AccountID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Account models
func (Account) TableName() string {
	return "accounts"
}
