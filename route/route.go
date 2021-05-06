package route

import "github.com/gofiber/fiber/v2"

// InitRoute 初始化路由
func InitRoute(app *fiber.App) {
	initLoginRoute(app)
}
