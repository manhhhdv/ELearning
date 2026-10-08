import { useCallback, useEffect, useState } from 'react'

import { api } from '../api/client'
import type { LessonPlanInput } from '../api/client'
import { useAIStatus } from '../api/useAIStatus'
import type { LessonPlan, LessonPlanContent } from '../api/types'
import { PageHeader } from '../components/Layout'
import { IconDownload, IconSparkles, IconTrash } from '../components/icons'
import { ErrorAlert, Loading, formatDate } from '../components/ui'

const emptyForm = (): LessonPlanInput => ({
  subject: '', grade: '', topic: '', content: '', objectives: '', durationMinutes: 45,
})

/**
 * Màn hình tạo giáo án bằng AI (mục 3.4.1): giáo viên nhập thông tin bài học,
 * AI trả về bản nháp theo từng hoạt động dạy học để xem, chỉnh sửa và lưu lại.
 */
export function LessonPlanPage() {
  const [form, setForm] = useState<LessonPlanInput>(emptyForm())
  const [plan, setPlan] = useState<LessonPlanContent | null>(null)
  const [savedId, setSavedId] = useState<string | null>(null)
  const [plans, setPlans] = useState<LessonPlan[]>([])
  const aiStatus = useAIStatus()
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  // Mã giáo án đang được xuất ra Word ('draft' = bản nháp trên màn hình).
  const [exporting, setExporting] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const set = <K extends keyof LessonPlanInput>(key: K, value: LessonPlanInput[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }))

  const loadPlans = useCallback(async () => {
    try {
      setPlans(await api.listLessonPlans())
    } catch {
      setPlans([])
    }
  }, [])

  useEffect(() => {
    loadPlans().finally(() => setLoading(false))
  }, [loadPlans])

  const generate = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setNotice(null)
    try {
      const res = await api.generateLessonPlan(form)
      setPlan(res.plan)
      // Bản nháp mới chưa gắn với giáo án đã lưu nào.
      setSavedId(null)
      if (res.attempts > 1) {
        setNotice(`Bản nháp đầu chưa đạt yêu cầu, AI đã soạn lại ${res.attempts - 1} lần.`)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tạo được giáo án')
    } finally {
      setBusy(false)
    }
  }

  const save = async () => {
    if (!plan) return
    setBusy(true)
    setError(null)
    try {
      if (savedId) {
        await api.updateLessonPlan(savedId, plan.title, plan)
        setNotice('Đã cập nhật giáo án.')
      } else {
        const created = await api.saveLessonPlan({
          subject: form.subject,
          grade: form.grade,
          topic: form.topic,
          objectives: form.objectives,
          durationMinutes: form.durationMinutes,
          title: plan.title,
          content: plan,
        })
        setSavedId(created.id)
        setNotice('Đã lưu giáo án vào danh sách của bạn.')
      }
      await loadPlans()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được giáo án')
    } finally {
      setBusy(false)
    }
  }

  /**
   * Xuất giáo án ra Word. Bản nháp gửi thẳng nội dung đang hiện trên màn hình
   * nên những chỉnh sửa chưa lưu vẫn có trong tệp.
   */
  const exportWord = async (saved?: LessonPlan) => {
    if (!saved && !plan) return
    setExporting(saved ? saved.id : 'draft')
    setError(null)
    try {
      if (saved) {
        await api.exportLessonPlan(saved.id)
      } else if (plan) {
        await api.exportLessonPlanDraft({
          subject: form.subject, grade: form.grade, durationMinutes: form.durationMinutes, content: plan,
        })
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không xuất được tệp Word')
    } finally {
      setExporting(null)
    }
  }

  const openSaved = async (p: LessonPlan) => {
    setPlan(p.content)
    setSavedId(p.id)
    setNotice(null)
    setError(null)
    setForm({
      subject: p.subject, grade: p.grade, topic: p.topic,
      content: '', objectives: p.objectives, durationMinutes: p.durationMinutes,
    })
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const remove = async (p: LessonPlan) => {
    if (!confirm(`Xoá giáo án "${p.title}"?`)) return
    try {
      await api.deleteLessonPlan(p.id)
      if (savedId === p.id) { setSavedId(null); setPlan(null) }
      await loadPlans()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không xoá được giáo án')
    }
  }

  /** Sửa một trường của bản nháp ngay trên màn hình trước khi lưu. */
  const editPlan = (patch: Partial<LessonPlanContent>) =>
    setPlan((prev) => (prev ? { ...prev, ...patch } : prev))

  const editActivity = (index: number, patch: Partial<LessonPlanContent['activities'][number]>) =>
    setPlan((prev) => prev
      ? { ...prev, activities: prev.activities.map((a, i) => (i === index ? { ...a, ...patch } : a)) }
      : prev)

  if (loading || aiStatus === null) return <Loading label="Đang tải…" />

  if (!aiStatus.enabled) {
    return (
      <>
        <PageHeader title="Soạn giáo án bằng AI" />
        <div className="callout warn">
          Chức năng AI chưa được bật. Quản trị viên cần chọn nhà cung cấp và nhập khoá API ở mục
          <b> Quản lý → Cấu hình AI</b>.
        </div>
      </>
    )
  }

  const totalMinutes = plan?.activities.reduce((sum, a) => sum + (a.durationMinutes || 0), 0) ?? 0

  return (
    <>
      <PageHeader
        title="Soạn giáo án bằng AI"
        subtitle="Nhập thông tin bài học, AI soạn bản nháp theo từng hoạt động dạy học để bạn chỉnh sửa."
      />

      <div className="page-body">
        <ErrorAlert message={error} />
        {notice && <div className="alert alert-success">{notice}</div>}

        <div className="panel">
          <form onSubmit={generate}>
            <div className="row">
              <div className="field">
                <label>Lớp học</label>
                <input value={form.subject} onChange={(e) => set('subject', e.target.value)} placeholder="VD: Vật lí" />
              </div>
              <div className="field">
                <label>Lớp</label>
                <input value={form.grade} onChange={(e) => set('grade', e.target.value)} placeholder="VD: lớp 11" />
              </div>
              <div className="field" style={{ maxWidth: 150 }}>
                <label>Thời lượng (phút)</label>
                <input
                  type="number" min={15} max={180} step={5}
                  value={form.durationMinutes}
                  onChange={(e) => set('durationMinutes', Number(e.target.value) || 45)}
                />
              </div>
            </div>

            <div className="field">
              <label>Chủ đề bài học<span style={{ color: 'var(--c-danger, #dc2626)' }}> *</span></label>
              <input
                value={form.topic}
                onChange={(e) => set('topic', e.target.value)}
                placeholder="VD: Định luật Ôm cho đoạn mạch"
                required
              />
            </div>

            <div className="field">
              <label>Nội dung trọng tâm</label>
              <textarea
                value={form.content}
                onChange={(e) => set('content', e.target.value)}
                placeholder="Những kiến thức chính cần truyền đạt…"
                rows={3}
              />
            </div>

            <div className="field">
              <label>Yêu cầu cần đạt</label>
              <textarea
                value={form.objectives}
                onChange={(e) => set('objectives', e.target.value)}
                placeholder="Sau bài học, học sinh có thể…"
                rows={3}
              />
            </div>

            <button className="btn btn-primary" disabled={busy || !form.topic.trim()}>
              <IconSparkles /> {busy ? 'AI đang soạn…' : 'Tạo giáo án'}
            </button>
          </form>
        </div>

        {plan && (
          <div className="panel">
            <div className="spread" style={{ marginBottom: 14 }}>
              <h3 style={{ margin: 0 }}>Bản nháp giáo án</h3>
              <div className="wrap-gap">
                <span className="tiny muted">
                  {plan.activities.length} hoạt động · {totalMinutes} phút
                </span>
                <button
                  className="btn btn-sm" onClick={() => exportWord()} disabled={exporting !== null}
                  title="Tải giáo án đang hiện trên màn hình về dạng Word (.docx)"
                >
                  <IconDownload /> {exporting === 'draft' ? 'Đang xuất…' : 'Xuất Word'}
                </button>
                <button className="btn btn-primary btn-sm" onClick={save} disabled={busy}>
                  {busy ? 'Đang lưu…' : savedId ? 'Cập nhật giáo án' : 'Lưu giáo án'}
                </button>
              </div>
            </div>

            <p className="hint" style={{ marginTop: 0 }}>
              Đây là bản nháp do AI soạn. Bạn hãy đọc lại, chỉnh sửa cho phù hợp với lớp mình rồi lưu.
            </p>

            <div className="field">
              <label>Tên bài dạy</label>
              <input value={plan.title} onChange={(e) => editPlan({ title: e.target.value })} />
            </div>

            <div className="field">
              <label>Yêu cầu cần đạt</label>
              <textarea
                rows={Math.max(3, plan.objectives.length)}
                value={plan.objectives.join('\n')}
                onChange={(e) => editPlan({ objectives: e.target.value.split('\n') })}
              />
              <div className="hint">Mỗi dòng là một yêu cầu cần đạt.</div>
            </div>

            {plan.materials.length > 0 && (
              <div className="field">
                <label>Thiết bị và học liệu</label>
                <textarea
                  rows={Math.max(2, plan.materials.length)}
                  value={plan.materials.join('\n')}
                  onChange={(e) => editPlan({ materials: e.target.value.split('\n') })}
                />
              </div>
            )}

            {plan.activities.map((a, i) => (
              <div className="question-card" key={i}>
                <div className="question-card-head">
                  <span className="badge badge-primary">{a.name}</span>
                  <span className="tiny muted grow">{a.durationMinutes} phút</span>
                </div>
                {a.objective && <div className="hint">Mục tiêu: {a.objective}</div>}
                <div className="field">
                  <label>Hoạt động của giáo viên</label>
                  <textarea
                    rows={3}
                    value={a.teacherActions}
                    onChange={(e) => editActivity(i, { teacherActions: e.target.value })}
                  />
                </div>
                <div className="field">
                  <label>Hoạt động của học sinh</label>
                  <textarea
                    rows={3}
                    value={a.studentActions}
                    onChange={(e) => editActivity(i, { studentActions: e.target.value })}
                  />
                </div>
                {a.output && <div className="hint">Sản phẩm cần đạt: {a.output}</div>}
              </div>
            ))}

            <div className="field">
              <label>Bài tập về nhà</label>
              <textarea rows={2} value={plan.homework} onChange={(e) => editPlan({ homework: e.target.value })} />
            </div>
            <div className="field">
              <label>Lưu ý khi tổ chức dạy học</label>
              <textarea rows={2} value={plan.notes} onChange={(e) => editPlan({ notes: e.target.value })} />
            </div>
          </div>
        )}

        <div className="panel">
          <h3 style={{ marginTop: 0 }}>Giáo án đã lưu</h3>
          {plans.length === 0 ? (
            <p className="muted tiny">Bạn chưa lưu giáo án nào.</p>
          ) : (
            <table className="ltable">
              <thead>
                <tr><th>Tên bài dạy</th><th>Môn · Lớp</th><th>Tạo lúc</th><th /></tr>
              </thead>
              <tbody>
                {plans.map((p) => (
                  <tr key={p.id}>
                    <td>{p.title}</td>
                    <td className="tiny muted">{[p.subject, p.grade].filter(Boolean).join(' · ') || '—'}</td>
                    <td className="tiny muted">{formatDate(p.createdAt)}</td>
                    <td style={{ textAlign: 'right' }}>
                      <button className="btn btn-sm" onClick={() => openSaved(p)}>Mở</button>{' '}
                      <button
                        className="btn btn-sm" onClick={() => exportWord(p)} disabled={exporting !== null}
                        title="Tải giáo án dạng Word (.docx)"
                      >
                        <IconDownload /> {exporting === p.id ? 'Đang xuất…' : 'Word'}
                      </button>{' '}
                      <button className="btn btn-ghost btn-sm" onClick={() => remove(p)} title="Xoá">
                        <IconTrash />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </>
  )
}
