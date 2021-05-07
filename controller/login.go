package controller

import (
	"github.com/labstack/echo/v4"
	"log"
	"web-wechat/core"
	"web-wechat/global"
	"web-wechat/protocol"
)

type loginResponse struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

// GetLoginUrlHandle 获取登录扫码连接
func GetLoginUrlHandle(ctx echo.Context) error {
	deviceId := ctx.QueryParam("deviceId")
	if len(deviceId) < 1 {
		return core.FailWithMessage("设备号必传", ctx)
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
	// 拼接URL
	url = url + *uuid

	// 保存Bot到实例
	global.SetBot(deviceId, bot)

	// 返回数据
	return core.OkWithData(loginResponse{Uuid: *uuid, Url: url}, ctx)
}

// LoginHandle 登录
func LoginHandle(ctx echo.Context) error {
	deviceId := ctx.QueryParam("deviceId")
	uuid := ctx.QueryParam("uuid")
	if len(deviceId) < 1 {
		return core.FailWithMessage("设备号必传", ctx)
	}
	bot := global.GetBot(deviceId)

	// 设置登录成功回调
	bot.LoginCallBack = func(body []byte) {
		log.Println("登录成功")
	}

	// 登录
	//hotLoginConfig := protocol.NewJsonFileHotReloadStorage("json文件路径，后面改成从Redis读取")
	//if err := bot.HotLoginWithUUID(uuid, hotLoginConfig, true); err != nil {
	if err := bot.LoginWithUUID(uuid); err != nil {
		log.Println(err)
		return err
	}

	// 阻塞主goroutine, 知道发生异常或者用户主动退出
	//err := bot.Block()
	//if err != nil {
	//	log.Println("Bot异常：", err.Error())
	//	return nil
	//}
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Println("获取登录用户信息失败")
		return err
	}
	log.Println("当前登录用户：", user.NickName, user.UserName)

	return core.OkWithMessage("登录成功", ctx)
}
