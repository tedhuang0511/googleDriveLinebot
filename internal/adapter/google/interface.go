package google

import (
	"context"

	"golang.org/x/oauth2"
)

type GoogleOAuthI interface {
	OAuthLoginURL() (oauthURL string)
	UserOAuthToken(authCode string) (*oauth2.Token, error)
	NewGoogleDrive(ctx context.Context, tok *oauth2.Token) (*GoogleDrive, error)
}
