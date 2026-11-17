package editorworkspace

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestFileStatusAndDate(t *testing.T) {
	day := time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		post models.Post
		want string
	}{
		{"draft wins over published", models.Post{Draft: true, Published: true}, "draft"},
		{"published", models.Post{Published: true}, "published"},
		{"neither", models.Post{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileStatus(&tt.post); got != tt.want {
				t.Errorf("fileStatus() = %q, want %q", got, tt.want)
			}
		})
	}
	if got := fileDate(&models.Post{Date: &day}); got != "2026-03-04" {
		t.Errorf("fileDate() = %q, want 2026-03-04", got)
	}
	if got := fileDate(&models.Post{}); got != "" {
		t.Errorf("fileDate() without date = %q, want empty", got)
	}
}
