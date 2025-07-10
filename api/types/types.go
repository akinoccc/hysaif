package types

import (
	"github.com/akinoccc/hysaif/api/models"
)

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// 认证相关类型
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// WebAuthn 相关类型
type WebAuthnBeginRegistrationRequest struct {
	CredentialName string `json:"credential_name" binding:"required,min=1,max=50"` // 凭证名称
}

type WebAuthnBeginRegistrationResponse struct {
	Options interface{} `json:"options"` // WebAuthn 注册选项
}

type WebAuthnFinishRegistrationRequest struct {
	Response       interface{} `json:"response" binding:"required"`                     // WebAuthn 注册响应
	CredentialName string      `json:"credential_name" binding:"required,min=1,max=50"` // 凭证名称
}

type WebAuthnBeginLoginResponse struct {
	Options interface{} `json:"options"` // WebAuthn 登录选项
}

type WebAuthnFinishLoginRequest struct {
	Response interface{} `json:"response" binding:"required"` // WebAuthn 登录响应
}

type WebAuthnCredentialResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreatedAt  uint64 `json:"created_at"`
	LastUsedAt uint64 `json:"last_used_at"`
}

type WebAuthnCredentialListResponse struct {
	Credentials []WebAuthnCredentialResponse `json:"credentials"`
}

// 企业微信登录相关类型
type WeWorkAuthRequest struct {
	Code  string `json:"code" binding:"required"` // 企业微信授权码
	State string `json:"state"`                   // 状态参数
}

type WeWorkUserInfo struct {
	UserID     string `json:"userid"`     // 企业微信用户ID
	Name       string `json:"name"`       // 用户姓名
	Email      string `json:"email"`      // 邮箱
	Mobile     string `json:"mobile"`     // 手机号
	Gender     string `json:"gender"`     // 性别
	Avatar     string `json:"avatar"`     // 头像
	Department []int  `json:"department"` // 部门ID列表
	Position   string `json:"position"`   // 职位
	CorpID     string `json:"corp_id"`    // 企业ID
}

type WeWorkAuthURLResponse struct {
	AuthURL string `json:"auth_url"` // 企业微信授权URL
}

// 用户相关类型
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

type CreateUserRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=50"`
	Password    string   `json:"password" binding:"required,min=8,max=128"`
	Email       string   `json:"email" binding:"required,email"`
	Role        string   `json:"role" binding:"required,oneof=super_admin sec_mgr dev auditor"`
	Permissions []string `json:"permissions"`
}

type UpdateUserRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=2,max=50"`
	Email       string   `json:"email" binding:"omitempty,email"`
	Role        string   `json:"role" binding:"omitempty,oneof=super_admin sec_mgr dev auditor"`
	Status      string   `json:"status" binding:"omitempty,oneof=active disabled locked expired"`
	Permissions []string `json:"permissions"`
	Password    string   `json:"password" binding:"omitempty,min=8,max=128"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

type UserListParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"`
	Role     string `form:"role" binding:"omitempty,oneof=super_admin sec_mgr admin user"`
	Status   string `form:"status" binding:"omitempty,oneof=active disabled locked expired"`
	SortBy   string `form:"sort_by"`
	SortDesc bool   `form:"sort_desc"`
}

// 信息项相关类型

type PostItemRequest struct {
	Name        string                `json:"name" binding:"required,min=1,max=100"`
	Type        string                `json:"type" binding:"required,oneof=password api_key access_key ssh_key certificate token kv"` // password, api_key, access_key, ssh_key, certificate, token, kv
	Description string                `json:"description,omitempty" binding:"max=500"`
	Category    string                `json:"category" binding:"required,min=1,max=50"`
	Environment string                `json:"environment" binding:"required,oneof=development staging production test"`
	Tags        []string              `json:"tags,omitempty" gorm:"type:text;serializer:json"`
	Data        models.SecretItemData `json:"data" binding:"required" gorm:"type:text;serializer:json"`
	ExpiresAt   uint64                `json:"expires_at,omitempty"`
}

// 通用响应类型
type ListResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Error   string              `json:"error"`
	Details map[string][]string `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// 访问申请相关类型
type CreateAccessRequestRequest struct {
	SecretItemID string `json:"secret_item_id" binding:"required"`
	Reason       string `json:"reason" binding:"required,min=5,max=500"`
}

type ApproveAccessRequestRequest struct {
	ValidDuration int    `json:"valid_duration" binding:"required,min=1,max=8760"` // 有效时长（小时），最大365天
	Note          string `json:"note" binding:"max=500"`                           // 审批备注
}

type RejectAccessRequestRequest struct {
	Reason string `json:"reason" binding:"required,min=5,max=500"` // 拒绝理由
}

type RevokeAccessRequestRequest struct {
	Reason string `json:"reason" binding:"required,min=5,max=500"` // 作废理由
}

type AccessRequestListParams struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`
	PageSize     int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status       string `form:"status" binding:"omitempty,oneof=pending approved rejected expired revoked"` // pending, approved, rejected, expired, revoked
	ApplicantID  string `form:"applicant_id"`                                                               // 申请人ID
	SecretItemID string `form:"secret_item_id"`                                                             // 密钥项ID
	SortBy       string `form:"sort_by"`
	SortDesc     bool   `form:"sort_desc"`
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
