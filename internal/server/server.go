package server

import (
	"context"
	"firstProject/internal/adapter/dynamodb"
	"firstProject/internal/adapter/ssm"
	"firstProject/internal/app"
	"firstProject/internal/router"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func initRouter(rootCtx context.Context, app *app.Application) (ginRouter *gin.Engine) {

	// create *gin.Engine
	ginRouter = gin.New()

	// RegisterHandlers
	router.RegisterHandlers(ginRouter, app)

	return ginRouter
}

func NewGinLambda() *ginadapter.GinLambda {
	rootCtx, _ := context.WithCancel(context.Background()) //nolint
	ssmsvc := ssm.NewSSM()

	lineSecret, err := ssmsvc.FindParameter(rootCtx, ssmsvc.Client, "CHANNEL_SECRET")
	if err != nil {
		log.Println(err)
	}

	lineAccessToken, err := ssmsvc.FindParameter(rootCtx, ssmsvc.Client, "CHANNEL_ACCESS_TOKEN")
	if err != nil {
		log.Println(err)
	}
	lineClientLambda, err := linebot.New(lineSecret, lineAccessToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("LineBot Create Success")

	db := dynamodb.NewTableBasics("google-oauth")

	app := app.NewApplication(rootCtx, db, lineClientLambda)
	ginRouter := initRouter(rootCtx, app)
	return ginadapter.New(ginRouter)
}

func StartNgrokServer() {
	// 初始化root上下文和取消函數
	rootCtx, rootCtxCancelFunc := context.WithCancel(context.Background())
	// 使用sync.WaitGroup等待所有goroutine的完成
	wg := sync.WaitGroup{}
	// 初始化LineBot客戶端
	lineClient, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// 初始化DynamoDB連接，然後切換到本地DynamoDB
	db := dynamodb.NewTableBasics("google-oauth")
	db.DynamoDbClient = dynamodb.CreateLocalClient(8000)
	// 初始化Application
	app := app.NewApplication(rootCtx, db, lineClient)
	// 初始化Gin路由
	ginRouter := initRouter(rootCtx, app)

	// 啟動 ngrok
	wg.Add(1)
	runNgrokServer(rootCtx, &wg, ginRouter)

	// 監聽SIGTERM/SIGINT信號來進行優雅的關閉
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	// 阻塞，當收到信號時執行下面的code觸發Ngrok的關閉
	<-gracefulStop
	rootCtxCancelFunc()

	// 使用goroutine等待所有服務的結束，等待最多10s
	var waitUntilDone = make(chan struct{})
	go func() {
		wg.Wait()
		close(waitUntilDone)
	}()
	select {
	case <-waitUntilDone:
		log.Println("success to close all services")
	case <-time.After(10 * time.Second):
		log.Println(context.DeadlineExceeded, "fail to close all services")
	}

}

func runNgrokServer(rootCtx context.Context, wg *sync.WaitGroup, ginRouter *gin.Engine) {

	tun, err := ngrok.Listen(rootCtx,
		config.HTTPEndpoint(config.WithDomain(os.Getenv("NGROK_DOMAIN"))),
		ngrok.WithAuthtokenFromEnv(),
	)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Application available at:", tun.URL())

	go func() {
		err = http.Serve(tun, ginRouter)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		<-rootCtx.Done()
		log.Println("Shutting down ngrok server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tun.CloseWithContext(ctx); err != nil {
			log.Printf("Error closing ngrok tunnel: %v\n", err)
		}
		log.Println("ngrok server gracefully stopped")
		wg.Done()
	}()

}
