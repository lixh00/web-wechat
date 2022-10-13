package main

import (
	"gitee.ltd/lxh/logger"
	"gitee.ltd/lxh/logger/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"web-wechat/core"
	"web-wechat/db"
	"web-wechat/global"
	"web-wechat/middleware"
	"web-wechat/oss"
	"web-wechat/route"
)

func init() {
	// 手动初始化日志
	if os.Getenv("RUN_MODE") == "dev" {
		logger.InitLogger(logger.LogConfig{Mode: logger.Dev, FileEnable: true})
	}

	vp := viper.New()
	// 初始化配置
	if os.Getenv("RUN_MODE") == "dev" {
		vp.AddConfigPath("E:/Lxh/web-wechat") // 设置配置文件路径
	} else {
		vp.AddConfigPath(".") // 设置配置文件路径
	}
	vp.SetConfigName("config") // 设置配置文件名
	vp.SetConfigType("yaml")   // 设置配置文件类型
	// 读取配置文件
	if err := vp.ReadInConfig(); err != nil {
		log.Panicf("读取配置文件失败: %v", err)
	}
	// 绑定配置文件
	if err := vp.Unmarshal(&core.SystemConfig); err != nil {
		log.Panicf("配置文件解析失败: %v", err)
	}
	log.Debugf("配置文件解析完成: %v", core.SystemConfig)

	// 初始化OSS
	oss.InitOssConnHandle()
	// 初始化MongoDB
	db.InitMongoConnHandle()
	// 初始化Redis连接
	db.InitRedisConnHandle()
}

// 程序启动入口
func main() {
	// 初始化Gin
	app := gin.Default()

	// 定义全局异常处理
	app.NoRoute(core.NotFoundErrorHandler())
	// AppKey预检
	app.Use(middleware.CheckAppKeyExistMiddleware(), middleware.CheckAppKeyIsLoggedInMiddleware())
	// 初始化路由
	route.InitRoute(app)

	// 初始化WechatBotMap
	global.InitWechatBotsMap()

	// 初始化Redis里登录的数据
	global.InitBotWithStart()

	// 定时更新 Bot 的热登录数据
	global.UpdateHotLoginData()

	// 保活
	//global.KeepAliveHandle()

	// 监听端口
	_ = app.Run(":8888")
}
