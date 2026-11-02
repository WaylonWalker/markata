package lsp

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func browserFixture(t *testing.T) (string, BrowserConfig) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "hello world.md"), []byte("---\ntitle: Hello World\nslug: hello\naliases: [greeting]\n---\nHello"), 0o600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "output")
	if err := os.MkdirAll(output, 0o700); err != nil {
		t.Fatal(err)
	}
	images := `{"$schema":"markata://schemas/image-index/v1","schema":"markata.image-index","schema_version":1,"generator":{"name":"test","version":"1"},"image_count":1,"images":[{"src":"/images/hello%20world.png","width":100,"height":50,"alt":"Hello","mime_type":"image/png","cover":false,"uses":[]}]}`
	if err := os.WriteFile(filepath.Join(output, "images.json"), []byte(images), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := BrowserWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	return root, config
}

func TestLoadManagedImages_UsesConfiguredOutputDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "markata-go.toml"), []byte("[markata-go]\noutput_dir = \"public\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "public")
	if err := os.MkdirAll(output, 0o700); err != nil {
		t.Fatal(err)
	}
	images := `{"$schema":"markata://schemas/image-index/v1","schema":"markata.image-index","schema_version":1,"generator":{"name":"test","version":"1"},"image_count":1,"images":[{"src":"/images/managed.png","width":100,"height":50,"alt":"Managed","mime_type":"image/png","cover":false,"uses":[]}]}`
	if err := os.WriteFile(filepath.Join(output, "images.json"), []byte(images), 0o600); err != nil {
		t.Fatal(err)
	}
	got := loadManagedImages(root)
	if len(got) != 1 || got[0].Src != "/images/managed.png" {
		t.Fatalf("loadManagedImages() = %#v, want configured image index", got)
	}
}

func browserConnect(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()
	conn, response, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), http.Header{"Origin": {server.URL}})
	if response != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if err = conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		t.Fatal(err)
	}
	return conn
}

func browserSend(t *testing.T, conn *websocket.Conn, method string, id int, params any) {
	t.Helper()
	msg := map[string]any{"jsonrpc": "2.0", "method": method, "params": params}
	if id > 0 {
		msg["id"] = id
	}
	if err := conn.WriteJSON(msg); err != nil {
		t.Fatal(err)
	}
}

func browserResponse(t *testing.T, conn *websocket.Conn, id int) Message {
	t.Helper()
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			t.Fatal(err)
		}
		if string(msg.ID) == jsonID(id) {
			return msg
		}
	}
}
func jsonID(id int) string { return strconv.Itoa(id) }

func browserInitialize(t *testing.T, conn *websocket.Conn) {
	t.Helper()
	// A client cannot redirect the server's workspace.
	browserSend(t, conn, "initialize", 1, map[string]any{"rootUri": "file:///not-the-workspace", "capabilities": map[string]any{}})
	if msg := browserResponse(t, conn, 1); msg.Error != nil {
		t.Fatal(msg.Error)
	}
	browserSend(t, conn, "initialized", 0, map[string]any{})
}

func TestBrowserWebSocket_CompletionsAndIsolation(t *testing.T) {
	root, config := browserFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler, err := NewWebSocketHandler(ctx, root, log.New(io.Discard, "", 0))
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	first := browserConnect(t, server)
	browserSend(t, first, "textDocument/completion", 99, CompletionParams{TextDocument: TextDocumentIdentifier{URI: config.DocumentURI}, Position: Position{}})
	if msg := browserResponse(t, first, 99); msg.Error == nil || msg.Error.Code != ServerNotInitialized {
		t.Fatal("language request ran before initialization")
	}
	browserInitialize(t, first)
	for i, tt := range []struct {
		text, want string
		start      int
	}{
		{"😀 [[gre", "greeting", 5},
		{"😀 [Hello](hello", "hello%20world.md", 11},
		{"[Hello](", "hello%20world.md", 8},
		{"![](", "/images/hello%20world.png", 4},
		{"![embed](", "", 9},
	} {
		browserSend(t, first, "textDocument/didOpen", 0, map[string]any{"textDocument": map[string]any{"uri": config.DocumentURI, "languageId": "markdown", "version": i + 1, "text": tt.text}})
		browserSend(t, first, "textDocument/completion", i+2, CompletionParams{TextDocument: TextDocumentIdentifier{URI: config.DocumentURI}, Position: Position{Line: 0, Character: utf16Column(tt.text, len(tt.text))}})
		msg := browserResponse(t, first, i+2)
		if msg.Error != nil {
			t.Fatal(msg.Error)
		}
		var result CompletionList
		if err = json.Unmarshal(msg.Result, &result); err != nil {
			t.Fatal(err)
		}
		found := false
		for _, item := range result.Items {
			if item.Label == tt.want {
				found = true
				if item.TextEdit == nil || item.TextEdit.Range.Start.Character != tt.start {
					t.Fatalf("%s: wrong edit: %+v", tt.text, item.TextEdit)
				}
			}
		}
		if tt.want == "" {
			if len(result.Items) != 0 {
				t.Fatalf("%s: expected no suggestions, got %+v", tt.text, result.Items)
			}
			continue
		}
		if !found {
			t.Fatalf("%s: missing %s in %+v", tt.text, tt.want, result)
		}
	}
	second := browserConnect(t, server)
	browserInitialize(t, second)
	browserSend(t, second, "textDocument/completion", 9, CompletionParams{TextDocument: TextDocumentIdentifier{URI: config.DocumentURI}, Position: Position{}})
	msg := browserResponse(t, second, 9)
	var empty CompletionList
	if err = json.Unmarshal(msg.Result, &empty); err != nil {
		t.Fatal(err)
	}
	if len(empty.Items) != 0 {
		t.Fatal("draft leaked between sessions")
	}
	browserSend(t, first, "textDocument/completion", 10, CompletionParams{TextDocument: TextDocumentIdentifier{URI: "file:///etc/passwd"}, Position: Position{}})
	if msg = browserResponse(t, first, 10); msg.Error == nil {
		t.Fatal("accepted foreign document")
	}
	browserSend(t, first, "textDocument/completion", 11, CompletionParams{TextDocument: TextDocumentIdentifier{URI: config.DocumentURI}, Position: Position{Line: -1, Character: -1}})
	if msg = browserResponse(t, first, 11); msg.Error == nil {
		t.Fatal("accepted negative position")
	}
	browserSend(t, first, "shutdown", 12, nil)
	if msg = browserResponse(t, first, 12); msg.Error != nil {
		t.Fatal(msg.Error)
	}
	if _, _, err = first.ReadMessage(); err == nil {
		t.Fatal("shutdown left connection open")
	}
	if _, err = os.Stat(uriToPath(config.DocumentURI)); !os.IsNotExist(err) {
		t.Fatalf("draft written to disk: %v", err)
	}
	cancel()
	if _, _, err = second.ReadMessage(); err == nil {
		t.Fatal("context cancellation left connection open")
	}
}

