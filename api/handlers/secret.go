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
	"github.com/akinoccc/hysaif/api/types"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// GetSecretItems 获取信息项列表
func GetSecretItems(c *gin.Context) {
	user := context.GetCurrentUser(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	itemType := c.Query("type")
	status := c.Query("status")
	environment := c.Query("environment")
	search := c.Query("search")
	tag := c.Query("tag")
	creator := c.Query("creator_name")

	createdAtFrom := c.Query("created_at_from")
	createdAtTo := c.Query("created_at_to")
	sortBy := c.Query("sort_by")

	offset := (page - 1) * pageSize

	query := models.DB.Model(&models.SecretItem{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if itemType != "" {
		query = query.Where("type = ?", itemType)
	}
	if status != "" {
		switch status {
		case "expired":
			query = query.Where("expires_at < ?", time.Now().UnixMilli())
		case "expiring":
			query = query.Where("expires_at > ? AND expires_at < ?", time.Now().UnixMilli(), time.Now().AddDate(0, 0, 7).UnixMilli())
		case "active":
			query = query.Where("expires_at > ?", time.Now().UnixMilli())
		}
	}
	if environment != "" {
		query = query.Where("environment = ?", environment)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if tag != "" {
		query = query.Where("tags LIKE %?%", tag)
	}
	if creator != "" {
		userIDs := []string{}
		models.DB.Model(&models.User{}).Where("name LIKE ?", "%"+creator+"%").Pluck("id", &userIDs)
		query = query.Where("created_by IN (?)", userIDs)
	}
	if createdAtFrom != "" {
		query = query.Where("created_at >= ?", createdAtFrom)
	}
	if createdAtTo != "" {
		query = query.Where("created_at <= ?", createdAtTo)
	}
	if sortBy != "" {
		query = query.Order(sortBy)
	}

	var total int64
	query.Count(&total)

	var items []models.SecretItem
	if err := query.
		Preload("Creator").
		Preload("Updater").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
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
			}
		} else {
			// 查询失败时，默认设置为无访问权限
			for i := range items {
				items[i].HasApprovedAccess = false
			}
		}
	}

	middleware.AuditLog(types.AuditLogActionRead, middleware.GetSecretResourceType(itemType))(c)

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

// CreateSecretItem 创建信息项
func CreateSecretItem(c *gin.Context) {
	user := context.GetCurrentUser(c)

	var req types.PostItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
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

	middleware.AuditLog(types.AuditLogActionRead, middleware.GetSecretResourceType(item.Type))(c)

	c.JSON(http.StatusOK, item)
}

// UpdateSecretItem 更新信息项
func UpdateSecretItem(c *gin.Context) {
	user := context.GetCurrentUser(c)
	id := c.Param("id")

	var req types.PostItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
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

	// 保存更新
	if err := models.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "更新失败"})
		return
	}

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

	var item models.SecretItem
	if err := models.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

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
