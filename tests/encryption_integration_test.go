// Package tests provides integration tests for markata-go.
package tests

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

// =============================================================================
// Encryption Integration Tests
// =============================================================================

// privateMarker is a unique string embedded in private post content.
// Tests verify this string NEVER appears in plaintext in the build output.
const privateMarker = "SUPER_SECRET_MARKER_TEXT_12345"

// buildWithEncryption creates a manager configured for encryption testing.
// It registers glob, load, render markdown, encryption, and publish plugins.
func buildWithEncryption(t *testing.T, site *testSite, modelsConfig *models.Config) (*lifecycle.Manager, error) {
	t.Helper()

	m := lifecycle.NewManager()

	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"

	if modelsConfig != nil {
		cfg.Extra["models_config"] = modelsConfig
	}

	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewDescriptionPlugin())
	m.RegisterPlugin(plugins.NewStructuredDataPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewEncryptionPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	err := m.Run()
	return m, err
}

// TestIntegration_Encryption_PrivateContentNeverPlaintext verifies that private
// content is encrypted and the marker text never appears in plaintext output.
func TestIntegration_Encryption_PrivateContentNeverPlaintext(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Set encryption key
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td") // pragma: allowlist secret

	site := newTestSite(t)

	// Private post with marker text
	site.addPost("secret-post.md", `---
title: Secret Post
slug: secret-post
published: true
private: true
---
This post contains `+privateMarker+` that should be encrypted.`)

	// Public post (no marker text)
	site.addPost("public-post.md", `---
title: Public Post
slug: public-post
published: true
---
This is a public post with nothing secret.`)

	// Another private post with explicit secret_key
	site.addPost("private-with-key.md", `---
title: Private With Key
slug: private-with-key
published: true
private: true
secret_key: default
---
Another private post with `+privateMarker+` inside it.`)

	modelsConfig := models.NewConfig()
	// Encryption is enabled by default with DefaultKey="default"

	m, err := buildWithEncryption(t, site, modelsConfig)
	if err != nil {
		t.Fatalf("build should succeed with encryption key set: %v", err)
	}

	// Verify private posts were encrypted
	for _, post := range m.Posts() {
		if !post.Private {
			continue
		}

		// The marker text must NOT appear in plaintext in ArticleHTML
		if strings.Contains(post.ArticleHTML, privateMarker) {
			t.Errorf("SECURITY: private post %q contains marker text in plaintext ArticleHTML", post.Path)
		}

		// Must have encryption wrapper
		if !strings.Contains(post.ArticleHTML, `class="encrypted-content"`) {
			t.Errorf("private post %q missing encrypted-content wrapper", post.Path)
		}
		if !strings.Contains(post.ArticleHTML, `data-encrypted="`) {
			t.Errorf("private post %q missing data-encrypted attribute", post.Path)
		}

		// Must be marked for script inclusion
		if hasEncrypted, ok := post.Extra["has_encrypted_content"].(bool); !ok || !hasEncrypted {
			t.Errorf("private post %q should have has_encrypted_content=true", post.Path)
		}
	}

	// Verify public post is NOT encrypted
	for _, post := range m.Posts() {
		if post.Private {
			continue
		}
		if strings.Contains(post.ArticleHTML, `class="encrypted-content"`) {
			t.Errorf("public post %q should NOT have encrypted-content wrapper", post.Path)
		}
		if strings.Contains(post.ArticleHTML, "data-encrypted=") {
			t.Errorf("public post %q should NOT have data-encrypted attribute", post.Path)
		}
	}

	// Verify the encrypted content can actually be decrypted
	for _, post := range m.Posts() {
		if !post.Private {
			continue
		}
		encryptedData := extractEncryptedData(t, post.ArticleHTML)
		if encryptedData == "" {
			t.Errorf("private post %q has no encrypted data", post.Path)
			continue
		}

		decrypted, decErr := encryption.Decrypt(encryptedData, "h7Qm!2Vx9#Lp4@Td") // pragma: allowlist secret
		if decErr != nil {
			t.Errorf("failed to decrypt private post %q: %v", post.Path, decErr)
			continue
		}

		// The decrypted content SHOULD contain the marker
		if !strings.Contains(string(decrypted), privateMarker) {
			t.Errorf("decrypted content of %q should contain marker text", post.Path)
		}
	}

	// Check output files if they were written
	if site.fileExists("secret-post/index.html") {
		content := site.readFile("secret-post/index.html")
		if strings.Contains(content, privateMarker) {
			t.Error("SECURITY: output file secret-post/index.html contains marker in plaintext")
		}
	}
	if site.fileExists("private-with-key/index.html") {
		content := site.readFile("private-with-key/index.html")
		if strings.Contains(content, privateMarker) {
			t.Error("SECURITY: output file private-with-key/index.html contains marker in plaintext")
		}
	}
}

