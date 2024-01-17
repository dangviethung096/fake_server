package main

import (
	"fake_web/route"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", route.HomePage)
	fmt.Println("Server running on port 80")
	go http.ListenAndServe("0.0.0.0:80", nil)

	fmt.Println("Server running on port 443")
	err := http.ListenAndServeTLS("0.0.0.0:443", "./cert/cert.pem", "./cert/key.pem", nil)
	if err != nil {
		fmt.Printf("server failed to start: %s\n", err.Error())
	}
}
