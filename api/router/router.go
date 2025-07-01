package router

import (
	"github.com/akinoccc/hysaif/api/handlers"
	"github.com/akinoccc/hysaif/api/middleware"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) {
	// 中间件
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)

			// WebAuthn 路由
			auth.POST("/webauthn/login/begin", handlers.WebAuthnBeginLogin)
			auth.POST("/webauthn/login/finish", handlers.WebAuthnFinishLogin)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// 用户管理
			users := protected.Group("/users")
			users.Use(middleware.AutoAuditLog(types.AuditLogResourceUser))
			{
				// 个人资料相关（所有用户都可以访问自己的资料）
				users.GET("/profile", handlers.GetProfile)
				users.PUT("/profile", handlers.UpdateProfile)
				users.GET("/login-history", handlers.GetLoginHistory)

				// WebAuthn 凭证管理
				users.POST("/webauthn/register/begin", handlers.WebAuthnBeginRegistration)
				users.POST("/webauthn/register/finish", handlers.WebAuthnFinishRegistration)
				users.GET("/webauthn/credentials", handlers.GetUserCredentials)
				users.DELETE("/webauthn/credentials/:id", handlers.DeleteUserCredential)

				// 用户CRUD操作（需要用户管理权限）
				users.GET("/", middleware.RequirePermission("user", "read"), handlers.GetUsers)
				users.POST("/", middleware.RequirePermission("user", "create"), handlers.CreateUser)
				users.GET("/:id", middleware.RequirePermission("user", "read"), handlers.GetUser)
				users.PUT("/:id", middleware.RequirePermission("user", "update"), handlers.UpdateUser)
				users.DELETE("/:id", middleware.RequirePermission("user", "delete"), handlers.DeleteUser)
			}

			// 信息项管理
			items := protected.Group("/items")
			{
				items.GET("/", middleware.RequirePermission("secret", "read"), handlers.GetSecretItems)
				items.POST("/", middleware.RequirePermission("secret", "create"), handlers.CreateSecretItem)
				items.GET("/:id", middleware.RequirePermission("secret", "read"), handlers.GetSecretItem)
				items.PUT("/:id", middleware.RequirePermission("secret", "update"), handlers.UpdateSecretItem)
				items.DELETE("/:id", middleware.RequirePermission("secret", "delete"), handlers.DeleteSecretItem)

				// 通过申请访问密钥项（所有用户都可以使用）
				items.GET("/:id/access",
					middleware.AuditLog(types.AuditLogActionAccess, types.AuditLogResourceCustom),
					handlers.GetItemWithAccessCheck)
			}

			// 访问申请管理
			access := protected.Group("/access-requests")
			{
				// 创建访问申请（所有用户都可以）
				access.POST("/",
					middleware.AuditLog(types.AuditLogActionRequest, types.AuditLogResourceAccessRequest),
					handlers.CreateAccessRequest)
				// 获取申请列表（用户只能看自己的，管理员可以看全部）
				access.GET("/",
					middleware.AutoAuditLog(types.AuditLogResourceAccessRequest),
					handlers.GetAccessRequests)
				// 审批申请（需要管理权限）
				access.PUT("/:id/approve",
					middleware.RequirePermission("access_request", "approve"),
					middleware.AuditLog(types.AuditLogActionApprove, types.AuditLogResourceAccessRequest),
					handlers.ApproveAccessRequest)
				access.PUT("/:id/reject",
					middleware.RequirePermission("access_request", "approve"),
					middleware.AuditLog(types.AuditLogActionReject, types.AuditLogResourceAccessRequest),
					handlers.RejectAccessRequest)
				access.PUT("/:id/revoke",
					middleware.RequirePermission("access_request", "approve"),
					middleware.AuditLog(types.AuditLogActionRevoke, types.AuditLogResourceAccessRequest),
					handlers.RevokeAccessRequest)
			}

			// 通知管理
			notifications := protected.Group("/notifications")
			{
				// 用户通知相关（所有用户都可以访问自己的通知）
				notifications.GET("/", handlers.GetNotifications)
				notifications.PUT("/:id/read", handlers.MarkNotificationAsRead)
				notifications.PUT("/read-all", handlers.MarkAllNotificationsAsRead)
				notifications.GET("/unread-count", handlers.GetUnreadNotificationCount)
				notifications.DELETE("/:id", handlers.DeleteNotification)
				notifications.GET("/stats", handlers.GetNotificationStats)

				// 管理员功能（需要通知管理权限）
				notifications.POST("/", middleware.RequirePermission("notification", "create"), handlers.CreateNotification)
				notifications.POST("/bulk", middleware.RequirePermission("notification", "bulk_send"), handlers.SendBulkNotification)
				notifications.GET("/templates", middleware.RequirePermission("notification", "view_templates"), handlers.GetNotificationTemplates)
			}

			// 权限管理（需要权限管理权限）
			permissions := protected.Group("/permissions")
			{
				permissions.POST("/check", handlers.CheckPermission)
				permissions.GET("/policies", middleware.RequirePermission("permission", "read"), handlers.GetPolicies)
				permissions.POST("/policies", middleware.RequirePermission("permission", "create"), handlers.AddPolicy)
				permissions.DELETE("/policies", middleware.RequirePermission("permission", "delete"), handlers.RemovePolicy)
				permissions.POST("/reload", middleware.RequirePermission("permission", "update"), handlers.ReloadPolicy)

				// 用户角色管理
				permissions.POST("/users/roles", middleware.RequirePermission("permission", "create"), handlers.AddRoleForUser)
				permissions.DELETE("/users/roles", middleware.RequirePermission("permission", "delete"), handlers.DeleteRoleForUser)
				permissions.GET("/users/:user/roles", middleware.RequirePermission("permission", "read"), handlers.GetRolesForUser)
				permissions.GET("/roles/:role/users", middleware.RequirePermission("permission", "read"), handlers.GetUsersForRole)
				permissions.GET("/roles/:role/permissions", handlers.GetPermissionsForRole)

				// 批量更新角色权限
				permissions.PUT("/roles/:role/permissions", middleware.RequirePermission("permission", "update"), handlers.UpdateRolePermissions)

				// 获取用户可访问的菜单列表
				permissions.GET("/menus", handlers.GetUserAccessibleMenus)
			}

			// 审计日志查询
			audit := protected.Group("/audit")
			{
				audit.GET("/logs", middleware.RequirePermission("audit", "read"), handlers.GetAuditLogs)
			}
		}
	}
}
