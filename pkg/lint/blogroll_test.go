package lint

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestBlogroll_DuplicateHandles(t *testing.T) {
	tests := []struct {
		name     string
		feeds    []models.ExternalFeedConfig
		wantLen  int
		wantCode string
	}{
		{
			name: "no duplicates",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://example.com/feed.xml", Handle: "example"},
				{URL: "https://other.com/feed.xml", Handle: "other"},
			},
			wantLen: 0,
		},
		{
			name: "duplicate handle",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://example.com/feed.xml", Handle: "example"},
				{URL: "https://other.com/feed.xml", Handle: "example"},
			},
			wantLen:  1,
			wantCode: "LBL001",
		},
		{
			name: "multiple duplicate handles",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://a.com/feed.xml", Handle: "example"},
				{URL: "https://b.com/feed.xml", Handle: "example"},
				{URL: "https://c.com/feed.xml", Handle: "example"},
			},
			wantLen:  2, // 2 duplicates of first
			wantCode: "LBL001",
		},
		{
			name: "empty handles not considered duplicates",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://a.com/feed.xml", Handle: ""},
				{URL: "https://b.com/feed.xml", Handle: ""},
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &models.BlogrollConfig{Feeds: tt.feeds}
			result := Blogroll(config)

			var handleIssues []BlogrollIssue
			for _, issue := range result.Issues {
				if issue.Code == "LBL001" {
					handleIssues = append(handleIssues, issue)
				}
			}

			if len(handleIssues) != tt.wantLen {
				t.Errorf("got %d LBL001 issues, want %d", len(handleIssues), tt.wantLen)
			}

			if tt.wantLen > 0 && handleIssues[0].Code != tt.wantCode {
				t.Errorf("got code %q, want %q", handleIssues[0].Code, tt.wantCode)
			}
		})
	}
}

func TestBlogroll_DuplicateURLs(t *testing.T) {
	tests := []struct {
		name     string
		feeds    []models.ExternalFeedConfig
		wantLen  int
		wantCode string
	}{
		{
			name: "no duplicates",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://example.com/feed.xml"},
				{URL: "https://other.com/feed.xml"},
			},
			wantLen: 0,
		},
		{
			name: "duplicate URL",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://example.com/feed.xml", Handle: "first"},
				{URL: "https://example.com/feed.xml", Handle: "second"},
			},
			wantLen:  1,
			wantCode: "LBL002",
		},
		{
			name: "empty URLs not considered duplicates",
			feeds: []models.ExternalFeedConfig{
				{URL: "", Handle: "first"},
				{URL: "", Handle: "second"},
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &models.BlogrollConfig{Feeds: tt.feeds}
			result := Blogroll(config)

			var urlIssues []BlogrollIssue
			for _, issue := range result.Issues {
				if issue.Code == "LBL002" {
					urlIssues = append(urlIssues, issue)
				}
			}

			if len(urlIssues) != tt.wantLen {
				t.Errorf("got %d LBL002 issues, want %d", len(urlIssues), tt.wantLen)
			}

			if tt.wantLen > 0 && urlIssues[0].Code != tt.wantCode {
				t.Errorf("got code %q, want %q", urlIssues[0].Code, tt.wantCode)
			}
		})
	}
}

func TestBlogroll_PrimaryPersonValidation(t *testing.T) {
	tests := []struct {
		name     string
		feeds    []models.ExternalFeedConfig
		wantLen  int
		wantCode string
	}{
		{
			name: "valid primary_person reference",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://dave.com/feed.xml", Handle: "daverupert", Primary: boolPtr(true)},
				{URL: "https://dave.social/feed.xml", Handle: "davesocial", PrimaryPerson: "daverupert"},
			},
			wantLen: 0,
		},
		{
			name: "invalid primary_person reference",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://dave.social/feed.xml", Handle: "davesocial", PrimaryPerson: "nonexistent"},
			},
			wantLen:  1,
			wantCode: "LBL003",
		},
		{
			name: "empty primary_person is valid",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://example.com/feed.xml", Handle: "example", PrimaryPerson: ""},
			},
			wantLen: 0,
		},
		{
			name: "multiple invalid references",
			feeds: []models.ExternalFeedConfig{
				{URL: "https://a.com/feed.xml", Handle: "a", PrimaryPerson: "missing1"},
				{URL: "https://b.com/feed.xml", Handle: "b", PrimaryPerson: "missing2"},
			},
			wantLen:  2,
			wantCode: "LBL003",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &models.BlogrollConfig{Feeds: tt.feeds}
			result := Blogroll(config)

			var primaryIssues []BlogrollIssue
			for _, issue := range result.Issues {
				if issue.Code == "LBL003" {
					primaryIssues = append(primaryIssues, issue)
				}
			}

			if len(primaryIssues) != tt.wantLen {
				t.Errorf("got %d LBL003 issues, want %d", len(primaryIssues), tt.wantLen)
			}

			if tt.wantLen > 0 && primaryIssues[0].Code != tt.wantCode {
				t.Errorf("got code %q, want %q", primaryIssues[0].Code, tt.wantCode)
			}
		})
	}
}

