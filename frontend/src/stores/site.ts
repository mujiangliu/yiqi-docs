// frontend/src/stores/site.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { publicApi } from '@/api/client'
import type { PublicSite, PublicPage } from '@/api/types'

export const useSiteStore = defineStore('site', () => {
  const site = ref<PublicSite | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

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

  return { site, loading, error, load, findPageByPath }
})
