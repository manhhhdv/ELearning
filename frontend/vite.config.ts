import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      // Bộ kết xuất nội dung (markdown/LaTeX) nặng và gần như không đổi giữa các lần
      // phát hành, tách riêng để trình duyệt còn dùng lại bản đã tải.
      output: { manualChunks: { richtext: ['marked', 'katex', 'dompurify'] } },
    },
  },
  server: {
    port: 3006,
    // Không cho Vite tự nhảy cổng khi 3006 bận: đổi cổng ngầm sẽ làm hỏng
    // CORS và redirect OAuth vì backend cấu hình cứng FRONTEND_URL.
    strictPort: true,
    proxy: {
      // Gọi API qua cùng origin để tránh vướng CORS lúc phát triển.
      '/api': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
    },
  },
})
