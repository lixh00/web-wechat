package core

import (
	"gitee.ltd/lxh/logger/log"
	"github.com/gin-gonic/gin"
)

// CustomHTTPErrorHandler 默认全局异常处理
func CustomHTTPErrorHandler() gin.HandlerFunc {
	//ctx.Logger().Error(err.Error())
	log.Error("粗现异常啦~")
	//_ = FailWithMessage(err.Error(), ctx)
	return func(ctx *gin.Context) {
		FailWithMessage("不知道出了啥异常", ctx)
		ctx.Abort()
	}
}

// NotFoundErrorHandler 404异常处理
func NotFoundErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		FailWithMessage("错误路径", ctx)
		ctx.Abort()
	}
}
