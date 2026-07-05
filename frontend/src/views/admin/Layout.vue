<template>
  <div class="admin-layout">
    <aside class="admin-side">
      <div class="brand">文档后台</div>
      <nav>
        <router-link to="/admin/sites">我的站点</router-link>
        <router-link v-if="auth.isSuper()" to="/admin/users">用户管理</router-link>
      </nav>
      <div class="user-info">
        <span>{{ auth.user?.username }} ({{ auth.user?.role }})</span>
        <button @click="onLogout">退出</button>
      </div>
    </aside>
    <main class="admin-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

onMounted(() => {
  if (!auth.user) auth.fetchMe()
})

async function onLogout() {
  await auth.logout()
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout { display: flex; min-height: 100vh; background: var(--vp-c-bg); }
.admin-side {
  width: 220px; background: var(--vp-c-bg-alt); padding: 1.5rem 1rem;
  display: flex; flex-direction: column; gap: 1rem; border-right: 1px solid var(--vp-c-border);
}
.brand { font-weight: 600; font-size: 1.1rem; margin-bottom: 1rem; }
nav { display: flex; flex-direction: column; gap: 4px; }
nav a { padding: 8px 12px; border-radius: 4px; text-decoration: none; color: var(--vp-c-text-2); }
nav a.router-link-active { background: var(--vp-c-brand-1); color: #fff; }
.user-info { margin-top: auto; font-size: 0.85rem; color: var(--vp-c-text-2); display: flex; flex-direction: column; gap: 8px; }
.user-info button { padding: 6px; background: none; border: 1px solid var(--vp-c-border); border-radius: 4px; cursor: pointer; color: var(--vp-c-text-1); }
.admin-main { flex: 1; padding: 2rem; overflow: auto; }
</style>
