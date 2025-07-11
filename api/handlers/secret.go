package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/akinoccc/hysaif/api/middleware"
	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/packages/context"
	"github.com/akinoccc/hysaif/api/packages/query"
	"github.com/akinoccc/hysaif/api/packages/validation"
	"github.com/akinoccc/hysaif/api/types"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// GetSecretItems 获取信息项列表
func GetSecretItems(c *gin.Context) {
	user := context.GetCurrentUser(c)

	// 使用查询构建器简化查询逻辑
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

	// 批量查询用户对这些密钥项的访问权限
	if len(items) > 0 {
		itemIDs := make([]string, len(items))
		for i, item := range items {
			itemIDs[i] = item.ID
		}

		var accessRequests []models.AccessRequest
		now := uint64(time.Now().UnixMilli())
		err := models.DB.Where(`secret_item_id IN (?) AND applicant_id = ? AND status = ? AND ? BETWEEN valid_from AND valid_until`,
			itemIDs, user.ID, models.RequestStatusApproved, now).
			Find(&accessRequests).Error

		if err == nil {
			// 创建映射表，快速查找用户对每个密钥项的访问权限
			accessMap := make(map[string]bool)
			for _, req := range accessRequests {
				accessMap[req.SecretItemID] = true
			}

			// 设置每个密钥项的访问权限标志
			for i := range items {
				items[i].HasApprovedAccess = accessMap[items[i].ID]
				// 加载历史信息
				items[i].LoadHistoryInfo()
			}
		} else {
			// 查询失败时，默认设置为无访问权限
			for i := range items {
				items[i].HasApprovedAccess = false
				// 加载历史信息
				items[i].LoadHistoryInfo()
			}
		}
	}

	// 获取itemType用于审计日志
	itemType := c.Query("type")
	middleware.AuditLog(types.AuditLogActionRead, middleware.GetSecretResourceType(itemType))(c)

	c.JSON(http.StatusOK, types.ListResponse[models.SecretItem]{
		Data:       items,
		Pagination: *pagination,
	})
}

// CreateSecretItem 创建信息项
func CreateSecretItem(c *gin.Context) {
	user := context.GetCurrentUser(c)

	var req types.PostItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.HandleValidationErrors(c, err)
		return
	}

	item := models.SecretItem{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Environment: req.Environment,
		Tags:        req.Tags,
		Data:        &req.Data,
		ExpiresAt:   req.ExpiresAt,
		CreatedByID: user.ID,
		UpdatedByID: user.ID,
	}

	if err := models.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "创建失败"})
		return
	}

	// 创建初始历史版本
	item.CreateHistory(models.HistoryChangeTypeCreated, "创建密钥项", user.ID)

	// 重新查询以获取关联数据
	models.DB.Preload("Creator").Preload("Updater").First(&item, "id = ?", item.ID)

	middleware.AuditLog(types.AuditLogActionCreate, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusCreated, item)
}

// GetSecretItem 获取单个信息项
func GetSecretItem(c *gin.Context) {
	id := c.Param("id")

	user := context.GetCurrentUser(c)
	var approvedRequests []models.AccessRequest
	now := uint64(time.Now().UnixMilli())
	err := models.DB.Where(`secret_item_id IN (?) AND applicant_id = ? AND status = ? AND ? BETWEEN valid_from AND valid_until`,
		id, user.ID, models.RequestStatusApproved, now).
		Find(&approvedRequests).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		return
	}

	var item models.SecretItem
	query := models.DB.
		Preload("Creator").
		Preload("Updater").
		Where("id = ?", id)

	if len(approvedRequests) == 0 {
		query = query.Where("created_by_id = ?", user.ID)
	}

	if err := query.
		First(&item).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "你无法访问此信息项"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

	// 加载历史信息
	item.LoadHistoryInfo()

	middleware.AuditLog(types.AuditLogActionRead, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, item)
}

// UpdateSecretItem 更新信息项
func UpdateSecretItem(c *gin.Context) {
	user := context.GetCurrentUser(c)
	id := c.Param("id")

	var req types.PostItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.HandleValidationErrors(c, err)
		return
	}

	// 先查询现有记录
	var item models.SecretItem
	if err := models.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

	// 更新字段
	item.Name = req.Name
	item.Description = req.Description
	item.Type = req.Type
	item.Category = req.Category
	item.Environment = req.Environment
	item.Tags = req.Tags
	item.Data = &req.Data // 这里会触发自定义序列化器
	item.ExpiresAt = req.ExpiresAt
	item.UpdatedByID = user.ID
	item.Version++ // 增加版本号

	// 保存更新
	if err := models.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "更新失败"})
		return
	}

	// 创建历史版本
	item.CreateHistory(models.HistoryChangeTypeUpdated, "更新密钥项", user.ID)

	// 重新查询以获取关联数据
	models.DB.
		Preload("Creator").
		Preload("Updater").
		First(&item, "id = ?", id)

	middleware.AuditLog(types.AuditLogActionUpdate, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, item)
}

// DeleteSecretItem 删除信息项
func DeleteSecretItem(c *gin.Context) {
	id := c.Param("id")
	user := context.GetCurrentUser(c)

	var item models.SecretItem
	if err := models.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

	// 创建删除历史版本
	item.CreateHistory(models.HistoryChangeTypeDeleted, "删除密钥项", user.ID)

	if err := models.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "删除失败"})
		return
	}

	middleware.AuditLog(types.AuditLogActionDelete, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, types.MessageResponse{Message: "删除成功"})
}

