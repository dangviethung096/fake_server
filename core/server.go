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
