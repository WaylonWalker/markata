package lsp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentLinkCompletion_RelativeAndWorkspacePaths(t *testing.T) {
	root, _ := browserFixture(t)
	server := New(log.New(io.Discard, "", 0))
	server.rootURI = pathToURI(root)
	if err := server.index.Build(root); err != nil {
		t.Fatal(err)
	}
	uri := pathToURI(filepath.Join(root, "posts", "draft.md"))
	for _, tt := range []struct{ source, want string }{
		{"[Doc](../he", "../hello%20world.md"},
		{"[Doc](/he", "/hello%20world.md"},
		{"[Doc](https://example.com/", ""},
		{"[Doc](#section", ""},
	} {
		t.Run(tt.source, func(t *testing.T) {
			server.documents[uri] = &Document{URI: uri, Content: tt.source, Version: 1}
			var output bytes.Buffer
			server.writer = &output
			params, err := json.Marshal(CompletionParams{TextDocument: TextDocumentIdentifier{URI: uri}, Position: Position{Line: 0, Character: len(tt.source)}})
			if err != nil {
				t.Fatal(err)
			}
			if err = server.handleCompletion(context.Background(), &Message{ID: json.RawMessage("1"), Params: params}); err != nil {
				t.Fatal(err)
			}
			response, err := server.readMessage(bufio.NewReader(&output))
			if err != nil {
				t.Fatal(err)
			}
			var list CompletionList
			if err = json.Unmarshal(response.Result, &list); err != nil {
				t.Fatal(err)
			}
			if tt.want == "" {
				if len(list.Items) != 0 {
					t.Fatalf("unexpected document suggestions: %+v", list.Items)
				}
				return
			}
			if len(list.Items) != 1 || list.Items[0].TextEdit == nil || list.Items[0].TextEdit.NewText != tt.want {
				t.Fatalf("unexpected completions: %+v", list)
			}
		})
	}
}

func TestBrowserWorkspace_EscapedURI(t *testing.T) {
	root := filepath.Join(t.TempDir(), "a workspace #1")
	// URI conversion does not depend on the path existing.
	uri := pathToURI(root)
	if !strings.Contains(uri, "%20") || !strings.Contains(uri, "%23") {
		t.Fatalf("unescaped URI: %s", uri)
	}
	if got := uriToPath(uri); got != root {
		t.Fatalf("roundtrip: got %s, want %s", got, root)
	}
}

func TestBrowserDiagnostics_UTF16AndRevision(t *testing.T) {
	root, config := browserFixture(t)
	server, err := NewBrowserSession(log.New(io.Discard, "", 0), root)
	if err != nil {
		t.Fatal(err)
	}
	source := "😀 [[missing]]"
	server.documents[config.DocumentURI] = &Document{URI: config.DocumentURI, Content: source, Version: 7}
	var output bytes.Buffer
	server.writer = &output
	if err = server.publishDiagnostics(config.DocumentURI, source); err != nil {
		t.Fatal(err)
	}
	response, err := server.readMessage(bufio.NewReader(&output))
	if err != nil {
		t.Fatal(err)
	}
	var params PublishDiagnosticsParams
	if err = json.Unmarshal(response.Params, &params); err != nil {
		t.Fatal(err)
	}
	if params.Version == nil || *params.Version != 7 {
		t.Fatalf("missing document revision: %+v", params)
	}
	for _, diagnostic := range params.Diagnostics {
		if strings.Contains(diagnostic.Message, "missing") {
			if diagnostic.Range.Start.Character != 3 || diagnostic.Range.End.Character != 14 {
				t.Fatalf("wrong UTF-16 range: %+v", diagnostic.Range)
			}
			return
		}
	}
	t.Fatal("missing broken-link diagnostic")
}
