package builderadmin

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/WaylonWalker/markata/internal/fileidentity"
	"github.com/WaylonWalker/markata/pkg/assets"
	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
	"github.com/WaylonWalker/markata/pkg/sourcegit"
)

const (
	noOpBaselineVersion = 2
	noOpKeySize         = 32
	noOpMaxAge          = 30 * time.Minute
	noOpBuildLogLimit   = 64 << 20
)

var noOpNowToken = regexp.MustCompile(`(?i)\bnow\b`)
var noOpDateFormat = regexp.MustCompile(`date\s*(?::|\()\s*["']([^"']+)["']`)
var noOpJinjaFlag = regexp.MustCompile(`(?im)^\s*jinja(?:_md)?\s*:\s*(?:true|yes|on)\s*(?:#.*)?$`)
var noOpMentionCandidate = regexp.MustCompile(`(?i)(?:^|[^@~\w])([@~])([a-zA-Z][a-zA-Z0-9_.-]*)(?:[^a-zA-Z0-9_.-]|$)`)

type noOpBaseline struct {
	Version            int       `json:"version"`
	Fingerprint        string    `json:"fingerprint"`
	ReleaseFingerprint string    `json:"release_fingerprint"`
	ReleaseID          string    `json:"release_id"`
	CreatedAt          time.Time `json:"created_at"`
	ExpiresAt          time.Time `json:"expires_at"`
}

type noOpSnapshot struct {
	Fingerprint            string
	SourceCommit           string
	SourceFingerprint      string
	BuilderFingerprint     string
	EnvironmentFingerprint string
	EncoderFingerprint     string
	InputFingerprint       string
	CalendarDay            string
	ConfigPath             string
	InputRoots             []string
	siteConfig             *models.Config
	SourceMS               int64
	InputsMS               int64
	InputEntries           int64
	ReleaseMS              int64
	ReleaseEntries         int64
	ExternalExpiresAt      time.Time
}

func (s *Service) canUseNoOp(req queueRequest) bool {
	return s.noOpEligibilityReason(req) == ""
}

func (s *Service) noOpEligibilityReason(req queueRequest) string {
	if req.Force || req.SyncSource || len(req.Changed) > 0 {
		return "request_requires_full_build"
	}
	switch req.TriggerType {
	case "manual-ui", "agent-api":
	case "scheduled-refresh":
		if req.Detail != "Build enqueued by refresh task reader-update" {
			return "refresh_task_not_eligible"
		}
	default:
		return "trigger_not_eligible"
	}
	s.stateMu.Lock()
	hasQueued := len(s.state.Queue) > 0
	s.stateMu.Unlock()
	if hasQueued {
		return "other_work_queued"
	}
	return ""
}

func (s *Service) verifiedNoOp(ctx context.Context, req queueRequest, record *BuildRecord, log io.Writer) (bool, int64) {
	attemptStart := time.Now()
	skip := func(reason string) (bool, int64) {
		if req.TriggerType == "manual-ui" || req.TriggerType == "agent-api" || req.TriggerType == "scheduled-refresh" {
			_, _ = fmt.Fprintf(log, "no_op=skipped reason=%s attempt_ms=%d\n", reason, time.Since(attemptStart).Milliseconds())
		}
		return false, 0
	}
	if reason := s.noOpEligibilityReason(req); reason != "" {
		return skip(reason)
	}
	if !fileidentity.Available {
		return skip("file_identity_unavailable")
	}
	if s.cfg.ConfigPath == "" {
		return skip("explicit_config_required")
	}
	s.releaseMu.Lock()
	defer s.releaseMu.Unlock()
	verification, reason := s.verifyNoOpState(ctx, record)
	if reason != "" {
		return skip(reason)
	}
	checkMS := time.Since(record.StartedAt).Milliseconds()
	_, _ = fmt.Fprintf(log, "result=unchanged\nrelease=%s\nsource_commit=%s\ncheck_ms=%d\nverification_passes=2\nsource_scan_ms=%d\ninput_scan_ms=%d\ninput_entries=%d\nrelease_scan_ms=%d\nrelease_entries=%d\n", verification.releaseID, verification.snapshot.SourceCommit, checkMS, verification.snapshot.SourceMS, verification.snapshot.InputsMS, verification.snapshot.InputEntries, verification.snapshot.ReleaseMS, verification.snapshot.ReleaseEntries)
	record.Status = "success"
	record.Result = "unchanged"
	record.SourceCommit = reliableBuilderRevision(verification.snapshot.SourceCommit)
	record.ReleaseID = verification.releaseID
	record.ReleasePath = verification.releasePath
	record.NoOpMS = checkMS
	record.NoOpSourceMS = verification.snapshot.SourceMS
	record.NoOpInputsMS = verification.snapshot.InputsMS
	record.NoOpInputEntries = verification.snapshot.InputEntries
	record.NoOpReleaseMS = verification.snapshot.ReleaseMS
	record.NoOpReleaseEntries = verification.snapshot.ReleaseEntries
	record.TotalMS = checkMS
	record.PerfSummary = []string{"Result: unchanged", fmt.Sprintf("No-op check: %dms", checkMS)}
	record.FinishedAt = time.Now().UTC()
	record.BecameLive = false
	s.finishBuild(*record)
	return true, checkMS
}

type noOpVerification struct {
	releaseID   string
	releasePath string
	snapshot    noOpSnapshot
}

type noOpReleaseIdentity struct {
	name string
	info fs.FileInfo
}

func (s *Service) verifyNoOpState(ctx context.Context, record *BuildRecord) (noOpVerification, string) {
	baseline, current, reason := s.currentNoOpBaseline(record)
	if reason != "" {
		return noOpVerification{}, reason
	}
	snapshot, _, releasePath, reason := s.verifyNoOpSnapshot(ctx, current, baseline)
	if reason != "" {
		return noOpVerification{}, reason
	}
	confirm, err := s.confirmNoOpSnapshot(ctx, snapshot)
	if err != nil {
		return noOpVerification{}, "input_recheck_failed"
	}
	if confirm.SourceFingerprint != snapshot.SourceFingerprint || confirm.InputFingerprint != snapshot.InputFingerprint || confirm.BuilderFingerprint != snapshot.BuilderFingerprint || confirm.EnvironmentFingerprint != snapshot.EnvironmentFingerprint || confirm.EncoderFingerprint != snapshot.EncoderFingerprint || confirm.CalendarDay != snapshot.CalendarDay {
		return noOpVerification{}, "inputs_changed_during_verification"
	}
	confirmReleasePath, err := s.noOpCurrentReleasePath(current)
	if err != nil || confirmReleasePath != releasePath {
		return noOpVerification{}, "live_release_changed_during_verification"
	}
	if reason := noOpCompletionReason(ctx, baseline.ExpiresAt, time.Now()); reason != "" {
		return noOpVerification{}, reason
	}
	snapshot.SourceMS += confirm.SourceMS
	snapshot.InputsMS += confirm.InputsMS
	snapshot.InputEntries += confirm.InputEntries
	return noOpVerification{releaseID: current, releasePath: releasePath, snapshot: snapshot}, ""
}

// A baseline may expire or the build may be canceled during verification.
func noOpCompletionReason(ctx context.Context, expiresAt, now time.Time) string {
	if ctx.Err() != nil {
		return "verification_canceled"
	}
	if !now.Before(expiresAt) {
		return "baseline_expired_during_verification"
	}
	return ""
}

func (s *Service) currentNoOpBaseline(record *BuildRecord) (noOpBaseline, string, string) {
	if record.ReleaseID != "" {
		return noOpBaseline{}, "", "release_already_selected"
	}
	current := s.currentReleaseID()
	if current == "" {
		return noOpBaseline{}, "", "no_current_release"
	}
	reusable, err := reusableWorkspaceMatches(s.cfg.WorkDir, current)
	if err != nil || !reusable {
		return noOpBaseline{}, "", "workspace_not_reusable"
	}
	baseline, err := s.loadNoOpBaseline()
	if err != nil {
		return noOpBaseline{}, "", "baseline_missing_or_invalid"
	}
	if baseline.Version != noOpBaselineVersion || baseline.ReleaseID != current || baseline.Fingerprint == "" {
		return noOpBaseline{}, "", "baseline_does_not_match_release"
	}
	if baseline.ReleaseFingerprint == "" {
		return noOpBaseline{}, "", "release_baseline_missing"
	}
	if !time.Now().Before(baseline.ExpiresAt) {
		return noOpBaseline{}, "", "baseline_expired"
	}
	return baseline, current, ""
}

