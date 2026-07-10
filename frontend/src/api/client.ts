// frontend/src/api/client.ts
import axios from 'axios'
import type {
  User,
  Site,
  Page,
  PublicSite,
  MediaUploadResult,
  Role,
  SiteStatus,
} from './types'

export const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
  timeout: 30000,
})

// 响应拦截：统一解包 data
api.interceptors.response.use(
  (resp) => resp,
  (err) => {
    if (err.response?.status === 401) {
      // 跳登录
      if (!location.pathname.startsWith('/admin/login')) {
        location.href = '/admin/login'
      }
    }
    return Promise.reject(err)
  }
)

// 鉴权
export const authApi = {
  login: (username: string, password: string) =>
    api.post<{ data: User }>('/auth/login', { username, password }),
  logout: () => api.post('/auth/logout'),
  me: () => api.get<{ data: User }>('/me'),
}

// 公开
export const publicApi = {
  getSite: (path: string) => api.get<{ data: PublicSite }>(`/sites/${path}`),
  getPage: (sitePath: string, pagePath: string) =>
    api.get<{ data: PublicSite['pages'][number] }>(`/sites/${sitePath}/pages/${pagePath}`),
}

// 站点管理
export const adminSiteApi = {
  list: () => api.get<{ data: Site[] }>('/admin/sites'),
  create: (body: { path: string; title: string; description?: string; status?: SiteStatus }) =>
    api.post<{ data: Site }>('/admin/sites', body),
  update: (id: number, body: Partial<Pick<Site, 'path' | 'title' | 'description' | 'status'>>) =>
    api.put<{ data: Site }>(`/admin/sites/${id}`, body),
  remove: (id: number) => api.delete(`/admin/sites/${id}`),
  uploadMedia: (siteId: number, file: File) => {
    const fd = new FormData()
    fd.append('file', file)
    return api.post<{ data: MediaUploadResult }>(`/admin/sites/${siteId}/media`, fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

// 页面管理
export const adminPageApi = {
  listBySite: (siteId: number) => api.get<{ data: Page[] }>(`/admin/sites/${siteId}/pages`),
  get: (id: number) => api.get<{ data: Page }>(`/admin/pages/${id}`),
  create: (siteId: number, body: { parent_id?: number | null; slug: string; title: string; sort?: number; content_md?: string }) =>
    api.post<{ data: Page }>(`/admin/sites/${siteId}/pages`, body),
  update: (id: number, body: Partial<Pick<Page, 'parent_id' | 'slug' | 'title' | 'sort' | 'content_md'>>) =>
    api.put<{ data: Page }>(`/admin/pages/${id}`, body),
  remove: (id: number) => api.delete(`/admin/pages/${id}`),
  reorder: (ids: number[]) => api.post('/admin/pages/reorder', { ids }),
}

// 用户管理
export const adminUserApi = {
  list: () => api.get<{ data: User[] }>('/admin/users'),
  create: (body: { username: string; password: string; role: Role }) =>
    api.post<{ data: User }>('/admin/users', body),
  update: (id: number, body: { role?: Role }) => api.put<{ data: User }>(`/admin/users/${id}`, body),
  resetPassword: (id: number, password: string) =>
    api.post(`/admin/users/${id}/reset-password`, { password }),
  remove: (id: number) => api.delete(`/admin/users/${id}`),
}
