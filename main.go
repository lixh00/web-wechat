package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"web-wechat/global"
	"web-wechat/route"
)

// 程序启动入口
func main() {
	app := echo.New()

	// 使用日志中间件
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} [${level}] [${remote_ip}] ${method} [${status}] ${uri}` + "\n",
	}))
	app.Use(middleware.Recover())

	// 初始化路由
	route.InitRoute(app)

	// 初始化WechatBotMap
	global.InitWechatBotsMap()

	// TODO 初始化数据库连接等

	// 监听端口
	app.Logger.Fatal(app.Start(":8888"))
}