// TestIntegration_Encryption_BuildFailsWithoutKey verifies the build fails with
// a CriticalError when private posts exist but no encryption key is available.
func TestIntegration_Encryption_BuildFailsWithoutKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Ensure NO encryption keys are set
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "") // pragma: allowlist secret

	site := newTestSite(t)

	site.addPost("secret.md", `---
title: Secret
slug: secret
published: true
private: true
---
This must not be published without encryption.`)

	modelsConfig := models.NewConfig()

	_, err := buildWithEncryption(t, site, modelsConfig)
	if err == nil {
		t.Fatal("build should FAIL when private posts exist but no encryption key is set")
	}

	// Verify it's a CriticalError via HookErrors.HasCritical()
	var hookErrors *lifecycle.HookErrors
	if !errors.As(err, &hookErrors) {
		t.Fatalf("error should be *lifecycle.HookErrors, got: %T", err)
	}
	if !hookErrors.HasCritical() {
		t.Error("HookErrors should contain at least one critical error")
	}

	// Error message should mention the missing key
	if !strings.Contains(err.Error(), "secret.md") {
		t.Errorf("error should reference the affected post, got: %v", err)
	}
}

// TestIntegration_Encryption_PrivateTags verifies that posts are automatically
// marked private based on configured private_tags.
func TestIntegration_Encryption_PrivateTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")   // pragma: allowlist secret
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_PERSONAL", "vN3$kP8!tR2@wZ6#q") // pragma: allowlist secret

	site := newTestSite(t)

	// Post tagged "diary" -- should be auto-private via private_tags config
	site.addPost("diary-entry.md", `---
title: Diary Entry
slug: diary-entry
published: true
tags:
  - diary
  - reflection
---
My diary has `+privateMarker+` in it.`)

	// Public post with no matching tags
	site.addPost("blog-post.md", `---
title: Blog Post
slug: blog-post
published: true
tags:
  - golang
  - tutorial
---
This is a public blog post.`)

	modelsConfig := models.NewConfig()
	modelsConfig.Encryption.PrivateTags = map[string]string{
		"diary": "personal",
	}

	m, err := buildWithEncryption(t, site, modelsConfig)
	if err != nil {
		t.Fatalf("build should succeed: %v", err)
	}

	// Find the diary post and verify it was encrypted
	var diaryPost *models.Post
	var blogPost *models.Post
	for _, post := range m.Posts() {
		if strings.Contains(post.Path, "diary-entry") {
			diaryPost = post
		}
		if strings.Contains(post.Path, "blog-post") {
			blogPost = post
		}
	}

	if diaryPost == nil {
		t.Fatal("diary post not found in build output")
	}
	if !diaryPost.Private {
		t.Error("diary post should be marked private via private_tags")
	}
	if strings.Contains(diaryPost.ArticleHTML, privateMarker) {
		t.Error("SECURITY: diary post contains marker text in plaintext")
	}
	if !strings.Contains(diaryPost.ArticleHTML, `class="encrypted-content"`) {
		t.Error("diary post should be encrypted")
	}

	if blogPost == nil {
		t.Fatal("blog post not found in build output")
	}
	if blogPost.Private {
		t.Error("blog post should NOT be marked private")
	}
}

