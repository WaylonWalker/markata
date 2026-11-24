// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// EncryptionEnvPrefix is the prefix for encryption key environment variables.
const EncryptionEnvPrefix = "MARKATA_GO_ENCRYPTION_KEY_"

// EncryptionBuildError is returned when a private post cannot be encrypted.
// It implements lifecycle.CriticalError to halt the build and prevent
// exposing unencrypted private content.
type EncryptionBuildError struct {
	Posts []string
	Msg   string
}

func (e *EncryptionBuildError) Error() string {
	if len(e.Posts) == 0 {
		return fmt.Sprintf("encryption error: %s", e.Msg)
	}
	return fmt.Sprintf("encryption error: %s\n%s", e.Msg, strings.Join(e.Posts, "\n"))
}

// IsCritical marks this error as a build-halting error.
// Private posts must never be published without encryption.
func (e *EncryptionBuildError) IsCritical() bool {
	return true
}

// EncryptionPlugin encrypts content for private posts.
// It runs during the Render stage (after markdown is converted to HTML) to encrypt
// the ArticleHTML content.
//
// # Encryption is enabled by default
//
// The plugin is enabled by default with default_key="default". Users only need to
// set MARKATA_GO_ENCRYPTION_KEY_DEFAULT in their environment or .env file.
//
// # How It Works
//
//  1. Posts with private: true are automatically encrypted
//  2. Posts with tags matching [encryption.private_tags] are treated as private
//  3. Frontmatter secret_key (or aliases: private_key, encryption_key) specifies which key to use
//  4. If no key is specified, the default_key is used
//  5. The build FAILS if a private post has no available encryption key
//
// # Client-Side Decryption
//
// The encrypted content is wrapped in a div with data attributes.
// The client-side JavaScript uses Web Crypto API with matching PBKDF2 parameters.
type EncryptionPlugin struct {
	enabled        bool
	defaultKey     string
	decryptionHint string
	privateTags    map[string]string // tag -> key name
	// keys maps key names to passwords (loaded from env vars)
	keys                      map[string]string
	enforceStrength           bool
	minPasswordLength         int
	minEstimatedCrackDuration time.Duration
}

// NewEncryptionPlugin creates a new EncryptionPlugin.
func NewEncryptionPlugin() *EncryptionPlugin {
	return &EncryptionPlugin{
		keys:                      make(map[string]string),
		privateTags:               make(map[string]string),
		enforceStrength:           true,
		minPasswordLength:         encryption.DefaultMinPasswordLength,
		minEstimatedCrackDuration: encryption.DefaultMinEstimatedCrackDuration,
	}
}

// Name returns the unique name of the plugin.
func (p *EncryptionPlugin) Name() string {
	return "encryption"
}

// Configure loads encryption configuration and encryption keys from environment.
func (p *EncryptionPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Check for models.Config via Extra or direct config access
	if modelsConfig, ok := getModelsConfig(config); ok {
		p.enabled = modelsConfig.Encryption.Enabled
		p.defaultKey = modelsConfig.Encryption.DefaultKey
		p.decryptionHint = modelsConfig.Encryption.DecryptionHint
		p.enforceStrength = modelsConfig.Encryption.EnforceStrength
		p.minPasswordLength = modelsConfig.Encryption.MinPasswordLength // pragma: allowlist secret
		if p.minPasswordLength == 0 {
			p.minPasswordLength = encryption.DefaultMinPasswordLength // pragma: allowlist secret
		}
		durationStr := modelsConfig.Encryption.MinEstimatedCrackTime
		if durationStr == "" {
			durationStr = encryption.DefaultMinEstimatedCrackTime
		}
		duration, err := encryption.ParseEstimatedCrackDuration(durationStr)
		if err != nil {
			return fmt.Errorf("invalid encryption.min_estimated_crack_time: %w", err)
		}
		p.minEstimatedCrackDuration = duration
		if modelsConfig.Encryption.PrivateTags != nil {
			for tag, key := range modelsConfig.Encryption.PrivateTags {
				p.privateTags[strings.ToLower(tag)] = key
			}
		}
	}

	// Also check Extra for backward compatibility
	if config.Extra != nil {
		if enabled, ok := config.Extra["encryption_enabled"].(bool); ok {
			p.enabled = enabled
		}
		if defaultKey, ok := config.Extra["encryption_default_key"].(string); ok {
			p.defaultKey = defaultKey
		}
		if hint, ok := config.Extra["encryption_decryption_hint"].(string); ok {
			p.decryptionHint = hint
		}
	}

	if !p.enabled {
		return nil
	}

	// Load encryption keys from environment variables
	p.loadKeysFromEnvironment()

	return nil
}

