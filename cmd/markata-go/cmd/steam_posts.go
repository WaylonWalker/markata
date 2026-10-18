package cmd

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"gopkg.in/yaml.v3"
)

// createAchievementPosts creates individual posts for unlocked achievements.
func createAchievementPosts(game *models.SteamGame, config *models.SteamConfig) int {
	postsCreated := 0

	for _, achievement := range game.UnlockedAchievementsList() {
		if createAchievementPost(game, &achievement, config) {
			postsCreated++
		}
	}

	return postsCreated
}

// createAchievementPost creates a single achievement post.
func createAchievementPost(game *models.SteamGame, achievement *models.SteamAchievement, config *models.SteamConfig) bool {
	unlockDate := achievement.UnlockDate()
	if unlockDate == nil {
		return false // Skip achievements without unlock time
	}

	// Create safe filename components
	gameName := sanitizeForFilename(game.Name)
	achievementName := sanitizeForFilename(coalesceStr(achievement.Name, "achievement"))

	dateStr := unlockDate.Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s-%s.md", dateStr, gameName, achievementName)
	postPath := filepath.Join(config.PostsDir, filename)

	// Check if file already exists
	if _, err := os.Stat(postPath); err == nil {
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "⏭️  Achievement post already exists: %s\n", filename)
		}
		return false
	}

	// Create game page URL for cross-linking
	gameURL := fmt.Sprintf("/%s/", gameName)

	// Create content
	content := fmt.Sprintf(`---
title: %s
description: %s
date: "%s"
published: true
templateKey: %s
steam:
  game: %s
  app_id: %d
  achievement:
    name: %s
    description: %s
    api_name: %s
    unlock_time: %d
    unlock_date: "%s"
    icon: %s
    icongray: %s
tags: %s
slug: "steam/%s"
---

%s

%s

Unlocked in **[%s](%s)** on %s.

---

*Achievement data automatically imported from Steam.*`,
		escapeYAML(coalesceStr(achievement.Name, fmt.Sprintf("Achievement in %s", game.Name))),
		escapeYAML(fmt.Sprintf("%s: %s", game.Name, coalesceStr(achievement.Description, ""))),
		unlockDate.Format("2006-01-02"),
		config.Template,
		escapeYAML(game.Name),
		game.AppID,
		escapeYAML(coalesceStr(achievement.Name, "")),
		escapeYAML(coalesceStr(achievement.Description, "")),
		escapeYAML(achievement.APIName),
		coalesceInt64(achievement.UnlockTime, 0),
		unlockDate.Format(time.RFC3339),
		escapeYAML(coalesceStr(achievement.Icon, "")),
		escapeYAML(coalesceStr(achievement.IconGray, "")),
		escapeYAMLArray([]string{"steam-achievement", "steam", "achievement", gameName}),
		gameName,
		achievementIconHTML(achievement.Icon, coalesceStr(achievement.Name, "Achievement")),
		coalesceStr(achievement.Description, ""),
		game.Name,
		gameURL,
		unlockDate.Format("January 2, 2006 at 3:04 PM"),
	)

	// Validate frontmatter before writing
	if !validateFrontmatter(content) {
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "❌ Invalid frontmatter in achievement post %s - skipping\n", filename)
		}
		return false
	}

	// Write file (0644 is appropriate for user-editable content files)
	if err := os.WriteFile(postPath, []byte(content), 0o644); err != nil { //nolint:gosec // G306
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "❌ Failed to write achievement post %s: %v\n", filename, err)
		}
		return false
	}

	if steamVerbose {
		fmt.Fprintf(os.Stderr, "📝 Created achievement post: %s\n", filename)
	}
	return true
}

