package handler

import (
	"web-wechat/logger"
	"web-wechat/protocol"
)

// 设置消息为已读
func setTheMessageAsRead(ctx *protocol.MessageContext) {
	err := ctx.AsRead()
	if err != nil {
		logger.Log.Errorf("设置消息为已读出错: %v", err)
	}
	ctx.Next()
}
