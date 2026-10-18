package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/spf13/cobra"
)

var (
	encryptionPasswordLength int
	encryptionCheckKey       string
	encryptionDryRun         bool
	decryptionDryRun         bool
	encryptionWorkers        int
)

var encryptionCmd = &cobra.Command{
	Use:   "encryption",
	Short: "Utilities for encryption keys and passwords",
	Long: `Encryption utilities help you manage passwords and source-encrypted private posts.
`}

var generatePasswordCmd = &cobra.Command{
	Use:     "generate-password",
	Aliases: []string{"gen"},
	Short:   "Generate a policy-compliant encryption password",
	Long: `Generate a password that satisfies the default encryption policy (>=14 chars, strong entropy).

The password is printed to stdout only, making it easy to pipe into secret stores or copy it from your terminal.
`,
	Args: cobra.NoArgs,
	RunE: runGeneratePasswordCommand,
}

var checkPasswordCmd = &cobra.Command{
	Use:   "check",
	Short: "Check configured encryption key strength",
	Long: `Check configured encryption keys against the active policy.

By default this checks every key required by your config (default_key and private_tags mappings).
Use --key to check one specific key name.
`,
	Args: cobra.NoArgs,
	RunE: runCheckPasswordCommand,
}

var encryptPostsCmd = &cobra.Command{
	Use:   "encrypt-posts",
	Short: "Encrypt all private Markdown source bodies",
	Long: `Encrypt all private Markdown source bodies matched by the active content glob configuration.

Posts are encrypted in place by default. Draft, skipped, public, and already
source-encrypted posts are reported but not rewritten. Use --dry-run to preview
changes without modifying files. Use --workers to bound concurrent preparation.
`,
	Args: cobra.NoArgs,
	RunE: runEncryptPostsCommand,
}

var decryptPostsCmd = &cobra.Command{
	Use:   "decrypt-posts [path...]",
	Short: "Decrypt source-encrypted Markdown bodies back to plaintext",
	Long: `Decrypt source-encrypted Markdown bodies back to plaintext.

This is the inverse of ` + "`encryption encrypt-posts`" + `. With no arguments it scans the
active content glob configuration; with explicit paths it only processes those files.

Files are decrypted in place by default. Files that are not source-encrypted are
reported but not rewritten. Use --dry-run to preview changes without modifying files.

The key name is read from the encrypted source marker, falling back to
encryption.default_key. The password is read from MARKATA_GO_ENCRYPTION_KEY_<KEY>.
Use --workers to bound concurrent preparation.
`,
	Args: cobra.ArbitraryArgs,
	RunE: runDecryptPostsCommand,
}

// encryptCmd is a deliberately hidden shortcut for encrypt-posts.
var encryptCmd = &cobra.Command{
	Use:    "encrypt",
	Short:  "Encrypt all private Markdown source bodies",
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE:   runEncryptPostsCommand,
}

// decryptCmd is a deliberately hidden shortcut for decrypt-posts.
var decryptCmd = &cobra.Command{
	Use:    "decrypt [path...]",
	Short:  "Decrypt source-encrypted Markdown bodies back to plaintext",
	Hidden: true,
	Args:   cobra.ArbitraryArgs,
	RunE:   runDecryptPostsCommand,
}

