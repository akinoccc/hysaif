import { SECRET_ITEM_TYPE } from './secretItem'

export const AUDIT_LOG_RESOURCE_LIST = {
  User: 'user',
  AccessRequest: 'access_request',
  Password: SECRET_ITEM_TYPE.Password,
  ApiKey: SECRET_ITEM_TYPE.ApiKey,
  AccessKey: SECRET_ITEM_TYPE.AccessKey,
  SshKey: SECRET_ITEM_TYPE.SshKey,
  Token: SECRET_ITEM_TYPE.Token,
  Custom: SECRET_ITEM_TYPE.Custom,
} as const

export const AUDIT_LOG_RESOURCE_MAP = {
  [AUDIT_LOG_RESOURCE_LIST.User]: '用户',
  [AUDIT_LOG_RESOURCE_LIST.AccessRequest]: '访问申请',
  [AUDIT_LOG_RESOURCE_LIST.Password]: '密码',
  [AUDIT_LOG_RESOURCE_LIST.ApiKey]: 'API密钥',
  [AUDIT_LOG_RESOURCE_LIST.AccessKey]: '访问密钥',
  [AUDIT_LOG_RESOURCE_LIST.SshKey]: 'SSH密钥',
  [AUDIT_LOG_RESOURCE_LIST.Token]: '令牌',
  [AUDIT_LOG_RESOURCE_LIST.Custom]: '自定义',
}

export const AUDIT_LOG_ACTION_LIST = {
  Login: 'login',
  Logout: 'logout',

  Create: 'create',
  Update: 'update',
  Delete: 'delete',
  Read: 'read',
}

export const AUDIT_LOG_ACTION_MAP = {
  [AUDIT_LOG_ACTION_LIST.Login]: '登录',
  [AUDIT_LOG_ACTION_LIST.Logout]: '退出登录',
  [AUDIT_LOG_ACTION_LIST.Create]: '创建',
  [AUDIT_LOG_ACTION_LIST.Update]: '更新',
  [AUDIT_LOG_ACTION_LIST.Delete]: '删除',
  [AUDIT_LOG_ACTION_LIST.Read]: '查看',
}
