package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/blogroll"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// Constants for blogroll add command.
const (
	defaultCategory     = "Uncategorized"
	promptDefaultYesStr = "Y/n"
	promptDefaultNoStr  = "y/N"
	youtubeFeedBaseURL  = "https://www.youtube.com/feeds/videos.xml?channel_id="
	categoryCustomValue = "__custom__"
)

var (
	youtubeChannelIDPattern      = regexp.MustCompile(`UC[a-zA-Z0-9_-]{22}`)
	youtubeVideoIDPattern        = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)
	youtubeOEmbedEndpoint        = "https://www.youtube.com/oembed"
	blogrollAddHTTPClientFactory = func(timeout time.Duration) *http.Client { return &http.Client{Timeout: timeout} }
)

// Flags for blogroll add command.
var (
	blogrollAddTitle       string
	blogrollAddDescription string
	blogrollAddCategory    string
	blogrollAddType        string
	blogrollAddTags        []string
	blogrollAddSiteURL     string
	blogrollAddHandle      string
	blogrollAddActive      bool
	blogrollAddDryRun      bool
	blogrollAddNoPrompt    bool
)

// blogrollAddCmd adds a new feed to the blogroll.
var blogrollAddCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add a new external feed to blogroll",
	Long: `Add a new RSS/Atom feed to your blogroll by providing the feed URL.

The command will:
1. Fetch and validate the feed
2. Extract metadata (title, description, site URL)
3. Prompt for optional fields (category, tags)
4. Add the feed to your config

Example:
  markata-go blogroll add https://example.com/rss.xml
  markata-go blogroll add https://example.com/atom.xml --category tech
  markata-go blogroll add https://example.com/feed --no-prompt
  markata-go blogroll add https://example.com/feed --dry-run
  markata-go blogroll add https://example.com/feed --title "My Blog" --tags "go,programming"`,
	Args: cobra.ExactArgs(1),
	RunE: runBlogrollAdd,
}

func init() {
	blogrollCmd.AddCommand(blogrollAddCmd)

	blogrollAddCmd.Flags().StringVar(&blogrollAddTitle, "title", "", "feed title (auto-fetched if not set)")
	blogrollAddCmd.Flags().StringVar(&blogrollAddDescription, "description", "", "feed description (auto-fetched if not set)")
	blogrollAddCmd.Flags().StringVar(&blogrollAddCategory, "category", "", "category for grouping (default: "+defaultCategory+")")
	blogrollAddCmd.Flags().StringVar(&blogrollAddType, "type", "", "reader type: written, video, or podcast")
	blogrollAddCmd.Flags().StringSliceVar(&blogrollAddTags, "tags", nil, "comma-separated tags")
	blogrollAddCmd.Flags().StringVar(&blogrollAddSiteURL, "site-url", "", "main website URL (auto-fetched if not set)")
	blogrollAddCmd.Flags().StringVar(&blogrollAddHandle, "handle", "", "handle for @mentions (auto-generated if not set)")
	blogrollAddCmd.Flags().BoolVar(&blogrollAddActive, "active", true, "include in reader page")
	blogrollAddCmd.Flags().BoolVar(&blogrollAddDryRun, "dry-run", false, "preview changes without modifying config")
	blogrollAddCmd.Flags().BoolVar(&blogrollAddNoPrompt, "no-prompt", false, "skip interactive prompts, use defaults")
}

func runBlogrollAdd(_ *cobra.Command, args []string) error {
	rawInput := args[0]

	// Load config and check for duplicates
	configPath, cfg, err := loadConfigForBlogrollAdd()
	if err != nil {
		return err
	}

	feedURL, err := normalizeBlogrollAddInput(cfg, rawInput)
	if err != nil {
		return err
	}
	if blogrollAddType != "" && parseReaderFeedType(blogrollAddType) == "" {
		return fmt.Errorf("invalid --type %q: expected written, video, or podcast", blogrollAddType)
	}

	// Validate URL format
	if err := validateFeedURL(feedURL); err != nil {
		return err
	}

	if err := checkDuplicateFeedURL(cfg, feedURL); err != nil {
		return err
	}

	// Fetch feed metadata
	metadata, err := fetchFeedMetadataForAdd(cfg, feedURL)
	if err != nil {
		return err
	}

	// Build feed values from flags and metadata
	feedValues := buildFeedValues(metadata, feedURL)

	// Check for duplicate handle
	if err := checkDuplicateHandle(cfg, feedValues.handle, feedValues.title); err != nil {
		return err
	}

	// Display fetched info
	displayFetchedInfo(feedValues.title, feedValues.description)

	// Interactive prompts if not disabled
	if !blogrollAddNoPrompt {
		feedValues, err = promptForFeedValues(cfg, metadata, feedValues)
		if err != nil {
			return err
		}
	}

	// Build the feed config
	feedConfig := buildFeedConfig(feedURL, feedValues, metadata)

	// Display what will be added
	fmt.Println()
	printFeedConfig(feedConfig)

	if blogrollAddDryRun {
		fmt.Println("\n(dry-run mode - no changes written)")
		return nil
	}

	// Add feed to config and write
	return saveFeedToConfig(configPath, cfg, feedConfig)
}

