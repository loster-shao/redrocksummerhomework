package gun

//处理错误
func (c *Context)Fail(code int, pattern string )  {
	c.JSON(code, S{"status": code, "message": pattern, "err" : "err"})
}
