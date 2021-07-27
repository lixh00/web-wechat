package db

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
	"web-wechat/core"
	"web-wechat/logger"
)

// RedisConn Redis连接对象
var redisConn redis.Conn

// InitRedisConnHandle 初始化Redis连接对象
func InitRedisConnHandle() {
	// 读取配置
	core.InitRedisConfig()
	// 初始化连接
	conn, err := redis.Dial("tcp",
		// Redis连接信息
		fmt.Sprintf("%s:%s", core.RedisConfig.Host, core.RedisConfig.Port),
		// 密码
		redis.DialPassword(core.RedisConfig.Password),
		// 默认使用数据库
		redis.DialDatabase(core.RedisConfig.Db),
		redis.DialKeepAlive(1*time.Second),
		redis.DialConnectTimeout(5*time.Second),
		redis.DialReadTimeout(1*time.Second),
		redis.DialWriteTimeout(1*time.Second))

	if err != nil {
		logger.Log.Errorf("Redis初始化连接失败: %v", err.Error())
		os.Exit(1)
	} else {
		redisConn = conn
	}

	//defer c.Close()
}

// GetRedis 获取数据
func GetRedis(key string) (string, error) {
	return redis.String(redisConn.Do("get", key))
}

// GetRedisKeys 获取key列表
func GetRedisKeys(key string) ([]string, error) {
	return redis.Strings(redisConn.Do("keys", key))
}

// SetRedis 保存数据
func SetRedis(key string, value string) error {
	_, err := redisConn.Do("set", key, value)
	if err != nil {
		return errors.New("Redis保存数据失败")
	}
	return nil
}

// SetRedisWithTimeout 保存带过期时间的数据(单位：秒)
func SetRedisWithTimeout(key string, value string, timeout string) error {
	_, err := redisConn.Do("set", key, value, "EX", timeout)
	if err != nil {
		return errors.New("Redis保存数据失败")
	}
	return nil
}

// DelRedis 根据key删除Redis数据
func DelRedis(key string) error {
	_, err := redisConn.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}