// TestIntegration_Encryption_PrivateTagAutoFeedIncludesPrivate verifies that
// private-tag posts are encrypted on their own pages and included in the
// auto-generated feed outputs for that tag with encrypted content.
func TestIntegration_Encryption_PrivateTagAutoFeedIncludesPrivate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "Lighthouse!Harbor!Cedar!42") // pragma: allowlist secret

	site := newTestSite(t)

	const privateTitle = "Private Diary Entry"
	const secondPrivateTitle = "Second Private Diary Entry"

	site.addPost("private-diary.md", `---
title: `+privateTitle+`
slug: private-diary
published: true
tags:
  - diary
---
This private diary entry contains `+privateMarker+` in the body.`)

	site.addPost("second-private-diary.md", `---
title: `+secondPrivateTitle+`
slug: second-private-diary
published: true
tags:
  - diary
---
This second private diary entry also contains `+privateMarker+` in the body.`)

	site.addPost("public-post.md", `---
title: Public Post
slug: public-post
published: true
tags:
  - tutorial
---
This is a public post outside the private tag.`)

	modelsConfig := models.NewConfig()
	modelsConfig.Encryption.PrivateTags = map[string]string{
		"diary": "default",
	}

	m := lifecycle.NewManager()
	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"
	cfg.Extra["models_config"] = modelsConfig
	cfg.Extra["auto_feeds"] = plugins.AutoFeedsConfig{
		Tags: plugins.AutoFeedTypeConfig{
			Enabled:    true,
			SlugPrefix: "tags",
			Formats: models.FeedFormats{
				HTML:       true,
				SimpleHTML: true,
				RSS:        true,
				Atom:       true,
				JSON:       true,
				Markdown:   true,
				Text:       true,
			},
		},
	}
	cfg.Extra["feed_defaults"] = models.FeedDefaults{
		ItemsPerPage:    1,
		OrphanThreshold: 0,
	}
	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewDescriptionPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewEncryptionPlugin())
	m.RegisterPlugin(plugins.NewAutoFeedsPlugin())
	m.RegisterPlugin(plugins.NewPublishFeedsPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	if err := m.Run(); err != nil {
		t.Fatalf("build should succeed: %v", err)
	}

	feedConfigsValue, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}
	feedConfigs, ok := feedConfigsValue.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs has wrong type: %T", feedConfigsValue)
	}

	var diaryFeed *models.FeedConfig
	for i := range feedConfigs {
		if feedConfigs[i].Slug == "tags/diary" {
			diaryFeed = &feedConfigs[i]
			break
		}
	}
	if diaryFeed == nil {
		t.Fatal("tags/diary feed config not found")
	}
	if !diaryFeed.IncludePrivate {
		t.Fatal("tags/diary feed should opt into private posts (private tag auto-feed)")
	}
	if len(diaryFeed.Posts) != 2 {
		t.Fatalf("tags/diary feed should include private-tag posts, got %d posts", len(diaryFeed.Posts))
	}
	if len(diaryFeed.Pages) == 0 {
		t.Fatal("tags/diary feed should paginate with private posts included")
	}

	privatePage := site.readFile("private-diary/index.html")
	if !strings.Contains(privatePage, `class="encrypted-content"`) {
		t.Fatal("private diary page should still render encrypted content")
	}

	// The diary tag HTML feed pages should exist with private posts visible
	htmlFeedFiles := []string{
		"tags/diary/index.html",
		"tags/diary/simple/index.html",
	}
	for _, path := range htmlFeedFiles {
		if !site.fileExists(path) {
			continue
		}
		content := site.readFile(path)
		if !strings.Contains(content, privateTitle) {
			t.Fatalf("private diary title should be visible in %s", path)
		}
		if strings.Contains(content, privateMarker) {
			t.Fatalf("private content marker leaked into %s", path)
		}
	}

	machineFeedFiles := []string{
		"tags/diary/rss.xml",
		"tags/diary/atom.xml",
		"tags/diary/feed.json",
	}
	for _, path := range machineFeedFiles {
		if !site.fileExists(path) {
			continue
		}
		content := site.readFile(path)
		if !strings.Contains(content, privateTitle) {
			t.Fatalf("private diary title should be visible in %s", path)
		}
		if !strings.Contains(content, "encrypted-content") {
			t.Fatalf("encrypted content wrapper should be present in %s", path)
		}
		if strings.Contains(content, privateMarker) {
			t.Fatalf("private content marker leaked into %s", path)
		}
	}

	// Second page should exist since itemsPerPage=1 and we have 2 private posts
	if !site.fileExists("tags/diary/page/2/index.html") {
		t.Fatal("private tag auto feed should generate page 2 for second private post")
	}

	// Tutorial tag should still exclude private posts (not a private tag)
	var tutorialFeed *models.FeedConfig
	for i := range feedConfigs {
		if feedConfigs[i].Slug == "tags/tutorial" {
			tutorialFeed = &feedConfigs[i]
			break
		}
	}
	if tutorialFeed == nil {
		t.Fatal("tags/tutorial feed config not found")
	}
	if tutorialFeed.IncludePrivate {
		t.Fatal("tags/tutorial feed should not opt into private posts")
	}
	if len(tutorialFeed.Posts) != 1 {
		t.Fatalf("tags/tutorial feed should include 1 public post")
	}
}

