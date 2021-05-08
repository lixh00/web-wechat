package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"web-wechat/core"
	"web-wechat/global"
	"web-wechat/protocol"
)

// 获取登录URL返回结构体
type loginUrlResponse struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

// GetLoginUrlHandle 获取登录扫码连接
func GetLoginUrlHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")
	if len(appKey) < 1 {
		core.FailWithMessage("AppKey为必传参数", ctx)
		return
	}

	// 获取一个微信机器人对象
	bot := global.InitWechatBotHandle()

	// 获取登录二维码链接
	url := "https://login.weixin.qq.com/qrcode/"
	bot.UUIDCallback = protocol.PrintlnQrcodeUrl
	uuid, err := bot.GetUUID()
	if err != nil {
		core.FailWithMessage("获取UUID失败", ctx)
		return
	}
	// 拼接URL
	url = url + *uuid

	// 保存Bot到实例
	global.SetBot(appKey, bot)

	// 返回数据
	core.OkWithData(loginUrlResponse{Uuid: *uuid, Url: url}, ctx)
}

// LoginHandle 登录
func LoginHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")
	uuid := ctx.Query("uuid")
	if len(appKey) < 1 {
		core.FailWithMessage("AppKey为必传参数", ctx)
		return
	}
	if len(uuid) < 1 {
		core.FailWithMessage("uuid为必传参数", ctx)
		return
	}
	// 获取Bot对象
	bot := global.GetBot(appKey)
	if bot == nil {
		bot = global.InitWechatBotHandle()
		global.SetBot(appKey, bot)
	}

	// 设置登录成功回调
	bot.LoginCallBack = func(body []byte) {
		log.Println("登录成功")
	}

	// 热登录
	storage := protocol.NewJsonFileHotReloadStorage("wechat:login:" + appKey)
	if err := bot.HotLoginWithUUID(uuid, storage, true); err != nil {
		log.Println(err)
		core.FailWithMessage("登录失败："+err.Error(), ctx)
		return
	}
	// 冷登录
	//if err := bot.LoginWithUUID(uuid); err != nil {
	//	log.Println(err)
	//	core.FailWithMessage("登录失败："+err.Error(), ctx)
	//	return
	//}
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Println("获取登录用户信息失败: ", err.Error())
		core.FailWithMessage("获取登录用户信息失败："+err.Error(), ctx)
		return
	}
	log.Println("当前登录用户：", user.NickName)
	core.OkWithMessage("登录成功", ctx)
}
