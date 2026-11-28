package builderadmin

import (
	"strings"
	"testing"
)

func TestBuildCommandArgsForceCleansRenderCache(t *testing.T) {
	service := &Service{
		cfg:         Config{ConfigPath: "markata-go.toml"},
		overrideDir: t.TempDir(),
	}
	for _, test := range []struct {
		name      string
		force     bool
		wantClean bool
	}{
		{name: "normal build retains cache"},
		{name: "forced build clears cache", force: true, wantClean: true},
	} {
		t.Run(test.name, func(t *testing.T) {
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
		})
	}
}
