package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/akinoccc/hysaif/api/internal/context"
	"github.com/akinoccc/hysaif/api/internal/notification"
	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
)

// GetNotifications 获取用户通知列表
func GetNotifications(c *gin.Context) {
	user := context.GetCurrentUser(c)

	var req types.GetNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 获取通知列表
	notifications, total, err := notification.GetUserNotifications(
		user.ID, req.Page, req.PageSize, req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "获取通知列表失败"})
		return
	}

	// 获取未读通知数量
	unreadCount, err := notification.GetUnreadNotificationCount(user.ID)
	if err != nil {
		unreadCount = 0
	}

	// 计算分页信息
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	pagination := types.Pagination{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}

	response := types.NotificationListResponse{
		Data:        notifications,
		Pagination:  pagination,
		UnreadCount: unreadCount,
	}

	c.JSON(http.StatusOK, response)
}

// MarkNotificationAsRead 标记通知为已读
func MarkNotificationAsRead(c *gin.Context) {
	user := context.GetCurrentUser(c)
	notificationID := c.Param("id")

	if notificationID == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "通知ID不能为空"})
		return
	}

	err := notification.MarkNotificationAsRead(notificationID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "标记通知失败"})
		return
	}

	c.JSON(http.StatusOK, types.MessageResponse{Message: "通知已标记为已读"})
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func MarkAllNotificationsAsRead(c *gin.Context) {
	user := context.GetCurrentUser(c)

	err := notification.MarkAllNotificationsAsRead(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "标记通知失败"})
		return
	}

	c.JSON(http.StatusOK, types.MessageResponse{Message: "所有通知已标记为已读"})
}

// GetUnreadNotificationCount 获取未读通知数量
func GetUnreadNotificationCount(c *gin.Context) {
	user := context.GetCurrentUser(c)

	count, err := notification.GetUnreadNotificationCount(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "获取未读通知数量失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// DeleteNotification 删除通知
func DeleteNotification(c *gin.Context) {
	user := context.GetCurrentUser(c)
	notificationID := c.Param("id")

	if notificationID == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "通知ID不能为空"})
		return
	}

	err := notification.DeleteNotification(notificationID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "删除通知失败"})
		return
	}

	c.JSON(http.StatusOK, types.MessageResponse{Message: "通知已删除"})
}

// GetNotificationStats 获取通知统计信息
func GetNotificationStats(c *gin.Context) {
	user := context.GetCurrentUser(c)

	// 获取总数
	var totalCount int64
	models.DB.Model(&models.Notification{}).Where("recipient_id = ?", user.ID).Count(&totalCount)

	// 获取未读数量
	var unreadCount int64
	models.DB.Model(&models.Notification{}).
		Where("recipient_id = ? AND status = ?", user.ID, models.NotificationStatusUnread).
		Count(&unreadCount)

	// 获取已读数量
	readCount := totalCount - unreadCount

	// 按类型统计
	byType := make(map[string]int64)
	var typeStats []struct {
		Type  string
		Count int64
	}
	models.DB.Model(&models.Notification{}).
		Select("type, count(*) as count").
		Where("recipient_id = ?", user.ID).
		Group("type").
		Scan(&typeStats)

	for _, stat := range typeStats {
		byType[stat.Type] = stat.Count
	}

	// 按优先级统计
	byPriority := make(map[string]int64)
	var priorityStats []struct {
		Priority string
		Count    int64
	}
	models.DB.Model(&models.Notification{}).
		Select("priority, count(*) as count").
		Where("recipient_id = ?", user.ID).
		Group("priority").
		Scan(&priorityStats)

	for _, stat := range priorityStats {
		byPriority[stat.Priority] = stat.Count
	}

	stats := types.NotificationStatsResponse{
		TotalCount:  totalCount,
		UnreadCount: unreadCount,
		ReadCount:   readCount,
		ByType:      byType,
		ByPriority:  byPriority,
	}

	c.JSON(http.StatusOK, stats)
}

// CreateNotification 创建通知（管理员功能）
func CreateNotification(c *gin.Context) {
	user := context.GetCurrentUser(c)

	// 检查权限
	if !user.HasPermission("notification", "create") {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "权限不足"})
		return
	}

	var req types.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 设置默认优先级
	if req.Priority == "" {
		req.Priority = models.NotificationPriorityNormal
	}

	// 为每个接收者创建通知
	var createdCount int
	for _, recipientID := range req.RecipientIDs {
		// 验证接收者是否存在
		var recipient models.User
		if err := models.DB.Where("id = ? AND status = ?", recipientID, models.StatusActive).First(&recipient).Error; err != nil {
			continue // 跳过不存在或非活跃用户
		}

		// 创建通知
		notification := models.Notification{
			RecipientID: recipientID,
			Type:        req.Type,
			Title:       req.Title,
			Content:     req.Content,
			Status:      models.NotificationStatusUnread,
			Priority:    req.Priority,
			ExpiresAt:   req.ExpiresAt,
		}

		// 如果没有设置过期时间，默认30天
		if notification.ExpiresAt == 0 {
			notification.ExpiresAt = uint64(time.Now().AddDate(0, 0, 30).UnixMilli())
		}

		if err := models.DB.Create(&notification).Error; err == nil {
			createdCount++
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "通知创建完成",
		"created_count": createdCount,
		"total_count":   len(req.RecipientIDs),
	})
}

