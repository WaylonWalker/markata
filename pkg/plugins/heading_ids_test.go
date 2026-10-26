package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestHeadingIDTransformer_LinkDestinationDoesNotLeakIntoID(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "## Recent [TIL](/til/)\n\n## Plain `code` heading\n\n## Custom {#custom-id}\n\n## Recent TIL\n"}
	if err := p.renderPost(post); err != nil {
		t.Fatalf("renderPost: %v", err)
	}
	for _, want := range []string{`id="recent-til"`, `id="plain-code-heading"`, `id="custom-id"`, `id="recent-til-1"`} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("expected %s in %s", want, post.ArticleHTML)
		}
	}

	if strings.Contains(post.ArticleHTML, "tiltil") {
		t.Errorf("link destination leaked into heading id: %s", post.ArticleHTML)
	}
}

func TestHeadingIDTransformer_CollisionOrderMatchesTOC(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "## Recent [TIL](/til/)\n\n## Recent TIL\n"}

	if err := p.renderPost(post); err != nil {
		t.Fatalf("renderPost: %v", err)
	}

	first := strings.Index(post.ArticleHTML, `<h2 id="recent-til">Recent <a href="/til/">TIL</a></h2>`)
	second := strings.Index(post.ArticleHTML, `<h2 id="recent-til-1">Recent TIL</h2>`)
	if first == -1 || second == -1 || first > second {
		t.Fatalf("heading IDs do not match source order: %s", post.ArticleHTML)
	}
}
