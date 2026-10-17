package builderadmin

import (
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/diagnostics"
	"github.com/WaylonWalker/markata-go/pkg/servecontrol"
	"github.com/WaylonWalker/markata-go/pkg/sourcegit"
	"github.com/fsnotify/fsnotify"
)

const (
	defaultLogDirName   = "logs"
	defaultStateName    = "state.json"
	defaultOverrideName = "overrides"
	defaultLeaderName   = "leader.json"
	defaultLockName     = "leader.lock"
	defaultListenHost   = "127.0.0.1"
	defaultListenPort   = 8080
	defaultReleaseKeep  = 25
	maxWebhookBodyBytes = 1 << 20
)

var ErrBuildQueueFull = errors.New("build queue is full")

type Config struct {
	Host                 string
	Port                 int
	SourceDir            string
	BuilderVersion       string
	BuilderCommit        string
	SiteDir              string
	WorkDir              string
	ConfigPath           string
	CacheMount           string
	HistoryDir           string
	WatchEnabled         bool
	WatchDebounce        time.Duration
	Fast                 bool
	MermaidMode          string
	ReleasesKeep         int
	SuccessfulBuildsKeep int
	FailedBuildsKeep     int
	RefreshRunsKeep      int
	RefreshTasks         []RefreshTaskConfig
	BuildTimeout         time.Duration
	TrustedProxyCIDRs    []string
	AuthHeaders          AuthHeaders
	PublicAuthOrigin     string
	PublicOrigin         string
	PreviewOrigin        string
	Webhook              WebhookConfig
}

// WebhookConfig configures a signed GitHub or Forgejo push webhook.
// Secret is deliberately excluded from JSON responses served by the admin API.
type WebhookConfig struct {
	Enabled bool   `json:"enabled"`
	Branch  string `json:"branch"`
	Secret  string `json:"-"`
}

type RefreshTaskConfig struct {
	Name                  string   `json:"name"`
	Every                 string   `json:"every"`
	EnqueueBuildOnSuccess bool     `json:"enqueue_build_on_success"`
	Args                  []string `json:"args"`

	interval time.Duration
}

type State struct {
	Queue   []QueuedOperation `json:"queue"`
	Running *RunningOperation `json:"running,omitempty"`
	Builds  []BuildRecord     `json:"builds"`
	Refresh []RefreshRecord   `json:"refresh"`
}

type QueuedOperation struct {
	ID          string    `json:"id"`
	Kind        string    `json:"kind"`
	ActorID     string    `json:"actor_id,omitempty"`
	Label       string    `json:"label"`
	TriggerType string    `json:"trigger_type"`
	Detail      string    `json:"detail,omitempty"`
	Changed     []string  `json:"changed,omitempty"`
	EnqueuedAt  time.Time `json:"enqueued_at"`
	ReleaseID   string    `json:"release_id,omitempty"`
	TaskName    string    `json:"task_name,omitempty"`
	SyncSource  bool      `json:"sync_source,omitempty"`
}

type RunningOperation struct {
	ID          string    `json:"id"`
	Kind        string    `json:"kind"`
	ActorID     string    `json:"actor_id,omitempty"`
	Label       string    `json:"label"`
	TriggerType string    `json:"trigger_type"`
	Detail      string    `json:"detail,omitempty"`
	StartedAt   time.Time `json:"started_at"`
	EnqueuedAt  time.Time `json:"enqueued_at,omitempty"`
	Phase       string    `json:"phase"`
	Impact      string    `json:"impact,omitempty"`
}

type BuildRecord struct {
	ID              string                    `json:"id"`
	Kind            string                    `json:"kind"`
	ActorID         string                    `json:"actor_id,omitempty"`
	Status          string                    `json:"status"`
	TriggerType     string                    `json:"trigger_type"`
	TriggerDetail   string                    `json:"trigger_detail,omitempty"`
	ChangedPaths    []string                  `json:"changed_paths,omitempty"`
	Impact          string                    `json:"impact,omitempty"`
	EnqueuedAt      time.Time                 `json:"enqueued_at"`
	StartedAt       time.Time                 `json:"started_at"`
	FinishedAt      time.Time                 `json:"finished_at"`
	QueueWaitMS     int64                     `json:"queue_wait_ms"`
	PrepareMS       int64                     `json:"prepare_ms"`
	BuildMS         int64                     `json:"build_ms"`
	PromoteMS       int64                     `json:"promote_ms"`
	PruneMS         int64                     `json:"prune_ms"`
	TotalMS         int64                     `json:"total_ms"`
	ReleaseID       string                    `json:"release_id,omitempty"`
	ReleasePath     string                    `json:"release_path,omitempty"`
	BecameLive      bool                      `json:"became_live"`
	LogPath         string                    `json:"log_path,omitempty"`
	LogDiagnostics  []servecontrol.Diagnostic `json:"log_diagnostics,omitempty"`
	PerfSummary     []string                  `json:"perf_summary,omitempty"`
	Error           string                    `json:"error,omitempty"`
	RollbackRelease string                    `json:"rollback_release,omitempty"`
	BuilderVersion  string                    `json:"builder_version,omitempty"`
	BuilderCommit   string                    `json:"builder_commit,omitempty"`
	SourceCommit    string                    `json:"source_commit,omitempty"`
}

type RefreshRecord struct {
	ID                    string                    `json:"id"`
	ActorID               string                    `json:"actor_id,omitempty"`
	TaskName              string                    `json:"task_name"`
	Status                string                    `json:"status"`
	TriggerType           string                    `json:"trigger_type"`
	TriggerDetail         string                    `json:"trigger_detail,omitempty"`
	EnqueuedAt            time.Time                 `json:"enqueued_at"`
	StartedAt             time.Time                 `json:"started_at"`
	FinishedAt            time.Time                 `json:"finished_at"`
	QueueWaitMS           int64                     `json:"queue_wait_ms"`
	RunMS                 int64                     `json:"run_ms"`
	TotalMS               int64                     `json:"total_ms"`
	LogPath               string                    `json:"log_path,omitempty"`
	LogDiagnostics        []servecontrol.Diagnostic `json:"log_diagnostics,omitempty"`
	EnqueuedBuildID       string                    `json:"enqueued_build_id,omitempty"`
	EnqueueBuildOnSuccess bool                      `json:"enqueue_build_on_success"`
	Command               []string                  `json:"command,omitempty"`
	Error                 string                    `json:"error,omitempty"`
}

type ReleaseView struct {
	ID            string    `json:"id"`
	Path          string    `json:"path"`
	CreatedAt     time.Time `json:"created_at"`
	Current       bool      `json:"current"`
	BuildID       string    `json:"build_id,omitempty"`
	BuildStatus   string    `json:"build_status,omitempty"`
	BuilderCommit string    `json:"builder_commit,omitempty"`
	SourceCommit  string    `json:"source_commit,omitempty"`
	RollbackOnly  bool      `json:"rollback_only"`
}

type completedJobView struct {
	ID          string
	Kind        string
	Status      string
	Trigger     string
	FinishedAt  time.Time
	QueueWaitMS int64
	RunMS       int64
	Release     string
	LogPath     string
	Build       *BuildRecord
}

type Service struct {
	cfg                  Config
	executable           string
	statePath            string
	logDir               string
	overrideDir          string
	leaderPath           string
	queueCh              chan queueRequest
	watchMu              sync.Mutex
	watchChanged         map[string]struct{}
	watchTimer           *time.Timer
	stateMu              sync.Mutex
	state                State
	controlRuntime       *servecontrol.Runtime
	releaseMu            sync.Mutex
	pruneScheduleMu      sync.Mutex
	pruneRunning         bool
	prunePending         []pruneRequest
	pruneCancel          context.CancelFunc
	pruneDone            chan struct{}
	leaderMu             sync.RWMutex
	leader               bool
	leaderCancel         context.CancelFunc
	leaderLock           *os.File
	instanceID           string
	instanceAddr         string
	server               *http.Server
	trustedProxyPrefixes []netip.Prefix
	theme                uiTheme
}

type leaderRecord struct {
	InstanceID string    `json:"instance_id"`
	Addr       string    `json:"addr"`
	AcquiredAt time.Time `json:"acquired_at"`
}

type queueRequest struct {
	QueuedOperation
	commandArgs []string
}

type pruneRequest struct {
	ctx     context.Context
	buildID string
}

//nolint:gocyclo // Configuration normalization is intentionally kept together at construction.
func New(cfg Config) (*Service, error) {
	trustedProxyPrefixes, err := parseTrustedProxyPrefixes(cfg.TrustedProxyCIDRs)
	if err != nil {
		return nil, err
	}
	authHeaders, err := normalizeAuthHeaders(cfg.AuthHeaders)
	if err != nil {
		return nil, err
	}
	cfg.AuthHeaders = authHeaders
	publicOrigin, err := normalizePublicAuthOrigin(cfg.PublicOrigin)
	if err != nil {
		return nil, fmt.Errorf("public origin: %w", err)
	}
	cfg.PublicOrigin = publicOrigin
	publicAuthOrigin, err := normalizePublicAuthOrigin(cfg.PublicAuthOrigin)
	if err != nil {
		return nil, err
	}
	cfg.PublicAuthOrigin = publicAuthOrigin
	if cfg.Host == "" {
		cfg.Host = defaultListenHost
	}
	if cfg.Port == 0 {
		cfg.Port = defaultListenPort
	}
	if cfg.SourceDir == "" {
		cfg.SourceDir = "."
	}
	if cfg.SiteDir == "" {
		cfg.SiteDir = "public"
	}
	if cfg.WorkDir == "" {
		cfg.WorkDir = filepath.Join(cfg.SiteDir, ".build-work")
	}
	if cfg.HistoryDir == "" {
		cfg.HistoryDir = filepath.Join(cfg.SiteDir, ".builder-admin")
	}
	if cfg.WatchDebounce <= 0 {
		cfg.WatchDebounce = 2 * time.Second
	}
	if cfg.ReleasesKeep <= 0 {
		cfg.ReleasesKeep = defaultReleaseKeep
	}
	if cfg.SuccessfulBuildsKeep <= 0 {
		cfg.SuccessfulBuildsKeep = 60
	}
	if cfg.FailedBuildsKeep <= 0 {
		cfg.FailedBuildsKeep = 100
	}
	if cfg.RefreshRunsKeep <= 0 {
		cfg.RefreshRunsKeep = 100
	}
	if cfg.BuildTimeout <= 0 {
		cfg.BuildTimeout = 2 * time.Hour
	}
	if cfg.Webhook.Enabled && strings.TrimSpace(cfg.Webhook.Secret) == "" {
		return nil, fmt.Errorf("webhook secret is required when webhooks are enabled")
	}
	if cfg.Webhook.Branch == "" {
		cfg.Webhook.Branch = "main"
	}
	for i := range cfg.RefreshTasks {
		if cfg.RefreshTasks[i].Name == "" {
			return nil, fmt.Errorf("refresh task name is required")
		}
		if len(cfg.RefreshTasks[i].Args) == 0 {
			return nil, fmt.Errorf("refresh task %q must define args", cfg.RefreshTasks[i].Name)
		}
		d, err := time.ParseDuration(cfg.RefreshTasks[i].Every)
		if err != nil {
			return nil, fmt.Errorf("refresh task %q invalid every duration: %w", cfg.RefreshTasks[i].Name, err)
		}
		cfg.RefreshTasks[i].interval = d
	}
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("resolve executable: %w", err)
	}
	s := &Service{
		cfg:                  cfg,
		executable:           execPath,
		statePath:            filepath.Join(cfg.HistoryDir, defaultStateName),
		logDir:               filepath.Join(cfg.HistoryDir, defaultLogDirName),
		overrideDir:          filepath.Join(cfg.HistoryDir, defaultOverrideName),
		leaderPath:           filepath.Join(cfg.HistoryDir, defaultLeaderName),
		queueCh:              make(chan queueRequest, 128),
		controlRuntime:       servecontrol.NewRuntime(),
		watchChanged:         make(map[string]struct{}),
		instanceID:           os.Getenv("POD_NAME"),
		trustedProxyPrefixes: trustedProxyPrefixes,
		theme:                loadUITheme(cfg),
	}
	if s.instanceID == "" {
		s.instanceID = hostSuffix()
	}
	instanceHost := os.Getenv("POD_IP")
	if instanceHost == "" {
		instanceHost = cfg.Host
	}
	s.instanceAddr = fmt.Sprintf("%s:%d", instanceHost, cfg.Port)
	if err := os.MkdirAll(s.logDir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}
	if err := os.MkdirAll(s.overrideDir, 0o755); err != nil {
		return nil, fmt.Errorf("create override dir: %w", err)
	}
	lockPath := filepath.Join(cfg.HistoryDir, defaultLockName)
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open leader lock: %w", err)
	}
	s.leaderLock = lockFile
	if err := s.loadState(); err != nil {
		_ = lockFile.Close()
		return nil, err
	}
	s.refreshControlStateWithFiles(s.state)
	return s, nil
}

func (s *Service) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	s.registerRoutes(mux)
	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
	}
	go s.runLeadershipLoop(ctx)
	go func() {
		<-ctx.Done()
		s.releaseLeadership()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.server.Shutdown(shutdownCtx)
	}()
	err := s.server.ListenAndServe()
	if s.leaderLock != nil {
		_ = s.leaderLock.Close()
	}
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Service) runLeadershipLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		s.tryBecomeLeader(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *Service) tryBecomeLeader(ctx context.Context) {
	if s.isLeader() || s.leaderLock == nil {
		return
	}
	if err := tryLockFile(s.leaderLock); err != nil {
		return
	}
	if err := s.writeLeaderRecord(); err != nil {
		_ = unlockFile(s.leaderLock)
		return
	}
	leaderCtx, cancel := context.WithCancel(ctx)
	s.leaderMu.Lock()
	s.leader = true
	s.leaderCancel = cancel
	s.leaderMu.Unlock()
	go s.worker(leaderCtx)
	s.resumeLeaderSession()
	for i := range s.cfg.RefreshTasks {
		go s.runRefreshScheduler(leaderCtx, s.cfg.RefreshTasks[i])
	}
	if s.cfg.WatchEnabled {
		go s.watchSource(leaderCtx)
	}
}

func (s *Service) releaseLeadership() {
	s.leaderMu.Lock()
	if !s.leader {
		s.leaderMu.Unlock()
		return
	}
	cancel := s.leaderCancel
	s.leader = false
	s.leaderCancel = nil
	s.leaderMu.Unlock()
	if cancel != nil {
		cancel()
	}
	s.pruneScheduleMu.Lock()
	pruneCancel := s.pruneCancel
	pruneDone := s.pruneDone
	s.pruneScheduleMu.Unlock()
	if pruneCancel != nil {
		pruneCancel()
	}
	if pruneDone != nil {
		<-pruneDone
	}
	if s.leaderLock != nil {
		_ = unlockFile(s.leaderLock)
	}
}

func (s *Service) isLeader() bool {
	s.leaderMu.RLock()
	defer s.leaderMu.RUnlock()
	return s.leader
}

func (s *Service) writeLeaderRecord() error {
	record := leaderRecord{
		InstanceID: s.instanceID,
		Addr:       s.instanceAddr,
		AcquiredAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.leaderPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.leaderPath)
}

func (s *Service) readLeaderRecord() (leaderRecord, error) {
	data, err := os.ReadFile(s.leaderPath)
	if err != nil {
		return leaderRecord{}, err
	}
	var record leaderRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return leaderRecord{}, err
	}
	return record, nil
}

func (s *Service) resumeLeaderSession() {
	queued := s.recoverQueuedState()
	for _, item := range queued {
		if req, ok := s.queueRequestFromQueued(item); ok {
			s.queueCh <- req
		}
	}
}

func (s *Service) recoverQueuedState() []QueuedOperation {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	queued := append([]QueuedOperation(nil), s.state.Queue...)
	if s.state.Running != nil {
		s.state.Running = nil
		s.saveStateLocked()
	}
	return queued
}

