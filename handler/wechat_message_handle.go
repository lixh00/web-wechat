package handler

import "web-wechat/protocol"

func HandleMessage(bot *protocol.Bot) {
	// 定义一个处理器
	dispatcher := protocol.NewMessageMatchDispatcher()
	// 设置为异步处理
	dispatcher.SetAsync(true)

	// 注册文本消息处理函数
	dispatcher.OnText(textMessageHandle)
	// 注册图片消息处理器
	dispatcher.OnImage(imageMessageHandle)
	// 未定义消息处理
	dispatcher.RegisterHandler(checkIsOther, otherMessageHandle)

	// 注册消息处理器
	bot.MessageHandler = protocol.DispatchMessage(dispatcher)
}
