package builderadmin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCommandArgsForceCleansRenderCache(t *testing.T) {
	for _, test := range []struct {
		name         string
		force        bool
		cacheMount   string
		sourceDir    string
		wantCacheDir string
		wantClean    bool
	}{
		{name: "normal build retains cache", cacheMount: "/cache", sourceDir: "/source"},
		{name: "forced build uses cache mount", force: true, cacheMount: "/cache", sourceDir: "/source", wantCacheDir: "/cache/build", wantClean: true},
		{name: "forced build uses source cache without mount", force: true, sourceDir: "/source", wantCacheDir: "/source/.markata", wantClean: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			overrideDir := t.TempDir()
			service := &Service{
				cfg: Config{
					ConfigPath: "markata-go.toml",
					CacheMount: test.cacheMount,
					SourceDir:  test.sourceDir,
				},
				overrideDir: overrideDir,
			}
			args, cleanup, err := service.buildCommandArgs("build-test", t.TempDir(), test.force)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			gotClean := false
			for _, arg := range args {
				if arg == "--clean" {
					gotClean = true
				}
			}
			if gotClean != test.wantClean {
				t.Fatalf("build args %q include --clean=%t, want %t", strings.Join(args, " "), gotClean, test.wantClean)
			}
			override, err := os.ReadFile(filepath.Join(overrideDir, "builder-admin.toml"))
			if err != nil {
				t.Fatal(err)
			}
			gotCacheSetting := ""
			if test.wantCacheDir != "" {
				gotCacheSetting = "cache_dir = \"" + test.wantCacheDir + "\""
			}
			if strings.Contains(string(override), "cache_dir = ") != (test.wantCacheDir != "") {
				t.Fatalf("override %q has cache_dir setting, wantCacheDir=%q", override, test.wantCacheDir)
			}
			if gotCacheSetting != "" && !strings.Contains(string(override), gotCacheSetting) {
				t.Fatalf("override %q missing expected cache setting %q", override, gotCacheSetting)
			}
		})
	}
}
