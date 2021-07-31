package handler

import (
	"encoding/xml"
	"testing"
)

func TestImageMessage(t *testing.T) {
	xmlStr := "<?xml version=\"1.0\"?>\n<msg>\n\t<img aeskey=\"4966c8d6596b9071edef84ae4d8c0886\" encryver=\"1\" cdnthumbaeskey=\"4966c8d6596b9071edef84ae4d8c0886\" cdnthumburl=\"3078020100046c306a020100020464d27d2a02032f54050204186a42b702046102691f044564316333393631323336383637616162636136366533353165373530656637315f61623739363464642d616363632d343263322d616139652d3331363566353038393235620204010800010201000405004c52ad00\" cdnthumblength=\"1814\" cdnthumbheight=\"0\" cdnthumbwidth=\"0\" cdnmidheight=\"0\" cdnmidwidth=\"0\" cdnhdheight=\"0\" cdnhdwidth=\"0\" cdnmidimgurl=\"3078020100046c306a020100020464d27d2a02032f54050204186a42b702046102691f044564316333393631323336383637616162636136366533353165373530656637315f61623739363464642d616363632d343263322d616139652d3331363566353038393235620204010800010201000405004c52ad00\" length=\"18377\" cdnbigimgurl=\"3078020100046c306a020100020464d27d2a02032f54050204186a42b702046102691f044564316333393631323336383637616162636136366533353165373530656637315f61623739363464642d616363632d343263322d616139652d3331363566353038393235620204010800010201000405004c52ad00\" hdlength=\"789750\" md5=\"811109bd6a97c37d8e60f1279f7a77c7\" />\n</msg>"
	t.Logf("原始内容:\n%v", xmlStr)
	// 解析xml为结构体
	var data ImageMessageData
	if err := xml.Unmarshal([]byte(xmlStr), &data); err != nil {
		t.Logf("消息解析失败: %v", err.Error())
	} else {
		t.Logf("解析出的AesKey：%v", data.Img.AesKey)
	}
}
