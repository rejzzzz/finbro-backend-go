// internal/services/user_service.go (Updated)
package services

import (
	"errors"
	"strings"

	"finbro-backend-go/internal/db"
	"finbro-backend-go/internal/db/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	var accountCount int64
	s.db.Model(&models.Account{}).Where("user_id = ?", userID).Count(&accountCount)
	stats["total_accounts"] = accountCount

	var transactionCount int64
	s.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&transactionCount)
	stats["total_transactions"] = transactionCount

	var totalBalance float64
	s.db.Model(&models.Account{}).Where("user_id = ?", userID).Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance)
	stats["total_balance"] = totalBalance

	return stats, nil
}

func (s *UserService) CreateOrUpdateOAuthUser(c *gin.Context, email, fullName string) (*models.User, error) {
	var user models.User

	// Try to find existing user by email
	err := s.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		// Check if it's a "not found" error vs actual database error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// User doesn't exist, create new one
		firstName, lastName := s.parseFullName(fullName)

		user = models.User{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			UserType:  s.determineUserType(email), // You can implement logic here
		}

		if createErr := s.db.Create(&user).Error; createErr != nil {
			return nil, createErr
		}
	}

	// Clear password field for security
	user.Password = ""
	return &user, nil
}

// Helper function to parse full name into first and last name
func (s *UserService) parseFullName(fullName string) (firstName, lastName string) {
	names := strings.Fields(strings.TrimSpace(fullName))
	if len(names) == 0 {
		return "", ""
	}

	firstName = names[0]
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	return firstName, lastName
}

func (s *UserService) determineUserType(email string) models.UserType {
	businessDomains := []string{"@company.com", "@business.org"}

	for _, domain := range businessDomains {
		if strings.Contains(email, domain) {
			return "business"
		}
	}

	return "individual" // default
}
