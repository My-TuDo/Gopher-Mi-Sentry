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

	// 尝试查询原神角色
	relos, err := miClient.GetGameRoles("hk4e_cn")
	if err != nil || relos.Retcode != 0 {
		log.Fatalf("获取游戏角色失败： %v, 消息： %s", err, relos.Message)
	}

	fmt.Printf("发现账号 [%s] 下共有 %d 个角色：\n", config.GlobalConfig.Mihoyo.Nickname, len(relos.Data.List))
	for _, role := range relos.Data.List {
		fmt.Printf("- %s (UID: %s, 等级: %d, 服务器: %s)\n", role.Nickname, role.GameUid, role.Level, role.RegionName)
	}

	// 遍历角色并签到
	for _, role := range relos.Data.List {
		fmt.Printf("正在为角色 [%s] 执行签到...\n", role.Nickname)

		res, err := miClient.DoSign(role)
		if err != nil {
			fmt.Printf("签到通信异常： %v\n", err)
			continue
		}

		if res.Retcode == 0 {
			fmt.Printf("签到成功！奖励已发往油箱。消息： %s\n", res.Message)
		} else {
			fmt.Printf("签到反馈： %s (错误码: %d)\n", res.Message, res.Retcode)
		}

	}

	// // 3. 执行一次探测
	// fmt.Println(" 正在尝试连接米哈游服务器...")
	// result, err := miClient.GetGameRoles()
	// if err != nil {
	// 	log.Fatalf("通信失败： %v", err)
	// }
	// fmt.Println("米哈游返回的数据：")
	// fmt.Println(result)
}
