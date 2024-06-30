package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quincy0/qpro/qTrace"
)

// Error 失败数据处理
func Error(c *gin.Context, code int, err error) {
	var res Response[any]
	res.TraceId = qTrace.TraceIdFromContext(c.Request.Context())
	traceError(err)
	c.Set("err", code) // 给promotheus用
	c.JSON(http.StatusOK, res.ReturnError(code, err.Error()))
	c.Abort()
}

// Fail 失败数据处理
func Fail(c *gin.Context, code int, msg string) {
	var res Response[any]
	res.TraceId = qTrace.TraceIdFromContext(c.Request.Context())
	trace(msg)
	c.Set("err", code) // 给promotheus用
	c.JSON(http.StatusOK, res.ReturnError(code, msg))
}

// OK 通常成功数据处理
func OK[T any](c *gin.Context, data T) {
	var res Response[T]
	res.Data = data
	res.TraceId = qTrace.TraceIdFromContext(c.Request.Context())
	c.JSON(http.StatusOK, res.ReturnOK())
}

// Void 通常成功数据处理
func Void(c *gin.Context) {
	var res Response[any]
	c.JSON(http.StatusOK, res.ReturnOK())
}

// PageOK 分页数据处理
func PageOK[T any](c *gin.Context, result []T, count int64, pageIndex int, pageSize int) {
	var res Response[Page[T]]
	var page Page[T]
	page.List = result
	page.Count = count
	page.PageIndex = pageIndex
	page.PageSize = pageSize
	res.Data = page
	c.JSON(http.StatusOK, res.ReturnOK())
}

func Forbidden(c *gin.Context, code int, err error) {
	var res Response[any]
	res.TraceId = qTrace.TraceIdFromContext(c.Request.Context())
	traceError(err)
	c.Set("err", code)
	c.JSON(http.StatusForbidden, res.ReturnError(code, err.Error()))
	c.Abort()
}

func ServerError(c *gin.Context, code int, err error) {
	var res Response[any]
	res.TraceId = qTrace.TraceIdFromContext(c.Request.Context())
	traceError(err)
	c.Set("err", code)
	c.JSON(http.StatusInternalServerError, res.ReturnError(code, err.Error()))
	c.Abort()
}

func NotImplemented(c *gin.Context, code int, err error) {
	var res Response[any]
	res.TraceId = qTrace.TraceIdFromContext(c.Request.Context())
	traceError(err)
	c.Set("err", code)
	c.JSON(http.StatusNotImplemented, res.ReturnError(code, err.Error()))
	c.Abort()
}

func GatewayTimeout(c *gin.Context, code int, err error) {
	var res Response[any]
	res.TraceId = qTrace.TraceIdFromContext(c.Request.Context())
	traceError(err)
	c.Set("err", code)
	c.JSON(http.StatusGatewayTimeout, res.ReturnError(code, err.Error()))
	c.Abort()
}
