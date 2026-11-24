package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

// Client is a Steam Web API client.
type Client struct {
	apiKey      string
	steamID     string
	httpClient  *http.Client
	cacheDir    string
	cacheExpiry time.Duration
}

// NewClient creates a new Steam API client.
func NewClient(apiKey, steamID, cacheDir string, cacheExpiry time.Duration) *Client {
	if cacheExpiry == 0 {
		cacheExpiry = time.Hour // Default 1 hour
	}

	return &Client{
		apiKey:  apiKey,
		steamID: steamID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheDir:    cacheDir,
		cacheExpiry: cacheExpiry,
	}
}

// cachePath returns the cache file path for a given key.
func (c *Client) cachePath(key string) string {
	// Sanitize the key to make it filesystem-safe
	safeKey := strings.NewReplacer(
		"/", "_",
		"?", "_",
		"&", "_",
		"=", "_",
		" ", "_",
	).Replace(key)

	// Remove leading/trailing underscores and collapse multiple underscores
	safeKey = strings.Trim(safeKey, "_")
	for strings.Contains(safeKey, "__") {
		safeKey = strings.ReplaceAll(safeKey, "__", "_")
	}

	return filepath.Join(c.cacheDir, fmt.Sprintf("steam_%s.json", safeKey))
}

// isCacheValid checks if a cache file is still valid.
func (c *Client) isCacheValid(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return time.Since(info.ModTime()) < c.cacheExpiry
}

