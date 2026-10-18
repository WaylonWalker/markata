package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/spf13/cobra"
)

func TestGeneratePasswordCommand_DefaultLength(t *testing.T) {
	originalLength := encryptionPasswordLength
	defer func() { encryptionPasswordLength = originalLength }()
	encryptionPasswordLength = encryption.DefaultMinPasswordLength

	buf := bytes.NewBuffer(nil)
	generatePasswordCmd.SetOut(buf)

	if err := runGeneratePasswordCommand(generatePasswordCmd, nil); err != nil {
		t.Fatalf("runGeneratePasswordCommand() error = %v", err)
	}

	password := strings.TrimSpace(buf.String())
	if len(password) != encryption.DefaultMinPasswordLength {
		t.Errorf("password length = %d, want %d", len(password), encryption.DefaultMinPasswordLength)
	}
	if err := encryption.ValidatePassword(password, encryption.DefaultMinPasswordLength, encryption.DefaultMinEstimatedCrackDuration); err != nil {
		t.Fatalf("generated password failed validation: %v", err)
	}
}

func TestEncryptionShortcuts_HiddenAndConfigured(t *testing.T) {
	for _, shortcut := range []*cobra.Command{encryptCmd, decryptCmd} {
		if !shortcut.Hidden {
			t.Errorf("%s shortcut must be hidden", shortcut.Name())
		}
		if shortcut.RunE == nil {
			t.Errorf("%s shortcut must have a command handler", shortcut.Name())
		}
		if shortcut.Flags().Lookup("dry-run") == nil {
			t.Errorf("%s shortcut is missing --dry-run", shortcut.Name())
		}
		if shortcut.Flags().Lookup("workers") == nil {
			t.Errorf("%s shortcut is missing --workers", shortcut.Name())
		}
	}
	if encryptCmd.RunE == nil || decryptCmd.RunE == nil {
		t.Fatal("encryption shortcuts must delegate to the bulk command handlers")
	}
}

func TestFormatEncryptionProgress_UsesColorWhenEnabled(t *testing.T) {
	originalForceColor := forceColor
	originalNoColor := noColor
	defer func() {
		forceColor = originalForceColor
		noColor = originalNoColor
	}()
	forceColor = true
	noColor = false

	colored := formatEncryptionProgress("ENCRYPTED", currentLogTheme.Success, "post.md", "default")
	if !strings.Contains(colored, "\033[") {
		t.Fatalf("expected ANSI color output, got %q", colored)
	}

	noColor = true
	plain := formatEncryptionProgress("ENCRYPTED", currentLogTheme.Success, "post.md", "default")
	if strings.Contains(plain, "\033[") {
		t.Fatalf("expected plain output with --no-color, got %q", plain)
	}
	if plain != "ENCRYPTED post.md key=default" {
		t.Fatalf("plain output = %q", plain)
	}
}

func TestGeneratePasswordCommand_LengthTooShort(t *testing.T) {
	originalLength := encryptionPasswordLength
	defer func() { encryptionPasswordLength = originalLength }()
	encryptionPasswordLength = encryption.DefaultMinPasswordLength - 1

	buf := bytes.NewBuffer(nil)
	generatePasswordCmd.SetOut(buf)

	if err := runGeneratePasswordCommand(generatePasswordCmd, nil); err == nil {
		t.Error("expected error when requested length < minimum")
	}
}

