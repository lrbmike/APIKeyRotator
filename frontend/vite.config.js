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
    server: {
      // 监听所有网络接口，以便 Docker 容器的端口可以被映射出去
      host: '0.0.0.0', 
      // 指定开发服务器端口，与 docker-compose.yml 中保持一致
      port: 5173, 
      proxy: {
        // 将所有 /api 开头的请求代理到后端服务
        '/api': {
          // 现在 env.VITE_API_TARGET_URL 可以被正确读取了
          target: env.VITE_API_TARGET_URL || 'http://localhost:8000',
          
          // changeOrigin: true 对于代理是必需的，它会修改请求头中的 Host，
          // 使其与目标服务器匹配，很多后端服务都需要这个配置。
          changeOrigin: true,

          // 如果您的后端 API 路由本身就包含了 /api (例如 /api/admin/login)，
          // 那么就不需要 rewrite。
          // 如果后端路由是 /admin/login，而前端请求的是 /api/admin/login，
          // 那么就需要取消下面的注释来移除 /api 前缀。
          rewrite: (path) => path.replace(/^\/api/, ''),
        }
      }
    }
  }
})