func init() {
	encryptionCmd.AddCommand(generatePasswordCmd, checkPasswordCmd, encryptPostsCmd, decryptPostsCmd)
	generatePasswordCmd.Flags().IntVar(&encryptionPasswordLength, "length", encryption.DefaultMinPasswordLength, "password length (must be at least the configured minimum)")
	checkPasswordCmd.Flags().StringVar(&encryptionCheckKey, "key", "", "specific key name to check (default: all required keys)")
	encryptPostsCmd.Flags().BoolVar(&encryptionDryRun, "dry-run", false, "report files that would be encrypted without modifying them")
	encryptPostsCmd.Flags().IntVar(&encryptionWorkers, "workers", 0, "number of concurrent workers (0 uses GOMAXPROCS)")
	decryptPostsCmd.Flags().BoolVar(&decryptionDryRun, "dry-run", false, "report files that would be decrypted without modifying them")
	decryptPostsCmd.Flags().IntVar(&encryptionWorkers, "workers", 0, "number of concurrent workers (0 uses GOMAXPROCS)")
	encryptCmd.Flags().BoolVar(&encryptionDryRun, "dry-run", false, "report files that would be encrypted without modifying them")
	encryptCmd.Flags().IntVar(&encryptionWorkers, "workers", 0, "number of concurrent workers (0 uses GOMAXPROCS)")
	decryptCmd.Flags().BoolVar(&decryptionDryRun, "dry-run", false, "report files that would be decrypted without modifying them")
	decryptCmd.Flags().IntVar(&encryptionWorkers, "workers", 0, "number of concurrent workers (0 uses GOMAXPROCS)")
	rootCmd.AddCommand(encryptionCmd, encryptCmd, decryptCmd)
}

func runGeneratePasswordCommand(cmd *cobra.Command, _ []string) error {
	password, err := encryption.GeneratePassword(encryptionPasswordLength)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), password)
	return nil
}

func runCheckPasswordCommand(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	results, minDuration, minLength, err := evaluateEncryptionKeyPolicy(cfg, encryptionCheckKey)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No encryption keys configured to check.")
		return nil
	}

	failures := 0
	fmt.Fprintf(cmd.OutOrStdout(), "Policy: min_length=%d, min_estimated_crack_time=%s\n", minLength, formatCrackDurationHuman(minDuration))
	for _, result := range results {
		if result.Err != nil {
			failures++
			if result.Configured {
				fmt.Fprintf(cmd.OutOrStdout(), "FAIL %s (%s): %s (estimated=%s)\n", result.KeyName, result.EnvName, result.Err, formatCrackDurationHuman(result.Estimated))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "FAIL %s (%s): %s\n", result.KeyName, result.EnvName, result.Err)
			}
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "PASS %s (%s): estimated=%s\n", result.KeyName, result.EnvName, formatCrackDurationHuman(result.Estimated))
	}

	if failures > 0 {
		return fmt.Errorf("%d encryption key(s) failed policy", failures)
	}

	if !cfg.Encryption.EnforceStrength {
		fmt.Fprintln(cmd.OutOrStdout(), "Warning: encryption.enforce_strength=false, builds will not enforce this policy.")
	}

	return nil
}

func runEncryptPostsCommand(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if !cfg.Encryption.Enabled {
		return fmt.Errorf("encryption is disabled in config")
	}

	results, minDuration, minLength, err := evaluateEncryptionKeyPolicy(cfg, "")
	if err != nil {
		return err
	}
	if err := failOnEncryptionKeyPolicyFailures(results, minDuration, minLength); err != nil {
		return err
	}

	files, err := encryptionContentFiles(cfg)
	if err != nil {
		return err
	}

	sourceCache, cacheErr := buildcache.LoadSourceEncryptionCache(buildcache.DefaultCacheDir)
	if cacheErr != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "Warning: %v; unchanged files may receive new ciphertext.\n", cacheErr)
	}
	prepared, err := prepareSourceFiles(files, encryptionWorkers, func(file string) (preparedEncryptPost, error) {
		return prepareEncryptPostSourceFileWithCache(file, cfg, sourceCache)
	})
	if err != nil {
		return err
	}

	stats := encryptPostsStats{}
	for index := range prepared {
		result := prepared[index].result
		stats.add(result)
		if result.Action == encryptPostActionEncrypted {
			if encryptionDryRun {
				fmt.Fprintln(cmd.OutOrStdout(), formatEncryptionProgress("WOULD ENCRYPT", currentLogTheme.Warning, result.Path, result.KeyName))
			}
		}
	}

	if !encryptionDryRun {
		documents := make([]sourceDocument, 0, stats.Encrypted)
		for index := range prepared {
			if prepared[index].result.Action == encryptPostActionEncrypted {
				documents = append(documents, prepared[index].document)
			}
		}
		cacheSourceEncryptionDocuments(sourceCache, documents)
		if err := sourceCache.Save(); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Warning: save source encryption cache: %v; future round trips may generate new ciphertext.\n", err)
		}
		if err := writeSourceDocuments(documents); err != nil {
			return err
		}
		for index := range prepared {
			if prepared[index].result.Action == encryptPostActionEncrypted {
				fmt.Fprintln(cmd.OutOrStdout(), formatEncryptionProgress("ENCRYPTED", currentLogTheme.Success, prepared[index].result.Path, prepared[index].result.KeyName))
			}
		}
	}

	action := "Encrypted"
	if encryptionDryRun {
		action = "Would encrypt"
	}
	fmt.Fprintln(cmd.OutOrStdout(), formatEncryptionSummary(action, currentLogTheme.Success, fmt.Sprintf("%d private post(s); skipped %d already encrypted, %d public, %d draft/skip.",
		stats.Encrypted, stats.AlreadyEncrypted, stats.Public, stats.DraftOrSkip)))

	return nil
}

func runDecryptPostsCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	files, err := decryptionTargetFiles(cfg, args)
	if err != nil {
		return err
	}

	sourceCache, cacheErr := buildcache.LoadSourceEncryptionCache(buildcache.DefaultCacheDir)
	if cacheErr != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "Warning: %v; unchanged files may receive new ciphertext.\n", cacheErr)
	}
	prepared, err := prepareSourceFiles(files, encryptionWorkers, func(file string) (preparedDecryptPost, error) {
		return prepareDecryptPostSourceFile(file, cfg)
	})
	if err != nil {
		return err
	}

	stats := decryptPostsStats{}
	for index := range prepared {
		result := prepared[index].result
		stats.add(result)
		if result.Action == decryptPostActionDecrypted {
			if decryptionDryRun {
				fmt.Fprintln(cmd.OutOrStdout(), formatEncryptionProgress("WOULD DECRYPT", currentLogTheme.Warning, result.Path, result.KeyName))
			}
		}
	}

	if !decryptionDryRun {
		documents := make([]sourceDocument, 0, stats.Decrypted)
		for index := range prepared {
			if prepared[index].result.Action == decryptPostActionDecrypted {
				documents = append(documents, prepared[index].document)
			}
		}
		cacheSourceEncryptionDocuments(sourceCache, documents)
		if err := sourceCache.Save(); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Warning: save source encryption cache: %v; future round trips may generate new ciphertext.\n", err)
		}
		if err := writeSourceDocuments(documents); err != nil {
			return err
		}
		for index := range prepared {
			if prepared[index].result.Action == decryptPostActionDecrypted {
				fmt.Fprintln(cmd.OutOrStdout(), formatEncryptionProgress("DECRYPTED", currentLogTheme.Success, prepared[index].result.Path, prepared[index].result.KeyName))
			}
		}
	}

	action := "Decrypted"
	if decryptionDryRun {
		action = "Would decrypt"
	}
	fmt.Fprintln(cmd.OutOrStdout(), formatEncryptionSummary(action, currentLogTheme.Success,
		fmt.Sprintf("%d post(s); skipped %d not encrypted.", stats.Decrypted, stats.NotEncrypted)))

	return nil
}

// decryptionTargetFiles resolves explicit path arguments, falling back to the
// active content glob configuration when no paths are supplied.
func decryptionTargetFiles(cfg *models.Config, args []string) ([]string, error) {
	if len(args) == 0 {
		return encryptionContentFiles(cfg)
	}

	paths := make([]string, 0, len(args))
	for _, arg := range args {
		abs, err := filepath.Abs(arg)
		if err != nil {
			return nil, fmt.Errorf("resolve %s: %w", arg, err)
		}
		info, err := os.Stat(abs)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", arg, err)
		}
		if !info.IsDir() {
			paths = append(paths, abs)
			continue
		}
		err = filepath.WalkDir(abs, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", arg, err)
		}
	}
	return paths, nil
}

type decryptPostAction string

const (
	decryptPostActionDecrypted    decryptPostAction = "decrypted"
	decryptPostActionNotEncrypted decryptPostAction = "not-encrypted"
)

