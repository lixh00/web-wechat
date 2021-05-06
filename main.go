package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"web-wechat/global"
	"web-wechat/route"
)

// 程序启动入口
func main() {
	app := fiber.New(fiber.Config{
		// 服务器Header
		ServerHeader: "wechat",
		// 接口地址是否区分大小写
		CaseSensitive: true,
	})

	// 使用日志中间件 - 使用默认配置
	app.Use(logger.New())

	// 初始化微信机器人插件
	global.InitWeChatBot()

	// 初始化路由
	route.InitRoute(app)

	// TODO 初始化数据库连接等

	// 监听端口
	log.Fatal(app.Listen(":8888"))
}
