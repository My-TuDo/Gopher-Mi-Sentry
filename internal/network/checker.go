package network

import (
	"fmt"
	"net"
	"time"
)

// CheckMihoyoStatus 探测米哈游服务器连通性
// 采用 SRE 级别的重试机制，防止因为瞬间网络波动（如 DNS 超时）导致误报
func CheckMihoyoStatus() (time.Duration, error) {
	hostname := "api-takumi.mihoyo.com"
	port := "443"

	// 使用职业写法： 防止 IPV6 报错
	target := net.JoinHostPort(hostname, port)

	var err error
	var conn net.Conn

	// --- 核心 SRE 逻辑： 重试 3 次 ---
	for i := 1; i <= 3; i++ {
		start := time.Now()
		// 每次探测限时 10s
		conn, err = net.DialTimeout("tcp", target, 10*time.Second)

		if err == nil {
			defer conn.Close()
			return time.Since(start), nil
		}

		fmt.Printf("第 %d 次网络探测失败: %v， 1秒后重试...\n", i, err)
		time.Sleep(1 * time.Second)
	}
	return 0, fmt.Errorf("米哈游服务器连接超时，请检查您的网络环境")
}
