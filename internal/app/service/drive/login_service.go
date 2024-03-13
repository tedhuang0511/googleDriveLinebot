package drive

import (
	"context"
	"firstProject/internal/adapter/dynamodb"
	domainDrive "firstProject/internal/domain/drive"
	"log"
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
		Expiry:       tok.Expiry,
	}
	err = dr.driveServiceDynamodb.AddGoogleOAuthToken(dToken)
	if err != nil {
		return err
	}

	return nil
}

func (dr *GoogleDriveService) LoginURL(ctx context.Context, lineID string) string {
	oauthURL := dr.driveServiceGoogleOA.OAuthLoginURL(lineID)
	resURL, err := domainDrive.AppendOpenExternalBrowserParam(oauthURL)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}
	return resURL
}
