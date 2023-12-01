package core

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

/*
* Info: Log Info with context information
* @params: format string, args ...interface{}
* @return: void
 */
func (ctx *Context) LogInfo(format string, args ...interface{}) {
	logStr := "[INFO] " + ctx.format(format, args...)
	log.Println(logStr)
}

/*
* Debug: Log Debug with context information
* @params: format string, args ...interface{}
* @return: void
 */
func (ctx *Context) LogDebug(format string, args ...interface{}) {
	logStr := "[DEBUG] " + ctx.format(format, args...)
	log.Println(logStr)
}

/*
* Error: Log Error with context information
* @params: format string, args ...interface{}
* @return: void
 */
func (ctx *Context) LogError(format string, args ...interface{}) {
	logStr := "[ERROR] " + ctx.format(format, args...)
	log.Println(logStr)
}

/*
* Warning: Log Warning with context information
* @params: format string, args ...interface{}
* @return: void
 */
func (ctx *Context) LogWarning(format string, args ...interface{}) {
	logStr := "[WARNING] " + ctx.format(format, args...)
	log.Println(logStr)
}

/*
* Fatal: Log Fatal with context information
* @params: format string, args ...interface{}
* @return: void
 */
func (ctx *Context) LogFatal(format string, args ...interface{}) {
	logStr := "[FATAL] " + ctx.format(format, args...)
	log.Fatalln(logStr)
}

/*
* format: Format the ctx: add to message log the file name
* and line number of the code that calls the ctx interface
* @params: format string, args ...interface{}
* @return: string
 */
func (ctx *Context) format(format string, args ...interface{}) string {
	// Format the ctx
	logStr := fmt.Sprintf(format, args...)

	// Get the file name and line number of the code that calls the ctx interface
	_, file, line, ok := runtime.Caller(2)
	if ok {
		path := strings.Split(file, "/")
		if len(path) > 3 {
			file = strings.Join(path[len(path)-3:], "/")
		}
		logStr = fmt.Sprintf("%s:%d: RequestID: %s, Message: %s", file, line, ctx.requestID, logStr)
	}

	// Return the formatted string
	return logStr
}
