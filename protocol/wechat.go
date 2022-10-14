package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitee.ltd/lxh/logger/log"
	. "github.com/eatmoreapple/openwechat"
	"net/url"
	"time"
	. "web-wechat/db"
)

type WechatBot struct {
	Bot
}

type RedisHotReloadStorage struct {
	Key    string
	reader *bytes.Reader
}

func NewRedisHotReloadStorage(key string) *RedisHotReloadStorage {
	return &RedisHotReloadStorage{Key: key}
}

// GetUUID 获取UUID
func (b *WechatBot) GetUUID() (*string, error) {
	uuid, err := b.Caller.GetLoginUUID()
	if err != nil {
		return nil, err
	}
	// 二维码获取回调
	if b.UUIDCallback != nil {
		b.UUIDCallback(uuid)
	}
	return &uuid, nil
}

// LoginWithUUID 根据传入的UUID登录
func (b *WechatBot) LoginWithUUID(uuid string) error {
	for {
		// 长轮询检查是否扫码登录
		resp, err := b.Caller.CheckLogin(uuid)
		if err != nil {
			return err
		}
		//log.Infof("CheckLogin: %v ==> %v", resp.Code, string(resp.Raw))
		switch resp.Code {
		case StatusSuccess:
			// 判断是否有登录回调，如果有执行它
			if b.LoginCallBack != nil {
				b.LoginCallBack(resp.Raw)
			}
			return b.HandleLogin(resp.Raw)

		case StatusScanned:
			// 执行扫码回调
			if b.ScanCallBack != nil {
				b.ScanCallBack(resp.Raw)
			}
		case StatusTimeout:
			return errors.New("login timeout")
		case StatusWait:
			continue
		}
	}
}

// HotLoginWithUUID 根据UUID热登录 TODO 这儿需要优化，不太能用的亚子
func (b *WechatBot) HotLoginWithUUID(uuid string, storage HotReloadStorage, retry ...bool) error {
	//b.IsHot = true
	b.HotReloadStorage = storage

	var err error

	// 如果load出错了,就执行正常登陆逻辑
	// 第一次没有数据load都会出错的
	var buffer bytes.Buffer
	if _, err = buffer.ReadFrom(storage); err != nil {
		return b.LoginWithUUID(uuid)
	}

	var item HotReloadStorageItem
	if err = json.NewDecoder(&buffer).Decode(&item); err != nil {
		return err
	}

	cookies := item.Cookies
	for u, ck := range cookies {
		path, err := url.Parse(u)
		if err != nil {
			return err
		}
		b.Caller.Client.Jar.SetCookies(path, ck)
	}
	b.Storage.LoginInfo = item.LoginInfo
	b.Storage.Request = item.BaseRequest
	b.Caller.Client.Domain = item.WechatDomain

	// 如果webInit出错,则说明可能身份信息已经失效
	// 如果retry为True的话,则进行正常登陆
	if err = b.WebInit(); err != nil {
		if len(retry) > 0 && retry[0] {
			return b.LoginWithUUID(uuid)
		}
	}
	return err
}

// Load 重写热登录数据加载，从Redis取数据
func (f *RedisHotReloadStorage) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		// 从Redis获取热登录数据
		data, err := RedisClient.GetData(f.Key)
		if err != nil {
			log.Errorf("读取热登录数据失败: %v", err)
			return 0, err
		}
		f.reader = bytes.NewReader([]byte(data))
	}
	return f.reader.Read(p)
}

// Dump 重写更新热登录数据，保存到Redis
func (f *RedisHotReloadStorage) Write(p []byte) (n int, err error) {
	err = RedisClient.SetWithTimeout(f.Key, string(p), 2*24*time.Hour)
	if err != nil {
		log.Errorf("保存微信热登录信息失败: %v", err.Error())
		return
	}
	return len(p), nil
}

// Close 需要关闭
func (f *RedisHotReloadStorage) Close() error {
	f.reader = nil
	return nil
}
