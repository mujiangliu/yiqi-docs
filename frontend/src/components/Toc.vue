<template>
  <nav class="toc" v-if="items.length">
    <div class="toc-title">本页目录</div>
    <ul>
      <li v-for="item in items" :key="item.id" :class="`level-${item.level}`">
        <a :href="`#${item.id}`" :class="{ active: item.id === activeId }">{{ item.text }}</a>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import type { TocItem } from '@/composables/useToc'
defineProps<{ items: TocItem[]; activeId: string }>()
</script>

<style scoped>
.toc {
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
  padding: 1.5rem 1rem;
  font-size: 0.85rem;
}
.toc-title { font-weight: 600; margin-bottom: 0.6rem; color: var(--vp-c-text-1); }
ul { list-style: none; padding-left: 0; }
li { margin: 4px 0; }
a { color: var(--vp-c-text-2); text-decoration: none; }
a:hover, a.active { color: var(--vp-c-brand-1); }
.level-3 { padding-left: 0.8rem; }
</style>
