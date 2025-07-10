package query

import (
	"math"
	"strconv"
	"time"

	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QueryBuilder 查询构建器
type QueryBuilder struct {
	db     *gorm.DB
	query  *gorm.DB
	ctx    *gin.Context
	page   int
	limit  int
	offset int
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(db *gorm.DB, ctx *gin.Context, model interface{}) *QueryBuilder {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	return &QueryBuilder{
		db:     db,
		query:  db.Model(model),
		ctx:    ctx,
		page:   page,
		limit:  pageSize,
		offset: offset,
	}
}

// Where 添加条件
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Where(condition, args...)
	return qb
}

// WhereIf 条件添加Where
func (qb *QueryBuilder) WhereIf(condition bool, query string, args ...interface{}) *QueryBuilder {
	if condition {
		qb.query = qb.query.Where(query, args...)
	}
	return qb
}

// StringFilter 字符串过滤器
func (qb *QueryBuilder) StringFilter(field, param string) *QueryBuilder {
	value := qb.ctx.Query(param)
	return qb.WhereIf(value != "", field+" = ?", value)
}

// LikeFilter 模糊搜索过滤器
func (qb *QueryBuilder) LikeFilter(field, param string) *QueryBuilder {
	value := qb.ctx.Query(param)
	return qb.WhereIf(value != "", field+" LIKE ?", "%"+value+"%")
}

// MultiLikeFilter 多字段模糊搜索
func (qb *QueryBuilder) MultiLikeFilter(fields []string, param string) *QueryBuilder {
	value := qb.ctx.Query(param)
	if value != "" && len(fields) > 0 {
		conditions := make([]string, len(fields))
		args := make([]interface{}, len(fields))
		for i, field := range fields {
			conditions[i] = field + " LIKE ?"
			args[i] = "%" + value + "%"
		}
		// 构建 OR 查询
		query := conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " OR " + conditions[i]
		}
		qb.query = qb.query.Where(query, args...)
	}
	return qb
}

// DateRangeFilter 日期范围过滤器
func (qb *QueryBuilder) DateRangeFilter(field, fromParam, toParam string) *QueryBuilder {
	from := qb.ctx.Query(fromParam)
	to := qb.ctx.Query(toParam)

	qb.WhereIf(from != "", field+" >= ?", from)
	qb.WhereIf(to != "", field+" <= ?", to)

	return qb
}

// InFilterByName 通过名称查询IN过滤器
func (qb *QueryBuilder) InFilterByName(field, param, nameField, table string) *QueryBuilder {
	value := qb.ctx.Query(param)
	if value != "" {
		var ids []string
		qb.db.Table(table).Where(nameField+" LIKE ?", "%"+value+"%").Pluck("id", &ids)
		if len(ids) > 0 {
			qb.query = qb.query.Where(field+" IN (?)", ids)
		}
	}
	return qb
}

// SecretStatusFilter 密钥状态过滤器
func (qb *QueryBuilder) SecretStatusFilter() *QueryBuilder {
	status := qb.ctx.Query("status")
	switch status {
	case "expired":
		qb.query = qb.query.Where("expires_at < ?", time.Now().UnixMilli())
	case "expiring":
		qb.query = qb.query.Where("expires_at > ? AND expires_at < ?",
			time.Now().UnixMilli(),
			time.Now().AddDate(0, 0, 7).UnixMilli())
	case "active":
		qb.query = qb.query.Where("expires_at > ?", time.Now().UnixMilli())
	}
	return qb
}

// Preload 预加载关联数据
func (qb *QueryBuilder) Preload(associations ...string) *QueryBuilder {
	for _, assoc := range associations {
		qb.query = qb.query.Preload(assoc)
	}
	return qb
}

// OrderBy 设置排序
func (qb *QueryBuilder) OrderBy(defaultOrder string) *QueryBuilder {
	sortBy := qb.ctx.Query("sort_by")
	if sortBy != "" {
		qb.query = qb.query.Order(sortBy)
	} else if defaultOrder != "" {
		qb.query = qb.query.Order(defaultOrder)
	}
	return qb
}

