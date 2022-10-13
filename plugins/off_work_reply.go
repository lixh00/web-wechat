package plugins

import (
	"fmt"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"github.com/golang-module/carbon"
	"time"
	"web-wechat/resource"
	"web-wechat/utils"
)

// checkOffWork
// @description: 下班倒计时
// @receiver weChatPlugin
// @param ctx
func (weChatPlugin) checkOffWork(ctx *openwechat.MessageContext) {
	// 如果不是工作日，跳过处理
	if isHoliday, h := utils.OffDuty().CheckIsHoliday(time.Now()); isHoliday {
		if _, err := ctx.ReplyText(fmt.Sprintf("不会吧不会吧，不会有人%v还在上班吧", h)); err != nil {
			log.Errorf("阴阳怪气失败: %v", err.Error())
		}
		return
	}
	// 非工作时间不执行
	if time.Now().Hour() < 9 || time.Now().Hour() >= 18 {
		if _, err := ctx.ReplyText("不会吧不会吧，不会有人这个点还没下班吧"); err != nil {
			log.Errorf("阴阳怪气失败: %v", err.Error())
		}
		return
	}

	lange := carbon.NewLanguage()
	// 读取本地化资源
	language, err := resource.LoadCarbonLanguageZhCn()
	if err != nil {
		log.Errorf("读取本地化资源失败: %v", err)
		return
	}
	lange.SetResources(language)
	car := carbon.SetLanguage(lange)
	offDutyTime := car.Now().StartOfDay().AddHours(18)
	now := car.Now()
	if _, err := ctx.ReplyText("距离下班还有 " + now.DiffInString(offDutyTime)); err != nil {
		log.Errorf("下班时间倒计时发送失败: %v", err.Error())
	}
	ctx.Next()
}
