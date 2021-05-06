package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"web-wechat/global"
)

// GetLoginUrlHandle 获取登录扫码连接
func GetLoginUrlHandle(ctx *fiber.Ctx) error {
	//fmt.Println("收到登录请求")
	log.Println("收到登录请求")
	// 获取登录二维码链接
	var url = "https://login.weixin.qq.com/qrcode/"
	global.WeChatBot.UUIDCallback = func(uuid string) {
		fmt.Println("UUID: " + uuid)
		url += uuid
	}
	return ctx.SendString(url)
}