type decryptPostResult struct {
	Path    string
	KeyName string
	Action  decryptPostAction
}

type decryptPostsStats struct {
	Decrypted    int
	NotEncrypted int
}

func (s *decryptPostsStats) add(result decryptPostResult) {
	switch result.Action {
	case decryptPostActionDecrypted:
		s.Decrypted++
	case decryptPostActionNotEncrypted:
		s.NotEncrypted++
	}
}

type preparedDecryptPost struct {
	result   decryptPostResult
	document sourceDocument
}

func decryptPostSourceFile(path string, cfg *models.Config, dryRun bool) (decryptPostResult, error) {
	prepared, err := prepareDecryptPostSourceFile(path, cfg)
	if err != nil {
		return decryptPostResult{}, err
	}
	if dryRun || prepared.result.Action != decryptPostActionDecrypted {
		return prepared.result, nil
	}
	if err := writeSourceDocuments([]sourceDocument{prepared.document}); err != nil {
		return decryptPostResult{}, err
	}
	return prepared.result, nil
}

func prepareDecryptPostSourceFile(path string, cfg *models.Config) (preparedDecryptPost, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return preparedDecryptPost{}, fmt.Errorf("read %s: %w", path, err)
	}
	content := string(contentBytes)

	_, body, rawFrontmatter, err := plugins.ParseFrontmatterWithRaw(content)
	if err != nil {
		return preparedDecryptPost{}, fmt.Errorf("parse frontmatter %s: %w", path, err)
	}
	if !encryption.IsSourceEncrypted(body) {
		return preparedDecryptPost{result: decryptPostResult{Path: path, Action: decryptPostActionNotEncrypted}}, nil
	}

	envelope, _, err := encryption.ParseSourceEnvelope(body)
	if err != nil {
		return preparedDecryptPost{}, fmt.Errorf("parse encrypted source %s: %w", path, err)
	}

	keyName := strings.TrimSpace(envelope.KeyName)
	if keyName == "" {
		keyName = strings.TrimSpace(cfg.Encryption.DefaultKey)
	}
	if keyName == "" {
		return preparedDecryptPost{}, fmt.Errorf("encrypted post %s has no key name; set encryption.default_key", path)
	}
	envName := plugins.EncryptionEnvPrefix + strings.ToUpper(keyName)
	password := os.Getenv(envName)
	if password == "" {
		return preparedDecryptPost{}, fmt.Errorf("encrypted post %s requires %s", path, envName)
	}

	plaintext, _, err := encryption.DecryptSourceMarkdown(body, password)
	if err != nil {
		return preparedDecryptPost{}, fmt.Errorf("decrypt source body for %s: %w", path, err)
	}
	return preparedDecryptPost{
		result: decryptPostResult{Path: path, KeyName: keyName, Action: decryptPostActionDecrypted},
		document: sourceDocument{
			path:                    path,
			original:                content,
			content:                 decryptedSourceDocument(rawFrontmatter, plaintext),
			mode:                    sourceFileMode(path),
			sourceEncryptionKeyName: keyName,
			sourceEncryptedBody:     body,
		},
	}, nil
}

func decryptedSourceDocument(rawFrontmatter, body string) string {
	if rawFrontmatter == "" {
		return body
	}
	return "---\n" + rawFrontmatter + "\n---\n" + body
}

func failOnEncryptionKeyPolicyFailures(results []encryptionKeyPolicyResult, minDuration time.Duration, minLength int) error {
	for _, result := range results {
		if result.Err != nil {
			return fmt.Errorf("encryption key %q failed policy (%s, min_length=%d, min_estimated_crack_time=%s): %w",
				result.KeyName, result.EnvName, minLength, formatCrackDurationHuman(minDuration), result.Err)
		}
	}
	return nil
}

