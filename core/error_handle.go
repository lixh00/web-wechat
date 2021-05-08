package core

import (
	"github.com/labstack/echo/v4"
	"log"
)

// CustomHTTPErrorHandler 默认全局异常处理
func CustomHTTPErrorHandler(err error, ctx echo.Context) {
	//ctx.Logger().Error(err.Error())
	log.Println("粗现异常啦~", err)
	_ = FailWithMessage(err.Error(), ctx)
}
