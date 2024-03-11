package middleware

import (
	"net/http"
	"qp/qConfig"
	"qp/qLog"

	"github.com/gin-contrib/pprof"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func InitMiddleware(r *gin.Engine) {
	r.Use(Cors())
	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	// 日志处理
	r.Use(LoggerToFile())
	// Set X-Request-Id header
	r.Use(RequestId())
	r.Use(panicApi)
	if qConfig.Settings.Application.Mode == "dev" || qConfig.Settings.Application.Mode == "local" {
		pprof.Register(r)
	}
	p := NewPrometheus("ac")
	p.Use(r)
	r.Use(otelgin.Middleware(qConfig.Settings.Application.Name))
}

func panicApi(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			qLog.TraceError(
				c.Request.Context(),
				"HttpPanic",
				zap.Any("panic", r),
				zap.Any("url", c.Request.URL),
				zap.Stack("stack"),
			)
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    500,
					"message": "Internal Server Error",
					"success": false,
				},
			)
		}
	}()
	c.Next()
}

func authMiddleware(c *gin.Context) {
	if false {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"code":    401,
				"message": "Invalid token",
				"success": false,
			},
		)
	}
	c.Next()
}
