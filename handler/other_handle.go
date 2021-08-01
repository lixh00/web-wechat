package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"web-wechat/logger"
)

func checkIsOther(message *openwechat.Message) bool {
	// 处理除文字消息和通知消息之外，并且不是自己发送的消息
	return !message.IsText() && !message.IsNotify() && !message.IsPicture() && !message.IsEmoticon() && !message.IsMedia() //  && !message.IsSendBySelf()
}

// 未定义消息处理
func otherMessageHandle(ctx *openwechat.MessageContext) {
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, err := ctx.SenderInGroup()
		if err != nil {
			logger.Log.Errorf("获取发信人失败: %v", err)
			return
		}
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}
	logger.Log.Info("========================================================================================")
	logger.Log.Infof("收到未定义消息\n消息类型: %v\n发信人: %v\n内容: %v", ctx.MsgType, senderUser,
		openwechat.XmlFormString(ctx.Content))
	logger.Log.Info("========================================================================================")
	ctx.Next()
}
