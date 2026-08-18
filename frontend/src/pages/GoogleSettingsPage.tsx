import { useEffect, useState } from 'react'

import { api } from '../api/client'
import type { GoogleSettings } from '../api/types'
import { PageHeader } from '../components/Layout'
import { IconGoogle } from '../components/icons'
import { ErrorAlert, Loading, SuccessAlert } from '../components/ui'

const SOURCE_LABEL: Record<GoogleSettings['source'], string> = {
  database: 'Đang dùng cấu hình lưu trong hệ thống',
  env: 'Đang dùng cấu hình mặc định từ máy chủ (.env)',
  none: 'Chưa cấu hình ở đâu cả',
}

export function GoogleSettingsPage() {
  const [settings, setSettings] = useState<GoogleSettings | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [enabled, setEnabled] = useState(false)
  const [clientId, setClientId] = useState('')
  const [clientSecret, setClientSecret] = useState('')
  const [allowedDomains, setAllowedDomains] = useState('')
  const [autoProvisionRole, setAutoProvisionRole] = useState('')
  const [saved, setSaved] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const load = () => {
    setLoading(true)
    api.getGoogleSettings()
      .then((s) => {
        setSettings(s)
        setEnabled(s.enabled)
        setClientId(s.clientId)
        setClientSecret('')
        setAllowedDomains(s.allowedDomains)
        setAutoProvisionRole(s.autoProvisionRole)
        setError(null)
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được cấu hình'))
      .finally(() => setLoading(false))
  }

  useEffect(load, [])

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setSaved(null)
    try {
      await api.saveGoogleSettings({
        enabled,
        clientId: clientId.trim(),
        clientSecret: clientSecret.trim(),
        allowedDomains,
        autoProvisionRole,
      })
      setSaved(enabled ? 'Đã lưu cấu hình đăng nhập Google' : 'Đã tắt — hệ thống dùng lại cấu hình từ .env (nếu có)')
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được cấu hình')
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <PageHeader title="Đăng nhập Google" subtitle="Cho phép học viên và giảng viên đăng nhập bằng tài khoản Google" />

      <div className="page-body">
        {loading ? <Loading /> : !settings ? <ErrorAlert message={error} /> : (
          <div className="card card-pad" style={{ maxWidth: 640 }}>
            <ErrorAlert message={error} />
            <SuccessAlert message={saved} />

            <div className="alert alert-info" style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <IconGoogle />
              <span>{SOURCE_LABEL[settings.source]}</span>
            </div>

            <form onSubmit={submit}>
              <label className="checkbox" style={{ marginBottom: 16 }}>
                <input type="checkbox" checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />
                Bật đăng nhập Google, quản lý cấu hình ngay tại đây
              </label>
              {!enabled && (
                <div className="hint" style={{ marginBottom: 16 }}>
                  Tắt sẽ xoá cấu hình đã lưu trong hệ thống. Nếu máy chủ có khai báo sẵn trong file
                  .env, đăng nhập Google vẫn hoạt động bằng cấu hình đó.
                </div>
              )}

              <fieldset disabled={!enabled} style={{ border: 'none', padding: 0, margin: 0 }}>
                <div className="field">
                  <label htmlFor="g-id">Client ID</label>
                  <input
                    id="g-id" type="text" className="mono" value={clientId}
                    onChange={(e) => setClientId(e.target.value)}
                    placeholder="xxxxxxxxxx.apps.googleusercontent.com"
                  />
                </div>

                <div className="field">
                  <label htmlFor="g-secret">Client Secret</label>
                  <input
                    id="g-secret" type="password" className="mono" value={clientSecret}
                    onChange={(e) => setClientSecret(e.target.value)}
                    placeholder={settings.hasSecret ? 'Để trống để giữ nguyên secret đã lưu' : 'GOCSPX-…'}
                  />
                  {settings.hasSecret && (
                    <div className="hint">Đã có secret được lưu. Chỉ nhập vào đây nếu muốn thay bằng secret mới.</div>
                  )}
                </div>

                <div className="field">
                  <label htmlFor="g-redirect">Authorized redirect URI</label>
                  <input id="g-redirect" type="text" className="mono" value={settings.redirectUrl} readOnly />
                  <div className="hint">
                    Khai báo đúng địa chỉ này trong Google Cloud Console → APIs &amp; Services →
                    Credentials → OAuth Client. Đây là địa chỉ cố định của máy chủ, không sửa được ở đây.
                  </div>
                </div>

                <div className="field">
                  <label htmlFor="g-domains">Giới hạn theo domain email</label>
                  <input
                    id="g-domains" type="text" value={allowedDomains}
                    onChange={(e) => setAllowedDomains(e.target.value)}
                    placeholder="congty.vn, chinhanh.congty.vn"
                  />
                  <div className="hint">Cách nhau bằng dấu phẩy. Để trống nghĩa là không giới hạn — mọi email Google đều đăng nhập được.</div>
                </div>

                <div className="field">
                  <label htmlFor="g-auto">Tự tạo tài khoản cho email lạ</label>
                  <select id="g-auto" value={autoProvisionRole} onChange={(e) => setAutoProvisionRole(e.target.value)}>
                    <option value="">Không — chỉ tài khoản admin đã cấp sẵn mới đăng nhập được</option>
                    <option value="student">Có — tạo với vai trò Học viên</option>
                    <option value="trainer">Có — tạo với vai trò Giảng viên</option>
                  </select>
                  <div className="hint">
                    Áp dụng khi một email đăng nhập Google lần đầu mà chưa có trong hệ thống.
                  </div>
                </div>
              </fieldset>

              <button className="btn btn-primary" disabled={busy}>
                {busy ? 'Đang lưu…' : 'Lưu cấu hình'}
              </button>
            </form>
          </div>
        )}
      </div>
    </>
  )
}