func TestCheckPasswordCommand_Pass(t *testing.T) {
	configPath := writeEncryptionConfigFile(t)
	originalCfg := cfgFile
	originalKey := encryptionCheckKey
	defer func() {
		cfgFile = originalCfg
		encryptionCheckKey = originalKey
	}()

	cfgFile = configPath
	encryptionCheckKey = ""
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")

	buf := bytes.NewBuffer(nil)
	checkPasswordCmd.SetOut(buf)

	if err := runCheckPasswordCommand(checkPasswordCmd, nil); err != nil {
		t.Fatalf("runCheckPasswordCommand() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "PASS default") {
		t.Fatalf("expected PASS output, got %q", output)
	}
}

func TestCheckPasswordCommand_Fail(t *testing.T) {
	configPath := writeEncryptionConfigFile(t)
	originalCfg := cfgFile
	originalKey := encryptionCheckKey
	defer func() {
		cfgFile = originalCfg
		encryptionCheckKey = originalKey
	}()

	cfgFile = configPath
	encryptionCheckKey = ""
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "weak")

	buf := bytes.NewBuffer(nil)
	checkPasswordCmd.SetOut(buf)

	err := runCheckPasswordCommand(checkPasswordCmd, nil)
	if err == nil {
		t.Fatal("expected runCheckPasswordCommand() to fail")
	}

	output := buf.String()
	if !strings.Contains(output, "FAIL default") {
		t.Fatalf("expected FAIL output, got %q", output)
	}
	if !strings.Contains(output, "estimated=") {
		t.Fatalf("expected estimated crack time in fail output, got %q", output)
	}
}

func TestEncryptPostSourceFile_PrivatePost(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	path := writeMarkdownFile(t, `---
title: Secret
private: true
---
# Secret
body
`)

	result, err := encryptPostSourceFile(path, cfg, false)
	if err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}
	if result.Action != encryptPostActionEncrypted {
		t.Fatalf("action = %q, want %q", result.Action, encryptPostActionEncrypted)
	}

	content := readFileString(t, path)
	_, body, err := plugins.ExtractFrontmatter(content)
	if err != nil {
		t.Fatalf("ExtractFrontmatter() error = %v", err)
	}
	if !encryption.IsSourceEncrypted(body) {
		t.Fatalf("expected source-encrypted body, got %q", body)
	}
	decrypted, keyName, err := encryption.DecryptSourceMarkdown(body, "h7Qm!2Vx9#Lp4@Td")
	if err != nil {
		t.Fatalf("DecryptSourceMarkdown() error = %v", err)
	}
	if keyName != "default" {
		t.Fatalf("keyName = %q, want default", keyName)
	}
	if decrypted != "# Secret\nbody\n" {
		t.Fatalf("decrypted body = %q", decrypted)
	}
}

func TestEncryptPostSourceFile_DryRunDoesNotWrite(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	original := `---
title: Secret
private: true
---
secret body
`
	path := writeMarkdownFile(t, original)

	result, err := encryptPostSourceFile(path, cfg, true)
	if err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}
	if result.Action != encryptPostActionEncrypted {
		t.Fatalf("action = %q, want %q", result.Action, encryptPostActionEncrypted)
	}
	if got := readFileString(t, path); got != original {
		t.Fatalf("dry run modified file: got %q", got)
	}
}

func TestEncryptPostSourceFile_PrivateTagUsesTagKey(t *testing.T) {
	cfg := testEncryptPostsConfig()
	cfg.Encryption.PrivateTags = map[string]string{"diary": "personal"}
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_PERSONAL", "h7Qm!2Vx9#Lp4@Td")
	path := writeMarkdownFile(t, `---
title: Diary
tags:
  - diary
---
tagged secret
`)

	result, err := encryptPostSourceFile(path, cfg, false)
	if err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}
	if result.KeyName != "personal" {
		t.Fatalf("key = %q, want personal", result.KeyName)
	}

	_, body, err := plugins.ExtractFrontmatter(readFileString(t, path))
	if err != nil {
		t.Fatalf("ExtractFrontmatter() error = %v", err)
	}
	_, keyName, err := encryption.DecryptSourceMarkdown(body, "h7Qm!2Vx9#Lp4@Td")
	if err != nil {
		t.Fatalf("DecryptSourceMarkdown() error = %v", err)
	}
	if keyName != "personal" {
		t.Fatalf("encrypted marker key = %q, want personal", keyName)
	}
}