// Execute 执行查询并返回分页结果
func (qb *QueryBuilder) Execute(result interface{}) (*types.Pagination, error) {
	// 计算总数
	var total int64
	if err := qb.query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 执行分页查询
	if err := qb.query.Offset(qb.offset).Limit(qb.limit).Find(result).Error; err != nil {
		return nil, err
	}

	// 构建分页信息
	pagination := &types.Pagination{
		Page:       qb.page,
		PageSize:   qb.limit,
		Total:      int(total),
		TotalPages: int(math.Ceil(float64(total) / float64(qb.limit))),
	}

	return pagination, nil
}

// 预定义的过滤器组合

// ApplySecretItemFilters 应用密钥项过滤器
func (qb *QueryBuilder) ApplySecretItemFilters() *QueryBuilder {
	return qb.
		StringFilter("category", "category").
		StringFilter("type", "type").
		StringFilter("environment", "environment").
		LikeFilter("tags", "tag").
		MultiLikeFilter([]string{"name", "description"}, "search").
		DateRangeFilter("created_at", "created_at_from", "created_at_to").
		SecretStatusFilter().
		InFilterByName("created_by", "creator_name", "name", "users")
}

// ApplyUserFilters 应用用户过滤器
func (qb *QueryBuilder) ApplyUserFilters() *QueryBuilder {
	return qb.
		LikeFilter("name", "name").
		StringFilter("role", "role").
		StringFilter("status", "status")
}

// ApplyAccessRequestFilters 应用访问申请过滤器
func (qb *QueryBuilder) ApplyAccessRequestFilters() *QueryBuilder {
	return qb.
		StringFilter("status", "status").
		DateRangeFilter("created_at", "created_at_from", "created_at_to").
		InFilterByName("applicant_id", "applicant_name", "name", "users").
		InFilterByName("secret_item_id", "secret_item_name", "name", "secret_items")
}

// ApplyAuditLogFilters 应用审计日志过滤器
func (qb *QueryBuilder) ApplyAuditLogFilters() *QueryBuilder {
	return qb.
		StringFilter("action", "action").
		StringFilter("resource", "resource")
}

// Helper 函数用于快速构建查询

// QuerySecretItems 查询密钥项
func QuerySecretItems(db *gorm.DB, ctx *gin.Context) (*types.ListResponse[models.SecretItem], error) {
	var items []models.SecretItem

	pagination, err := NewQueryBuilder(db, ctx, &models.SecretItem{}).
		ApplySecretItemFilters().
		Preload("Creator", "Updater").
		OrderBy("created_at DESC").
		Execute(&items)

	if err != nil {
		return nil, err
	}

	return &types.ListResponse[models.SecretItem]{
		Data:       items,
		Pagination: *pagination,
	}, nil
}

// QueryUsers 查询用户
func QueryUsers(db *gorm.DB, ctx *gin.Context) (*types.ListResponse[models.User], error) {
	var users []models.User

	pagination, err := NewQueryBuilder(db, ctx, &models.User{}).
		ApplyUserFilters().
		Preload("Creator", "Updater").
		OrderBy("created_at DESC").
		Execute(&users)

	if err != nil {
		return nil, err
	}

	return &types.ListResponse[models.User]{
		Data:       users,
		Pagination: *pagination,
	}, nil
}

// QueryAccessRequests 查询访问申请
func QueryAccessRequests(db *gorm.DB, ctx *gin.Context) (*types.ListResponse[models.AccessRequest], error) {
	var requests []models.AccessRequest

	pagination, err := NewQueryBuilder(db, ctx, &models.AccessRequest{}).
		ApplyAccessRequestFilters().
		Preload("SecretItem", "Applicant", "Approver").
		OrderBy("created_at DESC").
		Execute(&requests)

	if err != nil {
		return nil, err
	}

	return &types.ListResponse[models.AccessRequest]{
		Data:       requests,
		Pagination: *pagination,
	}, nil
}

// QueryAuditLogs 查询审计日志
func QueryAuditLogs(db *gorm.DB, ctx *gin.Context) (*types.ListResponse[models.AuditLog], error) {
	var logs []models.AuditLog

	pagination, err := NewQueryBuilder(db, ctx, &models.AuditLog{}).
		ApplyAuditLogFilters().
		Preload("User").
		OrderBy("created_at DESC").
		Execute(&logs)

	if err != nil {
		return nil, err
	}

	return &types.ListResponse[models.AuditLog]{
		Data:       logs,
		Pagination: *pagination,
	}, nil
}
