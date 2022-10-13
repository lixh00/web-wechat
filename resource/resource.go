package resource

import (
	"encoding/json"
	"errors"
	"os"
)

// LoadCarbonLanguageZhCn 加载Carbon本地化配置
func LoadCarbonLanguageZhCn() (data map[string]string, err error) {
	fileName := "resource/carbon-language-zh-CN.json"
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		err = errors.New("资源文件读取失败: " + err.Error())
		return
	}
	if json.Unmarshal(bytes, &data) != nil {
		err = errors.New("资源文件解析失败: " + err.Error())
	}
	return
}