func TestEncryptPostSourceFile_ExplicitPrivateFalseOverridesTag(t *testing.T) {
	cfg := testEncryptPostsConfig()
	cfg.Encryption.PrivateTags = map[string]string{"diary": "personal"}
	path := writeMarkdownFile(t, `---
title: Public diary
private: false
tags:
  - diary
---
public body
`)

	result, err := encryptPostSourceFile(path, cfg, false)
	if err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}
	if result.Action != encryptPostActionPublic {
		t.Fatalf("action = %q, want %q", result.Action, encryptPostActionPublic)
	}
	if encryption.IsSourceEncrypted(strings.TrimSpace(strings.SplitN(readFileString(t, path), "---", 3)[2])) {
		t.Fatal("explicit private: false post was encrypted")
	}
}

func TestEncryptPostSourceFile_AlreadyEncryptedSkipsWithoutKey(t *testing.T) {
	cfg := testEncryptPostsConfig()
	encryptedBody, err := encryption.EncryptSourceMarkdown("secret body\n", "default", "h7Qm!2Vx9#Lp4@Td")
	if err != nil {
		t.Fatalf("EncryptSourceMarkdown() error = %v", err)
	}
	path := writeMarkdownFile(t, "---\ntitle: Secret\nprivate: true\n---\n"+encryptedBody)

	result, err := encryptPostSourceFile(path, cfg, false)
	if err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}
	if result.Action != encryptPostActionAlreadyEncrypted {
		t.Fatalf("action = %q, want %q", result.Action, encryptPostActionAlreadyEncrypted)
	}
}

func TestEncryptPostSourceFile_WeakExplicitKeyFails(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_PERSONAL", "weak")
	path := writeMarkdownFile(t, `---
title: Secret
private: true
secret_key: personal
---
secret body
`)

	_, err := encryptPostSourceFile(path, cfg, true)
	if err == nil {
		t.Fatal("expected weak key error")
	}
	if !strings.Contains(err.Error(), "failed policy") {
		t.Fatalf("expected policy error, got %v", err)
	}
}

func writeEncryptionConfigFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "markata-go.toml")
	content := `[markata-go]
title = "test"

[markata-go.encryption]
enabled = true
default_key = "default"
enforce_strength = true
min_password_length = 14
min_estimated_crack_time = "10y"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func testEncryptPostsConfig() *models.Config {
	return &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns:     []string{"**/*.md"},
			UseGitignore: true,
		},
		Encryption: models.EncryptionConfig{
			Enabled:               true,
			DefaultKey:            "default",
			EnforceStrength:       true,
			MinPasswordLength:     14,
			MinEstimatedCrackTime: "10y",
		},
	}
}

func writeMarkdownFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "post.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}
	return path
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func TestFormatCrackDurationHuman(t *testing.T) {
	tests := []struct {
		name string
		in   time.Duration
		want string
	}{
		{name: "years", in: 10 * 365 * 24 * time.Hour, want: "10.0y"},
		{name: "days", in: 48 * time.Hour, want: "2.0d"},
		{name: "hours", in: 90 * time.Minute, want: "1.5h"},
		{name: "subsecond", in: 250 * time.Millisecond, want: "<1s"},
		{name: "zero", in: 0, want: "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatCrackDurationHuman(tt.in)
			if got != tt.want {
				t.Fatalf("formatCrackDurationHuman(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDecryptPostSourceFile_RoundTrip(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	original := `---
