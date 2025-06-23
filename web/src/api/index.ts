// API 模块化入口文件

export { auditAPI } from './audit'

// 导出各模块API
export { authAPI } from './auth'

// 导出配置
export { default as api, API_BASE_URL } from './config'
// 默认导出axios实例（保持向后兼容）
export { default } from './config'
// 导出权限API
export { permissionAPI } from './permission'
export { secretItemAPI as itemAPI } from './secret'
// 导出类型定义
export * from './types'

export { userAPI } from './user'
