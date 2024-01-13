package main

import (
	"fake_web/route"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", route.HomePage)
	fmt.Println("Server running on port 80")
	http.ListenAndServe(":80", nil)
}