type youtubeOEmbedResponse struct {
	AuthorURL string `json:"author_url"`
}

// feedValues holds the values for building a feed config.
type feedValues struct {
	title       string
	description string
	siteURL     string
	category    string
	feedType    models.ReaderFeedType
	handle      string
	tags        []string
	active      bool
}

// validateFeedURL validates that the URL is a valid HTTP/HTTPS URL.
func validateFeedURL(feedURL string) error {
	parsedURL, err := url.Parse(feedURL)
	if err != nil || (parsedURL.Scheme != subscriptionHTTP && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return fmt.Errorf("invalid feed URL: must be a valid HTTP/HTTPS URL")
	}
	return nil
}

func normalizeBlogrollAddInput(cfg *models.Config, input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", fmt.Errorf("invalid feed URL: must be a valid HTTP/HTTPS URL")
	}

	if strings.HasPrefix(strings.ToLower(trimmed), "yt:") {
		handle := strings.TrimSpace(trimmed[3:])
		if handle == "" {
			return "", fmt.Errorf("invalid youtube shortcut: expected yt:<handle>")
		}
		trimmed = "https://www.youtube.com/@" + strings.TrimPrefix(handle, "@")
	}

	parsedURL, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid feed URL: must be a valid HTTP/HTTPS URL")
	}

	if !isYouTubeURL(parsedURL) {
		return trimmed, nil
	}

	if strings.Contains(parsedURL.Path, "/feeds/videos.xml") {
		return trimmed, nil
	}

	if channelID := extractYouTubeChannelIDFromURL(trimmed); channelID != "" {
		return youtubeFeedBaseURL + channelID, nil
	}

	if isYouTubeVideoURL(parsedURL) {
		channelURL, err := fetchYouTubeAuthorURL(cfg, trimmed)
		if err != nil {
			return "", err
		}
		return resolveYouTubeChannelFeedURL(cfg, channelURL)
	}

	if isYouTubeChannelLikeURL(parsedURL) {
		return resolveYouTubeChannelFeedURL(cfg, trimmed)
	}

	return trimmed, nil
}

func resolveYouTubeChannelFeedURL(cfg *models.Config, rawURL string) (string, error) {
	if channelID := extractYouTubeChannelIDFromURL(rawURL); channelID != "" {
		return youtubeFeedBaseURL + channelID, nil
	}

	channelID, err := fetchYouTubeChannelID(cfg, rawURL)
	if err != nil {
		return "", err
	}

	return youtubeFeedBaseURL + channelID, nil
}

func fetchYouTubeAuthorURL(cfg *models.Config, videoURL string) (string, error) {
	timeout := blogrollAddTimeout(cfg)
	client := blogrollAddHTTPClientFactory(timeout)

	oembedURL := youtubeOEmbedEndpoint + "?url=" + url.QueryEscape(videoURL) + "&format=json"
	body, err := fetchBlogrollAddURL(client, oembedURL, "application/json")
	if err != nil {
		return "", fmt.Errorf("resolve youtube video channel: %w", err)
	}

	var response youtubeOEmbedResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("parse youtube oembed response: %w", err)
	}
	if response.AuthorURL == "" {
		return "", fmt.Errorf("resolve youtube video channel: author URL not found")
	}

	return response.AuthorURL, nil
}

