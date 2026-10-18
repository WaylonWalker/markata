package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/spf13/cobra"
)

func resetSiteDirState(t *testing.T) {
	t.Helper()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	oldSiteDir := siteDir
	oldActiveSiteDir := activeSiteDir
	oldActiveContentDir := activeContentDir
	oldSiteDirSelected := siteDirSelected
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
		siteDir = oldSiteDir
		activeSiteDir = oldActiveSiteDir
		activeContentDir = oldActiveContentDir
		siteDirSelected = oldSiteDirSelected
	})
}

func assertSameDirectory(t *testing.T, got, want string) {
	t.Helper()
	gotInfo, err := os.Stat(got)
	if err != nil {
		t.Fatalf("stat path %q: %v", got, err)
	}
	wantInfo, err := os.Stat(want)
	if err != nil {
		t.Fatalf("stat path %q: %v", want, err)
	}
	if !os.SameFile(gotInfo, wantInfo) {
		t.Fatalf("directory = %q, want %q", got, want)
	}
}

func TestActivateSiteDir_FlagTakesPrecedenceAndChangesDirectory(t *testing.T) {
	flagSite := t.TempDir()
	envSite := t.TempDir()
	resetSiteDirState(t)
	t.Setenv(siteDirEnv, envSite)
	siteDir = flagSite

	if err := activateSiteDir(); err != nil {
		t.Fatalf("activateSiteDir() error = %v", err)
	}
	gotWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	assertSameDirectory(t, gotWD, flagSite)
	if activeSiteDir != flagSite || !siteDirSelected {
		t.Fatalf("active site = %q, selected = %t", activeSiteDir, siteDirSelected)
	}
}

func TestActivateSiteDir_UsesEnvironment(t *testing.T) {
	site := t.TempDir()
	resetSiteDirState(t)
	t.Setenv(siteDirEnv, site)

	if err := activateSiteDir(); err != nil {
		t.Fatalf("activateSiteDir() error = %v", err)
	}
	gotWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	assertSameDirectory(t, gotWD, site)
}

func TestActivateSiteDir_RejectsNonDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-a-directory")
	resetSiteDirState(t)
	if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	siteDir = path

	err := activateSiteDir()
	if err == nil || !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("activateSiteDir() error = %v, want non-directory error", err)
	}
}

func TestSourcePathForOutput_ExplicitSiteIsAbsolute(t *testing.T) {
	site := t.TempDir()
	resetSiteDirState(t)
	siteDirSelected = true
	activeSiteDir = site

	got := sourcePathForOutput(filepath.Join("posts", "hello.md"))
	want := filepath.Join(activeSiteDir, "posts", "hello.md")
	if got != want {
		t.Fatalf("sourcePathForOutput() = %q, want %q", got, want)
	}
}

