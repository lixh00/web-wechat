package route

import (
	"github.com/gin-gonic/gin"
	"web-wechat/controller"
)

// 初始化消息相关路由
func initMessageRoute(app *gin.Engine) {
	group := app.Group("/message")

	// 向指定好友发送消息
	group.PUT("/user", controller.SendMessageToUser)

	// 向指定群组发送消息
	group.PUT("/group", controller.SendMessageToGroup)
}
