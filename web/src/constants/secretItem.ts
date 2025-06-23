import {
  Braces,
  Building,
  Cloud,
  Coins,
  Database,
  ExternalLink,
  GitBranch,
  Key,
  KeyRound,
  Lock,
  MessageCircle,
  MoreHorizontal,
  Server,
  Settings,
  Shield,
  Smartphone,
  Terminal,
  Ticket,
  User,
} from 'lucide-vue-next'

export type SecretItemTypeT = 'password' | 'api_key' | 'access_key' | 'ssh_key' | 'token' | 'custom'

export const SECRET_ITEM_TYPE = {
  Password: 'password',
  ApiKey: 'api_key',
  AccessKey: 'access_key',
  SshKey: 'ssh_key',
  Token: 'token',
  Custom: 'custom',
} as const

export const SECRET_ITEM_TYPE_MAP = {
  [SECRET_ITEM_TYPE.Password]: {
    label: '密码',
    icon: Lock,
  },
  [SECRET_ITEM_TYPE.ApiKey]: {
    label: 'API 密钥',
    icon: Key,
  },
  [SECRET_ITEM_TYPE.AccessKey]: {
    label: '访问密钥',
    icon: Cloud,
  },
  [SECRET_ITEM_TYPE.SshKey]: {
    label: 'SSH 密钥',
    icon: KeyRound,
  },
  [SECRET_ITEM_TYPE.Token]: {
    label: '令牌',
    icon: Coins,
  },
  [SECRET_ITEM_TYPE.Custom]: {
    label: '自定义',
    icon: Braces,
  },
}

export const SECRET_ITEM_CATEGORY = {
  // 密码类型分类
  WebsitePassword: {
    key: 'website_password',
    label: '网站密码',
    group: 'web_services',
    types: [SECRET_ITEM_TYPE.Password],
  },
  DatabasePassword: {
    key: 'database_password',
    label: '数据库密码',
    group: 'database',
    types: [SECRET_ITEM_TYPE.Password],
  },
  SystemPassword: {
    key: 'system_password',
    label: '系统密码',
    group: 'system',
    types: [SECRET_ITEM_TYPE.Password],
  },
  ApplicationPassword: {
    key: 'application_password',
    label: '应用密码',
    group: 'application',
    types: [SECRET_ITEM_TYPE.Password],
  },

  // API密钥类型分类
  CloudApiKey: {
    key: 'cloud_api_key',
    label: '云服务API密钥',
    group: 'cloud_services',
    types: [SECRET_ITEM_TYPE.ApiKey],
  },
  PaymentApiKey: {
    key: 'payment_api_key',
    label: '支付服务API密钥',
    group: 'payment',
    types: [SECRET_ITEM_TYPE.ApiKey],
  },
  SocialApiKey: {
    key: 'social_api_key',
    label: '社交平台API密钥',
    group: 'social',
    types: [SECRET_ITEM_TYPE.ApiKey],
  },
  ThirdPartyApiKey: {
    key: 'third_party_api_key',
    label: '第三方服务API密钥',
    group: 'third_party',
    types: [SECRET_ITEM_TYPE.ApiKey],
  },
  MiniProgramApiKey: {
    key: 'miniprogram_api_key',
    label: '小程序API密钥',
    group: 'mobile',
    types: [SECRET_ITEM_TYPE.ApiKey],
  },

  // 访问密钥类型分类
  AwsAccessKey: {
    key: 'aws_access_key',
    label: 'AWS访问密钥',
    group: 'cloud_services',
    types: [SECRET_ITEM_TYPE.AccessKey],
  },
  AzureAccessKey: {
    key: 'azure_access_key',
    label: 'Azure访问密钥',
    group: 'cloud_services',
    types: [SECRET_ITEM_TYPE.AccessKey],
  },
  GcpAccessKey: {
    key: 'gcp_access_key',
    label: 'GCP访问密钥',
    group: 'cloud_services',
    types: [SECRET_ITEM_TYPE.AccessKey],
  },
  AliCloudAccessKey: {
    key: 'alicloud_access_key',
    label: '阿里云访问密钥',
    group: 'cloud_services',
    types: [SECRET_ITEM_TYPE.AccessKey],
  },

  // SSH密钥类型分类
  ServerSshKey: {
    key: 'server_ssh_key',
    label: '服务器SSH密钥',
    group: 'server',
    types: [SECRET_ITEM_TYPE.SshKey],
  },
  GitSshKey: {
    key: 'git_ssh_key',
    label: 'Git仓库SSH密钥',
    group: 'development',
    types: [SECRET_ITEM_TYPE.SshKey],
  },
  DeviceSshKey: {
    key: 'device_ssh_key',
    label: '设备SSH密钥',
    group: 'device',
    types: [SECRET_ITEM_TYPE.SshKey],
  },

  // 令牌类型分类
  GitHubToken: {
    key: 'github_token',
    label: 'GitHub令牌',
    group: 'authentication',
    types: [SECRET_ITEM_TYPE.Token],
  },
  GoogleToken: {
    key: 'google_token',
    label: 'Google令牌',
    group: 'authentication',
    types: [SECRET_ITEM_TYPE.Token],
  },
  MicrosoftToken: {
    key: 'microsoft_token',
    label: 'Microsoft令牌',
    group: 'authentication',
    types: [SECRET_ITEM_TYPE.Token],
  },

  // 自定义类型分类
  CustomConfig: {
    key: 'custom_config',
    label: '自定义配置',
    group: 'custom',
    types: [SECRET_ITEM_TYPE.Custom],
  },
  CustomSecret: {
    key: 'custom_secret',
    label: '自定义密钥',
    group: 'custom',
    types: [SECRET_ITEM_TYPE.Custom],
  },
  CustomData: {
    key: 'custom_data',
    label: '自定义数据',
    group: 'custom',
    types: ['custom'],
  },
} as const

