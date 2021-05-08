package route

import (
	"github.com/gin-gonic/gin"
	"web-wechat/controller"
)

// initUserRoute 初始化登录路由信息
func initUserRoute(app *gin.Engine) {
	group := app.Group("/user")

	group.GET("/info", controller.GetCurrentUserInfoHandle)
}
