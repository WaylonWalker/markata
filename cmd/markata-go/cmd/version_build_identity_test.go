package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestBuildEntrypointsInjectVersionIdentity(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("get test working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("could not find module go.mod from test working directory")
		}
		root = parent
	}

	entrypoints := []string{"Dockerfile", "justfile"}
	flagsByEntryPoint := make(map[string]map[string]string, len(entrypoints))
	flagPattern := regexp.MustCompile(`-X\s+([^\s=]+)`)
	fields := []string{"Version", "Commit", "Date"}
	for _, entrypoint := range entrypoints {
		source, err := os.ReadFile(filepath.Join(root, entrypoint))
		if err != nil {
			t.Fatalf("read %s: %v", entrypoint, err)
		}
		flags := make(map[string]string, len(fields))
		for _, match := range flagPattern.FindAllStringSubmatch(string(source), -1) {
			target := match[1]
			for _, field := range fields {
				if strings.HasSuffix(target, "."+field) {
					if _, exists := flags[field]; exists {
						t.Fatalf("%s contains multiple linker targets for %s", entrypoint, field)
					}
					flags[field] = target
				}
			}
		}
		for _, field := range fields {
			if flags[field] == "" {
				t.Fatalf("%s does not set the %s linker target", entrypoint, field)
			}
		}
		flagsByEntryPoint[entrypoint] = flags
	}

	// Docker and just are separate production and local build paths. They must
	// target the same variables so either build reports the supplied identity.
	for _, field := range fields {
		if flagsByEntryPoint["Dockerfile"][field] != flagsByEntryPoint["justfile"][field] {
			t.Fatalf("Dockerfile and justfile target different %s variables: %q != %q", field, flagsByEntryPoint["Dockerfile"][field], flagsByEntryPoint["justfile"][field])
		}
	}

	binaryName := "version-identity-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binary := filepath.Join(t.TempDir(), binaryName)
	const (
		version = "9.8.7-test"
		commit  = "deadbeef1234"
		date    = "2026-10-18T12:34:56Z"
	)
	ldflags := strings.Join([]string{
		"-X " + flagsByEntryPoint["Dockerfile"]["Version"] + "=" + version,
		"-X " + flagsByEntryPoint["Dockerfile"]["Commit"] + "=" + commit,
		"-X " + flagsByEntryPoint["Dockerfile"]["Date"] + "=" + date,
	}, " ")

	build := exec.Command("go", "build", "-ldflags", ldflags, "-o", binary, "./cmd/markata-go")
	build.Dir = root
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI with version metadata: %v\n%s", err, output)
	}

	output, err := exec.Command(binary, "version").CombinedOutput()
	if err != nil {
		t.Fatalf("run version command: %v\n%s", err, output)
	}
	for _, want := range []string{
		"markata-go " + version,
		"commit:    " + commit,
		"built:     " + date,
	} {
		if !strings.Contains(string(output), want) {
			t.Errorf("version output %q does not contain %q", output, want)
		}
	}

	short, err := exec.Command(binary, "version", "--short").CombinedOutput()
	if err != nil {
		t.Fatalf("run short version command: %v\n%s", err, short)
	}
	if strings.TrimSpace(string(short)) != version {
		t.Errorf("short version = %q, want %q", strings.TrimSpace(string(short)), version)
	}
}
