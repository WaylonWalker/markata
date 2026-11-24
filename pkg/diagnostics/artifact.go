package diagnostics

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"
)

const (
	// ArtifactSchemaURL identifies the versioned diagnostics wire format.
	ArtifactSchemaURL = "markata://schemas/content-diagnostics/v1"

	// ArtifactSchema is the stable diagnostics artifact identity.
	ArtifactSchema = "markata.content-diagnostics"

	// ArtifactSchemaVersion is the current diagnostics artifact generation.
	ArtifactSchemaVersion = 1

	// DefaultArtifactPath is relative to a build's output directory.
	DefaultArtifactPath = ".markata/diagnostics.json"

	// CompressedArtifactPath is the losslessly compressed artifact path used by
	// Builder Admin releases. The decoded JSON has the same bytes and schema as
	// DefaultArtifactPath.
	CompressedArtifactPath = DefaultArtifactPath + ".gz"

	// ArtifactExecutorLegacy identifies the normal lifecycle build path.
	ArtifactExecutorLegacy = "legacy"
	// ArtifactExecutorDAG identifies the feature-flagged task-graph build path.
	ArtifactExecutorDAG = "dag"
)

// Artifact describes the content state observed during one successful build.
// Its content fields are copied from a ContentLedgerSnapshot; they are not
// reconstructed from posts, feeds, or plugin-local counters.
type Artifact struct {
	SchemaURL     string               `json:"$schema"`
	Schema        string               `json:"schema"`
	SchemaVersion int                  `json:"schema_version"`
	Generator     ArtifactGenerator    `json:"generator"`
	Source        *ArtifactSource      `json:"source,omitempty"`
	BuiltAt       time.Time            `json:"built_at"`
	Executor      string               `json:"executor,omitempty"`
	Summary       ContentSummary       `json:"summary"`
	TemplateCache *TemplateCacheStats  `json:"template_cache,omitempty"`
	PublishIO     *PublishIOStats      `json:"publish_io,omitempty"`
	Entries       []ContentDisposition `json:"entries"`
}

// ArtifactGenerator identifies the markata-go binary that produced an
// artifact. Commit is omitted when the binary was not built with a reliable
// source revision.
type ArtifactGenerator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
}

// ArtifactSource identifies the source Git revision used for the build.
type ArtifactSource struct {
	Commit string `json:"commit"`
}

// ArtifactBuildInfo supplies build identity values for an artifact.
type ArtifactBuildInfo struct {
	MarkataVersion string
	MarkataCommit  string
	SourceCommit   string
	BuiltAt        time.Time
	Executor       string
}

// NewArtifact creates a versioned artifact from a deterministic ledger
// snapshot. A zero BuiltAt is replaced with the current UTC time.
func NewArtifact(snapshot ContentLedgerSnapshot, info ArtifactBuildInfo) Artifact {
	builtAt := info.BuiltAt
	if builtAt.IsZero() {
		builtAt = time.Now().UTC()
	} else {
		builtAt = builtAt.UTC()
	}

	version := strings.TrimSpace(info.MarkataVersion)
	if version == "" {
		version = "dev"
	}

	artifact := Artifact{
		SchemaURL:     ArtifactSchemaURL,
		Schema:        ArtifactSchema,
		SchemaVersion: ArtifactSchemaVersion,
		Generator: ArtifactGenerator{
			Name:    "markata-go",
			Version: version,
			Commit:  reliableCommit(info.MarkataCommit),
		},
		BuiltAt:       builtAt,
		Executor:      strings.TrimSpace(info.Executor),
		Summary:       snapshot.Summary,
		TemplateCache: cloneTemplateCacheStats(snapshot.TemplateCache),
		PublishIO:     clonePublishIOStats(snapshot.PublishIO),
		Entries:       cloneArtifactEntries(snapshot.Entries),
	}
	if commit := reliableCommit(info.SourceCommit); commit != "" {
		artifact.Source = &ArtifactSource{Commit: commit}
	}
	return artifact
}

// MarshalArtifact encodes a diagnostics artifact as indented, stable JSON.
func MarshalArtifact(snapshot ContentLedgerSnapshot, info ArtifactBuildInfo) ([]byte, error) {
	artifact := NewArtifact(snapshot, info)
	if err := validateArtifact(artifact); err != nil {
		return nil, err
	}
	return json.MarshalIndent(artifact, "", "  ")
}

