// internal/services/transaction_service.go
package services

import (
	"finbro-backend-go/internal/db"
	"finbro-backend-go/internal/db/models"
	"time"
)

type TransactionService struct {
	db *db.DB
}

func NewTransactionService(db *db.DB) *TransactionService {
	return &TransactionService{db: db}
}

type TransactionFilter struct {
	UserID    uint
	AccountID uint
	Category  string
	StartDate time.Time
	EndDate   time.Time
	Limit     int
	Offset    int
}

func (s *TransactionService) GetTransactions(filter TransactionFilter) ([]models.Transaction, error) {
	query := s.db.Where("user_id = ?", filter.UserID)

	if filter.AccountID > 0 {
		query = query.Where("account_id = ?", filter.AccountID)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if !filter.StartDate.IsZero() {
		query = query.Where("transaction_date >= ?", filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		query = query.Where("transaction_date <= ?", filter.EndDate)
	}

	var transactions []models.Transaction
	err := query.Preload("Account").
		Order("transaction_date DESC").
		Limit(filter.Limit).Offset(filter.Offset).
		Find(&transactions).Error

	return transactions, err
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	tx := s.db.Begin()

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update account balance
	balanceChange := transaction.Amount
	if transaction.Type == "debit" {
		balanceChange = -balanceChange
	}

	if err := tx.Model(&models.Account{}).Where("id = ?", transaction.AccountID).
		Update("balance", models.Account{}.Balance+balanceChange).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *TransactionService) GetTransactionByID(transactionID, userID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := s.db.Preload("Account").
		Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error
	return &transaction, err
}

func (s *TransactionService) UpdateTransaction(transaction *models.Transaction) error {
	return s.db.Save(transaction).Error
}

func (s *TransactionService) DeleteTransaction(transactionID, userID uint) error {
	return s.db.Where("id = ? AND user_id = ?", transactionID, userID).
		Delete(&models.Transaction{}).Error
}

func (s *TransactionService) GetCategoryStats(userID uint, startDate, endDate time.Time) (map[string]float64, error) {
	var results []struct {
		Category string
		Total    float64
	}

	query := s.db.Model(&models.Transaction{}).
		Select("category, SUM(amount) as total").
		Where("user_id = ?", userID).
		Group("category")

	if !startDate.IsZero() {
		query = query.Where("transaction_date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("transaction_date <= ?", endDate)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[string]float64)
	for _, result := range results {
		stats[result.Category] = result.Total
	}

	return stats, nil
}
