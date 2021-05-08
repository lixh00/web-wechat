package global

import (
	"log"
	"web-wechat/protocol"
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

// InitWechatBotHandle 初始化微信机器人
func InitWechatBotHandle() *protocol.Bot {
	bot := protocol.DefaultBot(protocol.Desktop)

	// 定义读取消息错误回调函数
	bot.GetMessageErrorHandler = func(err error) {
		log.Println("获取消息发生错误：", err.Error())
	}
	// 注册消息处理函数
	bot.MessageHandler = func(msg *protocol.Message) {
		// TODO 保存消息到数据库
		// 取出发送者信息
		sender, err := msg.Sender()
		if err != nil {
			log.Println("获取消息发送者失败", err.Error())
		} else {
			log.Println("消息发送者：", sender.NickName)
		}
		if msg.IsText() {
			log.Printf("[收到新消息] 内容：%v \n", msg.Content)
		}
	}

	return bot
}
