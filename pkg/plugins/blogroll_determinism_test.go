package plugins

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestCompareEntries_DeterministicTieBreakers(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	a := &models.ExternalEntry{FeedURL: "https://b.example.com/rss", ID: "b", Title: "B", Published: &baseTime}
	b := &models.ExternalEntry{FeedURL: "https://a.example.com/rss", ID: "a", Title: "A", Published: &baseTime}

	if compareEntries(a, b) {
		t.Fatalf("expected feed URL tie-breaker to sort ascending")
	}
	if !compareEntries(b, a) {
		t.Fatalf("expected feed URL tie-breaker to sort ascending")
	}
}

func TestLatestEntryDate_UsesLatestPublishedOrUpdated(t *testing.T) {
	older := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	newer := time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC)

	entries := []*models.ExternalEntry{
		{Published: &older},
		{Updated: &newer},
	}

	latest := latestEntryDate(entries)
	if latest == nil || !latest.Equal(newer) {
		t.Fatalf("expected latest date %v, got %v", newer, latest)
	}
}
