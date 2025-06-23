package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// CasbinManager Casbin权限管理器
type CasbinManager struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
	mu       sync.RWMutex
}

var (
	casbinManager *CasbinManager
	once          sync.Once
)

// GetCasbinManager 获取Casbin管理器单例
func GetCasbinManager(db *gorm.DB) *CasbinManager {
	once.Do(func() {
		casbinManager = &CasbinManager{db: db}
		casbinManager.initCasbin()
	})
	return casbinManager
}

// initCasbin 初始化Casbin
func (cm *CasbinManager) initCasbin() {
	// 获取模型文件路径
	modelPath := filepath.Join(".", "rbac_model.conf")

	// 创建GORM适配器
	adapter, err := gormadapter.NewAdapterByDB(cm.db)
	if err != nil {
		log.Fatalf("Failed to create GORM adapter: %v", err)
	}

	// 创建模型
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}

	// 创建执行器
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}

	// 如果数据库中没有策略，则从CSV文件初始化
	policies, err := enforcer.GetPolicy()
	if err != nil {
		log.Printf("Error getting policy: %v", err)
	} else if len(policies) == 0 {
		cm.initPoliciesFromCSV(enforcer)
	}

	cm.enforcer = enforcer
	log.Println("Casbin initialized successfully with GORM adapter")
}

// initPoliciesFromCSV 从CSV文件初始化策略到数据库
func (cm *CasbinManager) initPoliciesFromCSV(enforcer *casbin.Enforcer) {
	log.Println("Initializing policies from CSV file...")

	// 定义初始策略 - 完整的菜单权限配置
	policies := [][]string{
		// 超级管理员拥有所有权限
		{"super_admin", "*", "*"},

		// 安全管理员权限
		{"sec_mgr", "dashboard", "read"},
		{"sec_mgr", "users", "read"},
		{"sec_mgr", "users", "create"},
		{"sec_mgr", "users", "update"},
		{"sec_mgr", "users", "delete"},
		{"sec_mgr", "permissions", "read"},
		{"sec_mgr", "permissions", "create"},
		{"sec_mgr", "permissions", "update"},
		{"sec_mgr", "permissions", "delete"},
		{"sec_mgr", "audit", "read"},
		{"sec_mgr", "secret", "read"},
		{"sec_mgr", "secret", "create"},
		{"sec_mgr", "secret", "update"},
		{"sec_mgr", "secret", "delete"},
		{"sec_mgr", "access_request", "read"},
		{"sec_mgr", "access_request", "approve"},
		{"sec_mgr", "access_request", "reject"},
		{"sec_mgr", "access_request", "cancel"},
		{"sec_mgr", "notification", "create"},
		{"sec_mgr", "notification", "bulk_send"},
		{"sec_mgr", "notification", "view_templates"},

		// 开发人员权限
		{"dev", "dashboard", "read"},
		{"dev", "secret", "read"},
		{"dev", "secret", "request"},
		{"dev", "access_request", "read"},
		{"dev", "access_request", "cancel"},

		// 审计员权限
		{"auditor", "dashboard", "read"},
		{"auditor", "audit", "read"},
		{"auditor", "notification", "view_templates"},

		// 机器人权限
		{"bot", "secret", "temp"},
	}

	// 添加策略
	for _, policy := range policies {
		_, err := enforcer.AddPolicy(policy)
		if err != nil {
			log.Printf("Failed to add policy %v: %v", policy, err)
		}
	}

	// 定义角色继承关系
	roleInheritance := [][]string{
		{"sec_mgr", "dev"},
		{"auditor", "dev"},
	}

	// 添加角色继承关系
	for _, role := range roleInheritance {
		_, err := enforcer.AddRoleForUser(role[0], role[1])
		if err != nil {
			log.Printf("Failed to add role inheritance %v: %v", role, err)
		}
	}

	// 保存策略到数据库
	err := enforcer.SavePolicy()
	if err != nil {
		log.Printf("Failed to save policies to database: %v", err)
	} else {
		log.Println("Policies initialized successfully from CSV")
	}
}

// CheckPermission 检查用户权限
func (cm *CasbinManager) CheckPermission(userRole, resource, action string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	result, err := cm.enforcer.Enforce(userRole, resource, action)
	log.Println(userRole, resource, action, result)
	if err != nil {
		log.Printf("Error checking permission: %v", err)
		return false
	}

	return result
}

// AddPolicy 添加权限策略
func (cm *CasbinManager) AddPolicy(role, resource, action string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	added, err := cm.enforcer.AddPolicy(role, resource, action)
	if err != nil {
		return fmt.Errorf("failed to add policy: %v", err)
	}

	if !added {
		return fmt.Errorf("policy already exists")
	}

	// 保存策略到文件
	err = cm.enforcer.SavePolicy()
	if err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	return nil
}

