package notification

import (
	"fmt"
	"log"
	"time"

	"github.com/akinoccc/hysaif/api/models"
)

var stopChan = make(chan struct{})

// Start 启动定时任务
func Start() {
	log.Println("启动定时任务服务")

	// 启动各种定时任务
	go startExpiredAccessRequestChecker()
	go startSecretItemExpirationChecker()
	go startNotificationCleanup()
}

// Stop 停止定时任务
func Stop() {
	log.Println("停止定时任务服务")
	close(stopChan)
}

// startExpiredAccessRequestChecker 检查过期的访问申请
func startExpiredAccessRequestChecker() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkExpiredAccessRequests()
		case <-stopChan:
			return
		}
	}
}

// startSecretItemExpirationChecker 检查即将过期的密钥项
func startSecretItemExpirationChecker() {
	ticker := time.NewTicker(6 * time.Hour) // 每6小时检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkSecretItemExpiration()
		case <-stopChan:
			return
		}
	}
}

// startNotificationCleanup 清理过期通知
func startNotificationCleanup() {
	ticker := time.NewTicker(24 * time.Hour) // 每天清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredNotifications()
		case <-stopChan:
			return
		}
	}
}

// checkExpiredAccessRequests 检查并处理过期的访问申请
func checkExpiredAccessRequests() {
	log.Println("检查过期的访问申请")

	now := uint64(time.Now().UnixMilli())
	var expiredRequests []models.AccessRequest

	// 查找已过期但状态仍为approved的申请
	err := models.DB.Where("status = ? AND valid_until < ?", models.RequestStatusApproved, now).
		Find(&expiredRequests).Error
	if err != nil {
		log.Printf("查询过期访问申请失败: %v", err)
		return
	}

	for _, request := range expiredRequests {
		// 更新状态为过期
		request.Status = models.RequestStatusExpired
		if err := models.DB.Save(&request).Error; err != nil {
			log.Printf("更新访问申请状态失败 (ID: %s): %v", request.ID, err)
			continue
		}

		// 通知用户访问权限已过期
		if err := NotifyAccessRequestExpired(&request); err != nil {
			log.Printf("发送过期通知失败: %v", err)
		}

		log.Printf("访问申请已过期 (ID: %s)", request.ID)
	}

	if len(expiredRequests) > 0 {
		log.Printf("处理了 %d 个过期的访问申请", len(expiredRequests))
	}
}

// checkSecretItemExpiration 检查即将过期的密钥项
func checkSecretItemExpiration() {
	log.Println("检查即将过期的密钥项")

	now := time.Now()
	// 检查7天内过期的密钥项
	warningTime := uint64(now.Add(7 * 24 * time.Hour).UnixMilli())
	// 检查已过期的密钥项
	expiredTime := uint64(now.UnixMilli())

	var expiringItems []models.SecretItem
	var expiredItems []models.SecretItem

	// 查找即将过期的密钥项（7天内过期且未发送过通知）
	err := models.DB.Where("expires_at > 0 AND expires_at > ? AND expires_at <= ?", expiredTime, warningTime).
		Find(&expiringItems).Error
	if err != nil {
		log.Printf("查询即将过期的密钥项失败: %v", err)
	} else {
		for _, item := range expiringItems {
			// 检查是否已经发送过即将过期的通知（避免重复发送）
			var existingNotification models.Notification
			err := models.DB.Where("type = ? AND related_id = ? AND created_at > ?",
				models.NotificationTypeSecretItemExpiring, item.ID, uint64(now.Add(-24*time.Hour).UnixMilli())).
				First(&existingNotification).Error
			if err == nil {
				continue // 已发送过通知，跳过
			}

			// 计算剩余时间
			expiresAt := time.UnixMilli(int64(item.ExpiresAt))
			remaining := time.Until(expiresAt)
			expiresIn := fmt.Sprintf("%.0f小时", remaining.Hours())

			// 发送即将过期通知
			if err := NotifySecretItemExpiring(&item, expiresIn); err != nil {
				log.Printf("发送密钥项过期通知失败: %v", err)
			}
		}
	}

	// 查找已过期的密钥项
	err = models.DB.Where("expires_at > 0 AND expires_at <= ?", expiredTime).
		Find(&expiredItems).Error
	if err != nil {
		log.Printf("查询已过期的密钥项失败: %v", err)
	} else {
		for _, item := range expiredItems {
			// 检查是否已经发送过过期通知（避免重复发送）
			var existingNotification models.Notification
			err := models.DB.Where("type = ? AND related_id = ? AND created_at > ?",
				models.NotificationTypeSecretItemExpired, item.ID, uint64(now.Add(-24*time.Hour).UnixMilli())).
				First(&existingNotification).Error
			if err == nil {
				continue // 已发送过通知，跳过
			}

			data := models.NotificationData{
				SecretItemName: item.Name,
			}

			// 获取管理员用户（简化实现）
			var users []models.User
			err = models.DB.Where("role IN (?, ?) AND status = ?",
				"super_admin", "sec_mgr", "active").Find(&users).Error
			if err != nil {
				log.Printf("获取管理员用户失败 (ID: %s): %v", item.ID, err)
				continue
			}

			// 为每个用户创建通知
			for _, user := range users {
				err := CreateNotification(
					user.ID,
					models.NotificationTypeSecretItemExpired,
					item.ID,
					"secret_item",
					data,
				)
				if err != nil {
					log.Printf("创建过期通知失败 (用户: %s): %v", user.ID, err)
				}
			}

			log.Printf("密钥项已过期 (ID: %s)", item.ID)
		}
	}

	if len(expiringItems) > 0 || len(expiredItems) > 0 {
		log.Printf("处理了 %d 个即将过期和 %d 个已过期的密钥项", len(expiringItems), len(expiredItems))
	}
}

// cleanupExpiredNotifications 清理过期通知
func cleanupExpiredNotifications() {
	log.Println("清理过期通知")

	if err := CleanupExpiredNotifications(); err != nil {
		log.Printf("清理过期通知失败: %v", err)
	} else {
		log.Println("过期通知清理完成")
	}
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%d天%d小时", days, hours)
		}
		return fmt.Sprintf("%d天", days)
	}
	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%d小时%d分钟", hours, minutes)
		}
		return fmt.Sprintf("%d小时", hours)
	}
	return fmt.Sprintf("%d分钟", minutes)
}
