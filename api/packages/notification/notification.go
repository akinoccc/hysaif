package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/akinoccc/hysaif/api/models"
)

// GetNotificationTemplates 获取通知模板
func GetNotificationTemplates() map[string]models.NotificationTemplate {
	return map[string]models.NotificationTemplate{
		models.NotificationTypeAccessRequestCreated: {
			Type:     models.NotificationTypeAccessRequestCreated,
			Title:    "新的访问申请",
			Content:  "用户 {{.ApplicantName}} 申请访问密钥项 {{.SecretItemName}}，申请理由：{{.Reason}}",
			Priority: models.NotificationPriorityNormal,
		},
		models.NotificationTypeAccessRequestApproved: {
			Type:     models.NotificationTypeAccessRequestApproved,
			Title:    "访问申请已批准",
			Content:  "您对密钥项 {{.SecretItemName}} 的访问申请已被 {{.ApproverName}} 批准，有效期至 {{.ValidUntil}}",
			Priority: models.NotificationPriorityNormal,
		},
		models.NotificationTypeAccessRequestRejected: {
			Type:     models.NotificationTypeAccessRequestRejected,
			Title:    "访问申请已拒绝",
			Content:  "您对密钥项 {{.SecretItemName}} 的访问申请已被拒绝，拒绝原因：{{.RejectReason}}",
			Priority: models.NotificationPriorityNormal,
		},
		models.NotificationTypeAccessRequestExpired: {
			Type:     models.NotificationTypeAccessRequestExpired,
			Title:    "访问权限已过期",
			Content:  "您对密钥项 {{.SecretItemName}} 的访问权限已过期",
			Priority: models.NotificationPriorityLow,
		},
		models.NotificationTypeSecretItemExpiring: {
			Type:     models.NotificationTypeSecretItemExpiring,
			Title:    "密钥项即将过期",
			Content:  "密钥项 {{.SecretItemName}} 将在 {{.ExpiresIn}} 后过期，请及时更新",
			Priority: models.NotificationPriorityHigh,
		},
		models.NotificationTypeSecretItemExpired: {
			Type:     models.NotificationTypeSecretItemExpired,
			Title:    "密钥项已过期",
			Content:  "密钥项 {{.SecretItemName}} 已过期，请立即更新",
			Priority: models.NotificationPriorityUrgent,
		},
		models.NotificationTypeSystemMaintenance: {
			Type:     models.NotificationTypeSystemMaintenance,
			Title:    "系统维护通知",
			Content:  "系统将于 {{.MaintenanceTime}} 进行维护，预计持续 {{.Duration}}",
			Priority: models.NotificationPriorityHigh,
		},
		models.NotificationTypeSecurityAlert: {
			Type:     models.NotificationTypeSecurityAlert,
			Title:    "安全警报",
			Content:  "检测到异常活动：{{.AlertDetails}}",
			Priority: models.NotificationPriorityUrgent,
		},
	}
}

// CreateNotification 创建通知
func CreateNotification(recipientID, notificationType, relatedID, relatedType string, data models.NotificationData) error {
	templates := GetNotificationTemplates()
	template, exists := templates[notificationType]
	if !exists {
		return fmt.Errorf("未知的通知类型: %s", notificationType)
	}

	// 渲染通知内容
	title, err := renderTemplate(template.Title, data)
	if err != nil {
		return fmt.Errorf("渲染标题失败: %v", err)
	}

	content, err := renderTemplate(template.Content, data)
	if err != nil {
		return fmt.Errorf("渲染内容失败: %v", err)
	}

	// 序列化元数据
	metadataBytes, _ := json.Marshal(data)

	// 创建通知记录
	notification := models.Notification{
		RecipientID: recipientID,
		Type:        notificationType,
		Title:       title,
		Content:     content,
		Status:      models.NotificationStatusUnread,
		Priority:    template.Priority,
		RelatedID:   relatedID,
		RelatedType: relatedType,
		Metadata:    string(metadataBytes),
	}

	// 设置过期时间（默认30天）
	notification.ExpiresAt = uint64(time.Now().AddDate(0, 0, 30).UnixMilli())

	return models.DB.Create(&notification).Error
}

