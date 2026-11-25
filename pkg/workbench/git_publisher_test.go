package workbench

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	publisherTestOrigin = "https://github.com/owner/site.git"
	publisherTestPush   = "ssh://git@github.com/owner/site.git"
	publisherTestBranch = "refs/heads/main"
)

func TestValidateGitPublisherOrigin(t *testing.T) {
	tests := []struct {
		name    string
		origin  string
		want    string
		wantErr bool
	}{
		{name: "github repository", origin: publisherTestOrigin, want: publisherTestPush},
		{name: "wrong scheme", origin: "http://github.com/owner/site", wantErr: true},
		{name: "wrong host", origin: "https://example.com/owner/site", wantErr: true},
		{name: "embedded user", origin: "https://user@github.com/owner/site", wantErr: true},
		{name: "extra path", origin: "https://github.com/owner/site/tree/main", wantErr: true},
		{name: "query", origin: "https://github.com/owner/site?redirect=elsewhere", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateGitPublisherOrigin(tt.origin)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateGitPublisherOrigin(%q) error = %v, wantErr %v", tt.origin, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("validateGitPublisherOrigin(%q) = %q, want %q", tt.origin, got, tt.want)
			}
		})
	}
}

func TestGitPublisherPush_RequiresSingleMarkdownFastForward(t *testing.T) {
	commit := strings.Repeat("a", 40)
	parent := strings.Repeat("b", 40)

	t.Run("pushes and confirms exact commit", func(t *testing.T) {
		var calls []string
		remoteReads := 0
		service := publisherTestService(func(_ context.Context, args ...string) (string, error) {
			call := strings.Join(args, " ")
			calls = append(calls, call)
			switch call {
			case "rev-parse HEAD":
				return commit, nil
			case "symbolic-ref -q HEAD":
				return publisherTestBranch, nil
			case "config --local --get remote.origin.url":
				return publisherTestOrigin, nil
			case "rev-list --parents -n 1 " + commit:
				return commit + " " + parent, nil
			case "diff-tree --no-commit-id --name-only -r " + commit:
				return "pages/post.md", nil
			case "ls-tree " + commit + " -- pages/post.md":
				return "100644 blob " + strings.Repeat("c", 40) + "\tpages/post.md", nil
			case "ls-remote --heads " + publisherTestPush + " " + publisherTestBranch:
				remoteReads++
				if remoteReads == 1 {
					return parent + "\t" + publisherTestBranch, nil
				}
				return commit + "\t" + publisherTestBranch, nil
			case "push --porcelain " + publisherTestPush + " " + commit + ":" + publisherTestBranch:
				return "", nil
			default:
				return "", errors.New("unexpected Git command: " + call)
			}
		})

		if err := service.push(context.Background(), publisherTestBranch, commit); err != nil {
			t.Fatalf("push() error = %v", err)
		}
		if remoteReads != 2 {
			t.Fatalf("remote branch checked %d times, want before and after push", remoteReads)
		}
		if calls[len(calls)-2] != "push --porcelain "+publisherTestPush+" "+commit+":"+publisherTestBranch {
			t.Fatalf("push command was %q", calls[len(calls)-2])
		}
	})

	for _, tt := range []struct {
		name    string
		parents string
		files   string
		remote  string
	}{
		{name: "merge commit", parents: commit + " " + parent + " " + strings.Repeat("c", 40), files: "pages/post.md", remote: parent},
		{name: "multiple files", parents: commit + " " + parent, files: "pages/a.md\npages/b.md", remote: parent},
		{name: "non markdown", parents: commit + " " + parent, files: "README.md\nsource.go", remote: parent},
		{name: "remote advanced", parents: commit + " " + parent, files: "pages/post.md", remote: strings.Repeat("c", 40)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pushed := false
			service := publisherTestService(func(_ context.Context, args ...string) (string, error) {
				call := strings.Join(args, " ")
				switch call {
				case "rev-parse HEAD":
					return commit, nil
				case "symbolic-ref -q HEAD":
					return publisherTestBranch, nil
				case "config --local --get remote.origin.url":
					return publisherTestOrigin, nil
				case "rev-list --parents -n 1 " + commit:
					return tt.parents, nil
				case "diff-tree --no-commit-id --name-only -r " + commit:
					return tt.files, nil
				case "ls-tree " + commit + " -- pages/post.md":
					return "100644 blob " + strings.Repeat("c", 40) + "\tpages/post.md", nil
				case "ls-remote --heads " + publisherTestPush + " " + publisherTestBranch:
					return tt.remote + "\t" + publisherTestBranch, nil
				default:
					if strings.HasPrefix(call, "push ") {
						pushed = true
					}
					return "", errors.New("unexpected Git command: " + call)
				}
			})
			if err := service.push(context.Background(), publisherTestBranch, commit); err == nil {
				t.Fatal("push() unexpectedly succeeded")
			}
			if pushed {
				t.Fatal("invalid commit reached git push")
			}
		})
	}
}

