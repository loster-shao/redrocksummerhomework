package gun

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type S map[string]interface{}

//引擎
type Context struct {
	Writer     http.ResponseWriter //http响应体(回复)
	Request    *http.Request       //http请求（请求）
	Path       string              //请求路径
	Method     string              //请求方法（Post,Get等）
	Params     map[string]string   //获取Handlerfunc()函数解析后的参数
	StatusCode int                 //状态码
	handlers   []HandlerFunc       //中间键
	index      int                 //中间键所需码（记录执行到第几个中间件）
	//index是记录当前执行到第几个中间件，
	// 当在中间件中调用Next方法时，
	// 控制权交给了下一个中间件，
	// 直到调用到最后一个中间件，
	// 然后再从后往前，调用每个中间件在Next方法之后定义的部分。
}

//构造函数，写入参数
func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Writer:  writer,            //http响应体
		Request: request,          //http请求写入
		Path:    request.URL.Path, //路径写入
		Method:  request.Method,   //请求方法
	}
}

//记录到上方中间件，在其他函数处理完后继续执行其中间件
func (c *Context) Next() {
	c.index ++
	for ; c.index < len(c.handlers); c.index ++ {
		c.handlers[c.index](c)
	}
}
//gin框架解析文档中有描述

//获取表单值
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
	//通过http.Request.FormValue这个函数取值，别问我这个函数原理，我只会用（才看的http）
}

//获取post请求头（url重值）
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
	//我没有太弄懂Query和PostForm的具体区别（应该就是一个是写在路由里面的，PostForm写在请求内容里的）。。。
}

//设置状态码（我从来没用过，我看gin框架里这个挺简单的。。。。）
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)//写请求码
}

//设置头信息（我从来没用过，我看gin框架里这个挺简单的。。。。）
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)//模仿浏览器设置头
}

//发送String格式响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//发送JSON数据
func (c *Context) JSON(code int, h interface{}) {
	c.SetHeader("Content-Type", "application/json")//设置请求头
	c.Status(code)//设置码
	encoder := json.NewEncoder(c.Writer)//json数据解码
	//fmt.Println("writer", c.Writer)
	if err := encoder.Encode(h)/*Json编码*/; err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//发送Data格式响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)//状态码
	c.Writer.Write(data)//发送byte
}

//发送HTML格式响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")//设置请求头
	c.Status(code)//状态码
	c.Writer.Write([]byte(html))//发送byte
}