// loadFromCache loads data from cache file.
func (c *Client) loadFromCache(key string, v interface{}) error {
	path := c.cachePath(key)
	if !c.isCacheValid(path) {
		return fmt.Errorf("cache not valid")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

// saveToCache saves data to cache file.
func (c *Client) saveToCache(key string, v interface{}) error {
	if c.cacheDir == "" {
		return nil // No caching configured
	}

	if err := os.MkdirAll(c.cacheDir, 0o755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	path := c.cachePath(key)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// buildURL builds a Steam API URL with the given endpoint and parameters.
func (c *Client) buildURL(endpoint string, params map[string]string) string {
	baseURL := "https://api.steampowered.com" + endpoint

	// Add required parameters
	if params == nil {
		params = make(map[string]string)
	}
	params["key"] = c.apiKey
	params["steamid"] = c.steamID

	// Build query string
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return baseURL + "?" + values.Encode()
}

// makeRequest makes an HTTP request to the Steam API.
func (c *Client) makeRequest(endpoint string, params map[string]string, v interface{}) error {
	cacheKey := endpoint
	for k, val := range params {
		cacheKey += "_" + k + "_" + val
	}

	// Try to load from cache first
	if err := c.loadFromCache(cacheKey, v); err == nil {
		return nil
	}

	requestURL := c.buildURL(endpoint, params)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("API request failed with status %d (unable to read body: %w)", resp.StatusCode, readErr)
		}
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Save to cache
	if err := c.saveToCache(cacheKey, v); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to save to cache: %v\n", err)
	}

	return nil
}

// GetOwnedGames retrieves the user's owned games.
func (c *Client) GetOwnedGames(includeAppInfo, includePlayedFreeGames bool) ([]GameResponse, error) {
	params := map[string]string{
		"format":                    "json",
		"include_appinfo":           strconv.FormatBool(includeAppInfo),
		"include_played_free_games": strconv.FormatBool(includePlayedFreeGames),
		"steamid":                   c.steamID,
	}

	var response OwnedGamesResponse
	if err := c.makeRequest("/IPlayerService/GetOwnedGames/v0001/", params, &response); err != nil {
		return nil, fmt.Errorf("failed to get owned games: %w", err)
	}

	return response.Response.Games, nil
}

// GetGameSchema retrieves the achievement schema for a game.
func (c *Client) GetGameSchema(appID int) (*GameSchemaResponse, error) {
	params := map[string]string{
		"appid": strconv.Itoa(appID),
		"l":     "english",
	}

	var response GameSchemaResponse
	if err := c.makeRequest("/ISteamUserStats/GetSchemaForGame/v2/", params, &response); err != nil {
		return nil, fmt.Errorf("failed to get game schema for app %d: %w", appID, err)
	}

	return &response, nil
}

// GetPlayerAchievements retrieves the player's achievements for a game.
func (c *Client) GetPlayerAchievements(appID int) (*PlayerAchievementsResponse, error) {
	params := map[string]string{
		"appid": strconv.Itoa(appID),
		"l":     "english",
	}

	var response PlayerAchievementsResponse
	if err := c.makeRequest("/ISteamUserStats/GetPlayerAchievements/v0001/", params, &response); err != nil {
		return nil, fmt.Errorf("failed to get player achievements for app %d: %w", appID, err)
	}

	return &response, nil
}

// GetGameWithAchievements retrieves a game with all achievement data.
func (c *Client) GetGameWithAchievements(game GameResponse) (*models.SteamGame, error) {
	// Get game schema
	schema, err := c.GetGameSchema(game.AppID)
	if err != nil {
		return nil, err
	}

	if schema.Game == nil || schema.Game.AvailableGameStats == nil {
		return nil, fmt.Errorf("no achievement data available for app %d", game.AppID)
	}

	// Get player achievements
	playerAchievements, err := c.GetPlayerAchievements(game.AppID)
	if err != nil {
		return nil, err
	}

	// Create achievement map from schema
	achievementSchema := make(map[string]AchievementSchema)
	for _, achievement := range schema.Game.AvailableGameStats.Achievements {
		achievementSchema[achievement.Name] = achievement
	}

	// Build achievements list
	achievements := make([]models.SteamAchievement, 0, len(playerAchievements.PlayerStats.Achievements))
	var unlockedCount int

	for _, playerAchievement := range playerAchievements.PlayerStats.Achievements {
		schema, exists := achievementSchema[playerAchievement.APIName]
		if !exists {
			continue
		}

		achievement := models.SteamAchievement{
			APIName:     playerAchievement.APIName,
			Name:        &schema.DisplayName,
			Description: &schema.Description,
			Icon:        &schema.Icon,
			IconGray:    &schema.IconGray,
			Achieved:    playerAchievement.Achieved,
		}

		if playerAchievement.UnlockTime != nil {
			achievement.UnlockTime = playerAchievement.UnlockTime
		}

		if achievement.IsUnlocked() {
			unlockedCount++
		}

		achievements = append(achievements, achievement)
	}

	// Calculate completion percentage
	completionPercentage := 0.0
	if len(achievements) > 0 {
		completionPercentage = float64(unlockedCount) / float64(len(achievements)) * 100
	}

	// Extract game info
	gameInfo := schema.Game
	var description *string
	if gameInfo.About != nil {
		description = &gameInfo.About.ShortDescription
	}

	developers := make([]string, 0, len(gameInfo.Developers))
	for _, dev := range gameInfo.Developers {
		developers = append(developers, dev.Name)
	}

	publishers := make([]string, 0, len(gameInfo.Publishers))
	for _, pub := range gameInfo.Publishers {
		publishers = append(publishers, pub.Name)
	}

	// Use the game name from the API if available
	gameName := game.Name
	if gameInfo.GameName != "" {
		gameName = gameInfo.GameName
	}

	steamGame := &models.SteamGame{
		AppID:                game.AppID,
		Name:                 gameName,
		Description:          description,
		Developers:           developers,
		Publishers:           publishers,
		Achievements:         achievements,
		TotalAchievements:    len(achievements),
		UnlockedAchievements: unlockedCount,
		CompletionPercentage: completionPercentage,
		PlaytimeForever:      &game.PlaytimeForever,
		Playtime2Weeks:       &game.Playtime2Weeks,
	}

	if game.RtimeLastPlayed != nil {
		steamGame.LastPlayed = game.RtimeLastPlayed
	}

	return steamGame, nil
}

// Response types for Steam API

type OwnedGamesResponse struct {
	Response struct {
		GameCount int            `json:"game_count"`
		Games     []GameResponse `json:"games"`
	} `json:"response"`
}

type GameResponse struct {
	AppID                    int    `json:"appid"`
	Name                     string `json:"name"`
	PlaytimeForever          int    `json:"playtime_forever"`
	Playtime2Weeks           int    `json:"playtime_2weeks"`
	RtimeLastPlayed          *int64 `json:"rtime_last_played"`
	ImgIconURL               string `json:"img_icon_url"`
	ImgLogoURL               string `json:"img_logo_url"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
}

type GameSchemaResponse struct {
	Game *GameInfo `json:"game"`
}

type GameInfo struct {
	GameName           string              `json:"gameName"`
	GameVersion        string              `json:"gameVersion"`
	AvailableGameStats *AvailableGameStats `json:"availableGameStats"`
	About              *GameAbout          `json:"about"`
	Developers         []GameDeveloper     `json:"developers"`
	Publishers         []GamePublisher     `json:"publishers"`
}

type AvailableGameStats struct {
	Achievements []AchievementSchema `json:"achievements"`
	Stats        []StatSchema        `json:"stats"`
}

type AchievementSchema struct {
	Name         string `json:"name"`
	DefaultValue int    `json:"defaultvalue"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	IconGray     string `json:"icongray"`
	Hidden       int    `json:"hidden"`
}

type StatSchema struct {
	Name         string `json:"name"`
	DefaultValue int    `json:"defaultvalue"`
	DisplayName  string `json:"displayName"`
}

type GameAbout struct {
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
}

type GameDeveloper struct {
	Name string `json:"name"`
}

type GamePublisher struct {
	Name string `json:"name"`
}

type PlayerAchievementsResponse struct {
	PlayerStats struct {
		SteamID      string              `json:"steamID"`
		GameName     string              `json:"gameName"`
		Achievements []PlayerAchievement `json:"achievements"`
		Success      bool                `json:"success"`
	} `json:"playerstats"`
}

type PlayerAchievement struct {
	APIName    string `json:"apiname"`
	Achieved   int    `json:"achieved"`
	UnlockTime *int64 `json:"unlocktime"`
}
