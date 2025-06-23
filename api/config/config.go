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
	EncryptionKey string `json:"encryption_key"`
	JWTSecret     string `json:"jwt_secret"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
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