func fetchYouTubeChannelID(cfg *models.Config, rawURL string) (string, error) {
	timeout := blogrollAddTimeout(cfg)
	client := blogrollAddHTTPClientFactory(timeout)
	body, err := fetchBlogrollAddURL(client, rawURL, "text/html,application/xhtml+xml")
	if err != nil {
		return "", fmt.Errorf("resolve youtube channel feed: %w", err)
	}

	channelID := extractYouTubeChannelIDFromHTML(string(body))
	if channelID == "" {
		return "", fmt.Errorf("resolve youtube channel feed: channel ID not found")
	}

	return channelID, nil
}

func fetchBlogrollAddURL(client *http.Client, targetURL, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, targetURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	req.Header.Set("User-Agent", "markata-go/1.0 (Blogroll Add)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

func blogrollAddTimeout(cfg *models.Config) time.Duration {
	timeout := 30 * time.Second
	if cfg != nil && cfg.Blogroll.Timeout > 0 {
		timeout = time.Duration(cfg.Blogroll.Timeout) * time.Second
	}
	return timeout
}

func isYouTubeURL(parsedURL *url.URL) bool {
	host := strings.ToLower(parsedURL.Host)
	host = strings.TrimPrefix(host, "www.")
	host = strings.TrimPrefix(host, "m.")
	return host == "youtube.com" || host == "youtu.be"
}

func isYouTubeVideoURL(parsedURL *url.URL) bool {
	host := strings.ToLower(parsedURL.Host)
	host = strings.TrimPrefix(host, "www.")
	host = strings.TrimPrefix(host, "m.")

	if host == "youtu.be" {
		videoID := strings.Trim(strings.TrimPrefix(parsedURL.Path, "/"), "/")
		return youtubeVideoIDPattern.MatchString(videoID)
	}

	if strings.HasPrefix(parsedURL.Path, "/watch") {
		return youtubeVideoIDPattern.MatchString(parsedURL.Query().Get("v"))
	}

	if strings.HasPrefix(parsedURL.Path, "/shorts/") {
		videoID := strings.TrimPrefix(parsedURL.Path, "/shorts/")
		videoID = strings.Trim(videoID, "/")
		return youtubeVideoIDPattern.MatchString(videoID)
	}

	return false
}

func isYouTubeChannelLikeURL(parsedURL *url.URL) bool {
	return strings.HasPrefix(parsedURL.Path, "/@") || strings.HasPrefix(parsedURL.Path, "/channel/") || strings.HasPrefix(parsedURL.Path, "/user/") || strings.HasPrefix(parsedURL.Path, "/c/")
}

func extractYouTubeChannelIDFromURL(rawURL string) string {
	if matches := youtubeChannelIDPattern.FindAllString(rawURL, -1); len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func extractYouTubeChannelIDFromHTML(html string) string {
	if idx := strings.Index(html, youtubeFeedBaseURL); idx >= 0 {
		candidate := html[idx+len(youtubeFeedBaseURL):]
		if match := youtubeChannelIDPattern.FindString(candidate); match != "" {
			return match
		}
	}

	return extractYouTubeChannelIDFromURL(html)
}

// loadConfigForBlogrollAdd loads the config file for adding a feed.
func loadConfigForBlogrollAdd() (string, *models.Config, error) {
	rootConfigPath := cfgFile
	if rootConfigPath == "" {
		var err error
		rootConfigPath, err = config.Discover()
		if err != nil {
			return "", nil, fmt.Errorf("no config file found: run 'markata-go init' first")
		}
	}

	cfg, err := config.Load(rootConfigPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load config: %w", err)
	}

	targetConfigPath, err := resolveBlogrollTargetConfigPath(rootConfigPath)
	if err != nil {
		return "", nil, err
	}

	return targetConfigPath, cfg, nil
}

func resolveBlogrollTargetConfigPath(rootConfigPath string) (string, error) {
	paths, err := config.DiscoverIncludedConfigPaths(rootConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve config includes: %w", err)
	}

	candidates := make([]string, 0)
	for _, path := range paths {
		ok, err := configFileHasBlogroll(path)
		if err != nil {
			return "", err
		}
		if ok {
			candidates = append(candidates, path)
		}
	}

	switch len(candidates) {
	case 0:
		return rootConfigPath, nil
	case 1:
		return candidates[0], nil
	default:
		if inputIsTerminal() && outputIsTerminal() {
			return promptForBlogrollConfigPath(candidates)
		}

		return "", fmt.Errorf("multiple config files contain blogroll config: %s (pass --config to choose the root config file)", strings.Join(candidates, ", "))
	}
}

func configFileHasBlogroll(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("failed to inspect config file %s: %w", path, err)
	}
	content := string(data)
	return strings.Contains(content, "[markata-go.blogroll]") ||
		strings.Contains(content, "[[markata-go.blogroll.feeds]]") ||
		strings.Contains(content, "[blogroll]") ||
		strings.Contains(content, "[[blogroll.feeds]]"), nil
}

func promptForBlogrollConfigPath(paths []string) (string, error) {
	options := make([]huh.Option[string], 0, len(paths))
	for _, path := range paths {
		options = append(options, huh.NewOption(path, path))
	}

	selected := paths[0]
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Blogroll Config File").
				Description("Multiple config files contain blogroll settings. Choose where to append the new feed.").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(createHuhThemeForPalette(currentCLIPalette))

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}

// checkDuplicateFeedURL checks if a feed with the given URL already exists.
func checkDuplicateFeedURL(cfg *models.Config, feedURL string) error {
	for i := range cfg.Blogroll.Feeds {
		if cfg.Blogroll.Feeds[i].URL == feedURL {
			return fmt.Errorf("feed with URL %q already exists (use 'blogroll update' to refresh metadata)", feedURL)
		}
	}
	return nil
}

// fetchFeedMetadataForAdd fetches metadata from the feed URL.
func fetchFeedMetadataForAdd(cfg *models.Config, feedURL string) (*blogroll.Metadata, error) {
	fmt.Printf("Fetching feed from %s...\n", feedURL)

	timeout := time.Duration(cfg.Blogroll.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	updater := blogroll.NewUpdater(timeout)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	metadata, err := updater.FetchMetadata(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed metadata: %w", err)
	}

	return metadata, nil
}

// buildFeedValues builds the feed values from flags and metadata.
func buildFeedValues(metadata *blogroll.Metadata, feedURL string) feedValues {
	title := blogrollAddTitle
	if title == "" {
		title = firstNonEmpty(metadata.Title, metadata.FeedTitle)
	}

	description := blogrollAddDescription
	if description == "" {
		description = firstNonEmpty(metadata.Description, metadata.FeedDescription)
	}

	siteURL := blogrollAddSiteURL
	if siteURL == "" {
		siteURL = metadata.SiteURL
	}

	category := blogrollAddCategory
	if category == "" {
		category = defaultCategory
	}

	feedType := parseReaderFeedType(blogrollAddType)
	if feedType == "" {
		feedType = inferReaderFeedType(metadata, feedURL)
	}

	handle := blogrollAddHandle
	if handle == "" {
		handle = generateHandle(title, feedURL)
	}

	tags := blogrollAddTags
	if len(tags) == 0 {
		tags = copyTags(firstNonEmptyTags(metadata.Tags, metadata.FeedTags))
	}

	return feedValues{
		title:       title,
		description: description,
		siteURL:     siteURL,
		category:    category,
		feedType:    feedType,
		handle:      handle,
		tags:        tags,
		active:      blogrollAddActive,
	}
}

// checkDuplicateHandle checks if a handle is already in use.
func checkDuplicateHandle(cfg *models.Config, handle, title string) error {
	for i := range cfg.Blogroll.Feeds {
		feed := &cfg.Blogroll.Feeds[i]
		if feed.Handle != "" && feed.Handle == handle {
			return fmt.Errorf("handle %q already in use by %q. Use --handle to specify a different handle", handle, feed.Title)
		}
	}
	_ = title // used in error message but handle check is separate
	return nil
}

// displayFetchedInfo displays the fetched feed info.
func displayFetchedInfo(title, description string) {
	if title != "" || description != "" {
		fmt.Printf("\nFetched feed: %q", title)
		if description != "" {
			truncDesc := description
			if len(truncDesc) > 60 {
				truncDesc = truncDesc[:57] + "..."
			}
			fmt.Printf(" - %q", truncDesc)
		}
		fmt.Println()
	}
}

// promptForFeedValues prompts the user for feed values interactively.
func promptForFeedValues(cfg *models.Config, _ *blogroll.Metadata, fv feedValues) (feedValues, error) {
	if inputIsTerminal() && outputIsTerminal() {
		return promptForFeedValuesHuh(cfg, fv)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println()

	fv.title = promptBlogroll(reader, "Title", fv.title)
	fv.description = promptBlogroll(reader, "Description", fv.description)
	fv.category = promptBlogroll(reader, "Category", fv.category)
	fv.feedType = parseReaderFeedType(promptBlogroll(reader, "Type (written/video/podcast)", string(fv.feedType)))
	if fv.feedType == "" {
		fv.feedType = models.ReaderFeedTypeWritten
	}

	if len(blogrollAddTags) == 0 {
		tagsInput := promptBlogroll(reader, "Tags (comma-separated)", strings.Join(fv.tags, ", "))
		fv.tags = parseTagsBlogroll(tagsInput)
	}

	fv.siteURL = promptBlogroll(reader, "Site URL", fv.siteURL)
	fv.handle = promptBlogroll(reader, "Handle (for @mentions)", fv.handle)
	fv.active = promptYesNoBlogroll(reader, "Include in reader page?", fv.active)

	return fv, nil
}

func promptForFeedValuesHuh(cfg *models.Config, fv feedValues) (feedValues, error) {
	theme := createHuhThemeForPalette(currentCLIPalette)

	tagsInput := strings.Join(fv.tags, ", ")
	typeChoice := string(fv.feedType)
	if typeChoice == "" {
		typeChoice = string(models.ReaderFeedTypeWritten)
	}
	categoryOptions := discoverExistingCategories(cfg)
	categoryChoice := fv.category
	customCategory := ""
	categoryGroup := huh.NewGroup(
		huh.NewInput().
			Title("Category").
			Description("Grouping label for this source").
			Value(&fv.category).
			Placeholder(defaultCategory),
	)
	groups := []*huh.Group{categoryGroup}
	if len(categoryOptions) > 0 {
		categoryChoice, customCategory = initialCategorySelection(fv.category, categoryOptions)
		categoryGroup = huh.NewGroup(
			huh.NewSelect[string]().
				Title("Category").
				Description("Type to filter existing categories or choose a custom one").
				Options(buildCategoryOptions(categoryOptions)...).
				Filtering(true).
				Value(&categoryChoice),
		)
		customCategoryGroup := huh.NewGroup(
			huh.NewInput().
				Title("Custom Category").
				Description("Enter a new category").
				Value(&customCategory).
				Placeholder(defaultCategory).
				Validate(func(s string) error {
					if categoryChoice == categoryCustomValue && strings.TrimSpace(s) == "" {
						return fmt.Errorf("category is required")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return categoryChoice != categoryCustomValue
		})
		groups = []*huh.Group{categoryGroup, customCategoryGroup}
	}

	groups = append([]*huh.Group{huh.NewGroup(
		huh.NewNote().
			Title("Add Blogroll Feed").
			Description("Review the fetched metadata and edit any fields before saving."),
		huh.NewInput().
			Title("Title").
			Description("Display name for this feed").
			Value(&fv.title).
			Placeholder("My Favorite Feed").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("title is required")
				}
				return nil
			}),
		huh.NewInput().
			Title("Description").
			Description("Short summary for your reader and feeds").
			Value(&fv.description),
		huh.NewSelect[string]().
			Title("Type").
			Description("Choose which reader stream this feed belongs to").
			Options(
				huh.NewOption("Written", string(models.ReaderFeedTypeWritten)),
				huh.NewOption("Video", string(models.ReaderFeedTypeVideo)),
				huh.NewOption("Podcast", string(models.ReaderFeedTypePodcast)),
			).
			Value(&typeChoice),
		huh.NewInput().
			Title("Tags").
			Description("Comma-separated tags").
			Value(&tagsInput),
		huh.NewInput().
			Title("Site URL").
			Description("Main website or channel URL").
			Value(&fv.siteURL),
		huh.NewInput().
			Title("Handle").
			Description("Used for @mentions").
			Value(&fv.handle).
			Placeholder("myfeed").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("handle is required")
				}
				return nil
			}),
		huh.NewConfirm().
			Title("Include in reader page?").
			Value(&fv.active),
	)}, groups...)

	form := huh.NewForm(groups...).WithTheme(theme)

	if err := form.Run(); err != nil {
		return fv, err
	}

	if len(categoryOptions) > 0 {
		if categoryChoice == categoryCustomValue {
			fv.category = strings.TrimSpace(customCategory)
		} else {
			fv.category = categoryChoice
		}
	}
	fv.feedType = parseReaderFeedType(typeChoice)
	if fv.feedType == "" {
		fv.feedType = models.ReaderFeedTypeWritten
	}

	if len(blogrollAddTags) == 0 {
		fv.tags = parseTagsBlogroll(tagsInput)
	}

	return fv, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func parseReaderFeedType(value string) models.ReaderFeedType {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(models.ReaderFeedTypeWritten):
		return models.ReaderFeedTypeWritten
	case string(models.ReaderFeedTypeVideo):
		return models.ReaderFeedTypeVideo
	case string(models.ReaderFeedTypePodcast):
		return models.ReaderFeedTypePodcast
	default:
		return ""
	}
}

