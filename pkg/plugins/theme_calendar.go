package plugins

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// ThemeCalendarPlugin applies theme rules based on the current date.
// It runs during the Configure stage to modify theme settings based on
// date-based rules defined in the configuration.
//
// Configuration example:
//
//	[markata-go.theme_calendar]
//	enabled = true
//
//	[[markata-go.theme_calendar.rules]]
//	name = "Christmas Season"
//	start_date = "12-15"
//	end_date = "12-26"
//	palette = "christmas"
//
//	[[markata-go.theme_calendar.rules]]
//	name = "Winter Frost"
//	start_date = "12-01"
//	end_date = "02-28"
//	palette = "winter-frost"
type ThemeCalendarPlugin struct {
	// nowFunc allows injecting a custom time function for testing
	nowFunc func() time.Time
}

// NewThemeCalendarPlugin creates a new ThemeCalendarPlugin.
func NewThemeCalendarPlugin() *ThemeCalendarPlugin {
	return &ThemeCalendarPlugin{
		nowFunc: time.Now,
	}
}

// Name returns the unique name of the plugin.
func (p *ThemeCalendarPlugin) Name() string {
	return "theme_calendar"
}

// Priority returns a high priority so theme calendar runs before other theme plugins.
// This ensures the theme is configured before palette_css and other plugins process it.
func (p *ThemeCalendarPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageConfigure {
		return -100 // Run very early in Configure stage
	}
	return 0 // Default priority for other stages
}

// Configure checks date-based theme rules and applies matching overrides.
func (p *ThemeCalendarPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Get theme calendar config from Config struct
	calendarConfig := p.getCalendarConfig(config)
	if !calendarConfig.IsEnabled() || len(calendarConfig.Rules) == 0 {
		return nil
	}

	// Get current date (or test date)
	now := p.nowFunc()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	log.Printf("[theme_calendar] Checking %d rules for date %02d-%02d", len(calendarConfig.Rules), currentMonth, currentDay)

	// Find matching rule
	var matchingRule *models.ThemeCalendarRule
	for i := range calendarConfig.Rules {
		rule := &calendarConfig.Rules[i]
		if p.isDateInRange(currentMonth, currentDay, rule.StartDate, rule.EndDate) {
			matchingRule = rule
			log.Printf("[theme_calendar] Matched rule: %s", rule.Name)
			break
		}
	}

	if matchingRule == nil {
		log.Printf("[theme_calendar] No matching rule found, using base theme")
		return nil
	}

	// Apply the matching rule to theme config
	p.applyRule(config, matchingRule)

	return nil
}

// getCalendarConfig retrieves the theme calendar configuration.
func (p *ThemeCalendarPlugin) getCalendarConfig(config *lifecycle.Config) *models.ThemeCalendarConfig {
	// Prefer the typed configuration projected by the CLI lifecycle setup.
	if config != nil && config.Extra != nil {
		switch calendar := config.Extra["theme_calendar"].(type) {
		case models.ThemeCalendarConfig:
			return &calendar
		case *models.ThemeCalendarConfig:
			if calendar != nil {
				return calendar
			}
		}
	}

	if config == nil || config.Extra == nil {
		return &models.ThemeCalendarConfig{}
	}

	// Check if theme_calendar exists in Extra (from TOML parsing)
	if cal, ok := config.Extra["theme_calendar"]; ok {
		if calMap, ok := cal.(map[string]interface{}); ok {
			return p.parseCalendarConfig(calMap)
		}
	}

	return &models.ThemeCalendarConfig{}
}

// parseCalendarConfig converts a map to ThemeCalendarConfig.
func (p *ThemeCalendarPlugin) parseCalendarConfig(m map[string]interface{}) *models.ThemeCalendarConfig {
	cfg := &models.ThemeCalendarConfig{}

	if enabled, ok := m["enabled"].(bool); ok {
		cfg.Enabled = &enabled
	}

	if defaultPalette, ok := m["default_palette"].(string); ok {
		cfg.DefaultPalette = defaultPalette
	}

	if rules, ok := m["rules"].([]interface{}); ok {
		for _, r := range rules {
			if ruleMap, ok := r.(map[string]interface{}); ok {
				rule := p.parseRule(ruleMap)
				cfg.Rules = append(cfg.Rules, rule)
			}
		}
	}

	return cfg
}

