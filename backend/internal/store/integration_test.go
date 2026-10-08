package store

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/manhnv/elearning/backend/internal/database"
	"github.com/manhnv/elearning/backend/internal/models"
)

// openTestDB mở kết nối tới database thật và chạy migration.
// Bỏ qua test khi chưa đặt TEST_DATABASE_URL để `go test ./...` vẫn chạy được
// trên máy không có Postgres.
func openTestDB(t *testing.T) (*Store, *pgxpool.Pool) {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("đặt TEST_DATABASE_URL để chạy test cần database")
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, url)
	if err != nil {
		t.Fatalf("không kết nối được database: %v", err)
	}
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("chạy migration thất bại: %v", err)
	}
	t.Cleanup(pool.Close)
	return New(pool), pool
}

// seedAssignment dựng sẵn người dùng, chương trình và một bài tập rỗng.
func seedAssignment(t *testing.T, st *Store) (assignmentID, userID uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	suffix := uuid.NewString()[:8]

	var programID uuid.UUID
	err := st.pool.QueryRow(ctx, `
		INSERT INTO users (email, full_name, role) VALUES ($1, 'Học viên test', 'student') RETURNING id`,
		"test-"+suffix+"@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("tạo người dùng: %v", err)
	}
	err = st.pool.QueryRow(ctx, `
		INSERT INTO programs (code, slug, title) VALUES ($1, $1, 'Chương trình test') RETURNING id`,
		"TEST-"+suffix).Scan(&programID)
	if err != nil {
		t.Fatalf("tạo chương trình: %v", err)
	}
	err = st.pool.QueryRow(ctx, `
		INSERT INTO nodes (program_id, kind, slug, title, is_published)
		VALUES ($1, 'assignment', $2, 'Bài kiểm tra test', true) RETURNING id`,
		programID, "bai-test-"+suffix).Scan(&assignmentID)
	if err != nil {
		t.Fatalf("tạo nút bài tập: %v", err)
	}
	if _, err := st.pool.Exec(ctx, `
		INSERT INTO assignments (node_id, max_attempts) VALUES ($1, 0)`, assignmentID); err != nil {
		t.Fatalf("tạo bài tập: %v", err)
	}

	t.Cleanup(func() {
		_, _ = st.pool.Exec(ctx, `DELETE FROM programs WHERE id = $1`, programID)
		_, _ = st.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})
	return assignmentID, userID
}

