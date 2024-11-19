package middleware

import (
	"hs/internal"
	"log"
	"net/http"
	"time"
)

// Logger 中间件，签名应为 func(internal.HandlerFunc) internal.HandlerFunc
func Logger(next internal.HandlerFunc) internal.HandlerFunc {
	return func(c *internal.Context) {
		start := time.Now()

		// 包装 ResponseWriter，用于捕获响应状态码
		lrw := &loggingResponseWriter{ResponseWriter: c.Writer, statusCode: http.StatusOK}

		// 将日志写入
		next(c) // 调用处理函数

		// 记录日志信息
		duration := time.Since(start)
		log.Printf(
			"%s - [%s] \"%s %s %s\" %d %s",
			c.Request.RemoteAddr,
			start.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			lrw.statusCode,
			duration,
		)
	}
}

// loggingResponseWriter 包装 http.ResponseWriter 用于捕获状态码
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 捕获状态码
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
