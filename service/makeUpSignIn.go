package service

import (
	"easy-signin/pkg/tools"
	"fmt"
)

// MakeUpSignIn 用户补签
func (s *Service) MakeUpSignIn(userID int64, day int) error {
	// 1. 获取用户签到表
	bitmap, err := getSigninForm(userID)
	if err != nil {
		return err
	}

	// 2. 判断用户补签的day,并模拟用户补签
	signed := checkAndSimulateSignin(bitmap, day)
	//	用户已签到，提醒用户不可补签
	if signed {
		return fmt.Errorf("用户[%d]已签到，不支持补签\n", userID)
	}

	// 3. 执行用户补签,更新用户签到信息
	if err := makeupCard(bitmap, day, userID); err != nil {
		return handleRedisError(err)
	}

	// 补签成功
	// 补签记录入influxDB
	// 后续业务逻辑...

	//4. 打印用户签到信息
	if err := s.PrintUserSignInData(userID); err != nil {
		return err
	}

	return nil
}

// checkAndSimulateSignin 判断用户补签的day,并模拟用户补签
func checkAndSimulateSignin(bitmap []byte, day int) bool {
	// bitmap是一个4字节，32位,只要定位到相应的字节，然后根据day构建一个字节，然后对两个字节位运算
	// setBit时从0开始，因此第day天在第day-1位上
	bit := day - 1
	// 1. 找到bit所属的byte
	index := bit / 8
	// 2. 根据bit,获取在byte中的相对位置，然后设置该位为1
	mask := byte(1 << (7 - bit%8))
	// 3. 判断bit位是否是1
	signed := bitmap[index]&mask != 0
	// 4. 模拟用户签到,将bit位置设置为1（后续要计算连续签到天数，因此在这里模拟用户签到）
	bitmap[index] = bitmap[index] | mask
	return signed
}

// makeupCard 执行用户补签,更新用户签到信息
func makeupCard(bitmap []byte, day int, userID int64) error {
	// 1. 获取补签后的连续签到次数
	continuousDays := getContinuousDays(bitmap)
	pipe := redisClient.Pipeline()
	// 设置签到位图
	pipe.SetBit(ctx, tools.GenBitmapKey(userID, timeStatus.YearMonth), int64(day-1), 1)
	// 更新签到总数
	pipe.Incr(ctx, tools.GenSignCountKey(userID))
	// 获取连续签到key
	continuousKey := tools.GenContinuousKey(userID)
	// 更新连续签到总数
	pipe.Set(ctx, continuousKey, continuousDays, 0)
	// 更新连续签到的过期时间
	pipe.ExpireAt(ctx, continuousKey, tools.NextExpireTime())
	// 执行
	if _, err := pipe.Exec(ctx); err != nil {
		return handleRedisError(err)
	}
	return nil
}

// getContinuousDays 获取用户连续签到的天数
func getContinuousDays(bitmap []byte) int {
	// 连续天数
	continuousDays := 0
	// 获取bitmap的长度
	length := len(bitmap)
	// bitmap的索引
	index := length - 1
	// bitmap中每个byte的位索引
	bit := 0

	// 1. 从最后一个字节开始往前遍历,找到非0值
Loop1:
	for ; index >= 0; index-- {
		// 获取当前字节的值
		byteValue := bitmap[index]
		// 遍历当前字节，找到非0值
		for ; bit < 8; bit++ {
			if byteValue>>bit&1 == 1 {
				break Loop1
			}
		}
		bit = 0 // 初始化bitIndex为0
	}

	// 2. 从当前位置往前找，如果有0值，则表示已找到连续天数
Loop2:
	for ; index >= 0; index-- {
		// 获取当前字节的值
		byteValue := bitmap[index]
		// 遍历当前字节
		for ; bit < 8; bit++ {
			if byteValue>>bit&1 == 0 {
				break Loop2
			}
			// 更新持续天数
			continuousDays++
		}
		bit = 0 // 初始化bitIndex为0
	}
	fmt.Println(continuousDays)
	return continuousDays
}
