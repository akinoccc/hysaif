package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/packages/context"
	"github.com/akinoccc/hysaif/api/types"
	"github.com/gin-gonic/gin"
)

// AuditLog 审计日志中间件
func AuditLog(action, resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 记录审计日志
		user := context.GetCurrentUser(c)
		if user == nil {
			fmt.Printf("用户未登录\n")
			return
		}

		resourceID := c.Param("id")

		// 创建审计日志记录
		auditLog := models.AuditLog{
			UserID:     user.ID,
			Action:     action,
			Resource:   resource,
			ResourceID: resourceID,
			Details:    getRequestDetails(c),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		}

		// 保存到数据库
		if err := models.DB.Create(&auditLog).Error; err != nil {
			// 记录错误，但不影响主要流程
			fmt.Printf("保存审计日志失败: %v\n", err)
		}
	}
}

// getRequestDetails 获取请求详情，用于审计日志
func getRequestDetails(c *gin.Context) string {
	details := map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"status": c.Writer.Status(),
		"query":  c.Request.URL.RawQuery,
	}

	// 对于非GET请求，记录请求体（敏感信息需要过滤）
	if c.Request.Method != "GET" && c.Request.Method != "DELETE" {
		if c.GetHeader("Content-Type") == "application/json" {
			details["content_type"] = "application/json"
		}
	}

	// 转换为JSON字符串
	detailsJSON := ""
	if jsonBytes, err := json.Marshal(details); err == nil {
		detailsJSON = string(jsonBytes)
	}

	return detailsJSON
}

// AutoAuditLog 智能审计日志中间件，自动识别操作类型
func AutoAuditLog(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 根据HTTP方法自动识别操作类型
		action := getActionFromMethod(c.Request.Method)

		// 记录审计日志
		user := context.GetCurrentUser(c)
		if user == nil {
			return
		}

		resourceID := c.Param("id")

		// 创建审计日志记录
		auditLog := models.AuditLog{
			UserID:     user.ID,
			Action:     action,
			Resource:   resource,
			ResourceID: resourceID,
			Details:    getRequestDetails(c),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		}

		// 保存到数据库
		if err := models.DB.Create(&auditLog).Error; err != nil {
			// 记录错误，但不影响主要流程
			fmt.Printf("保存审计日志失败: %v\n", err)
		}
	}
}

// getActionFromMethod 根据HTTP方法获取对应的操作类型
func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return types.AuditLogActionRead
	case "POST":
		return types.AuditLogActionCreate
	case "PUT", "PATCH":
		return types.AuditLogActionUpdate
	case "DELETE":
		return types.AuditLogActionDelete
	default:
		return "unknown"
	}
}

// GetSecretResourceType 根据请求路径获取资源类型
func GetSecretResourceType(resourceType string) string {
	switch resourceType {
	case "api_key":
		return types.AuditLogResourceApiKey
	case "access_key":
		return types.AuditLogResourceAccessKey
	case "ssh_key":
		return types.AuditLogResourceSshKey
	case "password":
		return types.AuditLogResourcePassword
	case "token":
		return types.AuditLogResourceToken
	default:
		return types.AuditLogResourceCustom
	}
}

// LogUserAction 记录特定用户操作的辅助函数
func LogUserAction(userID, action, resource, resourceID, details, ipAddress, userAgent string) {
	auditLog := models.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	if err := models.DB.Create(&auditLog).Error; err != nil {
		fmt.Printf("保存审计日志失败: %v\n", err)
	}
}

// LogAuthAction 记录认证相关操作的辅助函数
func LogAuthAction(c *gin.Context, user models.User, action string) {
	details := map[string]interface{}{
		"user_email": user.Email,
		"user_name":  user.Name,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	detailsJSON := ""
	if jsonBytes, err := json.Marshal(details); err == nil {
		detailsJSON = string(jsonBytes)
	}

	auditLog := models.AuditLog{
		UserID:     user.ID,
		Action:     action,
		Resource:   types.AuditLogResourceUser,
		ResourceID: user.ID,
		Details:    detailsJSON,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	if err := models.DB.Create(&auditLog).Error; err != nil {
		fmt.Printf("保存认证审计日志失败: %v\n", err)
	}
}
