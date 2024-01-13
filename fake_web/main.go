package main

import (
	"fake_web/route"
	"net/http"
)

func main() {
	http.HandleFunc("/", route.HomePage)
	http.ListenAndServe(":80", nil)
}
