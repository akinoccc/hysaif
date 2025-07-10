package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/packages/context"
	"github.com/akinoccc/hysaif/api/packages/notification"
	"github.com/akinoccc/hysaif/api/packages/query"
	"github.com/akinoccc/hysaif/api/types"
)

// CreateAccessRequest 创建访问申请
func CreateAccessRequest(c *gin.Context) {
	user := context.GetCurrentUser(c)

	var req types.CreateAccessRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 检查密钥项是否存在
	var secretItem models.SecretItem
	if err := models.DB.Where("id = ?", req.SecretItemID).First(&secretItem).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "密钥项不存在"})
		return
	}

	// 检查是否已有待审批的申请
	var existingRequest models.AccessRequest
	baseQuery := "secret_item_id = ? AND applicant_id = ?"
	if err := models.DB.Where(baseQuery+" AND status = ?",
		req.SecretItemID, user.ID, models.RequestStatusPending).
		Or(baseQuery+" AND status = ? AND valid_until > ?",
			req.SecretItemID, user.ID, models.RequestStatusApproved, uint64(time.Now().UnixMilli())).
		Or(baseQuery+" AND status = ? AND valid_until < ?",
			req.SecretItemID, user.ID, models.RequestStatusApproved, uint64(time.Now().UnixMilli())).
		First(&existingRequest).
		Error; err == nil {
		c.JSON(http.StatusConflict, types.ErrorResponse{
			Error: "已有待审批的申请或已审批的申请，请勿重复提交",
		})
		return
	}

	// 创建申请
	accessRequest := models.AccessRequest{
		SecretItemID: req.SecretItemID,
		ApplicantID:  user.ID,
		Reason:       req.Reason,
		Status:       models.RequestStatusPending,
	}

	if err := models.DB.Create(&accessRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "创建申请失败"})
		return
	}

	// 发送通知给管理员
	if err := notification.NotifyAccessRequestCreated(&accessRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建通知失败: " + err.Error()})
		return
	}

	// 重新查询以获取关联数据
	models.DB.Preload("SecretItem").Preload("Applicant").First(&accessRequest, "id = ?", accessRequest.ID)

	c.JSON(http.StatusCreated, accessRequest)
}

// GetAccessRequests 获取访问申请列表
func GetAccessRequests(c *gin.Context) {
	user := context.GetCurrentUser(c)

	// 使用查询构建器
	qb := query.NewQueryBuilder(models.DB, c, &models.AccessRequest{}).
		ApplyAccessRequestFilters().
		Preload("SecretItem", "Applicant", "Approver").
		OrderBy("created_at DESC")

	// 权限控制：普通用户只能查看自己的申请
	if !user.HasPermission("access_request", "approve") {
		qb = qb.Where("applicant_id = ?", user.ID)
	}

	var requests []models.AccessRequest
	pagination, err := qb.Execute(&requests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		return
	}

	c.JSON(http.StatusOK, types.ListResponse[models.AccessRequest]{
		Data:       requests,
		Pagination: *pagination,
	})
}

// ApproveAccessRequest 批准访问申请
func ApproveAccessRequest(c *gin.Context) {
	user := context.GetCurrentUser(c)
	requestID := c.Param("id")

	var req types.ApproveAccessRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 查询申请
	var accessRequest models.AccessRequest
	if err := models.DB.Preload("SecretItem").Preload("Applicant").Where("id = ?", requestID).First(&accessRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "申请不存在"})
		return
	}

	// 检查申请状态
	if accessRequest.Status != models.RequestStatusPending {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "申请状态不允许审批"})
		return
	}

	// 计算有效期
	now := uint64(time.Now().UnixMilli())
	validUntil := now + uint64(req.ValidDuration*3600*1000) // 转换为毫秒
	log.Println("validUntil", validUntil, uint64(req.ValidDuration*3600*1000))

	// 更新申请状态
	accessRequest.Status = models.RequestStatusApproved
	accessRequest.ApprovedByID = user.ID
	accessRequest.ApprovedAt = now
	accessRequest.ValidFrom = now
	accessRequest.ValidUntil = validUntil
	accessRequest.Note = req.Note

	if err := models.DB.Save(&accessRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "审批失败"})
		return
	}

	// 发送通知
	if err := notification.NotifyAccessRequestApproved(&accessRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送通知失败: " + err.Error()})
		return
	}

	// 重新查询以获取关联数据
	models.DB.Preload("SecretItem").Preload("Applicant").Preload("Approver").First(&accessRequest, "id = ?", accessRequest.ID)

	c.JSON(http.StatusOK, accessRequest)
}

