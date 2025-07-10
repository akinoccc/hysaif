# 查询构建器 (Query Builder)

一个通用的查询封装包，用于简化 API handlers 中重复的查询、过滤和分页逻辑。

## 特性

- 🚀 **自动分页处理** - 自动解析页码和页大小参数
- 🔍 **丰富的过滤器** - 支持字符串、模糊搜索、日期范围等多种过滤器
- 🔗 **关联数据预加载** - 链式调用预加载关联模型
- 📊 **统一响应格式** - 标准化的分页响应结构
- 🎯 **类型安全** - 完整的泛型支持

## 基础用法

### 1. 简单查询

```go
func GetUsers(c *gin.Context) {
    var users []models.User
    pagination, err := query.NewQueryBuilder(models.DB, c, &models.User{}).
        LikeFilter("name", "name").
        StringFilter("role", "role").
        StringFilter("status", "status").
        Preload("Creator", "Updater").
        OrderBy("created_at DESC").
        Execute(&users)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
        return
    }
    
    c.JSON(http.StatusOK, types.ListResponse[models.User]{
        Data:       users,
        Pagination: *pagination,
    })
}
```

### 2. 使用预定义过滤器组合

```go
func GetSecretItems(c *gin.Context) {
    var items []models.SecretItem
    pagination, err := query.NewQueryBuilder(models.DB, c, &models.SecretItem{}).
        ApplySecretItemFilters().
        Preload("Creator", "Updater").
        OrderBy("created_at DESC").
        Execute(&items)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
        return
    }
    
    c.JSON(http.StatusOK, types.ListResponse[models.SecretItem]{
        Data:       items,
        Pagination: *pagination,
    })
}
```

### 3. 使用快速查询函数

```go
func GetUsers(c *gin.Context) {
    response, err := query.QueryUsers(models.DB, c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
        return
    }
    
    c.JSON(http.StatusOK, response)
}
```

## 可用过滤器

### 基础过滤器

| 方法 | 说明 | 示例 |
|------|------|------|
| `StringFilter(field, param)` | 精确匹配字符串 | `StringFilter("status", "status")` |
| `LikeFilter(field, param)` | 模糊搜索 | `LikeFilter("name", "name")` |
| `MultiLikeFilter(fields, param)` | 多字段模糊搜索 | `MultiLikeFilter([]string{"name", "desc"}, "search")` |
| `DateRangeFilter(field, from, to)` | 日期范围过滤 | `DateRangeFilter("created_at", "from", "to")` |
| `InFilterByName(field, param, nameField, table)` | 通过名称查询ID列表 | `InFilterByName("user_id", "user_name", "name", "users")` |

### 特殊过滤器

| 方法 | 说明 |
|------|------|
| `SecretStatusFilter()` | 密钥状态过滤器 (expired/expiring/active) |
| `Where(condition, args...)` | 自定义WHERE条件 |
| `WhereIf(condition, query, args...)` | 条件性WHERE |

### 预定义过滤器组合

| 方法 | 适用模型 | 包含过滤器 |
|------|----------|------------|
| `ApplySecretItemFilters()` | SecretItem | category, type, environment, search, status, date_range, creator |
| `ApplyUserFilters()` | User | name, role, status |
| `ApplyAccessRequestFilters()` | AccessRequest | status, date_range, applicant_name, secret_item_name |
| `ApplyAuditLogFilters()` | AuditLog | action, resource |

## 查询参数说明

### 分页参数

- `page` - 页码，默认 1
- `page_size` - 每页大小，默认 10，最大 100

### 通用过滤参数

- `search` - 多字段模糊搜索
- `created_at_from` - 创建时间开始
- `created_at_to` - 创建时间结束
- `sort_by` - 排序字段

### 密钥项专用参数

- `category` - 分类
- `type` - 类型
- `environment` - 环境
- `status` - 状态 (expired/expiring/active)
- `tag` - 标签
- `creator_name` - 创建者名称

### 用户专用参数

- `name` - 用户名
- `role` - 角色
- `status` - 状态

### 访问申请专用参数

- `status` - 申请状态
- `applicant_name` - 申请人名称
- `secret_item_name` - 密钥项名称

## 扩展自定义过滤器

### 创建新的过滤器

```go
// 自定义状态过滤器
func CustomStatusFilter() *QueryBuilder {
    status := qb.ctx.Query("custom_status")
    switch status {
    case "special":
        return qb.Where("field = ? AND other_field > ?", "value", 100)
    default:
        return qb.StringFilter("status", "custom_status")
    }
}

// 使用
qb := query.NewQueryBuilder(db, ctx, &Model{}).
    CustomStatusFilter().
    Execute(&results)
```

### 添加新的过滤器组合

```go
func (qb *QueryBuilder) ApplyMyModelFilters() *QueryBuilder {
    return qb.
        StringFilter("field1", "param1").
        LikeFilter("field2", "param2").
        DateRangeFilter("created_at", "from", "to")
}
```
