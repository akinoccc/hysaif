package types

import "github.com/akinoccc/hysaif/api/models"

// 通知相关请求类型

// GetNotificationsRequest 获取通知列表请求
type GetNotificationsRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   string `form:"status"`   // unread, read, all
	Type     string `form:"type"`     // 通知类型过滤
	Priority string `form:"priority"` // 优先级过滤
}

// CreateNotificationRequest 创建通知请求（管理员使用）
type CreateNotificationRequest struct {
	RecipientIDs []string `json:"recipient_ids" binding:"required"`
	Type         string   `json:"type" binding:"required"`
	Title        string   `json:"title" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Priority     string   `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
	ExpiresAt    uint64   `json:"expires_at"` // 过期时间，0表示不过期
}

// MarkNotificationRequest 标记通知请求
type MarkNotificationRequest struct {
	NotificationIDs []string `json:"notification_ids" binding:"required"`
	Status          string   `json:"status" binding:"required,oneof=read unread"`
}

// NotificationListResponse 通知列表响应
type NotificationListResponse struct {
	Data        []models.Notification `json:"data"`
	Pagination  Pagination            `json:"pagination"`
	UnreadCount int64                 `json:"unread_count"`
}

// NotificationStatsResponse 通知统计响应
type NotificationStatsResponse struct {
	TotalCount  int64            `json:"total_count"`
	UnreadCount int64            `json:"unread_count"`
	ReadCount   int64            `json:"read_count"`
	ByType      map[string]int64 `json:"by_type"`
	ByPriority  map[string]int64 `json:"by_priority"`
}

// NotificationPreferencesRequest 通知偏好设置请求
type NotificationPreferencesRequest struct {
	EmailNotifications bool     `json:"email_notifications"`
	EnabledTypes       []string `json:"enabled_types"`
	QuietHours         struct {
		Enabled   bool   `json:"enabled"`
		StartTime string `json:"start_time"` // HH:MM 格式
		EndTime   string `json:"end_time"`   // HH:MM 格式
	} `json:"quiet_hours"`
}

// NotificationPreferencesResponse 通知偏好设置响应
type NotificationPreferencesResponse struct {
	UserID             string   `json:"user_id"`
	EmailNotifications bool     `json:"email_notifications"`
	EnabledTypes       []string `json:"enabled_types"`
	QuietHours         struct {
		Enabled   bool   `json:"enabled"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	} `json:"quiet_hours"`
	UpdatedAt uint64 `json:"updated_at"`
}

// BulkNotificationRequest 批量通知请求
type BulkNotificationRequest struct {
	UserRoles []string `json:"user_roles"` // 按角色发送
	UserIDs   []string `json:"user_ids"`   // 指定用户ID
	Type      string   `json:"type" binding:"required"`
	Title     string   `json:"title" binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Priority  string   `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
	ExpiresAt uint64   `json:"expires_at"`
}

// NotificationTemplateResponse 通知模板响应
type NotificationTemplateResponse struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Priority    string   `json:"priority"`
	Description string   `json:"description"`
	Variables   []string `json:"variables"` // 模板变量列表
}
