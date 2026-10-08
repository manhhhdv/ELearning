package store

// sampleAnhPrimary là các lớp học mẫu môn Tiếng Anh từ lớp 1 đến lớp 6 theo Chương trình GDPT 2018.
// Lớp 1–2 là giai đoạn làm quen (tự chọn); bài giảng giải thích bằng tiếng Việt, ví dụ bằng tiếng Anh kèm nghĩa.
var sampleAnhPrimary = []sampleCourse{
	{
		Code:        "TIENGANH1",
		Title:       "Tiếng Anh lớp 1",
		Description: "Lớp học mẫu làm quen môn Tiếng Anh lớp 1 (tự chọn) theo Chương trình GDPT 2018: chào hỏi, hỏi tên, đồ dùng học tập và màu sắc.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Chào hỏi",
				Lessons: []sampleLesson{
					{
						Title:           "Hello! – Xin chào!",
						DurationMinutes: 15,
						Body: `## 1. Em chào bạn

Khi gặp bạn, em nói: **Hello!** hoặc **Hi!**

Khi chia tay, em nói: **Goodbye!** hoặc **Bye!**

## 2. Từ mới

| Tiếng Anh | Nghĩa |
|---|---|
| Hello / Hi | Xin chào |
| Goodbye / Bye | Tạm biệt |
| I'm ... | Tớ là ... |
| teacher | cô giáo, thầy giáo |

## 3. Hội thoại mẫu

- **Nam:** Hello! I'm Nam. *(Xin chào! Tớ là Nam.)*
- **Mai:** Hi, Nam! I'm Mai. *(Chào Nam! Tớ là Mai.)*

Hết giờ học:

- **Nam:** Goodbye, Mai! *(Tạm biệt Mai!)*
- **Mai:** Bye, Nam! *(Tạm biệt Nam!)*

## 4. Em tập nói

Quay sang bạn bên cạnh. Em nói:

- Hello! I'm ... *(nói tên em)*.
- Goodbye!

> [!TIP]
> Khi chào, em nhìn bạn và mỉm cười nhé!`,
					},
					{
						Title:           "What's your name? – Bạn tên là gì?",
						DurationMinutes: 15,
						Body: `## 1. Em hỏi tên bạn

Muốn biết tên bạn, em hỏi:

- **What's your name?** *(Bạn tên là gì?)*

Bạn trả lời:

- **My name's Lan.** *(Tên tớ là Lan.)*
- Hoặc ngắn hơn: **I'm Lan.**

## 2. Từ mới

| Tiếng Anh | Nghĩa |
|---|---|
| name | tên |
| my | của tớ |
| your | của bạn |
| Nice to meet you. | Rất vui được gặp bạn. |
| Thank you. | Cảm ơn. |

## 3. Hội thoại mẫu

- **Lan:** Hello! What's your name? *(Xin chào! Bạn tên là gì?)*
- **Bin:** My name's Bin. *(Tên tớ là Bin.)*
- **Lan:** Nice to meet you, Bin. *(Rất vui được gặp bạn, Bin.)*
- **Bin:** Nice to meet you, too. *(Tớ cũng rất vui được gặp bạn.)*

## 4. Em tập nói

Em hỏi tên ba bạn trong lớp: "What's your name?".

> [!NOTE]
> "My name's" là cách nói ngắn của "My name is".`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Đồ dùng học tập và màu sắc",
				Lessons: []sampleLesson{
					{
						Title:           "School things – Đồ dùng học tập",
						DurationMinutes: 15,
						Body: `## 1. Từ mới

| Tiếng Anh | Nghĩa |
|---|---|
| pen | bút mực |
| pencil | bút chì |
| book | quyển sách |
| ruler | cái thước |
| bag | cái cặp |
| desk | cái bàn học |

## 2. Hỏi và trả lời

Cô chỉ vào một đồ vật và hỏi:

- **What's this?** *(Đây là cái gì?)*

Em trả lời:

- **It's a book.** *(Đó là quyển sách.)*

## 3. Ví dụ

- What's this? – It's a pen. *(Đó là cái bút mực.)*
- What's this? – It's a ruler. *(Đó là cái thước.)*
- What's this? – It's a bag. *(Đó là cái cặp.)*

## 4. Trò chơi

Em lấy một đồ dùng trong cặp ra. Bạn hỏi: "What's this?". Em trả lời: "It's a ...".

> [!TIP]
> Mỗi ngày, em cầm từng đồ dùng học tập lên và đọc to tên tiếng Anh của nó.`,
					},
					{
						Title:           "Colours – Màu sắc",
						DurationMinutes: 15,
						Body: `## 1. Từ mới

| Tiếng Anh | Nghĩa |
|---|---|
| red | màu đỏ |
| blue | màu xanh da trời |
| green | màu xanh lá cây |
| yellow | màu vàng |
| black | màu đen |
| white | màu trắng |
| pink | màu hồng |

## 2. Hỏi và trả lời

- **What colour is it?** *(Nó màu gì?)*
- **It's red.** *(Nó màu đỏ.)*

## 3. Ví dụ

- Quả táo: What colour is it? – It's red.
- Lá cây: What colour is it? – It's green.
- Bút chì vàng: What colour is it? – It's yellow.

## 4. Em tập nói

Em chỉ vào đồ vật trong lớp và hỏi bạn: "What colour is it?".

> [!TIP]
> Em có thể vừa tô màu tranh vừa đọc to tên màu bằng tiếng Anh.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Chào hỏi, đồ dùng học tập và màu sắc",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Khi gặp bạn ở trường, em chào bằng tiếng Anh thế nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Hello! nghĩa là Xin chào!, dùng khi gặp nhau; Goodbye! và Bye! dùng khi chia tay.",
					Options: []sampleOption{
						{Content: "Goodbye!"},
						{Content: "Hello!", IsCorrect: true},
						{Content: "Bye!"},
						{Content: "Thank you."},
					},
				},
				{
					Prompt: "Từ 'red' chỉ màu gì?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "red là màu đỏ; green là xanh lá cây, yellow là vàng, white là trắng.",
					Options: []sampleOption{
						{Content: "Màu xanh lá cây"},
						{Content: "Màu vàng"},
						{Content: "Màu trắng"},
						{Content: "Màu đỏ", IsCorrect: true},
					},
				},
				{
					Prompt: "Bạn hỏi em: 'What's your name?'. Em trả lời thế nào?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "What's your name? là hỏi tên, nên em trả lời bằng tên của mình: My name's Minh.",
					Options: []sampleOption{
						{Content: "My name's Minh.", IsCorrect: true},
						{Content: "It's a pen."},
						{Content: "Goodbye!"},
						{Content: "It's blue."},
					},
				},
				{
					Prompt: "Cô chỉ vào quyển sách và hỏi: 'What's this?'. Em trả lời: 'It's a ___.'",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "book nghĩa là quyển sách; pen là bút mực, ruler là thước, bag là cặp.",
					Options: []sampleOption{
						{Content: "pen"},
						{Content: "ruler"},
						{Content: "book", IsCorrect: true},
						{Content: "bag"},
					},
				},
				{
					Prompt: "Mai cầm một chiếc bút chì màu vàng. Bạn hỏi: 'What colour is it?'. Mai trả lời đúng là:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu hỏi What colour is it? hỏi về màu sắc; bút chì màu vàng nên trả lời It's yellow.",
					Options: []sampleOption{
						{Content: "It's a pencil."},
						{Content: "It's yellow.", IsCorrect: true},
						{Content: "It's green."},
						{Content: "I'm Mai."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH2",
		Title:       "Tiếng Anh lớp 2",
		Description: "Lớp học mẫu làm quen môn Tiếng Anh lớp 2 (tự chọn) theo Chương trình GDPT 2018: gia đình, số đếm, con vật và đồ chơi.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Gia đình em",
				Lessons: []sampleLesson{
					{
						Title:           "My family – Gia đình em",
						DurationMinutes: 15,
						Body: `## 1. Từ mới

| Tiếng Anh | Nghĩa |
|---|---|
| father (dad) | bố |
| mother (mum) | mẹ |
| brother | anh trai, em trai |
| sister | chị gái, em gái |
| grandfather | ông |
| grandmother | bà |

## 2. Giới thiệu người thân

Em chỉ vào ảnh và nói:

- **This is my mother.** *(Đây là mẹ của em.)*
- **This is my brother.** *(Đây là anh trai của em.)*

## 3. Hỏi và trả lời

- **Who's this?** *(Đây là ai?)*
- **It's my grandfather.** *(Đó là ông của em.)*

## 4. Hội thoại mẫu

- **Lan:** Who's this? *(Đây là ai?)*
- **Nam:** It's my sister. *(Đó là chị gái tớ.)*
- **Lan:** And who's this? *(Còn đây là ai?)*
- **Nam:** It's my dad. *(Đó là bố tớ.)*

> [!TIP]
> Em mang ảnh gia đình đến lớp và giới thiệu từng người: "This is my ...".`,
					},
					{
						Title:           "How old are you? – Số đếm từ 1 đến 10",
						DurationMinutes: 15,
						Body: `## 1. Số đếm

| Số | Tiếng Anh | Số | Tiếng Anh |
|---|---|---|---|
| 1 | one | 6 | six |
| 2 | two | 7 | seven |
| 3 | three | 8 | eight |
| 4 | four | 9 | nine |
| 5 | five | 10 | ten |

## 2. Hỏi tuổi

- **How old are you?** *(Bạn mấy tuổi?)*
- **I'm seven.** *(Tớ bảy tuổi.)*

## 3. Hội thoại mẫu

- **Mai:** How old are you, Bin? *(Bin ơi, bạn mấy tuổi?)*
- **Bin:** I'm eight. And you? *(Tớ tám tuổi. Còn bạn?)*
- **Mai:** I'm seven. *(Tớ bảy tuổi.)*

## 4. Đếm cùng nhau

Em đếm ngón tay: one, two, three, four, five, six, seven, eight, nine, ten!

> [!NOTE]
> Khi nói tuổi, em chỉ cần nói "I'm" và số tuổi, ví dụ: I'm six.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Con vật và đồ chơi",
				Lessons: []sampleLesson{
					{
						Title:           "Animals – Con vật",
						DurationMinutes: 15,
						Body: `## 1. Từ mới

| Tiếng Anh | Nghĩa |
|---|---|
| cat | con mèo |
| dog | con chó |
| bird | con chim |
| fish | con cá |
| duck | con vịt |
| rabbit | con thỏ |

## 2. Hỏi và trả lời

- **What's this?** *(Đây là con gì?)*
- **It's a cat.** *(Đó là con mèo.)*

## 3. Đoán con vật

- **Is it a dog?** *(Có phải con chó không?)*
- **Yes, it is.** *(Đúng rồi.)*
- **No, it isn't.** *(Không phải.)*

## 4. Ví dụ

Bạn giấu tranh con vịt. Em đoán:

- Is it a bird? – No, it isn't.
- Is it a duck? – Yes, it is!

> [!TIP]
> Em có thể vừa bắt chước tiếng kêu của con vật vừa nói tên tiếng Anh của nó.`,
					},
					{
						Title:           "Toys – Đồ chơi",
						DurationMinutes: 15,
						Body: `## 1. Từ mới

| Tiếng Anh | Nghĩa |
|---|---|
| ball | quả bóng |
| doll | búp bê |
| car | ô tô đồ chơi |
| kite | cái diều |
| robot | rô-bốt |
| teddy bear | gấu bông |

## 2. Em có đồ chơi gì?

- **I have a kite.** *(Tớ có một cái diều.)*
- **I have a doll.** *(Tớ có một con búp bê.)*

## 3. Nhiều đồ chơi

Khi có từ 2 đồ chơi trở lên, em thêm **s** vào cuối từ:

- one ball → two **balls**
- one car → three **cars**

Ví dụ: I have **three balls**. *(Tớ có ba quả bóng.)*

## 4. Em tập nói

Em kể đồ chơi của mình: "I have a robot. I have two cars."

> [!NOTE]
> Một đồ chơi: "a ball". Hai đồ chơi trở lên: "two balls", "five balls".`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Gia đình, số đếm, con vật và đồ chơi",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Từ 'mother' có nghĩa là gì?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "mother (mum) nghĩa là mẹ; father là bố, brother là anh/em trai, grandmother là bà.",
					Options: []sampleOption{
						{Content: "Bố"},
						{Content: "Mẹ", IsCorrect: true},
						{Content: "Anh trai"},
						{Content: "Bà"},
					},
				},
				{
					Prompt: "Từ 'eight' là số mấy?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "eight là số 8; six là 6, seven là 7, nine là 9.",
					Options: []sampleOption{
						{Content: "6"},
						{Content: "7"},
						{Content: "8", IsCorrect: true},
						{Content: "9"},
					},
				},
				{
					Prompt: "Bạn hỏi: 'How old are you?'. Em bảy tuổi. Em trả lời:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "How old are you? là hỏi tuổi; bảy tuổi nói là I'm seven.",
					Options: []sampleOption{
						{Content: "I'm seven.", IsCorrect: true},
						{Content: "I'm Nam."},
						{Content: "It's a cat."},
						{Content: "I have seven."},
					},
				},
				{
					Prompt: "Từ nào là tên một con vật?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "dog là con chó; kite, ball, doll đều là đồ chơi.",
					Options: []sampleOption{
						{Content: "kite"},
						{Content: "ball"},
						{Content: "doll"},
						{Content: "dog", IsCorrect: true},
					},
				},
				{
					Prompt: "Em có ba quả bóng. Em nói câu nào đúng?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Ba quả bóng là nhiều đồ vật nên thêm s: three balls.",
					Options: []sampleOption{
						{Content: "I have three ball."},
						{Content: "I have three balls.", IsCorrect: true},
						{Content: "I have a balls."},
						{Content: "I have two balls."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH3",
		Title:       "Tiếng Anh lớp 3",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 3 theo Chương trình GDPT 2018: chào hỏi và giới thiệu bạn bè, đồ dùng học tập, mệnh lệnh trong lớp học.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Bản thân và bạn bè",
				Lessons: []sampleLesson{
					{
						Title:           "Hello! How are you?",
						DurationMinutes: 20,
						Body: `## 1. Chào hỏi và hỏi thăm

Khi gặp nhau, ngoài "Hello!", em có thể hỏi thăm sức khỏe:

- **How are you?** *(Bạn có khỏe không?)*
- **I'm fine, thank you. And you?** *(Tớ khỏe, cảm ơn. Còn bạn?)*
- **Fine, thanks.** *(Tớ khỏe, cảm ơn.)*

## 2. Hỏi tên và đánh vần tên

- **What's your name?** – **My name's Linh.**
- **How do you spell your name?** *(Bạn đánh vần tên mình thế nào?)* – **L-I-N-H.**

## 3. Từ vựng

| Tiếng Anh | Nghĩa |
|---|---|
| How are you? | Bạn có khỏe không? |
| fine | khỏe |
| thank you / thanks | cảm ơn |
| spell | đánh vần |
| Goodbye. See you later. | Tạm biệt. Hẹn gặp lại. |

## 4. Hội thoại mẫu

- **Linh:** Hi, Nam. How are you?
- **Nam:** I'm fine, thank you. And you?
- **Linh:** Fine, thanks.

*(Linh: Chào Nam. Bạn khỏe không? – Nam: Tớ khỏe, cảm ơn. Còn bạn? – Linh: Tớ khỏe, cảm ơn.)*

## 5. Bài tập mẫu

Sắp xếp thành câu đúng: **name / your / what's / ?**

Lời giải: **What's your name?**

> [!TIP]
> Hãy tập đọc thuộc bảng chữ cái tiếng Anh (A, B, C...) để đánh vần tên mình thật nhanh.`,
					},
					{
						Title:           "This is my friend",
						DurationMinutes: 20,
						Body: `## 1. Giới thiệu bạn

Để giới thiệu một người bạn đứng gần, em dùng **This is ...**; bạn đứng xa thì dùng **That is ...** (viết gọn: That's).

- **This is Mai.** *(Đây là Mai.)*
- **That's Tom.** *(Kia là Tom.)*

## 2. He và She

| Từ | Dùng cho | Ví dụ |
|---|---|---|
| he (he's = he is) | bạn nam | He's my friend. *(Bạn ấy là bạn của tớ.)* |
| she (she's = she is) | bạn nữ | She's my friend. |

## 3. Hỏi xem đó có phải là ai không

- **Is that Mai?** *(Kia có phải là Mai không?)*
- **Yes, it is.** *(Đúng vậy.)* / **No, it isn't.** *(Không phải.)*

## 4. Hỏi tuổi của bạn

- **How old is he?** – **He's eight.** *(Bạn ấy tám tuổi.)*
- **How old is she?** – **She's nine.**

## 5. Bài tập mẫu

Điền **He** hoặc **She**: "This is Lan. ___'s my friend."

Lời giải: Lan là bạn nữ nên điền **She**: "This is Lan. She's my friend."

> [!NOTE]
> Bạn nam dùng "he", bạn nữ dùng "she". Đừng nhầm lẫn hai từ này nhé!`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Trường học của em",
				Lessons: []sampleLesson{
					{
						Title:           "School things – Is this your pen?",
						DurationMinutes: 20,
						Body: `## 1. Từ vựng

| Tiếng Anh | Nghĩa |
|---|---|
| pen | bút mực |
| pencil case | hộp bút |
| notebook | vở |
| eraser | cục tẩy |
| school bag | cặp sách |
| ruler | thước kẻ |

## 2. "a" hay "an"?

- Dùng **an** trước từ bắt đầu bằng nguyên âm a, e, i, o, u: **an eraser**, **an apple**.
- Dùng **a** trước các từ còn lại: **a pen**, **a ruler**, **a notebook**.

## 3. Hỏi đồ vật có phải của bạn không

- **Is this your pen?** *(Đây có phải bút của bạn không?)*
- **Yes, it is.** *(Đúng vậy.)* / **No, it isn't.** *(Không phải.)*

## 4. Hội thoại mẫu

- **Mai:** Is this your notebook, Nam?
- **Nam:** No, it isn't. It's Lan's notebook.
- **Mai:** Is that your school bag?
- **Nam:** Yes, it is.

## 5. Bài tập mẫu

Điền **a** hoặc **an**: ___ eraser; ___ ruler.

Lời giải: **an** eraser (eraser bắt đầu bằng nguyên âm e); **a** ruler.

> [!TIP]
> Khi trả lời ngắn câu hỏi "Is this ...?", em dùng "Yes, it is." hoặc "No, it isn't.".`,
					},
					{
						Title:           "Classroom instructions – Mệnh lệnh trong lớp",
						DurationMinutes: 20,
						Body: `## 1. Câu mệnh lệnh

Câu mệnh lệnh bắt đầu bằng **động từ**, thêm **please** cho lịch sự.

| Tiếng Anh | Nghĩa |
|---|---|
| Stand up, please. | Mời các em đứng lên. |
| Sit down, please. | Mời các em ngồi xuống. |
| Open your book, please. | Các em mở sách ra. |
| Close your book, please. | Các em gấp sách lại. |
| Be quiet, please. | Các em giữ trật tự. |
| Don't talk! | Không nói chuyện! |

## 2. Xin phép thầy cô

- **May I come in?** *(Em vào lớp được không ạ?)*
- **May I go out?** *(Em ra ngoài được không ạ?)*

Thầy cô trả lời:

- **Yes, you can.** *(Được, em có thể.)*
- **No, you can't.** *(Không, em không được.)*

## 3. Hội thoại mẫu

- **Nam:** May I come in, Miss Hoa?
- **Cô Hoa:** Yes, you can. Sit down, please.
- **Nam:** Thank you.

## 4. Bài tập mẫu

Cô muốn cả lớp gấp sách lại. Cô nói câu nào?

Lời giải: **Close your book, please.**

> [!NOTE]
> Muốn bảo ai đó KHÔNG làm gì, em thêm "Don't" trước động từ: Don't run! (Đừng chạy!)`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Bạn bè, đồ dùng học tập và lớp học",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Chọn câu trả lời phù hợp cho câu hỏi: 'How are you?'",
					Level:  "Nhận biết", Points: 2,
					Explanation: "How are you? là hỏi thăm sức khỏe, trả lời là I'm fine, thank you.",
					Options: []sampleOption{
						{Content: "I'm eight."},
						{Content: "I'm fine, thank you.", IsCorrect: true},
						{Content: "My name's Nam."},
						{Content: "Yes, it is."},
					},
				},
				{
					Prompt: "Câu mệnh lệnh 'Sit down, please.' có nghĩa là gì?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "sit down nghĩa là ngồi xuống; stand up mới là đứng lên.",
					Options: []sampleOption{
						{Content: "Mời em đứng lên."},
						{Content: "Mời em mở sách ra."},
						{Content: "Mời em ngồi xuống.", IsCorrect: true},
						{Content: "Mời em gấp sách lại."},
					},
				},
				{
					Prompt: "Chọn từ đúng: 'This is ___ eraser.'",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "eraser bắt đầu bằng nguyên âm e nên dùng an: an eraser.",
					Options: []sampleOption{
						{Content: "an", IsCorrect: true},
						{Content: "a"},
						{Content: "two"},
						{Content: "is"},
					},
				},
				{
					Prompt: "Chọn câu đúng để giới thiệu bạn nữ tên Mai:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Mai là bạn nữ nên dùng she: This is Mai. She's my friend.",
					Options: []sampleOption{
						{Content: "This is Mai. He's my friend."},
						{Content: "That are Mai. She's my friend."},
						{Content: "This Mai is. She's my friend."},
						{Content: "This is Mai. She's my friend.", IsCorrect: true},
					},
				},
				{
					Prompt: "Bạn hỏi: 'Is this your ruler?'. Chiếc thước đó không phải của em. Em trả lời thế nào?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu hỏi Is this ...? trả lời ngắn bằng Yes, it is. hoặc No, it isn't.; thước không phải của em nên chọn No, it isn't.",
					Options: []sampleOption{
						{Content: "Yes, it is."},
						{Content: "No, it isn't.", IsCorrect: true},
						{Content: "No, I'm not."},
						{Content: "Yes, you can."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH4",
		Title:       "Tiếng Anh lớp 4",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 4 theo Chương trình GDPT 2018: quốc gia và quốc tịch, ngày tháng sinh nhật, giờ giấc và khả năng với can.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Bạn bè quốc tế",
				Lessons: []sampleLesson{
					{
						Title:           "Where are you from?",
						DurationMinutes: 20,
						Body: `## 1. Hỏi bạn đến từ đâu

- **Where are you from?** *(Bạn đến từ đâu?)*
- **I'm from Viet Nam.** *(Tớ đến từ Việt Nam.)*

## 2. Hỏi quốc tịch

- **What nationality are you?** *(Bạn mang quốc tịch gì?)*
- **I'm Vietnamese.** *(Tớ là người Việt Nam.)*

## 3. Bảng từ vựng: quốc gia và quốc tịch

| Quốc gia | Quốc tịch | Nghĩa |
|---|---|---|
| Viet Nam | Vietnamese | Việt Nam |
| Japan | Japanese | Nhật Bản |
| Malaysia | Malaysian | Ma-lai-xi-a |
| Australia | Australian | Ô-xtrây-li-a |
| America | American | Mỹ |
| England | English | Anh |

## 4. Hỏi về người khác

- **Where is she from?** – **She's from Japan. She's Japanese.**
- **Where is he from?** – **He's from Australia. He's Australian.**

## 5. Bài tập mẫu

Hoàn thành câu: "Tom is from America. He's ___."

Lời giải: Tom đến từ Mỹ (America) nên quốc tịch là **American**.

> [!NOTE]
> Tên quốc gia và quốc tịch trong tiếng Anh luôn viết hoa chữ cái đầu: Japan, Japanese.`,
					},
					{
						Title:           "When's your birthday?",
						DurationMinutes: 20,
						Body: `## 1. Các tháng trong năm

January, February, March, April, May, June, July, August, September, October, November, December.

## 2. Số thứ tự

| Số | Viết tắt | Tiếng Anh |
|---|---|---|
| 1 | 1st | first |
| 2 | 2nd | second |
| 3 | 3rd | third |
| 4 | 4th | fourth |
| 5 | 5th | fifth |
| 8 | 8th | eighth |
| 12 | 12th | twelfth |
| 20 | 20th | twentieth |
| 21 | 21st | twenty-first |

## 3. Hỏi ngày sinh nhật

- **When's your birthday?** *(Sinh nhật bạn vào khi nào?)*
- **It's in May.** *(Vào tháng Năm.)*
- **It's on the fifth of May.** *(Vào ngày 5 tháng Năm.)*

## 4. Giới từ "in" và "on"

- **in** + tháng: in July, in December.
- **on** + ngày cụ thể: on the first of June.

## 5. Bài tập mẫu

Điền giới từ: "My birthday is ___ the twelfth of October."

Lời giải: có ngày cụ thể nên dùng **on**.

> [!TIP]
> Số thứ tự thường thêm "th", nhưng cần nhớ riêng: first, second, third, fifth, twelfth.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Thời gian và hoạt động",
				Lessons: []sampleLesson{
					{
						Title:           "What time is it?",
						DurationMinutes: 20,
						Body: `## 1. Hỏi giờ

- **What time is it?** *(Mấy giờ rồi?)*
- **It's seven o'clock.** *(Bảy giờ đúng.)*

## 2. Cách đọc giờ

Đọc **số giờ** trước, rồi **số phút** sau:

| Giờ | Cách đọc |
|---|---|
| 7:00 | seven o'clock |
| 7:15 | seven fifteen |
| 7:30 | seven thirty |
| 9:45 | nine forty-five |

## 3. Hỏi giờ làm việc gì đó

- **What time do you get up?** *(Bạn thức dậy lúc mấy giờ?)*
- **I get up at six o'clock.** *(Tớ thức dậy lúc 6 giờ.)*

Dùng **at** trước giờ: at six o'clock, at seven thirty.

## 4. Từ vựng hoạt động hằng ngày

| Tiếng Anh | Nghĩa |
|---|---|
| get up | thức dậy |
| have breakfast | ăn sáng |
| go to school | đi học |
| go to bed | đi ngủ |

## 5. Bài tập mẫu

Đồng hồ chỉ 6:30. Em nói thế nào?

Lời giải: **It's six thirty.**

> [!NOTE]
> "o'clock" chỉ dùng cho giờ đúng: 8:00 là eight o'clock, còn 8:15 là eight fifteen.`,
					},
					{
						Title:           "What can you do?",
						DurationMinutes: 20,
						Body: `## 1. Nói về khả năng với "can"

Cấu trúc: **chủ ngữ + can + động từ nguyên mẫu**

- **I can swim.** *(Tớ biết bơi.)*
- **She can dance.** *(Bạn ấy biết nhảy.)*

Phủ định: **can't** (= cannot)

- **He can't sing.** *(Bạn ấy không biết hát.)*

## 2. Hỏi và trả lời

- **What can you do?** – **I can ride a bike.** *(Tớ biết đi xe đạp.)*
- **Can you skate?** – **Yes, I can.** / **No, I can't.**

## 3. Từ vựng

| Tiếng Anh | Nghĩa |
|---|---|
| swim | bơi |
| dance | nhảy, múa |
| sing | hát |
| skate | trượt pa-tanh |
| ride a bike | đi xe đạp |
| play the piano | chơi đàn piano |

## 4. Lưu ý ngữ pháp

- Sau **can** dùng động từ nguyên mẫu, không thêm s, không có "to".
- Đúng: She can swim. Sai: She can swims. / She can to swim.

## 5. Bài tập mẫu

Sửa lỗi sai: "He can plays the piano."

Lời giải: **He can play the piano.**

> [!TIP]
> "can" giống nhau với mọi chủ ngữ: I can, you can, he can, she can, they can.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Quốc tịch, ngày tháng, giờ giấc và can",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Một bạn đến từ Nhật Bản (Japan). Quốc tịch của bạn ấy là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Japan (Nhật Bản) có quốc tịch tương ứng là Japanese.",
					Options: []sampleOption{
						{Content: "Vietnamese"},
						{Content: "English"},
						{Content: "Japanese", IsCorrect: true},
						{Content: "Malaysian"},
					},
				},
				{
					Prompt: "Số thứ tự 5th viết bằng chữ tiếng Anh là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "5th đọc và viết là fifth (đổi ve thành f rồi thêm th).",
					Options: []sampleOption{
						{Content: "fiveth"},
						{Content: "fifth", IsCorrect: true},
						{Content: "fivth"},
						{Content: "fifty"},
					},
				},
				{
					Prompt: "Chọn giới từ đúng: 'My birthday is ___ July.'",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Trước tên tháng (không có ngày cụ thể) dùng giới từ in: in July.",
					Options: []sampleOption{
						{Content: "on"},
						{Content: "at"},
						{Content: "of"},
						{Content: "in", IsCorrect: true},
					},
				},
				{
					Prompt: "Đồng hồ chỉ 7:30. Câu nào nói đúng giờ?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Đọc số giờ trước, số phút sau: 7:30 là seven thirty.",
					Options: []sampleOption{
						{Content: "It's seven thirty.", IsCorrect: true},
						{Content: "It's seven o'clock."},
						{Content: "It's thirty seven."},
						{Content: "It's three seventy."},
					},
				},
				{
					Prompt: "Chọn câu đúng ngữ pháp:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Sau can dùng động từ nguyên mẫu, không thêm s và không có to; can không đổi theo chủ ngữ.",
					Options: []sampleOption{
						{Content: "She can swims."},
						{Content: "She cans swim."},
						{Content: "She can swim.", IsCorrect: true},
						{Content: "She can to swim."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH5",
		Title:       "Tiếng Anh lớp 5",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 5 theo Chương trình GDPT 2018: thói quen hằng ngày với thì hiện tại đơn, nơi ở, và kể chuyện đã qua với thì quá khứ đơn.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Cuộc sống hằng ngày",
				Lessons: []sampleLesson{
					{
						Title:           "Daily routines – Thói quen hằng ngày",
						DurationMinutes: 20,
						Body: `## 1. Thì hiện tại đơn nói về thói quen

Dùng để nói việc em làm **thường xuyên, hằng ngày**.

- **I get up at six.** *(Tớ dậy lúc 6 giờ.)*
- **She goes to school by bike.** *(Bạn ấy đi học bằng xe đạp.)*

## 2. Thêm s/es khi chủ ngữ là he, she, it

| Quy tắc | Ví dụ |
|---|---|
| Đa số động từ: thêm s | get → gets, play → plays |
| Tận cùng o, s, sh, ch, x: thêm es | go → goes, watch → watches, brush → brushes |
| Phụ âm + y: đổi y thành i, thêm es | study → studies |
| Đặc biệt | have → has |

## 3. Trạng từ chỉ tần suất

| Tiếng Anh | Nghĩa |
|---|---|
| always | luôn luôn |
| usually | thường xuyên |
| often | hay, thường |
| sometimes | thỉnh thoảng |
| never | không bao giờ |

Vị trí: đứng **trước động từ thường**, **sau động từ to be**.

- I **always** do my homework. / He is **never** late.

## 4. Hỏi "bao lâu một lần"

- **How often do you go swimming?** – **Twice a week.** *(Hai lần một tuần.)*
- once a week *(một lần/tuần)*, every day *(mỗi ngày)*.

## 5. Bài tập mẫu

Chia động từ: "My mother (watch) TV every evening."

Lời giải: chủ ngữ là My mother (= she), watch tận cùng ch nên thêm es: **watches**.

> [!TIP]
> Nhớ quy tắc: he, she, it hoặc một người/vật số ít thì động từ thêm s hoặc es.`,
					},
					{
						Title:           "Where do you live?",
						DurationMinutes: 20,
						Body: `## 1. Hỏi địa chỉ

- **What's your address?** *(Địa chỉ của bạn là gì?)*
- **It's 25 Green Street.** *(Số 25 phố Green.)*

## 2. Hỏi nơi ở

- **Where do you live?** *(Bạn sống ở đâu?)*
- **I live in a flat in the city.** *(Tớ sống trong một căn hộ ở thành phố.)*
- **Where does she live?** – **She lives in the countryside.**

## 3. Từ vựng

| Tiếng Anh | Nghĩa |
|---|---|
| address | địa chỉ |
| street | phố, đường |
| flat | căn hộ |
| village | làng |
| city | thành phố |
| countryside | vùng nông thôn |
| quiet | yên tĩnh |
| busy | đông đúc, nhộn nhịp |

## 4. Hỏi nơi đó như thế nào

- **What's your village like?** *(Làng của bạn như thế nào?)*
- **It's quiet and beautiful.** *(Nó yên tĩnh và đẹp.)*

## 5. Bài tập mẫu

Chọn từ đúng: "Where (do / does) your grandparents live?"

Lời giải: your grandparents là số nhiều (= they) nên dùng **do**.

> [!NOTE]
> Câu hỏi với he, she, it dùng "does" và động từ giữ nguyên: Where does he live? (không viết "lives").`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Những ngày đã qua",
				Lessons: []sampleLesson{
					{
						Title:           "Where were you yesterday? – Was và were",
						DurationMinutes: 20,
						Body: `## 1. Quá khứ của động từ "be"

| Chủ ngữ | Hiện tại | Quá khứ |
|---|---|---|
| I | am | was |
| he, she, it | is | was |
| you, we, they | are | were |

## 2. Các dạng câu

- Khẳng định: **I was at home yesterday.** *(Hôm qua tớ ở nhà.)*
- Phủ định: **They weren't at school.** *(Họ đã không ở trường.)* — weren't = were not, wasn't = was not.
- Câu hỏi: **Were you at the zoo?** – **Yes, I was.** / **No, I wasn't.**

## 3. Từ chỉ thời gian quá khứ

| Tiếng Anh | Nghĩa |
|---|---|
| yesterday | hôm qua |
| last Sunday | Chủ nhật tuần trước |
| last week | tuần trước |
| two days ago | hai ngày trước |

## 4. Hội thoại mẫu

- **Mai:** Where were you last Sunday? *(Chủ nhật vừa rồi bạn ở đâu?)*
- **Nam:** I was at the beach with my family. *(Tớ ở bãi biển cùng gia đình.)*
- **Mai:** Was the weather nice? *(Thời tiết có đẹp không?)*
- **Nam:** Yes, it was. *(Có, đẹp lắm.)*

## 5. Bài tập mẫu

Điền **was** hoặc **were**: "My friends ___ at the library yesterday."

Lời giải: My friends là số nhiều (= they) nên dùng **were**.

> [!TIP]
> Thấy các từ yesterday, last..., ... ago là dấu hiệu câu nói về quá khứ.`,
					},
					{
						Title:           "What did you do? – Thì quá khứ đơn",
						DurationMinutes: 20,
						Body: `## 1. Động từ có quy tắc: thêm ed

| Nguyên mẫu | Quá khứ | Ghi chú |
|---|---|---|
| play | played | thêm ed |
| visit | visited | thêm ed |
| dance | danced | tận cùng e: thêm d |
| study | studied | phụ âm + y: đổi y thành i, thêm ed |
| stop | stopped | gấp đôi phụ âm cuối |

## 2. Động từ bất quy tắc (cần học thuộc)

| Nguyên mẫu | Quá khứ | Nghĩa |
|---|---|---|
| go | went | đi |
| have | had | có, ăn |
| see | saw | nhìn thấy |
| eat | ate | ăn |
| swim | swam | bơi |
| do | did | làm |

## 3. Câu hỏi và phủ định với "did"

- **What did you do last weekend?** – **I visited my grandparents.**
- **Did you go to the park?** – **Yes, I did.** / **No, I didn't.**
- **I didn't watch TV.** *(Tớ đã không xem TV.)*

Sau **did / didn't**, động từ trở về **nguyên mẫu**.

## 4. Bài tập mẫu

Sửa lỗi sai: "Did you went to the zoo yesterday?"

Lời giải: sau did dùng động từ nguyên mẫu → **Did you go to the zoo yesterday?**

> [!NOTE]
> Câu khẳng định dùng dạng quá khứ (went, played); câu hỏi và phủ định dùng did / didn't + động từ nguyên mẫu.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Hiện tại đơn, nơi ở và quá khứ đơn",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Dạng quá khứ của động từ 'go' là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "go là động từ bất quy tắc, quá khứ là went (gone là quá khứ phân từ).",
					Options: []sampleOption{
						{Content: "goed"},
						{Content: "went", IsCorrect: true},
						{Content: "gone"},
						{Content: "goes"},
					},
				},
				{
					Prompt: "Trạng từ nào có nghĩa là 'luôn luôn'?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "always là luôn luôn; never là không bao giờ, sometimes là thỉnh thoảng, often là thường.",
					Options: []sampleOption{
						{Content: "never"},
						{Content: "sometimes"},
						{Content: "often"},
						{Content: "always", IsCorrect: true},
					},
				},
				{
					Prompt: "Chọn từ đúng: 'My brother ___ his teeth every morning.'",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Chủ ngữ số ít (he), thói quen hằng ngày dùng hiện tại đơn; brush tận cùng sh nên thêm es: brushes.",
					Options: []sampleOption{
						{Content: "brush"},
						{Content: "brushs"},
						{Content: "brushes", IsCorrect: true},
						{Content: "brushing"},
					},
				},
				{
					Prompt: "Chọn từ đúng: 'They ___ at the zoo yesterday.'",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "yesterday là quá khứ, chủ ngữ they đi với were.",
					Options: []sampleOption{
						{Content: "were", IsCorrect: true},
						{Content: "was"},
						{Content: "are"},
						{Content: "is"},
					},
				},
				{
					Prompt: "Chọn câu đúng ngữ pháp:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu hỏi quá khứ đơn dùng Did + chủ ngữ + động từ nguyên mẫu: Did you visit ...?",
					Options: []sampleOption{
						{Content: "Did you visited your grandparents last Sunday?"},
						{Content: "Did you visit your grandparents last Sunday?", IsCorrect: true},
						{Content: "Do you visited your grandparents last Sunday?"},
						{Content: "Did you visits your grandparents last Sunday?"},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH6",
		Title:       "Tiếng Anh lớp 6",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 6 theo Chương trình GDPT 2018: thì hiện tại đơn và hiện tại tiếp diễn, there is/there are với giới từ chỉ vị trí, so sánh hơn của tính từ.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Trường học của em",
				Lessons: []sampleLesson{
					{
						Title:           "Thì hiện tại đơn và từ vựng trường học",
						DurationMinutes: 30,
						Body: `## 1. Cụm từ thường gặp ở trường

| Cụm từ | Nghĩa |
|---|---|
| do homework / do exercise | làm bài tập về nhà / tập thể dục |
| have lessons / have breakfast | có tiết học / ăn sáng |
| study English / study science | học tiếng Anh / học khoa học |
| play football / play the piano | chơi bóng đá / chơi đàn piano |
| wear a uniform | mặc đồng phục |

## 2. Cấu trúc thì hiện tại đơn

| Dạng | I / you / we / they | he / she / it |
|---|---|---|
| Khẳng định | I **study** maths. | She **studies** maths. |
| Phủ định | I **don't study** maths. | She **doesn't study** maths. |
| Câu hỏi | **Do** you **study** maths? | **Does** she **study** maths? |

Trả lời ngắn: Yes, I do. / No, she doesn't.

## 3. Cách dùng

- Thói quen, hoạt động lặp lại: We **have** English on Mondays.
- Sự thật hiển nhiên: The sun **rises** in the east.
- Lịch trình, thời gian biểu: The bus **leaves** at 6:45.

Dấu hiệu: always, usually, often, sometimes, never, every day, on Mondays...

## 4. Bài tập mẫu có lời giải

Chia động từ trong ngoặc:

1. Nam (not like) history. → **doesn't like** (chủ ngữ số ít, phủ định dùng doesn't + V nguyên mẫu).
2. (they / wear) a uniform? → **Do they wear** a uniform?
3. My sister (do) her homework after dinner. → **does** (do tận cùng o nên thêm es).

> [!TIP]
> Trong câu phủ định và câu hỏi, đã có does thì động từ chính không thêm s nữa: She doesn't like (không viết "doesn't likes").`,
					},
					{
						Title:           "Thì hiện tại tiếp diễn",
						DurationMinutes: 30,
						Body: `## 1. Cấu trúc

| Dạng | Cấu trúc | Ví dụ |
|---|---|---|
| Khẳng định | S + am/is/are + V-ing | They **are reading** in the library. |
| Phủ định | S + am/is/are + not + V-ing | He **isn't playing** football. |
| Câu hỏi | Am/Is/Are + S + V-ing? | **Are** you **doing** your homework? |

## 2. Quy tắc thêm -ing

- Đa số động từ: thêm ing → read → reading, study → studying.
- Tận cùng e câm: bỏ e, thêm ing → write → writing, make → making.
- Một nguyên âm + một phụ âm ở cuối (từ một âm tiết): gấp đôi phụ âm → swim → swimming, sit → sitting, run → running.
- Tận cùng ie: đổi thành y → lie → lying.

## 3. Cách dùng

- Hành động **đang diễn ra** lúc nói: Look! The students **are planting** trees.
- Kế hoạch đã sắp xếp trong **tương lai gần**: I **am visiting** my grandparents this weekend.

Dấu hiệu: now, right now, at the moment, Look!, Listen!

## 4. Phân biệt với hiện tại đơn

- Hiện tại đơn: thói quen → She **walks** to school every day.
- Hiện tại tiếp diễn: đang xảy ra → Today she **is riding** her bike to school.

Một số động từ chỉ trạng thái thường **không** dùng ở dạng tiếp diễn: like, love, know, want, understand.

## 5. Bài tập mẫu có lời giải

1. Listen! Someone (sing) in the music room. → **is singing**.
2. We (have) a picnic next Sunday — đã lên kế hoạch. → **are having**.
3. I (know) the answer now. → **know** (know là động từ trạng thái, không dùng tiếp diễn).

> [!NOTE]
> Khi thấy "Look!" hoặc "Listen!" ở đầu câu, hãy nghĩ ngay đến thì hiện tại tiếp diễn.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Ngôi nhà và khu phố",
				Lessons: []sampleLesson{
					{
						Title:           "There is / There are và giới từ chỉ vị trí",
						DurationMinutes: 30,
						Body: `## 1. Cấu trúc There is / There are

| Dạng | Số ít | Số nhiều |
|---|---|---|
| Khẳng định | There **is** a sofa in the living room. | There **are** two bedrooms in my flat. |
| Phủ định | There **isn't** a garden. | There **aren't any** posters on the wall. |
| Câu hỏi | **Is there** a lamp? – Yes, there is. | **Are there any** chairs? – No, there aren't. |

Lưu ý: khi liệt kê nhiều thứ, động từ thường chia theo danh từ **đứng ngay sau**: There **is** a bed and two chairs. / There **are** two chairs and a bed.

## 2. Từ vựng: các phòng và đồ đạc

| Tiếng Anh | Nghĩa |
|---|---|
| living room | phòng khách |
| bedroom | phòng ngủ |
| kitchen | nhà bếp |
| bathroom | phòng tắm |
| wardrobe | tủ quần áo |
| fridge | tủ lạnh |

## 3. Giới từ chỉ vị trí

| Giới từ | Nghĩa | Ví dụ |
|---|---|---|
| in | ở trong | The cat is **in** the box. |
| on | ở trên (tiếp xúc) | The book is **on** the desk. |
| under | ở dưới | The shoes are **under** the bed. |
| next to | bên cạnh | The lamp is **next to** the sofa. |
| behind | phía sau | The garden is **behind** the house. |
| in front of | phía trước | There is a tree **in front of** my house. |
| between ... and ... | ở giữa ... và ... | The table is **between** the sofa **and** the TV. |

## 4. Bài tập mẫu có lời giải

1. There (be) a fridge and a cooker in the kitchen. → **is** (danh từ đứng ngay sau là a fridge, số ít).
2. ___ there any pictures in your bedroom? → **Are** (pictures số nhiều).
3. Con mèo nằm dưới gầm bàn: The cat is ___ the table. → **under**.

> [!TIP]
> Dùng "any" trong câu phủ định và câu hỏi với danh từ số nhiều: Are there any books? There aren't any books.`,
					},
					{
						Title:           "So sánh hơn của tính từ",
						DurationMinutes: 30,
						Body: `## 1. Tính từ ngắn (một âm tiết, hoặc hai âm tiết tận cùng -y)

Cấu trúc: **S + be + adj-er + than + ...**

| Quy tắc | Ví dụ |
|---|---|
| Thêm er | cheap → cheaper, quiet → quieter |
| Tận cùng e: thêm r | large → larger, nice → nicer |
| Một nguyên âm + một phụ âm: gấp đôi phụ âm | big → bigger, hot → hotter |
| Tận cùng y: đổi y thành i, thêm er | noisy → noisier, busy → busier |

## 2. Tính từ dài (từ hai âm tiết trở lên, trừ loại tận cùng -y)

Cấu trúc: **S + be + more + adj + than + ...**

- modern → **more modern**, expensive → **more expensive**, convenient → **more convenient**.

## 3. Tính từ bất quy tắc

| Tính từ | So sánh hơn |
|---|---|
| good | better |
| bad | worse |
| far | farther / further |

## 4. Ví dụ về khu phố

- The city is **noisier than** the countryside.
- My new neighbourhood is **more convenient than** my old one.
- The air in the village is **better than** in the city.

## 5. Bài tập mẫu có lời giải

1. My house is (big) than Lan's house. → **bigger** (gấp đôi g).
2. This street is (busy) than that street. → **busier**.
3. Living in a flat is (expensive) than living in a village. → **more expensive**.

> [!NOTE]
> Không dùng đồng thời "more" và đuôi "-er": viết "bigger", không viết "more bigger".`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Các thì hiện tại, there is/there are và so sánh hơn",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Chọn dạng đúng của động từ: 'She ___ English on Mondays and Fridays.'",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Lịch học lặp lại dùng hiện tại đơn; chủ ngữ she, study tận cùng phụ âm + y nên đổi thành studies.",
					Options: []sampleOption{
						{Content: "study"},
						{Content: "studys"},
						{Content: "studies", IsCorrect: true},
						{Content: "studying"},
					},
				},
				{
					Prompt: "Giới từ nào có nghĩa là 'ở giữa (hai vật)'?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "between ... and ... nghĩa là ở giữa; behind là phía sau, under là ở dưới, next to là bên cạnh.",
					Options: []sampleOption{
						{Content: "behind"},
						{Content: "between", IsCorrect: true},
						{Content: "under"},
						{Content: "next to"},
					},
				},
				{
					Prompt: "Look! The children ___ football in the schoolyard.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Look! báo hiệu hành động đang diễn ra nên dùng hiện tại tiếp diễn; chủ ngữ số nhiều dùng are playing.",
					Options: []sampleOption{
						{Content: "are playing", IsCorrect: true},
						{Content: "play"},
						{Content: "plays"},
						{Content: "is playing"},
					},
				},
				{
					Prompt: "There ___ two bedrooms and a kitchen in my flat.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Danh từ đứng ngay sau là two bedrooms (số nhiều) nên dùng There are.",
					Options: []sampleOption{
						{Content: "is"},
						{Content: "be"},
						{Content: "am"},
						{Content: "are", IsCorrect: true},
					},
				},
				{
					Prompt: "Chọn câu đúng ngữ pháp:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "big là tính từ ngắn, tận cùng một nguyên âm + một phụ âm nên gấp đôi g rồi thêm er: bigger than; không dùng more với tính từ ngắn.",
					Options: []sampleOption{
						{Content: "My new school is more big than my old school."},
						{Content: "My new school is bigger than my old school.", IsCorrect: true},
						{Content: "My new school is biger than my old school."},
						{Content: "My new school is more bigger than my old school."},
					},
				},
			},
		},
	},
}
