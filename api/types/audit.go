package types

const (
	AuditLogResourceUser          = "user"
	AuditLogResourceApiKey        = "api_key"
	AuditLogResourceAccessKey     = "access_key"
	AuditLogResourceSshKey        = "ssh_key"
	AuditLogResourcePassword      = "password"
	AuditLogResourceToken         = "token"
	AuditLogResourceCustom        = "custom"
	AuditLogResourceAccessRequest = "access_request"
)

const (
	AuditLogActionLogin   = "login"
	AuditLogActionLogout  = "logout"
	AuditLogActionCreate  = "create"
	AuditLogActionUpdate  = "update"
	AuditLogActionDelete  = "delete"
	AuditLogActionRead    = "read"
	AuditLogActionRequest = "request" // 申请访问
	AuditLogActionApprove = "approve" // 批准申请
	AuditLogActionReject  = "reject"  // 拒绝申请
	AuditLogActionRevoke  = "revoke"  // 撤销申请
	AuditLogActionAccess  = "access"  // 通过申请访问密钥
)