// loadKeysFromEnvironment scans environment for MARKATA_GO_ENCRYPTION_KEY_* variables.
func (p *EncryptionPlugin) loadKeysFromEnvironment() {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, EncryptionEnvPrefix) {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				// Extract key name (lowercase for case-insensitive lookup)
				keyName := strings.ToLower(strings.TrimPrefix(parts[0], EncryptionEnvPrefix))
				password := parts[1]
				if password != "" {
					p.keys[keyName] = password
				}
			}
		}
	}
}

// getKeyPassword returns the password for a given key name.
// Falls back to default key if the specific key is not found.
func (p *EncryptionPlugin) getKeyPassword(keyName string) (string, error) {
	// Try exact key first (case-insensitive)
	if password, ok := p.keys[strings.ToLower(keyName)]; ok {
		return password, nil
	}

	// Try default key
	if p.defaultKey != "" {
		if password, ok := p.keys[strings.ToLower(p.defaultKey)]; ok {
			return password, nil
		}
	}

	return "", fmt.Errorf("encryption key %q not found in environment (expected %s%s)",
		keyName, EncryptionEnvPrefix, strings.ToUpper(keyName))
}

// Priority returns the plugin priority for the given stage.
//
// Transform stage: PriorityFirst (-1000)
//
//	Privacy marking (applyPrivateTags) must run before ALL other transform
//	plugins so that downstream plugins like Description (PriorityEarly = -100)
//	already see post.Private == true and skip private posts.
//
// Render stage: priority 50
//
//	Encryption of ArticleHTML runs after markdown rendering (PriorityDefault)
//	but before templates wrap the HTML (PriorityLate = 100).
func (p *EncryptionPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageTransform:
		// Must run before Description plugin (PriorityEarly = -100)
		return lifecycle.PriorityFirst
	case lifecycle.StageRender:
		// Run after markdown rendering but before templates
		return 50
	default:
		return lifecycle.PriorityDefault
	}
}

// applyPrivateTags marks posts as private based on configured private_tags.
// If a post has a tag or templateKey matching a private_tags entry, it is marked
// Private=true and assigned the tag's key (unless the post already has a
// frontmatter secret_key). This also checks the post's Template field (set from
// the templateKey frontmatter) to handle posts that use templateKey as their
// primary categorization without explicit tags.
func (p *EncryptionPlugin) applyPrivateTags(posts []*models.Post) {
	if len(p.privateTags) == 0 {
		return
	}

	for _, post := range posts {
		if post.Skip || post.Draft {
			continue
		}
		// An explicit private: false is an opt-out from tag/template defaults.
		if post.IsExplicitlyPublic() {
			continue
		}

		// Check tags
		matched := false
		for _, tag := range post.Tags {
			tagLower := strings.ToLower(tag)
			if keyName, ok := p.privateTags[tagLower]; ok {
				post.Private = true
				// Only set key from tag if frontmatter didn't specify one
				if post.SecretKey == "" {
					post.SecretKey = keyName // pragma: allowlist secret
				}
				matched = true
				break // one matching tag is enough
			}
		}

		// Also check templateKey (stored as post.Template) if tags didn't match.
		// Some posts use templateKey as their primary categorization without
		// including it in their tags list.
		if !matched && post.Template != "" {
			templateKeyLower := strings.ToLower(post.Template)
			if keyName, ok := p.privateTags[templateKeyLower]; ok {
				post.Private = true
				if post.SecretKey == "" {
					post.SecretKey = keyName // pragma: allowlist secret
				}
			}
		}
	}
}