func (s *Service) verifyNoOpSnapshot(ctx context.Context, current string, baseline noOpBaseline) (noOpSnapshot, string, string, string) {
	snapshot, err := s.noOpInputSnapshot(ctx)
	if err != nil {
		return noOpSnapshot{}, "", "", "input_snapshot_failed"
	}
	if snapshot.Fingerprint != baseline.Fingerprint {
		return noOpSnapshot{}, "", "", "inputs_changed"
	}
	releaseFingerprint, releasePath, releaseMS, releaseEntries, err := s.noOpReleaseFingerprintWithStats(ctx, current)
	if err != nil || releaseFingerprint != baseline.ReleaseFingerprint {
		return noOpSnapshot{}, "", "", "live_release_changed_or_unsafe"
	}
	snapshot.ReleaseMS = releaseMS
	snapshot.ReleaseEntries = releaseEntries
	return snapshot, releaseFingerprint, releasePath, ""
}

func (s *Service) captureNoOpSnapshot(ctx context.Context) (noOpSnapshot, error) {
	if !fileidentity.Available {
		return noOpSnapshot{}, errors.New("verified file identity is unavailable")
	}
	siteConfig, configPath, err := s.loadNoOpSiteConfig()
	if err != nil {
		return noOpSnapshot{}, err
	}
	if err := validateNoOpConfig(siteConfig, s.cfg.MermaidMode); err != nil {
		return noOpSnapshot{}, err
	}
	inputRoots := noOpInputRoots(s.cfg.SourceDir, configPath, s.cfg.CacheMount, siteConfig)
	cacheRoots, err := noOpManagedCacheInputRoots(s.cfg.SourceDir, s.cfg.CacheMount)
	if err != nil {
		return noOpSnapshot{}, err
	}
	inputRoots = append(inputRoots, cacheRoots...)
	themeRoots, err := noOpExternalThemeRoots(s.cfg.SourceDir, s.executable, siteConfig.Theme.Name)
	if err != nil {
		return noOpSnapshot{}, err
	}
	inputRoots = append(inputRoots, themeRoots...)
	source, sourceMS, inputRoots, err := s.noOpSourceSnapshot(ctx, siteConfig, inputRoots)
	if err != nil {
		return noOpSnapshot{}, err
	}
	if source.Dirty == nil || source.Fingerprint() == "" {
		return noOpSnapshot{}, errors.New("source snapshot is incomplete")
	}
	sourceFingerprint := source.Fingerprint()

	key, err := s.loadOrCreateNoOpKey()
	if err != nil {
		return noOpSnapshot{}, err
	}
	calendarDay := time.Now().Format("2006-01-02")
	h := hmac.New(sha256.New, key)
	writeNoOpPart(h, "version", fmt.Sprint(noOpBaselineVersion))
	writeNoOpPart(h, "source-commit", source.Commit)
	writeNoOpPart(h, "source-snapshot", sourceFingerprint)
	writeNoOpPart(h, "builder-version", s.cfg.BuilderVersion)
	writeNoOpPart(h, "builder-commit", reliableBuilderRevision(s.cfg.BuilderCommit))
	builderFingerprint, err := s.noOpBuilderFingerprint()
	if err != nil {
		return noOpSnapshot{}, err
	}
	writeNoOpPart(h, "builder-executable", builderFingerprint)
	writeNoOpPart(h, "config-path", configPath)
	writeNoOpPart(h, "fast", fmt.Sprint(s.cfg.Fast))
	writeNoOpPart(h, "skip-page-weight", fmt.Sprint(s.cfg.SkipPageWeight))
	writeNoOpPart(h, "mermaid-mode", s.cfg.MermaidMode)
	writeNoOpPart(h, "cache-mount", s.cfg.CacheMount)
	writeNoOpPart(h, "calendar-day", calendarDay)
	configJSON, err := json.Marshal(siteConfig)
	if err != nil {
		return noOpSnapshot{}, fmt.Errorf("encode effective site config: %w", err)
	}
	writeNoOpPart(h, "effective-config", string(configJSON))
	extraJSON, err := json.Marshal(siteConfig.Extra)
	if err != nil {
		return noOpSnapshot{}, fmt.Errorf("encode plugin configuration: %w", err)
	}
	writeNoOpPart(h, "plugin-config", string(extraJSON))
	encoderHash := sha256.New()
	if err := hashNoOpImageEncoders(encoderHash, s.cfg.SourceDir, siteConfig.Extra["image_optimization"]); err != nil {
		return noOpSnapshot{}, err
	}
	encoderFingerprint := hex.EncodeToString(encoderHash.Sum(nil))
	writeNoOpPart(h, "image-encoders", encoderFingerprint)
	environmentFingerprint := noOpEnvironmentFingerprint()
	writeNoOpPart(h, "environment", environmentFingerprint)
	sourceRoot, err := filepath.Abs(s.cfg.SourceDir)
	if err != nil {
		return noOpSnapshot{}, fmt.Errorf("resolve source root: %w", err)
	}
	inputHash := sha256.New()
	inputMS, inputEntries, err := hashNoOpInputRoots(ctx, inputHash, sourceRoot, inputRoots)
	if err != nil {
		return noOpSnapshot{}, err
	}
	inputFingerprint := hex.EncodeToString(inputHash.Sum(nil))
	writeNoOpPart(h, "input-roots", inputFingerprint)
	return noOpSnapshot{
		Fingerprint:            hex.EncodeToString(h.Sum(nil)),
		SourceCommit:           source.Commit,
		SourceFingerprint:      sourceFingerprint,
		BuilderFingerprint:     builderFingerprint,
		EnvironmentFingerprint: environmentFingerprint,
		EncoderFingerprint:     encoderFingerprint,
		InputFingerprint:       inputFingerprint,
		CalendarDay:            calendarDay,
		ConfigPath:             configPath,
		InputRoots:             append([]string(nil), inputRoots...),
		siteConfig:             siteConfig,
		SourceMS:               sourceMS,
		InputsMS:               inputMS,
		InputEntries:           inputEntries,
	}, nil
}

func (s *Service) confirmNoOpSnapshot(ctx context.Context, expected noOpSnapshot) (noOpSnapshot, error) {
	if expected.siteConfig == nil || expected.ConfigPath == "" || len(expected.InputRoots) == 0 {
		return noOpSnapshot{}, errors.New("no-op verification snapshot is incomplete")
	}
	if time.Now().Format("2006-01-02") != expected.CalendarDay {
		return noOpSnapshot{}, errors.New("no-op snapshot crossed a calendar day")
	}
	inputRoots := noOpInputRoots(s.cfg.SourceDir, expected.ConfigPath, s.cfg.CacheMount, expected.siteConfig)
	themeRoots, err := noOpExternalThemeRoots(s.cfg.SourceDir, s.executable, expected.siteConfig.Theme.Name)
	if err != nil {
		return noOpSnapshot{}, err
	}
	inputRoots = append(inputRoots, themeRoots...)
	source, sourceMS, roots, err := s.noOpSourceSnapshot(ctx, expected.siteConfig, inputRoots)
	if err != nil {
		return noOpSnapshot{}, err
	}
	if source.Dirty == nil || source.Commit != expected.SourceCommit || source.Fingerprint() != expected.SourceFingerprint || !equalNoOpRoots(roots, expected.InputRoots) {
		return noOpSnapshot{}, errors.New("source inputs changed during no-op verification")
	}
	builderFingerprint, err := s.noOpBuilderFingerprint()
	if err != nil {
		return noOpSnapshot{}, err
	}
	if builderFingerprint != expected.BuilderFingerprint {
		return noOpSnapshot{}, errors.New("builder executable changed during no-op verification")
	}
	encoderHash := sha256.New()
	if err := hashNoOpImageEncoders(encoderHash, s.cfg.SourceDir, expected.siteConfig.Extra["image_optimization"]); err != nil {
		return noOpSnapshot{}, err
	}
	encoderFingerprint := hex.EncodeToString(encoderHash.Sum(nil))
	environmentFingerprint := noOpEnvironmentFingerprint()
	sourceRoot, err := filepath.Abs(s.cfg.SourceDir)
	if err != nil {
		return noOpSnapshot{}, err
	}
	inputHash := sha256.New()
	inputMS, inputEntries, err := hashNoOpInputRoots(ctx, inputHash, sourceRoot, expected.InputRoots)
	if err != nil {
		return noOpSnapshot{}, err
	}
	return noOpSnapshot{
		SourceCommit:           source.Commit,
		SourceFingerprint:      source.Fingerprint(),
		BuilderFingerprint:     builderFingerprint,
		EnvironmentFingerprint: environmentFingerprint,
		EncoderFingerprint:     encoderFingerprint,
		InputFingerprint:       hex.EncodeToString(inputHash.Sum(nil)),
		CalendarDay:            time.Now().Format("2006-01-02"),
		SourceMS:               sourceMS,
		InputsMS:               inputMS,
		InputEntries:           inputEntries,
	}, nil
}

