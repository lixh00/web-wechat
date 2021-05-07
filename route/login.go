package route

import (
	"github.com/labstack/echo/v4"
	"web-wechat/controller"
)

// initLoginRoute 初始化登录路由信息
func initLoginRoute(app *echo.Echo) {
	// 获取登录二维码
	app.GET("/login", controller.GetLoginUrlHandle)
	// 检查登录状态
	app.POST("/login", controller.LoginHandle)
}
