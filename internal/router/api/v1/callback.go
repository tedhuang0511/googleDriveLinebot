package v1

import (
	"firstProject/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
)

func Callback(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		events, err := app.LineBotClient.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Println(err)
				c.JSON(http.StatusBadRequest, err)
			} else {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, err)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if message.Text == "login" {
						lineID := event.Source.UserID
						authURL := app.DriveService.LoginURL(ctx, lineID)
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(authURL)).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					samplePK, err := app.SampleService.Sample(ctx, message.Text)
					if err != nil {
						log.Println(err)
						return
					}
					if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(samplePK)).Do(); err != nil {
						log.Println(err)
					}
				}
			}
		}

	}

}