func noOpEnvironmentFingerprint() string {
	environment := append([]string(nil), os.Environ()...)
	sort.Strings(environment)
	h := sha256.New()
	for _, item := range environment {
		writeNoOpPart(h, "environment", item)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func equalNoOpRoots(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (s *Service) noOpSourceSnapshot(ctx context.Context, siteConfig *models.Config, inputRoots []string) (sourcegit.State, int64, []string, error) {
	siteGitignore, err := plugins.NewSiteGitignoreMatcher(s.cfg.SourceDir)
	if err != nil {
		return sourcegit.State{}, 0, nil, fmt.Errorf("load site .gitignore: %w", err)
	}
	start := time.Now()
	source, err := sourcegit.ReadMetadataWithOptions(ctx, s.cfg.SourceDir, sourcegit.ReadOptions{
		IncludeIgnoredContent: false,
		IgnoredPathFilter: func(path string) (bool, error) {
			directory := strings.HasSuffix(path, "/")
			path = strings.TrimSuffix(path, "/")
			absolutePath := filepath.Join(s.cfg.SourceDir, filepath.FromSlash(path))
			if !directory && isNoOpContentFile(siteConfig, path) {
				return siteGitignore(filepath.ToSlash(path)), nil
			}
			if directory && !siteGitignore(filepath.ToSlash(path)) {
				return false, nil
			}
			inputRoot := absolutePath
			if noOpManagedCacheRootPath(s.cfg.SourceDir, absolutePath) {
				inputRoot, err = resolveNoOpIgnoredCacheRoot(s.cfg.SourceDir, s.cfg.CacheMount, absolutePath)
				if err != nil {
					return false, fmt.Errorf("resolve ignored source input %q: %w", path, err)
				}
			}
			if !noOpPathCoveredByRoots(inputRoot, inputRoots) {
				inputRoots = append(inputRoots, inputRoot)
			}
			return true, nil
		},
	})
	if err != nil {
		return sourcegit.State{}, 0, nil, fmt.Errorf("read source snapshot: %w", err)
	}
	return source, time.Since(start).Milliseconds(), compactNoOpInputRoots(inputRoots), nil
}

func noOpManagedCacheRootPath(sourceDir, path string) bool {
	relative, err := filepath.Rel(filepath.Clean(sourceDir), filepath.Clean(path))
	if err != nil {
		return false
	}
	switch filepath.ToSlash(relative) {
	case ".markata", ".markata-cache":
		return true
	default:
		return false
	}
}

// noOpManagedCacheInputRoots verifies Builder Admin's cache-mount links before
// fingerprinting them. This also detects a missing link, which Git-ignore
// discovery cannot report as an input path.
func noOpManagedCacheInputRoots(sourceDir, cacheMount string) ([]string, error) {
	if strings.TrimSpace(cacheMount) == "" {
		return nil, nil
	}
	roots := make([]string, 0, 2)
	for _, name := range []string{".markata", ".markata-cache"} {
		root, err := resolveNoOpIgnoredCacheRoot(sourceDir, cacheMount, filepath.Join(sourceDir, name))
		if err != nil {
			return nil, err
		}
		roots = append(roots, root)
	}
	return roots, nil
}

// resolveNoOpIgnoredCacheRoot allows only the cache symlinks Builder Admin
// creates for a configured cache mount. Gitignore discovery sees the symlink
// itself as an ignored input directory; fingerprint its configured target
// instead of traversing an unverified alias.
func resolveNoOpIgnoredCacheRoot(sourceDir, cacheMount, path string) (string, error) {
	absoluteSource, err := filepath.Abs(sourceDir)
	if err != nil {
		return "", fmt.Errorf("resolve source root: %w", err)
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve cache symlink path: %w", err)
	}
	info, err := os.Lstat(absolutePath)
	if errors.Is(err, fs.ErrNotExist) {
		if strings.TrimSpace(cacheMount) == "" {
			return absolutePath, nil
		}
		return "", fmt.Errorf("managed cache symlink %q is missing", path)
	}
	if err != nil {
		return "", fmt.Errorf("inspect cache input: %w", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		if strings.TrimSpace(cacheMount) == "" {
			return absolutePath, nil
		}
		return "", fmt.Errorf("managed cache root %q is not a symlink", path)
	}

	relative, err := filepath.Rel(absoluteSource, absolutePath)
	if err != nil {
		return "", fmt.Errorf("locate cache symlink: %w", err)
	}
	var targetName string
	switch filepath.ToSlash(relative) {
	case ".markata":
		targetName = "build"
	case ".markata-cache":
		targetName = "plugin"
	default:
		return "", fmt.Errorf("unsupported symlink root %q", path)
	}
	if strings.TrimSpace(cacheMount) == "" {
		return "", fmt.Errorf("managed cache symlink %q has no configured cache mount", path)
	}
	mountPath, err := resolveNoOpCacheMount(absoluteSource, cacheMount)
	if err != nil {
		return "", err
	}
	targetPath := filepath.Join(mountPath, targetName)
	if err := requireNoOpRealDirectory(targetPath, "configured cache target"); err != nil {
		return "", err
	}
	expectedTarget, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		return "", fmt.Errorf("resolve expected cache target: %w", err)
	}
	actualTarget, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", fmt.Errorf("resolve cache symlink target: %w", err)
	}
	if filepath.Clean(actualTarget) != filepath.Clean(expectedTarget) {
		return "", fmt.Errorf("cache symlink %q targets %q, want configured cache target %q", path, actualTarget, expectedTarget)
	}
	actualInfo, err := os.Stat(actualTarget)
	if err != nil {
		return "", fmt.Errorf("inspect cache symlink target: %w", err)
	}
	if !actualInfo.IsDir() {
		return "", fmt.Errorf("cache symlink %q target is not a directory", path)
	}
	return actualTarget, nil
}

func resolveNoOpCacheMount(sourceDir, cacheMount string) (string, error) {
	if strings.TrimSpace(cacheMount) == "" {
		return "", fmt.Errorf("managed cache symlink has no configured cache mount")
	}
	mountPath := cacheMount
	if !filepath.IsAbs(mountPath) {
		mountPath = filepath.Join(sourceDir, mountPath)
	}
	mountPath, err := filepath.Abs(mountPath)
	if err != nil {
		return "", fmt.Errorf("resolve cache mount: %w", err)
	}
	if err := requireNoOpRealDirectory(mountPath, "cache mount"); err != nil {
		return "", err
	}
	resolvedMount, err := filepath.EvalSymlinks(mountPath)
	if err != nil {
		return "", fmt.Errorf("resolve cache mount path: %w", err)
	}
	if filepath.Clean(resolvedMount) != filepath.Clean(mountPath) {
		return "", fmt.Errorf("cache mount %q contains a symlink", mountPath)
	}
	return mountPath, nil
}

func requireNoOpRealDirectory(path, label string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect %s: %w", label, err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%s %q is not a regular directory", label, path)
	}
	return nil
}

func hashNoOpInputRoots(ctx context.Context, h hash.Hash, sourceRoot string, roots []string) (int64, int64, error) {
	start := time.Now()
	var entries int64
	for _, dir := range roots {
		relative, err := filepath.Rel(dir, sourceRoot)
		if err != nil || relative == "." || relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
			return 0, 0, fmt.Errorf("no-op cache %q contains the source root", dir)
		}
		if err := hashNoOpTree(ctx, h, dir, &entries); err != nil {
			return 0, 0, err
		}
	}
	return time.Since(start).Milliseconds(), entries, nil
}