func (s *Service) queueRequestFromQueued(queued QueuedOperation) (queueRequest, bool) {
	req := queueRequest{QueuedOperation: queued}
	switch queued.Kind {
	case "build", "rollback":
		return req, true
	case "refresh":
		task, ok := s.findRefreshTask(queued.TaskName)
		if !ok {
			return queueRequest{}, false
		}
		req.commandArgs = append([]string(nil), task.Args...)
		return req, true
	default:
		return queueRequest{}, false
	}
}

func (s *Service) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/webhook", s.handleWebhook)
	api := http.NewServeMux()
	api.HandleFunc("/llms.txt", s.handleLLMsTXT)
	api.HandleFunc("/state", s.handleAPIState)
	api.HandleFunc("/builds", s.handleAPIBuilds)
	api.HandleFunc("/builds/", s.handleAPIBuildDetail)
	api.HandleFunc("/refresh/", s.handleAPIRefresh)
	api.HandleFunc("/refreshes/", s.handleAPIRefreshDetail)
	api.HandleFunc("/releases/", s.handleAPIReleaseAction)
	api.HandleFunc("/logs/", s.handleAPILogs)
	api.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		writeAPIError(w, http.StatusNotFound, "not_found", "API resource not found.")
	})
	mux.Handle("/api/v1", s.requireBearerOperator(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeAPIError(w, http.StatusNotFound, "not_found", "API resource not found.")
	})))
	mux.Handle("/api/v1/", s.requireBearerOperator(http.StripPrefix("/api/v1", api)))
	protected := http.NewServeMux()
	protected.HandleFunc("/llms.txt", s.handleLLMsTXT)
	protected.HandleFunc("/builds/", s.handleBuildDetail)
	protected.HandleFunc("/", s.handleIndex)
	protected.HandleFunc("/api/state", s.handleState)
	protected.HandleFunc("/api/builds", s.handleBuilds)
	protected.HandleFunc("/api/refresh/", s.handleRefreshRun)
	protected.HandleFunc("/api/releases/", s.handleReleaseAction)
	protected.HandleFunc("/logs/", s.handleLogs)
	mux.Handle("/", s.requireTrustedOperator(protected))
}

type pushWebhookPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (s *Service) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.Webhook.Enabled {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxWebhookBodyBytes))
	if err != nil {
		http.Error(w, "invalid webhook body", http.StatusRequestEntityTooLarge)
		return
	}
	if !validWebhookSignature(r, body, s.cfg.Webhook.Secret) {
		http.Error(w, "invalid webhook signature", http.StatusUnauthorized)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if s.handleStandbyMutation(w, r) {
		return
	}
	if webhookEvent(r) != "push" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var payload pushWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid push webhook payload", http.StatusBadRequest)
		return
	}
	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
	if branch != s.cfg.Webhook.Branch {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	deliveryID := r.Header.Get("X-GitHub-Delivery")
	if deliveryID == "" {
		deliveryID = r.Header.Get("X-Gitea-Delivery")
	}
	detail := fmt.Sprintf("%s push %s on %s", webhookProvider(r), payload.After, branch)
	if payload.Repository.FullName != "" {
		detail += " from " + payload.Repository.FullName
	}
	if deliveryID != "" {
		detail += " (delivery " + deliveryID + ")"
	}
	if err := s.enqueueWebhookBuild(detail); err != nil {
		if errors.Is(err, ErrBuildQueueFull) {
			w.Header().Set("Retry-After", "30")
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func validWebhookSignature(r *http.Request, body []byte, secret string) bool {
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		signature = r.Header.Get("X-Gitea-Signature")
	}
	signature = strings.TrimPrefix(signature, "sha256=")
	if signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	expected := fmt.Sprintf("%x", mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func webhookEvent(r *http.Request) string {
	event := r.Header.Get("X-GitHub-Event")
	if event == "" {
		event = r.Header.Get("X-Gitea-Event")
	}
	return strings.ToLower(event)
}

func webhookProvider(r *http.Request) string {
	if r.Header.Get("X-GitHub-Event") != "" {
		return "GitHub"
	}
	return "Forgejo"
}

func (s *Service) pullSource(ctx context.Context) (bool, error) {
	branch, err := gitBranch(s.cfg.SourceDir)
	if err != nil {
		return false, err
	}
	if branch != s.cfg.Webhook.Branch {
		return false, fmt.Errorf("checked-out branch %q does not match webhook branch %q", branch, s.cfg.Webhook.Branch)
	}
	before, err := gitHead(s.cfg.SourceDir)
	if err != nil {
		return false, err
	}
	cmd := s.gitCommand(ctx, "pull", "--ff-only", "origin", s.cfg.Webhook.Branch)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if output, err := cmd.CombinedOutput(); err != nil {
		return false, fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	after, err := gitHead(s.cfg.SourceDir)
	if err != nil {
		return false, err
	}
	return before != after, nil
}

func gitBranch(sourceDir string) (string, error) {
	output, err := gitCommandForSource(context.Background(), sourceDir, "branch", "--show-current").Output()
	if err != nil {
		return "", fmt.Errorf("read checked-out git branch: %w", err)
	}
	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "", fmt.Errorf("source directory has a detached HEAD")
	}
	return branch, nil
}

func gitHead(sourceDir string) (string, error) {
	return sourcegit.Head(context.Background(), sourceDir)
}

func (s *Service) gitCommand(ctx context.Context, args ...string) *exec.Cmd {
	return gitCommandForSource(ctx, s.cfg.SourceDir, args...)
}

func gitCommandForSource(ctx context.Context, sourceDir string, args ...string) *exec.Cmd {
	return sourcegit.Command(ctx, sourceDir, args...)
}

func (s *Service) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	state := s.viewState()
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"queue":  len(state.Queue),
		"leader": s.isLeader(),
	})
}

func (s *Service) stateResponse() adminStateResponse {
	state := s.viewState()
	return adminStateResponse{
		State:         state,
		Revision:      s.revisionStatus(state),
		Releases:      s.discoverReleases(),
		CurrentID:     s.currentReleaseID(),
		CurrentPath:   s.currentReleasePath(),
		Config:        s.cfg,
		RefreshTasks:  s.refreshTaskNames(),
		ControlCenter: s.controlRuntime.Snapshot(),
	}
}

func (s *Service) handleState(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.stateResponse())
}

func (s *Service) handleBuilds(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if s.handleStandbyMutation(w, r) {
			return
		}
		if !s.validateCSRF(w, r) {
			return
		}
		profile, _ := operatorProfileFromHeaders(r.Header, s.cfg.AuthHeaders)
		if _, err := s.enqueueBuildAs("manual-ui", "Manual build from admin UI", nil, profile.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Service) handleRefreshRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.handleStandbyMutation(w, r) {
		return
	}
	if !s.validateCSRF(w, r) {
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/api/refresh/")
	if name == "" {
		http.Error(w, "missing task name", http.StatusBadRequest)
		return
	}
	profile, _ := operatorProfileFromHeaders(r.Header, s.cfg.AuthHeaders)
	if _, err := s.enqueueRefreshAs(name, "manual-ui", "Manual refresh from admin UI", profile.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) handleReleaseAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.handleStandbyMutation(w, r) {
		return
	}
	if !s.validateCSRF(w, r) {
		return
	}
	trimmed := strings.TrimPrefix(r.URL.Path, "/api/releases/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) != 2 || parts[1] != "rollback" {
		http.Error(w, "unsupported release action", http.StatusBadRequest)
		return
	}
	profile, _ := operatorProfileFromHeaders(r.Header, s.cfg.AuthHeaders)
	if _, err := s.enqueueRollbackAs(parts[0], "manual-ui", "Rollback from admin UI", profile.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) handleLogs(w http.ResponseWriter, r *http.Request) {
	rel := strings.TrimPrefix(r.URL.Path, "/logs/")
	if rel == "" || strings.Contains(rel, "..") {
		http.NotFound(w, r)
		return
	}
	path := filepath.Join(s.logDir, rel)
	data, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(data)
}

func (s *Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("builder-admin").Funcs(template.FuncMap{
		"msToSeconds": func(ms int64) string {
			return fmt.Sprintf("%.2fs", float64(ms)/1000)
		},
		"since": func(t time.Time) string {
			return formatUITimestamp(t, time.Now().UTC())
		},
		"releaseLabel": func(createdAt time.Time) string {
			if createdAt.IsZero() {
				return "Release"
			}
			return createdAt.UTC().Format("02 Jan · 15:04 UTC")
		},
		"summaryPreview": func(lines []string) []string {
			if len(lines) <= 6 {
				return lines
			}
			return lines[len(lines)-6:]
		},
		"statusClass": uiStatusClass,
		"json": func(value any) string {
			data, err := json.Marshal(value)
			if err != nil {
				return "{}"
			}
			return string(data)
		},
		"browserTokens": func(theme uiTheme) template.CSS {
			return template.CSS(servecontrol.BrowserTokenStylesheet(browserThemeMap(theme)))
		},
		"queueWait": func(items []QueuedOperation) string {
			if len(items) == 0 {
				return "none"
			}
			oldest := items[0].EnqueuedAt
			for _, item := range items[1:] {
				if item.EnqueuedAt.Before(oldest) {
					oldest = item.EnqueuedAt
				}
			}
			return formatQueueWait(time.Since(oldest))
		},
	}).Parse(indexHTML))
	csrfToken, err := newCSRFToken()
	if err != nil {
		http.Error(w, "create CSRF token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, csrfCookie(csrfToken))
	state := s.viewState()
	operator := s.mustOperatorProfile(r.Header)
	data := struct {
		State         State
		Releases      []ReleaseView
		Revision      RevisionStatus
		CurrentID     string
		CurrentPath   string
		CSRFToken     string
		Operator      OperatorProfile
		PictureURL    string
		RefreshTasks  []RefreshTaskConfig
		CompletedJobs []completedJobView
		PreviewOrigin string
		Theme         uiTheme
		Capabilities  map[string]bool
	}{
		State:         state,
		Releases:      s.discoverReleases(),
		Revision:      s.revisionStatus(state),
		CurrentID:     s.currentReleaseID(),
		CurrentPath:   s.currentReleasePath(),
		CSRFToken:     csrfToken,
		Operator:      operator,
		PictureURL:    profilePictureURL(s.cfg.PublicAuthOrigin, operator.UserID),
		RefreshTasks:  s.cfg.RefreshTasks,
		CompletedJobs: completedJobs(state),
		PreviewOrigin: s.cfg.PreviewOrigin,
		Theme:         s.theme,
		Capabilities: map[string]bool{
			"enqueueBuild": true, "enqueueRefresh": len(s.cfg.RefreshTasks) > 0,
			"promoteRelease": true, "sourceMutation": false,
		},
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, data)
}

func (s *Service) handleBuildDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/builds/"), "/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	state := s.viewState()
	for i := range state.Builds {
		build := state.Builds[i]
		if build.ID != id {
			continue
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = template.Must(template.New("build-detail").Funcs(template.FuncMap{
			"msToSeconds": func(ms int64) string { return fmt.Sprintf("%.2fs", float64(ms)/1000) },
			"since":       func(t time.Time) string { return formatUITimestamp(t, time.Now().UTC()) },
			"statusClass": uiStatusClass,
			"browserTokens": func(theme uiTheme) template.CSS {
				return template.CSS(servecontrol.BrowserTokenStylesheet(browserThemeMap(theme)))
			},
		}).Parse(buildDetailHTML)).Execute(w, struct {
			BuildRecord
			PreviewURL string
			Theme      uiTheme
		}{BuildRecord: build, PreviewURL: strings.TrimSuffix(s.cfg.PreviewOrigin, "/") + "/__preview/" + build.ReleaseID + "/", Theme: s.theme})
		return
	}
	http.NotFound(w, r)
}

func completedJobs(state State) []completedJobView {
	jobs := make([]completedJobView, 0, len(state.Builds)+len(state.Refresh))
	for i := range state.Builds {
		build := &state.Builds[i]
		kind := "Build"
		if build.RollbackRelease != "" {
			kind = "Promote release"
		}
		release := ""
		if build.BecameLive {
			release = "Live"
		}
		jobs = append(jobs, completedJobView{ID: build.ID, Kind: kind, Status: build.Status, Trigger: build.TriggerType, FinishedAt: build.FinishedAt, QueueWaitMS: build.QueueWaitMS, RunMS: build.BuildMS, Release: release, LogPath: build.LogPath, Build: build})
	}
	for i := range state.Refresh {
		refresh := &state.Refresh[i]
		jobs = append(jobs, completedJobView{ID: refresh.ID, Kind: "Refresh: " + refresh.TaskName, Status: refresh.Status, Trigger: refresh.TriggerType, FinishedAt: refresh.FinishedAt, QueueWaitMS: refresh.QueueWaitMS, RunMS: refresh.RunMS, LogPath: refresh.LogPath})
	}
	sort.SliceStable(jobs, func(i, j int) bool { return jobs[i].FinishedAt.After(jobs[j].FinishedAt) })
	return jobs
}

func (s *Service) mustOperatorProfile(headers http.Header) OperatorProfile {
	profile, _ := operatorProfileFromHeaders(headers, s.cfg.AuthHeaders)
	return profile
}

func (s *Service) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-s.queueCh:
			s.process(ctx, req)
		}
	}
}

func (s *Service) process(ctx context.Context, req queueRequest) {
	s.setRunning(req)
	defer s.clearRunning()
	s.removeQueued(req.ID)
	switch req.Kind {
	case "build":
		if req.SyncSource {
			s.updateRunningPhase("sync source")
			syncCtx, cancel := context.WithTimeout(ctx, s.cfg.BuildTimeout)
			changed, err := s.pullSource(syncCtx)
			cancel()
			if err != nil {
				s.finishWebhookSyncFailure(req, err)
				return
			}
			if !changed {
				return
			}
		}
		s.runBuild(ctx, req)
	case "refresh":
		s.runRefresh(ctx, req)
	case "rollback":
		s.runRollback(req)
	}
}

func (s *Service) enqueueBuild(triggerType, detail string, changed []string) error {
	_, err := s.enqueueBuildAs(triggerType, detail, changed, "")
	return err
}

func (s *Service) enqueueBuildAs(triggerType, detail string, changed []string, actorID string) (QueuedOperation, error) {
	queued := QueuedOperation{
		ID:          nextID("build"),
		Kind:        "build",
		ActorID:     actorID,
		Label:       "Build",
		TriggerType: triggerType,
		Detail:      detail,
		Changed:     append([]string(nil), changed...),
		EnqueuedAt:  time.Now().UTC(),
	}
	s.pushQueued(queued)
	request := queueRequest{QueuedOperation: queued}
	select {
	case s.queueCh <- request:
		return queued, nil
	default:
		s.removeQueued(queued.ID)
		return QueuedOperation{}, ErrBuildQueueFull
	}
}

func (s *Service) enqueueWebhookBuild(detail string) error {
	queued := QueuedOperation{
		ID:          nextID("build"),
		Kind:        "build",
		Label:       "Build",
		TriggerType: "webhook",
		Detail:      detail,
		EnqueuedAt:  time.Now().UTC(),
		SyncSource:  true,
	}
	s.pushQueued(queued)
	request := queueRequest{QueuedOperation: queued}
	select {
	case s.queueCh <- request:
		return nil
	default:
		s.removeQueued(queued.ID)
		return ErrBuildQueueFull
	}
}

func (s *Service) finishWebhookSyncFailure(req queueRequest, err error) {
	finished := time.Now().UTC()
	s.finishBuild(BuildRecord{
		ID:            req.ID,
		Kind:          req.Kind,
		Status:        "failed",
		TriggerType:   req.TriggerType,
		TriggerDetail: req.Detail,
		EnqueuedAt:    req.EnqueuedAt,
		StartedAt:     req.EnqueuedAt,
		FinishedAt:    finished,
		TotalMS:       finished.Sub(req.EnqueuedAt).Milliseconds(),
		Error:         fmt.Sprintf("git pull failed: %v", err),
	})
}

