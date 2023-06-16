package plugins

import (
	"context"
	"gitee.ltd/lxh/logger/log"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/eatmoreapple/openwechat"
	"net/http"
	"net/url"
	"strings"
	"time"
	"web-wechat/core"
)

// OpenGPT
// @description: 开放式GPT-3聊天机器人
// @receiver weChatPlugin
// @param ctx
func (weChatPlugin) OpenGPT(ctx *openwechat.MessageContext) {
	conf := core.SystemConfig.OpenAiConfig
	// 判断是否开启了GPT-3聊天机器人
	if !conf.Enable {
		return
	}

	msg := ctx.Content

	// 判断是不是引用消息
	//var beforeMsg string
	if strings.HasPrefix(ctx.Content, "「") && strings.Contains(ctx.Content, "\n- - - - - - - - - - - - - - -\n") {
		// 提取出前文
		//beforeMsg = strings.Split(ctx.Content, "\n- - - - - - - - - - - - - - -\n")[0]
		msg = strings.Split(ctx.Content, "\n- - - - - - - - - - - - - - -\n")[1]
	}

	// 取出消息第一行以及剩下的内容
	msgArray := strings.Split(msg, "\n")
	if len(msgArray) < 2 || strings.ToLower(msgArray[0]) != "@openai" {
		return
	}
	// 获取提问的内容
	question := strings.Join(msgArray[1:], "\n")

	log.Debugf("ChatGPT提问内容: %s", question)

	// 调用GPT-3聊天机器人
	// 如果配置了代理，就设置一下
	hc := http.Client{Timeout: 30 * time.Second}
	if conf.Proxy != "" {
		hc.Transport = &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(conf.Proxy)
			},
		}
	}

	// 创建OpenAI客户端
	client := gpt3.NewClient(conf.ApiKey, gpt3.WithHTTPClient(&hc))

	// 组装消息 TODO 懒得搞上下文联动，有想法的可以自己实现，只需要组装一下下面这个Message字段就行了，把之前的记录带过去
	request := gpt3.ChatCompletionRequest{
		Model:    gpt3.GPT3Dot5Turbo0301,
		Messages: []gpt3.ChatCompletionRequestMessage{{Role: "user", Content: question}},
	}
	// 调用聊天机器人
	resp, err := client.ChatCompletion(context.Background(), request)
	if err != nil {
		_, _ = ctx.ReplyText("ChatGPT AI 引擎出错了\n" + err.Error())
		return
	}
	log.Debugf("ChatGPT回答内容: %s", resp.Choices[0].Message.Content)
	// 发送回复
	_, _ = ctx.ReplyText(resp.Choices[0].Message.Content)
}
