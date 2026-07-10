<template>
  <div>
    <div class="page-header">
      <h2>{{ isEdit ? '编辑页面' : '新建页面' }}</h2>
      <router-link :to="`/admin/sites/${siteId}/pages`" class="back">返回列表</router-link>
    </div>
    <form @submit.prevent="onSubmit" class="page-form">
      <div class="row">
        <div class="field">
          <label>标题</label>
          <input v-model="form.title" required />
        </div>
        <div class="field">
          <label>slug</label>
          <input v-model="form.slug" required placeholder="intro" />
        </div>
      </div>
      <div class="field">
        <label>内容 (Markdown)</label>
        <MdEditor
          v-model="form.content_md"
          :theme="editorTheme"
          :style="{ height: '480px' }"
          :toolbars-exclude="['github', 'save']"
          @on-upload-img="onUploadImg"
        />
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
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { adminPageApi, adminSiteApi } from '@/api/client'
import { useMdEditorTheme } from '@/composables/useMdEditorTheme'

const route = useRoute()
const router = useRouter()
const siteId = Number(route.params.siteId)
const pageId = computed(() => route.params.id as string | undefined)
const isEdit = computed(() => !!pageId.value)
const submitting = ref(false)
const error = ref('')
const editorTheme = useMdEditorTheme()

const form = ref({
  title: '',
  slug: '',
  content_md: '',
  parent_id: null as number | null,
})

onMounted(async () => {
  if (isEdit.value) {
    const resp = await adminPageApi.get(Number(pageId.value))
    const page = resp.data.data
    form.value = {
      title: page.title,
      slug: page.slug,
      content_md: page.content_md || '',
      parent_id: page.parent_id,
    }
  }
})

async function onSubmit() {
  submitting.value = true
  error.value = ''
  try {
    if (isEdit.value) {
      await adminPageApi.update(Number(pageId.value), form.value)
    } else {
      await adminPageApi.create(siteId, { ...form.value, parent_id: form.value.parent_id })
    }
    router.push(`/admin/sites/${siteId}/pages`)
  } catch (e: any) {
    error.value = e.response?.data?.error || '保存失败'
  } finally {
    submitting.value = false
  }
}

// md-editor-v3 图片上传回调
async function onUploadImg(files: File[], callback: (urls: string[]) => void) {
  const urls: string[] = []
  for (const f of files) {
    const resp = await adminSiteApi.uploadMedia(siteId, f)
    urls.push(resp.data.data.url)
  }
  callback(urls)
}
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.back { font-size: 0.9rem; }
.page-form { max-width: 920px; }
.row { display: flex; gap: 1rem; }
.row .field { flex: 1; }
.field { margin-bottom: 1.2rem; }
.field label { display: block; margin-bottom: 6px; font-size: 0.9rem; color: var(--vp-c-text-2); }
.field input { width: 100%; padding: 8px 10px; border: 1px solid var(--vp-c-border); border-radius: 4px; background: var(--vp-c-bg); color: var(--vp-c-text-1); box-sizing: border-box; }
.actions { display: flex; align-items: center; gap: 1rem; }
button { padding: 8px 18px; background: var(--vp-c-brand-1); color: #fff; border: none; border-radius: 4px; cursor: pointer; }
button:disabled { opacity: 0.6; }
.error { color: #d33; font-size: 0.9rem; }
</style>
