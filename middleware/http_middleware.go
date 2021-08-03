package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"web-wechat/core"
	"web-wechat/global"
)

// CheckAppKeyIsLoggedInMiddleware 检查AppKey是否已登录微信
func CheckAppKeyIsLoggedInMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appKey := ctx.Request.Header.Get("AppKey")
		// TODO 从数据库判断AppKey是否存在
		// 如果不是登录请求，判断AppKey是否有效
		flag := true
		if !strings.Contains(ctx.Request.RequestURI, "login") {
			if err := global.CheckBot(appKey); err != nil {
				core.FailWithMessage("AppKey预检失败："+err.Error(), ctx)
				flag = false
			}
		}
		if flag {
			ctx.Next()
		} else {
			ctx.Abort()
		}
	}
}

// CheckAppKeyExistMiddleware 检查是否有appKey
func CheckAppKeyExistMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appKey := ctx.Request.Header.Get("AppKey")
		// 先判断AppKey是不是传了
		if len(appKey) < 1 {
			core.FailWithMessage("AppKey为必传参数", ctx)
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
