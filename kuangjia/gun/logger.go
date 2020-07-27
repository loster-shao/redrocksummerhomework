package gun

import (
	"log"
	"time"
)

type loggerMessage struct {
	Millisecond int64  `json:"timestamp"`
	Body        string `json:"body"`
	Postion     string `json:"postion"`  
}


func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

//func (logger *Logger) Writer(lever string, msg string)  {
//	funcName := "null"
//	pc, file, line, ok := runtime.Caller(2)//跳过前俩
//	if !ok{
//		file = "null"
//		line = 0
//	} else {
//		funcName = runtime.FuncForPC(pc).Name()
//	}
//	_, filename := path.Split(file)
//
//	loggerMsg := &loggerMessage{
//		Millisecond: time.Now().UnixNano(),
//		Body:        msg,
//		Postion:     fmt.Sprintf(filename, line, funcName),
//	}
//	logger.OutPut(loggerMsg)
//}
//
//func (logger *Logger)OutPut(msg *loggerMessage)  {
//	fmt.Println(msg)
//}
