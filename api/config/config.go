package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config 应用配置结构
type Config struct {
	Database DatabaseConfig `json:"database"`
	Security SecurityConfig `json:"security"`
	Server   ServerConfig   `json:"server"`
	WeCom    WeComConfig    `json:"wecom"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `json:"type"` // sqlite, postgres, mysql
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"ssl_mode"`
	Path     string `json:"path"` // SQLite 数据库文件路径
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EncryptionKey string         `json:"encryption_key"`
	JWTSecret     string         `json:"jwt_secret"`
	WebAuthn      WebAuthnConfig `json:"webauthn"`
	Vault         VaultConfig    `json:"vault"`
}

// WebAuthnConfig WebAuthn 配置
type WebAuthnConfig struct {
	RPDisplayName string   `json:"rp_display_name"` // 网站显示名称
	RPID          string   `json:"rp_id"`           // 网站域名
	RPOrigins     []string `json:"rp_origins"`      // 允许的源地址
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// VaultConfig Vault配置
type VaultConfig struct {
	Enabled   bool           `json:"enabled"`             // 是否启用Vault加密
	Address   string         `json:"address"`             // Vault服务器地址
	Token     string         `json:"token"`               // Vault访问令牌
	KeyName   string         `json:"key_name"`            // Transit引擎中的密钥名称
	MountPath string         `json:"mount_path"`          // Transit引擎挂载路径，默认为"transit"
	Namespace string         `json:"namespace,omitempty"` // Vault命名空间（企业版功能）
	TLSConfig VaultTLSConfig `json:"tls_config"`          // TLS配置
}

// VaultTLSConfig Vault TLS配置
type VaultTLSConfig struct {
	Insecure   bool   `json:"insecure"`    // 是否跳过TLS验证（仅用于开发环境）
	CACert     string `json:"ca_cert"`     // CA证书路径
	ClientCert string `json:"client_cert"` // 客户端证书路径
	ClientKey  string `json:"client_key"`  // 客户端私钥路径
}

// WeComConfig 企业微信配置
type WeComConfig struct {
	Enabled      bool   `json:"enabled"`        // 是否启用企业微信通知
	CorpID       string `json:"corp_id"`        // 企业ID
	AgentID      string `json:"agent_id"`       // 应用ID
	Secret       string `json:"secret"`         // 应用密钥
	RedirectURI  string `json:"redirect_uri"`   // 回调URL
	RobotHookKey string `json:"robot_hook_key"` // 企微机器人 webhook key
}

// AppConfig 全局配置实例
var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("未找到配置文件")
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置
	AppConfig = &Config{}
	if err := json.Unmarshal(data, AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 从环境变量覆盖敏感配置
	loadFromEnv()

	return nil
}

// loadFromEnv 从环境变量加载敏感配置
func loadFromEnv() {
	if key := os.Getenv("SIMS_ENCRYPTION_KEY"); key != "" {
		AppConfig.Security.EncryptionKey = key
	}

	if secret := os.Getenv("SIMS_JWT_SECRET"); secret != "" {
		AppConfig.Security.JWTSecret = secret
	}

	if dbHost := os.Getenv("SIMS_DB_HOST"); dbHost != "" {
		AppConfig.Database.Host = dbHost
	}

	if dbUser := os.Getenv("SIMS_DB_USER"); dbUser != "" {
		AppConfig.Database.UserName = dbUser
	}

	if dbPassword := os.Getenv("SIMS_DB_PASSWORD"); dbPassword != "" {
		AppConfig.Database.Password = dbPassword
	}

	if dbName := os.Getenv("SIMS_DB_NAME"); dbName != "" {
		AppConfig.Database.Database = dbName
	}

	// Vault相关环境变量
	if vaultEnabled := os.Getenv("SIMS_VAULT_ENABLED"); vaultEnabled != "" {
		AppConfig.Security.Vault.Enabled = vaultEnabled == "true"
	}

	if vaultAddress := os.Getenv("SIMS_VAULT_ADDRESS"); vaultAddress != "" {
		AppConfig.Security.Vault.Address = vaultAddress
	}

	if vaultToken := os.Getenv("SIMS_VAULT_TOKEN"); vaultToken != "" {
		AppConfig.Security.Vault.Token = vaultToken
	}

	if vaultKeyName := os.Getenv("SIMS_VAULT_KEY_NAME"); vaultKeyName != "" {
		AppConfig.Security.Vault.KeyName = vaultKeyName
	}

	if vaultMountPath := os.Getenv("SIMS_VAULT_MOUNT_PATH"); vaultMountPath != "" {
		AppConfig.Security.Vault.MountPath = vaultMountPath
	}

	if vaultNamespace := os.Getenv("SIMS_VAULT_NAMESPACE"); vaultNamespace != "" {
		AppConfig.Security.Vault.Namespace = vaultNamespace
	}

	// Vault TLS配置
	if vaultInsecure := os.Getenv("SIMS_VAULT_TLS_INSECURE"); vaultInsecure != "" {
		AppConfig.Security.Vault.TLSConfig.Insecure = vaultInsecure == "true"
	}

	if vaultCACert := os.Getenv("SIMS_VAULT_CA_CERT"); vaultCACert != "" {
		AppConfig.Security.Vault.TLSConfig.CACert = vaultCACert
	}

	if vaultClientCert := os.Getenv("SIMS_VAULT_CLIENT_CERT"); vaultClientCert != "" {
		AppConfig.Security.Vault.TLSConfig.ClientCert = vaultClientCert
	}

	if vaultClientKey := os.Getenv("SIMS_VAULT_CLIENT_KEY"); vaultClientKey != "" {
		AppConfig.Security.Vault.TLSConfig.ClientKey = vaultClientKey
	}
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	switch c.Database.Type {
	case "sqlite":
		return c.Database.Path
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host, c.Database.Port, c.Database.UserName,
			c.Database.Password, c.Database.Database, c.Database.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
			c.Database.UserName, c.Database.Password, c.Database.Host,
			c.Database.Port, c.Database.Database, c.Database.SSLMode)
	default:
		return c.Database.Path
	}
}
