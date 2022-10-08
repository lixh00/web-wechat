package controller

import (
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
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

	// 获取一个微信机器人对象
	bot := global.InitWechatBotHandle()

	// 获取登录二维码链接
	url := "https://login.weixin.qq.com/qrcode/"
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
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

	// 已扫码回调
	bot.ScanCallBack = func(body []byte) {
		log.Infof("[%v]已扫码", appKey)
	}

	// 设置登录成功回调
	bot.LoginCallBack = func(body []byte) {
		log.Infof("[%v]登录成功", appKey)
	}

	// 热登录
	storage := protocol.NewRedisHotReloadStorage("wechat:login:" + appKey)
	if err := bot.HotLoginWithUUID(uuid, storage, true); err != nil {
		log.Errorf("热登录失败: %v", err)
		core.FailWithMessage("登录失败："+err.Error(), ctx)
		return
	}

	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Errorf("获取登录用户信息失败: %v", err.Error())
		core.FailWithMessage("获取登录用户信息失败："+err.Error(), ctx)
		return
	}
	log.Infof("当前登录用户：%v", user.NickName)
	core.OkWithMessage("登录成功", ctx)
}