// TestFourQuestionTypesRoundTrip lưu đủ bốn dạng câu hỏi xuống database, cho
// một học viên nộp bài rồi kiểm tra điểm chấm tự động (Bảng 3.1 và mục 3.4.3).
func TestFourQuestionTypesRoundTrip(t *testing.T) {
	st, _ := openTestDB(t)
	ctx := context.Background()
	assignmentID, userID := seedAssignment(t, st)

	mc, err := st.CreateQuestion(ctx, SaveQuestionParams{
		AssignmentID: assignmentID, Type: models.QuestionSingleChoice,
		Prompt: "Đơn vị của điện trở là gì?", Points: 1,
		Level: models.LevelUnderstand, Explanation: "Điện trở đo bằng ôm.",
		Options: []QuestionOptionInput{
			{Content: "Vôn"}, {Content: "Ampe"}, {Content: "Ôm", IsCorrect: true}, {Content: "Oát"},
		},
	})
	if err != nil {
		t.Fatalf("tạo câu trắc nghiệm: %v", err)
	}
	if mc.Level != models.LevelUnderstand {
		t.Errorf("mức độ nhận thức không được lưu: %q", mc.Level)
	}

	tf, err := st.CreateQuestion(ctx, SaveQuestionParams{
		AssignmentID: assignmentID, Type: models.QuestionTrueFalse,
		Prompt: "Xét các phát biểu sau:", Points: 4, Explanation: "Xem lại định nghĩa.",
		Options: []QuestionOptionInput{
			{Content: "Dòng điện là dòng chuyển dời có hướng của điện tích.", IsCorrect: true},
			{Content: "Đơn vị cường độ dòng điện là vôn.", IsCorrect: false},
			{Content: "Định luật Ôm: I = U/R.", IsCorrect: true},
			{Content: "Điện trở đo bằng ampe.", IsCorrect: false},
		},
	})
	if err != nil {
		t.Fatalf("tạo câu đúng/sai: %v", err)
	}

	fb, err := st.CreateQuestion(ctx, SaveQuestionParams{
		AssignmentID: assignmentID, Type: models.QuestionFillBlank,
		Prompt: "Cường độ dòng điện tỉ lệ ___ với hiệu điện thế.", Points: 2,
		Explanation: "I = U/R.",
		Options: []QuestionOptionInput{
			{Content: "thuận", IsCorrect: true}, {Content: "tỉ lệ thuận", IsCorrect: true},
		},
	})
	if err != nil {
		t.Fatalf("tạo câu điền khuyết: %v", err)
	}

	essay, err := st.CreateQuestion(ctx, SaveQuestionParams{
		AssignmentID: assignmentID, Type: models.QuestionEssay,
		Prompt: "Trình bày cách đo điện trở.", Points: 3, Explanation: "Kiểm tra vận dụng.",
		SampleAnswer: "Dùng vôn kế và ampe kế, tính R = U/I.",
		Rubric:       "Sơ đồ mạch 1đ; công thức 1đ; các bước đo 1đ.",
	})
	if err != nil {
		t.Fatalf("tạo câu tự luận: %v", err)
	}
	if essay.SampleAnswer == "" || essay.Rubric == "" {
		t.Error("đáp án gợi ý và tiêu chí chấm của câu tự luận không được lưu")
	}

	// Mở lượt làm bài rồi nộp: đúng câu trắc nghiệm, đúng 3/4 phát biểu,
	// điền khuyết gõ thừa khoảng trắng và khác hoa thường.
	if _, err := st.StartAttempt(ctx, assignmentID, userID); err != nil {
		t.Fatalf("mở lượt làm bài: %v", err)
	}

	correctOption := func(q *models.Question, want string) uuid.UUID {
		for _, o := range q.Options {
			if o.Content == want {
				return o.ID
			}
		}
		t.Fatalf("không tìm thấy phương án %q", want)
		return uuid.Nil
	}

	sub, err := st.SubmitAssignment(ctx, assignmentID, userID, []AnswerInput{
		{QuestionID: mc.ID, SelectedOptionIDs: []uuid.UUID{correctOption(mc, "Ôm")}},
		{QuestionID: tf.ID, SelectedOptionIDs: []uuid.UUID{
			correctOption(tf, "Dòng điện là dòng chuyển dời có hướng của điện tích."),
			correctOption(tf, "Định luật Ôm: I = U/R."),
			// Đánh dấu nhầm một phát biểu sai là "Đúng" -> mất 1/4 số điểm.
			correctOption(tf, "Điện trở đo bằng ampe."),
		}},
		{QuestionID: fb.ID, EssayText: "  Tỉ Lệ   Thuận  "},
		{QuestionID: essay.ID, EssayText: "Mắc vôn kế song song, ampe kế nối tiếp."},
	})
	if err != nil {
		t.Fatalf("nộp bài: %v", err)
	}

	// 1 (trắc nghiệm) + 3 (đúng/sai: 3/4 × 4) + 2 (điền khuyết) = 6; tự luận chờ chấm.
	if sub.AutoScore != 6 {
		t.Errorf("điểm tự động = %v, mong đợi 6", sub.AutoScore)
	}
	if sub.MaxScore != 10 {
		t.Errorf("tổng điểm tối đa = %v, mong đợi 10", sub.MaxScore)
	}
	if sub.Status != "submitted" {
		t.Errorf("bài có câu tự luận phải ở trạng thái chờ chấm, nhận được %q", sub.Status)
	}

	// Câu tự luận phải xuất hiện trong danh sách chờ AI gợi ý điểm.
	essays, err := st.ListEssayAnswers(ctx, sub.ID)
	if err != nil {
		t.Fatalf("đọc câu tự luận: %v", err)
	}
	if len(essays) != 1 {
		t.Fatalf("mong đợi 1 câu tự luận, nhận được %d", len(essays))
	}
	if essays[0].Rubric == "" {
		t.Error("tiêu chí chấm phải được gửi kèm cho lớp Service AI")
	}

	// Gợi ý của AI không được tự ý đổi điểm chính thức.
	if err := st.SaveAISuggestion(ctx, essays[0].AnswerID, 2, "Thiếu công thức R = U/I."); err != nil {
		t.Fatalf("lưu gợi ý AI: %v", err)
	}
	after, err := st.GetSubmission(ctx, sub.ID)
	if err != nil {
		t.Fatalf("đọc lại bài nộp: %v", err)
	}
	for _, a := range after.Answers {
		if a.QuestionID != essay.ID {
			continue
		}
		if a.AIScore == nil || *a.AIScore != 2 {
			t.Errorf("điểm AI gợi ý = %v, mong đợi 2", a.AIScore)
		}
		if a.Score != 0 {
			t.Errorf("điểm chính thức phải giữ nguyên 0 cho tới khi giáo viên chấm, nhận được %v", a.Score)
		}
	}
	if after.Status != "submitted" {
		t.Errorf("gợi ý của AI không được chốt bài, trạng thái = %q", after.Status)
	}
}

