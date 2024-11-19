package internal

import (
	"net/http"
)

// RouterGroup 定义路由组
type RouterGroup struct {
	prefix      string
	router      *Router
	middlewares []Middleware
}

type IRouter interface {
	IRoutes
	Group(prefix string) *RouterGroup
}

type IRoutes interface {
	Use(mw ...Middleware) IRoutes
	Get(path string, handlers ...HandlerFunc) IRoutes
	Post(path string, handlers ...HandlerFunc) IRoutes
	Put(path string, handlers ...HandlerFunc) IRoutes
	Delete(path string, handlers ...HandlerFunc) IRoutes
}

var _ IRouter = (*RouterGroup)(nil)

// Group 创建一个路由组
func (r *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		prefix:      r.prefix + prefix, // 修正这里，将前缀拼接
		router:      r.router,
		middlewares: r.middlewares,
	}
}

// Use 注册中间件到路由组
func (g *RouterGroup) Use(mw ...Middleware) IRoutes {
	g.middlewares = append(g.middlewares, mw...)
	return g
}

// AddRoute 添加路由并将中间件应用到处理函数
func (g *RouterGroup) AddRoute(method, path string, handler HandlerFunc) {
	// 拼接路由的完整路径
	fullPath := g.prefix + path

	// 合并父路由和当前组的中间件
	allMiddlewares := append(g.router.middlewares, g.middlewares...)

	// 依次应用中间件
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		handler = wrapWithMiddleware(allMiddlewares[i], handler)
	}

	// 将处理器添加到路由表
	g.router.AddRoute(method, fullPath, handler)
}

// wrapWithMiddleware 包装中间件，将 HandlerFunc 转换为 http.Handler
func wrapWithMiddleware(mw Middleware, handler HandlerFunc) HandlerFunc {
	return func(c *Context) {
		// 执行中间件，并将 handler 作为参数传递进去
		mw(handler)(c)
	}
}

// 实现 IRoutes 接口的 Get 方法
func (g *RouterGroup) Get(path string, handlers ...HandlerFunc) IRoutes {
	// 定义一个用于执行传入处理函数的 handler
	handler := func(c *Context) {
		// 执行传入的处理函数
		for _, h := range handlers {
			h(c) // h 是 HandlerFunc，符合 HandlerFunc 的要求
		}
	}

	// 添加路由
	g.AddRoute(http.MethodGet, path, handler)
	return g
}

// 实现 IRoutes 接口的 Post 方法
func (g *RouterGroup) Post(path string, handlers ...HandlerFunc) IRoutes {
	// 定义一个用于执行传入处理函数的 handler
	handler := func(c *Context) {
		// 执行传入的处理函数
		for _, h := range handlers {
			h(c) // h 是 HandlerFunc，符合 HandlerFunc 的要求
		}
	}

	// 添加路由
	g.AddRoute(http.MethodPost, path, handler)
	return g
}

// 实现 IRoutes 接口的 Put 方法
func (g *RouterGroup) Put(path string, handlers ...HandlerFunc) IRoutes {
	// 定义一个用于执行传入处理函数的 handler
	handler := func(c *Context) {
		// 执行传入的处理函数
		for _, h := range handlers {
			h(c) // h 是 HandlerFunc，符合 HandlerFunc 的要求
		}
	}

	// 添加路由
	g.AddRoute(http.MethodPut, path, handler)
	return g
}

// 实现 IRoutes 接口的 Delete 方法
func (g *RouterGroup) Delete(path string, handlers ...HandlerFunc) IRoutes {
	// 定义一个用于执行传入处理函数的 handler
	handler := func(c *Context) {
		// 执行传入的处理函数
		for _, h := range handlers {
			h(c) // h 是 HandlerFunc，符合 HandlerFunc 的要求
		}
	}

	// 添加路由
	g.AddRoute(http.MethodDelete, path, handler)
	return g
}
