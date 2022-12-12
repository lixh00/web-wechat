package plugins

import (
	"context"
	"gitee.ltd/lxh/logger/log"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/eatmoreapple/openwechat"
	"strings"
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
	var beforeMsg string
	if strings.HasPrefix(ctx.Content, "「") && strings.Contains(ctx.Content, "\n- - - - - - - - - - - - - - -\n") {
		// 提取出前文
		beforeMsg = strings.Split(ctx.Content, "\n- - - - - - - - - - - - - - -\n")[0]
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
	client := gpt3.NewClient(conf.ApiKey)
	var answer string
	var prompt []string
	if beforeMsg != "" && question == "继续" {
		// 如果引用了前文且新的指令为继续，就直接传入前文
		question = beforeMsg
		// 去掉一下两个括号
		question = strings.ReplaceAll(question, "「", "")
		question = strings.ReplaceAll(question, "」", "")
	}
	prompt = []string{question}

	err := client.CompletionStreamWithEngine(context.Background(), gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt:    prompt,
		MaxTokens: gpt3.IntPtr(512),
		Echo:      true,
	}, func(response *gpt3.CompletionResponse) {
		answer += response.Choices[0].Text
	})
	if err != nil {
		_, _ = ctx.ReplyText("ChatGPT AI 引擎出错了\n" + err.Error())
		return
	}
	log.Debugf("ChatGPT回答内容: %s", answer)
	// 发送回复
	if answer != "" {
		_, _ = ctx.ReplyText(answer)
	}
}
