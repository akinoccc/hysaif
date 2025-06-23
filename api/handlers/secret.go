package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/akinoccc/hysaif/api/internal/context"
	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
)

// GetSecretItems 获取信息项列表
func GetSecretItems(c *gin.Context) {
	user := context.GetCurrentUser(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	itemType := c.Query("type")

	offset := (page - 1) * pageSize

	query := models.DB.Model(&models.SecretItem{}).Preload("Creator").Preload("Updater")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if itemType != "" {
		query = query.Where("type = ?", itemType)
	}

	var total int64
	query.Count(&total)

	var items []models.SecretItem
	if err := query.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
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
		Tags:        req.Tags,
		Data:        &req.Data,
		ExpiresAt:   req.ExpiresAt,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}

	if err := models.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "创建失败"})
		return
	}

	// 重新查询以获取关联数据
	models.DB.Preload("Creator").Preload("Updater").First(&item, "id = ?", item.ID)

	c.JSON(http.StatusCreated, item)
}

// GetSecretItem 获取单个信息项
func GetSecretItem(c *gin.Context) {
	id := c.Param("id")

	var item models.SecretItem
	if err := models.DB.
		Preload("Creator").
		Preload("Updater").
		Where("id = ?", id).
		First(&item).
		Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "信息项不存在"})
		return
	}

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
	item.Tags = req.Tags
	item.Data = &req.Data // 这里会触发自定义序列化器
	item.ExpiresAt = req.ExpiresAt
	item.UpdatedBy = user.ID

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

	c.JSON(http.StatusOK, types.MessageResponse{Message: "删除成功"})
}
