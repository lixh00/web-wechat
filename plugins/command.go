package plugins

import (
	"github.com/eatmoreapple/openwechat"
)

// Command
// @description:
// @receiver weChatPlugin
// @param ctx
func (p weChatPlugin) Command(ctx *openwechat.MessageContext) {
	switch ctx.Content {
	case "放假倒计时":
		p.checkHoliday(ctx)
	case "过节倒计时":
		p.checkFestivals(ctx)
	case "下班倒计时":
		p.checkOffWork(ctx)
	}

	ctx.Next()
}
