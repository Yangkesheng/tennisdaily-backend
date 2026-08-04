package storage

import "testing"

func TestFileID(t *testing.T) {
	s := New(Config{
		EnvID:  "prod-d6gkg1meqce4aaa7e",
		Bucket: "7072-prod-d6gkg1meqce4aaa7e-1440725645",
		Folder: "racket_library",
	})

	got := s.FileID("racket_library/Yonex/EZONE/12.jpg")
	want := "cloud://prod-d6gkg1meqce4aaa7e.7072-prod-d6gkg1meqce4aaa7e-1440725645/racket_library/Yonex/EZONE/12.jpg"
	if got != want {
		t.Fatalf("FileID() = %q, want %q", got, want)
	}
}

func TestObjectKey(t *testing.T) {
	s := New(Config{Folder: "racket_library"})

	cases := []struct {
		name   string
		brand  string
		series string
		id     int64
		ext    string
		want   string
	}{
		{"normal", "Yonex", "EZONE", 12, ".jpg", "racket_library/Yonex/EZONE/12.jpg"},
		{"space in brand", "Wilson Blade", "98 16x19", 7, ".png", "racket_library/Wilson-Blade/98-16x19/7.png"},
		{"chinese brand", "尤尼克斯", "EZONE", 1, ".jpg", "racket_library/尤尼克斯/EZONE/1.jpg"},
		{"empty series", "Yonex", "", 3, ".webp", "racket_library/Yonex/3.webp"},
		{"special chars", "Dunlop/Srixon", "CX 400 TOUR", 9, ".jpg", "racket_library/Dunlop-Srixon/CX-400-TOUR/9.jpg"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.ObjectKey(tc.id, tc.brand, tc.series, tc.ext)
			if got != tc.want {
				t.Fatalf("ObjectKey() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestExtensionFromURL(t *testing.T) {
	cases := []struct {
		name        string
		rawURL      string
		contentType string
		want        string
	}{
		{"jpeg url", "https://example.com/a.JPEG?x=1", "", ".jpg"},
		{"png url", "https://example.com/a.png", "", ".png"},
		{"no ext use content type", "https://example.com/img", "image/webp", ".webp"},
		{"no ext no content type", "https://example.com/img", "", ".jpg"},
		{"gif", "https://example.com/a.gif", "image/gif", ".gif"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtensionFromURL(tc.rawURL, tc.contentType)
			if got != tc.want {
				t.Fatalf("ExtensionFromURL() = %q, want %q", got, tc.want)
			}
		})
	}
}
