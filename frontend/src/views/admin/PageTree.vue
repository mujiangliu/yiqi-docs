<template>
  <div>
    <div class="page-header">
      <h2>页面管理 - {{ siteTitle }}</h2>
      <div>
        <router-link :to="`/admin/sites/${siteId}/pages/new`" class="btn">新建页面</router-link>
        <router-link :to="`/admin/sites`" class="back">返回站点</router-link>
      </div>
    </div>
    <div v-if="loading" class="loading">加载中…</div>
    <div v-else class="tree-container">
      <nested-draggable
        v-model="tree"
        :props="{ children: 'children' }"
        item-key="id"
        @end="onDragEnd"
      >
        <template #item="{ element }">
          <div class="tree-item">
            <span class="drag-handle">⋮⋮</span>
            <span class="title">{{ element.title }}</span>
            <code class="slug">/{{ element.slug }}</code>
            <div class="ops">
              <router-link :to="`/admin/sites/${siteId}/pages/${element.id}`">编辑</router-link>
              <button @click="onDelete(element)" class="link-btn">删除</button>
            </div>
          </div>
        </template>
      </nested-draggable>
      <div v-if="!pages.length" class="empty">还没有页面</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import draggable from 'vuedraggable'
import { adminPageApi, adminSiteApi } from '@/api/client'
import type { Page } from '@/api/types'

const route = useRoute()
const siteId = Number(route.params.id)
const pages = ref<Page[]>([])
const loading = ref(true)
const siteTitle = ref('')

interface TreeNode extends Page {
  children: TreeNode[]
}

const tree = computed<TreeNode[]>({
  get() {
    const byId = new Map<number, TreeNode>()
    pages.value.forEach((p) => byId.set(p.id, { ...p, children: [] }))
    const roots: TreeNode[] = []
    byId.forEach((node) => {
      if (node.parent_id == null) roots.push(node)
      else byId.get(node.parent_id)?.children.push(node)
    })
    return roots
  },
  set(newTree: TreeNode[]) {
    // 拖拽后扁平化
    const flat: Page[] = []
    function walk(nodes: TreeNode[], parentId: number | null) {
      nodes.forEach((n, i) => {
        flat.push({ ...n, parent_id: parentId, sort: i })
        walk(n.children, n.id)
      })
    }
    walk(newTree, null)
    pages.value = flat
  },
})

const NestedDraggable = draggable // 别名以便模板使用 kebab-case

async function load() {
  loading.value = true
  const [siteResp, pageResp] = await Promise.all([
    adminSiteApi.list(),
    adminPageApi.listBySite(siteId),
  ])
  const site = siteResp.data.data.find((s) => s.id === siteId)
  siteTitle.value = site?.title || ''
  pages.value = pageResp.data.data
  loading.value = false
}

async function onDragEnd() {
  // 收集所有 id，按当前顺序提交 reorder
  const ids = pages.value.map((p) => p.id)
  try {
    await adminPageApi.reorder(ids)
  } catch (e) {
    alert('排序保存失败')
    await load()
  }
}

async function onDelete(p: Page) {
  if (!confirm(`确认删除页面「${p.title}」？`)) return
  await adminPageApi.remove(p.id)
  await load()
}

onMounted(load)
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.page-header h2 { margin: 0; }
.btn { padding: 6px 14px; background: var(--vp-c-brand-1); color: #fff; border-radius: 4px; text-decoration: none; font-size: 0.9rem; margin-right: 8px; }
.back { font-size: 0.9rem; }
.tree-container { max-width: 720px; }
.tree-item {
  display: flex; align-items: center; gap: 0.6rem;
  padding: 8px 12px; border: 1px solid var(--vp-c-border); border-radius: 4px;
  margin-bottom: 4px; background: var(--vp-c-bg);
}
.drag-handle { cursor: grab; color: var(--vp-c-text-2); }
.title { font-weight: 500; }
.slug { color: var(--vp-c-text-2); font-size: 0.85rem; }
.ops { margin-left: auto; display: flex; gap: 8px; }
.ops a, .link-btn { font-size: 0.85rem; color: var(--vp-c-brand-1); }
.link-btn { background: none; border: none; cursor: pointer; padding: 0; }
.empty, .loading { text-align: center; color: var(--vp-c-text-2); padding: 2rem; }
code { background: var(--vp-c-bg-alt); padding: 2px 6px; border-radius: 4px; }
</style>
