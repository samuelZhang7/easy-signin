package service

import (
	"easy-signin/pkg/tools"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// SignIn 用户签到
func (s *Service) SignIn(userID int64) error {
	// 1. 检查用户是否已签到
	signed, err := checkSignIn(userID)
	if err != nil {
		return err
	}
	if signed {
		fmt.Printf("用户[%d]已签到，请勿重复签到\n", userID)
		return nil
	}
	// 2. 签到并更新统计
	return s.updateSignInfo(userID)
}

// checkSignIn 检查是否已经签到
func checkSignIn(userID int64) (bool, error) {
	// 通过GET BIT判断对应日期是否已签到
	signed, err := redisClient.GetBit(ctx, tools.GenBitmapKey(userID, timeStatus.YearMonth), timeStatus.Day-1).Result()
	if err != nil {
		return false, err
	}
	// 已签到
	if signed != 0 {
		return true, nil
	}
	// 未签到
	return false, nil
}

// updateSignInfo 执行签到更新
func (s *Service) updateSignInfo(userID int64) error {

	// 初始化pipeline
	pipe := redisClient.Pipeline()

	// 设置签到位图
	pipe.SetBit(ctx, tools.GenBitmapKey(userID, timeStatus.YearMonth), timeStatus.Day-1, 1)
	// 更新签到总数
	pipe.Incr(ctx, tools.GenSignCountKey(userID))
	// 获取连续签到key
	continuousKey := tools.GenContinuousKey(userID)
	// 更新连续签到总数
	pipe.Incr(ctx, continuousKey)
	// 更新连续签到的过期时间
	pipe.ExpireAt(ctx, continuousKey, tools.NextExpireTime())
	// 执行
	if _, err := pipe.Exec(ctx); err != nil {
		return handleRedisError(err)
	}
	// 签到成功
	s.PrintUserSignInData(userID) // 打印用户签到信息
	// 签到记录入influxDB
	// 后续业务逻辑
	return nil
}

// 错误处理函数
func handleRedisError(err error) error {
	// 日志记录
	// 返回统一错误
	// 尝试回滚操作
	return fmt.Errorf("%w", err)
}

// PrintUserSignInData 打印用户登录数据，连续签到天数、累计签到天数以及签到表
func (s *Service) PrintUserSignInData(userID int64) {
	// 获取连续签到key
	continuousKey := tools.GenContinuousKey(userID)

	// 将多个 Redis 命令添加到 pipeline 中
	pipe := redisClient.Pipeline()
	pipe.Get(ctx, continuousKey)                                    // 获取连续签到天数
	pipe.TTL(ctx, continuousKey)                                    // 获取连续签到天数的剩余过期时间
	pipe.Get(ctx, tools.GenSignCountKey(userID))                    // 获取累计签到天数
	pipe.Get(ctx, tools.GenBitmapKey(userID, timeStatus.YearMonth)) // 获取位图信息
	result, err := pipe.Exec(ctx)

	// 执行 pipeline 中的所有命令
	if err != nil {
		panic(err)
	}

	// 获取 pipeline 执行结果
	signInContinuous := result[0].(*redis.StringCmd)           // 连续签到天数
	signInContinuousExpireAt := result[1].(*redis.DurationCmd) // 连续签到天数的剩余过期时间
	signInCount := result[2].(*redis.StringCmd)                // 累计签到天数
	signInBitmap := result[3].(*redis.StringCmd)               // 位图信息

	// 打印用户签到数据
	fmt.Printf("用户[%d]签到成功，已连续签到：%s(天),连续签到到期时间:%s（ttl:%s）,累计签到：%s(天)\n",
		userID,
		signInContinuous.Val(),
		tools.NextExpireTime(),
		signInContinuousExpireAt.Val(),
		signInCount.Val())

	// 获取位图信息并转换为 JSON 格式
	bitmapBytes, _ := signInBitmap.Bytes()
	fmt.Println("用户签到表：")
	bitmapJson := tools.BitmapToJSON(bitmapBytes, timeStatus)
	fmt.Println(string(bitmapJson)) // 打印 JSON 字符串
}
