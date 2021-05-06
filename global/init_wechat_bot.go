package global

import (
	"fmt"
	"github.com/eatMoreApple/openwechat"
)

func InitWeChatBot() {
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		// TODO 保存消息到数据库
		fmt.Println(msg)
	}

	WeChatBot = bot
}