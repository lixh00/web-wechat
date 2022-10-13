package global

import (
	"encoding/json"
	"errors"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron"
	"time"
	. "web-wechat/db"
	"web-wechat/handler"
	"web-wechat/protocol"
)

// InitWechatBotsMap 初始化WechatBots
func InitWechatBotsMap() {
	wechatBots = make(map[string]*protocol.WechatBot)
}

// GetBot 获取Bot对象
func GetBot(appKey string) *protocol.WechatBot {
	return wechatBots[appKey]
}

// SetBot 保存Bot对象
func SetBot(appKey string, bot *protocol.WechatBot) {
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
func InitWechatBotHandle() *protocol.WechatBot {
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 定义读取消息错误回调函数
	//var getMessageErrorCount int32
	//bot.GetMessageErrorHandler = func(err error) {
	//	atomic.AddInt32(&getMessageErrorCount, 1)
	//	// 如果发生了三次错误,那么直接退出
	//	if getMessageErrorCount == 3 {
	//		log.Errorf("获取消息发生错误达到三次，直接退出。错误信息：%v", err.Error())
	//		_ = bot.Logout()
	//	}
	//}

	// 设置心跳回调
	bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
		if resp.RetCode == "1100" {
			log.Errorf("微信已退出")
			// do something
		}
		switch resp.Selector {
		case "0":
			log.Debugf("正常")
		case "2", "6":
			log.Debugf("有新消息")
		case "7":
			log.Debugf("进入/离开聊天界面")
			err := bot.WebInit()
			if err != nil {
				// 短信通知一下
				// do something
				//log.Panicf("重新初始化失败: %v", err)
			}
		default:
			log.Debugf("RetCode: %s  Selector: %s", resp.RetCode, resp.Selector)
		}
	}

	// 注册消息处理函数
	handler.HandleMessage(bot)
	// 获取消息发生错误
	//bot.MessageOnError()
	// 返回机器人对象
	return &protocol.WechatBot{Bot: *bot}
}

// UpdateHotLoginData 更新热登录数据
func UpdateHotLoginData() {
	// 创建一个新的定时任务管理器
	c := cron.New()
	// 添加一个每三十分钟执行一次的执行器
	_ = c.AddFunc("0 0/30 * * * ? ", dealUpdateHotLoginData)
	// 新启一个协程，运行定时任务
	go c.Start()
	// 等待停止信号结束任务
	defer c.Stop()
}

// 更新热登录数据函数
func dealUpdateHotLoginData() {
	for key, bot := range wechatBots {
		if bot.Alive() {
			user, _ := bot.GetCurrentUser()
			// 判断是否需要热登录
			if checkHotLogin(key) {
				storage := protocol.NewRedisHotReloadStorage("wechat:login:" + key)
				if err := bot.HotLogin(storage, false); err != nil {
					log.Errorf("[%v]定时热登录失败: %v", user.NickName, err)
					continue
				}
				log.Debugf("[%v]热登录成功", user.NickName)
			} else {
				log.Debugf("[%v]到期时间大于11小时，暂不重新登录", user.NickName)
			}
			if err := bot.DumpHotReloadStorage(); err != nil {
				log.Errorf("【%v】更新热登录数据失败，错误信息: %v", user.NickName, err)
			}
			log.Infof("【%v】热登录数据更新成功", user.NickName)
		}
		continue
	}
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
					log.Errorf("获取文件助手失败 ====> %v", err.Error())
					continue
				}
				if _, err := file.SendText(openwechat.ZombieText); err != nil {
					log.Errorf("【%v】保活失败 ====> %v", user.NickName, err.Error())
					errKey = append(errKey, k)
					continue
				}
				log.Infof("【%v】保活成功", user.NickName)
			}
			continue
		}
		// 清理掉无效的登录实例
		for _, key := range errKey {
			// 取出热登录信息登录一次，如果登录失败就删除实例
			bot := GetBot(key)
			storage := protocol.NewRedisHotReloadStorage("wechat:login:" + key)
			if err := bot.HotLogin(storage, false); err != nil {
				log.Errorf("[%v] 热登录失败，错误信息：%v", key, err.Error())
				// 登录失败，删除热登录数据
				if err := RedisClient.Del(key); err != nil {
					log.Errorf("[%v] Redis缓存删除失败，错误信息：%v", key, err.Error())
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

// 判断是否需要热登录
func checkHotLogin(appKey string) bool {
	jsonStr, err := RedisClient.GetData("wechat:login:" + appKey)
	if err != nil {
		log.Errorf("热登录数据获取失败: %v", err)
		return false
	}
	// 热登录数据转化为实体
	var hotLoginData openwechat.HotReloadStorageItem
	// 反序列化热登录数据
	err = json.Unmarshal([]byte(jsonStr), &hotLoginData)
	if err != nil {
		log.Errorf("反序列化热登录数据失败: %v", err)
		return false
	}

	for _, cookies := range hotLoginData.Cookies {
		if len(cookies) > 0 {
			//log.Debugf("保存的Cookie值: %v", cookies)
			for _, cookie := range cookies {
				if cookie.Name == "wxsid" {
					loc, _ := time.LoadLocation("GMT")
					expiresGMTTime, _ := time.ParseInLocation("Mon, 02-Jan-2006 15:04:05 GMT", cookie.RawExpires, loc)
					loc2, _ := time.LoadLocation("Local")
					expiresLocalTime := expiresGMTTime.In(loc2)
					//log.Debugf("登录有效到期时间: %v", expiresLocalTime)
					overHours := expiresLocalTime.Sub(time.Now().In(loc2)).Hours()
					if overHours < 1 {
						log.Errorf("[%v]状态异常", appKey)
					}
					log.Debugf("[%v]距离到期时间还剩 %v 小时", appKey, overHours)
					// 小于11小时就返回true，表示需要热登录
					return overHours < 11
				}
			}
		}
	}
	return false
}