title: Secret
private: true
---
# Secret
body
`
	path := writeMarkdownFile(t, original)

	if _, err := encryptPostSourceFile(path, cfg, false); err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}

	result, err := decryptPostSourceFile(path, cfg, false)
	if err != nil {
		t.Fatalf("decryptPostSourceFile() error = %v", err)
	}
	if result.Action != decryptPostActionDecrypted {
		t.Fatalf("action = %q, want %q", result.Action, decryptPostActionDecrypted)
	}
	if result.KeyName != "default" {
		t.Fatalf("keyName = %q, want default", result.KeyName)
	}
	if got := readFileString(t, path); got != original {
		t.Fatalf("round trip mismatch:\ngot  %q\nwant %q", got, original)
	}
}

func TestSourceEncryptionCache_ReusesCiphertextAfterDecrypt(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	path := writeMarkdownFile(t, `---
title: Secret
private: true
---
secret body
`)
	cache, err := buildcache.LoadSourceEncryptionCache(t.TempDir())
	if err != nil {
		t.Fatalf("load source encryption cache: %v", err)
	}

	first, err := prepareEncryptPostSourceFileWithCache(path, cfg, cache)
	if err != nil {
		t.Fatalf("initial preparation: %v", err)
	}
	if err := writeSourceDocuments([]sourceDocument{first.document}); err != nil {
		t.Fatalf("write encrypted source: %v", err)
	}
	cacheSourceEncryptionDocuments(cache, []sourceDocument{first.document})
	originalEncrypted := readFileString(t, path)

	decrypted, err := prepareDecryptPostSourceFile(path, cfg)
	if err != nil {
		t.Fatalf("decrypt preparation: %v", err)
	}
	if err := writeSourceDocuments([]sourceDocument{decrypted.document}); err != nil {
		t.Fatalf("write decrypted source: %v", err)
	}
	cacheSourceEncryptionDocuments(cache, []sourceDocument{decrypted.document})

	reencrypted, err := prepareEncryptPostSourceFileWithCache(path, cfg, cache)
	if err != nil {
		t.Fatalf("re-encrypt preparation: %v", err)
	}
	if err := writeSourceDocuments([]sourceDocument{reencrypted.document}); err != nil {
		t.Fatalf("write re-encrypted source: %v", err)
	}
	if got := readFileString(t, path); got != originalEncrypted {
		t.Fatalf("re-encryption changed unchanged ciphertext")
	}
}

func TestDecryptPostSourceFile_SkipsPlaintext(t *testing.T) {
	cfg := testEncryptPostsConfig()
	original := `---
