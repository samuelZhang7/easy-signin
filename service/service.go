package service

import (
	"context"
	"easy-signin/pkg/app"
	"easy-signin/pkg/tools"
	"github.com/redis/go-redis/v9"
)

type Service struct{}

var redisClient *redis.Client    // service包内的redis客户端，用于执行 Redis 操作
var timeStatus *tools.TimeStatus // service包内的时间状态对象，用于处理时间相关的信息
var ctx context.Context          // service包内的上下文环境，用于传递请求范围的数据。

// GetService 获取一个已初始化的 Service 对象。
func GetService(app *app.App) *Service {
	// 从 App 对象中获取必要的信息
	appClient := app
	redisClient = appClient.RedisClient // 设置 Redis 客户端
	timeStatus = appClient.TimeStatus   // 设置时间状态对象
	ctx = appClient.Ctx                 // 设置上下文环境
	return &Service{}                   // 返回初始化后的 Service 实例
}
