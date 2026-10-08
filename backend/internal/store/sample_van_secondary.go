package store

// sampleVanSecondary là các lớp học mẫu môn Ngữ văn lớp 7–12 theo Chương trình
// GDPT 2018: đọc hiểu theo thể loại, tiếng Việt và kỹ năng viết.
var sampleVanSecondary = []sampleCourse{
	{
		Code:        "NGUVAN7",
		Title:       "Ngữ văn lớp 7",
		Description: "Lớp học mẫu môn Ngữ văn lớp 7 theo Chương trình GDPT 2018: đọc hiểu tục ngữ và truyện ngụ ngôn, biện pháp tu từ nói quá – nói giảm nói tránh và viết bài văn biểu cảm.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Đọc hiểu: tục ngữ và truyện ngụ ngôn",
				Lessons: []sampleLesson{
					{
						Title:           "Tục ngữ về thiên nhiên, lao động và con người",
						DurationMinutes: 30,
						Body: `## 1. Tục ngữ là gì?

Tục ngữ là những câu nói dân gian **ngắn gọn, ổn định, có nhịp điệu và hình ảnh**, đúc kết kinh nghiệm của nhân dân về thiên nhiên, lao động sản xuất, con người và xã hội. Mỗi câu tục ngữ là một **câu trọn vẹn**, nêu một nhận định hay một lời khuyên.

## 2. Đặc điểm hình thức

- **Ngắn gọn**: có câu chỉ bốn chữ, như "Tấc đất tấc vàng".
- **Thường có vần** (hay gặp vần lưng): "Ráng mỡ **gà**, có **nhà** thì giữ".
- **Nhịp cân đối, hay dùng phép đối**: "Mau sao thì **nắng**, **vắng** sao thì mưa" (mau – vắng, nắng – mưa).
- **Giàu hình ảnh**, nhiều câu có cả nghĩa đen và nghĩa bóng.

## 3. Các nhóm tục ngữ tiêu biểu

| Nhóm | Ví dụ | Ý nghĩa |
|---|---|---|
| Thiên nhiên | Mau sao thì nắng, vắng sao thì mưa | Đêm nhiều sao thì hôm sau thường nắng, ít sao thì dễ mưa |
| Thiên nhiên | Ráng mỡ gà, có nhà thì giữ | Chân trời có ráng vàng như mỡ gà là dấu hiệu sắp có bão |
| Lao động | Nhất nước, nhì phân, tam cần, tứ giống | Thứ tự các yếu tố quan trọng khi trồng lúa |
| Con người | Một mặt người bằng mười mặt của | Con người quý hơn mọi của cải |
| Con người | Đói cho sạch, rách cho thơm | Dù nghèo khó vẫn phải giữ phẩm giá |

## 4. Ví dụ phân tích: "Tấc đất tấc vàng"

- **Nghĩa đen**: một mảnh đất rất nhỏ (một "tấc") cũng quý như vàng.
- **Nghĩa bóng**: đất đai nuôi sống con người nên vô cùng quý giá.
- **Lời khuyên**: phải biết quý trọng, sử dụng đất hợp lí, phê phán việc bỏ hoang, lãng phí đất.

## 5. Phân biệt tục ngữ và thành ngữ

- **Tục ngữ** là một câu, diễn đạt trọn vẹn một ý: "Có công mài sắt, có ngày nên kim".
- **Thành ngữ** là cụm từ cố định, chưa thành câu, thường dùng như một từ: "đầu voi đuôi chuột", "mẹ tròn con vuông".

> [!TIP]
> Khi tìm hiểu một câu tục ngữ, hãy tự hỏi ba điều: nghĩa đen là gì, nghĩa bóng là gì, và kinh nghiệm ấy còn dùng được trong hoàn cảnh nào hôm nay.`,
					},
					{
						Title:           "Truyện ngụ ngôn: Ếch ngồi đáy giếng và Thầy bói xem voi",
						DurationMinutes: 30,
						Body: `## 1. Truyện ngụ ngôn là gì?

Truyện ngụ ngôn là truyện kể ngắn, bằng văn xuôi hoặc văn vần, **mượn chuyện loài vật, đồ vật hoặc chính con người** để nói bóng gió, kín đáo chuyện đời, nhằm khuyên nhủ, răn dạy một bài học.

Đặc điểm thường gặp:

- Dung lượng ngắn, ít sự việc, nhân vật thường không có tên riêng.
- Tình huống gây bất ngờ hoặc gây cười.
- Bài học được gửi gắm qua câu chuyện, có khi được nói thẳng ở cuối truyện.

## 2. Ếch ngồi đáy giếng

**Tóm tắt**: Một con ếch sống lâu ngày trong giếng, xung quanh chỉ có vài con nhái, cua, ốc. Mỗi lần ếch kêu, cả giếng vang động khiến các con vật kia hoảng sợ, nên ếch tưởng bầu trời chỉ bé bằng chiếc vung và mình là chúa tể. Một năm mưa to, nước giếng dềnh lên đưa ếch ra ngoài. Quen thói cũ, ếch nghênh ngang đi lại, chẳng để ý xung quanh, nên bị một con trâu đi qua giẫm bẹp.

**Bài học**: Môi trường hạn hẹp dễ khiến hiểu biết nông cạn; phê phán thói kiêu ngạo, chủ quan; khuyên phải không ngừng mở rộng hiểu biết, khiêm tốn học hỏi.

## 3. Thầy bói xem voi

**Tóm tắt**: Năm ông thầy bói mù chưa biết hình thù con voi, góp tiền biếu người quản voi để được sờ thử. Mỗi ông chỉ sờ một bộ phận (vòi, ngà, tai, chân, đuôi) rồi khẳng định con voi giống thứ mình vừa sờ: người bảo như con đỉa, người bảo như cái đòn càn, người bảo như cái quạt thóc, người bảo như cái cột đình, người bảo như cái chổi sể. Không ai chịu ai, năm ông xô xát, đánh nhau toác đầu chảy máu.

**Bài học**: Muốn hiểu đúng sự vật phải xem xét **toàn diện**; phê phán cách nhìn phiến diện và thái độ bảo thủ, không chịu lắng nghe.

## 4. Nghệ thuật

- Mượn chuyện nhỏ, gần gũi để nói điều lớn.
- Cách so sánh dân dã, lặp lại cấu trúc lời nói của các nhân vật tạo tiếng cười.
- Kết thúc bất ngờ, làm nổi bật bài học.

> [!NOTE]
> Nhiều truyện ngụ ngôn đã trở thành thành ngữ quen thuộc như "ếch ngồi đáy giếng", "thầy bói xem voi", "đẽo cày giữa đường" — dùng để chê người hiểu biết hạn hẹp, nhìn phiến diện hoặc thiếu chủ kiến.`,
					},
				},
			},
			{
				Title: "Chương 2. Tiếng Việt và kỹ năng viết",
				Lessons: []sampleLesson{
					{
						Title:           "Nói quá và nói giảm nói tránh",
						DurationMinutes: 30,
						Body: `## 1. Nói quá

**Nói quá** là biện pháp tu từ **phóng đại** mức độ, quy mô, tính chất của sự vật, hiện tượng để nhấn mạnh, gây ấn tượng, tăng sức biểu cảm.

Ví dụ trong ca dao:

> Cày đồng đang buổi ban trưa,
> Mồ hôi thánh thót như mưa ruộng cày.

"Mồ hôi... như mưa" phóng đại để nhấn mạnh sự vất vả của người nông dân.

Một số thành ngữ nói quá: *nghĩ nát óc*, *ngáy như sấm*, *nhanh như cắt*, *chân cứng đá mềm*.

**Lưu ý**: nói quá khác **nói khoác**. Nói quá nhằm biểu cảm, người nghe đều hiểu đó là phóng đại; nói khoác nhằm khoe khoang, khiến người nghe tin vào điều không có thật.

## 2. Nói giảm nói tránh

**Nói giảm nói tránh** là cách diễn đạt **tế nhị, uyển chuyển** để tránh gây cảm giác quá đau buồn, nặng nề, hoặc để tránh thô tục, thiếu lịch sự.

Các cách thường dùng:

- Dùng từ Hán Việt đồng nghĩa: "từ trần", "tạ thế" thay cho "chết".
- Dùng cách nói vòng: "đã về với tổ tiên", "đã đi xa".
- Phủ định từ trái nghĩa: "Bài văn của em **chưa hay lắm**" thay cho "Bài văn của em dở".
- Dùng từ ngữ nhã nhặn: "người khiếm thị" thay cho "người mù".

## 3. Bài tập mẫu có lời giải

Xác định biện pháp tu từ:

1. "Nghĩ mãi không ra, em nghĩ nát cả óc." → **Nói quá** ("nát cả óc" phóng đại mức độ suy nghĩ).
2. "Bà cụ yếu lắm rồi, chắc sắp đi." → **Nói giảm nói tránh** ("đi" thay cho "chết").
3. Viết lại câu "Cậu lười quá!" theo cách nói giảm nói tránh → "Cậu **chưa chăm chỉ lắm**."

## 4. Dùng sao cho đúng?

- Nói giảm nói tránh thể hiện sự tôn trọng người nghe, nhưng **không nên** dùng khi cần nói thẳng, nói đúng sự thật (báo cáo, cảnh báo nguy hiểm).
- Nói quá cần phù hợp ngữ cảnh, tránh lạm dụng khiến lời nói thiếu chân thật.

> [!TIP]
> Mẹo nhận biết nhanh: nếu câu nói **làm to lên** mức độ sự việc thì là nói quá; nếu câu nói **làm nhẹ đi** cảm giác nặng nề thì là nói giảm nói tránh.`,
					},
					{
						Title:           "Viết bài văn biểu cảm về con người hoặc sự việc",
						DurationMinutes: 35,
						Body: `## 1. Yêu cầu của bài văn biểu cảm

- Giới thiệu được **đối tượng biểu cảm** (một người, một sự việc gắn bó với em).
- Bày tỏ được **tình cảm, cảm xúc chân thành** của người viết.
- Kết hợp **miêu tả, tự sự** để làm nổi bật cảm xúc, nhưng cảm xúc phải là chính.
- Người viết xưng "tôi" hoặc "em" (ngôi thứ nhất).

## 2. Bố cục

1. **Mở bài**: Giới thiệu người hoặc sự việc và ấn tượng, tình cảm chung.
2. **Thân bài**:
   - Nêu những đặc điểm nổi bật (ngoại hình, tính cách, việc làm) gắn với kỉ niệm cụ thể.
   - Bày tỏ cảm xúc của em qua từng chi tiết, từng kỉ niệm.
3. **Kết bài**: Khẳng định tình cảm, nêu mong ước hoặc lời hứa.

## 3. Ví dụ dàn ý

Đề: *Viết bài văn biểu cảm về người bà của em.*

- **Mở bài**: Bà là người gần gũi nhất với em từ thuở nhỏ; mỗi lần nghĩ đến bà, lòng em lại thấy ấm áp.
- **Thân bài**:
  - Hình ảnh bà: mái tóc bạc, đôi bàn tay gầy nhưng ấm, nụ cười hiền.
  - Kỉ niệm 1: những tối mất điện, bà quạt cho em ngủ và kể chuyện cổ tích.
  - Kỉ niệm 2: lần em bị ốm, bà thức suốt đêm nấu cháo, đắp khăn cho em.
  - Cảm xúc: biết ơn, yêu thương, có lúc ân hận vì từng làm bà buồn.
- **Kết bài**: Em mong bà luôn mạnh khoẻ; hứa sẽ chăm ngoan để bà vui lòng.

## 4. Đoạn văn mẫu (một phần thân bài)

Em nhớ nhất những buổi tối mùa hè mất điện. Bà ngồi bên giường, chiếc quạt nan trong tay bà cứ đều đều phe phẩy. Bà kể chuyện Tấm Cám, giọng chậm rãi, lúc trầm lúc bổng. Em chưa nghe hết chuyện đã thiếp đi, nhưng cảm giác mát lành từ chiếc quạt và giọng nói của bà thì em giữ mãi. Giờ đây, mỗi lần thấy bà ngồi lặng lẽ bên hiên, em lại muốn chạy đến ôm bà thật chặt.

## 5. Mẹo viết hay

- Chọn **một vài kỉ niệm thật**, đừng kể lan man cả cuộc đời nhân vật.
- Dùng câu cảm, từ ngữ biểu cảm, phép so sánh vừa phải.
- Không chép văn mẫu: cảm xúc vay mượn sẽ thiếu sức sống.

> [!TIP]
> Cảm xúc chỉ chạm đến người đọc khi gắn với chi tiết cụ thể. Thay vì viết "bà rất thương em", hãy kể một việc bà đã làm khiến em nhớ mãi.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: tục ngữ, truyện ngụ ngôn, biện pháp tu từ và văn biểu cảm",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Câu nào sau đây là tục ngữ (không phải thành ngữ)?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Tục ngữ là một câu trọn vẹn nêu một kinh nghiệm, lời khuyên. 'Có công mài sắt, có ngày nên kim' là một câu hoàn chỉnh; các phương án còn lại chỉ là cụm từ cố định (thành ngữ).",
					Options: []sampleOption{
						{Content: "Đầu voi đuôi chuột"},
						{Content: "Có công mài sắt, có ngày nên kim", IsCorrect: true},
						{Content: "Mẹ tròn con vuông"},
						{Content: "Nhanh như cắt"},
					},
				},
				{
					Prompt: "Đặc điểm nào sau đây đúng với truyện ngụ ngôn?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Truyện ngụ ngôn mượn chuyện loài vật, đồ vật hoặc con người để nói bóng gió chuyện đời và răn dạy một bài học; truyện thường ngắn, nhân vật ít khi có tên riêng.",
					Options: []sampleOption{
						{Content: "Kể về các vị thần sáng tạo ra trời đất"},
						{Content: "Có dung lượng dài, nhiều nhân vật có tên riêng"},
						{Content: "Luôn kết thúc có hậu, người tốt được đền đáp"},
						{Content: "Mượn chuyện loài vật, đồ vật hoặc con người để răn dạy một bài học", IsCorrect: true},
					},
				},
				{
					Prompt: "Bài học chính rút ra từ truyện 'Thầy bói xem voi' là gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Mỗi thầy bói chỉ sờ một bộ phận rồi khẳng định đó là cả con voi; truyện phê phán cách nhìn phiến diện và khuyên phải xem xét sự vật một cách toàn diện.",
					Options: []sampleOption{
						{Content: "Muốn hiểu đúng sự vật phải xem xét toàn diện, không nhìn phiến diện", IsCorrect: true},
						{Content: "Phải biết ơn những người đã giúp đỡ mình"},
						{Content: "Không nên kiêu ngạo vì môi trường sống chật hẹp"},
						{Content: "Phải siêng năng lao động mới có cái ăn"},
					},
				},
				{
					Prompt: "Câu 'Cụ ấy đã về với tổ tiên rồi.' sử dụng biện pháp tu từ nào?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "'Về với tổ tiên' là cách nói vòng thay cho từ 'chết', giúp giảm cảm giác đau buồn, nặng nề — đó là nói giảm nói tránh.",
					Options: []sampleOption{
						{Content: "Nói quá"},
						{Content: "So sánh"},
						{Content: "Nói giảm nói tránh", IsCorrect: true},
						{Content: "Nhân hoá"},
					},
				},
				{
					Prompt: "Khi viết bài văn biểu cảm về người bà, cách làm nào giúp cảm xúc chân thành, thuyết phục nhất?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Cảm xúc trong văn biểu cảm cần gắn với chi tiết, kỉ niệm cụ thể thì mới chân thật và chạm đến người đọc; liệt kê tính từ hay chép văn mẫu đều khiến bài viết sáo rỗng.",
					Options: []sampleOption{
						{Content: "Liệt kê thật nhiều tính từ khen ngợi bà"},
						{Content: "Gắn cảm xúc với một vài kỉ niệm, chi tiết cụ thể về bà", IsCorrect: true},
						{Content: "Kể lại toàn bộ cuộc đời của bà theo thứ tự năm tháng"},
						{Content: "Chép một bài văn mẫu hay rồi thay tên nhân vật"},
					},
				},
			},
		},
	},
	{
		Code:        "NGUVAN8",
		Title:       "Ngữ văn lớp 8",
		Description: "Lớp học mẫu môn Ngữ văn lớp 8 theo Chương trình GDPT 2018: đọc hiểu thơ thất ngôn bát cú Đường luật và truyện cười, từ tượng hình – tượng thanh và viết bài văn nghị luận về một vấn đề đời sống.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Đọc hiểu: thơ Đường luật và truyện cười",
				Lessons: []sampleLesson{
					{
						Title:           "Thơ thất ngôn bát cú Đường luật: Qua Đèo Ngang",
						DurationMinutes: 30,
						Body: `## 1. Văn bản

**Qua Đèo Ngang** — Bà Huyện Thanh Quan

> Bước tới Đèo Ngang bóng xế tà,
> Cỏ cây chen đá, lá chen hoa.
> Lom khom dưới núi tiều vài chú,
> Lác đác bên sông chợ mấy nhà.
> Nhớ nước đau lòng con quốc quốc,
> Thương nhà mỏi miệng cái gia gia.
> Dừng chân đứng lại trời non nước,
> Một mảnh tình riêng ta với ta.

## 2. Tác giả

Bà Huyện Thanh Quan tên thật là Nguyễn Thị Hinh, sống ở thế kỉ XIX, quê làng Nghi Tàm (nay thuộc Hà Nội). Chồng bà làm tri huyện Thanh Quan nên người đời quen gọi như vậy. Thơ bà trang nhã, đậm chất cổ điển và thường man mác buồn. Đèo Ngang thuộc dãy Hoành Sơn, nằm giữa hai tỉnh Hà Tĩnh và Quảng Bình.

## 3. Đặc điểm thể thơ thất ngôn bát cú Đường luật

- **Số câu, số chữ**: 8 câu, mỗi câu 7 chữ (56 chữ).
- **Bố cục**: đề (câu 1–2), thực (câu 3–4), luận (câu 5–6), kết (câu 7–8).
- **Vần**: một vần bằng, gieo ở cuối các câu 1, 2, 4, 6, 8. Trong bài: *tà – hoa – nhà – gia – ta*.
- **Đối**: hai câu thực và hai câu luận đối nhau từng cặp.
- **Luật bằng trắc, niêm** chặt chẽ.

## 4. Phân tích

- **Hai câu đề**: thời điểm "bóng xế tà" gợi buồn; điệp từ "chen" cho thấy cảnh hoang sơ, rậm rạp.
- **Hai câu thực**: đảo ngữ và từ tượng hình "lom khom", "lác đác" khắc hoạ con người nhỏ bé, thưa thớt giữa núi đèo.
- **Hai câu luận**: tiếng chim quốc quốc, gia gia gợi "nước" (quốc) và "nhà" (gia) — lối chơi chữ đồng âm diễn tả nỗi nhớ nước thương nhà.
- **Hai câu kết**: không gian "trời non nước" bao la đối lập với "một mảnh tình riêng"; cụm "ta với ta" diễn tả nỗi cô đơn tuyệt đối.

> [!NOTE]
> Hai cặp câu thực và luận đối nhau rất chỉnh: "lom khom" – "lác đác", "dưới núi" – "bên sông", "nhớ nước" – "thương nhà". Đây là dấu hiệu dễ nhận ra của thơ Đường luật.`,
					},
					{
						Title:           "Truyện cười dân gian: Treo biển và Lợn cưới, áo mới",
						DurationMinutes: 30,
						Body: `## 1. Truyện cười là gì?

Truyện cười là truyện kể dân gian ngắn, dùng **tiếng cười** để mua vui hoặc **phê phán** những thói hư tật xấu trong xã hội.

Các yếu tố gây cười thường gặp:

- **Mâu thuẫn** trái tự nhiên giữa lời nói và việc làm, giữa cái bên ngoài và bản chất.
- **Chi tiết gây cười** được đẩy lên cao dần.
- **Kết thúc bất ngờ**, thường bằng một câu nói "lỡ lời" của nhân vật.

## 2. Treo biển

**Tóm tắt**: Một nhà hàng treo tấm biển "Ở đây có bán cá tươi". Người thứ nhất góp ý bỏ chữ "tươi" vì bán cá thì ai chẳng bán cá tươi. Người thứ hai bảo bỏ "ở đây". Người thứ ba bảo bỏ "có bán". Người thứ tư nói chưa đến đầu phố đã ngửi thấy mùi cá, đề biển làm gì. Nghe ai góp ý nhà hàng cũng làm theo, cuối cùng cất luôn tấm biển.

**Ý nghĩa**: Phê phán sự **thiếu chủ kiến**, nghe ai cũng làm theo mà không suy xét. Bài học: tiếp thu ý kiến người khác phải có chọn lọc.

## 3. Lợn cưới, áo mới

**Tóm tắt**: Một anh có tính hay khoe may được chiếc áo mới, đứng ở cửa từ sáng đến chiều chờ người khen. Một anh khác cũng hay khoe, bị sổng lợn, chạy đến hỏi: "Bác có thấy con lợn cưới của tôi chạy qua đây không?" Anh kia liền giơ vạt áo ra trả lời rằng từ lúc mặc chiếc áo mới này, anh chẳng thấy con lợn nào chạy qua cả.

**Yếu tố gây cười**:

- Từ thừa "cưới" trong câu hỏi và cả câu trả lời dài dòng, khoe áo mới, đều là thông tin không cần thiết.
- Cái cười bật ra khi hai kẻ hay khoe "gặp nhau".

**Ý nghĩa**: Chế giễu thói **khoe khoang** — một tính xấu khá phổ biến.

## 4. Cách đọc hiểu một truyện cười

1. Xác định **đối tượng** bị cười (ai, thói xấu gì).
2. Chỉ ra **mâu thuẫn** và **chi tiết gây cười**.
3. Rút ra **ý nghĩa phê phán** và bài học cho bản thân.

> [!TIP]
> Muốn "bắt" được tiếng cười, hãy đọc to lời thoại của nhân vật — cái hài thường nằm ở chính câu nói thừa, nói hớ hoặc ngây ngô.`,
					},
				},
			},
			{
				Title: "Chương 2. Tiếng Việt và kỹ năng viết",
				Lessons: []sampleLesson{
					{
						Title:           "Từ tượng hình và từ tượng thanh",
						DurationMinutes: 25,
						Body: `## 1. Khái niệm

- **Từ tượng hình** là từ gợi tả hình ảnh, dáng vẻ, trạng thái của sự vật: *lom khom, lác đác, móm mém, lênh khênh, khấp khểnh, lả lướt*.
- **Từ tượng thanh** là từ mô phỏng âm thanh của tự nhiên, con người: *róc rách, ríu rít, ầm ầm, lộp độp, hu hu, lao xao*.

## 2. Tác dụng

Từ tượng hình, tượng thanh giúp câu văn, câu thơ **gợi hình, gợi thanh** cụ thể, sinh động và giàu sức biểu cảm. Chúng được dùng nhiều trong văn miêu tả, tự sự và thơ.

## 3. Ví dụ phân tích

> Lom khom dưới núi tiều vài chú,
> Lác đác bên sông chợ mấy nhà.
> (Bà Huyện Thanh Quan)

- "Lom khom" gợi dáng người tiều phu cúi mình vất vả.
- "Lác đác" gợi sự thưa thớt, ít ỏi của mấy ngôi nhà ven sông.
- Hai từ tượng hình đặt ở đầu câu (đảo ngữ) nhấn mạnh cảnh vắng vẻ, hiu quạnh nơi đèo núi.

Trong câu thơ "Lao xao chợ cá làng ngư phủ" của Nguyễn Trãi, từ tượng thanh "lao xao" gợi âm thanh huyên náo, đông vui của cảnh chợ quê.

## 4. Bài tập mẫu có lời giải

Tìm từ tượng hình, từ tượng thanh trong đoạn văn:

*Mưa rơi lộp độp trên mái tôn. Ngoài vườn, mấy tàu lá chuối rách tả tơi. Tiếng ếch nhái kêu ộp oạp khắp cánh đồng.*

- Từ tượng thanh: **lộp độp**, **ộp oạp**.
- Từ tượng hình: **tả tơi**.

## 5. Lưu ý

Phần lớn từ tượng hình, tượng thanh là **từ láy**, nhưng không phải từ láy nào cũng là từ tượng hình hay tượng thanh (ví dụ: *chăm chỉ*, *nhỏ nhẹ*).

> [!TIP]
> Khi viết văn miêu tả, thử thay một tính từ chung chung (như "đi chậm") bằng một từ tượng hình (như "lững thững", "lò dò") — câu văn sẽ sống động hơn hẳn.`,
					},
					{
						Title:           "Viết bài văn nghị luận về một vấn đề của đời sống",
						DurationMinutes: 35,
						Body: `## 1. Các yếu tố của bài nghị luận

- **Luận đề**: vấn đề chính được bàn luận trong toàn bài.
- **Luận điểm**: các ý kiến cụ thể triển khai luận đề.
- **Lí lẽ**: những lời giải thích, phân tích để làm sáng tỏ luận điểm.
- **Bằng chứng**: sự việc, số liệu, ví dụ xác thực minh hoạ cho lí lẽ.

## 2. Bố cục

1. **Mở bài**: Giới thiệu vấn đề và nêu ý kiến của người viết.
2. **Thân bài**: Trình bày các luận điểm, mỗi luận điểm có lí lẽ và bằng chứng; xem xét ý kiến trái chiều.
3. **Kết bài**: Khẳng định lại ý kiến, nêu bài học hoặc lời kêu gọi.

## 3. Ví dụ dàn ý

Đề: *Tác hại của việc nghiện điện thoại thông minh ở học sinh.*

- **Luận đề**: Nghiện điện thoại thông minh gây nhiều tác hại cho học sinh, cần được khắc phục.
- **Luận điểm 1 – Thực trạng**: nhiều bạn dùng điện thoại hàng giờ mỗi ngày, cả trong giờ ăn, giờ ngủ.
- **Luận điểm 2 – Tác hại**: ảnh hưởng thị lực, giấc ngủ; giảm tập trung, sa sút học tập; ít giao tiếp trực tiếp với gia đình, bạn bè.
- **Luận điểm 3 – Nguyên nhân**: sức hút của trò chơi, mạng xã hội; thiếu kĩ năng tự quản lí thời gian; thiếu hoạt động thay thế.
- **Ý kiến trái chiều**: điện thoại cũng có ích (học trực tuyến, tra cứu) — vấn đề nằm ở cách dùng, không phải ở thiết bị.
- **Giải pháp**: đặt giới hạn thời gian, tắt thông báo khi học, tham gia thể thao, câu lạc bộ.

## 4. Đoạn văn mẫu (luận điểm 2)

Trước hết, nghiện điện thoại khiến sức khoẻ của học sinh bị ảnh hưởng rõ rệt. Việc nhìn màn hình quá lâu làm mắt mỏi, khô và dễ dẫn đến cận thị. Nhiều bạn mang điện thoại lên giường, lướt mạng đến khuya, sáng hôm sau đến lớp với đôi mắt thâm quầng và đầu óc uể oải. Khi sức khoẻ giảm sút, việc học tập tất yếu cũng đi xuống.

## 5. Lưu ý khi dùng bằng chứng

- Bằng chứng phải **xác thực**: nếu dẫn số liệu, cần nêu rõ nguồn đáng tin cậy.
- Ưu tiên những ví dụ gần gũi, cụ thể trong đời sống học đường.

> [!NOTE]
> Một luận điểm thuyết phục luôn đi kèm "bộ đôi" lí lẽ và bằng chứng. Thiếu lí lẽ thì bài viết rời rạc, thiếu bằng chứng thì bài viết thiếu sức nặng.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: thơ Đường luật, truyện cười, từ tượng hình – tượng thanh và văn nghị luận",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Một bài thơ thất ngôn bát cú Đường luật có bao nhiêu câu và mỗi câu mấy chữ?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "'Thất ngôn' là bảy chữ, 'bát cú' là tám câu: bài thơ gồm 8 câu, mỗi câu 7 chữ.",
					Options: []sampleOption{
						{Content: "4 câu, mỗi câu 7 chữ"},
						{Content: "8 câu, mỗi câu 5 chữ"},
						{Content: "6 câu, mỗi câu 8 chữ"},
						{Content: "8 câu, mỗi câu 7 chữ", IsCorrect: true},
					},
				},
				{
					Prompt: "Từ nào sau đây là từ tượng thanh?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "'Ríu rít' mô phỏng âm thanh (tiếng chim, tiếng nói cười) nên là từ tượng thanh; 'lom khom', 'lác đác', 'lênh khênh' gợi hình dáng, trạng thái nên là từ tượng hình.",
					Options: []sampleOption{
						{Content: "lom khom"},
						{Content: "ríu rít", IsCorrect: true},
						{Content: "lác đác"},
						{Content: "lênh khênh"},
					},
				},
				{
					Prompt: "Câu thơ 'Nhớ nước đau lòng con quốc quốc' (Qua Đèo Ngang) gợi nỗi nhớ nước nhờ cách nào?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Tiếng chim 'quốc quốc' đồng âm với chữ 'quốc' (nghĩa là nước), nhà thơ mượn tiếng chim để gợi nỗi nhớ nước — đó là lối chơi chữ đồng âm.",
					Options: []sampleOption{
						{Content: "Phóng đại tiếng chim kêu thật to"},
						{Content: "So sánh tiếng chim với tiếng đàn"},
						{Content: "Dựa vào sự đồng âm giữa tiếng chim 'quốc' và chữ 'quốc' nghĩa là nước", IsCorrect: true},
						{Content: "Miêu tả hình dáng con chim quốc"},
					},
				},
				{
					Prompt: "Truyện cười 'Treo biển' chủ yếu phê phán điều gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Nhà hàng nghe ai góp ý cũng làm theo, cuối cùng cất luôn tấm biển: truyện phê phán sự thiếu chủ kiến, không biết suy xét khi tiếp thu ý kiến.",
					Options: []sampleOption{
						{Content: "Sự thiếu chủ kiến, ai góp ý gì cũng làm theo mà không suy xét", IsCorrect: true},
						{Content: "Thói khoe khoang của cải"},
						{Content: "Thói tham lam, keo kiệt"},
						{Content: "Tính lười biếng, ham chơi"},
					},
				},
				{
					Prompt: "Với đề 'Tác hại của việc nghiện điện thoại thông minh ở học sinh', câu nào sau đây có thể dùng làm một luận điểm?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Luận điểm là ý kiến cụ thể triển khai luận đề. 'Nghiện điện thoại làm giảm sút sức khoẻ và kết quả học tập' trực tiếp nêu một tác hại; các câu khác chỉ là thông tin chung hoặc cảm nhận cá nhân.",
					Options: []sampleOption{
						{Content: "Điện thoại thông minh đã xuất hiện từ nhiều năm trước."},
						{Content: "Bạn Nam lớp em vừa được mua một chiếc điện thoại mới."},
						{Content: "Em rất thích xem video trên điện thoại."},
						{Content: "Nghiện điện thoại làm học sinh giảm sút sức khoẻ và kết quả học tập.", IsCorrect: true},
					},
				},
			},
		},
	},
	{
		Code:        "NGUVAN9",
		Title:       "Ngữ văn lớp 9",
		Description: "Lớp học mẫu môn Ngữ văn lớp 9 theo Chương trình GDPT 2018: đọc hiểu truyện truyền kì và truyện thơ Nôm, cách dẫn trực tiếp – gián tiếp và viết bài văn phân tích tác phẩm truyện.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Đọc hiểu: truyện truyền kì và truyện thơ Nôm",
				Lessons: []sampleLesson{
					{
						Title:           "Truyện truyền kì: Chuyện người con gái Nam Xương",
						DurationMinutes: 35,
						Body: `## 1. Truyện truyền kì là gì?

Truyện truyền kì là thể loại văn xuôi tự sự thời trung đại, thường viết bằng chữ Hán, **kết hợp yếu tố hiện thực với yếu tố kì ảo** (hoang đường) để phản ánh hiện thực và gửi gắm khát vọng của con người. "Truyền kì" có nghĩa là ghi chép, lưu truyền những điều kì lạ.

## 2. Tác giả và tác phẩm

Nguyễn Dữ sống ở thế kỉ XVI, quê ở Hải Dương. Tác phẩm *Truyền kì mạn lục* của ông gồm 20 truyện, được đánh giá là "thiên cổ kì bút" (áng văn hay của muôn đời). "Chuyện người con gái Nam Xương" là một truyện trong tập, có nguồn gốc từ truyện dân gian "Vợ chàng Trương".

## 3. Tóm tắt

Vũ Nương (Vũ Thị Thiết), người con gái Nam Xương thuỳ mị, nết na, lấy Trương Sinh — con nhà hào phú, ít học, tính đa nghi. Trương Sinh phải đi lính; ở nhà, Vũ Nương sinh con, chăm sóc mẹ chồng chu đáo và lo ma chay khi bà mất. Đêm đêm, nàng chỉ bóng mình trên vách, bảo con đó là cha Đản.

Trương Sinh trở về, nghe con nói về "người cha" đêm nào cũng đến, liền nghi vợ thất tiết, mắng nhiếc rồi đuổi đi. Không thể minh oan, Vũ Nương gieo mình xuống sông Hoàng Giang và được các nàng tiên dưới thuỷ cung cứu. Ít lâu sau, một đêm bé Đản chỉ bóng Trương Sinh trên vách và gọi đó là cha, chàng mới hiểu ra nỗi oan của vợ.

Phan Lang, người cùng làng, gặp Vũ Nương dưới thuỷ cung và mang tín vật về. Trương Sinh lập đàn giải oan bên bến sông; Vũ Nương hiện về giữa dòng, nói lời tạ từ rồi biến mất.

## 4. Giá trị nội dung

- **Giá trị hiện thực**: phản ánh số phận bi kịch của người phụ nữ trong xã hội phong kiến nam quyền, chiến tranh phi nghĩa gây cảnh chia lìa.
- **Giá trị nhân đạo**: ca ngợi vẻ đẹp phẩm chất của Vũ Nương, bày tỏ niềm cảm thông sâu sắc với người phụ nữ.

## 5. Nghệ thuật

- **Chi tiết cái bóng**: vừa thể hiện tình thương con, nhớ chồng của Vũ Nương, vừa là "nút thắt" gây ra bi kịch và cũng là chi tiết "mở nút" giải oan.
- **Yếu tố kì ảo** (thuỷ cung, Vũ Nương hiện về) tạo kết thúc phần nào có hậu, thể hiện ước mơ công lí của nhân dân, nhưng nàng vẫn không thể trở về — bi kịch vẫn còn nguyên.

> [!TIP]
> Khi đọc truyện truyền kì, hãy tách riêng hai lớp: lớp hiện thực (con người, xã hội) và lớp kì ảo (thần tiên, cõi khác). Ý nghĩa sâu xa thường nằm ở sự đan cài giữa hai lớp ấy.`,
					},
					{
						Title:           "Truyện thơ Nôm: Truyện Kiều — đoạn trích Cảnh ngày xuân",
						DurationMinutes: 30,
						Body: `## 1. Tác giả và tác phẩm

Nguyễn Du (1765–1820), tự Tố Như, hiệu Thanh Hiên, quê làng Tiên Điền, huyện Nghi Xuân, tỉnh Hà Tĩnh — đại thi hào của dân tộc. *Truyện Kiều* (tên chữ: *Đoạn trường tân thanh*) là truyện thơ Nôm gồm 3254 câu thơ lục bát, dựa trên cốt truyện *Kim Vân Kiều truyện* của Thanh Tâm Tài Nhân (Trung Quốc), nhưng được Nguyễn Du sáng tạo lại với giá trị nghệ thuật vượt trội.

## 2. Thể thơ lục bát

- Từng cặp câu 6 chữ (lục) và 8 chữ (bát).
- Chữ thứ 6 câu lục vần với chữ thứ 6 câu bát; chữ thứ 8 câu bát vần với chữ thứ 6 câu lục tiếp theo.

## 3. Văn bản (bốn câu đầu đoạn trích)

> Ngày xuân con én đưa thoi,
> Thiều quang chín chục đã ngoài sáu mươi.
> Cỏ non xanh tận chân trời,
> Cành lê trắng điểm một vài bông hoa.

Thử kiểm tra vần: "thoi" (chữ 6 câu lục) vần với "ngoài" (chữ 6 câu bát); "mươi" (chữ 8 câu bát) vần với "trời" (chữ 6 câu lục kế tiếp).

## 4. Phân tích

- **Thời gian mùa xuân**: hình ảnh "con én đưa thoi" (chim én bay qua lại như con thoi dệt cửi) gợi thời gian trôi nhanh. "Thiều quang chín chục đã ngoài sáu mươi": chín mươi ngày xuân tươi đẹp đã qua hơn sáu mươi ngày, tức là đã sang tháng ba.
- **Bức hoạ mùa xuân**: nền cỏ non xanh trải tới tận chân trời, điểm xuyết vài bông lê trắng — màu sắc hài hoà, không gian khoáng đạt, trong trẻo. Chữ "điểm" khiến bức tranh thêm sinh động.
- **Khung cảnh lễ hội**: tiết Thanh minh có lễ tảo mộ (viếng, sửa sang mộ người thân) và hội đạp thanh (đi chơi xuân nơi đồng nội); không khí "nô nức" của nam thanh nữ tú.
- **Cảnh chị em du xuân trở về**: các từ láy "tà tà", "thanh thanh", "nao nao" gợi cảnh chiều nhạt dần và tâm trạng bâng khuâng, như báo trước điều sắp xảy ra.

## 5. Nghệ thuật

Bút pháp ước lệ, tả cảnh ngụ tình; ngôn ngữ giàu hình ảnh; sử dụng tài tình từ ghép, từ láy.

> [!NOTE]
> Năm 2015, nhân kỉ niệm 250 năm năm sinh, Nguyễn Du được UNESCO đưa vào danh sách các danh nhân văn hoá được vinh danh.`,
					},
				},
			},
			{
				Title: "Chương 2. Tiếng Việt và kỹ năng viết",
				Lessons: []sampleLesson{
					{
						Title:           "Cách dẫn trực tiếp và cách dẫn gián tiếp",
						DurationMinutes: 25,
						Body: `## 1. Cách dẫn trực tiếp

**Dẫn trực tiếp** là nhắc lại **nguyên văn** lời nói hay ý nghĩ của một người hoặc nhân vật. Lời dẫn trực tiếp được đặt trong **dấu ngoặc kép**, thường đứng sau dấu hai chấm.

Ví dụ: Bé Đản chỉ bóng trên vách và nói: "Cha Đản lại đến kia kìa!"

## 2. Cách dẫn gián tiếp

**Dẫn gián tiếp** là thuật lại lời nói hay ý nghĩ của người khác, **có điều chỉnh** cho phù hợp; lời dẫn **không** đặt trong dấu ngoặc kép, có thể thêm "rằng" hoặc "là" phía trước.

Ví dụ: Bé Đản chỉ bóng trên vách và bảo rằng cha Đản lại đến.

## 3. Chuyển lời dẫn trực tiếp thành gián tiếp

Các bước:

1. Bỏ dấu hai chấm và dấu ngoặc kép, có thể thêm "rằng" hoặc "là".
2. Thay đổi từ xưng hô cho phù hợp (tớ, tôi → tên người nói hoặc "bạn ấy"...).
3. Thay đổi từ chỉ thời gian, nơi chốn (ngày mai → hôm sau, đây → đó).
4. Lược bỏ từ hô gọi, từ cảm thán không cần thiết.

**Bài tập mẫu có lời giải**

Câu gốc: Lan nói với Hoa: "Ngày mai tớ sẽ mang sách đến đây cho cậu."

Chuyển thành gián tiếp: Lan nói với Hoa rằng hôm sau Lan sẽ mang sách đến đó cho Hoa.

## 4. Chuyển lời dẫn gián tiếp thành trực tiếp

Làm ngược lại: khôi phục nguyên văn lời nói, thêm dấu hai chấm và dấu ngoặc kép, điều chỉnh từ xưng hô theo người nói.

Ví dụ: Thầy dặn chúng tôi rằng phải ôn bài kĩ trước khi thi. → Thầy dặn chúng tôi: "Các em phải ôn bài kĩ trước khi thi."

## 5. Dẫn lời trong bài nghị luận văn học

- Khi trích thơ, văn **trực tiếp**, phải chép **chính xác** từng chữ và đặt trong ngoặc kép.
- Khi không nhớ nguyên văn, nên dùng cách dẫn **gián tiếp** (tóm ý) để tránh trích sai.

> [!TIP]
> Dấu hiệu nhanh nhất để phân biệt: có dấu ngoặc kép bao quanh nguyên văn lời nói là dẫn trực tiếp; không có ngoặc kép, lời nói đã được thuật lại là dẫn gián tiếp.`,
					},
					{
						Title:           "Viết bài văn nghị luận phân tích một tác phẩm truyện",
						DurationMinutes: 35,
						Body: `## 1. Yêu cầu

- Giới thiệu tác giả, tác phẩm và nêu **nhận định chung** về truyện.
- Phân tích **nội dung chủ đề** của truyện.
- Phân tích một số **nét đặc sắc nghệ thuật**: cốt truyện, tình huống, nhân vật, chi tiết, người kể chuyện.
- Mỗi luận điểm đều có **lí lẽ** và **bằng chứng** lấy từ tác phẩm.

## 2. Bố cục

1. **Mở bài**: Giới thiệu tác giả, tác phẩm; nêu ý kiến khái quát.
2. **Thân bài**: Lần lượt phân tích các luận điểm về nội dung và nghệ thuật.
3. **Kết bài**: Khẳng định giá trị của tác phẩm, nêu cảm nhận của bản thân.

## 3. Dàn ý mẫu

Đề: *Phân tích truyện "Chuyện người con gái Nam Xương" của Nguyễn Dữ.*

- **Mở bài**: Nguyễn Dữ và *Truyền kì mạn lục*; truyện là tiếng nói cảm thông với số phận người phụ nữ.
- **Luận điểm 1**: Vũ Nương — người phụ nữ thuỳ mị, nết na, thuỷ chung, hiếu thảo (dẫn chứng: cách cư xử khi chồng đi lính, chăm sóc mẹ chồng).
- **Luận điểm 2**: Bi kịch của Vũ Nương — bị chồng nghi oan, phải tìm đến cái chết; nguyên nhân từ sự đa nghi, gia trưởng của Trương Sinh và xã hội nam quyền.
- **Luận điểm 3**: Nghệ thuật — chi tiết cái bóng, yếu tố kì ảo, kết thúc vừa có hậu vừa bi kịch.
- **Kết bài**: Giá trị hiện thực và nhân đạo; bài học về niềm tin, sự thấu hiểu trong gia đình.

## 4. Đoạn văn mẫu (luận điểm 3)

Thành công nổi bật của truyện là chi tiết cái bóng. Lời nói đùa của Vũ Nương để dỗ con, xuất phát từ tình yêu thương chồng con, lại trở thành nguyên nhân trực tiếp đẩy nàng vào cái chết. Cũng chính cái bóng ấy, khi bé Đản chỉ vào bóng cha trên vách, đã giúp Trương Sinh nhận ra sự thật. Một chi tiết nhỏ mà vừa thắt nút, vừa mở nút câu chuyện, làm cho bi kịch càng thêm day dứt.

## 5. Lỗi thường gặp

- **Kể lại truyện** thay vì phân tích, đánh giá.
- Dẫn chứng chung chung, không gắn với lí lẽ.
- Trích dẫn sai nguyên văn tác phẩm.

> [!TIP]
> Sau mỗi dẫn chứng, hãy tự hỏi "Chi tiết này cho thấy điều gì?" và viết ra câu trả lời — đó chính là phần phân tích mà người chấm muốn đọc.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: truyện truyền kì, Truyện Kiều, cách dẫn lời và văn phân tích truyện",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "'Chuyện người con gái Nam Xương' là tác phẩm của ai, nằm trong tập truyện nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Truyện do Nguyễn Dữ (thế kỉ XVI) viết, nằm trong tập Truyền kì mạn lục gồm 20 truyện.",
					Options: []sampleOption{
						{Content: "Nguyễn Du — Truyện Kiều"},
						{Content: "Nguyễn Trãi — Quốc âm thi tập"},
						{Content: "Nguyễn Dữ — Truyền kì mạn lục", IsCorrect: true},
						{Content: "Bà Huyện Thanh Quan — Qua Đèo Ngang"},
					},
				},
				{
					Prompt: "Truyện Kiều của Nguyễn Du được viết bằng thể thơ nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Truyện Kiều là truyện thơ Nôm gồm 3254 câu thơ lục bát.",
					Options: []sampleOption{
						{Content: "Lục bát", IsCorrect: true},
						{Content: "Thất ngôn bát cú Đường luật"},
						{Content: "Song thất lục bát"},
						{Content: "Thơ tự do"},
					},
				},
				{
					Prompt: "Câu thơ 'Thiều quang chín chục đã ngoài sáu mươi' cho biết điều gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Mùa xuân có ba tháng, khoảng chín mươi ngày; 'đã ngoài sáu mươi' nghĩa là đã qua hơn sáu mươi ngày, tức là đã sang tháng thứ ba của mùa xuân.",
					Options: []sampleOption{
						{Content: "Nhân vật trong truyện đã ngoài sáu mươi tuổi"},
						{Content: "Lễ hội mùa xuân kéo dài chín mươi ngày"},
						{Content: "Mùa xuân năm ấy có sáu mươi ngày nắng"},
						{Content: "Chín mươi ngày xuân đã trôi qua hơn sáu mươi ngày, tức đã sang tháng ba", IsCorrect: true},
					},
				},
				{
					Prompt: "Câu nào sau đây sử dụng cách dẫn gián tiếp?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Câu 'Bé Đản bảo rằng cha Đản lại đến.' thuật lại lời nói, không dùng dấu ngoặc kép và có từ 'rằng' — đó là dẫn gián tiếp. Các câu còn lại đều giữ nguyên văn lời nói trong ngoặc kép.",
					Options: []sampleOption{
						{Content: "Bé Đản nói: \"Cha Đản lại đến kia kìa!\""},
						{Content: "Bé Đản bảo rằng cha Đản lại đến.", IsCorrect: true},
						{Content: "\"Cha Đản lại đến kia kìa!\", bé Đản reo lên."},
						{Content: "Bé Đản chỉ lên vách: \"Kia kìa!\""},
					},
				},
				{
					Prompt: "Khi viết bài phân tích truyện 'Chuyện người con gái Nam Xương', cách triển khai nào đúng yêu cầu của bài nghị luận phân tích tác phẩm?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Bài phân tích cần nêu luận điểm rõ ràng và dùng dẫn chứng trong truyện để phân tích, đánh giá; kể lại cốt truyện hay chép nguyên văn dài không phải là phân tích.",
					Options: []sampleOption{
						{Content: "Chỉ nêu cảm xúc yêu thích truyện, không cần dẫn chứng"},
						{Content: "Kể lại thật đầy đủ toàn bộ cốt truyện từ đầu đến cuối"},
						{Content: "Nêu luận điểm về nhân vật, chi tiết nghệ thuật rồi dùng dẫn chứng trong truyện để phân tích, đánh giá", IsCorrect: true},
						{Content: "Chép nguyên văn nhiều đoạn dài trong truyện rồi tóm tắt lại"},
					},
				},
			},
		},
	},
	{
		Code:        "NGUVAN10",
		Title:       "Ngữ văn lớp 10",
		Description: "Lớp học mẫu môn Ngữ văn lớp 10 theo Chương trình GDPT 2018: đọc hiểu thần thoại và sử thi, nhận diện – sửa lỗi dùng từ và viết bài văn nghị luận về một vấn đề xã hội.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Đọc hiểu: thần thoại và sử thi",
				Lessons: []sampleLesson{
					{
						Title:           "Thần thoại: Thần Trụ Trời",
						DurationMinutes: 30,
						Body: `## 1. Thần thoại là gì?

Thần thoại là thể loại tự sự dân gian ra đời sớm nhất, kể về các vị thần nhằm **giải thích nguồn gốc vũ trụ, thế giới tự nhiên và con người** theo quan niệm của người xưa. Có thể chia thành hai nhóm chính:

- **Thần thoại suy nguyên**: giải thích nguồn gốc trời đất, núi sông, muôn loài.
- **Thần thoại sáng tạo văn hoá**: kể về việc tìm ra lửa, trồng lúa, làm nhà...

"Thần Trụ Trời" là một thần thoại suy nguyên tiêu biểu của người Việt.

## 2. Tóm tắt truyện

Thuở trời đất còn chưa có muôn vật và loài người, mọi thứ chỉ là một khối hỗn độn, tối tăm và lạnh lẽo. Bỗng xuất hiện một vị thần khổng lồ. Thần đứng dậy, đội trời lên, rồi đào đất, khuân đá đắp thành một cột lớn để chống trời. Cột càng cao, trời càng bị đẩy lên xa đất. Khi trời đất đã phân đôi, thần phá cột, ném đất đá đi khắp nơi; từ đó mặt đất có chỗ cao, chỗ thấp, thành núi non, gò đồi.

## 3. Đặc điểm nghệ thuật

- **Nhân vật** là vị thần có hình dáng khổng lồ, sức mạnh phi thường, hành động mang tính sáng tạo thế giới.
- **Không gian, thời gian** mơ hồ, phiếm chỉ: thuở hỗn mang, khi chưa có thế gian.
- **Cách giải thích** hiện tượng tự nhiên (trời cao, đất thấp, núi đồi) hồn nhiên, chất phác, phản ánh tư duy của người cổ đại.

## 4. Giá trị

- Thể hiện **khát vọng nhận thức, lí giải thế giới** và trí tưởng tượng phong phú của người Việt cổ.
- Ca ngợi sức lao động sáng tạo: thần cũng "đào đất, khuân đá" như con người.
- Là nguồn tư liệu quý về vũ trụ quan, tín ngưỡng dân gian.

## 5. Câu hỏi gợi mở

Vì sao người xưa lại hình dung vị thần tạo lập trời đất bằng những công việc như đào đất, đắp cột? Gợi ý: vì người kể là cư dân nông nghiệp, quen với lao động chân tay, nên họ "đo" công việc của thần bằng chính kinh nghiệm lao động của mình.

> [!NOTE]
> Thần thoại suy nguyên và thần thoại sáng tạo văn hoá đều thuộc dòng thần thoại, nhưng khác nhau ở đối tượng được giải thích: một bên là thế giới tự nhiên, một bên là thành tựu văn hoá của con người.`,
					},
					{
						Title:           "Sử thi: Đăm Săn — đoạn trích Chiến thắng Mtao Mxây",
						DurationMinutes: 30,
						Body: `## 1. Giới thiệu chung

*Đăm Săn* là bộ sử thi anh hùng nổi tiếng của dân tộc **Ê-đê** (Tây Nguyên), kể về cuộc đời và chiến công của tù trưởng Đăm Săn. Đoạn trích "Chiến thắng Mtao Mxây" kể lại cuộc giao chiến giữa Đăm Săn và tù trưởng Mtao Mxây — kẻ đã cướp Hơ Nhị, vợ của Đăm Săn.

## 2. Tóm tắt đoạn trích

- Đăm Săn đến tận nhà Mtao Mxây, thách đấu.
- Mtao Mxây ra đánh, múa khiên trước nhưng vụng về, kêu lạch xạch; Đăm Săn múa khiên mạnh mẽ, nhanh như gió, khiến đối thủ hoảng sợ.
- Hơ Nhị ném cho chồng miếng trầu; Đăm Săn đớp được, sức mạnh tăng lên gấp bội.
- Đăm Săn đâm trúng Mtao Mxây nhưng không thủng vì hắn mặc áo giáp. Ông Trời hiện ra mách Đăm Săn lấy chày mòn ném vào vành tai kẻ thù. Mtao Mxây ngã gục và bị giết.
- Dân làng Mtao Mxây tự nguyện đi theo Đăm Săn; cả cộng đồng mở tiệc ăn mừng chiến thắng.

## 3. Đặc điểm nghệ thuật sử thi

- **Nhân vật anh hùng** mang vẻ đẹp lí tưởng của cộng đồng: dũng mãnh, tài giỏi, được thần linh trợ giúp.
- **Ngôn ngữ** giàu hình ảnh, dùng nhiều phép so sánh, phóng đại, so sánh trùng điệp khi miêu tả sức mạnh, vẻ đẹp nhân vật.
- **Ý nghĩa cộng đồng**: chiến công của Đăm Săn không chỉ vì danh dự cá nhân mà vì sự thịnh vượng của cả buôn làng.

## 4. Ý nghĩa

Đoạn trích ca ngợi vẻ đẹp, sức mạnh của người anh hùng sử thi, đồng thời thể hiện khát vọng về một cuộc sống ấm no, đoàn kết, hùng mạnh của các dân tộc Tây Nguyên xưa.

## 5. So sánh nhanh hai nhân vật

| Tiêu chí | Đăm Săn | Mtao Mxây |
|---|---|---|
| Thái độ | Đường hoàng, tự tin | Kiêu ngạo nhưng hèn nhát |
| Tài múa khiên | Mạnh mẽ, uyển chuyển | Vụng về, chậm chạp |
| Kết cục | Chiến thắng, được dân làng tôn vinh | Thất bại, bị tiêu diệt |

> [!TIP]
> Khi đọc sử thi, chú ý các đoạn miêu tả bằng so sánh trùng điệp, phóng đại — đây là dấu hiệu đặc trưng để nhận diện thể loại.`,
					},
				},
			},
			{
				Title: "Chương 2. Tiếng Việt và kỹ năng viết",
				Lessons: []sampleLesson{
					{
						Title:           "Lỗi dùng từ, lỗi trật tự từ và cách sửa",
						DurationMinutes: 25,
						Body: `## 1. Lỗi lặp từ

- Sai: *Truyện Thánh Gióng là truyện hay nên em rất thích đọc truyện Thánh Gióng.*
- Sửa: *Thánh Gióng là truyện hay nên em rất thích đọc.*

## 2. Lỗi lẫn lộn các từ gần âm

| Dùng sai | Dùng đúng | Giải thích |
|---|---|---|
| thăm quan | tham quan | "tham quan" là xem tận mắt để mở rộng hiểu biết |
| bàng quang (trước nỗi đau) | bàng quan | "bàng quan" là thờ ơ, đứng ngoài cuộc; "bàng quang" là một cơ quan trong cơ thể |
| bức tranh rất linh động | bức tranh rất sinh động | "linh động" là không cứng nhắc; "sinh động" là có sức sống, gợi cảm |

## 3. Lỗi dùng từ không đúng nghĩa

- Sai: *Bạn cần khắc phục những yếu điểm trong cách học.* ("yếu điểm" là điểm quan trọng)
- Sửa: *Bạn cần khắc phục những **điểm yếu** trong cách học.*
- Sai: *Anh ấy là người rất kiên cố.* ("kiên cố" dùng cho công trình vững chắc)
- Sửa: *Anh ấy là người rất **kiên định**.*

## 4. Lỗi thừa từ

- Sai: *Anh ấy là người đồng hương cùng quê với tôi.* ("đồng hương" đã có nghĩa là cùng quê)
- Sửa: *Anh ấy là người đồng hương với tôi.*

## 5. Lỗi trật tự từ

- Sai: *Nhà trường đã tổ chức cho học sinh tham quan khu di tích lịch sử rất đông.*
- Sửa: *Nhà trường đã tổ chức cho **rất đông học sinh** tham quan khu di tích lịch sử.*
- Sai: *Chị Lan đã đọc xong cuốn sách của thư viện rất dày.*
- Sửa: *Chị Lan đã đọc xong cuốn sách **rất dày** của thư viện.*

## 6. Cách phát hiện và sửa lỗi

1. Đọc lại câu, tự hỏi: từ này có đúng nghĩa không, có thừa không, có lặp không?
2. Tra từ điển khi phân vân về nghĩa của từ, nhất là từ Hán Việt.
3. Sửa bằng cách: bỏ từ thừa, thay từ đúng nghĩa, thay bằng đại từ hoặc từ đồng nghĩa, sắp xếp lại trật tự từ.

> [!TIP]
> Đọc to câu văn của mình lên: những chỗ nghe "vấp", "lủng củng" thường chính là nơi có lỗi lặp từ hoặc sai trật tự từ.`,
					},
					{
						Title:           "Viết bài văn nghị luận về một vấn đề xã hội",
						DurationMinutes: 30,
						Body: `## 1. Nghị luận về một vấn đề xã hội là gì?

Đây là kiểu bài yêu cầu người viết trình bày ý kiến về một **hiện tượng đời sống** hoặc một **tư tưởng, đạo lí**, dựa trên lí lẽ và bằng chứng thuyết phục.

## 2. Bố cục bài viết

1. **Mở bài**: Giới thiệu và nêu vấn đề cần bàn luận.
2. **Thân bài**:
   - Giải thích khái niệm, nội dung của vấn đề (nếu cần).
   - Phân tích, bàn luận các khía cạnh của vấn đề, có lí lẽ và bằng chứng cụ thể.
   - Xem xét vấn đề từ nhiều phía, phản biện những ý kiến trái chiều.
   - Rút ra bài học nhận thức và hành động cho bản thân.
3. **Kết bài**: Khẳng định lại ý kiến, để lại thông điệp.

## 3. Một số lưu ý khi viết

- Xác định đúng và bám sát yêu cầu của đề (vấn đề gì, phạm vi bằng chứng nào).
- Lập luận chặt chẽ: mỗi luận điểm có lí lẽ và bằng chứng đi kèm.
- Bằng chứng cụ thể, xác thực, tránh nêu chung chung.
- Diễn đạt rõ ràng, mạch lạc; có thể dùng câu hỏi tu từ để tăng sức thuyết phục.

## 4. Ví dụ đề bài

Đề: *Viết bài văn nghị luận trình bày suy nghĩ về ý nghĩa của tinh thần tự học đối với học sinh.*

Gợi ý triển khai:

- **Giải thích**: tự học là chủ động tìm tòi, tiếp thu và rèn luyện kiến thức, kĩ năng mà không cần ai thúc ép.
- **Bàn luận**: tự học giúp hiểu bài sâu, nhớ lâu; rèn tính kỉ luật và khả năng tư duy độc lập; là nền tảng để học suốt đời.
- **Phản biện**: phê phán thái độ ỷ lại, chỉ học khi bị nhắc nhở hoặc chép lời giải có sẵn; đồng thời lưu ý tự học không có nghĩa là không cần thầy cô, bạn bè.
- **Bài học**: lập thời gian biểu, đặt mục tiêu nhỏ, biết tìm nguồn tài liệu đáng tin cậy.

> [!TIP]
> Trước khi viết, hãy lập dàn ý nhanh ra giấy nháp theo cấu trúc ba phần để bài viết không thiếu ý hoặc lạc đề.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: thần thoại, sử thi, lỗi dùng từ và văn nghị luận xã hội",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Thần thoại 'Thần Trụ Trời' thuộc loại truyện nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Truyện kể về vị thần đội trời, đắp cột tách trời đất, qua đó giải thích nguồn gốc trời đất, núi đồi — đó là thần thoại suy nguyên.",
					Options: []sampleOption{
						{Content: "Truyện cười dân gian"},
						{Content: "Thần thoại suy nguyên, giải thích nguồn gốc trời đất, núi non", IsCorrect: true},
						{Content: "Truyện ngụ ngôn răn dạy con người"},
						{Content: "Sử thi anh hùng"},
					},
				},
				{
					Prompt: "Sử thi Đăm Săn là sử thi của dân tộc nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Đăm Săn là sử thi anh hùng của dân tộc Ê-đê ở Tây Nguyên.",
					Options: []sampleOption{
						{Content: "Mường"},
						{Content: "Ba Na"},
						{Content: "Gia Rai"},
						{Content: "Ê-đê", IsCorrect: true},
					},
				},
				{
					Prompt: "Trong đoạn trích 'Chiến thắng Mtao Mxây', chi tiết nào thể hiện sự giúp sức của thần linh đối với Đăm Săn?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Ông Trời hiện ra mách Đăm Săn lấy chày mòn ném vào vành tai Mtao Mxây — đây là sự trợ giúp của thần linh, một nét đặc trưng của sử thi.",
					Options: []sampleOption{
						{Content: "Ông Trời mách Đăm Săn lấy chày mòn ném vào vành tai Mtao Mxây", IsCorrect: true},
						{Content: "Hơ Nhị bị Mtao Mxây bắt đi"},
						{Content: "Dân làng mở tiệc mừng chiến thắng"},
						{Content: "Đăm Săn đến tận nhà Mtao Mxây thách đấu"},
					},
				},
				{
					Prompt: "Câu nào sau đây KHÔNG mắc lỗi dùng từ?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "'Tham quan' là cách dùng đúng. Các câu còn lại sai: 'bàng quang' phải là 'bàng quan'; 'yếu điểm' (điểm quan trọng) phải là 'điểm yếu'; 'đồng hương cùng quê' bị thừa từ.",
					Options: []sampleOption{
						{Content: "Chúng ta không được bàng quang trước nỗi đau của người khác."},
						{Content: "Bạn cần khắc phục những yếu điểm trong cách học."},
						{Content: "Cuối tuần, lớp em đi tham quan bảo tàng.", IsCorrect: true},
						{Content: "Anh ấy là người đồng hương cùng quê với tôi."},
					},
				},
				{
					Prompt: "Khi viết bài nghị luận về ý nghĩa của tinh thần tự học, ý nào thể hiện thao tác phản biện (xem xét mặt trái của vấn đề)?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Phản biện là nhìn vấn đề từ phía ngược lại; phê phán thái độ ỷ lại, chép lời giải sẵn chính là xem xét mặt trái của việc tự học. Các ý còn lại là giải thích, bàn luận mặt tích cực hoặc nêu bài học.",
					Options: []sampleOption{
						{Content: "Tự học là chủ động tìm tòi, tiếp thu kiến thức."},
						{Content: "Phê phán thái độ ỷ lại, chỉ học khi bị nhắc nhở hoặc chép lời giải có sẵn.", IsCorrect: true},
						{Content: "Tự học giúp rèn tính kỉ luật và khả năng tư duy độc lập."},
						{Content: "Em sẽ lập thời gian biểu tự học mỗi ngày."},
					},
				},
			},
		},
	},
	{
		Code:        "NGUVAN11",
		Title:       "Ngữ văn lớp 11",
		Description: "Lớp học mẫu môn Ngữ văn lớp 11 theo Chương trình GDPT 2018: đọc hiểu truyện ngắn hiện đại và thơ trữ tình, biện pháp tu từ đối – lặp cấu trúc và viết bài nghị luận phân tích, đánh giá tác phẩm thơ.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Đọc hiểu: truyện ngắn hiện đại và thơ trữ tình",
				Lessons: []sampleLesson{
					{
						Title:           "Người kể chuyện và điểm nhìn trong truyện Chí Phèo",
						DurationMinutes: 35,
						Body: `## 1. Người kể chuyện và điểm nhìn

- **Người kể chuyện** là nhân vật do nhà văn tạo ra để kể lại câu chuyện; có thể kể ở **ngôi thứ nhất** (xưng "tôi") hoặc **ngôi thứ ba** (giấu mình, gọi tên nhân vật).
- **Điểm nhìn** là vị trí, góc độ mà từ đó người kể quan sát, cảm nhận và kể lại sự việc. Điểm nhìn có thể ở **bên ngoài** nhân vật hoặc **dịch chuyển vào bên trong** nhân vật.
- **Lời nửa trực tiếp**: lời người kể chuyện hoà lẫn với giọng điệu, suy nghĩ của nhân vật.

## 2. Tác giả và tác phẩm

Nam Cao (1915–1951), tên khai sinh là Trần Hữu Tri, quê ở Lý Nhân, Hà Nam — nhà văn hiện thực xuất sắc viết về người nông dân nghèo và người trí thức tiểu tư sản. Truyện ngắn "Chí Phèo" (1941) lúc đầu có tên "Cái lò gạch cũ", khi in sách lần đầu bị nhà xuất bản đổi thành "Đôi lứa xứng đôi", về sau tác giả đặt lại là "Chí Phèo".

## 3. Tóm tắt

Chí Phèo là đứa trẻ bị bỏ rơi bên một lò gạch cũ, lớn lên đi ở cho nhiều nhà, rồi làm canh điền cho Bá Kiến. Vì một cơn ghen vô cớ, Bá Kiến đẩy Chí vào tù. Ra tù, Chí trở thành kẻ say rượu, rạch mặt ăn vạ, làm tay sai cho chính Bá Kiến — một "con quỷ dữ" của làng Vũ Đại. Gặp Thị Nở và được chăm sóc bằng bát cháo hành, Chí thức tỉnh, khao khát làm người lương thiện. Nhưng bà cô Thị Nở ngăn cản, Thị Nở từ chối Chí. Tuyệt vọng, Chí xách dao đến nhà Bá Kiến, đâm chết hắn rồi tự kết liễu đời mình.

## 4. Điểm nhìn trong truyện

- **Mở đầu**: điểm nhìn từ bên ngoài — cảnh Chí vừa đi vừa chửi, cả làng không ai thèm đáp lời — cho thấy Chí bị cộng đồng gạt bỏ.
- **Đoạn Chí tỉnh rượu**: điểm nhìn dịch chuyển vào bên trong nhân vật. Lần đầu tiên Chí nghe rõ những âm thanh bình dị của buổi sáng (tiếng chim, tiếng người đi chợ), nhớ lại ước mơ giản dị thuở trước và nhận ra mình đã già, cô độc.
- **Lời nửa trực tiếp** khiến người đọc vừa nghe giọng người kể, vừa nghe được tiếng lòng của Chí, nhờ đó thấy được phần người lương thiện còn ẩn sâu bên trong.

## 5. Giá trị

- **Hiện thực**: phản ánh tình trạng người nông dân lương thiện bị đẩy vào con đường tha hoá dưới chế độ thực dân nửa phong kiến.
- **Nhân đạo**: phát hiện và trân trọng bản chất lương thiện của con người; câu hỏi "Ai cho tao lương thiện?" là tiếng kêu đòi quyền làm người.

> [!NOTE]
> Khi đọc truyện, hãy trả lời hai câu hỏi tách biệt: "Ai kể?" và "Ai nhìn?". Người kể có thể là ngôi thứ ba nhưng điểm nhìn lại đặt vào chính nhân vật.`,
					},
					{
						Title:           "Thơ trữ tình: Câu cá mùa thu (Thu điếu) — Nguyễn Khuyến",
						DurationMinutes: 30,
						Body: `## 1. Văn bản

**Câu cá mùa thu** (Thu điếu) — Nguyễn Khuyến

> Ao thu lạnh lẽo nước trong veo,
> Một chiếc thuyền câu bé tẻo teo.
> Sóng biếc theo làn hơi gợn tí,
> Lá vàng trước gió khẽ đưa vèo.
> Tầng mây lơ lửng trời xanh ngắt,
> Ngõ trúc quanh co khách vắng teo.
> Tựa gối, ôm cần lâu chẳng được,
> Cá đâu đớp động dưới chân bèo.

## 2. Tác giả

Nguyễn Khuyến (1835–1909), hiệu Quế Sơn, quê làng Yên Đổ, huyện Bình Lục, tỉnh Hà Nam. Ông đỗ đầu cả ba kì thi Hương, Hội, Đình nên được gọi là **"Tam nguyên Yên Đổ"**. Làm quan một thời gian, ông cáo quan về quê sống thanh bạch. "Thu điếu" nằm trong chùm ba bài thơ thu viết bằng chữ Nôm: *Thu điếu*, *Thu vịnh*, *Thu ẩm*.

## 3. Các yếu tố của thơ trữ tình

- **Nhân vật trữ tình**: người ngồi câu cá — cũng chính là nhà thơ.
- **Cấu tứ**: điểm nhìn từ chiếc thuyền nhỏ trên ao, mở rộng lên bầu trời, theo ngõ trúc, rồi thu về với người ngồi câu.
- **Thể thơ**: thất ngôn bát cú Đường luật, gieo vần "eo" — một vần khó (tử vận) nhưng được dùng rất tự nhiên: *veo – teo – vèo – teo – bèo*.

## 4. Phân tích

- **Cảnh thu**: ao thu "lạnh lẽo", nước "trong veo", thuyền "bé tẻo teo" — không gian nhỏ, tĩnh, thanh sơ, đậm chất làng quê Bắc Bộ. Sóng "hơi gợn tí", lá "khẽ đưa vèo" — chuyển động rất nhẹ, càng làm nổi bật sự tĩnh lặng. Trời "xanh ngắt", ngõ trúc "quanh co", "khách vắng teo" — vắng lặng, quạnh hiu.
- **Tình thu**: hai câu kết cho thấy người câu cá không thực sự chú tâm vào việc câu. Âm thanh "cá đâu đớp động" nhỏ bé lại làm người đọc giật mình — nghệ thuật **lấy động tả tĩnh**. Ẩn sau cảnh là tâm trạng cô đơn, trầm ngâm và nỗi buồn thời thế của một nhà nho yêu nước.

> [!TIP]
> Hãy để ý cách vần "eo" khép dần không gian bài thơ: từ "trong veo", "tẻo teo" đến "vắng teo" — cảm giác hẹp lại, lặng đi, rất hợp với tâm trạng của nhân vật trữ tình.`,
					},
				},
			},
			{
				Title: "Chương 2. Tiếng Việt và kỹ năng viết",
				Lessons: []sampleLesson{
					{
						Title:           "Biện pháp tu từ đối và lặp cấu trúc",
						DurationMinutes: 25,
						Body: `## 1. Phép đối

**Đối** là cách sắp xếp từ ngữ, cụm từ, câu thành **từng cặp cân xứng** về âm thanh (bằng – trắc), từ loại và ý nghĩa (tương đồng hoặc tương phản).

Ví dụ trong "Thu điếu":

> Sóng biếc theo làn hơi gợn tí,
> Lá vàng trước gió khẽ đưa vèo.

- "sóng biếc" – "lá vàng"; "theo làn" – "trước gió"; "hơi gợn tí" – "khẽ đưa vèo".

Ví dụ trong tục ngữ: "Gần mực thì đen, gần đèn thì sáng" (mực – đèn, đen – sáng: đối tương phản).

**Tác dụng**: tạo nhịp điệu cân đối, hài hoà; làm nổi bật ý bằng sự tương đồng hoặc tương phản.

## 2. Lặp cấu trúc

**Lặp cấu trúc** là lặp lại một kiểu kết cấu cú pháp ở nhiều câu, nhiều vế liên tiếp.

Ví dụ trong ca dao:

> Thương thay thân phận con tằm,
> Kiếm ăn được mấy phải nằm nhả tơ.
> Thương thay lũ kiến li ti,
> Kiếm ăn được mấy phải đi tìm mồi.

Cấu trúc "Thương thay... / Kiếm ăn được mấy phải..." lặp lại, nhấn mạnh nỗi vất vả, thiệt thòi của người lao động nghèo.

**Tác dụng**: nhấn mạnh ý, tạo nhịp điệu, tăng sức biểu cảm và sức thuyết phục (rất hay gặp trong văn nghị luận, diễn thuyết).

## 3. Phân biệt

- **Đối** đòi hỏi các vế tương ứng **cân xứng** với nhau, không nhất thiết lặp từ.
- **Lặp cấu trúc** lặp lại **khuôn câu**, thường có từ ngữ được lặp lại.
- Hai biện pháp thường đi cùng nhau trong thơ ca, văn chính luận.

## 4. Bài tập mẫu có lời giải

Xác định biện pháp tu từ trong hai câu:

> Lom khom dưới núi tiều vài chú,
> Lác đác bên sông chợ mấy nhà.

**Lời giải**: Phép **đối** ("lom khom" – "lác đác", "dưới núi" – "bên sông", "tiều vài chú" – "chợ mấy nhà"), kết hợp **đảo ngữ** đưa từ tượng hình lên đầu câu, nhấn mạnh cảnh thưa thớt, hoang vắng.

> [!TIP]
> Khi phân tích phép đối, đừng chỉ liệt kê các cặp từ — hãy nói rõ chúng **tương đồng** hay **tương phản** và hiệu quả mà sự đối xứng ấy mang lại.`,
					},
					{
						Title:           "Viết văn bản nghị luận phân tích, đánh giá một tác phẩm thơ",
						DurationMinutes: 35,
						Body: `## 1. Yêu cầu

- Giới thiệu tác giả, bài thơ và nêu **nhận định khái quát**.
- **Phân tích** nội dung chủ đề và những nét đặc sắc nghệ thuật: hình ảnh, ngôn từ, vần, nhịp, biện pháp tu từ, cấu tứ.
- **Đánh giá** giá trị, vị trí của bài thơ: đóng góp gì, có gì độc đáo so với các bài thơ cùng đề tài.
- Có lí lẽ và bằng chứng (câu thơ, từ ngữ) cụ thể.

## 2. Quy trình

1. Đọc kĩ bài thơ, ghi lại những từ ngữ, hình ảnh gây ấn tượng.
2. Xác định luận đề (nhận định chung về bài thơ).
3. Lập dàn ý với các luận điểm rõ ràng.
4. Viết bài, kết hợp phân tích và đánh giá.
5. Đọc lại, sửa lỗi diễn đạt, kiểm tra trích dẫn.

## 3. Dàn ý mẫu

Đề: *Phân tích, đánh giá bài thơ "Câu cá mùa thu" của Nguyễn Khuyến.*

- **Mở bài**: Nguyễn Khuyến — nhà thơ của làng cảnh Việt Nam; "Câu cá mùa thu" là bức tranh thu đồng bằng Bắc Bộ đặc sắc.
- **Luận điểm 1**: Bức tranh thu tĩnh lặng, trong trẻo, thanh sơ (sáu câu đầu).
- **Luận điểm 2**: Tâm trạng nhân vật trữ tình: cô đơn, trầm ngâm, nỗi buồn thời thế (hai câu cuối).
- **Luận điểm 3**: Nghệ thuật: vần "eo" độc đáo, phép đối, bút pháp lấy động tả tĩnh, ngôn ngữ giản dị mà tinh tế.
- **Đánh giá**: bài thơ góp phần làm nên một hình ảnh mùa thu rất Việt Nam trong thơ Nôm trung đại.
- **Kết bài**: Khẳng định giá trị, nêu ấn tượng cá nhân.

## 4. Đoạn văn mẫu (luận điểm 2)

Hai câu kết hé lộ tâm trạng của người ngồi câu. Tư thế "tựa gối, ôm cần" kéo dài "lâu chẳng được" cho thấy ông không thật sự để tâm đến chuyện câu cá. Âm thanh "cá đâu đớp động dưới chân bèo" vốn rất khẽ, vậy mà đủ làm người câu giật mình — chứng tỏ cảnh vật tĩnh lặng đến tuyệt đối và tâm hồn nhà thơ đang chìm trong suy tư. Đó là nỗi buồn của một nhà nho thanh bạch trước thời thế đổi thay.

## 5. Lưu ý

- Tránh **diễn xuôi** bài thơ (chỉ kể lại ý từng câu).
- Luôn gắn phân tích nghệ thuật với nội dung: biện pháp ấy giúp thể hiện điều gì?
- Phần đánh giá nên có so sánh, liên hệ (với bài thơ khác cùng đề tài, với bối cảnh thời đại).

> [!NOTE]
> Phân tích trả lời câu hỏi "bài thơ viết gì và viết như thế nào", còn đánh giá trả lời câu hỏi "bài thơ có giá trị gì, hay ở đâu". Một bài viết tốt cần có cả hai.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: điểm nhìn trong truyện, thơ trữ tình, phép đối – lặp cấu trúc và nghị luận về thơ",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Vì sao Nguyễn Khuyến được gọi là 'Tam nguyên Yên Đổ'?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Nguyễn Khuyến đỗ đầu cả ba kì thi Hương, Hội, Đình (ba lần đỗ 'nguyên'), quê ở làng Yên Đổ nên được gọi là Tam nguyên Yên Đổ.",
					Options: []sampleOption{
						{Content: "Vì ông làm quan qua ba đời vua"},
						{Content: "Vì ông đỗ đầu cả ba kì thi Hương, Hội, Đình", IsCorrect: true},
						{Content: "Vì ông viết ba bài thơ thu nổi tiếng"},
						{Content: "Vì ông có ba người con đỗ đạt"},
					},
				},
				{
					Prompt: "Trong truyện kể, 'điểm nhìn' là gì?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Điểm nhìn là vị trí, góc độ mà từ đó người kể chuyện quan sát, cảm nhận và kể lại sự việc; điểm nhìn có thể ở bên ngoài hoặc bên trong nhân vật.",
					Options: []sampleOption{
						{Content: "Số lượng nhân vật xuất hiện trong truyện"},
						{Content: "Thời gian và địa điểm xảy ra câu chuyện"},
						{Content: "Bài học mà tác giả muốn gửi gắm"},
						{Content: "Vị trí, góc độ mà từ đó người kể quan sát, cảm nhận và kể lại sự việc", IsCorrect: true},
					},
				},
				{
					Prompt: "Hai câu thơ 'Sóng biếc theo làn hơi gợn tí, / Lá vàng trước gió khẽ đưa vèo.' sử dụng biện pháp tu từ nổi bật nào?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Các vế tương ứng cân xứng với nhau: 'sóng biếc' – 'lá vàng', 'theo làn' – 'trước gió', 'hơi gợn tí' – 'khẽ đưa vèo' — đó là phép đối.",
					Options: []sampleOption{
						{Content: "Phép đối", IsCorrect: true},
						{Content: "Nói quá"},
						{Content: "Hoán dụ"},
						{Content: "Nói giảm nói tránh"},
					},
				},
				{
					Prompt: "Câu hỏi 'Ai cho tao lương thiện?' của Chí Phèo thể hiện điều gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Đó là lúc Chí đã thức tỉnh, khao khát được làm người lương thiện nhưng bị cự tuyệt; câu hỏi vừa là tiếng kêu đau đớn, vừa là lời tố cáo xã hội đã cướp đi quyền làm người của Chí.",
					Options: []sampleOption{
						{Content: "Thói côn đồ, thích gây sự của Chí Phèo"},
						{Content: "Lời xin lỗi của Chí Phèo dành cho Thị Nở"},
						{Content: "Khát vọng làm người lương thiện và nỗi đau bị cự tuyệt quyền làm người", IsCorrect: true},
						{Content: "Sự hài lòng của Chí Phèo với cuộc sống hiện tại"},
					},
				},
				{
					Prompt: "Khi viết bài nghị luận về bài thơ 'Câu cá mùa thu', câu nào sau đây thể hiện yêu cầu ĐÁNH GIÁ chứ không chỉ mô tả, phân tích?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu nêu nhận xét về giá trị, vị trí của bài thơ ('một trong những bức tranh thu tiêu biểu') là đánh giá; các câu còn lại chỉ mô tả hình thức hoặc nội dung.",
					Options: []sampleOption{
						{Content: "Bài thơ có tám câu, mỗi câu bảy chữ."},
						{Content: "Với vần 'eo' độc đáo và bút pháp lấy động tả tĩnh, bài thơ là một trong những bức tranh thu tiêu biểu của thơ Nôm trung đại.", IsCorrect: true},
						{Content: "Câu thơ đầu tả ao thu lạnh lẽo, nước trong veo."},
						{Content: "Nhà thơ ngồi câu cá trên một chiếc thuyền nhỏ."},
					},
				},
			},
		},
	},
	{
		Code:        "NGUVAN12",
		Title:       "Ngữ văn lớp 12",
		Description: "Lớp học mẫu môn Ngữ văn lớp 12 theo Chương trình GDPT 2018: đọc hiểu tiểu thuyết hiện đại và văn bản chính luận, biện pháp nói mỉa – nghịch ngữ và viết bài nghị luận về vấn đề tuổi trẻ.",
		Cover:       "literature",
		Chapters: []sampleChapter{
			{
				Title: "Chương 1. Đọc hiểu: tiểu thuyết hiện đại và văn bản chính luận",
				Lessons: []sampleLesson{
					{
						Title:           "Tiểu thuyết Số đỏ và nghệ thuật trào phúng",
						DurationMinutes: 35,
						Body: `## 1. Tiểu thuyết và văn trào phúng

- **Tiểu thuyết** là tác phẩm tự sự cỡ lớn, phản ánh đời sống trên phạm vi rộng, có nhiều nhân vật, nhiều tuyến sự kiện đan xen.
- **Trào phúng** là cách dùng tiếng cười để châm biếm, đả kích những cái xấu, cái lố bịch trong xã hội.

## 2. Tác giả và tác phẩm

Vũ Trọng Phụng (1912–1939), quê gốc ở Hưng Yên, sinh ra và lớn lên ở Hà Nội, được mệnh danh là "ông vua phóng sự đất Bắc". Dù mất sớm, ông để lại nhiều tiểu thuyết, phóng sự có giá trị. Tiểu thuyết *Số đỏ* đăng báo năm 1936, là đỉnh cao của văn xuôi trào phúng Việt Nam trước Cách mạng tháng Tám.

## 3. Tóm tắt nội dung chính

Xuân Tóc Đỏ là một thanh niên mồ côi, lêu lổng, từng làm đủ nghề vặt như bán thuốc rong, nhặt bóng ở sân quần vợt. Nhờ sự láu lỉnh và những cơ may liên tiếp, Xuân lọt vào gia đình cụ cố Hồng, bà Phó Đoan và giới thượng lưu thành thị đang chạy theo phong trào "Âu hoá". Từ một kẻ vô học, Xuân được tung hô là "đốc tờ", "giáo sư quần vợt", "nhà cải cách xã hội". Cuối truyện, việc Xuân chịu thua đấu thủ nước ngoài trong một trận quần vợt — để "giữ hoà khí" — lại được ca ngợi như một hành động "cứu quốc", biến hắn thành "vĩ nhân".

## 4. Nghệ thuật trào phúng

- **Tình huống nghịch lí**: kẻ vô học thành "nhân tài", thua trận thành "cứu quốc".
- **Cường điệu, phóng đại** những thói lố lăng của xã hội thượng lưu.
- **Giọng văn mỉa mai**, lạnh lùng: người kể như đang khen mà thực chất là châm biếm.
- **Nhân vật biếm hoạ**: mỗi nhân vật được khắc hoạ bằng vài nét điển hình, gây cười.

Đoạn trích nổi tiếng "Hạnh phúc của một tang gia" cho thấy rõ điều đó: cái chết của cụ cố Tổ không đem lại nỗi buồn mà lại khiến con cháu "hạnh phúc" vì chúc thư sắp được thực hiện; đám tang biến thành dịp khoe khoang, phô trương.

## 5. Giá trị

*Số đỏ* phê phán sâu sắc xã hội thành thị tư sản đương thời: giả dối, lố lăng, chạy theo lối sống "văn minh" rởm, trong đó những kẻ cơ hội dễ dàng leo cao.

> [!NOTE]
> Cái cười trong *Số đỏ* là "cái cười lật mặt": càng tung hô Xuân bao nhiêu, nhà văn càng vạch trần sự rỗng tuếch của xã hội đã tung hô hắn bấy nhiêu.`,
					},
					{
						Title:           "Văn bản chính luận: Tuyên ngôn Độc lập — Hồ Chí Minh",
						DurationMinutes: 35,
						Body: `## 1. Văn bản chính luận

Văn bản chính luận bàn về các vấn đề chính trị, xã hội, nhằm tác động đến nhận thức và hành động của người đọc. Đặc điểm: **lập luận chặt chẽ, lí lẽ sắc bén, bằng chứng xác thực, giọng điệu đanh thép** nhưng giàu cảm xúc.

## 2. Hoàn cảnh ra đời

- Ngày 19/8/1945, Cách mạng tháng Tám thành công ở Hà Nội.
- Ngày 26/8/1945, Chủ tịch Hồ Chí Minh từ chiến khu Việt Bắc về Hà Nội; tại căn nhà số 48 Hàng Ngang, Người soạn thảo bản Tuyên ngôn.
- Ngày 2/9/1945, tại Quảng trường Ba Đình, Người đọc bản Tuyên ngôn Độc lập, khai sinh nước Việt Nam Dân chủ Cộng hoà.

## 3. Đối tượng và mục đích

- **Đối tượng**: đồng bào cả nước, nhân dân thế giới, và các thế lực đang âm mưu quay lại xâm lược.
- **Mục đích**: tuyên bố nền độc lập; bác bỏ luận điệu "khai hoá", "bảo hộ" mà thực dân dùng để biện minh cho việc trở lại Đông Dương.

## 4. Kết cấu lập luận

1. **Cơ sở pháp lí**: mở đầu bằng việc trích dẫn Tuyên ngôn Độc lập năm 1776 của Mỹ và Tuyên ngôn Nhân quyền và Dân quyền của Cách mạng Pháp, rồi "suy rộng ra" quyền bình đẳng, tự do của các dân tộc. Dùng chính lí lẽ của đối phương để bác bỏ họ — kiểu lập luận "gậy ông đập lưng ông".
2. **Cơ sở thực tế**: tố cáo tội ác của thực dân Pháp hơn 80 năm về chính trị, kinh tế; chỉ ra rằng chúng đã hai lần "bán" nước ta cho Nhật, và nhân dân ta đã giành chính quyền **từ tay Nhật** chứ không phải từ tay Pháp.
3. **Lời tuyên bố**: khẳng định quyền độc lập và quyết tâm bảo vệ độc lập. Câu then chốt: "Nước Việt Nam có quyền hưởng tự do và độc lập, và sự thật đã thành một nước tự do độc lập."

## 5. Nghệ thuật

- Lập luận chặt chẽ: từ lí lẽ chung (pháp lí) đến bằng chứng cụ thể (thực tế), rồi đi đến kết luận tất yếu.
- Bằng chứng xác thực, giàu sức tố cáo.
- Ngôn ngữ chính xác, hùng hồn; giọng văn vừa đanh thép vừa thấm đẫm tình cảm.

> [!TIP]
> Có thể học theo kết cấu của Tuyên ngôn khi viết bài nghị luận: nêu cơ sở lí lẽ → đưa bằng chứng thực tế → rút ra kết luận. Mỗi bước đều chuẩn bị cho bước sau.`,
					},
				},
			},
			{
				Title: "Chương 2. Tiếng Việt và kỹ năng viết",
				Lessons: []sampleLesson{
					{
						Title:           "Biện pháp tu từ nói mỉa và nghịch ngữ",
						DurationMinutes: 25,
						Body: `## 1. Nói mỉa

**Nói mỉa** là dùng lời lẽ có vẻ khen ngợi, đồng tình nhưng thực chất là để **chê bai, châm biếm**. Muốn nhận ra nói mỉa, cần dựa vào **ngữ cảnh** và **ngữ điệu**.

Ví dụ:

- "Giỏi thật! Một tuần đi học muộn bốn buổi." → bề ngoài là khen, thực chất là chê.
- Trong *Số đỏ*, những danh xưng như "đốc tờ", "giáo sư", "vĩ nhân cứu quốc" gán cho Xuân Tóc Đỏ là cách nói mỉa của nhà văn.

## 2. Nghịch ngữ

**Nghịch ngữ** là cách kết hợp những từ ngữ có ý nghĩa **trái ngược, tưởng như loại trừ nhau**, nhằm diễn tả một sự thật phức tạp và gây ấn tượng mạnh.

Ví dụ:

- Nhan đề "Hạnh phúc của một tang gia": "hạnh phúc" đặt cạnh "tang gia" (nhà có tang) — vạch trần sự giả dối, bất hiếu của những đứa con, đứa cháu.
- "im lặng đến chói tai", "một nụ cười đầy nước mắt".

## 3. Tác dụng trong văn chương

- Tạo tiếng cười châm biếm (nhất là trong văn trào phúng).
- Làm nổi bật mâu thuẫn, bản chất thật của sự việc.
- Gây bất ngờ, buộc người đọc phải suy ngẫm.

## 4. Bài tập mẫu có lời giải

Xác định biện pháp tu từ:

1. "Anh ấy hào phóng thật, rủ cả nhóm đi ăn rồi để bạn trả tiền." → **Nói mỉa** (khen "hào phóng" nhưng ý là chê keo kiệt).
2. "Đó là một chiến thắng cay đắng." → **Nghịch ngữ** ("chiến thắng" vốn vui, "cay đắng" vốn buồn).
3. "Mồ hôi thánh thót như mưa." → Không phải hai biện pháp trên; đây là **so sánh** kết hợp **nói quá**.

## 5. Lưu ý khi giao tiếp

Nói mỉa có thể gây tổn thương nếu dùng không đúng lúc, đúng người. Trong giao tiếp hằng ngày, nên góp ý thẳng thắn, tôn trọng thay vì mỉa mai.

> [!NOTE]
> Nói mỉa nằm ở **ý định** của người nói (nói ngược với điều mình nghĩ), còn nghịch ngữ nằm ở **sự kết hợp từ ngữ** trái nghĩa ngay trong câu chữ.`,
					},
					{
						Title:           "Viết bài văn nghị luận về một vấn đề liên quan đến tuổi trẻ",
						DurationMinutes: 35,
						Body: `## 1. Đặc điểm của kiểu bài

- Bàn về vấn đề **gần gũi với người trẻ**: ước mơ, lựa chọn nghề nghiệp, trách nhiệm với cộng đồng, lối sống trên mạng xã hội...
- Người viết phải có **quan điểm riêng**, rõ ràng.
- Kết hợp các thao tác: **giải thích, phân tích, chứng minh, bác bỏ**.

## 2. Bố cục

1. **Mở bài**: Nêu vấn đề và quan điểm của người viết.
2. **Thân bài**: Giải thích vấn đề; trình bày các luận điểm với lí lẽ và bằng chứng; bác bỏ ý kiến sai; đề xuất giải pháp.
3. **Kết bài**: Khẳng định quan điểm, liên hệ bản thân.

## 3. Ví dụ dàn ý

Đề: *Lựa chọn nghề nghiệp: theo đam mê hay theo xu hướng?*

- **Quan điểm**: chọn nghề cần dung hoà giữa **đam mê, năng lực bản thân và nhu cầu xã hội**.
- **Luận điểm 1**: Chọn nghề chỉ theo xu hướng mà thiếu năng lực, hứng thú thì dễ chán nản, bỏ dở.
- **Luận điểm 2**: Chọn nghề chỉ theo đam mê mà không tìm hiểu cơ hội việc làm, yêu cầu của nghề thì dễ gặp khó khăn khi lập nghiệp.
- **Luận điểm 3**: Hiểu rõ bản thân (sở thích, điểm mạnh) và tìm hiểu thị trường lao động giúp lựa chọn bền vững.
- **Bác bỏ**: ý kiến "nghề nào đang thịnh thì chọn nghề đó".
- **Giải pháp**: tham gia hoạt động hướng nghiệp, trải nghiệm thực tế, trò chuyện với người làm nghề.

## 4. Thao tác bác bỏ

Các bước: nêu ý kiến cần bác bỏ → chỉ ra chỗ sai bằng lí lẽ và bằng chứng → khẳng định ý kiến đúng.

**Đoạn văn mẫu**: Có người cho rằng cứ chọn nghề đang được nhiều người theo đuổi là chắc chắn thành công. Thực ra, một nghề "thịnh" hôm nay chưa chắc còn hấp dẫn sau vài năm, và khi quá nhiều người cùng đổ xô vào, sự cạnh tranh càng gay gắt. Nếu thiếu năng lực và hứng thú, người học dễ mệt mỏi, bỏ dở giữa chừng. Vì vậy, xu hướng xã hội chỉ nên là một yếu tố tham khảo, không thể thay thế cho sự hiểu biết về chính mình.

## 5. Lưu ý

- Tránh giọng răn dạy chung chung; hãy dùng bằng chứng gần gũi với đời sống học sinh.
- Khi dùng số liệu, cần nêu nguồn đáng tin cậy.

> [!TIP]
> Một bài nghị luận thuyết phục không né tránh ý kiến trái chiều. Dành một đoạn để bác bỏ hoặc thừa nhận có giới hạn quan điểm khác sẽ làm lập luận của bạn chặt chẽ hơn.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Số đỏ, Tuyên ngôn Độc lập, nói mỉa – nghịch ngữ và nghị luận về tuổi trẻ",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "Tác giả của tiểu thuyết 'Số đỏ' là ai?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Số đỏ là tiểu thuyết trào phúng của Vũ Trọng Phụng, đăng báo năm 1936.",
					Options: []sampleOption{
						{Content: "Nam Cao"},
						{Content: "Ngô Tất Tố"},
						{Content: "Vũ Trọng Phụng", IsCorrect: true},
						{Content: "Nguyễn Công Hoan"},
					},
				},
				{
					Prompt: "Chủ tịch Hồ Chí Minh đọc bản Tuyên ngôn Độc lập tại Quảng trường Ba Đình vào ngày nào?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Ngày 2/9/1945, tại Quảng trường Ba Đình, Chủ tịch Hồ Chí Minh đọc bản Tuyên ngôn Độc lập khai sinh nước Việt Nam Dân chủ Cộng hoà. Ngày 19/8/1945 là ngày khởi nghĩa giành chính quyền thành công ở Hà Nội.",
					Options: []sampleOption{
						{Content: "2/9/1945", IsCorrect: true},
						{Content: "19/8/1945"},
						{Content: "26/8/1945"},
						{Content: "7/5/1954"},
					},
				},
				{
					Prompt: "Việc mở đầu Tuyên ngôn Độc lập bằng lời trích tuyên ngôn của Mỹ và Pháp có tác dụng chủ yếu gì?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Lời trích tạo cơ sở pháp lí được thế giới thừa nhận, đồng thời dùng chính lí lẽ của các nước ấy để bác bỏ âm mưu xâm lược — kiểu lập luận 'gậy ông đập lưng ông'.",
					Options: []sampleOption{
						{Content: "Kể lại lịch sử hai nước Mỹ và Pháp"},
						{Content: "Thể hiện sự phụ thuộc vào các nước lớn"},
						{Content: "Giới thiệu văn hoá phương Tây cho nhân dân"},
						{Content: "Tạo cơ sở pháp lí vững chắc và dùng chính lí lẽ của đối phương để bác bỏ âm mưu xâm lược", IsCorrect: true},
					},
				},
				{
					Prompt: "Nhan đề 'Hạnh phúc của một tang gia' sử dụng biện pháp tu từ nào?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "'Hạnh phúc' và 'tang gia' (nhà có tang) là hai khái niệm trái ngược được đặt cạnh nhau — đó là nghịch ngữ, nhằm vạch trần sự giả dối của gia đình có tang.",
					Options: []sampleOption{
						{Content: "Nói giảm nói tránh"},
						{Content: "Nghịch ngữ", IsCorrect: true},
						{Content: "Điệp từ"},
						{Content: "So sánh"},
					},
				},
				{
					Prompt: "Trong bài nghị luận về chủ đề lựa chọn nghề nghiệp, đoạn nào sau đây sử dụng thao tác BÁC BỎ?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Bác bỏ là nêu một ý kiến rồi chỉ ra chỗ sai của nó bằng lí lẽ. Đoạn nêu quan niệm 'chọn nghề đang thịnh là chắc thành công' rồi phân tích vì sao chưa đúng chính là thao tác bác bỏ.",
					Options: []sampleOption{
						{Content: "Nghề nghiệp là công việc mỗi người gắn bó lâu dài để tạo ra thu nhập."},
						{Content: "Nhiều bạn trẻ hiện nay rất quan tâm đến việc chọn nghề."},
						{Content: "Có người cho rằng cứ chọn nghề đang thịnh là chắc thành công; thực ra nếu thiếu năng lực và hứng thú, người học dễ chán nản, bỏ dở giữa chừng.", IsCorrect: true},
						{Content: "Em mơ ước trở thành bác sĩ."},
					},
				},
			},
		},
	},
}
