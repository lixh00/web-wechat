package global

import "web-wechat/protocol"

var (
	// 登录用户的Bot对象
	wechatBots map[string]*protocol.Bot
)

// InitWechatBotsMap 初始化WechatBots
func InitWechatBotsMap() {
	wechatBots = make(map[string]*protocol.Bot)
}

// GetBot 获取Bot对象
func GetBot(uuid string) *protocol.Bot {
	return wechatBots[uuid]
}

// SetBot 保存Bot对象
func SetBot(uuid string, bot *protocol.Bot) {
	wechatBots[uuid] = bot
}