// createGamePost creates a comprehensive game post with all achievements.
func createGamePost(game *models.SteamGame, config *models.SteamConfig) bool {
	// Create safe filename
	gameName := sanitizeForFilename(game.Name)
	filename := fmt.Sprintf("%s.md", gameName)
	postPath := filepath.Join(config.PostsDir, filename)

	// Check if file already exists
	if _, err := os.Stat(postPath); err == nil {
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "⏭️  Game post already exists: %s\n", filename)
		}
		return false
	}

	// Get dates
	lastPlayed := game.LastPlayedDate()
	lastPlayedStr := ""
	if lastPlayed != nil {
		lastPlayedStr = lastPlayed.Format("2006-01-02")
	} else {
		lastPlayedStr = time.Now().Format("2006-01-02")
	}

	// Create achievements JSON string for frontmatter
	achievementsStr, err := json.Marshal(game.Achievements)
	if err != nil {
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "⚠️  Failed to marshal achievements for %s: %v\n", game.Name, err)
		}
		achievementsStr = []byte("[]")
	}

	// Create description with proper escaping
	description := fmt.Sprintf("Steam achievements and progress for %s - %.1f%% complete with %d/%d achievements unlocked.",
		game.Name, game.CompletionPercentage, game.UnlockedAchievements, game.TotalAchievements)

	// Create content
	content := fmt.Sprintf(`---
title: %s
description: %s
date: "%s"
published: true
templateKey: %s
steam:
  game: %s
  app_id: %d
  total_achievements: %d
  unlocked_achievements: %d
  completion_percentage: %.2f
  playtime_hours: %.1f
  last_played: %s
  description: %s
  developers: %s
  publishers: %s
  achievements: %s
tags: %s
slug: "steam/%s"
---

<style>
.game-header {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 30px;
  margin: 30px 0;
  padding: 20px;
  background: #1a1a1a;
  border-radius: 12px;
  border: 1px solid #333;
}

.game-header img {
  width: 200px;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  border: 1px solid #333;
  flex-shrink: 0;
}

.game-info {
  flex: 1;
}

.game-info h1 {
  margin: 0 0 15px 0;
  color: #fff;
  font-size: 2em;
}

.game-info p {
  margin: 0 0 15px 0;
  color: #ccc;
  line-height: 1.5;
}

.game-info .developers {
  font-size: 0.9em;
  color: #999;
}

.steam-game-progress {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 20px;
  margin: 20px 0;
  border: 1px solid #333;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin: 20px 0;
}

.stat-card {
  background: #2a2a2a;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
  border: 1px solid #444;
}

.stat-card h3 {
  margin: 0 0 15px 0;
  color: #4caf50;
  font-size: 1.1em;
}

.stat-value {
  font-size: 2em;
  font-weight: bold;
  color: #fff;
  margin: 10px 0;
}

.stat-card p {
  margin: 10px 0 0 0;
  color: #ccc;
  font-size: 0.9em;
}

.progress-bar {
  width: 100%%;
  height: 24px;
  background: #2a2a2a;
  border-radius: 12px;
  overflow: hidden;
  margin: 10px 0;
  position: relative;
}

.progress-fill {
  height: 100%%;
  background: linear-gradient(90deg, #4caf50, #8bc34a);
  border-radius: 12px;
  transition: width 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 12px;
}

.achievements-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(125px, 1fr));
  gap: 8px;
  margin: 20px 0;
}

.achievement-item {
  position: relative;
  text-align: center;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.achievement-item:hover {
  transform: scale(1.1);
  z-index: 10;
}

.achievement-icon-wrapper {
}

.achievement-icon {
  margin: 0;
  padding: 0;
  border-radius: 6px;
  border: 2px solid #444;
  transition: border-color 0.2s ease;
}

.achievement-item.unlocked .achievement-icon {
  border-color: #4caf50;
  box-shadow: 0 0 10px rgba(76, 175, 80, 0.3);
}

.achievement-item.locked .achievement-icon {
  filter: grayscale(100%%);
  opacity: 0.6;
}

.achievement-tooltip {
  position: absolute;
  bottom: 100%%;
  left: 50%%;
  transform: translateX(-50%%);
  background: rgba(0, 0, 0, 0.95);
  color: white;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  white-space: nowrap;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s ease;
  z-index: 100;
  margin-bottom: 5px;
  max-width: 200px;
  white-space: normal;
  text-align: center;
}

.achievement-item:hover .achievement-tooltip {
  opacity: 1;
}

.achievement-section {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 20px;
  margin: 20px 0;
  border: 1px solid #333;
}

.achievement-section h2 {
  margin-top: 0;
  color: #fff;
}
</style>

<div class="game-header">
  <img src="https://cdn.akamai.steamstatic.com/steam/apps/%d/library_600x900.jpg"
       alt="%s box art" loading="lazy"
       onerror="this.src='https://cdn.akamai.steamstatic.com/steam/apps/%d/header.jpg'">
  <div class="game-info">
    <h1>%s</h1>
    %s
    %s
  </div>
</div>

<div class="steam-game-progress">
<h2>📊 Game Progress & Stats</h2>

<div class="stats-grid">
  <div class="stat-card">
    <h3>Achievements</h3>
    <div class="progress-bar">
      <div class="progress-fill" style="width: %.1f%%">
        %.1f%%
      </div>
    </div>
    <p>%d/%d Unlocked</p>
  </div>

  <div class="stat-card">
    <h3>Playtime</h3>
    <div class="stat-value">%.1fh</div>
    <p>Total hours played</p>
  </div>

  %s
</div>
</div>

%s

%s

---

*Game data automatically imported from Steam. Achievement links will be created as individual posts when achievements are unlocked.*`,
		escapeYAML(game.Name),
		escapeYAML(description),
		lastPlayedStr,
		config.Template,
		escapeYAML(game.Name),
		game.AppID,
		game.TotalAchievements,
		game.UnlockedAchievements,
		game.CompletionPercentage,
		game.PlaytimeHours(),
		escapeYAML(lastPlayedStr),
		escapeYAML(coalesceStr(game.Description, "")),
		escapeYAMLArray(game.Developers),
		escapeYAMLArray(game.Publishers),
		escapeYAML(string(achievementsStr)),
		escapeYAMLArray([]string{"steam-game", "steam", "game", gameName}),
		gameName,
		game.AppID,
		html.EscapeString(game.Name),
		game.AppID,
		html.EscapeString(game.Name),
		gameDescriptionHTML(game.Description),
		gameDevelopersHTML(game.Developers),
		game.CompletionPercentage,
		game.CompletionPercentage,
		game.UnlockedAchievements,
		game.TotalAchievements,
		game.PlaytimeHours(),
		lastPlayedCardHTML(lastPlayed),
		unlockedAchievementsSection(game),
		lockedAchievementsSection(game),
	)

	// Validate frontmatter before writing
	if !validateFrontmatter(content) {
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "❌ Invalid frontmatter in game post %s - skipping\n", filename)
			// Debug: show the frontmatter that failed
			parts := strings.SplitN(content, "---", 3)
			if len(parts) >= 2 {
				fmt.Fprintf(os.Stderr, "Problematic frontmatter:\n%s\n", parts[1])
			}
		}
		return false
	}

	// Write file (0644 is appropriate for user-editable content files)
	if err := os.WriteFile(postPath, []byte(content), 0o644); err != nil { //nolint:gosec // G306
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "❌ Failed to write game post %s: %v\n", filename, err)
		}
		return false
	}

	if steamVerbose {
		fmt.Fprintf(os.Stderr, "📝 Created game post: %s\n", filename)
	}
	return true
}