func encryptionContentFiles(cfg *models.Config) ([]string, error) {
	lifecycleConfig := lifecycle.NewConfig()
	lifecycleConfig.ContentDir = "."
	lifecycleConfig.GlobPatterns = append([]string{}, cfg.GlobConfig.Patterns...)
	lifecycleConfig.Extra["use_gitignore"] = cfg.GlobConfig.UseGitignore

	manager := lifecycle.NewManager()
	manager.SetConfig(lifecycleConfig)

	globPlugin := plugins.NewGlobPlugin()
	if err := globPlugin.Configure(manager); err != nil {
		return nil, fmt.Errorf("configure content glob: %w", err)
	}
	if err := globPlugin.Glob(manager); err != nil {
		return nil, fmt.Errorf("scan content files: %w", err)
	}

	baseDir, err := filepath.Abs(lifecycleConfig.ContentDir)
	if err != nil {
		return nil, fmt.Errorf("resolve content dir: %w", err)
	}

	files := manager.Files()
	paths := make([]string, 0, len(files))
	for _, file := range files {
		if filepath.IsAbs(file) {
			paths = append(paths, file)
			continue
		}
		paths = append(paths, filepath.Join(baseDir, file))
	}
	return paths, nil
}

type encryptPostAction string

const (
	encryptPostActionEncrypted        encryptPostAction = "encrypted"
	encryptPostActionAlreadyEncrypted encryptPostAction = "already-encrypted"
	encryptPostActionPublic           encryptPostAction = "public"
	encryptPostActionDraftOrSkip      encryptPostAction = "draft-or-skip"
)

type encryptPostResult struct {
	Path    string
	KeyName string
	Action  encryptPostAction
}

type encryptPostsStats struct {
	Encrypted        int
	AlreadyEncrypted int
	Public           int
	DraftOrSkip      int
}

type preparedEncryptPost struct {
	result   encryptPostResult
	document sourceDocument
}

func (s *encryptPostsStats) add(result encryptPostResult) {
	switch result.Action {
	case encryptPostActionEncrypted:
		s.Encrypted++
	case encryptPostActionAlreadyEncrypted:
		s.AlreadyEncrypted++
	case encryptPostActionPublic:
		s.Public++
	case encryptPostActionDraftOrSkip:
		s.DraftOrSkip++
	}
}

func encryptPostSourceFile(path string, cfg *models.Config, dryRun bool) (encryptPostResult, error) {
	prepared, err := prepareEncryptPostSourceFile(path, cfg)
	if err != nil {
		return encryptPostResult{}, err
	}
	if dryRun || prepared.result.Action != encryptPostActionEncrypted {
		return prepared.result, nil
	}
	if err := writeSourceDocuments([]sourceDocument{prepared.document}); err != nil {
		return encryptPostResult{}, err
	}
	return prepared.result, nil
}

func prepareEncryptPostSourceFile(path string, cfg *models.Config) (preparedEncryptPost, error) {
	return prepareEncryptPostSourceFileWithCache(path, cfg, nil)
}

