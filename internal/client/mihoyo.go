// internal/client/mihoyo.go

// 重构

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/My-TuDo/gopher-mi-sentry/internal/config"
)

// MihoyoClient 米游社客户端
type MiClient struct {
	HttpClient *http.Client
}

// NewMiClient 初始化客户端
func NewMiClient() *MiClient {
	return &MiClient{
		HttpClient: &http.Client{
			Timeout: 10 * time.Second, // 设置 10s 超时
		},
	}
}

// MiResponse 定义统一的返回结构
type MiResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		List []GameRole `json:"list"`
	} `json:"data"`
}

type GameRole struct {
	Nickname   string `json:"nickname"`
	Level      int    `json:"level"`
	GameUid    string `json:"game_uid"`
	RegionName string `json:"region_name"`
}

// GetGameRoles 获取绑定的游戏角色（验证 Cookie 是否有效）
func (mc *MiClient) GetGameRoles(gameBiz string) (*MiResponse, error) {
	url := fmt.Sprintf("https://api-takumi.mihoyo.com/binding/api/getUserGameRolesByCookie?game_biz=%s", gameBiz)

	req, _ := http.NewRequest("GET", url, nil)

	// 从全局配置里读取 Cookie 和 UA
	req.Header.Set("Cookie", config.GlobalConfig.Mihoyo.Cookie)
	req.Header.Set("User-Agent", config.GlobalConfig.Mihoyo.UserAgent)

	resp, err := mc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// -- 反序列化 --
	var miResp MiResponse
	if err := json.Unmarshal(body, &miResp); err != nil {
		return nil, fmt.Errorf("解析失败： %w", err)
	}
	return &miResp, nil
}
