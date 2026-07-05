<template>
  <div>
    <div class="page-header">
      <h2>我的站点</h2>
      <router-link to="/admin/sites/new" class="btn">新建站点</router-link>
    </div>
    <div v-if="loading" class="loading">加载中…</div>
    <table v-else class="site-table">
      <thead>
        <tr><th>路径</th><th>标题</th><th>状态</th><th>操作</th></tr>
      </thead>
      <tbody>
        <tr v-for="s in sites" :key="s.id">
          <td><code>/{{ s.path }}</code></td>
          <td>{{ s.title }}</td>
          <td>
            <span :class="['badge', s.status]">{{ s.status === 'published' ? '已发布' : '草稿' }}</span>
          </td>
          <td class="actions">
            <router-link :to="`/admin/sites/${s.id}`" class="action-link">编辑</router-link>
            <router-link :to="`/admin/sites/${s.id}/pages`" class="action-link">页面</router-link>
            <button @click="onDelete(s)" class="link-btn">删除</button>
          </td>
        </tr>
        <tr v-if="!sites.length"><td colspan="4" class="empty">还没有站点，点击右上角新建</td></tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminSiteApi } from '@/api/client'
import type { Site } from '@/api/types'

const sites = ref<Site[]>([])
const loading = ref(true)

async function load() {
  loading.value = true
  const resp = await adminSiteApi.list()
  sites.value = resp.data.data
  loading.value = false
}

async function onDelete(s: Site) {
  if (!confirm(`确认删除站点「${s.title}」？所有页面与媒体将一并删除。`)) return
  await adminSiteApi.remove(s.id)
  await load()
}

onMounted(load)
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.page-header h2 { margin: 0; }
.btn { padding: 6px 14px; background: var(--vp-c-brand-1); color: #fff; border-radius: 4px; text-decoration: none; font-size: 0.9rem; }
.site-table { width: 100%; border-collapse: collapse; }
.site-table th, .site-table td { padding: 10px; border-bottom: 1px solid var(--vp-c-border); text-align: left; }
.site-table th { color: var(--vp-c-text-2); font-weight: 500; font-size: 0.9rem; }
.site-table td { font-size: 0.95rem; }
.badge { padding: 2px 8px; border-radius: 10px; font-size: 0.8rem; }
.badge.published { background: #e6f4e6; color: #2e7d32; }
.badge.draft { background: #fff4e6; color: #b86c00; }
.dark .badge.published { background: #1e3a1e; color: #81c784; }
.dark .badge.draft { background: #3a2e1e; color: #ffb74d; }
.link-btn { background: none; border: none; color: var(--vp-c-brand-1); cursor: pointer; font-size: 0.9rem; padding: 0; }
.empty, .loading { text-align: center; color: var(--vp-c-text-2); padding: 2rem; }
code { background: var(--vp-c-bg-alt); padding: 2px 6px; border-radius: 4px; }
.actions { white-space: nowrap; }
.action-link { display: inline-block; margin-right: 12px; color: var(--vp-c-brand-1); text-decoration: none; font-size: 0.9rem; }
</style>
