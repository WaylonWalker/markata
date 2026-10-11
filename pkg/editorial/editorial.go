// Package editorial provides deterministic, isolated Git authoring for post
// editorial work. It records the resulting revisions in YAML frontmatter.
package editorial

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DigestVersion identifies the digest algorithm used in editorial records.
const DigestVersion = "markata-go-editorial-v1"

var roles = map[string]bool{"copy-editor": true, "internal-link-editor": true}
var uuidRE = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

const (
	phaseStarted          = "started"
	phaseContentCommitted = "content_committed"
	phaseMetadataPending  = "metadata_pending"
	phaseDone             = "done"
	actionReviewed        = "reviewed"
	actionEdited          = "edited"
)

// Record is one append-only editorial history event.
type Record struct {
	SchemaVersion int    `json:"schema_version" yaml:"schema_version"`
	EventID       string `json:"event_id" yaml:"event_id"`
	Author        string `json:"author" yaml:"author"`
	Role          string `json:"role" yaml:"role"`
	Action        string `json:"action" yaml:"action"`
	RecordedBy    string `json:"recorded_by" yaml:"recorded_by"`
	RecordedAt    string `json:"recorded_at" yaml:"recorded_at"`
	BaseCommit    string `json:"base_commit" yaml:"base_commit"`
	Commit        string `json:"commit" yaml:"commit"`
	Path          string `json:"path" yaml:"path"`
	ContentDigest string `json:"content_digest" yaml:"content_digest"`
	DigestVersion string `json:"digest_version" yaml:"digest_version"`
	FileURL       string `json:"file_url" yaml:"file_url"`
	DiffURL       string `json:"diff_url" yaml:"diff_url"`
	Summary       string `json:"summary" yaml:"summary"`
}

// Session stores the durable state of an editorial operation.
type Session struct {
	ID            string                  `json:"id"`
	Author        string                  `json:"author"`
	Role          string                  `json:"role"`
	DefaultAuthor string                  `json:"default_author,omitempty"`
	RecordedBy    string                  `json:"recorded_by,omitempty"`
	BaseCommit    string                  `json:"base_commit"`
	Phase         string                  `json:"phase"`
	Posts         map[string]PostSnapshot `json:"posts"`
	Commits       []string                `json:"commits"`
	Actions       map[string]string       `json:"actions,omitempty"`
	ContentCommit string                  `json:"content_commit,omitempty"`
	Summary       string                  `json:"summary,omitempty"`
	Pending       map[string]string       `json:"pending,omitempty"`
}

// PostSnapshot captures the committed baseline for one path.
type PostSnapshot struct {
	Before string `json:"before"`
	Digest string `json:"digest"`
}

// Options controls repository location and mutation behavior.
type Options struct {
	Root       string
	RecordedBy string
	// DefaultAuthor preserves the configured primary writer when a post has no
	// explicit author list. It is never inferred from the editor identity.
	DefaultAuthor string
	DryRun        bool
}

