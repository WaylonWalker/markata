package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// =============================================================================
// WikilinksPlugin Tests based on tests.yaml
// =============================================================================

func TestWikilinksPlugin_Name(t *testing.T) {
	p := NewWikilinksPlugin()
	if p.Name() != "wikilinks" {
		t.Errorf("expected name 'wikilinks', got %q", p.Name())
	}
}

func TestWikilinksPlugin_Configure(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()
	if err := p.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}
}

func TestWikilinksPlugin_ConfigureWarnBroken(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.Extra = map[string]interface{}{
		"wikilinks_warn_broken": false,
	}

	if err := p.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}

	// Verify setting was applied (indirectly through behavior)
	if p.warnOnBroken != false {
		t.Error("expected warnOnBroken to be false after configuration")
	}
}

func TestWikilinksPlugin_BasicWikilink(t *testing.T) {
	// Test case from tests.yaml: "basic wikilink"
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Other Post"
	targetPost := &models.Post{
		Slug:  "other-post",
		Title: &targetTitle,
		Href:  "/other-post/",
	}
	sourcePost := &models.Post{
		Content: "Check out [[other-post]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Check that link is created with correct href, class, and data attributes
	if !strings.Contains(source.Content, `<a href="/other-post/" class="wikilink"`) {
		t.Errorf("expected wikilink anchor tag in content, got %q", source.Content)
	}
	if !strings.Contains(source.Content, `data-title="Other Post"`) {
		t.Errorf("expected data-title attribute in content, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">Other Post</a>") {
		t.Errorf("expected 'Other Post' as link text in content, got %q", source.Content)
	}
}

func TestWikilinksPlugin_WikilinkWithCustomText(t *testing.T) {
	// Test case from tests.yaml: "wikilink with custom text"
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Other Post"
	targetPost := &models.Post{
		Slug:  "other-post",
		Title: &targetTitle,
		Href:  "/other-post/",
	}
	sourcePost := &models.Post{
		Content: "See [[other-post|this article]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Check that custom display text is used
	if !strings.Contains(source.Content, `<a href="/other-post/" class="wikilink"`) {
		t.Errorf("expected wikilink anchor tag in content, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">this article</a>") {
		t.Errorf("expected 'this article' as link text in content, got %q", source.Content)
	}
}

func TestWikilinksPlugin_WikilinkToMissingPost(t *testing.T) {
	// Test case from tests.yaml: "wikilink to missing post"
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	sourcePost := &models.Post{
		Content: "Link to [[nonexistent]]",
		Slug:    "source-post",
		Extra:   make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Content should be unchanged - wikilink left as is
	if !strings.Contains(source.Content, "[[nonexistent]]") {
		t.Errorf("expected [[nonexistent]] to remain in content, got %q", source.Content)
	}

	// Check for warning
	warnings, ok := source.Extra["wikilink_warnings"].([]string)
	if !ok || len(warnings) == 0 {
		t.Error("expected wikilink warning for missing post")
	} else if !strings.Contains(warnings[0], "broken wikilink") {
		t.Errorf("expected broken wikilink warning, got %q", warnings[0])
	}
}

func TestWikilinksPlugin_CaseInsensitiveLookup(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "My Post"
	targetPost := &models.Post{
		Slug:  "my-post",
		Title: &targetTitle,
		Href:  "/my-post/",
	}
	sourcePost := &models.Post{
		Content: "Check out [[MY-POST]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Check that case-insensitive match works
	if !strings.Contains(source.Content, `<a href="/my-post/" class="wikilink"`) {
		t.Errorf("expected case-insensitive match in content, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">My Post</a>") {
		t.Errorf("expected 'My Post' as link text in content, got %q", source.Content)
	}
}

func TestWikilinksPlugin_SlugWithSpaces(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Hello World"
	targetPost := &models.Post{
		Slug:  "hello-world",
		Title: &targetTitle,
		Href:  "/hello-world/",
	}
	sourcePost := &models.Post{
		Content: "See [[hello world]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Check that space to hyphen conversion works
	if !strings.Contains(source.Content, `<a href="/hello-world/" class="wikilink"`) {
		t.Errorf("expected space to hyphen conversion in content, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">Hello World</a>") {
		t.Errorf("expected 'Hello World' as link text in content, got %q", source.Content)
	}
}

func TestWikilinksPlugin_MultipleWikilinks(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	post1Title := "Post One"
	post2Title := "Post Two"
	post1 := &models.Post{Slug: "post-one", Title: &post1Title, Href: "/post-one/"}
	post2 := &models.Post{Slug: "post-two", Title: &post2Title, Href: "/post-two/"}
	sourcePost := &models.Post{
		Content: "See [[post-one]] and [[post-two]] for more",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{post1, post2, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	if !strings.Contains(source.Content, `<a href="/post-one/" class="wikilink"`) {
		t.Errorf("expected first wikilink to be converted, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">Post One</a>") {
		t.Errorf("expected first wikilink text, got %q", source.Content)
	}
	if !strings.Contains(source.Content, `<a href="/post-two/" class="wikilink"`) {
		t.Errorf("expected second wikilink to be converted, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">Post Two</a>") {
		t.Errorf("expected second wikilink text, got %q", source.Content)
	}
}

func TestWikilinksPlugin_SkippedPost(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Target"
	targetPost := &models.Post{Slug: "target", Title: &targetTitle, Href: "/target/"}
	sourcePost := &models.Post{
		Content: "See [[target]]",
		Slug:    "source-post",
		Skip:    true,
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Content should be unchanged for skipped posts
	if !strings.Contains(source.Content, "[[target]]") {
		t.Errorf("expected skipped post content to be unchanged, got %q", source.Content)
	}
}

func TestWikilinksPlugin_EmptyContent(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	sourcePost := &models.Post{
		Content: "",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	if posts[0].Content != "" {
		t.Errorf("expected empty content to remain empty, got %q", posts[0].Content)
	}
}

func TestWikilinksPlugin_PostWithoutTitle(t *testing.T) {
	// When target post has no title, should use slug
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetPost := &models.Post{
		Slug: "my-post",
		Href: "/my-post/",
		// No Title set
	}
	sourcePost := &models.Post{
		Content: "See [[my-post]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should use slug as display text
	expected := `<a href="/my-post/" class="wikilink">my-post</a>`
	if !strings.Contains(source.Content, expected) {
		t.Errorf("expected slug as display text %q in content, got %q", expected, source.Content)
	}
}

func TestWikilinksPlugin_HTMLEscaping(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Test & Demo"
	targetPost := &models.Post{
		Slug:  "test-demo",
		Title: &targetTitle,
		Href:  "/test-demo/",
	}
	sourcePost := &models.Post{
		Content: "See [[test-demo]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Ampersand should be escaped
	if !strings.Contains(source.Content, "&amp;") {
		t.Errorf("expected HTML-escaped content, got %q", source.Content)
	}
}

func TestWikilinksPlugin_NoHref(t *testing.T) {
	// When target post has no href, should generate from slug
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Target"
	targetPost := &models.Post{
		Slug:  "target",
		Title: &targetTitle,
		// No Href set
	}
	sourcePost := &models.Post{
		Content: "See [[target]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should generate href from slug
	if !strings.Contains(source.Content, `<a href="/target/" class="wikilink"`) {
		t.Errorf("expected generated href in content, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">Target</a>") {
		t.Errorf("expected 'Target' as link text in content, got %q", source.Content)
	}
}

func TestWikilinksPlugin_SetWarnOnBroken(t *testing.T) {
	p := NewWikilinksPlugin()

	// Default should be true
	if p.warnOnBroken != true {
		t.Error("expected default warnOnBroken to be true")
	}

	p.SetWarnOnBroken(false)
	if p.warnOnBroken != false {
		t.Error("expected warnOnBroken to be false after SetWarnOnBroken(false)")
	}

	p.SetWarnOnBroken(true)
	if p.warnOnBroken != true {
		t.Error("expected warnOnBroken to be true after SetWarnOnBroken(true)")
	}
}

func TestWikilinksPlugin_MissingPostNoWarningWhenDisabled(t *testing.T) {
	p := NewWikilinksPlugin()
	p.SetWarnOnBroken(false)

	m := lifecycle.NewManager()
	sourcePost := &models.Post{
		Content: "Link to [[nonexistent]]",
		Slug:    "source-post",
		Extra:   make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should not have warning when disabled
	warnings, ok := source.Extra["wikilink_warnings"].([]string)
	if ok && len(warnings) > 0 {
		t.Error("expected no warnings when warnOnBroken is disabled")
	}
}

// Interface compliance tests
func TestWikilinksPlugin_Interfaces(_ *testing.T) {
	p := NewWikilinksPlugin()

	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.TransformPlugin = p
}

// =============================================================================
// Alias Resolution Tests (Issue #415)
// =============================================================================

func TestWikilinksPlugin_AliasResolution(t *testing.T) {
	// Test that wikilinks can resolve via aliases defined in frontmatter
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "ECMAScript Language Specification"
	targetPost := &models.Post{
		Slug:  "ecmascript",
		Title: &targetTitle,
		Href:  "/ecmascript/",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"js", "javascript", "JavaScript"},
		},
	}
	sourcePost := &models.Post{
		Content: "Check out [[js]] for more details.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Check that alias resolved to the target post
	if !strings.Contains(source.Content, `<a href="/ecmascript/" class="wikilink"`) {
		t.Errorf("expected wikilink via alias to resolve to /ecmascript/, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">ECMAScript Language Specification</a>") {
		t.Errorf("expected post title as link text, got %q", source.Content)
	}
}

func TestWikilinksPlugin_SlugTakesPrecedenceOverAlias(t *testing.T) {
	// Test that slug takes precedence over alias when both match
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	// Post 1 has slug "javascript"
	post1Title := "JavaScript Guide"
	post1 := &models.Post{
		Slug:  "javascript",
		Title: &post1Title,
		Href:  "/javascript/",
	}

	// Post 2 has "javascript" as an alias
	post2Title := "ECMAScript Spec"
	post2 := &models.Post{
		Slug:  "ecmascript",
		Title: &post2Title,
		Href:  "/ecmascript/",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"javascript", "js"},
		},
	}

	sourcePost := &models.Post{
		Content: "See [[javascript]] for details.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{post1, post2, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should link to post1 (slug) not post2 (alias)
	if !strings.Contains(source.Content, `<a href="/javascript/" class="wikilink"`) {
		t.Errorf("expected slug to take precedence over alias, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">JavaScript Guide</a>") {
		t.Errorf("expected 'JavaScript Guide' as link text (from slug match), got %q", source.Content)
	}
}

func TestWikilinksPlugin_AliasCaseInsensitive(t *testing.T) {
	// Test that alias matching is case-insensitive
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "ECMAScript Specification"
	targetPost := &models.Post{
		Slug:  "ecmascript",
		Title: &targetTitle,
		Href:  "/ecmascript/",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"JavaScript", "JS"},
		},
	}

	sourcePost := &models.Post{
		Content: "See [[JAVASCRIPT]] and [[js]] for info.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Both wikilinks should resolve (case-insensitive)
	linkCount := strings.Count(source.Content, `<a href="/ecmascript/" class="wikilink"`)
	if linkCount != 2 {
		t.Errorf("expected 2 resolved wikilinks, got %d in %q", linkCount, source.Content)
	}
}

func TestWikilinksPlugin_MultipleAliasesSamePost(t *testing.T) {
	// Test that multiple aliases can point to the same post
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "TypeScript Handbook"
	targetPost := &models.Post{
		Slug:  "typescript",
		Title: &targetTitle,
		Href:  "/typescript/",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"ts", "TS", "TypeScript"},
		},
	}

	sourcePost := &models.Post{
		Content: "Learn [[ts]], [[typescript]], and [[TypeScript]].",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// All three wikilinks should resolve to the same post
	linkCount := strings.Count(source.Content, `<a href="/typescript/" class="wikilink"`)
	if linkCount != 3 {
		t.Errorf("expected 3 resolved wikilinks, got %d in %q", linkCount, source.Content)
	}
}

func TestWikilinksPlugin_AliasWithCustomText(t *testing.T) {
	// Test that alias resolution works with custom display text
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "ECMAScript Specification"
	targetPost := &models.Post{
		Slug:  "ecmascript",
		Title: &targetTitle,
		Href:  "/ecmascript/",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"js", "javascript"},
		},
	}

	sourcePost := &models.Post{
		Content: "Learn about [[js|the JavaScript language]] here.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should resolve via alias and use custom display text
	if !strings.Contains(source.Content, `<a href="/ecmascript/" class="wikilink"`) {
		t.Errorf("expected alias to resolve to /ecmascript/, got %q", source.Content)
	}
	if !strings.Contains(source.Content, ">the JavaScript language</a>") {
		t.Errorf("expected custom display text, got %q", source.Content)
	}
}

func TestWikilinksPlugin_AliasNoMatchStillWarns(t *testing.T) {
	// Test that non-matching aliases still produce warnings for broken links
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetPost := &models.Post{
		Slug: "ecmascript",
		Href: "/ecmascript/",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"js", "javascript"},
		},
	}

	sourcePost := &models.Post{
		Content: "See [[python]] for details.",
		Slug:    "source-post",
		Extra:   make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Content should be unchanged - wikilink left as is
	if !strings.Contains(source.Content, "[[python]]") {
		t.Errorf("expected [[python]] to remain in content, got %q", source.Content)
	}

	// Check for warning
	warnings, ok := source.Extra["wikilink_warnings"].([]string)
	if !ok || len(warnings) == 0 {
		t.Error("expected wikilink warning for missing post")
	}
}

// TestWikilinksPlugin_PreservesCodeBlocks tests that wikilinks inside code blocks are not transformed.
func TestWikilinksPlugin_PreservesCodeBlocks(t *testing.T) {
	p := NewWikilinksPlugin()

	m := lifecycle.NewManager()
	config := m.Config()
	config.Extra = map[string]interface{}{
		"wikilinks_warn_broken": false,
	}

	// Create target posts
	title := "Getting Started"
	target := &models.Post{
		Slug:  "getting-started",
		Title: &title,
		Href:  "/getting-started/",
	}

	// Create source post with wikilinks inside and outside code blocks
	source := &models.Post{
		Slug: "source-post",
		Content: `Here is a link: [[getting-started]]

` + "```markdown" + `
This [[getting-started]] should NOT be converted.
` + "```" + `

And another link: [[getting-started|custom text]]

` + "~~~" + `
This [[getting-started]] in tilde fence should also be preserved.
` + "~~~" + `

Final link: [[getting-started]]`,
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{source, target})

	// Run the plugin
	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	posts := m.Posts()
	var result *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			result = post
			break
		}
	}

	// Check that wikilinks outside code blocks are converted
	if !strings.Contains(result.Content, `<a href="/getting-started/" class="wikilink"`) {
		t.Errorf("expected wikilinks outside code blocks to be converted, got %q", result.Content)
	}
	if !strings.Contains(result.Content, ">Getting Started</a>") {
		t.Errorf("expected 'Getting Started' as link text in converted wikilinks, got %q", result.Content)
	}

	// Check that wikilinks inside backtick code blocks are preserved
	if !strings.Contains(result.Content, "This [[getting-started]] should NOT be converted.") {
		t.Errorf("expected wikilinks inside backtick code blocks to be preserved, got %q", result.Content)
	}

	// Check that wikilinks inside tilde code blocks are preserved
	if !strings.Contains(result.Content, "This [[getting-started]] in tilde fence should also be preserved.") {
		t.Errorf("expected wikilinks inside tilde code blocks to be preserved, got %q", result.Content)
	}

	// Check custom text link is converted
	if !strings.Contains(result.Content, ">custom text</a>") {
		t.Errorf("expected custom text wikilink to be converted, got %q", result.Content)
	}
}

// TestWikilinksPlugin_DataAttributes tests that data attributes are added for tooltips.
func TestWikilinksPlugin_DataAttributes(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "My Article"
	targetDesc := "A great article about things"
	targetDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	targetPost := &models.Post{
		Slug:        "my-article",
		Title:       &targetTitle,
		Description: &targetDesc,
		Date:        &targetDate,
		Href:        "/my-article/",
	}
	sourcePost := &models.Post{
		Content: "Read [[my-article]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Check data-title attribute
	if !strings.Contains(source.Content, `data-title="My Article"`) {
		t.Errorf("expected data-title attribute, got %q", source.Content)
	}

	// Check data-description attribute
	if !strings.Contains(source.Content, `data-description="A great article about things"`) {
		t.Errorf("expected data-description attribute, got %q", source.Content)
	}

	// Check data-date attribute
	if !strings.Contains(source.Content, `data-date="2024-01-15"`) {
		t.Errorf("expected data-date attribute, got %q", source.Content)
	}
}

// TestWikilinksPlugin_DataAttributesPartial tests that only available attributes are added.
func TestWikilinksPlugin_DataAttributesPartial(t *testing.T) {
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	// Post with only title, no description or date
	targetTitle := "Simple Post"
	targetPost := &models.Post{
		Slug:  "simple-post",
		Title: &targetTitle,
		Href:  "/simple-post/",
	}
	sourcePost := &models.Post{
		Content: "Read [[simple-post]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Check data-title is present
	if !strings.Contains(source.Content, `data-title="Simple Post"`) {
		t.Errorf("expected data-title attribute, got %q", source.Content)
	}

	// Check data-description is NOT present (no description)
	if strings.Contains(source.Content, `data-description`) {
		t.Errorf("expected no data-description attribute when description is nil, got %q", source.Content)
	}

	// Check data-date is NOT present (no date)
	if strings.Contains(source.Content, `data-date`) {
		t.Errorf("expected no data-date attribute when date is nil, got %q", source.Content)
	}
}

// =============================================================================
// Fragment Identifier Tests (Issue #990)
// =============================================================================

func TestWikilinksPlugin_FragmentIdentifier(t *testing.T) {
	// Test that [[slug#section]] links to the correct post with fragment
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Getting Started"
	targetPost := &models.Post{
		Slug:  "getting-started",
		Title: &targetTitle,
		Href:  "/getting-started/",
	}
	sourcePost := &models.Post{
		Content: "See [[getting-started#installation]] for setup.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should link to the post with fragment appended after trailing slash
	if !strings.Contains(source.Content, `href="/getting-started/#installation"`) {
		t.Errorf("expected href with fragment, got %q", source.Content)
	}
	// Should use post title as display text
	if !strings.Contains(source.Content, ">Getting Started</a>") {
		t.Errorf("expected post title as display text, got %q", source.Content)
	}
}

func TestWikilinksPlugin_FragmentWithCustomText(t *testing.T) {
	// Test that [[slug#section|Custom Text]] works
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Configuration Guide"
	targetPost := &models.Post{
		Slug:  "config-guide",
		Title: &targetTitle,
		Href:  "/config-guide/",
	}
	sourcePost := &models.Post{
		Content: "See [[config-guide#advanced|Advanced Config]] for details.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should have fragment in href
	if !strings.Contains(source.Content, `href="/config-guide/#advanced"`) {
		t.Errorf("expected href with fragment, got %q", source.Content)
	}
	// Should use custom display text
	if !strings.Contains(source.Content, ">Advanced Config</a>") {
		t.Errorf("expected custom display text, got %q", source.Content)
	}
}

func TestWikilinksPlugin_FragmentWithNoHref(t *testing.T) {
	// Test fragment when target post has no Href (generated from slug)
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "My Post"
	targetPost := &models.Post{
		Slug:  "my-post",
		Title: &targetTitle,
		// No Href set - should be generated
	}
	sourcePost := &models.Post{
		Content: "See [[my-post#details]] for more.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should generate href from slug with fragment
	if !strings.Contains(source.Content, `href="/my-post/#details"`) {
		t.Errorf("expected generated href with fragment, got %q", source.Content)
	}
}

func TestWikilinksPlugin_FragmentOnlyDoesNotBreakLookup(t *testing.T) {
	// Test that fragment doesn't pollute slug lookup
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "Target Post"
	targetPost := &models.Post{
		Slug:  "target",
		Title: &targetTitle,
		Href:  "/target/",
	}
	sourcePost := &models.Post{
		Content: "See [[target#section-one]] and [[target#section-two]]",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Both should resolve to the same post with different fragments
	if !strings.Contains(source.Content, `href="/target/#section-one"`) {
		t.Errorf("expected first fragment link, got %q", source.Content)
	}
	if !strings.Contains(source.Content, `href="/target/#section-two"`) {
		t.Errorf("expected second fragment link, got %q", source.Content)
	}
	// Should not contain broken wikilink warnings
	if strings.Contains(source.Content, "[[target#") {
		t.Errorf("expected wikilinks to be resolved, got %q", source.Content)
	}
}

func TestWikilinksPlugin_FragmentWithAlias(t *testing.T) {
	// Test that fragment works with alias resolution
	p := NewWikilinksPlugin()
	m := lifecycle.NewManager()

	targetTitle := "ECMAScript Specification"
	targetPost := &models.Post{
		Slug:  "ecmascript",
		Title: &targetTitle,
		Href:  "/ecmascript/",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"js", "javascript"},
		},
	}
	sourcePost := &models.Post{
		Content: "See [[js#modules]] for module info.",
		Slug:    "source-post",
	}

	m.SetPosts([]*models.Post{targetPost, sourcePost})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	posts := m.Posts()
	var source *models.Post
	for _, post := range posts {
		if post.Slug == "source-post" {
			source = post
			break
		}
	}

	// Should resolve via alias and include fragment
	if !strings.Contains(source.Content, `href="/ecmascript/#modules"`) {
		t.Errorf("expected alias resolution with fragment, got %q", source.Content)
	}
}
