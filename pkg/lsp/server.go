// Package server provides the Language Server Protocol implementation for markata-go.
//
// # Overview
//
// The LSP server provides IDE features for markdown files with wikilink support:
//   - Autocomplete for [[wikilinks]] (textDocument/completion)
//   - Diagnostics for broken wikilinks (textDocument/publishDiagnostics)
//   - Hover information showing post title and description (textDocument/hover)
//   - Go to definition for navigating to linked posts (textDocument/definition)
//
// # Architecture
//
// The server maintains an index of all markdown posts in the workspace, including:
//   - File paths and slugs
//   - Titles and descriptions
//   - Wikilinks contained in each file
//
// The index is built on initialization and updated incrementally as files change.
package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
)

// Server is the LSP server implementation for markata-go.
type Server struct {
	logger   *log.Logger
	browser  *BrowserConfig
	sendJSON func([]byte) error

	// index holds the indexed posts for autocomplete and validation
	index *Index

	// documents tracks open documents by URI
	documents map[string]*Document
	docMu     sync.RWMutex

	// initialized indicates whether the server has been initialized
	initialized bool

	// rootURI is the workspace root URI
	rootURI string

	// shutdown indicates the server is shutting down
	shutdown bool

	// writer for sending responses
	writer io.Writer
	wrMu   sync.Mutex
}

// Document represents an open document in the editor.
type Document struct {
	URI     string
	Content string
	Version int
}

// New creates a new LSP server.
func New(logger *log.Logger) *Server {
	return &Server{
		logger:    logger,
		index:     NewIndex(logger),
		documents: make(map[string]*Document),
	}
}

// Run starts the LSP server, reading from reader and writing to writer.
// The server uses JSON-RPC 2.0 over the LSP protocol.
func (s *Server) Run(ctx context.Context, reader io.Reader, writer io.Writer) error {
	s.writer = writer
	bufReader := bufio.NewReaderSize(reader, 1024*1024) // 1MB buffer

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msg, err := s.readMessage(bufReader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			s.logger.Printf("Error reading message: %v", err)
			continue
		}

		if msg == nil {
			continue
		}

		if err := s.handleMessage(ctx, msg); err != nil {
			s.logger.Printf("Error handling message: %v", err)
		}

		if s.shutdown {
			return nil
		}
	}
}

// Message represents a JSON-RPC 2.0 message.
type Message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// ResponseError represents a JSON-RPC 2.0 error.
type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LSP error codes
const (
	// JSON-RPC errors
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603

	// LSP errors
	ServerNotInitialized = -32002
	RequestCancelled     = -32800
)

// LSP method names
const (
	methodExit       = "exit"
	methodInitialize = "initialize"
)

// methodHandler is a handler function for an LSP method.
type methodHandler func(ctx context.Context, msg *Message) error

// methodHandlers returns a map of method names to their handlers.
func (s *Server) methodHandlers() map[string]methodHandler {
	return map[string]methodHandler{
		// Lifecycle methods
		"initialize":  s.handleInitialize,
		"initialized": s.handleInitialized,
		"shutdown":    s.handleShutdown,
		"exit":        s.handleExit,

		// Document synchronization
		"textDocument/didOpen":   s.handleDidOpen,
		"textDocument/didChange": s.handleDidChange,
		"textDocument/didClose":  s.handleDidClose,
		"textDocument/didSave":   s.handleDidSave,

		// Language features
		"textDocument/completion": s.handleCompletion,
		"textDocument/hover":      s.handleHover,
		"textDocument/definition": s.handleDefinition,

		// Workspace
		"workspace/didChangeWatchedFiles": s.handleDidChangeWatchedFiles,
	}
}

