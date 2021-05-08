package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 返回数据包装
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// 定义状态码
const (
	ERROR   = 0
	SUCCESS = 1
)

// Result 手动组装返回结果
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

// Ok 返回无数据的成功
func Ok(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		SUCCESS,
		map[string]interface{}{},
		"操作成功",
	})
}

// OkWithMessage 返回自定义成功的消息
func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, nil, message, c)
}

// OkWithData 自定义内容的成功返回
func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

// OkDetailed 自定义消息和内容的成功返回
func OkDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

// Fail 返回默认失败
func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

// FailWithMessage 返回自定义消息的失败
func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}
