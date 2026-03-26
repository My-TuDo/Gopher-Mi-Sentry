package client

import (
	"io"
	"net/http"

	"github.com/My-TuDo/gopher-mi-sentry/internal/config"
)

// MihoyoClient 米游社客户端
type MiClient struct {
	HttpClient *http.Client
}

// NewMiClient 初始化客户端
func NewMiClient() *MiClient {
	return &MiClient{
		HttpClient: &http.Client{},
	}
}

// GetGameRoles 获取绑定的游戏角色（验证 Cookie 是否有效）
func (mc *MiClient) GetGameRoles() (string, error) {
	url := "https://api-takumi.mihoyo.com/binding/api/getUserGameRolesByCookie?game_biz=hk4e_cn"

	req, _ := http.NewRequest("GET", url, nil)

	// 从全局配置里读取 Cookie 和 UA
	req.Header.Set("Cookie", config.GlobalConfig.Mihoyo.Cookie)
	req.Header.Set("User-Agent", config.GlobalConfig.Mihoyo.UserAgent)

	resp, err := mc.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}
