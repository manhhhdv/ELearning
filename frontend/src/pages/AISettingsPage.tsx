import { useEffect, useState } from 'react'

import { api } from '../api/client'
import { invalidateAIStatus } from '../api/useAIStatus'
import type { AIChannelInput, AIProvider, AISettings } from '../api/types'
import { PageHeader } from '../components/Layout'
import { IconSparkles } from '../components/icons'
import { ErrorAlert, Loading, SuccessAlert } from '../components/ui'

const SOURCE_LABEL: Record<AISettings['source'], string> = {
  database: 'Đang dùng cấu hình lưu trong hệ thống',
  env: 'Đang dùng cấu hình mặc định từ máy chủ (.env)',
  none: 'Chưa cấu hình ở đâu cả — các chức năng AI đang tắt',
}

/** Một dòng trong biểu mẫu: kênh gửi lên máy chủ kèm gợi ý khoá đang lưu. */
type ChannelRow = AIChannelInput & { savedHint: string }

const emptyRow = (provider: AIProvider | undefined): ChannelRow => ({
  label: '',
  provider: provider?.id ?? 'gemini',
  apiKey: '',
  model: provider?.defaultModel ?? '',
  baseUrl: '',
  disabled: false,
  savedHint: '',
})

/**
 * Cấu hình các kênh AI. Chỉ admin vào được vì khoá API là thông tin nhạy cảm.
 *
 * Hệ thống gọi lần lượt theo thứ tự trong danh sách và tự chuyển sang kênh kế
 * tiếp khi kênh đang dùng bị hết hạn mức hoặc lỗi tạm thời, nên nên khai báo ít
 * nhất hai kênh của hai nhà cung cấp khác nhau.
 */
