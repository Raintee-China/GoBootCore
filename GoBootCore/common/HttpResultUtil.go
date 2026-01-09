package common

import "github.com/gin-gonic/gin"

type HttpResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Success(data interface{}, msg string) HttpResult {
	if msg == "" {
		msg = "操作成功"
	}
	return HttpResult{
		Code: 200,
		Data: data,
		Msg:  msg,
	}
}

func Error(code int, msg string) HttpResult {
	if code == 0 {
		code = 500
	}
	return HttpResult{
		Code: code,
		Data: nil,
		Msg:  msg,
	}
}

func Fail(msg string) HttpResult {
	return HttpResult{
		Code: 500,
		Data: nil,
		Msg:  msg,
	}
}

func (r HttpResult) Send(c *gin.Context) {
	c.JSON(r.Code, r)
}