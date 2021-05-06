package route

import (
	"github.com/gofiber/fiber/v2"
	"web-wechat/controller"
)

// initLoginRoute 初始化登录路由信息
func initLoginRoute(app *fiber.App) {
	app.Get("/login", controller.GetLoginUrlHandle)
}