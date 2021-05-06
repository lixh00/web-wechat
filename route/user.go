package route

import (
	"github.com/gofiber/fiber/v2"
	"web-wechat/controller"
)

// initUserRoute 初始化登录路由信息
func initUserRoute(app *fiber.App) {
	group := app.Group("/user")

	group.Get("/info", controller.GetCurrentUserInfoHandle)
}
