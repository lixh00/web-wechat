package utils

import (
	"os"
	"strconv"
)

// GetEnvVal 从环境变量获取字符串类型值
func GetEnvVal(key, defaultVal string) string {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defaultVal
	}
	return val
}

// GetEnvIntVal 从环境变量获取数字类型值
func GetEnvIntVal(key string, defaultVal int) int {
	valStr, exist := os.LookupEnv(key)
	if !exist {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