export type SecretItemCategoryGroupKey = keyof typeof SECRET_ITEM_CATEGORY_GROUPS

export const SECRET_ITEM_CATEGORY_GROUPS = {
  web_services: {
    key: 'web_services',
    label: '网站服务',
    icon: Cloud,
  },
  database: {
    key: 'database',
    label: '数据库',
    icon: Database,
  },
  system: {
    key: 'system',
    label: '系统',
    icon: Settings,
  },
  application: {
    key: 'application',
    label: '应用',
    icon: Building,
  },
  cloud_services: {
    key: 'cloud_services',
    label: '云服务',
    icon: Cloud,
  },
  payment: {
    key: 'payment',
    label: '支付服务',
    icon: Ticket,
  },
  social: {
    key: 'social',
    label: '社交平台',
    icon: MessageCircle,
  },
  third_party: {
    key: 'third_party',
    label: '第三方服务',
    icon: ExternalLink,
  },
  mobile: {
    key: 'mobile',
    label: '移动端',
    icon: Smartphone,
  },
  server: {
    key: 'server',
    label: '服务器',
    icon: Server,
  },
  development: {
    key: 'development',
    label: '开发工具',
    icon: GitBranch,
  },
  device: {
    key: 'device',
    label: '设备',
    icon: Terminal,
  },
  security: {
    key: 'security',
    label: '安全',
    icon: Shield,
  },
  authentication: {
    key: 'authentication',
    label: '认证',
    icon: User,
  },
  custom: {
    key: 'custom',
    label: '自定义',
    icon: MoreHorizontal,
  },
} as const

// 工具函数
// 根据类型获取可用的分类
export function getCategoriesByType(type: SecretItemTypeT) {
  return Object.values(SECRET_ITEM_CATEGORY).filter(category =>
    (category.types as unknown as SecretItemTypeT[]).includes(type),
  )
}

// 根据类型获取可用的分类组
export function getCategoryGroupsByType(type: SecretItemTypeT) {
  const categories = getCategoriesByType(type)
  const groupKeys = [...new Set(categories.map(cat => cat.group))]
  return groupKeys.map(key => SECRET_ITEM_CATEGORY_GROUPS[key as SecretItemCategoryGroupKey])
}

// 根据分组获取分类
export function getCategoriesByGroup(groupKey: SecretItemCategoryGroupKey, type?: SecretItemTypeT) {
  const allCategories = Object.values(SECRET_ITEM_CATEGORY).filter(category => category.group === groupKey)
  if (type) {
    return allCategories.filter(category => (category.types as unknown as SecretItemTypeT[]).includes(type))
  }
  return allCategories
}

// 获取所有类型
export function getAllTypes() {
  return Object.values(SECRET_ITEM_TYPE)
}

// 根据类型获取分组键
export function getGroupKeyByType(type: SecretItemTypeT) {
  const categories = getCategoriesByType(type)
  return [...new Set(categories.map(cat => cat.group))]
}

// 根据Key获取分类
export function getCategoryByKey(key?: string) {
  if (!key) {
    return null
  }
  return Object.values(SECRET_ITEM_CATEGORY).filter((category) => {
    return category.key === key
  })[0]
}
