package db

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
	"web-wechat/core"
	"web-wechat/logger"
)

// Redis连接对象
type redisConn struct {
	client redis.Conn
}

var RedisClient redisConn

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
		logger.Log.Panicf("Redis初始化连接失败: %v", err.Error())
		//os.Exit(1)
	} else {
		logger.Log.Info("Redis连接初始化成功")
		RedisClient = redisConn{
			client: conn,
		}
	}

	//defer c.Close()
}

// GetData 获取数据
func (r *redisConn) GetData(key string) (string, error) {
	return redis.String(r.client.Do("get", key))
}

// GetKeys 获取key列表
func (r *redisConn) GetKeys(key string) ([]string, error) {
	return redis.Strings(r.client.Do("keys", key))
}

// Set 保存数据
func (r *redisConn) Set(key string, value string) error {
	_, err := r.client.Do("set", key, value)
	if err != nil {
		return errors.New("Redis保存数据失败")
	}
	return nil
}

// SetWithTimeout 保存带过期时间的数据(单位：秒)
func (r *redisConn) SetWithTimeout(key string, value string, timeout string) error {
	_, err := r.client.Do("set", key, value, "EX", timeout)
	if err != nil {
		return errors.New("Redis保存数据失败")
	}
	return nil
}

// Del 根据key删除Redis数据
func (r *redisConn) Del(key string) error {
	_, err := r.client.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}
