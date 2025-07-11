package types

// 访问申请相关类型
type CreateAccessRequestRequest struct {
	SecretItemID string `json:"secret_item_id" binding:"required"`
	Reason       string `json:"reason" binding:"required,min=5,max=500"`
}

type ApproveAccessRequestRequest struct {
	ValidDuration int    `json:"valid_duration" binding:"required,min=1,max=8760"` // 有效时长（小时），最大365天
	Note          string `json:"note" binding:"max=500"`                           // 审批备注
}

type RejectAccessRequestRequest struct {
	Reason string `json:"reason" binding:"required,min=5,max=500"` // 拒绝理由
}

type RevokeAccessRequestRequest struct {
	Reason string `json:"reason" binding:"required,min=5,max=500"` // 作废理由
}

type AccessRequestListParams struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`
	PageSize     int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status       string `form:"status" binding:"omitempty,oneof=pending approved rejected expired revoked"` // pending, approved, rejected, expired, revoked
	ApplicantID  string `form:"applicant_id"`                                                               // 申请人ID
	SecretItemID string `form:"secret_item_id"`                                                             // 密钥项ID
	SortBy       string `form:"sort_by"`
	SortDesc     bool   `form:"sort_desc"`
}
