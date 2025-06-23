/**
 * 格式化日期时间
 * @param dateTimeStr 日期时间字符串或时间戳
 * @returns 格式化后的日期时间字符串
 */
export function formatDateTime(dateTimeStr: string | number): string {
  if (!dateTimeStr)
    return ''

  const date = new Date(dateTimeStr)

  // 检查日期是否有效
  if (Number.isNaN(date.getTime()))
    return ''

  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

/**
 * 格式化日期
 * @param dateStr 日期字符串或时间戳
 * @returns 格式化后的日期字符串
 */
export function formatDate(dateStr: string | number): string {
  if (!dateStr)
    return ''

  const date = new Date(dateStr)

  // 检查日期是否有效
  if (Number.isNaN(date.getTime()))
    return ''

  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}
