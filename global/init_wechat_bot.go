package global

import (
	"fmt"
	"web-wechat/protocol"
)

// InitWechatBotHandle 初始化微信机器人
func InitWechatBotHandle() *protocol.Bot {
	bot := protocol.DefaultBot(protocol.Desktop)

	// 注册消息处理函数
	bot.MessageHandler = func(msg *protocol.Message) {
		// TODO 保存消息到数据库
		fmt.Println(msg)
	}

	return bot
}
