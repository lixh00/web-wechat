package route

import (
	"github.com/labstack/echo/v4"
	"web-wechat/controller"
)

// initUserRoute 初始化登录路由信息
func initUserRoute(app *echo.Echo) {
	group := app.Group("/user")

	group.GET("/info", controller.GetCurrentUserInfoHandle)
}
