package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

type Route struct {
	URL     string
	Method  string
	handler func(writer http.ResponseWriter, request *http.Request)
}

type Middleware func(ctx *Context) HttpError

type Handler[T any] func(ctx *Context, request T) (HttpResponse, HttpError)

/*
* Register api: register api to routeMap
* @param url: url of api
* @param handler: handler of api
* @param middleware: middleware of api
* @return void
 */
func RegisterAPI[T any](url string, method string, handler Handler[T], middlewares ...Middleware) {
	LoggerInstance.Info("Register api: %s %s", method, url)
	// Create a new handler
	h := func(writer http.ResponseWriter, request *http.Request) {
		// Create a new context
		ctx := getContext()
		defer putContext(ctx)
		buildContext(ctx, writer, request)

		// Append to common middleware
		middlewareList := []Middleware{}
		middlewareList = append(middlewareList, commonMiddlewares...)
		middlewareList = append(middlewareList, middlewares...)

		// Call middleware of function
		for _, middleware := range middlewareList {
			ctx.isRequestEnd = true
			if err := middleware(ctx); ctx.isRequestEnd {
				if err != nil {
					ctx.writeError(err)
				}
				return
			}
		}

		// Unmarshal json request body to model T
		req := initRequest[T]()
		requestContentType := strings.ToLower(ctx.GetRequestHeader(CONTENT_TYPE_KEY))
		if len(ctx.requestBody) != 0 {
			if strings.Contains(requestContentType, JSON_CONTENT_TYPE) {
				if err := json.Unmarshal(ctx.requestBody, &req); err != nil {
					LoggerInstance.Info("Unmarshal request body fail. RequestId: %s, Error: %s", ctx.requestID, err.Error())
					ctx.writeError(NewDefaultHttpError(400, "Bad request (Marshal requeset body)"))
					return
				}
			} else if strings.Contains(requestContentType, FORMDATA_CONTENT_TYPE) {
				buffer := bytes.NewBuffer(ctx.requestBody)
				ctx.request.Body = io.NopCloser(buffer)
				ctx.request.ParseForm()
			}
		}

		// Validate go struct with tag
		errValidate := validate.Struct(req)
		if errValidate != nil {
			errMessage := "Request invalid: "
			for _, err := range errValidate.(validator.ValidationErrors) {
				errMessage = fmt.Sprintf("%s {Field: %s, Tag: %s, Value: %s}", errMessage, err.Field(), err.Tag(), err.Value())
			}
			ctx.writeError(NewHttpError(http.StatusBadRequest, ERROR_BAD_BODY_REQUEST, errMessage, nil))
			return
		}

		// Call handler
		ctx.LogInfo("Request: Url = %s, body = %+v", ctx.URL, req)
		res, err := handler(ctx, req)
		if err != nil {
			ctx.LogError("Response error: Url = %s, body = %s", ctx.URL, err.Error())
			ctx.writeError(err)
			return
		}

		if res != nil {
			ctx.LogInfo("Response: Url = %s, body = %+v", ctx.URL, res.GetBody())
			ctx.writeSuccess(res)
		}
	}

	routeSlice, ok := routeMap[url]
	if ok {
		routeSlice = append(routeSlice, Route{
			Method:  method,
			URL:     url,
			handler: h,
		})
		routeMap[url] = routeSlice
	} else {
		routeMap[url] = []Route{
			{
				Method:  method,
				URL:     url,
				handler: h,
			},
		}
	}

}

func initRequest[T any]() T {
	var request T
	ref := reflect.New(reflect.TypeOf(request).Elem())
	return ref.Interface().(T)
}

func buildContext(ctx *Context, writer http.ResponseWriter, request *http.Request) HttpError {
	// Assign response writer and request
	ctx.rw = writer
	ctx.request = request

	ctx.requestID = ID.GenerateID()
	// Get url
	ctx.URL = request.URL.Path
	ctx.Method = request.Method

	// Get request body
	buffer := bytes.NewBuffer(ctx.requestBody)
	buffer.Reset()

	if _, err := io.Copy(buffer, request.Body); err != nil {
		LoggerInstance.Error("Read request body fail. RequestId: %s, Error: %s", ctx.requestID, err.Error())
		return HTTP_ERROR_READ_BODY_REQUEST_FAIL
	}

	if err := request.Body.Close(); err != nil {
		LoggerInstance.Error("Close request body fail. RequestId: %s, Error: %s", ctx.requestID, err.Error())
		return HTTP_ERROR_CLOSE_BODY_REQUEST_FAIL
	}

	ctx.requestBody = buffer.Bytes()
	return nil
}