// TestIntegration_Encryption_FrontmatterAliases verifies that private_key and
// encryption_key frontmatter fields work as aliases for secret_key.
func TestIntegration_Encryption_FrontmatterAliases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td") // pragma: allowlist secret

	site := newTestSite(t)

	// Post using secret_key
	site.addPost("with-secret-key.md", `---
title: With SecretKey
slug: with-secret-key
published: true
private: true
secret_key: default
---
Secret content A `+privateMarker)

	// Post using private_key alias
	site.addPost("with-private-key.md", `---
title: With PrivateKey
slug: with-private-key
published: true
private: true
private_key: default
---
Secret content B `+privateMarker)

	// Post using encryption_key alias
	site.addPost("with-encryption-key.md", `---
title: With EncryptionKey
slug: with-encryption-key
published: true
private: true
encryption_key: default
---
Secret content C `+privateMarker)

	modelsConfig := models.NewConfig()

	m, err := buildWithEncryption(t, site, modelsConfig)
	if err != nil {
		t.Fatalf("build should succeed: %v", err)
	}

	// All three posts should be encrypted
	for _, post := range m.Posts() {
		if !post.Private {
			continue
		}
		if strings.Contains(post.ArticleHTML, privateMarker) {
			t.Errorf("SECURITY: post %q contains marker in plaintext (alias may not work)", post.Path)
		}
		if !strings.Contains(post.ArticleHTML, `class="encrypted-content"`) {
			t.Errorf("post %q should be encrypted", post.Path)
		}
	}
}

// TestIntegration_Encryption_FrontmatterKeyOverridesTagKey verifies that
// a frontmatter secret_key takes precedence over a tag-level key.
func TestIntegration_Encryption_FrontmatterKeyOverridesTagKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")       // pragma: allowlist secret
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_PERSONAL", "vN3$kP8!tR2@wZ6#q")     // pragma: allowlist secret
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_CUSTOM", "Cobalt-Drift-47!Lantern") // pragma: allowlist secret

	site := newTestSite(t)

	// Post with diary tag AND explicit secret_key -- frontmatter should win
	site.addPost("override.md", `---
title: Override Test
slug: override
published: true
private: true
secret_key: custom
tags:
  - diary
---
Override content `+privateMarker)

	modelsConfig := models.NewConfig()
	modelsConfig.Encryption.PrivateTags = map[string]string{
		"diary": "personal",
	}

	m, err := buildWithEncryption(t, site, modelsConfig)
	if err != nil {
		t.Fatalf("build should succeed: %v", err)
	}

	var overridePost *models.Post
	for _, post := range m.Posts() {
		if strings.Contains(post.Path, "override") {
			overridePost = post
		}
	}

	if overridePost == nil {
		t.Fatal("override post not found")
	}

	// The key should be "custom" (frontmatter), not "personal" (tag)
	keyName, ok := overridePost.Extra["encryption_key_name"].(string)
	if !ok {
		t.Fatal("override post should have encryption_key_name in Extra")
	}
	if keyName != "custom" {
		t.Errorf("encryption_key_name = %q, want %q (frontmatter should override tag)", keyName, "custom")
	}

	// Verify encryption happened with the custom key
	encryptedData := extractEncryptedData(t, overridePost.ArticleHTML)
	if encryptedData == "" {
		t.Fatal("override post should have encrypted data")
	}

	// Should decrypt with custom-pass, not personal-pass
	decrypted, decErr := encryption.Decrypt(encryptedData, "Cobalt-Drift-47!Lantern") // pragma: allowlist secret
	if decErr != nil {
		t.Fatalf("should decrypt with custom key password: %v", decErr)
	}
	if !strings.Contains(string(decrypted), privateMarker) {
		t.Error("decrypted content should contain marker")
	}
}

