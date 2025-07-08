# HashiCorp Vault 加密集成

## 概述

本系统已集成HashiCorp Vault来加强密文安全性，使用Vault的Transit引擎进行企业级加密。

## 主要特性

- **企业级加密**: 使用Vault Transit引擎
- **向后兼容**: 自动处理现有AES加密数据
- **智能回退**: Vault不可用时自动回退到AES加密
- **零停机升级**: 无需数据迁移即可启用

## 快速配置

### 1. 环境变量配置

```bash
# 启用Vault加密
export SIMS_VAULT_ENABLED=true
export SIMS_VAULT_ADDRESS=https://vault.example.com:8200
export SIMS_VAULT_TOKEN=s.your-vault-token-here
export SIMS_VAULT_KEY_NAME=sims-encrypt-key
export SIMS_VAULT_MOUNT_PATH=transit
```

### 2. Vault服务器设置

```bash
# 启用Transit引擎
vault secrets enable transit

# 创建加密密钥
vault write transit/keys/sims-encrypt-key type=aes256-gcm96
```

### 3. 权限策略

```hcl
# sims-policy.hcl
path "transit/encrypt/sims-encrypt-key" {
  capabilities = ["update"]
}
path "transit/decrypt/sims-encrypt-key" {
  capabilities = ["update"]
}
path "sys/health" {
  capabilities = ["read"]
}
```

## 开发测试

使用Docker快速启动开发环境：

```bash
# 启动Vault开发服务器
docker run --cap-add=IPC_LOCK -d --name=dev-vault \
  -p 8200:8200 \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=root' \
  vault:latest

# 配置Transit引擎
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=root
vault secrets enable transit
vault write transit/keys/sims-encrypt-key type=aes256-gcm96
```

## 工作原理

1. **加密流程**: 优先使用Vault，失败时自动回退到AES
2. **解密流程**: 自动检测加密方式并使用相应方法解密
3. **数据格式**: Vault加密数据带有"vault:"前缀标识

## 故障排除

- **连接失败**: 检查Vault地址和网络连接
- **权限错误**: 验证令牌策略权限
- **引擎未挂载**: 执行 `vault secrets enable transit`

## 安全建议

- 使用最小权限策略
- 启用TLS加密通信
- 定期轮换令牌和密钥
- 监控加密操作审计日志

更多详细信息请参考HashiCorp Vault官方文档。 