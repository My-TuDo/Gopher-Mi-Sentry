package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// 定义配置构造体（与 YAML 对应）
type Config struct {
	Mihoyo MihoyoConfig `mapstructure:"mihoyo"`
}

type MihoyoConfig struct {
	Cookie    string `mapstructure:"cookie"`
	UserAgent string `mapstructuere:"user_agent"`
}

var GlobalConfig *Config

// LoadConfig 加载配置（SRE 必备技能： 配置解耦）
func LoadConfig() error {
	viper.SetConfigName("config")    // 配置文件名
	viper.SetConfigType("yaml")      // 格式
	viper.AddConfigPath("./configs") // 配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 将配置映射到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}
