package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 自定义Handler
func Handler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&Ctx{c})
	}
}

// Success 成功返回
func (c *Ctx) Success(data ...interface{}) {
	var d interface{}
	d = ""
	if len(data) > 0 {
		d = data[0]
	}

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: d,
	})
}

// Error 失败返回
func (c *Ctx) Error(msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  msg,
	})
}

// ErrorWithCode 失败返回带错误码
func (c *Ctx) ErrorWithCode(code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}
