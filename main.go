package main

import (
	"easy-signin/pkg/app"
	"easy-signin/pkg/tools"
	"easy-signin/service"
	"fmt"
	"os"
	"strconv"
)

var appClient *app.App // app管理对象

func main() {
	// 1. 获取App对象
	appClient = app.GetApp()

	// 2. 启动Cron
	// 步骤1和2不能换，因为TimeStatus在步骤1才会初始化。
	// 如果先启动Corn，而恰巧scheduler正好执行，由于TimeStatus没有初始化，则会报错。
	tools.StartCron()

	// 3. 接收命令行参数
	userID, err := getArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4. 获取Service服务
	service := service.GetService(appClient)

	// 5. 签到
	if err := service.SignIn(userID); err != nil {
		fmt.Printf("用户[%d]签到失败,失败原因：%w", userID, err)
		return
	}
	// 手动打印用户的连续签到天数、累计签到天数以及签到表
	// service.PrintUserSignInData(userID)
}

func getArgs() (int64, error) {
	// 从命令行参数获取用户ID
	if len(os.Args) < 2 {
		return 0, fmt.Errorf("请传入userID作为命令行参数")
	}
	userIDStr := os.Args[1]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的userID：%s\n", userIDStr)
	}
	return userID, nil
}
