import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import App from './App'
import { AuthProvider } from './auth'

// Source Sans 3 được đóng gói cùng ứng dụng nên không cần gọi ra ngoài Internet.
import '@fontsource/source-sans-3/400.css'
import '@fontsource/source-sans-3/500.css'
import '@fontsource/source-sans-3/600.css'
import '@fontsource/source-sans-3/700.css'
// KaTeX đi kèm ứng dụng để công thức hiển thị được cả khi không có Internet.
import 'katex/dist/katex.min.css'
import './styles/app.css'
import './styles/learner.css'
import './styles/rich.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <App />
      </AuthProvider>
    </BrowserRouter>
  </StrictMode>,
)
