package global

import (
	"fmt"
	"web-wechat/protocol"
)

// InitWeChatBot 初始化微信机器人
func InitWeChatBot() {
	bot := protocol.DefaultBot(protocol.Desktop)

	// 注册消息处理函数
	bot.MessageHandler = func(msg *protocol.Message) {
		// TODO 保存消息到数据库
		fmt.Println(msg)
	}

	WeChatBot = bot
}