func inferReaderFeedType(metadata *blogroll.Metadata, feedURL string) models.ReaderFeedType {
	combined := strings.ToLower(strings.Join([]string{
		feedURL,
		metadata.Title,
		metadata.Description,
		metadata.FeedTitle,
		metadata.FeedDescription,
		metadata.SiteURL,
		strings.Join(metadata.Tags, " "),
		strings.Join(metadata.FeedTags, " "),
	}, " "))

	switch {
	case strings.Contains(combined, "youtube.com") || strings.Contains(combined, "youtu.be"):
		return models.ReaderFeedTypeVideo
	case strings.Contains(combined, "podcast") || strings.Contains(combined, "hx-pod"):
		return models.ReaderFeedTypePodcast
	default:
		return models.ReaderFeedTypeWritten
	}
}

func discoverExistingCategories(cfg *models.Config) []string {
	categoryCounts := map[string]int{
		defaultCategory: 1,
	}

	if cfg != nil {
		for i := range cfg.Blogroll.Feeds {
			feed := &cfg.Blogroll.Feeds[i]
			if category := strings.TrimSpace(feed.Category); category != "" {
				categoryCounts[category]++
			}
		}

		patterns := cfg.GlobConfig.Patterns
		if len(patterns) == 0 {
			patterns = []string{"**/*.md"}
		}
		for _, pattern := range patterns {
			matches, err := doublestar.FilepathGlob(pattern)
			if err != nil {
				continue
			}
			for _, path := range matches {
				content, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				metadata, _, err := plugins.ParseFrontmatter(string(content))
				if err != nil {
					continue
				}
				if category, ok := metadata["category"].(string); ok {
					if trimmed := strings.TrimSpace(category); trimmed != "" {
						categoryCounts[trimmed]++
					}
				}
			}
		}
	}

	categories := make([]string, 0, len(categoryCounts))
	for category := range categoryCounts {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		left := categories[i]
		right := categories[j]
		if categoryCounts[left] != categoryCounts[right] {
			return categoryCounts[left] > categoryCounts[right]
		}
		if left == defaultCategory {
			return false
		}
		if right == defaultCategory {
			return true
		}
		return strings.ToLower(left) < strings.ToLower(right)
	})
	return categories
}