func TestCreateManager_SourcePathsUseConfigDirectory(t *testing.T) {
	site := t.TempDir()
	resetSiteDirState(t)
	configDir := filepath.Join(site, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("create config directory: %v", err)
	}
	configPath := filepath.Join(configDir, "site.toml")
	if err := os.WriteFile(configPath, []byte("title = \"Site\"\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	oldCfgFile := cfgFile
	cfgFile = filepath.Join("config", "site.toml")
	t.Cleanup(func() { cfgFile = oldCfgFile })
	siteDir = site
	if err := activateSiteDir(); err != nil {
		t.Fatalf("activateSiteDir() error = %v", err)
	}
	if _, err := createManager(cfgFile); err != nil {
		t.Fatalf("createManager() error = %v", err)
	}

	got := sourcePathForOutput(filepath.Join("posts", "hello.md"))
	if filepath.Base(got) != "hello.md" || filepath.Base(filepath.Dir(got)) != "posts" {
		t.Fatalf("sourcePathForOutput() = %q, want config/posts/hello.md", got)
	}
	assertSameDirectory(t, filepath.Dir(filepath.Dir(got)), configDir)
}

func TestSteamPersistentPreRun_ActivatesSelectedSite(t *testing.T) {
	site := t.TempDir()
	resetSiteDirState(t)
	siteDir = site

	if err := steamCmd.PersistentPreRunE(steamCmd, nil); err != nil {
		t.Fatalf("steam persistent pre-run error = %v", err)
	}
	gotWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	assertSameDirectory(t, gotWD, site)
}

func TestBuilderAdmin_UsesReleaseDirFlag(t *testing.T) {
	if builderAdminCmd.Flags().Lookup("release-dir") == nil {
		t.Fatal("builder-admin must expose --release-dir for its release root")
	}
	if builderAdminCmd.Flags().Lookup("site-dir") != nil {
		t.Fatal("builder-admin must not shadow the global --site-dir flag")
	}
}

func TestRunNewCommand_UsesSelectedSiteAndRelativeConfig(t *testing.T) {
	site := t.TempDir()
	resetSiteDirState(t)
	resetNewCommandFlags()
	configDir := filepath.Join(site, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("create config directory: %v", err)
	}
	configPath := filepath.Join(configDir, "site.toml")
	if err := os.WriteFile(configPath, []byte("title = \"Site\"\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	mergePath := filepath.Join(configDir, "content.toml")
	if err := os.WriteFile(mergePath, []byte("[content_templates.placement]\npost = \"notes\"\n"), 0o600); err != nil {
		t.Fatalf("write merged config: %v", err)
	}

	oldCfgFile := cfgFile
	oldMergeConfigFiles := mergeConfigFiles
	oldOut := newCmd.OutOrStdout()
	cfgFile = filepath.Join("config", "site.toml")
	mergeConfigFiles = []string{filepath.Join("config", "content.toml")}
	stdout := bytes.NewBuffer(nil)
	newCmd.SetOut(stdout)
	t.Cleanup(func() {
		cfgFile = oldCfgFile
		mergeConfigFiles = oldMergeConfigFiles
		newCmd.SetOut(oldOut)
		resetNewCommandFlags()
	})

	siteDir = site
	if err := activateSiteDir(); err != nil {
		t.Fatalf("activateSiteDir() error = %v", err)
	}
	if err := runNewCommand(newCmd, []string{"Site Directory Post"}); err != nil {
		t.Fatalf("runNewCommand() error = %v", err)
	}

	want := filepath.Join(site, "notes", "site-directory-post.md")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("created file %q: %v", want, err)
	}
	if !strings.Contains(stdout.String(), "Created: "+want) {
		t.Fatalf("creation output = %q, want absolute path %q", stdout.String(), want)
	}
}

func TestRenderPosts_UsesAbsolutePathsForSelectedSite(t *testing.T) {
	site := t.TempDir()
	resetSiteDirState(t)
	siteDirSelected = true
	activeSiteDir = site
	want := filepath.Join(activeSiteDir, "posts", "hello.md")
	post := &models.Post{Path: filepath.Join("posts", "hello.md")}

	for _, format := range []string{listFormatTable, listFormatJSON, listFormatCSV, listFormatPath} {
		t.Run(format, func(t *testing.T) {
			stdout := bytes.NewBuffer(nil)
			command := &cobra.Command{Use: "posts"}
			command.SetOut(stdout)
			currentCmd = command
			t.Cleanup(func() { currentCmd = nil })

			if err := renderPosts(format, []*models.Post{post}); err != nil {
				t.Fatalf("renderPosts() error = %v", err)
			}
			output := stdout.String()
			expectedPath := want
			if format == listFormatJSON {
				expectedPath = strings.ReplaceAll(expectedPath, `\`, `\\`)
			}
			if !strings.Contains(output, expectedPath) {
				t.Fatalf("rendered output missing %q:\n%s", expectedPath, output)
			}
		})
	}

	stdout := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "search"}
	command.SetOut(stdout)
	currentCmd = command
	t.Cleanup(func() { currentCmd = nil })
	if err := renderSearchTable([]cliSearchResult{{post: post, score: 2}}, "hello"); err != nil {
		t.Fatalf("renderSearchTable() error = %v", err)
	}
	if !strings.Contains(stdout.String(), want) {
		t.Fatalf("search output missing %q:\n%s", want, stdout.String())
	}
}
