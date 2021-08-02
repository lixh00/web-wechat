package global

import (
	"web-wechat/db"
	"web-wechat/logger"
	"web-wechat/protocol"
)

// InitBotWithStart 系统启动的时候从Redis加载登录信息自动登录
func InitBotWithStart() {
	keys, err := db.GetRedisKeys("wechat:login:*")
	if err != nil {
		logger.Log.Error("获取Key失败")
		return
	}
	logger.Log.Infof("获取到登录用户信息数量：%v", len(keys))
	for _, key := range keys {
		// 提取出AppKey
		appKey := key[13:]
		// 调用热登录
		logger.Log.Debugf("当前热登录AppKey: %v", appKey)
		bot := InitWechatBotHandle()
		//storage := openwechat.NewJsonFileHotReloadStorage(key)
		storage := protocol.NewRedisHotReloadStorage("wechat:login:" + key)
		if err := bot.HotLogin(storage, false); err != nil {
			logger.Log.Infof("[%v] 热登录失败，错误信息：%v", appKey, err.Error())
			// 登录失败，删除热登录数据
			if err := db.DelRedis(key); err != nil {
				logger.Log.Errorf("[%v] Redis缓存删除失败，错误信息：%v", key, err.Error())
			}
			return
		}
		loginUser, _ := bot.GetCurrentUser()
		logger.Log.Infof("[%v]初始化自动登录成功，用户名：%v", appKey, loginUser.NickName)
		// 登录成功，写入到WechatBots
		SetBot(appKey, bot)
	}
}
