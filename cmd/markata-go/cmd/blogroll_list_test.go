package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestRenderBlogrollList(t *testing.T) {
	paused := false
	response := subscriptionResponse{
		Enabled: true,
		Feeds: []subscriptionView{
			{ExternalFeedConfig: models.ExternalFeedConfig{Title: "Engineering", URL: "https://example.com/feed", Tags: []string{"Go"}}, Source: "feeds.toml", Editable: true},
			{ExternalFeedConfig: models.ExternalFeedConfig{Title: "Podcast", URL: "https://example.org/feed", Active: &paused}, Source: "Read only"},
		},
	}
	var out bytes.Buffer
	if err := renderBlogrollList(&out, response, "", false); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Engineering [active]", "Podcast [paused]", "feeds.toml", "Read only"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing %q in %s", want, out.String())
		}
	}
	out.Reset()
	if err := renderBlogrollList(&out, response, "GO", true); err != nil {
		t.Fatal(err)
	}
	var got subscriptionResponse
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Feeds) != 1 || got.Feeds[0].Title != "Engineering" || !got.Feeds[0].Editable {
		t.Fatalf("unexpected filtered response: %+v", got)
	}
}

func TestRenderBlogrollListEmptyJSON(t *testing.T) {
	var out bytes.Buffer
	if err := renderBlogrollList(&out, subscriptionResponse{}, "", true); err != nil {
		t.Fatal(err)
	}
	var got struct {
		Feeds []subscriptionView `json:"feeds"`
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Feeds == nil || len(got.Feeds) != 0 {
		t.Fatalf("expected an empty array, got %s", out.String())
	}
}
