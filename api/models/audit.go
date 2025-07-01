package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ModelBase
	UserID     string `json:"-" gorm:"index"`
	Action     string `json:"action"`   // create, read, update, delete, login, logout
	Resource   string `json:"resource"` // user, secret_item
	ResourceID string `json:"resource_id"`
	Details    string `json:"details"` // JSON格式的详细信息
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`

	// 关联
	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

// BeforeCreate 钩子函数，在创建记录之前设置ID
func (al *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	al.ID = uuid.New().String()
	return
}