func buildCategoryOptions(categories []string) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(categories)+1)
	for _, category := range categories {
		options = append(options, huh.NewOption(category, category))
	}
	options = append(options, huh.NewOption("Custom...", categoryCustomValue))
	return options
}

func initialCategorySelection(current string, categories []string) (choice, custom string) {
	trimmedCurrent := strings.TrimSpace(current)
	if trimmedCurrent == "" || trimmedCurrent == defaultCategory {
		if len(categories) > 0 {
			return categories[0], ""
		}
		return defaultCategory, ""
	}

	for _, category := range categories {
		if category == trimmedCurrent {
			return category, ""
		}
	}
	return categoryCustomValue, trimmedCurrent
}

func firstNonEmptyTags(options ...[]string) []string {
	for _, option := range options {
		if len(option) > 0 {
			return option
		}
	}
	return nil
}

func copyTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}
	return append([]string{}, tags...)
}

// buildFeedConfig builds the ExternalFeedConfig from feed values.
func buildFeedConfig(feedURL string, fv feedValues, metadata *blogroll.Metadata) models.ExternalFeedConfig {
	feedConfig := models.ExternalFeedConfig{
		URL:         feedURL,
		Title:       fv.title,
		Description: fv.description,
		Category:    fv.category,
		Type:        fv.feedType,
		Tags:        fv.tags,
		SiteURL:     fv.siteURL,
		Handle:      fv.handle,
		Active:      &fv.active,
	}

	// Add image URL if fetched
	if metadata.ImageURL != "" {
		feedConfig.ImageURL = metadata.ImageURL
	}

	return feedConfig
}