// normalizePaths accepts Markdown files anywhere under the repository. The
// caller may provide paths relative to Options.Root (the selected site dir).
func normalizePaths(root, base string, in []string) ([]string, error) {
	if len(in) == 0 {
		return nil, errors.New("select at least one Markdown post")
	}
	out := make([]string, 0, len(in))
	seen := map[string]bool{}
	for _, raw := range in {
		candidate := raw
		if !filepath.IsAbs(candidate) {
			candidate = filepath.Join(base, candidate)
		}
		abs, err := filepath.Abs(candidate)
		if err != nil {
			return nil, err
		}
		// Git resolves directory aliases (including macOS /var -> /private/var).
		// Resolve the parent only so a symlinked post itself remains rejected.
		parent, err := filepath.EvalSymlinks(filepath.Dir(abs))
		if err != nil {
			return nil, err
		}
		abs = filepath.Join(parent, filepath.Base(abs))
		rel, err := filepath.Rel(root, abs)
		if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
			return nil, fmt.Errorf("path must be inside the Git repository: %s", raw)
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, ".git/") || filepath.Ext(rel) != ".md" {
			return nil, fmt.Errorf("expected a tracked Markdown file: %s", raw)
		}
		info, err := os.Lstat(abs)
		if err != nil || info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("invalid post path: %s", raw)
		}
		if _, err = git(root, "ls-files", "--error-unmatch", "--", rel); err != nil {
			return nil, fmt.Errorf("post is not tracked: %s", rel)
		}
		if !seen[rel] {
			out = append(out, rel)
			seen[rel] = true
		}
	}
	sort.Strings(out)
	return out, nil
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("generate editorial session ID: %v", err))
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
func validate(author, role string) error {
	if strings.TrimSpace(author) == "" {
		return errors.New("author must be a stable, nonempty ID")
	}
	if !roles[role] {
		return fmt.Errorf("unsupported role %q (use copy-editor or internal-link-editor)", role)
	}
	return nil
}
func repoRoot(root string) (string, error) {
	if root == "" {
		root = "."
	}
	out, err := run(root, nil, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not a Git repository: %w", err)
	}
	return filepath.EvalSymlinks(strings.TrimSpace(out))
}
func run(dir string, env []string, args ...string) (string, error) {
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Env = append(os.Environ(), env...)
	b, e := c.CombinedOutput()
	if e != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(string(b)))
	}
	return string(b), nil
}
func git(root string, args ...string) (string, error) { return run(root, nil, args...) }
func gitNames(root string, args ...string) ([]string, error) {
	out, err := git(root, append(args, "-z")...)
	if err != nil {
		return nil, err
	}
	raw := strings.Split(out, "\x00")
	names := make([]string, 0, len(raw))
	for _, name := range raw {
		if name != "" {
			names = append(names, name)
		}
	}
	return names, nil
}
func clean(root string, ps []string) error {
	stagedArgs := append([]string{"diff", "--cached", "--name-only", "--"}, ps...)
	o, e := git(root, stagedArgs...)
	if e != nil {
		return e
	}
	if strings.TrimSpace(o) != "" {
		return errors.New("selected posts have staged changes; unstage them first")
	}
	statusArgs := append([]string{"status", "--porcelain", "--"}, ps...)
	o, e = git(root, statusArgs...)
	if e != nil {
		return e
	}
	if strings.TrimSpace(o) != "" {
		return errors.New("selected posts must be clean; commit owned changes first")
	}
	return nil
}
func rejectStaged(root string, ps []string) error {
	args := append([]string{"diff", "--cached", "--name-only", "--"}, ps...)
	out, err := git(root, args...)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "" {
		return errors.New("selected posts have staged changes; unstage them first")
	}
	return nil
}
func split(text string) (body string, node yaml.Node, err error) {
	lines := strings.SplitAfter(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(strings.TrimSuffix(lines[0], "\n")) != "---" {
		return "", yaml.Node{}, errors.New("post must have YAML frontmatter")
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(strings.TrimSuffix(lines[i], "\n")) == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return "", yaml.Node{}, errors.New("unclosed frontmatter")
	}
	frontText := strings.Join(lines[1:end], "")
	var n yaml.Node
	if e := yaml.Unmarshal([]byte(frontText), &n); e != nil {
		return "", yaml.Node{}, fmt.Errorf("invalid frontmatter: %w", e)
	}
	if len(n.Content) == 0 || n.Content[0].Kind != yaml.MappingNode {
		return "", yaml.Node{}, errors.New("frontmatter must be a mapping")
	}
	return strings.Join(lines[end+1:], ""), n, nil
}
func mapping(n *yaml.Node, key string) *yaml.Node {
	if n == nil || len(n.Content) == 0 {
		return nil
	}
	m := n
	if n.Kind == yaml.DocumentNode {
		m = n.Content[0]
	}
	if m.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}
func set(n *yaml.Node, key string, v any) error {
	m := n.Content[0]
	valueNode := func(value any) (*yaml.Node, error) {
		if node, ok := value.(*yaml.Node); ok {
			return node, nil
		}
		var encoded yaml.Node
		if err := encoded.Encode(value); err != nil {
			return nil, err
		}
		if len(encoded.Content) == 0 {
			return nil, errors.New("encoded YAML value is empty")
		}
		return &encoded, nil
	}
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			x, e := valueNode(v)
			if e != nil {
				return e
			}
			m.Content[i+1] = x
			return nil
		}
	}
	k := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}
	x, e := valueNode(v)
	if e != nil {
		return e
	}
	m.Content = append(m.Content, k, x)
	return nil
}
func rewrite(text string, fields map[string]any) (string, error) {
	body, n, e := split(text)
	if e != nil {
		return "", e
	}
	for k, v := range fields {
		if err := set(&n, k, v); err != nil {
			return "", err
		}
	}
	b, e := yaml.Marshal(n.Content[0])
	if e != nil {
		return "", e
	}
	return "---\n" + string(b) + "---\n" + body, nil
}
func digest(text string) (string, error) {
	body, n, e := split(text)
	if e != nil {
		return "", e
	}
	m := map[string]any{}
	if err := n.Content[0].Decode(&m); err != nil {
		return "", err
	}
	delete(m, "authors")
	delete(m, "author")
	delete(m, "editorial_history")
	b, e := json.Marshal(m)
	if e != nil {
		return "", e
	}
	h := sha256.Sum256(append(append(b, '\n'), []byte(body)...))
	return "sha256:" + hex.EncodeToString(h[:]), nil
}

