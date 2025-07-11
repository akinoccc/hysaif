package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SecretItemHistory 密钥历史版本模型
type SecretItemHistory struct {
	ID           string          `json:"id" gorm:"primaryKey"`
	SecretItemID string          `json:"secret_item_id" gorm:"index;not null"`   // 关联的密钥项ID
	Version      int             `json:"version" gorm:"not null"`                // 版本号
	Name         string          `json:"name" gorm:"not null"`                   // 当时的名称
	Description  string          `json:"description"`                            // 当时的描述
	Type         string          `json:"type" gorm:"not null"`                   // 当时的类型
	Category     string          `json:"category"`                               // 当时的分类
	Tags         []string        `json:"tags" gorm:"serializer:json"`            // 当时的标签
	Data         *SecretItemData `json:"data" gorm:"type:text"`                  // 当时的敏感数据
	ExpiresAt    uint64          `json:"expires_at"`                             // 当时的过期时间
	Environment  string          `json:"environment"`                            // 当时的环境
	ChangeType   string          `json:"change_type" gorm:"not null"`            // 变更类型：created, updated, deleted
	ChangeReason string          `json:"change_reason"`                          // 变更原因
	CreatedAt    uint64          `json:"created_at" gorm:"autoCreateTime:milli"` // 创建时间
	CreatedByID  string          `json:"created_by_id" gorm:"index;not null"`    // 创建者ID

	// 关联数据
	SecretItem *SecretItem `json:"secret_item,omitempty" gorm:"foreignKey:SecretItemID;references:ID"`
	CreatedBy  *User       `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID;references:ID"`
}

// BeforeCreate 钩子函数，在创建记录之前设置ID
func (sih *SecretItemHistory) BeforeCreate(tx *gorm.DB) (err error) {
	sih.ID = uuid.New().String()
	return
}

// 变更类型常量
const (
	HistoryChangeTypeCreated = "created"
	HistoryChangeTypeUpdated = "updated"
	HistoryChangeTypeDeleted = "deleted"
)

// CreateSecretItemHistory 创建密钥历史版本记录
func CreateSecretItemHistory(secretItem *SecretItem, changeType, reason string, createdByID string) error {
	// 获取当前最大版本号
	var maxVersion int
	err := DB.Model(&SecretItemHistory{}).
		Where("secret_item_id = ?", secretItem.ID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error
	if err != nil {
		return fmt.Errorf("获取版本号失败: %w", err)
	}

	// 创建历史记录
	history := &SecretItemHistory{
		SecretItemID: secretItem.ID,
		Version:      maxVersion + 1,
		Name:         secretItem.Name,
		Description:  secretItem.Description,
		Type:         secretItem.Type,
		Category:     secretItem.Category,
		Tags:         secretItem.Tags,
		Data:         secretItem.Data,
		ExpiresAt:    secretItem.ExpiresAt,
		Environment:  secretItem.Environment,
		ChangeType:   changeType,
		ChangeReason: reason,
		CreatedByID:  createdByID,
	}

	return DB.Create(history).Error
}

// GetSecretItemHistory 获取密钥历史版本列表
func GetSecretItemHistory(secretItemID string) ([]SecretItemHistory, error) {
	var histories []SecretItemHistory
	err := DB.Where("secret_item_id = ?", secretItemID).
		Preload("CreatedBy").
		Order("version DESC").
		Find(&histories).Error
	return histories, err
}

// GetSecretItemHistoryByVersion 获取指定版本的密钥历史记录
func GetSecretItemHistoryByVersion(secretItemID string, version int) (*SecretItemHistory, error) {
	var history SecretItemHistory
	err := DB.Where("secret_item_id = ? AND version = ?", secretItemID, version).
		Preload("CreatedBy").
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// RestoreSecretItemFromHistory 从历史版本恢复密钥项
func RestoreSecretItemFromHistory(secretItemID string, version int, restoredByID string) (*SecretItem, error) {
	// 获取历史版本
	history, err := GetSecretItemHistoryByVersion(secretItemID, version)
	if err != nil {
		return nil, fmt.Errorf("获取历史版本失败: %w", err)
	}

	// 获取当前密钥项
	var currentItem SecretItem
	err = DB.Where("id = ?", secretItemID).First(&currentItem).Error
	if err != nil {
		return nil, fmt.Errorf("获取当前密钥项失败: %w", err)
	}

	// 先保存当前状态到历史记录
	err = CreateSecretItemHistory(&currentItem, HistoryChangeTypeUpdated, "恢复到历史版本", restoredByID)
	if err != nil {
		return nil, fmt.Errorf("保存当前状态到历史失败: %w", err)
	}

	// 更新密钥项为历史版本的数据
	currentItem.Name = history.Name
	currentItem.Description = history.Description
	currentItem.Type = history.Type
	currentItem.Category = history.Category
	currentItem.Tags = history.Tags
	currentItem.Data = history.Data
	currentItem.ExpiresAt = history.ExpiresAt
	currentItem.Environment = history.Environment
	currentItem.UpdatedByID = restoredByID

	err = DB.Save(&currentItem).Error
	if err != nil {
		return nil, fmt.Errorf("恢复密钥项失败: %w", err)
	}

	// 重新查询以获取关联数据
	DB.Preload("Creator").Preload("Updater").First(&currentItem, "id = ?", currentItem.ID)

	return &currentItem, nil
}

// CompareSecretItemVersions 比较两个版本的差异
func CompareSecretItemVersions(secretItemID string, version1, version2 int) (map[string]interface{}, error) {
	var history1, history2 SecretItemHistory

	err := DB.Where("secret_item_id = ? AND version = ?", secretItemID, version1).First(&history1).Error
	if err != nil {
		return nil, fmt.Errorf("获取版本%d失败: %w", version1, err)
	}

	err = DB.Where("secret_item_id = ? AND version = ?", secretItemID, version2).First(&history2).Error
	if err != nil {
		return nil, fmt.Errorf("获取版本%d失败: %w", version2, err)
	}

	// 比较字段差异
	diff := make(map[string]interface{})

	if history1.Name != history2.Name {
		diff["name"] = map[string]string{"old": history1.Name, "new": history2.Name}
	}
	if history1.Description != history2.Description {
		diff["description"] = map[string]string{"old": history1.Description, "new": history2.Description}
	}
	if history1.Type != history2.Type {
		diff["type"] = map[string]string{"old": history1.Type, "new": history2.Type}
	}
	if history1.Category != history2.Category {
		diff["category"] = map[string]string{"old": history1.Category, "new": history2.Category}
	}
	if history1.Environment != history2.Environment {
		diff["environment"] = map[string]string{"old": history1.Environment, "new": history2.Environment}
	}
	if history1.ExpiresAt != history2.ExpiresAt {
		diff["expires_at"] = map[string]uint64{"old": history1.ExpiresAt, "new": history2.ExpiresAt}
	}

	// 比较标签
	tags1Json, _ := json.Marshal(history1.Tags)
	tags2Json, _ := json.Marshal(history2.Tags)
	if string(tags1Json) != string(tags2Json) {
		diff["tags"] = map[string][]string{"old": history1.Tags, "new": history2.Tags}
	}

	// 比较数据（敏感数据不直接比较，只标记是否有变化）
	data1Json, _ := json.Marshal(history1.Data)
	data2Json, _ := json.Marshal(history2.Data)
	if string(data1Json) != string(data2Json) {
		diff["data"] = map[string]string{"old": "***", "new": "***"}
	}

	return diff, nil
}
