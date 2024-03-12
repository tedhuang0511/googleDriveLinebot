package main

import (
	"context"
	"firstProject/internal/adapter/ssm"
	"firstProject/internal/server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

var (
	ginLambda        *ginadapter.GinLambda
	lineClientLambda *linebot.Client
	ssmsvc           *ssm.SSM
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {
	// env GIN_MODE="release"
	if gin.Mode() == gin.ReleaseMode {
		log.Println("Run on Lambda")
		ginLambda = server.NewGinLambda()
		lambda.Start(Handler)
	} else if gin.Mode() == gin.DebugMode {
		log.Println("Debug mode run on local")
		server.StartNgrokServer()
	}
}