func (s *Service) enqueueRefresh(name, triggerType, detail string) error {
	_, err := s.enqueueRefreshAs(name, triggerType, detail, "")
	return err
}

func (s *Service) enqueueRefreshAs(name, triggerType, detail, actorID string) (QueuedOperation, error) {
	task, ok := s.findRefreshTask(name)
	if !ok {
		return QueuedOperation{}, fmt.Errorf("unknown refresh task %q", name)
	}
	queued := QueuedOperation{
		ID:          nextID("refresh"),
		Kind:        "refresh",
		ActorID:     actorID,
		Label:       "Refresh " + task.Name,
		TriggerType: triggerType,
		Detail:      detail,
		EnqueuedAt:  time.Now().UTC(),
		TaskName:    task.Name,
	}
	s.pushQueued(queued)
	select {
	case s.queueCh <- queueRequest{QueuedOperation: queued, commandArgs: append([]string(nil), task.Args...)}:
		return queued, nil
	default:
		s.removeQueued(queued.ID)
		return QueuedOperation{}, ErrBuildQueueFull
	}
}

func (s *Service) enqueueRollbackAs(releaseID, triggerType, detail, actorID string) (QueuedOperation, error) {
	if releaseID == "" {
		return QueuedOperation{}, fmt.Errorf("release id is required")
	}
	queued := QueuedOperation{
		ID:          nextID("rollback"),
		Kind:        "rollback",
		ActorID:     actorID,
		Label:       "Rollback " + releaseID,
		TriggerType: triggerType,
		Detail:      detail,
		EnqueuedAt:  time.Now().UTC(),
		ReleaseID:   releaseID,
	}
	s.pushQueued(queued)
	select {
	case s.queueCh <- queueRequest{QueuedOperation: queued}:
		return queued, nil
	default:
		s.removeQueued(queued.ID)
		return QueuedOperation{}, ErrBuildQueueFull
	}
}

func (s *Service) runBuild(ctx context.Context, req queueRequest) {
	started := time.Now().UTC()
	record := BuildRecord{
		ID:             req.ID,
		Kind:           req.Kind,
		ActorID:        req.ActorID,
		Status:         "running",
		BuilderVersion: s.cfg.BuilderVersion,
		BuilderCommit:  reliableBuilderRevision(s.cfg.BuilderCommit),
		TriggerType:    req.TriggerType,
		TriggerDetail:  req.Detail,
		ChangedPaths:   append([]string(nil), req.Changed...),
		Impact:         classifyBuildImpact(req.Changed),
		EnqueuedAt:     req.EnqueuedAt,
		StartedAt:      started,
		QueueWaitMS:    started.Sub(req.EnqueuedAt).Milliseconds(),
	}
	logPath, logFile, err := s.createLogFile(req.ID)
	if err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		s.finishBuild(record)
		return
	}
	defer logFile.Close()
	record.LogPath = logPath
	leaderCtx := ctx
	ctx, cancel := context.WithTimeout(ctx, s.cfg.BuildTimeout)
	defer cancel()

	phaseStart := time.Now()
	s.updateRunningPhase("prepare")
	if err := s.prepareBuild(logFile); err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.PrepareMS = time.Since(phaseStart).Milliseconds()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		record.PerfSummary = extractPerfSummaryFromFile(filepath.Join(s.logDir, logPath))
		s.finishBuild(record)
		return
	}
	record.PrepareMS = time.Since(phaseStart).Milliseconds()
	if sourceCommit, gitErr := sourcegit.Head(ctx, s.cfg.SourceDir); gitErr == nil {
		record.SourceCommit = reliableBuilderRevision(sourceCommit)
	}

	buildWork := s.cfg.WorkDir
	phaseStart = time.Now()
	s.updateRunningPhase("build")
	cmdArgs, cleanup, err := s.buildCommandArgs(req.ID, buildWork)
	if err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.BuildMS = time.Since(phaseStart).Milliseconds()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		record.PerfSummary = extractPerfSummaryFromFile(filepath.Join(s.logDir, logPath))
		s.finishBuild(record)
		return
	}
	defer cleanup()
	if err := s.runLoggedCommand(ctx, logFile, s.cfg.SourceDir, nil, cmdArgs...); err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.BuildMS = time.Since(phaseStart).Milliseconds()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		record.PerfSummary = extractPerfSummaryFromFile(filepath.Join(s.logDir, logPath))
		s.finishBuild(record)
		return
	}
	record.BuildMS = time.Since(phaseStart).Milliseconds()

	phaseStart = time.Now()
	s.updateRunningPhase("promote")
	releaseID, releasePath, err := s.promoteBuildLogged(buildWork, logFile)
	if err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.PromoteMS = time.Since(phaseStart).Milliseconds()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		record.PerfSummary = extractPerfSummaryFromFile(filepath.Join(s.logDir, logPath))
		s.finishBuild(record)
		return
	}
	record.PromoteMS = time.Since(phaseStart).Milliseconds()
	record.ReleaseID = releaseID
	record.ReleasePath = releasePath
	record.BecameLive = true
	if err := markReusableWorkspace(buildWork, releaseID); err != nil {
		_, _ = fmt.Fprintln(logFile, "workspace reuse unavailable:", err)
	}

	record.Status = "success"
	record.FinishedAt = time.Now().UTC()
	record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
	record.PerfSummary = extractPerfSummaryFromFile(filepath.Join(s.logDir, logPath))
	s.finishBuild(record)
	s.schedulePrune(leaderCtx, record.ID)
}

func (s *Service) runRefresh(ctx context.Context, req queueRequest) {
	started := time.Now().UTC()
	record := RefreshRecord{
		ID:                    req.ID,
		ActorID:               req.ActorID,
		TaskName:              req.TaskName,
		Status:                "running",
		TriggerType:           req.TriggerType,
		TriggerDetail:         req.Detail,
		EnqueuedAt:            req.EnqueuedAt,
		StartedAt:             started,
		QueueWaitMS:           started.Sub(req.EnqueuedAt).Milliseconds(),
		EnqueueBuildOnSuccess: false,
		Command:               append([]string(nil), req.commandArgs...),
	}
	task, ok := s.findRefreshTask(req.TaskName)
	if ok {
		record.EnqueueBuildOnSuccess = task.EnqueueBuildOnSuccess
	}
	logPath, logFile, err := s.createLogFile(req.ID)
	if err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		s.finishRefresh(record)
		return
	}
	defer logFile.Close()
	record.LogPath = logPath
	ctx, cancel := context.WithTimeout(ctx, s.cfg.BuildTimeout)
	defer cancel()
	s.updateRunningPhase("refresh")
	runStart := time.Now()
	if err := s.runLoggedCommand(ctx, logFile, s.cfg.SourceDir, nil, req.commandArgs...); err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.RunMS = time.Since(runStart).Milliseconds()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		s.finishRefresh(record)
		return
	}
	record.RunMS = time.Since(runStart).Milliseconds()
	record.Status = "success"
	if task.EnqueueBuildOnSuccess {
		buildID := nextID("build")
		queued := QueuedOperation{
			ID:          buildID,
			Kind:        "build",
			ActorID:     req.ActorID,
			Label:       "Build",
			TriggerType: "scheduled-refresh",
			Detail:      "Build enqueued by refresh task " + task.Name,
			EnqueuedAt:  time.Now().UTC(),
		}
		record.EnqueuedBuildID = buildID
		s.pushQueued(queued)
		s.queueCh <- queueRequest{QueuedOperation: queued}
	}
	record.FinishedAt = time.Now().UTC()
	record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
	s.finishRefresh(record)
}

func (s *Service) runRollback(req queueRequest) {
	started := time.Now().UTC()
	record := BuildRecord{
		ID:              req.ID,
		Kind:            req.Kind,
		ActorID:         req.ActorID,
		Status:          "running",
		TriggerType:     req.TriggerType,
		TriggerDetail:   req.Detail,
		EnqueuedAt:      req.EnqueuedAt,
		StartedAt:       started,
		QueueWaitMS:     started.Sub(req.EnqueuedAt).Milliseconds(),
		RollbackRelease: req.ReleaseID,
	}
	logPath, logFile, err := s.createLogFile(req.ID)
	if err != nil {
		record.Status = "failed"
		record.Error = err.Error()
		record.FinishedAt = time.Now().UTC()
		record.TotalMS = record.FinishedAt.Sub(started).Milliseconds()
		s.finishBuild(record)
		return
	}
	defer logFile.Close()
	record.LogPath = logPath
	s.updateRunningPhase("promote")
	s.releaseMu.Lock()
	defer s.releaseMu.Unlock()
	releasePath := filepath.Join(s.cfg.SiteDir, "releases", req.ReleaseID)
	if _, err := os.Stat(releasePath); err != nil {
		record.Status = "failed"
		record.Error = fmt.Sprintf("release %q not found: %v", req.ReleaseID, err)
	} else if err := s.switchCurrentRelease(req.ReleaseID); err != nil {
		record.Status = "failed"
		record.Error = err.Error()
	} else {
		_, _ = fmt.Fprintf(logFile, "promoted release %s\n", req.ReleaseID)
		record.Status = "success"
		record.ReleaseID = req.ReleaseID
		record.ReleasePath = releasePath
		record.BecameLive = true
	}
	record.FinishedAt = time.Now().UTC()
	record.PromoteMS = record.FinishedAt.Sub(started).Milliseconds()
	record.TotalMS = record.PromoteMS
	record.PerfSummary = extractPerfSummaryFromFile(filepath.Join(s.logDir, logPath))
	s.finishBuild(record)
}

func (s *Service) prepareBuild(log io.Writer) error {
	if err := os.MkdirAll(filepath.Join(s.cfg.SiteDir, "releases"), 0o755); err != nil {
		return err
	}
	if s.cfg.CacheMount != "" {
		for _, part := range []string{"build", "plugin", "xdg"} {
			if err := os.MkdirAll(filepath.Join(s.cfg.CacheMount, part), 0o755); err != nil {
				return err
			}
		}
		for _, linkName := range []string{".markata", ".markata-cache"} {
			_ = os.RemoveAll(filepath.Join(s.cfg.SourceDir, linkName))
		}
		if err := os.Symlink(filepath.Join(s.cfg.CacheMount, "build"), filepath.Join(s.cfg.SourceDir, ".markata")); err != nil {
			return err
		}
		if err := os.Symlink(filepath.Join(s.cfg.CacheMount, "plugin"), filepath.Join(s.cfg.SourceDir, ".markata-cache")); err != nil {
			return err
		}
	}
	buildWork := s.cfg.WorkDir
	reusable, err := claimReusableWorkspace(buildWork, s.currentReleaseID())
	if err != nil {
		return err
	}
	if reusable {
		_, _ = fmt.Fprintln(log, "reusing build work from current release")
		return nil
	}
	if err := invalidateReusableWorkspace(buildWork); err != nil {
		return err
	}
	if err := os.RemoveAll(buildWork); err != nil {
		return err
	}
	if err := os.MkdirAll(buildWork, 0o755); err != nil {
		return err
	}
	current := filepath.Join(s.cfg.SiteDir, "current")
	if _, err := os.Stat(current); err == nil {
		_, _ = fmt.Fprintln(log, "seeding build work from current release")
		// Copy file contents instead of creating hard links. The build mutates
		// files in .build-work, and hard links would mutate the historical
		// release that is still available for rollback.
		return s.runLoggedCommand(context.Background(), log, "", nil, "cp", "-a", current+"/.", buildWork+string(os.PathSeparator))
	}
	return nil
}

func (s *Service) buildCommandArgs(id, buildWork string) ([]string, func(), error) {
	args := make([]string, 0, 10)
	if s.cfg.ConfigPath != "" {
		args = append(args, "--config", s.cfg.ConfigPath)
	}
	cleanup := func() {}
	if s.cfg.MermaidMode != "" {
		overridePath := filepath.Join(s.overrideDir, "builder-admin.toml")
		contents := fmt.Sprintf("[markata-go.mermaid]\nmode = %q\n", s.cfg.MermaidMode)
		if err := os.WriteFile(overridePath, []byte(contents), 0o644); err != nil {
			return nil, cleanup, err
		}
		args = append(args, "-m", overridePath)
	}
	args = append(args, "build")
	if s.cfg.Fast {
		args = append(args, "--fast")
	}
	args = append(args, "--output", buildWork)
	return args, cleanup, nil
}

func (s *Service) promoteBuild(buildWork string) (string, string, error) {
	return s.promoteBuildLogged(buildWork, io.Discard)
}

func (s *Service) promoteBuildLogged(buildWork string, log io.Writer) (string, string, error) {
	s.releaseMu.Lock()
	defer s.releaseMu.Unlock()
	releaseID := time.Now().UTC().Format("20060102T150405Z") + "-" + hostSuffix()
	releasesDir := filepath.Join(s.cfg.SiteDir, "releases")
	releasePath, stats, err := stageIncrementalWorkspace(buildWork, releasesDir, releaseID, s.currentReleaseID())
	_, _ = fmt.Fprintf(log, "incremental publication: linked_files=%d copied_files=%d copied_bytes=%d compared_bytes=%d cached_source_hashes=%d cached_release_hashes=%d manifest_saved=%t\n", stats.LinkedFiles, stats.CopiedFiles, stats.CopiedBytes, stats.ComparedBytes, stats.CachedSourceHashes, stats.CachedReleaseHashes, stats.ManifestSaved)
	if err != nil {
		return "", "", err
	}
	if err := s.switchCurrentRelease(releaseID); err != nil {
		return "", "", err
	}
	return releaseID, releasePath, nil
}

func (s *Service) switchCurrentRelease(releaseID string) error {
	if strings.HasPrefix(releaseID, ".") {
		return fmt.Errorf("internal release directories cannot be promoted: %s", releaseID)
	}
	currentNext := filepath.Join(s.cfg.SiteDir, "current.next")
	_ = os.Remove(currentNext)
	if err := os.Symlink(filepath.Join("releases", releaseID), currentNext); err != nil {
		return err
	}
	return replaceCurrentRelease(currentNext, filepath.Join(s.cfg.SiteDir, "current"))
}

func (s *Service) schedulePrune(ctx context.Context, buildID string) {
	s.pruneScheduleMu.Lock()
	if s.pruneRunning {
		s.prunePending = append(s.prunePending, pruneRequest{ctx: ctx, buildID: buildID})
		s.pruneScheduleMu.Unlock()
		return
	}
	pruneCtx, cancel := context.WithCancel(ctx)
	s.pruneRunning = true
	s.pruneCancel = cancel
	s.pruneDone = make(chan struct{})
	done := s.pruneDone
	s.pruneScheduleMu.Unlock()
	go s.runScheduledPrunes(pruneCtx, buildID, done)
}

func (s *Service) runScheduledPrunes(ctx context.Context, buildID string, done chan struct{}) {
	defer func() {
		close(done)
		s.pruneScheduleMu.Lock()
		if s.pruneDone == done {
			s.pruneRunning = false
			s.pruneCancel = nil
			s.pruneDone = nil
		}
		s.pruneScheduleMu.Unlock()
	}()
	for {
		if ctx.Err() != nil {
			return
		}
		started := time.Now()
		if err := s.pruneReleases(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "builder-admin: release pruning failed:", err)
		}
		pruneMS := time.Since(started).Milliseconds()
		s.stateMu.Lock()
		for i := range s.state.Builds {
			if s.state.Builds[i].ID == buildID {
				s.state.Builds[i].PruneMS = pruneMS
				s.saveStateLocked()
				break
			}
		}
		s.stateMu.Unlock()

		s.pruneScheduleMu.Lock()
		if len(s.prunePending) == 0 {
			s.pruneRunning = false
			s.pruneCancel = nil
			s.pruneDone = nil
			s.pruneScheduleMu.Unlock()
			return
		}
		next := s.prunePending[0]
		s.prunePending = s.prunePending[1:]
		s.pruneScheduleMu.Unlock()
		if next.ctx == nil || next.ctx.Err() != nil {
			continue
		}
		buildID = next.buildID
		ctx = next.ctx
	}
}

