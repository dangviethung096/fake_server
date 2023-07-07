package main

import (
	"fake_server/core"
	"fake_server/handlers"
	"flag"
	"net/http"
)

var port = flag.Int64("port", 10015, "Port open in server")

func main() {
	flag.Parse()
	core.InitServer()
	defer core.ReleaseServer()

	core.Handle(http.MethodPost, "/login", handlers.Login)
	core.Handle(http.MethodPost, "/addAccount", handlers.AddAccount)
	core.Handle(http.MethodPost, "/listAccount", handlers.ListAccount)
	core.Handle(http.MethodPost, "/removeAccount", handlers.RemoveAccount)
	core.Handle(http.MethodPost, "/removeAllAccount", handlers.RemoveAllAccount)

	core.Start(*port)
}
