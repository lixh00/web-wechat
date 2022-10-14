package db

import (
	"context"
	"fmt"
	"gitee.ltd/lxh/logger/log"
	"github.com/go-redis/redis/v8"
	"time"
	"web-wechat/core"
)

// Redis连接对象
type redisConn struct {
	client *redis.Client
}

var RedisClient redisConn

// InitRedisConnHandle 初始化Redis连接对象
func InitRedisConnHandle() {
	// 初始化连接
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", core.SystemConfig.RedisConfig.Host, core.SystemConfig.RedisConfig.Port),
		Password: core.SystemConfig.RedisConfig.Password,
		DB:       core.SystemConfig.RedisConfig.Db,
	})
	if err := conn.Ping(context.Background()).Err(); err != nil {
		log.Panicf("Redis初始化失败: %v", err.Error())
	}
	log.Debug("Redis连接初始化成功")
	RedisClient = redisConn{client: conn}
}

// GetData 获取数据
func (r *redisConn) GetData(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

// GetKeys 获取key列表
func (r *redisConn) GetKeys(key string) ([]string, error) {
	return r.client.Keys(context.Background(), key).Result()
}

// Set 保存数据
func (r *redisConn) Set(key string, value string) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

// SetWithTimeout 保存带过期时间的数据(单位：秒)
func (r *redisConn) SetWithTimeout(key string, value string, timeout time.Duration) error {
	return r.client.Set(context.Background(), key, value, timeout).Err()
}

// Del 根据key删除Redis数据
func (r *redisConn) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
