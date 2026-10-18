package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestReadingTimeOwnership_SharedCalculatorCharacterization(t *testing.T) {
	cases := []struct {
		name      string
		content   string
		wordCount int
		minutes   int
		text      string
	}{
		{name: "empty content", content: "", wordCount: 0, minutes: 0, text: "< 1 min read"},
		{name: "short prose", content: "One two three four five", wordCount: 5, minutes: 1, text: "1 min read"},
		{name: "multiple paragraphs", content: "First paragraph.\n\nSecond paragraph with more words.", wordCount: 7, minutes: 1, text: "1 min read"},
		{name: "markdown formatting", content: "**Bold** and _italic_ text.", wordCount: 4, minutes: 1, text: "1 min read"},
		{name: "links", content: "Read [the guide](https://example.com) and https://example.org now.", wordCount: 5, minutes: 1, text: "1 min read"},
		{name: "headings", content: "# Heading\n\n## Subheading\n\nBody text.", wordCount: 4, minutes: 1, text: "1 min read"},
		{name: "inline code", content: "Use the `fmt.Println` function to print.", wordCount: 5, minutes: 1, text: "1 min read"},
		{name: "fenced code", content: "# Intro\n\nText before.\n\n```go\nfunc main() {\n  fmt.Println(\"hello\")\n}\n```\n\nText after.", wordCount: 5, minutes: 1, text: "1 min read"},
		{name: "code only", content: "```go\nfmt.Println(\"hello\")\n```", wordCount: 0, minutes: 0, text: "< 1 min read"},
		{name: "html", content: "<p>Readable <strong>HTML</strong> text.</p>", wordCount: 3, minutes: 1, text: "1 min read"},
		{name: "whitespace", content: "  many\twords\nwith\r\nvaried   whitespace  ", wordCount: 5, minutes: 1, text: "1 min read"},
		{name: "unicode", content: "こんにちは、世界! Привет мир café naïve", wordCount: 6, minutes: 1, text: "1 min read"},
		{name: "unclosed fence", content: "Before\n\n```go\ncode words remain", wordCount: 5, minutes: 1, text: "1 min read"},
		{name: "199 words", content: strings.TrimSpace(strings.Repeat("word ", 199)), wordCount: 199, minutes: 1, text: "1 min read"},
		{name: "200 words", content: strings.TrimSpace(strings.Repeat("word ", 200)), wordCount: 200, minutes: 1, text: "1 min read"},
		{name: "201 words", content: strings.TrimSpace(strings.Repeat("word ", 201)), wordCount: 201, minutes: 2, text: "2 min read"},
	}

	readingTime := NewReadingTimePlugin()
	stats := NewStatsPlugin()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			metrics := calculateReadingTimeMetrics(tc.content, defaultWordsPerMinute, false)
			if metrics.WordCount != tc.wordCount || metrics.ReadingTime != tc.minutes || metrics.ReadingTimeText != tc.text {
				t.Errorf("shared calculator = %d, %d, %q; want %d, %d, %q", metrics.WordCount, metrics.ReadingTime, metrics.ReadingTimeText, tc.wordCount, tc.minutes, tc.text)
			}

			if got := readingTime.countWords(tc.content); got != tc.wordCount {
				t.Errorf("ReadingTimePlugin.countWords() = %d, want %d", got, tc.wordCount)
			}
			statsResult := stats.calculatePostStats(tc.content)
			if statsResult.WordCount != tc.wordCount || statsResult.ReadingTime != tc.minutes || statsResult.ReadingTimeText != tc.text {
				t.Errorf("StatsPlugin reading metrics = %d, %d, %q; want %d, %d, %q", statsResult.WordCount, statsResult.ReadingTime, statsResult.ReadingTimeText, tc.wordCount, tc.minutes, tc.text)
			}
			if statsResult.CharCount != metrics.CharCount {
				t.Errorf("StatsPlugin.CharCount = %d, want shared normalized count %d", statsResult.CharCount, metrics.CharCount)
			}
		})
	}
}

