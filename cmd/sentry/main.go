package main

import (
	"fmt"
	"log"
	"time"

	"github.com/My-TuDo/gopher-mi-sentry/internal/client"
	"github.com/My-TuDo/gopher-mi-sentry/internal/config"
	"github.com/My-TuDo/gopher-mi-sentry/internal/database"
	"github.com/My-TuDo/gopher-mi-sentry/internal/network"
	"github.com/robfig/cron/v3"
)

func ExecuteTasks(miClient *client.MiClient) {
	fmt.Printf("\n [%s] 自动化任务开始执行...\n", time.Now().Format("15:04:05"))

	// 1. 网络预测
	latency, err := network.CheckMihoyoStatus()
	if err != nil {
		log.Printf("网络异常，跳过本次执行： %v", err)
		return
	}
	fmt.Printf("网络延迟： %v\n", latency)

	// 2. 获取账号
	relos, err := miClient.GetGameRoles("hk4e_cn")
	if err != nil || relos.Retcode != 0 {
		log.Printf("获取游戏角色失败： %v, 消息： %s", err, relos.Message)
		return
	}

	// 3. 遍历角色并签到
	for _, role := range relos.Data.List {
		// A. 执行签到
		fmt.Printf("正在为角色 [%s] 执行签到...\n", role.Nickname)
		res, err := miClient.DoSign(role)

		signStatus := "Success" // 初始化状态假设成功
		if err != nil {
			log.Printf("签到通信异常： %v\n", err)
			signStatus = "NetworkError" // 网络异常
		} else if res.Retcode != 0 {
			fmt.Printf("签到反馈： %s (错误码: %d)\n", res.Message, res.Retcode)
			signStatus = res.Message // 直接使用米哈游的反馈作为状态，方便后续分析
		} else {
			fmt.Println("签到成功！")
		}

		// B. 同步数据库
		acc := &database.Account{
			UID:       role.GameUid,
			Nickname:  role.Nickname,
			Cookie:    config.GlobalConfig.Mihoyo.Cookie,
			Status:    signStatus, // 将签到状态存入数据库，方便后续查询和分析
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := database.SyncAccount(acc); err != nil {
			log.Printf("同步UID： %s 失败： %v", role.GameUid, err)
		}
	}
	fmt.Println("本轮任务执行完毕！")
}

func main() {
	// 1. 系统初始化（这部分只跑一次）
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("配置加载失败： %v", err)
	}
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败： %v", err)
	}
	miClient := client.NewMiClient()

	fmt.Println("Gopher-Mi-Sentry 自动化引擎启动成功！")

	// --- 核心改动 2： 配置调度经理 ---
	// 使用 cron.WithSeconds() 来支持秒级调度，方便测试和演示
	c := cron.New(cron.WithSeconds())

	// 添加任务计划： 每分钟的第 0 秒执行一次
	// 语法含义： 秒 分 时 日 月 周
	_, err := c.AddFunc("0 */1 * * * *", func() {
		ExecuteTasks(miClient)
	})
	if err != nil {
		log.Fatalf("无法添加定时任务： %v", err)
	}

	// 2. 启动经理（它会自己开一个协程去后台盯着表）
	c.Start()

	// --- 核心改动 3： 守住大门 ---
	// 这一行非常关键！如果没有它， main跑完 Start() 后就直接退出了，定时任务根本来不及执行
	fmt.Println("正在后台守候任务清单，按 Ctrl+C 退出程序...")
	select {}
}
