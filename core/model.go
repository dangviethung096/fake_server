package core

type Request struct {
	Method int
	Body   string
	URL    string
}

type Response struct {
	StatusCode int
	Body       interface{}
}

type RequestProperty struct {
}

type Req interface{}