func TestReadingTimeOwnership_ReadingTimeFormattingAndWPM(t *testing.T) {
	plugin := NewReadingTimePlugin()
	plugin.SetWordsPerMinute(100)

	content := strings.TrimSpace(strings.Repeat("word ", 101))
	wordCount := plugin.countWords(content)
	readingTime := plugin.calculateReadingTime(wordCount)
	readingTimeText := plugin.formatReadingTime(readingTime)
	metrics := calculateReadingTimeMetrics(content, plugin.wordsPerMinute, false)
	if metrics.WordCount != 101 || metrics.ReadingTime != 2 || metrics.ReadingTimeText != "2 min read" {
		t.Fatalf("metrics = %+v; want 101 words, 2 minutes, %q", metrics, "2 min read")
	}
	if wordCount != metrics.WordCount || readingTime != metrics.ReadingTime || readingTimeText != metrics.ReadingTimeText {
		t.Fatalf("ReadingTimePlugin wrappers = %d, %d, %q; want shared metrics", wordCount, readingTime, readingTimeText)
	}

	for minutes, want := range map[int]string{
		0: "< 1 min read",
		1: "1 min read",
		2: "2 min read",
	} {
		if got := plugin.formatReadingTime(minutes); got != want {
			t.Errorf("formatReadingTime(%d) = %q, want %q", minutes, got, want)
		}
	}
}

func TestReadingTimeOwnership_StatsWPMPrecedence(t *testing.T) {
	cases := []struct {
		name  string
		extra map[string]interface{}
		want  int
	}{
		{
			name: "nested only",
			extra: map[string]interface{}{
				"stats": map[string]interface{}{"words_per_minute": 100},
			},
			want: 100,
		},
		{
			name: "top-level only",
			extra: map[string]interface{}{
				"words_per_minute": 150,
			},
			want: 150,
		},
		{
			name: "top-level overrides nested",
			extra: map[string]interface{}{
				"words_per_minute": 150,
				"stats":            map[string]interface{}{"words_per_minute": 100},
			},
			want: 150,
		},
		{
			name: "top-level overrides nested decoded types",
			extra: map[string]interface{}{
				"words_per_minute": int64(150),
				"stats":            map[string]interface{}{"words_per_minute": float64(100)},
			},
			want: 150,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			manager := lifecycle.NewManager()
			manager.Config().Extra = tc.extra
			plugin := NewStatsPlugin()
			if err := plugin.Configure(manager); err != nil {
				t.Fatal(err)
			}
			if plugin.wordsPerMinute != tc.want {
				t.Fatalf("wordsPerMinute = %d, want %d", plugin.wordsPerMinute, tc.want)
			}
		})
	}
}

func TestReadingTimeOwnership_WPMDefaultsForMissingAndNonPositiveValues(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]interface{}
	}{
		{name: "missing", extra: map[string]interface{}{}},
		{name: "top-level zero", extra: map[string]interface{}{"words_per_minute": int64(0)}},
		{name: "top-level negative", extra: map[string]interface{}{"words_per_minute": float64(-100)}},
		{name: "nested zero", extra: map[string]interface{}{"stats": map[string]interface{}{"words_per_minute": int64(0)}}},
		{name: "nested negative", extra: map[string]interface{}{"stats": map[string]interface{}{"words_per_minute": float64(-100)}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manager := lifecycle.NewManager()
			manager.Config().Extra = tc.extra

			stats := NewStatsPlugin()
			if err := stats.Configure(manager); err != nil {
				t.Fatal(err)
			}
			if stats.wordsPerMinute != defaultWordsPerMinute {
				t.Errorf("Stats wordsPerMinute = %d, want default %d", stats.wordsPerMinute, defaultWordsPerMinute)
			}

			readingTime := NewReadingTimePlugin()
			if err := readingTime.Configure(manager); err != nil {
				t.Fatal(err)
			}
			if readingTime.wordsPerMinute != defaultWordsPerMinute {
				t.Errorf("ReadingTime wordsPerMinute = %d, want default %d", readingTime.wordsPerMinute, defaultWordsPerMinute)
			}
		})
	}
}

