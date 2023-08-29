package service

import (
	"easy-signin/pkg/tools"
	"fmt"
)

// SignIn 用户签到
func (s *Service) SignIn(userID int64) error {
	// 1. 检查用户是否已签到
	signed, err := s.CheckSignIn(userID)
	if err != nil {
		return err
	}
	if signed {
		return fmt.Errorf("用户[%d]已签到，请勿重复签到\n", userID)
	}

	// 2. 执行签到操作
	if err := s.updateSignInfo(userID); err != nil {
		return err
	}

	// 签到成功
	// 签到记录入influxDB
	// 后续业务逻辑...

	// 3.打印用户签到信息
	if err := s.PrintUserSignInData(userID); err != nil {
		return err
	}

	return nil
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
	return nil
}
