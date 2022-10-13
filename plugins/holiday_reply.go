package plugins

import (
	"fmt"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"web-wechat/utils"
)

// checkHoliday
// @description: 放假倒计时
// @receiver weChatPlugin
// @param ctx
func (weChatPlugin) checkHoliday(ctx *openwechat.MessageContext) {
	dd := []string{"今天", "明天", "后天"}
	// 获取最近的节假日或周末
	d, t := utils.OffDuty().GetNextHolidayOrWeekend()
	replyStr := fmt.Sprintf("距离%v还有%v天", d, t)
	if t == 0 {
		replyStr = fmt.Sprintf("今天已经是%v了啊，不会你还没放假吧？", d)
	}
	if t > 0 && t < 3 {
		replyStr = fmt.Sprintf("%v就是%v啦，再坚持一下咯~", dd[t], d)
	}
	if _, err := ctx.ReplyText(replyStr); err != nil {
		log.Errorf("[放假倒计时]消息回复失败: %v", err.Error())
	}
	ctx.Next()
}

// checkFestivals
// @description: 过节倒计时
// @receiver weChatPlugin
// @param ctx
func (weChatPlugin) checkFestivals(ctx *openwechat.MessageContext) {
	dd := []string{"今天", "明天", "后天"}
	// 获取最近的节假日或周末
	d, t := utils.OffDuty().GetNextHoliday()
	replyStr := fmt.Sprintf("距离%v还有%v天", d, t)
	if t == 0 {
		replyStr = fmt.Sprintf("今天已经是%v了啊，不会你还没放假吧？", d)
	}
	if t > 0 && t < 3 {
		replyStr = fmt.Sprintf("%v就是%v啦，再坚持一下咯~", dd[t], d)
	}
	if _, err := ctx.ReplyText(replyStr); err != nil {
		log.Errorf("[过节倒计时]消息回复失败: %v", err.Error())
	}
	ctx.Next()
}