// renderTemplate 渲染模板
func renderTemplate(templateStr string, data models.NotificationData) (string, error) {
	tmpl, err := template.New("notification").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// NotifyAccessRequestCreated 通知新的访问申请
func NotifyAccessRequestCreated(accessRequest *models.AccessRequest) error {
	// 获取需要通知的管理员用户
	admins, err := getAccessRequestApprovers()
	if err != nil {
		return fmt.Errorf("获取管理员用户失败: %v", err)
	}

	// 加载关联数据
	models.DB.Preload("Applicant").Preload("SecretItem").First(accessRequest, "id = ?", accessRequest.ID)

	data := models.NotificationData{
		ApplicantName:  accessRequest.Applicant.Name,
		SecretItemName: accessRequest.SecretItem.Name,
		Reason:         accessRequest.Reason,
	}

	// 为每个管理员创建通知
	for _, admin := range admins {
		err := CreateNotification(
			admin.ID,
			models.NotificationTypeAccessRequestCreated,
			accessRequest.ID,
			"access_request",
			data,
		)
		if err != nil {
			return fmt.Errorf("创建通知失败 (用户: %s): %v", admin.ID, err)
		}
	}

	return nil
}

// NotifyAccessRequestApproved 通知申请已批准
func NotifyAccessRequestApproved(accessRequest *models.AccessRequest) error {
	// 加载关联数据
	models.DB.Preload("Applicant").Preload("SecretItem").Preload("Approver").First(accessRequest, "id = ?", accessRequest.ID)

	validUntil := time.UnixMilli(int64(accessRequest.ValidUntil)).Format("2006-01-02 15:04:05")

	data := models.NotificationData{
		ApproverName:   accessRequest.Approver.Name,
		SecretItemName: accessRequest.SecretItem.Name,
		ValidUntil:     validUntil,
	}

	return CreateNotification(
		accessRequest.ApplicantID,
		models.NotificationTypeAccessRequestApproved,
		accessRequest.ID,
		"access_request",
		data,
	)
}

// NotifyAccessRequestRejected 通知申请已拒绝
func NotifyAccessRequestRejected(accessRequest *models.AccessRequest) error {
	// 加载关联数据
	models.DB.Preload("Applicant").Preload("SecretItem").First(accessRequest, "id = ?", accessRequest.ID)

	data := models.NotificationData{
		SecretItemName: accessRequest.SecretItem.Name,
		RejectReason:   accessRequest.RejectReason,
	}

	return CreateNotification(
		accessRequest.ApplicantID,
		models.NotificationTypeAccessRequestRejected,
		accessRequest.ID,
		"access_request",
		data,
	)
}

// NotifyAccessRequestExpired 通知申请已过期
func NotifyAccessRequestExpired(accessRequest *models.AccessRequest) error {
	// 加载关联数据
	models.DB.Preload("Applicant").Preload("SecretItem").First(accessRequest, "id = ?", accessRequest.ID)

	data := models.NotificationData{
		SecretItemName: accessRequest.SecretItem.Name,
	}

	return CreateNotification(
		accessRequest.ApplicantID,
		models.NotificationTypeAccessRequestExpired,
		accessRequest.ID,
		"access_request",
		data,
	)
}

// NotifySecretItemExpiring 通知密钥项即将过期
func NotifySecretItemExpiring(secretItem *models.SecretItem, expiresIn string) error {
	// 获取有权限的用户
	users, err := getAuthorizedUsers(secretItem)
	if err != nil {
		return fmt.Errorf("获取授权用户失败: %v", err)
	}

	data := models.NotificationData{
		SecretItemName: secretItem.Name,
		ExpiresIn:      expiresIn,
	}

	// 为每个用户创建通知
	for _, user := range users {
		err := CreateNotification(
			user.ID,
			models.NotificationTypeSecretItemExpiring,
			secretItem.ID,
			"secret_item",
			data,
		)
		if err != nil {
			return fmt.Errorf("创建通知失败 (用户: %s): %v", user.ID, err)
		}
	}

	return nil
}

// getAccessRequestApprovers 获取有审核权限的用户
func getAccessRequestApprovers() ([]models.User, error) {
	var users []models.User
	err := models.DB.Raw(`
		SELECT * FROM users
		WHERE (
			role IN (
				SELECT DISTINCT v0 FROM casbin_rule 
				WHERE (v1 = 'access_request' AND (v2 = 'approve' OR v2 = 'reject'))
				OR (v0 = 'super_admin')
			) 
			AND status = 'active'
		)
	`).Scan(&users).Error
	return users, err
}

// getAuthorizedUsers 获取有权限访问密钥项的用户
func getAuthorizedUsers(secretItem *models.SecretItem) ([]models.User, error) {
	// 简化实现：返回管理员和创建者
	var users []models.User

	// 获取管理员
	admins, err := getAccessRequestApprovers()
	if err != nil {
		return nil, err
	}
	users = append(users, admins...)

	// 获取创建者
	if secretItem.CreatedByID != "" {
		var creator models.User
		if err := models.DB.Where("id = ? AND status = ?", secretItem.CreatedByID, "active").First(&creator).Error; err == nil {
			users = append(users, creator)
		}
	}

	return users, nil
}

// GetUserNotifications 获取用户通知列表
func GetUserNotifications(userID string, page, pageSize int, status string) ([]models.Notification, int64, error) {
	query := models.DB.Model(&models.Notification{}).Where("recipient_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var notifications []models.Notification
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error

	return notifications, total, err
}

// MarkNotificationAsRead 标记通知为已读
func MarkNotificationAsRead(notificationID, userID string) error {
	return models.DB.Model(&models.Notification{}).
		Where("id = ? AND recipient_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"status":  models.NotificationStatusRead,
			"read_at": uint64(time.Now().UnixMilli()),
		}).Error
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func MarkAllNotificationsAsRead(userID string) error {
	return models.DB.Model(&models.Notification{}).
		Where("recipient_id = ? AND status = ?", userID, models.NotificationStatusUnread).
		Updates(map[string]interface{}{
			"status":  models.NotificationStatusRead,
			"read_at": uint64(time.Now().UnixMilli()),
		}).Error
}

// GetUnreadNotificationCount 获取未读通知数量
func GetUnreadNotificationCount(userID string) (int64, error) {
	var count int64
	err := models.DB.Model(&models.Notification{}).
		Where("recipient_id = ? AND status = ?", userID, models.NotificationStatusUnread).
		Count(&count).Error
	return count, err
}

// DeleteNotification 删除通知
func DeleteNotification(notificationID, userID string) error {
	return models.DB.Where("id = ? AND recipient_id = ?", notificationID, userID).
		Delete(&models.Notification{}).Error
}

// CleanupExpiredNotifications 清理过期通知
func CleanupExpiredNotifications() error {
	now := uint64(time.Now().UnixMilli())
	return models.DB.Where("expires_at > 0 AND expires_at < ?", now).
		Delete(&models.Notification{}).Error
}
