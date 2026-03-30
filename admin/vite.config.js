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
      // FontAwesome 静态文件代理到后端
      '/fontawesome': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
    allowedHosts: true,
  },
  build: {
    // 输出到 dist 目录
    outDir: 'dist',
    emptyOutDir: true,
    // 代码分割优化
    rollupOptions: {
      output: {
        // 使用默认的分块策略，避免循环依赖问题
        manualChunks: undefined,
      },
    },
    // 启用 gzip 压缩大小报告
    reportCompressedSize: true,
    // chunk 大小警告阈值（提高到 2000kB）
    chunkSizeWarningLimit: 2000,
  },
  // 优化依赖预构建
  optimizeDeps: {
    include: ['vue', 'vue-router', 'pinia', 'element-plus', 'axios'],
  },
})
