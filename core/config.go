package core

import (
	"web-wechat/utils"
)

// RedisConfig Redis配置
var (
	RedisConfig redisConfig
	MySQLConfig mysqlConfig
	OssConfig   ossConfig
)

// Redis配置
type redisConfig struct {
	Host     string // Redis主机
	Port     string // Redis端口
	Password string // Redis密码
	Db       int    // Redis库
}

// MySQL配置
type mysqlConfig struct {
	Host     string // 主机
	Port     string // 端口
	Username string // 用户名
	Password string // 密码
	DbName   string // 数据库名称
}

type ossConfig struct {
	Endpoint        string // 接口地址
	AccessKeyID     string // 账号
	SecretAccessKey string // 密码
	BucketName      string // 桶名称
	UseSsl          bool   // 是否使用SSL
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

// InitOssConfig 初始化OSS配置
func InitOssConfig() {
	endpoint := utils.GetEnvVal("OSS_ENDPOINT", "wechat_oss")
	accessKeyID := utils.GetEnvVal("OSS_KEY", "minio")
	secretAccessKey := utils.GetEnvVal("OSS_SECRET", "minio")
	bucketName := utils.GetEnvVal("OSS_BUCKET", "wechat")
	useSSL := utils.GetEnvBoolVal("OSS_SSL", true)

	OssConfig = ossConfig{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
		UseSsl:          useSSL,
	}
}