// TestIntegration_Encryption_DisabledSkipsEncryption verifies that when
// encryption is disabled, private posts pass through unmodified.
func TestIntegration_Encryption_DisabledSkipsEncryption(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	site := newTestSite(t)

	site.addPost("secret.md", `---
title: Secret
slug: secret
published: true
private: true
---
Content here.`)

	modelsConfig := models.NewConfig()
	modelsConfig.Encryption.Enabled = false

	m, err := buildWithEncryption(t, site, modelsConfig)
	if err != nil {
		t.Fatalf("build should succeed with encryption disabled: %v", err)
	}

	for _, post := range m.Posts() {
		if strings.Contains(post.ArticleHTML, `class="encrypted-content"`) {
			t.Error("encryption should not run when disabled")
		}
	}
}

// TestIntegration_Encryption_DraftAndSkippedPostsIgnored verifies that
// draft and skipped posts are not subject to encryption requirements.
func TestIntegration_Encryption_DraftAndSkippedPostsIgnored(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// No encryption key -- normally this would fail if private posts exist
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "") // pragma: allowlist secret

	site := newTestSite(t)

	// Draft private post -- should be ignored
	site.addPost("draft.md", `---
title: Draft
slug: draft
published: true
private: true
draft: true
---
Draft content.`)

	// Non-private published post
	site.addPost("public.md", `---
title: Public
slug: public
published: true
---
Public content.`)

	modelsConfig := models.NewConfig()

	_, err := buildWithEncryption(t, site, modelsConfig)
	if err != nil {
		t.Fatalf("build should succeed when only draft/skipped posts are private: %v", err)
	}
}

// =============================================================================
// Tracer Scan Test — The Definitive Security Test
// =============================================================================

