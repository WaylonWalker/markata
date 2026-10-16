package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/builderadmin"
	"github.com/WaylonWalker/markata-go/pkg/config"
	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/spf13/cobra"
)

var (
	builderAdminHost                  string
	builderAdminPort                  int
	builderAdminSourceDir             string
	builderAdminSiteDir               string
	builderAdminWorkDir               string
	builderAdminCacheMount            string
	builderAdminHistoryDir            string
	builderAdminWatch                 bool
	builderAdminWatchDebounce         time.Duration
	builderAdminFast                  bool
	builderAdminMermaidMode           string
	builderAdminReleasesKeep          int
	builderAdminSuccessfulBuildsKeep  int
	builderAdminFailedBuildsKeep      int
	builderAdminRefreshRunsKeep       int
	builderAdminBuildTimeout          time.Duration
	builderAdminRefreshTaskSpecs      []string
	builderAdminTrustedProxyCIDRs     []string
	builderAdminAuthUserIDHeader      string
	builderAdminAuthUsernameHeader    string
	builderAdminAuthDisplayNameHeader string
	builderAdminAuthEmailHeader       string
	builderAdminAuthGroupsHeader      string
	builderAdminAuthRolesHeader       string
	builderAdminAuthScopesHeader      string
	builderAdminAuthMethodHeader      string
	builderAdminPublicAuthOrigin      string
	builderAdminPublicOrigin          string
	builderAdminPreviewOrigin         string
	builderAdminWebhookEnabled        bool
	builderAdminWebhookBranch         string
	builderAdminWebhookSecret         string
)

var builderAdminAutoBuildOnUpgrade bool

var builderAdminCmd = &cobra.Command{
	Use:   "builder-admin",
	Short: "Run the long-lived builder admin HTTP service",
	Long: `Run the long-lived builder admin HTTP service.

The service keeps a warm build worker running for hostPath and Kubernetes authoring loops.
It exposes a queue-driven UI/API for builds, logs, releases, rollback, and scheduled refresh tasks.`,
	RunE: runBuilderAdmin,
}

