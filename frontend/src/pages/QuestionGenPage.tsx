import { useEffect, useState } from 'react'

import { api } from '../api/client'
import type { Program, TreeNode } from '../api/types'
import { PageHeader } from '../components/Layout'
import { ExportQuestionsButtons } from '../components/ExportButtons'
import { GenerateQuestionsModal } from '../components/GenerateQuestionsModal'
import { IconSparkles } from '../components/icons'
import { ErrorAlert, Loading } from '../components/ui'

const flatten = (nodes: TreeNode[]): TreeNode[] => nodes.flatMap((n) => [n, ...flatten(n.children)])

/**
 * Màn hình ra đề bằng AI đứng riêng (mục 3.4.2): giáo viên chọn lớp học và bài
 * tập đích trước, sau đó dùng lại chính khối tạo câu hỏi AI vẫn dùng trong
 * trình soạn bài tập — tránh phải mở bài tập ra mới ra đề được.
 */
export function QuestionGenPage() {
  const [programs, setPrograms] = useState<Program[]>([])
  const [programId, setProgramId] = useState('')
  const [assignments, setAssignments] = useState<TreeNode[]>([])
  const [nodeId, setNodeId] = useState('')
  const [loading, setLoading] = useState(true)
  const [loadingTree, setLoadingTree] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)
  const [open, setOpen] = useState(false)
  // Bài tập vừa được lưu câu hỏi trong phiên này — cây nội dung đã tải chưa có số câu mới.
  const [imported, setImported] = useState<Set<string>>(new Set())

  useEffect(() => {
    api.listPrograms()
      .then((list) => {
        setPrograms(list)
        setProgramId(list[0]?.id || '')
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được lớp học'))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (!programId) { setAssignments([]); setNodeId(''); return }
    let active = true
    setLoadingTree(true)
    api.getTree(programId)
      .then((tree) => {
        if (!active) return
        const items = flatten(tree).filter((n) => n.kind === 'assignment')
        setAssignments(items)
        setNodeId(items[0]?.id || '')
      })
      .catch((err) => { if (active) setError(err instanceof Error ? err.message : 'Không tải được cây nội dung') })
      .finally(() => { if (active) setLoadingTree(false) })
    return () => { active = false }
  }, [programId])

  const selectedAssignment = assignments.find((n) => n.id === nodeId)
  const hasQuestions = !!selectedAssignment
    && ((selectedAssignment.assignment?.questionCount ?? 0) > 0 || imported.has(selectedAssignment.id))

  if (loading) return <Loading label="Đang tải…" />

  return (
    <>
      <PageHeader
        title="Tạo đề bằng AI"
        subtitle="Chọn lớp học và bài tập đích, AI sẽ soạn câu hỏi để bạn duyệt rồi lưu thẳng vào bài tập đó."
      />

      <div className="page-body">
        <ErrorAlert message={error} />
        {notice && <div className="alert alert-success">{notice}</div>}

        <div className="panel">
          {programs.length === 0 ? (
            <p className="muted tiny">Bạn chưa phụ trách lớp học nào để ra đề.</p>
          ) : (
            <>
              <div className="row">
                <div className="field">
                  <label>Lớp học</label>
                  <select value={programId} onChange={(e) => setProgramId(e.target.value)}>
                    {programs.map((p) => <option key={p.id} value={p.id}>{p.title}</option>)}
                  </select>
                </div>
                <div className="field grow">
                  <label>Bài tập đích</label>
                  <select value={nodeId} onChange={(e) => setNodeId(e.target.value)} disabled={loadingTree || assignments.length === 0}>
                    {assignments.length === 0 && <option value="">{loadingTree ? 'Đang tải…' : 'Chưa có bài tập nào trong lớp này'}</option>}
                    {assignments.map((n) => <option key={n.id} value={n.id}>{n.title}</option>)}
                  </select>
                </div>
              </div>

              <div className="wrap-gap" style={{ alignItems: 'center' }}>
                <button
                  className="btn btn-primary"
                  disabled={!selectedAssignment}
                  onClick={() => setOpen(true)}
                >
                  <IconSparkles /> Ra đề cho bài tập này
                </button>
                {selectedAssignment && hasQuestions && (
                  <ExportQuestionsButtons nodeId={selectedAssignment.id} onError={setError} />
                )}
              </div>
              <p className="hint">
                Cần bài tập khác? Vào <b>Lớp học</b> để tạo bài tập mới rồi quay lại đây chọn.
              </p>
            </>
          )}
        </div>
      </div>

      {open && selectedAssignment && (
        <GenerateQuestionsModal
          nodeId={selectedAssignment.id}
          defaultTopic={selectedAssignment.title}
          onClose={() => setOpen(false)}
          onImported={(count) => {
            setOpen(false)
            setImported((prev) => new Set(prev).add(selectedAssignment.id))
            setNotice(`Đã lưu ${count} câu vào bài tập "${selectedAssignment.title}". Bạn có thể xuất đề ra Word ngay bên dưới.`)
          }}
        />
      )}
    </>
  )
}
