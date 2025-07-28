// internal/utils/validation.go
package utils

import (
	"regexp"
	"strings"
)

var (
	emailRegex      = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	validCurrencies = map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "INR": true,
		"CAD": true, "AUD": true, "JPY": true, "CNY": true,
	}
	validAccountTypes = map[string]bool{
		"checking": true, "savings": true, "credit": true,
		"investment": true, "loan": true, "business": true,
	}
)

// Email validation
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	return len(email) > 0 && len(email) <= 254 && emailRegex.MatchString(email)
}

// Password validation
func IsValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 128
}

// Currency validation
func IsValidCurrency(currency string) bool {
	return validCurrencies[strings.ToUpper(strings.TrimSpace(currency))]
}

// Account type validation
func IsValidAccountType(accountType string) bool {
	return validAccountTypes[strings.ToLower(strings.TrimSpace(accountType))]
}

// Transaction type validation
func IsValidTransactionType(txType string) bool {
	txType = strings.ToLower(strings.TrimSpace(txType))
	return txType == "debit" || txType == "credit"
}

// Amount validation
func IsValidAmount(amount float64) bool {
	return amount > 0 && amount <= 999999999.99
}

// String sanitization
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

// Name validation (for first/last names)
func IsValidName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) >= 1 && len(name) <= 50
}