// Start captures a clean baseline for selected Markdown files.
func Start(o Options, author, role string, selected []string) (Session, error) {
	if e := validate(author, role); e != nil {
		return Session{}, e
	}
	root, e := repoRoot(o.Root)
	if e != nil {
		return Session{}, e
	}
	baseDir := o.Root
	if baseDir == "" {
		baseDir = root
	}
	ps, e := normalizePaths(root, baseDir, selected)
	if e != nil {
		return Session{}, e
	}
	if err := clean(root, ps); err != nil {
		return Session{}, err
	}
	head, e := git(root, "rev-parse", "HEAD")
	if e != nil {
		return Session{}, e
	}
	s := Session{ID: newID(), Author: author, Role: role, DefaultAuthor: o.DefaultAuthor, RecordedBy: o.RecordedBy, BaseCommit: strings.TrimSpace(head), Phase: phaseStarted, Posts: map[string]PostSnapshot{}, Commits: []string{}}
	for _, p := range ps {
		b, e := git(root, "show", s.BaseCommit+":"+p)
		if e != nil {
			return Session{}, e
		}
		d, e := digest(b)
		if e != nil {
			return Session{}, fmt.Errorf("%s: %w", p, e)
		}
		s.Posts[p] = PostSnapshot{Before: b, Digest: d}
	}
	if !o.DryRun {
		if err := save(root, s); err != nil {
			return Session{}, err
		}
	}
	return s, nil
}
func sessionDir(root string) string {
	g, e := git(root, "rev-parse", "--git-dir")
	if e != nil {
		return filepath.Join(root, ".git", "markata-go-editorial")
	}
	if !filepath.IsAbs(strings.TrimSpace(g)) {
		g = filepath.Join(root, strings.TrimSpace(g))
	}
	return filepath.Join(g, "markata-go-editorial")
}
func save(root string, s Session) error {
	d := sessionDir(root)
	if e := os.MkdirAll(d, 0o700); e != nil {
		return e
	}
	b, e := json.MarshalIndent(s, "", "  ")
	if e != nil {
		return e
	}
	dest := filepath.Join(d, s.ID+".json")
	tmp, err := os.CreateTemp(d, ".session-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err = tmp.Chmod(0o600); err == nil {
		_, err = tmp.Write(append(b, '\n'))
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(name, dest)
}

// Load reads a durable session by UUID.
func Load(root, id string) (Session, error) {
	if !uuidRE.MatchString(strings.ToLower(id)) {
		return Session{}, errors.New("session must be a UUID")
	}
	b, e := os.ReadFile(filepath.Join(sessionDir(root), id+".json"))
	if e != nil {
		return Session{}, fmt.Errorf("load session: %w", e)
	}
	var s Session
	e = json.Unmarshal(b, &s)
	return s, e
}

func commit(root string, ps []string, message, author string, dry bool) (string, error) {
	if e := rejectStaged(root, ps); e != nil {
		return "", e
	}
	if dry {
		return "(dry-run)", nil
	}
	base, e := git(root, "rev-parse", "HEAD")
	if e != nil {
		return "", e
	}
	base = strings.TrimSpace(base)
	tmp, e := os.MkdirTemp("", "markata-editorial-index-")
	if e != nil {
		return "", e
	}
	defer os.RemoveAll(tmp)
	idx := filepath.Join(tmp, "index")
	env := []string{"GIT_INDEX_FILE=" + idx}
	if _, e = run(root, env, "read-tree", base); e != nil {
		return "", e
	}
	if _, e = run(root, env, append([]string{"add", "--"}, ps...)...); e != nil {
		return "", e
	}
	changed, e := run(root, env, "diff", "--cached", "--name-only", "-z")
	if e != nil {
		return "", e
	}
	if strings.Trim(changed, "\x00") == "" {
		return "", errors.New("selected posts have no changes")
	}
	if _, e = run(root, env, "commit", "-m", message, "-m", "Editorial-Author: "+author); e != nil {
		return "", e
	}
	sha, e := git(root, "rev-parse", "HEAD")
	if e != nil {
		return "", e
	}
	actualNames, err := gitNames(root, "diff-tree", "--no-commit-id", "--name-only", "-r", strings.TrimSpace(sha))
	if err != nil {
		return "", err
	}
	allowed := map[string]bool{}
	for _, p := range ps {
		allowed[p] = true
	}
	for _, p := range actualNames {
		if !allowed[p] {
			return "", fmt.Errorf("commit %s changed unexpected file %s", strings.TrimSpace(sha), p)
		}
	}
	_, e = git(root, append([]string{"reset", "-q", "HEAD", "--"}, ps...)...)
	return strings.TrimSpace(sha), e
}
func github(root string) (string, error) {
	u, e := git(root, "remote", "get-url", "origin")
	if e != nil {
		return "", errors.New("origin must be a GitHub repository URL")
	}
	u = strings.TrimSpace(u)
	if strings.HasPrefix(u, "git@github.com:") {
		u = "https://github.com/" + strings.TrimPrefix(u, "git@github.com:")
	}
	u = strings.TrimSuffix(u, ".git")
	parsed, err := url.Parse(u)
	if err != nil || parsed.Scheme != "https" || parsed.Host != "github.com" || strings.Trim(parsed.Path, "/") == "" || strings.Contains(parsed.Path, "..") {
		return "", errors.New("origin must be a GitHub repository URL")
	}
	return "https://github.com/" + strings.Trim(parsed.Path, "/"), nil
}
func credit(text, author, role, defaultAuthor string) (string, error) {
	_, n, e := split(text)
	if e != nil {
		return "", e
	}
	var authors *yaml.Node
	if v := mapping(&n, "authors"); v != nil {
		authors = v
		if authors.Kind != yaml.SequenceNode {
			return "", errors.New("authors must be a list")
		}
	} else {
		authors = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		if v := mapping(&n, "author"); v != nil {
			authors.Content = append(authors.Content, v)
		} else if defaultAuthor != "" {
			authors.Content = append(authors.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: defaultAuthor})
		} else {
			return "", errors.New("post has no primary author; configure authors/default author or pass --primary-author")
		}
	}
	for _, a := range authors.Content {
		if a.Kind == yaml.MappingNode && mapping(a, "id") != nil && mapping(a, "id").Value == author && mapping(a, "role") != nil && strings.ReplaceAll(mapping(a, "role").Value, " ", "-") == role {
			return text, nil
		}
	}
	a := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{{Kind: yaml.ScalarNode, Tag: "!!str", Value: "id"}, {Kind: yaml.ScalarNode, Tag: "!!str", Value: author}, {Kind: yaml.ScalarNode, Tag: "!!str", Value: "role"}, {Kind: yaml.ScalarNode, Tag: "!!str", Value: strings.ReplaceAll(role, "-", " ")}}}
	authors.Content = append(authors.Content, a)
	return rewrite(text, map[string]any{"authors": authors})
}
func sessionRecorder(s Session) string {
	if s.RecordedBy != "" {
		return s.RecordedBy
	}
	return s.Author
}

