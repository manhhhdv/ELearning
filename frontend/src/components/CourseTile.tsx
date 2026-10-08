import { Link } from 'react-router-dom'

import type { Program } from '../api/types'
import { IconDoc, IconPin } from './icons'

export const COURSE_COVERS = [
  { id: 'math', label: 'Toán học' },
  { id: 'science', label: 'Khoa học' },
  { id: 'code', label: 'Công nghệ' },
  { id: 'literature', label: 'Văn học' },
  { id: 'language', label: 'Ngôn ngữ' },
  { id: 'nature', label: 'Thiên nhiên' },
  { id: 'geography', label: 'Địa lý' },
  { id: 'history', label: 'Lịch sử' },
].map((cover) => ({ ...cover, url: `/covers/${cover.id}.svg` }))

/** Stable local illustrations, including a fallback when a custom image fails. */
export function defaultCourseCover(code: string) {
  let hash = 0
  for (const ch of code) hash = (hash * 31 + ch.charCodeAt(0)) % 997
  return COURSE_COVERS[hash % COURSE_COVERS.length].url
}

export function courseCover(program: Pick<Program, 'coverUrl' | 'code'>): React.CSSProperties {
  const fallback = `url(${JSON.stringify(defaultCourseCover(program.code))})`
  return { backgroundImage: program.coverUrl.trim()
    ? `url(${JSON.stringify(program.coverUrl.trim())}), ${fallback}` : fallback }
}

export function CourseTile({ program }: { program: Program }) {
  const total = program.lessonCount
  const done = program.completedLessonCount
  const percent = total > 0 ? Math.round((done / total) * 100) : 0

  return (
    <Link to={`/hoc/${program.slug}`} className="tile">
      {program.isDefaultCourse && (
        <span className="tile-default-badge" title="Lớp học mặc định — tự động hiện với mọi người">
          <IconPin size={11} /> Bắt buộc
        </span>
      )}
      <div className="tile-cover" style={courseCover(program)}>
        <span className="tile-course-code">{program.code}</span>
      </div>

      <div className="tile-body">
        <h3>{program.title}</h3>
        {program.description && <p className="tile-desc">{program.description}</p>}
      </div>

      {done > 0 && (
        <div className="tile-progress">
          <div className="bar" role="progressbar" aria-label="Tiến độ lớp học" aria-valuenow={percent} aria-valuemin={0} aria-valuemax={100}><i style={{ width: `${percent}%` }} /></div>
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
