package internal

import (
	"fmt"
	"log"
	"net/http"
)

// Middleware 定义中间件类型
type Middleware func(HandlerFunc) HandlerFunc

// Router 定义路由器
type Router struct {
	RouterGroup
	middlewares []Middleware           // 中间件队列
	routes      map[string]HandlerFunc // 路由表，使用 HandlerFunc
}

// NewRouter 创建并返回一个新的 Router 实例
func NewRouter() *Router {
	return &Router{
		middlewares: []Middleware{},
		routes:      make(map[string]HandlerFunc),
	}
}

// Use 注册一个或多个中间件
func (r *Router) Use(mw ...Middleware) *Router {
	r.middlewares = append(r.middlewares, mw...)
	return r
}

// AddRoute 添加路由，并将所有中间件应用到处理器
func (r *Router) AddRoute(method, path string, handler HandlerFunc) {
	// 应用中间件
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	// 将处理器添加到路由表
	r.routes[method+path] = handler
}

// ServeHTTP 处理 HTTP 请求并根据路由表分发
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 根据方法和路径生成路由键
	key := req.Method + req.URL.Path
	if handler, ok := r.routes[key]; ok {
		// 创建 Context
		c := &Context{Writer: w, Request: req}

		// 执行路由处理函数
		handler(c)
		return
	}
	// 如果路由没有找到，则返回 404
	http.NotFound(w, req)
}

// Start 启动 HTTP 服务
func (r *Router) Start(addr string) {
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Printf("Starting server on %s...\n", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