// saveFeedToConfig adds the feed to the config and writes it.
func saveFeedToConfig(configPath string, cfg *models.Config, feedConfig models.ExternalFeedConfig) error {
	if strings.EqualFold(filepath.Ext(configPath), ".toml") {
		if err := appendFeedToTOMLConfig(configPath, feedConfig, !cfg.Blogroll.Enabled); err != nil {
			return fmt.Errorf("failed to append feed to config: %w", err)
		}
		if !cfg.Blogroll.Enabled {
			fmt.Println("\nEnabled blogroll in config (was disabled)")
		}
		fmt.Printf("\nFeed added to %s\n", configPath)
		return nil
	}

	cfg.Blogroll.Feeds = append(cfg.Blogroll.Feeds, feedConfig)

	// Enable blogroll if not already enabled
	if !cfg.Blogroll.Enabled {
		cfg.Blogroll.Enabled = true
		fmt.Println("\nEnabled blogroll in config (was disabled)")
	}

	// Write updated config
	if err := writeConfigUpdate(configPath, cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("\nFeed added to %s\n", configPath)
	return nil
}

func appendFeedToTOMLConfig(configPath string, feedConfig models.ExternalFeedConfig, ensureEnabled bool) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	content := string(data)
	prefix := detectBlogrollTablePrefix(content)
	var snippet strings.Builder

	if trimmed := strings.TrimRight(content, "\n"); trimmed != content {
		content = trimmed
	}

	if content != "" {
		snippet.WriteString("\n\n")
	}

	if ensureEnabled && !strings.Contains(content, "["+prefix+"]") {
		snippet.WriteString("[" + prefix + "]\n")
		snippet.WriteString("enabled = true\n\n")
	}

	snippet.WriteString("[[" + prefix + ".feeds]]\n")
	snippet.WriteString("url = " + formatTOMLString(feedConfig.URL) + "\n")
	if feedConfig.Title != "" {
		snippet.WriteString("title = " + formatTOMLString(feedConfig.Title) + "\n")
	}
	if feedConfig.Description != "" {
		snippet.WriteString("description = " + formatTOMLString(feedConfig.Description) + "\n")
	}
	if feedConfig.Category != "" {
		snippet.WriteString("category = " + formatTOMLString(feedConfig.Category) + "\n")
	}
	if feedConfig.Type != "" {
		snippet.WriteString("type = " + formatTOMLString(string(feedConfig.Type)) + "\n")
	}
	if len(feedConfig.Tags) > 0 {
		snippet.WriteString("tags = [" + formatTOMLStringArray(feedConfig.Tags) + "]\n")
	}
	if feedConfig.SiteURL != "" {
		snippet.WriteString("site_url = " + formatTOMLString(feedConfig.SiteURL) + "\n")
	}
	if feedConfig.Handle != "" {
		snippet.WriteString("handle = " + formatTOMLString(feedConfig.Handle) + "\n")
	}
	if feedConfig.ImageURL != "" {
		snippet.WriteString("image_url = " + formatTOMLString(feedConfig.ImageURL) + "\n")
	}
	if feedConfig.Active != nil {
		snippet.WriteString("active = " + strconv.FormatBool(*feedConfig.Active) + "\n")
	}

	updated := content + snippet.String() + "\n"
	info, err := os.Stat(configPath)
	if err != nil {
		return fmt.Errorf("stat config: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(updated), info.Mode().Perm()); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func detectBlogrollTablePrefix(content string) string {
	if strings.Contains(content, "[blogroll]") || strings.Contains(content, "[[blogroll.feeds]]") {
		return "blogroll"
	}
	return "markata-go.blogroll"
}

func formatTOMLString(value string) string {
	return strconv.Quote(value)
}

func formatTOMLStringArray(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, formatTOMLString(value))
	}
	return strings.Join(quoted, ", ")
}

