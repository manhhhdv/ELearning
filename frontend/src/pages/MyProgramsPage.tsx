import { useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

import { api } from '../api/client'
import type { Program, Submission } from '../api/types'
import { useAuth } from '../auth'
import { CourseTile } from '../components/CourseTile'
import { IconEnrolled, IconFinished, IconSubmitted } from '../components/icons'
import { Loading } from '../components/ui'

export function MyProgramsPage() {
  const { user } = useAuth()
  const [params, setParams] = useSearchParams()
  const [programs, setPrograms] = useState<Program[]>([])
  const [submissions, setSubmissions] = useState<Submission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const query = params.get('q') ?? ''

  useEffect(() => {
    Promise.all([api.myPrograms(), api.mySubmissions()])
      .then(([ps, subs]) => { setPrograms(ps); setSubmissions(subs); setError(null) })
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được dữ liệu'))
      .finally(() => setLoading(false))
  }, [])

  // Lọc tại chỗ: danh sách khoá của một học viên vốn ngắn, không cần gọi lại máy chủ.
  const shown = useMemo(() => {
    const term = query.trim().toLowerCase()
    if (!term) return programs
    return programs.filter((p) =>
      p.title.toLowerCase().includes(term)
      || p.code.toLowerCase().includes(term)
      || p.description.toLowerCase().includes(term))
  }, [programs, query])

  if (loading) return <Loading />

  // Khoá học coi là hoàn thành khi đã học hết số bài học của khoá đó.
  const finished = programs.filter((p) => p.lessonCount > 0 && p.completedLessonCount >= p.lessonCount).length

  return (
    <div className="dash">
      <h1 className="dash-hello">Chào mừng trở lại, {user?.fullName || user?.email}!</h1>

      {error && <div className="alert alert-error">{error}</div>}

      <div className="stat-row">
        <Stat icon={<IconEnrolled />} tone="blue" label="Đã ghi danh" value={`${programs.length} khoá học`} />
        <Stat icon={<IconFinished />} tone="green" label="Đã hoàn thành" value={`${finished} khoá học`} />
        <Stat icon={<IconSubmitted />} tone="amber" label="Số bài đã nộp" value={`${submissions.length} bài`} />
      </div>

      <div className="section-head">
        <h2>Khoá học của tôi</h2>
        {query && (
          <button className="btn btn-sm" onClick={() => setParams({})}>
            Xoá bộ lọc “{query}”
          </button>
        )}
      </div>

      <div className="tray">
        {shown.length === 0 ? (
          <div className="empty" style={{ padding: '32px 20px' }}>
            <h3>{query ? 'Không tìm thấy khoá học phù hợp' : 'Bạn chưa được ghi danh vào khoá học nào'}</h3>
            <p>
              {query
                ? 'Thử từ khoá khác, hoặc xoá bộ lọc để xem toàn bộ khoá học của bạn.'
                : 'Liên hệ quản trị viên hoặc giảng viên để được thêm vào chương trình đào tạo.'}
            </p>
          </div>
        ) : (
          <div className="tray-grid">
            {shown.map((p) => <CourseTile key={p.id} program={p} />)}
          </div>
        )}
      </div>
    </div>
  )
}

function Stat({
  icon, tone, label, value,
}: {
  icon: React.ReactNode
  tone: 'blue' | 'green' | 'amber'
  label: string
  value: string
}) {
  return (
    <div className="stat">
      <div className={`stat-icon ${tone}`} aria-hidden>{icon}</div>
      <div>
        <div className="label">{label}</div>
        <div className="value">{value}</div>
      </div>
    </div>
  )
}
