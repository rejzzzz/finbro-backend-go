// File 12: internal/services/user_service.go
package services

import (
	"finbro-backend-go/internal/db"
	"finbro-backend-go/internal/db/models"
)

type UserService struct {
	db *db.DB
}

func NewUserService(db *db.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.Preload("Accounts").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

func (s *UserService) GetUserStats(userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total accounts
	var accountCount int64
	s.db.Model(&models.Account{}).Where("user_id = ?", userID).Count(&accountCount)
	stats["total_accounts"] = accountCount

	// Total transactions
	var transactionCount int64
	s.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&transactionCount)
	stats["total_transactions"] = transactionCount

	// Total balance
	var totalBalance float64
	s.db.Model(&models.Account{}).Where("user_id = ?", userID).Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance)
	stats["total_balance"] = totalBalance

	return stats, nil
}