func (s *Service) runLoggedCommand(ctx context.Context, log io.Writer, cwd string, env []string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}
	cmdName := args[0]
	cmdArgs := args[1:]
	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	if strings.HasPrefix(cmdName, "-") || cmdName == "build" || strings.HasSuffix(cmdName, "markata-go") || filepath.Base(cmdName) == filepath.Base(s.executable) {
		if strings.HasSuffix(cmdName, "markata-go") || filepath.Base(cmdName) == filepath.Base(s.executable) {
			cmd = exec.CommandContext(ctx, s.executable, cmdArgs...)
		} else {
			cmd = exec.CommandContext(ctx, s.executable, args...)
		}
	}
	cmd.Stdout = log
	cmd.Stderr = log
	if cwd != "" {
		cmd.Dir = cwd
	}
	cmd.Env = os.Environ()
	if len(env) > 0 {
		cmd.Env = append(cmd.Env, env...)
	}
	_, _ = fmt.Fprintf(log, "$ %s\n", strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err == nil {
		return nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if code, ok := processExitCode(exitErr); ok {
			return fmt.Errorf("command failed with exit code %d", code)
		}
	}
	return err
}

func (s *Service) handleStandbyMutation(w http.ResponseWriter, r *http.Request) bool {
	if s.isLeader() {
		return false
	}
	record, err := s.readLeaderRecord()
	if err != nil || record.Addr == "" || record.InstanceID == s.instanceID {
		http.Error(w, "builder-admin standby is waiting for the active leader", http.StatusServiceUnavailable)
		return true
	}
	target, err := url.Parse("http://" + record.Addr)
	if err != nil {
		http.Error(w, "builder-admin leader address is unavailable", http.StatusServiceUnavailable)
		return true
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.Header.Set(forwardedLeaderHeader, "1")
	}
	proxy.ErrorHandler = func(rw http.ResponseWriter, _ *http.Request, proxyErr error) {
		http.Error(rw, fmt.Sprintf("builder-admin leader proxy failed: %v", proxyErr), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
	return true
}

func (s *Service) runRefreshScheduler(ctx context.Context, task RefreshTaskConfig) {
	ticker := time.NewTicker(task.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = s.enqueueRefresh(task.Name, "schedule", "Scheduled refresh")
		}
	}
}

func (s *Service) watchSource(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	_ = addDirRecursive(watcher, s.cfg.SourceDir)
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if ignoreWatchPath(s.cfg.SourceDir, event.Name) {
				continue
			}
			if event.Op&fsnotify.Create != 0 {
				_ = addDirRecursiveIfDir(watcher, event.Name)
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}
			s.recordWatchPath(event.Name)
		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func (s *Service) recordWatchPath(path string) {
	rel, err := filepath.Rel(s.cfg.SourceDir, path)
	if err != nil {
		rel = path
	}
	s.watchMu.Lock()
	defer s.watchMu.Unlock()
	s.watchChanged[filepath.ToSlash(rel)] = struct{}{}
	if s.watchTimer != nil {
		if !s.watchTimer.Stop() {
			select {
			case <-s.watchTimer.C:
			default:
			}
		}
	}
	s.watchTimer = time.AfterFunc(s.cfg.WatchDebounce, func() {
		s.flushWatchBuild()
	})
}

func (s *Service) flushWatchBuild() {
	s.watchMu.Lock()
	changed := make([]string, 0, len(s.watchChanged))
	for path := range s.watchChanged {
		changed = append(changed, path)
	}
	clear(s.watchChanged)
	s.watchMu.Unlock()
	if len(changed) == 0 {
		return
	}
	sort.Strings(changed)
	_ = s.enqueueBuild("file-watch", fmt.Sprintf("Debounced file-watch build (%d paths)", len(changed)), changed)
}

func (s *Service) findRefreshTask(name string) (RefreshTaskConfig, bool) {
	for _, task := range s.cfg.RefreshTasks {
		if task.Name == name {
			return task, true
		}
	}
	return RefreshTaskConfig{}, false
}

func (s *Service) refreshTaskNames() []string {
	names := make([]string, 0, len(s.cfg.RefreshTasks))
	for _, task := range s.cfg.RefreshTasks {
		names = append(names, task.Name)
	}
	return names
}

func (s *Service) createLogFile(id string) (string, *os.File, error) {
	rel := id + ".log"
	path := filepath.Join(s.logDir, rel)
	f, err := os.Create(path)
	if err != nil {
		return "", nil, err
	}
	return rel, f, nil
}

func (s *Service) snapshotState() State {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	clone := s.state
	clone.Queue = append([]QueuedOperation(nil), s.state.Queue...)
	clone.Builds = append([]BuildRecord(nil), s.state.Builds...)
	clone.Refresh = append([]RefreshRecord(nil), s.state.Refresh...)
	return clone
}

func (s *Service) readPersistedState() (State, error) {
	data, err := os.ReadFile(s.statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, err
	}
	return state, nil
}

func (s *Service) viewState() State {
	if s.isLeader() {
		return s.snapshotState()
	}
	state, err := s.readPersistedState()
	if err == nil {
		return state
	}
	return s.snapshotState()
}

func (s *Service) pushQueued(queued QueuedOperation) {
	s.stateMu.Lock()
	s.state.Queue = append(s.state.Queue, queued)
	s.saveStateLocked()
	s.stateMu.Unlock()
}

func (s *Service) removeQueued(id string) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	filtered := s.state.Queue[:0]
	for _, queued := range s.state.Queue {
		if queued.ID != id {
			filtered = append(filtered, queued)
		}
	}
	s.state.Queue = append([]QueuedOperation(nil), filtered...)
	s.saveStateLocked()
}

func (s *Service) setRunning(req queueRequest) {
	s.stateMu.Lock()
	s.state.Running = &RunningOperation{
		ID:          req.ID,
		Kind:        req.Kind,
		ActorID:     req.ActorID,
		Label:       req.Label,
		TriggerType: req.TriggerType,
		Detail:      req.Detail,
		StartedAt:   time.Now().UTC(),
		EnqueuedAt:  req.EnqueuedAt,
		Phase:       "starting",
		Impact:      classifyBuildImpact(req.Changed),
	}
	s.saveStateLocked()
	s.stateMu.Unlock()
}

func classifyBuildImpact(paths []string) string {
	if len(paths) == 0 {
		return "unknown"
	}
	impact := "content"
	for _, path := range paths {
		path = strings.ToLower(filepath.ToSlash(path))
		switch {
		case strings.Contains(path, "markata-go.") || strings.HasPrefix(path, "pkg/") || strings.HasPrefix(path, "plugins/"):
			return "config or plugin"
		case strings.HasPrefix(path, "templates/") || strings.HasPrefix(path, "themes/") || strings.HasPrefix(path, "palettes/") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js"):
			impact = "template or asset"
		}
	}
	return impact
}

func (s *Service) updateRunningPhase(phase string) {
	s.stateMu.Lock()
	if s.state.Running != nil {
		s.state.Running.Phase = phase
		s.saveStateLocked()
	}
	s.stateMu.Unlock()
}

func (s *Service) clearRunning() {
	s.stateMu.Lock()
	s.state.Running = nil
	s.saveStateLocked()
	s.stateMu.Unlock()
}

func (s *Service) finishBuild(record BuildRecord) {
	record.LogDiagnostics = s.readControlLogDiagnostics(record.LogPath, record.ID)
	s.stateMu.Lock()
	s.state.Builds = append([]BuildRecord{record}, s.state.Builds...)
	s.pruneBuildHistoryLocked()
	s.saveStateLocked()
	s.stateMu.Unlock()
}

func (s *Service) finishRefresh(record RefreshRecord) {
	record.LogDiagnostics = s.readControlLogDiagnostics(record.LogPath, record.ID)
	s.stateMu.Lock()
	s.state.Refresh = append([]RefreshRecord{record}, s.state.Refresh...)
	s.pruneRefreshHistoryLocked()
	s.saveStateLocked()
	s.stateMu.Unlock()
}

func (s *Service) pruneBuildHistoryLocked() {
	kept := make([]BuildRecord, 0, len(s.state.Builds))
	successCount := 0
	failureCount := 0
	for _, record := range s.state.Builds {
		keep := false
		switch record.Status {
		case "success":
			if successCount < s.cfg.SuccessfulBuildsKeep {
				keep = true
				successCount++
			}
		default:
			if failureCount < s.cfg.FailedBuildsKeep {
				keep = true
				failureCount++
			}
		}
		if keep {
			kept = append(kept, record)
			continue
		}
		if record.LogPath != "" {
			_ = os.Remove(filepath.Join(s.logDir, record.LogPath))
		}
	}
	s.state.Builds = kept
}

func (s *Service) pruneRefreshHistoryLocked() {
	if len(s.state.Refresh) <= s.cfg.RefreshRunsKeep {
		return
	}
	for _, record := range s.state.Refresh[s.cfg.RefreshRunsKeep:] {
		if record.LogPath != "" {
			_ = os.Remove(filepath.Join(s.logDir, record.LogPath))
		}
	}
	s.state.Refresh = append([]RefreshRecord(nil), s.state.Refresh[:s.cfg.RefreshRunsKeep]...)
}

func (s *Service) loadState() error {
	state, err := s.readPersistedState()
	if err != nil {
		return fmt.Errorf("read state: %w", err)
	}
	s.state = state
	return nil
}

func (s *Service) saveStateLocked() {
	// Refresh the richer projection when Builder Admin state changes. The state
	// endpoint is polled frequently, so it must only read the already-built
	// in-memory snapshot instead of reopening job logs and diagnostics artifacts.
	defer s.refreshControlStateWithFiles(s.state)
	data, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return
	}
	tmp := s.statePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return
	}
	_ = os.Rename(tmp, s.statePath)
}

func (s *Service) refreshControlStateWithFiles(state State) {
	if s.controlRuntime == nil {
		return
	}
	snapshot := s.controlSnapshot(state)
	s.addRecentControlLogs(&snapshot, state)
	s.addCurrentReleaseDiagnostics(&snapshot, state)
	s.controlRuntime.ReplaceSnapshot(snapshot)
}

func (s *Service) controlSnapshot(state State) servecontrol.Snapshot {
	snapshot := projectControlState(state)
	snapshot.Server = servecontrol.ServerState{
		Status:  servecontrol.StateSuccess,
		Address: fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
	}
	return snapshot
}

func projectControlState(state State) servecontrol.Snapshot {
	snapshot := servecontrol.Snapshot{
		Jobs:        make([]servecontrol.Job, 0, len(state.Queue)+len(state.Builds)+len(state.Refresh)+1),
		Diagnostics: []servecontrol.Diagnostic{},
		Pages:       []servecontrol.Page{},
		Logs:        []servecontrol.LogEntry{},
	}
	queuedByID := make(map[string]QueuedOperation, len(state.Queue))
	for _, queued := range state.Queue {
		queuedByID[queued.ID] = queued
		if state.Running != nil && state.Running.ID == queued.ID {
			continue
		}
		snapshot.Jobs = append(snapshot.Jobs, servecontrol.Job{
			ID: queued.ID, Name: queued.Label, Type: queued.Kind, Trigger: queued.TriggerType,
			State: servecontrol.StateQueued, QueuedAt: queued.EnqueuedAt,
			Pages: append([]string(nil), queued.Changed...),
		})
	}
	if running := state.Running; running != nil {
		queuedAt := running.EnqueuedAt
		pages := []string(nil)
		if queued, ok := queuedByID[running.ID]; ok {
			if queuedAt.IsZero() {
				queuedAt = queued.EnqueuedAt
			}
			pages = append(pages, queued.Changed...)
		}
		if queuedAt.IsZero() {
			queuedAt = running.StartedAt
		}
		job := servecontrol.Job{
			ID: running.ID, Name: running.Label, Type: running.Kind, Trigger: running.TriggerType,
			State: servecontrol.StateRunning, QueuedAt: queuedAt, StartedAt: running.StartedAt,
			Pages: pages,
		}
		if running.Phase != "" {
			job.Steps = []servecontrol.Step{{ID: running.ID + "/" + running.Phase, Name: running.Phase, State: servecontrol.StateRunning, StartedAt: running.StartedAt}}
		}
		snapshot.Jobs = append(snapshot.Jobs, job)
	}
	for _, record := range state.Builds {
		job := servecontrol.Job{
			ID: record.ID, Name: "Build", Type: record.Kind, Trigger: record.TriggerType,
			State: controlState(record.Status), QueuedAt: record.EnqueuedAt,
			StartedAt: record.StartedAt, EndedAt: record.FinishedAt,
			Pages: append([]string(nil), record.ChangedPaths...),
			Steps: buildControlSteps(record),
		}
		if record.Kind == "rollback" {
			job.Name = "Rollback " + record.RollbackRelease
		}
		if record.Error != "" {
			diagnostic := servecontrol.Diagnostic{
				Code: "MARKATA-E901", Severity: "error", Message: record.Error, JobID: record.ID,
				SuggestedFix: "Open this job's log, correct the failed build command, then trigger another build.",
			}
			job.Diagnostics = append(job.Diagnostics, diagnostic)
			snapshot.Diagnostics = append(snapshot.Diagnostics, diagnostic)
		}
		job.Diagnostics = append(job.Diagnostics, record.LogDiagnostics...)
		snapshot.Diagnostics = append(snapshot.Diagnostics, record.LogDiagnostics...)
		snapshot.Jobs = append(snapshot.Jobs, job)
	}
	for _, record := range state.Refresh {
		job := servecontrol.Job{
			ID: record.ID, Name: "Refresh " + record.TaskName, Type: "refresh", Trigger: record.TriggerType,
			State: controlState(record.Status), QueuedAt: record.EnqueuedAt,
			StartedAt: record.StartedAt, EndedAt: record.FinishedAt,
			Steps: []servecontrol.Step{{ID: record.ID + "/refresh", Name: "refresh", State: controlState(record.Status), StartedAt: record.StartedAt, EndedAt: record.FinishedAt}},
		}
		if record.Error != "" {
			diagnostic := servecontrol.Diagnostic{
				Code: "MARKATA-E902", Severity: "error", Message: record.Error, JobID: record.ID,
				SuggestedFix: "Open this refresh job's log, correct the command failure, then run the refresh again.",
			}
			job.Diagnostics = append(job.Diagnostics, diagnostic)
			snapshot.Diagnostics = append(snapshot.Diagnostics, diagnostic)
		}
		job.Diagnostics = append(job.Diagnostics, record.LogDiagnostics...)
		snapshot.Diagnostics = append(snapshot.Diagnostics, record.LogDiagnostics...)
		snapshot.Jobs = append(snapshot.Jobs, job)
	}
	sort.SliceStable(snapshot.Jobs, func(i, j int) bool {
		if snapshot.Jobs[i].QueuedAt.Equal(snapshot.Jobs[j].QueuedAt) {
			return snapshot.Jobs[i].ID < snapshot.Jobs[j].ID
		}
		return snapshot.Jobs[i].QueuedAt.Before(snapshot.Jobs[j].QueuedAt)
	})
	return snapshot
}

func controlState(status string) servecontrol.State {
	switch status {
	case "queued":
		return servecontrol.StateQueued
	case "running":
		return servecontrol.StateRunning
	case "success":
		return servecontrol.StateSuccess
	case "warning":
		return servecontrol.StateWarning
	//nolint:misspell // Accepts the legacy serialized spelling for compatibility.
	case "cancelled":
		return servecontrol.StateCancelled
	default:
		return servecontrol.StateFailed
	}
}

