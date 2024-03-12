package firstProject

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	ginLambda *ginadapter.GinLambda
	ssmsvc    *SSM
)

func init() {
	ssmsvc = NewSSMClient()
	lineSecret, err := ssmsvc.Param("CHANNEL_SECRET", true).GetValue()
	if err != nil {
		log.Println(err)
	}
	lineAccessToken, err := ssmsvc.Param("CHANNEL_ACCESS_TOKEN", true).GetValue()
	if err != nil {
		log.Println(err)
	}
	bot, err := linebot.New(
		lineSecret,
		lineAccessToken,
	)
	if err != nil {
		log.Fatal(err)
	}
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/callback", func(c *gin.Context) {
		// ctx := c.Request.Context()
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if errors.Is(err, linebot.ErrInvalidSignature) {
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
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

type SSM struct {
	client ssmiface.SSMAPI
}

func Sessions() (*session.Session, error) {
	sess, err := session.NewSession()
	svc := session.Must(sess, err)
	return svc, err
}

func NewSSMClient() *SSM {
	// Create AWS Session
	sess, err := Sessions()
	if err != nil {
		log.Println(err)
		return nil
	}
	ssmsvc := &SSM{ssm.New(sess)}
	// Return SSM client
	return ssmsvc
}

type Param struct {
	Name           string
	WithDecryption bool
	ssmsvc         *SSM
}

func (s *SSM) Param(name string, decryption bool) *Param {
	return &Param{
		Name:           name,
		WithDecryption: decryption,
		ssmsvc:         s,
	}
}

func (p *Param) GetValue() (string, error) {
	ssmsvc := p.ssmsvc.client
	parameter, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &p.Name,
		WithDecryption: &p.WithDecryption,
	})
	if err != nil {
		return "", err
	}
	value := *parameter.Parameter.Value
	return value, nil
}