func TestBrowserWebSocket_OriginAndMessageBounds(t *testing.T) {
	root, _ := browserFixture(t)
	handler, err := NewWebSocketHandler(context.Background(), root, log.New(io.Discard, "", 0))
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	for _, test := range []struct{ origin, host string }{
		{"", ""},
		{"https://foreign.invalid", ""},
		{"http://evil.invalid", "evil.invalid"}, // DNS-rebinding style Host and Origin.
	} {
		headers := http.Header{"Origin": {test.origin}}
		if test.host != "" {
			headers.Set("Host", test.host)
		}
		conn, response, dialErr := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), headers)
		if conn != nil {
			_ = conn.Close()
		}
		if response != nil {
			_ = response.Body.Close()
		}
		if dialErr == nil || response == nil || response.StatusCode != http.StatusForbidden {
			t.Fatalf("accepted origin %q with Host %q", test.origin, test.host)
		}
	}
	conn := browserConnect(t, server)
	if err = conn.WriteMessage(websocket.TextMessage, []byte("not JSON")); err != nil {
		t.Fatal(err)
	}
	var invalid Message
	if err = conn.ReadJSON(&invalid); err != nil || invalid.Error == nil {
		t.Fatalf("bad JSON response: %v %+v", err, invalid)
	}
	browserInitialize(t, conn)
	if err = conn.WriteMessage(websocket.TextMessage, []byte(strings.Repeat("x", MaxBrowserMessage+1))); err != nil {
		t.Fatal(err)
	}
	if _, _, err = conn.ReadMessage(); err == nil {
		t.Fatal("accepted oversized message")
	}
}

func TestBrowserSession_IgnoresWatchedFiles(t *testing.T) {
	root, _ := browserFixture(t)
	server, err := NewBrowserSession(log.New(io.Discard, "", 0), root)
	if err != nil {
		t.Fatal(err)
	}
	server.initialized = true
	if !server.allowBrowserMessage(&Message{Method: "textDocument/hover", Params: json.RawMessage(`{"textDocument":{"uri":"` + server.browser.DocumentURI + `"},"position":{"line":0,"character":0}}`)}) {
		t.Fatal("hover unavailable")
	}
	if server.allowBrowserMessage(&Message{Method: "workspace/didChangeWatchedFiles", Params: json.RawMessage(`{"changes":[{"uri":"file:///outside.md","type":1}]}`)}) {
		t.Fatal("browser can trigger arbitrary disk reads")
	}
}

func TestBrowserOrigin_RejectsSpoofedLocalhostFromRemotePeer(t *testing.T) {
	for _, tc := range []struct {
		peer    string
		allowed bool
	}{
		{"127.0.0.1:42300", true},
		{"[::1]:42300", true},
		{"192.0.2.44:42300", false},
		{"[2001:db8::44]:42300", false},
		{"malformed", false},
	} {
		t.Run(tc.peer, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8766/lsp", http.NoBody)
			request.Host = "localhost:8766"
			request.Header.Set("Origin", "http://localhost:8766")
			request.RemoteAddr = tc.peer
			if got := validBrowserOrigin(request); got != tc.allowed {
				t.Fatalf("validBrowserOrigin(peer=%q) = %v, want %v", tc.peer, got, tc.allowed)
			}
		})
	}
}