// TestIntegration_Encryption_TracerScan is the definitive security test.
// It plants a unique marker string in a private post, runs a comprehensive
// build with all output-producing plugins, then recursively scans every file
// in the output directory. The marker must NOT appear in ANY file in plaintext.
//
// This test catches leaks through ALL channels:
//   - HTML output (publish_html)
//   - Feed output (RSS, Atom, JSON, HTML feeds, Markdown, Text)
//   - 404 index (_404-index.json)
//   - Sitemaps (sitemap.xml)
//   - Search index (pagefind, if present)
//   - Any other file written to the output directory
func TestIntegration_Encryption_TracerScan(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Set encryption key
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "Maple!Orbit!Harbor!29") // pragma: allowlist secret

	site := newTestSite(t)

	// Use separate markers for content (private) vs metadata (public).
	// The CONTENT marker must never appear in any output file in plaintext.
	// The TITLE and DESCRIPTION markers are public metadata and SHOULD appear.
	const contentMarker = "TRACER_CONTENT_7f3a9b2e1d4c_PRIVATE_BODY_LEAK_DETECTOR"
	const titleText = "My Secret Encrypted Post Title"
	const descriptionText = "An explicit frontmatter description for the encrypted post"

	// Private post with content marker in body only, and safe title/description
	site.addPost("private-tracer.md", `---
title: `+titleText+`
slug: private-tracer
published: true
private: true
description: "`+descriptionText+`"
tags:
  - golang
  - secrets
---
# Private Heading

This private post body contains the tracer marker: `+contentMarker+`

It also has **bold** and [links](https://example.com) and other markdown.

Second paragraph also has `+contentMarker+` for good measure.
`)

	// Public post for reference (ensures build works and output is produced)
	site.addPost("public-post.md", `---
title: A Normal Public Post
slug: public-post
published: true
description: "This is a perfectly normal public post"
tags:
  - golang
  - tutorial
---
# Public Content

This is a public post. It should be output normally.
`)

	// Another public post to enable prev/next navigation
	site.addPost("another-public.md", `---
title: Another Public Post
slug: another-public
published: true
tags:
  - golang
---
Another public post for navigation testing.
`)

	modelsConfig := models.NewConfig()
	// Encryption is enabled by default

	// Build with a comprehensive set of plugins covering all output channels
	m := lifecycle.NewManager()

	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"
	cfg.Extra["models_config"] = modelsConfig

	// Configure a feed so feed output is generated
	cfg.Extra["feeds"] = []models.FeedConfig{
		{
			Slug:   "all",
			Title:  "All Posts",
			Filter: "published==true",
			Formats: models.FeedFormats{
				HTML: true,
				RSS:  true,
			},
		},
	}
	cfg.Extra["feed_defaults"] = models.FeedDefaults{
		ItemsPerPage:    10,
		OrphanThreshold: 3,
		Formats: models.FeedFormats{
			HTML: true,
			RSS:  true,
		},
	}

	m.SetConfig(cfg)

	// Register all output-producing plugins (excluding external deps like pagefind/blogroll)
	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewAutoTitlePlugin())
	m.RegisterPlugin(plugins.NewDescriptionPlugin())
	m.RegisterPlugin(plugins.NewStructuredDataPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewEncryptionPlugin())
	m.RegisterPlugin(plugins.NewTemplatesPlugin())
	m.RegisterPlugin(plugins.NewFeedsPlugin())
	m.RegisterPlugin(plugins.NewPrevNextPlugin())
	m.RegisterPlugin(plugins.NewPublishFeedsPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())
	m.RegisterPlugin(plugins.NewErrorPagesPlugin())
	m.RegisterPlugin(plugins.NewSitemapPlugin())

	err := m.Run()
	if err != nil {
		t.Fatalf("build should succeed: %v", err)
	}

	// Phase 1: Check in-memory post objects
	for _, post := range m.Posts() {
		if !post.Private {
			continue
		}

		// ArticleHTML must not contain the content marker in plaintext
		if strings.Contains(post.ArticleHTML, contentMarker) {
			t.Errorf("SECURITY LEAK: private post %q has content marker in ArticleHTML", post.Path)
		}

		// Content (raw markdown) must be scrubbed
		if strings.Contains(post.Content, contentMarker) {
			t.Errorf("SECURITY LEAK: private post %q has content marker in Content (raw markdown not scrubbed)", post.Path)
		}

		// Title should be PRESERVED (public metadata)
		if post.Title == nil || *post.Title != titleText {
			t.Errorf("private post %q: Title should be preserved as %q, got %v", post.Path, titleText, post.Title)
		}

		// Description should be PRESERVED (explicitly set in frontmatter)
		if post.Description == nil || *post.Description != descriptionText {
			t.Errorf("private post %q: explicit Description should be preserved as %q, got %v", post.Path, descriptionText, post.Description)
		}

		// Structured data should be PRESERVED (contains only title/description/dates)
		if _, ok := post.Extra["structured_data"]; !ok {
			t.Errorf("private post %q: structured_data should be preserved in Extra", post.Path)
		}

		// Must have data-pagefind-ignore to prevent search indexing
		if !strings.Contains(post.ArticleHTML, "data-pagefind-ignore") {
			t.Errorf("private post %q missing data-pagefind-ignore attribute", post.Path)
		}
	}

	// Phase 2: Recursively scan EVERY file in the output directory
	// The CONTENT marker must not appear in any file in plaintext.
	var filesScanned int
	var leakFiles []string

	scanErr := filepath.Walk(site.outputDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		filesScanned++

		// Read the file content
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Errorf("failed to read output file %s: %v", path, readErr)
			return nil
		}

		content := string(data)
		if strings.Contains(content, contentMarker) {
			relPath := path
			if rel, relErr := filepath.Rel(site.outputDir, path); relErr == nil {
				relPath = rel
			}
			leakFiles = append(leakFiles, relPath)
			t.Errorf("SECURITY LEAK: content marker found in output file: %s", relPath)
		}

		return nil
	})
	if scanErr != nil {
		t.Fatalf("failed to walk output directory: %v", scanErr)
	}

	if filesScanned == 0 {
		t.Fatal("no files found in output directory — build may not have produced output")
	}

	t.Logf("Tracer scan complete: scanned %d files, found %d leaks", filesScanned, len(leakFiles))

	if len(leakFiles) > 0 {
		t.Errorf("SECURITY FAILURE: content marker leaked to %d files: %s",
			len(leakFiles), strings.Join(leakFiles, ", "))
	}

	// Phase 3: Verify the encrypted content CAN be decrypted (proves encryption actually worked)
	for _, post := range m.Posts() {
		if !post.Private {
			continue
		}
		encData := extractEncryptedData(t, post.ArticleHTML)
		if encData == "" {
			t.Errorf("private post %q has no encrypted data attribute", post.Path)
			continue
		}
		decrypted, decErr := encryption.Decrypt(encData, "Maple!Orbit!Harbor!29") // pragma: allowlist secret
		if decErr != nil {
			t.Errorf("failed to decrypt private post %q: %v", post.Path, decErr)
			continue
		}
		if !strings.Contains(string(decrypted), contentMarker) {
			t.Errorf("decrypted content of %q should contain the content marker", post.Path)
		}
	}

	// Phase 4: Verify public posts are unaffected
	for _, post := range m.Posts() {
		if post.Private {
			continue
		}
		if strings.Contains(post.ArticleHTML, `class="encrypted-content"`) {
			t.Errorf("public post %q should NOT be encrypted", post.Path)
		}
		// Public posts should have their content intact
		if post.Content == "" && post.ArticleHTML == "" {
			t.Errorf("public post %q has empty content — scrubbing may have been too aggressive", post.Path)
		}
	}
}

