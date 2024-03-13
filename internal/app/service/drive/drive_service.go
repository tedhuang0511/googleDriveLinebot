package drive

import (
	"context"
	"log"
)

func (dr *GoogleDriveService) ListFiles(authCode string, ctx context.Context) (map[string]string, error) {
	tok, err := dr.driveServiceGoogleOA.UserOAuthToken(authCode)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, tok)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	result, err := d.ListFiles(10)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	return result, nil
}
