package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestGenerateRSS_WebSubDiscovery(t *testing.T) {
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"url": "https://example.com",
			"websub": models.WebSubConfig{
				Enabled: webSubBoolPtr(true),
				Hubs:    []string{"https://hub.example.com/"},
			},
		},
	}
	feed := &lifecycle.Feed{Path: "blog", Posts: []*models.Post{}}

	rss, err := GenerateRSS(feed, config)
	if err != nil {
		t.Fatalf("GenerateRSS() error = %v", err)
	}

	if !strings.Contains(rss, "rel=\"hub\"") {
		t.Fatalf("expected rel=\"hub\" in RSS output")
	}
	if !strings.Contains(rss, "https://hub.example.com/") {
		t.Fatalf("expected hub URL in RSS output")
	}
}

func TestGenerateAtom_WebSubDiscovery(t *testing.T) {
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"url": "https://example.com",
			"websub": models.WebSubConfig{
				Enabled: webSubBoolPtr(true),
				Hubs:    []string{"https://hub.example.com/"},
			},
		},
	}
	feed := &lifecycle.Feed{Path: "blog", Posts: []*models.Post{}}

	atom, err := GenerateAtom(feed, config)
	if err != nil {
		t.Fatalf("GenerateAtom() error = %v", err)
	}

	if !strings.Contains(atom, "rel=\"hub\"") {
		t.Fatalf("expected rel=\"hub\" in Atom output")
	}
	if !strings.Contains(atom, "https://hub.example.com/") {
		t.Fatalf("expected hub URL in Atom output")
	}
}

func webSubBoolPtr(value bool) *bool {
	return &value
}