func makeRecord(root string, s Session, p, sha, action, summary string) (Record, error) {
	base, e := github(root)
	if e != nil {
		return Record{}, e
	}
	b, e := os.ReadFile(filepath.Join(root, p))
	if e != nil {
		return Record{}, e
	}
	d, e := digest(string(b))
	if e != nil {
		return Record{}, e
	}
	anchor := sha256.Sum256([]byte(p))
	encodedPath := ""
	encodedParts := make([]string, 0)
	for _, part := range strings.Split(p, "/") {
		encodedParts = append(encodedParts, url.PathEscape(part))
	}
	encodedPath = strings.Join(encodedParts, "/")
	return Record{SchemaVersion: 1, EventID: s.ID + ":" + p, Author: s.Author, Role: s.Role, Action: action, RecordedBy: sessionRecorder(s), RecordedAt: time.Now().UTC().Format(time.RFC3339), BaseCommit: s.BaseCommit, Commit: sha, Path: p, ContentDigest: d, DigestVersion: DigestVersion, FileURL: base + "/blob/" + sha + "/" + encodedPath, DiffURL: base + "/commit/" + sha + "#diff-" + hex.EncodeToString(anchor[:]), Summary: summary}, nil
}
func recordNode(r Record) (*yaml.Node, error) {
	var n yaml.Node
	if err := n.Encode(r); err != nil {
		return nil, err
	}
	if len(n.Content) == 0 {
		return nil, errors.New("encoded record is empty")
	}
	return &n, nil
}

