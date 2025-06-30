package webauthn

import (
	"sync"

	"github.com/akinoccc/hysaif/api/config"
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	instance *webauthn.WebAuthn
	once     sync.Once
)

// GetWebAuthn 获取 WebAuthn 实例
func GetWebAuthn() (*webauthn.WebAuthn, error) {
	var initErr error
	once.Do(func() {
		wconfig := &webauthn.Config{
			RPDisplayName: config.AppConfig.Security.WebAuthn.RPDisplayName,
			RPID:          config.AppConfig.Security.WebAuthn.RPID,
			RPOrigins:     config.AppConfig.Security.WebAuthn.RPOrigins,
		}

		instance, initErr = webauthn.New(wconfig)
	})

	return instance, initErr
}
