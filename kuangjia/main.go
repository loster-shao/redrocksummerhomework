package main

import (
	"fmt"
	"gopkg.in/olahol/melody.v1"

	"kuangjia/gun"
	"net/http"
)

var M *melody.Melody

func main()  {
	c := gun.Default()

	g1 := c.Group("/szs")
	{
		g1.POST("/json", Jsons)
		g1.GET("/html", Html)
		g1.GET("/string",String)
		g2 := g1.Group("/szs")
		{
			g2.GET("/666", Html)
		}
	}
	//websocket还没完工，不能用
	c.GET("/s", Web)
	c.GET("/ws", Websockets)
	c.Run(":8080")
}

func Jsons(c *gun.Context) {
	szs := c.PostForm("szs")
	fmt.Println(szs)
	c.JSON(200, gun.S{
		"status"  : http.StatusOK,
		"message" : szs,
	})

}

func String(c *gun.Context)  {
	szs := c.Query("szs")
	c.String(200, szs)
}
//HTML静态页面
func Html(c *gun.Context){
	c.HTML(http.StatusOK, `<!doctype html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport"
	content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Hello World</title>
	<style>
	.div {
	width: 200px;
		font-size: 40px;
	color: red;
	margin: 0 auto;
	}
	</style>
	</head>
	<body>
	<div class="div">Hello world</div>
	</body>
	</html>"`)

}

//发送json消息


//发送静态页面
func Web(c *gun.Context)  {
	http.ServeFile(c.Writer, c.Request, "index.html")
}

//ws连接使用（还没写完）
func Websockets(c *gun.Context)  {
	M.HandleRequest(c.Writer, c.Request)
}
