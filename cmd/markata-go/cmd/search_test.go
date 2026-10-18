package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestRunSearchServer_InvalidMode(t *testing.T) {
	originalMode := searchServerMode
	defer func() { searchServerMode = originalMode }()

	searchServerMode = "invalid"
	err := runSearchServer(searchServerCmd, nil)
	if err == nil {
		t.Fatal("expected invalid mode error")
	}
	if err.Error() != "invalid --mode \"invalid\" (expected runtime-index, watch-content, or read-only-index)" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunSearchServer_ReadOnlyRequiresIndexDir(t *testing.T) {
	originalMode := searchServerMode
	originalIndexDir := searchServerIndexDir
	defer func() {
		searchServerMode = originalMode
		searchServerIndexDir = originalIndexDir
	}()

	searchServerMode = "read-only-index"
	searchServerIndexDir = ""
	err := runSearchServer(searchServerCmd, nil)
	if err == nil {
		t.Fatal("expected missing index-dir error")
	}
	if err.Error() != "--index-dir is required in read-only-index mode" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchBuildIndexCommand_BuildsArtifacts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	indexDir := filepath.Join(tmpDir, "search.bleve")
	hashPath := filepath.Join(tmpDir, "search.hash")

	originalCfg := cfgFile
	originalMerge := mergeConfigFiles
	defer func() {
		cfgFile = originalCfg
		mergeConfigFiles = originalMerge
	}()

	cfgFile = filepath.Join("..", "..", "..", "benchmarks", "site", "markata-go.toml")
	mergeConfigFiles = nil

	cmd := searchBuildIndexCmd()
	cmd.SetArgs([]string{"--index-dir", indexDir, "--hash-path", hashPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute build-index: %v", err)
	}

	if info, err := os.Stat(indexDir); err != nil || !info.IsDir() {
		t.Fatalf("expected index dir to exist: %v", err)
	}
	if _, err := os.Stat(hashPath); err != nil {
		t.Fatalf("expected hash file to exist: %v", err)
	}
}

func TestConfiguredSearchEndpoints(t *testing.T) {
	tests := []struct {
		name            string
		cfg             *models.Config
		wantClient      string
		wantHandlerPath string
	}{
		{
			name:            "defaults without config",
			cfg:             nil,
			wantClient:      defaultSearchEndpoint,
			wantHandlerPath: defaultSearchEndpoint,
		},
		{
			name:            "relative bleve endpoint",
			cfg:             &models.Config{Search: models.SearchConfig{Bleve: models.BleveSearchConfig{Endpoint: "/custom/search"}}},
			wantClient:      "/custom/search",
			wantHandlerPath: "/custom/search",
		},
		{
			name:            "absolute bleve endpoint uses path for handler",
			cfg:             &models.Config{Search: models.SearchConfig{Bleve: models.BleveSearchConfig{Endpoint: "https://search.example.com/api/search"}}},
			wantClient:      "https://search.example.com/api/search",
			wantHandlerPath: "/api/search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, handlerPath := configuredSearchEndpoints(tt.cfg)
			if client != tt.wantClient {
				t.Fatalf("client endpoint = %q, want %q", client, tt.wantClient)
			}
			if handlerPath != tt.wantHandlerPath {
				t.Fatalf("handler path = %q, want %q", handlerPath, tt.wantHandlerPath)
			}
		})
	}
}
