package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/akinoccc/hysaif/api/config"
	vaultapi "github.com/hashicorp/vault/api"
)

var (
	vaultClient *vaultapi.Client
	vaultOnce   sync.Once
	initErr     error
)

// getEncryptionKey 获取加密密钥（用于向后兼容）
func getEncryptionKey() []byte {
	key := config.AppConfig.Security.EncryptionKey
	if len(key) != 32 {
		panic(fmt.Sprintf("加密密钥长度必须为32字节，当前长度: %d", len(key)))
	}
	return []byte(key)
}

// getVaultClient 获取Vault客户端实例
func getVaultClient() (*vaultapi.Client, error) {
	vaultOnce.Do(func() {
		if !config.AppConfig.Security.Vault.Enabled {
			return
		}

		// 创建Vault客户端配置
		vaultConfig := vaultapi.DefaultConfig()
		vaultConfig.Address = config.AppConfig.Security.Vault.Address

		// 配置TLS
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.AppConfig.Security.Vault.TLSConfig.Insecure,
		}

		// 设置CA证书
		if config.AppConfig.Security.Vault.TLSConfig.CACert != "" {
			vaultConfig.ConfigureTLS(&vaultapi.TLSConfig{
				CACert: config.AppConfig.Security.Vault.TLSConfig.CACert,
			})
		}

		// 设置客户端证书
		if config.AppConfig.Security.Vault.TLSConfig.ClientCert != "" &&
			config.AppConfig.Security.Vault.TLSConfig.ClientKey != "" {
			vaultConfig.ConfigureTLS(&vaultapi.TLSConfig{
				ClientCert: config.AppConfig.Security.Vault.TLSConfig.ClientCert,
				ClientKey:  config.AppConfig.Security.Vault.TLSConfig.ClientKey,
			})
		}

		// 设置自定义TLS配置
		if config.AppConfig.Security.Vault.TLSConfig.Insecure {
			vaultConfig.HttpClient.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
		}

		// 创建客户端
		client, err := vaultapi.NewClient(vaultConfig)
		if err != nil {
			initErr = fmt.Errorf("创建Vault客户端失败: %w", err)
			return
		}

		// 设置认证令牌
		client.SetToken(config.AppConfig.Security.Vault.Token)

		// 设置命名空间（企业版功能）
		if config.AppConfig.Security.Vault.Namespace != "" {
			client.SetNamespace(config.AppConfig.Security.Vault.Namespace)
		}

		vaultClient = client

		// 测试连接和权限
		if err := testVaultConnection(client); err != nil {
			initErr = fmt.Errorf("Vault连接测试失败: %w", err)
			log.Printf("警告: Vault连接失败，将回退到AES加密: %v", err)
			return
		}

		log.Println("Vault客户端初始化成功")
	})

	return vaultClient, initErr
}

// testVaultConnection 测试Vault连接和权限
func testVaultConnection(client *vaultapi.Client) error {
	// 检查Health状态
	resp, err := client.Sys().Health()
	if err != nil {
		return fmt.Errorf("无法获取Vault健康状态: %w", err)
	}

	if !resp.Initialized {
		return fmt.Errorf("Vault服务器未初始化")
	}

	if resp.Sealed {
		return fmt.Errorf("Vault服务器已密封")
	}

	// 检查Transit引擎是否挂载
	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return fmt.Errorf("无法列出挂载点: %w", err)
	}

	mountPath := strings.TrimSuffix(config.AppConfig.Security.Vault.MountPath, "/") + "/"
	if _, exists := mounts[mountPath]; !exists {
		return fmt.Errorf("Transit引擎未在路径 '%s' 挂载", mountPath)
	}

	// 验证密钥是否存在，如果不存在则创建
	if err := ensureTransitKey(client); err != nil {
		return fmt.Errorf("无法确保Transit密钥存在: %w", err)
	}

	return nil
}

