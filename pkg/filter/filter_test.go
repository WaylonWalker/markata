package filter

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

// ptr returns a pointer to the given string
func ptr(s string) *string {
	return &s
}

// timePtr returns a pointer to the given time
func timePtr(t time.Time) *time.Time {
	return &t
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s) //nolint:errcheck // test helper, format is controlled
	return t
}

// =============================================================================
// Filter Expression Tests based on tests.yaml
// =============================================================================

func TestFilter_BooleanTrue(t *testing.T) {
	// Test case: "filter by boolean true"
	// filter: "published == True"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter by boolean true",
			filter: "published == True",
			posts: []*models.Post{
				{Title: ptr("Post 1"), Published: true},
				{Title: ptr("Post 2"), Published: false},
				{Title: ptr("Post 3"), Published: true},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched
			for _, p := range result {
				if !p.Published {
					t.Errorf("matched post %v should be published", *p.Title)
				}
			}
		})
	}
}

func TestFilter_BooleanFalse(t *testing.T) {
	// Test case: "filter by boolean false"
	// filter: "draft == False"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter by boolean false",
			filter: "draft == False",
			posts: []*models.Post{
				{Title: ptr("Post 1"), Draft: true},
				{Title: ptr("Post 2"), Draft: false},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched
			for _, p := range result {
				if p.Draft {
					t.Errorf("matched post %v should not be draft", *p.Title)
				}
			}
		})
	}
}

func TestFilter_StringInList(t *testing.T) {
	// Test case: "filter by string in list"
	// filter: "'python' in tags"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter by string in list",
			filter: "'python' in tags",
			posts: []*models.Post{
				{Title: ptr("Python Tips"), Tags: []string{"python", "tips"}},
				{Title: ptr("JS Guide"), Tags: []string{"javascript"}},
				{Title: ptr("Python Advanced"), Tags: []string{"python", "advanced"}},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched
			for _, p := range result {
				found := false
				for _, tag := range p.Tags {
					if tag == "python" {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("matched post %v should have 'python' tag", *p.Title)
				}
			}
		})
	}
}

func TestFilter_DateComparison(t *testing.T) {
	// Test case: "filter by date comparison"
	// filter: "date <= today"
	today := parseDate("2024-07-01")

	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		today    time.Time
		expected int
	}{
		{
			name:   "filter by date comparison",
			filter: "date <= today",
			posts: []*models.Post{
				{Title: ptr("Old Post"), Date: timePtr(parseDate("2023-01-01"))},
				{Title: ptr("New Post"), Date: timePtr(parseDate("2024-06-15"))},
				{Title: ptr("Future Post"), Date: timePtr(parseDate("2025-01-01"))},
			},
			today:    today,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			// Set context with specific today value
			ctx := &EvalContext{
				Today: tt.today,
				Now:   tt.today,
			}
			f.SetContext(ctx)
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched (dates <= today)
			for _, p := range result {
				if p.Date != nil && p.Date.After(tt.today) {
					t.Errorf("matched post %v date %v should be <= %v", *p.Title, p.Date, tt.today)
				}
			}
		})
	}
}

func TestFilter_CompoundAnd(t *testing.T) {
	// Test case: "filter with compound and"
	// filter: "published == True and featured == True"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter with compound and",
			filter: "published == True and featured == True",
			posts: []*models.Post{
				{Title: ptr("Post 1"), Published: true, Extra: map[string]interface{}{"featured": true}},
				{Title: ptr("Post 2"), Published: true, Extra: map[string]interface{}{"featured": false}},
				{Title: ptr("Post 3"), Published: false, Extra: map[string]interface{}{"featured": true}},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched
			for _, p := range result {
				if !p.Published {
					t.Errorf("matched post %v should be published", *p.Title)
				}
				if featured, ok := p.Extra["featured"].(bool); !ok || !featured {
					t.Errorf("matched post %v should be featured", *p.Title)
				}
			}
		})
	}
}

func TestFilter_CompoundOr(t *testing.T) {
	// Test case: "filter with compound or"
	// filter: "status == 'draft' or status == 'review'"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter with compound or",
			filter: "status == 'draft' or status == 'review'",
			posts: []*models.Post{
				{Title: ptr("Draft"), Extra: map[string]interface{}{"status": "draft"}},
				{Title: ptr("Review"), Extra: map[string]interface{}{"status": "review"}},
				{Title: ptr("Published"), Extra: map[string]interface{}{"status": "published"}},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched
			for _, p := range result {
				statusVal, ok := p.Extra["status"]
				if !ok {
					t.Errorf("matched post %v should have status field", *p.Title)
					continue
				}
				status, ok := statusVal.(string)
				if !ok {
					t.Errorf("matched post %v status should be string", *p.Title)
					continue
				}
				if status != "draft" && status != "review" {
					t.Errorf("matched post %v should have status 'draft' or 'review', got %q", *p.Title, status)
				}
			}
		})
	}
}