func init() {
	rootCmd.AddCommand(builderAdminCmd)
	builderAdminCmd.Flags().StringVar(&builderAdminHost, "host", "127.0.0.1", "host to bind to")
	builderAdminCmd.Flags().IntVar(&builderAdminPort, "port", 8080, "port to listen on")
	builderAdminCmd.Flags().StringVar(&builderAdminSourceDir, "source-dir", ".", "source directory to watch and build from")
	builderAdminCmd.Flags().StringVar(&builderAdminSiteDir, "release-dir", "public", "release root that contains releases/ and current")
	builderAdminCmd.Flags().StringVar(&builderAdminWorkDir, "work-dir", "", "optional build workspace directory; defaults to <release-dir>/.build-work")
	builderAdminCmd.Flags().StringVar(&builderAdminCacheMount, "cache-mount", "", "optional dedicated cache mount for .markata symlinks")
	builderAdminCmd.Flags().StringVar(&builderAdminHistoryDir, "history-dir", "", "directory for persisted builder-admin state and logs")
	builderAdminCmd.Flags().BoolVar(&builderAdminWatch, "watch", true, "enable recursive file watching")
	builderAdminCmd.Flags().BoolVar(&builderAdminAutoBuildOnUpgrade, "build-on-startup-revision-change", os.Getenv("MARKATA_GO_BUILDER_ADMIN_STARTUP_REBUILD") == "true", "build on leader startup when the builder revision differs from the last successful build")
	builderAdminCmd.Flags().DurationVar(&builderAdminWatchDebounce, "watch-debounce", 2*time.Second, "debounce window for file-watch rebuilds")
	builderAdminCmd.Flags().BoolVar(&builderAdminFast, "fast", false, "run queued builds with --fast")
	builderAdminCmd.Flags().StringVar(&builderAdminMermaidMode, "mermaid-mode", "", "override [markata-go.mermaid].mode for queued builds")
	builderAdminCmd.Flags().IntVar(&builderAdminReleasesKeep, "releases-keep", 25, "number of rendered releases to keep")
	builderAdminCmd.Flags().IntVar(&builderAdminSuccessfulBuildsKeep, "successful-builds-keep", 60, "number of successful build records to keep")
	builderAdminCmd.Flags().IntVar(&builderAdminFailedBuildsKeep, "failed-builds-keep", 100, "number of failed build records to keep")
	builderAdminCmd.Flags().IntVar(&builderAdminRefreshRunsKeep, "refresh-runs-keep", 100, "number of refresh run records to keep")
	builderAdminCmd.Flags().DurationVar(&builderAdminBuildTimeout, "build-timeout", 2*time.Hour, "maximum runtime for a queued build or refresh task")
	builderAdminCmd.Flags().StringArrayVar(&builderAdminRefreshTaskSpecs, "refresh-task", nil, "repeatable task spec: name|every|enqueue|arg1|arg2...")
	builderAdminCmd.Flags().StringArrayVar(&builderAdminTrustedProxyCIDRs, "trusted-proxy-cidr", nil, "repeatable CIDR permitted to supply trusted proxy identity headers")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthUserIDHeader, "auth-user-id-header", "X-Hlab-User-Id", "trusted proxy header containing the durable operator ID")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthUsernameHeader, "auth-username-header", "X-Hlab-Username", "trusted proxy header containing the operator username")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthDisplayNameHeader, "auth-display-name-header", "X-Hlab-Display-Name", "trusted proxy header containing the operator display name")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthEmailHeader, "auth-email-header", "X-Hlab-Email", "trusted proxy header containing the operator email")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthGroupsHeader, "auth-groups-header", "X-Hlab-Groups", "trusted proxy header containing operator groups")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthRolesHeader, "auth-roles-header", "X-Hlab-Roles", "trusted proxy header containing operator roles")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthScopesHeader, "auth-scopes-header", "X-Hlab-Scopes", "trusted proxy header containing operator scopes")
	builderAdminCmd.Flags().StringVar(&builderAdminAuthMethodHeader, "auth-method-header", "X-Hlab-Auth-Method", "trusted proxy header identifying bearer-authenticated machine API requests")
	builderAdminCmd.Flags().StringVar(&builderAdminPublicAuthOrigin, "public-auth-origin", "", "optional HTTPS hlab-auth origin used for the signed-in operator profile picture")
	builderAdminCmd.Flags().StringVar(&builderAdminPublicOrigin, "public-origin", "", "exact HTTPS public origin used to validate browser mutations")
	builderAdminCmd.Flags().StringVar(&builderAdminPreviewOrigin, "preview-origin", "", "HTTPS site origin used for retained release previews")
	builderAdminCmd.Flags().BoolVar(&builderAdminWebhookEnabled, "webhook-enabled", false, "enable signed GitHub and Forgejo push webhooks")
	builderAdminCmd.Flags().StringVar(&builderAdminWebhookBranch, "webhook-branch", "main", "Git branch accepted by the webhook")
	builderAdminCmd.Flags().StringVar(&builderAdminWebhookSecret, "webhook-secret", "", "HMAC-SHA256 secret for GitHub and Forgejo webhooks")
}

