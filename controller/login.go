package controller

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"web-wechat/core"
	"web-wechat/global"
	"web-wechat/protocol"
)

// GetLoginUrlHandle 获取登录扫码连接
func GetLoginUrlHandle(ctx *fiber.Ctx) error {
	deviceId := ctx.Query("deviceId")
	if len(deviceId) < 1 {
		return ctx.SendString("设备号必传")
	}
	log.Println("收到登录请求")

	// 获取一个微信机器人对象
	bot := global.InitWechatBotHandle()

	// 获取登录二维码链接
	url := "https://login.weixin.qq.com/qrcode/"
	bot.UUIDCallback = protocol.PrintlnQrcodeUrl
	uuid, err := bot.GetUUID()
	if err != nil {
		return core.FailWithMessage("获取UUID失败", ctx)
	}
	url = url + *uuid

	// 保存Bot到实例
	global.SetBot(deviceId, bot)

	// 返回数据
	return ctx.SendString(url)
}

// LoginHandle 登录
func LoginHandle(ctx *fiber.Ctx) error {
	deviceId := ctx.Query("deviceId")
	if len(deviceId) < 1 {
		return ctx.SendString("设备号必传")
	}
	bot := global.GetBot(deviceId)

	// 设置登录成功回调
	bot.LoginCallBack = func(body []byte) {
		log.Println(string(body))
	}

	// 登录
	if err := bot.Login(); err != nil {
		log.Println(err)
		return err
	}

	return ctx.SendString("登录成功")
}