// promptBlogroll displays a prompt and returns user input or default.
func promptBlogroll(reader *bufio.Reader, question, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", question, defaultVal)
	} else {
		fmt.Printf("%s: ", question)
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultVal
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

// promptYesNoBlogroll displays a yes/no prompt.
func promptYesNoBlogroll(reader *bufio.Reader, question string, defaultYes bool) bool {
	defaultStr := promptDefaultNoStr
	if defaultYes {
		defaultStr = promptDefaultYesStr
	}
	fmt.Printf("%s (%s): ", question, defaultStr)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultYes
	}
	return input == "y" || input == statusYes
}

// generateHandle creates a handle from title or URL.
func generateHandle(title, feedURL string) string {
	// Try to generate from title first
	if title != "" {
		handle := strings.ToLower(title)
		// Remove common suffixes
		handle = strings.TrimSuffix(handle, "'s blog")
		handle = strings.TrimSuffix(handle, " blog")
		// Replace spaces with nothing and remove non-alphanumeric
		var result strings.Builder
		for _, r := range handle {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				result.WriteRune(r)
			}
		}
		if result.Len() > 0 {
			return result.String()
		}
	}

	// Fall back to domain name
	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		return ""
	}

	host := parsedURL.Host
	host = strings.TrimPrefix(host, "www.")

	// Extract first part of domain
	parts := strings.Split(host, ".")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

