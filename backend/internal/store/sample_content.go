package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

// sampleLesson là một bài giảng mẫu (nội dung văn bản, không phải giáo án).
type sampleLesson struct {
	Title           string
	DurationMinutes int
	Body            string
}

// sampleOption là một phương án của câu hỏi trắc nghiệm mẫu.
type sampleOption struct {
	Content   string
	IsCorrect bool
}

// sampleQuestion là một câu hỏi trắc nghiệm mẫu (chỉ dùng single_choice cho đơn giản).
type sampleQuestion struct {
	Prompt      string
	Level       string
	Points      float64
	Explanation string
	Options     []sampleOption
}

// sampleQuiz là bài tập trắc nghiệm mẫu đi kèm các bài giảng của một môn.
type sampleQuiz struct {
	Title        string
	Instructions string
	Questions    []sampleQuestion
}

// sampleSubjectContent là toàn bộ nội dung mẫu (bài giảng + bài tập) của một môn học mẫu.
type sampleSubjectContent struct {
	Lessons []sampleLesson
	Quiz    sampleQuiz
}

// sampleContentByCode chứa nội dung học thật sự (không phải văn bản giữ chỗ)
// cho từng môn trong SampleSubjects, để lớp học mẫu có ngay bài giảng và bài
// tập dùng thử được chứ không chỉ là một lớp học rỗng.
var sampleContentByCode = map[string]sampleSubjectContent{
	"TOAN": {
		Lessons: []sampleLesson{
			{
				Title:           "Mệnh đề và tập hợp",
				DurationMinutes: 30,
				Body: `## 1. Mệnh đề

**Mệnh đề** là một câu khẳng định có tính đúng hoặc sai, không thể vừa đúng vừa sai.

Ví dụ: "$5$ là số nguyên tố" là một mệnh đề đúng; "$2 + 2 = 5$" là một mệnh đề sai; còn "Bạn có khỏe không?" không phải là mệnh đề vì không khẳng định đúng/sai.

## 2. Mệnh đề phủ định

Mệnh đề phủ định của mệnh đề $P$, ký hiệu $\overline{P}$ (hay "không $P$"), đúng khi $P$ sai và ngược lại.

Ví dụ: $P$: "$7$ là số chẵn" (sai) thì $\overline{P}$: "$7$ không là số chẵn" (đúng).

## 3. Mệnh đề kéo theo và mệnh đề đảo

Mệnh đề "**Nếu** $P$ **thì** $Q$" (ký hiệu $P \Rightarrow Q$) chỉ sai khi $P$ đúng mà $Q$ sai; các trường hợp còn lại đều đúng.

Mệnh đề $Q \Rightarrow P$ gọi là **mệnh đề đảo** của $P \Rightarrow Q$. Hai mệnh đề này nói chung không tương đương nhau. Nếu cả $P \Rightarrow Q$ và $Q \Rightarrow P$ đều đúng, ta viết $P \Leftrightarrow Q$ ("$P$ khi và chỉ khi $Q$").

## 4. Tập hợp và các phép toán

Cho hai tập hợp $A, B$:

- **Hợp**: $A \cup B = \{x \mid x \in A \text{ hoặc } x \in B\}$
- **Giao**: $A \cap B = \{x \mid x \in A \text{ và } x \in B\}$
- **Hiệu**: $A \setminus B = \{x \mid x \in A \text{ và } x \notin B\}$
- **Tập con**: $A \subset B$ khi mọi phần tử của $A$ đều thuộc $B$.

## 5. Ví dụ áp dụng

Cho $A = \{1, 2, 3, 4\}$ và $B = \{3, 4, 5, 6\}$. Khi đó:

$$A \cup B = \{1,2,3,4,5,6\} \qquad A \cap B = \{3,4\} \qquad A \setminus B = \{1,2\}$$

> [!TIP]
> Có thể minh hoạ các phép toán tập hợp bằng biểu đồ Ven để dễ hình dung phần giao, hợp, hiệu.`,
			},
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

Nếu $a > 0$: hàm số nghịch biến trên khoảng $\left(-\infty; -\dfrac{b}{2a}\right)$ và đồng biến trên khoảng $\left(-\dfrac{b}{2a}; +\infty\right)$; giá trị nhỏ nhất là $-\dfrac{\Delta}{4a}$ tại đỉnh.

Nếu $a < 0$: chiều biến thiên ngược lại; giá trị lớn nhất là $-\dfrac{\Delta}{4a}$ tại đỉnh.

## 4. Ví dụ áp dụng

**Ví dụ.** Xét hàm số $y = x^2 - 4x + 3$.

Ta có $a = 1, b = -4, c = 3$, nên toạ độ đỉnh: $x = -\dfrac{-4}{2 \cdot 1} = 2$, $y = 2^2 - 4\cdot2 + 3 = -1$.

Vậy đỉnh $I(2; -1)$, trục đối xứng $x = 2$; vì $a = 1 > 0$ nên hàm số có giá trị nhỏ nhất là $-1$.

## 5. Ứng dụng

Hàm số bậc hai mô tả nhiều hiện tượng thực tế: quỹ đạo vật ném xiên (bỏ qua sức cản không khí), lợi nhuận theo giá bán, hình dạng cổng parabol...

> [!NOTE]
> Dấu của $\Delta = b^2 - 4ac$ cho biết số giao điểm của đồ thị với trục hoành: $\Delta > 0$ (2 giao điểm), $\Delta = 0$ (1 giao điểm — đỉnh nằm trên trục hoành), $\Delta < 0$ (không có giao điểm).`,
			},
			{
				Title:           "Bất phương trình bậc nhất hai ẩn",
				DurationMinutes: 30,
				Body: `## 1. Bất phương trình bậc nhất hai ẩn

Bất phương trình bậc nhất hai ẩn $x, y$ có dạng $ax + by \le c$ (hoặc $<, \ge, >$), trong đó $a, b, c$ là các số cho trước, $a, b$ không đồng thời bằng $0$.

## 2. Miền nghiệm

Mỗi cặp số $(x_0; y_0)$ thoả mãn bất phương trình gọi là một nghiệm. Tập hợp các điểm biểu diễn nghiệm trên mặt phẳng toạ độ gọi là **miền nghiệm** của bất phương trình.

## 3. Cách xác định miền nghiệm

1. Vẽ đường thẳng $d: ax + by = c$ (đường ranh giới).
2. Lấy một điểm $M_0$ không thuộc $d$ (thường chọn gốc toạ độ $O(0;0)$ nếu $d$ không đi qua $O$) thay vào bất phương trình.
3. Nếu $M_0$ thoả mãn thì miền nghiệm là nửa mặt phẳng bờ $d$ chứa $M_0$; nếu không thì miền nghiệm là nửa mặt phẳng còn lại.

## 4. Ví dụ áp dụng

Xác định miền nghiệm của bất phương trình $x + 2y \le 4$.

Đường ranh giới: $x + 2y = 4$. Thay $O(0;0)$: $0 + 0 = 0 \le 4$ (đúng), nên miền nghiệm là nửa mặt phẳng bờ $d$ chứa gốc toạ độ (kể cả đường thẳng $d$ vì dấu $\le$).

## 5. Hệ bất phương trình bậc nhất hai ẩn

Miền nghiệm của một hệ bất phương trình bậc nhất hai ẩn là giao của các miền nghiệm từng bất phương trình trong hệ — dùng để giải các bài toán tối ưu (quy hoạch tuyến tính) trong thực tế: lập kế hoạch sản xuất, phân bổ nguồn lực sao cho lợi nhuận cao nhất.

> [!TIP]
> Đường ranh giới vẽ nét liền khi bất phương trình có dấu $\le$ hoặc $\ge$ (nghiệm bao gồm cả biên); vẽ nét đứt khi dấu là $<$ hoặc $>$.`,
			},
			{
				Title:           "Dấu của tam thức bậc hai",
				DurationMinutes: 30,
				Body: `## 1. Tam thức bậc hai

Tam thức bậc hai (theo $x$) là biểu thức có dạng $f(x) = ax^2 + bx + c$ với $a \neq 0$.

## 2. Định lý về dấu của tam thức bậc hai

Xét $\Delta = b^2 - 4ac$:

- Nếu $\Delta < 0$: $f(x)$ luôn cùng dấu với hệ số $a$ với mọi $x \in \mathbb{R}$.
- Nếu $\Delta = 0$: $f(x)$ cùng dấu với $a$ với mọi $x \neq -\dfrac{b}{2a}$ (bằng $0$ tại $x = -\dfrac{b}{2a}$).
- Nếu $\Delta > 0$: $f(x)$ có hai nghiệm $x_1 < x_2$; $f(x)$ trái dấu với $a$ khi $x \in (x_1; x_2)$ và cùng dấu với $a$ khi $x < x_1$ hoặc $x > x_2$.

## 3. Ví dụ áp dụng

Xét dấu của $f(x) = x^2 - 5x + 6$.

$\Delta = 25 - 24 = 1 > 0$, hai nghiệm $x_1 = 2, x_2 = 3$ (vì $a = 1$).

Vậy $f(x) > 0$ khi $x < 2$ hoặc $x > 3$; $f(x) < 0$ khi $2 < x < 3$; $f(x) = 0$ khi $x = 2$ hoặc $x = 3$.

## 4. Ứng dụng: giải bất phương trình bậc hai

Xét dấu của tam thức là công cụ chính để giải các bất phương trình dạng $ax^2 + bx + c > 0$ (hoặc $<, \le, \ge$ $0$) — chỉ cần lập bảng xét dấu rồi chọn khoảng nghiệm phù hợp với chiều bất phương trình.

> [!TIP]
> Học thuộc quy tắc "trong trái, ngoài cùng": giữa hai nghiệm thì $f(x)$ trái dấu $a$, ngoài khoảng hai nghiệm thì $f(x)$ cùng dấu $a$.`,
			},
			{
				Title:           "Phương trình chứa căn thức",
				DurationMinutes: 30,
				Body: `## 1. Dạng phương trình thường gặp

Phương trình chứa ẩn dưới dấu căn bậc hai thường gặp có dạng $\sqrt{f(x)} = g(x)$ hoặc $\sqrt{f(x)} = \sqrt{g(x)}$.

## 2. Cách giải $\sqrt{f(x)} = g(x)$

Bình phương hai vế với điều kiện $g(x) \ge 0$:

$$\sqrt{f(x)} = g(x) \iff \begin{cases} g(x) \ge 0 \\ f(x) = [g(x)]^2 \end{cases}$$

## 3. Cách giải $\sqrt{f(x)} = \sqrt{g(x)}$

$$\sqrt{f(x)} = \sqrt{g(x)} \iff \begin{cases} f(x) \ge 0 \\ f(x) = g(x) \end{cases}$$

## 4. Ví dụ áp dụng

Giải phương trình $\sqrt{2x + 3} = x$.

Điều kiện: $x \ge 0$. Bình phương: $2x + 3 = x^2 \iff x^2 - 2x - 3 = 0 \iff x = -1$ hoặc $x = 3$.

Đối chiếu điều kiện $x \ge 0$: loại $x = -1$, nhận $x = 3$. Vậy nghiệm là $x = 3$.

## 5. Lưu ý quan trọng

Bình phương hai vế có thể sinh ra **nghiệm ngoại lai** (nghiệm không thoả mãn phương trình ban đầu), nên bắt buộc phải đối chiếu điều kiện hoặc thử lại nghiệm sau khi giải.

> [!NOTE]
> Có thể thử lại bằng cách thay trực tiếp nghiệm tìm được vào phương trình gốc thay vì chỉ dựa vào điều kiện, để chắc chắn không sót sai sót khi lập điều kiện.`,
			},
			{
				Title:           "Giá trị lượng giác của một góc từ 0° đến 180°",
				DurationMinutes: 30,
				Body: `## 1. Định nghĩa

Trong mặt phẳng toạ độ, vẽ nửa đường tròn đơn vị tâm $O$, lấy điểm $M$ sao cho góc $xOM = \alpha$ ($0^\circ \le \alpha \le 180^\circ$), gọi $M(x_0; y_0)$. Khi đó:

$$\sin\alpha = y_0 \qquad \cos\alpha = x_0 \qquad \tan\alpha = \dfrac{y_0}{x_0}\ (x_0 \neq 0) \qquad \cot\alpha = \dfrac{x_0}{y_0}\ (y_0 \neq 0)$$

## 2. Một số giá trị đặc biệt

$$\sin 0^\circ = 0,\ \sin 30^\circ = \dfrac12,\ \sin 45^\circ = \dfrac{\sqrt2}{2},\ \sin 60^\circ = \dfrac{\sqrt3}{2},\ \sin 90^\circ = 1$$

$$\cos 0^\circ = 1,\ \cos 30^\circ = \dfrac{\sqrt3}{2},\ \cos 45^\circ = \dfrac{\sqrt2}{2},\ \cos 60^\circ = \dfrac12,\ \cos 90^\circ = 0$$

## 3. Quan hệ giữa các góc bù nhau

Với $0^\circ \le \alpha \le 180^\circ$:

$$\sin(180^\circ - \alpha) = \sin\alpha \qquad \cos(180^\circ - \alpha) = -\cos\alpha$$

Đây là lý do $\sin\alpha$ luôn không âm còn $\cos\alpha$ có thể âm khi $\alpha$ là góc tù ($90^\circ < \alpha < 180^\circ$).

## 4. Ví dụ áp dụng

Tính $\cos 120^\circ$. Ta có $120^\circ = 180^\circ - 60^\circ$, nên $\cos 120^\circ = -\cos 60^\circ = -\dfrac12$.

> [!TIP]
> Với góc tù, $\cos$ luôn âm còn $\sin$ luôn dương — nhớ điều này giúp kiểm tra nhanh dấu của kết quả.`,
			},
			{
				Title:           "Định lý côsin và định lý sin trong tam giác",
				DurationMinutes: 30,
				Body: `## 1. Định lý côsin

Trong tam giác $ABC$ với $BC = a, CA = b, AB = c$:

$$a^2 = b^2 + c^2 - 2bc\cos A$$

và các hệ thức tương tự cho $b^2, c^2$. Từ đó suy ra công thức tính góc:

$$\cos A = \dfrac{b^2 + c^2 - a^2}{2bc}$$

## 2. Định lý sin

$$\dfrac{a}{\sin A} = \dfrac{b}{\sin B} = \dfrac{c}{\sin C} = 2R$$

trong đó $R$ là bán kính đường tròn ngoại tiếp tam giác $ABC$.

## 3. Ví dụ áp dụng định lý côsin

Tam giác $ABC$ có $b = 6, c = 8, A = 60^\circ$. Tính cạnh $a$.

$$a^2 = 6^2 + 8^2 - 2\times6\times8\times\cos60^\circ = 36 + 64 - 48 = 52 \Rightarrow a = \sqrt{52} \approx 7{,}21$$

## 4. Ví dụ áp dụng định lý sin

Tam giác $ABC$ có $a = 10, A = 30^\circ$. Tính bán kính đường tròn ngoại tiếp $R$.

$$2R = \dfrac{a}{\sin A} = \dfrac{10}{\sin30^\circ} = \dfrac{10}{0{,}5} = 20 \Rightarrow R = 10$$

## 5. Khi nào dùng định lý nào?

Định lý côsin dùng khi biết hai cạnh và góc xen giữa (hoặc biết ba cạnh để tính góc); định lý sin dùng khi biết một cạnh và các góc đối diện, hoặc cần tính bán kính đường tròn ngoại tiếp.

> [!NOTE]
> Định lý côsin là mở rộng của định lý Pythagoras: khi $A = 90^\circ$ thì $\cos A = 0$ và công thức trở thành $a^2 = b^2 + c^2$.`,
			},
			{
				Title:           "Vectơ và các phép toán vectơ",
				DurationMinutes: 30,
				Body: `## 1. Khái niệm vectơ

Vectơ là một đoạn thẳng có hướng, ký hiệu $\vec{AB}$ (điểm đầu $A$, điểm cuối $B$). Độ dài đoạn thẳng $AB$ gọi là **độ dài** (hay **độ lớn**) của vectơ, ký hiệu $|\vec{AB}|$.

Hai vectơ **cùng phương** khi giá của chúng song song hoặc trùng nhau; **cùng hướng** / **ngược hướng** khi cùng phương và chiều mũi tên giống/khác nhau. Hai vectơ **bằng nhau** khi cùng hướng và cùng độ dài.

## 2. Phép cộng và phép trừ vectơ

**Quy tắc ba điểm**: $\vec{AB} + \vec{BC} = \vec{AC}$.

**Quy tắc hình bình hành**: nếu $ABCD$ là hình bình hành thì $\vec{AB} + \vec{AD} = \vec{AC}$.

**Phép trừ**: $\vec{AB} - \vec{AC} = \vec{CB}$.

## 3. Phép nhân vectơ với một số

Với số thực $k$ và vectơ $\vec{a}$, vectơ $k\vec{a}$ có độ dài $|k|\cdot|\vec{a}|$, cùng hướng với $\vec{a}$ nếu $k > 0$, ngược hướng nếu $k < 0$.

## 4. Ví dụ áp dụng

Cho tam giác $ABC$, gọi $M$ là trung điểm $BC$. Khi đó $\vec{AM} = \dfrac{1}{2}(\vec{AB} + \vec{AC})$ — công thức trung điểm quen thuộc, dùng nhiều trong chứng minh hình học bằng vectơ.

## 5. Ứng dụng

Vectơ là công cụ để chứng minh các quan hệ hình học (thẳng hàng, song song, trung điểm, trọng tâm...) mà không cần vẽ hình phức tạp, đồng thời là nền tảng để học vật lý (lực, vận tốc, gia tốc đều là đại lượng vectơ).

> [!TIP]
> Vectơ $\vec{0}$ (vectơ-không) có độ dài bằng $0$ và cùng phương, cùng hướng với mọi vectơ.`,
			},
			{
				Title:           "Tích vô hướng của hai vectơ",
				DurationMinutes: 30,
				Body: `## 1. Định nghĩa

Tích vô hướng của hai vectơ $\vec{a}$ và $\vec{b}$ khác $\vec{0}$ là một số, ký hiệu $\vec{a}\cdot\vec{b}$, xác định bởi:

$$\vec{a}\cdot\vec{b} = |\vec{a}|\cdot|\vec{b}|\cdot\cos(\vec{a}, \vec{b})$$

Nếu $\vec{a} = \vec{0}$ hoặc $\vec{b} = \vec{0}$ thì quy ước $\vec{a}\cdot\vec{b} = 0$.

## 2. Biểu thức toạ độ

Trong hệ trục toạ độ $Oxy$, nếu $\vec{a} = (x_1; y_1)$ và $\vec{b} = (x_2; y_2)$ thì:

$$\vec{a}\cdot\vec{b} = x_1x_2 + y_1y_2$$

## 3. Ứng dụng: tính góc giữa hai vectơ

$$\cos(\vec{a}, \vec{b}) = \dfrac{\vec{a}\cdot\vec{b}}{|\vec{a}|\cdot|\vec{b}|}$$

## 4. Điều kiện vuông góc

Hai vectơ $\vec{a}, \vec{b}$ (khác $\vec{0}$) vuông góc với nhau khi và chỉ khi $\vec{a}\cdot\vec{b} = 0$.

## 5. Ví dụ áp dụng

Cho $\vec{a} = (2; 1)$ và $\vec{b} = (-1; 2)$. Ta có $\vec{a}\cdot\vec{b} = 2\times(-1) + 1\times2 = 0$, vậy $\vec{a} \perp \vec{b}$.

> [!NOTE]
> Tích vô hướng còn dùng để tính độ dài vectơ: $|\vec{a}| = \sqrt{\vec{a}\cdot\vec{a}}$.`,
			},
			{
				Title:           "Các số đặc trưng đo xu thế trung tâm",
				DurationMinutes: 30,
				Body: `## 1. Số trung bình

Với mẫu số liệu $x_1, x_2, \ldots, x_n$, số trung bình là:

$$\bar{x} = \dfrac{x_1 + x_2 + \cdots + x_n}{n}$$

## 2. Trung vị

Sắp xếp mẫu số liệu theo thứ tự không giảm. **Trung vị** ($M_e$) là giá trị ở chính giữa dãy: nếu $n$ lẻ, trung vị là số ở vị trí $\dfrac{n+1}{2}$; nếu $n$ chẵn, trung vị là trung bình cộng của hai số ở giữa dãy.

## 3. Mốt

**Mốt** ($M_o$) là giá trị xuất hiện với tần số lớn nhất trong mẫu số liệu. Một mẫu số liệu có thể có một, nhiều mốt, hoặc không có mốt.

## 4. Ví dụ áp dụng

Điểm kiểm tra của 7 học sinh: $6, 7, 7, 8, 8, 8, 9$.

- Số trung bình: $\bar{x} = \dfrac{6+7+7+8+8+8+9}{7} = \dfrac{53}{7} \approx 7{,}57$.
- Trung vị: dãy đã sắp xếp, $n = 7$ lẻ, trung vị là số ở vị trí thứ $4$, tức $M_e = 8$.
- Mốt: giá trị $8$ xuất hiện nhiều nhất ($3$ lần), nên $M_o = 8$.

## 5. Khi nào dùng số đặc trưng nào?

Số trung bình phản ánh tốt xu thế chung nhưng nhạy với giá trị bất thường; trung vị ít bị ảnh hưởng bởi giá trị bất thường hơn; mốt phù hợp với dữ liệu định tính hoặc khi cần biết giá trị phổ biến nhất.

> [!TIP]
> Khi mẫu số liệu có giá trị quá lớn hoặc quá nhỏ bất thường (outlier), trung vị thường phản ánh xu thế trung tâm chính xác hơn số trung bình.`,
			},
		},
		Quiz: sampleQuiz{
			Title:        "Bài tập: Mệnh đề, tập hợp và hàm số bậc hai",
			Instructions: "Làm bài để tự kiểm tra kiến thức hai bài giảng vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Mệnh đề nào sau đây là mệnh đề đúng?",
					Level:  "Nhận biết", Points: 2.5,
					Explanation: "5 chỉ chia hết cho 1 và chính nó nên là số nguyên tố; các mệnh đề còn lại đều sai.",
					Options: []sampleOption{
						{Content: "5 là số nguyên tố", IsCorrect: true},
						{Content: "2 + 2 = 5"},
						{Content: "7 là số chẵn"},
						{Content: "10 nhỏ hơn 5"},
					},
				},
				{
					Prompt: "Cho A = {1, 2, 3, 4} và B = {3, 4, 5, 6}. Tập A ∩ B là?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "Giao của hai tập hợp gồm các phần tử chung: 3 và 4.",
					Options: []sampleOption{
						{Content: "{3, 4}", IsCorrect: true},
						{Content: "{1, 2, 5, 6}"},
						{Content: "{1, 2, 3, 4, 5, 6}"},
						{Content: "{1, 2}"},
					},
				},
				{
					Prompt: "Toạ độ đỉnh của parabol y = x² - 4x + 3 là?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "x = -b/2a = -(-4)/2 = 2; y = 2² - 4×2 + 3 = -1. Vậy đỉnh là (2; -1).",
					Options: []sampleOption{
						{Content: "(2; -1)", IsCorrect: true},
						{Content: "(2; 1)"},
						{Content: "(-2; -1)"},
						{Content: "(4; 3)"},
					},
				},
				{
					Prompt: "Hàm số y = -2x² + 4x + 1 có giá trị lớn nhất bằng bao nhiêu?",
					Level:  "Vận dụng", Points: 2.5,
					Explanation: "a = -2 < 0 nên hàm số có giá trị lớn nhất tại đỉnh: x = -4/(2×-2) = 1, y = -2(1)² + 4(1) + 1 = 3.",
					Options: []sampleOption{
						{Content: "3", IsCorrect: true},
						{Content: "1"},
						{Content: "-3"},
						{Content: "24"},
					},
				},
			},
		},
	},
	"NGUVAN": {
		Lessons: []sampleLesson{
			{
				Title:           "Kỹ năng viết bài văn nghị luận về một vấn đề xã hội",
				DurationMinutes: 30,
				Body: `## 1. Nghị luận về một vấn đề xã hội là gì?

Đây là kiểu bài yêu cầu người viết trình bày ý kiến, lập luận về một hiện tượng đời sống hoặc một tư tưởng, đạo lý, dựa trên lí lẽ và dẫn chứng thuyết phục — một dạng bài quen thuộc trong chương trình Ngữ văn lớp 10.

## 2. Bố cục bài viết

1. **Mở bài**: Giới thiệu và nêu vấn đề cần bàn luận.
2. **Thân bài**:
   - Giải thích khái niệm, nội dung của vấn đề (nếu cần).
   - Phân tích, bàn luận các khía cạnh của vấn đề, có lí lẽ và dẫn chứng cụ thể.
   - Xem xét vấn đề từ nhiều phía, phản biện những ý kiến trái chiều (nếu có).
   - Rút ra bài học nhận thức và hành động cho bản thân.
3. **Kết bài**: Khẳng định lại ý kiến, để lại ấn tượng, thông điệp.

## 3. Một số lưu ý khi viết

- Xác định đúng và bám sát yêu cầu của đề (vấn đề gì, phạm vi dẫn chứng nào).
- Lập luận cần chặt chẽ: mỗi luận điểm có lí lẽ và dẫn chứng đi kèm.
- Dẫn chứng nên cụ thể, xác thực, tránh nêu chung chung.
- Diễn đạt rõ ràng, mạch lạc; có thể dùng câu hỏi tu từ, trích dẫn để tăng sức thuyết phục.

## 4. Ví dụ đề bài

Đề: *Viết bài văn nghị luận trình bày suy nghĩ về ý nghĩa của tinh thần tự học đối với học sinh.*

Gợi ý triển khai: mở bài nêu vấn đề "tự học" → giải thích tự học là gì → bàn luận: tự học giúp chủ động tiếp thu kiến thức, rèn tính kỷ luật, dẫn chứng những tấm gương tự học thành công → phản đề: phê phán thái độ ỷ lại, lười tự học → bài học: xây dựng thói quen tự học mỗi ngày → kết bài khẳng định lại vai trò của tự học.

> [!TIP]
> Trước khi viết, hãy lập dàn ý nhanh ra giấy nháp theo cấu trúc 3 phần trên để bài viết không thiếu ý hoặc lạc đề.`,
			},
			{
				Title:           "Đọc hiểu bài thơ Cảnh ngày hè — Nguyễn Trãi",
				DurationMinutes: 30,
				Body: `## 1. Văn bản

**Cảnh ngày hè** (trích *Bảo kính cảnh giới*, bài 43) — Nguyễn Trãi

> Rồi hóng mát thuở ngày trường,
> Hoè lục đùn đùn tán rợp giương.
> Thạch lựu hiên còn phun thức đỏ,
> Hồng liên trì đã tiễn mùi hương.
> Lao xao chợ cá làng ngư phủ,
> Dắng dỏi cầm ve lầu tịch dương.
> Dẽ có Ngu cầm đàn một tiếng,
> Dân giàu đủ khắp đòi phương.

## 2. Vài nét về tác giả và tác phẩm

Nguyễn Trãi (1380–1442) là nhà chính trị, quân sự, ngoại giao và là nhà thơ lớn của dân tộc. *Quốc âm thi tập* của ông là tập thơ Nôm sớm nhất còn lại của văn học Việt Nam. "Cảnh ngày hè" nằm trong chùm thơ *Bảo kính cảnh giới* (gương báu răn mình), viết theo thể **thất ngôn xen lục ngôn** — nét sáng tạo riêng của Nguyễn Trãi so với thơ Đường luật thuần túy.

## 3. Phân tích

**Bức tranh thiên nhiên ngày hè (6 câu đầu)**: Cảnh vật hiện lên sống động, rực rỡ qua các động từ mạnh và tính từ chỉ màu sắc: hòe "đùn đùn" tán lá rợp xanh, thạch lựu "phun" sắc đỏ trước hiên, sen hồng toả hương nơi ao. Âm thanh cuộc sống hiện diện qua tiếng "lao xao" của chợ cá và tiếng ve kêu "dắng dỏi" lúc chiều tà — bức tranh vừa có màu sắc, hương thơm, vừa có âm thanh, tràn đầy sức sống.

**Tâm sự của nhà thơ (2 câu cuối)**: Từ cảnh sắc thanh bình, no ấm ấy, tác giả ước có được cây đàn của vua Ngu Thuấn để gảy khúc "Nam phong", cầu mong nhân dân khắp nơi được ấm no, giàu đủ. Đó là tấm lòng của một người luôn đau đáu tâm sự vì dân vì nước, dù đang sống ẩn dật.

## 4. Giá trị nội dung và nghệ thuật

- **Nội dung**: Bài thơ vừa là bức tranh thiên nhiên, cuộc sống ngày hè tràn đầy sức sống, vừa thể hiện tình yêu thiên nhiên và tấm lòng yêu nước thương dân sâu nặng của Nguyễn Trãi.
- **Nghệ thuật**: Thể thơ thất ngôn xen lục ngôn (câu 1 và câu 8 chỉ có 6 chữ) phá cách so với thơ Đường luật; hệ thống động từ, tính từ giàu sức gợi tả; hình ảnh gần gũi, bình dị của thiên nhiên và cuộc sống lao động.

> [!NOTE]
> Câu thơ đầu và câu cuối chỉ có 6 chữ (lục ngôn) thay vì 7 chữ như các câu còn lại — đây là dấu ấn riêng, thể hiện sự sáng tạo của Nguyễn Trãi khi Việt hoá thể thơ Đường luật.`,
			},
			{
				Title:           "Thần thoại Thần Trụ Trời",
				DurationMinutes: 30,
				Body: `## 1. Thần thoại là gì?

Thần thoại là thể loại tự sự dân gian ra đời sớm nhất, kể về các vị thần nhằm giải thích nguồn gốc vũ trụ, thế giới tự nhiên và con người theo quan niệm của người xưa. "Thần Trụ Trời" là một thần thoại suy nguyên (giải thích nguồn gốc) tiêu biểu của người Việt.

## 2. Tóm tắt truyện

Thuở khai thiên lập địa, trời đất còn là một khối hỗn độn, tối tăm. Có một vị thần khổng lồ (thần Trụ Trời) xuất hiện, tự mình đội trời lên rồi đào đất, đá đắp thành một cái cột chống trời, tách trời và đất ra làm hai. Khi trời đã cao và khô cứng, thần phá cột đi, đất đá văng ra bốn phương thành núi, đảo, cao nguyên; chỗ thần đào đất thành biển. Sau đó, các vị thần khác tiếp tục công việc: thần Sao, thần Sông, thần Biển... để hoàn thiện thế giới.

## 3. Đặc điểm nghệ thuật của thần thoại suy nguyên

- **Nhân vật** là các vị thần có sức mạnh phi thường, hành động mang tính sáng tạo vũ trụ.
- **Không gian, thời gian** phiếm chỉ, mơ hồ ("thuở xưa", "chưa có thế gian").
- **Cách giải thích** hiện tượng tự nhiên (núi, biển, cột chống trời...) mang tính hồn nhiên, chất phác — phản ánh tư duy của người nguyên thủy khi lý giải thế giới.

## 4. Giá trị

Thần thoại "Thần Trụ Trời" thể hiện khát vọng nhận thức, lý giải thế giới tự nhiên và trí tưởng tượng phong phú của người Việt cổ, đồng thời là nguồn tư liệu quý về tín ngưỡng, vũ trụ quan dân gian.

> [!NOTE]
> Thần thoại suy nguyên khác thần thoại sáng tạo (kể về phát minh văn hoá như lửa, nghề dệt...) — cả hai đều thuộc dòng thần thoại nhưng chức năng giải thích khác nhau.`,
			},
			{
				Title:           "Sử thi Đăm Săn — đoạn trích Chiến thắng Mtao Mxây",
				DurationMinutes: 30,
				Body: `## 1. Giới thiệu chung

*Đăm Săn* là bộ sử thi anh hùng nổi tiếng của dân tộc Ê-đê (Tây Nguyên), kể về cuộc đời và chiến công của tù trưởng Đăm Săn. Đoạn trích "Chiến thắng Mtao Mxây" kể lại việc Đăm Săn giao chiến với Mtao Mxây để cứu vợ là Hơ Nhị bị bắt đi.

## 2. Tóm tắt đoạn trích

Mtao Mxây cướp vợ của Đăm Săn. Đăm Săn đến nhà Mtao Mxây khiêu chiến. Ban đầu Mtao Mxây tỏ ra kiêu ngạo, múa khiên yếu ớt; Đăm Săn múa khiên mạnh mẽ, uy dũng khiến Mtao Mxây khiếp sợ phải cầu cứu thần linh. Cuối cùng, nhờ sự giúp đỡ của ông Trời (ném cho miếng trầu và chỉ cách dùng chày mòn), Đăm Săn đánh bại và giết chết Mtao Mxây, cứu được vợ, thu phục dân làng và của cải của kẻ thù.

## 3. Đặc điểm nghệ thuật sử thi

- **Nhân vật anh hùng** mang vẻ đẹp lý tưởng của cộng đồng: dũng mãnh, tài giỏi, được thần linh trợ giúp.
- **Ngôn ngữ** giàu hình ảnh so sánh, phóng đại (miêu tả sức mạnh, hình dáng nhân vật rất kỳ vĩ).
- **Không gian sử thi** rộng lớn, mang tính cộng đồng (chiến đấu không chỉ vì cá nhân mà vì danh dự, sự thịnh vượng của cả buôn làng).

## 4. Ý nghĩa

Đoạn trích ca ngợi vẻ đẹp, sức mạnh và tài năng của người anh hùng sử thi, đồng thời thể hiện khát vọng về một cuộc sống ấm no, đoàn kết, thịnh vượng của cộng đồng các dân tộc Tây Nguyên xưa.

> [!TIP]
> Khi đọc sử thi, chú ý các đoạn miêu tả bằng thủ pháp so sánh trùng điệp, phóng đại — đây là dấu hiệu đặc trưng để nhận diện thể loại.`,
			},
			{
				Title:           "Bài thơ Tỏ lòng (Thuật hoài) — Phạm Ngũ Lão",
				DurationMinutes: 30,
				Body: `## 1. Văn bản (dịch thơ)

**Tỏ lòng** (nguyên tác chữ Hán: *Thuật hoài*) — Phạm Ngũ Lão

> Múa giáo non sông trải mấy thu,
> Ba quân khí mạnh nuốt trôi trâu.
> Công danh nam tử còn vương nợ,
> Luống thẹn tai nghe chuyện Vũ hầu.

## 2. Vài nét về tác giả

Phạm Ngũ Lão (1255–1320) là tướng lĩnh nhà Trần, có nhiều công lao trong hai cuộc kháng chiến chống quân Nguyên – Mông. Bài thơ được viết bằng chữ Hán theo thể thất ngôn tứ tuyệt Đường luật, thể hiện khí phách của một vị tướng thời Trần.

## 3. Phân tích

**Hai câu đầu**: khắc hoạ tầm vóc, khí thế hào hùng của người tráng sĩ và quân đội nhà Trần: hình ảnh "múa giáo" trấn giữ non sông qua năm tháng, sức mạnh "ba quân" được so sánh phóng đại "nuốt trôi trâu" — thể hiện "hào khí Đông A" (hào khí thời Trần) ngút trời.

**Hai câu sau**: là nỗi lòng riêng của tác giả — "công danh" theo quan niệm "chí làm trai" thời phong kiến là phải lập công giúp nước; nhắc đến Vũ hầu (Gia Cát Lượng — bậc mưu thần lập nhiều công lớn thời Tam Quốc) để tự thấy hổ thẹn vì công danh của mình còn "vương nợ", chưa sánh được — qua đó bộc lộ khát vọng lập công lớn, ý thức trách nhiệm với đất nước.

## 4. Giá trị nội dung và nghệ thuật

- **Nội dung**: Bài thơ thể hiện lý tưởng và khí phách anh hùng, tinh thần trách nhiệm với đất nước của trang nam nhi thời Trần.
- **Nghệ thuật**: Bút pháp hoành tráng, giàu tính biểu tượng, ngôn ngữ hàm súc, sử dụng điển tích (Vũ hầu) — đặc trưng thơ trung đại.

> [!NOTE]
> "Hào khí Đông A" là cách gọi tinh thần hào hùng, quật cường của quân dân thời Trần, thường được giải thích là do chữ Hán "Trần" (陳) có thể tách thành hai chữ "Đông" (東) và "A" (阿).`,
			},
			{
				Title:           "Đoạn trích Chí khí anh hùng (Truyện Kiều) — Nguyễn Du",
				DurationMinutes: 30,
				Body: `## 1. Vị trí đoạn trích

"Chí khí anh hùng" là đoạn trích trong *Truyện Kiều* của Nguyễn Du, kể về việc Từ Hải — một người anh hùng "đầu đội trời, chân đạp đất" — từ biệt Thuý Kiều để ra đi lập nghiệp lớn sau nửa năm chung sống.

## 2. Một đoạn tiêu biểu

> Trượng phu thoắt đã động lòng bốn phương,
> Trông vời trời bể mênh mang,
> Thanh gươm yên ngựa lên đường thẳng rong.

## 3. Phân tích

**Hình tượng Từ Hải**: Hiện lên qua chí khí "động lòng bốn phương" — khát vọng tung hoành thiên hạ, không cam chịu cuộc sống êm đềm, quyến luyến thê nhi. Không gian "trời bể mênh mang" tương xứng với tầm vóc phi thường của người anh hùng.

**Cuộc chia tay với Thuý Kiều**: Trước lời xin đi theo của Kiều, Từ Hải khuyên nàng ở lại chờ ngày mình "làm nên", thể hiện bản lĩnh dứt khoát của người tráng sĩ, đồng thời cũng bộc lộ tình cảm trân trọng dành cho Kiều qua lời hẹn ước rõ ràng, quả quyết.

## 4. Giá trị nội dung và nghệ thuật

- **Nội dung**: Đoạn trích khắc hoạ lý tưởng người anh hùng mang khát vọng tự do, công danh lớn lao — hình tượng khác biệt hẳn với các nhân vật khác trong *Truyện Kiều*.
- **Nghệ thuật**: Bút pháp lý tưởng hoá nhân vật (ước lệ, phóng đại) đặc trưng khi Nguyễn Du xây dựng hình tượng anh hùng; ngôn ngữ trang trọng, hình ảnh kỳ vĩ mang màu sắc sử thi.

> [!TIP]
> Từ Hải là nhân vật hiếm hoi trong *Truyện Kiều* được xây dựng theo bút pháp lý tưởng hoá thay vì bút pháp tả thực — nên hình ảnh của chàng luôn gắn với không gian rộng lớn, kỳ vĩ.`,
			},
			{
				Title:           "Chèo cổ — đoạn trích Thị Mầu lên chùa",
				DurationMinutes: 30,
				Body: `## 1. Giới thiệu chung

*Quan Âm Thị Kính* là một vở chèo cổ nổi tiếng trong kho tàng sân khấu dân gian Việt Nam. Đoạn trích "Thị Mầu lên chùa" khắc hoạ nhân vật Thị Mầu — một cô gái lẳng lơ, phá cách — lên chùa ve vãn tiểu Kính Tâm (thực chất là Thị Kính giả trai đi tu).

## 2. Tóm tắt đoạn trích

Thị Mầu lên chùa vãn cảnh, trông thấy tiểu Kính Tâm liền nảy sinh tình ý, buông lời chọc ghẹo, ve vãn công khai bất chấp lễ giáo. Kính Tâm một mực giữ ý, khước từ, chuyên tâm tụng kinh niệm Phật. Thị Mầu không nản lòng, càng thể hiện sự táo bạo, bất chấp khuôn phép của lễ giáo phong kiến.

## 3. Nhân vật Thị Mầu

Thị Mầu là kiểu nhân vật "nữ lệch" (đào lệch) đặc trưng của chèo — đối lập với hình mẫu người phụ nữ đoan trang, khuôn phép. Qua ngôn ngữ, cử chỉ táo bạo, hồn nhiên của Thị Mầu, vở chèo thể hiện khát vọng tự do yêu đương, đồng thời phê phán những trói buộc khắc nghiệt của lễ giáo phong kiến đối với người phụ nữ.

## 4. Đặc trưng nghệ thuật chèo

- Kết hợp giữa **lời thoại**, **hát** (các làn điệu chèo) và **động tác múa** trên sân khấu ước lệ.
- Nhân vật được phân loại theo **loại vai** cố định (đào, kép, mụ, lão, hề...), mỗi loại có lối diễn, giọng hát riêng.
- Đan xen yếu tố **trữ tình** và **trào lộng** (hài hước, châm biếm) — Thị Mầu là nhân vật mang đậm chất trào lộng dân gian.

> [!NOTE]
> Chèo là loại hình sân khấu dân gian gắn với hội làng ở đồng bằng Bắc Bộ, khác với tuồng (thường gắn với đề tài cung đình, lịch sử) về đề tài và phong cách biểu diễn.`,
			},
			{
				Title:           "Kỹ năng đọc hiểu văn bản thông tin",
				DurationMinutes: 30,
				Body: `## 1. Văn bản thông tin là gì?

Văn bản thông tin là loại văn bản cung cấp tri thức, dữ liệu về một sự vật, hiện tượng nhằm mục đích chính là truyền đạt thông tin khách quan, khác với văn bản văn học (thiên về biểu đạt cảm xúc, hư cấu).

## 2. Đặc điểm của văn bản thông tin

- **Nhan đề, đề mục**: thường rõ ràng, khái quát nội dung.
- **Sa-pô (đoạn mở đầu)**: tóm tắt ý chính của toàn văn bản.
- **Số liệu, hình ảnh, sơ đồ, chú thích**: hỗ trợ minh hoạ, tăng độ tin cậy.
- **Ngôn ngữ**: chính xác, khách quan, ít dùng biện pháp tu từ.

## 3. Các bước đọc hiểu văn bản thông tin

1. Đọc nhan đề, sa-pô, đề mục để nắm khái quát nội dung.
2. Xác định thông tin chính (ai, cái gì, khi nào, ở đâu, vì sao, như thế nào).
3. Chú ý cách trình bày: trật tự thời gian, quan hệ nhân quả, so sánh - đối chiếu, hay phân loại.
4. Đánh giá độ tin cậy: nguồn thông tin, số liệu có được dẫn chứng rõ ràng không.

## 4. Ví dụ áp dụng

Khi đọc một văn bản thuyết minh về lễ hội truyền thống, cần chú ý: sa-pô giới thiệu tên và ý nghĩa lễ hội; các đề mục nhỏ có thể trình bày theo trình tự thời gian diễn ra lễ hội (phần lễ, phần hội); số liệu (thời gian, địa điểm, số lượng người tham gia) giúp tăng tính xác thực.

> [!TIP]
> Với văn bản thông tin có sơ đồ hoặc bảng biểu, hãy đọc kỹ phần chú thích trước khi diễn giải số liệu để tránh hiểu sai.`,
			},
			{
				Title:           "Kỹ năng viết bài luận thuyết phục từ bỏ một thói quen hay quan niệm",
				DurationMinutes: 30,
				Body: `## 1. Mục đích của kiểu bài

Bài luận thuyết phục nhằm tác động đến nhận thức, thái độ của người đọc/người nghe để họ từ bỏ một thói quen chưa tốt hoặc một quan niệm chưa đúng, hướng tới thay đổi tích cực.

## 2. Bố cục bài viết

1. **Mở bài**: Nêu thói quen/quan niệm cần thuyết phục từ bỏ.
2. **Thân bài**:
   - Chỉ ra tác hại, hạn chế của thói quen/quan niệm đó (có dẫn chứng, lí lẽ cụ thể).
   - Đề xuất giải pháp, hướng thay đổi tích cực.
   - Dự đoán và phản hồi những ý kiến phản đối có thể có.
3. **Kết bài**: Khẳng định lại quan điểm, kêu gọi hành động.

## 3. Một số lưu ý khi viết

- Chọn lí lẽ, dẫn chứng có sức thuyết phục, tránh áp đặt cảm tính.
- Giọng văn chân thành, tôn trọng người đọc, tránh lên án gay gắt.
- Có thể dùng câu hỏi tu từ, hình ảnh so sánh để tăng hiệu quả thuyết phục.

## 4. Ví dụ đề bài

Đề: *Viết bài luận thuyết phục các bạn từ bỏ thói quen thức khuya sử dụng điện thoại.*

Gợi ý: nêu thói quen → chỉ ra tác hại (ảnh hưởng sức khoẻ, kết quả học tập, giấc ngủ) → đề xuất giải pháp (đặt giờ giới nghiêm dùng điện thoại, thay bằng hoạt động thư giãn khác) → phản hồi ý kiến "chỉ dùng một chút không sao" bằng dẫn chứng khoa học → kêu gọi thay đổi.

> [!TIP]
> Thuyết phục hiệu quả nhất khi kết hợp cả lí lẽ (logic) và cảm xúc (sự đồng cảm, thấu hiểu) chứ không chỉ đơn thuần liệt kê tác hại.`,
			},
			{
				Title:           "Kỹ năng tóm tắt văn bản",
				DurationMinutes: 30,
				Body: `## 1. Tóm tắt văn bản là gì?

Tóm tắt là trình bày lại ngắn gọn nội dung chính của một văn bản gốc, giữ đúng ý chính và tinh thần của văn bản, bằng lời văn của người tóm tắt.

## 2. Các bước tóm tắt văn bản

1. Đọc kỹ toàn bộ văn bản để nắm nội dung tổng quát.
2. Xác định luận điểm/sự kiện chính, bỏ qua chi tiết phụ, ví dụ minh hoạ.
3. Sắp xếp lại các ý chính theo trình tự hợp lý (thường theo trình tự văn bản gốc).
4. Viết lại bằng lời văn ngắn gọn, súc tích, đảm bảo đúng nội dung gốc, không thêm ý kiến chủ quan.

## 3. Yêu cầu đối với một bản tóm tắt tốt

- Đầy đủ ý chính, không bỏ sót luận điểm quan trọng.
- Ngắn gọn hơn đáng kể so với văn bản gốc (thường bằng 1/4 đến 1/3 độ dài).
- Khách quan, trung thành với nội dung và quan điểm của văn bản gốc, không xen ý kiến cá nhân.
- Diễn đạt mạch lạc, các ý liên kết chặt chẽ.

## 4. Phân biệt tóm tắt và phân tích

Tóm tắt chỉ tái hiện lại nội dung chính một cách khách quan; còn phân tích, bình luận đòi hỏi người viết đưa ra đánh giá, nhận xét, lí giải riêng về văn bản — hai kỹ năng này thường được kết hợp khi viết bài nghị luận văn học.

> [!NOTE]
> Một mẹo hữu ích: sau khi đọc xong mỗi đoạn/phần, hãy tự đặt câu hỏi "Đoạn này nói về điều gì?" và ghi lại câu trả lời bằng một câu ngắn gọn — đó chính là các ý chính để ghép thành bản tóm tắt.`,
			},
		},
		Quiz: sampleQuiz{
			Title:        "Bài tập: Đọc hiểu và nghị luận xã hội",
			Instructions: "Làm bài để tự kiểm tra kiến thức hai bài giảng vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Bài thơ 'Cảnh ngày hè' được viết theo thể thơ nào?",
					Level:  "Nhận biết", Points: 2.5,
					Explanation: "Câu 1 và câu 8 chỉ có 6 chữ (lục ngôn), khác với thất ngôn bát cú Đường luật thuần túy — đây là thể thất ngôn xen lục ngôn.",
					Options: []sampleOption{
						{Content: "Thất ngôn xen lục ngôn", IsCorrect: true},
						{Content: "Lục bát"},
						{Content: "Song thất lục bát"},
						{Content: "Thất ngôn bát cú Đường luật thuần túy"},
					},
				},
				{
					Prompt: "Hai câu thơ cuối bài 'Cảnh ngày hè' thể hiện điều gì ở Nguyễn Trãi?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "Tác giả ước có đàn của vua Ngu Thuấn để cầu mong dân chúng khắp nơi được ấm no, giàu đủ.",
					Options: []sampleOption{
						{Content: "Khát vọng nhân dân được ấm no, giàu đủ", IsCorrect: true},
						{Content: "Niềm vui thú điền viên, quên hết sự đời"},
						{Content: "Nỗi buồn vì thời gian trôi nhanh"},
						{Content: "Mong muốn được trọng dụng ở triều đình"},
					},
				},
				{
					Prompt: "Phần nào của bài văn nghị luận có nhiệm vụ giải thích, phân tích và phản biện vấn đề?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "Thân bài là phần triển khai giải thích, bàn luận, phản biện và rút ra bài học.",
					Options: []sampleOption{
						{Content: "Thân bài", IsCorrect: true},
						{Content: "Mở bài"},
						{Content: "Kết bài"},
						{Content: "Không phần nào"},
					},
				},
				{
					Prompt: "Hình ảnh 'thạch lựu hiên còn phun thức đỏ' trong bài thơ gợi tả điều gì?",
					Level:  "Vận dụng", Points: 2.5,
					Explanation: "Động từ 'phun' kết hợp màu đỏ của hoa lựu gợi sức sống mãnh liệt, rực rỡ của cảnh vật ngày hè.",
					Options: []sampleOption{
						{Content: "Màu hoa lựu đỏ rực trước hiên nhà vào mùa hè", IsCorrect: true},
						{Content: "Cảnh chợ cá tấp nập"},
						{Content: "Tiếng ve kêu râm ran"},
						{Content: "Nỗi buồn của tác giả"},
					},
				},
			},
		},
	},
	"VATLY": {
		Lessons: []sampleLesson{
			{
				Title:           "Ba định luật Newton",
				DurationMinutes: 30,
				Body: `## 1. Định luật I Newton (định luật quán tính)

Nếu một vật không chịu tác dụng của lực nào, hoặc chịu tác dụng của các lực có hợp lực bằng không, thì vật đang đứng yên sẽ tiếp tục đứng yên, vật đang chuyển động sẽ tiếp tục chuyển động thẳng đều.

Tính chất giữ nguyên trạng thái chuyển động (hoặc đứng yên) của vật gọi là **quán tính**.

## 2. Định luật II Newton

Gia tốc của một vật cùng hướng với lực tác dụng lên vật; độ lớn của gia tốc tỉ lệ thuận với độ lớn của lực và tỉ lệ nghịch với khối lượng của vật:

$$\vec{a} = \dfrac{\vec{F}}{m} \quad\text{hay}\quad \vec{F} = m\vec{a}$$

trong đó $F$ là hợp lực tác dụng lên vật (N), $m$ là khối lượng của vật (kg), $a$ là gia tốc (m/s²).

## 3. Định luật III Newton

Khi vật A tác dụng lên vật B một lực, thì vật B cũng tác dụng trở lại vật A một lực. Hai lực này có cùng giá, cùng độ lớn nhưng ngược chiều:

$$\vec{F}_{AB} = -\vec{F}_{BA}$$

Hai lực này gọi là **lực và phản lực**, xuất hiện và mất đi đồng thời, luôn cùng loại nhưng đặt vào hai vật khác nhau nên **không cân bằng nhau**.

## 4. Ví dụ áp dụng

**Ví dụ.** Một vật có khối lượng $m = 2\ kg$ chịu tác dụng của một lực không đổi $F = 6\ N$. Tính gia tốc của vật.

$$a = \dfrac{F}{m} = \dfrac{6}{2} = 3\ m/s^2$$

## 5. Ứng dụng

Ba định luật Newton là nền tảng của cơ học cổ điển, giải thích được hầu hết chuyển động trong đời sống: vì sao hành khách ngả người về sau khi xe tăng tốc (quán tính), vì sao tên lửa bay được nhờ đẩy khí ra phía sau (lực và phản lực)...

> [!NOTE]
> Định luật II Newton chỉ áp dụng đúng khi khối lượng của vật không đổi trong quá trình chuyển động.`,
			},
			{
				Title:           "Chuyển động thẳng biến đổi đều",
				DurationMinutes: 30,
				Body: `## 1. Gia tốc

Gia tốc đặc trưng cho sự thay đổi nhanh hay chậm của vận tốc, xác định bởi:

$$a = \dfrac{v - v_0}{t}$$

trong đó $v_0$ là vận tốc đầu, $v$ là vận tốc tại thời điểm $t$. Đơn vị của gia tốc là m/s².

## 2. Chuyển động thẳng biến đổi đều

Là chuyển động thẳng có độ lớn gia tốc không đổi theo thời gian. Nếu tốc độ tăng dần, đó là chuyển động **nhanh dần đều**; nếu tốc độ giảm dần, đó là chuyển động **chậm dần đều**.

## 3. Các công thức cơ bản

$$v = v_0 + at$$

$$s = v_0t + \dfrac{1}{2}at^2$$

$$v^2 - v_0^2 = 2as$$

## 4. Ví dụ áp dụng

**Ví dụ.** Một xe bắt đầu chuyển động nhanh dần đều từ trạng thái đứng yên ($v_0 = 0$) với gia tốc $a = 2\ m/s^2$. Tính vận tốc và quãng đường xe đi được sau $5$ giây.

$$v = v_0 + at = 0 + 2 \times 5 = 10\ m/s$$

$$s = v_0t + \dfrac{1}{2}at^2 = 0 + \dfrac{1}{2}\times2\times5^2 = 25\ m$$

## 5. Sự rơi tự do — một trường hợp riêng

Rơi tự do là chuyển động nhanh dần đều theo phương thẳng đứng, với gia tốc rơi tự do $g \approx 9{,}8\ m/s^2$ (thường lấy tròn $g = 10\ m/s^2$ khi tính gần đúng).

> [!TIP]
> Nếu vật chuyển động chậm dần đều, gia tốc $a$ ngược dấu với vận tốc $v_0$ — chọn đúng chiều dương giúp tránh nhầm dấu khi thay vào công thức.`,
			},
			{
				Title:           "Tốc độ và vận tốc trong chuyển động thẳng",
				DurationMinutes: 30,
				Body: `## 1. Độ dịch chuyển và quãng đường

**Quãng đường** ($s$) là độ dài thực tế mà vật đã đi được, luôn là một số không âm. **Độ dịch chuyển** ($\vec{d}$) là một đại lượng vectơ, có gốc là vị trí đầu, ngọn là vị trí cuối của vật — chỉ phụ thuộc vị trí đầu và cuối, không phụ thuộc đường đi.

Khi vật chuyển động thẳng, không đổi chiều thì độ lớn độ dịch chuyển bằng quãng đường; nếu vật đổi chiều, hai đại lượng này khác nhau.

## 2. Tốc độ

**Tốc độ trung bình** là đại lượng đặc trưng cho độ nhanh chậm của chuyển động, tính bằng:

$$v_{tb} = \dfrac{s}{t}$$

trong đó $s$ là quãng đường đi được trong thời gian $t$.

## 3. Vận tốc

**Vận tốc** là đại lượng vectơ, đặc trưng cho sự nhanh chậm và cả hướng chuyển động:

$$\vec{v} = \dfrac{\vec{d}}{t}$$

Độ lớn của vận tốc trung bình bằng $\dfrac{|\vec{d}|}{t}$, có thể khác tốc độ trung bình nếu vật đi không theo đường thẳng một chiều.

## 4. Ví dụ áp dụng

Một người đi bộ $3$ km về phía Đông rồi quay lại đi $1$ km về phía Tây, hết tất cả $1$ giờ. Quãng đường đi được là $s = 3 + 1 = 4$ km, nhưng độ dịch chuyển chỉ là $d = 3 - 1 = 2$ km (về phía Đông). Tốc độ trung bình là $4$ km/h, còn độ lớn vận tốc trung bình chỉ là $2$ km/h.

> [!NOTE]
> Đây là điểm khác biệt quan trọng của chương trình Vật lý 10 (2018) so với trước: phân biệt rõ "quãng đường – tốc độ" (đại lượng vô hướng) với "độ dịch chuyển – vận tốc" (đại lượng vectơ).`,
			},
			{
				Title:           "Sự rơi tự do",
				DurationMinutes: 30,
				Body: `## 1. Định nghĩa

Sự rơi tự do là chuyển động của một vật chỉ dưới tác dụng của trọng lực, khi bỏ qua mọi lực cản của không khí.

## 2. Đặc điểm của chuyển động rơi tự do

- Có phương thẳng đứng, chiều từ trên xuống dưới.
- Là chuyển động thẳng nhanh dần đều, với gia tốc rơi tự do $g$.
- Tại cùng một nơi trên Trái Đất, mọi vật đều rơi tự do với cùng gia tốc $g$, không phụ thuộc khối lượng.

## 3. Công thức

Chọn gốc thời gian lúc vật bắt đầu rơi ($v_0 = 0$):

$$v = gt \qquad\qquad h = \dfrac{1}{2}gt^2 \qquad\qquad v^2 = 2gh$$

Giá trị $g$ thường lấy $g \approx 9{,}8\ m/s^2$, hoặc $g = 10\ m/s^2$ khi tính gần đúng.

## 4. Ví dụ áp dụng

Thả rơi tự do một vật từ độ cao $h = 20\ m$ (lấy $g = 10\ m/s^2$). Tính thời gian rơi và vận tốc lúc chạm đất.

$$h = \dfrac{1}{2}gt^2 \Rightarrow t = \sqrt{\dfrac{2h}{g}} = \sqrt{\dfrac{2\times20}{10}} = 2\ s$$

$$v = gt = 10 \times 2 = 20\ m/s$$

## 5. Thí nghiệm kiểm chứng

Nhà bác học Galileo Galilei là người đầu tiên chứng minh (qua thí nghiệm và suy luận) rằng các vật rơi tự do với cùng gia tốc bất kể khối lượng — bác bỏ quan niệm sai lầm trước đó cho rằng vật nặng luôn rơi nhanh hơn vật nhẹ.

> [!TIP]
> Trong không khí, vật nhẹ như tờ giấy rơi chậm hơn vật nặng là do lực cản không khí, không phải do bản chất sự rơi tự do lý tưởng.`,
			},
			{
				Title:           "Chuyển động ném ngang",
				DurationMinutes: 30,
				Body: `## 1. Định nghĩa

Chuyển động ném ngang là chuyển động của một vật được ném với vận tốc ban đầu theo phương nằm ngang, chỉ chịu tác dụng của trọng lực (bỏ qua sức cản không khí).

## 2. Phân tích chuyển động

Chuyển động ném ngang được phân tích thành hai chuyển động thành phần độc lập:

- Theo phương ngang ($Ox$): chuyển động thẳng đều với vận tốc $v_0$.
- Theo phương thẳng đứng ($Oy$): chuyển động rơi tự do.

## 3. Phương trình chuyển động

Chọn gốc toạ độ tại vị trí ném, gốc thời gian lúc bắt đầu ném:

$$x = v_0t \qquad\qquad y = \dfrac{1}{2}gt^2$$

## 4. Phương trình quỹ đạo và tầm xa

Khử $t$ từ hai phương trình trên, quỹ đạo của vật là một **parabol**:

$$y = \dfrac{g}{2v_0^2}x^2$$

Thời gian rơi (đến khi chạm đất, từ độ cao $h$): $t = \sqrt{\dfrac{2h}{g}}$. Tầm xa (theo phương ngang): $L = v_0t = v_0\sqrt{\dfrac{2h}{g}}$.

## 5. Ví dụ áp dụng

Một vật được ném ngang từ độ cao $h = 45\ m$ với vận tốc đầu $v_0 = 20\ m/s$ (lấy $g = 10\ m/s^2$). Tính thời gian rơi và tầm xa.

$$t = \sqrt{\dfrac{2\times45}{10}} = 3\ s \qquad L = 20 \times 3 = 60\ m$$

> [!NOTE]
> Thời gian rơi của vật ném ngang chỉ phụ thuộc độ cao $h$, không phụ thuộc vận tốc ném $v_0$ — vì chuyển động theo phương thẳng đứng vẫn là rơi tự do, độc lập với chuyển động ngang.`,
			},
			{
				Title:           "Trọng lực và lực căng dây",
				DurationMinutes: 30,
				Body: `## 1. Trọng lực

Trọng lực là lực hấp dẫn mà Trái Đất tác dụng lên một vật, có phương thẳng đứng, chiều hướng xuống, điểm đặt tại trọng tâm của vật:

$$P = mg$$

trong đó $P$ là trọng lượng (N), $m$ là khối lượng (kg), $g$ là gia tốc rơi tự do (m/s²).

## 2. Phân biệt khối lượng và trọng lượng

**Khối lượng** ($m$) là đại lượng đặc trưng cho mức quán tính của vật, không đổi theo vị trí. **Trọng lượng** ($P$) là độ lớn của trọng lực, phụ thuộc vào $g$ tại nơi đó — ví dụ một vật có cùng khối lượng nhưng trọng lượng trên Mặt Trăng nhỏ hơn nhiều so với trên Trái Đất vì $g$ trên Mặt Trăng nhỏ hơn.

## 3. Lực căng dây

Khi một sợi dây (không giãn, khối lượng không đáng kể) bị kéo căng, nó tác dụng vào vật ở hai đầu một lực gọi là **lực căng dây**, có phương trùng với dây, chiều hướng từ hai đầu dây vào phía trong (kéo vật lại gần dây).

## 4. Ví dụ áp dụng

Một vật khối lượng $m = 2\ kg$ được treo cân bằng bởi một sợi dây thẳng đứng (lấy $g = 10\ m/s^2$). Tính lực căng dây.

Vì vật cân bằng nên lực căng dây cân bằng với trọng lực: $T = P = mg = 2 \times 10 = 20\ N$.

## 5. Ứng dụng

Việc phân tích trọng lực và lực căng dây là cơ sở để giải các bài toán về vật cân bằng, hệ ròng rọc, con lắc — những mô hình phổ biến trong kỹ thuật và đời sống.

> [!TIP]
> Khi vật đứng yên hoặc chuyển động thẳng đều, hợp lực tác dụng lên vật bằng $0$ — đây là chìa khoá để lập phương trình cân bằng lực.`,
			},
			{
				Title:           "Lực ma sát",
				DurationMinutes: 30,
				Body: `## 1. Lực ma sát trượt

Lực ma sát trượt xuất hiện khi một vật trượt trên bề mặt một vật khác, có phương tiếp tuyến với bề mặt tiếp xúc, chiều ngược với chiều chuyển động (hoặc xu hướng chuyển động) của vật:

$$F_{mst} = \mu_t N$$

trong đó $\mu_t$ là hệ số ma sát trượt (phụ thuộc bản chất, tình trạng bề mặt tiếp xúc), $N$ là áp lực (độ lớn phản lực pháp tuyến) vuông góc với bề mặt tiếp xúc.

## 2. Lực ma sát nghỉ

Lực ma sát nghỉ xuất hiện khi vật chịu tác dụng của ngoại lực nhưng vẫn đứng yên trên bề mặt, có tác dụng giữ cho vật không trượt. Độ lớn lực ma sát nghỉ bằng độ lớn ngoại lực (theo phương tiếp tuyến) cho tới khi đạt giá trị cực đại, vượt qua giá trị đó vật bắt đầu trượt.

## 3. Đặc điểm chung

Lực ma sát trượt hầu như không phụ thuộc diện tích tiếp xúc và tốc độ trượt, chỉ phụ thuộc vào bản chất, tình trạng bề mặt (qua hệ số $\mu_t$) và áp lực $N$.

## 4. Ví dụ áp dụng

Một vật khối lượng $m = 5\ kg$ trượt trên mặt sàn nằm ngang với hệ số ma sát trượt $\mu_t = 0{,}2$ (lấy $g = 10\ m/s^2$). Tính độ lớn lực ma sát trượt.

Áp lực $N = P = mg = 50\ N$ (vì sàn ngang). Lực ma sát trượt: $F_{mst} = \mu_t N = 0{,}2 \times 50 = 10\ N$.

## 5. Vai trò của ma sát

Ma sát vừa có lợi (giúp người đi lại không trơn trượt, xe bám đường, phanh xe hoạt động...) vừa có hại (gây hao mòn máy móc, sinh nhiệt làm giảm hiệu suất) — cần được tính toán, kiểm soát hợp lý tùy tình huống.

> [!NOTE]
> Hệ số ma sát nghỉ cực đại thường lớn hơn một chút so với hệ số ma sát trượt của cùng một cặp bề mặt.`,
			},
			{
				Title:           "Momen lực. Điều kiện cân bằng của vật rắn",
				DurationMinutes: 30,
				Body: `## 1. Momen lực

Momen lực của một lực $\vec{F}$ đối với một trục quay là đại lượng đặc trưng cho tác dụng làm quay của lực quanh trục đó, xác định bởi:

$$M = F \cdot d$$

trong đó $F$ là độ lớn lực (N), $d$ là **cánh tay đòn** — khoảng cách từ trục quay đến giá của lực (m). Đơn vị của momen lực là N·m.

## 2. Quy tắc momen lực

Muốn một vật rắn có trục quay cố định ở trạng thái cân bằng, tổng momen lực làm vật quay theo chiều kim đồng hồ phải bằng tổng momen lực làm vật quay ngược chiều kim đồng hồ:

$$M_1 + M_2 + \cdots = M_1' + M_2' + \cdots$$

## 3. Ví dụ áp dụng

Một thanh cứng nhẹ có trục quay tại $O$. Đầu $A$ cách $O$ một đoạn $d_1 = 0{,}5\ m$ treo vật nặng tạo lực $F_1 = 20\ N$. Cần treo tại đầu $B$ cách $O$ một đoạn $d_2 = 0{,}2\ m$ (phía đối diện) một lực $F_2$ bao nhiêu để thanh cân bằng?

Theo quy tắc momen: $F_1 d_1 = F_2 d_2 \Rightarrow F_2 = \dfrac{F_1 d_1}{d_2} = \dfrac{20 \times 0{,}5}{0{,}2} = 50\ N$.

## 4. Ứng dụng

Quy tắc momen lực là nguyên lý hoạt động của đòn bẩy, cân đòn, cờ-lê vặn ốc... — vật càng xa trục quay thì cùng một lực tạo ra momen (tác dụng làm quay) càng lớn.

> [!TIP]
> Muốn nới lỏng một con ốc chặt, nên cầm cờ-lê ở vị trí xa trục ốc nhất có thể — cánh tay đòn dài hơn giúp tạo momen lớn hơn với cùng lực tay.`,
			},
			{
				Title:           "Công và công suất",
				DurationMinutes: 30,
				Body: `## 1. Công của một lực

Khi lực $\vec{F}$ không đổi tác dụng lên một vật và làm vật dịch chuyển một quãng đường $s$ theo hướng hợp với lực một góc $\alpha$, công của lực đó là:

$$A = F \cdot s \cdot \cos\alpha$$

Đơn vị của công là **jun** (J).

## 2. Các trường hợp đặc biệt

- $\alpha = 0^\circ$ (lực cùng hướng chuyển động): $A = Fs > 0$ — công phát động.
- $\alpha = 90^\circ$ (lực vuông góc với chuyển động): $A = 0$ — lực không sinh công.
- $\alpha = 180^\circ$ (lực ngược hướng chuyển động): $A = -Fs < 0$ — công cản.

## 3. Công suất

Công suất đặc trưng cho tốc độ thực hiện công, xác định bởi:

$$P = \dfrac{A}{t}$$

Đơn vị của công suất là **oát** (W); $1\ W = 1\ J/s$.

## 4. Ví dụ áp dụng

Một lực $F = 100\ N$ kéo một vật trượt trên mặt phẳng ngang, đi được quãng đường $s = 5\ m$ theo đúng hướng của lực, trong thời gian $t = 10\ s$. Tính công và công suất.

$$A = F \cdot s = 100 \times 5 = 500\ J \qquad P = \dfrac{A}{t} = \dfrac{500}{10} = 50\ W$$

## 5. Ứng dụng

Công suất là thông số quan trọng để so sánh khả năng sinh công theo thời gian của các loại máy móc, động cơ (VD: công suất động cơ xe máy, ô tô tính bằng mã lực hoặc kW).

> [!NOTE]
> Một lực có thể rất lớn nhưng nếu không làm vật dịch chuyển (hoặc vuông góc với chuyển động) thì công của lực đó vẫn bằng $0$.`,
			},
			{
				Title:           "Động lượng. Định luật bảo toàn động lượng",
				DurationMinutes: 30,
				Body: `## 1. Động lượng

Động lượng của một vật là đại lượng vectơ, bằng tích khối lượng và vận tốc của vật:

$$\vec{p} = m\vec{v}$$

Đơn vị động lượng là kg·m/s. Động lượng cùng hướng với vận tốc của vật.

## 2. Hệ kín (hệ cô lập)

Hệ kín là hệ chỉ có các vật trong hệ tương tác với nhau (nội lực), không chịu tác dụng của ngoại lực, hoặc các ngoại lực cân bằng nhau.

## 3. Định luật bảo toàn động lượng

Trong một hệ kín, tổng động lượng của hệ được bảo toàn (không đổi theo thời gian). Với hệ hai vật va chạm:

$$m_1\vec{v}_1 + m_2\vec{v}_2 = m_1\vec{v}_1' + m_2\vec{v}_2'$$

(dấu phẩy chỉ vận tốc sau va chạm).

## 4. Ví dụ áp dụng

Một xe khối lượng $m_1 = 2\ kg$ chuyển động với vận tốc $v_1 = 3\ m/s$ đến va chạm và dính vào một xe khối lượng $m_2 = 1\ kg$ đang đứng yên (va chạm mềm). Tính vận tốc hai xe sau va chạm.

Bảo toàn động lượng: $m_1v_1 = (m_1 + m_2)v \Rightarrow v = \dfrac{2 \times 3}{2 + 1} = 2\ m/s$.

## 5. Ứng dụng

Định luật bảo toàn động lượng giải thích nguyên lý chuyển động của tên lửa (đẩy khí ra sau để tiến lên phía trước), súng giật khi bắn, và là công cụ chính để giải các bài toán va chạm.

> [!TIP]
> Định luật bảo toàn động lượng luôn đúng với hệ kín, kể cả trong va chạm mềm (không bảo toàn động năng) lẫn va chạm đàn hồi (bảo toàn cả động năng).`,
			},
		},
		Quiz: sampleQuiz{
			Title:        "Bài tập: Ba định luật Newton và chuyển động thẳng biến đổi đều",
			Instructions: "Làm bài để tự kiểm tra kiến thức hai bài giảng vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Theo định luật I Newton, nếu hợp lực tác dụng lên một vật bằng 0 thì vật sẽ như thế nào?",
					Level:  "Nhận biết", Points: 2.5,
					Explanation: "Vật đang đứng yên tiếp tục đứng yên, vật đang chuyển động tiếp tục chuyển động thẳng đều — đây là tính chất quán tính.",
					Options: []sampleOption{
						{Content: "Đứng yên hoặc chuyển động thẳng đều", IsCorrect: true},
						{Content: "Luôn đứng yên"},
						{Content: "Chuyển động nhanh dần đều"},
						{Content: "Dừng lại ngay lập tức"},
					},
				},
				{
					Prompt: "Một vật khối lượng 4 kg chịu tác dụng của lực 8 N. Gia tốc của vật là?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "a = F/m = 8/4 = 2 m/s².",
					Options: []sampleOption{
						{Content: "2 m/s²", IsCorrect: true},
						{Content: "0,5 m/s²"},
						{Content: "32 m/s²"},
						{Content: "4 m/s²"},
					},
				},
				{
					Prompt: "Cặp lực và phản lực theo định luật III Newton có đặc điểm gì?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "Lực và phản lực cùng giá, cùng độ lớn, ngược chiều, đặt vào hai vật khác nhau nên không cân bằng nhau.",
					Options: []sampleOption{
						{Content: "Cùng độ lớn, ngược chiều, đặt vào hai vật khác nhau", IsCorrect: true},
						{Content: "Cùng độ lớn, cùng chiều, đặt vào cùng một vật"},
						{Content: "Khác độ lớn, ngược chiều"},
						{Content: "Luôn triệt tiêu lẫn nhau"},
					},
				},
				{
					Prompt: "Một xe chuyển động nhanh dần đều từ v₀ = 0 với gia tốc 2 m/s². Sau 4 giây, quãng đường xe đi được là?",
					Level:  "Vận dụng", Points: 2.5,
					Explanation: "s = v₀t + ½at² = 0 + ½×2×4² = 16 m.",
					Options: []sampleOption{
						{Content: "16 m", IsCorrect: true},
						{Content: "8 m"},
						{Content: "32 m"},
						{Content: "4 m"},
					},
				},
			},
		},
	},
	"HOAHOC": {
		Lessons: []sampleLesson{
			{
				Title:           "Bảng tuần hoàn các nguyên tố hóa học",
				DurationMinutes: 30,
				Body: `## 1. Nguyên tắc sắp xếp

Các nguyên tố hoá học trong bảng tuần hoàn được sắp xếp theo chiều tăng dần của **điện tích hạt nhân** (số proton, còn gọi là số hiệu nguyên tử $Z$).

## 2. Cấu tạo bảng tuần hoàn

- **Ô nguyên tố**: mỗi nguyên tố chiếm một ô, cho biết số hiệu nguyên tử, ký hiệu hoá học, tên nguyên tố và nguyên tử khối trung bình.
- **Chu kỳ**: là dãy các nguyên tố mà nguyên tử có cùng số lớp electron, xếp theo hàng ngang. Bảng tuần hoàn hiện có 7 chu kỳ.
- **Nhóm**: gồm các nguyên tố mà nguyên tử có cấu hình electron lớp ngoài cùng tương tự nhau, tính chất hoá học gần giống nhau, xếp theo cột dọc. Có 18 nhóm (nhóm A và nhóm B).

## 3. Sự biến đổi tính chất theo chu kỳ và nhóm

- **Trong một chu kỳ** (trái → phải): tính kim loại giảm dần, tính phi kim tăng dần.
- **Trong một nhóm A** (trên → dưới): tính kim loại tăng dần, tính phi kim giảm dần.

## 4. Một số nhóm nguyên tố quan trọng

- **Nhóm IA** (trừ H): kim loại kiềm — Li, Na, K... rất hoạt động hoá học, dễ nhường 1 electron.
- **Nhóm VIIA**: halogen — F, Cl, Br, I... phi kim mạnh, dễ nhận 1 electron.
- **Nhóm VIIIA**: khí hiếm — He, Ne, Ar... hầu như trơ về mặt hoá học.

## 5. Ví dụ áp dụng

Nguyên tố Natri (Na) có số hiệu nguyên tử $Z = 11$, thuộc chu kỳ 3, nhóm IA — là kim loại kiềm, có 1 electron lớp ngoài cùng nên dễ nhường electron để tạo ion $Na^+$.

> [!NOTE]
> Bảng tuần hoàn hiện đại do nhà hoá học người Nga Dmitri Mendeleev đề xuất năm 1869, dựa trên chiều tăng dần của nguyên tử khối; sau này được hoàn thiện lại theo điện tích hạt nhân.`,
			},
			{
				Title:           "Phản ứng oxi hoá - khử",
				DurationMinutes: 30,
				Body: `## 1. Một số khái niệm

- **Số oxi hoá**: điện tích quy ước của nguyên tử trong phân tử nếu giả định các liên kết đều là liên kết ion.
- **Sự oxi hoá**: quá trình một chất **nhường electron** (số oxi hoá tăng).
- **Sự khử**: quá trình một chất **nhận electron** (số oxi hoá giảm).
- **Chất khử**: chất nhường electron (bị oxi hoá). **Chất oxi hoá**: chất nhận electron (bị khử).

## 2. Phản ứng oxi hoá - khử

Là phản ứng hoá học trong đó có sự chuyển electron giữa các chất phản ứng, tức là có sự thay đổi số oxi hoá của một số nguyên tố.

## 3. Ví dụ minh hoạ

Xét phản ứng:

$$Zn + CuSO_4 \rightarrow ZnSO_4 + Cu$$

- Kẽm ($Zn$) nhường 2 electron: $Zn \rightarrow Zn^{2+} + 2e$ → Zn là **chất khử**, bị oxi hoá.
- Đồng ($Cu^{2+}$) nhận 2 electron: $Cu^{2+} + 2e \rightarrow Cu$ → $Cu^{2+}$ là **chất oxi hoá**, bị khử.

## 4. Các bước cân bằng phản ứng oxi hoá - khử (phương pháp thăng bằng electron)

1. Xác định số oxi hoá của các nguyên tố để tìm chất khử, chất oxi hoá.
2. Viết quá trình oxi hoá và quá trình khử, cân bằng số electron nhường - nhận.
3. Đặt hệ số thích hợp vào chất khử, chất oxi hoá sao cho tổng electron nhường bằng tổng electron nhận.
4. Kiểm tra lại và cân bằng các nguyên tố còn lại, hoàn thành phương trình.

## 5. Ứng dụng

Phản ứng oxi hoá - khử có mặt trong rất nhiều quá trình: sự cháy, sự gỉ sét của kim loại, pin và ắc quy, luyện kim (điều chế kim loại từ quặng), quang hợp và hô hấp ở sinh vật...

> [!TIP]
> Mẹo nhớ: "**Khử cho** – **O nhận**" (chất khử cho electron, chất oxi hoá nhận electron).`,
			},
			{
				Title:           "Thành phần nguyên tử",
				DurationMinutes: 30,
				Body: `## 1. Nguyên tử

Nguyên tử là hạt vô cùng nhỏ, trung hoà về điện, cấu tạo nên các chất. Nguyên tử gồm hạt nhân mang điện tích dương ở tâm và lớp vỏ electron mang điện tích âm chuyển động xung quanh.

## 2. Các hạt cấu tạo nên nguyên tử

- **Proton** ($p$): điện tích $+1$, khối lượng xấp xỉ $1$ u, nằm trong hạt nhân.
- **Neutron** ($n$): không mang điện, khối lượng xấp xỉ $1$ u, nằm trong hạt nhân.
- **Electron** ($e$): điện tích $-1$, khối lượng rất nhỏ (xấp xỉ $0{,}00055$ u), chuyển động trong lớp vỏ.

## 3. Tính trung hoà về điện

Vì nguyên tử trung hoà về điện nên số proton (trong hạt nhân) luôn bằng số electron (ở lớp vỏ): $Z = P = E$.

## 4. Kích thước và khối lượng nguyên tử

Nguyên tử có kích thước vô cùng nhỏ (đường kính cỡ $10^{-10}\ m$), trong khi hạt nhân còn nhỏ hơn rất nhiều (cỡ $10^{-14}\ m$) — nghĩa là nguyên tử có cấu tạo rỗng, hầu hết khối lượng tập trung ở hạt nhân vì electron có khối lượng rất nhỏ so với proton, neutron.

> [!NOTE]
> Mô hình nguyên tử hiện đại được xây dựng qua nhiều thí nghiệm nổi tiếng, trong đó thí nghiệm bắn phá lá vàng của Rutherford (1911) đã chứng minh hạt nhân mang điện dương, kích thước rất nhỏ so với toàn bộ nguyên tử.`,
			},
			{
				Title:           "Cấu tạo hạt nhân. Nguyên tố hóa học, đồng vị",
				DurationMinutes: 30,
				Body: `## 1. Số hiệu nguyên tử, số khối

**Số hiệu nguyên tử** ($Z$) là số proton trong hạt nhân, cũng là số electron của nguyên tử trung hoà — đặc trưng riêng cho mỗi nguyên tố hoá học.

**Số khối** ($A$) là tổng số proton và neutron trong hạt nhân: $A = Z + N$ (với $N$ là số neutron). Một nguyên tử được ký hiệu đầy đủ là $^A_ZX$.

## 2. Nguyên tố hoá học

Nguyên tố hoá học là tập hợp các nguyên tử có cùng số proton trong hạt nhân (cùng số hiệu nguyên tử $Z$). Các nguyên tử của cùng một nguyên tố có tính chất hoá học giống nhau.

## 3. Đồng vị

Đồng vị là những nguyên tử có cùng số proton ($Z$) nhưng khác số neutron ($N$), do đó có số khối $A$ khác nhau. Các đồng vị của cùng một nguyên tố có tính chất hoá học gần như giống nhau (vì cùng số electron) nhưng tính chất vật lý có thể khác (VD: khối lượng).

## 4. Ví dụ áp dụng

Nguyên tố Carbon có ba đồng vị phổ biến: $^{12}_6C$, $^{13}_6C$, $^{14}_6C$ — cả ba đều có $Z = 6$ (6 proton) nhưng số neutron lần lượt là $6, 7, 8$.

## 5. Nguyên tử khối trung bình

Vì các nguyên tố tồn tại trong tự nhiên dưới dạng hỗn hợp nhiều đồng vị theo tỉ lệ phần trăm nhất định, nguyên tử khối ghi trong bảng tuần hoàn là **nguyên tử khối trung bình**, tính theo tỉ lệ phần trăm số nguyên tử của mỗi đồng vị.

> [!TIP]
> Đồng vị phóng xạ $^{14}_6C$ được dùng trong phương pháp xác định tuổi các di vật khảo cổ (phương pháp xác định niên đại bằng carbon phóng xạ).`,
			},
			{
				Title:           "Cấu trúc lớp vỏ electron nguyên tử",
				DurationMinutes: 30,
				Body: `## 1. Sự chuyển động của electron trong nguyên tử

Electron chuyển động rất nhanh quanh hạt nhân, không theo quỹ đạo xác định như hành tinh quanh Mặt Trời, mà tạo thành một vùng không gian gọi là **orbital nguyên tử** — nơi có xác suất tìm thấy electron lớn nhất.

## 2. Lớp và phân lớp electron

Các electron được sắp xếp thành các **lớp** electron (ký hiệu $n = 1, 2, 3...$, hay $K, L, M, N...$), mỗi lớp lại chia thành các **phân lớp** ($s, p, d, f$). Số electron tối đa của lớp thứ $n$ là $2n^2$.

## 3. Nguyên lý sắp xếp electron

Electron được sắp xếp vào các lớp, phân lớp theo thứ tự mức năng lượng tăng dần, tuân theo thứ tự gần đúng: $1s\ 2s\ 2p\ 3s\ 3p\ 4s\ 3d\ 4p...$

## 4. Cấu hình electron

Cấu hình electron biểu diễn sự phân bố electron trên các phân lớp. Ví dụ, nguyên tử Natri ($Z = 11$) có cấu hình electron: $1s^2 2s^2 2p^6 3s^1$.

## 5. Electron lớp ngoài cùng

Số electron ở lớp ngoài cùng quyết định phần lớn tính chất hoá học của nguyên tố: nguyên tử có $1, 2, 3$ electron lớp ngoài cùng thường là kim loại; có $5, 6, 7$ electron lớp ngoài cùng thường là phi kim; có đúng $8$ electron (hoặc $2$ với heli) lớp ngoài cùng là khí hiếm, bền vững.

> [!NOTE]
> Với Natri ($1s^2 2s^2 2p^6 3s^1$), electron lớp ngoài cùng $3s^1$ dễ bị nhường đi để đạt cấu hình bền của khí hiếm gần nhất — giải thích vì sao Na dễ tạo ion $Na^+$.`,
			},
			{
				Title:           "Liên kết ion",
				DurationMinutes: 30,
				Body: `## 1. Ion là gì?

Khi nguyên tử nhường hoặc nhận electron, nó trở thành phần tử mang điện gọi là **ion**. Nguyên tử nhường electron tạo thành **cation** (ion dương); nguyên tử nhận electron tạo thành **anion** (ion âm).

Ví dụ: $Na \rightarrow Na^+ + 1e$ (cation natri); $Cl + 1e \rightarrow Cl^-$ (anion clorua).

## 2. Liên kết ion

Liên kết ion là liên kết được hình thành bởi lực hút tĩnh điện giữa các ion mang điện tích trái dấu, thường xảy ra giữa nguyên tử kim loại điển hình (dễ nhường electron) và nguyên tử phi kim điển hình (dễ nhận electron).

## 3. Quá trình hình thành liên kết ion — ví dụ NaCl

$$Na \rightarrow Na^+ + 1e \qquad\qquad Cl + 1e \rightarrow Cl^-$$

$$Na^+ + Cl^- \rightarrow NaCl$$

Ion $Na^+$ và $Cl^-$ hút nhau bằng lực hút tĩnh điện, tạo thành tinh thể ion NaCl bền vững.

## 4. Tính chất chung của hợp chất ion

- Là chất rắn ở điều kiện thường, có nhiệt độ nóng chảy và nhiệt độ sôi cao (do lực hút tĩnh điện giữa các ion khá mạnh).
- Thường tan tốt trong nước, khi tan hoặc nóng chảy thì dẫn được điện (do các ion trở nên linh động).

## 5. Điều kiện hình thành liên kết ion

Liên kết ion thường hình thành giữa các nguyên tố có **hiệu độ âm điện lớn** (thường $\ge 1{,}7$), phổ biến giữa kim loại điển hình nhóm IA, IIA với phi kim điển hình nhóm VIIA, VIA.

> [!TIP]
> Nhớ nhanh: kim loại "cho" electron thành ion dương, phi kim "nhận" electron thành ion âm — hai loại ion trái dấu hút nhau tạo liên kết ion.`,
			},
			{
				Title:           "Liên kết cộng hóa trị",
				DurationMinutes: 30,
				Body: `## 1. Định nghĩa

Liên kết cộng hoá trị là liên kết được hình thành giữa hai nguyên tử bằng một hay nhiều cặp electron dùng chung, thường xảy ra giữa các nguyên tử phi kim (hoặc giữa các nguyên tử của cùng một nguyên tố).

## 2. Liên kết cộng hoá trị không cực và có cực

- **Không cực**: cặp electron dùng chung không lệch về phía nguyên tử nào — thường giữa hai nguyên tử của cùng một nguyên tố (VD: $H_2$, $Cl_2$, $N_2$).
- **Có cực**: cặp electron dùng chung lệch về phía nguyên tử có độ âm điện lớn hơn — giữa hai nguyên tử phi kim khác nguyên tố (VD: $HCl$, $H_2O$).

## 3. Công thức Lewis và công thức cấu tạo

Công thức Lewis biểu diễn các electron hoá trị bằng dấu chấm; mỗi cặp electron dùng chung có thể thay bằng một gạch nối trong công thức cấu tạo. Ví dụ, phân tử $H_2$: $H{-}H$ (một cặp electron dùng chung, liên kết đơn); phân tử $N_2$: $N{\equiv}N$ (ba cặp electron dùng chung, liên kết ba).

## 4. Ví dụ áp dụng — phân tử nước $H_2O$

Nguyên tử oxygen góp chung mỗi nguyên tử hydrogen một cặp electron, tạo thành hai liên kết cộng hoá trị $O{-}H$. Vì độ âm điện của oxygen lớn hơn hydrogen nên đây là liên kết cộng hoá trị có cực.

## 5. So sánh với liên kết ion

Liên kết cộng hoá trị hình thành do dùng chung electron (không có sự chuyển hẳn electron như liên kết ion), thường xảy ra khi hiệu độ âm điện giữa hai nguyên tử nhỏ (thường $< 1{,}7$).

> [!NOTE]
> Một nguyên tử có thể tạo nhiều liên kết cộng hoá trị cùng lúc, miễn là tổng số electron dùng chung giúp mỗi nguyên tử đạt cấu hình electron bền (thường là đủ 8 electron lớp ngoài cùng — quy tắc octet).`,
			},
			{
				Title:           "Biến thiên enthalpy trong các phản ứng hóa học",
				DurationMinutes: 30,
				Body: `## 1. Phản ứng toả nhiệt và thu nhiệt

**Phản ứng toả nhiệt** là phản ứng giải phóng năng lượng dưới dạng nhiệt ra môi trường xung quanh (VD: phản ứng đốt cháy nhiên liệu). **Phản ứng thu nhiệt** là phản ứng hấp thụ năng lượng từ môi trường (VD: phản ứng nung vôi).

## 2. Biến thiên enthalpy

Biến thiên enthalpy của phản ứng ($\Delta_rH$) là nhiệt lượng toả ra hay thu vào của phản ứng ở điều kiện áp suất không đổi.

- Nếu $\Delta_rH < 0$: phản ứng toả nhiệt.
- Nếu $\Delta_rH > 0$: phản ứng thu nhiệt.

Đơn vị thường dùng là kJ (kilojun).

## 3. Enthalpy tạo thành chuẩn

Enthalpy tạo thành chuẩn ($\Delta_fH^{\circ}_{298}$) của một chất là biến thiên enthalpy của phản ứng tạo thành $1$ mol chất đó từ các đơn chất bền ở điều kiện chuẩn ($25^\circ C$, $1\ bar$).

## 4. Ví dụ áp dụng

Phản ứng đốt cháy khí methane: $CH_4(g) + 2O_2(g) \rightarrow CO_2(g) + 2H_2O(l)$ có $\Delta_rH^{\circ}_{298} = -890\ kJ$ — đây là phản ứng toả nhiệt mạnh, giải phóng $890\ kJ$ nhiệt lượng khi đốt cháy hoàn toàn $1$ mol khí methane.

## 5. Ứng dụng

Tính toán biến thiên enthalpy giúp đánh giá hiệu quả năng lượng của nhiên liệu, thiết kế các quá trình công nghiệp (luyện kim, sản xuất hoá chất), và hiểu các quá trình sinh học (hô hấp tế bào toả năng lượng).

> [!TIP]
> Phản ứng toả nhiệt ($\Delta_rH < 0$) thường xảy ra dễ dàng và tự phát hơn phản ứng thu nhiệt, nhưng không phải yếu tố duy nhất quyết định phản ứng có tự xảy ra hay không.`,
			},
			{
				Title:           "Tốc độ phản ứng hóa học",
				DurationMinutes: 30,
				Body: `## 1. Khái niệm tốc độ phản ứng

Tốc độ phản ứng là đại lượng đặc trưng cho sự nhanh, chậm của một phản ứng hoá học, được xác định bằng độ biến thiên nồng độ của một chất phản ứng hoặc sản phẩm trong một đơn vị thời gian.

## 2. Các yếu tố ảnh hưởng đến tốc độ phản ứng

- **Nồng độ**: nồng độ chất phản ứng càng lớn, tốc độ phản ứng càng nhanh.
- **Nhiệt độ**: nhiệt độ càng cao, tốc độ phản ứng càng nhanh (các phân tử chuyển động nhanh hơn, va chạm hiệu quả hơn).
- **Áp suất** (đối với phản ứng có chất khí): áp suất tăng làm tăng nồng độ chất khí, thường làm tăng tốc độ phản ứng.
- **Diện tích bề mặt tiếp xúc**: chất rắn ở dạng bột mịn phản ứng nhanh hơn dạng khối lớn (diện tích tiếp xúc lớn hơn).
- **Chất xúc tác**: làm tăng tốc độ phản ứng nhưng không bị biến đổi sau phản ứng.

## 3. Ví dụ áp dụng

Khi cho cùng một lượng đá vôi (CaCO₃) phản ứng với dung dịch HCl: đá vôi dạng bột mịn phản ứng nhanh và sinh khí CO₂ mạnh hơn hẳn so với đá vôi ở dạng viên lớn, do diện tích bề mặt tiếp xúc lớn hơn.

## 4. Ý nghĩa thực tiễn

Hiểu các yếu tố ảnh hưởng đến tốc độ phản ứng giúp điều chỉnh quá trình sản xuất công nghiệp (tăng tốc độ để nâng cao hiệu suất), bảo quản thực phẩm (làm chậm phản ứng ôi thiu bằng cách hạ nhiệt độ - để tủ lạnh), và nhiều ứng dụng khác trong đời sống.

> [!NOTE]
> Chất xúc tác chỉ làm thay đổi tốc độ phản ứng, không làm thay đổi trạng thái cân bằng cuối cùng hay lượng sản phẩm tạo thành của phản ứng.`,
			},
			{
				Title:           "Nguyên tố nhóm VIIA (Halogen)",
				DurationMinutes: 30,
				Body: `## 1. Vị trí trong bảng tuần hoàn

Nhóm VIIA gồm các nguyên tố: Fluorine (F), Chlorine (Cl), Bromine (Br), Iodine (I)... gọi chung là **halogen**. Nguyên tử các nguyên tố này đều có $7$ electron lớp ngoài cùng.

## 2. Tính chất hoá học đặc trưng

Do có $7$ electron lớp ngoài cùng, nguyên tử halogen dễ nhận thêm $1$ electron để đạt cấu hình bền của khí hiếm, thể hiện tính **phi kim mạnh**, tính oxi hoá mạnh: $X + 1e \rightarrow X^-$.

## 3. Sự biến đổi tính chất trong nhóm

Đi từ Fluorine đến Iodine (từ trên xuống dưới trong nhóm), bán kính nguyên tử tăng dần, độ âm điện giảm dần nên tính phi kim (tính oxi hoá) giảm dần: $F_2 > Cl_2 > Br_2 > I_2$.

## 4. Đơn chất và hợp chất tiêu biểu

- $Cl_2$ (khí chlorine): màu vàng lục, dùng khử trùng nước sinh hoạt, sản xuất chất tẩy trắng.
- $HCl$ (acid clohydric): là acid mạnh, có nhiều ứng dụng trong công nghiệp và có sẵn (nồng độ loãng) trong dịch vị dạ dày người.
- $NaCl$ (muối ăn): hợp chất ion phổ biến nhất của Chlorine, thiết yếu trong đời sống.

## 5. Ứng dụng

Các hợp chất của halogen được dùng rộng rãi: fluoride trong kem đánh răng (ngừa sâu răng), chlorine khử trùng nước, iodine bổ sung trong muối ăn (phòng bệnh bướu cổ)...

> [!TIP]
> Halogen luôn tồn tại trong tự nhiên ở dạng hợp chất (ion âm hoặc phân tử $X_2$) chứ không tồn tại ở dạng đơn nguyên tử tự do, vì tính hoạt động hoá học rất mạnh.`,
			},
		},
		Quiz: sampleQuiz{
			Title:        "Bài tập: Bảng tuần hoàn và phản ứng oxi hoá - khử",
			Instructions: "Làm bài để tự kiểm tra kiến thức hai bài giảng vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Các nguyên tố trong bảng tuần hoàn được sắp xếp theo chiều tăng dần của đại lượng nào?",
					Level:  "Nhận biết", Points: 2.5,
					Explanation: "Nguyên tắc sắp xếp là theo chiều tăng dần của điện tích hạt nhân (số hiệu nguyên tử).",
					Options: []sampleOption{
						{Content: "Điện tích hạt nhân (số hiệu nguyên tử)", IsCorrect: true},
						{Content: "Khối lượng riêng"},
						{Content: "Nhiệt độ nóng chảy"},
						{Content: "Bán kính nguyên tử"},
					},
				},
				{
					Prompt: "Trong một chu kỳ, đi từ trái sang phải, tính chất của các nguyên tố biến đổi như thế nào?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "Trong một chu kỳ, tính kim loại giảm dần còn tính phi kim tăng dần.",
					Options: []sampleOption{
						{Content: "Tính kim loại giảm dần, tính phi kim tăng dần", IsCorrect: true},
						{Content: "Tính kim loại tăng dần, tính phi kim giảm dần"},
						{Content: "Cả tính kim loại và phi kim đều tăng"},
						{Content: "Không có sự biến đổi"},
					},
				},
				{
					Prompt: "Trong phản ứng Zn + CuSO₄ → ZnSO₄ + Cu, chất nào đóng vai trò chất khử?",
					Level:  "Thông hiểu", Points: 2.5,
					Explanation: "Zn nhường electron (Zn → Zn²⁺ + 2e) nên Zn là chất khử.",
					Options: []sampleOption{
						{Content: "Zn", IsCorrect: true},
						{Content: "Cu²⁺"},
						{Content: "SO₄²⁻"},
						{Content: "Cu"},
					},
				},
				{
					Prompt: "Quá trình Cu²⁺ + 2e → Cu được gọi là gì?",
					Level:  "Vận dụng", Points: 2.5,
					Explanation: "Cu²⁺ nhận electron để thành Cu, đây là sự khử.",
					Options: []sampleOption{
						{Content: "Sự khử", IsCorrect: true},
						{Content: "Sự oxi hoá"},
						{Content: "Sự trung hoà"},
						{Content: "Sự thuỷ phân"},
					},
				},
			},
		},
	},
}

// ensureSampleContent thêm bài giảng và bài tập mẫu cho một lớp học mẫu nếu
// lớp đó hiện chưa có nội dung nào. Chỉ xét khi mã có trong sampleContentByCode;
// nếu chương trình đã có ít nhất một nút (do seed trước, hoặc do giáo viên tự
// soạn/xoá) thì bỏ qua hẳn — seed không bao giờ ghi đè nội dung đã tồn tại.
func (s *Store) ensureSampleContent(ctx context.Context, programID uuid.UUID, code string) error {
	content, ok := sampleContentByCode[code]
	if !ok {
		return nil
	}

	var nodeCount int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM nodes WHERE program_id = $1`, programID).
		Scan(&nodeCount); err != nil {
		return fmt.Errorf("đếm nội dung hiện có: %w", err)
	}
	if nodeCount > 0 {
		return nil
	}

	for _, lesson := range content.Lessons {
		if _, err := s.CreateNode(ctx, SaveNodeParams{
			ProgramID:   programID,
			Kind:        models.KindLesson,
			Title:       lesson.Title,
			IsPublished: true,
			Lesson: &models.Lesson{
				ContentType:     "richtext",
				Body:            lesson.Body,
				DurationMinutes: lesson.DurationMinutes,
				Attachments:     []models.LessonAttachment{},
			},
		}); err != nil {
			return fmt.Errorf("tạo bài giảng mẫu %q: %w", lesson.Title, err)
		}
	}

	quiz := content.Quiz
	assignmentNode, err := s.CreateNode(ctx, SaveNodeParams{
		ProgramID:   programID,
		Kind:        models.KindAssignment,
		Title:       quiz.Title,
		IsPublished: true,
		Assignment: &models.Assignment{
			Instructions: quiz.Instructions,
			MaxAttempts:  3,
			PassScore:    5,
		},
	})
	if err != nil {
		return fmt.Errorf("tạo bài tập mẫu %q: %w", quiz.Title, err)
	}

	for _, q := range quiz.Questions {
		opts := make([]QuestionOptionInput, len(q.Options))
		for i, o := range q.Options {
			opts[i] = QuestionOptionInput{Content: o.Content, IsCorrect: o.IsCorrect}
		}
		if _, err := s.CreateQuestion(ctx, SaveQuestionParams{
			AssignmentID: assignmentNode.ID,
			Type:         models.QuestionSingleChoice,
			Prompt:       q.Prompt,
			Points:       q.Points,
			Level:        q.Level,
			Explanation:  q.Explanation,
			Options:      opts,
		}); err != nil {
			return fmt.Errorf("tạo câu hỏi mẫu %q: %w", q.Prompt, err)
		}
	}

	return nil
}