// parseRule converts a map to ThemeCalendarRule.
func (p *ThemeCalendarPlugin) parseRule(m map[string]interface{}) models.ThemeCalendarRule {
	rule := models.ThemeCalendarRule{}

	if name, ok := m["name"].(string); ok {
		rule.Name = name
	}
	if startDate, ok := m["start_date"].(string); ok {
		rule.StartDate = startDate
	}
	if endDate, ok := m["end_date"].(string); ok {
		rule.EndDate = endDate
	}
	if palette, ok := m["palette"].(string); ok {
		rule.Palette = palette
	}
	if paletteLight, ok := m["palette_light"].(string); ok {
		rule.PaletteLight = paletteLight
	}
	if paletteDark, ok := m["palette_dark"].(string); ok {
		rule.PaletteDark = paletteDark
	}
	if customCSS, ok := m["custom_css"].(string); ok {
		rule.CustomCSS = customCSS
	}
	if variables, ok := m["variables"].(map[string]interface{}); ok {
		rule.Variables = make(map[string]string)
		for k, v := range variables {
			if str, ok := v.(string); ok {
				rule.Variables[k] = str
			}
		}
	}

	// Parse background config if present
	if bg, ok := m["background"].(map[string]interface{}); ok {
		rule.Background = p.parseBackgroundConfig(bg)
	}

	// Parse font config if present
	if font, ok := m["font"].(map[string]interface{}); ok {
		rule.Font = p.parseFontConfig(font)
	}

	return rule
}

// parseBackgroundConfig converts a map to BackgroundConfig.
func (p *ThemeCalendarPlugin) parseBackgroundConfig(m map[string]interface{}) *models.BackgroundConfig {
	bg := &models.BackgroundConfig{}

	if enabled, ok := m["enabled"].(bool); ok {
		bg.Enabled = &enabled
	}
	if css, ok := m["css"].(string); ok {
		bg.CSS = css
	}
	if scripts, ok := m["scripts"].([]interface{}); ok {
		for _, s := range scripts {
			if str, ok := s.(string); ok {
				bg.Scripts = append(bg.Scripts, str)
			}
		}
	}
	if backgrounds, ok := m["backgrounds"].([]interface{}); ok {
		for _, b := range backgrounds {
			if bgMap, ok := b.(map[string]interface{}); ok {
				elem := models.BackgroundElement{}
				if html, ok := bgMap["html"].(string); ok {
					elem.HTML = html
				}
				if zIndex, ok := parseIntFromInterface(bgMap["z_index"]); ok {
					elem.ZIndex = zIndex
				}
				bg.Backgrounds = append(bg.Backgrounds, elem)
			}
		}
	}

	return bg
}

// parseFontConfig converts a map to FontConfig.
func (p *ThemeCalendarPlugin) parseFontConfig(m map[string]interface{}) *models.FontConfig {
	font := &models.FontConfig{}

	if family, ok := m["family"].(string); ok {
		font.Family = family
	}
	if headingFamily, ok := m["heading_family"].(string); ok {
		font.HeadingFamily = headingFamily
	}
	if codeFamily, ok := m["code_family"].(string); ok {
		font.CodeFamily = codeFamily
	}
	if size, ok := m["size"].(string); ok {
		font.Size = size
	}
	if lineHeight, ok := m["line_height"].(string); ok {
		font.LineHeight = lineHeight
	}
	if googleFonts, ok := m["google_fonts"].([]interface{}); ok {
		for _, f := range googleFonts {
			if str, ok := f.(string); ok {
				font.GoogleFonts = append(font.GoogleFonts, str)
			}
		}
	}
	if customURLs, ok := m["custom_urls"].([]interface{}); ok {
		for _, u := range customURLs {
			if str, ok := u.(string); ok {
				font.CustomURLs = append(font.CustomURLs, str)
			}
		}
	}

	return font
}

// isDateInRange checks if the given month/day falls within the start-end range.
// Handles year-boundary crossings (e.g., Dec 1 to Feb 28).
func (p *ThemeCalendarPlugin) isDateInRange(month, day int, startDate, endDate string) bool {
	startMonth, startDay, err := p.parseMMDD(startDate)
	if err != nil {
		log.Printf("[theme_calendar] Invalid start_date %q: %v", startDate, err)
		return false
	}

	endMonth, endDay, err := p.parseMMDD(endDate)
	if err != nil {
		log.Printf("[theme_calendar] Invalid end_date %q: %v", endDate, err)
		return false
	}

	// Convert to day-of-year-like number for comparison (month*100 + day)
	current := month*100 + day
	start := startMonth*100 + startDay
	end := endMonth*100 + endDay

	if start <= end {
		// Simple range (e.g., Mar 1 to May 31)
		return current >= start && current <= end
	}

	// Range crosses year boundary (e.g., Dec 1 to Feb 28)
	// Either current >= start (Dec-onwards) OR current <= end (Jan-Feb)
	return current >= start || current <= end
}

// parseMMDD parses a date string in MM-DD format.
func (p *ThemeCalendarPlugin) parseMMDD(s string) (month, day int, err error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected MM-DD format, got %q", s)
	}

	month, err = strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid month in %q", s)
	}

	day, err = strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 31 {
		return 0, 0, fmt.Errorf("invalid day in %q", s)
	}

	return month, day, nil
}

