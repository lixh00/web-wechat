package handler

import (
	"encoding/json"
	"github.com/eatmoreapple/openwechat"
	"time"
	. "web-wechat/db"
)

// 检查是否需要保存
func checkNeedSave(message *openwechat.Message) bool {
	return message.IsText() || message.IsEmoticon() || message.IsPicture() || message.IsMedia()
}

// 保存消息到MongoDB
func saveToDb(ctx *openwechat.MessageContext) {
	// TODO 需要解析成支持的结构体

	type weChatMsg struct {
		Uin          int
		MsgId        string
		MsgType      openwechat.MessageType
		Content      string
		SendUserName string
		GroupName    string
		IsRead       int
		BaseStr      string
		DateTime     string
	}

	slew, _ := ctx.Bot.GetCurrentUser()
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	groupName := ""
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = senderInGroup.NickName
		groupName = sender.NickName
	}
	msgStr, _ := json.Marshal(ctx)
	msg := weChatMsg{
		Uin:          slew.Uin,
		MsgId:        ctx.MsgId,
		MsgType:      ctx.MsgType,
		Content:      ctx.Content,
		SendUserName: senderUser,
		GroupName:    groupName,
		IsRead:       0,
		BaseStr:      string(msgStr),
		DateTime:     time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
	}

	MongoClient.Save(msg, "message")
	ctx.Next()
}
