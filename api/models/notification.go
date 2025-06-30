package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 通知类型常量
const (
	NotificationTypeAccessRequestCreated  = "access_request_created"  // 新的访问申请
	NotificationTypeAccessRequestApproved = "access_request_approved" // 申请已批准
	NotificationTypeAccessRequestRejected = "access_request_rejected" // 申请已拒绝
	NotificationTypeAccessRequestExpired  = "access_request_expired"  // 申请已过期
	NotificationTypeSecretItemExpiring    = "secret_item_expiring"    // 密钥项即将过期
	NotificationTypeSecretItemExpired     = "secret_item_expired"     // 密钥项已过期
	NotificationTypeSystemMaintenance     = "system_maintenance"      // 系统维护通知
	NotificationTypeSecurityAlert         = "security_alert"          // 安全警报
)

// 通知状态常量
const (
	NotificationStatusUnread = "unread" // 未读
	NotificationStatusRead   = "read"   // 已读
)

// 通知优先级常量
const (
	NotificationPriorityLow    = "low"    // 低优先级
	NotificationPriorityNormal = "normal" // 普通优先级
	NotificationPriorityHigh   = "high"   // 高优先级
	NotificationPriorityUrgent = "urgent" // 紧急
)

// Notification 通知模型
type Notification struct {
	ModelBase
	RecipientID string `json:"-" gorm:"not null;index"`          // 接收者ID
	Type        string `json:"type" gorm:"not null;index"`       // 通知类型
	Title       string `json:"title" gorm:"not null"`            // 通知标题
	Content     string `json:"content" gorm:"type:text"`         // 通知内容
	Status      string `json:"status" gorm:"default:'unread'"`   // 通知状态
	Priority    string `json:"priority" gorm:"default:'normal'"` // 优先级
	RelatedID   string `json:"related_id" gorm:"index"`          // 相关资源ID（如申请ID、密钥项ID等）
	RelatedType string `json:"related_type"`                     // 相关资源类型
	ReadAt      uint64 `json:"read_at"`                          // 阅读时间
	ExpiresAt   uint64 `json:"expires_at"`                       // 过期时间
	Metadata    string `json:"metadata" gorm:"type:text"`        // 额外元数据（JSON格式）

	// 关联
	Recipient User `json:"recipient" gorm:"foreignKey:RecipientID"`
}

// BeforeCreate 钩子函数，在创建记录之前设置ID
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID = uuid.New().String()
	return
}

// IsRead 检查通知是否已读
func (n *Notification) IsRead() bool {
	return n.Status == NotificationStatusRead
}

// IsExpired 检查通知是否已过期
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == 0 {
		return false
	}
	now := uint64(time.Now().UnixMilli())
	return now > n.ExpiresAt
}

// MarkAsRead 标记通知为已读
func (n *Notification) MarkAsRead() {
	n.Status = NotificationStatusRead
	n.ReadAt = uint64(time.Now().UnixMilli())
}

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	Type     string
	Title    string
	Content  string
	Priority string
}

// NotificationData 通知数据结构
type NotificationData struct {
	ApplicantName   string
	ApproverName    string
	SecretItemName  string
	Reason          string
	RejectReason    string
	ValidUntil      string
	ExpiresIn       string
	MaintenanceTime string
	Duration        string
	AlertDetails    string
}