// GetAccessedSecretItems 获取用户有访问权限的信息项
func GetAccessedSecretItems(c *gin.Context) {
	user := context.GetCurrentUser(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var approvedSecretIDs []string
	now := uint64(time.Now().UnixMilli())
	err := models.DB.Model(&models.AccessRequest{}).Where(`applicant_id = ? AND status = ? AND ? BETWEEN valid_from AND valid_until`,
		user.ID, models.RequestStatusApproved, now).Pluck("secret_item_id", &approvedSecretIDs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		return
	}

	query := models.DB.Model(&models.SecretItem{}).Where("id IN (?)", approvedSecretIDs)

	var total int64
	query.Count(&total)

	var items []models.SecretItem
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		return
	}

	c.JSON(http.StatusOK, types.ListResponse[models.SecretItem]{
		Data: items,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      int(total),
			TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
		},
	})
}

// GetSecretItemHistory 获取信息项历史版本列表
func GetSecretItemHistory(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 检查用户对该密钥项的访问权限
	user := context.GetCurrentUser(c)
	var item models.SecretItem
	query := models.DB.Where("id = ?", id)

	// 如果用户不是创建者，需要检查是否有访问权限
	var approvedRequests []models.AccessRequest
	now := uint64(time.Now().UnixMilli())
	err := models.DB.Where(`secret_item_id = ? AND applicant_id = ? AND status = ? AND ? BETWEEN valid_from AND valid_until`,
		id, user.ID, models.RequestStatusApproved, now).
		Find(&approvedRequests).Error

	if err == nil && len(approvedRequests) == 0 {
		query = query.Where("created_by_id = ?", user.ID)
	}

	if err := query.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "你无法访问此信息项"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

	// 获取历史版本列表
	histories, err := models.GetSecretItemHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "获取历史版本失败"})
		return
	}

	// 分页处理
	total := len(histories)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		histories = []models.SecretItemHistory{}
	} else {
		if end > total {
			end = total
		}
		histories = histories[start:end]
	}

	middleware.AuditLog(types.AuditLogActionRead, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, types.ListResponse[models.SecretItemHistory]{
		Data: histories,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
		},
	})
}

// GetSecretItemHistoryByVersion 获取指定版本的信息项历史记录
func GetSecretItemHistoryByVersion(c *gin.Context) {
	id := c.Param("id")
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "版本号格式错误"})
		return
	}

	// 检查用户对该密钥项的访问权限
	user := context.GetCurrentUser(c)
	var item models.SecretItem
	query := models.DB.Where("id = ?", id)

	// 如果用户不是创建者，需要检查是否有访问权限
	var approvedRequests []models.AccessRequest
	now := uint64(time.Now().UnixMilli())
	err = models.DB.Where(`secret_item_id = ? AND applicant_id = ? AND status = ? AND ? BETWEEN valid_from AND valid_until`,
		id, user.ID, models.RequestStatusApproved, now).
		Find(&approvedRequests).Error

	if err == nil && len(approvedRequests) == 0 {
		query = query.Where("created_by_id = ?", user.ID)
	}

	if err := query.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "你无法访问此信息项"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

	// 获取指定版本的历史记录
	history, err := models.GetSecretItemHistoryByVersion(id, version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "历史版本不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "获取历史版本失败"})
		return
	}

	middleware.AuditLog(types.AuditLogActionRead, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, history)
}

// RestoreSecretItemFromHistory 从历史版本恢复信息项
func RestoreSecretItemFromHistory(c *gin.Context) {
	id := c.Param("id")
	user := context.GetCurrentUser(c)

	var req types.RestoreSecretItemFromHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.HandleValidationErrors(c, err)
		return
	}

	// 检查用户对该密钥项的访问权限（只有创建者才能恢复）
	var item models.SecretItem
	if err := models.DB.Where("id = ? AND created_by_id = ?", id, user.ID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "你无法恢复此信息项"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

	// 恢复历史版本
	restoredItem, err := models.RestoreSecretItemFromHistory(id, req.Version, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: fmt.Sprintf("恢复失败: %v", err)})
		return
	}

	middleware.AuditLog(types.AuditLogActionUpdate, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, restoredItem)
}

// CompareSecretItemVersions 比较两个版本的差异
func CompareSecretItemVersions(c *gin.Context) {
	id := c.Param("id")
	user := context.GetCurrentUser(c)

	var req types.CompareVersionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validation.HandleValidationErrors(c, err)
		return
	}

	// 检查用户对该密钥项的访问权限
	var item models.SecretItem
	query := models.DB.Where("id = ?", id)

	// 如果用户不是创建者，需要检查是否有访问权限
	var approvedRequests []models.AccessRequest
	now := uint64(time.Now().UnixMilli())
	err := models.DB.Where(`secret_item_id = ? AND applicant_id = ? AND status = ? AND ? BETWEEN valid_from AND valid_until`,
		id, user.ID, models.RequestStatusApproved, now).
		Find(&approvedRequests).Error

	if err == nil && len(approvedRequests) == 0 {
		query = query.Where("created_by_id = ?", user.ID)
	}

	if err := query.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "你无法访问此信息项"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

	// 比较版本差异
	changes, err := models.CompareSecretItemVersions(id, req.Version1, req.Version2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: fmt.Sprintf("版本比较失败: %v", err)})
		return
	}

	middleware.AuditLog(types.AuditLogActionRead, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, types.VersionComparisonResponse{
		Version1: req.Version1,
		Version2: req.Version2,
		Changes:  changes,
	})
}
