package plugins

import (
	"io"
	"log"
	"regexp"
	"strconv"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestAutoFeedsCollectionPhaseRecord(t *testing.T) {
	original := log.Writer()
	flags, prefix := log.Flags(), log.Prefix()
	t.Cleanup(func() {
		log.SetOutput(original)
		log.SetFlags(flags)
		log.SetPrefix(prefix)
	})
	log.SetFlags(0)
	log.SetPrefix("")
	var messages []string
	log.SetOutput(logging.NewWriter(logging.Options{
		Writer: io.Discard, Format: logging.FormatPlain,
		Observer: func(entry logging.Entry, message string) {
			if entry.Component == "auto_feeds" && entry.Phase == "collect" && entry.Level == "debug" {
				messages = append(messages, message)
			}
		},
	}))
	m := lifecycle.NewManager()
	m.SetPosts([]*models.Post{
		{Path: "sensitive-source.md", Published: true, Tags: []string{"sensitive-tag"}},
		{Path: "private-source.md", Published: true, Private: true, Tags: []string{"other-tag"}},
		{Path: "excluded.md", Draft: true, Skip: true, Tags: []string{"sensitive-tag"}},
	})
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"auto_feeds": AutoFeedsConfig{Tags: AutoFeedTypeConfig{Enabled: true}}}
	m.SetConfig(config)
	if err := NewAutoFeedsPlugin().Collect(m); err != nil {
		t.Fatal(err)
	}
	if len(messages) != 1 {
		t.Fatalf("got %d phase records", len(messages))
	}
	pattern := regexp.MustCompile(`^generation_ns=(\d+) filtering_sorting_ns=(\d+) selection_recording_ns=(\d+) pagination_preparation_ns=(\d+) feeds=(\d+) matched=(\d+) selected=(\d+) observations=(\d+)$`)
	fields := pattern.FindStringSubmatch(messages[0])
	if fields == nil {
		t.Fatalf("unbounded or unexpected phase record: %q", messages[0])
	}
	for _, field := range fields[1:5] {
		value, err := strconv.ParseInt(field, 10, 64)
		// Short phases can finish within one clock tick, especially on Windows.
		if err != nil || value < 0 {
			t.Fatalf("phase duration must be a non-negative integer: %q", field)
		}
	}
	if fields[5] != "2" || fields[8] != "6" {
		t.Fatalf("incorrect feed/occurrence counts: %q", messages[0])
	}
	for _, entry := range m.ContentLedger().Snapshot().Entries {
		if len(entry.Feeds) != 2 {
			t.Fatalf("lost exclusion observations: %q", entry.Path)
		}
	}
}
