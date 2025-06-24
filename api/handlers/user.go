package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/packages/context"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetProfile 获取当前用户信息
func GetProfile(c *gin.Context) {
	user := context.GetCurrentUser(c)
	c.JSON(http.StatusOK, user)
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	// 获取当前用户信息，检查权限
	user := context.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 查询参数
	name := c.Query("name")
	role := c.Query("role")
	status := c.Query("status")

	// 构建查询
	var users []models.User
	var total int64
	query := models.DB.Model(&models.User{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	err := query.
		Preload("Creator").
		Preload("Updater").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		return
	}

	c.JSON(http.StatusOK, types.ListResponse[models.User]{
		Data: users,
		Pagination: types.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      int(total),
			TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
		},
	})
}

// GetUser 获取指定用户信息
func GetUser(c *gin.Context) {
	// 获取当前用户信息，检查权限
	user := context.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 获取用户ID
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "用户ID不能为空"})
		return
	}

	// 查询用户
	var targetUser models.User
	if err := models.DB.Where("id = ?", userID).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		}
		return
	}

	c.JSON(http.StatusOK, targetUser)
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	// 获取当前用户信息，检查权限
	user := context.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 解析请求
	var req types.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	// 验证角色
	if err := models.ValidateUserRole(req.Role); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// 检查用户名是否已存在
	var count int64
	models.DB.Model(&models.User{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	models.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "邮箱已存在"})
		return
	}

	// 角色权限检查：只有超级管理员可以创建超级管理员和安全管理员
	if (req.Role == models.RoleSuperAdmin || req.Role == models.RoleSecMgr) && user.Role != models.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "只有超级管理员可以创建超级管理员和安全管理员"})
		return
	}

	// 创建用户
	newUser := models.User{
		Name:        req.Name,
		Password:    req.Password,
		Email:       req.Email,
		Role:        req.Role,
		Status:      models.StatusActive,
		Permissions: req.Permissions,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}

	if err := models.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "创建用户失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// UpdateProfile 更新当前用户信息
func UpdateProfile(c *gin.Context) {
	user := context.GetCurrentUser(c)

	var req types.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	user.Email = req.Email

	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "更新失败"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新指定用户信息
func UpdateUser(c *gin.Context) {
	// 获取当前用户信息，检查权限
	user := context.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 获取目标用户ID
	targetID := c.Param("id")
	if targetID == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "用户ID不能为空"})
		return
	}

	// 解析请求
	var req types.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误: " + err.Error()})
		return
	}

	// 查询目标用户
	var targetUser models.User
	if err := models.DB.Where("id = ?", targetID).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		}
		return
	}

	// 角色权限检查
	// 1. 只有超级管理员可以修改超级管理员和安全管理员
	if (targetUser.Role == models.RoleSuperAdmin || targetUser.Role == models.RoleSecMgr) && user.Role != models.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "只有超级管理员可以修改超级管理员和安全管理员"})
		return
	}

	// 2. 如果要修改角色，需要检查权限
	if req.Role != "" && req.Role != targetUser.Role {
		// 验证角色有效性
		if err := models.ValidateUserRole(req.Role); err != nil {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
			return
		}

		// 只有超级管理员可以将用户提升为超级管理员或安全管理员
		if (req.Role == models.RoleSuperAdmin || req.Role == models.RoleSecMgr) && user.Role != models.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "只有超级管理员可以提升用户为超级管理员或安全管理员"})
			return
		}

		targetUser.Role = req.Role
	}

	// 更新用户信息
	if req.Name != "" {
		targetUser.Name = req.Name
	}

	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var count int64
		models.DB.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, targetID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "邮箱已被其他用户使用"})
			return
		}

		targetUser.Email = req.Email
	}

	if req.Status != "" {
		// 验证状态有效性
		switch req.Status {
		case models.StatusActive, models.StatusDisabled, models.StatusLocked, models.StatusExpired:
			targetUser.Status = req.Status
		default:
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "无效的用户状态"})
			return
		}
	}

	// 更新权限列表
	if req.Permissions != nil {
		targetUser.Permissions = req.Permissions
	}

	// 更新密码（如果提供）
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "密码加密失败"})
			return
		}
		targetUser.Password = string(hashedPassword)
	}

	// 更新修改者信息
	targetUser.UpdatedBy = user.ID

	// 保存更新
	if err := models.DB.Save(&targetUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, targetUser)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	// 获取当前用户信息，检查权限
	user := context.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 获取目标用户ID
	targetID := c.Param("id")
	if targetID == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "用户ID不能为空"})
		return
	}

	// 不能删除自己
	if targetID == user.ID {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "不能删除自己"})
		return
	}

	// 查询目标用户
	var targetUser models.User
	if err := models.DB.Where("id = ?", targetID).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "查询失败"})
		}
		return
	}

	// 角色权限检查：只有超级管理员可以删除超级管理员和安全管理员
	if (targetUser.Role == models.RoleSuperAdmin || targetUser.Role == models.RoleSecMgr) && user.Role != models.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "只有超级管理员可以删除超级管理员和安全管理员"})
		return
	}

	// 删除用户
	if err := models.DB.Delete(&targetUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}
