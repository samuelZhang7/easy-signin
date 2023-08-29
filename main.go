package main

import (
	"easy-signin/pkg/app"
	"easy-signin/pkg/tools"
	"easy-signin/service"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var appClient *app.App // app管理对象

func main() {
	// 1. 获取App对象
	appClient = app.GetApp()

	// 2. 启动Cron
	// 步骤1和2不能换，因为TimeStatus在步骤1才会初始化。
	// 如果先启动Corn，而恰巧scheduler正好执行，由于TimeStatus没有初始化，则会报错。
	tools.StartCron()

	// 3. 获取Service服务
	service := service.GetService(appClient)

	// 4. 接收命令行参数
	args, err := getArgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd := args[0]
	userID := args[1]
	day := args[2]

	// 5. 处理cmd
	switch cmd {
	case 1:
		if err := service.SignIn(userID); err != nil {
			fmt.Printf("用户[%d]签到失败,失败原因：%v", userID, err)
			return
		}
	case 2:
		if err := service.MakeUpSignIn(userID, int(day)); err != nil {
			fmt.Printf("用户[%d]补签失败,失败原因：%v", userID, err)
			return
		}
	case 3:
		if err := service.PrintUserSignInData(userID); err != nil {
			fmt.Printf("用户[%d]签到数据打印失败,失败原因：%v", userID, err)
			return
		}

	}

}

// getArgs 获取用户输入的参数
func getArgs() ([]int64, error) {

	args := make([]int64, 3)
	// 从命令行参数获取操作命令和userID
	if len(os.Args) < 4 {
		return args, fmt.Errorf("请传入参数1：signin|makeup|print，参数2：userID, 参数3：补签的日期（如非补卡命令，请任意输入1~9任意一个数字）")
	}

	// 第一个参数是可执行文件的路径，忽略

	// 获取第二个参数，操作命令
	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "signin":
		args[0] = 1
	case "makeup":
		args[0] = 2
	case "print":
		args[0] = 3
	default:
		return args, fmt.Errorf("：请传入正确的参数1：signin|makeup|print，参数2：userID")
	}

	// 获取第三个参数，userID
	userID, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		return args, fmt.Errorf("无效的userID：%s\n", os.Args[2])
	}
	args[1] = userID

	// 获取第4个参数，补卡日期
	day, err := strconv.ParseInt(os.Args[3], 10, 64)
	// 此处只做了基础的日期判断，实际还应该做日期是否大于当月的最大天数
	if err != nil || day < 1 || day > 31 {
		return args, fmt.Errorf("无效的补签日期（如非补卡命令，请任意输入1~9任意一个数字）：%s\n", os.Args[2])
	}
	args[2] = day

	return args, nil
}
