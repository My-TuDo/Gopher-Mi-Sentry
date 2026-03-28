package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/My-TuDo/gopher-mi-sentry/internal/config"
)

// SignRequest 签到请求参数
type SignRequest struct {
	ActID  string `json:"act_id"`
	Region string `json:"region"`
	Uid    string `json:"uid"`
}

// DoSign 执行签到
func (mc *MiClient) DoSign(role GameRole) (*MiResponse, error) {
	// 米游社原神福利签到地址
	url := "https://api-takumi.mihoyo.com/event/bbs_sign_reward/sign"

	// 1. 构造请求体
	reqBody := SignRequest{
		ActID:  "e38215756454", // 固定活动 ID, 建议以后放进 config.yaml
		Region: role.RegionName,
		Uid:    role.GameUid,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	// 2. 创建请求
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))

	// 3. 注入Header
	req.Header.Set("Cookie", config.GlobalConfig.Mihoyo.Cookie)
	req.Header.Set("User-Agent", config.GlobalConfig.Mihoyo.UserAgent)
	req.Header.Set("Referer", "https://act.mihoyo.com/")
	req.Header.Set("Accept", "application/json")

	// 4. 发送
	resq, err := mc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resq.Body.Close()

	// 5. 解析回执
	body, _ := io.ReadAll(resq.Body)
	var miResp MiResponse
	json.Unmarshal(body, &miResp)

	return &miResp, nil
}