func (s *Service) loadNoOpSiteConfig() (*models.Config, string, error) {
	dotEnvPath := filepath.Join(s.cfg.SourceDir, ".env")
	if info, err := os.Lstat(dotEnvPath); err == nil {
		if !info.Mode().IsRegular() {
			return nil, "", errors.New("source .env is not a regular file")
		}
		if err := config.LoadDotEnvFile(dotEnvPath); err != nil {
			return nil, "", fmt.Errorf("load source .env: %w", err)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, "", fmt.Errorf("inspect source .env: %w", err)
	}
	if s.cfg.ConfigPath == "" {
		return nil, "", errors.New("explicit site config is required")
	}
	configPath := s.cfg.ConfigPath
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(s.cfg.SourceDir, configPath)
	}
	siteConfig, err := config.LoadWithMergeOptions(config.LoadOptions{DisableDotEnv: true}, configPath)
	if err != nil {
		return nil, "", fmt.Errorf("load effective site config: %w", err)
	}
	return siteConfig, configPath, nil
}

func (s *Service) validateNoOpSemantics(ctx context.Context) (time.Time, error) {
	if !fileidentity.Available {
		return time.Time{}, errors.New("verified file identity is unavailable")
	}
	siteConfig, _, err := s.loadNoOpSiteConfig()
	if err != nil {
		return time.Time{}, err
	}
	if err := validateNoOpConfig(siteConfig, s.cfg.MermaidMode); err != nil {
		return time.Time{}, err
	}
	if err := checkTemplateTimeUse(s.cfg.SourceDir, siteConfig.TemplatesDir); err != nil {
		return time.Time{}, err
	}
	if err := checkJinjaTimeUse(ctx, s.cfg.SourceDir, siteConfig); err != nil {
		return time.Time{}, err
	}
	if err := checkNoOpEmbedCache(ctx, s.cfg.SourceDir, s.cfg.CacheMount, siteConfig); err != nil {
		return time.Time{}, err
	}
	if noOpHookEnabled(siteConfig, "cdn_assets") {
		if err := checkNoOpCDNAssetCache(s.cfg.SourceDir, siteConfig); err != nil {
			return time.Time{}, err
		}
	}
	if noOpHookEnabled(siteConfig, "link_avatars") {
		if err := checkNoOpLinkAvatarCache(ctx, s.cfg.SourceDir, s.cfg.CacheMount, siteConfig); err != nil {
			return time.Time{}, err
		}
	}
	return checkNoOpMentionFreshness(ctx, s.cfg.SourceDir, siteConfig, time.Now())
}

func validateNoOpConfig(siteConfig *models.Config, builderMermaidMode ...string) error {
	if !siteConfig.GlobConfig.UseGitignore {
		return errors.New("no-op snapshots require use_gitignore=true")
	}
	if err := validateNoOpContentPatterns(siteConfig.GlobConfig.Patterns); err != nil {
		return err
	}
	buildCacheHook := false
	for _, hook := range siteConfig.Hooks {
		if strings.EqualFold(strings.TrimSpace(hook), "default") || strings.EqualFold(strings.TrimSpace(hook), "build_cache") {
			buildCacheHook = true
		}
	}
	for _, hook := range siteConfig.DisabledHooks {
		if strings.EqualFold(strings.TrimSpace(hook), "build_cache") {
			buildCacheHook = false
		}
	}
	if !buildCacheHook {
		return errors.New("no-op snapshots require the build_cache hook")
	}
	if extraBoolIsFalseMap(siteConfig.Extra["build_cache"], "enabled") {
		return errors.New("no-op snapshots require the build cache to be enabled")
	}
	if err := validateNoOpExternalInputs(siteConfig, builderMermaidMode); err != nil {
		return err
	}
	for _, hook := range siteConfig.Hooks {
		if hook == "default" {
			continue
		}
		if _, registered := plugins.PluginByName(hook); !registered {
			// A hook that is not registered in this builder binary cannot affect
			// its build output. Its source/config files remain part of the input
			// fingerprint, and queued refresh tasks still force a regular build.
			continue
		}
		switch hook {
		case "wikilink_hover", "webawesome", "image_optimization":
			// These built-in hooks are covered by the source/config/assets and
			// encoder fingerprints below.
		default:
			return fmt.Errorf("registered hook %q has no no-op input contract", hook)
		}
	}
	if siteConfig.Blogroll.GetRefreshOnBuild() {
		return errors.New("blogroll refresh_on_build is enabled")
	}
	return nil
}

func validateNoOpContentPatterns(patterns []string) error {
	for _, pattern := range patterns {
		if filepath.IsAbs(pattern) {
			return fmt.Errorf("absolute content pattern %q is not supported by no-op snapshots", pattern)
		}
		switch strings.ToLower(filepath.Ext(pattern)) {
		case ".md", ".markdown", ".mdx":
		default:
			return fmt.Errorf("content pattern %q has an unsupported extension for no-op checks", pattern)
		}
	}
	return nil
}

func validateNoOpExternalInputs(siteConfig *models.Config, builderMermaidMode []string) error {
	// Build-time outputs such as .well-known/time are preserved by the
	// no-op release check. Only inputs that can change independently of the
	// captured source, configuration, and cache fingerprints need guards here.
	if noOpTailwindBuildEnabled(siteConfig) {
		return errors.New("tailwind builds use an external executable and cache that are not covered by no-op snapshots")
	}
	if noOpIconVendorEnabled(siteConfig) {
		return errors.New("icon vendor downloads and publishes assets outside the no-op input snapshot")
	}
	if noOpPythonDocsEnabled(siteConfig) {
		return errors.New("python docs use an external interpreter that is not covered by no-op snapshots")
	}
	if noOpExternalLinkHoverEnabled(siteConfig) {
		return errors.New("external link hover uses expiring remote metadata that is not covered by no-op snapshots")
	}
	if siteConfig.Search.IsEnabled() && !strings.EqualFold(strings.TrimSpace(siteConfig.Search.SearchBackend()), "bleve") {
		return errors.New("pagefind can install a variable external binary and change the search index")
	}
	if noOpExternalMermaidRenderer(siteConfig, builderMermaidMode) {
		return errors.New("no-op snapshots support only the native Go Mermaid renderer")
	}
	if noOpExternalExcalidrawRenderer(siteConfig) {
		return errors.New("excalidraw chromium rendering uses an external browser not covered by no-op snapshots")
	}
	return nil
}

func (s *Service) noOpInputSnapshot(ctx context.Context) (noOpSnapshot, error) {
	return s.captureNoOpSnapshot(ctx)
}

func (s *Service) saveNoOpBaseline(ctx context.Context, snapshot noOpSnapshot, releaseID string, now time.Time) error {
	if snapshot.Fingerprint == "" || releaseID == "" {
		return errors.New("incomplete no-op baseline")
	}
	releaseFingerprint, _, err := s.noOpReleaseFingerprint(ctx, releaseID)
	if err != nil {
		return err
	}
	expires := now.Add(noOpMaxAge)
	local := now.In(now.Location())
	midnight := time.Date(local.Year(), local.Month(), local.Day()+1, 0, 0, 0, 0, local.Location())
	if midnight.Before(expires) {
		expires = midnight
	}
	if !snapshot.ExternalExpiresAt.IsZero() && snapshot.ExternalExpiresAt.Before(expires) {
		expires = snapshot.ExternalExpiresAt
	}
	if !now.Before(expires) {
		return errors.New("no-op baseline would already be expired")
	}
	data, err := json.Marshal(noOpBaseline{Version: noOpBaselineVersion, Fingerprint: snapshot.Fingerprint, ReleaseFingerprint: releaseFingerprint, ReleaseID: releaseID, CreatedAt: now.UTC(), ExpiresAt: expires.UTC()})
	if err != nil {
		return err
	}
	return writePrivateAtomic(s.noOpBaselinePath(), data)
}

