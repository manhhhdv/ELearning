package store

// sampleVanPrimary là các lớp học mẫu môn Tiếng Việt lớp 1–5 và Ngữ văn lớp 6
// theo Chương trình GDPT 2018: âm vần, đọc hiểu, luyện từ và câu, viết.
var sampleVanPrimary = []sampleCourse{
	{
		Code:        "TIENGVIET1",
		Title:       "Tiếng Việt lớp 1",
		Description: "Lớp học mẫu môn Tiếng Việt lớp 1 theo Chương trình GDPT 2018: làm quen âm và chữ, học vần, đánh vần và tập đọc câu ngắn.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Làm quen với âm và chữ",
				Lessons: []sampleLesson{
					{
						Title:           "Âm, chữ và dấu thanh",
						DurationMinutes: 15,
						Body: `## 1. Âm và chữ

Khi nói, ta phát ra **âm**. Khi viết, ta dùng **chữ** để ghi lại âm.

- Chữ **a** ghi âm "a".
- Chữ **b** ghi âm "bờ".
- Chữ **c** ghi âm "cờ".

## 2. Sáu thanh của tiếng Việt

Tiếng Việt có 6 thanh. Thanh ngang không có dấu.

| Thanh | Dấu | Ví dụ |
|---|---|---|
| Ngang | không có dấu | ba |
| Huyền | dấu huyền | bà |
| Sắc | dấu sắc | bá |
| Hỏi | dấu hỏi | bả |
| Ngã | dấu ngã | bã |
| Nặng | dấu nặng | bạ |

## 3. Luyện đọc

Em đọc to, rõ từng tiếng:

- ba – bà – bá
- ca – cà – cá
- bà – cá – con ba ba

## 4. Bài tập mẫu

**Đề:** Tiếng "cá" có dấu gì?

**Lời giải:** Tiếng "cá" có **dấu sắc**, đặt trên chữ a.

> [!TIP]
> Dấu nặng đặt ở dưới chữ. Các dấu còn lại (huyền, sắc, hỏi, ngã) đặt ở trên chữ.`,
					},
					{
						Title:           "Ghép âm thành tiếng",
						DurationMinutes: 15,
						Body: `## 1. Tiếng gồm những phần nào?

Tiếng "ba" có âm đầu **b** và vần **a**.

Ghép lại, em đánh vần: **bờ – a – ba**.

## 2. Đánh vần tiếng có dấu thanh

Em đánh vần tiếng không dấu trước, rồi thêm thanh:

- bờ – a – ba – huyền – **bà**
- cờ – a – ca – sắc – **cá**
- cờ – o – co – nặng – **cọ**

## 3. Viết c hay k?

Âm "cờ" thường được viết bằng chữ **c** hoặc chữ **k**:

- Viết **k** khi đứng trước **e, ê, i**: kẻ, kể, kì.
- Viết **c** khi đứng trước các chữ khác: ca, cô, cụ.

Tương tự:

- **gh** đứng trước e, ê, i (ghế, ghi); còn lại viết **g** (gà, gõ).
- **ngh** đứng trước e, ê, i (nghe, nghỉ); còn lại viết **ng** (ngõ, ngà).

## 4. Bài tập mẫu

**Đề:** Điền c hay k: ...ì cọ, ...á cờ.

**Lời giải:** Trước "i" viết **k**: **kì cọ**. Trước "a" viết **c**: **cá cờ**.

> [!NOTE]
> Mẹo nhớ: k, gh, ngh là ba bạn thân, chỉ đứng trước e, ê, i.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Học vần và tập đọc",
				Lessons: []sampleLesson{
					{
						Title:           "Học vần: an, at, anh",
						DurationMinutes: 20,
						Body: `## 1. Vần an

Vần **an** có âm **a** đứng trước, âm **n** đứng sau.

Đánh vần: **a – nờ – an**.

Ghép tiếng: bờ – an – ban – huyền – **bàn**; lờ – an – **lan**.

## 2. Vần at

Vần **at** có âm **a** và âm **t**. Đánh vần: **a – tờ – at**.

Ghép tiếng: hờ – at – **hát**; mờ – at – **mát**.

Lưu ý: tiếng có vần kết thúc bằng t chỉ mang thanh **sắc** hoặc **nặng**: hát, hạt.

## 3. Vần anh

Vần **anh** có âm **a** và âm **nh**. Đánh vần: **a – nhờ – anh**.

Ghép tiếng: xờ – anh – **xanh**; chờ – anh – **chanh**.

## 4. Luyện đọc từ

bàn học, hoa lan, ca hát, gió mát, quả chanh, lá xanh

## 5. Bài tập mẫu

**Đề:** Tiếng "lan" có vần gì?

**Lời giải:** Bỏ âm đầu **l**, phần còn lại là vần **an**.

> [!TIP]
> Muốn tìm vần của một tiếng, em bỏ âm đầu và dấu thanh, phần còn lại chính là vần.`,
					},
					{
						Title:           "Tập đọc câu ngắn",
						DurationMinutes: 20,
						Body: `## 1. Câu là gì?

Câu nói trọn một ý.

- Chữ đầu câu viết **hoa**.
- Cuối câu có **dấu chấm**.

Ví dụ: **B**é Hà có con mèo.

## 2. Đọc câu ngắn

Em đọc chậm, rõ từng tiếng, nghỉ hơi ở dấu chấm:

Bé Hà có con mèo. Mèo có bộ lông mềm. Bé cho mèo ăn cá. Mèo kêu: "Meo meo!"

## 3. Trả lời câu hỏi

- Bé Hà có con gì? Bé Hà có **con mèo**.
- Bé cho mèo ăn gì? Bé cho mèo ăn **cá**.
- Mèo kêu thế nào? Mèo kêu **meo meo**.

## 4. Bài tập mẫu

**Đề:** Viết lại cho đúng: "bé hà đi học"

**Lời giải:** Viết hoa chữ đầu câu và tên bạn Hà, thêm dấu chấm cuối câu: **Bé Hà đi học.**

> [!NOTE]
> Tên người (Hà, Nam, Lan...) luôn viết hoa, dù đứng ở đầu hay giữa câu.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Âm, vần, dấu thanh và câu ngắn",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Tiếng “bà” có dấu thanh gì?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Tiếng “bà” có dấu huyền đặt trên chữ a.",
					Options: []sampleOption{
						{Content: "Dấu sắc"},
						{Content: "Dấu huyền", IsCorrect: true},
						{Content: "Dấu hỏi"},
						{Content: "Dấu nặng"},
					},
				},
				{
					Prompt: "Từ nào viết đúng chính tả?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Trước chữ i phải viết k, nên viết đúng là “kì cọ”.",
					Options: []sampleOption{
						{Content: "cì cọ"},
						{Content: "kì kọ"},
						{Content: "cì kọ"},
						{Content: "kì cọ", IsCorrect: true},
					},
				},
				{
					Prompt: "Tiếng “lan” có vần gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Bỏ âm đầu l, phần còn lại của tiếng “lan” là vần an.",
					Options: []sampleOption{
						{Content: "an", IsCorrect: true},
						{Content: "la"},
						{Content: "n"},
						{Content: "lan"},
					},
				},
				{
					Prompt: "Đánh vần tiếng “xanh” thế nào là đúng?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Tiếng “xanh” gồm âm đầu x và vần anh: xờ – anh – xanh.",
					Options: []sampleOption{
						{Content: "xờ – an – xanh"},
						{Content: "xờ – a – xanh"},
						{Content: "xờ – anh – xanh", IsCorrect: true},
						{Content: "anh – xờ – xanh"},
					},
				},
				{
					Prompt: "Câu nào viết đúng?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu đúng phải viết hoa chữ đầu câu, viết hoa tên riêng Hà và có dấu chấm ở cuối câu.",
					Options: []sampleOption{
						{Content: "bé Hà có con mèo."},
						{Content: "Bé hà có con mèo."},
						{Content: "Bé Hà có con mèo.", IsCorrect: true},
						{Content: "bé hà có con mèo"},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGVIET2",
		Title:       "Tiếng Việt lớp 2",
		Description: "Lớp học mẫu môn Tiếng Việt lớp 2 theo Chương trình GDPT 2018: đọc hiểu bài đọc ngắn, từ chỉ sự vật – hoạt động – đặc điểm, các kiểu câu và viết đoạn văn ngắn.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Em lớn lên từng ngày",
				Lessons: []sampleLesson{
					{
						Title:           "Đọc hiểu: Ngày đầu đến lớp",
						DurationMinutes: 20,
						Body: `## 1. Đọc bài

Sáng nay, mẹ đưa Minh đến trường mới. Sân trường rộng, có hàng cây bàng xanh mát. Lúc đầu, Minh hơi rụt rè. Cô giáo mỉm cười, dắt Minh vào lớp và giới thiệu bạn với cả lớp. Bạn Na ngồi bên cạnh cho Minh mượn bút chì màu. Giờ ra chơi, hai bạn cùng chơi nhảy dây. Chiều về, Minh kể với mẹ: "Con có thêm một người bạn mới!"

## 2. Tìm hiểu bài

- Ai đưa Minh đến trường? **Mẹ** đưa Minh đến trường.
- Lúc đầu, Minh cảm thấy thế nào? Minh thấy **hơi rụt rè**.
- Ai giúp Minh bớt rụt rè? **Cô giáo** và **bạn Na**.
- Vì sao chiều về Minh rất vui? Vì Minh **có thêm bạn mới**.

## 3. Cách đọc hiểu một bài

1. Đọc chậm toàn bài một lần.
2. Chú ý tên nhân vật, nơi chốn, việc làm.
3. Đọc câu hỏi, tìm câu trả lời ngay trong bài.
4. Trả lời thành câu đầy đủ.

## 4. Bài tập mẫu

**Đề:** Bạn Na đã làm gì để giúp Minh?

**Lời giải:** Bạn Na **cho Minh mượn bút chì màu** và **cùng Minh chơi nhảy dây** lúc ra chơi.

> [!TIP]
> Khi trả lời, em nhắc lại một phần câu hỏi để câu trả lời rõ ràng, đầy đủ.`,
					},
					{
						Title:           "Từ chỉ sự vật, hoạt động, đặc điểm",
						DurationMinutes: 20,
						Body: `## 1. Từ chỉ sự vật

Là từ gọi tên người, con vật, đồ vật, cây cối...

Ví dụ: cô giáo, học sinh, con mèo, quyển vở, cây bàng.

## 2. Từ chỉ hoạt động

Là từ chỉ việc làm của người, con vật.

Ví dụ: đọc, viết, chạy, nhảy dây, hót.

## 3. Từ chỉ đặc điểm

Là từ chỉ màu sắc, hình dáng, tính nết...

Ví dụ: xanh mát, đỏ tươi, tròn xoe, cao, chăm chỉ.

## 4. Bài tập mẫu

**Đề:** Xếp các từ sau vào ba nhóm: bút chì, viết, vàng tươi, cái bàn, chạy, nhỏ xinh.

**Lời giải:**

- Từ chỉ sự vật: bút chì, cái bàn.
- Từ chỉ hoạt động: viết, chạy.
- Từ chỉ đặc điểm: vàng tươi, nhỏ xinh.

> [!TIP]
> Từ chỉ sự vật trả lời câu hỏi "Ai? Cái gì? Con gì?". Từ chỉ hoạt động trả lời câu hỏi "Làm gì?". Từ chỉ đặc điểm trả lời câu hỏi "Thế nào?".`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Đi học vui sao",
				Lessons: []sampleLesson{
					{
						Title:           "Câu giới thiệu và câu nêu hoạt động",
						DurationMinutes: 20,
						Body: `## 1. Câu giới thiệu (Ai là gì?)

Dùng để giới thiệu về một người, một vật.

- Đây là cô giáo của em.
- Bạn Lan là lớp trưởng lớp em.

## 2. Câu nêu hoạt động (Ai làm gì?)

Dùng để kể việc làm của người, con vật.

- Mẹ em đang nấu cơm.
- Các bạn chơi đá cầu.

## 3. Câu nêu đặc điểm (Ai thế nào?)

Dùng để nói về đặc điểm của người, vật.

- Chiếc cặp của em rất mới.

## 4. Dấu câu

- **Dấu chấm** đặt cuối câu kể: Em thích đọc truyện.
- **Dấu chấm hỏi** đặt cuối câu hỏi: Bạn tên là gì?

## 5. Bài tập mẫu

**Đề:** Câu "Chú mèo đang rửa mặt." là kiểu câu nào?

**Lời giải:** "Rửa mặt" là một hoạt động, nên đây là **câu nêu hoạt động (Ai làm gì?)**.

> [!NOTE]
> Câu giới thiệu thường có từ "là". Câu nêu hoạt động có từ chỉ hoạt động như đọc, viết, chạy, nấu...`,
					},
					{
						Title:           "Viết đoạn văn tả đồ dùng học tập",
						DurationMinutes: 20,
						Body: `## 1. Đoạn văn gồm mấy câu?

Em viết từ 3 đến 5 câu, các câu cùng nói về **một đồ dùng học tập**.

## 2. Gợi ý các câu

1. Giới thiệu: Đó là đồ dùng gì? Ai tặng hay mua cho em?
2. Tả: Đồ dùng có hình dáng, màu sắc thế nào?
3. Công dụng: Em dùng nó để làm gì?
4. Tình cảm: Em giữ gìn nó ra sao?

## 3. Đoạn văn mẫu

Em có một chiếc hộp bút mới do bố tặng. Hộp bút hình chữ nhật, màu xanh da trời. Bên trong có hai ngăn để bút chì và thước kẻ. Nhờ có hộp bút, đồ dùng của em luôn gọn gàng. Em luôn cất hộp bút cẩn thận vào cặp sau khi học.

## 4. Tự kiểm tra bài viết

- Chữ đầu câu đã viết hoa chưa?
- Cuối mỗi câu đã có dấu chấm chưa?
- Các câu có cùng nói về một đồ dùng không?

## 5. Bài tập mẫu

**Đề:** Viết một câu tả màu sắc chiếc thước kẻ của em.

**Lời giải (gợi ý):** Chiếc thước kẻ của em **màu hồng nhạt**, trong suốt.

> [!TIP]
> Dùng từ chỉ đặc điểm (màu đỏ, tròn, nhỏ xinh...) giúp đoạn văn của em rõ ràng và sinh động hơn.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Đọc hiểu, từ và câu",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Từ nào là từ chỉ hoạt động?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "“Chạy” chỉ việc làm của người, con vật nên là từ chỉ hoạt động; các từ còn lại chỉ sự vật.",
					Options: []sampleOption{
						{Content: "cái bút"},
						{Content: "quyển vở"},
						{Content: "chạy", IsCorrect: true},
						{Content: "cặp sách"},
					},
				},
				{
					Prompt: "Từ nào là từ chỉ đặc điểm?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "“Đỏ tươi” chỉ màu sắc, trả lời câu hỏi “thế nào?” nên là từ chỉ đặc điểm.",
					Options: []sampleOption{
						{Content: "quả táo"},
						{Content: "ăn"},
						{Content: "nhảy dây"},
						{Content: "đỏ tươi", IsCorrect: true},
					},
				},
				{
					Prompt: "Câu “Bạn Lan là lớp trưởng lớp em.” thuộc kiểu câu nào?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Câu dùng từ “là” để giới thiệu bạn Lan, nên là câu giới thiệu (Ai là gì?).",
					Options: []sampleOption{
						{Content: "Câu giới thiệu (Ai là gì?)", IsCorrect: true},
						{Content: "Câu nêu hoạt động (Ai làm gì?)"},
						{Content: "Câu nêu đặc điểm (Ai thế nào?)"},
						{Content: "Câu hỏi"},
					},
				},
				{
					Prompt: "Câu nào là câu nêu hoạt động?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "“Nấu cơm” là việc làm, nên “Mẹ em đang nấu cơm.” là câu nêu hoạt động (Ai làm gì?).",
					Options: []sampleOption{
						{Content: "Bố em là bác sĩ."},
						{Content: "Mẹ em đang nấu cơm.", IsCorrect: true},
						{Content: "Chiếc cặp rất mới."},
						{Content: "Đây là cô giáo em."},
					},
				},
				{
					Prompt: "Chọn từ thích hợp điền vào chỗ trống để tả chiếc bút: “Em có một chiếc bút chì ... . Em dùng nó để vẽ tranh.”",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Chỗ trống cần một từ chỉ đặc điểm của chiếc bút; “màu vàng” chỉ màu sắc nên phù hợp.",
					Options: []sampleOption{
						{Content: "chạy nhanh"},
						{Content: "đang ngủ"},
						{Content: "màu vàng", IsCorrect: true},
						{Content: "cái cặp"},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGVIET3",
		Title:       "Tiếng Việt lớp 3",
		Description: "Lớp học mẫu môn Tiếng Việt lớp 3 theo Chương trình GDPT 2018: đọc hiểu tìm ý chính, biện pháp so sánh, các kiểu câu và viết đoạn văn nêu tình cảm.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Những trải nghiệm thú vị",
				Lessons: []sampleLesson{
					{
						Title:           "Đọc hiểu: Tìm ý chính của đoạn văn",
						DurationMinutes: 20,
						Body: `## 1. Đọc đoạn văn

Mùa hè năm nay, em về quê ngoại. Buổi sáng, em theo ông ra vườn tưới rau. Buổi trưa, bà kể chuyện cổ tích cho em nghe dưới gốc nhãn. Chiều chiều, em cùng các bạn thả diều trên đê. Những ngày ở quê ngoại thật vui và đáng nhớ.

## 2. Ý chính là gì?

**Ý chính** là điều quan trọng nhất mà đoạn văn muốn nói. Ý chính thường nằm ở câu mở đầu hoặc câu kết đoạn.

Ý chính của đoạn văn trên: **Những ngày hè vui vẻ, đáng nhớ của em ở quê ngoại.**

## 3. Cách tìm ý chính

1. Đọc kĩ toàn đoạn.
2. Tự hỏi: Đoạn văn nói về ai? Về việc gì?
3. Tìm câu nêu ý khái quát (thường ở đầu hoặc cuối đoạn).
4. Nói lại ý chính bằng một câu ngắn của em.

## 4. Bài tập mẫu

**Đề:** Các câu bắt đầu bằng "Buổi sáng", "Buổi trưa", "Chiều chiều" có tác dụng gì?

**Lời giải:** Các câu này kể những việc em làm **theo trình tự thời gian trong ngày**, giúp người đọc hiểu vì sao những ngày ở quê ngoại lại vui như vậy.

> [!TIP]
> Câu kết "Những ngày ở quê ngoại thật vui và đáng nhớ" gói lại ý của cả đoạn. Hãy chú ý những câu như thế khi tìm ý chính.`,
					},
					{
						Title:           "Biện pháp so sánh",
						DurationMinutes: 20,
						Body: `## 1. So sánh là gì?

**So sánh** là đối chiếu sự vật này với sự vật khác có nét giống nhau, giúp câu văn sinh động, dễ hình dung.

## 2. Các bộ phận của một hình ảnh so sánh

| Sự vật 1 | Đặc điểm | Từ so sánh | Sự vật 2 |
|---|---|---|---|
| Mặt trăng | tròn | như | cái đĩa bạc |

Một số từ so sánh thường gặp: như, là, giống như, tựa, như là.

## 3. So sánh trong ca dao

Công cha như núi Thái Sơn,
Nghĩa mẹ như nước trong nguồn chảy ra.

- **Công cha** được so sánh với **núi Thái Sơn**: to lớn, vững chãi.
- **Nghĩa mẹ** được so sánh với **nước trong nguồn**: dồi dào, không bao giờ cạn.

## 4. Bài tập mẫu

**Đề:** Tìm hình ảnh so sánh trong câu: "Những chùm hoa phượng đỏ rực như những ngọn lửa."

**Lời giải:** Chùm hoa phượng được so sánh với **ngọn lửa**; từ so sánh là **như**; nét giống nhau là màu **đỏ rực**.

> [!NOTE]
> Có khi câu so sánh không nêu đặc điểm, như "Công cha như núi Thái Sơn". Khi đó, em tự suy nghĩ xem hai sự vật giống nhau ở điểm nào.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Mái nhà yêu thương",
				Lessons: []sampleLesson{
					{
						Title:           "Câu kể, câu hỏi, câu cảm, câu khiến",
						DurationMinutes: 20,
						Body: `## 1. Câu kể

Dùng để kể, tả hoặc giới thiệu. Cuối câu có **dấu chấm**.

Ví dụ: Bà em trồng một vườn rau.

## 2. Câu hỏi

Dùng để hỏi điều mình chưa biết. Thường có các từ: ai, gì, nào, sao, không... Cuối câu có **dấu chấm hỏi**.

Ví dụ: Bà ơi, bà trồng rau gì thế?

## 3. Câu cảm

Dùng để bộc lộ cảm xúc (vui, buồn, ngạc nhiên...). Thường có các từ: ôi, chao, quá, thật... Cuối câu có **dấu chấm than**.

Ví dụ: Ôi, vườn rau xanh quá!

## 4. Câu khiến

Dùng để yêu cầu, đề nghị, khuyên bảo. Thường có các từ: hãy, đừng, chớ, nhé, nào... Cuối câu thường có **dấu chấm than** (hoặc dấu chấm).

Ví dụ: Cháu hãy tưới rau giúp bà nhé!

## 5. Bài tập mẫu

**Đề:** Câu "Đừng ngắt lá non nhé!" là kiểu câu gì?

**Lời giải:** Câu có từ "đừng", dùng để khuyên bảo, nên là **câu khiến**.

> [!TIP]
> Muốn xác định kiểu câu, em hãy hỏi: Câu này dùng để làm gì? Kể, hỏi, bộc lộ cảm xúc hay yêu cầu?`,
					},
					{
						Title:           "Viết đoạn văn nêu tình cảm với người thân",
						DurationMinutes: 20,
						Body: `## 1. Bố cục đoạn văn

1. **Câu mở đầu:** giới thiệu người thân em muốn viết.
2. **Các câu giữa:** kể một vài đặc điểm hoặc kỉ niệm khiến em yêu quý người đó.
3. **Câu kết:** nêu tình cảm, mong ước của em.

## 2. Đoạn văn mẫu

Người em yêu quý nhất là bà nội. Tóc bà đã bạc nhưng đôi mắt bà vẫn hiền hậu. Tối nào bà cũng kể chuyện cổ tích cho em nghe. Có lần em bị ốm, bà thức suốt đêm để chăm sóc em. Em rất yêu bà và mong bà luôn khỏe mạnh.

## 3. Nhận xét đoạn văn mẫu

- Câu đầu giới thiệu người thân: bà nội.
- Câu 2–4 nêu đặc điểm và kỉ niệm về bà.
- Câu cuối nói lên tình cảm, mong ước của em.

## 4. Bài tập mẫu

**Đề:** Viết một câu kể một việc làm của mẹ khiến em cảm động.

**Lời giải (gợi ý):** Ngày mưa, mẹ đội mưa đến trường đón em và che ô cho em suốt đường về.

> [!NOTE]
> Hãy viết bằng những kỉ niệm có thật của em. Chi tiết càng cụ thể, tình cảm trong đoạn văn càng chân thành.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Ý chính, so sánh và các kiểu câu",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Trong câu ca dao “Công cha như núi Thái Sơn”, sự vật nào được so sánh với núi Thái Sơn?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Câu ca dao so sánh công cha với núi Thái Sơn, qua từ so sánh “như”.",
					Options: []sampleOption{
						{Content: "Nghĩa mẹ"},
						{Content: "Nước trong nguồn"},
						{Content: "Công cha", IsCorrect: true},
						{Content: "Con cái"},
					},
				},
				{
					Prompt: "Câu “Ôi, bông hoa đẹp quá!” là kiểu câu gì?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Câu có từ “ôi”, “quá” và dấu chấm than, dùng để bộc lộ cảm xúc nên là câu cảm.",
					Options: []sampleOption{
						{Content: "Câu kể"},
						{Content: "Câu cảm", IsCorrect: true},
						{Content: "Câu hỏi"},
						{Content: "Câu khiến"},
					},
				},
				{
					Prompt: "Câu nào có hình ảnh so sánh?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Câu này so sánh mặt trăng với cái đĩa bạc bằng từ “như”, nét giống nhau là “tròn”.",
					Options: []sampleOption{
						{Content: "Mặt trăng rất tròn."},
						{Content: "Mặt trăng chiếu sáng sân nhà."},
						{Content: "Em ngắm trăng cùng bà."},
						{Content: "Mặt trăng tròn như cái đĩa bạc.", IsCorrect: true},
					},
				},
				{
					Prompt: "Câu “Con hãy cất sách vở vào cặp.” là kiểu câu gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Câu có từ “hãy”, dùng để yêu cầu người khác làm một việc nên là câu khiến.",
					Options: []sampleOption{
						{Content: "Câu khiến", IsCorrect: true},
						{Content: "Câu kể"},
						{Content: "Câu cảm"},
						{Content: "Câu hỏi"},
					},
				},
				{
					Prompt: "Khi viết đoạn văn nêu tình cảm với bà, câu nào phù hợp nhất để làm câu kết?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu kết cần nêu tình cảm, mong ước của em; các câu khác chỉ giới thiệu hoặc kể chi tiết.",
					Options: []sampleOption{
						{Content: "Bà em năm nay bảy mươi tuổi."},
						{Content: "Tối nào bà cũng kể chuyện cho em nghe."},
						{Content: "Em rất yêu bà và mong bà luôn khỏe mạnh.", IsCorrect: true},
						{Content: "Nhà bà ở gần cánh đồng lúa."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGVIET4",
		Title:       "Tiếng Việt lớp 4",
		Description: "Lớp học mẫu môn Tiếng Việt lớp 4 theo Chương trình GDPT 2018: đọc hiểu truyện cổ tích, danh từ – động từ – tính từ và viết bài văn miêu tả cây cối.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Câu chuyện và nhân vật",
				Lessons: []sampleLesson{
					{
						Title:           "Đọc hiểu truyện: Ăn khế trả vàng",
						DurationMinutes: 20,
						Body: `## 1. Tóm tắt câu chuyện

Ngày xưa, có hai anh em mồ côi cha mẹ. Khi chia gia tài, người anh tham lam lấy hết nhà cửa, ruộng vườn, chỉ để cho em một túp lều và một cây khế.

Một hôm, có con chim lạ đến ăn khế. Chim hứa sẽ trả vàng và dặn người em may túi ba gang để đựng. Người em làm đúng lời chim, chỉ lấy một ít vàng vừa đủ dùng, cuộc sống trở nên khá giả.

Người anh biết chuyện, đòi đổi cả gia tài lấy cây khế. Hắn may túi to hơn nhiều, lấy thật nhiều vàng. Trên đường về, chim mỏi cánh vì quá nặng, người anh rơi xuống biển.

## 2. Nhận xét nhân vật

| Nhân vật | Việc làm | Tính cách |
|---|---|---|
| Người anh | Chiếm hết gia tài, may túi to, lấy quá nhiều vàng | Tham lam, ích kỉ |
| Người em | Nhận cây khế, làm đúng lời chim, lấy vàng vừa đủ | Hiền lành, chăm chỉ, không tham |

## 3. Bài học rút ra

Truyện khen ngợi người hiền lành, biết đủ; phê phán thói tham lam. Truyện thể hiện niềm tin "ở hiền gặp lành", "tham thì thâm".

## 4. Bài tập mẫu

**Đề:** Vì sao người em được chim trả vàng mà vẫn bình an?

**Lời giải:** Vì người em **hiền lành, không tham lam**, chỉ may túi ba gang như lời chim dặn và lấy vàng vừa đủ.

> [!TIP]
> Muốn nhận xét tính cách nhân vật, em hãy dựa vào việc làm, lời nói và cách ứng xử của nhân vật trong truyện.`,
					},
					{
						Title:           "Danh từ",
						DurationMinutes: 20,
						Body: `## 1. Danh từ là gì?

**Danh từ** là từ chỉ sự vật: người, vật, hiện tượng tự nhiên, thời gian, đơn vị...

| Loại | Ví dụ |
|---|---|
| Chỉ người | bác sĩ, học sinh, ông bà |
| Chỉ vật | cái bàn, cây khế, quyển sách |
| Chỉ hiện tượng | mưa, gió, sấm, cầu vồng |
| Chỉ thời gian | buổi sáng, mùa hè, năm học |
| Chỉ đơn vị | cái, con, chiếc, mét, ki-lô-gam |

## 2. Danh từ chung và danh từ riêng

- **Danh từ chung** là tên gọi chung của một loại sự vật: sông, núi, thành phố.
- **Danh từ riêng** là tên riêng của một sự vật, **phải viết hoa**: sông Hồng, núi Ba Vì, Hà Nội.

## 3. Bài tập mẫu

**Đề:** Tìm danh từ trong câu: "Mùa thu, lá vàng rơi đầy sân trường."

**Lời giải:** Các danh từ: **mùa thu** (thời gian), **lá** (vật), **sân trường** (vật, nơi chốn).

**Đề:** Viết lại cho đúng: "em sống ở thành phố đà nẵng."

**Lời giải:** Em sống ở thành phố **Đà Nẵng**. (Chữ đầu câu và tên riêng phải viết hoa.)

> [!NOTE]
> Với tên người, tên địa lí Việt Nam, em viết hoa chữ cái đầu của mỗi tiếng: Nguyễn Văn An, Đà Nẵng, Phú Quốc.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Thiên nhiên quanh em",
				Lessons: []sampleLesson{
					{
						Title:           "Động từ và tính từ",
						DurationMinutes: 20,
						Body: `## 1. Động từ

**Động từ** là từ chỉ hoạt động, trạng thái của sự vật.

- Chỉ hoạt động: chạy, đọc, viết, nở, bay.
- Chỉ trạng thái: ngủ, thức, yêu, ghét, lo lắng.

## 2. Tính từ

**Tính từ** là từ chỉ đặc điểm, tính chất của sự vật, hoạt động.

- Màu sắc: xanh, đỏ rực, vàng óng.
- Hình dáng, kích thước: cao, thấp, tròn, to lớn.
- Tính chất: hiền lành, chăm chỉ, thơm ngát.

Có thể thêm các từ chỉ mức độ: **rất** đẹp, đẹp **quá**, đẹp **lắm**.

## 3. Phân tích ví dụ

Câu: "Cây phượng nở hoa đỏ rực."

- Danh từ: cây phượng, hoa.
- Động từ: **nở**.
- Tính từ: **đỏ rực**.

## 4. Bài tập mẫu

**Đề:** Tìm động từ và tính từ trong câu: "Đàn chim hót líu lo trên cành cây cao."

**Lời giải:** Động từ: **hót**. Tính từ: **líu lo** (đặc điểm của tiếng hót), **cao** (đặc điểm của cành cây).

> [!TIP]
> Mẹo: Từ nào kết hợp được với "đang", "sẽ", "đã" để chỉ việc làm thường là động từ; từ nào kết hợp được với "rất", "quá" thường là tính từ.`,
					},
					{
						Title:           "Viết bài văn miêu tả cây cối",
						DurationMinutes: 20,
						Body: `## 1. Bố cục bài văn

1. **Mở bài:** giới thiệu cây định tả (cây gì, ở đâu, có từ bao giờ).
2. **Thân bài:** tả bao quát rồi tả từng bộ phận (thân, cành, lá, hoa, quả); tả cây thay đổi theo mùa; nêu ích lợi của cây.
3. **Kết bài:** nêu tình cảm, suy nghĩ của em về cây.

## 2. Cách quan sát

- Quan sát từ xa đến gần, từ bao quát đến từng bộ phận.
- Dùng nhiều giác quan: mắt nhìn màu sắc, tai nghe tiếng lá, mũi ngửi mùi hoa.
- Dùng tính từ gợi tả và hình ảnh so sánh.

## 3. Đoạn văn mẫu

**Mở bài:** Trước sân trường em có một cây bàng đã nhiều năm tuổi.

**Thân bài (trích):** Thân cây to, vỏ sần sùi. Tán lá xòe rộng như một chiếc ô khổng lồ che mát cả góc sân. Mùa hè, lá bàng xanh thẫm; sang đông, lá chuyển sang màu đỏ thắm rồi rụng dần.

**Kết bài:** Em rất yêu cây bàng và mong cây mãi xanh tốt để che mát cho chúng em.

## 4. Bài tập mẫu

**Đề:** Viết một câu tả lá cây có dùng hình ảnh so sánh.

**Lời giải (gợi ý):** Lá sen to, tròn **như** chiếc nón úp ngược trên mặt hồ.

> [!NOTE]
> Bài văn hay không cần tả mọi thứ, mà cần chọn những nét nổi bật, riêng biệt của cây em tả.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Đọc hiểu truyện, từ loại và văn tả cây cối",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Từ nào là danh từ riêng?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "“Hà Nội” là tên riêng của một thành phố nên là danh từ riêng và phải viết hoa.",
					Options: []sampleOption{
						{Content: "thành phố"},
						{Content: "dòng sông"},
						{Content: "Hà Nội", IsCorrect: true},
						{Content: "ngọn núi"},
					},
				},
				{
					Prompt: "Từ nào là động từ?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "“Đọc” chỉ hoạt động nên là động từ; “quyển sách” là danh từ, “xanh” và “chăm chỉ” là tính từ.",
					Options: []sampleOption{
						{Content: "quyển sách"},
						{Content: "xanh"},
						{Content: "chăm chỉ"},
						{Content: "đọc", IsCorrect: true},
					},
				},
				{
					Prompt: "Trong câu “Cây phượng nở hoa đỏ rực.”, từ nào là tính từ?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "“Đỏ rực” chỉ màu sắc của hoa nên là tính từ; “nở” là động từ, “cây phượng” và “hoa” là danh từ.",
					Options: []sampleOption{
						{Content: "đỏ rực", IsCorrect: true},
						{Content: "nở"},
						{Content: "cây phượng"},
						{Content: "hoa"},
					},
				},
				{
					Prompt: "Trong truyện “Ăn khế trả vàng”, vì sao người em được chim trả vàng mà vẫn bình an?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Người em hiền lành, không tham, làm đúng lời chim dặn nên được đền đáp và bình an.",
					Options: []sampleOption{
						{Content: "Vì người em giàu có nhất làng."},
						{Content: "Vì người em hiền lành, không tham, chỉ may túi ba gang như lời chim dặn.", IsCorrect: true},
						{Content: "Vì người em may được chiếc túi rất to."},
						{Content: "Vì người anh nhường cây khế cho em."},
					},
				},
				{
					Prompt: "Khi viết bài văn tả cây bàng, câu nào phù hợp để làm mở bài?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Mở bài cần giới thiệu cây định tả và nơi cây mọc; các câu khác tả bộ phận hoặc nêu tình cảm.",
					Options: []sampleOption{
						{Content: "Thân cây to, vỏ sần sùi."},
						{Content: "Mùa đông, lá bàng chuyển sang màu đỏ thắm."},
						{Content: "Em mong cây bàng mãi xanh tốt."},
						{Content: "Trước sân trường em có một cây bàng đã nhiều năm tuổi.", IsCorrect: true},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGVIET5",
		Title:       "Tiếng Việt lớp 5",
		Description: "Lớp học mẫu môn Tiếng Việt lớp 5 theo Chương trình GDPT 2018: đọc hiểu ca dao về quê hương, từ đồng nghĩa, câu ghép – kết từ và viết bài văn tả phong cảnh.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Quê hương tươi đẹp",
				Lessons: []sampleLesson{
					{
						Title:           "Đọc hiểu: Vẻ đẹp quê hương trong ca dao",
						DurationMinutes: 20,
						Body: `## 1. Đọc các câu ca dao

Đường vô xứ Nghệ quanh quanh,
Non xanh nước biếc như tranh họa đồ.

Gió đưa cành trúc la đà,
Tiếng chuông Trấn Vũ, canh gà Thọ Xương.

## 2. Tìm hiểu nội dung

**Câu thứ nhất:**

- Con đường vào xứ Nghệ "quanh quanh", uốn lượn.
- Cảnh "non xanh nước biếc" được so sánh **như tranh họa đồ**: đẹp như một bức tranh vẽ.

**Câu thứ hai:**

- Gợi cảnh kinh thành Thăng Long xưa lúc sớm mai: cành trúc đung đưa trong gió.
- Âm thanh tiếng chuông đền Trấn Vũ, tiếng gà gáy ở Thọ Xương tạo không khí yên bình.

## 3. Tình cảm của tác giả dân gian

Qua ca dao, người xưa thể hiện **niềm tự hào và tình yêu** đối với cảnh đẹp quê hương, đất nước.

## 4. Bài tập mẫu

**Đề:** Trong câu "Non xanh nước biếc như tranh họa đồ", tác giả dùng biện pháp gì? Tác dụng?

**Lời giải:** Biện pháp **so sánh** (từ "như"): cảnh núi sông xứ Nghệ được ví với bức tranh, làm nổi bật vẻ đẹp hài hòa, thơ mộng.

> [!TIP]
> Khi đọc ca dao, hãy chú ý từ ngữ gợi hình (xanh, biếc, la đà) và gợi âm thanh (tiếng chuông, canh gà) để cảm nhận bức tranh quê hương.`,
					},
					{
						Title:           "Từ đồng nghĩa",
						DurationMinutes: 20,
						Body: `## 1. Từ đồng nghĩa là gì?

**Từ đồng nghĩa** là những từ có nghĩa giống nhau hoặc gần giống nhau.

Ví dụ: chăm chỉ – siêng năng – cần cù; đất nước – Tổ quốc – non sông.

## 2. Hai loại từ đồng nghĩa

- **Đồng nghĩa hoàn toàn:** có thể thay thế cho nhau trong lời nói. Ví dụ: quả – trái, lợn – heo, máy bay – tàu bay.
- **Đồng nghĩa không hoàn toàn:** nghĩa gần giống nhưng khác nhau về sắc thái, cần cân nhắc khi dùng. Ví dụ: ăn – xơi – chén; chết – mất – hi sinh.

## 3. Chọn từ cho phù hợp

- Mời ông bà dùng bữa: "Mời ông bà **xơi** cơm ạ." (lễ phép, kính trọng)
- Nói về người lính ngã xuống vì Tổ quốc: "Anh đã **hi sinh**." (trang trọng)

## 4. Bài tập mẫu

**Đề:** Tìm từ đồng nghĩa với "xanh" trong câu: "Non xanh nước biếc như tranh họa đồ."

**Lời giải:** Từ **biếc** cũng chỉ màu xanh (xanh biếc), đồng nghĩa không hoàn toàn với "xanh". Dùng hai từ khác nhau giúp câu thơ không bị lặp và giàu hình ảnh.

> [!NOTE]
> Dùng từ đồng nghĩa hợp lí giúp tránh lặp từ và thể hiện đúng thái độ, tình cảm của người nói.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Câu và bài văn",
				Lessons: []sampleLesson{
					{
						Title:           "Câu ghép và kết từ",
						DurationMinutes: 20,
						Body: `## 1. Câu ghép là gì?

**Câu ghép** là câu do hai vế câu trở lên ghép lại. Mỗi vế có cấu tạo giống một câu đơn (có chủ ngữ và vị ngữ).

Ví dụ: **Trời** / mưa to **nên** **đường** / rất trơn.

- Vế 1: Trời (CN) – mưa to (VN).
- Vế 2: đường (CN) – rất trơn (VN).

## 2. Cách nối các vế câu ghép

- Dùng **kết từ**: và, rồi, nhưng, còn, hoặc.
- Dùng **cặp kết từ**: vì... nên..., nếu... thì..., tuy... nhưng...
- Dùng **dấu câu**: dấu phẩy, dấu chấm phẩy.

## 3. Ví dụ

- Vì em chăm học **nên** em đạt kết quả tốt. (nguyên nhân – kết quả)
- **Nếu** trời nắng **thì** chúng em đi cắm trại. (điều kiện – kết quả)
- **Tuy** nhà xa **nhưng** bạn Nam không bao giờ đi học muộn. (tương phản)

## 4. Bài tập mẫu

**Đề:** Điền kết từ thích hợp: "Em đã cố gắng nhiều ... kết quả chưa cao như mong muốn."

**Lời giải:** Hai vế có ý trái ngược nhau, nên dùng kết từ **nhưng**.

> [!TIP]
> Phân biệt câu ghép với câu đơn có chủ ngữ ghép: "Lan và Hoa đều học giỏi" chỉ có một vị ngữ chung nên là câu đơn.`,
					},
					{
						Title:           "Viết bài văn tả phong cảnh",
						DurationMinutes: 20,
						Body: `## 1. Bố cục bài văn

1. **Mở bài:** giới thiệu cảnh định tả (cảnh gì, ở đâu, vào lúc nào).
2. **Thân bài:** tả bao quát rồi tả từng phần của cảnh; tả sự thay đổi của cảnh theo thời gian; có cả hoạt động của con người, con vật.
3. **Kết bài:** nêu cảm nghĩ của em về cảnh.

## 2. Trình tự miêu tả

- **Theo không gian:** từ xa đến gần, từ trên xuống dưới.
- **Theo thời gian:** sáng sớm, buổi trưa, chiều tối.

## 3. Đoạn văn mẫu (tả cánh đồng lúc bình minh)

Mặt trời từ từ nhô lên sau lũy tre làng. Cánh đồng lúa trải rộng mênh mông, những bông lúa chín vàng óng như được dát vàng. Gió nhẹ thổi qua, cả biển lúa nhấp nhô gợn sóng. Trên bờ ruộng, mấy chú cò trắng thong thả tìm mồi.

## 4. Bài tập mẫu

**Đề:** Chỉ ra một hình ảnh so sánh và một từ gợi tả trong đoạn văn mẫu.

**Lời giải:** So sánh: bông lúa chín vàng óng **như được dát vàng**. Từ gợi tả: **mênh mông** (gợi sự rộng lớn), **nhấp nhô** (gợi chuyển động của lúa).

> [!NOTE]
> Hãy dùng nhiều giác quan khi quan sát: màu sắc, âm thanh, mùi hương... giúp bài văn tả cảnh sinh động và chân thật.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Ca dao quê hương, từ đồng nghĩa và câu ghép",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Từ nào đồng nghĩa với “chăm chỉ”?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "“Siêng năng” có nghĩa giống “chăm chỉ”: chịu khó, làm việc đều đặn.",
					Options: []sampleOption{
						{Content: "lười biếng"},
						{Content: "siêng năng", IsCorrect: true},
						{Content: "thông minh"},
						{Content: "vui vẻ"},
					},
				},
				{
					Prompt: "Câu ghép là câu như thế nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Câu ghép gồm từ hai vế trở lên, mỗi vế có chủ ngữ và vị ngữ riêng như một câu đơn.",
					Options: []sampleOption{
						{Content: "Câu chỉ có một chủ ngữ và một vị ngữ."},
						{Content: "Câu có nhiều tính từ."},
						{Content: "Câu luôn kết thúc bằng dấu chấm than."},
						{Content: "Câu gồm hai vế trở lên, mỗi vế có chủ ngữ và vị ngữ.", IsCorrect: true},
					},
				},
				{
					Prompt: "Trong câu ca dao “Đường vô xứ Nghệ quanh quanh / Non xanh nước biếc như tranh họa đồ”, cảnh xứ Nghệ được so sánh với gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Từ so sánh “như” ví cảnh non xanh nước biếc với bức tranh họa đồ, làm nổi bật vẻ đẹp của xứ Nghệ.",
					Options: []sampleOption{
						{Content: "Bức tranh họa đồ", IsCorrect: true},
						{Content: "Con đường quanh co"},
						{Content: "Tiếng chuông chùa"},
						{Content: "Cành trúc la đà"},
					},
				},
				{
					Prompt: "Câu nào là câu ghép?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Câu có hai vế: “Trời / mưa to” và “đường / rất trơn”, nối bằng kết từ “nên”, nên là câu ghép.",
					Options: []sampleOption{
						{Content: "Sáng nay, em đi học sớm."},
						{Content: "Những bông hoa trong vườn nở rộ."},
						{Content: "Trời mưa to nên đường rất trơn.", IsCorrect: true},
						{Content: "Bạn Minh và bạn Lan đều học giỏi."},
					},
				},
				{
					Prompt: "Chọn kết từ thích hợp: “Em đã cố gắng nhiều ... kết quả chưa cao như mong muốn.”",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Hai vế có ý trái ngược (cố gắng nhiều – kết quả chưa cao) nên dùng kết từ “nhưng”.",
					Options: []sampleOption{
						{Content: "nên"},
						{Content: "vì"},
						{Content: "và"},
						{Content: "nhưng", IsCorrect: true},
					},
				},
			},
		},
	},
	{
		Code:        "NGUVAN6",
		Title:       "Ngữ văn lớp 6",
		Description: "Lớp học mẫu môn Ngữ văn lớp 6 theo Chương trình GDPT 2018: truyện dân gian (truyền thuyết, cổ tích), thơ lục bát và từ đơn – từ phức.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Truyện dân gian Việt Nam",
				Lessons: []sampleLesson{
					{
						Title:           "Truyền thuyết: Thánh Gióng",
						DurationMinutes: 30,
						Body: `## 1. Đặc điểm của truyền thuyết

**Truyền thuyết** là loại truyện dân gian kể về các nhân vật, sự kiện có liên quan đến lịch sử, thường có yếu tố kì ảo. Qua đó, nhân dân thể hiện thái độ, cách đánh giá đối với nhân vật, sự kiện ấy. Truyền thuyết thường gắn với các dấu tích, lễ hội còn lưu lại đến nay.

## 2. Tóm tắt truyện Thánh Gióng

Vào đời Hùng Vương thứ sáu, ở làng Gióng có một bà mẹ ướm chân vào một vết chân lạ ngoài đồng rồi mang thai, sinh ra một cậu bé. Lên ba tuổi, cậu vẫn không biết nói, biết cười. Khi giặc Ân xâm lược, nghe sứ giả đi tìm người tài, cậu bỗng cất tiếng nói, xin vua làm ngựa sắt, roi sắt, áo giáp sắt để đánh giặc.

Từ đó, Gióng lớn nhanh như thổi; bà con làng xóm góp gạo nuôi cậu. Gióng vươn vai thành tráng sĩ, cưỡi ngựa sắt xông ra trận. Roi sắt gãy, Gióng nhổ những bụi tre bên đường quật vào giặc. Giặc tan, Gióng lên đỉnh núi Sóc, cởi áo giáp sắt, rồi cả người lẫn ngựa bay về trời. Vua phong Gióng là Phù Đổng Thiên Vương.

## 3. Ý nghĩa một số chi tiết

| Chi tiết | Ý nghĩa |
|---|---|
| Tiếng nói đầu tiên là xin đi đánh giặc | Lòng yêu nước luôn thường trực trong mỗi người dân |
| Bà con góp gạo nuôi Gióng | Sức mạnh của Gióng là sức mạnh của cả cộng đồng |
| Nhổ tre đánh giặc | Cây cỏ quê hương cũng thành vũ khí chống giặc |
| Gióng bay về trời | Người anh hùng trở thành bất tử trong lòng nhân dân |

## 4. Dấu tích và lễ hội

Truyện gắn với các dấu tích như tre đằng ngà, những ao hồ liên tiếp (tương truyền là vết chân ngựa). Hội Gióng ở đền Phù Đổng và đền Sóc (Hà Nội) được UNESCO ghi danh là Di sản văn hóa phi vật thể đại diện của nhân loại năm 2010.

## 5. Bài tập mẫu

**Đề:** Chỉ ra hai yếu tố kì ảo trong truyện.

**Lời giải:** Bà mẹ ướm chân vào vết chân lạ rồi mang thai; Gióng vươn vai thành tráng sĩ; Gióng cùng ngựa sắt bay về trời (chọn hai trong số đó).

> [!NOTE]
> Yếu tố kì ảo trong truyền thuyết không làm câu chuyện xa rời lịch sử, mà giúp tôn vinh, lí tưởng hóa người anh hùng.`,
					},
					{
						Title:           "Truyện cổ tích: Thạch Sanh",
						DurationMinutes: 30,
						Body: `## 1. Đặc điểm của truyện cổ tích

**Truyện cổ tích** là loại truyện dân gian kể về cuộc đời của một số kiểu nhân vật quen thuộc: người mồ côi, người em, người dũng sĩ, người thông minh... Truyện thường có yếu tố kì ảo và thể hiện ước mơ của nhân dân về công lí: cái thiện thắng cái ác, "ở hiền gặp lành".

## 2. Tóm tắt truyện Thạch Sanh

Thạch Sanh mồ côi, sống một mình dưới gốc đa, hằng ngày đốn củi kiếm sống, được thiên thần dạy cho võ nghệ và phép thần thông. Lí Thông, một người bán rượu, kết nghĩa anh em với Thạch Sanh để lợi dụng chàng.

Thạch Sanh lần lượt giết chằn tinh, bắn đại bàng cứu công chúa, cứu thái tử con vua Thủy Tề và được tặng cây đàn thần. Nhưng chàng nhiều lần bị Lí Thông cướp công, lại bị vu oan phải vào ngục. Tiếng đàn của chàng giúp công chúa khỏi câm, sự thật được làm sáng tỏ.

Về sau, quân mười tám nước chư hầu kéo sang đánh. Thạch Sanh gảy đàn làm quân giặc mất hết tinh thần, rồi đem niêu cơm nhỏ ra đãi, quân lính ăn mãi không hết. Thạch Sanh lấy công chúa và được nối ngôi vua.

## 3. Phân tích nhân vật

- **Thạch Sanh:** thật thà, dũng cảm, giàu lòng nhân ái, yêu hòa bình.
- **Lí Thông:** gian xảo, ích kỉ, vong ân bội nghĩa.

## 4. Ý nghĩa vật thần

- **Tiếng đàn thần:** giúp nhận ra sự thật, đòi lại công lí; làm quân giặc chịu hàng, thể hiện tinh thần yêu hòa bình.
- **Niêu cơm thần:** ăn mãi không hết, thể hiện lòng nhân đạo, sự bao dung với kẻ thua trận.

## 5. Bài tập mẫu

**Đề:** Chỉ ra một chi tiết kì ảo và nêu ý nghĩa.

**Lời giải:** Niêu cơm nhỏ mà quân lính mười tám nước ăn mãi không hết: thể hiện tấm lòng nhân hậu, yêu hòa bình của Thạch Sanh và của nhân dân.

> [!TIP]
> Khi đọc truyện cổ tích, hãy so sánh các cặp nhân vật đối lập (thiện – ác) để hiểu rõ bài học mà tác giả dân gian gửi gắm.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Thơ lục bát và từ tiếng Việt",
				Lessons: []sampleLesson{
					{
						Title:           "Thơ lục bát qua ca dao",
						DurationMinutes: 30,
						Body: `## 1. Thể thơ lục bát

**Lục bát** là thể thơ truyền thống của dân tộc, gồm các cặp câu: câu sáu tiếng (câu lục) và câu tám tiếng (câu bát).

## 2. Luật thơ lục bát

- **Vần:** tiếng thứ 6 của câu lục vần với tiếng thứ 6 của câu bát; tiếng thứ 8 của câu bát vần với tiếng thứ 6 của câu lục tiếp theo.
- **Thanh:** các tiếng thứ 2, 4, 6 của câu lục lần lượt là bằng – trắc – bằng; câu bát: tiếng 2 bằng, tiếng 4 trắc, tiếng 6 và 8 đều bằng (thường một tiếng thanh ngang, một tiếng thanh huyền).
- **Nhịp:** thường là nhịp chẵn (2/2/2, 4/4).

## 3. Phân tích ví dụ

Anh đi anh nhớ quê nhà,
Nhớ canh rau muống, nhớ cà dầm tương.

- Câu lục: "nhà" (tiếng 6) vần với "cà" (tiếng 6 câu bát) – vần "a".
- Thanh câu lục: đi (B) – nhớ (T) – nhà (B).
- Thanh câu bát: canh (B) – muống (T) – cà (B, huyền) – tương (B, ngang).
- Nội dung: nỗi nhớ quê hương gắn với những món ăn bình dị.

## 4. Một ví dụ khác

Trong đầm gì đẹp bằng sen,
Lá xanh bông trắng lại chen nhị vàng.

Tiếng "sen" (tiếng 6 câu lục) vần với tiếng "chen" (tiếng 6 câu bát).

## 5. Bài tập mẫu

**Đề:** Đếm số tiếng và chỉ ra cách ngắt nhịp câu "Lá xanh bông trắng lại chen nhị vàng".

**Lời giải:** Câu có **8 tiếng**, ngắt nhịp **2/2/2/2**: Lá xanh / bông trắng / lại chen / nhị vàng.

> [!TIP]
> Đọc to bài thơ lục bát sẽ giúp em cảm nhận rõ nhịp điệu nhẹ nhàng, êm ái gần với lời ru, lời hát dân gian.`,
					},
					{
						Title:           "Từ đơn và từ phức",
						DurationMinutes: 25,
						Body: `## 1. Từ đơn

**Từ đơn** là từ chỉ gồm **một tiếng**.

Ví dụ: nhà, sông, ăn, đẹp, mẹ.

## 2. Từ phức

**Từ phức** là từ gồm **hai tiếng trở lên**. Từ phức gồm hai loại:

- **Từ ghép:** các tiếng có quan hệ với nhau về nghĩa. Ví dụ: sách vở, bàn ghế, xe đạp, quê hương.
- **Từ láy:** các tiếng có quan hệ láy âm (lặp lại âm đầu, vần hoặc cả tiếng). Ví dụ: lom khom, xinh xắn, rì rào, long lanh, xanh xanh.

## 3. Phân tích ví dụ

Câu: "Mẹ mua cho em một chiếc xe đạp xinh xắn."

| Loại từ | Các từ |
|---|---|
| Từ đơn | mẹ, mua, cho, em, một, chiếc |
| Từ ghép | xe đạp |
| Từ láy | xinh xắn |

## 4. Tác dụng của từ láy

Từ láy giúp gợi tả hình ảnh, âm thanh cụ thể, sinh động:

- "rì rào" gợi âm thanh của gió, của sóng.
- "lom khom" gợi dáng người cúi xuống.

## 5. Bài tập mẫu

**Đề:** Xếp các từ sau vào nhóm từ ghép hoặc từ láy: quê hương, long lanh, bàn ghế, rì rào.

**Lời giải:** Từ ghép: **quê hương, bàn ghế**. Từ láy: **long lanh, rì rào**.

> [!NOTE]
> Không phải từ nào có hai tiếng giống âm đầu cũng là từ láy: "tươi tốt" có hai tiếng đều có nghĩa nên thường được xếp vào từ ghép.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Truyện dân gian, thơ lục bát, từ đơn và từ phức",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Truyện “Thánh Gióng” thuộc thể loại nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Thánh Gióng kể về nhân vật, sự kiện có liên quan đến lịch sử (chống giặc Ân), có yếu tố kì ảo và gắn với dấu tích, lễ hội nên là truyền thuyết.",
					Options: []sampleOption{
						{Content: "Truyện cổ tích"},
						{Content: "Truyền thuyết", IsCorrect: true},
						{Content: "Truyện ngụ ngôn"},
						{Content: "Truyện cười"},
					},
				},
				{
					Prompt: "Mỗi cặp câu thơ lục bát gồm những câu nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Lục bát gồm câu lục (6 tiếng) và câu bát (8 tiếng).",
					Options: []sampleOption{
						{Content: "Câu 5 tiếng và câu 7 tiếng"},
						{Content: "Câu 7 tiếng và câu 8 tiếng"},
						{Content: "Câu 6 tiếng và câu 8 tiếng", IsCorrect: true},
						{Content: "Hai câu đều 7 tiếng"},
					},
				},
				{
					Prompt: "Chi tiết nào trong truyện “Thạch Sanh” là yếu tố kì ảo?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Niêu cơm nhỏ mà ăn mãi không hết là điều không có thật ngoài đời, nên là chi tiết kì ảo.",
					Options: []sampleOption{
						{Content: "Thạch Sanh sống một mình dưới gốc đa."},
						{Content: "Thạch Sanh hằng ngày đi đốn củi kiếm sống."},
						{Content: "Lí Thông kết nghĩa anh em với Thạch Sanh."},
						{Content: "Niêu cơm nhỏ mà quân lính mười tám nước chư hầu ăn mãi không hết.", IsCorrect: true},
					},
				},
				{
					Prompt: "Từ nào là từ láy?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "“Lom khom” có hai tiếng láy vần “om”, gợi tả dáng người cúi xuống; các từ còn lại có các tiếng quan hệ về nghĩa nên là từ ghép.",
					Options: []sampleOption{
						{Content: "lom khom", IsCorrect: true},
						{Content: "sách vở"},
						{Content: "bàn ghế"},
						{Content: "xe đạp"},
					},
				},
				{
					Prompt: "Trong câu ca dao “Trong đầm gì đẹp bằng sen / Lá xanh bông trắng lại chen nhị vàng”, tiếng nào ở câu bát vần với tiếng “sen”?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Theo luật lục bát, tiếng thứ 6 câu lục (“sen”) vần với tiếng thứ 6 câu bát (“chen”), cùng vần “en”.",
					Options: []sampleOption{
						{Content: "xanh"},
						{Content: "trắng"},
						{Content: "chen", IsCorrect: true},
						{Content: "vàng"},
					},
				},
			},
		},
	},
}
