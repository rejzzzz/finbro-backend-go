// internal/api/handlers/account.go
package handlers

import (
	"net/http"
	"strconv"

	"finbro-backend-go/internal/db"
	"finbro-backend-go/internal/db/models"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	db *db.DB
}

func NewAccountHandler(db *db.DB) *AccountHandler {
	return &AccountHandler{db: db}
}

type CreateAccountRequest struct {
	AccountName   string  `json:"account_name" binding:"required"`
	AccountType   string  `json:"account_type"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
	BankName      string  `json:"bank_name"`
	AccountNumber string  `json:"account_number"`
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var accounts []models.Account
	if err := h.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := &models.Account{
		UserID:        userID.(uint),
		AccountName:   req.AccountName,
		AccountType:   req.AccountType,
		Balance:       req.Balance,
		Currency:      req.Currency,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
	}

	if err := h.db.Create(account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	accountID, _ := strconv.Atoi(c.Param("id"))

	var account models.Account
	if err := h.db.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	accountID, _ := strconv.Atoi(c.Param("id"))

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account models.Account
	if err := h.db.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	account.AccountName = req.AccountName
	account.AccountType = req.AccountType
	account.BankName = req.BankName
	account.AccountNumber = req.AccountNumber

	if err := h.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	accountID, _ := strconv.Atoi(c.Param("id"))

	if err := h.db.Where("id = ? AND user_id = ?", accountID, userID).Delete(&models.Account{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