// sessionCommitAtHead recognizes a completed commit when a process stops before
// persisting the next session phase. The session marker, exact parent, and
// selected paths must all match; unrelated HEAD movement is never adopted.
func sessionCommitAtHead(root string, s Session, parent, stage string, ps []string) (sha string, recognized bool, err error) {
	sha, err = git(root, "rev-parse", "HEAD")
	if err != nil {
		return "", false, err
	}
	sha = strings.TrimSpace(sha)
	if sha == parent {
		return sha, false, nil
	}
	actualParent, err := git(root, "rev-parse", "HEAD^")
	if err != nil || strings.TrimSpace(actualParent) != parent {
		return sha, false, nil
	}
	message, err := git(root, "show", "-s", "--format=%B", "HEAD")
	if err != nil {
		return "", false, err
	}
	if !strings.Contains(message, "Editorial-Session: "+s.ID+"\n") ||
		!strings.Contains(message, "Editorial-Phase: "+stage+"\n") {
		return sha, false, nil
	}
	changed, err := gitNames(root, "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return "", false, err
	}
	allowed := make(map[string]bool, len(ps))
	for _, p := range ps {
		allowed[p] = true
	}
	if len(changed) == 0 {
		return sha, false, nil
	}
	for _, p := range changed {
		if !allowed[p] {
			return sha, false, nil
		}
	}
	for _, p := range ps {
		committed, err := git(root, "show", sha+":"+p)
		if err != nil {
			return "", false, err
		}
		current, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			return "", false, err
		}
		if string(current) != committed {
			return "", false, fmt.Errorf("selected post changed after the %s commit: %s", stage, p)
		}
	}
	return sha, true, nil
}

// reconcileCommittedIndex handles a crash between Git advancing HEAD with
// our isolated index and the subsequent reset of selected paths in the user's
// index. Only reset entries that still match the exact old HEAD; actual user
// staging remains protected by rejectStaged.
func reconcileCommittedIndex(root string, s Session, ps []string) error {
	parent, stage := "", ""
	switch s.Phase {
	case phaseStarted:
		parent, stage = s.BaseCommit, "content"
	case phaseMetadataPending:
		parent, stage = s.ContentCommit, "metadata"
	default:
		return nil
	}
	_, recognized, err := sessionCommitAtHead(root, s, parent, stage, ps)
	if err != nil || !recognized {
		return err
	}
	for _, p := range ps {
		indexed, err := git(root, "show", ":"+p)
		if err != nil {
			return err
		}
		before, err := git(root, "show", parent+":"+p)
		if err != nil {
			return err
		}
		if indexed != before {
			// The caller may have staged something new; do not touch it.
			return nil
		}
	}
	_, err = git(root, append([]string{"reset", "-q", "HEAD", "--"}, ps...)...)
	return err
}

type preparedEditorialPost struct {
	Action string
	Before string
	After  string
}

