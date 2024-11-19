package middleware

import (
	"hs/internal"
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery 是一个中间件，用于捕获 panic 并记录日志，同时返回错误响应
func Recovery(next internal.HandlerFunc) internal.HandlerFunc {
	return internal.HandlerFunc(func(c *internal.Context) {
		defer func() {
			// 捕获 panic
			if err := recover(); err != nil {
				// 记录错误日志
				log.Printf(
					"panic recovered: %v\n%s",
					err,
					string(debug.Stack()), // 打印堆栈信息
				)

				// 返回 500 错误响应
				c.Writer.WriteHeader(http.StatusInternalServerError)
				c.Writer.Write([]byte(`{"error": "Internal Server Error", "message": "Something went wrong!"}`))
			}
		}()

		// 执行下一个处理器
		next(c)
	})
}

// JSONRecovery 是一个中间件，用于捕获 panic 并返回 JSON 格式的错误信息
func JSONRecovery(next internal.HandlerFunc) internal.HandlerFunc {
	return internal.HandlerFunc(func(c *internal.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v\n", err)

				// 返回 JSON 格式错误信息
				c.Writer.Header().Set("Content-Type", "application/json")
				c.Writer.WriteHeader(http.StatusInternalServerError)
				c.Writer.Write([]byte(`{"error": "Internal Server Error", "message": "Something went wrong!"}`))
			}
		}()
		next(c) // 调用处理函数
	})
}