func TestBlogroll_NilConfig(t *testing.T) {
	result := Blogroll(nil)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues for nil config, got %d", len(result.Issues))
	}
}

func TestBlogroll_EmptyFeeds(t *testing.T) {
	config := &models.BlogrollConfig{Feeds: []models.ExternalFeedConfig{}}
	result := Blogroll(config)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues for empty feeds, got %d", len(result.Issues))
	}
}

func TestBlogrollResult_HasErrors(t *testing.T) {
	tests := []struct {
		name   string
		issues []BlogrollIssue
		want   bool
	}{
		{
			name:   "no issues",
			issues: nil,
			want:   false,
		},
		{
			name:   "only warnings",
			issues: []BlogrollIssue{{Code: "LBL001", Severity: SeverityWarning}},
			want:   false,
		},
		{
			name:   "has error",
			issues: []BlogrollIssue{{Code: "LBL001", Severity: SeverityError}},
			want:   true,
		},
		{
			name: "mixed",
			issues: []BlogrollIssue{
				{Code: "LBL001", Severity: SeverityWarning},
				{Code: "LBL002", Severity: SeverityError},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BlogrollResult{Issues: tt.issues}
			if got := r.HasErrors(); got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlogrollResult_ErrorCount(t *testing.T) {
	result := &BlogrollResult{
		Issues: []BlogrollIssue{
			{Code: "LBL001", Severity: SeverityError},
			{Code: "LBL002", Severity: SeverityWarning},
			{Code: "LBL003", Severity: SeverityError},
		},
	}

	if got := result.ErrorCount(); got != 2 {
		t.Errorf("ErrorCount() = %d, want 2", got)
	}
}

func TestBlogrollResult_WarningCount(t *testing.T) {
	result := &BlogrollResult{
		Issues: []BlogrollIssue{
			{Code: "LBL001", Severity: SeverityError},
			{Code: "LBL002", Severity: SeverityWarning},
			{Code: "LBL003", Severity: SeverityWarning},
		},
	}

	if got := result.WarningCount(); got != 2 {
		t.Errorf("WarningCount() = %d, want 2", got)
	}
}

func TestExternalFeedConfig_IsPrimary(t *testing.T) {
	tests := []struct {
		name string
		feed models.ExternalFeedConfig
		want bool
	}{
		{
			name: "default is primary",
			feed: models.ExternalFeedConfig{URL: "https://example.com/feed.xml"},
			want: true,
		},
		{
			name: "explicit primary true",
			feed: models.ExternalFeedConfig{URL: "https://example.com/feed.xml", Primary: boolPtr(true)},
			want: true,
		},
		{
			name: "explicit primary false",
			feed: models.ExternalFeedConfig{URL: "https://example.com/feed.xml", Primary: boolPtr(false)},
			want: false,
		},
		{
			name: "has primary_person without explicit primary",
			feed: models.ExternalFeedConfig{URL: "https://example.com/feed.xml", PrimaryPerson: "someone"},
			want: false,
		},
		{
			name: "has primary_person with explicit primary true",
			feed: models.ExternalFeedConfig{URL: "https://example.com/feed.xml", PrimaryPerson: "someone", Primary: boolPtr(true)},
			want: true, // Explicit value wins
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.feed.IsPrimary()
			if got != tt.want {
				t.Errorf("IsPrimary() = %v, want %v", got, tt.want)
			}
		})
	}
}
