import { createI18n } from 'vue-i18n';
import en from './locales/en.json';
import zhCN from './locales/zh-CN.json';

const i18n = createI18n({
  legacy: false, // 使用 Composition API，必须设置为 false
  locale: localStorage.getItem('locale') || 'zh-CN', // 默认语言
  fallbackLocale: 'en', // 回退语言
  messages: {
    'en': en,
    'zh-CN': zhCN,
  },
});

export default i18n;