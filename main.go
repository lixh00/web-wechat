package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"web-wechat/core"
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
		// 全局错误处理
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(core.Response{
					Code: core.ERROR,
					Msg:  "系统错误",
				})
		},
	})

	// 使用日志中间件 - 使用默认配置
	//app.Use(logger.New())
	// 输出日志到文件
	file, err := os.OpenFile("./logs/run.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	// 配置日志中间件
	app.Use(logger.New(logger.Config{
		Output: file,
	}))

	// 初始化路由
	route.InitRoute(app)

	// 初始化WechatBotMap
	global.InitWechatBotsMap()

	// TODO 初始化数据库连接等

	// 监听端口
	log.Fatal(app.Listen(":8888"))
}
