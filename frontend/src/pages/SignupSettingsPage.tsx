import { useEffect, useState } from 'react'

import { api } from '../api/client'
import { PageHeader } from '../components/Layout'
import { ErrorAlert, Loading, SuccessAlert } from '../components/ui'

/**
 * Bật/tắt việc người dùng tự đăng ký tài khoản ở trang đăng nhập.
 *
 * Tài khoản tự đăng ký luôn có vai trò Học sinh và dùng được ngay. Quyền soạn
 * nội dung hay quản trị vẫn phải do admin nâng cấp ở mục Người dùng.
 */
export function SignupSettingsPage() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const [enabled, setEnabled] = useState(false)
  const [allowedDomains, setAllowedDomains] = useState('')

  const load = () => {
    setLoading(true)
    api.getSignupSettings()
      .then((s) => {
        setEnabled(s.enabled)
        setAllowedDomains(s.allowedDomains)
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
      await api.saveSignupSettings({ enabled, allowedDomains })
      setSaved(enabled
        ? 'Đã bật đăng ký. Trang đăng nhập hiện nút "Đăng ký tài khoản mới".'
        : 'Đã tắt đăng ký. Chỉ quản trị viên cấp được tài khoản.')
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được cấu hình')
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <PageHeader
        title="Đăng ký tài khoản"
        subtitle="Cho phép người dùng tự tạo tài khoản ở trang đăng nhập"
      />

      <div className="page-body">
        {loading ? <Loading /> : (
          <div className="card card-pad" style={{ maxWidth: 640 }}>
            <ErrorAlert message={error} />
            <SuccessAlert message={saved} />

            <form onSubmit={submit}>
              <label className="checkbox" style={{ marginBottom: 16 }}>
                <input type="checkbox" checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />
                Cho phép tự đăng ký tài khoản
              </label>

              <div className={enabled ? '' : 'hint'} style={{ marginBottom: 16 }}>
                {enabled ? (
                  <div className="alert alert-info">
                    Tài khoản tự đăng ký có vai trò <b>Học sinh</b> và dùng được ngay. Muốn cấp quyền
                    Giáo viên hoặc Quản trị viên, hãy đổi vai trò ở mục <b>Người dùng</b>.
                  </div>
                ) : (
                  'Đang tắt: trang đăng nhập chỉ hiện hướng dẫn liên hệ quản trị viên.'
                )}
              </div>

              <fieldset disabled={!enabled} style={{ border: 'none', padding: 0, margin: 0 }}>
                <div className="field">
                  <label htmlFor="su-domains">Giới hạn theo domain email</label>
                  <input
                    id="su-domains" type="text" value={allowedDomains}
                    onChange={(e) => setAllowedDomains(e.target.value)}
                    placeholder="truonghoc.edu.vn, hs.truonghoc.edu.vn"
                  />
                  <div className="hint">
                    Cách nhau bằng dấu phẩy. <b>Để trống nghĩa là bất kỳ ai biết địa chỉ trang web
                    cũng tạo được tài khoản</b> — nên đặt giới hạn domain của trường.
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
