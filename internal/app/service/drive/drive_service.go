package drive

import (
	"context"
	domainDrive "firstProject/internal/domain/drive"
	"golang.org/x/oauth2"
	"io"
	"log"
)

func (dr *GoogleDriveService) ListFiles(ctx context.Context, lineID string) (map[string]string, error) {
	// token改成去db取
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// 把token轉成oauth2的格式
	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := d.ListFiles(10)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (dr *GoogleDriveService) ListMyDriveFolders(ctx context.Context, lineID string) (map[string]string, error) {

	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := d.ListMyDriveFolders()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (dr *GoogleDriveService) ListSharedFolders(ctx context.Context, lineID string) (map[string]string, error) {

	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := d.ListSharedFolders()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (dr *GoogleDriveService) UploadFile(ctx context.Context, lineID string, fileName string, content io.ReadCloser) error {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)

	if err != nil {
		log.Println(err)
		return err
	}

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}

	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return err
	}

	file, err := domainDrive.SaveContent(content)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("START Upload File To Drive")

	folderID := dToken.Info["upload_folder_id"].(string)
	// 假設預設的儲存路徑
	folderID = "1a17lQuvZCjPcBj_UoNryg0tdLr7lb1islzFtWNOZnqTxzuW6Am0nJ67HzxdzPBsp4gP1jPmQ"

	err = d.UploadFile(folderID, fileName, file)
	if err != nil {
		log.Println("err:", err)
		return err
	}
	return nil

}
