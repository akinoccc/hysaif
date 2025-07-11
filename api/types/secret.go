package types

import "github.com/akinoccc/hysaif/api/models"

// 密钥项相关类型
type PostItemRequest struct {
	Name        string                `json:"name" binding:"required,min=1,max=100"`
	Type        string                `json:"type" binding:"required,oneof=password api_key access_key ssh_key certificate token kv"` // password, api_key, access_key, ssh_key, certificate, token, kv
	Description string                `json:"description,omitempty" binding:"max=500"`
	Category    string                `json:"category" binding:"required,min=1,max=50"`
	Environment string                `json:"environment" binding:"required,oneof=dev staging prod test"`
	Tags        []string              `json:"tags,omitempty" gorm:"type:text;serializer:json"`
	Data        models.SecretItemData `json:"data" binding:"required" gorm:"type:text;serializer:json"`
	ExpiresAt   uint64                `json:"expires_at,omitempty"`
}

type ItemsListParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Category string `form:"category"`
	Type     string `form:"type" binding:"omitempty,oneof=password api_key access_key ssh_key certificate token kv"`
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	SortDesc bool   `form:"sort_desc"`
}

// 版本历史相关类型
type SecretItemHistoryListParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type RestoreSecretItemFromHistoryRequest struct {
	Version int    `json:"version" binding:"required,min=1"`
	Reason  string `json:"reason" binding:"max=500"`
}

type CompareVersionsRequest struct {
	Version1 int `json:"version1" binding:"required,min=1"`
	Version2 int `json:"version2" binding:"required,min=1"`
}

type VersionComparisonResponse struct {
	Version1 int                    `json:"version1"`
	Version2 int                    `json:"version2"`
	Changes  map[string]interface{} `json:"changes"`
}
