import { Link } from 'react-router-dom'

import type { Program } from '../api/types'
import { IconDoc, IconPin } from './icons'

/** Dải màu dự phòng khi khoá học chưa có ảnh bìa, suy từ mã khoá nên luôn ổn định. */
const GRADIENTS = [
  'linear-gradient(135deg, #0b5cd5 0%, #4a9df5 100%)',
  'linear-gradient(135deg, #0b7a44 0%, #37b978 100%)',
  'linear-gradient(135deg, #6d28d9 0%, #a855f7 100%)',
  'linear-gradient(135deg, #b45309 0%, #f5a524 100%)',
  'linear-gradient(135deg, #0f766e 0%, #2dd4bf 100%)',
  'linear-gradient(135deg, #9f1239 0%, #f43f5e 100%)',
]

function fallbackCover(code: string) {
  let hash = 0
  for (const ch of code) hash = (hash * 31 + ch.charCodeAt(0)) % 997
  return GRADIENTS[hash % GRADIENTS.length]
}

/** Kiểu nền cho ảnh bìa: dùng ảnh nếu có, không thì lấy dải màu theo mã khoá. */
export function courseCover(program: Program): React.CSSProperties {
  return program.coverUrl
    ? { backgroundImage: `url(${program.coverUrl})` }
    : { background: fallbackCover(program.code) }
}

export function CourseTile({ program }: { program: Program }) {
  const total = program.lessonCount
  const done = program.completedLessonCount
  const percent = total > 0 ? Math.round((done / total) * 100) : 0

  return (
    <Link to={`/hoc/${program.slug}`} className="tile">
      {program.isDefaultCourse && (
        <span className="tile-default-badge" title="Khoá học mặc định — tự động hiện với mọi người">
          <IconPin size={11} /> Bắt buộc
        </span>
      )}
      <div className="tile-cover" style={courseCover(program)}>
        {!program.coverUrl && program.code}
      </div>

      <div className="tile-body">
        <h3>{program.title}</h3>
      </div>

      {done > 0 && (
        <div className="tile-progress">
          <div className="bar"><i style={{ width: `${percent}%` }} /></div>
          <span>{percent === 100 ? 'Đã hoàn thành' : `Đã học ${done}/${total} bài · ${percent}%`}</span>
        </div>
      )}

      <div className="tile-foot">
        <IconDoc size={16} />
        {program.lessonCount} bài học
        {program.assignmentCount > 0 && ` · ${program.assignmentCount} bài tập`}
      </div>
    </Link>
  )
}
