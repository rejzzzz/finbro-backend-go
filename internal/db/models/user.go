// File 4: internal/db/models/user.go
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserType string

const (
	Individual UserType = "individual"
	Business   UserType = "business"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	Password    string    `json:"-" gorm:"not null"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserType    UserType  `json:"user_type" gorm:"default:individual"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GoogleID    string    `json:"-" gorm:"index"`
	Provider    string    `json:"provider" gorm:"default:email"`
	IsOAuthUser bool      `json:"is_oauth_user" gorm:"default:false"`

	// Relationships
	Accounts     []Account     `json:"accounts,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
	Budgets      []Budget      `json:"budgets,omitempty"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserType == "" {
		u.UserType = Individual
	}
	return nil
}

type Account struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	AccountName   string    `json:"account_name" gorm:"not null"`
	AccountType   string    `json:"account_type"`
	Balance       float64   `json:"balance" gorm:"default:0"`
	Currency      string    `json:"currency" gorm:"default:USD"`
	BankName      string    `json:"bank_name"`
	AccountNumber string    `json:"account_number"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	User         User          `json:"user,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

type Transaction struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	AccountID       uint      `json:"account_id" gorm:"not null"`
	Amount          float64   `json:"amount" gorm:"not null"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	TransactionDate time.Time `json:"transaction_date"`
	Type            string    `json:"type"` // debit, credit
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	User    User    `json:"user,omitempty"`
	Account Account `json:"account,omitempty"`
}

type Budget struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Category  string    `json:"category"`
	Amount    float64   `json:"amount" gorm:"not null"`
	Spent     float64   `json:"spent" gorm:"default:0"`
	Period    string    `json:"period" gorm:"default:monthly"` // monthly, weekly, yearly
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty"`
}