// WriteArtifact streams the same indented JSON as MarshalArtifact, without a
// trailing newline or modifying snapshot. Temporary JSON buffers and sanitized
// copies are limited to one entry, plus a capped per-call feed fragment cache;
// the ordering index scales with entry count.
// Use a buffered writer when publishing to a file, and flush it before syncing.
func WriteArtifact(writer io.Writer, snapshot ContentLedgerSnapshot, info ArtifactBuildInfo) error {
	header := NewArtifact(ContentLedgerSnapshot{Summary: snapshot.Summary, TemplateCache: snapshot.TemplateCache, PublishIO: snapshot.PublishIO}, info)
	if err := validateArtifact(header); err != nil {
		return err
	}
	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return err
	}
	if len(snapshot.Entries) == 0 {
		return writeArtifactBytes(writer, data)
	}

	// Entries is the final field of Artifact. Verify the generated suffix rather
	// than silently corrupting output if that layout changes in the future.
	const emptyEntriesSuffix = "  \"entries\": []\n}"
	if !bytes.HasSuffix(data, []byte(emptyEntriesSuffix)) {
		return fmt.Errorf("diagnostics artifact header must end with empty entries")
	}
	if err := writeArtifactBytes(writer, data[:len(data)-len("]\n}")]); err != nil {
		return err
	}

	type entryIndex struct {
		path  string
		index int
	}
	order := make([]entryIndex, len(snapshot.Entries))
	for index, entry := range snapshot.Entries {
		order[index] = entryIndex{path: normalizeContentPath(entry.Path), index: index}
	}
	sort.SliceStable(order, func(i, j int) bool {
		return order[i].path < order[j].path
	})
	// Reuse both buffers across entries. Encoder supplies the same escaping as
	// Marshal, while Indent supplies exactly MarshalIndent's layout. Neither
	// buffer grows with the total number of entries.
	var compact, indented bytes.Buffer
	encoder := json.NewEncoder(&compact)
	var scratch artifactEntryScratch
	var fragments artifactFeedFragments
	fragmentLayout := artifactFeedLayoutSupported(reflect.TypeFor[ContentDisposition]())
	for index, item := range order {
		entry := scratch.clone(snapshot.Entries[item.index])
		feeds := entry.Feeds
		if fragmentLayout {
			entry.Feeds = nil
		}
		compact.Reset()
		if err := encoder.Encode(entry); err != nil {
			return err
		}
		indented.Reset()
		// Encode appends a newline; MarshalArtifact does not.
		encoded := bytes.TrimSuffix(compact.Bytes(), []byte("\n"))
		if err := json.Indent(&indented, encoded, "    ", "  "); err != nil {
			return err
		}
		if fragmentLayout && len(feeds) > 0 {
			const entryEnd = "\n    }"
			if !bytes.HasSuffix(indented.Bytes(), []byte(entryEnd)) {
				return fmt.Errorf("diagnostics artifact entry must end with object close")
			}
			indented.Truncate(indented.Len() - len(entryEnd))
			indented.WriteString(",\n      \"feeds\": [")
			for feedIndex, feed := range feeds {
				if feedIndex > 0 {
					indented.WriteByte(',')
				}
				indented.WriteString("\n        ")
				fragment, err := fragments.encode(feed)
				if err != nil {
					return err
				}
				indented.Write(fragment)
			}
			indented.WriteString("\n      ]\n    }")
		}
		separator := "\n    "
		if index > 0 {
			separator = ",\n    "
		}
		if err := writeArtifactBytes(writer, []byte(separator)); err != nil {
			return err
		}
		if err := writeArtifactBytes(writer, indented.Bytes()); err != nil {
			return err
		}
	}
	return writeArtifactBytes(writer, []byte("\n  ]\n}"))
}

func writeArtifactBytes(writer io.Writer, data []byte) error {
	written, err := writer.Write(data)
	if err != nil {
		return err
	}
	if written != len(data) {
		return io.ErrShortWrite
	}
	return nil
}

// ParseArtifact decodes and validates a diagnostics artifact.
func ParseArtifact(data []byte) (Artifact, error) {
	var artifact Artifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		return Artifact{}, fmt.Errorf("decode diagnostics artifact JSON: %w", err)
	}
	if err := validateArtifact(artifact); err != nil {
		return Artifact{}, err
	}
	return artifact, nil
}

