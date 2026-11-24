package filter

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func makePost(opts ...func(*models.Post)) *models.Post {
	p := models.NewPost("test.md")
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func withTitle(title string) func(*models.Post) {
	return func(p *models.Post) {
		p.Title = &title
	}
}

func withPublished(published bool) func(*models.Post) {
	return func(p *models.Post) {
		p.Published = published
	}
}

func withDraft(draft bool) func(*models.Post) {
	return func(p *models.Post) {
		p.Draft = draft
	}
}

func withSkip(skip bool) func(*models.Post) {
	return func(p *models.Post) {
		p.Skip = skip
	}
}

func withTags(tags ...string) func(*models.Post) {
	return func(p *models.Post) {
		p.Tags = tags
	}
}

func withAuthors(authors ...string) func(*models.Post) {
	return func(p *models.Post) {
		p.Authors = authors
	}
}

func withAuthor(author string) func(*models.Post) {
	return func(p *models.Post) {
		p.Author = &author
	}
}

func withDate(date time.Time) func(*models.Post) {
	return func(p *models.Post) {
		p.Date = &date
	}
}

func withTemplate(template string) func(*models.Post) {
	return func(p *models.Post) {
		p.Template = template
	}
}

func withExtra(key string, value interface{}) func(*models.Post) {
	return func(p *models.Post) {
		p.Set(key, value)
	}
}

func TestEvaluate_BooleanComparison(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "published == True (true)",
			expr:     "published == True",
			post:     makePost(withPublished(true)),
			expected: true,
		},
		{
			name:     "published == True (false)",
			expr:     "published == True",
			post:     makePost(withPublished(false)),
			expected: false,
		},
		{
			name:     "draft == False",
			expr:     "draft == False",
			post:     makePost(withDraft(false)),
			expected: true,
		},
		{
			name:     "published != False",
			expr:     "published != False",
			post:     makePost(withPublished(true)),
			expected: true,
		},
		{
			name:     "author == value (from authors array)",
			expr:     "author == 'waylon'",
			post:     makePost(withAuthors("waylon", "wyatt")),
			expected: true,
		},
		{
			name:     "author == value (legacy author field)",
			expr:     "author == 'waylon'",
			post:     makePost(withAuthor("waylon")),
			expected: true,
		},
		{
			name:     "templateKey aliases template",
			expr:     `templateKey == "photo"`,
			post:     makePost(withTemplate("shots")),
			expected: true,
		},
		{
			name:     "template equality respects aliases",
			expr:     `template == "shots"`,
			post:     makePost(withTemplate("photo")),
			expected: true,
		},
		{
			name:     "template inequality respects aliases",
			expr:     `template != "shots"`,
			post:     makePost(withTemplate("photo")),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_StringInList(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "'python' in tags (present)",
			expr:     "'python' in tags",
			post:     makePost(withTags("go", "python", "rust")),
			expected: true,
		},
		{
			name:     "'python' in tags (absent)",
			expr:     "'python' in tags",
			post:     makePost(withTags("go", "rust")),
			expected: false,
		},
		{
			name:     "'go' in tags (empty)",
			expr:     "'go' in tags",
			post:     makePost(),
			expected: false,
		},
		{
			name:     "double quotes",
			expr:     `"python" in tags`,
			post:     makePost(withTags("python")),
			expected: true,
		},
		{
			name:     "author in authors (present)",
			expr:     "'waylon' in authors",
			post:     makePost(withAuthors("waylon", "wyatt")),
			expected: true,
		},
		{
			name:     "author in authors (absent)",
			expr:     "'waylon' in authors",
			post:     makePost(withAuthors("wyatt")),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_DateComparison(t *testing.T) {
	// today is end-of-day (23:59:59.999999999) so that "date <= today"
	// includes posts from today regardless of their time component.
	ctx := &EvalContext{
		Today: time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
		Now:   time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "date <= today (past)",
			expr:     "date <= today",
			post:     makePost(withDate(time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))),
			expected: true,
		},
		{
			name:     "date <= today (today at midnight)",
			expr:     "date <= today",
			post:     makePost(withDate(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))),
			expected: true,
		},
		{
			name:     "date <= today (today with time component)",
			expr:     "date <= today",
			post:     makePost(withDate(time.Date(2024, 1, 15, 9, 24, 52, 0, time.UTC))),
			expected: true,
		},
		{
			name:     "date <= today (today at 23:59)",
			expr:     "date <= today",
			post:     makePost(withDate(time.Date(2024, 1, 15, 23, 59, 0, 0, time.UTC))),
			expected: true,
		},
		{
			name:     "date <= today (future)",
			expr:     "date <= today",
			post:     makePost(withDate(time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC))),
			expected: false,
		},
		{
			name:     "date > today",
			expr:     "date > today",
			post:     makePost(withDate(time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC))),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			f.SetContext(ctx)
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_CompoundExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "and (both true)",
			expr:     "published == True and draft == False",
			post:     makePost(withPublished(true), withDraft(false)),
			expected: true,
		},
		{
			name:     "and (first false)",
			expr:     "published == True and draft == False",
			post:     makePost(withPublished(false), withDraft(false)),
			expected: false,
		},
		{
			name:     "and (second false)",
			expr:     "published == True and draft == False",
			post:     makePost(withPublished(true), withDraft(true)),
			expected: false,
		},
		{
			name:     "or (first true)",
			expr:     "published == True or draft == True",
			post:     makePost(withPublished(true), withDraft(false)),
			expected: true,
		},
		{
			name:     "or (second true)",
			expr:     "published == True or draft == True",
			post:     makePost(withPublished(false), withDraft(true)),
			expected: true,
		},
		{
			name:     "or (both false)",
			expr:     "published == True or draft == True",
			post:     makePost(withPublished(false), withDraft(false)),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_Negation(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "not skip (false)",
			expr:     "not skip",
			post:     makePost(withSkip(false)),
			expected: true,
		},
		{
			name:     "not skip (true)",
			expr:     "not skip",
			post:     makePost(withSkip(true)),
			expected: false,
		},
		{
			name:     "not draft",
			expr:     "not draft",
			post:     makePost(withDraft(false)),
			expected: true,
		},
		{
			name:     "not not published",
			expr:     "not not published",
			post:     makePost(withPublished(true)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_StringMethods(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "startswith (match)",
			expr:     "title.startswith('How')",
			post:     makePost(withTitle("How to write tests")),
			expected: true,
		},
		{
			name:     "startswith (no match)",
			expr:     "title.startswith('How')",
			post:     makePost(withTitle("Writing tests")),
			expected: false,
		},
		{
			name:     "endswith (match)",
			expr:     "title.endswith('tests')",
			post:     makePost(withTitle("How to write tests")),
			expected: true,
		},
		{
			name:     "contains (match)",
			expr:     "title.contains('write')",
			post:     makePost(withTitle("How to write tests")),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_NumericComparison(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "word_count > 400 (true)",
			expr:     "word_count > 400",
			post:     makePost(withExtra("word_count", 500)),
			expected: true,
		},
		{
			name:     "word_count > 400 (false)",
			expr:     "word_count > 400",
			post:     makePost(withExtra("word_count", 300)),
			expected: false,
		},
		{
			name:     "word_count >= 400",
			expr:     "word_count >= 400",
			post:     makePost(withExtra("word_count", 400)),
			expected: true,
		},
		{
			name:     "word_count < 1000",
			expr:     "word_count < 1000",
			post:     makePost(withExtra("word_count", 500)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_SpecialLiterals(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "True always matches",
			expr:     "True",
			post:     makePost(),
			expected: true,
		},
		{
			name:     "False never matches",
			expr:     "False",
			post:     makePost(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvaluate_ExtraFields(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "status == 'draft'",
			expr:     "status == 'draft'",
			post:     makePost(withExtra("status", "draft")),
			expected: true,
		},
		{
			name:     "status == 'draft' or status == 'review' (first)",
			expr:     "status == 'draft' or status == 'review'",
			post:     makePost(withExtra("status", "draft")),
			expected: true,
		},
		{
			name:     "status == 'draft' or status == 'review' (second)",
			expr:     "status == 'draft' or status == 'review'",
			post:     makePost(withExtra("status", "review")),
			expected: true,
		},
		{
			name:     "status == 'draft' or status == 'review' (neither)",
			expr:     "status == 'draft' or status == 'review'",
			post:     makePost(withExtra("status", "published")),
			expected: false,
		},
		{
			name:     "featured == True",
			expr:     "featured == True",
			post:     makePost(withExtra("featured", true)),
			expected: true,
		},
		{
			name:     "'wip' in tags",
			expr:     "'wip' in tags",
			post:     makePost(withTags("wip", "draft")),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result, err := f.Match(tt.post)
			if err != nil {
				t.Fatalf("failed to evaluate: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFilter_MatchAll(t *testing.T) {
	posts := []*models.Post{
		makePost(withTitle("Post 1"), withPublished(true), withTags("go")),
		makePost(withTitle("Post 2"), withPublished(false), withTags("python")),
		makePost(withTitle("Post 3"), withPublished(true), withTags("go", "python")),
		makePost(withTitle("Post 4"), withPublished(true), withTags("rust")),
	}

	tests := []struct {
		name     string
		expr     string
		expected int
	}{
		{"all published", "published == True", 3},
		{"has python tag", "'python' in tags", 2},
		{"published with go", "published == True and 'go' in tags", 2},
		{"all", "True", 4},
		{"none", "False", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			result := f.MatchAll(posts)
			if len(result) != tt.expected {
				t.Errorf("expected %d posts, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestFilter_Combinators(t *testing.T) {
	post := makePost(withPublished(true), withDraft(false))

	t.Run("And", func(t *testing.T) {
		f1 := MustParse("published == True")
		f2 := MustParse("draft == False")
		combined := And(f1, f2)

		result, err := combined.Match(post)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result {
			t.Error("expected match")
		}
	})

	t.Run("Or", func(t *testing.T) {
		f1 := MustParse("published == False")
		f2 := MustParse("draft == False")
		combined := Or(f1, f2)

		result, err := combined.Match(post)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result {
			t.Error("expected match")
		}
	})

	t.Run("Not", func(t *testing.T) {
		f := MustParse("draft == True")
		negated := Not(f)

		result, err := negated.Match(post)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result {
			t.Error("expected match")
		}
	})
}

func TestFilter_ConvenienceFunctions(t *testing.T) {
	post := makePost(withPublished(true))

	t.Run("Always", func(t *testing.T) {
		f := Always()
		result, _ := f.Match(post) //nolint:errcheck // Always() never errors
		if !result {
			t.Error("Always() should match")
		}
	})

	t.Run("Never", func(t *testing.T) {
		f := Never()
		result, _ := f.Match(post) //nolint:errcheck // Never() never errors
		if result {
			t.Error("Never() should not match")
		}
	})

	t.Run("Posts", func(t *testing.T) {
		posts := []*models.Post{
			makePost(withPublished(true)),
			makePost(withPublished(false)),
		}
		result, err := Posts("published == True", posts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 post, got %d", len(result))
		}
	})

	t.Run("MatchPost", func(t *testing.T) {
		result, err := MatchPost("published == True", post)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result {
			t.Error("expected match")
		}
	})
}
