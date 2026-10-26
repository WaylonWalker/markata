package plugins

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestParsePostFromContent_PrivateOverride(t *testing.T) {
	public, err := ParsePostFromContent("public.md", "---\nprivate: false\n---\nbody\n")
	if err != nil {
		t.Fatalf("parse explicit false: %v", err)
	}
	if !public.IsExplicitlyPublic() || public.Private {
		t.Fatal("explicit private: false was not preserved")
	}

	private, err := ParsePostFromContent("private.md", "---\nprivate: true\n---\nbody\n")
	if err != nil {
		t.Fatalf("parse explicit true: %v", err)
	}
	if private.IsExplicitlyPublic() || !private.Private {
		t.Fatal("explicit private: true was not preserved")
	}
}

func TestParsePostFromContent_InvalidPrivateValue(t *testing.T) {
	for _, value := range []string{"\"false\"", "0", "null"} {
		_, err := ParsePostFromContent("invalid.md", "---\nprivate: "+value+"\n---\nbody\n")
		if err == nil {
			t.Errorf("private: %s should fail parsing", value)
		}
	}
}

func TestComputePostFeedItemHash_ChangesWithPrivateState(t *testing.T) {
	post := &models.Post{
		Slug:    "post",
		Href:    "/post/",
		Content: "same content",
	}

	publicHash := computePostFeedItemHash(post)
	post.Private = true
	privateHash := computePostFeedItemHash(post)

	if publicHash == privateHash {
		t.Fatal("feed item hash should change when private state changes")
	}
}

func TestLoadPlugin_ReservesPrivateMetadataMarkers(t *testing.T) {
	loader := NewLoadPlugin()
	post := models.NewPost("post.md")
	if err := loader.applyMetadata(post, map[string]interface{}{
		"_title_explicit":       true,
		"_description_explicit": true,
	}); err != nil {
		t.Fatalf("applyMetadata() error = %v", err)
	}
	if post.Has("_title_explicit") || post.Has("_description_explicit") {
		t.Fatal("reserved metadata markers were copied from frontmatter")
	}

	post = models.NewPost("post.md")
	if err := loader.applyMetadata(post, map[string]interface{}{
		"title":       "Authored title",
		"description": "Authored description",
	}); err != nil {
		t.Fatalf("applyMetadata() error = %v", err)
	}
	if post.Get("_title_explicit") != true || post.Get("_description_explicit") != true {
		t.Fatal("authored metadata did not set internal provenance markers")
	}
}

func TestLoadPlugin_RestoresCachedFrontmatterDiagnostics(t *testing.T) {
	loader := NewLoadPlugin()
	post := models.NewPost("post.md")
	post.Published = true
	post.Set("_frontmatter_present", true)
	post.Set("_frontmatter_valid", true)
	issue := diagnostics.Issue{
		File:     "post.md",
		Code:     diagnostics.ReasonFrontmatterSuspiciousDelimiter,
		Severity: diagnostics.SeverityWarning,
		Message:  "suspicious delimiter",
	}

	cached := loader.postToCachedData(post, []diagnostics.Issue{issue})
	if len(cached.FrontmatterIssues) != 1 {
		t.Fatalf("cached frontmatter issues = %d, want 1", len(cached.FrontmatterIssues))
	}
	encoded, err := json.Marshal(cached)
	if err != nil {
		t.Fatalf("marshal cached post = %v", err)
	}
	var decoded buildcache.CachedPostData
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal cached post = %v", err)
	}

	restored := loader.restorePostFromCache(&decoded, nil)
	m := lifecycle.NewManager()
	m.SetFiles([]string{"post.md"})
	recordLoadedPost(m, restored)

	snapshot := m.ContentDiagnostics()
	if len(snapshot.Entries) != 1 || len(snapshot.Entries[0].Diagnostics) != 1 {
		t.Fatalf("restored diagnostics = %+v, want one issue", snapshot.Entries)
	}
	if snapshot.Entries[0].Diagnostics[0].Code != diagnostics.ReasonFrontmatterSuspiciousDelimiter {
		t.Fatalf("restored diagnostics = %+v", snapshot.Entries[0].Diagnostics)
	}
}

