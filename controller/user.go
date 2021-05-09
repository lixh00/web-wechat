package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"web-wechat/core"
	"web-wechat/global"
	"web-wechat/protocol"
)

// 返回用户信息包装类
type responseUserInfo struct {
	Uin int `json:"uin"`
	// 性别
	Sex int `json:"sex"`
	// 省
	Province string `json:"province"`
	// 市
	City string `json:"city"`
	// 别名
	Alias string `json:"alias"`
	// 显示名称
	DisplayName string `json:"display_name"`
	// 昵称
	NickName string `json:"nick_name"`
	// 备注
	RemarkName string `json:"remark_name"`
	// 头像
	HeadImgUrl string `json:"head_img_url"`
	// 当前登录中用户的唯一标识
	UserName string `json:"user_name"`
}

// 返回的好友列表的实体
type friendsResponse struct {
	Count   int                `json:"count"`
	Friends []responseUserInfo `json:"friends"`
	Groups  []responseUserInfo `json:"groups"`
}

// GetCurrentUserInfoHandle 获取当前登录用户
func GetCurrentUserInfoHandle(ctx *gin.Context) {
	// 预检AppKey
	appKey := ctx.Request.Header.Get("AppKey")
	if err := checkBot(appKey); err != nil {
		core.FailWithMessage("AppKey预检失败："+err.Error(), ctx)
		return
	}
	bot := global.GetBot(appKey)
	// 获取登录用户信息
	user, err := bot.GetCurrentUser()
	if err != nil {
		core.FailWithMessage("获取登录用户信息失败", ctx)
		return
	}

	log.Println("登录用户：", user.NickName)
	// TODO 这儿的返回数据后面改成struct
	core.OkWithData(map[string]string{"nickName": user.NickName}, ctx)
}

// GetFriendsListHandle 获取好友列表
func GetFriendsListHandle(ctx *gin.Context) {
	// 预检AppKey
	appKey := ctx.Request.Header.Get("AppKey")
	if err := checkBot(appKey); err != nil {
		core.FailWithMessage("AppKey预检失败："+err.Error(), ctx)
		return
	}
	bot := global.GetBot(appKey)
	// 获取好友列表
	user, _ := bot.GetCurrentUser()
	friends, err := user.Friends(true)
	if err != nil {
		core.FailWithMessage("获取好友列表失败", ctx)
		return
	}

	groups, err := user.Groups(true)
	if err != nil {
		core.FailWithMessage("获取群聊列表失败", ctx)
		return
	}

	// 循环处理数据
	var friendList []responseUserInfo
	for _, friend := range friends {
		friendList = append(friendList, responseUserInfo{
			Uin:         friend.Uin,
			Sex:         friend.Sex,
			Province:    friend.Province,
			City:        friend.City,
			Alias:       friend.Alias,
			DisplayName: protocol.FormatEmoji(friend.DisplayName),
			NickName:    protocol.FormatEmoji(friend.NickName),
			RemarkName:  protocol.FormatEmoji(friend.RemarkName),
			HeadImgUrl:  friend.HeadImgUrl,
			UserName:    friend.UserName,
		})
	}

	// 循环处理数据
	var groupList []responseUserInfo
	for _, group := range groups {
		groupList = append(groupList, responseUserInfo{
			Uin:         group.Uin,
			Sex:         group.Sex,
			Province:    group.Province,
			City:        group.City,
			Alias:       group.Alias,
			DisplayName: protocol.FormatEmoji(group.DisplayName),
			NickName:    protocol.FormatEmoji(group.NickName),
			RemarkName:  protocol.FormatEmoji(group.RemarkName),
			HeadImgUrl:  group.HeadImgUrl,
			UserName:    group.UserName,
		})
	}

	// 返回给前端
	core.OkWithData(friendsResponse{Count: friends.Count(), Friends: friendList, Groups: groupList}, ctx)
}
