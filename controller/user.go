package controller

import (
	"github.com/labstack/echo/v4"
	"web-wechat/core"
	"web-wechat/global"
)

// GetCurrentUserInfoHandle 获取当前登录用户
func GetCurrentUserInfoHandle(ctx echo.Context) error {
	deviceId := ctx.QueryParam("deviceId")
	bot := global.GetBot(deviceId)
	if nil == bot {
		return core.FailWithMessage("设备ID无登录记录", ctx)
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		return core.FailWithMessage("用户登录失败", ctx)
	}

	return core.OkWithData(user, ctx)
}