func buildControlSteps(record BuildRecord) []servecontrol.Step {
	phases := []struct {
		name string
		ms   int64
	}{{"prepare", record.PrepareMS}, {"build", record.BuildMS}, {"promote", record.PromoteMS}, {"prune", record.PruneMS}}
	steps := make([]servecontrol.Step, 0, len(phases))
	started := record.StartedAt
	for _, phase := range phases {
		if phase.ms == 0 {
			continue
		}
		ended := started.Add(time.Duration(phase.ms) * time.Millisecond)
		steps = append(steps, servecontrol.Step{
			ID: record.ID + "/" + phase.name, Name: phase.name, State: servecontrol.StateSuccess,
			StartedAt: started, EndedAt: ended,
		})
		started = ended
	}
	if record.Status == "failed" && len(steps) > 0 {
		steps[len(steps)-1].State = servecontrol.StateFailed
	}
	return steps
}

// readControlLogDiagnostics scans the completed file once so warnings and
// errors survive later log output and subsequent builds in persisted history.
func (s *Service) readControlLogDiagnostics(logPath, jobID string) []servecontrol.Diagnostic {
	if logPath == "" || filepath.Base(logPath) != logPath {
		return nil
	}
	file, err := os.Open(filepath.Join(s.logDir, logPath))
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64<<10), 16<<20)
	var found []servecontrol.Diagnostic
	seen := make(map[string]struct{})
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		severity := controlLogSeverity(line)
		if severity == "" {
			continue
		}
		key := severity + "\x00" + line
		if _, duplicate := seen[key]; duplicate {
			continue
		}
		seen[key] = struct{}{}
		code := "MARKATA-W900"
		if severity == "error" {
			code = "MARKATA-E900"
		}
		found = append(found, servecontrol.Diagnostic{
			Code: code, Severity: severity, Message: line, JobID: jobID,
			Explanation:  "This message came from the builder command output.",
			SuggestedFix: "Review the full job log and the source file named in this message, then rerun the job.",
		})
	}
	return found
}

func controlLogSeverity(line string) string {
	lower := strings.ToLower(strings.TrimSpace(line))
	switch {
	case strings.HasPrefix(lower, "error:"), strings.HasPrefix(lower, "error "),
		strings.HasPrefix(lower, "fatal:"), strings.HasPrefix(lower, "fatal "),
		strings.Contains(lower, "level=error"), strings.Contains(lower, "level=fatal"):
		return "error"
	case strings.HasPrefix(lower, "warning:"), strings.HasPrefix(lower, "warning "),
		strings.HasPrefix(lower, "warn:"), strings.HasPrefix(lower, "warn "),
		strings.Contains(lower, "level=warn"):
		return "warning"
	default:
		return ""
	}
}

// addRecentControlLogs imports a bounded tail of active and recent job files.
// Older full logs remain available through Builder Admin's existing /logs URL.
func (s *Service) addRecentControlLogs(snapshot *servecontrol.Snapshot, state State) {
	logPaths := make(map[string]string, len(state.Builds)+len(state.Refresh)+1)
	for _, record := range state.Builds {
		logPaths[record.ID] = record.LogPath
	}
	for _, record := range state.Refresh {
		logPaths[record.ID] = record.LogPath
	}
	if state.Running != nil {
		logPaths[state.Running.ID] = state.Running.ID + ".log"
	}
	imported := 0
	for i := len(snapshot.Jobs) - 1; i >= 0 && imported < 12; i-- {
		job := &snapshot.Jobs[i]
		path := logPaths[job.ID]
		if path == "" {
			continue
		}
		logs := s.readControlLogTail(path, job.ID)
		if len(logs) == 0 {
			continue
		}
		job.Logs = logs
		snapshot.Logs = append(snapshot.Logs, logs...)
		imported++
	}
}

func (s *Service) readControlLogTail(logPath, jobID string) []servecontrol.LogEntry {
	if logPath == "" || filepath.Base(logPath) != logPath {
		return nil
	}
	file, err := os.Open(filepath.Join(s.logDir, logPath))
	if err != nil {
		return nil
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil
	}
	start := max(int64(0), info.Size()-(16<<10))
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return nil
	}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 16<<10), 1<<20)
	if start > 0 {
		scanner.Scan() // Discard the partial first line.
	}
	logs := make([]servecontrol.LogEntry, 0, 64)
	for scanner.Scan() {
		line := scanner.Text()
		logs = append(logs, servecontrol.LogEntry{JobID: jobID, Level: controlLogSeverity(line), Message: line})
		if len(logs) > 200 {
			logs = logs[1:]
		}
	}
	return logs
}

// addCurrentReleaseDiagnostics reads only the live release's bounded,
// validated diagnostics artifact. Old releases retain their own artifact.
func (s *Service) addCurrentReleaseDiagnostics(snapshot *servecontrol.Snapshot, state State) {
	if s.cfg.SiteDir == "" {
		return
	}
	releaseID := s.currentReleaseID()
	if releaseID == "" {
		return
	}
	path := filepath.Join(s.cfg.SiteDir, "current", diagnostics.DefaultArtifactPath)
	info, err := os.Stat(path)
	if err != nil || info.Size() > 16<<20 {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	artifact, err := diagnostics.ParseArtifact(data)
	if err != nil {
		return
	}
	jobID := ""
	for _, build := range state.Builds {
		if build.ReleaseID == releaseID && build.Kind == "build" && build.Status == "success" {
			jobID = build.ID
			break
		}
	}
	jobIndex := -1
	for i := range snapshot.Jobs {
		if snapshot.Jobs[i].ID == jobID {
			jobIndex = i
			break
		}
	}
	for _, entry := range artifact.Entries {
		page := servecontrol.Page{Path: entry.Path, LastJobID: jobID}
		if entry.Emitted {
			page.Status = servecontrol.StateSuccess
		}
		for _, issue := range entry.Diagnostics {
			diagnostic := servecontrol.Diagnostic{
				Code: issue.Code, Severity: issue.Severity.String(), Message: issue.Message,
				File: issue.File, Page: entry.Path, Line: issue.Range.StartLine + 1,
				Column: issue.Range.StartCol + 1, JobID: jobID,
			}
			page.Diagnostics = append(page.Diagnostics, diagnostic)
			snapshot.Diagnostics = append(snapshot.Diagnostics, diagnostic)
			if jobIndex >= 0 {
				snapshot.Jobs[jobIndex].Diagnostics = append(snapshot.Jobs[jobIndex].Diagnostics, diagnostic)
			}
			if diagnostic.Severity == "error" {
				page.Status = servecontrol.StateFailed
			} else if diagnostic.Severity == "warning" && page.Status != servecontrol.StateFailed {
				page.Status = servecontrol.StateWarning
			}
		}
		snapshot.Pages = append(snapshot.Pages, page)
	}
}

func (s *Service) discoverReleases() []ReleaseView {
	releasesDir := filepath.Join(s.cfg.SiteDir, "releases")
	entries, err := os.ReadDir(releasesDir)
	if err != nil {
		return nil
	}
	current := s.currentReleaseID()
	type releaseBuildMeta struct {
		id            string
		status        string
		finished      time.Time
		builderCommit string
		sourceCommit  string
	}
	buildByRelease := make(map[string]releaseBuildMeta)
	for _, record := range s.viewState().Builds {
		if record.ReleaseID != "" && buildByRelease[record.ReleaseID].id == "" {
			buildByRelease[record.ReleaseID] = releaseBuildMeta{
				id:            record.ID,
				status:        record.Status,
				finished:      record.FinishedAt,
				builderCommit: record.BuilderCommit,
				sourceCommit:  record.SourceCommit,
			}
		}
	}
	views := make([]ReleaseView, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		meta := buildByRelease[entry.Name()]
		createdAt := info.ModTime().UTC()
		if !meta.finished.IsZero() {
			createdAt = meta.finished.UTC()
		} else if parsed, ok := releaseTimestampFromID(entry.Name()); ok {
			createdAt = parsed
		}
		views = append(views, ReleaseView{
			ID:            entry.Name(),
			Path:          filepath.Join(releasesDir, entry.Name()),
			CreatedAt:     createdAt,
			Current:       entry.Name() == current,
			BuildID:       meta.id,
			BuildStatus:   meta.status,
			BuilderCommit: meta.builderCommit,
			SourceCommit:  meta.sourceCommit,
			RollbackOnly:  true,
		})
	}
	sort.Slice(views, func(i, j int) bool {
		if views[i].Current != views[j].Current {
			return views[i].Current
		}
		if !views[i].CreatedAt.Equal(views[j].CreatedAt) {
			return views[i].CreatedAt.After(views[j].CreatedAt)
		}
		return views[i].ID > views[j].ID
	})
	return views
}

func (s *Service) currentReleasePath() string {
	target, err := os.Readlink(filepath.Join(s.cfg.SiteDir, "current"))
	if err != nil {
		return ""
	}
	if filepath.IsAbs(target) {
		return target
	}
	return filepath.Join(s.cfg.SiteDir, target)
}

func (s *Service) currentReleaseID() string {
	path := s.currentReleasePath()
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

func addDirRecursive(w *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if ignoreWatchPath(root, path) && path != root {
			return filepath.SkipDir
		}
		return w.Add(path)
	})
}

func addDirRecursiveIfDir(w *fsnotify.Watcher, path string) error {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return nil
	}
	return addDirRecursive(w, path)
}

func ignoreWatchPath(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	if rel == "." {
		return false
	}
	parts := strings.Split(rel, "/")
	for _, part := range parts {
		switch part {
		case ".git", ".markata", ".markata-cache", ".builder-admin":
			return true
		}
	}
	return false
}

func extractPerfSummaryFromFile(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	summary := make([]string, 0, 24)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "] took ") || strings.Contains(line, "Hotspots:") || strings.Contains(line, "Duration:") {
			summary = append(summary, line)
		}
	}
	if len(summary) > 24 {
		return summary[len(summary)-24:]
	}
	return summary
}

func nextID(prefix string) string {
	return prefix + "-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

func hostSuffix() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "host"
	}
	return host
}

func formatUITimestamp(ts, now time.Time) string {
	if ts.IsZero() {
		return ""
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return fmt.Sprintf("%s (%s ago)", ts.UTC().Format(time.RFC3339), humanizeAge(now.Sub(ts)))
}

func humanizeAge(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Round(time.Second)/time.Second))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Round(time.Minute)/time.Minute))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Round(time.Hour)/time.Hour))
	default:
		return fmt.Sprintf("%dd", int(d.Round(24*time.Hour)/(24*time.Hour)))
	}
}

func formatQueueWait(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	return humanizeAge(d)
}

func uiStatusClass(value string) string {
	switch value {
	case "success", "live", "ready":
		return "status-success"
	case "idle":
		return "status-neutral"
	case "running", "build", "refresh", "promote", "prepare", "prune":
		return "status-running"
	case "queued", "pending", "starting":
		return "status-queued"
	case "failed", "error", "cancelled": //nolint:misspell // Matches the persisted status contract.
		return "status-failed"
	default:
		return "status-neutral"
	}
}