func (s *Service) noOpBuilderFingerprint() (string, error) {
	path := s.executable
	if path == "" {
		var err error
		path, err = os.Executable()
		if err != nil {
			return "", fmt.Errorf("resolve builder executable: %w", err)
		}
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve builder executable path: %w", err)
	}
	file, err := os.Open(resolved)
	if err != nil {
		return "", fmt.Errorf("open builder executable: %w", err)
	}
	defer file.Close()
	before, err := file.Stat()
	if err != nil || !before.Mode().IsRegular() {
		return "", errors.New("builder executable is not a regular file")
	}
	identity := fileidentity.Read(before)
	if identity.Inode == 0 {
		return "", errors.New("builder executable identity is unavailable")
	}
	identityKey := fmt.Sprintf("%d:%d:%d:%d:%d:%d", identity.Device, identity.Inode, identity.Size, identity.Mode, identity.Modified, identity.Changed)
	s.noOpBuilderMu.Lock()
	defer s.noOpBuilderMu.Unlock()
	if s.noOpBuilderDigest != "" && s.noOpBuilderIdentity == identityKey {
		return s.noOpBuilderDigest, nil
	}
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("hash builder executable: %w", err)
	}
	after, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("recheck builder executable: %w", err)
	}
	if fileidentity.Read(before) != fileidentity.Read(after) {
		return "", errors.New("builder executable changed during fingerprinting")
	}
	digest := hex.EncodeToString(h.Sum(nil))
	s.noOpBuilderDigest = digest
	s.noOpBuilderIdentity = identityKey
	return digest, nil
}

func (s *Service) noOpCurrentReleasePath(releaseID string) (string, error) {
	if !plainReleaseID(releaseID) {
		return "", errors.New("current release ID is unsafe")
	}
	currentPath := filepath.Join(s.cfg.SiteDir, "current")
	target, err := os.Readlink(currentPath)
	if err != nil {
		return "", fmt.Errorf("read current release link: %w", err)
	}
	if filepath.IsAbs(target) || filepath.Clean(target) != filepath.Join("releases", releaseID) {
		return "", errors.New("current release link does not target its release directory")
	}
	releasePath := filepath.Join(s.cfg.SiteDir, "releases", releaseID)
	info, err := os.Lstat(releasePath)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", errors.New("current release is not a regular directory")
	}
	return releasePath, nil
}

func (s *Service) noOpReleaseFingerprint(ctx context.Context, releaseID string) (string, string, error) {
	fingerprint, releasePath, _, _, err := s.noOpReleaseFingerprintWithStats(ctx, releaseID)
	return fingerprint, releasePath, err
}

