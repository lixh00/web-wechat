package main

import (
	"github.com/gin-gonic/gin"
	"web-wechat/global"
	"web-wechat/route"
)

// 程序启动入口
func main() {
	app := gin.Default()

	// 定义全局异常处理
	//app.Use(core.CustomHTTPErrorHandler)
	// 初始化路由
	route.InitRoute(app)

	// 初始化WechatBotMap
	global.InitWechatBotsMap()

	// TODO 初始化数据库连接等

	// 监听端口
	app.Run(":8888")
}
