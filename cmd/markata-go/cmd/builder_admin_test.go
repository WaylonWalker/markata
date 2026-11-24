package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/builderadmin"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestApplyBuilderAdminConfigHeaders(t *testing.T) {
	userID := "X-Authentik-Uid"
	empty := ""
	headers := builderadmin.DefaultAuthHeaders()
	applyBuilderAdminConfigHeaders(&headers, models.BuilderAdminAuthHeadersConfig{
		UserID:      &userID,
		DisplayName: &empty,
	})
	if headers.UserID != userID {
		t.Fatalf("UserID = %q, want %q", headers.UserID, userID)
	}
	if headers.DisplayName != "" {
		t.Fatalf("DisplayName = %q, want empty", headers.DisplayName)
	}
	if headers.Username != "X-Hlab-Username" {
		t.Fatalf("Username = %q, want default", headers.Username)
	}
}

func TestResolveBuilderAdminAuthHeaders_ConfigFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	if err := os.WriteFile(configPath, []byte(`
[markata-go.builder_admin.auth.headers]
user_id = "X-Authentik-Uid"
display_name = ""
`), 0o600); err != nil {
		t.Fatal(err)
	}
	headers, err := resolveBuilderAdminAuthHeaders(builderAdminCmd, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if headers.UserID != "X-Authentik-Uid" || headers.DisplayName != "" {
		t.Fatalf("headers = %#v", headers)
	}
}

func TestResolveBuilderAdminConfigPath(t *testing.T) {
	sourceDir := t.TempDir()
	configPath := filepath.Join(sourceDir, "markata-go.toml")
	if err := os.WriteFile(configPath, []byte("[markata-go]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := resolveBuilderAdminConfigPath("", sourceDir); got != configPath {
		t.Fatalf("resolveBuilderAdminConfigPath() = %q, want %q", got, configPath)
	}
	if got := resolveBuilderAdminConfigPath("config/site.toml", sourceDir); got != filepath.Join(sourceDir, "config", "site.toml") {
		t.Fatalf("relative config path = %q", got)
	}
}

func TestResolveBuilderAdminWebhook_ConfigFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	if err := os.WriteFile(configPath, []byte(`
[markata-go.builder_admin.webhook]
enabled = true
branch = "dev"
secret = "test-secret"
`), 0o600); err != nil {
		t.Fatal(err)
	}
	webhook, err := resolveBuilderAdminWebhook(builderAdminCmd, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if !webhook.Enabled || webhook.Branch != "dev" || webhook.Secret != "test-secret" {
		t.Fatalf("webhook = %#v", webhook)
	}
}

func TestResolveBuilderAdminWebhook_EnvironmentOverridesConfigFile(t *testing.T) {
	t.Setenv("MARKATA_GO_BUILDER_ADMIN_WEBHOOK_BRANCH", "qa")
	t.Setenv("MARKATA_GO_BUILDER_ADMIN_WEBHOOK_SECRET", "environment-secret")
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	if err := os.WriteFile(configPath, []byte(`
[markata-go.builder_admin.webhook]
enabled = true
branch = "dev"
secret = "file-secret"
`), 0o600); err != nil {
		t.Fatal(err)
	}
	webhook, err := resolveBuilderAdminWebhook(builderAdminCmd, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if !webhook.Enabled || webhook.Branch != "qa" || webhook.Secret != "environment-secret" {
		t.Fatalf("webhook = %#v", webhook)
	}
}

func TestResolveBuilderAdminWebhook_ExplicitCLIOverridesEnvironment(t *testing.T) {
	t.Setenv("MARKATA_GO_BUILDER_ADMIN_WEBHOOK_BRANCH", "environment")
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	if err := os.WriteFile(configPath, []byte(`
[markata-go.builder_admin.webhook]
branch = "file"
`), 0o600); err != nil {
		t.Fatal(err)
	}

	flag := builderAdminCmd.Flags().Lookup("webhook-branch")
	previousValue, previousChanged := flag.Value.String(), flag.Changed
	if err := flag.Value.Set("cli"); err != nil {
		t.Fatal(err)
	}
	flag.Changed = true
	t.Cleanup(func() {
		if err := flag.Value.Set(previousValue); err != nil {
			t.Errorf("restore flag %s: %v", flag.Name, err)
		}
		flag.Changed = previousChanged
	})

	webhook, err := resolveBuilderAdminWebhook(builderAdminCmd, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if webhook.Branch != "cli" {
		t.Fatalf("Branch = %q, want explicit CLI value", webhook.Branch)
	}
}
