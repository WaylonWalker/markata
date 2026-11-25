package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	gitPublisherRemotePath = "/v1/remote-head"
	gitPublisherPushPath   = "/v1/push"
	gitPublisherMaxBody    = 4096
)

var gitPublisherBranchPattern = regexp.MustCompile(`^refs/heads/[A-Za-z0-9][A-Za-z0-9._/-]*$`)
var gitPublisherRepoPartPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

// GitPublisher provides the narrow remote operations used by hosted editor
// publication. Implementations must not accept a caller-selected remote.
type GitPublisher interface {
	RemoteHead(context.Context, string, string) (string, error)
	Push(context.Context, string, string) error
}

type unixGitPublisher struct {
	socket string
	client *http.Client
}

// NewUnixGitPublisher connects the workbench to a Git publisher over a Unix
// domain socket. The socket path is local deployment configuration.
func NewUnixGitPublisher(socket string) GitPublisher {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", socket)
		},
	}
	return &unixGitPublisher{
		socket: socket,
		client: &http.Client{Transport: transport, Timeout: 35 * time.Second},
	}
}

type gitPublisherRemoteResponse struct {
	Head string `json:"head"`
}

type gitPublisherPushRequest struct {
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}

func (p *unixGitPublisher) RemoteHead(ctx context.Context, branch, head string) (string, error) {
	if p.socket == "" {
		return "", errors.New("git publisher socket is not configured")
	}
	endpoint := "http://git-publisher" + gitPublisherRemotePath + "?branch=" + url.QueryEscape(branch) + "&head=" + url.QueryEscape(head)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return "", err
	}
	var response gitPublisherRemoteResponse
	if err := p.do(req, &response); err != nil {
		return "", err
	}
	if !gitHashPattern.MatchString(response.Head) {
		return "", errors.New("git publisher returned an invalid remote revision")
	}
	return response.Head, nil
}

func (p *unixGitPublisher) Push(ctx context.Context, branch, commit string) error {
	if p.socket == "" {
		return errors.New("git publisher socket is not configured")
	}
	body, err := json.Marshal(gitPublisherPushRequest{Branch: branch, Commit: commit})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://git-publisher"+gitPublisherPushPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	var response struct {
		Status string `json:"status"`
		Commit string `json:"commit"`
	}
	if err := p.do(req, &response); err != nil {
		return err
	}
	if response.Status != "pushed" || response.Commit != commit {
		return errors.New("git publisher did not confirm the requested commit")
	}
	return nil
}

func (p *unixGitPublisher) do(req *http.Request, dst any) error {
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("contact local Git publisher: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, gitPublisherMaxBody+1))
	if err != nil {
		return fmt.Errorf("read local Git publisher response: %w", err)
	}
	if len(body) > gitPublisherMaxBody {
		return errors.New("local Git publisher response exceeded its limit")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = http.StatusText(resp.StatusCode)
		}
		return fmt.Errorf("git publisher rejected request (%d): %s", resp.StatusCode, message)
	}
	if dst == nil {
		return nil
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("decode local Git publisher response: %w", err)
	}
	return nil
}

// GitPublisherServerConfig fixes the publisher to one checkout, GitHub
// repository, branch, SSH key, and verified known_hosts file.
type GitPublisherServerConfig struct {
	SourceDir  string
	SocketPath string
	OriginURL  string
	Branch     string
	PrivateKey string
	KnownHosts string
}

// ServeGitPublisher serves the credential-bearing Git operation broker on a
// Unix socket. It never accepts a remote, path, or refspec from a client.
func ServeGitPublisher(ctx context.Context, config GitPublisherServerConfig) error {
	service, err := newGitPublisherService(config)
	if err != nil {
		return err
	}
	if config.SocketPath == "" {
		return errors.New("git publisher socket path is required")
	}
	socket, err := filepath.Abs(config.SocketPath)
	if err != nil {
		return fmt.Errorf("resolve Git publisher socket: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(socket), 0o770); err != nil {
		return fmt.Errorf("create Git publisher socket directory: %w", err)
	}
	if info, err := os.Lstat(socket); err == nil {
		if info.Mode()&os.ModeSocket == 0 {
			return fmt.Errorf("refusing to replace non-socket path %s", socket)
		}
		if err := os.Remove(socket); err != nil {
			return fmt.Errorf("remove stale Git publisher socket: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect Git publisher socket: %w", err)
	}
	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "unix", socket)
	if err != nil {
		return fmt.Errorf("listen on Git publisher socket: %w", err)
	}
	defer func() {
		_ = listener.Close()
		_ = os.Remove(socket)
	}()
	if err := os.Chmod(socket, 0o660); err != nil {
		return fmt.Errorf("set Git publisher socket permissions: %w", err)
	}
	server := &http.Server{
		Handler:           service.handler(),
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    8192,
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("git publisher shutdown: %v", err)
		}
	}()
	err = server.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) || errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}

type gitPublisherService struct {
	root       string
	originURL  string
	pushURL    string
	branch     string
	privateKey string
	knownHosts string
	run        func(context.Context, ...string) (string, error)
}

