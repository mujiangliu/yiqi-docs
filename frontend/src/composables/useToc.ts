// frontend/src/composables/useToc.ts
import { ref, nextTick } from 'vue'

export interface TocItem {
  id: string
  text: string
  level: number
}

export function useToc() {
  const toc = ref<TocItem[]>([])
  const activeId = ref<string>('')

  function extractFrom(el: HTMLElement) {
    const headings = el.querySelectorAll('h2, h3')
    const items: TocItem[] = []
    headings.forEach((h) => {
      const id = h.id || slugify(h.textContent || '')
      h.id = id
      items.push({
        id,
        text: h.textContent || '',
        level: h.tagName === 'H2' ? 2 : 3,
      })
    })
    toc.value = items
    nextTick(setupScrollSpy)
  }

  function slugify(s: string): string {
    return s.trim().toLowerCase().replace(/\s+/g, '-')
  }

  function setupScrollSpy() {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((e) => {
          if (e.isIntersecting) activeId.value = e.target.id
        })
      },
      { rootMargin: '-80px 0px -70% 0px' }
    )
    toc.value.forEach((item) => {
      const el = document.getElementById(item.id)
      if (el) observer.observe(el)
    })
  }

  return { toc, activeId, extractFrom }
}