// Transform marks posts as private based on configured private_tags.
// This runs at PriorityFirst (-1000) so that all subsequent transform plugins
// (e.g., Description at PriorityEarly) already see post.Private == true and
// can skip private posts, preventing content-derived metadata leaks.
func (p *EncryptionPlugin) Transform(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}

	p.applyPrivateTags(m.Posts())
	return p.validatePrivatePosts(m.Posts())
}

// Render encrypts content for private posts with encryption keys.
// Returns a CriticalError if any private post cannot be encrypted,
// preventing unencrypted private content from being published.
func (p *EncryptionPlugin) Render(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}

	cache := GetBuildCache(m)

	// Find all private posts (whether they need encryption or not)
	privatePosts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && !post.Draft && post.Private
	})

	if len(privatePosts) == 0 {
		return nil
	}

	// Check that every private post can be encrypted.
	// This is a safety check: we must NEVER expose private content unencrypted.
	if err := p.validatePrivatePosts(privatePosts); err != nil {
		return err
	}

	// Filter to posts that actually have content to encrypt
	postsToEncrypt := m.FilterPosts(func(post *models.Post) bool {
		return p.shouldEncrypt(post)
	})

	postsToEncrypt = filterEncryptedPostsForServe(m, postsToEncrypt)

	if len(postsToEncrypt) == 0 {
		return nil
	}

	regenerated := make(chan *models.Post, len(postsToEncrypt))
	err := m.ProcessPostsSliceConcurrently(postsToEncrypt, func(post *models.Post) error {
		changed, err := p.encryptPostWithCacheResult(post, cache)
		if changed {
			regenerated <- post
		}
		return err
	})
	close(regenerated)

	// Merge only after all workers finish: serve affected paths is an ordinary
	// map. Keep successful siblings invalidated even if another worker failed.
	affected := lifecycle.GetServeAffectedPaths(m)
	wrapperChanged := false
	for post := range regenerated {
		wrapperChanged = true
		if affected == nil {
			affected = make(map[string]bool)
		}
		affected[post.Path] = true
		if cache != nil {
			for _, identity := range postDependencyIdentities(post) {
				cache.MarkSlugChanged(identity)
			}
		}
	}
	if cache != nil && wrapperChanged {
		changed := cache.GetChangedSlugs()
		for _, path := range cache.GetAffectedPosts(changed) {
			affected[path] = true
		}
		cache.MarkAffectedDependents(changed)
	}
	if wrapperChanged {
		lifecycle.SetServeAffectedPaths(m, affected)
	}
	return err
}

func (p *EncryptionPlugin) validatePrivatePosts(posts []*models.Post) error {
	if len(posts) == 0 {
		return nil
	}

	var failures []encryptionFailure
	passwordCache := make(map[string]string)
	policyCache := make(map[string]error)
	for _, post := range posts {
		if post.Skip || post.Draft || !post.Private {
			continue
		}
		keyName := post.SecretKey
		if keyName == "" {
			keyName = p.defaultKey
		}
		if keyName == "" {
			failures = append(failures, encryptionFailure{
				postPath: post.Path,
				keyName:  "(none)",
				reason:   "no encryption key specified and no default key configured",
			})
			continue
		}
		normalized := strings.ToLower(keyName)
		if _, cached := passwordCache[normalized]; !cached {
			var password string
			var err error
			password, err = p.getKeyPassword(keyName)
			if err != nil {
				failures = append(failures, encryptionFailure{
					postPath: post.Path,
					keyName:  keyName,
					reason:   fmt.Sprintf("set %s%s in environment or .env", EncryptionEnvPrefix, strings.ToUpper(keyName)),
				})
				continue
			}
			passwordCache[normalized] = password
			if p.enforceStrength {
				policyCache[normalized] = p.validatePasswordPolicy(password)
			} else {
				policyCache[normalized] = nil
			}
		}
		if err := policyCache[normalized]; err != nil {
			failures = append(failures, encryptionFailure{
				postPath: post.Path,
				keyName:  keyName,
				reason:   err.Error(),
			})
		}
	}

	if len(failures) == 0 {
		return nil
	}

	return &EncryptionBuildError{
		Posts: summarizeEncryptionFailures(failures),
		Msg:   "private posts found without available encryption keys. Build halted to prevent exposing private content",
	}
}

