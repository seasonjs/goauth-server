package ginConfig

import (
	"github.com/gin-gonic/gin"
	"oauthServer/pkg/logger"
)

// Default 自定义的gin初始化配置
func Default() *gin.Engine {
	engine := gin.New()
	//仅仅是移除了logger中间件，以方便自定义中间件注入
	engine.Use(gin.Recovery())
	// 注入自定义logger
	logger.GinLog(engine)
	return engine
}
