package app

import (
	"context"
	"firstProject/internal/adapter/dynamodb"
	"firstProject/internal/adapter/google"
	serviceDrive "firstProject/internal/app/service/drive"
	serviceSample "firstProject/internal/app/service/sample"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Application struct {
	SampleService *serviceSample.SampleService
	DriveService  *serviceDrive.GoogleDriveService
	LineBotClient *linebot.Client
}

// NewApplication 建立並回傳一個新的 Application 實例
// dynamodb: 輸入實現了 dynamodb.DynamodbI interface的對象
func NewApplication(ctx context.Context, dynamodb dynamodb.DynamodbI, oauth google.GoogleOAuthI, lineBotClient *linebot.Client) *Application {

	app := &Application{
		LineBotClient: lineBotClient,
		SampleService: serviceSample.NewSampleService(ctx, serviceSample.SampleServiceParam{
			SampleServiceDynamodb: dynamodb,
		}),
		DriveService: serviceDrive.NewGoogleDriveService(ctx, serviceDrive.GoogleDriveServiceParam{
			DriveServiceGoogleOA: oauth,
			DriveServiceDynamodb: dynamodb,
		}),
	}
	return app
}
