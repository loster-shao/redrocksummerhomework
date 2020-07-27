package gun

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)


//基本写废了，别看了这个
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))

				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

//这个看网上的。。。。
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	//skip如果是0，返回当前调用Caller函数的函数名、文件、程序指针PC，1是上一层函数，以此类推
	//Callers 用来返回调用栈的程序计数器,
	// 第 0 个 Caller 是 Callers 本身，
	// 第 1 个是上一层 trace，
	// 第 2 个是再上一层的 defer func
	var str strings.Builder//声明string
	str.WriteString(message + "\nTraceback:")//WriteString向w中写入s的所有替换进行完后的拷贝。
	//Go语言虽然支持+=操作符来追加字符串，但更好的方式是使用bytes.Buffer，这种方式在节省内存和效率方面有更好的表现。
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)//FuncForPC返回一个表示调用栈标识符pc对应的调用栈的*Func；
		// 如果该调用栈标识符没有对应的调用栈，函数会返回nil。每一个调用栈必然是对某个函数的调用。
		file, line := fn.FileLine(pc)//FileLine返回该调用栈所调用的函数的源代码文件名和行号。
		// 如果pc不是f内的调用栈标识符，结果是不精确的
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
