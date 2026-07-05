<template>
  <button class="theme-toggle" @click="toggle" :title="theme === 'dark' ? '切换浅色' : '切换暗色'">
    {{ theme === 'dark' ? '☀️' : '🌙' }}
  </button>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'
const { theme, toggle } = useTheme()
</script>

<style scoped>
.theme-toggle {
  background: none;
  border: 1px solid var(--vp-c-border);
  border-radius: 6px;
  padding: 4px 10px;
  cursor: pointer;
  color: var(--vp-c-text-1);
}
</style>
