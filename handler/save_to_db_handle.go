package handler

import (
	"encoding/json"
	"github.com/eatmoreapple/openwechat"
	"web-wechat/db"
)

// 检查是否需要保存
func checkNeedSave(message *openwechat.Message) bool {
	return true
}

// 保存消息到MongoDB
func saveToDb(ctx *openwechat.MessageContext) {
	// TODO 需要解析成支持的结构体
	msgStr, _ := json.Marshal(ctx)
	db.SaveToMongo(string(msgStr), "message")
	ctx.Next()
}
