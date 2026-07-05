// frontend/src/api/types.ts
export type Role = 'super_admin' | 'admin'
export type SiteStatus = 'published' | 'draft'

export interface User {
  id: number
  username: string
  role: Role
  created_at?: string
}

export interface Site {
  id: number
  owner_id: number
  path: string
  title: string
  description: string
  status: SiteStatus
  created_at?: string
  updated_at?: string
}

export interface Page {
  id: number
  site_id: number
  parent_id: number | null
  slug: string
  title: string
  sort: number
  content_md: string
  created_at?: string
  updated_at?: string
}

// 公开 API 返回的页面节点（多一个 path 字段）
export interface PublicPage {
  id: number
  parent_id: number | null
  slug: string
  title: string
  sort: number
  content_md: string
  path: string
}

export interface PublicSite {
  title: string
  description: string
  pages: PublicPage[]
}

export interface MediaUploadResult {
  hash: string
  url: string
}