export function AISettingsPage() {
  const [settings, setSettings] = useState<AISettings | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [enabled, setEnabled] = useState(false)
  const [rows, setRows] = useState<ChannelRow[]>([])
  const [saved, setSaved] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [testing, setTesting] = useState<number | null>(null)
  const [testResult, setTestResult] = useState<Record<number, { ok: boolean; message: string }>>({})

  const load = () => {
    setLoading(true)
    api.getAISettings()
      .then((s) => {
        setSettings(s)
        setEnabled(s.enabled && s.source === 'database')
        // Không bao giờ điền lại khoá: máy chủ chỉ gửi về 4 ký tự cuối.
        setRows(s.channels.map((c) => ({
          label: c.label, provider: c.provider, apiKey: '', model: c.model,
          baseUrl: c.baseUrl, disabled: c.disabled, savedHint: c.apiKeyHint,
        })))
        setTestResult({})
        setError(null)
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được cấu hình'))
      .finally(() => setLoading(false))
  }

  useEffect(load, [])

  const providerOf = (id: string): AIProvider | undefined =>
    settings?.providers.find((p) => p.id === id) ?? settings?.providers[0]

  const patch = (index: number, change: Partial<ChannelRow>) =>
    setRows((prev) => prev.map((row, i) => (i === index ? { ...row, ...change } : row)))

  const move = (index: number, delta: number) =>
    setRows((prev) => {
      const next = [...prev]
      const target = index + delta
      if (target < 0 || target >= next.length) return prev
      ;[next[index], next[target]] = [next[target], next[index]]
      return next
    })

  const remove = (index: number) => setRows((prev) => prev.filter((_, i) => i !== index))

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setSaved(null)
    try {
      await api.saveAISettings({
        enabled,
        channels: rows.map(({ savedHint: _hint, ...channel }) => ({
          ...channel,
          label: channel.label.trim(),
          apiKey: channel.apiKey.trim(),
          model: channel.model.trim(),
          baseUrl: channel.baseUrl.trim(),
        })),
      })
      // Các màn hình khác đang nhớ trạng thái cũ — buộc chúng hỏi lại máy chủ.
      invalidateAIStatus()
      setSaved(enabled
        ? 'Đã lưu cấu hình AI. Các chức năng AI dùng ngay danh sách kênh mới, không cần khởi động lại máy chủ.'
        : 'Đã tắt — hệ thống dùng lại cấu hình từ .env (nếu có).')
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được cấu hình')
    } finally {
      setBusy(false)
    }
  }

  /** Gọi thử một kênh bằng thông tin đang gõ, trước khi lưu. */
  const test = async (index: number) => {
    const row = rows[index]
    setTesting(index)
    setError(null)
    try {
      const result = await api.testAISettings({
        index,
        provider: row.provider,
        apiKey: row.apiKey.trim(),
        model: row.model.trim(),
        baseUrl: row.baseUrl.trim(),
        label: row.label.trim(),
      })
      setTestResult((prev) => ({ ...prev, [index]: result }))
    } catch (err) {
      setTestResult((prev) => ({
        ...prev,
        [index]: { ok: false, message: err instanceof Error ? err.message : 'Không kiểm tra được kết nối' },
      }))
    } finally {
      setTesting(null)
    }
  }

  return (
    <>
      <PageHeader
        title="Cấu hình AI"
        subtitle="Các kênh AI dùng cho soạn giáo án, tạo đề, chấm tự luận và trợ lý AI — tự xoay vòng khi một kênh hết hạn mức"
      />

      <div className="page-body">
        {loading ? <Loading /> : !settings ? <ErrorAlert message={error} /> : (
          <div className="card card-pad" style={{ maxWidth: 760 }}>
            <ErrorAlert message={error} />
            <SuccessAlert message={saved} />

            <div className="alert alert-info" style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <IconSparkles />
              <span>
                {SOURCE_LABEL[settings.source]}
                {settings.enabled && <> · {settings.channels.filter((c) => !c.disabled && c.hasApiKey).length} kênh đang bật · kênh đầu tiên dùng model <code>{settings.model}</code></>}
              </span>
            </div>

            <form onSubmit={submit}>
              <label className="checkbox" style={{ marginBottom: 16 }}>
                <input type="checkbox" checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />
                Bật các chức năng AI, quản lý các kênh ngay tại đây
              </label>
              {!enabled && (
                <div className="hint" style={{ marginBottom: 16 }}>
                  Tắt sẽ xoá cấu hình đã lưu trong hệ thống. Nếu máy chủ có khai báo sẵn
                  <code> GEMINI_API_KEY </code> trong file .env, các chức năng AI vẫn chạy bằng khoá đó.
                  Không có ở đâu cả thì mọi nút AI tự ẩn khỏi giao diện.
                </div>
              )}

              <fieldset disabled={!enabled} style={{ border: 'none', padding: 0, margin: 0 }}>
                <div className="hint" style={{ marginBottom: 16 }}>
                  Hệ thống gọi lần lượt từ kênh trên cùng xuống. Kênh nào bị báo hết hạn mức
                  (lỗi 429) hoặc lỗi tạm thời sẽ bị tạm ngưng vài phút và lượt gọi tự chuyển
                  sang kênh kế tiếp. Nên khai báo ít nhất hai kênh của hai nhà cung cấp khác
                  nhau để giờ cao điểm không bị gián đoạn.
                </div>

                {rows.length === 0 && (
                  <div className="hint" style={{ marginBottom: 16 }}>Chưa có kênh nào. Bấm “Thêm kênh” để bắt đầu.</div>
                )}

                {rows.map((row, index) => {
                  const prov = providerOf(row.provider)
                  const result = testResult[index]
                  return (
                    <div
                      key={index}
                      className="card card-pad"
                      style={{ marginBottom: 14, opacity: row.disabled ? 0.6 : 1 }}
                    >
                      <div className="wrap-gap" style={{ justifyContent: 'space-between', marginBottom: 10 }}>
                        <b>Kênh {index + 1}{index === 0 ? ' · ưu tiên cao nhất' : ''}</b>
                        <span className="wrap-gap">
                          <button type="button" className="btn btn-sm" onClick={() => move(index, -1)} disabled={index === 0}>↑</button>
                          <button type="button" className="btn btn-sm" onClick={() => move(index, 1)} disabled={index === rows.length - 1}>↓</button>
                          <button type="button" className="btn btn-sm btn-danger" onClick={() => remove(index)}>Xoá</button>
                        </span>
                      </div>

                      <div className="field">
                        <label htmlFor={`ai-label-${index}`}>Tên gợi nhớ</label>
                        <input
                          id={`ai-label-${index}`} type="text" value={row.label}
                          onChange={(e) => patch(index, { label: e.target.value })}
                          placeholder={`${prov?.label ?? ''} – khoá chính`}
                        />
                        <div className="hint">Chỉ để phân biệt khi có nhiều khoá của cùng một nhà cung cấp; hiện trong nhật ký lỗi.</div>
                      </div>

                      <div className="field">
                        <label htmlFor={`ai-provider-${index}`}>Nhà cung cấp</label>
                        <select
                          id={`ai-provider-${index}`}
                          value={row.provider}
                          onChange={(e) => {
                            // Đổi nhà cung cấp thì model và endpoint cũ không còn hợp lệ.
                            const next = providerOf(e.target.value)
                            patch(index, {
                              provider: e.target.value,
                              model: next?.defaultModel ?? '',
                              baseUrl: '',
                            })
                          }}
                        >
                          {settings.providers.map((p) => (
                            <option key={p.id} value={p.id}>{p.label}</option>
                          ))}
                        </select>
                      </div>

                      <div className="field">
                        <label htmlFor={`ai-key-${index}`}>Khoá API</label>
                        <input
                          id={`ai-key-${index}`} type="password" className="mono" value={row.apiKey}
                          onChange={(e) => patch(index, { apiKey: e.target.value })}
                          placeholder={row.savedHint ? 'Để trống để giữ nguyên khoá đã lưu' : 'Dán khoá API vào đây'}
                          autoComplete="off"
                        />
                        <div className="hint">
                          {row.savedHint
                            ? `Đang dùng khoá ${row.savedHint}. Chỉ nhập vào đây nếu muốn thay bằng khoá mới.`
                            : 'Khoá chỉ được lưu ở máy chủ và không bao giờ gửi xuống trình duyệt.'}
                          {prov?.apiKeyUrl && (
                            <> Lấy khoá tại <a href={prov.apiKeyUrl} target="_blank" rel="noreferrer">{prov.apiKeyUrl}</a>.</>
                          )}
                        </div>
                      </div>

                      <div className="field">
                        <label htmlFor={`ai-model-${index}`}>Model</label>
                        <input
                          id={`ai-model-${index}`} type="text" className="mono" value={row.model}
                          onChange={(e) => patch(index, { model: e.target.value })}
                          list={`ai-model-list-${index}`}
                          placeholder={prov?.defaultModel}
                        />
                        <datalist id={`ai-model-list-${index}`}>
                          {prov?.models.map((m) => <option key={m} value={m} />)}
                        </datalist>
                        <div className="hint">
                          Chọn trong danh sách gợi ý hoặc gõ tên model khác. Model nhanh hợp với việc
                          tạo đề số lượng lớn; model mạnh hơn cho chất lượng tốt hơn nhưng chậm và tốn
                          hạn mức hơn.
                        </div>
                      </div>

                      <div className="field">
                        <label htmlFor={`ai-base-${index}`}>
                          Địa chỉ API {prov?.requiresBaseUrl ? '(bắt buộc)' : '(tuỳ chọn)'}
                        </label>
                        <input
                          id={`ai-base-${index}`} type="text" className="mono" value={row.baseUrl}
                          onChange={(e) => patch(index, { baseUrl: e.target.value })}
                          placeholder={prov?.endpoint || 'https://…/v1'}
                        />
                        <div className="hint">
                          Để trống sẽ dùng địa chỉ chính thức của nhà cung cấp. Chỉ nhập khi đi qua
                          proxy nội bộ hoặc dùng dịch vụ tương thích OpenAI khác.
                        </div>
                      </div>

                      <label className="checkbox" style={{ marginBottom: 12 }}>
                        <input
                          type="checkbox" checked={row.disabled}
                          onChange={(e) => patch(index, { disabled: e.target.checked })}
                        />
                        Tạm ngưng kênh này (giữ lại cấu hình nhưng không gọi tới)
                      </label>

                      <div className="wrap-gap">
                        <button
                          type="button" className="btn btn-sm"
                          onClick={() => test(index)} disabled={testing !== null || busy}
                        >
                          {testing === index ? 'Đang gọi thử…' : 'Kiểm tra kết nối'}
                        </button>
                        {result && (
                          <span className={result.ok ? 'badge badge-success' : 'badge badge-danger'}>
                            {result.message}
                          </span>
                        )}
                      </div>
                    </div>
                  )
                })}

                <button
                  type="button" className="btn btn-sm" style={{ marginBottom: 16 }}
                  onClick={() => setRows((prev) => [...prev, emptyRow(settings.providers[0])])}
                >
                  + Thêm kênh
                </button>
              </fieldset>

              <div>
                <button className="btn btn-primary" disabled={busy}>
                  {busy ? 'Đang lưu…' : 'Lưu cấu hình'}
                </button>
              </div>
            </form>

            <div className="hint" style={{ marginTop: 18 }}>
              Nội dung do AI tạo ra là <b>bản nháp</b>. Hệ thống tự kiểm tra cấu trúc và số lượng câu
              hỏi, nhưng giáo viên vẫn phải duyệt trước khi giao cho học sinh, và điểm tự luận do AI
              gợi ý chỉ được ghi nhận sau khi giáo viên xác nhận.
            </div>
          </div>
        )}
      </div>
    </>
  )
}
