package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator"
)

var sqliteSession dbSession
var LoggerInstance Logger
var routeMap map[string][]Route
var httpContextPool sync.Pool
var commonMiddlewares []Middleware
var Config CoreConfig
var redisClient cacheClient
var rabbitMQClient *messageQueue
var coreContext *Context
var validate *validator.Validate
var contextTimeout time.Duration

func Init(configFile string) {
	// Init core context
	coreContext = &Context{
		Context: context.Background(),
	}

	// Init config
	// Read config from file
	Config = loadConfigFile(configFile)
	if Config.Context.Timeout > 0 {
		contextTimeout = time.Second * time.Duration(Config.Context.Timeout)
	} else {
		// Default
		contextTimeout = time.Second * 60
	}

	// Set default if it is not config
	if Config.HttpClient.RetryTimes == 0 {
		Config.HttpClient.RetryTimes = 3
	}

	if Config.HttpClient.WaitTimes == 0 {
		Config.HttpClient.WaitTimes = 2000
	}

	LoggerInstance = initLogger()
	sqliteSession = openDBConnection(DBInfo{
		FilePath: Config.Database.FilePath,
	})

	// Init id generator
	initIdGenerator()
	// Core context will hold first id from instance
	coreContext.requestID = ID.GenerateID()

	routeMap = make(map[string][]Route)
	httpContextPool = sync.Pool{
		New: func() interface{} {
			return &Context{
				requestBody: make([]byte, 16384),
			}
		},
	}

	commonMiddlewares = make([]Middleware, 0)
	validate = validator.New()

}

/*
* Release: Release all resources
* @return void
 */
func Release() {
	closeLogger()
	closeDB()
	releaseCacheDB()
	releaseMessageQueue()
	stopScheduler()
}

func closeDB() {
	sqliteSession.Close()
}

func closeLogger() {
	LoggerInstance = nil
}

func releaseCacheDB() {
	redisClient.Close()
}

func releaseMessageQueue() {
	rabbitMQClient.connection.Close()
}

/*
* Start: Start server
* Register all routes and listen to port
* @return void
 */
func Start() {
	// Register all routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		routeList := routeMap[r.URL.Path]
		for _, route := range routeList {
			if route.Method == r.Method || r.Method == http.MethodOptions {
				route.handler(w, r)
				return
			}
		}
		http.NotFound(w, r)
	})

	// Listen and serve
	LoggerInstance.Info("Start server at port: %d", Config.Server.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Server.Port), nil)
	if err != nil {
		log.Fatalln("ListenAndServe fail: ", err)
	}
}

/*
* CacheClient: Get cache client
* @return cacheClient
 */
func CacheClient() cacheClient {
	return redisClient
}

/*
* MessageQueue: Get message queue client
* @return messageQueue
 */
func MessageQueue() *messageQueue {
	return rabbitMQClient
}

/*
* DBSession: Get database session
* @return dbSession
 */
func DBSession() dbSession {
	return sqliteSession
}
