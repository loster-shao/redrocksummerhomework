package gun

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

//新建路由
func newRouter() *router{
	//fmt.Println("make", make(map[string]HandlerFunc))
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

//获取路由
func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) //允许只有一个 *存在，匹配有没有*
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)//寻找节点数

	if n != nil {
		parts := parsePattern(n.pattern)//允许只有一个 *存在，找寻*
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

//添加路由
func(r *router) addRouter(method string, pattern string, handler HandlerFunc){
	//fmt.Println("handler", &handler)
	parts := parsePattern(pattern)  //允许只有一个 *存在，匹配有没有*
	log.Printf("Route %4s - %s", method, pattern)//打印日志
	key := method + "-" + pattern//方法+地址
	_, ok := r.roots[method]
	if !ok{
		r.roots[method] =&node{}
	}
	r.roots[method].insert(pattern, parts, 0)//插入新路由
	r.handlers[key] = handler
}


func (r *router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)//获取路由

	//里面就是点字符串拼接等基础操作
	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 Not Found！你访问的页面不存在%s\n", c.Path)
			//这里我在考虑使用啥返回，string还是json
		})
	}
	c.Next()
}

//允许只有一个 *存在
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")//通过/拆分

	parts := make([]string, 0)//创建map
	//for循环
	for _, item := range vs {
		//如果为不为空，计入map
		if item != "" {
			parts = append(parts, item)
			//第一个就为0的话可以直接退出
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}