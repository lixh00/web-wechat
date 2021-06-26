package protocol

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

// RedisConn Redis连接对象
var RedisConn redis.Conn

// InitRedisConnHandle 初始化Redis连接对象
func InitRedisConnHandle() {
	conn, err := redis.Dial("tcp",
		// Redis连接信息
		"10.0.0.30:6379",
		// 密码
		redis.DialPassword("pGhQKwj7DE7FbFL1"),
		// 默认使用数据库
		redis.DialDatabase(5),
		redis.DialKeepAlive(1*time.Second),
		redis.DialConnectTimeout(5*time.Second),
		redis.DialReadTimeout(1*time.Second),
		redis.DialWriteTimeout(1*time.Second))

	if err != nil {
		fmt.Println("Redis初始化连接失败: ", err.Error())
		os.Exit(1)
	} else {
		RedisConn = conn
	}

	//defer c.Close()
}

// 获取数据
func get(key string) (string, error) {
	return redis.String(RedisConn.Do("get", key))
}

// 保存数据
func set(key string, value string) error {
	_, err := RedisConn.Do("set", key, value)
	if err != nil {
		return errors.New("Redis保存数据失败")
	}
	return nil
}

// 保存带过期时间的数据(单位：秒)
func setWithTimeout(key string, value string, timeout string) error {
	_, err := RedisConn.Do("set", key, value, "EX", timeout)
	if err != nil {
		return errors.New("Redis保存数据失败")
	}
	return nil
}
