<template>
  <el-dropdown @command="handleCommand" trigger="click">
    <el-button type="primary" text>
      {{ $t('langSwitcher.label') }}
      <el-icon class="el-icon--right"><arrow-down /></el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="zh-CN" :disabled="locale === 'zh-CN'">中文</el-dropdown-item>
        <el-dropdown-item command="en" :disabled="locale === 'en'">English</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import { ArrowDown } from '@element-plus/icons-vue'

const { locale } = useI18n();

const handleCommand = (command) => {
  if (locale.value === command) return;
  locale.value = command;
  localStorage.setItem('locale', command);
};
</script>

<style scoped>
.el-button.el-button--text {
  color: white;
}
.el-button.el-button--text:hover,
.el-button.el-button--text:focus {
  color: white;
  background-color: transparent;
}
</style>