func TestReadingTimePlugin_LoadedConfigFormatsWordsPerMinute(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		format config.Format
	}{
		{
			name:   "TOML",
			data:   "[markata-go]\nwords_per_minute = 100\n",
			format: config.FormatTOML,
		},
		{
			name:   "YAML",
			data:   "markata-go:\n  words_per_minute: 100\n",
			format: config.FormatYAML,
		},
		{
			name:   "JSON",
			data:   `{"markata-go":{"words_per_minute":100}}`,
			format: config.FormatJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loaded, err := config.LoadFromString(tt.data, tt.format)
			if err != nil {
				t.Fatalf("load config: %v", err)
			}

			manager := lifecycle.NewManager()
			manager.SetConcurrency(1)
			manager.Config().Extra = loaded.Extra
			post := models.NewPost("reading-time-" + tt.name + ".md")
			post.Content = strings.TrimSpace(strings.Repeat("word ", 401))
			manager.SetPosts([]*models.Post{post})

			plugin := NewReadingTimePlugin()
			if err := plugin.Configure(manager); err != nil {
				t.Fatalf("configure reading time: %v", err)
			}
			if err := plugin.Transform(manager); err != nil {
				t.Fatalf("transform reading time: %v", err)
			}

			if plugin.wordsPerMinute != 100 {
				t.Errorf("wordsPerMinute = %d, want 100", plugin.wordsPerMinute)
			}
			assertReadingTimeFields(t, post, 401, 5, "5 min read")
		})
	}
}

func TestReadingTimeOwnership_StandalonePluginContracts(t *testing.T) {
	t.Run("reading_time uses top-level WPM and excludes code", func(t *testing.T) {
		manager := lifecycle.NewManager()
		manager.Config().Extra["words_per_minute"] = 1
		post := models.NewPost("reading-time.md")
		post.Content = "one two\n\n```go\nthree four\n```"
		manager.SetPosts([]*models.Post{post})

		plugin := NewReadingTimePlugin()
		if err := plugin.Configure(manager); err != nil {
			t.Fatal(err)
		}
		if err := plugin.Transform(manager); err != nil {
			t.Fatal(err)
		}
		assertReadingTimeFields(t, post, 2, 2, "2 min read")
	})

	t.Run("stats uses nested WPM and includes code", func(t *testing.T) {
		manager := lifecycle.NewManager()
		manager.Config().Extra["stats"] = map[string]interface{}{
			"words_per_minute":      1,
			"include_code_in_count": true,
		}
		post := models.NewPost("stats.md")
		post.Content = "one two\n\n```go\nthree four\n```"
		manager.SetPosts([]*models.Post{post})

		plugin := NewStatsPlugin()
		if err := plugin.Configure(manager); err != nil {
			t.Fatal(err)
		}
		if err := plugin.Transform(manager); err != nil {
			t.Fatal(err)
		}
		assertReadingTimeFields(t, post, 4, 4, "4 min read")
		statsResult, ok := post.Get(statsPluginName).(*PostStats)
		if !ok {
			t.Fatalf("stats field = %T, want *PostStats", post.Get(statsPluginName))
		}
		if statsResult.WordCount != 4 || statsResult.ReadingTime != 4 {
			t.Fatalf("PostStats reading metrics = %d, %d; want 4, 4", statsResult.WordCount, statsResult.ReadingTime)
		}
	})
}

