package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEditorialOptions_DiscoveredConfigErrors(t *testing.T) {
	directory := t.TempDir()
	t.Chdir(directory)
	previousConfig := cfgFile
	cfgFile = ""
	t.Cleanup(func() { cfgFile = previousConfig })
	if err := os.WriteFile(filepath.Join(directory, "markata-go.toml"), []byte("[markata-go\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := editorialOptions(); err == nil || !strings.Contains(err.Error(), "load config") {
		t.Fatalf("malformed discovered config ignored: %v", err)
	}
}

func TestEditorialRoot_UsesActivatedSiteDirectory(t *testing.T) {
	previousSelected, previousActive, previousRequested := siteDirSelected, activeSiteDir, siteDir
	t.Cleanup(func() { siteDirSelected, activeSiteDir, siteDir = previousSelected, previousActive, previousRequested })
	siteDirSelected, activeSiteDir, siteDir = true, t.TempDir(), "relative/site"
	if editorialRoot() != activeSiteDir {
		t.Fatalf("editorial root = %q, want %q", editorialRoot(), activeSiteDir)
	}
}