func prepareStartedPost(root string, s Session, path string) (preparedEditorialPost, error) {
	b, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		return preparedEditorialPost{}, err
	}
	d, err := digest(string(b))
	if err != nil {
		return preparedEditorialPost{}, err
	}
	action := actionReviewed
	if d != s.Posts[path].Digest {
		action = actionEdited
	}
	_, before, err := split(s.Posts[path].Before)
	if err != nil {
		return preparedEditorialPost{}, fmt.Errorf("%s: invalid frontmatter: %w", path, err)
	}
	_, after, err := split(string(b))
	if err != nil {
		return preparedEditorialPost{}, fmt.Errorf("%s: invalid frontmatter: %w", path, err)
	}
	if history := mapping(&after, "editorial_history"); history != nil && history.Kind != yaml.SequenceNode {
		return preparedEditorialPost{}, fmt.Errorf("%s: editorial_history must be a list", path)
	}
	oldBytes, err := yaml.Marshal(mapping(&before, "editorial_history"))
	if err != nil {
		return preparedEditorialPost{}, err
	}
	newBytes, err := yaml.Marshal(mapping(&after, "editorial_history"))
	if err != nil {
		return preparedEditorialPost{}, err
	}
	if !bytes.Equal(oldBytes, newBytes) {
		return preparedEditorialPost{}, fmt.Errorf("session must preserve prior editorial history: %s", path)
	}
	x, err := credit(string(b), s.Author, s.Role, s.DefaultAuthor)
	if err != nil {
		return preparedEditorialPost{}, err
	}
	return preparedEditorialPost{Action: action, Before: string(b), After: x}, nil
}

func recoverStartedCommit(o Options, root string, s Session, ps []string, summary string) (Session, error) {
	sha, recognized, err := sessionCommitAtHead(root, s, s.BaseCommit, "content", ps)
	if err != nil {
		return Session{}, err
	}
	if !recognized {
		return Session{}, errors.New("HEAD changed during session; inspect and start a new session")
	}
	s.Actions = make(map[string]string, len(ps))
	for _, p := range ps {
		committed, err := git(root, "show", sha+":"+p)
		if err != nil {
			return Session{}, err
		}
		d, err := digest(committed)
		if err != nil {
			return Session{}, err
		}
		s.Actions[p] = actionReviewed
		if d != s.Posts[p].Digest {
			s.Actions[p] = actionEdited
		}
	}
	s.ContentCommit = sha
	s.Commits = append(s.Commits, sha)
	s.Summary = summary
	s.Phase = phaseContentCommitted
	if !o.DryRun {
		if err := save(root, s); err != nil {
			return Session{}, err
		}
	}
	return s, nil
}

func finishStarted(o Options, root string, s Session, ps []string, summary string) (Session, error) {
	head, err := git(root, "rev-parse", "HEAD")
	if err != nil {
		return Session{}, err
	}
	head = strings.TrimSpace(head)
	if head != s.BaseCommit {
		return recoverStartedCommit(o, root, s, ps, summary)
	}

	// Prepare *all* posts and validate the entire selection before changing
	// any files. A malformed later post must not leave earlier files rewritten.
	prepared := make(map[string]preparedEditorialPost, len(ps))
	s.Actions = make(map[string]string, len(ps))
	for _, p := range ps {
		pending, err := prepareStartedPost(root, s, p)
		if err != nil {
			return Session{}, err
		}
		prepared[p] = pending
		s.Actions[p] = pending.Action
	}
	if !o.DryRun {
		for _, p := range ps {
			current, err := os.ReadFile(filepath.Join(root, p))
			if err != nil {
				return Session{}, err
			}
			if string(current) != prepared[p].Before {
				return Session{}, fmt.Errorf("selected post changed during validation: %s", p)
			}
		}
		for _, p := range ps {
			if prepared[p].After != prepared[p].Before {
				if err := os.WriteFile(filepath.Join(root, p), []byte(prepared[p].After), 0o600); err != nil {
					return Session{}, err
				}
			}
		}
	}

	sha := s.BaseCommit
	if !o.DryRun {
		diffArgs := append([]string{"diff", "--name-only", "--"}, ps...)
		out, err := git(root, diffArgs...)
		if err != nil {
			return Session{}, err
		}
		if strings.TrimSpace(out) != "" {
			message := "editorial: " + s.Role + " by " + s.Author +
				"\n\nEditorial-Session: " + s.ID + "\nEditorial-Phase: content"
			sha, err = commit(root, ps, message, s.Author, false)
			if err != nil {
				return Session{}, err
			}
			s.Commits = append(s.Commits, sha)
		}
	}
	s.ContentCommit = sha
	s.Summary = summary
	s.Phase = phaseContentCommitted
	if !o.DryRun {
		if err := save(root, s); err != nil {
			return Session{}, err
		}
	}
	return s, nil
}

