package v1

import (
	"firstProject/internal/app"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.RouterGroup, app *app.Application) {

	v1 := router.Group("/v1")
	{
		// 將 /v1/callback 請求路由到 Callback 函數來處理
		v1.POST("/callback", Callback(app))
		v1.GET("/oauth-login", OAuthLogin(app))

	}
}
