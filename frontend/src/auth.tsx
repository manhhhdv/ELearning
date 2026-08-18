import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'

import { ApiError, api, getToken, setToken } from './api/client'
import type { User } from './api/types'

interface AuthState {
  user: User | null
  loading: boolean
  signIn: (email: string, password: string) => Promise<User>
  signInWithToken: (token: string) => Promise<User>
  signOut: () => void
  refresh: () => Promise<void>
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  // Khôi phục phiên từ token đã lưu khi tải lại trang.
  useEffect(() => {
    if (!getToken()) {
      setLoading(false)
      return
    }
    api.me()
      .then(setUser)
      .catch((err) => {
        if (err instanceof ApiError && (err.status === 401 || err.status === 403)) setToken(null)
      })
      .finally(() => setLoading(false))
  }, [])

  const signIn = useCallback(async (email: string, password: string) => {
    const res = await api.login(email, password)
    setToken(res.token)
    setUser(res.user)
    return res.user
  }, [])

  const signInWithToken = useCallback(async (token: string) => {
    setToken(token)
    const me = await api.me()
    setUser(me)
    return me
  }, [])

  const signOut = useCallback(() => {
    setToken(null)
    setUser(null)
  }, [])

  const refresh = useCallback(async () => {
    setUser(await api.me())
  }, [])

  const value = useMemo(
    () => ({ user, loading, signIn, signInWithToken, signOut, refresh }),
    [user, loading, signIn, signInWithToken, signOut, refresh],
  )
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth phải nằm trong AuthProvider')
  return ctx
}

/** Người quản lý nội dung: quản trị viên và giảng viên — thấy nút sửa/xoá/tạo mới. */
export function canManageContent(user: User | null) {
  return user?.role === 'admin' || user?.role === 'trainer'
}

/**
 * Vào được khu vực quản lý: quản lý nội dung cộng thêm vai trò Giám sát —
 * họ xem được mọi thứ nhưng giao diện phải ẩn hết nút sửa/xoá/tạo mới.
 * Dùng canManageContent() cho các thao tác ghi để phân biệt hai nhóm này.
 */
export function canAccessAdminArea(user: User | null) {
  return canManageContent(user) || user?.role === 'supervisor'
}

/** Chỉ Giám sát mới cần chế độ đọc — admin/trainer luôn sửa được nội dung họ đang xem. */
export function isReadOnlyViewer(user: User | null) {
  return user?.role === 'supervisor'
}
