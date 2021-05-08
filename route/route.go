package route

import (
	"github.com/gin-gonic/gin"
)

// InitRoute 初始化路由
func InitRoute(app *gin.Engine) {
	// 初始化登录相关路由
	initLoginRoute(app)

	// 初始化用户相关路由
	initUserRoute(app)
}