// parseTagsBlogroll parses a comma-separated string into a slice of tags.
func parseTagsBlogroll(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	tags := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}

// printFeedConfig displays the feed configuration that will be added.
func printFeedConfig(feed models.ExternalFeedConfig) {
	ext := strings.ToLower(filepath.Ext(cfgFile))
	if ext == "" {
		// Auto-detect from discovered config
		configPath, err := config.Discover()
		if err == nil {
			ext = strings.ToLower(filepath.Ext(configPath))
		}
	}

	switch ext {
	case extYAML, extYML:
		printFeedConfigYAML(feed)
	default:
		printFeedConfigTOML(feed)
	}
}

// printFeedConfigTOML prints the feed config in TOML format.
func printFeedConfigTOML(feed models.ExternalFeedConfig) {
	fmt.Println("[[blogroll.feeds]]")
	fmt.Printf("url = %q\n", feed.URL)
	if feed.Title != "" {
		fmt.Printf("title = %q\n", feed.Title)
	}
	if feed.Description != "" {
		fmt.Printf("description = %q\n", feed.Description)
	}
	if feed.Category != "" && feed.Category != defaultCategory {
		fmt.Printf("category = %q\n", feed.Category)
	}
	if feed.Type != "" {
		fmt.Printf("type = %q\n", feed.Type)
	}
	if len(feed.Tags) > 0 {
		fmt.Printf("tags = %s\n", formatTOMLArray(feed.Tags))
	}
	if feed.SiteURL != "" {
		fmt.Printf("site_url = %q\n", feed.SiteURL)
	}
	if feed.Handle != "" {
		fmt.Printf("handle = %q\n", feed.Handle)
	}
	if feed.ImageURL != "" {
		fmt.Printf("image_url = %q\n", feed.ImageURL)
	}
	if feed.Active != nil && !*feed.Active {
		fmt.Printf("active = %v\n", *feed.Active)
	}
}

// printFeedConfigYAML prints the feed config in YAML format.
func printFeedConfigYAML(feed models.ExternalFeedConfig) {
	fmt.Println("- url:", feed.URL)
	if feed.Title != "" {
		fmt.Println("  title:", feed.Title)
	}
	if feed.Description != "" {
		fmt.Println("  description:", feed.Description)
	}
	if feed.Category != "" && feed.Category != defaultCategory {
		fmt.Println("  category:", feed.Category)
	}
	if feed.Type != "" {
		fmt.Println("  type:", feed.Type)
	}
	if len(feed.Tags) > 0 {
		fmt.Print("  tags: [")
		for i, t := range feed.Tags {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(t)
		}
		fmt.Println("]")
	}
	if feed.SiteURL != "" {
		fmt.Println("  site_url:", feed.SiteURL)
	}
	if feed.Handle != "" {
		fmt.Println("  handle:", feed.Handle)
	}
	if feed.ImageURL != "" {
		fmt.Println("  image_url:", feed.ImageURL)
	}
	if feed.Active != nil && !*feed.Active {
		fmt.Printf("  active: %v\n", *feed.Active)
	}
}

// formatTOMLArray formats a string slice as a TOML array.
func formatTOMLArray(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("%q", item)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}
