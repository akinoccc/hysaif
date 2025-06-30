package handlers

import (
	"net/http"

	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/packages/context"
	"github.com/akinoccc/hysaif/api/packages/permission"

	"github.com/gin-gonic/gin"
)

// PermissionRequest 权限请求结构
type PermissionRequest struct {
	Role     string `json:"role" binding:"required"`
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

// RoleRequest 角色请求结构
type RoleRequest struct {
	User string `json:"user" binding:"required"`
	Role string `json:"role" binding:"required"`
}

// CheckPermission 检查权限
func CheckPermission(c *gin.Context) {
	var req PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户
	user := context.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 检查指定权限
	casbinManager := permission.GetCasbinManager(models.DB)
	hasPermission := casbinManager.CheckPermission(user.Role, req.Resource, req.Action)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"has_permission": hasPermission,
			"role":           req.Role,
			"resource":       req.Resource,
			"action":         req.Action,
		},
	})
}

// AddPolicy 添加权限策略
func AddPolicy(c *gin.Context) {
	var req PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 添加权限策略
	casbinManager := permission.GetCasbinManager(models.DB)
	err := casbinManager.AddPolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"role":     req.Role,
			"resource": req.Resource,
			"action":   req.Action,
		},
		"message": "权限策略添加成功",
	})
}

// RemovePolicy 移除权限策略
func RemovePolicy(c *gin.Context) {
	var req PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 移除权限策略
	casbinManager := permission.GetCasbinManager(models.DB)
	err := casbinManager.RemovePolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"role":     req.Role,
			"resource": req.Resource,
			"action":   req.Action,
		},
		"message": "权限策略移除成功",
	})
}

// GetPolicies 获取所有权限策略
func GetPolicies(c *gin.Context) {
	// 获取所有权限策略
	casbinManager := permission.GetCasbinManager(models.DB)
	policies := casbinManager.GetPolicy()

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"policies": policies,
		},
	})
}

// AddRoleForUser 为用户添加角色
func AddRoleForUser(c *gin.Context) {
	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证角色是否有效
	if err := models.ValidateUserRole(req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 为用户添加角色
	casbinManager := permission.GetCasbinManager(models.DB)
	err := casbinManager.AddRoleForUser(req.User, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": req.User,
			"role": req.Role,
		},
		"message": "用户角色添加成功",
	})
}

// DeleteRoleForUser 删除用户角色
func DeleteRoleForUser(c *gin.Context) {
	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 删除用户角色
	casbinManager := permission.GetCasbinManager(models.DB)
	err := casbinManager.DeleteRoleForUser(req.User, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": req.User,
			"role": req.Role,
		},
		"message": "用户角色删除成功",
	})
}

// GetRolesForUser 获取用户的所有角色
func GetRolesForUser(c *gin.Context) {
	userParam := c.Param("user")
	if userParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户参数不能为空"})
		return
	}

	// 获取用户的所有角色
	casbinManager := permission.GetCasbinManager(models.DB)
	roles := casbinManager.GetRolesForUser(userParam)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user":  userParam,
			"roles": roles,
		},
	})
}

// GetUsersForRole 获取角色下的所有用户
func GetUsersForRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}

	// 获取角色下的所有用户
	casbinManager := permission.GetCasbinManager(models.DB)
	users := casbinManager.GetUsersForRole(role)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"role":  role,
			"users": users,
		},
	})
}

// GetPermissionsForRole 获取角色的所有权限
func GetPermissionsForRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}

	// 获取角色的所有权限
	casbinManager := permission.GetCasbinManager(models.DB)
	permissions := casbinManager.GetPermissionsForRole(role)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"role":        role,
			"permissions": permissions,
		},
	})
}

// ReloadPolicy 重新加载权限策略
func ReloadPolicy(c *gin.Context) {
	// 重新加载权限策略
	casbinManager := permission.GetCasbinManager(models.DB)
	err := casbinManager.ReloadPolicy()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "权限策略重新加载成功",
	})
}

// UpdateRolePermissionsRequest 批量更新角色权限请求结构
type UpdateRolePermissionsRequest struct {
	Permissions map[string][]string `json:"permissions" binding:"required"`
}

// UpdateRolePermissions 批量更新角色权限
func UpdateRolePermissions(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色参数不能为空"})
		return
	}

	var req UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证角色是否有效
	if err := models.ValidateUserRole(role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 批量更新角色权限
	casbinManager := permission.GetCasbinManager(models.DB)
	err := casbinManager.UpdateRolePermissions(role, req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"role":        role,
			"permissions": req.Permissions,
		},
		"message": "角色权限更新成功",
	})
}

// MenuItemResponse 菜单项响应结构
type MenuItemResponse struct {
	Path  string `json:"path"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Order int    `json:"order"`
}

// GetUserAccessibleMenus 获取用户可访问的菜单列表
func GetUserAccessibleMenus(c *gin.Context) {
	// 获取当前用户
	user := context.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	casbinManager := permission.GetCasbinManager(models.DB)

	// 定义所有可能的菜单项
	allMenus := []MenuItemResponse{
		{Path: "/dashboard", Title: "仪表板", Icon: "LayoutDashboard", Order: 1},
		{Path: "/users", Title: "用户管理", Icon: "Users", Order: 2},
		{Path: "/permission", Title: "权限管理", Icon: "Shield", Order: 3},
		{Path: "/audit", Title: "审计日志", Icon: "FileText", Order: 4},
		{Path: "/access_requests", Title: "访问申请", Icon: "FileText", Order: 5},
		{Path: "/api_key", Title: "API密钥", Icon: "Key", Order: 6},
		{Path: "/access_key", Title: "访问密钥", Icon: "KeyRound", Order: 7},
		{Path: "/ssh_key", Title: "SSH密钥", Icon: "Terminal", Order: 8},
		{Path: "/password", Title: "密码", Icon: "Lock", Order: 9},
		{Path: "/token", Title: "令牌", Icon: "Coins", Order: 10},
		{Path: "/custom", Title: "自定义", Icon: "Settings", Order: 11},
	}

	// 定义菜单权限映射 - 移除静态角色列表，只保留资源和动作映射
	menuPermissions := map[string]struct {
		resource string
		action   string
	}{
		"/dashboard":       {resource: "dashboard", action: "read"},
		"/users":           {resource: "users", action: "read"},
		"/permission":      {resource: "permissions", action: "read"},
		"/audit":           {resource: "audit", action: "read"},
		"/access_requests": {resource: "access_request", action: "read"},
		"/api_key":         {resource: "secret", action: "read"},
		"/access_key":      {resource: "secret", action: "read"},
		"/ssh_key":         {resource: "secret", action: "read"},
		"/password":        {resource: "secret", action: "read"},
		"/token":           {resource: "secret", action: "read"},
		"/custom":          {resource: "secret", action: "read"},
		"/notifications":   {resource: "notification", action: "read"},
	}

	// 过滤用户可访问的菜单 - 完全基于Casbin动态权限检查
	var accessibleMenus []MenuItemResponse
	for _, menu := range allMenus {
		permission, exists := menuPermissions[menu.Path]
		if !exists {
			continue
		}

		// 直接使用Casbin检查权限，不再进行静态角色检查
		if casbinManager.CheckPermission(user.Role, permission.resource, permission.action) || user.IsAdmin() {
			accessibleMenus = append(accessibleMenus, menu)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"menus": accessibleMenus,
		},
	})
}
