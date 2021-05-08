package global

import (
	"web-wechat/protocol"
)

var (
	// 登录用户的Bot对象
	wechatBots map[string]*protocol.Bot
)