// SendBulkNotification 批量发送通知（管理员功能）
func SendBulkNotification(c *gin.Context) {
	user := context.GetCurrentUser(c)

	// 检查权限
	if !user.HasPermission("notification", "bulk_send") {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "权限不足"})
		return
	}

	var req types.BulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 设置默认优先级
	if req.Priority == "" {
		req.Priority = models.NotificationPriorityNormal
	}

	// 收集接收者ID
	var recipientIDs []string

	// 按角色获取用户
	if len(req.UserRoles) > 0 {
		var users []models.User
		models.DB.Where("role IN (?) AND status = ?", req.UserRoles, models.StatusActive).Find(&users)
		for _, user := range users {
			recipientIDs = append(recipientIDs, user.ID)
		}
	}

	// 添加指定用户ID
	recipientIDs = append(recipientIDs, req.UserIDs...)

	// 去重
	uniqueIDs := make(map[string]bool)
	var finalRecipientIDs []string
	for _, id := range recipientIDs {
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
			finalRecipientIDs = append(finalRecipientIDs, id)
		}
	}

	// 创建通知
	var createdCount int
	for _, recipientID := range finalRecipientIDs {
		notification := models.Notification{
			RecipientID: recipientID,
			Type:        req.Type,
			Title:       req.Title,
			Content:     req.Content,
			Status:      models.NotificationStatusUnread,
			Priority:    req.Priority,
			ExpiresAt:   req.ExpiresAt,
		}

		// 如果没有设置过期时间，默认30天
		if notification.ExpiresAt == 0 {
			notification.ExpiresAt = uint64(time.Now().AddDate(0, 0, 30).UnixMilli())
		}

		if err := models.DB.Create(&notification).Error; err == nil {
			createdCount++
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "批量通知发送完成",
		"created_count": createdCount,
		"total_count":   len(finalRecipientIDs),
	})
}

// GetNotificationTemplates 获取通知模板列表
func GetNotificationTemplates(c *gin.Context) {
	user := context.GetCurrentUser(c)

	// 检查权限
	if !user.HasPermission("notification", "view_templates") {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "权限不足"})
		return
	}

	templates := notification.GetNotificationTemplates()
	var response []types.NotificationTemplateResponse

	for notificationType, template := range templates {
		response = append(response, types.NotificationTemplateResponse{
			Type:        notificationType,
			Title:       template.Title,
			Content:     template.Content,
			Priority:    template.Priority,
			Description: getTemplateDescription(notificationType),
			Variables:   getTemplateVariables(notificationType),
		})
	}

	c.JSON(http.StatusOK, response)
}

// getTemplateDescription 获取模板描述
func getTemplateDescription(notificationType string) string {
	descriptions := map[string]string{
		models.NotificationTypeAccessRequestCreated:  "当用户创建新的访问申请时发送给管理员",
		models.NotificationTypeAccessRequestApproved: "当访问申请被批准时发送给申请人",
		models.NotificationTypeAccessRequestRejected: "当访问申请被拒绝时发送给申请人",
		models.NotificationTypeAccessRequestExpired:  "当访问权限过期时发送给用户",
		models.NotificationTypeSecretItemExpiring:    "当密钥项即将过期时发送给相关用户",
		models.NotificationTypeSecretItemExpired:     "当密钥项已过期时发送给相关用户",
		models.NotificationTypeSystemMaintenance:     "系统维护通知",
		models.NotificationTypeSecurityAlert:         "安全警报通知",
	}
	return descriptions[notificationType]
}

// getTemplateVariables 获取模板变量
func getTemplateVariables(notificationType string) []string {
	variables := map[string][]string{
		models.NotificationTypeAccessRequestCreated:  {"ApplicantName", "SecretItemName", "Reason"},
		models.NotificationTypeAccessRequestApproved: {"ApproverName", "SecretItemName", "ValidUntil"},
		models.NotificationTypeAccessRequestRejected: {"SecretItemName", "RejectReason"},
		models.NotificationTypeAccessRequestExpired:  {"SecretItemName"},
		models.NotificationTypeSecretItemExpiring:    {"SecretItemName", "ExpiresIn"},
		models.NotificationTypeSecretItemExpired:     {"SecretItemName"},
		models.NotificationTypeSystemMaintenance:     {"MaintenanceTime", "Duration"},
		models.NotificationTypeSecurityAlert:         {"AlertDetails"},
	}
	return variables[notificationType]
}
