package types

// 用户相关类型
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

type CreateUserRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=50"`
	Password    string   `json:"password" binding:"required,min=8,max=128"`
	Email       string   `json:"email" binding:"required,email"`
	Role        string   `json:"role" binding:"required,oneof=super_admin sec_mgr dev auditor"`
	Permissions []string `json:"permissions"`
}

type UpdateUserRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=2,max=50"`
	Email       string   `json:"email" binding:"omitempty,email"`
	Role        string   `json:"role" binding:"omitempty,oneof=super_admin sec_mgr dev auditor"`
	Status      string   `json:"status" binding:"omitempty,oneof=active disabled locked expired"`
	Permissions []string `json:"permissions"`
	Password    string   `json:"password" binding:"omitempty,min=8,max=128"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

type UserListParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"`
	Role     string `form:"role" binding:"omitempty,oneof=super_admin sec_mgr admin user"`
	Status   string `form:"status" binding:"omitempty,oneof=active disabled locked expired"`
	SortBy   string `form:"sort_by"`
	SortDesc bool   `form:"sort_desc"`
}