// RemovePolicy 移除权限策略
func (cm *CasbinManager) RemovePolicy(role, resource, action string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	removed, err := cm.enforcer.RemovePolicy(role, resource, action)
	if err != nil {
		return fmt.Errorf("failed to remove policy: %v", err)
	}

	if !removed {
		return fmt.Errorf("policy does not exist")
	}

	// 保存策略到文件
	err = cm.enforcer.SavePolicy()
	if err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	return nil
}

// AddRoleForUser 为用户添加角色
func (cm *CasbinManager) AddRoleForUser(user, role string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	added, err := cm.enforcer.AddRoleForUser(user, role)
	if err != nil {
		return fmt.Errorf("failed to add role for user: %v", err)
	}

	if !added {
		return fmt.Errorf("role already exists for user")
	}

	// 保存策略到文件
	err = cm.enforcer.SavePolicy()
	if err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	return nil
}

// DeleteRoleForUser 删除用户角色
func (cm *CasbinManager) DeleteRoleForUser(user, role string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	deleted, err := cm.enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return fmt.Errorf("failed to delete role for user: %v", err)
	}

	if !deleted {
		return fmt.Errorf("role does not exist for user")
	}

	// 保存策略到文件
	err = cm.enforcer.SavePolicy()
	if err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	return nil
}

// GetRolesForUser 获取用户的所有角色
func (cm *CasbinManager) GetRolesForUser(user string) []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	roles, err := cm.enforcer.GetRolesForUser(user)
	if err != nil {
		log.Printf("Error getting roles for user: %v", err)
		return []string{}
	}

	return roles
}

// GetUsersForRole 获取角色下的所有用户
func (cm *CasbinManager) GetUsersForRole(role string) []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	users, err := cm.enforcer.GetUsersForRole(role)
	if err != nil {
		log.Printf("Error getting users for role: %v", err)
		return []string{}
	}

	return users
}

// GetPolicy 获取所有权限策略
func (cm *CasbinManager) GetPolicy() [][]string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	policy, err := cm.enforcer.GetPolicy()
	if err != nil {
		log.Printf("Error getting policy: %v", err)
		return [][]string{}
	}

	return policy
}

// GetPermissionsForRole 获取角色的所有权限（包括继承的权限）
func (cm *CasbinManager) GetPermissionsForRole(role string) [][]string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// 获取角色的直接权限
	directPermissions, err := cm.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		log.Printf("Error getting direct permissions for role: %v", err)
		return [][]string{}
	}

	// 获取角色继承的所有角色
	inheritedRoles, err := cm.enforcer.GetRolesForUser(role)
	if err != nil {
		log.Printf("Error getting inherited roles for role: %v", err)
		return directPermissions
	}

	// 创建一个map来去重权限
	permissionMap := make(map[string][]string)

	// 添加直接权限
	for _, perm := range directPermissions {
		if len(perm) >= 3 {
			key := perm[1] + ":" + perm[2] // resource:action
			permissionMap[key] = perm
		}
	}

	// 递归获取继承角色的权限
	for _, inheritedRole := range inheritedRoles {
		inheritedPermissions := cm.GetPermissionsForRole(inheritedRole)
		for _, perm := range inheritedPermissions {
			if len(perm) >= 3 {
				key := perm[1] + ":" + perm[2] // resource:action
				if _, exists := permissionMap[key]; !exists {
					// 创建新的权限条目，保持原角色信息但标记为继承
					newPerm := make([]string, len(perm))
					copy(newPerm, perm)
					permissionMap[key] = newPerm
				}
			}
		}
	}

	// 将map转换回slice
	allPermissions := make([][]string, 0, len(permissionMap))
	for _, perm := range permissionMap {
		allPermissions = append(allPermissions, perm)
	}

	log.Println(allPermissions)

	return allPermissions
}

// ReloadPolicy 重新加载策略
func (cm *CasbinManager) ReloadPolicy() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	err := cm.enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("failed to reload policy: %v", err)
	}

	return nil
}

// UpdateRolePermissions 批量更新角色权限
func (cm *CasbinManager) UpdateRolePermissions(role string, permissions map[string][]string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 首先删除该角色的所有现有权限
	_, err := cm.enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		return fmt.Errorf("failed to remove existing permissions: %v", err)
	}

	// 添加新的权限
	for resource, actions := range permissions {
		for _, action := range actions {
			_, err := cm.enforcer.AddPolicy(role, resource, action)
			if err != nil {
				return fmt.Errorf("failed to add permission %s:%s:%s: %v", role, resource, action, err)
			}
		}
	}

	// 保存策略到数据库
	err = cm.enforcer.SavePolicy()
	if err != nil {
		return fmt.Errorf("failed to save policy: %v", err)
	}

	return nil
}