// TestIntegration_Encryption_MetadataScrubbing verifies that private post metadata
// is handled correctly after encryption: title and explicit description are preserved,
// raw content is scrubbed, and structured data is kept.
func TestIntegration_Encryption_MetadataScrubbing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "7fR@2mL#9pQ!4xT$8k") // pragma: allowlist secret

	site := newTestSite(t)

	const contentMarker = "METADATA_SCRUB_CONTENT_abc123"
	const titleText = "Private Meta Test"
	const descriptionText = "An explicit description from frontmatter"

	site.addPost("private-meta.md", `---
title: `+titleText+`
slug: private-meta
published: true
private: true
description: "`+descriptionText+`"
tags:
  - golang
---
Body content with `+contentMarker+` inside.
`)

	// Private post WITHOUT explicit description
	site.addPost("private-no-desc.md", `---
title: Private No Description
slug: private-no-desc
published: true
private: true
tags:
  - golang
---
Body with `+contentMarker+` and no frontmatter description.
`)

	site.addPost("public-meta.md", `---
title: Public Meta Test
slug: public-meta
published: true
---
Normal public content.
`)

	modelsConfig := models.NewConfig()

	m, err := buildWithEncryption(t, site, modelsConfig)
	if err != nil {
		t.Fatalf("build should succeed: %v", err)
	}

	for _, post := range m.Posts() {
		if !post.Private {
			continue
		}

		// Content (raw markdown) must be empty
		if post.Content != "" {
			t.Errorf("private post %q: Content should be empty after scrubbing, got %q", post.Path, post.Content)
		}

		if strings.HasSuffix(post.Path, "private-meta.md") {
			// Title should be preserved
			if post.Title == nil || *post.Title != titleText {
				t.Errorf("private post %q: Title should be preserved as %q, got %v", post.Path, titleText, post.Title)
			}

			// Explicit frontmatter description should be preserved
			if post.Description == nil || *post.Description != descriptionText {
				t.Errorf("private post %q: explicit Description should be preserved as %q, got %v", post.Path, descriptionText, post.Description)
			}

			// Structured data should be preserved
			if _, ok := post.Extra["structured_data"]; !ok {
				t.Errorf("private post %q: structured_data should be preserved in Extra", post.Path)
			}
		}

		if strings.HasSuffix(post.Path, "private-no-desc.md") {
			// Title should be preserved
			if post.Title == nil || *post.Title != "Private No Description" {
				t.Errorf("private post %q: Title should be preserved, got %v", post.Path, post.Title)
			}

			// No explicit description was provided, so it should remain nil
			// (description plugin skips private posts, so no auto-generation)
			if post.Description != nil {
				t.Errorf("private post %q: Description should be nil (no explicit frontmatter description), got %q", post.Path, *post.Description)
			}
		}
	}

	// Verify public posts still have their metadata
	for _, post := range m.Posts() {
		if post.Private {
			continue
		}
		if post.Content == "" {
			t.Errorf("public post %q: Content should not be empty", post.Path)
		}
	}
}

// =============================================================================
// Helpers
// =============================================================================

// extractEncryptedData extracts the base64-encoded encrypted data from the
// encrypted-content div's data-encrypted attribute.
func extractEncryptedData(t *testing.T, html string) string {
	t.Helper()
	marker := `data-encrypted="`
	start := strings.Index(html, marker)
	if start == -1 {
		return ""
	}
	start += len(marker)
	end := strings.Index(html[start:], `"`)
	if end == -1 {
		return ""
	}
	return html[start : start+end]
}