func finishContent(o Options, root string, s Session, ps []string) (Session, error) {
	currentHead, err := git(root, "rev-parse", "HEAD")
	if err != nil {
		return Session{}, err
	}
	if strings.TrimSpace(currentHead) != s.ContentCommit {
		return Session{}, errors.New("HEAD changed after content commit; restore the session before retrying")
	}
	if o.DryRun {
		return s, nil
	}
	if len(s.Pending) == 0 {
		diffArgs := append([]string{"diff", "--name-only", "--"}, ps...)
		changed, err := git(root, diffArgs...)
		if err != nil {
			return Session{}, err
		}
		if strings.TrimSpace(changed) != "" {
			return Session{}, errors.New("selected posts changed after the content commit; restore the session before retrying")
		}
		s.Pending, err = buildPending(root, s, ps)
		if err != nil {
			return Session{}, err
		}
		s.Phase = phaseMetadataPending
		// Journal the intended bytes before writing them so a partial write can
		// be safely completed on retry without overwriting intervening edits.
		if err := save(root, s); err != nil {
			return Session{}, err
		}
	}
	// Verify the complete selection before writing any file. A file is safe
	// to complete only when it is the content commit or our own pending bytes.
	for _, p := range ps {
		current, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			return Session{}, err
		}
		if string(current) == s.Pending[p] {
			continue
		}
		baseline, err := git(root, "show", s.ContentCommit+":"+p)
		if err != nil {
			return Session{}, err
		}
		if string(current) != baseline {
			return Session{}, fmt.Errorf("attribution interrupted or post changed: %s", p)
		}
	}
	for _, p := range ps {
		current, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			return Session{}, err
		}
		if string(current) != s.Pending[p] {
			if err := os.WriteFile(filepath.Join(root, p), []byte(s.Pending[p]), 0o600); err != nil {
				return Session{}, err
			}
		}
	}
	s.Phase = phaseMetadataPending
	return s, nil
}

func buildPending(root string, s Session, ps []string) (map[string]string, error) {
	pending := make(map[string]string, len(ps))
	for _, p := range ps {
		b, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			return nil, err
		}
		r, err := makeRecord(root, s, p, s.ContentCommit, s.Actions[p], s.Summary)
		if err != nil {
			return nil, err
		}
		_, n, err := split(string(b))
		if err != nil {
			return nil, err
		}
		history := mapping(&n, "editorial_history")
		if history == nil {
			history = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		}
		if history.Kind != yaml.SequenceNode {
			return nil, errors.New("editorial_history must be a list")
		}
		rn, err := recordNode(r)
		if err != nil {
			return nil, err
		}
		history.Content = append(history.Content, rn)
		x, err := rewrite(string(b), map[string]any{"editorial_history": history})
		if err != nil {
			return nil, err
		}
		pending[p] = x
	}
	return pending, nil
}

func finishMetadata(root string, s Session, ps []string) (Session, error) {
	currentHead, err := git(root, "rev-parse", "HEAD")
	if err != nil {
		return Session{}, err
	}
	if strings.TrimSpace(currentHead) != s.ContentCommit {
		sha, recognized, err := sessionCommitAtHead(root, s, s.ContentCommit, "metadata", ps)
		if err != nil {
			return Session{}, err
		}
		if !recognized {
			return Session{}, errors.New("HEAD changed after content commit; restore the session before retrying")
		}
		for _, p := range ps {
			committed, err := git(root, "show", sha+":"+p)
			if err != nil {
				return Session{}, err
			}
			if committed != s.Pending[p] {
				return Session{}, fmt.Errorf("metadata commit does not match pending attribution: %s", p)
			}
		}
		s.Commits = append(s.Commits, sha)
		s.Phase = phaseDone
		if err := save(root, s); err != nil {
			return Session{}, err
		}
		return s, nil
	}
	for _, p := range ps {
		actual, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			return Session{}, err
		}
		if string(actual) != s.Pending[p] {
			return Session{}, fmt.Errorf("attribution interrupted or post changed: %s", p)
		}
	}
	message := "editorial: record " + s.Role + " revision links" +
		"\n\nEditorial-Session: " + s.ID + "\nEditorial-Phase: metadata"
	sha, err := commit(root, ps, message, sessionRecorder(s), false)
	if err != nil {
		return Session{}, err
	}
	s.Commits = append(s.Commits, sha)
	s.Phase = phaseDone
	if err := save(root, s); err != nil {
		return Session{}, err
	}
	return s, nil
}