func runBuilderAdmin(cmd *cobra.Command, _ []string) error {
	refreshTasks, err := parseRefreshTasks(builderAdminRefreshTaskSpecs)
	if err != nil {
		return err
	}
	if err := validateBuilderAdminWorkDir(builderAdminWorkDir, builderAdminSourceDir, builderAdminSiteDir); err != nil {
		return err
	}
	configPath := resolveBuilderAdminConfigPath(cfgFile, builderAdminSourceDir)
	authHeaders, err := resolveBuilderAdminAuthHeaders(cmd, configPath)
	if err != nil {
		return err
	}
	webhook, err := resolveBuilderAdminWebhook(cmd, configPath)
	if err != nil {
		return err
	}
	svc, err := builderadmin.New(builderadmin.Config{
		Host:                 builderAdminHost,
		Port:                 builderAdminPort,
		SourceDir:            builderAdminSourceDir,
		BuilderVersion:       Version,
		BuilderCommit:        Commit,
		SiteDir:              builderAdminSiteDir,
		WorkDir:              builderAdminWorkDir,
		ConfigPath:           configPath,
		CacheMount:           builderAdminCacheMount,
		HistoryDir:           builderAdminHistoryDir,
		WatchEnabled:         builderAdminWatch,
		AutoBuildOnUpgrade:   builderAdminAutoBuildOnUpgrade,
		WatchDebounce:        builderAdminWatchDebounce,
		Fast:                 builderAdminFast,
		MermaidMode:          builderAdminMermaidMode,
		ReleasesKeep:         builderAdminReleasesKeep,
		SuccessfulBuildsKeep: builderAdminSuccessfulBuildsKeep,
		FailedBuildsKeep:     builderAdminFailedBuildsKeep,
		RefreshRunsKeep:      builderAdminRefreshRunsKeep,
		RefreshTasks:         refreshTasks,
		BuildTimeout:         builderAdminBuildTimeout,
		TrustedProxyCIDRs:    builderAdminTrustedProxyCIDRs,
		AuthHeaders:          authHeaders,
		PublicAuthOrigin:     builderAdminPublicAuthOrigin,
		PublicOrigin:         builderAdminPublicOrigin,
		PreviewOrigin:        builderAdminPreviewOrigin,
		Webhook:              webhook,
	})
	if err != nil {
		return err
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	return svc.Start(ctx)
}

func validateBuilderAdminWorkDir(workDir, sourceDir, siteDir string) error {
	if strings.TrimSpace(workDir) == "" {
		return nil
	}
	workAbs, err := filepath.Abs(workDir)
	if err != nil {
		return fmt.Errorf("resolve builder-admin work directory: %w", err)
	}
	if filepath.Dir(workAbs) == workAbs {
		return fmt.Errorf("builder-admin work directory must not be a filesystem root: %s", workDir)
	}
	sourceAbs, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("resolve builder-admin source directory: %w", err)
	}
	siteAbs, err := filepath.Abs(siteDir)
	if err != nil {
		return fmt.Errorf("resolve builder-admin release directory: %w", err)
	}
	if workAbs == sourceAbs {
		return fmt.Errorf("builder-admin work directory must not equal the source directory: %s", sourceDir)
	}
	if workAbs == siteAbs {
		return fmt.Errorf("builder-admin work directory must not equal the release directory: %s", siteDir)
	}
	if pathWithin(workAbs, siteAbs) {
		rel, err := filepath.Rel(siteAbs, workAbs)
		if err != nil {
			return fmt.Errorf("resolve builder-admin work directory relative to release directory: %w", err)
		}
		first := rel
		if separator := strings.IndexRune(rel, os.PathSeparator); separator >= 0 {
			first = rel[:separator]
		}
		switch first {
		case "releases", "current", ".builder-admin":
			return fmt.Errorf("builder-admin work directory must not use reserved release path %q: %s", first, workDir)
		}
		return nil
	}
	if pathWithin(workAbs, sourceAbs) {
		return fmt.Errorf("builder-admin work directory must not be inside the source directory: %s", workDir)
	}
	return nil
}

