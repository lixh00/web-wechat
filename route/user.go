package route

import (
	"github.com/gin-gonic/gin"
	"web-wechat/controller"
)

// initUserRoute 初始化登录路由信息
func initUserRoute(app *gin.Engine) {
	group := app.Group("/user")

	// 获取登录的用户信息
	group.GET("/info", controller.GetCurrentUserInfoHandle)
	// 获取好友列表
	group.GET("/friends", controller.GetFriendsListHandle)
}
