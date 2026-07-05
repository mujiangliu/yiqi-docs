<template>
  <nav class="sidebar">
    <div class="sidebar-title">{{ title }}</div>
    <ul class="sidebar-list">
      <li v-for="item in tree" :key="item.id">
        <a
          :href="`/${sitePath}/${item.path}`"
          :class="{ active: item.path === currentPath }"
          @click.prevent="onNav(item.path)"
        >{{ item.title }}</a>
        <SidebarBranch
          v-if="item.children.length"
          :items="item.children"
          :site-path="sitePath"
          :current-path="currentPath"
          @navigate="onNav"
        />
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, type PropType, type VNode } from 'vue'
import type { PublicPage } from '@/api/types'

const props = defineProps<{
  title: string
  sitePath: string
  pages: PublicPage[]
  currentPath: string
}>()

const emit = defineEmits<{ (e: 'navigate', path: string): void }>()

interface TreeNode extends PublicPage {
  children: TreeNode[]
}

const SidebarBranch = defineComponent({
  name: 'SidebarBranch',
  props: {
    items: { type: Array as PropType<TreeNode[]>, required: true },
    sitePath: { type: String, required: true },
    currentPath: { type: String, required: true },
  },
  emits: ['navigate'],
  setup(branchProps, { emit }) {
    const renderItems = (items: TreeNode[]): VNode =>
      h('ul', { class: 'sidebar-list nested' }, items.map((item) =>
        h('li', { key: item.id }, [
          h('a', {
            href: `/${branchProps.sitePath}/${item.path}`,
            class: { active: item.path === branchProps.currentPath },
            onClick: (event: MouseEvent) => {
              event.preventDefault()
              emit('navigate', item.path)
            },
          }, item.title),
          item.children.length ? renderItems(item.children) : null,
        ]),
      ))
    return () => renderItems(branchProps.items)
  },
})

const tree = computed<TreeNode[]>(() => {
  const byId = new Map<number, TreeNode>()
  props.pages.forEach((p) => byId.set(p.id, { ...p, children: [] }))
  const roots: TreeNode[] = []
  byId.forEach((node) => {
    if (node.parent_id == null) {
      roots.push(node)
    } else {
      const parent = byId.get(node.parent_id)
      if (parent) parent.children.push(node)
    }
  })
  return roots
})

function onNav(path: string) {
  emit('navigate', path)
}
</script>

<style scoped>
.sidebar {
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
  padding: 1.5rem 1rem;
  background: var(--vp-sidebar-bg);
  border-right: 1px solid var(--vp-c-border);
}
.sidebar-title { font-weight: 600; font-size: 1.05rem; margin-bottom: 1rem; padding: 0 0.5rem; }
.sidebar-list { list-style: none; padding-left: 0.5rem; }
li { margin: 4px 0; }
a {
  display: block;
  padding: 4px 8px;
  border-radius: 4px;
  color: var(--vp-c-text-2);
  text-decoration: none;
  font-size: 0.92rem;
}
a:hover { color: var(--vp-c-brand-1); }
a.active { color: var(--vp-c-brand-1); font-weight: 500; background: var(--vp-c-bg-alt); }
.nested { padding-left: 0.8rem; }
.nested a { font-size: 0.88rem; }
</style>
