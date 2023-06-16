package handler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"gitee.ltd/lxh/logger/log"
	"github.com/eatmoreapple/openwechat"
	"io"
	"net/http"
	"strings"
	"web-wechat/core"
	"web-wechat/oss"
)

// VideoMessageData 图片消息结构体
type VideoMessageData struct {
	XMlName  xml.Name `xml:"msg"`
	VideoMsg struct {
		AesKey            string `xml:"aeskey,attr"`
		CdnVideoUrl       string `xml:"cdnvideourl,attr"`
		CdnThumbAesKey    string `xml:"cdnthumbaeskey,attr"`
		CdnThumbUrl       string `xml:"cdnthumburl,attr"`
		Length            int64  `xml:"length,attr"`
		PlayLength        int64  `xml:"playlength,attr"`
		CdnThumbLength    int64  `xml:"cdnthumblength,attr"`
		CdnThumbWidth     int64  `xml:"cdnthumbwidth,attr"`
		CdnThumbHeight    int64  `xml:"cdnthumbheight,attr"`
		FromUserName      string `xml:"fromusername,attr"`
		Md5               string `xml:"md5,attr"`
		NewMd5            string `xml:"newmd5,attr"`
		IsPlaceHolder     int64  `xml:"isplaceholder,attr"`
		RawMd5            string `xml:"rawmd5,attr"`
		RawLength         int64  `xml:"rawlength,attr"`
		CdnRawVideoUrl    string `xml:"cdnrawvideourl,attr"`
		CdnRawVideoAesKey string `xml:"cdnrawvideoaeskey,attr"`
		OverwriteNewMsgId int64  `xml:"overwritenewmsgid,attr"`
		IsAd              int64  `xml:"isad,attr"`
	} `xml:"videomsg"`
}

func videoMessageHandle(ctx *openwechat.MessageContext) {
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}
	// 解析xml为结构体
	var data VideoMessageData
	if err := xml.Unmarshal([]byte(ctx.Content), &data); err != nil {
		log.Errorf("消息解析失败: %v", err.Error())
		log.Debugf("发信人: %v ==> 原始内容: %v", senderUser, ctx.Content)
		return
	}
	log.Infof("[收到新视频消息] == 发信人：%v", senderUser)
	fileResp, err := ctx.GetVideo()
	if err != nil {
		log.Errorf("视频下载失败: %v", err.Error())
		return
	}
	defer fileResp.Body.Close()
	videoFileByte, err := io.ReadAll(fileResp.Body)
	if err != nil {
		log.Errorf("视频读取错误: %v", err.Error())
		return
	} else {
		// 读取文件相关信息
		contentType := http.DetectContentType(videoFileByte)
		fileType := strings.Split(contentType, "/")[1]
		fileName := fmt.Sprintf("%v.%v", ctx.MsgId, fileType)
		if user, err := ctx.Bot().GetCurrentUser(); err == nil {
			uin := user.Uin
			fileName = fmt.Sprintf("%v/%v", uin, fileName)
		}

		// 上传文件
		reader2 := io.NopCloser(bytes.NewReader(videoFileByte))
		flag := oss.SaveToOss(reader2, contentType, fileName)
		if flag {
			fileUrl := fmt.Sprintf("http://%v/%v/%v", core.SystemConfig.OssConfig.Endpoint, core.SystemConfig.OssConfig.BucketName, fileName)
			log.Infof("视频保存成功，视频链接: %v", fileUrl)
			ctx.Content = fileUrl
		} else {
			log.Error("视频保存失败")
		}
	}
	ctx.Next()
}
