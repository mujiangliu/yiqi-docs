<template>
  <div>
    <div class="page-header">
      <h2>{{ isEdit ? '编辑站点' : '新建站点' }}</h2>
      <router-link :to="`/admin/sites`" class="back">返回</router-link>
    </div>
    <form @submit.prevent="onSubmit" class="site-form">
      <div class="field">
        <label>路径 (slug)</label>
        <input v-model="form.path" required placeholder="如 yiqicuqu-api" />
        <div class="hint">访问路径：<code>/{{ form.path }}</code></div>
      </div>
      <div class="field">
        <label>标题</label>
        <input v-model="form.title" required />
      </div>
      <div class="field">
        <label>简介</label>
        <textarea v-model="form.description" rows="3"></textarea>
      </div>
      <div class="field">
        <label>状态</label>
        <select v-model="form.status">
          <option value="draft">草稿</option>
          <option value="published">已发布</option>
        </select>
      </div>
      <div class="actions">
        <button type="submit" :disabled="submitting">{{ submitting ? '保存中…' : '保存' }}</button>
        <span v-if="error" class="error">{{ error }}</span>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { adminSiteApi } from '@/api/client'
import type { SiteStatus } from '@/api/types'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const submitting = ref(false)
const error = ref('')

const form = ref({
  path: '',
  title: '',
  description: '',
  status: 'draft' as SiteStatus,
})

onMounted(async () => {
  if (isEdit.value) {
    const id = Number(route.params.id)
    const resp = await adminSiteApi.list()
    const site = resp.data.data.find((s) => s.id === id)
    if (site) {
      form.value = { path: site.path, title: site.title, description: site.description, status: site.status }
    }
  }
})

async function onSubmit() {
  submitting.value = true
  error.value = ''
  try {
    if (isEdit.value) {
      await adminSiteApi.update(Number(route.params.id), form.value)
    } else {
      await adminSiteApi.create(form.value)
    }
    router.push('/admin/sites')
  } catch (e: any) {
    error.value = e.response?.data?.error || '保存失败'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.back { font-size: 0.9rem; }
.site-form { max-width: 560px; }
.field { margin-bottom: 1.2rem; }
.field label { display: block; margin-bottom: 6px; font-size: 0.9rem; color: var(--vp-c-text-2); }
.field input, .field textarea, .field select {
  width: 100%; padding: 8px 10px; border: 1px solid var(--vp-c-border);
  border-radius: 4px; background: var(--vp-c-bg); color: var(--vp-c-text-1);
  font-size: 0.95rem; box-sizing: border-box;
}
.hint { margin-top: 4px; font-size: 0.82rem; color: var(--vp-c-text-2); }
.actions { display: flex; align-items: center; gap: 1rem; }
button { padding: 8px 18px; background: var(--vp-c-brand-1); color: #fff; border: none; border-radius: 4px; cursor: pointer; }
button:disabled { opacity: 0.6; }
.error { color: #d33; font-size: 0.9rem; }
code { background: var(--vp-c-bg-alt); padding: 2px 6px; border-radius: 4px; }
</style>
