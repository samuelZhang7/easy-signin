package tools

import (
	"encoding/json"
	"fmt"
)

type SignRecord struct {
	Date   string `json:"date"`
	Signed int    `json:"signed"`
}

const februaryOfLeapYear = 29

// 每月的天数
var daysInMonth = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// bitmapToIntArray 将一个字节数组转成一个整形数组，也就是将0和1转成整形的0和1
func bitmapToIntArray(bitmap []byte) []int {
	binaryArray := make([]int, len(bitmap)*8) // 创建了一个长度和容量都为 len(bitmap)*8，即32的切片
	for i, b := range bitmap {                // 遍历字节数组，i 是索引，b 是字节值
		for j := 7; j >= 0; j-- { // 遍历每个字节的8个位，从高位到低位，因为在bitmap里面SETBIT是从高位开始
			bitValue := int((b >> j) & 1)   // 用右移和与运算来获取第j位的值，然后转换为整数，（&1：同为1，则返回1，否则返回0）
			binaryArray[i*8+7-j] = bitValue // 将位值存入整形数组，注意索引的计算方式（从高位获取，要从低位写）
		}
	}
	return binaryArray // 返回整形数组
}

// BitmapToJSON 将一个字节数组转成一个 JSON 格式的字符串，表示每天的签到情况
func BitmapToJSON(bitmap []byte, status *TimeStatus) []byte {
	// 将bitmap转为int数组，也就是将0和1转成整形的0和1
	bitmapArray := bitmapToIntArray(bitmap)
	// 根据每月的天数设置SignRecord 切片的大小
	signRecordSize := daysInMonth[status.Month-1]
	// 处理闰年情况
	if status.Month == 2 && status.IsLeapYear {
		signRecordSize = februaryOfLeapYear // 29天
	}
	// 根据每月的天数创建 SignRecord 切片，用来存储每天的签到记录
	signRecordArray := make([]SignRecord, signRecordSize)
	// 拼接year和month
	yearMonth := fmt.Sprintf("%d-%02d", status.Year, status.Month)
	// 遍历bitmapArray生成数据
	for i, v := range bitmapArray {
		if i == signRecordSize {
			break // bitmapArray数组共有32个元素，而每月最多只有31天，因此当i==signRecordSize，表示该月已结束
		}
		// 根据年月和索引生成日期字符串，例如 "2021-01-01"
		date := fmt.Sprintf("%s-%02d", yearMonth, i+1)
		// 创建一个 SignRecord 结构体，用来表示一天的签到情况
		item := SignRecord{
			Date:   date, // 日期字段
			Signed: v,    // 签到状态字段，0 表示未签到，1 表示已签到
		}
		// 将 SignRecord 结构体存入 signRecordArray 切片中
		signRecordArray[i] = item
	}

	// 生成JSON
	jsonBytes, _ := json.Marshal(signRecordArray) // 使用 json 包的 Marshal 函数将 signRecordArray 切片转换为 JSON 格式的字节数组
	return jsonBytes                              // 返回 JSON 字节数组
}
