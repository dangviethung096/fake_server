package core

type Error struct {
	Message string `json:"message"`
}

var (
	SERVER_ERROR = Response{StatusCode: 500, Body: Error{Message: "Server Error"}}
)
