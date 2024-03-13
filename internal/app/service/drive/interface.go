package drive

import (
	"context"
	"firstProject/internal/adapter/google"

	"golang.org/x/oauth2"
)

type DriveServiceGoogleOAuthI interface {
	OAuthLoginURL() (oauthURL string)
	UserOAuthToken(authCode string) (*oauth2.Token, error)
	NewGoogleDrive(ctx context.Context, tok *oauth2.Token) (*google.GoogleDrive, error)
}