// validateFrontmatter checks if the frontmatter is valid YAML
func validateFrontmatter(content string) bool {
	// Extract frontmatter between --- markers
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return false
	}

	frontmatter := strings.TrimSpace(parts[1])
	if frontmatter == "" {
		return false
	}

	// Try to parse as YAML
	var data interface{}
	err := yaml.Unmarshal([]byte(frontmatter), &data)
	if err != nil {
		if steamVerbose {
			fmt.Fprintf(os.Stderr, "YAML validation error: %v\n", err)
		}
		return false
	}

	return true
}

// Helper functions

func sanitizeForFilename(s string) string {
	// Replace invalid characters with dashes
	result := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, s)

	// Replace multiple consecutive dashes with a single dash
	result = strings.Join(strings.Fields(result), "-")
	return strings.ToLower(result)
}

func coalesceStr(s *string, def string) string {
	if s != nil {
		return *s
	}
	return def
}

func coalesceInt64(i *int64, def int64) int64 {
	if i != nil {
		return *i
	}
	return def
}

// yamlSpecialChars contains characters that require quoting in YAML values
var yamlSpecialChars = map[rune]bool{
	':': true, '#': true, '[': true, ']': true, '{': true, '}': true,
	'&': true, '*': true, '!': true, '|': true, '>': true, '\'': true,
	'"': true, '%': true, '@': true, '`': true, ',': true, '?': true,
	'\n': true, '\r': true, '\t': true,
}

// yamlStartChars contains characters that require quoting when at string start
var yamlStartChars = map[byte]bool{
	' ': true, '\t': true, '-': true, '?': true, ':': true, '[': true,
	'{': true, '!': true, '&': true, '*': true, '#': true, '|': true,
	'>': true, '\'': true, '"': true, '%': true, '@': true, '`': true,
}

// yamlReservedWords contains YAML boolean/null values that need quoting
var yamlReservedWords = map[string]bool{
	"true": true, "false": true, "null": true, "yes": true,
	"no": true, "on": true, "off": true, "~": true,
}

// escapeYAML properly escapes a string for safe YAML output.
// It handles all YAML special characters including colons, hashes, brackets,
// quotes, and other reserved characters that could cause parsing failures.
func escapeYAML(s string) string {
	if s == "" {
		return `""`
	}

	needsQuoting := yamlNeedsQuoting(s)
	if !needsQuoting {
		return s
	}
	return quoteYAMLString(s)
}

// yamlNeedsQuoting checks if a string needs to be quoted for YAML
func yamlNeedsQuoting(s string) bool {
	// Check for special characters in content
	for _, r := range s {
		if yamlSpecialChars[r] {
			return true
		}
	}

	// Check first/last character
	if yamlStartChars[s[0]] || s[len(s)-1] == ' ' || s[len(s)-1] == '\t' {
		return true
	}

	// Check for reserved words (case insensitive)
	return yamlReservedWords[strings.ToLower(s)]
}

