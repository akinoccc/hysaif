package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 申请状态常量
const (
	RequestStatusPending  = "pending"  // 待审批
	RequestStatusApproved = "approved" // 已批准
	RequestStatusRejected = "rejected" // 已拒绝
	RequestStatusExpired  = "expired"  // 已过期
	RequestStatusRevoked  = "revoked"  // 已撤销
)

// AccessRequest 访问申请模型
type AccessRequest struct {
	ModelBase
	SecretItemID string `json:"-" gorm:"not null"`               // 申请访问的密钥项ID
	ApplicantID  string `json:"-" gorm:"not null"`               // 申请人ID
	Reason       string `json:"reason" gorm:"not null"`          // 申请理由
	Status       string `json:"status" gorm:"default:'pending'"` // 申请状态
	ApprovedByID string `json:"-" gorm:"index"`                  // 审批人ID
	ApprovedAt   uint64 `json:"approved_at"`                     // 审批时间
	Note         string `json:"note"`                            // 备注
	RejectReason string `json:"reject_reason"`                   // 拒绝理由
	ValidFrom    uint64 `json:"valid_from"`                      // 有效期开始时间
	ValidUntil   uint64 `json:"valid_until"`                     // 有效期结束时间
	AccessCount  int    `json:"access_count" gorm:"default:0"`   // 访问次数
	LastAccessed uint64 `json:"last_accessed"`                   // 最后访问时间

	// 关联
	SecretItem SecretItem `json:"secret_item" gorm:"foreignKey:SecretItemID"`
	Applicant  User       `json:"applicant" gorm:"foreignKey:ApplicantID"`
	Approver   User       `json:"approver" gorm:"foreignKey:ApprovedByID"`
}

// BeforeCreate 钩子函数，在创建记录之前设置ID
func (ar *AccessRequest) BeforeCreate(tx *gorm.DB) (err error) {
	ar.ID = uuid.New().String()
	return
}

// IsValid 检查申请是否有效
func (ar *AccessRequest) IsValid() bool {
	if ar.Status != RequestStatusApproved {
		return false
	}

	now := uint64(time.Now().UnixMilli())
	return now >= ar.ValidFrom && now <= ar.ValidUntil
}

// IsExpired 检查申请是否已过期
func (ar *AccessRequest) IsExpired() bool {
	if ar.Status != RequestStatusApproved {
		return false
	}

	now := uint64(time.Now().UnixMilli())
	return now > ar.ValidUntil
}

// CanAccess 检查是否可以访问
func (ar *AccessRequest) CanAccess() bool {
	return ar.IsValid() && !ar.IsExpired()
}
