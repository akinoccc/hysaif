// API 请求和响应类型定义

// 通用响应类型
export interface ApiResponse<T = any> {
  data?: T
  message?: string
  error?: string
}

export interface ApiListResponse<T = any> extends ApiResponse<T[]> {
  pagination: Pagination
}

// 错误响应类型
export interface ErrorResponse {
  error: string
  message?: string
  details?: any
}

export interface Pagination {
  page: number
  page_size: number
  total: number
  total_pages: number
}

export interface ModelBase {
  id: string
  created_at: string
  updated_at: string
  created_by: string
  updated_by: string
  deleted_at: string
  deleted_by: string
  creator?: User
  updater?: User
  deleter?: User
}

// 分页参数
export interface PaginationParams {
  page?: number
  page_size?: number
}

// 认证相关类型
export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
  message: string
}

export interface LogoutResponse {
  message: string
}

// 用户相关类型
export interface User extends ModelBase {
  name: string
  avatar?: string
  email: string
  role: string
  status: string
  last_login_at: number
  last_login_ip: string
  failed_attempts?: number
  permissions?: string[]
}

export interface UpdateProfileRequest {
  email?: string
  name?: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

// 用户管理相关类型
export interface CreateUserRequest {
  name: string
  password: string
  email: string
  role: string
  permissions?: string[]
}

export interface UpdateUserRequest {
  name?: string
  email?: string
  role?: string
  status?: string
  permissions?: string[]
  password?: string
}

export interface UserListParams extends PaginationParams {
  username?: string
  role?: string
  status?: string
  page_size?: number
}

// 敏感信息项元数据类型
export interface SecretItemMeta {
  website_url?: string
  server_address?: string
  domain?: string
  api_docs?: string
  environment?: string
  owner?: string
  project?: string
  last_rotated?: string
}

// 敏感信息项数据类型
export interface SecretItemData {
  // 通用字段
  username?: string
  password?: string
  email?: string
  notes?: string

  // API Key 相关
  api_key?: string
  api_secret?: string

  // Access Key 相关
  access_key?: string
  secret_key?: string
  region?: string

  // SSH Key 相关
  private_key?: string
  public_key?: string
  passphrase?: string
  key_type?: string

  // Token 相关
  token?: string
  refresh_token?: string
  token_type?: string
  scope?: string

  // 自定义字段
  custom_data?: Record<string, any>
}

export interface SecretBaseInfo {
  name: string
  description?: string
  type: string
  environment?: string
  category?: string
  tags?: string[]
  expires_at?: number
}

// 信息项相关类型
export interface SecretItem<D = any> extends SecretBaseInfo {
  id: string
  data: D
  created_at: string
  updated_at: string
  created_by: string
  updated_by: string

  // 访问权限相关
  has_approved_access?: boolean // 是否有已批准的访问申请

  // 关联字段
  creator: User
  updater: User
}

export type PostItemRequest<D> = Omit<
  SecretItem<D>,
  'id' | 'created_at' | 'updated_at' | 'created_by' | 'updated_by' | 'creator' | 'updater'
>

export interface ItemsListParams extends PaginationParams {
  category?: string
  search?: string
  tag?: string
  creator_name?: string
  environment?: string
  status?: string
  created_at_from?: number
  created_at_to?: number
  sort_by?: string
  page_size?: number
}

// 审计日志相关类型
export interface AuditLog extends ModelBase {
  user_id: string
  user: User
  action: string
  resource: string
  resource_id: string
  details?: string
  ip_address: string
  user_agent: string
}

export interface AuditLogsParams extends PaginationParams {
  user_id?: string
  user?: string
  action?: string
  resource?: string
  start_date?: number
  end_date?: number
  page_size?: number
}

// 访问申请相关类型
export interface AccessRequest extends ModelBase {
  secret_item_id: string
  applicant_id: string
  reason: string
  status: string // pending, approved, rejected, expired, revoked
  approved_by?: string
  approved_at?: number
  reject_reason?: string
  valid_from?: number
  valid_until?: number
  access_count: number
  last_accessed?: number

  // 关联字段
  secret_item: SecretItem
  applicant: User
  approver?: User
}

export interface CreateAccessRequestRequest {
  secret_item_id: string
  reason: string
}

export interface ApproveAccessRequestRequest {
  valid_duration: number // 有效时长（小时）
  note?: string // 审批备注
}

export interface RejectAccessRequestRequest {
  reason: string // 拒绝理由
}

export interface RevokeAccessRequestRequest {
  reason: string // 作废理由
}

export interface AccessRequestListParams extends PaginationParams {
  status?: string
  applicant_name?: string
  secret_item_name?: string
  sort_by?: string
  created_at_from?: string
  created_at_to?: string
  page?: number
  page_size?: number
}

// API 方法返回类型
export type ApiMethod<T = any> = Promise<T>