func pathWithin(path, parent string) bool {
	rel, err := filepath.Rel(parent, path)
	if err != nil || rel == "." || rel == ".." {
		return rel == "."
	}
	return !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// resolveBuilderAdminWebhook applies site configuration and MARKATA_GO_ overrides,
// then applies explicitly supplied CLI flags as the highest-precedence source.
func resolveBuilderAdminWebhook(cmd *cobra.Command, configPath string) (builderadmin.WebhookConfig, error) {
	siteConfig, err := config.Load(configPath)
	if err != nil {
		return builderadmin.WebhookConfig{}, fmt.Errorf("load builder-admin configuration: %w", err)
	}
	configured := siteConfig.BuilderAdmin.Webhook
	webhook := builderadmin.WebhookConfig{
		Branch: "main",
	}
	if configured.Enabled != nil {
		webhook.Enabled = *configured.Enabled
	}
	if configured.Branch != nil {
		webhook.Branch = *configured.Branch
	}
	if configured.Secret != nil {
		webhook.Secret = *configured.Secret
	}
	if webhook.Branch == "" {
		webhook.Branch = "main"
	}
	for _, flag := range []struct {
		name string
		set  func()
	}{
		{"webhook-enabled", func() { webhook.Enabled = builderAdminWebhookEnabled }},
		{"webhook-branch", func() { webhook.Branch = builderAdminWebhookBranch }},
		{"webhook-secret", func() { webhook.Secret = builderAdminWebhookSecret }},
	} {
		if cmd.Flags().Changed(flag.name) {
			flag.set()
		}
	}
	return webhook, nil
}

func resolveBuilderAdminConfigPath(configPath, sourceDir string) string {
	if configPath != "" {
		if filepath.IsAbs(configPath) {
			return configPath
		}
		return filepath.Join(sourceDir, configPath)
	}
	for _, name := range []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"} {
		path := filepath.Join(sourceDir, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}
	return ""
}

// resolveBuilderAdminAuthHeaders applies site configuration and MARKATA_GO_ overrides,
// then applies explicitly supplied CLI flags as the highest-precedence source.
func resolveBuilderAdminAuthHeaders(cmd *cobra.Command, configPath string) (builderadmin.AuthHeaders, error) {
	siteConfig, err := config.Load(configPath)
	if err != nil {
		return builderadmin.AuthHeaders{}, fmt.Errorf("load builder-admin configuration: %w", err)
	}
	headers := builderadmin.DefaultAuthHeaders()
	applyBuilderAdminConfigHeaders(&headers, siteConfig.BuilderAdmin.Auth.Headers)
	for _, flag := range []struct {
		name   string
		value  string
		target *string
	}{
		{"auth-user-id-header", builderAdminAuthUserIDHeader, &headers.UserID},
		{"auth-username-header", builderAdminAuthUsernameHeader, &headers.Username},
		{"auth-display-name-header", builderAdminAuthDisplayNameHeader, &headers.DisplayName},
		{"auth-email-header", builderAdminAuthEmailHeader, &headers.Email},
		{"auth-groups-header", builderAdminAuthGroupsHeader, &headers.Groups},
		{"auth-roles-header", builderAdminAuthRolesHeader, &headers.Roles},
		{"auth-scopes-header", builderAdminAuthScopesHeader, &headers.Scopes},
		{"auth-method-header", builderAdminAuthMethodHeader, &headers.AuthMethod},
	} {
		if cmd.Flags().Changed(flag.name) {
			*flag.target = flag.value
		}
	}
	return headers, nil
}

func applyBuilderAdminConfigHeaders(headers *builderadmin.AuthHeaders, configured models.BuilderAdminAuthHeadersConfig) {
	for _, field := range []struct {
		source *string
		target *string
	}{
		{configured.UserID, &headers.UserID},
		{configured.Username, &headers.Username},
		{configured.DisplayName, &headers.DisplayName},
		{configured.Email, &headers.Email},
		{configured.Groups, &headers.Groups},
		{configured.Roles, &headers.Roles},
		{configured.Scopes, &headers.Scopes},
		{configured.AuthMethod, &headers.AuthMethod},
	} {
		if field.source != nil {
			*field.target = *field.source
		}
	}
}

func parseRefreshTasks(specs []string) ([]builderadmin.RefreshTaskConfig, error) {
	tasks := make([]builderadmin.RefreshTaskConfig, 0, len(specs))
	for _, spec := range specs {
		parts := strings.Split(spec, "|")
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid --refresh-task %q: expected name|every|enqueue|arg1|arg2", spec)
		}
		enqueue := strings.EqualFold(parts[2], "true") || parts[2] == "1" || strings.EqualFold(parts[2], "yes")
		tasks = append(tasks, builderadmin.RefreshTaskConfig{
			Name:                  parts[0],
			Every:                 parts[1],
			EnqueueBuildOnSuccess: enqueue,
			Args:                  append([]string(nil), parts[3:]...),
		})
	}
	return tasks, nil
}
