package tools

import "fmt"

const signInContinuousKey = "cc_uid_%d:sign_in_continuous"
const signInCountKey = "cc_uid_%d:sign_in_count"
const signInBitmapKey = "cc_uid_%d_%s:sign_in_bitmap"

// GenContinuousKey  生成用户连续签到总数的Key
func GenContinuousKey(userID int64) string {
	key := fmt.Sprintf(signInContinuousKey, userID)
	return key
}

// GenSignCountKey 生成用户签到总数的Key
func GenSignCountKey(userID int64) string {
	key := fmt.Sprintf(signInCountKey, userID)
	return key
}

// GenBitmapKey 生成用户签到位图的Key
func GenBitmapKey(userID int64, yearMonth string) string {
	key := fmt.Sprintf(signInBitmapKey, userID, yearMonth)
	return key
}