func prepareEncryptPostSourceFileWithCache(path string, cfg *models.Config, cache *buildcache.SourceEncryptionCache) (preparedEncryptPost, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return preparedEncryptPost{}, fmt.Errorf("read %s: %w", path, err)
	}
	content := string(contentBytes)

	_, body, rawFrontmatter, err := plugins.ParseFrontmatterWithRaw(content)
	if err != nil {
		return preparedEncryptPost{}, fmt.Errorf("parse frontmatter %s: %w", path, err)
	}
	if encryption.IsSourceEncrypted(body) {
		return preparedEncryptPost{result: encryptPostResult{Path: path, Action: encryptPostActionAlreadyEncrypted}}, nil
	}

	post, err := plugins.ParsePostFromContentWithConfig(path, content, cfg)
	if err != nil {
		return preparedEncryptPost{}, fmt.Errorf("parse %s: %w", path, err)
	}
	applyEncryptPostsPrivateTags(post, cfg)
	if post.Draft || post.Skip {
		return preparedEncryptPost{result: encryptPostResult{Path: path, Action: encryptPostActionDraftOrSkip}}, nil
	}
	if !post.Private {
		return preparedEncryptPost{result: encryptPostResult{Path: path, Action: encryptPostActionPublic}}, nil
	}

	keyName := strings.TrimSpace(post.SecretKey)
	if keyName == "" {
		keyName = strings.TrimSpace(cfg.Encryption.DefaultKey)
	}
	if keyName == "" {
		return preparedEncryptPost{}, fmt.Errorf("private post %s has no encryption key; set secret_key or encryption.default_key", path)
	}
	password := os.Getenv(plugins.EncryptionEnvPrefix + strings.ToUpper(keyName))
	if password == "" {
		return preparedEncryptPost{}, fmt.Errorf("private post %s requires %s%s", path, plugins.EncryptionEnvPrefix, strings.ToUpper(keyName))
	}
	if err := validateEncryptPostsPassword(password, cfg); err != nil {
		return preparedEncryptPost{}, fmt.Errorf("private post %s key %q failed policy: %w", path, keyName, err)
	}

	encryptedBody := ""
	if cache != nil {
		encryptedBody, _ = cache.Get(path, body, keyName, password)
	}
	if encryptedBody == "" {
		encryptedBody, err = encryption.EncryptSourceMarkdown(body, keyName, password)
		if err != nil {
			return preparedEncryptPost{}, fmt.Errorf("encrypt source body for %s: %w", path, err)
		}
	}
	return preparedEncryptPost{
		result: encryptPostResult{Path: path, KeyName: keyName, Action: encryptPostActionEncrypted},
		document: sourceDocument{
			path:                    path,
			original:                content,
			content:                 encryptedSourceDocument(rawFrontmatter, encryptedBody),
			mode:                    sourceFileMode(path),
			sourceEncryptionKeyName: keyName,
			sourceEncryptedBody:     encryptedBody,
		},
	}, nil
}

func sourceFileMode(path string) os.FileMode {
	if stat, err := os.Stat(path); err == nil {
		return stat.Mode().Perm()
	}
	return 0o600
}

type sourceDocument struct {
	path                    string
	original                string
	content                 string
	mode                    os.FileMode
	sourceEncryptionKeyName string
	sourceEncryptedBody     string
}

func cacheSourceEncryptionDocuments(cache *buildcache.SourceEncryptionCache, documents []sourceDocument) {
	if cache == nil {
		return
	}
	for _, document := range documents {
		cache.Put(document.path, document.sourceEncryptionKeyName, document.sourceEncryptedBody)
	}
}

// writeSourceDocuments writes a prepared batch only when every source still
// matches the content used during preparation. Each replacement is atomic; if
// a later write fails, previously written documents are restored.
func writeSourceDocuments(documents []sourceDocument) error {
	for _, document := range documents {
		if err := ensureSourceUnchanged(document); err != nil {
			return err
		}
	}

	written := make([]sourceDocument, 0, len(documents))
	for _, document := range documents {
		if err := ensureSourceUnchanged(document); err != nil {
			rollbackErr := restoreSourceDocuments(written)
			if rollbackErr != nil {
				return fmt.Errorf("source changed; rollback failed: %w", errors.Join(err, rollbackErr))
			}
			return err
		}
		if err := writeSourceDocumentAtomically(document.path, document.content, document.mode); err != nil {
			rollbackErr := restoreSourceDocuments(written)
			if rollbackErr != nil {
				return fmt.Errorf("write source %s; rollback failed: %w", document.path, errors.Join(err, rollbackErr))
			}
			return fmt.Errorf("write source %s: %w", document.path, err)
		}
		written = append(written, document)
	}
	return nil
}

func restoreSourceDocuments(documents []sourceDocument) error {
	for i := len(documents) - 1; i >= 0; i-- {
		document := documents[i]
		current, err := os.ReadFile(document.path)
		if err != nil {
			return fmt.Errorf("read source before restore %s: %w", document.path, err)
		}
		if string(current) != document.content {
			return fmt.Errorf("source changed before restore: %s", document.path)
		}
		if err := writeSourceDocumentAtomically(document.path, document.original, document.mode); err != nil {
			return fmt.Errorf("restore source %s: %w", document.path, err)
		}
	}
	return nil
}

