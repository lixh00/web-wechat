package handler

import (
	"fmt"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
)

// 处理文本消息
func textMessageHandle(ctx *openwechat.MessageContext) {
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}

	log.Infof("[收到新文字消息] == 发信人：%v ==> 内容：%v", senderUser, ctx.Content)
	//if !ctx.IsSendBySelf() {
	//	sender, _ := ctx.Sender()
	//	if ctx.IsSendByGroup() {
	//		// 取出消息在群里面的发送者
	//		senderInGroup, _ := ctx.SenderInGroup()
	//		log.Infof("[群聊][收到新文字消息] == 发信人：%v[%v] ==> 内容：%v", sender.NickName,
	//			senderInGroup.NickName, ctx.Content)
	//	} else {
	//		log.Infof("[好友][收到新文字消息] == 发信人：%v ==> 内容：%v", sender.NickName, ctx.Content)
	//	}
	//}
	ctx.Next()
}
