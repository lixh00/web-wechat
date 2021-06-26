package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"web-wechat/core"
	"web-wechat/global"
)

// CheckAppKeyMiddleware 检查是否有appKey
func CheckAppKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appKey := ctx.Request.Header.Get("AppKey")
		// 先判断AppKey是不是传了
		if len(appKey) < 1 {
			core.FailWithMessage("AppKey为必传参数", ctx)
			ctx.Abort()
		}
		// 如果不是登录请求，判断AppKey是否有效
		if !strings.Contains(ctx.Request.RequestURI, "login") {
			if err := global.CheckBot(appKey); err != nil {
				core.FailWithMessage("AppKey预检失败："+err.Error(), ctx)
				ctx.Abort()
			}
		}
		ctx.Next()
	}
}
