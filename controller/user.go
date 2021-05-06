package controller

import (
	"github.com/gofiber/fiber/v2"
	"web-wechat/global"
)

// GetCurrentUserInfoHandle 获取当前登录用户
func GetCurrentUserInfoHandle(ctx *fiber.Ctx) error {
	deviceId := ctx.Query("deviceId")
	bot := global.GetBot(deviceId)

	user, err := bot.GetCurrentUser()
	if err != nil {
		return ctx.SendString("用戶登錄失敗")
	}

	return ctx.SendString(user.UserName)
}
