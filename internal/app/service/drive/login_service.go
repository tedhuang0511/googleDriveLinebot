package drive

import (
	"context"
	"firstProject/internal/adapter/dynamodb"
)

func (dr *GoogleDriveService) Login(ctx context.Context, lineID string, authCode string) error {
	tok, err := dr.driveServiceGoogleOA.UserOAuthToken(authCode)
	if err != nil {
		return err
	}

	dToken := dynamodb.GoogleOAuthToken{
		PK:           lineID,
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry.String(),
	}

	err = dr.driveServiceDynamodb.AddGoogleOAuthToken(dToken)
	if err != nil {
		return err
	}

	return nil
}

func (dr *GoogleDriveService) LoginURL(ctx context.Context, lineID string) string {
	oauthURL := dr.driveServiceGoogleOA.OAuthLoginURL(lineID)
	return oauthURL
}
