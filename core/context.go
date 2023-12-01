package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

/*
* Context type: which carries deadlines, cancellation signals,
* and other request-scoped values across API boundaries and between processes.
 */
type Context struct {
	context.Context
	cancelFunc    context.CancelFunc
	URL           string
	Method        string
	Timeout       time.Duration
	requestID     string
	requestBody   []byte
	responseBody  responseBody
	isRequestEnd  bool
	request       *http.Request
	rw            http.ResponseWriter
	isResponseEnd bool
}

/*
* GetContext: Get context from pool
* @return: *Context
 */
func getContext() *Context {
	ctx := httpContextPool.Get().(*Context)
	ctx.Context, ctx.cancelFunc = context.WithTimeout(coreContext, contextTimeout)
	ctx.Timeout = contextTimeout
	ctx.isResponseEnd = false
	return ctx
}

/*
* PutContext: Put context to pool
* @params: *Context
* @return: void
 */
func putContext(ctx *Context) {
	ctx.cancelFunc()
	httpContextPool.Put(ctx)
}

/*
* Next: Set isRequestEnd to false
* This funciton must to be called when you want to call next middleware
* @return: void
 */
func (ctx *Context) Next() {
	ctx.isRequestEnd = false
}

/*
* GetRequestHeader: Get request header by key
* @params: key string
* @return: string
 */
func (ctx *Context) GetRequestHeader(key string) string {
	return ctx.request.Header.Get(key)
}

/*
* GetQueryParam: Get query param by key
* @params: key string
* @return: string
 */
func (ctx *Context) GetQueryParam(key string) string {
	return ctx.request.URL.Query().Get(key)
}

/*
* ListQueryParam: Get list query param by key
* @params: key string
* @return: []string
 */
func (ctx *Context) GetArrayQueryParam(key string) []string {
	return ctx.request.URL.Query()[key]
}

/*
* GetFormData: get data in body when context/type is application/x-www-form-urlencoded
* in header request
* @return string
* if key exist in form data return value of key, otherwise return empty string
 */
func (ctx *Context) GetFormData(key string) string {
	return ctx.request.PostForm.Get(key)
}

/*
* Redirect url
 */
func (ctx *Context) RedirectURL(url string) {
	ctx.isResponseEnd = true
	http.Redirect(ctx.rw, ctx.request, url, http.StatusSeeOther)
}

/*
* writeError: write error http response to user
 */
func (ctx *Context) writeError(httpErr HttpError) {
	ctx.rw.Header().Set("Content-Type", "application/json")
	ctx.rw.Header().Set("Request-Id", ctx.requestID)
	ctx.responseBody.Code = httpErr.GetCode()
	ctx.responseBody.Message = httpErr.GetMessage()
	ctx.responseBody.Data = httpErr.GetErrorData()

	body, err := json.Marshal(ctx.responseBody)
	if err != nil {
		ctx.LogError("Marshal error json. RequestId: %s, Error: %s", ctx.requestID, err.Error())
		ctx.endResponse(http.StatusInternalServerError, `{"code":500,"message":"Internal server error(Marshal error response data)","errorData":null,"data":null}`)
		return
	}

	ctx.endResponse(int(httpErr.GetStatusCode()), string(body))
}

/*
* writeSuccess: write success http response to user
 */
func (ctx *Context) writeSuccess(httpRes HttpResponse) {
	ctx.rw.Header().Set("Content-Type", "application/json")
	ctx.rw.Header().Set("Request-Id", ctx.requestID)
	ctx.responseBody.Code = httpRes.GetReponseCode()
	ctx.responseBody.Message = BLANK
	ctx.responseBody.Data = httpRes.GetBody()

	body, err := json.Marshal(ctx.responseBody)
	if err != nil {
		ctx.LogError("Marshal json. RequestId: %s, Error: %s", ctx.requestID, err.Error())
		ctx.endResponse(http.StatusInternalServerError, `{"code":500,"message":"Internal server error(Marshal response data)","errorData":null,"data":null}`)
		return
	}

	ctx.endResponse(int(httpRes.GetStatusCode()), string(body))
}

/*
* endResponse: call write header if it is not called before and write body to writer
 */
func (ctx *Context) endResponse(statusCode int, body string) {
	if !ctx.isResponseEnd {
		ctx.isResponseEnd = true
		// end response
		ctx.rw.WriteHeader(statusCode)
		fmt.Fprint(ctx.rw, body)
	}
}

/*
* GetContextForTest: Get context for test
* Caution: This function is only used for test
* @return: *Context
 */
func GetContextForTest() *Context {
	ctx := httpContextPool.Get().(*Context)
	ctx.Context, ctx.cancelFunc = context.WithTimeout(coreContext, contextTimeout)
	ctx.requestID = ID.GenerateID()
	return ctx
}

/*
* Get child of core context with timeout as a parameter
 */
func GetContextWithTimeout(timeout time.Duration) *Context {
	ctx := httpContextPool.Get().(*Context)
	ctx.Context, ctx.cancelFunc = context.WithTimeout(coreContext, timeout)
	ctx.Timeout = timeout
	ctx.requestID = ID.GenerateID()
	ctx.URL = "GetContextWithTimeout"
	return ctx
}

/*
* Return context to http context pool
 */
func PutContext(ctx *Context) {
	ctx.cancelFunc()
	httpContextPool.Put(ctx)
}
