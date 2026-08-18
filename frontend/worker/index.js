// Worker phục vụ frontend (static assets trong dist/) và proxy /api/* sang backend.
//
// `run_worker_first: ["/api/*"]` trong wrangler.jsonc đảm bảo mọi request /api/* luôn chạy qua
// đây trước, không bị assets binding trả 404 do không khớp file tĩnh nào. Các request còn lại rơi
// xuống `env.ASSETS.fetch()`, gồm cả fallback SPA (not_found_handling: single-page-application)
// nên tải lại trang ở một route con của React Router (vd /hoc/:maKhoa) vẫn ra đúng index.html.
export default {
  async fetch(request, env) {
    const url = new URL(request.url)

    if (url.pathname.startsWith('/api/')) {
      const target = new URL(url.pathname + url.search, env.BACKEND_ORIGIN)
      return fetch(new Request(target, request))
    }

    return env.ASSETS.fetch(request)
  },
}
