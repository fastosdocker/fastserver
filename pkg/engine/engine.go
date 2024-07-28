package engine

import (
	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine
)

func init() {
	engine = gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
}

func Get() *gin.Engine {
	return engine
}