func newGitPublisherService(config GitPublisherServerConfig) (*gitPublisherService, error) {
	if config.SourceDir == "" || config.PrivateKey == "" || config.KnownHosts == "" {
		return nil, errors.New("git publisher source, private key, and known_hosts paths are required")
	}
	if err := validatePublisherFile(config.PrivateKey, true); err != nil {
		return nil, fmt.Errorf("invalid Git publisher private key: %w", err)
	}
	if err := validatePublisherFile(config.KnownHosts, false); err != nil {
		return nil, fmt.Errorf("invalid Git publisher known_hosts file: %w", err)
	}
	pushURL, err := validateGitPublisherOrigin(config.OriginURL)
	if err != nil {
		return nil, err
	}
	if !gitPublisherBranchPattern.MatchString(config.Branch) || strings.Contains(config.Branch, "..") || strings.Contains(config.Branch, "//") || strings.HasSuffix(config.Branch, ".") {
		return nil, errors.New("git publisher branch must be a valid fully qualified branch")
	}
	root, err := filepath.Abs(config.SourceDir)
	if err != nil {
		return nil, fmt.Errorf("resolve Git publisher source: %w", err)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, fmt.Errorf("resolve Git publisher checkout: %w", err)
	}
	return &gitPublisherService{
		root: root, originURL: config.OriginURL, pushURL: pushURL,
		branch: config.Branch, privateKey: config.PrivateKey, knownHosts: config.KnownHosts,
	}, nil
}

func validatePublisherFile(path string, ownerOnly bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return errors.New("must be a regular file")
	}
	if info.Size() == 0 {
		return errors.New("must not be empty")
	}
	if ownerOnly && info.Mode().Perm()&0o077 != 0 {
		return errors.New("private key permissions must be owner-only for OpenSSH")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return file.Close()
}

func validateGitPublisherOrigin(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != secureScheme || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("git publisher origin must be an HTTPS GitHub repository URL")
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	owner, repo := "", ""
	if len(parts) == 2 {
		owner, repo = parts[0], strings.TrimSuffix(parts[1], ".git")
	}
	if len(parts) != 2 || owner == "." || owner == ".." || repo == "." || repo == ".." || !gitPublisherRepoPartPattern.MatchString(owner) || !gitPublisherRepoPartPattern.MatchString(repo) {
		return "", errors.New("git publisher origin must name one GitHub owner and repository")
	}
	return "ssh://git@github.com/" + strings.TrimPrefix(u.EscapedPath(), "/"), nil
}

func (s *gitPublisherService) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(gitPublisherRemotePath, s.handleRemoteHead)
	mux.HandleFunc(gitPublisherPushPath, s.handlePush)
	return mux
}

func (s *gitPublisherService) handleRemoteHead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return
	}
	branch, head := r.URL.Query().Get("branch"), r.URL.Query().Get("head")
	if branch != s.branch || !gitHashPattern.MatchString(head) {
		http.Error(w, "requested branch or revision is invalid", http.StatusUnprocessableEntity)
		return
	}
	localHead, localBranch, err := s.localRevision(r.Context())
	if err != nil || localHead != head || localBranch != branch {
		http.Error(w, "checkout changed; review Git status again", http.StatusConflict)
		return
	}
	if err := s.validateOrigin(r.Context()); err != nil {
		http.Error(w, "configured Git origin does not match the publisher allowlist", http.StatusUnprocessableEntity)
		return
	}
	remote, err := s.remoteHead(r.Context())
	if err != nil {
		http.Error(w, "cannot verify remote Git state", http.StatusServiceUnavailable)
		return
	}
	writePublisherJSON(w, gitPublisherRemoteResponse{Head: remote})
}

func (s *gitPublisherService) handlePush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, gitPublisherMaxBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request gitPublisherPushRequest
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid push request", http.StatusBadRequest)
		return
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		http.Error(w, "unexpected additional request data", http.StatusBadRequest)
		return
	}
	if request.Branch != s.branch || !gitHashPattern.MatchString(request.Commit) {
		http.Error(w, "requested branch or commit is invalid", http.StatusUnprocessableEntity)
		return
	}
	if err := s.push(r.Context(), request.Branch, request.Commit); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	writePublisherJSON(w, map[string]string{"status": "pushed", "commit": request.Commit})
}

func (s *gitPublisherService) push(ctx context.Context, branch, commit string) error {
	parent, err := s.validatePushCommit(ctx, branch, commit)
	if err != nil {
		return err
	}
	remote, err := s.remoteHead(ctx)
	if err != nil || remote != parent {
		return errors.New("remote branch changed; sync and review again")
	}
	if _, err := s.runGit(ctx, "push", "--porcelain", s.pushURL, commit+":"+s.branch); err != nil {
		return errors.New("remote push failed; local commit is retained")
	}
	remote, err = s.remoteHead(ctx)
	if err != nil || remote != commit {
		return errors.New("push completed but the remote revision could not be confirmed")
	}
	return nil
}

