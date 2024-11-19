package internal

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

// Context 封装 HTTP 请求与响应的上下文信息
type Context struct {
	Writer  http.ResponseWriter // HTTP 响应写入器
	Request *http.Request       // HTTP 请求
	Params  map[string]string   // 路由中解析的参数
	Status  int                 // 响应的 HTTP 状态码
}

// NewContext 创建新的 Context
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: req,
		Params:  make(map[string]string),
	}
}

// SetStatus 设置 HTTP 状态码
func (c *Context) SetStatus(code int) {
	c.Status = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// JSON 返回 JSON 格式的响应
func (c *Context) JSON(statusCode int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(statusCode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// HTML 返回 HTML 格式的响应
func (c *Context) HTML(statusCode int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(statusCode)
	c.Writer.Write([]byte(html))
}

// String 返回纯文本响应
func (c *Context) String(statusCode int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(statusCode)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Query 获取 URL 中的查询参数
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// PostForm 获取 POST 表单参数
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

// Param 获取路由参数
func (c *Context) Param(key string) string {
	if value, ok := c.Params[key]; ok {
		return value
	}
	return ""
}

// Redirect 实现重定向
func (c *Context) Redirect(statusCode int, location string) {
	http.Redirect(c.Writer, c.Request, location, statusCode)
}

// FormFile 处理文件上传
func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	file, fileHeader, err := c.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	file.Close()
	return fileHeader, nil
}

// Abort 提前终止处理
func (c *Context) Abort(statusCode int, msg string) {
	c.SetStatus(statusCode)
	c.Writer.Write([]byte(msg))
}
