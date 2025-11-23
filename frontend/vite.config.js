import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
// 将整个导出对象包裹在一个函数中，这样 Vite 会将 mode 传递进来
export default defineConfig(({ mode }) => {

  // process.cwd() 获取当前工作目录，loadEnv 会在这里寻找 .env 文件
  const env = loadEnv(mode, process.cwd());

  console.log('Vite Proxy Target URL:', env.VITE_API_TARGET_URL);

  // 返回最终的配置对象
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src')
      }
    },
    build: {
      rollupOptions: {
        external: [],
        output: {
          manualChunks: undefined
        }
      },
      commonjsOptions: {
        include: [/node_modules/]
      }
    },
    server: {
      // 监听所有网络接口，以便 Docker 容器的端口可以被映射出去
      host: '0.0.0.0',
      // 指定开发服务器端口，与 docker-compose.yml 中保持一致
      port: 5173,
      proxy: {
        // 将所有 /admin 开头的请求代理到后端服务
        '/admin': {
          target: env.VITE_API_TARGET_URL || 'http://localhost:8000',
          changeOrigin: true,
        },
        // 将所有 /llm 开头的请求代理到后端服务
        '/llm': {
          target: env.VITE_API_TARGET_URL || 'http://localhost:8000',
          changeOrigin: true,
        },
        // 保留 /api 代理规则以实现兼容性
        '/api': {
          target: env.VITE_API_TARGET_URL || 'http://localhost:8000',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        }
      }
    }
  }
})
