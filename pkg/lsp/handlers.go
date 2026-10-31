package lsp

import (
	"context"
	"encoding/json"
	"strings"
)

// handleInitialize handles the initialize request.
func (s *Server) handleInitialize(_ context.Context, msg *Message) error {
	var params InitializeParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		return s.sendError(msg.ID, InvalidParams, "invalid initialize params")
	}

	// Store root URI
	if params.RootURI != nil {
		s.rootURI = *params.RootURI
	} else if params.RootPath != nil {
		s.rootURI = pathToURI(*params.RootPath)
	}

	s.logger.Printf("Initializing with root: %s", s.rootURI)

	// Return server capabilities
	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: &TextDocumentSyncOptions{
				OpenClose: true,
				Change:    TextDocumentSyncKindFull,
				Save: &SaveOptions{
					IncludeText: true,
				},
			},
			CompletionProvider: &CompletionOptions{
				TriggerCharacters: []string{"[", "@", "!", "?", " ", "(", "/"},
				ResolveProvider:   false,
			},
			HoverProvider:      true,
			DefinitionProvider: true,
		},
		ServerInfo: &ServerInfo{
			Name:    "markata-go-lsp",
			Version: "0.1.0",
		},
	}

	s.initialized = true
	return s.sendResponse(msg.ID, result)
}

// handleInitialized handles the initialized notification.
func (s *Server) handleInitialized(_ context.Context, _ *Message) error {
	s.initialized = true
	s.logger.Println("Server initialized")

	// Build the index
	if s.rootURI != "" && s.sharedIndex == nil {
		rootPath := uriToPath(s.rootURI)
		s.logger.Printf("Building index from: %s", rootPath)
		if err := s.index.Build(rootPath); err != nil {
			s.logger.Printf("Failed to build index: %v", err)
		} else {
			posts := s.index.AllPosts()
			s.logger.Printf("Indexed %d posts", len(posts))
		}
	}

	return nil
}

// handleShutdown handles the shutdown request.
func (s *Server) handleShutdown(_ context.Context, msg *Message) error {
	s.shutdown = true
	s.logger.Println("Shutting down")
	return s.sendResponse(msg.ID, nil)
}

// handleExit handles the exit notification.
func (s *Server) handleExit(_ context.Context, _ *Message) error {
	// Exit is handled by the main loop checking s.shutdown
	return nil
}

// handleDidOpen handles textDocument/didOpen notification.
func (s *Server) handleDidOpen(_ context.Context, msg *Message) error {
	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.logger.Printf("Failed to parse didOpen params: %v", err)
		return nil
	}

	// Only track markdown files
	if !strings.HasSuffix(strings.ToLower(params.TextDocument.URI), ".md") {
		return nil
	}

	// Store the document
	doc := &Document{
		URI:     params.TextDocument.URI,
		Content: params.TextDocument.Text,
		Version: params.TextDocument.Version,
	}

	s.docMu.Lock()
	s.documents[params.TextDocument.URI] = doc
	s.docMu.Unlock()

	// Update index and publish diagnostics
	if err := s.index.Update(params.TextDocument.URI, params.TextDocument.Text); err != nil {
		s.logger.Printf("Failed to update index: %v", err)
	}

	return s.publishDiagnostics(params.TextDocument.URI, params.TextDocument.Text)
}

// handleDidChange handles textDocument/didChange notification.
func (s *Server) handleDidChange(_ context.Context, msg *Message) error {
	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.logger.Printf("Failed to parse didChange params: %v", err)
		return nil
	}

	s.docMu.Lock()
	doc, ok := s.documents[params.TextDocument.URI]
	if !ok {
		doc = &Document{URI: params.TextDocument.URI}
		s.documents[params.TextDocument.URI] = doc
	}

	// Apply changes (we use full sync, so just take the last content)
	for _, change := range params.ContentChanges {
		doc.Content = change.Text
	}
	doc.Version = params.TextDocument.Version
	s.docMu.Unlock()

	// Update index and publish diagnostics
	if err := s.index.Update(params.TextDocument.URI, doc.Content); err != nil {
		s.logger.Printf("Failed to update index: %v", err)
	}

	return s.publishDiagnostics(params.TextDocument.URI, doc.Content)
}

// handleDidClose handles textDocument/didClose notification.
func (s *Server) handleDidClose(_ context.Context, msg *Message) error {
	var params DidCloseTextDocumentParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.logger.Printf("Failed to parse didClose params: %v", err)
		return nil
	}

	s.docMu.Lock()
	delete(s.documents, params.TextDocument.URI)
	s.docMu.Unlock()

	// Clear diagnostics
	if s.sharedIndex != nil {
		s.sharedLoaded = false
	}
	return s.clearDiagnostics(params.TextDocument.URI)
}

// handleDidSave handles textDocument/didSave notification.
func (s *Server) handleDidSave(_ context.Context, msg *Message) error {
	var params DidSaveTextDocumentParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.logger.Printf("Failed to parse didSave params: %v", err)
		return nil
	}

	// Re-index on save
	if params.Text != nil {
		if err := s.index.Update(params.TextDocument.URI, *params.Text); err != nil {
			s.logger.Printf("Failed to update index on save: %v", err)
		}

		// Publish diagnostics for all open documents
		// (a save might fix broken links in other files)
		s.docMu.RLock()
		docs := make([]*Document, 0, len(s.documents))
		for _, doc := range s.documents {
			docs = append(docs, doc)
		}
		s.docMu.RUnlock()

		for _, doc := range docs {
			if err := s.publishDiagnostics(doc.URI, doc.Content); err != nil {
				s.logger.Printf("Failed to publish diagnostics: %v", err)
			}
		}
	}

	return nil
}

// handleDidChangeWatchedFiles handles workspace/didChangeWatchedFiles notification.
func (s *Server) handleDidChangeWatchedFiles(_ context.Context, msg *Message) error {
	var params DidChangeWatchedFilesParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		s.logger.Printf("Failed to parse didChangeWatchedFiles params: %v", err)
		return nil
	}

	for _, change := range params.Changes {
		// Only handle markdown files
		if !strings.HasSuffix(strings.ToLower(change.URI), ".md") {
			continue
		}

		switch change.Type {
		case FileChangeTypeCreated, FileChangeTypeChanged:
			// Re-index the file
			path := uriToPath(change.URI)
			if err := s.index.indexFile(path); err != nil {
				s.logger.Printf("Failed to re-index %s: %v", path, err)
			}
		case FileChangeTypeDeleted:
			// Remove from index
			s.index.Remove(change.URI)
		}
	}

	// Republish diagnostics for all open documents
	s.docMu.RLock()
	docs := make([]*Document, 0, len(s.documents))
	for _, doc := range s.documents {
		docs = append(docs, doc)
	}
	s.docMu.RUnlock()

	for _, doc := range docs {
		if err := s.publishDiagnostics(doc.URI, doc.Content); err != nil {
			s.logger.Printf("Failed to publish diagnostics: %v", err)
		}
	}

	return nil
}
