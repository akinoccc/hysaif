package models

import (
	"encoding/base64"
	"encoding/json"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WebAuthnCredential WebAuthn 凭证模型
type WebAuthnCredential struct {
	ModelBase
	UserID          string `json:"-" gorm:"type:varchar(36);index"`
	CredentialID    string `json:"credential_id" gorm:"type:varchar(255);uniqueIndex"` // Base64 编码的凭证ID
	PublicKey       string `json:"public_key" gorm:"type:text"`                        // Base64 编码的公钥
	AttestationType string `json:"attestation_type"`
	Transport       string `json:"transport" gorm:"type:text"` // JSON 编码的传输方式数组
	Flags           uint8  `json:"flags"`
	SignCounter     uint32 `json:"sign_counter"`
	AAGUID          string `json:"aaguid"`
	CredentialName  string `json:"credential_name"` // 用户自定义的凭证名称
	LastUsedAt      uint64 `json:"last_used_at"`

	// 关联
	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

// BeforeCreate 创建前钩子
func (w *WebAuthnCredential) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New().String()
	return nil
}

// ToWebAuthnCredential 转换为 webauthn.Credential
func (w *WebAuthnCredential) ToWebAuthnCredential() (*webauthn.Credential, error) {
	credentialID, err := base64.RawURLEncoding.DecodeString(w.CredentialID)
	if err != nil {
		return nil, err
	}

	publicKey, err := base64.RawURLEncoding.DecodeString(w.PublicKey)
	if err != nil {
		return nil, err
	}

	aaguid, err := base64.RawURLEncoding.DecodeString(w.AAGUID)
	if err != nil {
		return nil, err
	}

	// 解码 Transport
	var transport []protocol.AuthenticatorTransport
	if w.Transport != "" && w.Transport != "[]" {
		if err := json.Unmarshal([]byte(w.Transport), &transport); err != nil {
			// 如果解码失败，使用空数组
			transport = []protocol.AuthenticatorTransport{}
		}
	}

	credential := &webauthn.Credential{
		ID:              credentialID,
		PublicKey:       publicKey,
		AttestationType: w.AttestationType,
		Transport:       transport,
		Flags: webauthn.CredentialFlags{
			UserPresent:    (w.Flags & 0x01) != 0,
			UserVerified:   (w.Flags & 0x04) != 0,
			BackupEligible: (w.Flags & 0x08) != 0,
			BackupState:    (w.Flags & 0x10) != 0,
		},
		Authenticator: webauthn.Authenticator{
			SignCount: w.SignCounter,
			AAGUID:    aaguid,
		},
	}

	return credential, nil
}

// FromWebAuthnCredential 从 webauthn.Credential 创建
func FromWebAuthnCredential(userID string, credential *webauthn.Credential, name string) *WebAuthnCredential {
	flags := uint8(0)
	if credential.Flags.UserPresent {
		flags |= 0x01
	}
	if credential.Flags.UserVerified {
		flags |= 0x04
	}
	if credential.Flags.BackupEligible {
		flags |= 0x08
	}
	if credential.Flags.BackupState {
		flags |= 0x10
	}

	// 将 Transport 数组编码为 JSON
	transportJSON := "[]"
	if len(credential.Transport) > 0 {
		if data, err := json.Marshal(credential.Transport); err == nil {
			transportJSON = string(data)
		}
	}

	return &WebAuthnCredential{
		UserID:          userID,
		CredentialID:    base64.RawURLEncoding.EncodeToString(credential.ID),
		PublicKey:       base64.RawURLEncoding.EncodeToString(credential.PublicKey),
		AttestationType: credential.AttestationType,
		Transport:       transportJSON,
		Flags:           flags,
		SignCounter:     credential.Authenticator.SignCount,
		AAGUID:          base64.RawURLEncoding.EncodeToString(credential.Authenticator.AAGUID),
		CredentialName:  name,
	}
}