func TestFilter_AlwaysTrue(t *testing.T) {
	// Test case: "filter all (True)"
	// filter: "True"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter all (True)",
			filter: "True",
			posts: []*models.Post{
				{Title: ptr("Post 1")},
				{Title: ptr("Post 2")},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFilter_AlwaysFalse(t *testing.T) {
	// Test case: "filter none (False)"
	// filter: "False"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter none (False)",
			filter: "False",
			posts: []*models.Post{
				{Title: ptr("Post 1")},
				{Title: ptr("Post 2")},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFilter_StringMethod(t *testing.T) {
	// Test case: "filter by string method"
	// filter: "title.startswith('How')"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter by string method startswith",
			filter: "title.startswith('How')",
			posts: []*models.Post{
				{Title: ptr("How to Cook")},
				{Title: ptr("Why Python")},
				{Title: ptr("How to Code")},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched
			for _, p := range result {
				if p.Title == nil || (*p.Title)[:3] != "How" {
					titleStr := ""
					if p.Title != nil {
						titleStr = *p.Title
					}
					t.Errorf("matched post title %q should start with 'How'", titleStr)
				}
			}
		})
	}
}

func TestFilter_NumericComparison(t *testing.T) {
	// Test case: "filter with numeric comparison"
	// filter: "word_count > 400"
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter with numeric comparison",
			filter: "word_count > 400",
			posts: []*models.Post{
				{Title: ptr("Short"), Extra: map[string]interface{}{"word_count": 100}},
				{Title: ptr("Medium"), Extra: map[string]interface{}{"word_count": 500}},
				{Title: ptr("Long"), Extra: map[string]interface{}{"word_count": 2000}},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
			// Verify correct posts were matched
			for _, p := range result {
				wordCount, ok := p.Extra["word_count"].(int)
				if !ok || wordCount <= 400 {
					t.Errorf("matched post %v word_count %d should be > 400", *p.Title, wordCount)
				}
			}
		})
	}
}

// =============================================================================
// Additional Filter Tests
// =============================================================================

func TestFilter_EmptyExpression(t *testing.T) {
	// Empty expression should match all
	f, err := Parse("")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	posts := []*models.Post{
		{Title: ptr("Post 1")},
		{Title: ptr("Post 2")},
	}

	result := f.MatchAll(posts)
	if len(result) != 2 {
		t.Errorf("empty expression should match all, got %d", len(result))
	}
}

