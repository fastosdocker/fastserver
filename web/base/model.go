package base

import "github.com/gin-gonic/gin"

// Ctx 自定义上下文
type Ctx struct {
	*gin.Context
}

// HandlerFunc 自定义handler函数
type HandlerFunc func(c *Ctx)

// Response 返回客户端的结构体
type Response struct {
	Code int
	Msg  string
	Data interface{}
}
