package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认路由引擎
	r := gin.Default()

	// 定义 GET 路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	// 启动 HTTP 服务（默认端口 8080）
	r.Run(":3500")
}
