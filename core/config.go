package core

import (
	"fmt"
)

// SystemConfig 系统配置
var SystemConfig systemConfig

// 系统配置
type systemConfig struct {
	RedisConfig   redisConfig `mapstructure:"redis"`
	MySQLConfig   mysqlConfig `mapstructure:"mysql"`
	OssConfig     ossConfig   `mapstructure:"oss"`
	MongoDbConfig mongoConfig `mapstructure:"mongodb"`
}

// Redis配置
type redisConfig struct {
	Host     string `mapstructure:"host"`     // Redis主机
	Port     string `mapstructure:"port"`     // Redis端口
	Password string `mapstructure:"password"` // Redis密码
	Db       int    `mapstructure:"db"`       // Redis库
}

// MySQL配置
type mysqlConfig struct {
	Host     string `mapstructure:"host"`     // 主机
	Port     string `mapstructure:"port"`     // 端口
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
	DbName   string `mapstructure:"database"` // 数据库名称
}

type ossConfig struct {
	Endpoint        string `mapstructure:"endpoint"`        // 接口地址
	AccessKeyID     string `mapstructure:"accessKeyID"`     // 账号
	SecretAccessKey string `mapstructure:"secretAccessKey"` // 密码
	BucketName      string `mapstructure:"bucket"`          // 桶名称
	UseSsl          bool   `mapstructure:"ssl"`             // 是否使用SSL
}

type mongoConfig struct {
	Host     string `mapstructure:"host"`     // 地址
	Port     int    `mapstructure:"port"`     // 端口
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
	DbName   string `mapstructure:"database"` // 数据库名称
}

func (c mongoConfig) GetClientUri() string {
	return fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?ssl=false&authSource=admin", c.Username, c.Password, c.Host, c.Port, c.DbName)
}

//
//// InitRedisConfig 初始化Redis配置
//func InitRedisConfig() {
//	// RedisHost Redis主机
//	host := utils.GetEnvVal("REDIS_HOST", "wechat_redis")
//	// RedisPort Redis端口
//	port := utils.GetEnvVal("REDIS_PORT", "6379")
//	// RedisPassword Redis密码
//	password := utils.GetEnvVal("REDIS_PWD", "")
//	// Redis库
//	db := utils.GetEnvIntVal("REDIS_DB", 0)
//
//	RedisConfig = redisConfig{
//		Host:     host,
//		Port:     port,
//		Password: password,
//		Db:       db,
//	}
//}
//
//// InitOssConfig 初始化OSS配置
//func InitOssConfig() {
//	endpoint := utils.GetEnvVal("OSS_ENDPOINT", "wechat_oss")
//	accessKeyID := utils.GetEnvVal("OSS_KEY", "minio")
//	secretAccessKey := utils.GetEnvVal("OSS_SECRET", "minio")
//	bucketName := utils.GetEnvVal("OSS_BUCKET", "wechat")
//	useSSL := utils.GetEnvBoolVal("OSS_SSL", true)
//
//	OssConfig = ossConfig{
//		Endpoint:        endpoint,
//		AccessKeyID:     accessKeyID,
//		SecretAccessKey: secretAccessKey,
//		BucketName:      bucketName,
//		UseSsl:          useSSL,
//	}
//}
//
//// InitMongoConfig 初始化MongoDB配置
//func InitMongoConfig() {
//	host := utils.GetEnvVal("MONGO_HOST", "wechat_mongo")
//	port := utils.GetEnvIntVal("MONGO_PORT", 27017)
//	user := utils.GetEnvVal("MONGO_USER", "wechat")
//	password := utils.GetEnvVal("MONGO_PWD", "wechat")
//	dbName := utils.GetEnvVal("MONGO_DB", "web-wechat")
//
//	MongoDbConfig = mongoConfig{host, port, user, password, dbName}
//}
