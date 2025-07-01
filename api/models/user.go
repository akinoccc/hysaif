package models

import (
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/akinoccc/hysaif/api/packages/permission"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户角色常量
const (
	RoleSuperAdmin = "super_admin" // 超级管理员
	RoleSecMgr     = "sec_mgr"     // 安全管理员
	RoleDev        = "dev"         // 开发人员
	RoleAuditor    = "auditor"     // 审计员
	RoleBot        = "bot"         // 服务账号
)

// 用户状态常量
const (
	StatusActive   = "active"   // 活跃
	StatusDisabled = "disabled" // 禁用
	StatusLocked   = "locked"   // 锁定
	StatusExpired  = "expired"  // 过期
)

// User 用户模型
type User struct {
	ModelBase
	Name            string   `json:"name"`
	Password        string   `json:"-" gorm:"not null"`
	Email           string   `json:"email" gorm:"type:varchar(191);uniqueIndex;not null"`
	Role            string   `json:"role" gorm:"default:'dev'"`      // super_admin, sec_mgr, dev, auditor, bot
	Status          string   `json:"status" gorm:"default:'active'"` // active, disabled, locked, expired
	LastLoginAt     uint64   `json:"last_login_at"`
	LastLoginIP     string   `json:"last_login_ip"`
	FailedAttempts  int      `json:"failed_attempts" gorm:"default:0"`             // 登录失败次数
	Permissions     []string `json:"permissions" gorm:"type:text;serializer:json"` // 特殊权限列表
	CreatedByUserID string   `json:"-" gorm:"index"`
	UpdatedByUserID string   `json:"-" gorm:"index"`

	// 企业微信相关字段
	WeWorkUserID string `json:"wework_user_id" gorm:"index"`       // 企业微信用户ID
	WeWorkCorpID string `json:"wework_corp_id"`                    // 企业微信企业ID
	Avatar       string `json:"avatar"`                            // 头像URL
	Mobile       string `json:"mobile"`                            // 手机号
	Department   string `json:"department"`                        // 部门
	Position     string `json:"position"`                          // 职位
	LoginType    string `json:"login_type" gorm:"default:'local'"` // 登录类型：local, wework

	Creator *User `json:"creator" gorm:"foreignKey:CreatedByUserID;references:ID"`
	Updater *User `json:"updater" gorm:"foreignKey:UpdatedByUserID;references:ID"`

	// WebAuthn 凭证（不在JSON中序列化）
	Credentials []WebAuthnCredential `json:"-"`
}

// BeforeCreate 密码加密
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return
}

// BeforeUpdate 密码加密
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Password == "" {
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return
}

// HasPermission 检查用户是否拥有指定权限（使用Casbin）
func (u *User) HasPermission(resource, action string) bool {
	// 超级管理员拥有所有权限
	if u.Role == RoleSuperAdmin {
		return true
	}

	// 使用Casbin检查权限
	casbinManager := permission.GetCasbinManager(DB)
	hasPermission := casbinManager.CheckPermission(u.Role, resource, action)

	return hasPermission
}

// IsActive 检查用户是否处于活跃状态
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// ValidateUserRole 验证用户角色是否有效
func ValidateUserRole(role string) error {
	switch role {
	case RoleSuperAdmin, RoleSecMgr, RoleDev, RoleAuditor, RoleBot:
		return nil
	default:
		return errors.New("无效的用户角色")
	}
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// WebAuthn 接口实现

// WebAuthnID 返回用户的 WebAuthn ID
func (u *User) WebAuthnID() []byte {
	return []byte(u.ID)
}

// WebAuthnName 返回用户名
func (u *User) WebAuthnName() string {
	return u.Email
}

// WebAuthnDisplayName 返回显示名称
func (u *User) WebAuthnDisplayName() string {
	return u.Name
}

// WebAuthnCredentials 返回用户的所有凭证
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	var credentials []webauthn.Credential

	// 加载用户的 WebAuthn 凭证
	if len(u.Credentials) == 0 {
		DB.Where("user_id = ?", u.ID).Find(&u.Credentials)
	}

	for _, cred := range u.Credentials {
		if webauthnCred, err := cred.ToWebAuthnCredential(); err == nil {
			credentials = append(credentials, *webauthnCred)
		}
	}

	return credentials
}

// AddWebAuthnCredential 添加 WebAuthn 凭证
func (u *User) AddWebAuthnCredential(credential *webauthn.Credential, name string) error {
	webauthnCred := FromWebAuthnCredential(u.ID, credential, name)
	log.Printf("webauthnCred: %+v", webauthnCred.CredentialName)
	return DB.Create(webauthnCred).Error
}

// UpdateWebAuthnCredential 更新 WebAuthn 凭证（主要用于更新签名计数）
func (u *User) UpdateWebAuthnCredential(credential *webauthn.Credential) error {
	credentialID := base64.RawURLEncoding.EncodeToString(credential.ID)

	return DB.Model(&WebAuthnCredential{}).
		Where("user_id = ? AND credential_id = ?", u.ID, credentialID).
		Updates(map[string]interface{}{
			"sign_counter": credential.Authenticator.SignCount,
			"last_used_at": time.Now(),
		}).Error
}