func (s *Service) noOpReleaseFingerprintWithStats(ctx context.Context, releaseID string) (string, string, int64, int64, error) {
	releasePath, err := s.noOpCurrentReleasePath(releaseID)
	if err != nil {
		return "", "", 0, 0, err
	}
	releasesDir := filepath.Join(s.cfg.SiteDir, "releases")
	rootInfo, err := os.Lstat(releasesDir)
	if err != nil || !rootInfo.IsDir() || rootInfo.Mode()&os.ModeSymlink != 0 {
		return "", "", 0, 0, errors.New("release root is not a regular directory")
	}
	info, err := os.Lstat(releasePath)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", "", 0, 0, errors.New("current release is not a regular directory")
	}
	rootResolved, err := filepath.EvalSymlinks(releasesDir)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("resolve release root: %w", err)
	}
	releaseResolved, err := filepath.EvalSymlinks(releasePath)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("resolve current release: %w", err)
	}
	relative, err := filepath.Rel(rootResolved, releaseResolved)
	if err != nil || relative != releaseID {
		return "", "", 0, 0, errors.New("current release resolves outside the release root")
	}
	if err := ctx.Err(); err != nil {
		return "", "", 0, 0, err
	}
	h := sha256.New()
	if err := writeNoOpIdentity(h, "release-root", releaseID, info); err != nil {
		return "", "", 0, 0, err
	}
	releaseScanStart := time.Now()
	releaseEntries, err := hashNoOpReleaseTree(ctx, h, releasePath)
	releaseMS := time.Since(releaseScanStart).Milliseconds()
	if err != nil {
		return "", "", releaseMS, releaseEntries, fmt.Errorf("fingerprint current release contents: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), releasePath, releaseMS, releaseEntries, nil
}

func hashNoOpReleaseTree(ctx context.Context, h hash.Hash, releasePath string) (int64, error) {
	paths, err := listNoOpReleasePaths(ctx, releasePath)
	if err != nil {
		return int64(len(paths)), err
	}
	entries, err := statNoOpReleasePaths(ctx, releasePath, paths)
	if err != nil {
		return int64(len(paths)), err
	}
	for _, entry := range entries {
		if err := writeNoOpIdentity(h, "release-entry", entry.name, entry.info); err != nil {
			return int64(len(paths)), err
		}
	}
	return int64(len(paths)), nil
}

func listNoOpReleasePaths(ctx context.Context, releasePath string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var paths []string
	err := filepath.WalkDir(releasePath, func(path string, _ fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == releasePath {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}

func statNoOpReleasePaths(ctx context.Context, releasePath string, paths []string) ([]noOpReleaseIdentity, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	entries := make([]noOpReleaseIdentity, len(paths))
	workers := runtime.GOMAXPROCS(0) * 2
	if workers < 8 {
		workers = 8
	}
	if workers > 32 {
		workers = 32
	}
	jobs := make(chan int)
	var wg sync.WaitGroup
	var errorOnce sync.Once
	var firstErr error
	setError := func(err error) {
		if err != nil {
			errorOnce.Do(func() { firstErr = err })
		}
	}
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				entry, err := statNoOpReleasePath(ctx, releasePath, paths[index])
				if err != nil {
					setError(err)
					continue
				}
				entries[index] = entry
			}
		}()
	}
	for index := range paths {
		select {
		case <-ctx.Done():
			setError(ctx.Err())
			close(jobs)
			wg.Wait()
			return nil, firstErr
		case jobs <- index:
		}
	}
	close(jobs)
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func statNoOpReleasePath(ctx context.Context, releaseRoot, path string) (noOpReleaseIdentity, error) {
	if err := ctx.Err(); err != nil {
		return noOpReleaseIdentity{}, err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return noOpReleaseIdentity{}, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return noOpReleaseIdentity{}, fmt.Errorf("release contains unsupported symlink %q", path)
	}
	if !info.IsDir() && !info.Mode().IsRegular() {
		return noOpReleaseIdentity{}, fmt.Errorf("release contains unsupported file type %q", path)
	}
	relative, err := filepath.Rel(releaseRoot, path)
	if err != nil {
		return noOpReleaseIdentity{}, err
	}
	return noOpReleaseIdentity{name: filepath.ToSlash(relative), info: info}, nil
}

func forEachNoOpContentFile(ctx context.Context, sourceDir string, siteConfig *models.Config, visit func(string, []byte) (bool, error)) error {
	files, err := sourcegit.Command(ctx, sourceDir, "ls-files", "--cached", "--others", "--exclude-standard", "-z").Output()
	if err != nil {
		return fmt.Errorf("list content files for no-op verification: %w", err)
	}
	for _, name := range strings.Split(string(files), "\x00") {
		if !isNoOpContentFile(siteConfig, name) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		path := filepath.Join(sourceDir, filepath.FromSlash(name))
		info, err := os.Lstat(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("inspect content source %q: %w", name, err)
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() > 32<<20 {
			return fmt.Errorf("content source %q is unsafe", name)
		}
		data, err := os.ReadFile(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read content source %q: %w", name, err)
		}
		done, err := visit(name, data)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
	return nil
}

func isNoOpContentFile(siteConfig *models.Config, name string) bool {
	if name == "" || !noOpContentPathMatches(siteConfig, name) {
		return false
	}
	switch strings.ToLower(filepath.Ext(name)) {
	case ".md", ".markdown", ".mdx":
		return true
	default:
		return false
	}
}

func checkNoOpMentionFreshness(ctx context.Context, sourceDir string, siteConfig *models.Config, now time.Time) (time.Time, error) {
	if siteConfig == nil || extraBoolIsTrue(siteConfig.Extra, "mentions_disabled") {
		return time.Time{}, nil
	}
	mentionsConfig := noOpMentionsConfig(siteConfig)
	if !mentionsConfig.IsEnabled() {
		return time.Time{}, nil
	}
	hasCandidates, err := hasNoOpMentionCandidates(ctx, sourceDir, siteConfig)
	if err != nil {
		return time.Time{}, err
	}
	if !hasCandidates {
		return time.Time{}, nil
	}
	return time.Time{}, errors.New("mention cache domain coverage cannot be proven")
}

func noOpMentionsConfig(siteConfig *models.Config) models.MentionsConfig {
	if siteConfig != nil && siteConfig.Extra != nil {
		switch value := siteConfig.Extra["mentions"].(type) {
		case models.MentionsConfig:
			return value
		case *models.MentionsConfig:
			if value != nil {
				return *value
			}
		default:
			if _, exists := siteConfig.Extra["mentions"]; exists {
				return models.NewMentionsConfig()
			}
		}
	}
	if siteConfig != nil {
		return siteConfig.Mentions
	}
	return models.NewMentionsConfig()
}

func hasNoOpMentionCandidates(ctx context.Context, sourceDir string, siteConfig *models.Config) (bool, error) {
	found := false
	err := forEachNoOpContentFile(ctx, sourceDir, siteConfig, func(_ string, data []byte) (bool, error) {
		found = noOpMentionCandidate.Match(data)
		return found, nil
	})
	return found, err
}

func extraBoolIsTrue(extra map[string]any, name string) bool {
	if extra == nil {
		return false
	}
	value, ok := extra[name].(bool)
	return ok && value
}

func extraBoolIsFalseMap(value any, name string) bool {
	config, ok := value.(map[string]any)
	if !ok {
		return false
	}
	field, ok := config[name].(bool)
	return ok && !field
}

func checkNoOpEmbedCache(ctx context.Context, sourceDir, cacheMount string, siteConfig *models.Config) error {
	if siteConfig == nil || embedsDisabled(siteConfig.Extra) {
		return nil
	}
	cache, err := buildcache.Load(noOpBuildCacheDir(sourceDir, cacheMount, siteConfig))
	if err != nil {
		return fmt.Errorf("load build cache for embed verification: %w", err)
	}
	return forEachNoOpContentFile(ctx, sourceDir, siteConfig, func(name string, data []byte) (bool, error) {
		return false, checkNoOpEmbedCacheEntry(name, data, cache, siteConfig)
	})
}

func checkNoOpEmbedCacheEntry(name string, data []byte, cache *buildcache.Cache, siteConfig *models.Config) error {
	metadata, body, err := plugins.ParseFrontmatter(string(data))
	if err != nil {
		return fmt.Errorf("parse embed source %q: %w", name, err)
	}
	post := models.NewPost(filepath.ToSlash(name))
	post.Content = body
	post.Template = noOpMetadataString(metadata, "template")
	post.Private = noOpMetadataBool(metadata, "private")
	post.Skip = noOpMetadataBool(metadata, "skip")
	for _, key := range []string{"link", "url"} {
		if value := noOpMetadataString(metadata, key); value != "" {
			post.Set(key, value)
		}
	}
	if !plugins.HasExternalEmbedCandidate(post) {
		return nil
	}
	post, err = plugins.ParsePostFromContentWithConfig(filepath.ToSlash(name), string(data), siteConfig)
	if err != nil {
		return fmt.Errorf("load embed source %q: %w", name, err)
	}
	if post.Skip || !plugins.HasExternalEmbedCandidate(post) {
		return nil
	}
	if !plugins.HasCachedEmbedTransform(cache, post, siteConfig.Extra) {
		return fmt.Errorf("external embed transform cache is missing or stale for %q", name)
	}
	return nil
}

func noOpHookEnabled(siteConfig *models.Config, name string) bool {
	if siteConfig == nil {
		return false
	}
	for _, disabled := range siteConfig.DisabledHooks {
		if strings.EqualFold(strings.TrimSpace(disabled), name) {
			return false
		}
	}
	if len(siteConfig.Hooks) == 0 {
		return true
	}
	for _, hook := range siteConfig.Hooks {
		hook = strings.ToLower(strings.TrimSpace(hook))
		if hook == "default" || hook == name {
			return true
		}
	}
	return false
}

func noOpContentPathMatches(siteConfig *models.Config, name string) bool {
	if siteConfig == nil {
		return false
	}
	path := filepath.ToSlash(name)
	for _, pattern := range siteConfig.GlobConfig.Patterns {
		if doublestar.MatchUnvalidated(filepath.ToSlash(pattern), path) {
			return true
		}
	}
	return false
}

func checkNoOpCDNAssetCache(sourceDir string, siteConfig *models.Config) error {
	if siteConfig == nil {
		return errors.New("site config is required for CDN asset verification")
	}
	requested, _ := siteConfig.Extra["cdn_assets_extra"].([]assets.Asset)
	if !siteConfig.Assets.IsSelfHosted() && len(requested) == 0 {
		return nil
	}
	required := requested
	if siteConfig.Assets.IsSelfHosted() {
		base := assets.Registry()
		byName := make(map[string]int, len(base)+len(requested))
		for index, asset := range base {
			byName[asset.Name] = index
		}
		required = append([]assets.Asset(nil), base...)
		for _, asset := range requested {
			if index, exists := byName[asset.Name]; exists {
				required[index] = asset
				continue
			}
			byName[asset.Name] = len(required)
			required = append(required, asset)
		}
	}
	cacheDir := siteConfig.Assets.GetCacheDir()
	if !filepath.IsAbs(cacheDir) {
		cacheDir = filepath.Join(sourceDir, cacheDir)
	}
	downloader := assets.NewDownloader(cacheDir, siteConfig.Assets.IsVerifyIntegrityEnabled())
	for _, asset := range required {
		if !downloader.IsCached(asset) {
			return fmt.Errorf("CDN asset cache is incomplete for %q", asset.Name)
		}
	}
	return nil
}

func checkNoOpLinkAvatarCache(ctx context.Context, sourceDir, cacheMount string, siteConfig *models.Config) error {
	if siteConfig == nil || !plugins.LinkAvatarsRequiresBuildTimeRender(siteConfig.Extra) {
		return nil
	}
	cache, err := buildcache.Load(noOpBuildCacheDir(sourceDir, cacheMount, siteConfig))
	if err != nil {
		return fmt.Errorf("load build cache for link-avatar verification: %w", err)
	}
	return forEachNoOpContentFile(ctx, sourceDir, siteConfig, func(name string, data []byte) (bool, error) {
		return false, checkNoOpLinkAvatarCacheEntry(name, data, cache, siteConfig)
	})
}

func checkNoOpLinkAvatarCacheEntry(name string, data []byte, cache *buildcache.Cache, siteConfig *models.Config) error {
	post, err := plugins.ParsePostFromContentWithConfig(filepath.ToSlash(name), string(data), siteConfig)
	if err != nil {
		return fmt.Errorf("load link-avatar source %q: %w", name, err)
	}
	if post.Skip || (strings.TrimSpace(post.Content) == "" && !noOpPostHasLink(post)) {
		return nil
	}
	cached := cache.GetCachedPost(post.Path)
	if cached == nil || cached.LinkAvatarsHash == "" || cached.LinkAvatarsHTML == "" {
		return fmt.Errorf("link-avatar transform cache is missing for %q", name)
	}
	return nil
}

func noOpPostHasLink(post *models.Post) bool {
	if post == nil {
		return false
	}
	for _, key := range []string{"link", "url"} {
		if value, ok := post.Get(key).(string); ok && strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}

func embedsDisabled(extra map[string]any) bool {
	if extra == nil {
		return false
	}
	config, ok := extra["embeds"].(map[string]any)
	if !ok {
		return false
	}
	enabled, ok := config["enabled"].(bool)
	return ok && !enabled
}

func noOpBuildCacheDir(sourceDir, cacheMount string, siteConfig *models.Config) string {
	configured := ""
	if siteConfig != nil && siteConfig.Extra != nil {
		configured, _ = siteConfig.Extra["cache_dir"].(string)
	}
	if strings.TrimSpace(configured) != "" {
		if filepath.IsAbs(configured) {
			return filepath.Clean(configured)
		}
		return filepath.Join(sourceDir, configured)
	}
	if cacheMount != "" {
		if filepath.IsAbs(cacheMount) {
			return filepath.Join(cacheMount, "build")
		}
		return filepath.Join(sourceDir, cacheMount, "build")
	}
	return filepath.Join(sourceDir, ".markata")
}

func noOpBuildCacheFile(sourceDir, cacheMount string, siteConfig *models.Config) string {
	return filepath.Join(noOpBuildCacheDir(sourceDir, cacheMount, siteConfig), buildcache.CacheFileName)
}

func noOpMetadataString(metadata map[string]any, key string) string {
	value, _ := metadata[key].(string)
	return value
}

func noOpMetadataBool(metadata map[string]any, key string) bool {
	value, _ := metadata[key].(bool)
	return value
}

func (s *Service) loadNoOpBaseline() (noOpBaseline, error) {
	var baseline noOpBaseline
	path := s.noOpBaselinePath()
	info, err := os.Lstat(path)
	if err != nil {
		return baseline, err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 || info.Size() > 4096 {
		return baseline, errors.New("no-op baseline is invalid or not private")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return baseline, err
	}
	if err := json.Unmarshal(data, &baseline); err != nil {
		return baseline, err
	}
	if baseline.Version != noOpBaselineVersion || baseline.CreatedAt.IsZero() || !baseline.ExpiresAt.After(baseline.CreatedAt) || baseline.ReleaseFingerprint == "" {
		return noOpBaseline{}, errors.New("no-op baseline has invalid timestamps or version")
	}
	return baseline, nil
}

func (s *Service) noOpKeyPath() string {
	return s.statePath + ".noop-key"
}

func (s *Service) noOpBaselinePath() string {
	return s.statePath + ".noop-baseline.json"
}

func (s *Service) loadOrCreateNoOpKey() ([]byte, error) {
	path := s.noOpKeyPath()
	data, err := os.ReadFile(path)
	if err == nil {
		info, statErr := os.Lstat(path)
		if statErr != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 || len(data) != noOpKeySize {
			return nil, errors.New("no-op key is invalid or not private")
		}
		return data, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("read no-op key: %w", err)
	}
	key := make([]byte, noOpKeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("generate no-op key: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return s.loadOrCreateNoOpKey()
	}
	if err != nil {
		return nil, fmt.Errorf("create no-op key: %w", err)
	}
	_, writeErr := file.Write(key)
	syncErr := file.Sync()
	closeErr := file.Close()
	if writeErr != nil || syncErr != nil || closeErr != nil {
		_ = os.Remove(path)
		return nil, errors.Join(writeErr, syncErr, closeErr)
	}
	return key, nil
}

func writePrivateAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-")
	if err != nil {
		return err
	}
	tmp := file.Name()
	defer os.Remove(tmp)
	if err := file.Chmod(0o600); err != nil {
		_ = file.Close()
		return err
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func noOpInputRoots(sourceDir, configPath, cacheMount string, siteConfig *models.Config) []string {
	seen := make(map[string]struct{})
	var dirs []string
	add := func(value string) {
		if strings.TrimSpace(value) == "" {
			return
		}
		path := value
		if !filepath.IsAbs(path) {
			path = filepath.Join(sourceDir, path)
		}
		path = filepath.Clean(path)
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		dirs = append(dirs, path)
	}
	add(configPath)
	add(filepath.Join(sourceDir, ".env"))
	add(configPath)
	for _, path := range siteConfig.ExternalCacheDirs() {
		add(path)
	}
	add(noOpBuildCacheFile(sourceDir, cacheMount, siteConfig))
	mentionsConfig := noOpMentionsConfig(siteConfig)
	add(mentionsConfig.GetCacheDir())
	if siteConfig.Extra != nil {
		if cacheDir := extraStringField(siteConfig.Extra["webmentions"], "cache_dir"); cacheDir != "" {
			add(cacheDir)
		}
	}
	add(siteConfig.Assets.GetCacheDir())
	add(runtimeenv.BundledAssetsCacheDir())
	imageCacheDir := extraStringField(siteConfig.Extra["image_optimization"], "cache_dir")
	if imageCacheDir == "" {
		imageCacheDir = ".markata/image-cache"
	}
	add(imageCacheDir)
	add(siteConfig.AssetsDir)
	add(siteConfig.TemplatesDir)
	if siteConfig.FontpacksFile != "" {
		fontpacksPath := siteConfig.FontpacksFile
		if !filepath.IsAbs(fontpacksPath) {
			fontpacksPath = filepath.Join(filepath.Dir(configPath), fontpacksPath)
		}
		add(fontpacksPath)
	}
	sort.Strings(dirs)
	return dirs
}

func noOpPathCoveredByRoots(path string, roots []string) bool {
	path = filepath.Clean(path)
	for _, root := range roots {
		relative, err := filepath.Rel(filepath.Clean(root), path)
		if err != nil {
			continue
		}
		if relative == "." || relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func noOpExternalThemeRoots(sourceDir, executable, themeName string) ([]string, error) {
	if executable == "" {
		var err error
		executable, err = os.Executable()
		if err != nil {
			return nil, fmt.Errorf("resolve builder executable for theme inputs: %w", err)
		}
	}
	executable, err := filepath.EvalSymlinks(executable)
	if err != nil {
		return nil, fmt.Errorf("resolve builder executable for theme inputs: %w", err)
	}
	executableDir := filepath.Dir(executable)
	shareDir := filepath.Join(filepath.Dir(executableDir), "share", "markata-go", "themes")
	if themeName == "" {
		themeName = "default"
	}
	candidates := []string{
		filepath.Join(executableDir, "themes", themeName, "templates"),
		filepath.Join(executableDir, "themes", "default", "templates"),
		filepath.Join(shareDir, "default", "templates"),
		filepath.Join(executableDir, "themes", themeName, "static"),
		filepath.Join(shareDir, themeName, "static"),
	}
	roots := make([]string, 0, len(candidates))
	for _, path := range candidates {
		info, err := os.Stat(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("inspect external theme input %q: %w", path, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("external theme input %q is not a directory", path)
		}
		if !pathWithin(sourceDir, path) {
			roots = append(roots, path)
		}
	}
	return roots, nil
}

func pathWithin(root, path string) bool {
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	return err == nil && (relative == "." || relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)))
}

func compactNoOpInputRoots(roots []string) []string {
	unique := make(map[string]struct{}, len(roots))
	paths := make([]string, 0, len(roots))
	for _, path := range roots {
		path = filepath.Clean(path)
		if _, ok := unique[path]; ok {
			continue
		}
		unique[path] = struct{}{}
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) == len(paths[j]) {
			return paths[i] < paths[j]
		}
		return len(paths[i]) < len(paths[j])
	})
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if !noOpPathCoveredByRoots(path, result) {
			result = append(result, path)
		}
	}
	sort.Strings(result)
	return result
}

func hashNoOpImageEncoders(h hash.Hash, sourceDir string, rawConfig any) error {
	if enabled, ok := extraBoolField(rawConfig, "enabled"); ok && !enabled {
		return nil
	}
	for _, encoder := range []struct {
		name      string
		configKey string
	}{
		{name: "avifenc", configKey: "avifenc_path"},
		{name: "cwebp", configKey: "cwebp_path"},
	} {
		path := extraStringField(rawConfig, encoder.configKey)
		if path == "" {
			found, err := exec.LookPath(encoder.name)
			if err == nil {
				path = found
			}
		}
		writeNoOpPart(h, "image-encoder-name", encoder.name)
		if path == "" {
			writeNoOpPart(h, "image-encoder", "missing")
			continue
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(sourceDir, path)
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("resolve image encoder %q: %w", encoder.name, err)
		}
		info, err := os.Stat(resolved)
		if err != nil {
			return fmt.Errorf("inspect image encoder %q: %w", encoder.name, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("image encoder %q is not a regular file", encoder.name)
		}
		writeNoOpPart(h, "image-encoder-path", resolved)
		if err := writeNoOpIdentity(h, "image-encoder", resolved, info); err != nil {
			return err
		}
	}
	return nil
}

func extraBoolField(value any, name string) (bool, bool) {
	config, ok := value.(map[string]any)
	if !ok {
		return false, false
	}
	field, ok := config[name].(bool)
	return field, ok
}

func noOpTailwindBuildEnabled(siteConfig *models.Config) bool {
	if siteConfig == nil {
		return true
	}
	if enabled, ok := extraBoolField(siteConfig.Extra["tailwind"], "build"); ok {
		return enabled
	}
	// Tailwind builds by default. Require an explicit opt-out until the
	// executable and installer cache are included in the no-op input contract.
	return true
}

func noOpIconVendorEnabled(siteConfig *models.Config) bool {
	if siteConfig == nil || siteConfig.Extra == nil {
		return false
	}
	icons, ok := siteConfig.Extra["icons"].(map[string]any)
	if !ok {
		return false
	}
	vendor, ok := icons["vendor"].(map[string]any)
	if !ok {
		return false
	}
	enabled, _ := vendor["enabled"].(bool)
	return enabled
}

func noOpPythonDocsEnabled(siteConfig *models.Config) bool {
	if siteConfig == nil || siteConfig.Extra == nil {
		return false
	}
	enabled, _ := extraBoolField(siteConfig.Extra["python_docs"], "enabled")
	return enabled
}

func noOpExternalLinkHoverEnabled(siteConfig *models.Config) bool {
	if siteConfig == nil || siteConfig.Extra == nil {
		return false
	}
	enabled, _ := extraBoolField(siteConfig.Extra["external_link_hover"], "enabled")
	return enabled
}

func noOpExternalMermaidRenderer(siteConfig *models.Config, builderModes []string) bool {
	if siteConfig == nil {
		return true
	}
	mode := "go"
	if siteConfig.Extra != nil {
		raw := siteConfig.Extra["mermaid"]
		switch config := raw.(type) {
		case models.MermaidConfig:
			if !config.Enabled {
				return false
			}
			mode = config.Mode
		case *models.MermaidConfig:
			if config != nil {
				if !config.Enabled {
					return false
				}
				mode = config.Mode
			}
		default:
			if enabled, ok := extraBoolField(raw, "enabled"); ok && !enabled {
				return false
			}
			mode = extraStringField(raw, "mode")
		}
	}
	if strings.TrimSpace(mode) == "" {
		mode = "go"
	}
	if len(builderModes) > 0 && strings.TrimSpace(builderModes[0]) != "" {
		mode = strings.TrimSpace(builderModes[0])
	}
	// Go mode renders deterministic SVGs in-process. Client mode only keeps
	// the Mermaid source in the page for the browser script, so its build-time
	// output is also deterministic and covered by the source/config snapshots.
	// CLI and Chromium modes depend on external renderers whose output is not
	// covered by the no-op fingerprints.
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "go", "client":
		return false
	default:
		return true
	}
}

func noOpExternalExcalidrawRenderer(siteConfig *models.Config) bool {
	if siteConfig == nil || siteConfig.Extra == nil {
		return false
	}
	raw := siteConfig.Extra["excalidraw"]
	switch config := raw.(type) {
	case models.ExcalidrawConfig:
		return config.Enabled && strings.EqualFold(strings.TrimSpace(config.Mode), "chromium")
	case *models.ExcalidrawConfig:
		return config != nil && config.Enabled && strings.EqualFold(strings.TrimSpace(config.Mode), "chromium")
	}
	if enabled, ok := extraBoolField(raw, "enabled"); ok && !enabled {
		return false
	}
	mode := extraStringField(raw, "mode")
	return strings.EqualFold(strings.TrimSpace(mode), "chromium")
}

func extraStringField(value any, name string) string {
	switch config := value.(type) {
	case map[string]any:
		if field, ok := config[name].(string); ok {
			return field
		}
	case models.WebMentionsConfig:
		if name == "cache_dir" {
			return config.CacheDir
		}
	case *models.WebMentionsConfig:
		if config != nil && name == "cache_dir" {
			return config.CacheDir
		}
	}
	return ""
}

func hashNoOpTree(ctx context.Context, h hash.Hash, root string, entryCount *int64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if filepath.Clean(root) == filepath.Clean(".") || filepath.Clean(root) == filepath.Clean("/") {
		return fmt.Errorf("unsafe no-op cache path %q", root)
	}
	info, err := os.Lstat(root)
	if errors.Is(err, fs.ErrNotExist) {
		writeNoOpPart(h, "cache-missing", root)
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect no-op cache %q: %w", root, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("no-op cache root %q is a symlink", root)
	}
	if !info.IsDir() {
		return writeNoOpIdentity(h, "input-file", root, info)
	}
	if err := writeNoOpIdentity(h, "input-dir", root, info); err != nil {
		return err
	}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if path == root {
			return nil
		}
		if entryCount != nil {
			*entryCount++
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if err := writeNoOpIdentity(h, "input-entry", filepath.ToSlash(strings.TrimPrefix(path, root+string(os.PathSeparator))), info); err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("no-op input %q contains a symlink", path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk no-op cache %q: %w", root, err)
	}
	return nil
}

func writeNoOpIdentity(h hash.Hash, label, name string, info fs.FileInfo) error {
	identity := fileidentity.Read(info)
	if identity.Inode == 0 {
		return fmt.Errorf("file identity unavailable for no-op input %q", name)
	}
	writeNoOpPart(h, label+"-path", name)
	var encoded [6 * 8]byte
	binary.BigEndian.PutUint64(encoded[0:8], identity.Device)
	binary.BigEndian.PutUint64(encoded[8:16], identity.Inode)
	binary.BigEndian.PutUint64(encoded[16:24], uint64(identity.Size))
	binary.BigEndian.PutUint64(encoded[24:32], uint64(identity.Mode))
	binary.BigEndian.PutUint64(encoded[32:40], uint64(identity.Modified))
	binary.BigEndian.PutUint64(encoded[40:48], uint64(identity.Changed))
	_, _ = h.Write(encoded[:])
	return nil
}

func writeNoOpPart(h hash.Hash, name, value string) {
	_, _ = fmt.Fprintf(h, "%d:%s%d:%s", len(name), name, len(value), value)
}

func buildCacheSaveFailed(logPath string) (bool, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return false, fmt.Errorf("open build log for cache verification: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, noOpBuildLogLimit+1))
	if err != nil {
		return false, fmt.Errorf("read build log for cache verification: %w", err)
	}
	if len(data) > noOpBuildLogLimit {
		return false, errors.New("build log exceeds the cache verification limit")
	}
	return bytes.Contains(data, []byte("failed to save build cache")), nil
}

func checkTemplateTimeUse(sourceDir, templatesDir string) error {
	if templatesDir == "" {
		templatesDir = "templates"
	}
	root := templatesDir
	if !filepath.IsAbs(root) {
		root = filepath.Join(sourceDir, root)
	}
	if _, err := os.Stat(root); errors.Is(err, fs.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("inspect site templates: %w", err)
	}
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !isTemplateText(path) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := validateNoOpTimeExpressions(data); err != nil {
			return fmt.Errorf("template %q has unsupported time-dependent expression: %w", path, err)
		}
		return nil
	})
}

func checkJinjaTimeUse(ctx context.Context, sourceDir string, siteConfig *models.Config) error {
	return forEachNoOpContentFile(ctx, sourceDir, siteConfig, func(name string, data []byte) (bool, error) {
		return false, checkNoOpJinjaSource(name, data)
	})
}

func checkNoOpJinjaSource(name string, data []byte) error {
	frontmatter, ok := markdownFrontmatter(data)
	if !ok || !noOpJinjaFlag.Match(frontmatter) || !noOpNowToken.Match(data) {
		return nil
	}
	if err := validateNoOpTimeExpressions(data); err != nil {
		return fmt.Errorf("jinja-enabled source %q has unsupported time-dependent expression: %w", name, err)
	}
	return nil
}

func validateNoOpTimeExpressions(data []byte) error {
	for offset := 0; offset < len(data); {
		openOffset := bytes.Index(data[offset:], []byte("{{"))
		if openOffset < 0 {
			break
		}
		open := offset + openOffset
		closeOffset := bytes.Index(data[open+2:], []byte("}}"))
		if closeOffset < 0 {
			break
		}
		close := open + 2 + closeOffset + 2
		expression := data[open:close]
		if !noOpNowToken.Match(expression) {
			offset = close
			continue
		}
		if len(noOpNowToken.FindAllIndex(expression, -1)) != 1 {
			return errors.New("multiple now references are unsupported")
		}
		formats := noOpDateFormat.FindAllSubmatch(expression, -1)
		if len(formats) != 1 || len(formats[0]) != 2 || (string(formats[0][1]) != "2006" && string(formats[0][1]) != "2006-01-02") {
			return errors.New("only year or calendar-date formatting is supported")
		}
		offset = close
	}
	return nil
}

func markdownFrontmatter(data []byte) ([]byte, bool) {
	lines := strings.SplitN(string(data), "\n", 2)
	if len(lines) != 2 || strings.TrimSpace(strings.TrimPrefix(lines[0], "\ufeff")) != "---" {
		return nil, false
	}
	end := strings.Index(lines[1], "\n---")
	if end < 0 {
		return nil, false
	}
	return []byte(lines[1][:end]), true
}

func isTemplateText(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".html", ".htm", ".xml", ".txt", ".pongo":
		return true
	default:
		return false
	}
}
