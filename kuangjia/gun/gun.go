////实现ServeHTTP方法 根据请求的方法及路径来匹配Handler
//func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	key := r.Method + "-" + r.URL.Path
//	if handler, ok := engine.router[key]; ok {
//		handler(w, r)
//	} else {
//		fmt.Fprintf(w, "404 NOT FOUND %s ", r.URL.Path)
//	}
//}
package gun

import (
	"log"
	"net/http"
	"strings"
)

type RouterGroup struct {
	prefix      string        //前缀  举例： /s/sz/szs前缀就是/s/sz
	middlewares []HandlerFunc //支持中间键
	parent      *RouterGroup  //支持嵌套
	engine      *Engine       //所有路由共享一个引擎
}

//定义函数必须为HandlerFunc类型（就是你写POST("/szs",SZS)这个SZS函数必须用这个函数类型）
type HandlerFunc func(c *Context)

//定义Engine结构体
type Engine struct {
	*RouterGroup         //路由组
	router *router       //路由
	groups []*RouterGroup//路由组嵌套（和去年头疼死的消息嵌套一样）
}

//使用中间键
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//实现ServeHTTP方法 根据请求的方法及路径来匹配Handler
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc//声明
	//遍历group
	for _, group := range engine.groups {
		//如果req.URL.Path前缀含有group.prefix执行（即判断/s/ss前缀有没有/s）
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)//添加中间键
		}
	}
	c := NewContext(w, req)//构造函数，写入参数
	c.handlers = middlewares
	engine.router.handle(c)
}

//启动服务
func (engine *Engine) Run(addr string) (err error) {
	log.Printf("Listening and serving HTTP on %s\n", addr)
	return http.ListenAndServe(addr, engine)
}

//外部调用框架入口
// 说是New()加上中间键。。。不太会用，
//不用中间键其实也没啥问题
func Default() *Engine {
	//不太会用中间键。。。。。。
	//反正不用中间键也是能用的。。。。
	engine := New()
	engine.Use(Logger(),Recovery())//使用中间键部分（我鸽了）
	return engine
}

//外部调用框架入口
func New() *Engine {
	//engine赋值
	engine := &Engine{
		router: newRouter(),//新建路由
	}

	//获取引擎
	engine.RouterGroup =&RouterGroup{
		engine:      engine,
	}

	//添加路由组
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//路由组写法(写炸了,好像又没写炸)，好像也能用
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine//获取引擎
	//赋值
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,//前缀+后缀  举例： /s/sz/szs前缀就是/s/sz， 这个就是/s/sz/szs
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)//添加map
	//fmt.Println("groups", engine.groups, "new", *newGroup)
	return newGroup
}

//框架新增路由
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp//前缀+后缀  举例： /s/sz/szs前缀就是/s/sz，这个就是/s/sz/szs
	//fmt.Println("5", group.prefix)
	//log.Printf("Route %s - %s", method, pattern)//打印日志
	group.engine.router.addRouter(method, pattern, handler)//添加路由
}

//匹配get方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)//添加路由
}

//匹配post方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)//添加路由
}
