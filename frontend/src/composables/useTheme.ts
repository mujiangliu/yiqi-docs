// frontend/src/composables/useTheme.ts
import { ref, watch } from 'vue'

const theme = ref<'light' | 'dark'>(localStorage.getItem('theme') as 'light' | 'dark' || 'light')

function applyClass() {
  if (theme.value === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

export function useTheme() {
  watch(theme, (v) => {
    localStorage.setItem('theme', v)
    applyClass()
  })
  applyClass()

  function toggle() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
  }

  return { theme, toggle }
}
