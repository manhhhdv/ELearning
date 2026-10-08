package store

// sampleAnhSecondary là các lớp học mẫu môn Tiếng Anh từ lớp 7 đến lớp 12
// theo Chương trình GDPT 2018: mỗi lớp 2 Unit, mỗi Unit 2 bài giảng và một bài ôn tập.
var sampleAnhSecondary = []sampleCourse{
	{
		Code:        "TIENGANH7",
		Title:       "Tiếng Anh lớp 7",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 7 theo Chương trình GDPT 2018: từ vựng về sở thích và lối sống lành mạnh, thì hiện tại đơn, câu đơn và lời khuyên với should/shouldn't.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Unit 1. Hobbies",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng về sở thích và thì hiện tại đơn",
						DurationMinutes: 30,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| hobby | n | sở thích |
| collect stamps | v phr | sưu tầm tem |
| gardening | n | làm vườn |
| horse riding | n | cưỡi ngựa |
| build dollhouses | v phr | làm nhà búp bê |
| make models | v phr | làm mô hình |
| unusual | adj | khác thường, độc đáo |
| patient | adj | kiên nhẫn |

## 2. Thì hiện tại đơn (Present simple)

Dùng để nói về **thói quen, sở thích** và **sự thật hiển nhiên**.

- Khẳng định: I / You / We / They **play**; He / She / It **plays**.
- Phủ định: She **doesn't (does not) play**; They **don't play**.
- Nghi vấn: **Does** he **play**...? – Yes, he does. / No, he doesn't.

Quy tắc thêm -s/-es với ngôi thứ ba số ít: watch → watches, go → goes, study → studies, play → plays.

## 3. Trạng từ chỉ tần suất

always > usually > often > sometimes > rarely > never. Trạng từ đứng **trước động từ thường** và **sau động từ to be**.

- My sister **often goes** horse riding at weekends.
- He **is always** patient when he makes models.

## 4. Bài tập mẫu

Chia động từ: "My father (collect) old coins. He (not like) football."

**Lời giải:** My father **collects** old coins. He **doesn't like** football.

> [!TIP]
> Nhớ thêm -s/-es cho động từ khi chủ ngữ là he, she, it hoặc một danh từ số ít như my brother, Lan.`,
					},
					{
						Title:           "Đọc – viết: giới thiệu sở thích của em",
						DurationMinutes: 30,
						Body: `## 1. Đọc hiểu

Đọc đoạn văn sau:

*My name is Minh. My hobby is gardening. I started it two years ago with my grandmother. Every afternoon, I water the plants and pull up the weeds. Gardening is not difficult, but you need to be patient. It makes me feel relaxed, and I can give fresh vegetables to my family.*

**Câu hỏi:**

1. What is Minh's hobby? → *His hobby is gardening.*
2. Who did he start it with? → *He started it with his grandmother.*
3. Why does he like it? → *Because it makes him feel relaxed and he can give fresh vegetables to his family.*

## 2. Dàn ý đoạn văn giới thiệu sở thích

- **Câu chủ đề:** My hobby is... / I enjoy... in my free time.
- **Khi nào bắt đầu, với ai:** I started it when I was... / with...
- **Làm như thế nào, cần gì:** I usually... / You need...
- **Lợi ích, cảm nhận:** It helps me... / It makes me feel...
- **Kế hoạch tương lai:** In the future, I will...

## 3. Mẫu câu hữu ích

| Mục đích | Mẫu câu |
|---|---|
| Nêu sở thích | I love / like / enjoy + V-ing |
| Nêu lợi ích | It helps me (to) + V; It is good for... |
| Nêu cảm xúc | It makes me feel happy / relaxed |

## 4. Bài viết mẫu (khoảng 60 từ)

*I enjoy making models in my free time. I started this hobby when I was ten. I usually use paper, glue and old boxes to make houses and cars. It takes a lot of time, so I have to be careful and patient. Making models helps me to be creative. In the future, I want to make a model of my school.*

> [!NOTE]
> Một đoạn văn tốt có câu chủ đề rõ ràng ở đầu, các câu triển khai ở giữa và một câu kết nêu cảm nghĩ hoặc kế hoạch.`,
					},
				},
			},
			{
				Title: "Unit 2. Healthy living",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng về sống khoẻ và câu đơn",
						DurationMinutes: 30,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| healthy | adj | khoẻ mạnh, lành mạnh |
| junk food | n | đồ ăn vặt kém dinh dưỡng |
| stay up late | v phr | thức khuya |
| do exercise | v phr | tập thể dục |
| sunburn | n | cháy nắng |
| acne | n | mụn trứng cá |
| suncream | n | kem chống nắng |
| tofu | n | đậu phụ |

## 2. Câu đơn (Simple sentence)

Câu đơn là câu chỉ có **một mệnh đề độc lập**, tức là có một chủ ngữ và một vị ngữ. Chủ ngữ hoặc động từ có thể được ghép bằng "and".

- I **eat** a lot of vegetables. (1 chủ ngữ, 1 động từ)
- **My brother and I** go jogging every morning. (chủ ngữ ghép)
- She **washes** her face and **goes** to bed early. (vị ngữ ghép)

## 3. Phân biệt với câu ghép

Câu có **hai mệnh đề** nối bằng and, but, so (có hai cặp chủ ngữ – động từ) **không phải** câu đơn:

- I like fruit, **but** my sister likes sweets. → câu ghép
- He stays up late, **so** he feels tired. → câu ghép

## 4. Bài tập mẫu

Câu nào là câu đơn? (a) Lan and Mai drink a lot of water. (b) I was tired, so I went to bed.

**Lời giải:** (a) là câu đơn vì chỉ có một mệnh đề với chủ ngữ ghép "Lan and Mai". Câu (b) có hai mệnh đề nối bằng "so" nên là câu ghép.

> [!TIP]
> Đếm số cặp "chủ ngữ + động từ chia thì": chỉ có một cặp (dù chủ ngữ hay động từ được ghép) thì đó là câu đơn.`,
					},
					{
						Title:           "Đọc – viết: lời khuyên để sống khoẻ",
						DurationMinutes: 30,
						Body: `## 1. Cấu trúc đưa lời khuyên

- **should + V (nguyên mẫu)**: nên làm gì. *You should eat more vegetables.*
- **shouldn't (should not) + V**: không nên làm gì. *You shouldn't stay up late.*
- Hỏi xin lời khuyên: *What should I do to keep fit?*

## 2. Đọc hiểu

*Nam is a healthy boy. He gets up at 6 a.m. and jogs for 30 minutes. He drinks a lot of water and eats fruit every day. He doesn't eat junk food, and he goes to bed before 10 p.m. He also wears a hat and suncream when he goes out in the sun.*

**Câu hỏi:**

1. How long does Nam jog every morning? → *For 30 minutes.*
2. What does he do to avoid sunburn? → *He wears a hat and suncream.*
3. Does he eat junk food? → *No, he doesn't.*

## 3. Lời khuyên cho từng vấn đề

| Vấn đề | Lời khuyên |
|---|---|
| I often feel tired. | You should go to bed early. |
| I have acne. | You should wash your face regularly and eat less fried food. |
| I'm putting on weight. | You shouldn't eat junk food. You should do more exercise. |

## 4. Bài viết mẫu: Tips for a healthy life

*There are some simple ways to stay healthy. First, you should eat a lot of fruit and vegetables. Second, you should do exercise for at least 30 minutes a day. Third, you shouldn't stay up late. Finally, you should drink enough water. Follow these tips, and you will feel better every day.*

> [!NOTE]
> Dùng các từ nối First, Second, Third, Finally để sắp xếp ý cho bài viết mạch lạc.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Sở thích và lối sống lành mạnh",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "My brother ___ the guitar every evening.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Câu nói về thói quen (every evening), chủ ngữ ngôi thứ ba số ít \"my brother\" nên dùng hiện tại đơn \"plays\".",
					Options: []sampleOption{
						{Content: "play"},
						{Content: "plays", IsCorrect: true},
						{Content: "playing"},
						{Content: "is play"},
					},
				},
				{
					Prompt: "Which hobby is about growing flowers and vegetables?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Gardening\" (làm vườn) là trồng và chăm sóc hoa, rau.",
					Options: []sampleOption{
						{Content: "dancing"},
						{Content: "cooking"},
						{Content: "gardening", IsCorrect: true},
						{Content: "skating"},
					},
				},
				{
					Prompt: "You ___ eat too much junk food. It's bad for your health.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Vế sau cho biết việc này có hại cho sức khoẻ nên đây là lời khuyên không nên làm: \"shouldn't\".",
					Options: []sampleOption{
						{Content: "should"},
						{Content: "must"},
						{Content: "can"},
						{Content: "shouldn't", IsCorrect: true},
					},
				},
				{
					Prompt: "Which of the following is a simple sentence?",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Câu này chỉ có một mệnh đề với chủ ngữ ghép \"My sister and I\"; các câu còn lại có hai mệnh đề nên không phải câu đơn.",
					Options: []sampleOption{
						{Content: "My sister and I go swimming every Sunday.", IsCorrect: true},
						{Content: "I like milk, but my brother likes juice."},
						{Content: "She stays up late, so she feels tired."},
						{Content: "I eat vegetables because they are healthy."},
					},
				},
				{
					Prompt: "Read: \"Nam gets up at 6 a.m. and jogs for 30 minutes. He drinks a lot of water and goes to bed early.\" Which habit is NOT mentioned?",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Đoạn văn nhắc đến dậy sớm, chạy bộ, uống nhiều nước và đi ngủ sớm; không nhắc đến việc ăn trái cây.",
					Options: []sampleOption{
						{Content: "jogging in the morning"},
						{Content: "eating fruit every day", IsCorrect: true},
						{Content: "drinking a lot of water"},
						{Content: "going to bed early"},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH8",
		Title:       "Tiếng Anh lớp 8",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 8 theo Chương trình GDPT 2018: hoạt động thời gian rảnh, cuộc sống ở nông thôn, động từ chỉ sự yêu thích + V-ing và so sánh hơn của trạng từ.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Unit 1. Leisure time",
				Lessons: []sampleLesson{
					{
						Title:           "Động từ chỉ sự yêu thích + V-ing",
						DurationMinutes: 30,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| leisure time | n | thời gian rảnh rỗi |
| hang out (with friends) | v phr | đi chơi (với bạn bè) |
| DIY (do-it-yourself) | n | tự làm đồ thủ công |
| surf the net | v phr | lướt mạng |
| message | v | nhắn tin |
| cruel | adj | độc ác, tàn nhẫn |
| addicted (to) | adj | nghiện |

## 2. Động từ chỉ sự yêu thích / không thích

| Mức độ | Động từ |
|---|---|
| Rất thích | adore, love |
| Thích | like, enjoy, fancy |
| Không phiền | don't mind |
| Không thích | dislike, hate, detest |

## 3. Quy tắc dùng

- **adore, enjoy, fancy, don't mind, dislike, detest** + **V-ing**:
  *I don't mind helping my mum with the cooking.*
- **like, love, hate, prefer** + **V-ing** hoặc **to V** (nghĩa gần như nhau):
  *She loves reading / to read comics.*

Lưu ý: KHÔNG nói *I enjoy to play*; phải nói *I enjoy playing*.

## 4. Bài tập mẫu

Sửa lỗi: "My brother detests to do the washing-up."

**Lời giải:** detests to do → **detests doing**, vì "detest" chỉ đi với V-ing.

> [!TIP]
> Học thuộc nhóm "chỉ đi với V-ing": adore, enjoy, fancy, don't mind, dislike, detest. Các động từ like, love, hate, prefer dùng được cả hai dạng.`,
					},
					{
						Title:           "Đọc – viết: em làm gì khi rảnh rỗi?",
						DurationMinutes: 30,
						Body: `## 1. Đọc hiểu

*Teenagers today have many ways to spend their leisure time. Some enjoy playing sports or hanging out with friends at the park. Others prefer doing DIY projects, such as making bracelets or painting old T-shirts. However, many teens spend too much time surfing the net and messaging. Experts say that young people should balance screen time with outdoor activities.*

**Câu hỏi:**

1. Name two DIY projects in the passage. → *Making bracelets and painting old T-shirts.*
2. What do experts advise? → *Young people should balance screen time with outdoor activities.*

## 2. Từ nối trong đoạn văn

- Thêm ý: **also, in addition, besides**
- Đối lập: **however, but**
- Đưa ví dụ: **such as, for example**
- Kết luận: **in short, in conclusion**

## 3. Dàn ý viết về hoạt động yêu thích

1. Giới thiệu hoạt động: *In my leisure time, I love...*
2. Lý do thứ nhất: *First, it helps me...*
3. Lý do thứ hai: *In addition, ...*
4. Kết luận: *In short, ... is a great way to...*

## 4. Bài viết mẫu

*In my leisure time, I love doing DIY. First, it helps me relax after hours of studying. I enjoy turning old things into useful ones, for example, a jar into a pencil holder. In addition, I can make nice gifts for my friends. In short, DIY is a creative and useful way to spend free time.*

> [!NOTE]
> Hãy dùng đa dạng động từ chỉ sự yêu thích (love, enjoy, fancy, don't mind) để bài viết sinh động hơn.`,
					},
				},
			},
			{
				Title: "Unit 2. Life in the countryside",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng nông thôn và so sánh hơn của trạng từ",
						DurationMinutes: 30,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| countryside | n | vùng nông thôn |
| harvest | v, n | thu hoạch; vụ mùa |
| herd | v | chăn (gia súc) |
| paddy field | n | ruộng lúa |
| fly a kite | v phr | thả diều |
| peaceful | adj | yên bình |
| hospitable | adj | hiếu khách |
| vast | adj | rộng lớn, bao la |

## 2. So sánh hơn của trạng từ

| Loại trạng từ | Cách tạo | Ví dụ |
|---|---|---|
| Ngắn (cùng dạng tính từ) | thêm -er | fast → faster, hard → harder, early → earlier |
| Dài (đuôi -ly) | more + adv | carefully → more carefully |
| Bất quy tắc | — | well → better, badly → worse |

Cấu trúc: **S + V + adv-er / more + adv + than + ...**

- *Farmers get up **earlier than** city people.*
- *People in the countryside live **more peacefully than** those in the city.*

## 3. So sánh không bằng

**S + V + not ... as + adv + as + ...**

- *My brother doesn't run **as fast as** me.* = *I run **faster than** my brother.*

## 4. Bài tập mẫu

Viết lại: "Lan sings more beautifully than Mai." → Mai doesn't sing as beautifully as Lan.

> [!TIP]
> Không có dạng "fastly" hay "more fast": "fast" vừa là tính từ vừa là trạng từ, so sánh hơn là "faster".`,
					},
					{
						Title:           "Đọc – viết: một ngày ở quê",
						DurationMinutes: 30,
						Body: `## 1. Đọc hiểu

*Last summer, Hoa visited her grandparents in a small village. Life there was very different from life in the city. People got up very early to work in the paddy fields. In the afternoon, children herded buffaloes and flew kites in the open fields. The villagers were hospitable and always smiled at her. Hoa found life there slower but more peaceful.*

**Câu hỏi:**

1. Where did Hoa spend last summer? → *In a small village with her grandparents.*
2. What did the children do in the afternoon? → *They herded buffaloes and flew kites.*
3. How did Hoa find life in the village? → *Slower but more peaceful.*

## 2. Ưu – nhược điểm của cuộc sống nông thôn

| Advantages | Disadvantages |
|---|---|
| fresh air, peaceful | fewer entertainment facilities |
| friendly, hospitable people | fewer job opportunities |
| cheap, fresh food | poor public transport |

## 3. Bài viết mẫu

*I prefer living in the countryside. First, the air there is much fresher than in the city. Second, people are friendly and help one another. In addition, life is less stressful, and I can enjoy nature every day. However, there are fewer modern facilities. Overall, I think the countryside is a great place to live.*

> [!NOTE]
> Khi viết đoạn văn nêu quan điểm, hãy nêu ý kiến ở câu đầu, đưa 2–3 lý do và có thể nhắc một hạn chế trước khi kết luận.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Thời gian rảnh và cuộc sống nông thôn",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "My grandmother doesn't mind ___ for the whole family.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Don't mind\" luôn đi với V-ing nên chọn \"cooking\".",
					Options: []sampleOption{
						{Content: "cook"},
						{Content: "to cook"},
						{Content: "cooking", IsCorrect: true},
						{Content: "cooked"},
					},
				},
				{
					Prompt: "In the afternoon, children in the village often fly ___ in the open fields.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Cụm từ \"fly a kite / fly kites\" nghĩa là thả diều, một hoạt động quen thuộc ở nông thôn.",
					Options: []sampleOption{
						{Content: "computers"},
						{Content: "buses"},
						{Content: "buffaloes"},
						{Content: "kites", IsCorrect: true},
					},
				},
				{
					Prompt: "Tom runs ___ than his brother.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "\"Fast\" là trạng từ ngắn nên so sánh hơn là \"faster\"; không có dạng \"fastly\" hay \"more fast\".",
					Options: []sampleOption{
						{Content: "fast"},
						{Content: "faster", IsCorrect: true},
						{Content: "more fast"},
						{Content: "fastly"},
					},
				},
				{
					Prompt: "Choose the correct sentence.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "\"Fancy\" đi với V-ing nên \"I fancy going...\" đúng. Các câu khác sai vì detest, enjoy, don't mind phải đi với V-ing.",
					Options: []sampleOption{
						{Content: "I fancy going to the cinema tonight.", IsCorrect: true},
						{Content: "She detests to do housework."},
						{Content: "He enjoys to play chess."},
						{Content: "They don't mind help us."},
					},
				},
				{
					Prompt: "Choose the sentence closest in meaning to: \"Lan sings more beautifully than Mai.\"",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Lan hát hay hơn Mai nghĩa là Mai hát không hay bằng Lan: \"Mai doesn't sing as beautifully as Lan.\"",
					Options: []sampleOption{
						{Content: "Mai sings more beautifully than Lan."},
						{Content: "Lan doesn't sing as beautifully as Mai."},
						{Content: "Mai doesn't sing as beautifully as Lan.", IsCorrect: true},
						{Content: "Lan and Mai sing equally beautifully."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH9",
		Title:       "Tiếng Anh lớp 9",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 9 theo Chương trình GDPT 2018: cộng đồng địa phương và cuộc sống đô thị, từ để hỏi + to V, cụm động từ và so sánh kép.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Unit 1. Local community",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng cộng đồng và từ để hỏi + to V",
						DurationMinutes: 30,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| local community | n | cộng đồng địa phương |
| community helper | n | người phục vụ cộng đồng |
| artisan | n | nghệ nhân, thợ thủ công |
| handicraft | n | đồ thủ công |
| craft village | n | làng nghề |
| firefighter | n | lính cứu hoả |
| delivery person | n | người giao hàng |
| facility | n | cơ sở vật chất, tiện ích |

## 2. Từ để hỏi + động từ nguyên mẫu có "to"

Cấu trúc **what / how / where / when / who / which + to V** dùng sau các động từ như know, ask, decide, wonder, show, tell... để nói về việc **nên làm gì / làm thế nào**.

- *I don't know **how to get** to the craft village.* (không biết đi đến làng nghề bằng cách nào)
- *Can you tell me **where to buy** handicrafts?*
- *She hasn't decided **what to make** for the festival.*

Lưu ý: KHÔNG dùng "why + to V".

## 3. Viết lại câu

*I don't know what I should say.* → *I don't know **what to say**.*

*Could you show me how I can use this tool?* → *Could you show me **how to use** this tool?*

## 4. Bài tập mẫu

Điền từ: "The tourists asked the artisan ___ to make a conical hat."

**Lời giải:** **how** – vì hỏi về cách làm nón lá.

> [!TIP]
> Cấu trúc "từ để hỏi + to V" thường thay thế cho mệnh đề có "should / can": what I should do = what to do.`,
					},
					{
						Title:           "Cụm động từ và đọc hiểu về làng nghề",
						DurationMinutes: 30,
						Body: `## 1. Cụm động từ (phrasal verbs) thường gặp

| Cụm động từ | Nghĩa | Ví dụ |
|---|---|---|
| pass down | truyền lại (đời sau) | The skills are passed down from parents to children. |
| set up | thành lập | They set up a pottery workshop in 2010. |
| take over | tiếp quản | He took over his father's business. |
| look after | chăm sóc | Volunteers look after the elderly. |
| close down | đóng cửa (vĩnh viễn) | The old cinema closed down last year. |
| turn down | từ chối | She turned down the job offer. |
| find out | tìm ra, phát hiện | Let's find out more about this craft. |

## 2. Đọc hiểu

*In many craft villages, traditional skills are passed down from generation to generation. Artisans make beautiful products such as pottery, silk and conical hats. These villages attract many tourists, who want to see how the products are made. However, some young people leave their villages to find jobs in big cities, so some crafts may disappear. Local authorities are trying to help artisans sell their products online.*

**Câu hỏi:**

1. How are traditional skills kept alive? → *They are passed down from generation to generation.*
2. Why may some crafts disappear? → *Because young people leave their villages to work in cities.*

## 3. Kỹ năng đọc: đoán nghĩa từ ngữ cảnh

Từ "attract" trong câu "These villages attract many tourists" có thể đoán là **thu hút** vì kết quả là có nhiều khách du lịch đến.

> [!NOTE]
> Cụm động từ có nghĩa khác hẳn động từ gốc, ví dụ "turn down" không phải "quay xuống" mà là "từ chối". Hãy học chúng kèm câu ví dụ.`,
					},
				},
			},
			{
				Title: "Unit 2. City life",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng đô thị và so sánh kép",
						DurationMinutes: 30,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| crowded | adj | đông đúc |
| traffic jam | n | tắc đường |
| public transport | n | phương tiện công cộng |
| skyscraper | n | nhà chọc trời |
| pollution | n | sự ô nhiễm |
| convenient | adj | tiện lợi |
| liveable | adj | đáng sống |
| noisy | adj | ồn ào |

## 2. So sánh kép dạng "càng... càng..."

**The + so sánh hơn + S + V, the + so sánh hơn + S + V**

- *The bigger the city is, the more crowded it becomes.* (Thành phố càng lớn thì càng đông đúc.)
- *The more cars there are, the more polluted the air gets.*
- *The earlier you leave, the less traffic you meet.*

## 3. So sánh kép dạng "ngày càng..."

- Tính từ ngắn: **adj-er and adj-er** → *The city is getting bigger and bigger.*
- Tính từ dài: **more and more + adj** → *Life is becoming more and more expensive.*

## 4. Bài tập mẫu

Viết lại: "As the traffic gets heavier, the air becomes more polluted."

**Lời giải:** *The heavier the traffic gets, the more polluted the air becomes.*

> [!TIP]
> Trong cấu trúc "The..., the...", phần so sánh luôn đứng ngay sau "the" ở đầu mỗi vế, sau đó mới đến chủ ngữ và động từ.`,
					},
					{
						Title:           "Đọc – viết: sống ở thành phố, ưu và nhược",
						DurationMinutes: 30,
						Body: `## 1. Đọc hiểu

*City life has many advantages. There are good schools, modern hospitals and lots of entertainment, such as cinemas and shopping malls. Public transport is also convenient. However, cities are getting more and more crowded. Traffic jams happen every day, and the air is polluted. The more people move to the city, the more pressure there is on housing and services.*

**Câu hỏi:**

1. Name two advantages of city life. → *Good schools and modern hospitals (also entertainment and convenient public transport).*
2. What problems do cities have? → *Crowds, traffic jams and air pollution.*

## 2. Dàn ý bài viết so sánh

| Phần | Nội dung |
|---|---|
| Mở đoạn | Giới thiệu: City life has both advantages and disadvantages. |
| Thân đoạn 1 | Advantages: On the one hand, ... |
| Thân đoạn 2 | Disadvantages: On the other hand, ... |
| Kết đoạn | Quan điểm cá nhân: In my opinion, ... |

## 3. Bài viết mẫu

*City life has both advantages and disadvantages. On the one hand, people in cities have better access to education, healthcare and entertainment. On the other hand, they have to deal with traffic jams, noise and pollution. In my opinion, the city is a good place to study and work, but we should use more public transport to make it more liveable.*

> [!NOTE]
> Cặp từ nối "On the one hand... On the other hand..." giúp trình bày hai mặt của vấn đề một cách cân bằng.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Cộng đồng địa phương và cuộc sống đô thị",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "A person who makes traditional products skillfully with their hands is called a(n) ___.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Artisan\" là nghệ nhân, thợ thủ công làm sản phẩm bằng tay.",
					Options: []sampleOption{
						{Content: "tourist"},
						{Content: "customer"},
						{Content: "officer"},
						{Content: "artisan", IsCorrect: true},
					},
				},
				{
					Prompt: "I don't know ___ to get to the craft village. Can you show me the way?",
					Level:  "Nhận biết", Points: 2,
					Explanation: "Hỏi về cách đi đến làng nghề nên dùng \"how to get\".",
					Options: []sampleOption{
						{Content: "how", IsCorrect: true},
						{Content: "what"},
						{Content: "who"},
						{Content: "why"},
					},
				},
				{
					Prompt: "Traditional skills in the village are ___ from parents to children.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "\"Pass down\" nghĩa là truyền lại cho thế hệ sau, phù hợp với \"from parents to children\".",
					Options: []sampleOption{
						{Content: "given up"},
						{Content: "turned down"},
						{Content: "passed down", IsCorrect: true},
						{Content: "closed down"},
					},
				},
				{
					Prompt: "The bigger the city is, ___.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "So sánh kép \"The + so sánh hơn + S + V, the + so sánh hơn + S + V\" nên vế sau là \"the more crowded it becomes\".",
					Options: []sampleOption{
						{Content: "it becomes more crowded"},
						{Content: "the more crowded it becomes", IsCorrect: true},
						{Content: "the most crowded it becomes"},
						{Content: "more crowded it becomes"},
					},
				},
				{
					Prompt: "Choose the sentence closest in meaning to: \"As the traffic gets heavier, the air becomes more polluted.\"",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu gốc diễn tả quan hệ \"càng... càng...\", viết lại bằng so sánh kép: \"The heavier the traffic gets, the more polluted the air becomes.\"",
					Options: []sampleOption{
						{Content: "The traffic gets heavier and the air becomes cleaner."},
						{Content: "The heavier the traffic gets, the less polluted the air becomes."},
						{Content: "The heavier the traffic gets, the more polluted the air becomes.", IsCorrect: true},
						{Content: "The more polluted the air is, the lighter the traffic gets."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH10",
		Title:       "Tiếng Anh lớp 10",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 10 theo Chương trình GDPT 2018: đời sống gia đình và con người với môi trường, phân biệt hiện tại đơn/hiện tại tiếp diễn, will/be going to và câu bị động.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Unit 1. Family life",
				Lessons: []sampleLesson{
					{
						Title:           "Việc nhà và hiện tại đơn – hiện tại tiếp diễn",
						DurationMinutes: 35,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| household chores | n | việc nhà |
| breadwinner | n | trụ cột kinh tế trong gia đình |
| homemaker | n | người nội trợ |
| do the laundry | v phr | giặt giũ |
| do the washing-up | v phr | rửa bát |
| put out the rubbish | v phr | đổ rác |
| heavy lifting | n | việc mang vác nặng |
| share | v | chia sẻ, cùng làm |

## 2. Hiện tại đơn và hiện tại tiếp diễn

| | Hiện tại đơn | Hiện tại tiếp diễn |
|---|---|---|
| Cấu trúc | S + V(s/es) | S + am/is/are + V-ing |
| Cách dùng | thói quen, sự thật, lịch trình | hành động đang diễn ra lúc nói hoặc quanh thời điểm nói |
| Dấu hiệu | always, usually, every day | now, at the moment, Look!, Listen! |

- *My father usually **does** the heavy lifting.*
- *Look! My mother **is hanging** the laundry in the yard.*

## 3. Động từ chỉ trạng thái

Các động từ như **know, like, love, want, need, believe, understand, belong** thường **không** dùng ở thì tiếp diễn.

- Đúng: *I **know** the answer now.* – Sai: *I am knowing the answer.*

## 4. Bài tập mẫu

Chia động từ: "My sister (cook) dinner now, but she usually (cook) at 6 p.m."

**Lời giải:** My sister **is cooking** dinner now, but she usually **cooks** at 6 p.m.

> [!TIP]
> Dấu hiệu "now, at the moment, Look!" gợi ý hiện tại tiếp diễn; "usually, often, every day" gợi ý hiện tại đơn.`,
					},
					{
						Title:           "Đọc – viết: lợi ích của việc chia sẻ việc nhà",
						DurationMinutes: 35,
						Body: `## 1. Đọc hiểu

*In many families today, both parents work, so sharing household chores is very important. When children help with the housework, they learn life skills such as cooking and cleaning. They also learn to be responsible and to appreciate their parents' hard work. Moreover, doing chores together helps family members spend more time with one another and strengthens family bonds.*

**Câu hỏi:**

1. Why is sharing chores important today? → *Because both parents work in many families.*
2. What do children learn from doing housework? → *Life skills, responsibility and appreciation of their parents' hard work.*

## 2. Cấu trúc đoạn văn nêu lợi ích

- **Câu chủ đề:** There are several benefits of children doing housework.
- **Lợi ích 1 + giải thích:** Firstly, ...
- **Lợi ích 2 + giải thích:** Secondly, ...
- **Lợi ích 3 + giải thích:** Finally, ...
- **Câu kết:** In conclusion, ...

## 3. Bài viết mẫu (khoảng 100 từ)

*There are several benefits of children doing housework. Firstly, they develop life skills. For example, they can cook simple meals and keep their rooms tidy. Secondly, they become more responsible because they have their own tasks to do. Finally, sharing chores strengthens family bonds, as everyone works together and talks more. In conclusion, parents should encourage their children to help with the housework.*

## 4. Từ vựng mở rộng

| Từ | Nghĩa |
|---|---|
| responsible | có trách nhiệm |
| appreciate | trân trọng |
| strengthen family bonds | gắn kết tình cảm gia đình |

> [!NOTE]
> Mỗi lợi ích nên có một câu giải thích hoặc ví dụ đi kèm để đoạn văn thuyết phục hơn.`,
					},
				},
			},
			{
				Title: "Unit 2. Humans and the environment",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng môi trường và will – be going to",
						DurationMinutes: 35,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| environment | n | môi trường |
| carbon footprint | n | lượng khí thải CO₂ do một người/hoạt động tạo ra |
| eco-friendly | adj | thân thiện với môi trường |
| energy-saving | adj | tiết kiệm năng lượng |
| recycle | v | tái chế |
| reduce | v | giảm |
| household appliance | n | thiết bị gia dụng |
| adopt | v | áp dụng, chấp nhận |

## 2. Will và be going to

| | will + V | be going to + V |
|---|---|---|
| Quyết định | quyết định ngay lúc nói | dự định đã có từ trước |
| Dự đoán | dựa trên ý kiến, niềm tin (I think, probably) | dựa trên bằng chứng ở hiện tại |

- *A: The bin is full. – B: I'**ll** take it out.* (quyết định tức thì)
- *We **are going to** plant trees in the schoolyard next week.* (kế hoạch)
- *I think people **will** use more electric vehicles in the future.* (ý kiến)
- *Look at those dark clouds! It**'s going to** rain.* (có bằng chứng)

## 3. Bài tập mẫu

Chọn đáp án: "I've bought some cloth bags. I (will / am going to) use them instead of plastic bags."

**Lời giải:** **am going to** – vì đã mua túi vải từ trước, đây là dự định có sẵn.

> [!TIP]
> Hỏi bản thân: "Có bằng chứng hoặc kế hoạch từ trước không?" Có thì dùng be going to; nếu chỉ là ý kiến hoặc quyết định tức thì thì dùng will.`,
					},
					{
						Title:           "Câu bị động và bài viết bảo vệ môi trường",
						DurationMinutes: 35,
						Body: `## 1. Câu bị động (Passive voice)

Dùng khi muốn nhấn mạnh **đối tượng chịu tác động** hơn là người thực hiện.

**S + be (chia thì) + V3/V-ed (+ by O)**

| Thì | Chủ động | Bị động |
|---|---|---|
| Hiện tại đơn | People recycle bottles. | Bottles are recycled. |
| Quá khứ đơn | They planted trees. | Trees were planted. |
| Tương lai đơn | We will reduce waste. | Waste will be reduced. |
| be going to | They are going to clean the beach. | The beach is going to be cleaned. |

Thường lược bỏ "by people, by them, by someone" vì không cần thiết.

## 2. Đọc hiểu

*Every year, millions of tonnes of plastic are thrown away. Much of it is not recycled and ends up in the ocean. In many schools, students are encouraged to bring their own water bottles. Old paper is collected and sent to recycling centres. Small actions like these can help reduce our carbon footprint.*

**Câu hỏi:** What happens to old paper in many schools? → *It is collected and sent to recycling centres.*

## 3. Bài viết mẫu

*There are many ways to protect the environment at home. First, we should turn off lights and household appliances when they are not used. Second, rubbish should be sorted so that paper and plastic can be recycled. Finally, we can grow plants to make the air fresher. If everyone takes small actions, our environment will be cleaner.*

> [!NOTE]
> Câu bị động giúp văn phong trang trọng, khách quan, rất hợp với các bài viết về môi trường và khoa học.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Đời sống gia đình và môi trường",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "In my family, my father is the ___; he earns money to support all of us.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Breadwinner\" là trụ cột kinh tế, người kiếm tiền nuôi gia đình.",
					Options: []sampleOption{
						{Content: "homemaker"},
						{Content: "breadwinner", IsCorrect: true},
						{Content: "housework"},
						{Content: "chore"},
					},
				},
				{
					Prompt: "Look! My mother ___ the laundry in the yard.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Look!\" báo hiệu hành động đang diễn ra nên dùng hiện tại tiếp diễn \"is hanging\".",
					Options: []sampleOption{
						{Content: "hangs"},
						{Content: "hanged"},
						{Content: "is hanging", IsCorrect: true},
						{Content: "hang"},
					},
				},
				{
					Prompt: "Look at those dark clouds! It ___ rain.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Dự đoán dựa trên bằng chứng ở hiện tại (mây đen) nên dùng \"is going to\".",
					Options: []sampleOption{
						{Content: "is going to", IsCorrect: true},
						{Content: "rains"},
						{Content: "is raining to"},
						{Content: "was going"},
					},
				},
				{
					Prompt: "I think people ___ more electric vehicles in the future.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Dự đoán tương lai dựa trên ý kiến cá nhân (I think) nên dùng \"will use\".",
					Options: []sampleOption{
						{Content: "are using"},
						{Content: "used"},
						{Content: "were using"},
						{Content: "will use", IsCorrect: true},
					},
				},
				{
					Prompt: "Choose the correct passive form of: \"People recycle plastic bottles in this factory.\"",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Câu chủ động ở hiện tại đơn, tân ngữ số nhiều nên bị động là \"are recycled\".",
					Options: []sampleOption{
						{Content: "Plastic bottles recycle in this factory."},
						{Content: "Plastic bottles are recycled in this factory.", IsCorrect: true},
						{Content: "Plastic bottles were recycled in this factory."},
						{Content: "Plastic bottles is recycling in this factory."},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH11",
		Title:       "Tiếng Anh lớp 11",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 11 theo Chương trình GDPT 2018: sống khoẻ và sống lâu, khoảng cách thế hệ, phân biệt quá khứ đơn/hiện tại hoàn thành và động từ khuyết thiếu must, have to, should.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Unit 1. A long and healthy life",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng sức khoẻ và quá khứ đơn – hiện tại hoàn thành",
						DurationMinutes: 35,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| life expectancy | n | tuổi thọ trung bình |
| longevity | n | sự sống lâu |
| balanced diet | n | chế độ ăn cân bằng |
| nutrition | n | dinh dưỡng |
| work out | v phr | tập luyện thể thao |
| boost the immune system | v phr | tăng cường hệ miễn dịch |
| fitness | n | sự khoẻ mạnh, thể lực |
| ingredient | n | thành phần, nguyên liệu |

## 2. Quá khứ đơn và hiện tại hoàn thành

| | Quá khứ đơn | Hiện tại hoàn thành |
|---|---|---|
| Cấu trúc | S + V2/V-ed | S + have/has + V3/V-ed |
| Cách dùng | hành động đã kết thúc tại thời điểm xác định trong quá khứ | hành động bắt đầu trong quá khứ, kéo dài đến hiện tại hoặc kết quả còn liên quan đến hiện tại |
| Dấu hiệu | yesterday, last year, ago, in 2020 | since, for, already, yet, ever, never, so far |

- *My grandfather **retired** in 2018.*
- *He **has gone** jogging every morning **since** he retired.*
- *I **have** never **tried** this healthy recipe before.*

## 3. Lỗi thường gặp

Không dùng hiện tại hoàn thành với mốc thời gian quá khứ xác định:

- Sai: *I have visited my grandparents last summer.*
- Đúng: *I **visited** my grandparents last summer.*

## 4. Bài tập mẫu

Chia động từ: "She (start) eating more vegetables two months ago, and she (lose) three kilos so far."

**Lời giải:** She **started** ... two months ago, and she **has lost** three kilos so far.

> [!TIP]
> Thấy "ago, last, yesterday, in + năm" thì dùng quá khứ đơn; thấy "since, for, so far, yet, already" thì nghĩ đến hiện tại hoàn thành.`,
					},
					{
						Title:           "Đọc – viết: bí quyết sống khoẻ, sống lâu",
						DurationMinutes: 35,
						Body: `## 1. Đọc hiểu

*Scientists have studied people in many parts of the world who live to 100 or more. They have found that these people share some common habits. They eat a balanced diet with lots of vegetables, beans and little meat. They stay active every day, often by walking or gardening rather than going to the gym. They also have strong relationships with family and friends, which reduces stress. Life expectancy depends not only on genes but also on lifestyle.*

**Câu hỏi:**

1. What kind of diet do these people have? → *A balanced diet with lots of vegetables, beans and little meat.*
2. How do they stay active? → *By walking or gardening every day.*
3. Besides genes, what affects life expectancy? → *Lifestyle.*

## 2. Kỹ năng đọc: tìm ý chính

Ý chính thường nằm ở **câu đầu** hoặc **câu cuối** đoạn. Ở đoạn trên, ý chính là: những người sống thọ có chung một số thói quen tốt, và lối sống ảnh hưởng đến tuổi thọ.

## 3. Bài viết mẫu: My healthy habits

*I have tried to live a healthier life since last year. First, I have cut down on sugary drinks and eaten more fruit. Second, I work out three times a week, which has improved my fitness. Finally, I go to bed before 11 p.m. to get enough sleep. Thanks to these habits, I have felt more energetic and have caught fewer colds.*

> [!NOTE]
> Khi kể về những thay đổi từ quá khứ đến nay, hiện tại hoàn thành là thì phù hợp nhất: I have tried..., I have felt...`,
					},
				},
			},
			{
				Title: "Unit 2. The generation gap",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng gia đình và động từ khuyết thiếu",
						DurationMinutes: 35,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| generation gap | n | khoảng cách thế hệ |
| nuclear family | n | gia đình hạt nhân (bố mẹ và con) |
| extended family | n | gia đình nhiều thế hệ |
| conflict | n | sự xung đột |
| viewpoint | n | quan điểm |
| privacy | n | sự riêng tư |
| respect | v, n | tôn trọng; sự tôn trọng |
| open-minded | adj | cởi mở |

## 2. Must, have to, should

| Động từ | Ý nghĩa | Ví dụ |
|---|---|---|
| must | bắt buộc do người nói cho là cần thiết, hoặc nội quy | You must be home before 10 p.m. |
| have to | bắt buộc do hoàn cảnh, quy định bên ngoài | I have to wear a uniform at school. |
| should | lời khuyên | You should listen to your grandparents. |

## 3. Mustn't và don't have to

Hai dạng phủ định này có nghĩa **rất khác nhau**:

- **mustn't** = cấm, không được phép: *You mustn't use your phone during the exam.*
- **don't have to** = không cần thiết: *Students don't have to wear uniforms on Saturdays.* (có thể mặc hoặc không)

Lưu ý: "must" không có dạng quá khứ; dùng **had to**: *Last year I had to share a room with my brother.*

## 4. Bài tập mẫu

Điền mustn't hoặc don't have to: "Tomorrow is Sunday, so you ___ get up early."

**Lời giải:** **don't have to** – vì không bắt buộc phải dậy sớm, chứ không phải bị cấm.

> [!TIP]
> Ghi nhớ: mustn't = "không được"; don't have to = "không cần".`,
					},
					{
						Title:           "Đọc – viết: thu hẹp khoảng cách thế hệ",
						DurationMinutes: 35,
						Body: `## 1. Đọc hiểu

*In an extended family, three generations often live under the same roof. This can lead to conflicts because each generation has different viewpoints. Grandparents may think teenagers spend too much time on their phones, while teenagers may feel that their grandparents do not respect their privacy. However, living together also has benefits: grandparents can share their experience, and young people can teach them to use new technology.*

**Câu hỏi:**

1. Why can conflicts happen in extended families? → *Because each generation has different viewpoints.*
2. What can young people teach their grandparents? → *How to use new technology.*

## 2. Từ nối chỉ sự tương phản

- **while**: trong khi (đối lập hai vế cùng câu)
- **however**: tuy nhiên (đầu câu, sau đó có dấu phẩy)
- **although + mệnh đề**: mặc dù

## 3. Dàn ý đoạn văn giải pháp

1. Nêu vấn đề: *The generation gap can cause conflicts in families.*
2. Giải pháp 1: *Firstly, family members should...*
3. Giải pháp 2: *Secondly, ...*
4. Kết luận: *If we..., our families will be happier.*

## 4. Bài viết mẫu

*The generation gap can cause conflicts in many families. Firstly, family members should spend time talking to one another to understand different viewpoints. Secondly, parents and grandparents should respect teenagers' privacy, while young people must show respect for their elders. Finally, families can enjoy activities together, such as cooking or travelling. If we are open-minded, our families will be happier.*

> [!NOTE]
> Khi viết giải pháp, dùng "should" để đưa lời khuyên và "must" khi muốn nhấn mạnh điều thật sự cần thiết.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Sống khoẻ và khoảng cách thế hệ",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "___ is the average number of years a person is expected to live.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Life expectancy\" là tuổi thọ trung bình.",
					Options: []sampleOption{
						{Content: "Fitness"},
						{Content: "Nutrition"},
						{Content: "Life expectancy", IsCorrect: true},
						{Content: "Ingredient"},
					},
				},
				{
					Prompt: "A family in which grandparents, parents and children live together is called a(n) ___ family.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Extended family\" là gia đình nhiều thế hệ cùng chung sống; \"nuclear family\" chỉ gồm bố mẹ và con.",
					Options: []sampleOption{
						{Content: "extended", IsCorrect: true},
						{Content: "nuclear"},
						{Content: "single"},
						{Content: "modern"},
					},
				},
				{
					Prompt: "My grandfather ___ jogging every morning since he retired.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "\"Since + mốc thời gian\" diễn tả hành động kéo dài đến hiện tại nên dùng hiện tại hoàn thành \"has gone\".",
					Options: []sampleOption{
						{Content: "went"},
						{Content: "goes"},
						{Content: "was going"},
						{Content: "has gone", IsCorrect: true},
					},
				},
				{
					Prompt: "Students ___ wear uniforms on Saturdays. They can wear casual clothes if they like.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Câu sau cho biết được tự do lựa chọn, tức là không cần thiết, nên dùng \"don't have to\"; \"mustn't\" mang nghĩa cấm.",
					Options: []sampleOption{
						{Content: "mustn't"},
						{Content: "don't have to", IsCorrect: true},
						{Content: "must"},
						{Content: "have to"},
					},
				},
				{
					Prompt: "Find the mistake: \"I have visited my grandparents in the countryside last summer.\"",
					Level:  "Vận dụng", Points: 2,
					Explanation: "\"Last summer\" là mốc thời gian quá khứ xác định nên phải dùng quá khứ đơn: \"visited\" thay cho \"have visited\".",
					Options: []sampleOption{
						{Content: "in"},
						{Content: "the countryside"},
						{Content: "have visited", IsCorrect: true},
						{Content: "last summer"},
					},
				},
			},
		},
	},
	{
		Code:        "TIENGANH12",
		Title:       "Tiếng Anh lớp 12",
		Description: "Lớp học mẫu môn Tiếng Anh lớp 12 theo Chương trình GDPT 2018: những câu chuyện cuộc đời đáng ngưỡng mộ, thế giới đa văn hoá, quá khứ đơn/quá khứ tiếp diễn và cách dùng mạo từ.",
		Cover:       "language",
		Chapters: []sampleChapter{
			{
				Title: "Unit 1. Life stories we admire",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng tiểu sử và quá khứ đơn – quá khứ tiếp diễn",
						DurationMinutes: 35,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| biography | n | tiểu sử (do người khác viết) |
| autobiography | n | tự truyện |
| achievement | n | thành tựu |
| admire | v | ngưỡng mộ |
| dedicate / devote (oneself) to | v | cống hiến cho |
| inspire | v | truyền cảm hứng |
| talented | adj | tài năng |
| determined | adj | quyết tâm |

## 2. Quá khứ đơn và quá khứ tiếp diễn

| | Quá khứ đơn | Quá khứ tiếp diễn |
|---|---|---|
| Cấu trúc | S + V2/V-ed | S + was/were + V-ing |
| Cách dùng | hành động đã hoàn tất; chuỗi hành động nối tiếp | hành động đang diễn ra tại một thời điểm trong quá khứ |

**Kết hợp hai thì:** hành động đang diễn ra (quá khứ tiếp diễn) thì một hành động khác xen vào (quá khứ đơn).

- *She **was studying** in Paris **when** she **met** her future husband.*
- ***While** I **was reading** her biography, my friend **called** me.*

Hai hành động song song: *While my mother **was cooking**, my father **was watching** TV.*

## 3. Bài tập mẫu

Nối câu: "He was walking in the park. He had a great idea."

**Lời giải:** *While he **was walking** in the park, he **had** a great idea.*

> [!TIP]
> "While" thường đi với quá khứ tiếp diễn (hành động kéo dài); "when" thường đi với quá khứ đơn (hành động xen vào, ngắn).`,
					},
					{
						Title:           "Đọc – viết: tiểu sử một nhà khoa học",
						DurationMinutes: 35,
						Body: `## 1. Đọc hiểu

*Marie Curie was born in Warsaw in 1867. In 1891, she moved to Paris to study physics and mathematics at the Sorbonne. While she was working in a laboratory, she met the physicist Pierre Curie, and they married in 1895. Together they studied radioactivity. In 1903, she became the first woman to win a Nobel Prize, which she shared with Pierre Curie and Henri Becquerel. In 1911, she won a second Nobel Prize, this time in Chemistry. She devoted her whole life to science and died in 1934.*

**Câu hỏi:**

1. Where did Marie Curie study? → *At the Sorbonne in Paris.*
2. What was special about her Nobel Prize in 1903? → *She was the first woman to win a Nobel Prize.*
3. In which fields did she win Nobel Prizes? → *Physics (1903) and Chemistry (1911).*

## 2. Cấu trúc một đoạn tiểu sử

| Phần | Nội dung | Mẫu câu |
|---|---|---|
| Mở đầu | năm sinh, nơi sinh | X was born in... in... |
| Học tập, sự nghiệp | các mốc quan trọng | In..., he/she studied / worked... |
| Thành tựu | giải thưởng, đóng góp | He/She is famous for... / won... |
| Kết thúc | ý nghĩa, lý do ngưỡng mộ | He/She has inspired... |

## 3. Bài viết mẫu (rút gọn)

*Marie Curie was a talented and determined scientist. Born in Warsaw in 1867, she later studied in Paris. She is famous for her research on radioactivity, which won her two Nobel Prizes. She is the only person to have won Nobel Prizes in two different sciences. Her dedication has inspired generations of young scientists, especially women.*

> [!NOTE]
> Khi viết tiểu sử, hãy sắp xếp các sự kiện theo trình tự thời gian và dùng quá khứ đơn cho các mốc đã qua.`,
					},
				},
			},
			{
				Title: "Unit 2. A multicultural world",
				Lessons: []sampleLesson{
					{
						Title:           "Từ vựng đa văn hoá và cách dùng mạo từ",
						DurationMinutes: 35,
						Body: `## 1. Từ vựng trọng tâm

| Từ / cụm từ | Loại từ | Nghĩa |
|---|---|---|
| multicultural | adj | đa văn hoá |
| cultural identity | n | bản sắc văn hoá |
| cultural diversity | n | sự đa dạng văn hoá |
| custom | n | phong tục |
| tradition | n | truyền thống |
| cuisine | n | ẩm thực |
| globalisation | n | toàn cầu hoá |
| heritage | n | di sản |

## 2. Mạo từ a / an

Dùng với danh từ đếm được số ít, nhắc đến lần đầu hoặc không xác định:

- **a** + âm phụ âm: *a festival, a university* (u đọc /juː/)
- **an** + âm nguyên âm: *an artist, an hour* (h câm)

## 3. Mạo từ the

- Vật đã được nhắc đến hoặc người nghe đã biết: *I bought a dress. **The** dress is from India.*
- Vật duy nhất: *the sun, the world*
- Trước tên một số quốc gia có từ chỉ thể chế, dạng số nhiều và các đại dương, sông: *the United Kingdom, the Philippines, the Pacific, the Mekong River*
- So sánh nhất: *the most famous dish*

## 4. Không dùng mạo từ (Ø)

- Danh từ số nhiều hoặc không đếm được nói chung: *Ø Music connects people.*
- Tên phần lớn quốc gia, thành phố: *Ø Viet Nam, Ø Japan, Ø Paris*
- Bữa ăn, môn học, ngôn ngữ: *have Ø lunch, study Ø English*

## 5. Bài tập mẫu

Điền mạo từ: "It took me ___ hour to cook ___ traditional dish that my grandmother taught me."

**Lời giải:** **an** hour (h câm) ... **the** traditional dish (đã được xác định bởi mệnh đề "that my grandmother taught me").

> [!TIP]
> Chọn a hay an theo **âm** đầu tiên chứ không theo chữ cái: an hour, a university, an MP3 player.`,
					},
					{
						Title:           "Đọc – viết: giữ gìn bản sắc trong thời toàn cầu hoá",
						DurationMinutes: 35,
						Body: `## 1. Đọc hiểu

*Thanks to globalisation, people today can enjoy food, music and films from all over the world. Cities have become more multicultural, and young people can easily learn about other cultures online. However, some people worry that local traditions may disappear as global trends spread. To protect their cultural identity, many communities organise traditional festivals, teach folk songs at school and promote local cuisine to tourists.*

**Câu hỏi:**

1. What is one benefit of globalisation mentioned in the text? → *People can enjoy food, music and films from all over the world.*
2. What do some people worry about? → *That local traditions may disappear.*
3. How do communities protect their cultural identity? → *By organising festivals, teaching folk songs and promoting local cuisine.*

## 2. Cấu trúc bài viết nêu quan điểm

1. **Mở bài:** nêu vấn đề và quan điểm (In my opinion, ...).
2. **Thân bài:** 2 lý do, mỗi lý do có ví dụ (For example, ...).
3. **Kết bài:** khẳng định lại quan điểm (In conclusion, ...).

## 3. Đoạn văn mẫu

*In my opinion, young people play an important role in preserving cultural identity. First, they can learn about traditional customs from their grandparents and share them on social media. For example, many students post videos about the Lunar New Year and folk games. Second, they can introduce Vietnamese cuisine and festivals to international friends. In conclusion, we can be global citizens while still keeping our own traditions.*

> [!NOTE]
> Sự đa dạng văn hoá không có nghĩa là đánh mất bản sắc: hội nhập tốt nhất là vừa học hỏi, vừa giữ gìn giá trị riêng.`,
					},
				},
			},
		},
		Quiz: sampleQuiz{
			Title:        "Ôn tập: Câu chuyện cuộc đời và thế giới đa văn hoá",
			Instructions: "Làm bài để tự kiểm tra kiến thức các bài vừa học. Có thể làm lại nhiều lần.",
			Questions: []sampleQuestion{
				{
					Prompt: "A ___ is the story of a person's life written by someone else.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Biography\" là tiểu sử do người khác viết; \"autobiography\" là tự truyện do chính người đó viết.",
					Options: []sampleOption{
						{Content: "autobiography"},
						{Content: "biography", IsCorrect: true},
						{Content: "novel"},
						{Content: "diary"},
					},
				},
				{
					Prompt: "It took me ___ hour to cook this traditional dish.",
					Level:  "Nhận biết", Points: 2,
					Explanation: "\"Hour\" có chữ h câm, bắt đầu bằng âm nguyên âm nên dùng \"an\".",
					Options: []sampleOption{
						{Content: "a"},
						{Content: "the"},
						{Content: "an", IsCorrect: true},
						{Content: "Ø (no article)"},
					},
				},
				{
					Prompt: "When my friend called me last night, I ___ a biography of Marie Curie.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "Hành động đang diễn ra trong quá khứ thì có hành động khác xen vào (called) nên dùng quá khứ tiếp diễn \"was reading\".",
					Options: []sampleOption{
						{Content: "am reading"},
						{Content: "was reading", IsCorrect: true},
						{Content: "have read"},
						{Content: "will read"},
					},
				},
				{
					Prompt: "Cultural ___ means the variety of different cultures existing together in a society.",
					Level:  "Thông hiểu", Points: 2,
					Explanation: "\"Cultural diversity\" là sự đa dạng văn hoá; \"cultural identity\" là bản sắc văn hoá của một cộng đồng.",
					Options: []sampleOption{
						{Content: "identity"},
						{Content: "heritage"},
						{Content: "tradition"},
						{Content: "diversity", IsCorrect: true},
					},
				},
				{
					Prompt: "Choose the best way to combine: \"She was walking home. She saw an accident.\"",
					Level:  "Vận dụng", Points: 2,
					Explanation: "Hành động kéo dài \"was walking\" dùng quá khứ tiếp diễn sau \"while\"; hành động xen vào \"saw\" dùng quá khứ đơn.",
					Options: []sampleOption{
						{Content: "While she was walking home, she saw an accident.", IsCorrect: true},
						{Content: "When she was walking home, she was seeing an accident."},
						{Content: "While she walked home, she was seeing an accident."},
						{Content: "When she saw an accident, she was walking home after."},
					},
				},
			},
		},
	},
}
