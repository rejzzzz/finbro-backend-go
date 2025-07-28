// internal/auth/oauth.go
package auth

import (
	"context"

	"golang.org/x/oauth2"
)

type OAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error)
}

type OAuthService struct {
	Google OAuthProvider
}

func NewOAuthService(google OAuthProvider) *OAuthService {
	return &OAuthService{
		Google: google,
	}
}
