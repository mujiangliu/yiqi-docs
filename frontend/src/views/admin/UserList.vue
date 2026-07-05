<template>
  <div>
    <div class="page-header">
      <h2>用户管理</h2>
      <button @click="showCreate = true" class="btn">新建用户</button>
    </div>
    <table class="user-table">
      <thead>
        <tr><th>ID</th><th>用户名</th><th>角色</th><th>操作</th></tr>
      </thead>
      <tbody>
        <tr v-for="u in users" :key="u.id">
          <td>{{ u.id }}</td>
          <td>{{ u.username }}</td>
          <td>{{ u.role === 'super_admin' ? '总管' : '内容管理员' }}</td>
          <td>
            <button @click="onReset(u)" class="link-btn">重置密码</button>
            <button @click="onDelete(u)" class="link-btn" v-if="u.id !== auth.user?.id">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showCreate" class="modal-bg" @click.self="showCreate = false">
      <div class="modal">
        <h3>新建用户</h3>
        <div class="field">
          <label>用户名</label>
          <input v-model="newUser.username" />
        </div>
        <div class="field">
          <label>密码</label>
          <input v-model="newUser.password" type="password" />
        </div>
        <div class="field">
          <label>角色</label>
          <select v-model="newUser.role">
            <option value="admin">内容管理员</option>
            <option value="super_admin">总管</option>
          </select>
        </div>
        <div class="actions">
          <button @click="onCreate" :disabled="creating">创建</button>
          <button @click="showCreate = false" class="btn-text">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminUserApi } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { User, Role } from '@/api/types'

const auth = useAuthStore()
const users = ref<User[]>([])
const showCreate = ref(false)
const creating = ref(false)
const newUser = ref({ username: '', password: '', role: 'admin' as Role })

async function load() {
  const resp = await adminUserApi.list()
  users.value = resp.data.data
}

async function onCreate() {
  creating.value = true
  try {
    await adminUserApi.create(newUser.value)
    showCreate.value = false
    newUser.value = { username: '', password: '', role: 'admin' }
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '创建失败')
  } finally {
    creating.value = false
  }
}

async function onReset(u: User) {
  const pw = prompt(`为 ${u.username} 设置新密码：`)
  if (!pw) return
  await adminUserApi.resetPassword(u.id, pw)
  alert('密码已重置')
}

async function onDelete(u: User) {
  if (!confirm(`确认删除用户「${u.username}」？`)) return
  try {
    await adminUserApi.remove(u.id)
    await load()
  } catch (e: any) {
    alert(e.response?.data?.error || '删除失败')
  }
}

onMounted(load)
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.btn { padding: 6px 14px; background: var(--vp-c-brand-1); color: #fff; border: none; border-radius: 4px; cursor: pointer; }
.user-table { width: 100%; border-collapse: collapse; }
.user-table th, .user-table td { padding: 10px; border-bottom: 1px solid var(--vp-c-border); text-align: left; }
.link-btn { background: none; border: none; color: var(--vp-c-brand-1); cursor: pointer; font-size: 0.9rem; margin-right: 8px; }
.modal-bg { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: var(--vp-c-bg); padding: 1.5rem 2rem; border-radius: 8px; min-width: 360px; }
.modal h3 { margin: 0 0 1rem; }
.field { margin-bottom: 1rem; }
.field label { display: block; margin-bottom: 4px; font-size: 0.9rem; color: var(--vp-c-text-2); }
.field input, .field select { width: 100%; padding: 8px; border: 1px solid var(--vp-c-border); border-radius: 4px; background: var(--vp-c-bg); color: var(--vp-c-text-1); box-sizing: border-box; }
.actions { display: flex; gap: 1rem; }
.actions button { padding: 6px 14px; background: var(--vp-c-brand-1); color: #fff; border: none; border-radius: 4px; cursor: pointer; }
.btn-text { background: none !important; color: var(--vp-c-text-2) !important; }
</style>
