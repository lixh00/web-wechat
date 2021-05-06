package core

import (
	"github.com/gofiber/fiber/v2"
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
	ERROR   = 7
	SUCCESS = 0
)

// Result 手动组装返回结果
func Result(code int, data interface{}, msg string, c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(Response{
		code,
		data,
		msg,
	})
}

// Ok 返回无数据的成功
func Ok(ctx *fiber.Ctx) error {
	return ctx.JSON(Response{
		SUCCESS,
		map[string]interface{}{},
		"操作成功",
	})
}

// OkWithMessage 返回自定义成功的消息
func OkWithMessage(message string, c *fiber.Ctx) error {
	return Result(SUCCESS, nil, message, c)
}

// OkWithData 自定义内容的成功返回
func OkWithData(data interface{}, c *fiber.Ctx) error {
	return Result(SUCCESS, data, "操作成功", c)
}

// OkDetailed 自定义消息和内容的成功返回
func OkDetailed(data interface{}, message string, c *fiber.Ctx) error {
	return Result(SUCCESS, data, message, c)
}

// Fail 返回默认失败
func Fail(c *fiber.Ctx) error {
	return Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

// FailWithMessage 返回自定义消息的失败
func FailWithMessage(message string, c *fiber.Ctx) error {
	return Result(ERROR, map[string]interface{}{}, message, c)
}
