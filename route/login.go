package route

import (
	"github.com/gofiber/fiber/v2"
	"web-wechat/controller"
)

// initLoginRoute 初始化登录路由信息
func initLoginRoute(app *fiber.App) {
	// 获取登录二维码
	app.Get("/login", controller.GetLoginUrlHandle)
	// 检查登录状态
	app.Post("/login", controller.LoginHandle)
}
