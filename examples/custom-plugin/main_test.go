package main

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestGreetingPluginTransformsPostsDeterministically(t *testing.T) {
	plugin := &greetingPlugin{}
	if plugin.Name() != "greeting" {
		t.Fatalf("unexpected plugin name %q", plugin.Name())
	}
	m := lifecycle.NewManager()
	visible := &models.Post{}
	skipped := &models.Post{Skip: true}
	m.SetPosts([]*models.Post{visible, skipped})
	if err := plugin.Transform(m); err != nil {
		t.Fatal(err)
	}
	if got := visible.Get("greeting"); got != "Hello from a site-owned Markata plugin" {
		t.Fatalf("unexpected post metadata: %v", got)
	}
	if skipped.Has("greeting") {
		t.Fatal("plugin modified a skipped post")
	}
	// A repeat build with identical inputs should produce identical metadata.
	if err := plugin.Transform(m); err != nil {
		t.Fatal(err)
	}
	if got := visible.Get("greeting"); got != "Hello from a site-owned Markata plugin" {
		t.Fatalf("non-deterministic repeated transform: %v", got)
	}
}