func TestUnixGitPublisher_RemoteHeadUsesLocalSocket(t *testing.T) {
	commit := strings.Repeat("a", 40)
	remote := strings.Repeat("b", 40)
	service := publisherTestService(func(_ context.Context, args ...string) (string, error) {
		switch strings.Join(args, " ") {
		case "rev-parse HEAD":
			return commit, nil
		case "symbolic-ref -q HEAD":
			return publisherTestBranch, nil
		case "config --local --get remote.origin.url":
			return publisherTestOrigin, nil
		case "ls-remote --heads " + publisherTestPush + " " + publisherTestBranch:
			return remote + "\t" + publisherTestBranch, nil
		default:
			return "", errors.New("unexpected Git command")
		}
	})
	socketDir, err := os.MkdirTemp("", "mp-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(socketDir); err != nil {
			t.Errorf("remove short socket directory: %v", err)
		}
	})
	socket := filepath.Join(socketDir, "s")
	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(context.Background(), "unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: service.handler(), ReadHeaderTimeout: 5 * time.Second}
	serveErrors := make(chan error, 1)
	go func() { serveErrors <- server.Serve(listener) }()
	t.Cleanup(func() {
		if err := server.Close(); err != nil {
			t.Errorf("close test HTTP server: %v", err)
		}
		if err := <-serveErrors; !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("test HTTP server Serve() error = %v", err)
		}
		_ = os.Remove(socket)
	})

	client := NewUnixGitPublisher(socket)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, err := client.RemoteHead(ctx, publisherTestBranch, commit)
	if err != nil {
		t.Fatalf("RemoteHead() error = %v", err)
	}
	if got != remote {
		t.Fatalf("RemoteHead() = %q, want %q", got, remote)
	}
}

func publisherTestService(run func(context.Context, ...string) (string, error)) *gitPublisherService {
	return &gitPublisherService{
		root:       "/data/source",
		originURL:  publisherTestOrigin,
		pushURL:    publisherTestPush,
		branch:     publisherTestBranch,
		privateKey: "/run/key",
		knownHosts: "/etc/ssh/known_hosts",
		run:        run,
	}
}

func TestGitPublisherServerConfig_RejectsArbitraryBranch(t *testing.T) {
	key, knownHosts := publisherTestCredentialFiles(t)
	config := GitPublisherServerConfig{
		SourceDir:  t.TempDir(),
		SocketPath: filepath.Join(t.TempDir(), "publisher.sock"),
		OriginURL:  publisherTestOrigin,
		Branch:     "refs/heads/main;git-config-core.sshCommand=evil",
		PrivateKey: key,
		KnownHosts: knownHosts,
	}
	if _, err := newGitPublisherService(config); err == nil {
		t.Fatal("newGitPublisherService() accepted an invalid branch")
	}
}

func TestGitPublisherServerConfig_RequiresOpenSSHPrivateKeyPermissions(t *testing.T) {
	key, knownHosts := publisherTestCredentialFiles(t)
	if err := os.Chmod(key, 0o440); err != nil {
		t.Fatal(err)
	}
	config := GitPublisherServerConfig{
		SourceDir:  t.TempDir(),
		SocketPath: filepath.Join(t.TempDir(), "publisher.sock"),
		OriginURL:  publisherTestOrigin,
		Branch:     publisherTestBranch,
		PrivateKey: key,
		KnownHosts: knownHosts,
	}
	if _, err := newGitPublisherService(config); err == nil {
		t.Fatal("newGitPublisherService() accepted a group-readable SSH private key")
	}
}

func TestGitPublisherValidatesRawFetchOrigin(t *testing.T) {
	root := t.TempDir()
	gitForReviewTest(t, root, "init", "-q", "-b", "main")
	gitForReviewTest(t, root, "remote", "add", "origin", publisherTestOrigin)
	service := publisherTestService(nil)
	service.root = root
	if err := service.validateOrigin(context.Background()); err != nil {
		t.Fatalf("validateOrigin() rejected configured origin: %v", err)
	}
	gitForReviewTest(t, root, "remote", "set-url", "origin", "https://github.com/attacker/other.git")
	if err := service.validateOrigin(context.Background()); err == nil {
		t.Fatal("validateOrigin() accepted a different repository")
	}
}

func TestGitPushAcceptsReviewedCommitHashAsSource(t *testing.T) {
	root := t.TempDir()
	gitForReviewTest(t, root, "init", "-q", "-b", "main")
	gitForReviewTest(t, root, "config", "user.name", "Publisher Test")
	gitForReviewTest(t, root, "config", "user.email", "publisher-test@example.com")
	if err := os.WriteFile(filepath.Join(root, "post.md"), []byte("# Post\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gitForReviewTest(t, root, "add", "post.md")
	gitForReviewTest(t, root, "commit", "-m", "publish one post")
	commit, err := gitOutputForTest(root, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	commit = strings.TrimSpace(commit)
	remote := t.TempDir()
	gitForReviewTest(t, remote, "init", "--bare")
	gitForReviewTest(t, root, "push", "--porcelain", remote, commit+":refs/heads/main")
	remoteHead, err := gitOutputForTest(remote, "rev-parse", "refs/heads/main")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(remoteHead) != commit {
		t.Fatalf("remote head %q does not match pushed commit %q", remoteHead, commit)
	}
}

func publisherTestCredentialFiles(t *testing.T) (key, knownHosts string) {
	t.Helper()
	dir := t.TempDir()
	key = filepath.Join(dir, "id_ed25519")
	knownHosts = filepath.Join(dir, "known_hosts")
	if err := os.WriteFile(key, []byte("test key"), 0o400); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(knownHosts, []byte("github.com ssh-ed25519 test-key\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return key, knownHosts
}
