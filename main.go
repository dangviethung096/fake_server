package main

import (
	"fake_server/core"
	"fake_server/handlers"
	"net/http"
)

func main() {
	core.InitServer()
	defer core.ReleaseServer()

	core.Handle(http.MethodPost, "/login", handlers.Login)
	core.Handle(http.MethodPost, "/addAccount", handlers.AddAccount)
	core.Handle(http.MethodPost, "/listAccount", handlers.ListAccount)

	core.Start(10015)
}