title: Public
---
plain body
`
	path := writeMarkdownFile(t, original)

	result, err := decryptPostSourceFile(path, cfg, false)
	if err != nil {
		t.Fatalf("decryptPostSourceFile() error = %v", err)
	}
	if result.Action != decryptPostActionNotEncrypted {
		t.Fatalf("action = %q, want %q", result.Action, decryptPostActionNotEncrypted)
	}
	if got := readFileString(t, path); got != original {
		t.Fatalf("plaintext file was modified: got %q", got)
	}
}

func TestDecryptPostSourceFile_DryRunDoesNotWrite(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	path := writeMarkdownFile(t, `---
title: Secret
private: true
---
secret body
`)
	if _, err := encryptPostSourceFile(path, cfg, false); err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}
	encrypted := readFileString(t, path)

	result, err := decryptPostSourceFile(path, cfg, true)
	if err != nil {
		t.Fatalf("decryptPostSourceFile() error = %v", err)
	}
	if result.Action != decryptPostActionDecrypted {
		t.Fatalf("action = %q, want %q", result.Action, decryptPostActionDecrypted)
	}
	if got := readFileString(t, path); got != encrypted {
		t.Fatalf("dry run modified file: got %q", got)
	}
}

func TestDecryptPostSourceFile_MissingKeyEnvFails(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	path := writeMarkdownFile(t, `---
title: Secret
private: true
---
secret body
`)
	if _, err := encryptPostSourceFile(path, cfg, false); err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}

	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "")
	if _, err := decryptPostSourceFile(path, cfg, true); err == nil {
		t.Fatal("expected error when key env var is unset")
	}
}

func TestDecryptPostSourceFile_WrongPasswordFails(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	path := writeMarkdownFile(t, `---
title: Secret
private: true
---
secret body
`)
	if _, err := encryptPostSourceFile(path, cfg, false); err != nil {
		t.Fatalf("encryptPostSourceFile() error = %v", err)
	}

	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "wrong-password-value-1")
	if _, err := decryptPostSourceFile(path, cfg, true); err == nil {
		t.Fatal("expected error when password is wrong")
	}
}

func TestPrepareEncryptPostSourceFiles_FailureDoesNotModifyAnyFile(t *testing.T) {
	cfg := testEncryptPostsConfig()
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")
	dir := t.TempDir()
	validPath := filepath.Join(dir, "valid.md")
	invalidPath := filepath.Join(dir, "invalid.md")
	valid := "---\ntitle: Valid\nprivate: true\n---\nprivate body\n"
	invalid := "---\ntitle: Invalid\nprivate: true\nsecret_key: missing\n---\nprivate body\n"
	for path, content := range map[string]string{validPath: valid, invalidPath: invalid} {
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	_, err := prepareSourceFiles([]string{validPath, invalidPath}, 2, func(path string) (preparedEncryptPost, error) {
		return prepareEncryptPostSourceFile(path, cfg)
	})
	if err == nil {
		t.Fatal("expected preparation to fail for missing key")
	}
	if got := readFileString(t, validPath); got != valid {
		t.Fatalf("valid file changed after another preparation failed: got %q, want %q", got, valid)
	}
	if got := readFileString(t, invalidPath); got != invalid {
		t.Fatalf("invalid file changed after preparation failed: got %q, want %q", got, invalid)
	}
}

func TestPrepareSourceFiles_PreservesInputOrder(t *testing.T) {
	paths := []string{"first", "second", "third"}
	results, err := prepareSourceFiles(paths, 2, func(path string) (string, error) {
		return path + "-done", nil
	})
	if err != nil {
		t.Fatalf("prepareSourceFiles() error = %v", err)
	}
	want := []string{"first-done", "second-done", "third-done"}
	for i := range want {
		if results[i] != want[i] {
			t.Errorf("results[%d] = %q, want %q", i, results[i], want[i])
		}
	}
}

func TestPrepareSourceFiles_NegativeWorkersFails(t *testing.T) {
	_, err := prepareSourceFiles([]string{"post.md"}, -1, func(path string) (string, error) {
		return path, nil
	})
	if err == nil {
		t.Fatal("expected negative workers to fail")
	}
}

func TestWriteSourceDocuments_ChangedSourceLeavesBatchUnchanged(t *testing.T) {
	dir := t.TempDir()
	firstPath := filepath.Join(dir, "first.md")
	secondPath := filepath.Join(dir, "second.md")
	for path, content := range map[string]string{firstPath: "first original", secondPath: "second original"} {
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	// Simulate an edit after preparation, before the batch write begins.
	if err := os.WriteFile(secondPath, []byte("second edited"), 0o600); err != nil {
		t.Fatalf("edit %s: %v", secondPath, err)
	}
	err := writeSourceDocuments([]sourceDocument{
		{path: firstPath, original: "first original", content: "first transformed", mode: 0o600},
		{path: secondPath, original: "second original", content: "second transformed", mode: 0o600},
	})
	if err == nil {
		t.Fatal("expected changed source to fail")
	}
	if got := readFileString(t, firstPath); got != "first original" {
		t.Fatalf("first file changed after stale source detection: got %q", got)
	}
	if got := readFileString(t, secondPath); got != "second edited" {
		t.Fatalf("edited file was overwritten: got %q", got)
	}
}

func TestDecryptionTargetFiles_ExplicitPathsAndDirs(t *testing.T) {
	cfg := testEncryptPostsConfig()
	dir := t.TempDir()
	nested := filepath.Join(dir, "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	mdFile := filepath.Join(nested, "a.md")
	txtFile := filepath.Join(nested, "b.txt")
	for _, p := range []string{mdFile, txtFile} {
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatalf("write %s: %v", p, err)
		}
	}

	files, err := decryptionTargetFiles(cfg, []string{dir})
	if err != nil {
		t.Fatalf("decryptionTargetFiles() error = %v", err)
	}
	if len(files) != 1 || files[0] != mdFile {
		t.Fatalf("files = %v, want [%s]", files, mdFile)
	}
}