// quoteYAMLString wraps a string in double quotes and escapes special chars
func quoteYAMLString(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 10) // Pre-allocate with some extra space for escapes
	sb.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			sb.WriteRune(r)
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

func escapeYAMLArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}

	escaped := make([]string, len(arr))
	for i, s := range arr {
		escaped[i] = escapeYAML(s)
	}

	return fmt.Sprintf("[%s]", strings.Join(escaped, ", "))
}

func achievementIconHTML(icon *string, name string) string {
	iconURL := coalesceStr(icon, "")
	if iconURL == "" {
		return ""
	}
	// Using string concat to avoid gocritic sprintfQuotedString lint warning
	// We need literal HTML attribute quotes, not Go %q escaping
	return `<img src="` + html.EscapeString(iconURL) + `" alt="` + html.EscapeString(name) + `" style="width: 64px; height: 64px;">`
}

func gameDescriptionHTML(desc *string) string {
	descStr := coalesceStr(desc, "")
	if descStr == "" {
		return ""
	}
	return fmt.Sprintf(`<p><em>%s</em></p>`, html.EscapeString(descStr))
}

func gameDevelopersHTML(devs []string) string {
	if len(devs) == 0 {
		return ""
	}

	// Show max 2 developers
	displayDevs := devs
	if len(devs) > 2 {
		displayDevs = devs[:2]
	}

	// Escape each developer name
	escapedDevs := make([]string, len(displayDevs))
	for i, dev := range displayDevs {
		escapedDevs[i] = html.EscapeString(dev)
	}
	devStr := strings.Join(escapedDevs, ", ")
	if len(devs) > 2 {
		devStr += "..."
	}

	return fmt.Sprintf(`<p class="developers">Developed by %s</p>`, devStr)
}

func lastPlayedCardHTML(lastPlayed *time.Time) string {
	if lastPlayed == nil {
		return ""
	}

	return fmt.Sprintf(`
<div class="stat-card">
    <h3>Last Played</h3>
    <div class="stat-value">%s</div>
    <p>Most recent session</p>
  </div>`, lastPlayed.Format("2006-01-02"))
}

func unlockedAchievementsSection(game *models.SteamGame) string {
	unlocked := game.UnlockedAchievementsList()
	if len(unlocked) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
<div class="achievement-section">
<h2>🏆 Unlocked Achievements (%d)</h2>

<div class="achievements-grid">
`, len(unlocked)))

	for _, achievement := range unlocked {
		unlockDate := achievement.UnlockDate()
		dateStr := ""
		if unlockDate != nil {
			dateStr = unlockDate.Format("January 2, 2006")
		}

		badgeURL := html.EscapeString(coalesceStr(achievement.Icon, ""))
		achievementName := html.EscapeString(coalesceStr(achievement.Name, "Unknown Achievement"))
		description := html.EscapeString(coalesceStr(achievement.Description, "No description"))

		sb.WriteString(fmt.Sprintf(`
<div class="achievement-item unlocked">
  <span class="achievement-icon-wrapper">
    <img src="%s" alt="%s" class="achievement-icon">
  </span>
  <div class="achievement-tooltip">
    <strong>%s</strong><br>
    %s<br>
    <small>Unlocked: %s</small>
  </div>
</div>`, badgeURL, achievementName, achievementName, description, dateStr))
	}

	sb.WriteString(`
</div>
</div>
`)
	return sb.String()
}

func lockedAchievementsSection(game *models.SteamGame) string {
	locked := game.LockedAchievementsList()
	if len(locked) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
<div class="achievement-section">
<h2>🔒 Locked Achievements (%d)</h2>

<div class="achievements-grid">
`, len(locked)))

	for _, achievement := range locked {
		badgeURL := html.EscapeString(coalesceStr(achievement.IconGray, coalesceStr(achievement.Icon, "")))
		achievementName := html.EscapeString(coalesceStr(achievement.Name, "Unknown Achievement"))
		description := html.EscapeString(coalesceStr(achievement.Description, "No description"))

		sb.WriteString(fmt.Sprintf(`
<div class="achievement-item locked">
  <span class="achievement-icon-wrapper">
    <img src="%s" alt="%s" class="achievement-icon">
  </span>
  <div class="achievement-tooltip">
    <strong>%s</strong><br>
    %s
  </div>
</div>`, badgeURL, achievementName, achievementName, description))
	}

	sb.WriteString(`
</div>
</div>
`)
	return sb.String()
}
