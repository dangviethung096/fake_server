package core

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

var server Server

func InitServer() {
	ginEngine := gin.Default()
	ginEngine.Use(corsMiddleware())
	OpenDB()
	server = Server{
		engine: ginEngine,
	}
}

func Start(port int64) error {
	fmt.Println("Run server in port: ", port)
	return server.engine.Run(fmt.Sprintf(":%d", port))
}

func ReleaseServer() {
	CloseDB()
}

func Handle(method string, url string, handler ...gin.HandlerFunc) {
	server.engine.Handle(method, url, handler...)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
