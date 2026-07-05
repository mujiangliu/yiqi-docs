// frontend/src/composables/useMdEditorTheme.ts
import { computed } from 'vue'
import { useTheme } from './useTheme'

export function useMdEditorTheme() {
  const { theme } = useTheme()
  const editorTheme = computed<'light' | 'dark'>(() => theme.value)
  return editorTheme
}
