package logger

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http/httputil"
	"os"
	"time"
)

// GinLog 默认配置方法
func GinLog(r *gin.Engine) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
	// 携带requestId
	r.Use(requestid.New())
	// 自定义格式
	r.Use(
		logger.SetLogger(
			logger.WithLogger(
				func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
					reqDump, _ := httputil.DumpRequest(c.Request, true)
					return log.Logger.With().
						//Str().
						Bytes("request", reqDump).
						Logger()
				},
			),
		),
	)
}
