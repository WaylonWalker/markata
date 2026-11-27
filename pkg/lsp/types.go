package lsp

// LSP Protocol Types
// These types follow the Language Server Protocol specification.
// See: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/

// Position represents a position in a text document.
type Position struct {
	Line      int `json:"line"`      // 0-based line number
	Character int `json:"character"` // 0-based character offset
}

// Range represents a range in a text document.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location represents a location inside a resource.
type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// TextDocumentIdentifier identifies a text document.
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// VersionedTextDocumentIdentifier identifies a specific version of a text document.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// TextDocumentItem represents a text document item.
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// TextEdit represents a text edit operation.
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// MarkupContent represents content with optional markup.
type MarkupContent struct {
	Kind  string `json:"kind"` // "plaintext" or "markdown"
	Value string `json:"value"`
}

// Diagnostic represents a diagnostic message.
type Diagnostic struct {
	Range              Range         `json:"range"`
	Severity           int           `json:"severity,omitempty"`
	Code               interface{}   `json:"code,omitempty"`
	CodeDescription    *CodeDesc     `json:"codeDescription,omitempty"`
	Source             string        `json:"source,omitempty"`
	Message            string        `json:"message"`
	Tags               []int         `json:"tags,omitempty"`
	RelatedInformation []DiagRelated `json:"relatedInformation,omitempty"`
	Data               interface{}   `json:"data,omitempty"`
}

// CodeDesc provides additional information about a diagnostic code.
type CodeDesc struct {
	Href string `json:"href"`
}

// DiagRelated represents related diagnostic information.
type DiagRelated struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// DiagnosticSeverity constants.
const (
	DiagnosticSeverityError       = 1
	DiagnosticSeverityWarning     = 2
	DiagnosticSeverityInformation = 3
	DiagnosticSeverityHint        = 4
)

// PublishDiagnosticsParams contains the parameters for publishDiagnostics.
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     *int         `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// InitializeParams contains the parameters for initialize.
type InitializeParams struct {
	ProcessID        *int               `json:"processId"`
	RootURI          *string            `json:"rootUri"`
	RootPath         *string            `json:"rootPath"`
	Capabilities     ClientCapabilities `json:"capabilities"`
	Trace            string             `json:"trace,omitempty"`
	WorkspaceFolders []WorkspaceFolder  `json:"workspaceFolders,omitempty"`
}

// WorkspaceFolder represents a workspace folder.
type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// ClientCapabilities represents the client's capabilities.
type ClientCapabilities struct {
	TextDocument TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Workspace    WorkspaceClientCapabilities    `json:"workspace,omitempty"`
}

// TextDocumentClientCapabilities represents text document client capabilities.
type TextDocumentClientCapabilities struct {
	Completion CompletionClientCapabilities `json:"completion,omitempty"`
	Hover      HoverClientCapabilities      `json:"hover,omitempty"`
}

// CompletionClientCapabilities represents completion client capabilities.
type CompletionClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// HoverClientCapabilities represents hover client capabilities.
type HoverClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// WorkspaceClientCapabilities represents workspace client capabilities.
type WorkspaceClientCapabilities struct {
	DidChangeWatchedFiles DidChangeWatchedFilesClientCapabilities `json:"didChangeWatchedFiles,omitempty"`
}

// DidChangeWatchedFilesClientCapabilities represents file watching capabilities.
type DidChangeWatchedFilesClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// InitializeResult is the result of the initialize request.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

// ServerInfo provides information about the server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// ServerCapabilities represents the server's capabilities.
type ServerCapabilities struct {
	TextDocumentSync   *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
	CompletionProvider *CompletionOptions       `json:"completionProvider,omitempty"`
	HoverProvider      bool                     `json:"hoverProvider,omitempty"`
	DefinitionProvider bool                     `json:"definitionProvider,omitempty"`
	CodeActionProvider bool                    `json:"codeActionProvider,omitempty"`
	ExecuteCommandProvider *ExecuteCommandOptions `json:"executeCommandProvider,omitempty"`
	Workspace          *WorkspaceOptions        `json:"workspace,omitempty"`
}

// ExecuteCommandOptions advertises native prose ignore actions.
type ExecuteCommandOptions struct {
	Commands []string `json:"commands"`
}

// TextDocumentSyncOptions represents text document sync options.
type TextDocumentSyncOptions struct {
	OpenClose bool         `json:"openClose,omitempty"`
	Change    int          `json:"change,omitempty"` // 0=None, 1=Full, 2=Incremental
	Save      *SaveOptions `json:"save,omitempty"`
}

// SaveOptions represents save options.
type SaveOptions struct {
	IncludeText bool `json:"includeText,omitempty"`
}

// CompletionOptions represents completion options.
type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
}

// WorkspaceOptions represents workspace options.
type WorkspaceOptions struct {
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
}

// WorkspaceFoldersServerCapabilities represents workspace folder capabilities.
type WorkspaceFoldersServerCapabilities struct {
	Supported           bool `json:"supported,omitempty"`
	ChangeNotifications bool `json:"changeNotifications,omitempty"`
}

// TextDocumentSyncKind constants.
const (
	TextDocumentSyncKindNone        = 0
	TextDocumentSyncKindFull        = 1
	TextDocumentSyncKindIncremental = 2
)

// DidOpenTextDocumentParams contains the parameters for textDocument/didOpen.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidChangeTextDocumentParams contains the parameters for textDocument/didChange.
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// TextDocumentContentChangeEvent represents a text document content change event.
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidCloseTextDocumentParams contains the parameters for textDocument/didClose.
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DidSaveTextDocumentParams contains the parameters for textDocument/didSave.
type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Text         *string                `json:"text,omitempty"`
}

// DidChangeWatchedFilesParams contains the parameters for workspace/didChangeWatchedFiles.
type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

// FileEvent represents a file event.
type FileEvent struct {
	URI  string `json:"uri"`
	Type int    `json:"type"` // 1=Created, 2=Changed, 3=Deleted
}

// FileChangeType constants.
const (
	FileChangeTypeCreated = 1
	FileChangeTypeChanged = 2
	FileChangeTypeDeleted = 3
)

// HoverParams contains the parameters for textDocument/hover.
type HoverParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// Hover represents hover information.
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// DefinitionParams contains the parameters for textDocument/definition.
type DefinitionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}
