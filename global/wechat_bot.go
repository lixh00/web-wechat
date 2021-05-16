package global

import (
	"github.com/robfig/cron"
	"log"
	"web-wechat/protocol"
)

// InitWechatBotsMap 初始化WechatBots
func InitWechatBotsMap() {
	wechatBots = make(map[string]*protocol.Bot)
}

// GetBot 获取Bot对象
func GetBot(uuid string) *protocol.Bot {
	return wechatBots[uuid]
}

// SetBot 保存Bot对象
func SetBot(uuid string, bot *protocol.Bot) {
	wechatBots[uuid] = bot
}

// InitWechatBotHandle 初始化微信机器人
func InitWechatBotHandle() *protocol.Bot {
	bot := protocol.DefaultBot(protocol.Desktop)

	// 定义读取消息错误回调函数
	bot.GetMessageErrorHandler = func(err error) {
		log.Println("获取消息发生错误：", err.Error())
	}
	// 注册消息处理函数
	wechatMessageHandle(bot)
	// 返回机器人对象
	return bot
}

// 微信消息处理函数
func wechatMessageHandle(bot *protocol.Bot) {
	dispatcher := protocol.NewMessageMatchDispatcher()
	// 设置为异步处理
	dispatcher.SetAsync(true)

	// 处理文字消息
	textHandle := func(ctx *protocol.MessageContext) {
		if !ctx.IsSendBySelf() {
			sender, _ := ctx.Sender()
			log.Printf("[收到新文字消息] == 发信人：%v ==> 内容：%v \n", sender.NickName, ctx.Content)
		}
		ctx.Next()
	}
	dispatcher.OnText(textHandle)

	// 处理其他消息
	otherHandle := func(ctx *protocol.MessageContext) {
		sender, _ := ctx.Sender()
		log.Printf("[收到新消息] 发送者：%v == 消息类型: %v ==> 内容：%v \n",
			sender.NickName, ctx.MsgType, protocol.XmlFormString(ctx.Content))
	}
	dispatcher.RegisterHandler(func(message *protocol.Message) bool {
		// 处理除文字消息和通知消息之外，并且不是自己发送的消息
		return !message.IsText() && !message.IsNotify() && !message.IsSendBySelf()
	}, otherHandle)

	// 注册消息处理函数
	bot.MessageHandler = protocol.DispatchMessage(dispatcher)
}

// UpdateHotLoginData 更新热登录数据
func UpdateHotLoginData() {
	// 创建一个新的定时任务管理器
	c := cron.New()
	// 添加一个每小时执行一次的执行器
	_ = c.AddFunc("0 5 * * * ? ", func() {
		for _, bot := range wechatBots {
			if bot.Alive() {
				user, _ := bot.GetCurrentUser()
				if err := bot.DumpHotReloadStorage(); err != nil {
					log.Printf("【%v】更新热登录数据失败 \n", user.NickName)
				}
				log.Printf("【%v】热登录数据更新成功 \n", user.NickName)
			}
			continue
		}
	})
	// 新启一个协程，运行定时任务
	go c.Start()
	// 等待停止信号结束任务
	defer c.Stop()
}

// KeepAliveHandle 保活，每三时自动给文件传输助手发一条消息
func KeepAliveHandle() {
	// 创建一个新的定时任务管理器
	c := cron.New()
	// 添加一个每小时执行一次的执行器
	_ = c.AddFunc("0 0 * * * ? ", func() {
		for _, bot := range wechatBots {
			if bot.Alive() {
				user, _ := bot.GetCurrentUser()
				file, _ := user.FileHelper()
				if _, err := file.SendText("芜湖"); err != nil {
					log.Printf("【%v】保活失败 \n", user.NickName)
				}
				log.Printf("【%v】保活成功 \n", user.NickName)
			}
			continue
		}
	})
	// 新启一个协程，运行定时任务
	go c.Start()
	// 等待停止信号结束任务
	defer c.Stop()
}
