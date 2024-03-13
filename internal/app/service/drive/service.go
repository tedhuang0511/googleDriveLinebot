package drive

import "context"

type GoogleDriveService struct {
	driveServiceGoogleOA DriveServiceGoogleOAuthI
}

type GoogleDriveServiceParam struct {
	DriveServiceGoogleOA DriveServiceGoogleOAuthI
}

func NewGoogleDriveService(_ context.Context, param GoogleDriveServiceParam) *GoogleDriveService {
	return &GoogleDriveService{
		driveServiceGoogleOA: param.DriveServiceGoogleOA,
	}
}
