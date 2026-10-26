package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestMDVideoPlugin_ProcessPost_CentersSizedVideo(t *testing.T) {
	p := NewMDVideoPlugin()
	post := &models.Post{
		ArticleHTML: `<p><img src="https://dropper.waylonwalker.com/file/demo.mp4?width=400" alt="demo"></p>`,
	}

	if err := p.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `margin-inline:auto`) {
		t.Errorf("expected generated video to be centered, got: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `src="https://dropper.waylonwalker.com/file/demo.mp4?width=400"`) {
		t.Errorf("expected width query parameter to be preserved, got: %s", post.ArticleHTML)
	}
}
