package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/steam"
	"github.com/spf13/cobra"
)

var (
	steamAPIKey   string
	steamID       string
	steamPostsDir string
	steamTemplate string
	steamCacheDur int
	steamMinHours float64
	steamVerbose  bool
)

// steamCmd represents the steam command group
var steamCmd = &cobra.Command{
	Use:   "steam",
	Short: "Steam achievements commands",
	Long: `Commands for importing Steam games and achievements as markdown posts.

These commands fetch data from the Steam Web API and create markdown posts
for games and individual achievements. Requires Steam API key and Steam ID.

Get your Steam Web API key from: https://steamcommunity.com/dev/apikey`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		if err := activateSiteDir(); err != nil {
			return err
		}
		// Load environment variables from .env file if it exists
		if err := loadEnvFile(); err != nil && steamVerbose {
			fmt.Fprintf(os.Stderr, "Warning: failed to load .env file: %v\n", err)
		}
		return nil
	},
}

// steamAchievementsCmd represents the steam achievements command
var steamAchievementsCmd = &cobra.Command{
	Use:   "achievements",
	Short: "Create Steam achievement posts from unlocked achievements",
	Long: `Fetches your Steam games and creates individual markdown posts for each
unlocked achievement. Posts are organized by date and game name.`,
	RunE: runSteamAchievements,
}

// steamGamesCmd represents the steam games command
var steamGamesCmd = &cobra.Command{
	Use:   "games",
	Short: "Create Steam game posts with all achievements",
	Long: `Fetches your Steam games and creates comprehensive markdown posts for each
game, showing all achievements with completion status and cross-links.`,
	RunE: runSteamGames,
}

// init adds the steam commands and their flags
func init() {
	// Add steam command group to root
	rootCmd.AddCommand(steamCmd)

	// Add subcommands
	steamCmd.AddCommand(steamAchievementsCmd)
	steamCmd.AddCommand(steamGamesCmd)

	// Global steam flags
	steamCmd.PersistentFlags().StringVar(&steamAPIKey, "api-key", "", "Steam Web API key (env: STEAM_API_KEY)")
	steamCmd.PersistentFlags().StringVar(&steamID, "steam-id", "", "Steam ID (env: STEAM_ID)")
	steamCmd.PersistentFlags().StringVar(&steamPostsDir, "posts-dir", "pages/steam", "Directory to create posts")
	steamCmd.PersistentFlags().StringVar(&steamTemplate, "template", "steam_achievement", "Template name for posts")
	steamCmd.PersistentFlags().IntVar(&steamCacheDur, "cache-duration", 3600, "Cache duration in seconds")
	steamCmd.PersistentFlags().Float64Var(&steamMinHours, "min-hours", 3.0, "Minimum playtime in hours")
	steamCmd.PersistentFlags().BoolVar(&steamVerbose, "verbose", false, "Verbose output")
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	envPath := ".env"
	if _, err := os.Stat(envPath); err != nil {
		return nil // No .env file, that's OK
	}

	file, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Simple .env parser
	content, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	return nil
}

// getSteamConfig creates a SteamConfig from flags and environment
func getSteamConfig() (*models.SteamConfig, error) {
	config := models.NewSteamConfig()

	// Use flags first, then environment variables
	if steamAPIKey != "" {
		config.APIKey = steamAPIKey
	} else if env := os.Getenv("STEAM_API_KEY"); env != "" {
		config.APIKey = env
	}

	if steamID != "" {
		config.SteamID = steamID
	} else if env := os.Getenv("STEAM_ID"); env != "" {
		config.SteamID = env
	}

	// Override with other flags
	if steamPostsDir != "pages/steam" {
		config.PostsDir = steamPostsDir
	}
	if steamTemplate != "steam_achievement" {
		config.Template = steamTemplate
	}
	if steamCacheDur != 3600 {
		config.CacheDuration = steamCacheDur
	}
	if steamMinHours != 3.0 {
		config.MinPlaytimeHours = steamMinHours
	}

	// Validate required fields
	if config.APIKey == "" {
		return nil, fmt.Errorf("steam API key is required (use --api-key or STEAM_API_KEY)")
	}
	if config.SteamID == "" {
		return nil, fmt.Errorf("steam ID is required (use --steam-id or STEAM_ID)")
	}

	return &config, nil
}

