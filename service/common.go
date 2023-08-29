package service

import (
	"easy-signin/pkg/tools"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// CheckSignIn 检查是否已经签到
func (s *Service) CheckSignIn(userID int64) (bool, error) {
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

// PrintUserSignInData 打印用户登录数据，连续签到天数、累计签到天数以及签到表
func (s *Service) PrintUserSignInData(userID int64) error {
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
		return handleRedisError(err)
	}

	// 获取 pipeline 执行结果
	signInContinuous := result[0].(*redis.StringCmd)           // 连续签到天数
	signInContinuousExpireAt := result[1].(*redis.DurationCmd) // 连续签到天数的剩余过期时间
	signInCount := result[2].(*redis.StringCmd)                // 累计签到天数
	signInBitmap := result[3].(*redis.StringCmd)               // 位图信息

	// 打印用户签到数据
	fmt.Printf("用户[%d]操作成功，已连续签到：%s(天),连续签到到期时间:%s（ttl:%s）,累计签到：%s(天)\n",
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

	return nil
}

// getSigninForm 获取用户签到表
func getSigninForm(userID int64) ([]byte, error) {
	if bitmap, err := redisClient.Get(ctx, tools.GenBitmapKey(userID, timeStatus.YearMonth)).Bytes(); err != nil {
		return nil, err
	} else {
		return bitmap, nil
	}
}

// 错误处理函数
func handleRedisError(err error) error {
	// 日志记录
	// 返回统一错误
	// 尝试回滚操作
	return fmt.Errorf("%v", err)
}
