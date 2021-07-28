package global

import (
	"errors"
	"github.com/robfig/cron"
	"sync/atomic"
	"web-wechat/db"
	"web-wechat/logger"
	"web-wechat/protocol"
)

// InitWechatBotsMap 初始化WechatBots
func InitWechatBotsMap() {
	wechatBots = make(map[string]*protocol.Bot)
}

// GetBot 获取Bot对象
func GetBot(appKey string) *protocol.Bot {
	return wechatBots[appKey]
}

// SetBot 保存Bot对象
func SetBot(appKey string, bot *protocol.Bot) {
	wechatBots[appKey] = bot
}

// CheckBot 预检AppKey是否存在登录记录且登录状态是否正常
func CheckBot(appKey string) error {
	// 判断指定AppKey是不是有登录信息
	bot := GetBot(appKey)
	if nil == bot {
		return errors.New("未获取到登录记录")
	}
	// 判断在线状态是否正常
	if !bot.Alive() {
		return errors.New("微信在线状态异常，请重新登录")
	}
	return nil
}

// InitWechatBotHandle 初始化微信机器人
func InitWechatBotHandle() *protocol.Bot {
	bot := protocol.DefaultBot(protocol.Desktop)

	// 定义读取消息错误回调函数
	var getMessageErrorCount int32
	bot.GetMessageErrorHandler = func(err error) {
		atomic.AddInt32(&getMessageErrorCount, 1)
		// 如果发生了三次错误,那么直接退出
		if getMessageErrorCount == 3 {
			logger.Log.Errorf("获取消息发生错误达到三次，直接退出。错误信息：%v", err.Error())
			_ = bot.Logout()
		}
	}
	// 注册消息处理函数
	wechatMessageHandle(bot)
	// 获取消息发生错误
	//bot.MessageOnError()
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
			if ctx.IsSendByGroup() {
				// 取出消息在群里面的发送者
				senderInGroup, _ := ctx.SenderInGroup()
				logger.Log.Infof("[群聊][收到新文字消息] == 发信人：%v[%v] ==> 内容：%v", sender.NickName,
					senderInGroup.NickName, ctx.Content)
			} else {
				logger.Log.Infof("[好友][收到新文字消息] == 发信人：%v ==> 内容：%v", sender.NickName, ctx.Content)
			}
		}
		ctx.Next()
	}
	dispatcher.OnText(textHandle)

	// 处理其他消息
	otherHandle := func(ctx *protocol.MessageContext) {
		sender, _ := ctx.Sender()
		logger.Log.Infof("[收到新消息] 发送者：%v == 消息类型: %v ==> 内容：%v",
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
					logger.Log.Errorf("【%v】更新热登录数据失败", user.NickName)
				}
				logger.Log.Infof("【%v】热登录数据更新成功", user.NickName)
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
	// 添加一个每半小时执行一次的执行器
	_ = c.AddFunc("0 0/30 * * * ? ", func() {
		var errKey []string
		for k, bot := range wechatBots {
			if bot.Alive() {
				user, _ := bot.GetCurrentUser()
				file, err := user.FileHelper()
				if err != nil {
					logger.Log.Errorf("获取文件助手失败 ====> %v", err.Error())
					continue
				}
				if _, err := file.SendText(protocol.ZombieText); err != nil {
					logger.Log.Errorf("【%v】保活失败 ====> %v", user.NickName, err.Error())
					errKey = append(errKey, k)
					continue
				}
				logger.Log.Infof("【%v】保活成功", user.NickName)
			}
			continue
		}
		// 清理掉无效的登录实例
		for _, key := range errKey {
			// 取出热登录信息登录一次，如果登录失败就删除实例
			bot := GetBot(key)
			storage := protocol.NewJsonFileHotReloadStorage("wechat:login:" + key)
			if err := bot.HotLogin(storage, false); err != nil {
				logger.Log.Errorf("[%v] 热登录失败，错误信息：%v", key, err.Error())
				// 登录失败，删除热登录数据
				if err := db.DelRedis(key); err != nil {
					logger.Log.Errorf("[%v] Redis缓存删除失败，错误信息：%v", key, err.Error())
				}
				delete(wechatBots, key)
			} else {
				// 登录成功，更新实例
				SetBot(key, bot)
			}
		}
	})
	// 新启一个协程，运行定时任务
	go c.Start()
	// 等待停止信号结束任务
	defer c.Stop()
}
