package main

import (
	"fmt"
	"log"

	"github.com/My-TuDo/gopher-mi-sentry/internal/client"
	"github.com/My-TuDo/gopher-mi-sentry/internal/config"
)

func main() {
	// 1. 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("初始化失败： %v", err)
	}
	fmt.Println("配置加载成功！")

	// 2. 初始化米游社客户端
	miClient := client.NewMiClient()

	// 3. 执行一次探测
	fmt.Println(" 正在尝试连接米哈游服务器...")
	result, err := miClient.GetGameRoles()
	if err != nil {
		log.Fatalf("通信失败： %v", err)
	}
	fmt.Println("米哈游返回的数据：")
	fmt.Println(result)
}