// applyRule applies the matching rule's theme overrides to the config.
func (p *ThemeCalendarPlugin) applyRule(config *lifecycle.Config, rule *models.ThemeCalendarRule) {
	// Get or create theme map in Extra
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}

	themeMap := p.getOrCreateThemeMap(config)

	// Apply palette overrides
	p.applyPaletteOverrides(themeMap, rule)

	// Apply custom CSS override
	if rule.CustomCSS != "" {
		themeMap["custom_css"] = rule.CustomCSS
		log.Printf("[theme_calendar] Setting custom_css to %q", rule.CustomCSS)
	}

	// Merge variables (deep merge with existing)
	p.applyVariablesOverride(themeMap, rule)

	// Apply background override (replace entirely if set)
	p.applyBackgroundOverride(themeMap, rule)

	// Apply font override (merge with existing)
	p.applyFontOverride(themeMap, rule)

	// Store the updated theme map
	config.Extra["theme"] = themeMap

	// Log the active rule
	log.Printf("[theme_calendar] Applied rule %q for current date", rule.Name)
}

// getOrCreateThemeMap retrieves existing theme map or creates a new one.
func (p *ThemeCalendarPlugin) getOrCreateThemeMap(config *lifecycle.Config) map[string]interface{} {
	if existing, ok := config.Extra["theme"].(map[string]interface{}); ok {
		return existing
	}
	return make(map[string]interface{})
}

// applyPaletteOverrides applies palette-related overrides from the rule.
func (p *ThemeCalendarPlugin) applyPaletteOverrides(themeMap map[string]interface{}, rule *models.ThemeCalendarRule) {
	if rule.Palette != "" {
		themeMap["palette"] = rule.Palette
		log.Printf("[theme_calendar] Setting palette to %q", rule.Palette)
	}
	if rule.PaletteLight != "" {
		themeMap["palette_light"] = rule.PaletteLight
		log.Printf("[theme_calendar] Setting palette_light to %q", rule.PaletteLight)
	}
	if rule.PaletteDark != "" {
		themeMap["palette_dark"] = rule.PaletteDark
		log.Printf("[theme_calendar] Setting palette_dark to %q", rule.PaletteDark)
	}
}

// applyVariablesOverride merges CSS variables from the rule.
func (p *ThemeCalendarPlugin) applyVariablesOverride(themeMap map[string]interface{}, rule *models.ThemeCalendarRule) {
	if len(rule.Variables) == 0 {
		return
	}
	existingVars, ok := themeMap["variables"].(map[string]interface{})
	if !ok || existingVars == nil {
		existingVars = make(map[string]interface{})
	}
	for k, v := range rule.Variables {
		existingVars[k] = v
	}
	themeMap["variables"] = existingVars
	log.Printf("[theme_calendar] Merged %d CSS variables", len(rule.Variables))
}

// applyBackgroundOverride applies background decoration override from the rule.
func (p *ThemeCalendarPlugin) applyBackgroundOverride(themeMap map[string]interface{}, rule *models.ThemeCalendarRule) {
	if rule.Background == nil || rule.Background.Enabled == nil || !*rule.Background.Enabled {
		return
	}
	bgMap := make(map[string]interface{})
	bgMap["enabled"] = true
	if rule.Background.CSS != "" {
		bgMap["css"] = rule.Background.CSS
	}
	if len(rule.Background.Scripts) > 0 {
		bgMap["scripts"] = rule.Background.Scripts
	}
	if len(rule.Background.Backgrounds) > 0 {
		bgs := make([]map[string]interface{}, 0, len(rule.Background.Backgrounds))
		for _, b := range rule.Background.Backgrounds {
			bgs = append(bgs, map[string]interface{}{
				"html":    b.HTML,
				"z_index": b.ZIndex,
			})
		}
		bgMap["backgrounds"] = bgs
	}
	themeMap["background"] = bgMap
	log.Printf("[theme_calendar] Applied background override")
}

// applyFontOverride applies font configuration override from the rule.
func (p *ThemeCalendarPlugin) applyFontOverride(themeMap map[string]interface{}, rule *models.ThemeCalendarRule) {
	if rule.Font == nil {
		return
	}
	existingFont, ok := themeMap["font"].(map[string]interface{})
	if !ok || existingFont == nil {
		existingFont = make(map[string]interface{})
	}
	if rule.Font.Family != "" {
		existingFont["family"] = rule.Font.Family
	}
	if rule.Font.HeadingFamily != "" {
		existingFont["heading_family"] = rule.Font.HeadingFamily
	}
	if rule.Font.CodeFamily != "" {
		existingFont["code_family"] = rule.Font.CodeFamily
	}
	if rule.Font.Size != "" {
		existingFont["size"] = rule.Font.Size
	}
	if rule.Font.LineHeight != "" {
		existingFont["line_height"] = rule.Font.LineHeight
	}
	if len(rule.Font.GoogleFonts) > 0 {
		existingFont["google_fonts"] = rule.Font.GoogleFonts
	}
	if len(rule.Font.CustomURLs) > 0 {
		existingFont["custom_urls"] = rule.Font.CustomURLs
	}
	themeMap["font"] = existingFont
	log.Printf("[theme_calendar] Applied font override")
}

// Compile-time interface verification.
var (
	_ lifecycle.Plugin          = (*ThemeCalendarPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ThemeCalendarPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*ThemeCalendarPlugin)(nil)
)
