import { Edit, Eye, LogIn, LogOut, Plus, Shield, Trash2 } from 'lucide-vue-next'
import { AUDIT_LOG_ACTION_LIST, AUDIT_LOG_ACTION_MAP, AUDIT_LOG_RESOURCE_MAP } from '@/constants'

export function getActionIcon(action: string) {
  const icons: Record<string, any> = {
    [AUDIT_LOG_ACTION_LIST.Login]: LogIn,
    [AUDIT_LOG_ACTION_LIST.Logout]: LogOut,
    [AUDIT_LOG_ACTION_LIST.Create]: Plus,
    [AUDIT_LOG_ACTION_LIST.Read]: Eye,
    [AUDIT_LOG_ACTION_LIST.Update]: Edit,
    [AUDIT_LOG_ACTION_LIST.Delete]: Trash2,
  }
  return icons[action] || Shield
}

export function getActionDisplayName(action: string) {
  return AUDIT_LOG_ACTION_MAP[action as keyof typeof AUDIT_LOG_ACTION_MAP] || action
}

export function getResourceDisplayName(resource: string) {
  return AUDIT_LOG_RESOURCE_MAP[resource as keyof typeof AUDIT_LOG_RESOURCE_MAP] || resource
}

export function getActionColor(action: string) {
  const colors: Record<string, string> = {
    [AUDIT_LOG_ACTION_LIST.Login]: 'bg-success/10 text-success',
    [AUDIT_LOG_ACTION_LIST.Logout]: 'bg-muted text-muted-foreground',
    [AUDIT_LOG_ACTION_LIST.Create]: 'bg-info/10 text-info',
    [AUDIT_LOG_ACTION_LIST.Read]: 'bg-primary/10 text-primary',
    [AUDIT_LOG_ACTION_LIST.Update]: 'bg-warning/10 text-warning',
    [AUDIT_LOG_ACTION_LIST.Delete]: 'bg-destructive/10 text-destructive',
  }
  return colors[action] || 'bg-muted text-muted-foreground'
}

export function getUserAgentInfo(userAgent: string) {
  if (userAgent.includes('Chrome'))
    return 'Chrome'
  if (userAgent.includes('Firefox'))
    return 'Firefox'
  if (userAgent.includes('Safari'))
    return 'Safari'
  if (userAgent.includes('Edge'))
    return 'Edge'
  if (userAgent.includes('Mobile'))
    return '移动设备'
  return '未知浏览器'
}
