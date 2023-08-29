package app

import (
	"context"
	"easy-signin/pkg/tools"
	"github.com/redis/go-redis/v9"
	"sync"
)

// App 用于存放 redis 客户端，上下文环境和时间状态对象
type App struct {
	RedisClient *redis.Client
	Ctx         context.Context
	TimeStatus  *tools.TimeStatus
}

var app App               // app对象
var appOnce = sync.Once{} // 确保app对象是单例

// GetApp 用于创建一个 App 对象，并初始化它的字段
func GetApp() *App {
	appOnce.Do(func() {
		app = App{
			RedisClient: initRedisClient(),
			Ctx:         initContext(),
			TimeStatus:  initTimeStatus(),
		}
	})
	return &app
}

// initRedisClient 用于初始化 redis 客户端
func initRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "justplan329",
	})
}

// initContext 用于初始化上下文环境
func initContext() context.Context {
	return context.Background()
}

// initTimeStatus 用于初始化时间状态对象
func initTimeStatus() *tools.TimeStatus {
	return tools.GetTimeStatus()
}
