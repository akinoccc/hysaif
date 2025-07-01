package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/akinoccc/hysaif/api/packages/crypto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SecretItem 敏感信息项模型
type SecretItem struct {
	ModelBase
	Name        string          `json:"name" gorm:"not null"`
	Description string          `json:"description"`
	Type        string          `json:"type" gorm:"not null"`        // password, api_key, access_key, ssh_key, certificate, token, custom
	Category    string          `json:"category"`                    // 分类：aws, aliyun, github等
	Tags        []string        `json:"tags" gorm:"serializer:json"` // JSON格式的标签数组
	Data        *SecretItemData `json:"data" gorm:"type:text"`       // 加密后的敏感数据
	ExpiresAt   uint64          `json:"expires_at"`                  // 过期时间, 0表示永不过期
	Environment string          `json:"environment"`                 // 环境：dev, test, prod
	CreatedByID string          `json:"-" gorm:"index"`              // 创建者ID
	UpdatedByID string          `json:"-" gorm:"index"`              // 更新者ID

	// 访问权限相关（不存储在数据库中）
	HasApprovedAccess bool `json:"has_approved_access" gorm:"-"` // 是否有已批准的访问申请

	// 关联用户
	Creator *User `json:"creator" gorm:"foreignKey:CreatedByID;references:ID"`
	Updater *User `json:"updater" gorm:"foreignKey:UpdatedByID;references:ID"`
}

// BeforeCreate 钩子函数，在创建记录之前设置ID
func (si *SecretItem) BeforeCreate(tx *gorm.DB) (err error) {
	si.ID = uuid.New().String()
	return
}

// SecretItemData 敏感信息数据结构（用于前端展示，不包含加密数据）
type SecretItemData struct {
	// 密码相关
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Address  string `json:"address,omitempty"`
	Notes    string `json:"notes,omitempty"`

	// API Key 相关
	APIKey    string `json:"api_key,omitempty"`
	APISecret string `json:"api_secret,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`

	// Access Key 相关
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	Region    string `json:"region,omitempty"`

	// SSH Key 相关
	PrivateKey string `json:"private_key,omitempty"`
	PublicKey  string `json:"public_key,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`

	// Token 相关
	Token        string `json:"token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	// 自定义数据
	CustomData []map[string]string `json:"custom_data,omitempty"`
}

// Value 实现 driver.Valuer 接口，用于将 SecretItemData 序列化为数据库存储格式
func (s SecretItemData) Value() (driver.Value, error) {
	// 将结构体转换为JSON
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SecretItemData: %w", err)
	}

	// 加密JSON数据
	encryptedData, err := crypto.Encrypt(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt SecretItemData: %w", err)
	}

	return encryptedData, nil
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取并反序列化 SecretItemData
func (s *SecretItemData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	// 获取加密的字符串数据
	var encryptedData string
	switch v := value.(type) {
	case string:
		encryptedData = v
	case []byte:
		encryptedData = string(v)
	default:
		return fmt.Errorf("cannot scan %T into SecretItemData", value)
	}

	// 如果是空字符串，返回空结构体
	if encryptedData == "" {
		return nil
	}

	// 解密数据
	decryptedData, err := crypto.Decrypt(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decrypt SecretItemData: %w", err)
	}

	// 反序列化JSON到结构体
	err = json.Unmarshal(decryptedData, s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal SecretItemData: %w", err)
	}

	return nil
}
