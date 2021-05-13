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
	wechatMessageHandle(bot)
	// 返回机器人对象
	return bot
}

// 微信消息处理函数
func wechatMessageHandle(bot *protocol.Bot) {
	dispatcher := protocol.NewMessageMatchDispatcher()
	// 设置为异步处理
	dispatcher.SetAsync(true)

	// 处理文字消息
	textHandle := func(ctx *protocol.MessageContext) {
		sender, _ := ctx.Sender()
		log.Printf("[收到新文字消息] == 发信人：%v ==> 内容：%v \n", sender.NickName, ctx.Content)
		ctx.Next()
	}
	dispatcher.OnText(textHandle)

	// 处理其他消息
	otherHandle := func(ctx *protocol.MessageContext) {
		sender, _ := ctx.Sender()
		log.Printf("[收到新消息] 发送者：%v == 消息类型: %v ==> 内容：%v \n",
			sender.NickName, ctx.MsgType, protocol.XmlFormString(ctx.Content))
	}
	dispatcher.RegisterHandler(func(message *protocol.Message) bool {
		return !message.IsText()
	}, otherHandle)

	// 注册消息处理函数
	bot.MessageHandler = protocol.DispatchMessage(dispatcher)
}
