package blogroll

import (
	"strings"
	"testing"
)

func TestPlanTOMLMutationReviewWorkflow(t *testing.T) {
	const address = "https://example.test/feed"
	original := "[blogroll]\nenabled = true\n\n[[blogroll.feeds]]\n" +
		"url = \"https://example.test/feed\"\ntitle = \"Original\"\nactive = true # keep comment\nmax_entries = 10\n"
	staged, err := PlanTOMLMutation(original, TOMLMutation{
		Action: MutationActive, URL: address, Active: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Replace(original, "active = true", "active = false", 1); staged != want {
		t.Fatalf("unexpected proposed change:\n%s", staged)
	}
	edited, err := PlanTOMLMutation(staged, TOMLMutation{
		Action: MutationUpdate, URL: address,
		Fields: map[string]string{"title": `"New title"`}, SelectedFields: []string{"title"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(edited, "active = false # keep comment") || !strings.Contains(edited, `title = "New title"`) ||
		!strings.Contains(edited, "max_entries = 10") {
		t.Fatalf("staged edits lost unrelated values: %s", edited)
	}
	removed, err := PlanTOMLMutation(edited, TOMLMutation{
		Action: MutationRemove, URL: address, RemoveLast: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(removed, address) || !strings.Contains(removed, "feeds = []") {
		t.Fatalf("final removal must preserve an empty list: %s", removed)
	}
}

func TestPlanTOMLMutationInitialAdd(t *testing.T) {
	const address = "https://example.test/feed"
	fields := map[string]string{
		"title":             `"Example"`,
		"description":       `""`,
		"category":          `"Tech"`,
		"type":              `"written"`,
		"tags":              `["go"]`,
		"handle":            `"example"`,
		mutationFieldActive: "true",
	}
	for _, original := range []string{"", "[markata-go]\n", "[blogroll]\nenabled = false\nfeeds = []\n"} {
		proposed, err := PlanTOMLMutation(original, TOMLMutation{
			Action: MutationAdd, URL: address, Fields: fields,
		})
		if err != nil {
			t.Fatalf("initial add from %q: %v", original, err)
		}
		if !strings.Contains(proposed, "enabled = true") || !strings.Contains(proposed, address) {
			t.Fatalf("unexpected add: %s", proposed)
		}
	}
}

func TestPlanTOMLMutationRejectsUnexpectedValues(t *testing.T) {
	const address = "https://example.test/feed"
	original := "[blogroll]\nenabled = true\n\n[[blogroll.feeds]]\n" +
		"url = \"https://example.test/feed\"\ntitle = \"Original\"\n\n[theme]\nname = \"Keep\"\n"
	req := TOMLMutation{
		Action: MutationUpdate, URL: address,
		Fields: map[string]string{"title": `"New title"`}, SelectedFields: []string{"title"},
	}
	proposed := strings.ReplaceAll(strings.ReplaceAll(original, "Original", "New title"), "Keep", "Changed")
	if err := verifyMutationDocument(original, proposed, req); err == nil {
		t.Fatal("must reject unrelated semantic changes")
	}
	req.Fields["title"] = "\"New title\"\nunknown = 1"
	if _, err := PlanTOMLMutation(original, req); err == nil {
		t.Fatal("must reject injected unrelated fields")
	}
	req = TOMLMutation{Action: MutationRemove, URL: address}
	if _, err := PlanTOMLMutation(original, req); err == nil {
		t.Fatal("must reject final removal without explicit empty array")
	}
}