func TestReadingTimeOwnership_MalformedMultilineInlineCodeUsesCanonicalNormalization(t *testing.T) {
	content := "before `foo\nbar` after"

	// origin/main's Stats scanner stops malformed inline code at a newline and
	// reports 4 words and 17 letters/digits for this input. The shared
	// calculator intentionally follows ReadingTime's canonical normalization:
	// a backtick-delimited span may cross a newline, so the result is 2 words
	// and 11 letters/digits.
	stats := NewStatsPlugin().calculatePostStats(content)
	if stats.WordCount != 2 || stats.CharCount != 11 || stats.ReadingTime != 1 || stats.ReadingTimeText != "1 min read" {
		t.Fatalf("Stats metrics = %+v; want 2 words, 11 chars, 1 minute, %q", stats, "1 min read")
	}

	metrics := calculateReadingTimeMetrics(content, defaultWordsPerMinute, false)
	if stats.WordCount != metrics.WordCount || stats.CharCount != metrics.CharCount || stats.ReadingTime != metrics.ReadingTime || stats.ReadingTimeText != metrics.ReadingTimeText {
		t.Fatalf("Stats metrics = %+v, shared metrics = %+v; want canonical convergence", stats, metrics)
	}
}

func TestReadingTimeOwnership_DivergentDefaultContracts(t *testing.T) {
	content := strings.TrimSpace(strings.Repeat("prose ", 101)) +
		"\n\n```go\n" + strings.TrimSpace(strings.Repeat("code ", 100)) + "\n```"

	manager := lifecycle.NewManager()
	manager.SetConcurrency(1)
	manager.Config().Extra = map[string]interface{}{
		// Leave the top-level value absent so ReadingTime uses its default
		// contract while Stats uses its nested WPM configuration. Top-level
		// precedence is covered separately by TestReadingTimeOwnership_StatsWPMPrecedence.
		"stats": map[string]interface{}{
			"words_per_minute":      100,
			"include_code_in_count": true,
		},
	}
	post := models.NewPost("divergent.md")
	post.Content = content
	manager.SetPosts([]*models.Post{post})
	manager.SetFeeds([]*lifecycle.Feed{{Name: "posts", Posts: []*models.Post{post}}})
	manager.RegisterPlugins(NewReadingTimePlugin(), NewStatsPlugin())

	if err := manager.RunTo(lifecycle.StageCollect); err != nil {
		t.Fatalf("RunTo(collect) error = %v", err)
	}

	// ReadingTimePlugin runs later in the normal transform lifecycle and owns
	// these public fields, using top-level words_per_minute and excluding code.
	assertReadingTimeFields(t, post, 101, 1, "1 min read")

	// StatsPlugin calculated its own result early with Stats configuration.
	statsResult, ok := post.Get(statsPluginName).(*PostStats)
	if !ok {
		t.Fatalf("stats field = %T, want *PostStats", post.Get(statsPluginName))
	}
	if statsResult.WordCount != 201 || statsResult.ReadingTime != 3 || statsResult.ReadingTimeText != "3 min read" {
		t.Fatalf("PostStats = %+v; want 201 words, 3 minutes, %q", statsResult, "3 min read")
	}
	if statsResult.CodeBlocks != 1 {
		t.Fatalf("PostStats.CodeBlocks = %d, want 1", statsResult.CodeBlocks)
	}

	// Feed and site aggregates must use PostStats, not the later public fields.
	feedStats := GetFeedStats(manager, "posts")
	if feedStats == nil {
		t.Fatal("feed stats were not stored")
	}
	if feedStats.TotalWords != 201 || feedStats.TotalReadingTime != 3 || feedStats.AverageWords != 201 || feedStats.AverageReadingTime != 3 {
		t.Fatalf("feed stats = %+v; want Stats-derived 201 words and 3 minutes", feedStats)
	}

	siteStats := GetSiteStats(manager)
	if siteStats == nil {
		t.Fatal("site stats were not stored")
	}
	if siteStats.TotalWords != 201 || siteStats.TotalReadingTime != 3 || siteStats.AverageWords != 201 || siteStats.AverageReadingTime != 3 {
		t.Fatalf("site stats = %+v; want Stats-derived 201 words and 3 minutes", siteStats)
	}
	if helper := GetStatsHelper(manager); helper == nil || helper.TotalWords() != 201 || helper.TotalReadingTime() != 3 {
		t.Fatalf("StatsHelper = %#v; want Stats-derived totals", helper)
	}
}

