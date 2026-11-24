package templates

import "testing"

func TestMediaDimensionsFromURL(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantWidth  int
		wantHeight int
		wantOK     bool
	}{
		{name: "w and h", input: "https://example.com/image.jpg?w=1200&h=675", wantWidth: 1200, wantHeight: 675, wantOK: true},
		{name: "width and height", input: "https://example.com/image.jpg?width=800&height=600", wantWidth: 800, wantHeight: 600, wantOK: true},
		{name: "width only", input: "https://example.com/image.jpg?w=900", wantWidth: 900, wantHeight: 0, wantOK: true},
		{name: "unrelated query", input: "https://example.com/image.jpg?foo=bar", wantWidth: 0, wantHeight: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height, ok := MediaDimensionsFromURL(tt.input)
			if width != tt.wantWidth || height != tt.wantHeight || ok != tt.wantOK {
				t.Fatalf("MediaDimensionsFromURL(%q) = (%d, %d, %v), want (%d, %d, %v)", tt.input, width, height, ok, tt.wantWidth, tt.wantHeight, tt.wantOK)
			}
		})
	}
}

func TestResponsiveSrcsetUsesTrustedRemoteWidthsAndPreservesQuery(t *testing.T) {
	got := ResponsiveSrcset("http://dropper.waylonwalker.com/image.webp?token=abc&w=600#preview", []int{480, 160, 320, 320, -1})
	want := "https://dropper.waylonwalker.com/image.webp?token=abc&w=160#preview 160w, " +
		"https://dropper.waylonwalker.com/image.webp?token=abc&w=320#preview 320w, " +
		"https://dropper.waylonwalker.com/image.webp?token=abc&w=480#preview 480w"
	if got != want {
		t.Fatalf("ResponsiveSrcset() = %q, want %q", got, want)
	}
}

func TestResponsiveSrcsetSkipsLocalAndUntrustedURLs(t *testing.T) {
	for _, raw := range []string{"/images/photo.jpg", "https://example.com/photo.jpg", "ftp://dropper.waylonwalker.com/photo.jpg", "data:image/png;base64,abc"} {
		if got := ResponsiveSrcset(raw, []int{160, 320}); got != "" {
			t.Errorf("ResponsiveSrcset(%q) = %q, want empty", raw, got)
		}
	}
}

func TestResponsiveSrcsetRequiresTwoDistinctPositiveWidths(t *testing.T) {
	for _, widths := range [][]int{{320}, {0, -1}, {320, 320}} {
		if got := ResponsiveSrcset("https://dropper.waylonwalker.com/photo.jpg", widths); got != "" {
			t.Errorf("ResponsiveSrcset(%v) = %q, want empty", widths, got)
		}
	}
}

func TestOriginalMediaURLRemovesSizeAndPreservesOtherQuery(t *testing.T) {
	got := OriginalMediaURL("http://dropper.waylonwalker.com/photo.webp?token=abc&w=600&h=400&width=500#preview")
	want := "https://dropper.waylonwalker.com/photo.webp?token=abc#preview"
	if got != want {
		t.Fatalf("OriginalMediaURL() = %q, want %q", got, want)
	}
}

func TestOriginalMediaURLLeavesUnknownHostsUntouched(t *testing.T) {
	raw := "https://example.com/photo.webp?w=600&token=abc"
	if got := OriginalMediaURL(raw); got != raw {
		t.Fatalf("OriginalMediaURL() = %q, want %q", got, raw)
	}
}

func TestDerivedVideoPosterDropsSourceDimensions(t *testing.T) {
	got := PosterURLFromMap(map[string]interface{}{},
		"https://dropper.wayl.one/file/clip.webm?width=400&height=225&w=320&h=180&token=abc#preview")
	want := "https://dropper.wayl.one/file/clip.webp?token=abc#preview"
	if got != want {
		t.Fatalf("PosterURLFromMap() = %q, want %q", got, want)
	}
	sized := WithSize(got, 1200, 675)
	wantSized := "https://dropper.wayl.one/file/clip.webp?h=675&token=abc&w=1200#preview"
	if sized != wantSized {
		t.Fatalf("WithSize(poster) = %q, want %q", sized, wantSized)
	}
}

func TestExplicitVideoPosterKeepsAuthorURL(t *testing.T) {
	explicit := "https://dropper.wayl.one/file/custom.webp?width=400&token=abc"
	got := PosterURLFromMap(map[string]interface{}{"poster_image": explicit}, "https://dropper.wayl.one/file/clip.webm?width=400")
	if got != explicit {
		t.Fatalf("explicit poster changed: got %q, want %q", got, explicit)
	}
}

func TestVideoMediaHelpers_HandleWindowsPaths(t *testing.T) {
	for _, path := range []string{`C:\Users\runner\static\clip.Mp4`, `C:/Users/runner/static/clip.Mp4`} {
		if !IsVideoURL(path) {
			t.Errorf("IsVideoURL(%q) = false, want true", path)
		}
		if got := VideoMIMEType(path); got != "video/mp4" {
			t.Errorf("VideoMIMEType(%q) = %q, want video/mp4", path, got)
		}
	}
	if IsVideoURL("javascript:clip.mp4") {
		t.Fatal("IsVideoURL classified an opaque javascript URL as video")
	}
}
