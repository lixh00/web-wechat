package handler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"io/ioutil"
	"net/http"
	"strings"
	"web-wechat/core"
	"web-wechat/oss"
)

// EmoticonMessageData 表情包消息结构体
type EmoticonMessageData struct {
	XMLName xml.Name `xml:"msg"`
	Emoji   struct {
		Fromusername      string `xml:"fromusername,attr"`
		Tousername        string `xml:"tousername,attr"`
		Type              string `xml:"type,attr"`
		Idbuffer          string `xml:"idbuffer,attr"`
		Md5               string `xml:"md5,attr"`
		Len               string `xml:"len,attr"`
		Productid         string `xml:"productid,attr"`
		Androidmd5        string `xml:"androidmd5,attr"`
		Androidlen        string `xml:"androidlen,attr"`
		S60v3md5          string `xml:"s60v3md5,attr"`
		S60v3len          string `xml:"s60v3len,attr"`
		S60v5md5          string `xml:"s60v5md5,attr"`
		S60v5len          string `xml:"s60v5len,attr"`
		Cdnurl            string `xml:"cdnurl,attr"`
		Designerid        string `xml:"designerid,attr"`
		Thumburl          string `xml:"thumburl,attr"`
		Encrypturl        string `xml:"encrypturl,attr"`
		Aeskey            string `xml:"aeskey,attr"`
		Externurl         string `xml:"externurl,attr"`
		Externmd5         string `xml:"externmd5,attr"`
		Width             string `xml:"width,attr"`
		Height            string `xml:"height,attr"`
		Tpurl             string `xml:"tpurl,attr"`
		Tpauthkey         string `xml:"tpauthkey,attr"`
		Attachedtext      string `xml:"attachedtext,attr"`
		Attachedtextcolor string `xml:"attachedtextcolor,attr"`
		Lensid            string `xml:"lensid,attr"`
		Emojiattr         string `xml:"emojiattr,attr"`
		Linkid            string `xml:"linkid,attr"`
	} `xml:"emoji"`
	Gameext struct {
		Type    string `xml:"type,attr"`
		Content string `xml:"content,attr"`
	} `xml:"gameext"`
}

// 表情包消息处理
func emoticonMessageHandle(ctx *openwechat.MessageContext) {
	// 取出发送者
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}

	// 判断消息是不是表情商店的，如果是，不支持解析
	if !strings.Contains(ctx.Content, "<msg>") {
		log.Debugf("原始数据: %v", ctx.Content)
		log.Infof("[收到新表情包消息] == 发信人：%v ==> 内容：「收到了一个表情，请在手机上查看」", senderUser)
	} else {
		// 解析表情包
		var data EmoticonMessageData
		if err := xml.Unmarshal([]byte(ctx.Content), &data); err != nil {
			log.Errorf("消息解析失败: %v", err.Error())
			log.Debugf("原始内容: %v", ctx.Content)
			return
		} else {
			log.Infof("[收到新表情包消息] == 发信人：%v ==> 内容：%v", senderUser, data.Emoji.Md5)
			// 下载图片资源
			fileResp, err := ctx.GetFile()
			if err != nil {
				log.Errorf("表情包下载失败: %v", err.Error())
				return
			}
			defer fileResp.Body.Close()
			imgFileByte, err := ioutil.ReadAll(fileResp.Body)
			if err != nil {
				log.Errorf("表情包读取错误: %v", err.Error())
				return
			} else {
				// 读取文件相关信息
				contentType := http.DetectContentType(imgFileByte)
				fileType := strings.Split(contentType, "/")[1]
				fileName := fmt.Sprintf("%v.%v", ctx.MsgId, fileType)
				if user, err := ctx.Bot.GetCurrentUser(); err == nil {
					uin := user.Uin
					fileName = fmt.Sprintf("%v/%v", uin, fileName)
				}
				// 上传文件(reader2解决上传空文件的BUG,因为http.Response.Body只允许读一次)
				reader2 := ioutil.NopCloser(bytes.NewReader(imgFileByte))
				flag := oss.SaveToOss(reader2, contentType, fileName)
				if flag {
					fileUrl := fmt.Sprintf("https://%v/%v/%v", core.OssConfig.Endpoint, core.OssConfig.BucketName, fileName)
					log.Infof("表情包保存成功，图片链接: %v", fileUrl)
					ctx.Content = fileUrl
				} else {
					log.Error("表情包保存失败")
				}
			}
		}
	}
	ctx.Next()
}
