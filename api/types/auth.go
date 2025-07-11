package types

import "github.com/akinoccc/hysaif/api/models"

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
