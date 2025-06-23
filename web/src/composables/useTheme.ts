import { useColorMode, usePreferredDark } from '@vueuse/core'
import { computed, watch } from 'vue'

export function useTheme() {
  const colorMode = useColorMode({
    emitAuto: true,
    modes: {
      light: 'light',
      dark: 'dark',
    },
  })

  const isDark = computed(() => colorMode.value === 'dark')
  const isLight = computed(() => colorMode.value === 'light')
  const isSystem = computed(() => colorMode.value === 'auto')

  const systemPrefersDark = usePreferredDark()

  // 当前实际使用的主题（考虑系统偏好）
  const currentTheme = computed(() => {
    if (colorMode.value === 'auto') {
      return systemPrefersDark.value ? 'dark' : 'light'
    }
    return colorMode.value
  })

  function toggleTheme() {
    colorMode.value = colorMode.value === 'dark' ? 'light' : 'dark'
  }

  function setTheme(theme: 'light' | 'dark' | 'auto') {
    colorMode.value = theme
  }

  function setLight() {
    colorMode.value = 'light'
  }

  function setDark() {
    colorMode.value = 'dark'
  }

  function setAuto() {
    colorMode.value = 'auto'
  }

  // 监听主题变化，更新文档类名
  watch(currentTheme, (theme) => {
    const root = document.documentElement
    root.classList.remove('light', 'dark')
    root.classList.add(theme)
  }, { immediate: true })

  return {
    colorMode,
    isDark,
    isLight,
    isSystem,
    currentTheme,
    systemPrefersDark,
    toggleTheme,
    setTheme,
    setLight,
    setDark,
    setAuto,
  }
}
