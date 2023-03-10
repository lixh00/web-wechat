package controller

import (
	"fmt"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"unicode/utf8"
	"web-wechat/core"
	"web-wechat/global"
	"web-wechat/protocol"
)

// loginUrlResponse
// @description: 获取登录URL返回结构体
type loginUrlResponse struct {
	Uuid string `json:"uuid"`
	Url  string `json:"url"`
}

// GetLoginUrlHandle 获取登录扫码连接
func GetLoginUrlHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")

	// 获取一个微信机器人对象
	bot := global.InitWechatBotHandle()
	// 已扫码回调
	bot.ScanCallBack = func(body openwechat.CheckLoginResponse) {
		log.Infof("[%v]已扫码", appKey)
	}

	// 设置登录成功回调
	bot.LoginCallBack = func(body openwechat.CheckLoginResponse) {
		log.Infof("[%v]登录成功", appKey)
	}

	// 获取登录二维码链接
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	uuid, err := bot.Caller.GetLoginUUID()
	if err != nil {
		log.Errorf("获取登录二维码失败: %v", err.Error())
		core.FailWithMessage("获取登录二维码失败："+err.Error(), ctx)
		return
	}
	log.Infof("获取到uuid: %v", uuid)
	// 拼接URL
	url := fmt.Sprintf("https://login.weixin.qq.com/qrcode/%s", uuid)

	// 保存Bot到实例
	global.SetBot(appKey, bot)

	// 返回数据
	core.OkWithData(loginUrlResponse{Uuid: uuid, Url: url}, ctx)
}

// LoginHandle 登录
func LoginHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")
	uuid := ctx.Query("uuid")
	if utf8.RuneCountInString(uuid) < 1 {
		core.FailWithMessage("uuid为必传参数", ctx)
		return
	}
	//usePush := ctx.Query("usePush") // 是否使用免扫码登录
	//isPush := usePush == "1" || usePush == "true" || usePush == "yes"

	// 获取Bot对象
	bot := global.GetBot(appKey)
	if bot == nil {
		//bot = global.InitWechatBotHandle()
		//global.SetBot(appKey, bot)
		core.FailWithMessage("请先获取登录二维码", ctx)
		return
	}

	// 设置UUID
	bot.SetUUID(uuid)

	// 定义登录数据缓存
	storage := protocol.NewRedisHotReloadStorage("wechat:login:" + appKey)

	// 热登录
	var opts []openwechat.BotLoginOption
	opts = append(opts, openwechat.NewRetryLoginOption()) // 热登录失败使用扫码登录，适配第一次登录的时候无热登录数据
	//opts = append(opts, openwechat.NewSyncReloadDataLoginOption(10*time.Minute)) // 十分钟同步一次热登录数据

	// 登录
	if err := bot.HotLogin(storage, opts...); err != nil {
		log.Errorf("登录失败: %v", err)
		core.FailWithMessage("登录失败："+err.Error(), ctx)
		return
	}

	// 获取登录用户信息
	user, err := bot.GetCurrentUser()
	if err != nil {
		log.Errorf("获取登录用户信息失败: %v", err.Error())
		core.FailWithMessage("获取登录用户信息失败："+err.Error(), ctx)
		return
	}
	log.Infof("当前登录用户：%v", user.NickName)
	core.OkWithMessage("登录成功", ctx)
}
