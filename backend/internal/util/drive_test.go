package util

import "testing"

func TestBuildEmbedURL(t *testing.T) {
	cases := []struct {
		name        string
		contentType string
		input       string
		wantID      string
		wantEmbed   string
	}{
		{
			name:        "link chia sẻ video",
			contentType: "video",
			input:       "https://drive.google.com/file/d/1A2b3C4d5E6f7G8h9I0j/view?usp=sharing",
			wantID:      "1A2b3C4d5E6f7G8h9I0j",
			wantEmbed:   "https://drive.google.com/file/d/1A2b3C4d5E6f7G8h9I0j/preview",
		},
		{
			name:        "link Google Slides",
			contentType: "slide",
			input:       "https://docs.google.com/presentation/d/1A2b3C4d5E6f7G8h9I0j/edit#slide=id.p1",
			wantID:      "1A2b3C4d5E6f7G8h9I0j",
			wantEmbed:   "https://docs.google.com/presentation/d/1A2b3C4d5E6f7G8h9I0j/embed?start=false&loop=false&delayms=5000",
		},
		{
			name:        "dán thẳng ID file",
			contentType: "pdf",
			input:       "1A2b3C4d5E6f7G8h9I0j",
			wantID:      "1A2b3C4d5E6f7G8h9I0j",
			wantEmbed:   "https://drive.google.com/file/d/1A2b3C4d5E6f7G8h9I0j/preview",
		},
		{
			name:        "link open?id=",
			contentType: "video",
			input:       "https://drive.google.com/open?id=1A2b3C4d5E6f7G8h9I0j",
			wantID:      "1A2b3C4d5E6f7G8h9I0j",
			wantEmbed:   "https://drive.google.com/file/d/1A2b3C4d5E6f7G8h9I0j/preview",
		},
		{
			name:        "link ngoài Drive giữ nguyên",
			contentType: "link",
			input:       "https://www.youtube.com/embed/abc",
			wantID:      "",
			wantEmbed:   "https://www.youtube.com/embed/abc",
		},
		{
			name:        "chuỗi rỗng",
			contentType: "video",
			input:       "   ",
			wantID:      "",
			wantEmbed:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, embed := BuildEmbedURL(tc.contentType, tc.input)
			if id != tc.wantID {
				t.Errorf("driveID = %q, mong đợi %q", id, tc.wantID)
			}
			if embed != tc.wantEmbed {
				t.Errorf("embedURL = %q, mong đợi %q", embed, tc.wantEmbed)
			}
		})
	}
}