// readMessage reads a single LSP message from the reader.
// LSP uses Content-Length based framing with \r\n line endings in headers.
func (s *Server) readMessage(reader *bufio.Reader) (*Message, error) {
	// Read headers until empty line
	var contentLength int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("failed to read header: %w", err)
		}

		// Trim \r\n or \n
		line = strings.TrimRight(line, "\r\n")

		// Empty line signals end of headers
		if line == "" {
			break
		}

		// Parse Content-Length header
		if strings.HasPrefix(line, "Content-Length: ") {
			var parseErr error
			contentLength, parseErr = strconv.Atoi(strings.TrimPrefix(line, "Content-Length: "))
			if parseErr != nil {
				return nil, fmt.Errorf("invalid Content-Length: %w", parseErr)
			}
		}
		// Ignore other headers (like Content-Type)
	}

	if contentLength <= 0 || (s.browser != nil && contentLength > MaxBrowserMessage) {
		return nil, fmt.Errorf("missing, invalid, or excessive Content-Length header")
	}

	// Read exactly contentLength bytes for the body
	content := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, content); err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	var msg Message
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	return &msg, nil
}

// handleMessage dispatches a message to the appropriate handler.
func (s *Server) handleMessage(ctx context.Context, msg *Message) error {
	isRequest := s.isRequest(msg)
	if !s.allowBrowserMessage(msg) {
		if isRequest {
			return s.sendError(msg.ID, InvalidParams, "method or document is outside the browser session")
		}
		return nil
	}

	if (s.shutdown && msg.Method != methodExit) || (!s.initialized && msg.Method != methodInitialize && msg.Method != methodExit) {
		return s.checkMessageAllowed(msg, isRequest)
	}

	return s.dispatchMethod(ctx, msg, isRequest)
}

// isRequest returns true if the message is a request (has an ID).
func (s *Server) isRequest(msg *Message) bool {
	return msg.ID != nil && string(msg.ID) != "null"
}

// checkMessageAllowed validates that the message can be processed in the current server state.
// Returns an error response if the message is not allowed, nil otherwise.
func (s *Server) checkMessageAllowed(msg *Message, isRequest bool) error {
	// Handle shutdown state - only exit is allowed
	if s.shutdown && msg.Method != methodExit {
		if isRequest {
			return s.sendError(msg.ID, InvalidRequest, "server is shutting down")
		}
		return nil
	}

	// Check initialization - only initialize and exit are allowed before initialization
	if !s.initialized && msg.Method != methodInitialize && msg.Method != methodExit {
		if isRequest {
			return s.sendError(msg.ID, ServerNotInitialized, "server not initialized")
		}
		return nil
	}

	return nil
}

// dispatchMethod routes the message to the appropriate handler.
func (s *Server) dispatchMethod(ctx context.Context, msg *Message, isRequest bool) error {
	handlers := s.methodHandlers()

	if handler, ok := handlers[msg.Method]; ok {
		return handler(ctx, msg)
	}

	// Unknown method
	if isRequest {
		return s.sendError(msg.ID, MethodNotFound, fmt.Sprintf("method not found: %s", msg.Method))
	}
	// Ignore unknown notifications
	return nil
}

// sendResponse sends a JSON-RPC response.
func (s *Server) sendResponse(id json.RawMessage, result interface{}) error {
	response := Message{
		JSONRPC: "2.0",
		ID:      id,
	}

	if result != nil {
		data, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}
		response.Result = data
	} else {
		response.Result = json.RawMessage("null")
	}

	return s.writeMessage(&response)
}

// sendError sends a JSON-RPC error response.
func (s *Server) sendError(id json.RawMessage, code int, message string) error {
	response := Message{
		JSONRPC: "2.0",
		ID:      id,
		Error: &ResponseError{
			Code:    code,
			Message: message,
		},
	}
	return s.writeMessage(&response)
}

// sendNotification sends a JSON-RPC notification (no ID).
func (s *Server) sendNotification(method string, params interface{}) error {
	msg := Message{
		JSONRPC: "2.0",
		Method:  method,
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal params: %w", err)
		}
		msg.Params = data
	}

	return s.writeMessage(&msg)
}

// writeMessage writes a message to the output stream.
func (s *Server) writeMessage(msg *Message) error {
	s.wrMu.Lock()
	defer s.wrMu.Unlock()

	content, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if s.sendJSON != nil {
		return s.sendJSON(content)
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	if _, err := io.WriteString(s.writer, header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := s.writer.Write(content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return nil
}
