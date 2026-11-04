package editorsync

import (
	"context"
	"strings"
	"testing"
)

func testMutation() Mutation {
	return Mutation{
		SiteID: "family-blog",
		DocumentID: "post-001",
		SubjectID: "subject-42",
		MutationID: "8eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa",
		BaseRevision: 3,
		Visibility: Private,
		Source: []byte("# Hello\n"),
	}
}

func TestMutationValidation(t *testing.T) {
	base := testMutation()
	if err := base.Validate(); err != nil {
		t.Fatalf("valid mutation: %v", err)
	}
	tests := []struct {
		name string
		change func(*Mutation)
	}{
		{"invalid site", func(m *Mutation) { m.SiteID = "../bad" }},
		{"empty document", func(m *Mutation) { m.DocumentID = "" }},
		{"invalid subject", func(m *Mutation) { m.SubjectID = "x\ny" }},
		{"missing ID", func(m *Mutation) { m.MutationID = "" }},
		{"bad ID", func(m *Mutation) { m.MutationID = "duplicate" }},
		{"unknown visibility", func(m *Mutation) { m.Visibility = "admin" }},
		{"oversized source", func(m *Mutation) { m.Source = make([]byte, 2<<20+1) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := base
			tt.change(&value)
			if err := value.Validate(); err == nil {
				t.Error("invalid mutation accepted")
			}
		})
	}
}

func TestDigestCoversEntireMutationAndIsStable(t *testing.T) {
	base := testMutation()
	first := base.Digest()
	if len(first) != 64 || base.Digest() != first {
		t.Fatal("digest is not stable SHA-256")
	}
	for _, tt := range []struct {
		name string
		change func(*Mutation)
	}{
		{"site", func(m *Mutation) { m.SiteID = "other-site" }},
		{"document", func(m *Mutation) { m.DocumentID = "post-002" }},
		{"subject", func(m *Mutation) { m.SubjectID = "subject-43" }},
		{"mutation", func(m *Mutation) { m.MutationID = "9eb7a4ab-4b5d-4ba7-b189-40ff3dd0eafa" }},
		{"base revision", func(m *Mutation) { m.BaseRevision++ }},
		{"visibility", func(m *Mutation) { m.Visibility = SharedDraft }},
		{"source", func(m *Mutation) { m.Source = []byte("# Changed\n") }},
		{"deletion", func(m *Mutation) { m.Deleted = true }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			value := base
			tt.change(&value)
			if value.Digest() == first {
				t.Error("changed mutation intent produced same digest")
			}
		})
	}
	// Explicit length framing prevents ambiguity from concatenated fields.
	a := base
	a.DocumentID, a.SubjectID = "ab", "c"
	b := base
	b.DocumentID, b.SubjectID = "a", "bc"
	if a.Digest() == b.Digest() {
		t.Error("digest lacks field boundaries")
	}
}

func TestMigrationSourceAndNilDB(t *testing.T) {
	if err := Migrate(context.Background(), nil); err == nil {
		t.Error("nil SQLite connection accepted")
	}
	raw, err := migrations.ReadFile("schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"editor_sites", "editor_documents", "editor_document_versions",
		"editor_mutation_receipts", "editor_changes", "editor_conflicts",
	} {
		if !strings.Contains(string(raw), "CREATE TABLE IF NOT EXISTS "+name) {
			t.Errorf("missing table %s", name)
		}
	}
}
