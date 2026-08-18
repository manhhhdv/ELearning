import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

import { api } from '../api/client'
import type { Program } from '../api/types'
import { PageHeader } from '../components/Layout'
import { courseCover } from '../components/CourseTile'
import { IconCheck, IconDoc, IconPlus } from '../components/icons'
import { Loading } from '../components/ui'

export function CatalogPage() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const [programs, setPrograms] = useState<Program[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [busyId, setBusyId] = useState<string | null>(null)

  const query = params.get('q') ?? ''

  const load = useCallback(async () => {
    try {
      setPrograms(await api.catalog(query))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được danh mục khoá học')
    } finally {
      setLoading(false)
    }
  }, [query])

  useEffect(() => { load() }, [load])

  const enroll = async (p: Program) => {
    setBusyId(p.id)
    setError(null)
    try {
      await api.selfEnroll(p.id)
      navigate(`/hoc/${p.slug}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không ghi danh được')
      setBusyId(null)
    }
  }

  return (
    <>
      <PageHeader
        title="Khám phá khoá học"
        subtitle="Các khoá học mở, bạn có thể tự đăng ký mà không cần chờ quản trị viên"
      />

      <div className="page-body">
        {error && <div className="alert alert-error">{error}</div>}

        {loading ? <Loading /> : programs.length === 0 ? (
          <div className="card">
            <div className="empty">
              <h3>Hiện chưa có khoá học nào mở đăng ký</h3>
              <p>Khi có khoá mới mở, khoá đó sẽ xuất hiện tại đây.</p>
            </div>
          </div>
        ) : (
          <div className="tray-grid">
            {programs.map((p) => (
              <div className="tile" key={p.id} style={{ cursor: 'default' }}>
                <div className="tile-cover" style={courseCover(p)}>
                  {!p.coverUrl && p.code}
                </div>
                <div className="tile-body">
                  <h3>{p.title}</h3>
                  {p.description && (
                    <p style={{ fontSize: 14, color: 'var(--text-muted)', marginTop: 6 }}>{p.description}</p>
                  )}
                </div>
                <div className="tile-foot">
                  <IconDoc size={16} />
                  {p.lessonCount} bài học
                  {p.assignmentCount > 0 && ` · ${p.assignmentCount} bài tập`}
                </div>
                <div style={{ padding: '0 18px 18px' }}>
                  {p.enrolled ? (
                    <button className="btn btn-block" onClick={() => navigate(`/hoc/${p.slug}`)}>
                      <IconCheck size={15} /> Đã ghi danh — Vào học
                    </button>
                  ) : (
                    <button
                      className="btn btn-primary btn-block"
                      disabled={busyId === p.id}
                      onClick={() => enroll(p)}
                    >
                      <IconPlus size={15} /> {busyId === p.id ? 'Đang ghi danh…' : 'Đăng ký học'}
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  )
}
