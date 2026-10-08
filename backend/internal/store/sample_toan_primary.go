package store

// sampleToanPrimary là các lớp học mẫu môn Toán từ lớp 1 đến lớp 6 theo
// Chương trình GDPT 2018: mỗi lớp 2 chương × 2 bài giảng và một bài ôn tập trắc nghiệm.
var sampleToanPrimary = []sampleCourse{
	{
		Code:        "TOAN1",
		Title:       "Toán lớp 1",
		Description: "Lớp học mẫu môn Toán lớp 1 theo Chương trình GDPT 2018: đếm, đọc, viết, so sánh các số từ 0 đến 10 và phép cộng, phép trừ trong phạm vi 10.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Các số từ 0 đến 10",
				Lessons: []sampleLesson{
					{
						Title:           "Các số từ 0 đến 10",
						DurationMinutes: 15,
						Body: `## 1. Đếm từ 0 đến 10

Em cùng đếm nhé:

$$0,\ 1,\ 2,\ 3,\ 4,\ 5,\ 6,\ 7,\ 8,\ 9,\ 10$$

## 2. Đọc và viết số

- $0$: không
- $1$: một
- $2$: hai
- $3$: ba
- $4$: bốn
- $5$: năm
- $6$: sáu
- $7$: bảy
- $8$: tám
- $9$: chín
- $10$: mười

## 3. Đếm đồ vật

Có ba quả cam: ● ● ●. Ta viết số $3$.

Có năm bông hoa: ● ● ● ● ●. Ta viết số $5$.

Không có cái kẹo nào. Ta viết số $0$.

## 4. Số liền trước, số liền sau

- Số liền sau của $4$ là $5$ (đếm thêm $1$).
- Số liền trước của $4$ là $3$ (đếm bớt $1$).

**Bài tập mẫu.** Số liền sau của $7$ là số nào?

**Lời giải.** Đếm tiếp sau $7$ là $8$. Vậy số liền sau của $7$ là $8$.

> [!TIP]
> Em hãy đếm ngón tay để kiểm tra: mỗi ngón tay là một đơn vị.`,
					},
					{
						Title:           "So sánh các số trong phạm vi 10",
						DurationMinutes: 15,
						Body: `## 1. Dấu lớn hơn, bé hơn, bằng nhau

- Dấu $>$ đọc là "lớn hơn".
- Dấu $<$ đọc là "bé hơn".
- Dấu $=$ đọc là "bằng".

## 2. Ví dụ

Hàng trên có $3$ quả bóng, hàng dưới có $5$ quả bóng.

$3$ quả ít hơn $5$ quả. Ta viết: $3 < 5$.

Ngược lại: $5 > 3$.

Hai bạn mỗi bạn có $4$ cái bút. Ta viết: $4 = 4$.

## 3. Cách nhớ

Khi đếm, số nào đếm sau thì lớn hơn.

- $7$ đếm sau $2$ nên $7 > 2$.
- $6$ đếm trước $9$ nên $6 < 9$.

## 4. Sắp xếp các số

**Bài tập mẫu.** Xếp các số $6,\ 2,\ 9$ theo thứ tự từ bé đến lớn.

**Lời giải.** Số bé nhất là $2$, rồi đến $6$, lớn nhất là $9$.

Vậy ta có: $2,\ 6,\ 9$.

> [!TIP]
> Đầu nhọn của dấu $>$ hay $<$ luôn chỉ về số bé hơn.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Phép cộng, phép trừ trong phạm vi 10",
				Lessons: []sampleLesson{
					{
						Title:           "Phép cộng trong phạm vi 10",
						DurationMinutes: 20,
						Body: `## 1. Gộp lại thì làm phép cộng

Có $3$ con vịt, thêm $2$ con vịt nữa. Có tất cả $5$ con vịt.

$$3 + 2 = 5$$

Đọc là: ba cộng hai bằng năm.

## 2. Đổi chỗ các số

Đổi chỗ hai số trong phép cộng, kết quả không đổi:

$$2 + 3 = 5 \qquad 3 + 2 = 5$$

## 3. Cộng với 0

Một số cộng với $0$ thì bằng chính số đó:

$$4 + 0 = 4 \qquad 0 + 7 = 7$$

## 4. Một số phép cộng cần nhớ

- $1 + 9 = 10$
- $2 + 8 = 10$
- $3 + 7 = 10$
- $4 + 6 = 10$
- $5 + 5 = 10$

## 5. Bài tập mẫu

Lan có $4$ bông hoa. Mai cho Lan thêm $3$ bông hoa. Hỏi Lan có tất cả mấy bông hoa?

**Lời giải.** Lan có tất cả: $4 + 3 = 7$ (bông hoa).

> [!TIP]
> Khi thấy chữ "thêm", "tất cả", "gộp lại", em thường làm phép cộng.`,
					},
					{
						Title:           "Phép trừ trong phạm vi 10",
						DurationMinutes: 20,
						Body: `## 1. Bớt đi thì làm phép trừ

Có $7$ quả táo, ăn hết $2$ quả. Còn lại $5$ quả táo.

$$7 - 2 = 5$$

Đọc là: bảy trừ hai bằng năm.

## 2. Trừ đi 0 và trừ cho chính nó

- Một số trừ đi $0$ thì bằng chính số đó: $6 - 0 = 6$.
- Một số trừ đi chính nó thì bằng $0$: $8 - 8 = 0$.

## 3. Phép cộng và phép trừ là bạn

Từ phép cộng $5 + 3 = 8$, ta có hai phép trừ:

$$8 - 3 = 5 \qquad 8 - 5 = 3$$

Em có thể dùng phép cộng để kiểm tra phép trừ.

## 4. Bài tập mẫu

Trên cành có $9$ con chim. Có $4$ con bay đi. Hỏi trên cành còn lại mấy con chim?

**Lời giải.** Trên cành còn lại: $9 - 4 = 5$ (con chim).

**Kiểm tra:** $5 + 4 = 9$. Đúng rồi!

## 5. Em tự làm

- $10 - 3 = 7$
- $6 - 4 = 2$
- $9 - 9 = 0$

> [!NOTE]
> Khi thấy chữ "bớt", "bay đi", "còn lại", em thường làm phép trừ.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Các số từ 0 đến 10, phép cộng và phép trừ",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Số liền sau của số 6 là số nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Đếm tiếp sau 6 là 7, nên số liền sau của 6 là 7.",
					Options: []sampleOption{
						{Content: "5"},
						{Content: "7", IsCorrect: true},
						{Content: "8"},
						{Content: "6"},
					},
				},
				{
					Prompt: "Chọn dấu thích hợp điền vào chỗ trống: 8 ... 5",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Số 8 đếm sau số 5 nên 8 lớn hơn 5, ta viết 8 > 5.",
					Options: []sampleOption{
						{Content: "<"},
						{Content: "="},
						{Content: ">", IsCorrect: true},
						{Content: "Không điền được dấu nào"},
					},
				},
				{
					Prompt: "Kết quả của phép tính 4 + 5 là bao nhiêu?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Đếm thêm 5 từ 4: 5, 6, 7, 8, 9. Vậy 4 + 5 = 9.",
					Options: []sampleOption{
						{Content: "8"},
						{Content: "10"},
						{Content: "1"},
						{Content: "9", IsCorrect: true},
					},
				},
				{
					Prompt: "Phép tính nào dưới đây có kết quả bằng 6?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "9 − 3 = 6. Các phép tính còn lại đều có kết quả bằng 5.",
					Options: []sampleOption{
						{Content: "9 − 3", IsCorrect: true},
						{Content: "8 − 3"},
						{Content: "2 + 3"},
						{Content: "10 − 5"},
					},
				},
				{
					Prompt: "Trong vườn có 10 con gà, 3 con gà chạy vào chuồng. Hỏi ngoài vườn còn lại mấy con gà?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Có 10 con, bớt đi 3 con nên làm phép trừ: 10 − 3 = 7 (con gà).",
					Options: []sampleOption{
						{Content: "6"},
						{Content: "7", IsCorrect: true},
						{Content: "13"},
						{Content: "8"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN2",
		Title:       "Toán lớp 2",
		Description: "Lớp học mẫu môn Toán lớp 2 theo Chương trình GDPT 2018: tên gọi thành phần của phép cộng, phép trừ, bài toán nhiều hơn – ít hơn và phép cộng, phép trừ qua 10 trong phạm vi 20.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Ôn tập và bổ sung",
				Lessons: []sampleLesson{
					{
						Title:           "Số hạng, tổng – Số bị trừ, số trừ, hiệu",
						DurationMinutes: 15,
						Body: `## 1. Tên gọi trong phép cộng

$$\underbrace{12}_{\text{số hạng}} + \underbrace{5}_{\text{số hạng}} = \underbrace{17}_{\text{tổng}}$$

- $12$ và $5$ là **số hạng**.
- $17$ là **tổng**.
- $12 + 5$ cũng gọi là tổng.

## 2. Tên gọi trong phép trừ

$$\underbrace{18}_{\text{số bị trừ}} - \underbrace{6}_{\text{số trừ}} = \underbrace{12}_{\text{hiệu}}$$

- $18$ là **số bị trừ**.
- $6$ là **số trừ**.
- $12$ là **hiệu**. $18 - 6$ cũng gọi là hiệu.

## 3. Bài tập mẫu

**Bài 1.** Tính tổng, biết hai số hạng là $23$ và $14$.

**Lời giải.** Tổng là: $23 + 14 = 37$.

**Bài 2.** Tính hiệu, biết số bị trừ là $48$, số trừ là $25$.

**Lời giải.** Hiệu là: $48 - 25 = 23$.

## 4. Kiểm tra lại

Lấy hiệu cộng với số trừ thì được số bị trừ:

$$23 + 25 = 48$$

> [!TIP]
> Muốn kiểm tra phép trừ, em lấy hiệu cộng với số trừ. Nếu được số bị trừ là em làm đúng.`,
					},
					{
						Title:           "Bài toán về nhiều hơn, ít hơn",
						DurationMinutes: 20,
						Body: `## 1. Bài toán về nhiều hơn

**Ví dụ.** Hoa có $12$ cái nhãn vở. Lan có nhiều hơn Hoa $5$ cái nhãn vở. Hỏi Lan có bao nhiêu cái nhãn vở?

**Lời giải.**

Lan có số nhãn vở là:

$$12 + 5 = 17 \text{ (cái)}$$

Đáp số: $17$ cái nhãn vở.

## 2. Bài toán về ít hơn

**Ví dụ.** Bình có $15$ viên bi. An có ít hơn Bình $4$ viên bi. Hỏi An có bao nhiêu viên bi?

**Lời giải.**

An có số viên bi là:

$$15 - 4 = 11 \text{ (viên)}$$

Đáp số: $11$ viên bi.

## 3. Cách nhận biết

- "Nhiều hơn ... " thì thường làm **phép cộng**.
- "Ít hơn ... " thì thường làm **phép trừ**.

## 4. Các bước giải bài toán có lời văn

1. Đọc kĩ đề, tìm cái đã cho và cái cần tìm.
2. Viết câu lời giải.
3. Viết phép tính, ghi đơn vị trong ngoặc.
4. Viết đáp số.

> [!NOTE]
> Em nên vẽ sơ đồ đoạn thẳng: đoạn dài hơn chỉ số lớn hơn. Sơ đồ giúp em chọn đúng phép tính.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Phép cộng, phép trừ trong phạm vi 20",
				Lessons: []sampleLesson{
					{
						Title:           "Phép cộng (qua 10) trong phạm vi 20",
						DurationMinutes: 20,
						Body: `## 1. Cách làm tròn 10

Muốn tính $9 + 5$, em tách $5$ thành $1$ và $4$:

$$9 + 5 = 9 + 1 + 4 = 10 + 4 = 14$$

Vì $9$ thêm $1$ là tròn $10$, rồi cộng tiếp $4$.

## 2. Thêm ví dụ

$$8 + 6 = 8 + 2 + 4 = 10 + 4 = 14$$

$$7 + 5 = 7 + 3 + 2 = 10 + 2 = 12$$

## 3. Bảng cộng cần nhớ

- $9 + 2 = 11$, $9 + 3 = 12$, $9 + 4 = 13$
- $8 + 3 = 11$, $8 + 4 = 12$, $8 + 5 = 13$
- $7 + 4 = 11$, $7 + 6 = 13$, $6 + 6 = 12$

## 4. Bài tập mẫu

Lớp em có $8$ bạn nam đi trồng cây và $7$ bạn nữ. Hỏi có tất cả bao nhiêu bạn đi trồng cây?

**Lời giải.**

Có tất cả số bạn là:

$$8 + 7 = 15 \text{ (bạn)}$$

Đáp số: $15$ bạn.

> [!TIP]
> Hãy tìm xem số lớn hơn cần thêm mấy để tròn $10$, rồi tách số còn lại cho phù hợp.`,
					},
					{
						Title:           "Phép trừ (qua 10) trong phạm vi 20",
						DurationMinutes: 20,
						Body: `## 1. Trừ để về 10

Muốn tính $13 - 5$, em tách $5$ thành $3$ và $2$:

$$13 - 5 = 13 - 3 - 2 = 10 - 2 = 8$$

Vì $13$ bớt $3$ là còn $10$, rồi bớt tiếp $2$.

## 2. Thêm ví dụ

$$15 - 7 = 15 - 5 - 2 = 10 - 2 = 8$$

$$12 - 4 = 12 - 2 - 2 = 10 - 2 = 8$$

## 3. Dùng phép cộng để tính phép trừ

Em nhớ $9 + 7 = 16$, nên:

$$16 - 7 = 9 \qquad 16 - 9 = 7$$

## 4. Bài tập mẫu

Mẹ mua $14$ quả trứng, đã dùng $6$ quả. Hỏi còn lại bao nhiêu quả trứng?

**Lời giải.**

Số quả trứng còn lại là:

$$14 - 6 = 8 \text{ (quả)}$$

Đáp số: $8$ quả trứng.

**Kiểm tra:** $8 + 6 = 14$. Đúng!

> [!NOTE]
> Trừ qua $10$: em bớt trước để được tròn $10$, sau đó bớt tiếp phần còn lại.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Thành phần phép tính và cộng, trừ qua 10",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Trong phép tính 15 − 6 = 9, số 9 được gọi là gì?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Trong phép trừ, kết quả được gọi là hiệu. Vậy 9 là hiệu.",
					Options: []sampleOption{
						{Content: "Số bị trừ"},
						{Content: "Số trừ"},
						{Content: "Tổng"},
						{Content: "Hiệu", IsCorrect: true},
					},
				},
				{
					Prompt: "Kết quả của phép tính 9 + 6 là bao nhiêu?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "9 + 6 = 9 + 1 + 5 = 10 + 5 = 15.",
					Options: []sampleOption{
						{Content: "14"},
						{Content: "15", IsCorrect: true},
						{Content: "16"},
						{Content: "13"},
					},
				},
				{
					Prompt: "Để tính 7 + 5 bằng cách làm tròn 10, ta tách như thế nào là đúng?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "7 cần thêm 3 để tròn 10, nên tách 5 = 3 + 2: 7 + 3 + 2 = 10 + 2 = 12.",
					Options: []sampleOption{
						{Content: "7 + 3 + 2", IsCorrect: true},
						{Content: "7 + 5 + 3"},
						{Content: "7 + 2 + 2"},
						{Content: "7 + 4 + 2"},
					},
				},
				{
					Prompt: "Kết quả của phép tính 14 − 8 là bao nhiêu?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "14 − 8 = 14 − 4 − 4 = 10 − 4 = 6. Kiểm tra: 6 + 8 = 14.",
					Options: []sampleOption{
						{Content: "5"},
						{Content: "7"},
						{Content: "6", IsCorrect: true},
						{Content: "8"},
					},
				},
				{
					Prompt: "Tổ Một trồng được 9 cây. Tổ Hai trồng được nhiều hơn tổ Một 4 cây. Hỏi tổ Hai trồng được bao nhiêu cây?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Bài toán về nhiều hơn nên làm phép cộng: 9 + 4 = 13 (cây).",
					Options: []sampleOption{
						{Content: "5"},
						{Content: "13", IsCorrect: true},
						{Content: "14"},
						{Content: "12"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN3",
		Title:       "Toán lớp 3",
		Description: "Lớp học mẫu môn Toán lớp 3 theo Chương trình GDPT 2018: bảng nhân, bảng chia 6 và 7, nhân số có hai chữ số với số có một chữ số, phép chia hết và phép chia có dư.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Bảng nhân, bảng chia",
				Lessons: []sampleLesson{
					{
						Title:           "Bảng nhân 6, bảng nhân 7",
						DurationMinutes: 20,
						Body: `## 1. Bảng nhân 6

$$\begin{array}{ll}
6 \times 1 = 6 & 6 \times 6 = 36 \\
6 \times 2 = 12 & 6 \times 7 = 42 \\
6 \times 3 = 18 & 6 \times 8 = 48 \\
6 \times 4 = 24 & 6 \times 9 = 54 \\
6 \times 5 = 30 & 6 \times 10 = 60
\end{array}$$

Mỗi lần thêm một thừa số, tích tăng thêm $6$.

## 2. Bảng nhân 7

$$\begin{array}{ll}
7 \times 1 = 7 & 7 \times 6 = 42 \\
7 \times 2 = 14 & 7 \times 7 = 49 \\
7 \times 3 = 21 & 7 \times 8 = 56 \\
7 \times 4 = 28 & 7 \times 9 = 63 \\
7 \times 5 = 35 & 7 \times 10 = 70
\end{array}$$

## 3. Đổi chỗ thừa số

Đổi chỗ các thừa số thì tích không đổi: $6 \times 7 = 7 \times 6 = 42$.

## 4. Bài tập mẫu

**Bài 1.** Mỗi hộp có $6$ cái bút. Hỏi $4$ hộp như thế có bao nhiêu cái bút?

**Lời giải.** Số bút có là: $6 \times 4 = 24$ (cái). Đáp số: $24$ cái bút.

**Bài 2.** Mỗi tuần lễ có $7$ ngày. Hỏi $3$ tuần lễ có bao nhiêu ngày?

**Lời giải.** Số ngày là: $7 \times 3 = 21$ (ngày). Đáp số: $21$ ngày.

> [!TIP]
> Quên một kết quả? Hãy lấy kết quả liền trước cộng thêm: $7 \times 8 = 7 \times 7 + 7 = 49 + 7 = 56$.`,
					},
					{
						Title:           "Bảng chia 6, bảng chia 7",
						DurationMinutes: 20,
						Body: `## 1. Từ phép nhân đến phép chia

Từ $6 \times 5 = 30$ ta có hai phép chia:

$$30 : 6 = 5 \qquad 30 : 5 = 6$$

Lấy tích chia cho thừa số này thì được thừa số kia.

## 2. Bảng chia 6

$$\begin{array}{ll}
6 : 6 = 1 & 36 : 6 = 6 \\
12 : 6 = 2 & 42 : 6 = 7 \\
18 : 6 = 3 & 48 : 6 = 8 \\
24 : 6 = 4 & 54 : 6 = 9 \\
30 : 6 = 5 & 60 : 6 = 10
\end{array}$$

## 3. Bảng chia 7

$$\begin{array}{ll}
7 : 7 = 1 & 42 : 7 = 6 \\
14 : 7 = 2 & 49 : 7 = 7 \\
21 : 7 = 3 & 56 : 7 = 8 \\
28 : 7 = 4 & 63 : 7 = 9 \\
35 : 7 = 5 & 70 : 7 = 10
\end{array}$$

## 4. Bài tập mẫu

Có $42$ quả cam chia đều vào $7$ đĩa. Hỏi mỗi đĩa có mấy quả cam?

**Lời giải.** Mỗi đĩa có số quả cam là: $42 : 7 = 6$ (quả). Đáp số: $6$ quả cam.

**Kiểm tra:** $6 \times 7 = 42$.

> [!NOTE]
> Thuộc bảng nhân là sẽ thuộc bảng chia. Muốn tính $56 : 7$, em nhẩm: "$7$ nhân mấy bằng $56$?" và được $8$.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Phép nhân, phép chia trong phạm vi 100",
				Lessons: []sampleLesson{
					{
						Title:           "Nhân số có hai chữ số với số có một chữ số",
						DurationMinutes: 20,
						Body: `## 1. Phép nhân không nhớ

Tính $23 \times 3$:

$$\begin{array}{r}
23 \\
\times \quad 3 \\
\hline
69
\end{array}$$

- $3$ nhân $3$ bằng $9$, viết $9$.
- $3$ nhân $2$ bằng $6$, viết $6$.

Vậy $23 \times 3 = 69$.

## 2. Phép nhân có nhớ

Tính $26 \times 3$:

$$\begin{array}{r}
26 \\
\times \quad 3 \\
\hline
78
\end{array}$$

- $3$ nhân $6$ bằng $18$, viết $8$ nhớ $1$.
- $3$ nhân $2$ bằng $6$, thêm $1$ bằng $7$, viết $7$.

Vậy $26 \times 3 = 78$.

## 3. Bài tập mẫu

Mỗi thùng có $15$ chai nước. Hỏi $4$ thùng như thế có bao nhiêu chai nước?

**Lời giải.** Số chai nước là: $15 \times 4 = 60$ (chai).

Cách tính: $4$ nhân $5$ bằng $20$, viết $0$ nhớ $2$; $4$ nhân $1$ bằng $4$, thêm $2$ bằng $6$, viết $6$.

Đáp số: $60$ chai nước.

> [!TIP]
> Luôn nhân từ phải sang trái (bắt đầu từ hàng đơn vị) và đừng quên cộng phần nhớ vào hàng chục.`,
					},
					{
						Title:           "Phép chia hết và phép chia có dư",
						DurationMinutes: 20,
						Body: `## 1. Phép chia hết

Chia $12$ cái kẹo cho $3$ bạn, mỗi bạn được $4$ cái, không thừa cái nào:

$$12 : 3 = 4$$

Đây là **phép chia hết** (số dư bằng $0$).

## 2. Phép chia có dư

Chia $14$ cái kẹo cho $3$ bạn, mỗi bạn được $4$ cái, còn thừa $2$ cái:

$$14 : 3 = 4 \text{ (dư } 2)$$

Đây là **phép chia có dư**, $2$ là **số dư**.

## 3. Chú ý về số dư

Số dư luôn **bé hơn** số chia. Nếu số dư bằng hoặc lớn hơn số chia thì ta còn chia tiếp được.

## 4. Kiểm tra phép chia có dư

Lấy thương nhân với số chia, rồi cộng số dư, phải được số bị chia:

$$4 \times 3 + 2 = 14$$

## 5. Bài tập mẫu

Lớp có $31$ học sinh, mỗi bàn ngồi $2$ bạn. Cần ít nhất bao nhiêu cái bàn?

**Lời giải.** Ta có $31 : 2 = 15$ (dư $1$).

$15$ bàn ngồi được $30$ bạn, còn $1$ bạn cần thêm $1$ bàn nữa.

Vậy cần ít nhất: $15 + 1 = 16$ (cái bàn).

> [!NOTE]
> Trong bài toán thực tế, hãy đọc kĩ câu hỏi: có khi phải lấy thương, có khi phải lấy thương cộng thêm $1$ như bài trên.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Bảng nhân, bảng chia và phép chia có dư",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Kết quả của phép tính 7 × 8 là bao nhiêu?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Theo bảng nhân 7: 7 × 8 = 56.",
					Options: []sampleOption{
						{Content: "48"},
						{Content: "54"},
						{Content: "63"},
						{Content: "56", IsCorrect: true},
					},
				},
				{
					Prompt: "Trong phép chia có dư, số dư phải như thế nào so với số chia?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Số dư luôn bé hơn số chia; nếu không thì vẫn còn chia tiếp được.",
					Options: []sampleOption{
						{Content: "Lớn hơn số chia"},
						{Content: "Bằng số chia"},
						{Content: "Bé hơn số chia", IsCorrect: true},
						{Content: "Luôn bằng 0"},
					},
				},
				{
					Prompt: "Kết quả của phép tính 24 × 3 là bao nhiêu?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "3 × 4 = 12, viết 2 nhớ 1; 3 × 2 = 6, thêm 1 bằng 7. Vậy 24 × 3 = 72.",
					Options: []sampleOption{
						{Content: "72", IsCorrect: true},
						{Content: "62"},
						{Content: "612"},
						{Content: "27"},
					},
				},
				{
					Prompt: "Phép chia 38 : 6 có kết quả là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "6 × 6 = 36, 38 − 36 = 2 và 2 < 6. Vậy 38 : 6 = 6 (dư 2).",
					Options: []sampleOption{
						{Content: "6 (dư 1)"},
						{Content: "6 (dư 2)", IsCorrect: true},
						{Content: "5 (dư 8)"},
						{Content: "7 (dư 4)"},
					},
				},
				{
					Prompt: "Có 45 quả cam xếp vào các túi, mỗi túi 7 quả. Cần ít nhất bao nhiêu cái túi để xếp hết số cam?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "45 : 7 = 6 (dư 3). 6 túi chứa 42 quả, còn 3 quả cần thêm 1 túi nữa, nên cần 7 túi.",
					Options: []sampleOption{
						{Content: "6"},
						{Content: "5"},
						{Content: "7", IsCorrect: true},
						{Content: "8"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN4",
		Title:       "Toán lớp 4",
		Description: "Lớp học mẫu môn Toán lớp 4 theo Chương trình GDPT 2018: hàng và lớp của số có nhiều chữ số, so sánh và làm tròn số, cộng trừ số có nhiều chữ số và bài toán tìm hai số khi biết tổng và hiệu.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Số có nhiều chữ số",
				Lessons: []sampleLesson{
					{
						Title:           "Hàng và lớp",
						DurationMinutes: 20,
						Body: `## 1. Các hàng và các lớp

Các hàng được xếp thành từng lớp, mỗi lớp có ba hàng:

- **Lớp đơn vị**: hàng đơn vị, hàng chục, hàng trăm.
- **Lớp nghìn**: hàng nghìn, hàng chục nghìn, hàng trăm nghìn.
- **Lớp triệu**: hàng triệu, hàng chục triệu, hàng trăm triệu.

Khi viết số, ta để một khoảng trống giữa các lớp cho dễ đọc, ví dụ $352\ 418$.

## 2. Đọc số

Ta đọc từ trái sang phải, lần lượt từng lớp.

- $352\ 418$: ba trăm năm mươi hai nghìn bốn trăm mười tám.
- $4\ 025\ 600$: bốn triệu không trăm hai mươi lăm nghìn sáu trăm.

## 3. Giá trị của chữ số

Giá trị của một chữ số phụ thuộc vào vị trí (hàng) của nó.

Trong số $352\ 418$:

- Chữ số $3$ thuộc hàng trăm nghìn, có giá trị $300\ 000$.
- Chữ số $5$ thuộc hàng chục nghìn, có giá trị $50\ 000$.
- Chữ số $4$ thuộc hàng trăm, có giá trị $400$.

## 4. Bài tập mẫu

Viết số gồm: $7$ triệu, $3$ trăm nghìn, $6$ nghìn và $9$ đơn vị.

**Lời giải.** Các hàng còn thiếu ta viết chữ số $0$:

$$7\ 306\ 009$$

Đọc là: bảy triệu ba trăm linh sáu nghìn không trăm linh chín.

> [!TIP]
> Mỗi lớp luôn có đủ ba chữ số (trừ lớp đầu tiên bên trái). Hàng nào không có, em viết chữ số $0$.`,
					},
					{
						Title:           "So sánh và làm tròn số",
						DurationMinutes: 20,
						Body: `## 1. So sánh hai số tự nhiên

- Số nào có **nhiều chữ số hơn** thì lớn hơn: $102\ 300 > 98\ 765$.
- Nếu hai số có số chữ số bằng nhau, ta so sánh từng cặp chữ số **từ trái sang phải**.

**Ví dụ.** So sánh $45\ 312$ và $45\ 298$.

Hai số đều có năm chữ số; hàng chục nghìn và hàng nghìn giống nhau; ở hàng trăm có $3 > 2$.

Vậy $45\ 312 > 45\ 298$.

## 2. Làm tròn số

Muốn làm tròn đến một hàng, ta nhìn chữ số ở **hàng ngay bên phải**:

- Nếu chữ số đó bé hơn $5$: giữ nguyên hàng cần làm tròn.
- Nếu chữ số đó lớn hơn hoặc bằng $5$: cộng thêm $1$ vào hàng cần làm tròn.

Các chữ số phía sau hàng làm tròn đều thay bằng $0$.

## 3. Ví dụ làm tròn

- Làm tròn $72\ 486$ đến hàng nghìn: chữ số hàng trăm là $4 < 5$, được $72\ 000$.
- Làm tròn $72\ 586$ đến hàng nghìn: chữ số hàng trăm là $5$, được $73\ 000$.
- Làm tròn $3\ 456\ 000$ đến hàng trăm nghìn: chữ số hàng chục nghìn là $5$, được $3\ 500\ 000$.

## 4. Bài tập mẫu

Một thành phố có $1\ 284\ 530$ người. Làm tròn số dân đến hàng trăm nghìn.

**Lời giải.** Chữ số hàng chục nghìn là $8 > 5$, nên hàng trăm nghìn tăng từ $2$ lên $3$.

Số dân khoảng $1\ 300\ 000$ người.

> [!NOTE]
> Làm tròn số giúp ta ước lượng và nói về số lớn nhanh hơn, nhưng kết quả chỉ là giá trị gần đúng.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Các phép tính với số tự nhiên",
				Lessons: []sampleLesson{
					{
						Title:           "Cộng, trừ các số có nhiều chữ số",
						DurationMinutes: 20,
						Body: `## 1. Phép cộng

Đặt tính thẳng hàng rồi cộng từ phải sang trái:

$$\begin{array}{r}
37\ 584 \\
+\ 26\ 739 \\
\hline
64\ 323
\end{array}$$

Các bước: $4 + 9 = 13$, viết $3$ nhớ $1$; $8 + 3 + 1 = 12$, viết $2$ nhớ $1$; $5 + 7 + 1 = 13$, viết $3$ nhớ $1$; $7 + 6 + 1 = 14$, viết $4$ nhớ $1$; $3 + 2 + 1 = 6$, viết $6$.

## 2. Phép trừ

$$\begin{array}{r}
82\ 405 \\
-\ 37\ 618 \\
\hline
44\ 787
\end{array}$$

**Kiểm tra:** $44\ 787 + 37\ 618 = 82\ 405$.

## 3. Tính chất của phép cộng

- Giao hoán: $a + b = b + a$.
- Kết hợp: $(a + b) + c = a + (b + c)$.

**Ví dụ.** Tính thuận tiện:

$$125 + 368 + 75 = (125 + 75) + 368 = 200 + 368 = 568$$

## 4. Bài tập mẫu

Một kho có $45\ 250$ kg gạo, đã xuất đi $18\ 760$ kg. Hỏi kho còn lại bao nhiêu ki-lô-gam gạo?

**Lời giải.** Số gạo còn lại là:

$$45\ 250 - 18\ 760 = 26\ 490 \text{ (kg)}$$

Đáp số: $26\ 490$ kg gạo.

> [!TIP]
> Trước khi tính, hãy ước lượng: $45\ 000 - 19\ 000 = 26\ 000$. Kết quả thật phải gần số này.`,
					},
					{
						Title:           "Tìm hai số khi biết tổng và hiệu",
						DurationMinutes: 20,
						Body: `## 1. Công thức

Khi biết tổng và hiệu của hai số:

$$\text{Số lớn} = (\text{Tổng} + \text{Hiệu}) : 2$$

$$\text{Số bé} = (\text{Tổng} - \text{Hiệu}) : 2$$

## 2. Vì sao lại như vậy?

Vẽ sơ đồ: đoạn số lớn dài hơn đoạn số bé một phần đúng bằng hiệu.

Nếu bớt phần hơn đó đi, hai đoạn bằng nhau, tổng còn lại là $\text{Tổng} - \text{Hiệu}$, gồm hai lần số bé.

## 3. Ví dụ

Tổng hai số là $70$, hiệu là $16$. Tìm hai số đó.

**Lời giải.**

- Số lớn là: $(70 + 16) : 2 = 43$.
- Số bé là: $43 - 16 = 27$.

**Kiểm tra:** $43 + 27 = 70$ và $43 - 27 = 16$.

## 4. Bài toán có lời văn

Hai lớp 4A và 4B trồng được $150$ cây. Lớp 4A trồng nhiều hơn lớp 4B $12$ cây. Hỏi mỗi lớp trồng được bao nhiêu cây?

**Lời giải.**

Lớp 4A trồng được số cây là:

$$(150 + 12) : 2 = 81 \text{ (cây)}$$

Lớp 4B trồng được số cây là:

$$150 - 81 = 69 \text{ (cây)}$$

Đáp số: Lớp 4A: $81$ cây; lớp 4B: $69$ cây.

> [!NOTE]
> Tìm được một số rồi, em có thể lấy tổng trừ đi số đó (hoặc lấy số lớn trừ hiệu) để tìm số còn lại.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Số có nhiều chữ số và tìm hai số khi biết tổng và hiệu",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Trong số 6 830 512, chữ số 8 thuộc hàng nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Số 6 830 512 có 6 ở hàng triệu, 8 ở hàng trăm nghìn, 3 ở hàng chục nghìn.",
					Options: []sampleOption{
						{Content: "Hàng chục nghìn"},
						{Content: "Hàng triệu"},
						{Content: "Hàng nghìn"},
						{Content: "Hàng trăm nghìn", IsCorrect: true},
					},
				},
				{
					Prompt: "Làm tròn số 47 650 đến hàng nghìn ta được số nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Chữ số hàng trăm là 6 ≥ 5 nên hàng nghìn tăng thêm 1: được 48 000.",
					Options: []sampleOption{
						{Content: "47 000"},
						{Content: "48 000", IsCorrect: true},
						{Content: "47 700"},
						{Content: "50 000"},
					},
				},
				{
					Prompt: "Kết quả của phép tính 52 316 + 28 795 là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Cộng từ phải sang trái có nhớ: 52 316 + 28 795 = 81 111.",
					Options: []sampleOption{
						{Content: "81 111", IsCorrect: true},
						{Content: "80 011"},
						{Content: "81 101"},
						{Content: "71 111"},
					},
				},
				{
					Prompt: "Số lớn nhất trong các số 109 875; 190 758; 190 785; 109 857 là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Các số đều có sáu chữ số. 190 758 và 190 785 lớn hơn hai số kia; so tiếp hàng chục: 8 > 5 nên 190 785 lớn nhất.",
					Options: []sampleOption{
						{Content: "109 875"},
						{Content: "190 758"},
						{Content: "190 785", IsCorrect: true},
						{Content: "109 857"},
					},
				},
				{
					Prompt: "Tổng của hai số là 96, hiệu của chúng là 18. Số lớn là bao nhiêu?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Số lớn = (96 + 18) : 2 = 57. Kiểm tra: số bé là 39, 57 + 39 = 96 và 57 − 39 = 18.",
					Options: []sampleOption{
						{Content: "39"},
						{Content: "57", IsCorrect: true},
						{Content: "78"},
						{Content: "114"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN5",
		Title:       "Toán lớp 5",
		Description: "Lớp học mẫu môn Toán lớp 5 theo Chương trình GDPT 2018: phân số thập phân, hỗn số, khái niệm số thập phân và so sánh số thập phân.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chủ đề 1. Ôn tập và bổ sung về phân số",
				Lessons: []sampleLesson{
					{
						Title:           "Phân số thập phân",
						DurationMinutes: 20,
						Body: `## 1. Khái niệm

Các phân số có mẫu số là $10, 100, 1000, \ldots$ được gọi là **phân số thập phân**.

**Ví dụ:** $\dfrac{3}{10}$, $\dfrac{47}{100}$, $\dfrac{125}{1000}$ là các phân số thập phân.

Phân số $\dfrac{3}{5}$ không phải phân số thập phân vì mẫu số là $5$.

## 2. Chuyển phân số thành phân số thập phân

Một số phân số có thể viết thành phân số thập phân bằng cách nhân cả tử số và mẫu số với cùng một số thích hợp:

$$\frac{3}{5} = \frac{3 \times 2}{5 \times 2} = \frac{6}{10}$$

$$\frac{7}{25} = \frac{7 \times 4}{25 \times 4} = \frac{28}{100}$$

$$\frac{3}{4} = \frac{3 \times 25}{4 \times 25} = \frac{75}{100}$$

## 3. Rút gọn về phân số thập phân

Có khi phải chia cả tử và mẫu:

$$\frac{18}{200} = \frac{18 : 2}{200 : 2} = \frac{9}{100}$$

## 4. Bài tập mẫu

Viết $\dfrac{11}{20}$ thành phân số thập phân có mẫu số là $100$.

**Lời giải.** Vì $20 \times 5 = 100$ nên:

$$\frac{11}{20} = \frac{11 \times 5}{20 \times 5} = \frac{55}{100}$$

> [!TIP]
> Phân số mà mẫu số chỉ có thừa số $2$ và $5$ (như $2, 4, 5, 8, 20, 25, 50$) đều viết được thành phân số thập phân.`,
					},
					{
						Title:           "Hỗn số",
						DurationMinutes: 20,
						Body: `## 1. Khái niệm hỗn số

Có $2$ cái bánh nguyên và $\dfrac{3}{4}$ cái bánh. Ta viết gọn là $2\dfrac{3}{4}$ cái bánh.

$2\dfrac{3}{4}$ là một **hỗn số**, đọc là "hai và ba phần tư".

- $2$ là **phần nguyên**.
- $\dfrac{3}{4}$ là **phần phân số** (luôn bé hơn $1$).

Hỗn số có nghĩa là: $2\dfrac{3}{4} = 2 + \dfrac{3}{4}$.

## 2. Chuyển hỗn số thành phân số

Tử số mới bằng phần nguyên nhân mẫu số rồi cộng tử số; mẫu số giữ nguyên:

$$2\frac{3}{4} = \frac{2 \times 4 + 3}{4} = \frac{11}{4}$$

## 3. Chuyển phân số thành hỗn số

Lấy tử số chia cho mẫu số:

$$\frac{17}{5}: \quad 17 : 5 = 3 \text{ (dư } 2) \quad \Rightarrow \quad \frac{17}{5} = 3\frac{2}{5}$$

## 4. Bài tập mẫu

Tính $1\dfrac{1}{2} + 2\dfrac{1}{4}$.

**Lời giải.** Chuyển về phân số:

$$1\frac{1}{2} = \frac{3}{2} = \frac{6}{4}, \qquad 2\frac{1}{4} = \frac{9}{4}$$

$$\frac{6}{4} + \frac{9}{4} = \frac{15}{4} = 3\frac{3}{4}$$

> [!NOTE]
> Khi tính với hỗn số, cách an toàn là chuyển hết về phân số, tính xong mới đổi lại thành hỗn số nếu cần.`,
					},
				},
			},
			{
				Title: "Chủ đề 2. Số thập phân",
				Lessons: []sampleLesson{
					{
						Title:           "Khái niệm số thập phân",
						DurationMinutes: 20,
						Body: `## 1. Từ phân số thập phân đến số thập phân

$$\frac{1}{10} = 0{,}1 \qquad \frac{1}{100} = 0{,}01 \qquad \frac{1}{1000} = 0{,}001$$

Các số $0{,}1$; $0{,}01$; $0{,}001$ là **số thập phân**.

Tương tự: $\dfrac{3}{10} = 0{,}3$ và $\dfrac{47}{100} = 0{,}47$.

## 2. Cấu tạo số thập phân

Mỗi số thập phân gồm hai phần, ngăn cách bởi dấu phẩy:

- Bên trái dấu phẩy là **phần nguyên**.
- Bên phải dấu phẩy là **phần thập phân**.

Ví dụ: số $2{,}35$ có phần nguyên là $2$, phần thập phân là $35$.

## 3. Hàng của số thập phân

Trong số $2{,}358$:

- $3$ ở hàng **phần mười**, có giá trị $\dfrac{3}{10}$.
- $5$ ở hàng **phần trăm**, có giá trị $\dfrac{5}{100}$.
- $8$ ở hàng **phần nghìn**, có giá trị $\dfrac{8}{1000}$.

## 4. Đọc số thập phân

Đọc phần nguyên, đọc "phẩy", rồi đọc phần thập phân:

- $2{,}35$: hai phẩy ba mươi lăm.
- $6{,}08$: sáu phẩy không tám.

## 5. Bài tập mẫu

Viết hỗn số $5\dfrac{7}{10}$ thành số thập phân.

**Lời giải.** $5\dfrac{7}{10} = 5 + 0{,}7 = 5{,}7$.

> [!TIP]
> Mẫu số có bao nhiêu chữ số $0$ thì phần thập phân có bấy nhiêu chữ số: $\dfrac{9}{100} = 0{,}09$.`,
					},
					{
						Title:           "So sánh số thập phân",
						DurationMinutes: 20,
						Body: `## 1. Số thập phân bằng nhau

Viết thêm (hoặc bỏ đi) chữ số $0$ ở tận cùng bên phải phần thập phân thì giá trị không đổi:

$$0{,}5 = 0{,}50 = 0{,}500 \qquad 7{,}300 = 7{,}3$$

## 2. Quy tắc so sánh

- So sánh **phần nguyên** trước: số nào có phần nguyên lớn hơn thì lớn hơn.
- Nếu phần nguyên bằng nhau, so sánh lần lượt hàng **phần mười**, **phần trăm**, **phần nghìn**...

## 3. Ví dụ

- $12{,}4 > 9{,}87$ vì phần nguyên $12 > 9$.
- So sánh $3{,}72$ và $3{,}8$: phần nguyên bằng nhau, hàng phần mười có $7 < 8$, nên $3{,}72 < 3{,}8$.

Chú ý: $3{,}72$ có nhiều chữ số hơn nhưng vẫn **bé hơn** $3{,}8$.

## 4. Bài tập mẫu

Sắp xếp các số $4{,}05$; $4{,}5$; $4{,}15$; $3{,}99$ theo thứ tự từ bé đến lớn.

**Lời giải.**

- $3{,}99$ có phần nguyên bé nhất.
- Ba số còn lại có phần nguyên là $4$; so hàng phần mười: $0 < 1 < 5$.

Vậy: $3{,}99 < 4{,}05 < 4{,}15 < 4{,}5$.

> [!NOTE]
> Để dễ so sánh, hãy viết thêm chữ số $0$ cho các phần thập phân dài bằng nhau: $4{,}50$ và $4{,}15$.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Phân số thập phân, hỗn số và số thập phân",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Phân số nào sau đây là phân số thập phân?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Phân số thập phân có mẫu số là 10, 100, 1000,... nên 7/100 là phân số thập phân.",
					Options: []sampleOption{
						{Content: "3/5"},
						{Content: "10/7"},
						{Content: "9/20"},
						{Content: "7/100", IsCorrect: true},
					},
				},
				{
					Prompt: "Số thập phân 6,08 được đọc là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Phần nguyên là 6, phần thập phân là 08, nên đọc là sáu phẩy không tám.",
					Options: []sampleOption{
						{Content: "Sáu phẩy tám"},
						{Content: "Sáu phẩy tám mươi"},
						{Content: "Sáu phẩy không tám", IsCorrect: true},
						{Content: "Sáu mươi tám"},
					},
				},
				{
					Prompt: "Hỗn số 3 2/5 (ba và hai phần năm) viết dưới dạng phân số là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Tử số mới bằng 3 × 5 + 2 = 17, mẫu số giữ nguyên 5, nên được 17/5.",
					Options: []sampleOption{
						{Content: "17/5", IsCorrect: true},
						{Content: "13/5"},
						{Content: "32/5"},
						{Content: "5/17"},
					},
				},
				{
					Prompt: "Phân số 3/4 viết dưới dạng số thập phân là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "3/4 = 75/100 = 0,75.",
					Options: []sampleOption{
						{Content: "0,34"},
						{Content: "0,75", IsCorrect: true},
						{Content: "3,4"},
						{Content: "0,43"},
					},
				},
				{
					Prompt: "Sắp xếp các số 5,3; 5,09; 5,29; 5,1 theo thứ tự từ bé đến lớn.",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Phần nguyên đều là 5; so hàng phần mười: 0 < 1 < 2 < 3, nên 5,09 < 5,1 < 5,29 < 5,3.",
					Options: []sampleOption{
						{Content: "5,1; 5,3; 5,09; 5,29"},
						{Content: "5,3; 5,1; 5,29; 5,09"},
						{Content: "5,09; 5,1; 5,29; 5,3", IsCorrect: true},
						{Content: "5,09; 5,29; 5,1; 5,3"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN6",
		Title:       "Toán lớp 6",
		Description: "Lớp học mẫu môn Toán lớp 6 theo Chương trình GDPT 2018: tập hợp các số tự nhiên, lũy thừa và thứ tự thực hiện phép tính, dấu hiệu chia hết, số nguyên tố và ước chung lớn nhất.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Tập hợp các số tự nhiên",
				Lessons: []sampleLesson{
					{
						Title:           "Tập hợp và cách cho tập hợp",
						DurationMinutes: 30,
						Body: `## 1. Tập hợp và phần tử

Các đối tượng có chung một tính chất nào đó tạo thành một **tập hợp**. Mỗi đối tượng là một **phần tử** của tập hợp.

Người ta thường đặt tên tập hợp bằng chữ cái in hoa: $A, B, C, \ldots$

- $x$ là phần tử của $A$, ký hiệu $x \in A$ (đọc: $x$ thuộc $A$).
- $y$ không là phần tử của $A$, ký hiệu $y \notin A$ (đọc: $y$ không thuộc $A$).

## 2. Hai cách cho một tập hợp

**Cách 1. Liệt kê các phần tử**, đặt trong dấu ngoặc nhọn, cách nhau bởi dấu chấm phẩy:

$$A = \{0;\ 1;\ 2;\ 3\}$$

**Cách 2. Chỉ ra tính chất đặc trưng** của các phần tử:

$$A = \{x \in \mathbb{N} \mid x < 4\}$$

Mỗi phần tử chỉ được liệt kê một lần, thứ tự liệt kê tùy ý.

## 3. Tập hợp số tự nhiên

- $\mathbb{N} = \{0;\ 1;\ 2;\ 3;\ \ldots\}$ là tập hợp các số tự nhiên.
- $\mathbb{N}^* = \{1;\ 2;\ 3;\ \ldots\}$ là tập hợp các số tự nhiên khác $0$.

## 4. Ví dụ

**Ví dụ 1.** Viết tập hợp $B$ các chữ cái trong cụm từ "NHA TRANG".

**Lời giải.** $B = \{N;\ H;\ A;\ T;\ R;\ G\}$ (các chữ N, A lặp lại chỉ viết một lần).

**Ví dụ 2.** Viết tập hợp $C = \{x \in \mathbb{N}^* \mid x \le 5\}$ bằng cách liệt kê.

**Lời giải.** $C = \{1;\ 2;\ 3;\ 4;\ 5\}$. Khi đó $3 \in C$ nhưng $0 \notin C$.

> [!TIP]
> Có thể minh họa tập hợp bằng một vòng kín (biểu đồ Ven), mỗi phần tử là một dấu chấm bên trong vòng.`,
					},
					{
						Title:           "Lũy thừa với số mũ tự nhiên",
						DurationMinutes: 30,
						Body: `## 1. Định nghĩa

Lũy thừa bậc $n$ của $a$ là tích của $n$ thừa số bằng nhau, mỗi thừa số bằng $a$:

$$a^n = \underbrace{a \cdot a \cdot \ldots \cdot a}_{n \text{ thừa số}} \quad (n \in \mathbb{N}^*)$$

$a$ là **cơ số**, $n$ là **số mũ**. Quy ước $a^1 = a$ và $a^0 = 1$ (với $a \neq 0$).

$a^2$ đọc là "$a$ bình phương", $a^3$ đọc là "$a$ lập phương".

**Ví dụ:** $2^3 = 2 \cdot 2 \cdot 2 = 8$; $\quad 10^4 = 10\ 000$.

## 2. Nhân và chia hai lũy thừa cùng cơ số

$$a^m \cdot a^n = a^{m+n}$$

$$a^m : a^n = a^{m-n} \quad (a \neq 0,\ m \ge n)$$

**Ví dụ:** $5^3 \cdot 5^2 = 5^5 = 3125$; $\quad 7^6 : 7^4 = 7^2 = 49$.

## 3. Thứ tự thực hiện phép tính

- Biểu thức không có dấu ngoặc: **lũy thừa** $\rightarrow$ **nhân, chia** $\rightarrow$ **cộng, trừ**.
- Biểu thức có dấu ngoặc: thực hiện trong ngoặc tròn $( )$, rồi ngoặc vuông $[ ]$, cuối cùng ngoặc nhọn $\{ \}$.

## 4. Ví dụ áp dụng

**Ví dụ 1.**

$$3 \cdot 2^3 - 18 : 3^2 = 3 \cdot 8 - 18 : 9 = 24 - 2 = 22$$

**Ví dụ 2.**

$$80 - [130 - (12 - 4)^2] = 80 - [130 - 64] = 80 - 66 = 14$$

> [!NOTE]
> Lỗi hay gặp: $2^3 \neq 2 \cdot 3$. Lũy thừa là phép nhân lặp lại cơ số, không phải nhân cơ số với số mũ.`,
					},
				},
			},
			{
				Title: "Chương 2. Tính chia hết trong tập hợp các số tự nhiên",
				Lessons: []sampleLesson{
					{
						Title:           "Dấu hiệu chia hết cho 2, 5, 3, 9",
						DurationMinutes: 30,
						Body: `## 1. Dấu hiệu chia hết cho 2 và cho 5

- Số có chữ số tận cùng là $0; 2; 4; 6; 8$ thì chia hết cho $2$.
- Số có chữ số tận cùng là $0$ hoặc $5$ thì chia hết cho $5$.
- Số có chữ số tận cùng là $0$ thì chia hết cho cả $2$ và $5$.

## 2. Dấu hiệu chia hết cho 9 và cho 3

- Số có **tổng các chữ số** chia hết cho $9$ thì chia hết cho $9$.
- Số có **tổng các chữ số** chia hết cho $3$ thì chia hết cho $3$.

Số chia hết cho $9$ thì cũng chia hết cho $3$, nhưng điều ngược lại không luôn đúng (ví dụ $12$).

## 3. Ví dụ

Xét số $2\ 754$:

- Chữ số tận cùng là $4$ nên $2\ 754$ chia hết cho $2$, không chia hết cho $5$.
- Tổng các chữ số: $2 + 7 + 5 + 4 = 18$, chia hết cho $9$.

Vậy $2\ 754$ chia hết cho $2$, $3$ và $9$. Thật vậy, $2\ 754 : 9 = 306$.

## 4. Tìm chữ số chưa biết

**Bài toán.** Tìm chữ số $*$ để số $\overline{4*5}$ chia hết cho $9$.

**Lời giải.** Tổng các chữ số là $4 + * + 5 = 9 + *$.

Để tổng chia hết cho $9$ thì $* \in \{0;\ 9\}$. Ta được các số $405$ và $495$.

## 5. Tính chất chia hết của một tổng

Nếu mọi số hạng của tổng đều chia hết cho $m$ thì tổng chia hết cho $m$:

$$a \ \vdots\ m,\ b \ \vdots\ m \Rightarrow (a + b) \ \vdots\ m$$

Ví dụ: $36 \ \vdots\ 9$ và $45 \ \vdots\ 9$ nên $36 + 45 = 81$ chia hết cho $9$.

> [!TIP]
> Với dấu hiệu chia hết cho $2$ và $5$, chỉ cần nhìn chữ số cuối; với $3$ và $9$, phải cộng tất cả các chữ số.`,
					},
					{
						Title:           "Số nguyên tố và ước chung lớn nhất",
						DurationMinutes: 35,
						Body: `## 1. Số nguyên tố, hợp số

- **Số nguyên tố** là số tự nhiên lớn hơn $1$, chỉ có hai ước là $1$ và chính nó. Ví dụ: $2; 3; 5; 7; 11; 13$.
- **Hợp số** là số tự nhiên lớn hơn $1$, có nhiều hơn hai ước. Ví dụ: $4; 6; 9; 15$.
- Số $0$ và số $1$ không là số nguyên tố, cũng không là hợp số.

$2$ là số nguyên tố chẵn duy nhất.

## 2. Phân tích một số ra thừa số nguyên tố

Viết số đó dưới dạng tích các thừa số nguyên tố, ví dụ:

$$60 = 2 \cdot 2 \cdot 3 \cdot 5 = 2^2 \cdot 3 \cdot 5$$

## 3. Ước chung lớn nhất (ƯCLN)

ƯCLN của hai hay nhiều số là số lớn nhất trong tập hợp các ước chung của các số đó.

**Cách tìm ƯCLN bằng phân tích ra thừa số nguyên tố:**

1. Phân tích mỗi số ra thừa số nguyên tố.
2. Chọn các thừa số nguyên tố **chung**.
3. Lập tích các thừa số đã chọn, mỗi thừa số lấy với **số mũ nhỏ nhất**.

**Ví dụ.** Tìm ƯCLN$(36, 48)$.

$$36 = 2^2 \cdot 3^2 \qquad 48 = 2^4 \cdot 3$$

Thừa số chung là $2$ và $3$, nên ƯCLN$(36, 48) = 2^2 \cdot 3 = 12$.

Nếu ƯCLN$(a, b) = 1$ thì $a$ và $b$ gọi là hai số **nguyên tố cùng nhau**, ví dụ $8$ và $15$.

## 4. Bài toán thực tế

Có $36$ cái bút và $48$ quyển vở, chia đều vào các phần thưởng sao cho mỗi phần đều có cả bút và vở. Chia được nhiều nhất bao nhiêu phần thưởng?

**Lời giải.** Số phần thưởng là ước chung của $36$ và $48$; nhiều nhất là ƯCLN$(36, 48) = 12$.

Khi đó mỗi phần có $36 : 12 = 3$ cái bút và $48 : 12 = 4$ quyển vở.

> [!NOTE]
> Bài toán "chia đều nhiều nhất", "cắt thành các phần lớn nhất" thường dẫn đến việc tìm ƯCLN.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Tập hợp, lũy thừa và tính chia hết",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Cho tập hợp A = {2; 4; 6; 8}. Khẳng định nào sau đây đúng?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "8 là một phần tử của A nên 8 ∈ A. Còn 3 ∉ A, 6 ∈ A và 4 ∈ A.",
					Options: []sampleOption{
						{Content: "3 ∈ A"},
						{Content: "6 ∉ A"},
						{Content: "8 ∈ A", IsCorrect: true},
						{Content: "4 ∉ A"},
					},
				},
				{
					Prompt: "Trong các số sau, số nào là số nguyên tố?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "17 chỉ có hai ước là 1 và 17. Số 1 không là số nguyên tố; 21 = 3 · 7 và 27 = 3³ là hợp số.",
					Options: []sampleOption{
						{Content: "17", IsCorrect: true},
						{Content: "1"},
						{Content: "21"},
						{Content: "27"},
					},
				},
				{
					Prompt: "Giá trị của biểu thức 2³ · 3 − 12 : 2² là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Tính lũy thừa trước: 8 · 3 − 12 : 4; rồi nhân, chia: 24 − 3; cuối cùng trừ: 21.",
					Options: []sampleOption{
						{Content: "3"},
						{Content: "21", IsCorrect: true},
						{Content: "18"},
						{Content: "6"},
					},
				},
				{
					Prompt: "Số nào sau đây chia hết cho cả 3 và 5?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "2 535 có chữ số tận cùng là 5 nên chia hết cho 5, và tổng các chữ số 2 + 5 + 3 + 5 = 15 chia hết cho 3.",
					Options: []sampleOption{
						{Content: "1 250"},
						{Content: "3 456"},
						{Content: "4 205"},
						{Content: "2 535", IsCorrect: true},
					},
				},
				{
					Prompt: "Cô giáo có 24 quyển vở và 40 cái bút, muốn chia đều vào các phần thưởng sao cho mỗi phần đều có cả vở và bút. Có thể chia được nhiều nhất bao nhiêu phần thưởng?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Số phần thưởng nhiều nhất là ƯCLN(24, 40). Vì 24 = 2³ · 3 và 40 = 2³ · 5 nên ƯCLN = 2³ = 8.",
					Options: []sampleOption{
						{Content: "4"},
						{Content: "8", IsCorrect: true},
						{Content: "16"},
						{Content: "120"},
					},
				},
			},
		},
	},
}