// ensureTransitKey 确保Transit密钥存在
func ensureTransitKey(client *vaultapi.Client) error {
	keyName := config.AppConfig.Security.Vault.KeyName
	mountPath := config.AppConfig.Security.Vault.MountPath

	// 检查密钥是否存在
	path := fmt.Sprintf("%s/keys/%s", mountPath, keyName)
	_, err := client.Logical().Read(path)
	if err == nil {
		// 密钥存在
		return nil
	}

	// 尝试创建密钥
	createPath := fmt.Sprintf("%s/keys/%s", mountPath, keyName)
	data := map[string]interface{}{
		"type": "aes256-gcm96",
	}

	_, err = client.Logical().Write(createPath, data)
	if err != nil {
		return fmt.Errorf("创建Transit密钥失败: %w", err)
	}

	log.Printf("成功创建Transit密钥: %s", keyName)
	return nil
}

// Encrypt 加密数据 - 优先使用Vault，失败时回退到AES
func Encrypt(data []byte) (string, error) {
	// 尝试使用Vault加密
	if config.AppConfig.Security.Vault.Enabled {
		if encrypted, err := encryptWithVault(data); err == nil {
			return "vault:" + encrypted, nil
		} else {
			log.Printf("Vault加密失败，回退到AES: %v", err)
		}
	}

	// 回退到AES加密
	return encryptWithAES(data)
}

// Decrypt 解密数据 - 自动检测加密方式并解密
func Decrypt(encryptedData string) ([]byte, error) {
	// 检查是否为Vault加密的数据
	if strings.HasPrefix(encryptedData, "vault:") {
		vaultData := strings.TrimPrefix(encryptedData, "vault:")
		return decryptWithVault(vaultData)
	}

	// 向后兼容：解密AES加密的数据
	return decryptWithAES(encryptedData)
}

// encryptWithVault 使用Vault Transit引擎加密
func encryptWithVault(data []byte) (string, error) {
	client, err := getVaultClient()
	if err != nil {
		return "", err
	}

	keyName := config.AppConfig.Security.Vault.KeyName
	mountPath := config.AppConfig.Security.Vault.MountPath

	// 对数据进行base64编码
	encodedData := base64.StdEncoding.EncodeToString(data)

	// 构造加密请求
	path := fmt.Sprintf("%s/encrypt/%s", mountPath, keyName)
	requestData := map[string]interface{}{
		"plaintext": encodedData,
	}

	// 执行加密
	resp, err := client.Logical().Write(path, requestData)
	if err != nil {
		return "", fmt.Errorf("Vault加密请求失败: %w", err)
	}

	if resp == nil || resp.Data == nil {
		return "", fmt.Errorf("Vault加密响应为空")
	}

	ciphertext, ok := resp.Data["ciphertext"].(string)
	if !ok {
		return "", fmt.Errorf("无效的Vault加密响应格式")
	}

	return ciphertext, nil
}

// decryptWithVault 使用Vault Transit引擎解密
func decryptWithVault(encryptedData string) ([]byte, error) {
	client, err := getVaultClient()
	if err != nil {
		return nil, err
	}

	keyName := config.AppConfig.Security.Vault.KeyName
	mountPath := config.AppConfig.Security.Vault.MountPath

	// 构造解密请求
	path := fmt.Sprintf("%s/decrypt/%s", mountPath, keyName)
	requestData := map[string]interface{}{
		"ciphertext": encryptedData,
	}

	// 执行解密
	resp, err := client.Logical().Write(path, requestData)
	if err != nil {
		return nil, fmt.Errorf("Vault解密请求失败: %w", err)
	}

	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("Vault解密响应为空")
	}

	plaintext, ok := resp.Data["plaintext"].(string)
	if !ok {
		return nil, fmt.Errorf("无效的Vault解密响应格式")
	}

	// 解码base64
	data, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}

	return data, nil
}

// encryptWithAES 使用AES加密（向后兼容）
func encryptWithAES(data []byte) (string, error) {
	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptWithAES 使用AES解密（向后兼容）
func decryptWithAES(encryptedData string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(getEncryptionKey())
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("加密数据长度不足")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// IsVaultEnabled 检查是否启用了Vault加密
func IsVaultEnabled() bool {
	return config.AppConfig.Security.Vault.Enabled
}

// GetEncryptionMethod 获取当前使用的加密方法
func GetEncryptionMethod() string {
	if IsVaultEnabled() {
		if _, err := getVaultClient(); err == nil {
			return "Vault Transit Engine"
		}
	}
	return "AES-256-GCM"
}
