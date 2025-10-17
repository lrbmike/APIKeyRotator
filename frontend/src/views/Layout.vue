<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div>API Key Rotator</div>
      <el-button type="danger" @click="handleLogout">退出登录</el-button>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { useRouter } from 'vue-router'

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
  background-color: #303133;
  color: #fff;
  font-size: 20px;
}
</style>