package controller

import (
	"errors"
	"web-wechat/global"
)

// 预检AppKey是否存在登录记录且登录状态是否正常
func checkBot(appKey string) error {
	if len(appKey) < 1 {
		return errors.New("AppKey为必传参数")
	}
	// 判断指定AppKey是不是有登录信息
	bot := global.GetBot(appKey)
	if nil == bot {
		return errors.New("未获取到登录记录")
	}
	// 判断在线状态是否正常
	if !bot.Alive() {
		return errors.New("微信在线状态异常，请重新登录")
	}
	return nil
}
