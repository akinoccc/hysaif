package handlers

import (
	"net/http"

	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/packages/query"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
)

// GetAuditLogs 获取审计日志
func GetAuditLogs(c *gin.Context) {
	// 使用查询构建器
	response, err := query.QueryAuditLogs(models.DB, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}
