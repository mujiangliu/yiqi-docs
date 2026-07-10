<template>
  <div v-if="siteStore.loading" class="loading">加载中…</div>
  <div v-else-if="siteStore.error" class="error">{{ siteStore.error }}</div>
  <DocLayout
    v-else-if="siteStore.site"
    :title="siteStore.site.title"
    :site-path="sitePath"
    :pages="siteStore.site.pages"
    :current-path="currentPagePath"
    @navigate="onNavigate"
  >
    <div v-if="siteStore.pageLoading" class="empty">页面加载中…</div>
    <div v-else-if="siteStore.pageError" class="error">{{ siteStore.pageError }}</div>
    <article v-else-if="currentPage" v-html="renderedHtml"></article>
    <div v-else class="empty">请从左侧选择一个页面</div>

    <template #toc>
      <Toc :items="toc.toc.value" :active-id="toc.activeId.value" />
    </template>
  </DocLayout>
</template>

<script setup lang="ts">
import { watch, computed, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import DocLayout from '@/components/DocLayout.vue'
import Toc from '@/components/Toc.vue'
import { useSiteStore } from '@/stores/site'
import { useToc } from '@/composables/useToc'

const route = useRoute()
const router = useRouter()
const siteStore = useSiteStore()
const toc = useToc()

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

const md = new MarkdownIt({
  html: true,
  linkify: true,
  highlight(str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre><code class="hljs language-${lang}">${hljs.highlight(str, { language: lang }).value}</code></pre>`
      } catch {}
    }
    return `<pre><code class="hljs">${escapeHtml(str)}</code></pre>`
  },
})

const sitePath = computed(() => route.params.sitePath as string)
const currentPagePath = computed(() => (route.params.pagePath as string) || '')
const currentPage = computed(() => siteStore.findPageByPath(currentPagePath.value))

const renderedHtml = computed(() => {
  if (!currentPage.value) return ''
  return md.render(currentPage.value.content_md || '')
})

async function loadSite() {
  await siteStore.load(sitePath.value)
  if (!currentPagePath.value && siteStore.site?.pages.length) {
    await router.replace(`/${sitePath.value}/${siteStore.site.pages[0].path}`)
    return
  }
  await loadCurrentPage()
}

async function loadCurrentPage() {
  if (!currentPagePath.value || !siteStore.site) return
  await siteStore.loadPage(sitePath.value, currentPagePath.value)
  await nextTick()
  extractToc()
}

function extractToc() {
  const el = document.querySelector('.markdown-body')
  if (el) toc.extractFrom(el as HTMLElement)
}

function onNavigate(path: string) {
  router.push(`/${sitePath.value}/${path}`)
}

watch(currentPagePath, () => {
  loadCurrentPage()
})

watch(renderedHtml, () => {
  nextTick(extractToc)
})

onMounted(loadSite)
watch(sitePath, loadSite)
</script>

<style scoped>
.loading, .error, .empty {
  max-width: 720px; margin: 4rem auto; text-align: center; color: var(--vp-c-text-2);
}
</style>
