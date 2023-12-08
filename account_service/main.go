package main

import (
	"core"
	"fake_server/handlers"
	"flag"
	"net/http"
)

var configFile = flag.String("config", "core-config.yaml", "Core config file")

func main() {
	flag.Parse()
	core.Init(*configFile)
	defer core.Release()

	core.UserCorsMiddleware()

	core.RegisterAPI("/login", http.MethodPost, handlers.Login)
	core.RegisterAPI("/addAccount", http.MethodPost, handlers.AddAccount)
	core.RegisterAPI("/listAccount", http.MethodPost, handlers.ListAccount)
	core.RegisterAPI("/removeAccount", http.MethodPost, handlers.RemoveAccount)
	core.RegisterAPI("/removeAllAccount", http.MethodPost, handlers.RemoveAllAccount)

	core.Start()
}
