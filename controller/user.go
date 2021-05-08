package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"web-wechat/core"
	"web-wechat/global"
)

// GetCurrentUserInfoHandle 获取当前登录用户
func GetCurrentUserInfoHandle(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("AppKey")
	if len(appKey) < 1 {
		core.FailWithMessage("AppKey为必传参数", ctx)
		return
	}
	// 判断指定AppKey是不是有登录信息
	bot := global.GetBot(appKey)
	if nil == bot {
		core.FailWithMessage("未获取到登录记录", ctx)
		return
	}
	// 判断在线状态是否正常
	if !bot.Alive() {
		core.FailWithMessage("微信在线状态异常，请重新登录", ctx)
		return
	}
	// 获取登录用户信息
	user, err := bot.GetCurrentUser()
	if err != nil {
		core.FailWithMessage("获取登录用户信息失败", ctx)
		return
	}

	log.Println("登录用户：", user.NickName)

	core.OkWithData(user, ctx)
}