func TestLoadPlugin_ContinuesForRecoverableFrontmatterErrors(t *testing.T) {
	contentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "good.md"), []byte("---\ntitle: Good\npublished: true\n---\nGood"), 0o600); err != nil {
		t.Fatalf("write good.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "broken.md"), []byte("---\ntitle: [broken\npublished: true\n---\nBroken"), 0o600); err != nil {
		t.Fatalf("write broken.md: %v", err)
	}

	loader := NewLoadPlugin()
	m := lifecycle.NewManager()
	m.Config().ContentDir = contentDir
	m.SetFiles([]string{"broken.md", "good.md"})

	if _, err := loader.parseFile("broken.md", "---\ntitle: [broken\npublished: true\n---\nBroken"); !errors.Is(err, ErrRecoverableContent) {
		t.Fatalf("parseFile() error = %v, want ErrRecoverableContent", err)
	}
	if _, err := loader.parseFile("broken.md", "---\ntitle: [broken\npublished: true\n---\nBroken"); !errors.Is(err, ErrInvalidFrontmatter) {
		t.Fatalf("parseFile() error = %v, want ErrInvalidFrontmatter", err)
	}

	if err := loader.Load(m); err != nil {
		t.Fatalf("Load() error = %v, want nil for recoverable content error", err)
	}
	posts := m.Posts()
	if len(posts) != 1 || posts[0].Path != "good.md" {
		t.Fatalf("loaded posts = %+v, want only good.md", posts)
	}

	snapshot := m.ContentDiagnostics()
	entry := snapshot.Entries[0]
	if entry.Path != "broken.md" || !entry.FrontmatterPresent || entry.FrontmatterValid {
		t.Fatalf("broken frontmatter state = %+v", entry)
	}
	if !containsDiagnosticCode(entry.Diagnostics, diagnostics.ReasonFrontmatterParseError) {
		t.Fatalf("broken diagnostics = %+v", entry.Diagnostics)
	}
	if containsReason(entry.Reasons, diagnostics.ReasonContentLoadError) {
		t.Fatalf("recoverable content received load error: %v", entry.Reasons)
	}
}

func TestLoadPlugin_PropagatesOperationalSourceErrors(t *testing.T) {
	tests := []struct {
		name       string
		brokenPath string
		makeBroken func(t *testing.T, path string)
	}{
		{
			name:       "stat failure",
			brokenPath: "missing.md",
			makeBroken: func(*testing.T, string) {},
		},
		{
			name:       "read failure",
			brokenPath: "directory.md",
			makeBroken: func(t *testing.T, path string) {
				t.Helper()
				if err := os.Mkdir(path, 0o700); err != nil {
					t.Fatalf("mkdir directory.md: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentDir := t.TempDir()
			if err := os.WriteFile(filepath.Join(contentDir, "good.md"), []byte("---\ntitle: Good\npublished: true\n---\nGood"), 0o600); err != nil {
				t.Fatalf("write good.md: %v", err)
			}
			tt.makeBroken(t, filepath.Join(contentDir, tt.brokenPath))

			m := lifecycle.NewManager()
			m.Config().ContentDir = contentDir
			m.SetFiles([]string{"good.md", tt.brokenPath})

			if err := NewLoadPlugin().Load(m); err == nil {
				t.Fatal("Load() returned nil for operational source failure")
			}
			if posts := m.Posts(); len(posts) != 1 || posts[0].Path != "good.md" {
				t.Fatalf("loaded posts = %+v, want valid sibling retained", posts)
			}

			entry := findContentDisposition(m.ContentDiagnostics(), tt.brokenPath)
			if entry == nil {
				t.Fatalf("missing ledger entry for %s", tt.brokenPath)
			}
			if !containsDiagnosticCode(entry.Diagnostics, diagnostics.ReasonContentLoadError) {
				t.Fatalf("operational diagnostics = %+v", entry.Diagnostics)
			}
		})
	}
}

func containsDiagnosticCode(issues []diagnostics.Issue, code string) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func containsReason(reasons []string, want string) bool {
	for _, reason := range reasons {
		if reason == want {
			return true
		}
	}
	return false
}

func findContentDisposition(snapshot diagnostics.ContentLedgerSnapshot, path string) *diagnostics.ContentDisposition {
	for index := range snapshot.Entries {
		if snapshot.Entries[index].Path == path {
			return &snapshot.Entries[index]
		}
	}
	return nil
}

func TestResolveFileModTime_UsesGlobCachedModTime(t *testing.T) {
	m := lifecycle.NewManager()
	m.Cache().Set(cacheKeyGlobFileModTimes, map[string]int64{"post.md": 12345})

	got, err := resolveFileModTime(m, "post.md", "/does/not/exist")
	if err != nil {
		t.Fatalf("resolveFileModTime() error = %v", err)
	}
	if got != 12345 {
		t.Fatalf("resolveFileModTime() = %d, want 12345", got)
	}
}

func TestResolveFileModTime_FallsBackToStat(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "post.md")
	if err := os.WriteFile(path, []byte("# post"), 0o600); err != nil {
		t.Fatalf("WriteFile(post.md) error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(post.md) error = %v", err)
	}

	got, err := resolveFileModTime(nil, "post.md", tmpDir)
	if err != nil {
		t.Fatalf("resolveFileModTime() error = %v", err)
	}
	if got != info.ModTime().UnixNano() {
		t.Fatalf("resolveFileModTime() = %d, want %d", got, info.ModTime().UnixNano())
	}
}

func TestLoadPlugin_MarksCachedPostsMissing(t *testing.T) {
	contentDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(contentDir, "keep.md"), []byte("---\ntitle: Keep\n---\nkeep\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "source.md"), []byte("---\ntitle: Source\n---\nsource\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	m := lifecycle.NewManager()
	m.Config().ContentDir = contentDir
	m.Config().GlobPatterns = []string{"**/*.md"}
	m.SetFiles([]string{"keep.md", "source.md"})
	cache := buildcache.New(filepath.Join(t.TempDir(), ".markata"))
	cache.UpdateModTime("keep.md", 1, "keep")
	cache.UpdateModTime("source.md", 1, "source")
	cache.UpdateModTime("deleted.md", 1, "deleted")
	cache.SetDependencies("source.md", "source", []string{"deleted"})
	m.Cache().Set("build_cache", cache)

	if err := NewLoadPlugin().Load(m); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	removed := lifecycle.GetServeRemovedPaths(m)
	if len(removed) != 1 || removed[0] != "deleted.md" {
		t.Fatalf("removed paths = %v, want [deleted.md]", removed)
	}
	if affected := lifecycle.GetServeAffectedPaths(m); !affected["source.md"] {
		t.Fatalf("affected paths = %v, want source.md", affected)
	}
}

func TestParseDateString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		// Standard formats
		{
			name:  "RFC3339",
			input: "2024-01-15T10:30:00Z",
			want:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "ISO datetime with T",
			input: "2024-01-15T10:30:00",
			want:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "datetime with space",
			input: "2024-01-15 10:30:00",
			want:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "date only",
			input: "2024-01-15",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},

		// Single-digit hours (issue #34)
		{
			name:  "single-digit hour",
			input: "2025-02-08 1:00:00",
			want:  time.Date(2025, 2, 8, 1, 0, 0, 0, time.UTC),
		},
		{
			name:  "single-digit hour with T",
			input: "2025-02-08T1:00:00",
			want:  time.Date(2025, 2, 8, 1, 0, 0, 0, time.UTC),
		},
		{
			name:  "single-digit hour 9am",
			input: "2024-06-20 9:30:00",
			want:  time.Date(2024, 6, 20, 9, 30, 0, 0, time.UTC),
		},

		// Malformed time components (issue #34)
		{
			name:  "malformed time with extra zero",
			input: "2025-07-14 8:011:00",
			want:  time.Date(2025, 7, 14, 8, 11, 0, 0, time.UTC),
		},
		{
			name:  "malformed time multiple extra zeros",
			input: "2025-07-14 08:001:030",
			want:  time.Date(2025, 7, 14, 8, 1, 30, 0, time.UTC),
		},

		// Date without time
		{
			name:  "date with slashes",
			input: "2024/01/15",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "US date format",
			input: "01/15/2024",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},

		// Named month formats
		{
			name:  "full month name",
			input: "January 15, 2024",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "abbreviated month name",
			input: "Jan 15, 2024",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "day first with full month",
			input: "15 January 2024",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},

		// Datetime with slashes
		{
			name:  "datetime with slashes",
			input: "2024/01/15 10:30:00",
			want:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},

		// Without seconds
		{
			name:  "datetime without seconds",
			input: "2024-01-15 10:30",
			want:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "datetime with T without seconds",
			input: "2024-01-15T10:30",
			want:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},

		// Whitespace handling
		{
			name:  "leading whitespace",
			input: "  2024-01-15",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "trailing whitespace",
			input: "2024-01-15  ",
			want:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},

		// Error cases
		{
			name:    "invalid date",
			input:   "not a date",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDateString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDateString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("parseDateString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeDateString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "already normalized",
			input: "2024-01-15 10:30:00",
			want:  "2024-01-15 10:30:00",
		},
		{
			name:  "single-digit hour with space",
			input: "2024-01-15 1:30:00",
			want:  "2024-01-15 01:30:00",
		},
		{
			name:  "single-digit hour with T",
			input: "2024-01-15T1:30:00",
			want:  "2024-01-15T01:30:00",
		},
		{
			name:  "malformed minutes",
			input: "2024-01-15 8:011:00",
			want:  "2024-01-15 08:11:00",
		},
		{
			name:  "whitespace trimmed",
			input: "  2024-01-15 10:30:00  ",
			want:  "2024-01-15 10:30:00",
		},
		{
			name:  "no time component",
			input: "2024-01-15",
			want:  "2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeDateString(tt.input)
			if got != tt.want {
				t.Errorf("normalizeDateString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestApplyMetadata_AuthorAliases(t *testing.T) {
	p := &LoadPlugin{}

	tests := []struct {
		name       string
		metadata   map[string]interface{}
		wantAuthor string
	}{
		{
			name:       "author field takes priority",
			metadata:   map[string]interface{}{"author": "alice", "by": "bob", "writer": "charlie"},
			wantAuthor: "alice",
		},
		{
			name:       "by alias resolves to author",
			metadata:   map[string]interface{}{"by": "bob"},
			wantAuthor: "bob",
		},
		{
			name:       "writer alias resolves to author",
			metadata:   map[string]interface{}{"writer": "charlie"},
			wantAuthor: "charlie",
		},
		{
			name:       "by takes priority over writer",
			metadata:   map[string]interface{}{"by": "bob", "writer": "charlie"},
			wantAuthor: "bob",
		},
		{
			name:       "no author fields leaves author nil",
			metadata:   map[string]interface{}{"title": "some post"},
			wantAuthor: "",
		},
		{
			name:       "by and writer not stored in Extra",
			metadata:   map[string]interface{}{"by": "bob"},
			wantAuthor: "bob",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{}
			if err := p.applyMetadata(post, tt.metadata); err != nil {
				t.Fatalf("applyMetadata() error = %v", err)
			}

			if tt.wantAuthor == "" {
				if post.Author != nil {
					t.Errorf("expected Author to be nil, got %q", *post.Author)
				}
			} else {
				if post.Author == nil {
					t.Fatalf("expected Author to be %q, got nil", tt.wantAuthor)
				}
				if *post.Author != tt.wantAuthor {
					t.Errorf("Author = %q, want %q", *post.Author, tt.wantAuthor)
				}
			}

			// Verify aliases are not stored in Extra
			if tt.name == "by and writer not stored in Extra" {
				if val := post.Get("by"); val != nil {
					t.Errorf("'by' should not be in Extra, got %v", val)
				}
				if val := post.Get("writer"); val != nil {
					t.Errorf("'writer' should not be in Extra, got %v", val)
				}
			}
		})
	}
}

func TestApplyMetadata_DateAliases(t *testing.T) {
	p := &LoadPlugin{}

	tests := []struct {
		name        string
		metadata    map[string]interface{}
		wantDate    time.Time
		wantMod     time.Time
		wantErr     bool
		wantDateSet bool
		wantModSet  bool
	}{
		{
			name:        "publishdate wins over date and pubdate",
			metadata:    map[string]interface{}{"date": "2024-01-01", "publishdate": "2024-02-01", "pubdate": "2024-03-01"},
			wantDate:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			wantDateSet: true,
		},
		{
			name:        "date used when publishdate absent",
			metadata:    map[string]interface{}{"date": "2024-01-15"},
			wantDate:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			wantDateSet: true,
		},
		{
			name:        "pubdate used when others absent",
			metadata:    map[string]interface{}{"pubdate": "2024-03-10"},
			wantDate:    time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
			wantDateSet: true,
		},
		{
			name:       "lastmod wins over modified and updated",
			metadata:   map[string]interface{}{"modified": "2024-01-01", "lastmod": "2024-02-01", "updated": "2024-03-01"},
			wantMod:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			wantModSet: true,
		},
		{
			name:       "updated_at used when higher precedence absent",
			metadata:   map[string]interface{}{"updated_at": "2024-04-05"},
			wantMod:    time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC),
			wantModSet: true,
		},
		{
			name:     "invalid date returns error",
			metadata: map[string]interface{}{"date": "not-a-date"},
			wantErr:  true,
		},
		{
			name:     "invalid modified returns error",
			metadata: map[string]interface{}{"lastmod": "not-a-date"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{}
			err := p.applyMetadata(post, tt.metadata)
			if (err != nil) != tt.wantErr {
				t.Fatalf("applyMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.wantDateSet {
				if post.Date == nil {
					t.Fatalf("expected Date to be set")
				}
				if !post.Date.Equal(tt.wantDate) {
					t.Errorf("Date = %v, want %v", post.Date, tt.wantDate)
				}
			} else if post.Date != nil {
				t.Errorf("expected Date to be nil")
			}
			if tt.wantModSet {
				if post.Modified == nil {
					t.Fatalf("expected Modified to be set")
				}
				if !post.Modified.Equal(tt.wantMod) {
					t.Errorf("Modified = %v, want %v", post.Modified, tt.wantMod)
				}
			} else if post.Modified != nil {
				t.Errorf("expected Modified to be nil")
			}
		})
	}
}
