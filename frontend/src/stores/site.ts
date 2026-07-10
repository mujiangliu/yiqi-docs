// frontend/src/stores/site.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { publicApi } from '@/api/client'
import type { PublicSite, PublicPage } from '@/api/types'

export const useSiteStore = defineStore('site', () => {
  const site = ref<PublicSite | null>(null)
  const loading = ref(false)
  const pageLoading = ref(false)
  const error = ref<string | null>(null)
  const pageError = ref<string | null>(null)

  async function load(path: string) {
    loading.value = true
    error.value = null
    try {
      const resp = await publicApi.getSite(path)
      site.value = resp.data.data
    } catch (e: any) {
      site.value = null
      error.value = e.response?.data?.error || '加载失败'
    } finally {
      loading.value = false
    }
  }

  function findPageByPath(pagePath: string): PublicPage | undefined {
    return site.value?.pages.find((p) => p.path === pagePath)
  }

  async function loadPage(sitePath: string, pagePath: string) {
    const page = findPageByPath(pagePath)
    if (!page || page.content_md !== undefined) return

    pageLoading.value = true
    pageError.value = null
    try {
      const resp = await publicApi.getPage(sitePath, pagePath)
      Object.assign(page, resp.data.data)
    } catch (e: any) {
      pageError.value = e.response?.data?.error || '页面加载失败'
    } finally {
      pageLoading.value = false
    }
  }

  return { site, loading, pageLoading, error, pageError, load, loadPage, findPageByPath }
})
