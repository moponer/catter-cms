import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vite.dev/config/
export default defineConfig({
  // 管理后台基础路径
  base: '/admin/',
  plugins: [
    vue(),
    // Element Plus 按需自动导入
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: ['vue', 'vue-router', 'pinia'],
      dts: 'src/auto-imports.d.ts',
    }),
    Components({
      resolvers: [ElementPlusResolver()],
      dts: 'src/components.d.ts',
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
    allowedHosts: true,
  },
  build: {
    // 输出到 public 目录（Buffalo 静态资源目录）
    outDir: '../public',
    emptyOutDir: false, // 不清空目录，避免文件锁定问题
    // 代码分割优化
    rollupOptions: {
      output: {
        manualChunks: {
          // Vue 核心库
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          // Element Plus UI 库
          'element-plus': ['element-plus', '@element-plus/icons-vue'],
          // 工具库
          'utils': ['axios', '@vueuse/core'],
        },
      },
    },
    // 启用 gzip 压缩大小报告
    reportCompressedSize: true,
    // chunk 大小警告阈值
    chunkSizeWarningLimit: 1000,
  },
  // 优化依赖预构建
  optimizeDeps: {
    include: ['vue', 'vue-router', 'pinia', 'element-plus', 'axios'],
  },
})