func TestFilter_NotExpression(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "not draft",
			filter: "not draft",
			posts: []*models.Post{
				{Title: ptr("Post 1"), Draft: true},
				{Title: ptr("Post 2"), Draft: false},
				{Title: ptr("Post 3"), Draft: false},
			},
			expected: 2,
		},
		{
			name:   "not published",
			filter: "not published",
			posts: []*models.Post{
				{Title: ptr("Post 1"), Published: true},
				{Title: ptr("Post 2"), Published: false},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFilter_StringEndswith(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter by endswith",
			filter: "title.endswith('Tutorial')",
			posts: []*models.Post{
				{Title: ptr("Python Tutorial")},
				{Title: ptr("Go Guide")},
				{Title: ptr("JavaScript Tutorial")},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFilter_StringContains(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "filter by contains",
			filter: "title.contains('Python')",
			posts: []*models.Post{
				{Title: ptr("Learn Python Fast")},
				{Title: ptr("Go Guide")},
				{Title: ptr("Python Advanced")},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFilter_ComparisonOperators(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		posts    []*models.Post
		expected int
	}{
		{
			name:   "less than",
			filter: "count < 50",
			posts: []*models.Post{
				{Title: ptr("Low"), Extra: map[string]interface{}{"count": 10}},
				{Title: ptr("Mid"), Extra: map[string]interface{}{"count": 50}},
				{Title: ptr("High"), Extra: map[string]interface{}{"count": 100}},
			},
			expected: 1,
		},
		{
			name:   "less than or equal",
			filter: "count <= 50",
			posts: []*models.Post{
				{Title: ptr("Low"), Extra: map[string]interface{}{"count": 10}},
				{Title: ptr("Mid"), Extra: map[string]interface{}{"count": 50}},
				{Title: ptr("High"), Extra: map[string]interface{}{"count": 100}},
			},
			expected: 2,
		},
		{
			name:   "greater than or equal",
			filter: "count >= 50",
			posts: []*models.Post{
				{Title: ptr("Low"), Extra: map[string]interface{}{"count": 10}},
				{Title: ptr("Mid"), Extra: map[string]interface{}{"count": 50}},
				{Title: ptr("High"), Extra: map[string]interface{}{"count": 100}},
			},
			expected: 2,
		},
		{
			name:   "not equal",
			filter: "status != 'published'",
			posts: []*models.Post{
				{Title: ptr("Draft"), Extra: map[string]interface{}{"status": "draft"}},
				{Title: ptr("Published"), Extra: map[string]interface{}{"status": "published"}},
				{Title: ptr("Review"), Extra: map[string]interface{}{"status": "review"}},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.filter)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			result := f.MatchAll(tt.posts)
			if len(result) != tt.expected {
				t.Errorf("got %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFilter_InvalidExpression(t *testing.T) {
	// Test invalid filter expressions
	invalidFilters := []string{
		"published = True",   // Single equals
		"published ==",       // Missing right operand
		"== True",            // Missing left operand
		"(published == True", // Unclosed parenthesis
	}

	for _, filter := range invalidFilters {
		t.Run(filter, func(t *testing.T) {
			_, err := Parse(filter)
			if err == nil {
				t.Errorf("expected error for invalid filter %q", filter)
			}
		})
	}
}

func TestFilterSpec_ConvenienceFunctions(t *testing.T) {
	posts := []*models.Post{
		{Title: ptr("Post 1"), Published: true},
		{Title: ptr("Post 2"), Published: false},
	}

	// Test Posts
	result, err := Posts("published == True", posts)
	if err != nil {
		t.Fatalf("Posts error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Posts: got %d, want 1", len(result))
	}

	// Test MatchPost
	match, err := MatchPost("published == True", posts[0])
	if err != nil {
		t.Fatalf("MatchPost error: %v", err)
	}
	if !match {
		t.Error("MatchPost: expected true")
	}
}

func TestFilter_CombineFilters(t *testing.T) {
	posts := []*models.Post{
		{Title: ptr("Post 1"), Published: true, Draft: false},
		{Title: ptr("Post 2"), Published: true, Draft: true},
		{Title: ptr("Post 3"), Published: false, Draft: false},
	}

	// Test And combinator
	f1 := MustParse("published == True")
	f2 := MustParse("draft == False")
	combined := And(f1, f2)

	result := combined.MatchAll(posts)
	if len(result) != 1 {
		t.Errorf("And: got %d, want 1", len(result))
	}

	// Test Or combinator
	f3 := MustParse("published == True")
	f4 := MustParse("draft == True")
	orCombined := Or(f3, f4)

	result = orCombined.MatchAll(posts)
	if len(result) != 2 {
		t.Errorf("Or: got %d, want 2", len(result))
	}

	// Test Not combinator
	f5 := MustParse("published == True")
	notFilter := Not(f5)

	result = notFilter.MatchAll(posts)
	if len(result) != 1 {
		t.Errorf("Not: got %d, want 1", len(result))
	}
}

func TestFilter_Always_Never(t *testing.T) {
	posts := []*models.Post{
		{Title: ptr("Post 1")},
		{Title: ptr("Post 2")},
	}

	// Test Always
	always := Always()
	result := always.MatchAll(posts)
	if len(result) != 2 {
		t.Errorf("Always: got %d, want 2", len(result))
	}

	// Test Never
	never := Never()
	result = never.MatchAll(posts)
	if len(result) != 0 {
		t.Errorf("Never: got %d, want 0", len(result))
	}
}

func TestFilter_MatchAllWithErrors(t *testing.T) {
	posts := []*models.Post{
		{Title: ptr("Post 1"), Published: true},
		{Title: ptr("Post 2"), Published: false},
	}

	f := MustParse("published == True")
	result, errors := f.MatchAllWithErrors(posts)

	if len(errors) != 0 {
		t.Errorf("expected no errors, got %d", len(errors))
	}
	if len(result) != 1 {
		t.Errorf("got %d, want 1", len(result))
	}
}

func TestFilter_NilHandling(t *testing.T) {
	// Test filtering with nil values
	posts := []*models.Post{
		{Title: ptr("Post 1"), Date: timePtr(parseDate("2024-01-01"))},
		{Title: ptr("Post 2"), Date: nil},
		{Title: nil, Date: timePtr(parseDate("2024-02-01"))},
	}

	// Filter by date should handle nil dates
	f, err := Parse("date <= today")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	result := f.MatchAll(posts)
	// Posts with nil date should not match date comparisons
	if len(result) < 1 {
		t.Errorf("expected at least 1 match, got %d", len(result))
	}
}
