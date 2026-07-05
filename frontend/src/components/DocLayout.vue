<template>
  <div class="doc-layout">
    <header class="doc-header">
      <div class="header-inner">
        <a :href="`/${sitePath}`" class="brand">{{ title }}</a>
        <div class="header-right">
          <a href="/admin">后台</a>
          <ThemeToggle />
        </div>
      </div>
    </header>
    <div class="doc-body">
      <Sidebar
        :title="title"
        :site-path="sitePath"
        :pages="pages"
        :current-path="currentPath"
        @navigate="onNavigate"
      />
      <main class="doc-main">
        <div class="markdown-body" ref="contentRef"><slot /></div>
      </main>
      <aside class="doc-toc"><slot name="toc" /></aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './Sidebar.vue'
import ThemeToggle from './ThemeToggle.vue'
import type { PublicPage } from '@/api/types'

const props = defineProps<{
  title: string
  sitePath: string
  pages: PublicPage[]
  currentPath: string
}>()

void props

const emit = defineEmits<{ (e: 'navigate', path: string): void }>()
const contentRef = ref<HTMLElement | null>(null)

function onNavigate(path: string) {
  emit('navigate', path)
}

defineExpose({ contentRef })
</script>

<style scoped>
.doc-layout { min-height: 100vh; background: var(--vp-c-bg); }
.doc-header {
  position: sticky; top: 0; z-index: 10;
  background: var(--vp-c-bg);
  border-bottom: 1px solid var(--vp-c-border);
}
.header-inner {
  max-width: 1440px; margin: 0 auto; padding: 0 1.5rem;
  height: 56px; display: flex; align-items: center; justify-content: space-between;
}
.brand { font-weight: 600; font-size: 1.1rem; text-decoration: none; color: var(--vp-c-text-1); }
.header-right { display: flex; gap: 1rem; align-items: center; }
.header-right a { font-size: 0.9rem; }
.doc-body {
  max-width: 1440px; margin: 0 auto; display: grid;
  grid-template-columns: 272px 1fr 240px; gap: 1rem; padding: 0 1rem;
}
.doc-main { padding: 2rem 1.5rem; min-width: 0; }
.doc-toc { min-width: 0; }
@media (max-width: 1024px) {
  .doc-body { grid-template-columns: 272px 1fr; }
  .doc-toc { display: none; }
}
@media (max-width: 768px) {
  .doc-body { grid-template-columns: 1fr; }
  :deep(.sidebar) { display: none; }
}
</style>