func validateArtifact(artifact Artifact) error {
	if artifact.SchemaURL != ArtifactSchemaURL {
		return fmt.Errorf("diagnostics artifact $schema must be %q", ArtifactSchemaURL)
	}
	if artifact.Schema != ArtifactSchema {
		return fmt.Errorf("diagnostics artifact schema must be %q", ArtifactSchema)
	}
	if artifact.SchemaVersion != ArtifactSchemaVersion {
		return fmt.Errorf("unsupported diagnostics artifact schema version %d", artifact.SchemaVersion)
	}
	if artifact.Generator.Name == "" {
		return fmt.Errorf("diagnostics artifact generator.name is required")
	}
	if artifact.Generator.Version == "" {
		return fmt.Errorf("diagnostics artifact generator.version is required")
	}
	if artifact.BuiltAt.IsZero() {
		return fmt.Errorf("diagnostics artifact built_at is required")
	}
	if artifact.Executor != "" && artifact.Executor != ArtifactExecutorLegacy && artifact.Executor != ArtifactExecutorDAG {
		return fmt.Errorf("diagnostics artifact executor must be %q or %q", ArtifactExecutorLegacy, ArtifactExecutorDAG)
	}
	if artifact.Source != nil && artifact.Source.Commit == "" {
		return fmt.Errorf("diagnostics artifact source.commit must not be empty")
	}
	if artifact.Entries == nil {
		return fmt.Errorf("diagnostics artifact entries is required")
	}
	return nil
}

func reliableCommit(commit string) string {
	commit = strings.TrimSpace(commit)
	switch commit {
	case "", "none", "unknown", "dev":
		return ""
	default:
		return commit
	}
}

func cloneArtifactEntries(entries []ContentDisposition) []ContentDisposition {
	if entries == nil {
		return []ContentDisposition{}
	}
	result := make([]ContentDisposition, len(entries))
	for index, entry := range entries {
		result[index] = cloneArtifactEntry(entry)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Path < result[j].Path
	})
	return result
}

// cloneArtifactEntry returns owned mutable storage, never shared with scratch
// from another entry. Both encoders use artifactEntryScratch's sanitation.
func cloneArtifactEntry(entry ContentDisposition) ContentDisposition {
	var scratch artifactEntryScratch
	return scratch.clone(entry)
}

type artifactEntryScratch struct {
	reasons     []string
	diagnostics []Issue
	feeds       []ContentFeedDisposition
	feedReasons []string
}

func (s *artifactEntryScratch) clone(entry ContentDisposition) ContentDisposition {
	// Clear old references even when the next entry shrinks. Never append to or
	// sort caller storage, including shared backing slices between feeds.
	clear(s.reasons)
	clear(s.diagnostics)
	clear(s.feeds)
	clear(s.feedReasons)
	entry.Path = normalizeContentPath(entry.Path)
	s.reasons = append(s.reasons[:0], entry.Reasons...)
	entry.Reasons = sortedUnique(s.reasons)
	s.diagnostics = append(s.diagnostics[:0], entry.Diagnostics...)
	entry.Diagnostics = s.diagnostics
	for index := range entry.Diagnostics {
		entry.Diagnostics[index] = sanitizeArtifactIssue(entry.Diagnostics[index])
	}
	entry.Diagnostics = sortedIssues(entry.Diagnostics)
	s.feeds = append(s.feeds[:0], entry.Feeds...)
	reasonCount := 0
	for _, feed := range entry.Feeds {
		reasonCount += len(feed.Reasons)
	}
	s.feedReasons = slices.Grow(s.feedReasons[:0], reasonCount)
	entry.Feeds = s.feeds
	for index := range entry.Feeds {
		start := len(s.feedReasons)
		s.feedReasons = append(s.feedReasons, entry.Feeds[index].Reasons...)
		entry.Feeds[index].Reasons = sortedUnique(s.feedReasons[start:len(s.feedReasons):len(s.feedReasons)])
	}
	slices.SortStableFunc(entry.Feeds, func(a, b ContentFeedDisposition) int {
		return strings.Compare(a.Feed, b.Feed)
	})
	return entry
}

// Feeds must be the last, omittable field for assembly to match MarshalIndent.
// A future custom marshaler or layout change uses whole-entry encoding instead.
func artifactFeedLayoutSupported(t reflect.Type) bool {
	// Explicit v1 layout: embedding, duplicate tags or added/reordered fields
	// require review rather than silently changing the assembly assumptions.
	tags := [...]string{
		"path", "candidate", "loaded", "frontmatter_present", "frontmatter_valid",
		"post_created", "eligible", "rendered", "emitted", "excluded", "disposition",
		"reasons,omitempty", "diagnostics,omitempty", "feeds,omitempty",
	}
	if t.Kind() != reflect.Struct || t.NumField() != len(tags) {
		return false
	}
	marshaler := reflect.TypeFor[json.Marshaler]()
	if t.Implements(marshaler) || reflect.PointerTo(t).Implements(marshaler) {
		return false
	}
	for i, tag := range tags {
		field := t.Field(i)
		if field.Anonymous || field.PkgPath != "" || field.Tag.Get("json") != tag {
			return false
		}
	}
	field := t.Field(t.NumField() - 1)
	return field.Name == "Feeds" && field.Tag.Get("json") == "feeds,omitempty" &&
		field.Type == reflect.TypeFor[[]ContentFeedDisposition]()
}

