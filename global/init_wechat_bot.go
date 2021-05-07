package global

import (
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
		//fmt.Println(msg)

		//m, err := json.Marshal(msg)
		//if err != nil {
		//	log.Println("消息转换失败：", err.Error())
		//}
		//log.Println(m)

		//if msg.MsgType == 1 {
		//	_, err := msg.ReplyText("您说：" + msg.Content)
		//	if err != nil {
		//		log.Println("回复消息发生错误: ", err.Error())
		//	}
		//}
		// TODO 保存消息到数据库

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
