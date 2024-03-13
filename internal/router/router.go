package router

import (
	"firstProject/internal/app"
	v1 "firstProject/internal/router/api/v1"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, app *app.Application) {
	registerAPIHandlers(router, app)
}

func registerAPIHandlers(router *gin.Engine, app *app.Application) {

	api := router.Group("/api")
	{
		v1.RegisterRouter(api, app)
	}

}
