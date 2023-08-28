package tools

import (
	"fmt"
	"sync"
	"time"
)

// TimeStatus 时间状态结构体
type TimeStatus struct {
	Year       int
	Month      int
	Day        int64
	IsLeapYear bool
	YearMonth  string
}

// 2024-2100年中的所有闰年
var leapYears = map[int]struct{}{
	2024: {},
	2028: {},
	2032: {},
	2036: {},
	2040: {},
	2044: {},
	2048: {},
	2052: {},
	2056: {},
	2060: {},
}

var timeStatus *TimeStatus   // 时间状态对象
var timeStatusOnce sync.Once // 确保timeStatus单例

// GetTimeStatus 获取单例时间状态对象
func GetTimeStatus() *TimeStatus {
	timeStatusOnce.Do(
		func() {
			// 使用 new 函数创建一个 TimeStatus 对象，并赋值给 timeStatus
			timeStatus = new(TimeStatus)
			// 调用 UpdateTimeStatus 函数来更新时间状态对象
			UpdateTimeStatus()
		})

	return timeStatus
}

// UpdateTimeStatus 更新时间状态任务
func UpdateTimeStatus() {
	// 获取当前时间
	now := time.Now()
	// 获取时间状态对象
	status := timeStatus
	// 更新相关字段
	status.Year = now.Year()
	status.Month = int(now.Month())
	status.Day = int64(now.Day())
	// 拼接年和月，月格式化为2位,前面补0
	status.YearMonth = fmt.Sprintf("%d%02d", status.Year, status.Month)
	// 判断是否是闰年
	status.IsLeapYear = isLeapYear(status.Year)
}

// beginningOfDay 获取今天0点时间
func beginningOfDay() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// NextExpireTime 获取下一次过期的时间
func NextExpireTime() time.Time {
	// 获取后天零点时间
	return beginningOfDay().Add(48 * time.Hour)
}

// isLeapYear，判断给定的年份是否是闰年
func isLeapYear(year int) bool {
	// 如果map中存在该年份，返回true，否则返回false
	_, ok := leapYears[year]
	return ok
}
