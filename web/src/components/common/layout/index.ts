export { default as DataFilter } from './DataFilter.vue'
export { default as EmptyState } from './EmptyState.vue'
export { default as PageHeader } from './PageHeader.vue'

export interface FilterField {
  key: string
  label: string
  type: 'text' | 'select' | 'date-range' | 'user-search'
  placeholder?: string
  icon?: any
  options?: Array<{
    value: string
    label: string
    group?: string
  }>
  groups?: Array<{
    key: string
    label: string
    icon?: any
    options: Array<{
      value: string
      label: string
    }>
  }>
}

export interface ActionButton {
  text: string
  icon?: any
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  permission?: { resource: string, action: string }
  onClick?: () => void
  to?: string
}
