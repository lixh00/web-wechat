package global

import (
	"errors"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"web-wechat/handler"
)

// InitWechatBotsMap 初始化WechatBots
func InitWechatBotsMap() {
	wechatBots = make(map[string]*openwechat.Bot)
}

// GetBot 获取Bot对象
func GetBot(appKey string) *openwechat.Bot {
	return wechatBots[appKey]
}

// SetBot 保存Bot对象
func SetBot(appKey string, bot *openwechat.Bot) {
	wechatBots[appKey] = bot
}

// CheckBot 预检AppKey是否存在登录记录且登录状态是否正常
func CheckBot(appKey string) error {
	// 判断指定AppKey是不是有登录信息
	bot := GetBot(appKey)
	if nil == bot {
		return errors.New("未获取到登录记录")
	}
	// 判断在线状态是否正常
	if !bot.Alive() {
		return errors.New("微信在线状态异常，请重新登录")
	}
	return nil
}

// InitWechatBotHandle 初始化微信机器人
func InitWechatBotHandle() *openwechat.Bot {
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 定义读取消息错误回调函数
	//var getMessageErrorCount int32
	//bot.GetMessageErrorHandler = func(err error) {
	//	atomic.AddInt32(&getMessageErrorCount, 1)
	//	// 如果发生了三次错误,那么直接退出
	//	if getMessageErrorCount == 3 {
	//		log.Errorf("获取消息发生错误达到三次，直接退出。错误信息：%v", err.Error())
	//		_ = bot.Logout()
	//	}
	//}

	// 设置心跳回调
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		if resp.RetCode == "1100" {
			log.Errorf("微信已退出")
			// do something
		}
		switch resp.Selector {
		case "0":
			log.Debugf("正常")
		case "2", "6":
			log.Debugf("有新消息")
		case "7":
			log.Debugf("进入/离开聊天界面")
			err := bot.WebInit()
			if err != nil {
				// 短信通知一下
				// do something
				//log.Panicf("重新初始化失败: %v", err)
			}
		default:
			log.Debugf("RetCode: %s  Selector: %s", resp.RetCode, resp.Selector)
		}
	}

	// 注册消息处理函数
	handler.HandleMessage(bot)
	// 获取消息发生错误
	//bot.MessageOnError()
	// 返回机器人对象
	return bot
}