func TestReadingTimeOwnership_JinjaExpandedContentUsesReadingTimeOrdering(t *testing.T) {
	manager := lifecycle.NewManager()
	manager.SetConcurrency(1)
	post := models.NewPost("jinja.md")
	post.Tags = []string{"alpha", "beta"}
	post.Content = "{% for tag in post.tags %}{{ tag }} {% endfor %}"
	post.Set("jinja", true)
	manager.SetPosts([]*models.Post{post})
	manager.RegisterPlugins(NewStatsPlugin(), NewReadingTimePlugin(), NewJinjaMdPlugin())

	for _, plugin := range manager.Plugins() {
		configured, ok := plugin.(lifecycle.ConfigurePlugin)
		if !ok {
			continue
		}
		if err := configured.Configure(manager); err != nil {
			t.Fatalf("configure %s: %v", plugin.Name(), err)
		}
	}
	if err := lifecycle.RunTransformHooksSubset(manager, manager.Posts()); err != nil {
		t.Fatalf("run transform hooks: %v", err)
	}

	if post.Content != "alpha beta " {
		t.Fatalf("Jinja content = %q, want %q", post.Content, "alpha beta ")
	}
	assertReadingTimeFields(t, post, 2, 1, "1 min read")
	statsResult, ok := post.Get(statsPluginName).(*PostStats)
	if !ok {
		t.Fatalf("stats field = %T, want *PostStats", post.Get(statsPluginName))
	}
	if statsResult.WordCount != 7 {
		t.Fatalf("Stats word count = %d, want 7 from the pre-Jinja content", statsResult.WordCount)
	}
}

func TestReadingTimeOwnership_FrontmatterIsExcludedByLoadLifecycle(t *testing.T) {
	dir := t.TempDir()
	content := "---\ntitle: Frontmatter has many words\ndescription: More metadata words\n---\n\nBody only."
	path := filepath.Join(dir, "post.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	manager := lifecycle.NewManager()
	manager.SetConcurrency(1)
	manager.Config().ContentDir = dir
	manager.SetFiles([]string{"post.md"})
	load := NewLoadPlugin()
	if err := load.Configure(manager); err != nil {
		t.Fatal(err)
	}
	if err := load.Load(manager); err != nil {
		t.Fatal(err)
	}
	posts := manager.Posts()
	if len(posts) != 1 {
		t.Fatalf("loaded posts = %d, want 1", len(posts))
	}
	if posts[0].Content != "\nBody only." {
		t.Fatalf("loaded content = %q, want body without frontmatter", posts[0].Content)
	}

	manager.RegisterPlugins(NewReadingTimePlugin(), NewStatsPlugin())
	for _, plugin := range manager.Plugins() {
		configured, ok := plugin.(lifecycle.ConfigurePlugin)
		if !ok {
			continue
		}
		if err := configured.Configure(manager); err != nil {
			t.Fatalf("configure %s: %v", plugin.Name(), err)
		}
	}
	if err := lifecycle.RunTransformHooksSubset(manager, manager.Posts()); err != nil {
		t.Fatalf("run transform hooks: %v", err)
	}
	assertReadingTimeFields(t, posts[0], 2, 1, "1 min read")
}

func assertReadingTimeFields(t *testing.T, post *models.Post, wantWordCount, wantReadingTime int, wantText string) {
	t.Helper()
	if post == nil || post.Extra == nil {
		t.Fatal("reading-time fields are missing")
	}
	wordCount, wordCountOK := post.Extra["word_count"].(int)
	readingTime, readingTimeOK := post.Extra["reading_time"].(int)
	readingTimeText, readingTimeTextOK := post.Extra["reading_time_text"].(string)
	if !wordCountOK || !readingTimeOK || !readingTimeTextOK {
		t.Fatalf("reading-time fields are incomplete: %#v", post.Extra)
	}
	if wordCount != wantWordCount || readingTime != wantReadingTime || readingTimeText != wantText {
		t.Fatalf("reading-time fields = %d, %d, %q; want %d, %d, %q", wordCount, readingTime, readingTimeText, wantWordCount, wantReadingTime, wantText)
	}
}
