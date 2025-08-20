# Hysaif API 服务

## 启动服务

### 基本启动
```bash
# 使用默认配置文件 config.json
go run main.go

# 或者编译后运行
go build -o hysaif-api
./hysaif-api
```

### 指定配置文件路径
```bash
# 使用 --config 参数指定配置文件
go run main.go --config /path/to/config.json

# 使用简写 -c 参数
go run main.go -c /path/to/config.json

# 编译后运行
./hysaif-api --config /path/to/config.json
./hysaif-api -c /path/to/config.json
```

### 命令行参数说明
- `--config, -c`: 指定配置文件路径（默认: config.json）

### 配置文件示例
参考 `config.example.json` 文件创建你的配置文件。

### 环境变量支持
系统支持通过环境变量覆盖配置文件中的敏感信息：

- `SIMS_ENCRYPTION_KEY`: 加密密钥
- `SIMS_JWT_SECRET`: JWT密钥
- `SIMS_DB_HOST`: 数据库主机
- `SIMS_DB_USER`: 数据库用户名
- `SIMS_DB_PASSWORD`: 数据库密码
- `SIMS_DB_NAME`: 数据库名称
- `SIMS_VAULT_ENABLED`: 是否启用Vault
- `SIMS_VAULT_ADDRESS`: Vault服务器地址
- `SIMS_VAULT_TOKEN`: Vault访问令牌
- `SIMS_WECOM_ENABLED`: 是否启用企业微信
- `SIMS_WECOM_CORP_ID`: 企业微信企业ID
- `SIMS_WECOM_AGENT_ID`: 企业微信应用ID
- `SIMS_WECOM_SECRET`: 企业微信应用密钥

### 部署建议
在生产环境中，建议：
1. 使用绝对路径指定配置文件
2. 通过环境变量设置敏感信息
3. 确保配置文件权限设置正确
4. 使用systemd或supervisor管理服务进程