// Finish commits session edits and appends revision links to frontmatter.
func Finish(o Options, id, summary string) (Session, error) {
	root, err := repoRoot(o.Root)
	if err != nil {
		return Session{}, err
	}
	s, err := Load(root, id)
	if err != nil {
		return Session{}, err
	}
	if strings.TrimSpace(summary) == "" {
		return Session{}, errors.New("summary must not be empty")
	}
	if s.Phase == phaseDone {
		return s, nil
	}
	ps := make([]string, 0, len(s.Posts))
	for p := range s.Posts {
		ps = append(ps, p)
	}
	sort.Strings(ps)
	if err := reconcileCommittedIndex(root, s, ps); err != nil {
		return Session{}, err
	}
	if err := rejectStaged(root, ps); err != nil {
		return Session{}, err
	}
	if _, err = github(root); err != nil {
		return Session{}, err
	}
	s.Summary = summary
	if s.Phase == phaseStarted {
		s, err = finishStarted(o, root, s, ps, summary)
		if err != nil {
			return Session{}, err
		}
	}
	if s.Phase == phaseContentCommitted || s.Phase == phaseMetadataPending {
		return finishAfterContent(o, root, s, ps)
	}
	return s, nil
}

func finishAfterContent(o Options, root string, s Session, ps []string) (Session, error) {
	head, err := git(root, "rev-parse", "HEAD")
	if err != nil {
		return Session{}, err
	}
	// A successful metadata commit might precede the final session save.
	// Recognize it without rewriting any pending files.
	if s.Phase == phaseMetadataPending && strings.TrimSpace(head) != s.ContentCommit {
		if o.DryRun {
			return s, nil
		}
		return finishMetadata(root, s, ps)
	}
	s, err = finishContent(o, root, s, ps)
	if err != nil {
		return Session{}, err
	}
	if s.Phase == phaseMetadataPending && !o.DryRun {
		return finishMetadata(root, s, ps)
	}
	return s, nil
}

// Commit attributes explicitly owned current changes in selected files.
func Commit(o Options, author, role, summary string, selected []string) (Session, error) {
	if e := validate(author, role); e != nil {
		return Session{}, e
	}
	root, e := repoRoot(o.Root)
	if e != nil {
		return Session{}, e
	}
	baseDir := o.Root
	if baseDir == "" {
		baseDir = root
	}
	ps, e := normalizePaths(root, baseDir, selected)
	if e != nil {
		return Session{}, e
	}
	if strings.TrimSpace(summary) == "" {
		return Session{}, errors.New("summary must not be empty")
	}
	base, e := git(root, "rev-parse", "HEAD")
	if e != nil {
		return Session{}, e
	}
	if err := rejectStaged(root, ps); err != nil {
		return Session{}, err
	}
	s := Session{ID: newID(), Author: author, Role: role, DefaultAuthor: o.DefaultAuthor, RecordedBy: o.RecordedBy, BaseCommit: strings.TrimSpace(base), Phase: phaseStarted, Posts: map[string]PostSnapshot{}, Commits: []string{}, Summary: summary}
	for _, p := range ps {
		b, e := git(root, "show", strings.TrimSpace(base)+":"+p)
		if e != nil {
			return Session{}, e
		}
		d, e := digest(b)
		if e != nil {
			return Session{}, e
		}
		s.Posts[p] = PostSnapshot{Before: b, Digest: d}
	}
	if o.DryRun {
		return finishStarted(o, root, s, ps, summary)
	}
	if err := save(root, s); err != nil {
		return Session{}, err
	}
	return Finish(o, s.ID, summary)
}