const maxEncryptionFailureSamples = 3

type encryptionFailure struct {
	postPath string
	keyName  string
	reason   string
}

type encryptionFailureGroup struct {
	keyName string
	reason  string
	posts   []string
}

func summarizeEncryptionFailures(failures []encryptionFailure) []string {
	if len(failures) == 0 {
		return nil
	}

	groups := make(map[string]*encryptionFailureGroup)
	keys := make([]string, 0, len(failures))
	for _, failure := range failures {
		groupKey := failure.keyName + "\x00" + failure.reason
		group, ok := groups[groupKey]
		if !ok {
			group = &encryptionFailureGroup{keyName: failure.keyName, reason: failure.reason}
			groups[groupKey] = group
			keys = append(keys, groupKey)
		}
		group.posts = append(group.posts, failure.postPath)
	}

	sort.Strings(keys)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		group := groups[key]
		sample := group.posts
		more := false
		if len(sample) > maxEncryptionFailureSamples {
			sample = sample[:maxEncryptionFailureSamples]
			more = true
		}
		line := fmt.Sprintf("key %q: %s (posts: %d, sample: %s", group.keyName, group.reason, len(group.posts), strings.Join(sample, ", "))
		if more {
			line += ", ..."
		}
		line += ")"
		lines = append(lines, line)
	}

	return lines
}

func (p *EncryptionPlugin) validatePasswordPolicy(password string) error {
	if !p.enforceStrength {
		return nil
	}
	return encryption.ValidatePassword(password, p.minPasswordLength, p.minEstimatedCrackDuration)
}

func filterEncryptedPostsForServe(m *lifecycle.Manager, posts []*models.Post) []*models.Post {
	if !lifecycle.IsServeIncremental(m) {
		return posts
	}
	affected := lifecycle.GetServeAffectedPaths(m)
	if len(affected) == 0 {
		return posts
	}
	filtered := posts[:0]
	for _, post := range posts {
		if affected[post.Path] {
			filtered = append(filtered, post)
		}
	}
	return filtered
}

// encryptPostWithCacheResult reports a real wrapper generation, independently
// of best-effort cache storage. Nonempty hits do not invalidate full pages.
func (p *EncryptionPlugin) encryptPostWithCacheResult(post *models.Post, cache *buildcache.Cache) (bool, error) {
	keyName := post.SecretKey
	if keyName == "" {
		keyName = p.defaultKey
	}
	password, err := p.getKeyPassword(keyName)
	if err != nil {
		return false, err
	}
	encryptedHash := computeEncryptedHash(post.ArticleHTML, keyName, password, p.decryptionHint, post.Path)

	if cache != nil {
		if cached := cache.GetCachedEncryptedHTML(post.Path, encryptedHash); cached != "" {
			post.ArticleHTML = cached
			// A warm build loads raw Content again before restoring the cached
			// wrapper. Apply the same scrub as a cold build so feed/template
			// consumers cannot observe the plaintext through the cache path.
			p.scrubPrivateMetadata(post)
			post.Set("has_encrypted_content", true)
			if keyName != "" {
				post.Set("encryption_key_name", keyName)
			}
			templates.InvalidatePost(post)
			return false, nil
		}
	}

	err = p.encryptPost(post)
	if err != nil {
		return false, err
	}
	if cache != nil {
		// A failed later full-page write must not leave old-password HTML
		// available when the next build hits this newly generated wrapper.
		cache.InvalidateFullHTML(post.Path)
		//nolint:errcheck // best-effort caching
		cache.CacheEncryptedHTML(post.Path, encryptedHash, post.ArticleHTML)
	}
	return true, nil
}

// Bump whenever the wrapper markup or browser crypto format changes.
// This is a derived identity revision, not a persisted cache schema version.
const encryptedWrapperRevision = "encrypted-wrapper-v2"

