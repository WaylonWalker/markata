package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// MaxBrowserMessage bounds each JSON-RPC message, including a full draft.
const MaxBrowserMessage = 2 << 20

// BrowserConfig identifies the server-owned workspace and in-memory draft.
type BrowserConfig struct {
	RootURI     string `json:"rootUri"`
	DocumentURI string `json:"documentUri"`
}

// BrowserWorkspace validates a workspace without creating or changing files.
func BrowserWorkspace(root string) (BrowserConfig, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return BrowserConfig{}, fmt.Errorf("resolve LSP workspace: %w", err)
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return BrowserConfig{}, fmt.Errorf("open LSP workspace: %w", err)
	}
	if !info.IsDir() {
		return BrowserConfig{}, fmt.Errorf("LSP workspace must be a directory: %s", absolute)
	}
	return BrowserConfig{RootURI: pathToURI(absolute), DocumentURI: pathToURI(filepath.Join(absolute, ".markata", "browser-draft.md"))}, nil
}

// NewBrowserSession confines a session to a server-selected workspace and draft.
// It accepts unsaved buffers but ignores client requests to read watched files.
func NewBrowserSession(logger *log.Logger, root string) (*Server, error) {
	config, err := BrowserWorkspace(root)
	if err != nil {
		return nil, err
	}
	server := New(logger)
	server.browser = &config
	return server, nil
}

// RunMessages runs one session using one JSON-RPC object per transport message.
// The caller must interrupt receive when its context is canceled.
func (s *Server) RunMessages(ctx context.Context, receive func() ([]byte, error), send func([]byte) error) error {
	s.sendJSON = send
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		data, err := receive()
		if err != nil {
			return err
		}
		var message Message
		if len(data) > MaxBrowserMessage {
			return fmt.Errorf("LSP message exceeds %d bytes", MaxBrowserMessage)
		}
		if err = json.Unmarshal(data, &message); err != nil || message.JSONRPC != "2.0" || message.Method == "" {
			if err := s.sendError(json.RawMessage("null"), InvalidRequest, "expected one JSON-RPC 2.0 request or notification"); err != nil {
				return err
			}
			continue
		}
		s.refreshSharedIndex()
		if err := s.handleMessage(ctx, &message); err != nil {
			return err
		}
		if s.shutdown || message.Method == methodExit {
			return nil
		}
	}
}

// allowBrowserMessage pins initialization and accepts only the synthetic draft.
func (s *Server) allowBrowserMessage(msg *Message) bool {
	if s.browser == nil {
		return true
	}
	if msg.Method == methodInitialize {
		var params InitializeParams
		if json.Unmarshal(msg.Params, &params) != nil {
			return true
		} // normal invalid-params response
		root := s.browser.RootURI
		params.RootURI = &root
		params.RootPath = nil
		params.WorkspaceFolders = nil
		data, err := json.Marshal(params)
		if err != nil {
			return false
		}
		msg.Params = data
		return true
	}
	switch msg.Method {
	case "initialized", "shutdown", "exit", "$/cancelRequest", "$/setTrace":
		return true
	case "textDocument/didOpen", "textDocument/didChange", "textDocument/didClose", "textDocument/didSave", "textDocument/completion", "textDocument/hover", "textDocument/definition":
		var params struct {
			TextDocument TextDocumentIdentifier `json:"textDocument"`
			Position     *Position              `json:"position"`
		}
		if json.Unmarshal(msg.Params, &params) != nil {
			return false
		}
		return params.TextDocument.URI == s.browser.DocumentURI && (params.Position == nil || (params.Position.Line >= 0 && params.Position.Character >= 0))
	default:
		return false
	}
}

// NewWebSocketHandler exposes a local, same-origin LSP endpoint. Every connection
// gets its own session. Context cancellation closes active upgraded connections.
func NewWebSocketHandler(ctx context.Context, root string, logger *log.Logger) (http.Handler, error) {
	return NewWebSocketHandlerWithIndex(ctx, root, logger, nil)
}

// NewWebSocketHandlerWithIndex shares workspace parsing while isolating draft state.
func NewWebSocketHandlerWithIndex(ctx context.Context, root string, logger *log.Logger, index IndexProvider) (http.Handler, error) {
	config, err := BrowserWorkspace(root)
	if err != nil {
		return nil, err
	}
	root = uriToPath(config.RootURI)
	upgrader := websocket.Upgrader{HandshakeTimeout: 5 * time.Second, CheckOrigin: validBrowserOrigin}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server, err := NewBrowserSession(logger, root)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		server.sharedIndex = index
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
		defer stop()
		conn.SetReadLimit(MaxBrowserMessage)
		receive := func() ([]byte, error) {
			kind, data, readErr := conn.ReadMessage()
			if readErr == nil && kind != websocket.TextMessage {
				readErr = fmt.Errorf("expected text JSON-RPC message")
			}
			return data, readErr
		}
		send := func(data []byte) error {
			if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				return err
			}
			return conn.WriteMessage(websocket.TextMessage, data)
		}
		// Disconnects are normal. Protocol failures terminate only this session.
		if err := server.RunMessages(ctx, receive, send); err != nil {
			return
		}
	}), nil
}

func validBrowserOrigin(r *http.Request) bool {
	// Host and Origin are client-controlled headers. Require the actual TCP
	// peer to be loopback as well, even if --listen is configured broadly.
	// Never trust X-Forwarded-For to grant filesystem-aware LSP access.
	peer, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	address := net.ParseIP(peer)
	if address == nil || !address.IsLoopback() {
		return false
	}
	if !isLocalBrowserHost(r.Host) || r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		return false
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	origin, err := url.Parse(r.Header.Get("Origin"))
	return err == nil && origin.Scheme == scheme && origin.Host == r.Host && origin.User == nil &&
		origin.Path == "" && origin.RawQuery == "" && origin.Fragment == ""
}

func isLocalBrowserHost(hostport string) bool {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		host = strings.Trim(hostport, "[]")
	}
	if strings.EqualFold(strings.TrimSuffix(host, "."), "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
