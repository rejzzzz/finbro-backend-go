// File 11: internal/api/handlers/transaction.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"finbro-backend-go/internal/db"
	"finbro-backend-go/internal/db/models"
)

type TransactionHandler struct {
	db *db.DB
}

func NewTransactionHandler(db *db.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

type CreateTransactionRequest struct {
	AccountID       uint      `json:"account_id" binding:"required"`
	Amount          float64   `json:"amount" binding:"required"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Type            string    `json:"type" binding:"required,oneof=debit credit"`
	TransactionDate time.Time `json:"transaction_date"`
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Query parameters
	accountID := c.Query("account_id")
	category := c.Query("category")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	query := h.db.Where("user_id = ?", userID)
	
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var transactions []models.Transaction
	if err := query.Preload("Account").
		Order("transaction_date DESC").
		Limit(limit).Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify account belongs to user
	var account models.Account
	if err := h.db.Where("id = ? AND user_id = ?", req.AccountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	transaction := &models.Transaction{
		UserID:          userID.(uint),
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		Description:     req.Description,
		Category:        req.Category,
		Type:            req.Type,
		TransactionDate: req.TransactionDate,
	}

	if transaction.TransactionDate.IsZero() {
		transaction.TransactionDate = time.Now()
	}

	// Start transaction to update account balance
	tx := h.db.Begin()
	
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Update account balance
	balanceChange := req.Amount
	if req.Type == "debit" {
		balanceChange = -balanceChange
	}
	
	if err := tx.Model(&account).Update("balance", account.Balance+balanceChange).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")
	transactionID, _ := strconv.Atoi(c.Param("id"))

	var transaction models.Transaction
	if err := h.db.Preload("Account").
		Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")
	transactionID, _ := strconv.Atoi(c.Param("id"))

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	if err := h.db.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	transaction.Description = req.Description
	transaction.Category = req.Category
	transaction.TransactionDate = req.TransactionDate

	if err := h.db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")
	transactionID, _ := strconv.Atoi(c.Param("id"))

	if err := h.db.Where("id = ? AND user_id = ?", transactionID, userID).Delete(&models.Transaction{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}