func ensureSourceUnchanged(document sourceDocument) error {
	current, err := os.ReadFile(document.path)
	if err != nil {
		return fmt.Errorf("read source before write %s: %w", document.path, err)
	}
	if string(current) != document.original {
		return fmt.Errorf("source changed during encryption preparation: %s", document.path)
	}
	return nil
}

func writeSourceDocumentAtomically(path, content string, mode os.FileMode) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".markata-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if err := temporary.Chmod(mode); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.WriteString(content); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

func formatEncryptionProgress(action, actionColor, path, keyName string) string {
	return fmt.Sprintf("%s %s key=%s",
		colorizeOutput(action, actionColor),
		colorizeOutput(path, currentLogTheme.Component),
		colorizeOutput(keyName, currentLogTheme.Warning))
}

func formatEncryptionSummary(action, actionColor, summary string) string {
	return fmt.Sprintf("%s %s", colorizeOutput(action, actionColor), summary)
}

// prepareSourceFiles runs CPU-bound source transformations concurrently while
// retaining input order. Callers must complete this phase before writing files
// so a preparation error never leaves a partially transformed repository.
func prepareSourceFiles[T any](paths []string, workers int, prepare func(string) (T, error)) ([]T, error) {
	results := make([]T, len(paths))
	errs := make([]error, len(paths))
	if len(paths) == 0 {
		return results, nil
	}
	if workers < 0 {
		return nil, fmt.Errorf("workers must be zero or greater")
	}
	if workers == 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	if workers > len(paths) {
		workers = len(paths)
	}

	jobs := make(chan int)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				results[index], errs[index] = prepare(paths[index])
			}
		}()
	}
	for index := range paths {
		jobs <- index
	}
	close(jobs)
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func validateEncryptPostsPassword(password string, cfg *models.Config) error {
	minLength := cfg.Encryption.MinPasswordLength
	if minLength == 0 {
		minLength = encryption.DefaultMinPasswordLength
	}
	minDurationValue := cfg.Encryption.MinEstimatedCrackTime
	if minDurationValue == "" {
		minDurationValue = encryption.DefaultMinEstimatedCrackTime
	}
	minDuration, err := encryption.ParseEstimatedCrackDuration(minDurationValue)
	if err != nil {
		return fmt.Errorf("invalid encryption.min_estimated_crack_time: %w", err)
	}
	return encryption.ValidatePassword(password, minLength, minDuration)
}

func applyEncryptPostsPrivateTags(post *models.Post, cfg *models.Config) {
	if post == nil || cfg == nil || post.Skip || post.Draft {
		return
	}
	if post.IsExplicitlyPublic() {
		return
	}
	for _, tag := range post.Tags {
		if keyName, ok := cfg.Encryption.PrivateTags[strings.ToLower(tag)]; ok {
			post.Private = true
			if post.SecretKey == "" {
				post.SecretKey = keyName // pragma: allowlist secret
			}
			return
		}
	}
	if post.Template == "" {
		return
	}
	if keyName, ok := cfg.Encryption.PrivateTags[strings.ToLower(post.Template)]; ok {
		post.Private = true
		if post.SecretKey == "" {
			post.SecretKey = keyName // pragma: allowlist secret
		}
	}
}

func encryptedSourceDocument(rawFrontmatter, encryptedBody string) string {
	if rawFrontmatter == "" {
		return encryptedBody
	}
	return "---\n" + rawFrontmatter + "\n---\n" + encryptedBody
}

func formatCrackDurationHuman(d time.Duration) string {
	if d <= 0 {
		return "0s"
	}
	if d < time.Second {
		return "<1s"
	}

	const year = 365 * 24 * time.Hour
	if d >= year {
		years := float64(d) / float64(year)
		if years >= 100 {
			return fmt.Sprintf("%.0fy", years)
		}
		return fmt.Sprintf("%.1fy", years)
	}

	const day = 24 * time.Hour
	if d >= day {
		days := float64(d) / float64(day)
		return fmt.Sprintf("%.1fd", days)
	}

	if d >= time.Hour {
		hours := float64(d) / float64(time.Hour)
		return fmt.Sprintf("%.1fh", hours)
	}

	return d.Round(time.Second).String()
}
