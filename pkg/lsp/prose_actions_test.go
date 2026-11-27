package lsp

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/prose"
)

func TestProseLSPRangeUsesUTF16Columns(t *testing.T) {
	source := "🌱 We very unique."
	issues := prose.Check(source, prose.Options{})
	if len(issues) != 1 {
		t.Fatalf("expected one style issue, got %+v", issues)
	}
	got := proseLSPRange(source, issues[0])
	if got.Start.Character != 6 {
		t.Fatalf("expected UTF-16 column 6, got %+v", got)
	}
	if got.End.Character != 17 {
		t.Fatalf("expected UTF-16 end column 17, got %+v", got)
	}
	if issues[0].Suggestion != "unique" {
		t.Fatalf("expected a deterministic replacement, got %+v", issues[0])
	}
}

func TestProseIgnoreCommandAndCodeActions(t *testing.T) {
	t.Setenv("MARKATA_PROSE_ENABLED", "1")
	root := t.TempDir()
	uri := pathToURI(filepath.Join(root, "post.md"))
	s := New(log.New(io.Discard, "", 0))
	s.rootURI = pathToURI(root)
	s.documents[uri] = &Document{URI: uri, Content: "This is very unique.", Version: 2}
	var messages []Message
	s.sendJSON = func(data []byte) error {
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		messages = append(messages, msg)
		return nil
	}
	request, err := json.Marshal(proseCodeActionParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
		Range:        Range{Start: Position{Line: 0, Character: 8}, End: Position{Line: 0, Character: 8}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.handleProseCodeAction(context.Background(), &Message{
		ID: json.RawMessage("1"), Params: request,
	}); err != nil {
		t.Fatal(err)
	}
	var actions []proseCodeAction
	if err := json.Unmarshal(messages[0].Result, &actions); err != nil {
		t.Fatal(err)
	}
	if len(actions) != 3 || actions[0].Edit == nil || actions[0].Edit.DocumentChanges[0].Edits[0].NewText != "unique" {
		t.Fatalf("expected replacement and two ignore actions, got %+v", actions)
	}
	args := make([]json.RawMessage, 4)
	for i, value := range []interface{}{uri, "prose/intensifier", 0, 8} {
		args[i], err = json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
	}
	ignore, err := json.Marshal(proseExecuteParams{
		Command:   "markata.prose.ignoreOnce",
		Arguments: args,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.handleProseExecuteCommand(context.Background(), &Message{
		ID: json.RawMessage("2"), Params: ignore,
	}); err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range s.computeDiagnostics(uri, s.documents[uri].Content) {
		if diagnostic.Source == "markata-prose" {
			t.Fatalf("ignored occurrence is still visible: %+v", diagnostic)
		}
	}
}
