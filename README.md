# easy-signin
基于redis的bitmap实现签到和补卡

## 目录结构
```
├── LICENSE
├── README.md
├── go.mod
├── go.sum
├── main.go
├── pkg
│   ├── app
│   │   └── app.go
│   └── tools
│       ├── bitmapTools.go
│       ├── cornTools.go
│       ├── keyTools.go
│       └── timeTools.go
└── service
├── common.go
├── makeUpSignIn.go
├── service.go
└── signIn.go


```

## 运行方式
1. 修改pkg/app/app.go中initRedisClient()函数里redis连接的密码，如没有密码可以删去密码一行
2. 启动redis 
3. `go run main.go cmd userID day`
> 需要三个参数：
> 1. 参数1：cmd（signin|makeup|print）用于签到、补签和打印签到信息
> 2. 参数2：userID 用户ID
> 3. 参数3：day 补签日期（如非补卡命令，请任意输入1~9任意一个数字）
  

