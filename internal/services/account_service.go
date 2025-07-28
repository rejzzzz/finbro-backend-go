// File 13: internal/services/account_service.go
package services

import (
	"finbro-backend-go/internal/db"
	"finbro-backend-go/internal/db/models"
)

type AccountService struct {
	db *db.DB
}

func NewAccountService(db *db.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) GetUserAccounts(userID uint) ([]models.Account, error) {
	var accounts []models.Account
	err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&accounts).Error
	return accounts, err
}

func (s *AccountService) GetAccountByID(accountID, userID uint) (*models.Account, error) {
	var account models.Account
	err := s.db.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error
	return &account, err
}

func (s *AccountService) CreateAccount(account *models.Account) error {
	return s.db.Create(account).Error
}

func (s *AccountService) UpdateAccount(account *models.Account) error {
	return s.db.Save(account).Error
}

func (s *AccountService) DeleteAccount(accountID, userID uint) error {
	return s.db.Where("id = ? AND user_id = ?", accountID, userID).Delete(&models.Account{}).Error
}

func (s *AccountService) UpdateBalance(accountID uint, amount float64) error {
	return s.db.Model(&models.Account{}).Where("id = ?", accountID).
		Update("balance", amount).Error
}

func (s *AccountService) GetAccountBalance(accountID, userID uint) (float64, error) {
	var balance float64
	err := s.db.Model(&models.Account{}).
		Where("id = ? AND user_id = ?", accountID, userID).
		Select("balance").Scan(&balance).Error
	return balance, err
}