// runSteamAchievements executes the steam achievements command
func runSteamAchievements(_ *cobra.Command, _ []string) error {
	config, err := getSteamConfig()
	if err != nil {
		return err
	}

	if steamVerbose {
		fmt.Fprintf(os.Stderr, "🎮 Fetching Steam achievements...\n")
		fmt.Fprintf(os.Stderr, "📁 Posts directory: %s\n", config.PostsDir)
		fmt.Fprintf(os.Stderr, "⏱️  Cache duration: %d seconds\n", config.CacheDuration)
		fmt.Fprintf(os.Stderr, "🕐 Minimum playtime: %.1f hours\n", config.MinPlaytimeHours)
	}

	// Create Steam client
	client := steam.NewClient(config.APIKey, config.SteamID,
		filepath.Join(".markata", "steam_cache"),
		0) // Use default cache duration

	// Get owned games
	games, err := client.GetOwnedGames(true, true)
	if err != nil {
		return fmt.Errorf("failed to get owned games: %w", err)
	}

	if steamVerbose {
		fmt.Fprintf(os.Stderr, "📚 Found %d games in Steam library\n", len(games))
	}

	// Create posts directory
	if err := os.MkdirAll(config.PostsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create posts directory: %w", err)
	}

	totalPosts := 0
	gamesProcessed := 0
	gamesFiltered := 0

	for _, game := range games {
		playtimeHours := float64(game.PlaytimeForever) / 60.0

		if playtimeHours < config.MinPlaytimeHours {
			gamesFiltered++
			if steamVerbose {
				fmt.Fprintf(os.Stderr, "⏭️  Skipping %s - %.1fh (minimum: %.1fh)\n",
					game.Name, playtimeHours, config.MinPlaytimeHours)
			}
			continue
		}

		if steamVerbose {
			fmt.Fprintf(os.Stderr, "🎯 Processing: %s (ID: %d) - %.1fh\n",
				game.Name, game.AppID, playtimeHours)
		}

		// Get game with achievements
		steamGame, err := client.GetGameWithAchievements(game)
		if err != nil {
			if steamVerbose {
				fmt.Fprintf(os.Stderr, "⚠️  No achievement data for %s: %v\n", game.Name, err)
			}
			continue
		}

		// Create posts for unlocked achievements
		postsCreated := createAchievementPosts(steamGame, config)
		totalPosts += postsCreated
		gamesProcessed++

		if steamVerbose && postsCreated > 0 {
			fmt.Fprintf(os.Stderr, "✅ Created %d achievement posts for %s\n",
				postsCreated, steamGame.Name)
		}
	}

	if steamVerbose {
		if gamesFiltered > 0 {
			fmt.Fprintf(os.Stderr, "📊 Filtered out %d games with less than %.1fh playtime\n",
				gamesFiltered, config.MinPlaytimeHours)
		}
		fmt.Fprintf(os.Stderr, "🎉 Created %d achievement posts from %d games\n",
			totalPosts, gamesProcessed)
	}

	return nil
}

// runSteamGames executes the steam games command
func runSteamGames(_ *cobra.Command, _ []string) error {
	config, err := getSteamConfig()
	if err != nil {
		return err
	}

	if steamVerbose {
		fmt.Fprintf(os.Stderr, "🎮 Creating Steam game posts...\n")
		fmt.Fprintf(os.Stderr, "📁 Posts directory: %s\n", config.PostsDir)
		fmt.Fprintf(os.Stderr, "⏱️  Cache duration: %d seconds\n", config.CacheDuration)
		fmt.Fprintf(os.Stderr, "🕐 Minimum playtime: %.1f hours\n", config.MinPlaytimeHours)
	}

	// Create Steam client
	client := steam.NewClient(config.APIKey, config.SteamID,
		filepath.Join(".markata", "steam_cache"),
		0) // Use default cache duration

	// Get owned games
	games, err := client.GetOwnedGames(true, true)
	if err != nil {
		return fmt.Errorf("failed to get owned games: %w", err)
	}

	if steamVerbose {
		fmt.Fprintf(os.Stderr, "📚 Found %d games in Steam library\n", len(games))
	}

	// Create posts directory
	if err := os.MkdirAll(config.PostsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create posts directory: %w", err)
	}

	gamePostsCreated := 0
	gamesFiltered := 0

	for _, game := range games {
		playtimeHours := float64(game.PlaytimeForever) / 60.0

		if playtimeHours < config.MinPlaytimeHours {
			gamesFiltered++
			if steamVerbose {
				fmt.Fprintf(os.Stderr, "⏭️  Skipping %s - %.1fh (minimum: %.1fh)\n",
					game.Name, playtimeHours, config.MinPlaytimeHours)
			}
			continue
		}

		if steamVerbose {
			fmt.Fprintf(os.Stderr, "🎯 Processing: %s (ID: %d) - %.1fh\n",
				game.Name, game.AppID, playtimeHours)
		}

		// Get game with achievements
		steamGame, err := client.GetGameWithAchievements(game)
		if err != nil {
			if steamVerbose {
				fmt.Fprintf(os.Stderr, "⚠️  No achievement data for %s: %v\n", game.Name, err)
			}
			continue
		}

		// Create game post
		if createGamePost(steamGame, config) {
			gamePostsCreated++
		}
	}

	if steamVerbose {
		if gamesFiltered > 0 {
			fmt.Fprintf(os.Stderr, "📊 Filtered out %d games with less than %.1fh playtime\n",
				gamesFiltered, config.MinPlaytimeHours)
		}
		fmt.Fprintf(os.Stderr, "🎉 Created %d game posts\n", gamePostsCreated)
	}

	return nil
}