// TestLessonPlanAndAITaskPersistence kiểm tra các bảng mới của lớp chức năng AI.
func TestLessonPlanAndAITaskPersistence(t *testing.T) {
	st, _ := openTestDB(t)
	ctx := context.Background()
	_, userID := seedAssignment(t, st)

	plan, err := st.CreateLessonPlan(ctx, SaveLessonPlanParams{
		OwnerID: userID, Subject: "Vật lí", Grade: "lớp 11", Topic: "Định luật Ôm",
		DurationMinutes: 45, Title: "Định luật Ôm",
		Content: []byte(`{"title":"Định luật Ôm","activities":[]}`),
	})
	if err != nil {
		t.Fatalf("lưu giáo án: %v", err)
	}
	if plans, err := st.ListLessonPlans(ctx, userID); err != nil || len(plans) != 1 {
		t.Fatalf("liệt kê giáo án: %d bản, lỗi %v", len(plans), err)
	}
	if err := st.DeleteLessonPlan(ctx, plan.ID, userID); err != nil {
		t.Fatalf("xoá giáo án: %v", err)
	}

	if err := st.LogAITask(ctx, LogAITaskParams{
		UserID: userID, Kind: models.AIKindQuestions, Title: "Định luật Ôm (5 câu)", Attempts: 2,
	}); err != nil {
		t.Fatalf("ghi nhật ký tác vụ AI: %v", err)
	}
	tasks, err := st.ListAITasks(ctx, userID, 10)
	if err != nil || len(tasks) != 1 {
		t.Fatalf("đọc nhật ký: %d dòng, lỗi %v", len(tasks), err)
	}
	if tasks[0].Attempts != 2 {
		t.Errorf("số lần sinh lại = %d, mong đợi 2", tasks[0].Attempts)
	}

	// Hội thoại trợ lý AI: lưu được và không đọc chéo sang người dùng khác.
	convID, err := st.CreateConversation(ctx, userID, "Giải thích giúp em định luật Ôm", nil, nil)
	if err != nil {
		t.Fatalf("tạo hội thoại: %v", err)
	}
	if err := st.AppendMessages(ctx, convID, "Định luật Ôm là gì?", "Định luật Ôm cho biết…"); err != nil {
		t.Fatalf("lưu tin nhắn: %v", err)
	}
	conv, err := st.GetConversation(ctx, convID, userID)
	if err != nil {
		t.Fatalf("đọc hội thoại: %v", err)
	}
	if len(conv.Messages) != 2 {
		t.Errorf("mong đợi 2 tin nhắn, nhận được %d", len(conv.Messages))
	}
	if _, err := st.GetConversation(ctx, convID, uuid.New()); err == nil {
		t.Error("người dùng khác không được đọc hội thoại này")
	}
}

