package drive

import (
	"context"
	"errors"
	domainDrive "firstProject/internal/domain/drive"
	"golang.org/x/oauth2"
	"io"
	"log"
	"strings"
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

func (dr *GoogleDriveService) TestFolderCarousel(ctx context.Context, lineID string) (*domainDrive.FolderCarousel, error) {
	insideFolderM := map[string]string{
		"001": "F1",
		"002": "F2",
	}
	fileM := map[string]string{
		"001": "file1",
		"002": "file2",
	}

	var params domainDrive.NewFolderCarouselParam
	params.BubbleParams = append(params.BubbleParams,
		domainDrive.NewFolderBubbleParam{
			Type:          "我的雲端硬碟",
			Name:          "Folder1",
			Path:          "/xx/xx",
			ID:            "123",
			InsideFolderM: insideFolderM,
			FileM:         fileM,
		},
		domainDrive.NewFolderBubbleParam{
			Type:          "我的雲端硬碟",
			Name:          "Folder2",
			Path:          "/yy/yy",
			ID:            "1234",
			InsideFolderM: insideFolderM,
			FileM:         fileM,
		},
	)
	carousel := domainDrive.NewFolderCarousel(params)
	return &carousel, nil
}

func (dr *GoogleDriveService) ListFolderCarousel(ctx context.Context, lineID string, folderType domainDrive.FolderType) (*domainDrive.FolderCarousel, error) {
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

	var folderList map[string]string
	var folderTypeString string
	switch folderType {
	case domainDrive.PersonalFolder:
		folderList, err = d.ListMyDriveFolders()
		folderTypeString = "我的雲端硬碟"
	case domainDrive.SharedFolder:
		folderList, err = d.ListSharedFolders()
		folderTypeString = "與我共用"
	default:
		return nil, errors.New("unsupported folder type")
	}
	if err != nil {
		return nil, err
	}

	var params domainDrive.NewFolderCarouselParam

	for folderID, name := range folderList {
		path, err := d.FindFolderPathByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		insideFolderM, err := d.ListFolderByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fileM, err := d.ListFilesByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		param := domainDrive.NewFolderBubbleParam{
			Type:          folderTypeString,
			Name:          name,
			Path:          path,
			ID:            folderID,
			InsideFolderM: insideFolderM,
			FileM:         fileM,
		}
		params.BubbleParams = append(params.BubbleParams, param)
	}

	carousel := domainDrive.NewFolderCarousel(params)

	return &carousel, err
}

func (dr *GoogleDriveService) ListSelectedFolderCarousel(ctx context.Context, lineID string, folderID string) (*domainDrive.FolderCarousel, error) {
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

	var params domainDrive.NewFolderCarouselParam

	// 當前目錄下的自己，才能看到資料夾下的檔案
	folderList, err := d.ListFolderByID(folderID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	currentPath, err := d.FindFolderPathByID(folderID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	currentFile, err := d.ListFilesByID(folderID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	segments := strings.Split(currentPath, "/")
	currentName := segments[len(segments)-2]

	params.BubbleParams = append(params.BubbleParams, domainDrive.NewFolderBubbleParam{
		Type:          "打開資料夾",
		Name:          currentName,
		Path:          currentPath,
		ID:            folderID,
		InsideFolderM: folderList,
		FileM:         currentFile,
	})

	for folderID, name := range folderList {
		path, err := d.FindFolderPathByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		insideFolderM, err := d.ListFolderByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fileM, err := d.ListFilesByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		param := domainDrive.NewFolderBubbleParam{
			Type:          "子資料夾",
			Name:          name,
			Path:          path,
			ID:            folderID,
			InsideFolderM: insideFolderM,
			FileM:         fileM,
		}
		params.BubbleParams = append(params.BubbleParams, param)
	}

	carousel := domainDrive.NewFolderCarousel(params)

	return &carousel, err
}

func (dr *GoogleDriveService) SetUploadPath(ctx context.Context, lineID string, folderID string) error {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return err
	}
	dToken.PK = lineID
	dToken.Info = map[string]interface{}{
		"upload_folder_id": folderID,
	}

	_, err = dr.driveServiceDynamodb.TxUpdateGoogleOAuthToken(dToken)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (dr *GoogleDriveService) GetUploadPath(ctx context.Context, lineID string) (string, error) {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	folderID := dToken.Info["upload_folder_id"]

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}

	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return "", err
	}
	path, err := d.FindFolderPathByID(folderID.(string))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return path, nil
}
