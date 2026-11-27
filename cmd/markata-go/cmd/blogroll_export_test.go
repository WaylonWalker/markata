package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteNewOPMLFileNoClobber(t *testing.T) {
	path := filepath.Join(t.TempDir(), "subscriptions.opml")
	if err := writeNewOPMLFile(path, strings.NewReader("original")); err != nil {
		t.Fatal(err)
	}
	if err := writeNewOPMLFile(path, strings.NewReader("replacement")); err == nil {
		t.Fatal("expected an existing export to be protected")
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != "original" {
		t.Fatalf("unexpected export: %q, %v", got, err)
	}
}
