package handler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"web-wechat/core"
	"web-wechat/logger"
	"web-wechat/oss"
	"web-wechat/protocol"
)

// ImageMessageData 图片消息结构体
type ImageMessageData struct {
	XMLName xml.Name `xml:"msg"`
	Img     struct {
		Text           string `xml:",chardata"`
		AesKey         string `xml:"aeskey,attr"`
		EnCryVer       string `xml:"encryver,attr"`
		CdnThumbAesKey string `xml:"cdnthumbaeskey,attr"`
		CdnThumbUrl    string `xml:"cdnthumburl,attr"`
		CdnThumbLength string `xml:"cdnthumblength,attr"`
		CdnThumbHeight string `xml:"cdnthumbheight,attr"`
		CdnThumbWidth  string `xml:"cdnthumbwidth,attr"`
		CdnMidHeight   string `xml:"cdnmidheight,attr"`
		CdnMidWidth    string `xml:"cdnmidwidth,attr"`
		CdnHdHeight    string `xml:"cdnhdheight,attr"`
		CdnHdWidth     string `xml:"cdnhdwidth,attr"`
		CdnMidImgUrl   string `xml:"cdnmidimgurl,attr"`
		Length         int64  `xml:"length,attr"`
		CdnBigImgUrl   string `xml:"cdnbigimgurl,attr"`
		HdLength       int64  `xml:"hdlength,attr"`
		Md5            string `xml:"md5,attr"`
	} `xml:"img"`
}

// 处理图片消息
func imageMessageHandle(ctx *protocol.MessageContext) {
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}
	// 解析xml为结构体
	var data ImageMessageData
	if err := xml.Unmarshal([]byte(ctx.Content), &data); err != nil {
		logger.Log.Errorf("消息解析失败: %v", err.Error())
		logger.Log.Debugf("发信人: %v ==> 原始内容: %v", senderUser, protocol.XmlFormString(ctx.Content))
		return
	} else {
		logger.Log.Infof("[收到新图片消息] == 发信人：%v ==> 内容：%v", senderUser, data.Img.AesKey)
		// 下载图片资源
		fileResp, err := ctx.GetFile()
		if err != nil {
			logger.Log.Errorf("图片下载失败: %v", err.Error())
			return
		}
		defer fileResp.Body.Close()
		imgFileByte, err := ioutil.ReadAll(fileResp.Body)
		if err != nil {
			logger.Log.Errorf("图片读取错误: %v", err.Error())
			return
		} else {
			// 读取文件相关信息
			contentType := http.DetectContentType(imgFileByte)
			fileType := strings.Split(contentType, "/")[1]
			fileName := fmt.Sprintf("%v.%v", ctx.MsgId, fileType)
			// 上传文件
			reader2 := ioutil.NopCloser(bytes.NewReader(imgFileByte))
			flag := oss.SaveToOss(reader2, contentType, fileName)
			if flag {
				fileUrl := fmt.Sprintf("https://%v/%v/%v", core.OssConfig.Endpoint, core.OssConfig.BucketName, fileName)
				logger.Log.Infof("图片保存成功，图片链接: %v", fileUrl)
			} else {
				logger.Log.Error("图片保存失败")
			}
		}
	}
	ctx.Next()
}
