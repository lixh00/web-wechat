package core

import (
	"web-wechat/utils"
)

// RedisConfig Redis配置
var (
	RedisConfig redisConfig
	MySQLConfig mysqlConfig
)

// Redis配置
type redisConfig struct {
	// Redis主机
	Host string
	// Redis端口
	Port string
	// Redis密码
	Password string
	// Redis库
	Db int
}

// MySQL配置
type mysqlConfig struct {
	// 主机
	Host string
	// 端口
	Port string
	// 用户名
	Username string
	// 密码
	Password string
	// 数据库名称
	DbName string
}

// InitRedisConfig 初始化Redis配置
func InitRedisConfig() {
	// RedisHost Redis主机
	host := utils.GetEnvVal("REDIS_HOST", "wechat_redis")
	// RedisPort Redis端口
	port := utils.GetEnvVal("REDIS_PORT", "6379")
	// RedisPassword Redis密码
	password := utils.GetEnvVal("REDIS_PWD", "")
	// Redis库
	db := utils.GetEnvIntVal("REDIS_DB", 0)

	RedisConfig = redisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		Db:       db,
	}
}
