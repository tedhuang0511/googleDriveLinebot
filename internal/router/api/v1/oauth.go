package v1

import (
	"context"
	"firstProject/internal/app"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3" //go get google.golang.org/api/drive/v3
	"google.golang.org/api/option"
	"log"
	"os"
)

func OAuthLogin(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 URL 參數中獲取授權碼
		authCode := c.Query("code")
		// 建立 OAuth2 Config
		config := &oauth2.Config{
			ClientID:     os.Getenv("ClientID"),
			ClientSecret: os.Getenv("ClientSecret"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
			Endpoint:     google.Endpoint,
			Scopes:       []string{drive.DriveScope},
			RedirectURL:  os.Getenv("RedirectURL"),
		}
		// 交換授權碼以獲取token
		tok, err := config.Exchange(context.TODO(), authCode)
		if err != nil {
			log.Printf("Unable to retrieve token from web %v", err)
		}
		// 使用token初始化 Google Drive 服務
		client := config.Client(context.Background(), tok)
		srv, err := drive.NewService(c, option.WithHTTPClient(client))
		if err != nil {
			log.Printf("Unable to retrieve Drive client: %v", err)
		}
		// 列出 Google Drive 上的文件
		r, err := srv.Files.List().PageSize(10).
			Fields("nextPageToken, files(id, name)").Do()
		if err != nil {
			log.Printf("Unable to retrieve files: %v", err)
		}
		fmt.Println("Files:")
		if len(r.Files) == 0 {
			fmt.Println("No files found.")
		} else {
			for _, i := range r.Files {
				fmt.Printf("%s (%s)\n", i.Name, i.Id)
			}
		}
		// 將 HTML 寫入回應，顯示授權成功的消息
		_, err = c.Writer.Write([]byte("<html><title>Login</title> <body> Authorized successfully, please close this window</body></html>"))
		if err != nil {
			log.Printf("Unable to write HTML: %v", err)
		}

	}
}