func computeEncryptedHash(articleHTML, keyName, password, hint, sourcePath string) string {
	var b strings.Builder
	b.WriteString(articleHTML)
	b.WriteByte('\x00')
	b.WriteString(keyName)
	b.WriteByte('\x00')
	b.WriteString(password)
	b.WriteByte('\x00')
	b.WriteString(hint)
	b.WriteByte('\x00')
	b.WriteString(sourcePath)
	b.WriteByte('\x00')
	b.WriteString(encryptedWrapperRevision)
	return buildcache.ContentHash(b.String())
}

// shouldEncrypt determines if a post should have its content encrypted.
func (p *EncryptionPlugin) shouldEncrypt(post *models.Post) bool {
	if post.Skip || post.Draft {
		return false
	}

	// Must be private
	if !post.Private {
		return false
	}

	// Must have a secret_key or default key must be set
	if post.SecretKey == "" && p.defaultKey == "" {
		return false
	}

	// Must have content to encrypt
	if post.ArticleHTML == "" {
		return false
	}

	return true
}

// encryptPost encrypts the post's ArticleHTML content and scrubs all
// plaintext metadata to prevent content leaks through descriptions,
// structured data, search indexes, and other output channels.
func (p *EncryptionPlugin) encryptPost(post *models.Post) error {
	keyName := post.SecretKey
	if keyName == "" {
		keyName = p.defaultKey
	}

	password, err := p.getKeyPassword(keyName)
	if err != nil {
		// This should not happen because we already validated all keys in Render(),
		// but if it does, fail hard to protect private content.
		return &EncryptionBuildError{
			Posts: []string{post.Path},
			Msg:   fmt.Sprintf("encryption key %q not found", keyName),
		}
	}

	// Encrypt the article HTML
	encryptedContent, err := encryption.Encrypt([]byte(post.ArticleHTML), password)
	if err != nil {
		return fmt.Errorf("failed to encrypt post %q: %w", post.Path, err)
	}

	// Build the decryption hint HTML
	hintHTML := ""
	if p.decryptionHint != "" {
		hintHTML = fmt.Sprintf(`<p class="encrypted-content__hint">%s</p>`, escapeHTML(p.decryptionHint))
	}

	// Generate unique ID for accessibility (based on post path hash)
	inputID := fmt.Sprintf("decrypt-input-%d", hashString(post.Path))
	checkboxID := fmt.Sprintf("decrypt-remember-%d", hashString(post.Path))

	// Replace ArticleHTML with encrypted wrapper
	// Includes:
	// - data-key-name for multi-post unlock feature
	// - data-pagefind-ignore to prevent search indexing of encrypted content
	// - ARIA labels for accessibility
	// - Remember me checkbox for explicit session storage opt-in
	post.ArticleHTML = fmt.Sprintf(`<div class="encrypted-content" data-encrypted="%s" data-key-name="%s" data-pagefind-ignore role="region" aria-label="Encrypted content">
  <div class="encrypted-content__locked">
    <div class="encrypted-content__icon" aria-hidden="true">🔒</div>
    <h3 class="encrypted-content__title" id="encrypt-title-%d">Encrypted Content</h3>
    <p class="encrypted-content__message">This content is encrypted. Enter the password to decrypt and view.</p>
    %s
    <form class="encrypted-content__form" aria-labelledby="encrypt-title-%d" action="#" method="post" autocomplete="on">
      <label for="%s" class="sr-only">Password</label>
      <input type="password" id="%s" name="password" class="encrypted-content__input" placeholder="Enter password" autocomplete="current-password" autocapitalize="off" autocorrect="off" spellcheck="false" aria-describedby="encrypt-error-%d">
      <button type="button" class="encrypted-content__button" aria-busy="false">Decrypt</button>
    </form>
    <label class="encrypted-content__remember-label">
      <input type="checkbox" id="%s" name="remember" class="encrypted-content__remember" aria-describedby="remember-desc-%d">
      <span>Remember for this session</span>
    </label>
    <span id="remember-desc-%d" class="sr-only">If checked, the password will be saved in your browser for this session only. It will be cleared when you close the browser.</span>
    <p id="encrypt-error-%d" class="encrypted-content__error" style="display: none;" role="alert" aria-live="assertive"></p>
  </div>
  <div class="encrypted-content__decrypted" style="display: none;" tabindex="-1"></div>
</div>`, encryptedContent, escapeHTML(keyName), hashString(post.Path), hintHTML,
		hashString(post.Path), inputID, inputID, hashString(post.Path),
		checkboxID, hashString(post.Path), hashString(post.Path), hashString(post.Path))

	// Scrub all plaintext metadata to prevent content leaks.
	// This is critical: without scrubbing, private content leaks through
	// descriptions, feeds, search indexes, structured data, and 404 indexes.
	p.scrubPrivateMetadata(post)

	// Mark post as having encrypted content (for template to include decryption script)
	post.Set("has_encrypted_content", true)
	post.Set("encryption_key_name", keyName)

	// Invalidate the post map cache so templates pick up the new Extra fields
	templates.InvalidatePost(post)

	return nil
}