// TestPasswordResetRequestFlow kiểm tra luồng quên mật khẩu: gửi yêu cầu,
// không tạo trùng, và tự đóng khi admin cấp mật khẩu mới.
func TestPasswordResetRequestFlow(t *testing.T) {
	st, _ := openTestDB(t)
	ctx := context.Background()
	_, userID := seedAssignment(t, st)

	user, err := st.GetUserByID(ctx, userID)
	if err != nil {
		t.Fatalf("đọc người dùng: %v", err)
	}

	if err := st.CreatePasswordResetRequest(ctx, user.Email, "quên mật khẩu ạ"); err != nil {
		t.Fatalf("gửi yêu cầu: %v", err)
	}
	// Bấm lại khi đang chờ xử lý thì không tạo thêm dòng mới.
	if err := st.CreatePasswordResetRequest(ctx, user.Email, "gửi lại"); err != nil {
		t.Fatalf("gửi lại yêu cầu: %v", err)
	}

	pending, err := st.ListPasswordResetRequests(ctx, true)
	if err != nil {
		t.Fatalf("liệt kê yêu cầu: %v", err)
	}
	mine := 0
	for _, q := range pending {
		if q.Email == user.Email {
			mine++
			if q.UserID == nil || *q.UserID != userID {
				t.Error("yêu cầu phải khớp với tài khoản có email đó")
			}
		}
	}
	if mine != 1 {
		t.Fatalf("mong đợi đúng 1 yêu cầu đang chờ, nhận được %d", mine)
	}

	// Admin cấp mật khẩu mới -> yêu cầu tự đóng.
	if err := st.ResolvePasswordResetsForUser(ctx, userID, userID); err != nil {
		t.Fatalf("đóng yêu cầu: %v", err)
	}
	pending, err = st.ListPasswordResetRequests(ctx, true)
	if err != nil {
		t.Fatalf("liệt kê lại: %v", err)
	}
	for _, q := range pending {
		if q.Email == user.Email {
			t.Error("yêu cầu phải được đóng sau khi cấp mật khẩu mới")
		}
	}

	// Đóng rồi thì gửi yêu cầu mới lại được.
	if err := st.CreatePasswordResetRequest(ctx, user.Email, "lại quên nữa"); err != nil {
		t.Fatalf("gửi yêu cầu lần hai: %v", err)
	}
}

// TestWorkspaceQueries chạy các truy vấn của trang tổng quan để chắc chắn
// chúng hợp lệ về cú pháp và trả về dữ liệu rỗng an toàn.
func TestWorkspaceQueries(t *testing.T) {
	st, _ := openTestDB(t)
	ctx := context.Background()
	_, userID := seedAssignment(t, st)

	ws, err := st.StudentWorkspace(ctx, userID)
	if err != nil {
		t.Fatalf("StudentWorkspace: %v", err)
	}
	if ws.Schedule == nil || ws.Progress == nil {
		t.Error("các danh sách phải rỗng chứ không được nil, để JSON trả về [] thay vì null")
	}

	tw, err := st.TrainerWorkspace(ctx, userID, true)
	if err != nil {
		t.Fatalf("TrainerWorkspace: %v", err)
	}
	if tw.Alerts == nil {
		t.Error("danh sách cảnh báo phải rỗng chứ không được nil")
	}
}
