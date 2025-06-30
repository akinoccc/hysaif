package models

import (
	"fmt"

	"github.com/akinoccc/hysaif/api/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ModelBase 模型基类
type ModelBase struct {
	ID        string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CreatedAt uint64 `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt uint64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// DB 数据库实例
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() {
	var err error
	var dialector gorm.Dialector

	// 根据配置选择数据库驱动
	switch config.AppConfig.Database.Type {
	case "sqlite":
		dialector = sqlite.Open(config.AppConfig.GetDSN())
	case "postgres":
		dialector = postgres.Open(config.AppConfig.GetDSN())
	case "mysql":
		dialector = mysql.Open(config.AppConfig.GetDSN())
	default:
		panic(fmt.Sprintf("不支持的数据库类型: %s", config.AppConfig.Database.Type))
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
	})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	// 自动迁移 - User 模型必须首先创建，因为其他模型都依赖于它
	err = DB.AutoMigrate(&User{}, &SecretItem{}, &AccessRequest{}, &Notification{}, &WebAuthnCredential{}, &AuditLog{})
	if err != nil {
		panic("failed to migrate database")
	}

	// 创建默认管理员用户
	createDefaultAdmin()
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin() {
	var count int64
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		admin := User{
			Name:     "Super Admin",
			Password: "admin123",
			Email:    "admin@uozi.org",
			Role:     "super_admin",
			Status:   "active",
			// 不设置CreatedBy和UpdatedBy，避免外键约束冲突
		}
		// 使用Omit明确跳过CreatedBy和UpdatedBy字段
		err := DB.Omit("created_by", "updated_by").Create(&admin).Error
		if err != nil {
			panic(err)
		}
	}
}
