# 企业敏感信息管理系统（SIMS）项目规则

## 项目概述

这是一个企业级敏感信息管理系统（Sensitive Information Management System, SIMS），用于安全存储和管理企业的各类敏感信息，包括密码、API密钥、SSH密钥、证书等。系统采用前后端分离架构，前端使用Vue3 + TypeScript，后端使用Go + Gin框架。

## 项目结构

```
hysaif/
├── api/                    # 后端API服务
│   ├── handlers/          # 路由处理器
│   ├── middleware/        # 中间件
│   ├── models/           # 数据模型
│   ├── types/            # 类型定义
│   ├── utils/            # 工具函数
│   ├── main.go           # 主入口文件
│   └── go.mod            # Go模块依赖
├── web/                   # 前端Web应用
│   ├── src/
│   │   ├── api/          # API接口定义
│   │   ├── components/   # Vue组件
│   │   ├── views/        # 页面视图
│   │   ├── router/       # 路由配置
│   │   ├── stores/       # 状态管理
│   │   └── utils/        # 工具函数
│   ├── package.json      # 前端依赖
│   └── vite.config.ts    # Vite配置
└── README.md             # 项目文档
```

## 技术栈

### 后端（API）
- **语言**: Go 1.24
- **框架**: Gin (Web框架)
- **数据库**: SQLite (开发环境) / PostgreSQL (生产环境)
- **ORM**: GORM
- **认证**: JWT
- **加密**: golang.org/x/crypto
- **CORS**: gin-contrib/cors

### 前端（Web）
- **框架**: Vue 3.5.13 + TypeScript
- **构建工具**: Vite 6.3.5
- **路由**: Vue Router 4.2.5
- **状态管理**: Pinia 3.0.2
- **UI组件**: Reka UI 2.3.0
- **样式**: Tailwind CSS 4.1.8
- **表单验证**: Vee-Validate 4.15.0 + Zod 3.25.42
- **HTTP客户端**: Axios 1.6.0
- **图标**: Lucide Vue Next 0.511.0

## 核心功能模块

### 1. 用户管理
- 用户认证与授权
- 角色权限控制（super_admin, sec_mgr, dev, auditor）
- 用户配置文件管理

### 2. 敏感信息管理
支持的信息类型：
- **密码**: 用户名+密码组合
- **API密钥**: 各种服务的API密钥
- **访问密钥**: AWS AccessKey/SecretKey等
- **SSH密钥**: 公钥/私钥对
- **证书**: TLS/SSL证书
- **令牌**: GitHub Token、JWT等
- **自定义**: 键值对形式的自定义数据

### 3. 安全特性
- AES-256加密存储
- JWT身份认证
- 操作审计日志
- 权限控制
- 数据过期管理

### 4. 审计日志
- 用户操作记录
- 资源访问追踪
- IP地址和用户代理记录
- 详细操作信息

## 开发规范

### 代码风格

#### Go代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 结构体字段使用驼峰命名
- 包名使用小写字母
- 错误处理必须显式检查

#### TypeScript/Vue代码规范
- 使用TypeScript严格模式
- 组件名使用PascalCase
- 文件名使用kebab-case
- 使用Composition API
- 遵循Vue 3最佳实践

### 数据库设计

#### 核心表结构
1. **users**: 用户表
   - id (主键)
   - username (唯一)
   - password (加密)
   - email
   - role (角色)
   - status (状态)

2. **secret_items**: 敏感信息表
   - id (主键)
   - name (名称)
   - type (类型)
   - category (分类)
   - encrypted_data (加密数据)
   - expires_at (过期时间)
   - created_by/updated_by (创建者/更新者)

3. **audit_logs**: 审计日志表
   - id (主键)
   - user_id (用户ID)
   - action (操作类型)
   - resource (资源类型)
   - details (详细信息)
   - ip_address (IP地址)

### API设计规范

#### RESTful API
- 使用标准HTTP方法（GET, POST, PUT, DELETE）
- 统一的响应格式
- 适当的HTTP状态码
- API版本控制（/api/v1/）

#### 认证授权
- JWT Token认证
- Bearer Token格式
- 中间件统一处理认证
- 基于角色的访问控制

### 前端开发规范

#### 组件设计
- 单一职责原则
- 可复用组件抽象
- Props类型定义
- 事件命名规范

#### 状态管理
- 使用Pinia进行状态管理
- 按功能模块划分store
- 持久化敏感状态

#### 路由设计
- 嵌套路由结构
- 路由守卫认证
- 动态菜单生成
- 权限控制

## 安全要求

### 数据安全
- 敏感数据必须加密存储
- 传输过程使用HTTPS
- 密码使用bcrypt哈希
- 定期密钥轮换

### 访问控制
- 最小权限原则
- 会话超时控制
- 操作审计记录
- 异常访问检测

### 代码安全
- 输入验证和清理
- SQL注入防护
- XSS攻击防护
- CSRF保护

## 部署要求

### 开发环境
- Go 1.24+
- Node.js 18+
- SQLite 3
- Git

### 生产环境
- 容器化部署（Docker）
- 负载均衡
- 数据库集群
- 监控告警
- 备份策略

## 测试要求

### 后端测试
- 单元测试覆盖率 > 80%
- 集成测试
- API测试
- 安全测试

### 前端测试
- 组件单元测试
- E2E测试
- 用户交互测试
- 兼容性测试

## 文档要求

### 代码文档
- 函数和方法注释
- 复杂逻辑说明
- API文档
- 数据库设计文档

### 用户文档
- 安装部署指南
- 用户操作手册
- 管理员指南
- 故障排除指南

## 版本控制

### Git工作流
- 使用Git Flow分支模型
- 功能分支开发
- 代码审查机制
- 提交信息规范

### 发布管理
- 语义化版本控制
- 变更日志维护
- 发布标签管理
- 回滚策略

## 监控和维护

### 系统监控
- 应用性能监控
- 错误日志收集
- 资源使用监控
- 安全事件监控

### 维护策略
- 定期安全更新
- 性能优化
- 数据备份验证
- 灾难恢复演练

## 合规要求

### 数据保护
- GDPR合规
- 数据分类标记
- 数据保留策略
- 数据删除机制

### 审计要求
- 操作日志完整性
- 访问记录可追溯
- 合规报告生成
- 第三方审计支持

---

**注意**: 本项目涉及敏感信息管理，开发过程中必须严格遵循安全开发生命周期（SDLC）要求，确保系统安全性和合规性。