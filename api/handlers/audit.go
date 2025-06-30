package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
)

// GetAuditLogs 获取审计日志
func GetAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	action := c.Query("action")
	resource := c.Query("resource")

	offset := (page - 1) * pageSize

	query := models.DB.Model(&models.AuditLog{}).Preload("User")

	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	var total int64
	query.Count(&total)

	var logs []models.AuditLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		return
	}

	c.JSON(http.StatusOK, types.ListResponse[models.AuditLog]{
		Data: logs,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      int(total),
			TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
		},
	})
}
