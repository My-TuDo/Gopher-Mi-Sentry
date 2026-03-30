package database

import (
	"fmt"

	"github.com/My-TuDo/gopher-mi-sentry/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Account 对应的数据库里的表
// Account 对应的数据库里的表
type Account struct {
	ID uint `gorm:"primaryKey"`
	// 重点：必须增加 size:100，将其转为 varchar(100)，MySQL 才能建立索引
	UID       string `gorm:"size:100;uniqueIndex"`
	Nickname  string `gorm:"size:100"`
	Cookie    string `gorm:"type:text"`
	Status    string `gorm:"size:255"`
	UpdatedAt string `gorm:"size:100"`
}

// InitDB 初始化数据库
func InitDB() error {
	dsn := config.GlobalConfig.Database.DSN // 我们要在 YAML 里加上这个
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// --- 核心 SRE 逻辑： 自动同步表结构 ---
	err = db.AutoMigrate(&Account{})
	if err != nil {
		return fmt.Errorf("表结构同步失败：%w", err)
	}
	DB = db
	return nil
}
