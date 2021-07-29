package handler

import (
	"fmt"
	"web-wechat/logger"
	"web-wechat/protocol"
)

// 处理文本消息
func textMessageHandle(ctx *protocol.MessageContext) {
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}

	logger.Log.Infof("[收到新文字消息] == 发信人：%v ==> 内容：%v", senderUser, ctx.Content)
	//if !ctx.IsSendBySelf() {
	//	sender, _ := ctx.Sender()
	//	if ctx.IsSendByGroup() {
	//		// 取出消息在群里面的发送者
	//		senderInGroup, _ := ctx.SenderInGroup()
	//		logger.Log.Infof("[群聊][收到新文字消息] == 发信人：%v[%v] ==> 内容：%v", sender.NickName,
	//			senderInGroup.NickName, ctx.Content)
	//	} else {
	//		logger.Log.Infof("[好友][收到新文字消息] == 发信人：%v ==> 内容：%v", sender.NickName, ctx.Content)
	//	}
	//}
	ctx.Next()
}