func (s *gitPublisherService) validatePushCommit(ctx context.Context, branch, commit string) (string, error) {
	head, localBranch, err := s.localRevision(ctx)
	if err != nil || head != commit || localBranch != branch {
		return "", errors.New("checkout changed; review Git status again")
	}
	if err := s.validateOrigin(ctx); err != nil {
		return "", errors.New("configured Git origin does not match the publisher allowlist")
	}
	parents, err := s.runGit(ctx, "rev-list", "--parents", "-n", "1", commit)
	if err != nil {
		return "", errors.New("cannot verify publish commit")
	}
	parentFields := strings.Fields(parents)
	if len(parentFields) != 2 || parentFields[0] != commit || !gitHashPattern.MatchString(parentFields[1]) {
		return "", errors.New("publish commit must have exactly one parent")
	}
	changed, err := s.runGit(ctx, "diff-tree", "--no-commit-id", "--name-only", "-r", commit)
	if err != nil {
		return "", errors.New("cannot inspect publish commit")
	}
	changed = strings.TrimSuffix(changed, "\n")
	if changed == "" || strings.ContainsAny(changed, "\r\n\t\x00\"") || !filepath.IsLocal(filepath.FromSlash(changed)) || !strings.EqualFold(filepath.Ext(changed), ".md") {
		return "", errors.New("publish commit must change exactly one Markdown file")
	}
	treeEntry, err := s.runGit(ctx, "ls-tree", commit, "--", changed)
	if err != nil {
		return "", errors.New("cannot verify published file type")
	}
	entryParts := strings.SplitN(treeEntry, "\t", 2)
	entryFields := strings.Fields(entryParts[0])
	if len(entryParts) != 2 || entryParts[1] != changed || len(entryFields) != 3 || entryFields[0] != "100644" || entryFields[1] != "blob" {
		return "", errors.New("publish commit must contain one regular Markdown file")
	}
	return parentFields[1], nil
}

func (s *gitPublisherService) localRevision(ctx context.Context) (head, branch string, err error) {
	head, err = s.runGit(ctx, "rev-parse", "HEAD")
	if err != nil || !gitHashPattern.MatchString(head) {
		return "", "", errors.New("cannot read checkout revision")
	}
	branch, err = s.runGit(ctx, "symbolic-ref", "-q", "HEAD")
	if err != nil || !gitPublisherBranchPattern.MatchString(branch) {
		return "", "", errors.New("checkout must use a named branch")
	}
	return head, branch, nil
}

func (s *gitPublisherService) validateOrigin(ctx context.Context) error {
	// Read the raw fetch URL without Git's url.insteadOf rewriting. Network
	// operations use the fixed SSH URL derived from the configured allowlist,
	// never a remote URL or pushurl from the shared checkout config.
	fetchURL, err := s.runGit(ctx, "config", "--local", "--get", "remote.origin.url")
	if err != nil || fetchURL != s.originURL {
		return errors.New("fetch remote is not allowlisted")
	}
	return nil
}

func (s *gitPublisherService) remoteHead(ctx context.Context) (string, error) {
	output, err := s.runGit(ctx, "ls-remote", "--heads", s.pushURL, s.branch)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(output)
	if len(fields) != 2 || fields[1] != s.branch || !gitHashPattern.MatchString(fields[0]) {
		return "", errors.New("remote branch is missing or ambiguous")
	}
	return fields[0], nil
}

func (s *gitPublisherService) runGit(ctx context.Context, args ...string) (string, error) {
	if s.run != nil {
		return s.run(ctx, args...)
	}
	// #nosec G204 G702 -- Git subcommands and arguments are fixed or validated above; no shell is used.
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", s.root, "-c", "core.hooksPath=/dev/null"}, args...)...)
	cmd.Env = gitPublisherEnvironment(s.root, s.privateKey, s.knownHosts)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git operation failed: %w", err)
	}
	if len(output) > 2<<20 {
		return "", errors.New("git response exceeds limit")
	}
	return strings.TrimSpace(string(output)), nil
}

func gitPublisherEnvironment(root, privateKey, knownHosts string) []string {
	env := make([]string, 0, len(os.Environ())+9)
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		if strings.HasPrefix(key, "GIT_") || strings.HasPrefix(key, "SSH_") {
			continue
		}
		env = append(env, item)
	}
	sshCommand := "ssh -i " + shellQuote(privateKey) + " -o BatchMode=yes -o ConnectTimeout=10 -o IdentitiesOnly=yes -o IdentityAgent=none -o StrictHostKeyChecking=yes -o UserKnownHostsFile=" + shellQuote(knownHosts) + " -o GlobalKnownHostsFile=/dev/null"
	return append(env,
		"GIT_TERMINAL_PROMPT=0",
		"GIT_OPTIONAL_LOCKS=0",
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_COUNT=1",
		"GIT_CONFIG_KEY_0=safe.directory",
		"GIT_CONFIG_VALUE_0="+root,
		"GIT_SSH_COMMAND="+sshCommand,
	)
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func writePublisherJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("git publisher response encoding failed: %v", err)
	}
}
