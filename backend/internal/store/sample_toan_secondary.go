package store

// sampleToanSecondary là các lớp học mẫu môn Toán từ lớp 7 đến lớp 12 theo
// Chương trình GDPT 2018: mỗi lớp 2 chương × 2 bài giảng và một bài ôn tập 5 câu.
var sampleToanSecondary = []sampleCourse{
	{
		Code:        "TOAN7",
		Title:       "Toán lớp 7",
		Description: "Lớp học mẫu môn Toán lớp 7 theo Chương trình GDPT 2018: số hữu tỉ và các phép tính, lũy thừa của số hữu tỉ, số vô tỉ, căn bậc hai số học và tập hợp số thực.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Số hữu tỉ",
				Lessons: []sampleLesson{
					{
						Title:           "Tập hợp các số hữu tỉ",
						DurationMinutes: 30,
						Body: `## 1. Số hữu tỉ

**Số hữu tỉ** là số viết được dưới dạng phân số $\dfrac{a}{b}$ với $a, b \in \mathbb{Z}$ và $b \neq 0$. Tập hợp các số hữu tỉ kí hiệu là $\mathbb{Q}$.

Ví dụ: $-0{,}5 = \dfrac{-1}{2}$; $\ 3 = \dfrac{3}{1}$; $\ 1\dfrac{2}{3} = \dfrac{5}{3}$ đều là các số hữu tỉ.

## 2. Biểu diễn số hữu tỉ trên trục số

Mỗi số hữu tỉ được biểu diễn bởi một điểm trên trục số. Để biểu diễn $\dfrac{3}{4}$, ta chia đoạn thẳng từ $0$ đến $1$ thành $4$ phần bằng nhau rồi lấy $3$ phần.

Hai số $\dfrac{3}{4}$ và $-\dfrac{3}{4}$ là hai **số đối** nhau: chúng nằm về hai phía của điểm $0$ và cách đều điểm $0$.

## 3. So sánh hai số hữu tỉ

Muốn so sánh hai số hữu tỉ, ta viết chúng dưới dạng phân số có cùng mẫu dương rồi so sánh các tử số.

Ví dụ: so sánh $-\dfrac{2}{3}$ và $-\dfrac{3}{4}$. Ta có $-\dfrac{2}{3} = -\dfrac{8}{12}$ và $-\dfrac{3}{4} = -\dfrac{9}{12}$. Vì $-8 > -9$ nên $-\dfrac{2}{3} > -\dfrac{3}{4}$.

- Số hữu tỉ lớn hơn $0$ gọi là số hữu tỉ dương; nhỏ hơn $0$ gọi là số hữu tỉ âm.
- Số $0$ không là số hữu tỉ dương cũng không là số hữu tỉ âm.

## 4. Bài tập mẫu

**Đề bài.** Sắp xếp các số $0{,}5;\ -\dfrac{1}{3};\ \dfrac{2}{5};\ -1$ theo thứ tự tăng dần.

**Lời giải.** Các số âm: $-1 < -\dfrac{1}{3}$. Các số dương: $\dfrac{2}{5} = 0{,}4 < 0{,}5$.

Vậy thứ tự tăng dần là: $-1;\ -\dfrac{1}{3};\ \dfrac{2}{5};\ 0{,}5$.

> [!TIP]
> Khi so sánh nhiều số hữu tỉ, có thể đổi tất cả về số thập phân hoặc về phân số cùng mẫu dương cho dễ nhìn.`,
					},
					{
						Title:           "Các phép tính với số hữu tỉ",
						DurationMinutes: 30,
						Body: `## 1. Cộng, trừ số hữu tỉ

Viết các số hữu tỉ dưới dạng phân số rồi cộng, trừ như cộng, trừ phân số.

Ví dụ: $\dfrac{-2}{3} + \dfrac{3}{4} = \dfrac{-8}{12} + \dfrac{9}{12} = \dfrac{1}{12}$.

**Quy tắc chuyển vế:** khi chuyển một số hạng từ vế này sang vế kia của một đẳng thức, ta phải đổi dấu số hạng đó: $x + a = b \Rightarrow x = b - a$.

## 2. Nhân, chia số hữu tỉ

$$\frac{a}{b} \cdot \frac{c}{d} = \frac{a \cdot c}{b \cdot d} \qquad \frac{a}{b} : \frac{c}{d} = \frac{a}{b} \cdot \frac{d}{c} \quad (c \neq 0)$$

Ví dụ: $-0{,}75 \cdot \dfrac{4}{9} = \dfrac{-3}{4} \cdot \dfrac{4}{9} = \dfrac{-3}{9} = -\dfrac{1}{3}$.

## 3. Lũy thừa với số mũ tự nhiên

Lũy thừa bậc $n$ của $x$ là tích của $n$ thừa số $x$: $x^n = x \cdot x \cdots x$. Quy ước $x^1 = x$ và $x^0 = 1$ (với $x \neq 0$).

- $x^m \cdot x^n = x^{m+n}$
- $x^m : x^n = x^{m-n}$ (với $x \neq 0$, $m \ge n$)
- $(x^m)^n = x^{m \cdot n}$

Ví dụ: $\left(-\dfrac{1}{2}\right)^3 = -\dfrac{1}{8}$; $\ (0{,}2)^5 : (0{,}2)^3 = (0{,}2)^2 = 0{,}04$.

## 4. Thứ tự thực hiện phép tính và bài tập mẫu

Trong ngoặc trước; ngoài ngoặc thì lũy thừa, rồi nhân chia, cuối cùng cộng trừ.

**Bài 1.** Tính $\left(\dfrac{2}{3}\right)^2 - \dfrac{1}{3} : 3 = \dfrac{4}{9} - \dfrac{1}{9} = \dfrac{3}{9} = \dfrac{1}{3}$.

**Bài 2.** Tìm $x$ biết $x - \dfrac{1}{2} = \dfrac{1}{3}$. Chuyển vế: $x = \dfrac{1}{3} + \dfrac{1}{2} = \dfrac{5}{6}$.

> [!NOTE]
> Lũy thừa bậc chẵn của một số âm là số dương, lũy thừa bậc lẻ của một số âm là số âm.`,
					},
				},
			},
			{
				Title: "Chương 2. Số thực",
				Lessons: []sampleLesson{
					{
						Title:           "Số vô tỉ. Căn bậc hai số học",
						DurationMinutes: 30,
						Body: `## 1. Số thập phân hữu hạn và vô hạn tuần hoàn

Mỗi số hữu tỉ được biểu diễn bởi một số thập phân **hữu hạn** hoặc **vô hạn tuần hoàn**.

- $\dfrac{3}{8} = 0{,}375$ (hữu hạn)
- $\dfrac{1}{3} = 0{,}333\ldots = 0{,}(3)$ (vô hạn tuần hoàn, chu kì $3$)
- $\dfrac{5}{11} = 0{,}4545\ldots = 0{,}(45)$ (chu kì $45$)

## 2. Số vô tỉ

**Số vô tỉ** là số viết được dưới dạng số thập phân vô hạn **không tuần hoàn**. Tập hợp các số vô tỉ kí hiệu là $\mathbb{I}$.

Ví dụ: $\pi = 3{,}14159\ldots$; $\ \sqrt{2} = 1{,}41421\ldots$ là các số vô tỉ.

## 3. Căn bậc hai số học

Căn bậc hai số học của số $a \ge 0$ là số $x \ge 0$ sao cho $x^2 = a$, kí hiệu $\sqrt{a}$.

Ví dụ: $\sqrt{49} = 7$ vì $7 \ge 0$ và $7^2 = 49$; $\ \sqrt{0{,}25} = 0{,}5$; $\ \sqrt{0} = 0$.

- Số âm không có căn bậc hai số học.
- Với $a \ge 0$ thì $\left(\sqrt{a}\right)^2 = a$.
- Các số $\sqrt{2}, \sqrt{3}, \sqrt{5}, \sqrt{7}$ là số vô tỉ.

## 4. Bài tập mẫu

**Bài 1.** Tính $\sqrt{16} + \sqrt{\dfrac{9}{25}} = 4 + \dfrac{3}{5} = 4{,}6$.

**Bài 2.** Ước lượng $\sqrt{10}$. Vì $3^2 = 9 < 10 < 16 = 4^2$ nên $3 < \sqrt{10} < 4$. Dùng máy tính cầm tay: $\sqrt{10} \approx 3{,}16$ (làm tròn đến hàng phần trăm).

> [!TIP]
> Để biết $\sqrt{a}$ nằm giữa hai số tự nhiên nào, hãy tìm hai số chính phương liền nhau kẹp số $a$ ở giữa.`,
					},
					{
						Title:           "Số thực. Giá trị tuyệt đối của một số thực",
						DurationMinutes: 30,
						Body: `## 1. Tập hợp số thực

Số hữu tỉ và số vô tỉ được gọi chung là **số thực**. Tập hợp các số thực kí hiệu là $\mathbb{R}$, và ta có:

$$\mathbb{N} \subset \mathbb{Z} \subset \mathbb{Q} \subset \mathbb{R}$$

Mỗi số thực được biểu diễn bởi đúng một điểm trên trục số; ngược lại mỗi điểm trên trục số biểu diễn một số thực. Vì vậy trục số còn gọi là **trục số thực**.

## 2. So sánh hai số thực

So sánh hai số thực giống như so sánh hai số thập phân. Ví dụ: $\sqrt{5} = 2{,}236\ldots > 2{,}2$.

Với hai số $a, b$ không âm: nếu $a < b$ thì $\sqrt{a} < \sqrt{b}$. Ví dụ: $\sqrt{7} < \sqrt{8}$.

## 3. Giá trị tuyệt đối

Giá trị tuyệt đối của số thực $x$, kí hiệu $|x|$, là khoảng cách từ điểm $x$ đến điểm $0$ trên trục số:

$$|x| = \begin{cases} x & \text{nếu } x \ge 0 \\ -x & \text{nếu } x < 0 \end{cases}$$

Ví dụ: $|-3{,}5| = 3{,}5$; $\ |\sqrt{2}| = \sqrt{2}$; $\ |1 - \sqrt{2}| = \sqrt{2} - 1$ (vì $1 < \sqrt{2}$ nên $1 - \sqrt{2} < 0$).

## 4. Bài tập mẫu

**Bài 1.** Tìm $x$ biết $|x| = 5$. Ta có $x = 5$ hoặc $x = -5$.

**Bài 2.** Tìm $x$ biết $|x - 1| = 3$.

Ta có $x - 1 = 3$ hoặc $x - 1 = -3$, suy ra $x = 4$ hoặc $x = -2$.

**Bài 3.** Tìm $x$ biết $|x| = -2$. Vì $|x| \ge 0$ với mọi $x$ nên không có giá trị nào của $x$ thoả mãn.

> [!NOTE]
> Với mọi số thực $x$: $|x| \ge 0$ và $|x| = |-x|$. Hai số đối nhau có giá trị tuyệt đối bằng nhau.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Số hữu tỉ và số thực",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Số nào sau đây là số vô tỉ?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "√7 = 2,6457... là số thập phân vô hạn không tuần hoàn nên là số vô tỉ; 0,(3) = 1/3, √9 = 3 và −5/4 đều là số hữu tỉ.",
					Options: []sampleOption{
						{Content: "0,(3)"},
						{Content: "√9"},
						{Content: "√7", IsCorrect: true},
						{Content: "−5/4"},
					},
				},
				{
					Prompt: "Giá trị tuyệt đối của −2,5 là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Giá trị tuyệt đối là khoảng cách đến điểm 0 nên luôn không âm: |−2,5| = 2,5.",
					Options: []sampleOption{
						{Content: "−2,5"},
						{Content: "2,5", IsCorrect: true},
						{Content: "0"},
						{Content: "±2,5"},
					},
				},
				{
					Prompt: "Kết quả của phép tính −2/3 + 3/4 là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Quy đồng mẫu 12: −2/3 + 3/4 = −8/12 + 9/12 = 1/12.",
					Options: []sampleOption{
						{Content: "1/7"},
						{Content: "−1/12"},
						{Content: "17/12"},
						{Content: "1/12", IsCorrect: true},
					},
				},
				{
					Prompt: "Kết quả của phép tính (0,2)⁵ : (0,2)³ là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Chia hai lũy thừa cùng cơ số thì trừ số mũ: (0,2)⁵ : (0,2)³ = (0,2)² = 0,04.",
					Options: []sampleOption{
						{Content: "0,04", IsCorrect: true},
						{Content: "0,4"},
						{Content: "0,008"},
						{Content: "0,2"},
					},
				},
				{
					Prompt: "Tìm x biết |x − 1| = 3.",
					Level:  "Vận dụng", Points: 2,
					Explanation: "|x − 1| = 3 nên x − 1 = 3 hoặc x − 1 = −3, suy ra x = 4 hoặc x = −2.",
					Options: []sampleOption{
						{Content: "x = 4"},
						{Content: "x = −2"},
						{Content: "x = 4 hoặc x = −2", IsCorrect: true},
						{Content: "x = 2 hoặc x = −4"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN8",
		Title:       "Toán lớp 8",
		Description: "Lớp học mẫu môn Toán lớp 8 theo Chương trình GDPT 2018: đơn thức và đa thức nhiều biến, các phép tính với đa thức, hằng đẳng thức đáng nhớ và phân tích đa thức thành nhân tử.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Đa thức",
				Lessons: []sampleLesson{
					{
						Title:           "Đơn thức và đa thức nhiều biến",
						DurationMinutes: 30,
						Body: `## 1. Đơn thức

**Đơn thức** là biểu thức đại số chỉ gồm một số, hoặc một biến, hoặc một tích giữa các số và các biến.

Ví dụ: $3x^2y$, $\ -\dfrac{1}{2}xy^3$, $\ 5$, $\ x$ là các đơn thức; còn $x + y$ hay $\dfrac{1}{x}$ không phải là đơn thức.

Trong đơn thức thu gọn $-4x^3y^2$: **hệ số** là $-4$, **phần biến** là $x^3y^2$, **bậc** là tổng số mũ các biến: $3 + 2 = 5$.

## 2. Đơn thức đồng dạng

Hai đơn thức **đồng dạng** là hai đơn thức có hệ số khác $0$ và có cùng phần biến.

Để cộng (trừ) các đơn thức đồng dạng, ta cộng (trừ) các hệ số và giữ nguyên phần biến:

$$3x^2y + 5x^2y - x^2y = (3 + 5 - 1)x^2y = 7x^2y$$

## 3. Đa thức

**Đa thức** là một tổng của những đơn thức; mỗi đơn thức trong tổng là một **hạng tử**.

**Thu gọn** đa thức là cộng các hạng tử đồng dạng với nhau. **Bậc** của đa thức thu gọn là bậc cao nhất trong các hạng tử của nó.

Ví dụ: $P = 2x^2y - 3xy + 5 - x^2y + xy$.

Thu gọn: $P = (2x^2y - x^2y) + (-3xy + xy) + 5 = x^2y - 2xy + 5$. Đa thức $P$ có bậc $3$.

## 4. Bài tập mẫu

**Đề bài.** Tính giá trị của đa thức $P$ ở trên tại $x = 1$, $y = -2$.

**Lời giải.** Thay vào dạng thu gọn:

$$P = 1^2 \cdot (-2) - 2 \cdot 1 \cdot (-2) + 5 = -2 + 4 + 5 = 7$$

> [!TIP]
> Luôn thu gọn đa thức trước khi tìm bậc hoặc tính giá trị — vừa nhanh hơn, vừa tránh nhầm lẫn.`,
					},
					{
						Title:           "Phép cộng, trừ và nhân đa thức",
						DurationMinutes: 30,
						Body: `## 1. Cộng, trừ hai đa thức

Viết hai đa thức trong ngoặc, bỏ ngoặc (nếu trước ngoặc có dấu trừ thì đổi dấu mọi hạng tử trong ngoặc), rồi thu gọn.

Ví dụ: cho $A = 3x^2 - 2xy + y^2$ và $B = x^2 + xy - 2y^2$.

- $A + B = 3x^2 - 2xy + y^2 + x^2 + xy - 2y^2 = 4x^2 - xy - y^2$
- $A - B = 3x^2 - 2xy + y^2 - x^2 - xy + 2y^2 = 2x^2 - 3xy + 3y^2$

## 2. Nhân đơn thức với đa thức

Nhân đơn thức với từng hạng tử của đa thức rồi cộng các tích: $A(B + C) = AB + AC$.

Ví dụ: $2x(3x^2 - xy + 1) = 6x^3 - 2x^2y + 2x$.

## 3. Nhân đa thức với đa thức

Nhân mỗi hạng tử của đa thức này với từng hạng tử của đa thức kia rồi cộng các tích:

$$(A + B)(C + D) = AC + AD + BC + BD$$

Ví dụ 1: $(x + 2)(x - 3) = x^2 - 3x + 2x - 6 = x^2 - x - 6$.

Ví dụ 2: $(2x - y)(x + 3y) = 2x^2 + 6xy - xy - 3y^2 = 2x^2 + 5xy - 3y^2$.

## 4. Chia đa thức cho đơn thức

Chia từng hạng tử của đa thức cho đơn thức (khi chia hết) rồi cộng các kết quả:

$$(6x^3y^2 - 9x^2y) : 3x^2y = 2xy - 3$$

> [!NOTE]
> Khi nhân, chú ý quy tắc dấu: hai thừa số cùng dấu cho tích dương, khác dấu cho tích âm; nhân hai lũy thừa cùng biến thì cộng số mũ.`,
					},
				},
			},
			{
				Title: "Chương 2. Hằng đẳng thức đáng nhớ và ứng dụng",
				Lessons: []sampleLesson{
					{
						Title:           "Hằng đẳng thức đáng nhớ",
						DurationMinutes: 35,
						Body: `## 1. Bình phương của một tổng, một hiệu

$$(A + B)^2 = A^2 + 2AB + B^2$$

$$(A - B)^2 = A^2 - 2AB + B^2$$

## 2. Hiệu hai bình phương

$$A^2 - B^2 = (A - B)(A + B)$$

## 3. Lập phương của một tổng, một hiệu

$$(A + B)^3 = A^3 + 3A^2B + 3AB^2 + B^3$$

$$(A - B)^3 = A^3 - 3A^2B + 3AB^2 - B^3$$

## 4. Tổng và hiệu hai lập phương

$$A^3 + B^3 = (A + B)(A^2 - AB + B^2)$$

$$A^3 - B^3 = (A - B)(A^2 + AB + B^2)$$

## 5. Ví dụ áp dụng

- Khai triển: $(2x + 3)^2 = (2x)^2 + 2 \cdot 2x \cdot 3 + 3^2 = 4x^2 + 12x + 9$.
- Khai triển: $(x - 2)^3 = x^3 - 6x^2 + 12x - 8$.
- Tính nhanh: $99^2 = (100 - 1)^2 = 10000 - 200 + 1 = 9801$.
- Tính nhanh: $51 \cdot 49 = (50 + 1)(50 - 1) = 2500 - 1 = 2499$.

> [!TIP]
> Hãy đọc mỗi hằng đẳng thức theo cả hai chiều: chiều thuận để khai triển, chiều ngược để viết gọn thành tích (phân tích thành nhân tử).`,
					},
					{
						Title:           "Phân tích đa thức thành nhân tử",
						DurationMinutes: 35,
						Body: `## 1. Khái niệm

**Phân tích đa thức thành nhân tử** là biến đổi đa thức đó thành một tích của những đa thức.

## 2. Phương pháp đặt nhân tử chung

Tìm nhân tử chung của các hạng tử (phần số là ƯCLN các hệ số, phần biến là các biến chung với số mũ nhỏ nhất):

$$6x^2y - 9xy^2 = 3xy(2x - 3y)$$

## 3. Phương pháp dùng hằng đẳng thức

- $x^2 - 6x + 9 = (x - 3)^2$
- $4x^2 - 25 = (2x)^2 - 5^2 = (2x - 5)(2x + 5)$
- $x^3 - 8 = (x - 2)(x^2 + 2x + 4)$

## 4. Phương pháp nhóm hạng tử

Nhóm các hạng tử thích hợp để xuất hiện nhân tử chung:

$$x^2 - xy + 3x - 3y = x(x - y) + 3(x - y) = (x - y)(x + 3)$$

## 5. Ứng dụng

**Tính nhanh.** $37^2 - 13^2 = (37 - 13)(37 + 13) = 24 \cdot 50 = 1200$.

**Tìm $x$.** Giải $x^2 - 4x = 0$. Ta có $x(x - 4) = 0$, nên $x = 0$ hoặc $x = 4$.

> [!NOTE]
> Thứ tự nên thử: đặt nhân tử chung trước, sau đó xem có dùng được hằng đẳng thức không, cuối cùng mới nhóm hạng tử. Nhiều bài cần phối hợp cả ba cách.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Đa thức và hằng đẳng thức đáng nhớ",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Bậc của đơn thức −4x³y² là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Bậc của đơn thức là tổng số mũ các biến: 3 + 2 = 5; hệ số −4 không ảnh hưởng đến bậc.",
					Options: []sampleOption{
						{Content: "3"},
						{Content: "−4"},
						{Content: "5", IsCorrect: true},
						{Content: "6"},
					},
				},
				{
					Prompt: "Khai triển (A − B)² ta được:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Theo hằng đẳng thức bình phương của một hiệu: (A − B)² = A² − 2AB + B².",
					Options: []sampleOption{
						{Content: "A² − B²"},
						{Content: "A² − 2AB + B²", IsCorrect: true},
						{Content: "A² + 2AB + B²"},
						{Content: "A² − AB + B²"},
					},
				},
				{
					Prompt: "Kết quả của phép nhân (x + 2)(x − 3) là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "(x + 2)(x − 3) = x² − 3x + 2x − 6 = x² − x − 6.",
					Options: []sampleOption{
						{Content: "x² − 6"},
						{Content: "x² + x − 6"},
						{Content: "x² − 5x − 6"},
						{Content: "x² − x − 6", IsCorrect: true},
					},
				},
				{
					Prompt: "Phân tích đa thức 4x² − 25 thành nhân tử ta được:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "4x² − 25 = (2x)² − 5² là hiệu hai bình phương nên bằng (2x − 5)(2x + 5).",
					Options: []sampleOption{
						{Content: "(2x − 5)(2x + 5)", IsCorrect: true},
						{Content: "(2x − 5)²"},
						{Content: "(4x − 5)(4x + 5)"},
						{Content: "(x − 5)(4x + 5)"},
					},
				},
				{
					Prompt: "Tính nhanh giá trị của 37² − 13².",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Dùng hiệu hai bình phương: 37² − 13² = (37 − 13)(37 + 13) = 24 × 50 = 1200.",
					Options: []sampleOption{
						{Content: "576"},
						{Content: "2400"},
						{Content: "1200", IsCorrect: true},
						{Content: "1000"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN9",
		Title:       "Toán lớp 9",
		Description: "Lớp học mẫu môn Toán lớp 9 theo Chương trình GDPT 2018: phương trình và hệ hai phương trình bậc nhất hai ẩn, căn bậc hai, căn bậc ba và biến đổi biểu thức chứa căn.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Phương trình và hệ hai phương trình bậc nhất hai ẩn",
				Lessons: []sampleLesson{
					{
						Title:           "Hệ hai phương trình bậc nhất hai ẩn",
						DurationMinutes: 30,
						Body: `## 1. Phương trình bậc nhất hai ẩn

Phương trình bậc nhất hai ẩn $x, y$ có dạng $ax + by = c$, trong đó $a, b, c$ là các số đã biết và $a, b$ không đồng thời bằng $0$.

Ví dụ: với phương trình $2x - y = 3$, cặp số $(2; 1)$ là một nghiệm vì $2 \cdot 2 - 1 = 3$.

Phương trình bậc nhất hai ẩn luôn có **vô số nghiệm**. Trên mặt phẳng toạ độ, tập nghiệm của $2x - y = 3$ được biểu diễn bởi đường thẳng $y = 2x - 3$.

## 2. Hệ hai phương trình bậc nhất hai ẩn

$$\begin{cases} ax + by = c \\ a'x + b'y = c' \end{cases}$$

Mỗi cặp số $(x_0; y_0)$ là nghiệm chung của cả hai phương trình gọi là một **nghiệm của hệ**.

Ví dụ: cặp $(2; 1)$ là nghiệm của hệ $\begin{cases} 2x - y = 3 \\ x + y = 3 \end{cases}$ vì $2 \cdot 2 - 1 = 3$ và $2 + 1 = 3$.

## 3. Minh hoạ hình học số nghiệm

Mỗi phương trình của hệ biểu diễn một đường thẳng. Khi đó:

- Hai đường thẳng **cắt nhau**: hệ có một nghiệm duy nhất.
- Hai đường thẳng **song song**: hệ vô nghiệm.
- Hai đường thẳng **trùng nhau**: hệ có vô số nghiệm.

## 4. Bài tập mẫu

**Đề bài.** Trong hai cặp số $(1; 2)$ và $(2; 1)$, cặp nào là nghiệm của phương trình $3x + 2y = 7$?

**Lời giải.** Với $(1; 2)$: $3 \cdot 1 + 2 \cdot 2 = 7$ (đúng). Với $(2; 1)$: $3 \cdot 2 + 2 \cdot 1 = 8 \neq 7$. Vậy chỉ $(1; 2)$ là nghiệm.

> [!TIP]
> Để kiểm tra một cặp số có là nghiệm hay không, chỉ cần thay vào **từng** phương trình — phải thoả mãn tất cả mới là nghiệm của hệ.`,
					},
					{
						Title:           "Giải hệ phương trình bằng phương pháp thế và cộng đại số",
						DurationMinutes: 35,
						Body: `## 1. Phương pháp thế

Từ một phương trình, biểu diễn một ẩn theo ẩn kia rồi thế vào phương trình còn lại.

Ví dụ: giải hệ $\begin{cases} x + y = 5 \\ 2x - y = 1 \end{cases}$

Từ phương trình thứ nhất: $y = 5 - x$. Thế vào phương trình thứ hai: $2x - (5 - x) = 1 \Rightarrow 3x = 6 \Rightarrow x = 2$, suy ra $y = 3$.

Vậy hệ có nghiệm duy nhất $(2; 3)$.

## 2. Phương pháp cộng đại số

Nhân hai vế mỗi phương trình với số thích hợp (nếu cần) để hệ số của một ẩn bằng nhau hoặc đối nhau, rồi cộng hoặc trừ từng vế để khử ẩn đó.

Ví dụ: giải hệ $\begin{cases} 3x + 2y = 8 \\ x - 2y = 0 \end{cases}$

Cộng từng vế: $4x = 8 \Rightarrow x = 2$. Thay vào $x - 2y = 0$ được $y = 1$. Nghiệm của hệ là $(2; 1)$.

## 3. Giải bài toán bằng cách lập hệ phương trình

**Đề bài.** Mua 3 quyển vở loại A và 2 quyển vở loại B hết 34 nghìn đồng; mua 1 quyển loại A và 4 quyển loại B hết 38 nghìn đồng. Tính giá mỗi loại vở.

**Lời giải.** Gọi giá một quyển loại A là $x$, loại B là $y$ (nghìn đồng, $x, y > 0$). Ta có hệ:

$$\begin{cases} 3x + 2y = 34 \\ x + 4y = 38 \end{cases}$$

Nhân hai vế phương trình đầu với $2$: $6x + 4y = 68$. Trừ từng vế cho phương trình thứ hai: $5x = 30 \Rightarrow x = 6$, suy ra $y = \dfrac{38 - 6}{4} = 8$.

Vậy vở loại A giá 6 nghìn đồng, loại B giá 8 nghìn đồng.

> [!NOTE]
> Nên dùng phương pháp thế khi có một ẩn mang hệ số $1$ hoặc $-1$; dùng cộng đại số khi hệ số của một ẩn bằng nhau hoặc đối nhau.`,
					},
				},
			},
			{
				Title: "Chương 2. Căn bậc hai và căn bậc ba",
				Lessons: []sampleLesson{
					{
						Title:           "Căn bậc hai và căn thức bậc hai",
						DurationMinutes: 30,
						Body: `## 1. Căn bậc hai

Căn bậc hai của số thực $a$ là số $x$ sao cho $x^2 = a$.

- Số dương $a$ có đúng hai căn bậc hai là $\sqrt{a}$ (căn bậc hai số học) và $-\sqrt{a}$. Ví dụ: $16$ có hai căn bậc hai là $4$ và $-4$.
- Số $0$ có đúng một căn bậc hai là $0$.
- Số âm không có căn bậc hai.

## 2. Căn thức bậc hai

Với $A$ là một biểu thức đại số, $\sqrt{A}$ gọi là căn thức bậc hai; $\sqrt{A}$ **xác định** khi $A \ge 0$.

Ví dụ: $\sqrt{2x - 6}$ xác định khi $2x - 6 \ge 0 \Leftrightarrow x \ge 3$.

## 3. Hằng đẳng thức $\sqrt{A^2} = |A|$

Ví dụ: $\sqrt{(3 - \sqrt{10})^2} = |3 - \sqrt{10}| = \sqrt{10} - 3$ (vì $3 = \sqrt{9} < \sqrt{10}$).

## 4. Khai phương một tích, một thương

$$\sqrt{AB} = \sqrt{A} \cdot \sqrt{B} \ (A, B \ge 0) \qquad \sqrt{\frac{A}{B}} = \frac{\sqrt{A}}{\sqrt{B}} \ (A \ge 0, B > 0)$$

Ví dụ: $\sqrt{2} \cdot \sqrt{18} = \sqrt{36} = 6$; $\ \dfrac{\sqrt{75}}{\sqrt{3}} = \sqrt{25} = 5$.

## 5. Căn bậc ba

Căn bậc ba của số $a$ là số $x$ sao cho $x^3 = a$, kí hiệu $\sqrt[3]{a}$. Mọi số thực đều có đúng một căn bậc ba.

Ví dụ: $\sqrt[3]{8} = 2$; $\ \sqrt[3]{-27} = -3$.

> [!TIP]
> $\sqrt{A^2}$ không phải lúc nào cũng bằng $A$: nếu $A < 0$ thì $\sqrt{A^2} = -A$. Hãy xét dấu của $A$ trước khi bỏ căn.`,
					},
					{
						Title:           "Biến đổi đơn giản biểu thức chứa căn bậc hai",
						DurationMinutes: 35,
						Body: `## 1. Đưa thừa số ra ngoài dấu căn

Với $B \ge 0$: $\sqrt{A^2B} = |A|\sqrt{B}$.

Ví dụ: $\sqrt{50} = \sqrt{25 \cdot 2} = 5\sqrt{2}$; $\ \sqrt{12} = \sqrt{4 \cdot 3} = 2\sqrt{3}$.

## 2. Đưa thừa số vào trong dấu căn

- Với $A \ge 0, B \ge 0$: $A\sqrt{B} = \sqrt{A^2B}$.
- Với $A < 0, B \ge 0$: $A\sqrt{B} = -\sqrt{A^2B}$.

Ví dụ: so sánh $3\sqrt{5}$ và $\sqrt{40}$. Ta có $3\sqrt{5} = \sqrt{45} > \sqrt{40}$.

## 3. Trục căn thức ở mẫu

- $\dfrac{A}{\sqrt{B}} = \dfrac{A\sqrt{B}}{B}$ với $B > 0$. Ví dụ: $\dfrac{6}{\sqrt{3}} = \dfrac{6\sqrt{3}}{3} = 2\sqrt{3}$.
- Nhân cả tử và mẫu với biểu thức liên hợp:

$$\frac{1}{\sqrt{3} - 1} = \frac{\sqrt{3} + 1}{(\sqrt{3} - 1)(\sqrt{3} + 1)} = \frac{\sqrt{3} + 1}{2}$$

## 4. Rút gọn biểu thức

**Đề bài.** Rút gọn $\sqrt{12} + \sqrt{27} - \sqrt{48}$.

**Lời giải.** Đưa thừa số ra ngoài dấu căn:

$$\sqrt{12} + \sqrt{27} - \sqrt{48} = 2\sqrt{3} + 3\sqrt{3} - 4\sqrt{3} = \sqrt{3}$$

> [!NOTE]
> Muốn cộng, trừ các căn thức, hãy đưa chúng về cùng một căn (căn đồng dạng) rồi cộng, trừ các hệ số phía trước.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Hệ phương trình bậc nhất hai ẩn và căn bậc hai",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Cặp số nào sau đây là nghiệm của phương trình 3x + 2y = 7?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Thay (1; 2): 3 × 1 + 2 × 2 = 7, đúng. Các cặp còn lại cho kết quả 8, 6 và 5, khác 7.",
					Options: []sampleOption{
						{Content: "(2; 1)"},
						{Content: "(1; 2)", IsCorrect: true},
						{Content: "(0; 3)"},
						{Content: "(3; −2)"},
					},
				},
				{
					Prompt: "Biểu thức √(2x − 6) xác định khi:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "√(2x − 6) xác định khi 2x − 6 ≥ 0, tức là x ≥ 3.",
					Options: []sampleOption{
						{Content: "x > 3"},
						{Content: "x ≤ 3"},
						{Content: "x ≥ −3"},
						{Content: "x ≥ 3", IsCorrect: true},
					},
				},
				{
					Prompt: "Nghiệm của hệ phương trình gồm x + y = 5 và 2x − y = 1 là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Cộng từng vế hai phương trình được 3x = 6, nên x = 2 và y = 5 − 2 = 3.",
					Options: []sampleOption{
						{Content: "(2; 3)", IsCorrect: true},
						{Content: "(3; 2)"},
						{Content: "(1; 4)"},
						{Content: "(4; 1)"},
					},
				},
				{
					Prompt: "Rút gọn biểu thức √12 + √27 − √48 ta được:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "√12 = 2√3, √27 = 3√3, √48 = 4√3 nên biểu thức bằng 2√3 + 3√3 − 4√3 = √3.",
					Options: []sampleOption{
						{Content: "0"},
						{Content: "√3", IsCorrect: true},
						{Content: "2√3"},
						{Content: "5√3"},
					},
				},
				{
					Prompt: "Mua 2 quyển vở loại A và 3 quyển loại B hết 31 nghìn đồng; mua 1 quyển loại A và 2 quyển loại B hết 18 nghìn đồng. Giá một quyển vở loại A là:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Gọi giá vở A, B là x, y (nghìn đồng): 2x + 3y = 31 và x + 2y = 18. Từ x = 18 − 2y thế vào được 36 − y = 31, nên y = 5 và x = 8.",
					Options: []sampleOption{
						{Content: "5 nghìn đồng"},
						{Content: "6 nghìn đồng"},
						{Content: "8 nghìn đồng", IsCorrect: true},
						{Content: "9 nghìn đồng"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN10",
		Title:       "Toán lớp 10",
		Description: "Lớp học mẫu môn Toán lớp 10 theo Chương trình GDPT 2018: mệnh đề, tập hợp và các phép toán trên tập hợp, hàm số bậc hai và dấu của tam thức bậc hai.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Mệnh đề và tập hợp",
				Lessons: []sampleLesson{
					{
						Title:           "Mệnh đề",
						DurationMinutes: 30,
						Body: `## 1. Mệnh đề

**Mệnh đề** là một câu khẳng định có tính đúng hoặc sai, không thể vừa đúng vừa sai.

Ví dụ: "$5$ là số nguyên tố" là một mệnh đề đúng; "$2 + 2 = 5$" là một mệnh đề sai; còn "Bạn có khỏe không?" không phải là mệnh đề vì không khẳng định đúng/sai.

## 2. Mệnh đề phủ định

Mệnh đề phủ định của mệnh đề $P$, ký hiệu $\overline{P}$ (hay "không $P$"), đúng khi $P$ sai và sai khi $P$ đúng.

Ví dụ: $P$: "$7$ là số chẵn" (sai) thì $\overline{P}$: "$7$ không là số chẵn" (đúng).

## 3. Mệnh đề kéo theo, mệnh đề đảo, mệnh đề tương đương

Mệnh đề "**Nếu** $P$ **thì** $Q$" (ký hiệu $P \Rightarrow Q$) chỉ sai khi $P$ đúng mà $Q$ sai.

Mệnh đề $Q \Rightarrow P$ gọi là **mệnh đề đảo** của $P \Rightarrow Q$. Nếu cả hai đều đúng, ta viết $P \Leftrightarrow Q$ ("$P$ khi và chỉ khi $Q$").

Ví dụ: "Nếu tam giác $ABC$ đều thì $ABC$ cân" là mệnh đề đúng, nhưng mệnh đề đảo "Nếu tam giác $ABC$ cân thì $ABC$ đều" là sai.

## 4. Mệnh đề chứa biến và các ký hiệu $\forall$, $\exists$

Câu "$x > 3$" chưa là mệnh đề, nhưng trở thành mệnh đề khi thay $x$ bằng một số cụ thể: với $x = 5$ ta được mệnh đề đúng.

- "$\forall x \in \mathbb{R},\ x^2 \ge 0$" (mọi số thực có bình phương không âm) là mệnh đề đúng.
- "$\exists x \in \mathbb{Z},\ x^2 = 2$" là mệnh đề sai.

Phủ định của "$\forall x \in X,\ P(x)$" là "$\exists x \in X,\ \overline{P(x)}$". Ví dụ: phủ định của "$\forall x \in \mathbb{R},\ x^2 + 1 > 0$" là "$\exists x \in \mathbb{R},\ x^2 + 1 \le 0$".

> [!TIP]
> Khi phủ định mệnh đề có lượng từ: đổi "$\forall$" thành "$\exists$" (và ngược lại), đồng thời phủ định phần tính chất phía sau.`,
					},
					{
						Title:           "Tập hợp và các phép toán trên tập hợp",
						DurationMinutes: 30,
						Body: `## 1. Tập hợp và cách cho tập hợp

Có hai cách cho một tập hợp: **liệt kê** các phần tử, hoặc **chỉ ra tính chất đặc trưng** của các phần tử.

Ví dụ: $A = \{x \in \mathbb{N} \mid x < 5\} = \{0; 1; 2; 3; 4\}$. Tập hợp không có phần tử nào là **tập rỗng**, ký hiệu $\varnothing$.

## 2. Tập con và hai tập hợp bằng nhau

$A \subset B$ khi mọi phần tử của $A$ đều thuộc $B$. Nếu $A \subset B$ và $B \subset A$ thì $A = B$.

Các tập con thường dùng của $\mathbb{R}$: khoảng $(a; b)$, đoạn $[a; b]$, nửa khoảng $[a; b)$, $(a; b]$, và $(-\infty; a)$, $[a; +\infty)$...

## 3. Các phép toán trên tập hợp

- **Giao**: $A \cap B = \{x \mid x \in A \text{ và } x \in B\}$
- **Hợp**: $A \cup B = \{x \mid x \in A \text{ hoặc } x \in B\}$
- **Hiệu**: $A \setminus B = \{x \mid x \in A \text{ và } x \notin B\}$
- **Phần bù**: nếu $A \subset E$ thì $C_E A = E \setminus A$.

## 4. Ví dụ áp dụng

**Ví dụ 1.** Cho $A = \{1; 2; 3; 4\}$ và $B = \{3; 4; 5; 6\}$:

$$A \cup B = \{1; 2; 3; 4; 5; 6\} \qquad A \cap B = \{3; 4\} \qquad A \setminus B = \{1; 2\}$$

**Ví dụ 2.** Cho $A = [1; 5)$ và $B = (3; 7]$. Biểu diễn trên trục số, ta được:

$$A \cap B = (3; 5) \qquad A \cup B = [1; 7] \qquad A \setminus B = [1; 3]$$

> [!TIP]
> Với tập hữu hạn, hãy vẽ biểu đồ Ven; với khoảng, đoạn, hãy biểu diễn trên trục số và chú ý đầu mút nào được lấy (ngoặc vuông) hay không (ngoặc tròn).`,
					},
				},
			},
			{
				Title: "Chương 2. Hàm số, đồ thị và ứng dụng",
				Lessons: []sampleLesson{
					{
						Title:           "Hàm số bậc hai",
						DurationMinutes: 30,
						Body: `## 1. Định nghĩa

Hàm số bậc hai là hàm số có dạng:

$$y = ax^2 + bx + c \quad (a \neq 0)$$

Tập xác định của hàm số bậc hai là $\mathbb{R}$.

## 2. Đồ thị

Đồ thị của hàm số bậc hai là một đường **parabol** có:

- **Đỉnh**: $I\left(-\dfrac{b}{2a};\ -\dfrac{\Delta}{4a}\right)$, với $\Delta = b^2 - 4ac$.
- **Trục đối xứng**: đường thẳng $x = -\dfrac{b}{2a}$.
- Bề lõm quay lên trên nếu $a > 0$, quay xuống dưới nếu $a < 0$.

## 3. Sự biến thiên

Nếu $a > 0$: hàm số nghịch biến trên $\left(-\infty; -\dfrac{b}{2a}\right)$, đồng biến trên $\left(-\dfrac{b}{2a}; +\infty\right)$ và có giá trị nhỏ nhất $-\dfrac{\Delta}{4a}$ tại đỉnh.

Nếu $a < 0$: chiều biến thiên ngược lại; hàm số có giá trị lớn nhất $-\dfrac{\Delta}{4a}$ tại đỉnh.

## 4. Ví dụ áp dụng

**Ví dụ.** Xét hàm số $y = x^2 - 4x + 3$.

Ta có $a = 1, b = -4, c = 3$, nên hoành độ đỉnh $x = -\dfrac{-4}{2 \cdot 1} = 2$, tung độ đỉnh $y = 2^2 - 4 \cdot 2 + 3 = -1$.

Vậy đỉnh $I(2; -1)$, trục đối xứng $x = 2$; vì $a = 1 > 0$ nên hàm số nghịch biến trên $(-\infty; 2)$, đồng biến trên $(2; +\infty)$ và có giá trị nhỏ nhất là $-1$.

Đồ thị cắt trục tung tại $(0; 3)$ và cắt trục hoành tại $(1; 0)$, $(3; 0)$.

> [!NOTE]
> Dấu của $\Delta = b^2 - 4ac$ cho biết số giao điểm của parabol với trục hoành: $\Delta > 0$ có 2 giao điểm, $\Delta = 0$ có 1 giao điểm (đỉnh nằm trên trục hoành), $\Delta < 0$ không có giao điểm.`,
					},
					{
						Title:           "Dấu của tam thức bậc hai",
						DurationMinutes: 30,
						Body: `## 1. Tam thức bậc hai

Tam thức bậc hai (theo $x$) là biểu thức có dạng $f(x) = ax^2 + bx + c$ với $a \neq 0$.

## 2. Định lý về dấu của tam thức bậc hai

Xét $\Delta = b^2 - 4ac$:

- Nếu $\Delta < 0$: $f(x)$ luôn cùng dấu với hệ số $a$ với mọi $x \in \mathbb{R}$.
- Nếu $\Delta = 0$: $f(x)$ cùng dấu với $a$ với mọi $x \neq -\dfrac{b}{2a}$ và $f\left(-\dfrac{b}{2a}\right) = 0$.
- Nếu $\Delta > 0$: $f(x)$ có hai nghiệm $x_1 < x_2$; $f(x)$ trái dấu với $a$ khi $x_1 < x < x_2$ và cùng dấu với $a$ khi $x < x_1$ hoặc $x > x_2$.

## 3. Ví dụ áp dụng

Xét dấu của $f(x) = x^2 - 5x + 6$.

$\Delta = 25 - 24 = 1 > 0$, hai nghiệm $x_1 = 2$, $x_2 = 3$; hệ số $a = 1 > 0$.

Vậy $f(x) > 0$ khi $x < 2$ hoặc $x > 3$; $f(x) < 0$ khi $2 < x < 3$; $f(x) = 0$ khi $x = 2$ hoặc $x = 3$.

## 4. Ứng dụng: giải bất phương trình bậc hai

**Đề bài.** Giải bất phương trình $x^2 - 5x + 6 < 0$.

**Lời giải.** Theo bảng xét dấu ở trên, $f(x) < 0$ khi và chỉ khi $2 < x < 3$. Tập nghiệm $S = (2; 3)$.

Tương tự, bất phương trình $x^2 - 5x + 6 \ge 0$ có tập nghiệm $S = (-\infty; 2] \cup [3; +\infty)$.

> [!TIP]
> Ghi nhớ quy tắc "trong trái, ngoài cùng": giữa hai nghiệm thì $f(x)$ trái dấu với $a$, ngoài khoảng hai nghiệm thì $f(x)$ cùng dấu với $a$.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Mệnh đề, tập hợp và hàm số bậc hai",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Câu nào sau đây là một mệnh đề?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Số 7 là số nguyên tố.\" là câu khẳng định có tính đúng sai xác định (đúng). Câu cảm thán, câu hỏi và câu chứa biến x chưa xác định được đúng hay sai.",
					Options: []sampleOption{
						{Content: "Hôm nay trời đẹp quá!"},
						{Content: "Bạn đã làm bài tập chưa?"},
						{Content: "Số 7 là số nguyên tố.", IsCorrect: true},
						{Content: "x + 1 > 3"},
					},
				},
				{
					Prompt: "Parabol y = ax² + bx + c (a ≠ 0) có đỉnh với hoành độ bằng:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Đỉnh của parabol là I(−b/(2a); −Δ/(4a)), nên hoành độ đỉnh là −b/(2a).",
					Options: []sampleOption{
						{Content: "b/(2a)"},
						{Content: "−b/(2a)", IsCorrect: true},
						{Content: "−b/a"},
						{Content: "−Δ/(4a)"},
					},
				},
				{
					Prompt: "Cho A = [1; 5) và B = (3; 7]. Tập hợp A ∩ B là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Phần chung của [1; 5) và (3; 7] là các số lớn hơn 3 và nhỏ hơn 5, tức là (3; 5).",
					Options: []sampleOption{
						{Content: "[1; 7]"},
						{Content: "[3; 5]"},
						{Content: "[1; 3]"},
						{Content: "(3; 5)", IsCorrect: true},
					},
				},
				{
					Prompt: "Toạ độ đỉnh của parabol y = x² − 4x + 3 là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Hoành độ đỉnh x = −(−4)/(2 × 1) = 2; tung độ y = 2² − 4 × 2 + 3 = −1. Vậy đỉnh I(2; −1).",
					Options: []sampleOption{
						{Content: "I(2; −1)", IsCorrect: true},
						{Content: "I(−2; 15)"},
						{Content: "I(2; 1)"},
						{Content: "I(4; 3)"},
					},
				},
				{
					Prompt: "Tập nghiệm của bất phương trình x² − 5x + 6 < 0 là:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "x² − 5x + 6 có hai nghiệm 2 và 3, hệ số a = 1 > 0 nên tam thức âm khi 2 < x < 3 (trong trái). Tập nghiệm là (2; 3).",
					Options: []sampleOption{
						{Content: "(−∞; 2) ∪ (3; +∞)"},
						{Content: "[2; 3]"},
						{Content: "(2; 3)", IsCorrect: true},
						{Content: "(−3; −2)"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN11",
		Title:       "Toán lớp 11",
		Description: "Lớp học mẫu môn Toán lớp 11 theo Chương trình GDPT 2018: góc lượng giác và giá trị lượng giác, phương trình lượng giác cơ bản, cấp số cộng và cấp số nhân.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Hàm số lượng giác và phương trình lượng giác",
				Lessons: []sampleLesson{
					{
						Title:           "Góc lượng giác và giá trị lượng giác",
						DurationMinutes: 35,
						Body: `## 1. Đơn vị radian

Ngoài độ, góc còn được đo bằng **radian** (rad), với quan hệ:

$$180^\circ = \pi \text{ rad}$$

Ví dụ: $60^\circ = \dfrac{\pi}{3}$; $\ 90^\circ = \dfrac{\pi}{2}$; $\ \dfrac{3\pi}{4} = 135^\circ$.

## 2. Góc lượng giác

Góc lượng giác $(Ou, Ov)$ có vô số số đo, sai khác nhau một bội của $2\pi$: nếu một số đo là $\alpha$ thì mọi số đo có dạng $\alpha + k2\pi$ ($k \in \mathbb{Z}$).

## 3. Giá trị lượng giác trên đường tròn lượng giác

Trên đường tròn lượng giác (tâm $O$, bán kính $1$), gọi $M(x; y)$ là điểm biểu diễn góc $\alpha$. Khi đó:

- $\cos\alpha = x$, $\ \sin\alpha = y$
- $\tan\alpha = \dfrac{\sin\alpha}{\cos\alpha}$ (khi $\cos\alpha \neq 0$), $\ \cot\alpha = \dfrac{\cos\alpha}{\sin\alpha}$ (khi $\sin\alpha \neq 0$)

Một số giá trị đặc biệt: $\sin\dfrac{\pi}{6} = \dfrac{1}{2}$, $\cos\dfrac{\pi}{3} = \dfrac{1}{2}$, $\sin\dfrac{\pi}{4} = \cos\dfrac{\pi}{4} = \dfrac{\sqrt{2}}{2}$, $\tan\dfrac{\pi}{3} = \sqrt{3}$.

## 4. Các hệ thức cơ bản

$$\sin^2\alpha + \cos^2\alpha = 1 \qquad 1 + \tan^2\alpha = \frac{1}{\cos^2\alpha} \qquad \tan\alpha \cdot \cot\alpha = 1$$

## 5. Bài tập mẫu

**Đề bài.** Cho $\sin\alpha = \dfrac{3}{5}$ với $\dfrac{\pi}{2} < \alpha < \pi$. Tính $\cos\alpha$ và $\tan\alpha$.

**Lời giải.** $\cos^2\alpha = 1 - \sin^2\alpha = 1 - \dfrac{9}{25} = \dfrac{16}{25}$. Vì $\dfrac{\pi}{2} < \alpha < \pi$ nên $\cos\alpha < 0$, do đó $\cos\alpha = -\dfrac{4}{5}$.

Suy ra $\tan\alpha = \dfrac{\sin\alpha}{\cos\alpha} = \dfrac{3/5}{-4/5} = -\dfrac{3}{4}$.

> [!TIP]
> Nhớ dấu theo góc phần tư: góc phần tư I tất cả dương; II chỉ $\sin$ dương; III chỉ $\tan$, $\cot$ dương; IV chỉ $\cos$ dương.`,
					},
					{
						Title:           "Phương trình lượng giác cơ bản",
						DurationMinutes: 35,
						Body: `## 1. Phương trình $\sin x = m$

- Nếu $|m| > 1$: phương trình vô nghiệm.
- Nếu $|m| \le 1$: gọi $\alpha$ là góc thoả mãn $\sin\alpha = m$, khi đó

$$\sin x = m \Leftrightarrow x = \alpha + k2\pi \ \text{ hoặc } \ x = \pi - \alpha + k2\pi \quad (k \in \mathbb{Z})$$

## 2. Phương trình $\cos x = m$

- Nếu $|m| > 1$: phương trình vô nghiệm.
- Nếu $|m| \le 1$: gọi $\alpha$ là góc thoả mãn $\cos\alpha = m$, khi đó $x = \pm\alpha + k2\pi$ ($k \in \mathbb{Z}$).

## 3. Phương trình $\tan x = m$ và $\cot x = m$

- $\tan x = m \Leftrightarrow x = \alpha + k\pi$ với $\tan\alpha = m$ (điều kiện $x \neq \dfrac{\pi}{2} + k\pi$).
- $\cot x = m \Leftrightarrow x = \alpha + k\pi$ với $\cot\alpha = m$ (điều kiện $x \neq k\pi$).

## 4. Ví dụ áp dụng

a) $\sin x = \dfrac{1}{2} \Leftrightarrow x = \dfrac{\pi}{6} + k2\pi$ hoặc $x = \dfrac{5\pi}{6} + k2\pi$.

b) $\cos x = -\dfrac{\sqrt{2}}{2}$. Vì $\cos\dfrac{3\pi}{4} = -\dfrac{\sqrt{2}}{2}$ nên $x = \pm\dfrac{3\pi}{4} + k2\pi$.

c) $\tan x = \sqrt{3} \Leftrightarrow x = \dfrac{\pi}{3} + k\pi$.

d) $2\cos x - 1 = 0 \Leftrightarrow \cos x = \dfrac{1}{2} \Leftrightarrow x = \pm\dfrac{\pi}{3} + k2\pi$.

(Trong các ví dụ trên, $k \in \mathbb{Z}$.)

> [!NOTE]
> Các trường hợp đặc biệt hay gặp: $\sin x = 0 \Leftrightarrow x = k\pi$; $\ \cos x = 0 \Leftrightarrow x = \dfrac{\pi}{2} + k\pi$; $\ \sin x = 1 \Leftrightarrow x = \dfrac{\pi}{2} + k2\pi$.`,
					},
				},
			},
			{
				Title: "Chương 2. Dãy số. Cấp số cộng và cấp số nhân",
				Lessons: []sampleLesson{
					{
						Title:           "Cấp số cộng",
						DurationMinutes: 30,
						Body: `## 1. Định nghĩa

**Cấp số cộng** là dãy số (hữu hạn hoặc vô hạn) mà kể từ số hạng thứ hai, mỗi số hạng bằng số hạng đứng ngay trước nó cộng với một số không đổi $d$ gọi là **công sai**:

$$u_{n+1} = u_n + d \quad (n \ge 1)$$

Ví dụ: dãy $2; 5; 8; 11; \ldots$ là cấp số cộng với $u_1 = 2$, $d = 3$.

## 2. Số hạng tổng quát

$$u_n = u_1 + (n - 1)d$$

Ví dụ: với $u_1 = 2$, $d = 3$ thì $u_{20} = 2 + 19 \cdot 3 = 59$.

## 3. Tính chất

Kể từ số hạng thứ hai, mỗi số hạng (trừ số hạng cuối nếu dãy hữu hạn) là trung bình cộng của hai số hạng kề nó:

$$u_k = \frac{u_{k-1} + u_{k+1}}{2} \quad (k \ge 2)$$

## 4. Tổng $n$ số hạng đầu

$$S_n = u_1 + u_2 + \cdots + u_n = \frac{n(u_1 + u_n)}{2} = \frac{n\left[2u_1 + (n - 1)d\right]}{2}$$

**Ví dụ 1.** $1 + 2 + 3 + \cdots + 100 = \dfrac{100 \cdot (1 + 100)}{2} = 5050$.

**Ví dụ 2.** Cấp số cộng có $u_1 = 2$, $d = 3$. Ta có $u_{10} = 2 + 9 \cdot 3 = 29$, nên $S_{10} = \dfrac{10 \cdot (2 + 29)}{2} = 155$.

> [!TIP]
> Mọi bài toán cấp số cộng đều quy về hai "ẩn" $u_1$ và $d$: hãy viết các dữ kiện theo $u_1$, $d$ rồi giải hệ phương trình.`,
					},
					{
						Title:           "Cấp số nhân",
						DurationMinutes: 30,
						Body: `## 1. Định nghĩa

**Cấp số nhân** là dãy số mà kể từ số hạng thứ hai, mỗi số hạng bằng số hạng đứng ngay trước nó nhân với một số không đổi $q$ gọi là **công bội**:

$$u_{n+1} = u_n \cdot q \quad (n \ge 1)$$

Ví dụ: dãy $3; 6; 12; 24; \ldots$ là cấp số nhân với $u_1 = 3$, $q = 2$.

## 2. Số hạng tổng quát

$$u_n = u_1 \cdot q^{n-1}$$

Ví dụ: với $u_1 = 3$, $q = 2$ thì $u_6 = 3 \cdot 2^5 = 96$.

## 3. Tính chất

Kể từ số hạng thứ hai: $u_k^2 = u_{k-1} \cdot u_{k+1}$ ($k \ge 2$).

## 4. Tổng $n$ số hạng đầu

Với $q \neq 1$:

$$S_n = \frac{u_1(1 - q^n)}{1 - q}$$

Ví dụ: với $u_1 = 3$, $q = 2$: $S_5 = \dfrac{3(1 - 2^5)}{1 - 2} = \dfrac{3 \cdot (-31)}{-1} = 93$. Kiểm tra: $3 + 6 + 12 + 24 + 48 = 93$.

## 5. Ứng dụng: lãi kép

Gửi $100$ triệu đồng với lãi suất $6\%$/năm theo thể thức lãi kép. Số tiền sau mỗi năm lập thành cấp số nhân với công bội $q = 1{,}06$:

$$T_n = 100 \cdot (1{,}06)^n \text{ (triệu đồng)}$$

Sau 2 năm: $T_2 = 100 \cdot 1{,}1236 = 112{,}36$ triệu đồng.

> [!NOTE]
> Khi $q = 1$, mọi số hạng đều bằng $u_1$ nên $S_n = n \cdot u_1$ — không dùng được công thức có mẫu $1 - q$.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Lượng giác, cấp số cộng và cấp số nhân",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Góc 60° đổi sang đơn vị radian bằng:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Vì 180° = π rad nên 60° = 60π/180 = π/3 rad.",
					Options: []sampleOption{
						{Content: "π/6"},
						{Content: "π/4"},
						{Content: "π/3", IsCorrect: true},
						{Content: "2π/3"},
					},
				},
				{
					Prompt: "Cấp số cộng (uₙ) có số hạng đầu u₁ và công sai d. Số hạng tổng quát của cấp số cộng là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Mỗi bước cộng thêm d, từ u₁ đến uₙ có n − 1 bước nên uₙ = u₁ + (n − 1)d.",
					Options: []sampleOption{
						{Content: "uₙ = u₁ + nd"},
						{Content: "uₙ = u₁ × dⁿ⁻¹"},
						{Content: "uₙ = u₁ − (n − 1)d"},
						{Content: "uₙ = u₁ + (n − 1)d", IsCorrect: true},
					},
				},
				{
					Prompt: "Nghiệm của phương trình sin x = 1/2 là (k ∈ ℤ):",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "sin(π/6) = 1/2 nên sin x = sin(π/6) khi x = π/6 + k2π hoặc x = π − π/6 + k2π = 5π/6 + k2π.",
					Options: []sampleOption{
						{Content: "x = ±π/6 + k2π"},
						{Content: "x = π/6 + k2π hoặc x = 5π/6 + k2π", IsCorrect: true},
						{Content: "x = π/3 + kπ"},
						{Content: "x = ±π/3 + k2π"},
					},
				},
				{
					Prompt: "Cho sin α = 3/5 với π/2 < α < π. Giá trị của cos α là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "cos²α = 1 − 9/25 = 16/25. Góc α thuộc góc phần tư thứ II nên cos α < 0, do đó cos α = −4/5.",
					Options: []sampleOption{
						{Content: "−4/5", IsCorrect: true},
						{Content: "4/5"},
						{Content: "−3/4"},
						{Content: "16/25"},
					},
				},
				{
					Prompt: "Cấp số nhân (uₙ) có u₁ = 3 và công bội q = 2. Tổng 5 số hạng đầu tiên của cấp số nhân là:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "S₅ = u₁(1 − q⁵)/(1 − q) = 3 × (1 − 32)/(1 − 2) = 93; kiểm tra: 3 + 6 + 12 + 24 + 48 = 93.",
					Options: []sampleOption{
						{Content: "48"},
						{Content: "96"},
						{Content: "93", IsCorrect: true},
						{Content: "45"},
					},
				},
			},
		},
	},
	{
		Code:        "TOAN12",
		Title:       "Toán lớp 12",
		Description: "Lớp học mẫu môn Toán lớp 12 theo Chương trình GDPT 2018: ứng dụng đạo hàm để xét tính đơn điệu, cực trị và giá trị lớn nhất, nhỏ nhất của hàm số; nguyên hàm và tích phân.",
		Cover:       "math",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Ứng dụng đạo hàm để khảo sát và vẽ đồ thị hàm số",
				Lessons: []sampleLesson{
					{
						Title:           "Tính đơn điệu và cực trị của hàm số",
						DurationMinutes: 35,
						Body: `## 1. Tính đơn điệu

Cho hàm số $y = f(x)$ có đạo hàm trên khoảng $K$:

- Nếu $f'(x) > 0$ với mọi $x \in K$ thì hàm số **đồng biến** trên $K$.
- Nếu $f'(x) < 0$ với mọi $x \in K$ thì hàm số **nghịch biến** trên $K$.

Nếu $f'(x) \ge 0$ (hoặc $\le 0$) trên $K$ và $f'(x) = 0$ chỉ tại hữu hạn điểm thì kết luận vẫn đúng.

## 2. Cực trị

Giả sử $f$ liên tục trên khoảng chứa $x_0$ và có đạo hàm ở hai bên $x_0$:

- Nếu $f'(x)$ đổi dấu từ **dương sang âm** khi $x$ đi qua $x_0$ thì $x_0$ là **điểm cực đại**.
- Nếu $f'(x)$ đổi dấu từ **âm sang dương** khi $x$ đi qua $x_0$ thì $x_0$ là **điểm cực tiểu**.

## 3. Các bước xét tính đơn điệu và tìm cực trị

1. Tìm tập xác định.
2. Tính $f'(x)$; tìm các điểm tại đó $f'(x) = 0$ hoặc không xác định.
3. Lập bảng biến thiên và kết luận.

## 4. Ví dụ áp dụng

Xét hàm số $y = x^3 - 3x^2 + 2$ trên $\mathbb{R}$.

$y' = 3x^2 - 6x = 3x(x - 2)$; $\ y' = 0 \Leftrightarrow x = 0$ hoặc $x = 2$.

- $y' > 0$ trên $(-\infty; 0)$ và $(2; +\infty)$: hàm số đồng biến trên mỗi khoảng này.
- $y' < 0$ trên $(0; 2)$: hàm số nghịch biến.

Hàm số đạt cực đại tại $x = 0$, $y_{\text{CĐ}} = 2$; đạt cực tiểu tại $x = 2$, $y_{\text{CT}} = 8 - 12 + 2 = -2$.

> [!TIP]
> $f'(x_0) = 0$ chưa đủ để kết luận $x_0$ là điểm cực trị. Ví dụ $y = x^3$ có $y'(0) = 0$ nhưng $y' \ge 0$ ở cả hai phía nên $x = 0$ không phải điểm cực trị.`,
					},
					{
						Title:           "Giá trị lớn nhất và giá trị nhỏ nhất của hàm số",
						DurationMinutes: 35,
						Body: `## 1. Định nghĩa

Cho hàm số $y = f(x)$ xác định trên tập $D$. Số $M$ là **giá trị lớn nhất** của hàm số trên $D$ nếu $f(x) \le M$ với mọi $x \in D$ và tồn tại $x_0 \in D$ sao cho $f(x_0) = M$. Tương tự với **giá trị nhỏ nhất** $m$.

## 2. Cách tìm trên một đoạn

Nếu $f$ liên tục trên đoạn $[a; b]$:

1. Tìm các điểm $x_1, x_2, \ldots$ thuộc $(a; b)$ mà tại đó $f'(x) = 0$ hoặc $f'(x)$ không xác định.
2. Tính $f(a)$, $f(b)$, $f(x_1)$, $f(x_2)$, ...
3. Số lớn nhất trong các giá trị trên là $\max\limits_{[a;b]} f(x)$, số nhỏ nhất là $\min\limits_{[a;b]} f(x)$.

## 3. Ví dụ áp dụng

Tìm giá trị lớn nhất, nhỏ nhất của $f(x) = x^3 - 3x^2 + 2$ trên đoạn $[-1; 3]$.

$f'(x) = 3x^2 - 6x = 0 \Leftrightarrow x = 0$ hoặc $x = 2$ (đều thuộc $(-1; 3)$).

$f(-1) = -2$; $\ f(0) = 2$; $\ f(2) = -2$; $\ f(3) = 2$.

Vậy $\max\limits_{[-1;3]} f(x) = 2$ (tại $x = 0$ và $x = 3$), $\min\limits_{[-1;3]} f(x) = -2$ (tại $x = -1$ và $x = 2$).

## 4. Bài toán thực tế

Từ một tấm bìa hình vuông cạnh $12$ cm, cắt ở bốn góc bốn hình vuông bằng nhau cạnh $x$ cm rồi gấp thành chiếc hộp không nắp. Tìm $x$ để thể tích hộp lớn nhất.

Thể tích: $V(x) = x(12 - 2x)^2$ với $0 < x < 6$.

$V'(x) = (12 - 2x)^2 - 4x(12 - 2x) = (12 - 2x)(12 - 6x)$; trên $(0; 6)$, $V'(x) = 0 \Leftrightarrow x = 2$.

$V'(x) > 0$ trên $(0; 2)$ và $V'(x) < 0$ trên $(2; 6)$, nên $V$ lớn nhất tại $x = 2$: $V(2) = 2 \cdot 8^2 = 128$ cm³.

> [!NOTE]
> Trên một khoảng (không phải đoạn), hàm số có thể không có giá trị lớn nhất hoặc nhỏ nhất; khi đó hãy lập bảng biến thiên để kết luận.`,
					},
				},
			},
			{
				Title: "Chương 2. Nguyên hàm và tích phân",
				Lessons: []sampleLesson{
					{
						Title:           "Nguyên hàm",
						DurationMinutes: 30,
						Body: `## 1. Định nghĩa

Hàm số $F(x)$ là một **nguyên hàm** của $f(x)$ trên khoảng $K$ nếu $F'(x) = f(x)$ với mọi $x \in K$.

Nếu $F(x)$ là một nguyên hàm của $f(x)$ thì mọi nguyên hàm của $f(x)$ có dạng $F(x) + C$ ($C$ là hằng số), ký hiệu:

$$\int f(x)\,dx = F(x) + C$$

## 2. Tính chất

- $\displaystyle\int k f(x)\,dx = k\int f(x)\,dx$ với $k \neq 0$.
- $\displaystyle\int \left[f(x) \pm g(x)\right]dx = \int f(x)\,dx \pm \int g(x)\,dx$.

## 3. Bảng nguyên hàm cơ bản

- $\displaystyle\int x^\alpha\,dx = \frac{x^{\alpha+1}}{\alpha+1} + C$ với $\alpha \neq -1$
- $\displaystyle\int \frac{1}{x}\,dx = \ln|x| + C$
- $\displaystyle\int e^x\,dx = e^x + C$; $\ \displaystyle\int a^x\,dx = \frac{a^x}{\ln a} + C$ ($0 < a \neq 1$)
- $\displaystyle\int \cos x\,dx = \sin x + C$; $\ \displaystyle\int \sin x\,dx = -\cos x + C$
- $\displaystyle\int \frac{1}{\cos^2 x}\,dx = \tan x + C$

## 4. Ví dụ áp dụng

**Ví dụ 1.** $\displaystyle\int (3x^2 - 4x + 5)\,dx = x^3 - 2x^2 + 5x + C$.

**Ví dụ 2.** $\displaystyle\int (2\cos x - e^x)\,dx = 2\sin x - e^x + C$.

**Ví dụ 3.** Tìm nguyên hàm $F(x)$ của $f(x) = 2x + 1$ thoả mãn $F(1) = 5$.

Ta có $F(x) = x^2 + x + C$. Từ $F(1) = 1 + 1 + C = 5$ suy ra $C = 3$. Vậy $F(x) = x^2 + x + 3$.

> [!TIP]
> Luôn có thể tự kiểm tra kết quả: lấy đạo hàm của đáp án, nếu ra đúng hàm số dưới dấu tích phân thì bạn đã làm đúng.`,
					},
					{
						Title:           "Tích phân",
						DurationMinutes: 35,
						Body: `## 1. Định nghĩa

Cho $f$ liên tục trên đoạn $[a; b]$ và $F$ là một nguyên hàm của $f$ trên đoạn đó. Khi đó:

$$\int_a^b f(x)\,dx = F(x)\Big|_a^b = F(b) - F(a)$$

## 2. Tính chất

- $\displaystyle\int_a^a f(x)\,dx = 0$; $\ \displaystyle\int_a^b f(x)\,dx = -\int_b^a f(x)\,dx$.
- $\displaystyle\int_a^b f(x)\,dx = \int_a^c f(x)\,dx + \int_c^b f(x)\,dx$ với $a < c < b$.
- $\displaystyle\int_a^b \left[k f(x) \pm g(x)\right]dx = k\int_a^b f(x)\,dx \pm \int_a^b g(x)\,dx$.

## 3. Ví dụ tính tích phân

- $\displaystyle\int_0^2 (3x^2 + 1)\,dx = (x^3 + x)\Big|_0^2 = (8 + 2) - 0 = 10$.
- $\displaystyle\int_0^{\pi/2} \cos x\,dx = \sin x\Big|_0^{\pi/2} = 1 - 0 = 1$.
- $\displaystyle\int_1^e \frac{1}{x}\,dx = \ln x\Big|_1^e = 1 - 0 = 1$.

## 4. Ý nghĩa hình học

Diện tích hình phẳng giới hạn bởi đồ thị $y = f(x)$, trục hoành và hai đường thẳng $x = a$, $x = b$ là:

$$S = \int_a^b |f(x)|\,dx$$

Ví dụ: diện tích hình phẳng giới hạn bởi $y = x^2$, trục hoành, $x = 0$ và $x = 3$ là $\displaystyle\int_0^3 x^2\,dx = \frac{x^3}{3}\Big|_0^3 = 9$.

## 5. Ứng dụng trong vật lí

Vật chuyển động với vận tốc $v(t)$ thì quãng đường đi được từ thời điểm $a$ đến $b$ là $s = \displaystyle\int_a^b v(t)\,dt$.

Ví dụ: $v(t) = 2t + 3$ (m/s), từ $t = 0$ đến $t = 4$ giây: $s = (t^2 + 3t)\Big|_0^4 = 16 + 12 = 28$ (m).

> [!NOTE]
> Giá trị tích phân không phụ thuộc vào việc chọn nguyên hàm nào: hằng số $C$ luôn bị triệt tiêu khi lấy $F(b) - F(a)$.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Ứng dụng đạo hàm, nguyên hàm và tích phân",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Nếu f′(x) > 0 với mọi x thuộc khoảng K thì hàm số y = f(x):",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Đạo hàm dương trên K nghĩa là hàm số tăng trên K, tức là đồng biến trên K.",
					Options: []sampleOption{
						{Content: "nghịch biến trên K"},
						{Content: "đồng biến trên K", IsCorrect: true},
						{Content: "không đổi trên K"},
						{Content: "có cực trị trên K"},
					},
				},
				{
					Prompt: "Họ nguyên hàm của hàm số f(x) = cos x là:",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Vì (sin x)′ = cos x nên ∫cos x dx = sin x + C.",
					Options: []sampleOption{
						{Content: "−sin x + C"},
						{Content: "cos x + C"},
						{Content: "−cos x + C"},
						{Content: "sin x + C", IsCorrect: true},
					},
				},
				{
					Prompt: "Hàm số y = x³ − 3x² + 2 đạt cực tiểu tại:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "y′ = 3x² − 6x = 3x(x − 2) đổi dấu từ âm sang dương khi x đi qua 2, nên hàm số đạt cực tiểu tại x = 2 (còn x = 0 là điểm cực đại).",
					Options: []sampleOption{
						{Content: "x = 2", IsCorrect: true},
						{Content: "x = 0"},
						{Content: "x = −2"},
						{Content: "x = 1"},
					},
				},
				{
					Prompt: "Giá trị của tích phân ∫ từ 0 đến 2 của (3x² + 1) dx là:",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Một nguyên hàm là x³ + x, nên tích phân bằng (2³ + 2) − (0 + 0) = 10.",
					Options: []sampleOption{
						{Content: "8"},
						{Content: "12"},
						{Content: "10", IsCorrect: true},
						{Content: "14"},
					},
				},
				{
					Prompt: "Một vật chuyển động với vận tốc v(t) = 2t + 3 (m/s). Quãng đường vật đi được từ t = 0 đến t = 4 giây là:",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Quãng đường bằng tích phân của vận tốc: s = (t² + 3t) lấy từ 0 đến 4 = 16 + 12 = 28 m.",
					Options: []sampleOption{
						{Content: "11 m"},
						{Content: "28 m", IsCorrect: true},
						{Content: "19 m"},
						{Content: "32 m"},
					},
				},
			},
		},
	},
}
