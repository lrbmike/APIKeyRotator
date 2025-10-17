<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div>{{ $t('layout.title') }}</div>
      <div class="header-right">
        <LangSwitcher />
        <el-button type="danger" @click="handleLogout" style="margin-left: 20px;">{{ $t('layout.logout') }}</el-button>
      </div>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { useRouter } from 'vue-router'
import LangSwitcher from '../components/LangSwitcher.vue'

const router = useRouter()

const handleLogout = () => {
  localStorage.removeItem('authToken')
  router.push({ name: 'Login' })
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-right {
  display: flex;
  align-items: center;
}
.header {
  background-color: #303133;
  color: #fff;
  font-size: 20px;
}
.header-right .el-button.el-button--danger:hover,
.header-right .el-button.el-button--danger:focus {
    background-color: var(--el-color-danger);
    border-color: var(--el-color-danger);
}
</style>