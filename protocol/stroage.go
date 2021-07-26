package protocol

import (
	"encoding/json"
	"log"
	"net/http"
	"web-wechat/db"
)

// Storage 身份信息, 维持整个登陆的Session会话
type Storage struct {
	LoginInfo *LoginInfo
	Request   *BaseRequest
	Response  *WebInitResponse
}

type HotReloadStorageItem struct {
	Cookies      map[string][]*http.Cookie
	BaseRequest  *BaseRequest
	LoginInfo    *LoginInfo
	WechatDomain WechatDomain
}

// HotReloadStorage 热登陆存储接口
type HotReloadStorage interface {
	GetHotReloadStorageItem() HotReloadStorageItem // 获取HotReloadStorageItem
	Dump(item HotReloadStorageItem) error          // 实现该方法, 将必要信息进行序列化
	Load() error                                   // 实现该方法, 将存储媒介的内容反序列化
}

// JsonFileHotReloadStorage 实现HotReloadStorage接口
// 默认以json文件的形式存储
type JsonFileHotReloadStorage struct {
	item     HotReloadStorageItem
	filename string
}

// Dump 将信息写入json文件
//func (f *JsonFileHotReloadStorage) Dump(cookies map[string][]*http.Cookie, req *BaseRequest, info *LoginInfo) error {
//
//	file, err := os.OpenFile(f.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
//
//	if err != nil {
//		return err
//	}
//
//	defer file.Close()
//
//	f.Cookie = cookies
//	f.Req = req
//	f.Info = info
//
//	data, err := json.Marshal(f)
//	if err != nil {
//		return err
//	}
//	_, err = file.Write(data)
//	return err
//}

// Dump 将信息写入Redis
func (f *JsonFileHotReloadStorage) Dump(item HotReloadStorageItem) error {
	f.item = item
	data, err := json.Marshal(f.item)
	if err != nil {
		log.Println("序列化微信热登录信息失败：", err.Error())
		return err
	}
	// 保存信息到Redis
	//err = set(f.filename, string(data))
	err = db.SetRedisWithTimeout(f.filename, string(data), "86400")
	if err != nil {
		log.Println("保存微信热登录信息失败：", err.Error())
		return err
	}
	return nil
}

// Load 从文件中读取信息
//func (f *JsonFileHotReloadStorage) Load() error {
//	file, err := os.Open(f.filename)
//
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//	var buffer bytes.Buffer
//	if _, err := buffer.ReadFrom(file); err != nil {
//		return err
//	}
//	err = json.Unmarshal(buffer.Bytes(), f)
//	return err
//}

// Load 从Redis读取信息
func (f *JsonFileHotReloadStorage) Load() error {
	// 从Redis获取热登录数据
	data, err := db.GetRedis(f.filename)
	if err != nil {
		log.Println("读取微信热登录数据失败：", err.Error())
		return err
	}
	// 反序列化热登录数据
	err = json.Unmarshal([]byte(data), &f.item)
	return err
}

// GetHotReloadStorageItem 获取登录信息
func (f *JsonFileHotReloadStorage) GetHotReloadStorageItem() HotReloadStorageItem {
	return f.item
}

// NewJsonFileHotReloadStorage 新建一个JsonStorage对象
func NewJsonFileHotReloadStorage(filename string) *JsonFileHotReloadStorage {
	return &JsonFileHotReloadStorage{filename: filename}
}
