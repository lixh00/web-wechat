package controller

import (
	"github.com/labstack/echo/v4"
	"web-wechat/core"
	"web-wechat/global"
)

// GetCurrentUserInfoHandle 获取当前登录用户
func GetCurrentUserInfoHandle(ctx echo.Context) error {
	appKey := ctx.Request().Header.Get("AppKey")
	if len(appKey) < 1 {
		return core.FailWithMessage("AppKey为必传参数", ctx)
	}
	bot := global.GetBot(appKey)
	if nil == bot {
		return core.FailWithMessage("未获取到登录记录", ctx)
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		return core.FailWithMessage("用户登录失败", ctx)
	}

	return core.OkWithData(user, ctx)
}
