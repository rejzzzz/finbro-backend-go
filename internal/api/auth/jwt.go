// internal/auth/jwt.go
package auth

import (
	"fmt"
	"strconv"
	"time"

	"finbro-backend-go/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

type JWTAuth struct {
	secret     []byte
	expiry     time.Duration
	signingAlg *jwt.SigningMethodHMAC
}

func NewJWTAuth(cfg *config.Config) *JWTAuth {
	return &JWTAuth{
		secret:     []byte(cfg.JWT.Secret),
		expiry:     cfg.JWT.Expiry,
		signingAlg: jwt.SigningMethodHS256,
	}
}

func (j *JWTAuth) GenerateToken(userID uint, email string) (string, error) {
	subject := strconv.FormatUint(uint64(userID), 10)

	claims := &jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "finbro-backend",
		Audience:  []string{"client"},
	}

	token := jwt.NewWithClaims(j.signingAlg, claims)
	return token.SignedString(j.secret)
}

func (j *JWTAuth) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != j.signingAlg {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims or token not valid")
	}

	return claims, nil
}