// scrubPrivateMetadata removes content-derived plaintext from a private post
// while preserving user-provided metadata that is safe to display publicly.
//
// Preserved (intentionally public):
//   - Title -> kept for page cards, feed listings, HTML <title>, and SEO
//   - Description (if explicitly set in frontmatter) -> kept for cards and meta tags
//   - Structured data (JSON-LD, OpenGraph, Twitter) -> kept; derived from title/description only
//
// Scrubbed (contains private content):
//   - Content (raw markdown) -> leaked via Atom/JSON feeds, description generation
//   - Description (if NOT set in frontmatter) -> cleared to nil since the description
//     plugin skips private posts, so a nil description here means none was provided
//   - Inlinks/Outlinks text -> leaked via link analysis output
func (p *EncryptionPlugin) scrubPrivateMetadata(post *models.Post) {
	// Title is preserved — users expect encrypted pages to show their title
	// in cards, feed listings, HTML <title>, and navigation. The title is
	// considered public metadata, not private content.

	// Clear raw markdown content — prevents leaks through feed content_text,
	// Atom <content type="text">, and any future content consumers.
	post.Content = ""

	// Description: preserve if explicitly set in frontmatter (user chose to
	// make it public). The description plugin already skips private posts,
	// so if Description is non-nil here, it was explicitly provided in
	// frontmatter. If nil, leave it nil — no auto-generated description
	// should be created for encrypted posts.

	// Structured data is preserved — it contains only title, description,
	// dates, author, and URL. No article body or raw content is included.
	// Since we preserve title and explicit description, the structured data
	// remains accurate and useful for SEO.

	// Clear inlinks/outlinks source text — prevents content-derived
	// anchor text from leaking through link analysis.
	for _, link := range post.Inlinks {
		if link != nil {
			link.SourceText = ""
			link.TargetText = ""
		}
	}
	for _, link := range post.Outlinks {
		if link != nil {
			link.SourceText = ""
			link.TargetText = ""
		}
	}
}

// escapeHTML escapes special characters for safe HTML output.
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// hashString generates a simple numeric hash for a string.
// Used for generating unique IDs in HTML elements.
func hashString(s string) uint32 {
	var hash uint32
	for _, c := range s {
		hash = hash*31 + uint32(c)
	}
	return hash
}

// getModelsConfig attempts to extract models.Config from lifecycle.Config.
// This handles the conversion between the lifecycle config and models config.
func getModelsConfig(config *lifecycle.Config) (*models.Config, bool) {
	if config == nil || config.Extra == nil {
		return nil, false
	}

	// The models.Config is stored in Extra["models_config"] (set by core.go)
	if modelsConfig, ok := config.Extra["models_config"].(*models.Config); ok {
		return modelsConfig, true
	}

	return nil, false
}

// Ensure EncryptionPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*EncryptionPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*EncryptionPlugin)(nil)
	_ lifecycle.TransformPlugin = (*EncryptionPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*EncryptionPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*EncryptionPlugin)(nil)
)
