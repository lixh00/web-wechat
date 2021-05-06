package global

import (
	"fmt"
	"log"
	"web-wechat/protocol"
)

// InitWechatBotHandle 初始化微信机器人
func InitWechatBotHandle() *protocol.Bot {
	bot := protocol.DefaultBot(protocol.Desktop)

	// 定义读取消息错误回调函数
	bot.GetMessageErrorHandler = func(err error) {
		log.Println("获取消息发生错误：", err.Error())
	}
	// 注册消息处理函数
	bot.MessageHandler = func(msg *protocol.Message) {
		fmt.Println(msg)
		if msg.MsgType == 1 {
			_, err := msg.ReplyText("您说：" + msg.Content)
			if err != nil {
				log.Println("回复消息发生错误: ", err.Error())
			}
		}
		// TODO 保存消息到数据库
		//log.Fatalf("[收到新消息] 消息ID：%v, 消息类型：%v, 发件人：%v, 收件人：%v, 正文：%v",
		//	msg.MsgId,
		//	msg.MsgType,
		//	msg.FromUserName,
		//	msg.ToUserName,
		//	msg.Content)
	}

	return bot
}
