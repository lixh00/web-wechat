package plugins

import "github.com/eatmoreapple/openwechat"

// Status
// @description: 处理插件相关指令
// @receiver p
// @param ctx
func (p weChatPlugin) Status(ctx *openwechat.MessageContext) {
	switch ctx.Content {
	case "开启插件":
		ChangePluginStatus(true)
		_, _ = ctx.ReplyText("插件已开启")
	case "关闭插件":
		ChangePluginStatus(false)
		_, _ = ctx.ReplyText("插件已关闭")
	}
	ctx.Next()
}