// RevokeAccessRequest 作废访问申请
func RevokeAccessRequest(c *gin.Context) {
	requestID := c.Param("id")

	var req types.RevokeAccessRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 查询申请
	var accessRequest models.AccessRequest
	if err := models.DB.Preload("SecretItem").Preload("Applicant").Where("id = ?", requestID).First(&accessRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "申请不存在"})
		return
	}

	// 检查申请状态 - 只有已批准的申请才能作废
	if accessRequest.Status != models.RequestStatusApproved {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "只有已批准的申请才能作废"})
		return
	}

	// 更新申请状态
	accessRequest.Status = models.RequestStatusRevoked
	accessRequest.RejectReason = req.Reason // 复用拒绝理由字段存储作废理由

	if err := models.DB.Save(&accessRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "作废失败"})
		return
	}

	// 重新查询以获取关联数据
	models.DB.Preload("SecretItem").Preload("Applicant").Preload("Approver").First(&accessRequest, "id = ?", accessRequest.ID)

	c.JSON(http.StatusOK, accessRequest)
}

// RejectAccessRequest 拒绝访问申请
func RejectAccessRequest(c *gin.Context) {
	user := context.GetCurrentUser(c)
	requestID := c.Param("id")

	var req types.RejectAccessRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 查询申请
	var accessRequest models.AccessRequest
	if err := models.DB.Preload("SecretItem").Preload("Applicant").Where("id = ?", requestID).First(&accessRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "申请不存在"})
		return
	}

	// 检查申请状态
	if accessRequest.Status != models.RequestStatusPending {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "申请状态不允许审批"})
		return
	}

	// 更新申请状态
	accessRequest.Status = models.RequestStatusRejected
	accessRequest.ApprovedByID = user.ID
	accessRequest.ApprovedAt = uint64(time.Now().Unix())
	accessRequest.RejectReason = req.Reason

	if err := models.DB.Save(&accessRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "审批失败"})
		return
	}

	// 发送通知
	if err := notification.NotifyAccessRequestRejected(&accessRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送通知失败: " + err.Error()})
		return
	}

	// 重新查询以获取关联数据
	models.DB.Preload("SecretItem").Preload("Applicant").Preload("Approver").First(&accessRequest, "id = ?", accessRequest.ID)

	c.JSON(http.StatusOK, accessRequest)
}

// GetItemWithAccessCheck 通过申请访问密钥项详情
func GetItemWithAccessCheck(c *gin.Context) {
	user := context.GetCurrentUser(c)
	id := c.Param("id")

	if user.HasPermission("secret", "update") {
		GetSecretItem(c)
		return
	}

	// 检查是否有有效的访问申请
	var accessRequest models.AccessRequest
	if err := models.DB.Where("secret_item_id = ? AND applicant_id = ? AND status = ?",
		id, user.ID, models.RequestStatusApproved).First(&accessRequest).Error; err != nil {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "无访问权限，请先申请访问"})
		return
	}

	// 检查申请是否有效
	if !accessRequest.CanAccess() {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "访问权限已过期，请重新申请"})
		return
	}

	// 更新访问记录
	accessRequest.AccessCount++
	accessRequest.LastAccessed = uint64(time.Now().Unix())
	models.DB.Save(&accessRequest)

	// 获取密钥项详情
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
