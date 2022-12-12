package core

import (
	"fmt"
)

// SystemConfig 系统配置
var SystemConfig systemConfig

// 系统配置
type systemConfig struct {
	RedisConfig   redisConfig  `mapstructure:"redis"`
	MySQLConfig   mysqlConfig  `mapstructure:"mysql"`
	OssConfig     ossConfig    `mapstructure:"oss"`
	MongoDbConfig mongoConfig  `mapstructure:"mongodb"`
	OpenAiConfig  openAiConfig `mapstructure:"openai"`
}

//	openAiConfig
//	@description: openai配置
type openAiConfig struct {
	Enable bool   `mapstructure:"enable"` // 是否启用
	ApiKey string `mapstructure:"apikey"` // apikey
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
