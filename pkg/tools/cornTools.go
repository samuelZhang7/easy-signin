package tools

import (
	"github.com/robfig/cron"
)

// StartCron 用于创建和启动一个 cron 实例
func StartCron() {
	c := cron.New() // 创建一个新的 cron 实例
	// 每天0点更新时间状态
	c.AddFunc("@daily", UpdateTimeStatus) // 添加一个定时任务，每天零点更新时间状态对象
	// 启动 Cron
	c.Start() // 启动 cron 实例
}