func releaseTimestampFromID(id string) (time.Time, bool) {
	formats := []struct {
		layout string
		length int
	}{
		{layout: "20060102T150405Z", length: len("20060102T150405Z")},
		{layout: "20060102150405", length: len("20060102150405")},
	}
	for _, format := range formats {
		if len(id) < format.length {
			continue
		}
		candidate := id[:format.length]
		ts, err := time.Parse(format.layout, candidate)
		if err == nil {
			return ts.UTC(), true
		}
	}
	return time.Time{}, false
}

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Builder Admin</title>
  <link id="app-favicon" rel="icon" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 64 64'%3E%3Crect width='64' height='64' rx='16' fill='%23525562'/%3E%3Ccircle cx='32' cy='32' r='13' fill='none' stroke='white' stroke-width='6' stroke-linecap='round' stroke-dasharray='0.01 82'/%3E%3C/svg%3E">
  <style>{{ browserTokens .Theme }}</style>
  <style>
    :root {
	  --bg: var(--markata-background);
	  --panel: var(--markata-panel);
	  --panel-strong: var(--markata-surface);
	  --line: var(--markata-border);
	  --line-soft: var(--markata-elevated);
	  --text: var(--markata-text-primary);
	  --muted: var(--markata-text-secondary);
	  --accent: var(--markata-accent);
	  --link: var(--markata-link);
	  --focus: var(--markata-focus);
	  --success: var(--markata-success);
	  --warning: var(--markata-warning);
	  --error: var(--markata-error);
	  --info: var(--markata-info);
	  --code-bg: var(--markata-code-background);
	  --code-text: var(--markata-code-text);
	  --button-bg: var(--markata-button-background);
	  --button-text: var(--markata-button-text);
    }
    * { box-sizing: border-box; }
    html { background: var(--bg); }
    body {
      margin: 0;
      color: var(--text);
      font-family: Inter, ui-sans-serif, system-ui, sans-serif;
      background: var(--bg);
    }
	 a { color: var(--link); text-decoration: none; }
    a:hover { text-decoration: underline; }
    main {
      width: 100%;
      max-width: 1800px;
      margin: 0 auto;
      padding: 24px 28px 56px;
    }
    h1, h2, h3, p { margin: 0; }
    .topbar {
      display: flex;
      justify-content: space-between;
      gap: 20px;
      align-items: flex-start;
      margin-bottom: 28px;
      padding-bottom: 18px;
      border-bottom: 1px solid var(--line-soft);
    }
    .titleblock h1 {
      font-size: clamp(1.6rem, 3vw, 2.3rem);
      line-height: 1;
      letter-spacing: -0.035em;
      text-transform: none;
    }
    .titleblock p {
      margin-top: 10px;
      color: var(--muted);
      max-width: 62ch;
    }
    .title-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      margin-top: 16px;
    }
    .meta-chip {
      display: inline-flex;
      gap: 8px;
      align-items: center;
      padding: 5px 9px;
      border: 1px solid var(--line-soft);
      border-radius: 999px;
      color: var(--muted);
      font-size: 0.78rem;
    }
    .meta-chip strong {
      margin: 0;
      color: var(--text);
      letter-spacing: 0;
      font-size: 0.78rem;
    }
    .revision-banner { border: 1px solid var(--line-soft); border-left: 4px solid var(--accent); border-radius: 8px; background: var(--panel); padding: 16px 18px; margin: 0 0 20px; }
    .revision-banner[data-revision-status="unbuilt"] { border-left-color: var(--warning); background: color-mix(in srgb, var(--warning) 7%, var(--panel)); }
    .revision-banner[data-revision-status="built"] { border-left-color: var(--success); }
    .revision-head { display: flex; align-items: center; flex-wrap: wrap; gap: 10px; justify-content: space-between; margin-bottom: 13px; }
    .revision-head strong { font-size: 1rem; }
    .revision-comparison { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 13px; }
    .revision-comparison > div { min-width: 0; display: flex; flex-direction: column; gap: 5px; }
    .revision-comparison span { font-size: .67rem; font-weight: 700; color: var(--muted); letter-spacing: .08em; }
    .revision-comparison code { align-self: flex-start; font-size: .86rem; max-width: 100%; }
    .revision-comparison small, #revision-detail { color: var(--muted); font-size: .78rem; }
    #revision-detail { margin-top: 12px; }
    .hero {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 24px;
      align-items: center;
      padding: 20px 0;
      border-top: 1px solid var(--line-soft);
      border-bottom: 1px solid var(--line-soft);
      margin-bottom: 28px;
    }
    .section-grid {
      display: grid;
      grid-template-columns: minmax(0, 1.3fr) minmax(18rem, .7fr);
      gap: 28px;
      margin-bottom: 32px;
    }
    .support-panel {
      min-width: 0;
    }
    .hero strong, .support-panel strong, .muted-label {
      display: block;
      font-size: 0.72rem;
      letter-spacing: 0.14em;
      text-transform: uppercase;
      color: var(--muted);
      margin-bottom: 6px;
    }
    .value {
      font-size: 1.15rem;
      line-height: 1.25;
      word-break: break-word;
    }
    .actions {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      align-content: center;
      justify-content: flex-end;
    }
    .actions form { margin: 0; }
    button {
      background: var(--text);
	  color: var(--button-text);
	  background: var(--button-bg);
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 11px 15px;
      cursor: pointer;
      text-transform: uppercase;
      letter-spacing: 0.02em;
      font-size: 0.82rem;
      font-weight: 700;
    }
    button.secondary { background: transparent; color: var(--text); }
	button:hover { background: var(--accent); }
    button.secondary:hover { background: var(--panel); }
    .stack { display: flex; flex-direction: column; gap: 10px; }
    .panel-head {
      display: flex;
      justify-content: space-between;
      align-items: baseline;
      gap: 12px;
      margin-bottom: 12px;
    }
    .panel-head h2 { font-size: 1rem; letter-spacing: -0.01em; }
    .panel-head span { color: var(--muted); font-size: 0.8rem; }
    .workspace-head {
      margin-bottom: 10px;
    }
    table { width: 100%; border-collapse: collapse; }
    th, td {
      text-align: left;
      padding: 13px 10px;
      border-top: 1px solid var(--line-soft);
      vertical-align: top;
      font-size: 0.9rem;
    }
    th {
      color: var(--muted);
      letter-spacing: 0.04em;
      font-size: 0.72rem;
    }
    code {
      display: inline-block;
	  background: var(--code-bg);
	  border: 1px solid var(--line-soft);
      border-radius: 6px;
      padding: 3px 8px;
      white-space: nowrap;
      max-width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
	  color: var(--code-text);
    }
    pre {
      margin: 0;
      padding: 10px 12px;
      overflow: auto;
      white-space: pre-wrap;
	  background: var(--code-bg);
      border: 1px solid var(--line-soft);
      border-radius: 8px;
      max-height: 11rem;
      line-height: 1.45;
	  color: var(--code-text);
      font-size: 0.82rem;
    }
    .pill { --status-ink: var(--muted); display: inline-flex; align-items: center; min-height: 1.7rem; padding: 3px 8px 3px 10px; border: 0; border-left: 2px solid var(--status-ink); border-radius: 3px; background: color-mix(in srgb, var(--status-ink) 7%, var(--panel)); color: var(--text); font-size: 0.72rem; font-weight: 650; letter-spacing: 0.01em; text-transform: none; }
	.status-success { --status-ink: var(--success); }
	.status-running { --status-ink: var(--info); }
	.status-queued { --status-ink: var(--warning); }
	.status-failed { --status-ink: var(--error); }
	.status-neutral { --status-ink: var(--muted); }
    .work-row { cursor: pointer; }
    .work-row:hover > td { background: color-mix(in srgb, var(--panel-strong) 55%, transparent); }
    .work-row:focus-visible { outline: 2px solid var(--focus); outline-offset: -2px; }
    .work-row > td:first-child { padding-left: 14px; }
    .detail-row { display: none; }
    .work-row.is-expanded + .detail-row { display: table-row; }
    .detail-row td { padding: 0 14px 16px; border-top: 0; }
    .row-detail { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 16px; padding: 14px 16px; background: var(--panel-strong); border-left: 2px solid var(--line); }
    .row-detail-meta { display: flex; flex-wrap: wrap; gap: 8px 16px; color: var(--muted); font-size: .82rem; }
    .row-detail-meta code { max-width: min(42rem, 100%); }
    .record-link { display: inline-flex; align-items: center; min-height: 2.25rem; padding: 0 11px; border: 1px solid var(--line); border-radius: 5px; color: var(--text); font-weight: 700; white-space: nowrap; background: color-mix(in srgb, var(--accent) 8%, var(--panel)); }
    .record-link:hover { text-decoration: none; background: color-mix(in srgb, var(--accent) 16%, var(--panel)); }
    .release-name { display: block; font-weight: 650; white-space: nowrap; }
    .release-ref { display: block; margin-top: 3px; color: var(--muted); font-size: .76rem; }
    .build-time { min-width: 9rem; }
    .build-time-label { display: flex; justify-content: space-between; gap: 8px; font-variant-numeric: tabular-nums; font-size: .82rem; }
    .build-time-label span { color: var(--muted); font-size: .72rem; }
    .build-time-bar { position: relative; height: 4px; margin-top: 7px; overflow: visible; border-radius: 999px; background: var(--line-soft); }
    .build-time-fill { position: absolute; inset: 0 auto 0 0; width: 0; border-radius: inherit; background: var(--accent); }
    .build-time-fill.is-slow { background: var(--warning); }
    .build-time-fill.is-extreme { background: var(--error); }
    .build-time-fill.is-fast { background: var(--info); }
    .build-time-marker { position: absolute; top: -3px; bottom: -3px; width: 1px; background: var(--muted); opacity: .8; }
    .build-time-marker.max { opacity: .35; }
    .impact-label { margin-left: 8px; color: var(--muted); font-size: .78rem; }
    .eta-bar { height: 4px; max-width: 22rem; margin: 8px 0 4px; background: var(--line-soft); }
    .eta-bar i { display: block; height: 100%; background: var(--info); }
    .sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0; }
    .build-details { min-width: 9rem; }
    .build-details summary { cursor: pointer; color: var(--text); font-weight: 600; }
    .build-details[open] summary { margin-bottom: 10px; }
    .detail-list { display: grid; gap: 6px; color: var(--muted); font-size: 0.82rem; }
    .detail-list strong { color: var(--text); }
    .summary-meta { color: var(--muted); font-size: 0.76rem; margin-bottom: 6px; }
    .summary-list { display: grid; gap: 6px; }
	.summary-list div { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: var(--code-text); }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; }
    .wide { overflow-x: auto; }
    .muted { color: var(--muted); }
    .time-stamp { white-space: nowrap; }
    .tabs {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-bottom: 14px;
    }
    .tabs a {
      display: inline-flex;
      align-items: center;
      border: 1px solid var(--line-soft);
      border-radius: 7px;
      padding: 9px 14px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      font-size: 0.72rem;
      color: var(--muted);
	  background: var(--panel);
    }
    .tabs a.active {
      color: var(--text);
      background: var(--panel-strong);
      border-color: var(--line);
    }
    .sync-status { color: var(--muted); font-size: 0.78rem; }
    .operator-profile { text-align: right; max-width: 28rem; }
    .operator-profile strong, .operator-profile code { display: block; }
    .operator-profile span { color: var(--muted); font-size: 0.78rem; }
    .operator-identity { display: flex; justify-content: flex-end; align-items: center; gap: 10px; }
    .operator-avatar { width: 36px; height: 36px; border: 1px solid var(--line); border-radius: 50%; object-fit: cover; }
    .operator-avatar-fallback { display: inline-grid; place-items: center; width: 36px; height: 36px; border: 1px solid var(--line); border-radius: 50%; color: var(--muted); background: var(--panel-strong); }
    .operator-avatar-fallback[hidden] { display: none; }
    .operator-avatar-fallback svg { width: 18px; height: 18px; fill: currentColor; }
    .tab-shell { border-top: 1px solid var(--line-soft); padding-top: 22px; }
    .tab-panel { display: none; }
    .tab-panel.is-active { display: block; }
    @media (max-width: 1200px) {
      .hero, .section-grid { grid-template-columns: 1fr 1fr; }
      .revision-comparison { grid-template-columns: 1fr 1fr; }
    }
    @media (max-width: 800px) {
      main { padding: 18px 14px 36px; }
      .topbar, .hero, .section-grid { grid-template-columns: 1fr; display: grid; }
      .revision-comparison { grid-template-columns: 1fr; }
      .topbar { gap: 14px; }
      .operator-identity { justify-content: flex-start; text-align: left; }
      .operator-profile { text-align: left; }
      .actions { justify-content: flex-start; }
      .wide { overflow-x: visible; }
      .responsive-table thead { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0; }
      .responsive-table, .responsive-table tbody, .responsive-table tr, .responsive-table td { display: block; width: 100%; }
      .responsive-table tr { padding: 12px 0; border-top: 1px solid var(--line-soft); }
      .responsive-table td { border: 0; padding: 5px 0; display: grid; grid-template-columns: minmax(7rem, 35%) 1fr; gap: 12px; }
      .responsive-table td::before { content: attr(data-label); color: var(--muted); font-size: 0.72rem; font-weight: 700; letter-spacing: 0.04em; }
      .responsive-table td:empty { display: none; }
      .responsive-table td:empty::before { content: none; }
      .time-stamp { white-space: normal; }
      .work-row.is-expanded + .detail-row { display: block; padding: 0; }
      .detail-row td { display: block; padding: 0 0 14px; }
      .detail-row td::before { content: none; }
      .row-detail { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body data-capabilities='{{ json .Capabilities }}'>
<main>
  <div class="topbar">
    <div class="titleblock">
      <h1>Builder admin</h1>
      <p>Run and monitor site builds. Start with the active work and recent build results.</p>
      <div class="title-meta">
        <div class="meta-chip">Site <strong id="current-release">Live</strong></div>
        <div class="meta-chip">Queued <strong id="queue-count">{{ len .State.Queue }}</strong></div>
        <div class="meta-chip">Recent builds <strong id="build-count">{{ len .State.Builds }}</strong></div>
      </div>
    </div>
    <div class="operator-identity" aria-label="Authenticated operator">
      {{ if .PictureURL }}<img class="operator-avatar" src="{{ .PictureURL }}" alt="" referrerpolicy="no-referrer" onerror="this.hidden=true;this.nextElementSibling.hidden=false"><span class="operator-avatar-fallback" hidden role="img" aria-label="Profile picture unavailable"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 12a4.5 4.5 0 1 0 0-9 4.5 4.5 0 0 0 0 9Zm0 2c-4.5 0-8.25 2.3-8.25 5.25V21h16.5v-1.75C20.25 16.3 16.5 14 12 14Z"/></svg></span>{{ else }}<span class="operator-avatar-fallback" role="img" aria-label="Profile picture unavailable"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 12a4.5 4.5 0 1 0 0-9 4.5 4.5 0 0 0 0 9Zm0 2c-4.5 0-8.25 2.3-8.25 5.25V21h16.5v-1.75C20.25 16.3 16.5 14 12 14Z"/></svg></span>{{ end }}
      <div class="operator-profile">
        {{ if .Operator.DisplayName }}<strong>{{ .Operator.DisplayName }}</strong>{{ else if .Operator.Username }}<strong>{{ .Operator.Username }}</strong>{{ else }}<strong>Signed in</strong>{{ end }}
      </div>
    </div>
    <div class="sync-status" id="sync-status">Live polling every 2s</div>
  </div>

  <section id="revision-banner" class="revision-banner" role="status" aria-live="polite" data-revision-status="{{ .Revision.Status }}">
    <div class="revision-head"><strong id="revision-title">{{ if eq .Revision.Status "unbuilt" }}New Builder revision has not built yet{{ else if eq .Revision.Status "built" }}Current Builder revision has built successfully{{ else }}Builder revision status not yet verified{{ end }}</strong><span id="revision-marker" class="pill {{ if eq .Revision.Status "unbuilt" }}status-queued{{ else if eq .Revision.Status "built" }}status-success{{ else }}status-neutral{{ end }}">{{ .Revision.Status }}</span></div>
    <div class="revision-comparison"><div><span>RUNNING BUILDER</span><code id="revision-current">{{ if .Revision.BuilderCommit }}{{ .Revision.BuilderCommit }}{{ else }}Unknown{{ end }}</code><small id="revision-version">{{ .Revision.BuilderVersion }}</small></div><div><span>LAST SUCCESSFUL BUILD</span><code id="revision-last-built">{{ if .Revision.LastBuiltCommit }}{{ .Revision.LastBuiltCommit }}{{ else }}Unknown{{ end }}</code><small id="revision-built-time">{{ if .Revision.LastBuiltAt.IsZero }}No dated build record{{ else }}{{ since .Revision.LastBuiltAt }}{{ end }}</small></div><div><span>SITE SOURCE</span><code id="revision-source">{{ if .Revision.SourceCommit }}{{ .Revision.SourceCommit }}{{ else }}Unavailable{{ end }}</code><small id="revision-source-status"></small></div></div>
    <p id="revision-detail" class="muted">A successful build from this Builder revision will clear the unbuilt notice.</p>
  </section>

  <div class="hero" aria-labelledby="live-work-heading">
    <section>
      <div class="panel-head"><h2 id="live-work-heading">Active work</h2><span>one operation runs at a time</span></div>
      <div style="display:grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px;">
        <div>
          <strong>Queued</strong>
          <div class="value" id="queue-overview">{{ len .State.Queue }}</div>
        </div>
        <div>
          <strong>Current work</strong>
          <div class="value" id="active-work">{{ if .State.Running }}{{ .State.Running.Kind }} <span class="pill {{ statusClass .State.Running.Phase }}">{{ .State.Running.Phase }}</span>{{ else }}<span class="pill {{ statusClass "idle" }}">idle</span>{{ end }}</div>
          <p class="muted" id="active-work-detail">{{ if .State.Running }}Started {{ since .State.Running.StartedAt }}{{ else }}No build or refresh is running.{{ end }}</p>
        </div>
      </div>
    </section>
    <section class="actions" aria-label="Build actions">
      <form method="post" action="/api/builds" data-capability="enqueueBuild"><input type="hidden" name="csrf_token" value="{{ .CSRFToken }}"><button type="submit">Enqueue Build</button></form>
      {{ range .RefreshTasks }}
      <form method="post" action="/api/refresh/{{ .Name }}" data-capability="enqueueRefresh"><input type="hidden" name="csrf_token" value="{{ $.CSRFToken }}"><button class="secondary" type="submit">Run {{ .Name }}</button></form>
      {{ end }}
    </section>
  </div>

  <section class="wide tab-shell">
    <div class="panel-head workspace-head"><h2>Jobs</h2><span>Running and queued work first; open details for timings and logs.</span></div>
    <nav class="tabs">
      <a href="#jobs" data-tab-link="jobs">Jobs</a>
      <a href="#releases" data-tab-link="releases">Releases</a>
    </nav>

    <section id="jobs" class="tab-panel" data-tab-panel="jobs">
      <table class="responsive-table">
        <thead><tr><th>Status</th><th>Job</th><th>Source</th><th>When</th><th>Wait</th><th>Run</th><th>Release</th><th>Details</th></tr></thead>
        <tbody id="builds-body">
        {{ if .State.Running }}
        <tr><td data-label="Status"><span class="pill status-running">running</span></td><td data-label="Job">{{ .State.Running.Label }}</td><td data-label="Source">{{ .State.Running.TriggerType }}</td><td data-label="When">Started {{ since .State.Running.StartedAt }}</td><td data-label="Wait">—</td><td data-label="Run">—</td><td data-label="Release"></td><td data-label="Details"><details class="build-details"><summary>View</summary><div class="detail-list"><div><strong>Phase:</strong> {{ .State.Running.Phase }}</div><div><strong>Job ID:</strong> <code>{{ .State.Running.ID }}</code></div></div></details></td></tr>
        {{ end }}
        {{ range .State.Queue }}
        <tr><td data-label="Status"><span class="pill status-queued">queued</span></td><td data-label="Job">{{ .Label }}</td><td data-label="Source">{{ .TriggerType }}</td><td data-label="When">Waiting</td><td data-label="Wait">{{ since .EnqueuedAt }}</td><td data-label="Run">—</td><td data-label="Release"></td><td data-label="Details"><details class="build-details"><summary>View</summary><div class="detail-list"><div><strong>Reason:</strong> {{ .Detail }}</div><div><strong>Job ID:</strong> <code>{{ .ID }}</code></div></div></details></td></tr>
        {{ end }}
        {{ range .CompletedJobs }}
        <tr data-finished="{{ .FinishedAt }}"><td data-label="Status"><span class="pill {{ statusClass .Status }}">{{ .Status }}</span></td><td data-label="Job">{{ if .Build }}<a class="record-link" href="/builds/{{ .ID }}">View build record →</a>{{ else }}{{ .Kind }}{{ end }}</td><td data-label="Source">{{ .Trigger }}</td><td data-label="When">{{ since .FinishedAt }}</td><td data-label="Wait">{{ msToSeconds .QueueWaitMS }}</td><td data-label="Run">{{ msToSeconds .RunMS }}</td><td data-label="Release">{{ .Release }}</td><td data-label="Details"><details class="build-details"><summary>More</summary><div class="detail-list"><div><strong>Job ID:</strong> <code>{{ .ID }}</code></div>{{ if .Build }}<div><a class="record-link" href="/builds/{{ .ID }}">Open full build details →</a></div>{{ end }}{{ if .LogPath }}<div><a href="/logs/{{ .LogPath }}">Open raw log</a></div>{{ end }}</div></details></td></tr>
        {{ end }}
        </tbody>
      </table>
    </section>

    <section id="refresh-runs" class="tab-panel" data-tab-panel="refresh-runs">
      <table class="responsive-table">
        <thead><tr><th>ID</th><th>Task</th><th>Status</th><th>Total</th><th>Logs</th><th>Build</th><th>Command</th></tr></thead>
        <tbody id="refresh-body">
        {{ range .State.Refresh }}
        <tr>
          <td data-label="Run"><code>{{ .ID }}</code></td>
          <td data-label="Task">{{ .TaskName }}</td>
          <td data-label="Status"><span class="pill {{ statusClass .Status }}">{{ .Status }}</span></td>
          <td data-label="Duration">{{ msToSeconds .TotalMS }}</td>
          <td data-label="Log">{{ if .LogPath }}<a href="/logs/{{ .LogPath }}">Open</a>{{ end }}</td>
          <td data-label="Build">{{ if .EnqueuedBuildID }}<code>{{ .EnqueuedBuildID }}</code>{{ end }}</td>
          <td data-label="Command" class="mono muted">{{ if .Command }}{{ index .Command 0 }} {{ end }}</td>
        </tr>
        {{ else }}
        <tr><td colspan="7">No refresh runs yet.</td></tr>
        {{ end }}
        </tbody>
      </table>
    </section>

    <section id="releases" class="tab-panel" data-tab-panel="releases">
      <table class="responsive-table">
        <thead><tr><th>ID</th><th>Current</th><th>Created</th><th>Build</th><th>Status</th><th>Action</th></tr></thead>
        <tbody id="releases-body">
        {{ range .Releases }}
        <tr class="work-row" data-expandable tabindex="0" aria-expanded="false">
          <td data-label="Release"><span class="release-name">{{ releaseLabel .CreatedAt }}</span><span class="release-ref">release record</span></td>
          <td data-label="Current">{{ if .Current }}<span class="pill {{ statusClass "live" }}">live</span>{{ end }}</td>
          <td data-label="Created" class="time-stamp">{{ since .CreatedAt }}</td>
          <td data-label="Build">{{ if .BuildID }}<a class="record-link" href="/builds/{{ .BuildID }}">View build record →</a>{{ end }}</td>
          <td data-label="Status">{{ if .BuildStatus }}<span class="pill {{ statusClass .BuildStatus }}">{{ .BuildStatus }}</span>{{ end }}</td>
          <td data-label="Action">{{ if not .Current }}<form method="post" action="/api/releases/{{ .ID }}/rollback" data-capability="promoteRelease"><input type="hidden" name="csrf_token" value="{{ $.CSRFToken }}"><button class="secondary" type="submit">Promote</button></form>{{ end }}</td>
        </tr>
        <tr class="detail-row"><td colspan="6"><div class="row-detail"><div class="row-detail-meta"><span>Created {{ .CreatedAt.UTC.Format "02 Jan 2006 15:04:05 UTC" }}</span><code>{{ .ID }}</code>{{ if .BuildID }}<code>{{ .BuildID }}</code>{{ end }}</div>{{ if .BuildID }}<a class="record-link" href="/builds/{{ .BuildID }}">Open full build details →</a>{{ end }}</div></td></tr>
        {{ else }}
        <tr><td colspan="6">No releases found.</td></tr>
        {{ end }}
        </tbody>
      </table>
    </section>
  </section>
</main>
<script>
  const csrfToken = {{ printf "%q" .CSRFToken }};
  const capabilities = JSON.parse(document.body.dataset.capabilities || '{}');
  document.querySelectorAll('[data-capability]').forEach((node) => {
    if (!capabilities[node.dataset.capability]) node.remove();
  });
  const favicon = document.getElementById('app-favicon');
  const syncStatus = document.getElementById('sync-status');
  const currentRelease = document.getElementById('current-release');
  const revisionBanner = document.getElementById('revision-banner');
  const activeWork = document.getElementById('active-work');
  const activeWorkDetail = document.getElementById('active-work-detail');
  const queueCount = document.getElementById('queue-count');
  const queueOverview = document.getElementById('queue-overview');
  const queueSummary = document.getElementById('queue-summary');
  const buildCount = document.getElementById('build-count');
  const latestBuild = document.getElementById('latest-build');
  const liveReleaseValue = document.getElementById('live-release-value');
  const queueBody = document.getElementById('queue-body');
  const buildsBody = document.getElementById('builds-body');
  const refreshBody = document.getElementById('refresh-body');
  const releasesBody = document.getElementById('releases-body');
  let buildsFingerprint = '';
  let refreshFingerprint = '';
  let releasesFingerprint = '';
  let pollInFlight = false;

  function escapeHtml(value) {
    return String(value ?? '')
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#39;');
  }

  function fmtTime(value) {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toISOString().replace('.000', '') + ' (' + timeAgo(date) + ' ago)';
  }

  function fmtSeconds(ms) {
    return ((ms || 0) / 1000).toFixed(2) + 's';
  }

  function timeAgo(date) {
    const delta = Math.max(0, Date.now() - date.getTime());
    const seconds = Math.round(delta / 1000);
    if (seconds < 60) return seconds + 's';
    const minutes = Math.round(seconds / 60);
    if (minutes < 60) return minutes + 'm';
    const hours = Math.round(minutes / 60);
    if (hours < 24) return hours + 'h';
    return Math.round(hours / 24) + 'd';
  }

  function statusClass(value) {
    switch (value) {
      case 'idle':
		return 'status-neutral';
      case 'success':
      case 'live':
      case 'ready':
        return 'status-success';
      case 'running':
      case 'build':
      case 'refresh':
      case 'promote':
      case 'prepare':
      case 'prune':
        return 'status-running';
      case 'queued':
      case 'pending':
      case 'starting':
        return 'status-queued';
      case 'failed':
      case 'error':
      case 'cancelled':
        return 'status-failed';
      default:
        return 'status-neutral';
    }
  }

  function statusPill(value) {
    return '<span class="pill ' + statusClass(value) + '">' + escapeHtml(value) + '</span>';
  }

  function median(values) {
    const sorted = values.slice().sort((a, b) => a - b);
    const middle = Math.floor(sorted.length / 2);
    return sorted.length % 2 ? sorted[middle] : (sorted[middle - 1] + sorted[middle]) / 2;
  }

  function buildTimeBaseline(builds) {
    const cutoff = Date.now() - 30 * 24 * 60 * 60 * 1000;
    const values = (builds || []).filter((build) => build.status === 'success' && Number(build.build_ms) > 0 && new Date(build.finished_at).getTime() >= cutoff).slice(0, 60).map((build) => Number(build.build_ms));
    if (values.length < 8) return null;
    const center = median(values);
    const mad = median(values.map((value) => Math.abs(value - center)));
    const spread = Math.max(mad * 1.4826, center * 0.03);
    return { mean: values.reduce((sum, value) => sum + value, 0) / values.length, max: Math.max(...values), low: Math.max(0, center - 3 * spread), high: center + 3 * spread };
  }

  function buildTimeCell(value, baseline) {
    const duration = Number(value) || 0;
    if (!baseline) return '<span>' + escapeHtml(fmtSeconds(duration)) + '</span>';
    const scale = Math.max(baseline.max, baseline.high, duration, 1);
    const ratio = Math.min(100, duration / scale * 100);
    const mean = Math.min(100, baseline.mean / scale * 100);
    const max = Math.min(100, baseline.max / scale * 100);
    const state = duration > baseline.high * 1.5 ? ' is-extreme' : duration > baseline.high ? ' is-slow' : duration < baseline.low ? ' is-fast' : '';
    const classification = state === ' is-extreme' ? 'extreme slow outlier' : state === ' is-slow' ? 'slow outlier' : state === ' is-fast' ? 'faster than the normal range' : 'within the normal range';
    return '<div class="build-time" title="Mean ' + escapeHtml(fmtSeconds(baseline.mean)) + ' · recorded max ' + escapeHtml(fmtSeconds(baseline.max)) + '"><div class="build-time-label"><strong>' + escapeHtml(fmtSeconds(duration)) + '</strong><span>mean ' + escapeHtml(fmtSeconds(baseline.mean)) + '</span></div><div class="build-time-bar"><i class="build-time-fill' + state + '" style="width:' + ratio.toFixed(2) + '%"></i><i class="build-time-marker" style="left:' + mean.toFixed(2) + '%"></i><i class="build-time-marker max" style="left:' + max.toFixed(2) + '%"></i></div><span class="sr-only">Mean ' + escapeHtml(fmtSeconds(baseline.mean)) + ', recorded maximum ' + escapeHtml(fmtSeconds(baseline.max)) + ', ' + classification + '.</span></div>';
  }

  function makeExpandableRows(scope) {
    scope.querySelectorAll('details.build-details').forEach((details) => {
      const row = details.closest('tr');
      if (!row) return;
      row.classList.add('work-row');
      row.dataset.expandable = '';
      row.tabIndex = 0;
      row.setAttribute('aria-expanded', details.open ? 'true' : 'false');
    });
  }

  document.addEventListener('click', (event) => {
    const row = event.target.closest('tr[data-expandable]');
    if (!row || event.target.closest('a, button, input, form, summary')) return;
    const details = row.querySelector('details');
    if (details) {
      details.open = !details.open;
    } else {
      row.classList.toggle('is-expanded');
    }
    row.setAttribute('aria-expanded', row.classList.contains('is-expanded') || details?.open ? 'true' : 'false');
  });
  document.addEventListener('keydown', (event) => {
    if (event.key !== 'Enter' && event.key !== ' ') return;
    const row = event.target.closest('tr[data-expandable]');
    if (!row || event.target.closest('a, button, input, form, summary')) return;
    event.preventDefault();
    row.click();
  });
  document.addEventListener('toggle', (event) => {
    if (event.target.matches('details.build-details')) {
      event.target.closest('tr')?.setAttribute('aria-expanded', event.target.open ? 'true' : 'false');
    }
  }, true);

  function faviconDataURL(svg) {
    return 'data:image/svg+xml,' + encodeURIComponent(svg);
  }

  function buildFaviconSVG(stateName) {
    const base = {
      idle: '#525562',
      queued: '#d97706',
      build: '#16a34a',
      refresh: '#2563eb',
      error: '#dc2626'
    }[stateName] || '#7c3aed';
    const icon = {
      idle: '<path d="M18 33l9 9 19-20" fill="none" stroke="white" stroke-width="6" stroke-linecap="round" stroke-linejoin="round"/>',
      queued: '<circle cx="22" cy="32" r="4" fill="white"/><circle cx="32" cy="32" r="4" fill="white" opacity="0.85"/><circle cx="42" cy="32" r="4" fill="white" opacity="0.7"/>',
      build: '<circle cx="32" cy="32" r="12" fill="none" stroke="white" stroke-width="6"/><path d="M32 12v9M32 43v9M12 32h9M43 32h9M18 18l6 6M40 40l6 6M46 18l-6 6M18 46l6-6" fill="none" stroke="white" stroke-width="4" stroke-linecap="round"/>',
      refresh: '<path d="M18 28a14 14 0 0 1 24-8" fill="none" stroke="white" stroke-width="6" stroke-linecap="round"/><path d="M43 13v12H31" fill="none" stroke="white" stroke-width="6" stroke-linecap="round" stroke-linejoin="round"/><path d="M46 36a14 14 0 0 1-24 8" fill="none" stroke="white" stroke-width="6" stroke-linecap="round"/><path d="M21 51V39h12" fill="none" stroke="white" stroke-width="6" stroke-linecap="round" stroke-linejoin="round"/>',
      error: '<path d="M32 18v18" fill="none" stroke="white" stroke-width="6" stroke-linecap="round"/><circle cx="32" cy="44" r="3.5" fill="white"/>'
    }[stateName] || '<path d="M20 20l24 24M44 20L20 44" fill="none" stroke="white" stroke-width="6" stroke-linecap="round"/>';
    const svg = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">' +
      '<rect width="64" height="64" rx="16" fill="' + base + '" />' +
      icon +
      '</svg>';
    return svg;
  }

  function updateFavicon(stateName) {
    if (!favicon) {
      return;
    }
    favicon.href = faviconDataURL(buildFaviconSVG(stateName));
  }

  function faviconState(state) {
    if (state && state.running) {
      return state.running.kind === 'refresh' ? 'refresh' : 'build';
    }
    if (state && Array.isArray(state.queue) && state.queue.length > 0) {
      return 'queued';
    }
    return 'idle';
  }

  function summaryPreview(lines) {
    if (!Array.isArray(lines)) return [];
    return lines.slice(-6);
  }

  function activateTabs() {
    const requested = (location.hash || '#jobs').slice(1);
    const active = requested === 'builds' || requested === 'refresh-runs' || requested !== 'releases' && requested !== 'jobs' ? 'jobs' : requested;
    if (requested !== active) {
      history.replaceState(null, '', '#' + active);
    }
    document.querySelectorAll('[data-tab-link]').forEach((link) => {
      link.classList.toggle('active', link.dataset.tabLink === active);
    });
    document.querySelectorAll('[data-tab-panel]').forEach((panel) => {
      panel.classList.toggle('is-active', panel.dataset.tabPanel === active);
    });
  }

  function renderQueue(items) {
    const queued = items || [];
    const oldest = queued.reduce((oldestValue, item) => {
      const timestamp = new Date(item.enqueued_at);
      return !oldestValue || timestamp < oldestValue ? timestamp : oldestValue;
    }, null);
    queueSummary.textContent = queued.length + ' waiting · oldest ' + (oldest ? timeAgo(oldest) : 'none');
    if (!queued.length) {
      queueBody.innerHTML = '<tr><td colspan="4">Queue is empty.</td></tr>';
      return;
    }
    queueBody.innerHTML = queued.map((item) => {
      return '<tr>' +
        '<td data-label="Work"><code>' + escapeHtml(item.id) + '</code> ' + escapeHtml(item.label) + '</td>' +
        '<td data-label="Source">' + escapeHtml(item.trigger_type) + '</td>' +
        '<td data-label="Queued" class="time-stamp">' + escapeHtml(fmtTime(item.enqueued_at)) + '</td>' +
        '<td data-label="Details">' + escapeHtml(item.detail) + '</td>' +
      '</tr>';
    }).join('');
  }

  function renderRunning(running, builds) {
    if (!running) {
      activeWork.innerHTML = statusPill('idle');
      activeWorkDetail.textContent = 'No build or refresh is running.';
      return;
    }
    const baseline = buildTimeBaseline(builds || []);
    const elapsed = Math.max(0, Date.now() - new Date(running.started_at).getTime());
    const expected = baseline ? baseline.mean : 0;
    const progress = expected ? Math.min(100, elapsed / expected * 100) : 0;
    const eta = expected ? Math.max(0, expected - elapsed) : 0;
    activeWork.innerHTML = '<div class="active-build"><div>' + escapeHtml(running.kind) + ' ' + statusPill(running.phase) + '<span class="impact-label">' + escapeHtml(running.impact || 'unknown') + '</span></div>' + (expected ? '<div class="eta-bar"><i style="width:' + progress.toFixed(1) + '%"></i></div><small>Estimated ' + escapeHtml(fmtSeconds(expected)) + ' · about ' + escapeHtml(fmtSeconds(eta)) + ' remaining</small>' : '<small>Learning from recent completed builds</small>') + '</div>';
    activeWorkDetail.textContent = 'Started ' + fmtTime(running.started_at) + ' · ' + (running.detail || running.trigger_type || '');
  }

  function renderBuilds(state) {
    const queued = state.queue || [];
    const builds = state.builds || [];
    const refreshes = state.refresh || [];
    const fingerprint = JSON.stringify({ running: state.running, queued, builds, refreshes });
    if (fingerprint === buildsFingerprint) {
      return;
    }
    const openBuilds = new Set(Array.from(buildsBody.querySelectorAll('details[open]')).map((details) => details.closest('tr')?.querySelector('td code')?.textContent));
    const focusedElement = document.activeElement;
    const focusedBuildID = focusedElement?.closest('tr')?.querySelector('td code')?.textContent;
    const focusedTag = focusedElement?.tagName;
    if (!state.running && !queued.length && !builds.length && !refreshes.length) {
      buildsBody.innerHTML = '<tr><td colspan="6">No jobs yet.</td></tr>';
      buildsFingerprint = fingerprint;
      return;
    }
    const runningRow = state.running ? '<tr><td data-label="Status">' + statusPill('running') + '</td><td data-label="Job">' + escapeHtml(state.running.label) + '</td><td data-label="Source">' + escapeHtml(state.running.trigger_type) + '</td><td data-label="When">Started ' + escapeHtml(fmtTime(state.running.started_at)) + '</td><td data-label="Wait">—</td><td data-label="Run">—</td><td data-label="Release"></td><td data-label="Details"><details class="build-details"><summary>View</summary><div class="detail-list"><div><strong>Phase:</strong> ' + escapeHtml(state.running.phase) + '</div><div><strong>Job ID:</strong> <code>' + escapeHtml(state.running.id) + '</code></div></div></details></td></tr>' : '';
    const queuedRows = queued.map((item) => '<tr><td data-label="Status">' + statusPill('queued') + '</td><td data-label="Job">' + escapeHtml(item.label) + '</td><td data-label="Source">' + escapeHtml(item.trigger_type) + '</td><td data-label="When">Waiting</td><td data-label="Wait">' + escapeHtml(timeAgo(new Date(item.enqueued_at))) + '</td><td data-label="Run">—</td><td data-label="Release"></td><td data-label="Details"><details class="build-details"><summary>View</summary><div class="detail-list"><div><strong>Reason:</strong> ' + escapeHtml(item.detail) + '</div><div><strong>Job ID:</strong> <code>' + escapeHtml(item.id) + '</code></div></div></details></td></tr>').join('');
    const baseline = buildTimeBaseline(builds);
    const buildRows = builds.map((item) => {
      const timings = '<div><strong>Total:</strong> ' + escapeHtml(fmtSeconds(item.total_ms)) + '</div>' +
        '<div><strong>Prepare:</strong> ' + escapeHtml(fmtSeconds(item.prepare_ms)) + ' · <strong>Promote:</strong> ' + escapeHtml(fmtSeconds(item.promote_ms)) + ' · <strong>Prune:</strong> ' + escapeHtml(fmtSeconds(item.prune_ms)) + '</div>';
      const changed = (item.changed_paths || []).map((path) => '<code>' + escapeHtml(path) + '</code>').join(' ');
      const summary = summaryPreview(item.perf_summary || []).map(escapeHtml).join('\n');
      const details = '<details class="build-details"><summary>More</summary><div class="detail-list"><div><strong>Build ID:</strong> <code>' + escapeHtml(item.id) + '</code></div>' + (item.builder_commit ? '<div><strong>Builder revision:</strong> <code title="' + escapeHtml(item.builder_commit) + '">' + escapeHtml(revisionAbbrev(item.builder_commit)) + '</code></div>' : '') + (item.source_commit ? '<div><strong>Source revision:</strong> <code title="' + escapeHtml(item.source_commit) + '">' + escapeHtml(revisionAbbrev(item.source_commit)) + '</code></div>' : '') + timings +
        (item.trigger_detail ? '<div><strong>Reason:</strong> ' + escapeHtml(item.trigger_detail) + '</div>' : '') +
        (changed ? '<div><strong>Changed:</strong> ' + changed + '</div>' : '') +
        (item.log_path ? '<div><a href="/logs/' + encodeURIComponent(item.log_path) + '">Open raw log</a></div>' : '') +
        (summary ? '<pre>' + summary + '</pre>' : '') + '</div></details>';
      return '<tr class="work-row" data-expandable tabindex="0" aria-expanded="false" data-finished="' + escapeHtml(item.finished_at) + '">' +
        '<td data-label="Status">' + statusPill(item.status) + '</td>' +
        '<td data-label="Job"><a class="record-link" href="/builds/' + encodeURIComponent(item.id) + '">' + (item.rollback_release ? 'View promotion record →' : 'View build record →') + '</a></td>' +
        '<td data-label="Source">' + escapeHtml(item.trigger_type) + (item.builder_commit ? '<br><code title="Builder ' + escapeHtml(item.builder_commit) + '">' + escapeHtml(revisionAbbrev(item.builder_commit)) + '</code>' : '') + '</td>' +
        '<td data-label="When" class="time-stamp">' + escapeHtml(fmtTime(item.finished_at)) + '</td>' +
        '<td data-label="Wait">' + escapeHtml(fmtSeconds(item.queue_wait_ms)) + '</td>' +
        '<td data-label="Run">' + (item.status === 'success' ? buildTimeCell(item.build_ms, baseline) : escapeHtml(fmtSeconds(item.build_ms))) + '</td>' +
        '<td data-label="Release">' + (item.became_live ? 'Live' : '') + '</td>' +
        '<td data-label="Details">' + details + '</td>' +
      '</tr>';
    }).join('');
    const refreshRows = refreshes.map((item) => '<tr data-finished="' + escapeHtml(item.finished_at) + '"><td data-label="Status">' + statusPill(item.status) + '</td><td data-label="Job">Refresh: ' + escapeHtml(item.task_name) + '</td><td data-label="Source">' + escapeHtml(item.trigger_type) + '</td><td data-label="When">' + escapeHtml(fmtTime(item.finished_at)) + '</td><td data-label="Wait">' + escapeHtml(fmtSeconds(item.queue_wait_ms)) + '</td><td data-label="Run">' + escapeHtml(fmtSeconds(item.run_ms)) + '</td><td data-label="Release"></td><td data-label="Details"><details class="build-details"><summary>View</summary><div class="detail-list"><div><strong>Job ID:</strong> <code>' + escapeHtml(item.id) + '</code></div>' + (item.log_path ? '<div><a href="/logs/' + encodeURIComponent(item.log_path) + '">Open raw log</a></div>' : '') + '</div></details></td></tr>').join('');
    buildsBody.innerHTML = runningRow + queuedRows + buildRows + refreshRows;
    makeExpandableRows(buildsBody);
    Array.from(buildsBody.querySelectorAll('tr[data-finished]')).sort((a, b) => new Date(b.dataset.finished) - new Date(a.dataset.finished)).forEach((row) => buildsBody.append(row));
    buildsFingerprint = fingerprint;
    openBuilds.forEach((buildID) => {
      const details = Array.from(buildsBody.querySelectorAll('details')).find((item) => item.closest('tr')?.querySelector('td code')?.textContent === buildID);
      if (details) {
        details.open = true;
      }
    });
    if (focusedBuildID && focusedTag) {
      const row = Array.from(buildsBody.querySelectorAll('tr')).find((item) => item.querySelector('td code')?.textContent === focusedBuildID);
      row?.querySelector(focusedTag.toLowerCase())?.focus();
    }
  }

  function renderRefresh(items) {
    const fingerprint = JSON.stringify(items || []);
    if (fingerprint === refreshFingerprint) {
      return;
    }
    const focusedElement = document.activeElement;
    const focusedRefreshID = focusedElement?.closest('tr')?.querySelector('td code')?.textContent;
    const focusedTag = focusedElement?.tagName;
    if (!items || !items.length) {
      refreshBody.innerHTML = '<tr><td colspan="7">No refresh runs yet.</td></tr>';
      refreshFingerprint = fingerprint;
      return;
    }
    refreshBody.innerHTML = items.map((item) => {
      const command = Array.isArray(item.command) && item.command.length ? item.command.join(' ') : '';
      return '<tr>' +
        '<td data-label="Run"><code>' + escapeHtml(item.id) + '</code></td>' +
        '<td data-label="Task">' + escapeHtml(item.task_name) + '</td>' +
        '<td data-label="Status">' + statusPill(item.status) + '</td>' +
        '<td data-label="Duration">' + escapeHtml(fmtSeconds(item.total_ms)) + '</td>' +
        '<td data-label="Log">' + (item.log_path ? '<a href="/logs/' + encodeURIComponent(item.log_path) + '">Open</a>' : '') + '</td>' +
        '<td data-label="Build">' + (item.enqueued_build_id ? '<code>' + escapeHtml(item.enqueued_build_id) + '</code>' : '') + '</td>' +
        '<td data-label="Command" class="mono muted">' + escapeHtml(command) + '</td>' +
      '</tr>';
    }).join('');
    refreshFingerprint = fingerprint;
    if (focusedRefreshID && focusedTag) {
      const row = Array.from(refreshBody.querySelectorAll('tr')).find((item) => item.querySelector('td code')?.textContent === focusedRefreshID);
      row?.querySelector(focusedTag.toLowerCase())?.focus();
    }
  }

  function renderReleases(items) {
    const fingerprint = JSON.stringify(items || []);
    if (fingerprint === releasesFingerprint) {
      return;
    }
    const focusedElement = document.activeElement;
    const focusedReleaseID = focusedElement?.closest('tr')?.querySelector('td code')?.textContent;
    const focusedTag = focusedElement?.tagName;
    if (!items || !items.length) {
      releasesBody.innerHTML = '<tr><td colspan="6">No releases found.</td></tr>';
      releasesFingerprint = fingerprint;
      return;
    }
    releasesBody.innerHTML = items.map((item) => {
      const action = item.current || !capabilities.promoteRelease ? '' : '<form method="post" data-capability="promoteRelease" action="/api/releases/' + encodeURIComponent(item.id) + '/rollback"><input type="hidden" name="csrf_token" value="' + csrfToken + '"><button class="secondary" type="submit">Promote</button></form>';
      const created = new Date(item.created_at);
      const label = Number.isNaN(created.getTime()) ? 'Release' : created.toLocaleString(undefined, { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' });
      const detail = '<tr class="detail-row"><td colspan="6"><div class="row-detail"><div class="row-detail-meta"><span>Created ' + escapeHtml(fmtTime(item.created_at)) + '</span><code>' + escapeHtml(item.id) + '</code>' + (item.build_id ? '<code>' + escapeHtml(item.build_id) + '</code>' : '') + (item.builder_commit ? '<code title="Builder revision ' + escapeHtml(item.builder_commit) + '">' + escapeHtml(revisionAbbrev(item.builder_commit)) + '</code>' : '') + (item.source_commit ? '<code title="Source revision ' + escapeHtml(item.source_commit) + '">' + escapeHtml(revisionAbbrev(item.source_commit)) + '</code>' : '') + '</div>' + (item.build_id ? '<a class="record-link" href="/builds/' + encodeURIComponent(item.build_id) + '">Open full build details →</a>' : '') + '</div></td></tr>';
      return '<tr class="work-row" data-expandable tabindex="0" aria-expanded="false">' +
        '<td data-label="Release"><span class="release-name">' + escapeHtml(label) + '</span><span class="release-ref">release record</span></td>' +
        '<td data-label="Current">' + (item.current ? statusPill('live') : '') + '</td>' +
        '<td data-label="Created" class="time-stamp">' + escapeHtml(fmtTime(item.created_at)) + '</td>' +
        '<td data-label="Build">' + (item.build_id ? '<a class="record-link" href="/builds/' + encodeURIComponent(item.build_id) + '">View build record →</a>' : '') + '</td>' +
        '<td data-label="Status">' + (item.build_status ? statusPill(item.build_status) : '') + '</td>' +
        '<td data-label="Action">' + action + '</td>' +
      '</tr>' + detail;
    }).join('');
    releasesFingerprint = fingerprint;
    if (focusedReleaseID && focusedTag) {
      const row = Array.from(releasesBody.querySelectorAll('tr')).find((item) => item.querySelector('td code')?.textContent === focusedReleaseID);
      row?.querySelector(focusedTag.toLowerCase())?.focus();
    }
  }

  function revisionAbbrev(value) { return value && value.length > 12 ? value.slice(0, 12) : (value || 'Unknown'); }

  function renderRevision(revision, state) {
    if (!revision) return;
    revisionBanner.dataset.revisionStatus = revision.status || 'unknown';
    const pending = (state.running && state.running.kind === 'build') || (state.queue || []).some((item) => item.kind === 'build');
    const unbuilt = revision.status === 'unbuilt';
    document.getElementById('revision-title').textContent = unbuilt ? (pending ? 'New Builder revision — build pending' : 'New Builder revision has not built yet') : revision.status === 'built' ? 'Current Builder revision has built successfully' : 'Builder revision status not yet verified';
    const marker = document.getElementById('revision-marker');
    marker.textContent = unbuilt ? (pending ? 'waiting for build' : 'not built yet') : revision.status === 'built' ? 'built' : 'unknown';
    marker.className = 'pill ' + (unbuilt ? 'status-queued' : revision.status === 'built' ? 'status-success' : 'status-neutral');
    const current = document.getElementById('revision-current');
    current.textContent = revisionAbbrev(revision.builder_commit);
    current.title = revision.builder_commit || 'Unknown Builder revision';
    const previous = document.getElementById('revision-last-built');
    previous.textContent = revisionAbbrev(revision.last_built_commit);
    previous.title = revision.last_built_commit || 'Previous builds did not record a Builder revision';
    document.getElementById('revision-version').textContent = revision.builder_version || '';
    document.getElementById('revision-built-time').textContent = revision.last_built_at ? fmtTime(revision.last_built_at) : 'No dated build record';
    const source = document.getElementById('revision-source');
    source.textContent = revision.source_commit ? revisionAbbrev(revision.source_commit) : 'Unavailable';
    source.title = revision.source_commit || 'No Git source revision available';
    document.getElementById('revision-source-status').textContent = revision.source_changed ? 'Not yet built into the live release' : revision.live_source_commit ? 'Matches published source' : 'No published source revision';
    document.getElementById('revision-detail').textContent = unbuilt ? 'No retained successful build used the running Builder revision; trigger a build to validate it.' : revision.status === 'built' ? 'A successful build with this Builder revision is recorded.' : 'Builder or previous build revision is unavailable; this is not proof of an unbuilt revision.';
  }

  function renderState(payload) {
    const state = payload.state || {};
    currentRelease.textContent = payload.current_release_id ? 'Live' : 'Not published';
    renderRevision(payload.revision, state);
    queueCount.textContent = (state.queue || []).length;
    queueOverview.textContent = (state.queue || []).length;
    buildCount.textContent = (state.builds || []).length;
    renderRunning(state.running || null, state.builds || []);
    renderBuilds(state);
    renderReleases(payload.releases || []);
    syncStatus.textContent = 'Live polling every 2s';
    updateFavicon(faviconState(state));
  }

  async function pollState() {
    if (pollInFlight) {
      return;
    }
    pollInFlight = true;
    try {
      const response = await fetch('/api/state', { headers: { 'Accept': 'application/json' }, cache: 'no-store' });
      if (!response.ok) {
        throw new Error('HTTP ' + response.status);
      }
      const payload = await response.json();
      renderState(payload);
    } catch (error) {
      syncStatus.textContent = 'Sync stalled: ' + error.message;
      updateFavicon('error');
    } finally {
      pollInFlight = false;
    }
  }

  window.addEventListener('hashchange', activateTabs);
  makeExpandableRows(document);
  activateTabs();
  updateFavicon('idle');
  pollState();
  window.setInterval(pollState, 2000);
</script>
</body>
</html>`