const (
	artifactFeedCacheEntries = 1024
	artifactFeedCacheBytes   = 1 << 20
)

type artifactFeedFragments struct {
	cache    map[string][]byte
	retained int // key plus fragment bytes; map overhead is bounded by entry cap
	key      []byte
	compact  bytes.Buffer
	indented bytes.Buffer
}

func (c *artifactFeedFragments) encode(feed ContentFeedDisposition) ([]byte, error) {
	// Length-prefixed raw bytes avoid delimiter ambiguity and hashing collisions.
	// nil and empty reasons are identical on the wire (omitempty).
	c.key = binary.AppendUvarint(c.key[:0], uint64(len(feed.Feed)))
	c.key = append(c.key, feed.Feed...)
	if feed.Included {
		c.key = append(c.key, 1)
	} else {
		c.key = append(c.key, 0)
	}
	c.key = binary.AppendUvarint(c.key, uint64(len(feed.Reasons)))
	for _, reason := range feed.Reasons {
		c.key = binary.AppendUvarint(c.key, uint64(len(reason)))
		c.key = append(c.key, reason...)
	}
	if fragment, ok := c.cache[string(c.key)]; ok {
		return fragment, nil
	}
	c.compact.Reset()
	if err := json.NewEncoder(&c.compact).Encode(feed); err != nil {
		return nil, err
	}
	c.indented.Reset()
	if err := json.Indent(&c.indented, bytes.TrimSuffix(c.compact.Bytes(), []byte("\n")), "        ", "  "); err != nil {
		return nil, err
	}
	fragment := c.indented.Bytes()
	size := len(c.key) + len(fragment)
	if len(c.cache) < artifactFeedCacheEntries && size <= artifactFeedCacheBytes-c.retained {
		if c.cache == nil {
			c.cache = make(map[string][]byte)
		}
		fragment = bytes.Clone(fragment)
		c.cache[string(c.key)] = fragment
		c.retained += size
	}
	return fragment, nil
}

func sanitizeArtifactIssue(issue Issue) Issue {
	issue.File = normalizeContentPath(issue.File)
	// Built-in diagnostics use safe, concise messages. Do not copy arbitrary
	// producer-provided text into a public artifact because a third-party plugin
	// could otherwise publish raw content, configuration, or secrets.
	issue.Message = safeArtifactDiagnosticMessage(issue.Code)
	return issue
}

var safeArtifactDiagnosticMessages = map[string]string{
	ReasonFrontmatterSuspiciousDelimiter: "frontmatter opening delimiter is suspicious",
	ReasonFrontmatterLeadingWhitespace:   "frontmatter opening delimiter has leading whitespace",
	ReasonFrontmatterMalformedClosing:    "frontmatter closing delimiter is malformed",
	ReasonFrontmatterMissingClosing:      "frontmatter closing delimiter is missing",
	ReasonFrontmatterParseError:          "frontmatter could not be parsed",
	ReasonFrontmatterDuplicateKey:        "frontmatter contains a duplicate key",
	diagnosticCodeDuplicateKey:           "frontmatter contains a duplicate key",
	ReasonFrontmatterInvalidType:         "frontmatter field has an invalid type",
	ReasonFrontmatterInvalidDate:         "frontmatter date has an invalid format",
	diagnosticCodeInvalidDate:            "frontmatter date has an invalid format",
	ReasonContentPublishedFalse:          "content is not published",
	ReasonContentDraft:                   "content is a draft",
	ReasonContentSkip:                    "content is explicitly skipped",
	ReasonContentPrivate:                 "content is private",
	ReasonContentFiltered:                "content was filtered from the feed",
	ReasonContentDuplicateSlug:           "content has a duplicate slug",
	ReasonContentNoOutput:                "content has no output",
	ReasonContentLoadError:               "content could not be loaded",
	ReasonContentRenderError:             "content could not be rendered",
	ReasonContentWriteError:              "content output could not be written",
	ReasonFeedOffset:                     "content is outside the feed offset",
	ReasonFeedLimit:                      "content is outside the feed limit",
	diagnosticCodeMissingAltText:         "image link is missing alt text",
	diagnosticCodeProtocolLessURL:        "URL is missing a protocol",
	diagnosticCodeH1InContent:            "content contains an H1 heading",
	diagnosticCodeAdmonitionFence:        "fenced code follows an admonition without a blank line",
	diagnosticCodeBrokenWikilink:         "wikilink target was not found",
	diagnosticCodeUnknownMention:         "mention target was not found",
}

func safeArtifactDiagnosticMessage(code string) string {
	return safeArtifactDiagnosticMessages[code]
}
