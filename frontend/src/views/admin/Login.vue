<template>
  <div class="login-page">
    <form @submit.prevent="onSubmit" class="login-form">
      <h1>后台登录</h1>
      <div class="field">
        <label>用户名</label>
        <input v-model="username" required autocomplete="username" />
      </div>
      <div class="field">
        <label>密码</label>
        <input v-model="password" type="password" required autocomplete="current-password" />
      </div>
      <button type="submit" :disabled="submitting">{{ submitting ? '登录中…' : '登录' }}</button>
      <div v-if="error" class="error">{{ error }}</div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const submitting = ref(false)
const error = ref('')

async function onSubmit() {
  submitting.value = true
  error.value = ''
  try {
    await auth.login(username.value, password.value)
    const redirect = (route.query.redirect as string) || '/admin/sites'
    router.push(redirect)
  } catch (e: any) {
    error.value = e.response?.data?.error || '登录失败'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.login-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: var(--vp-c-bg-alt); }
.login-form { background: var(--vp-c-bg); padding: 2rem; border-radius: 8px; border: 1px solid var(--vp-c-border); min-width: 320px; }
.login-form h1 { margin: 0 0 1.5rem; font-size: 1.4rem; }
.field { margin-bottom: 1rem; }
.field label { display: block; margin-bottom: 4px; font-size: 0.9rem; color: var(--vp-c-text-2); }
.field input { width: 100%; padding: 8px 10px; border: 1px solid var(--vp-c-border); border-radius: 4px; font-size: 0.95rem; background: var(--vp-c-bg); color: var(--vp-c-text-1); box-sizing: border-box; }
button { width: 100%; padding: 10px; background: var(--vp-c-brand-1); color: #fff; border: none; border-radius: 4px; font-size: 0.95rem; cursor: pointer; }
button:disabled { opacity: 0.6; }
.error { margin-top: 1rem; color: #d33; font-size: 0.9rem; }
</style>
