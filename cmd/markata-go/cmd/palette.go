package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/palettes"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

// Common string constants to avoid goconst warnings.
const (
	paletteFormatCSS      = "css"
	paletteFormatSCSS     = "scss"
	paletteFormatJSON     = "json"
	paletteFormatTailwind = "tailwind"

	// Default palette colors (shared with theme.go).
	defaultTextColor = "#e0e0e0"
	defaultBgColor   = "#1e1e1e"

	// Semantic key for text color mapping.
	colorKeyText = "text"
)

// paletteCmd represents the palette command group.
var paletteCmd = &cobra.Command{
	Use:   "palette",
	Short: "Color palette commands",
	Long: `Commands for managing color palettes.

Subcommands:
  list     - List available palettes
  info     - Show palette details
  check    - Validate palette contrast ratios
  preview  - Generate HTML preview
  export   - Export palette to different formats
  new      - Create a new palette
  clone    - Clone an existing palette with fuzzy picker
  fetch    - Fetch palette from Lospec URL`,
}

// paletteListCmd lists available palettes.
var paletteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available palettes",
	Long: `List all available color palettes from built-in, user, and project sources.

Palettes are discovered from:
  1. Built-in palettes (embedded in binary)
  2. User palettes (~/.config/markata-go/palettes/)
  3. Project palettes (./palettes/)

Example usage:
  markata-go palette list              # List all palettes
  markata-go palette list --variant dark  # List only dark palettes
  markata-go palette list --json       # Output as JSON`,
	RunE: runPaletteListCommand,
}

// paletteInfoCmd shows palette details.
var paletteInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show palette details",
	Long: `Display detailed information about a specific palette.

Shows metadata, color counts, and all color values.

Example usage:
  markata-go palette info catppuccin-mocha
  markata-go palette info "Catppuccin Mocha" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runPaletteInfoCommand,
}

// paletteCheckCmd validates contrast ratios.
var paletteCheckCmd = &cobra.Command{
	Use:   "check <name>",
	Short: "Validate palette contrast",
	Long: `Check WCAG contrast ratios for a palette.

Validates that foreground/background color combinations meet
accessibility requirements:
  - AA: 4.5:1 for normal text, 3:1 for large text
  - AAA: 7:1 for normal text, 4.5:1 for large text

Exit codes:
  0 - All checks passed
  1 - One or more checks failed

Example usage:
  markata-go palette check catppuccin-mocha
  markata-go palette check catppuccin-mocha --strict  # Include AAA checks
  markata-go palette check --all                      # Check all palettes`,
	RunE: runPaletteCheckCommand,
}

// paletteExportCmd exports palette to different formats.
var paletteExportCmd = &cobra.Command{
	Use:   "export <name>",
	Short: "Export palette",
	Long: `Export a palette to different formats.

Supported formats:
  css      - CSS custom properties
  scss     - SCSS/Sass variables
  json     - JSON with resolved colors
  tailwind - Tailwind CSS config

Example usage:
  markata-go palette export catppuccin-mocha --format css
  markata-go palette export nord-dark --format tailwind > tailwind.colors.js`,
	Args: cobra.ExactArgs(1),
	RunE: runPaletteExportCommand,
}

// palettePreviewCmd generates HTML preview.
var palettePreviewCmd = &cobra.Command{
	Use:   "preview [name]",
	Short: "Generate palette preview",
	Long: `Generate an HTML preview page for a palette.

Creates a visual preview showing all colors with their names and hex values.

Example usage:
  markata-go palette preview catppuccin-mocha
  markata-go palette preview catppuccin-mocha --open  # Open in browser
  markata-go palette preview catppuccin-mocha -o preview.html
  markata-go palette preview --all                    # Preview all palettes`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPalettePreviewCommand,
}

// paletteNewCmd creates a new palette.
var paletteNewCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create new palette",
	Long: `Create a new palette file from a template.

The new palette can be based on an existing palette or start from scratch.

Example usage:
  markata-go palette new my-theme
  markata-go palette new my-theme --variant dark
  markata-go palette new my-theme --from catppuccin-mocha`,
	Args: cobra.ExactArgs(1),
	RunE: runPaletteNewCommand,
}

// paletteCloneCmd clones an existing palette with a fuzzy picker.
var paletteCloneCmd = &cobra.Command{
	Use:   "clone [source-palette]",
	Short: "Clone an existing palette",
	Long: `Clone an existing palette to create a customized version.

If no source palette is specified, opens a fuzzy picker to select one.
After selection, prompts for a name for the new palette.

The cloned palette is saved to ~/.config/markata-go/palettes/.

Example usage:
  markata-go palette clone                        # Fuzzy picker -> prompt for name
  markata-go palette clone catppuccin-mocha       # Clone specific palette, prompt for name
  markata-go palette clone catppuccin-mocha --name "my-custom"  # Clone with explicit name`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPaletteCloneCommand,
}

// paletteFetchCmd fetches a palette from a Lospec URL.
var paletteFetchCmd = &cobra.Command{
	Use:   "fetch <url>",
	Short: "Fetch palette from Lospec",
	Long: `Fetch a color palette from a Lospec.com URL.

Downloads the palette, generates semantic mappings, and saves it to your
user palettes directory (~/.config/markata-go/palettes/).

Supported URL format:
  https://lospec.com/palette-list/<palette-name>.txt

Example usage:
  markata-go palette fetch https://lospec.com/palette-list/cheese-palette.txt
  markata-go palette fetch https://lospec.com/palette-list/sweetie-16.txt --name "My Sweetie"
  markata-go palette fetch https://lospec.com/palette-list/tokyo-night.txt -o palettes/`,
	Args: cobra.ExactArgs(1),
	RunE: runPaletteFetchCommand,
}

var (
	// paletteListVariant filters list by light/dark variant.
	paletteListVariant string

	// paletteNewVariant sets variant for new palette.
	paletteNewVariant string

	// paletteJSON outputs as JSON.
	paletteJSON bool

	// paletteStrict includes AAA checks.
	paletteStrict bool

	// paletteAll checks all palettes.
	paletteAll bool

	// paletteFormat export format.
	paletteFormat string

	// paletteOutput output file.
	paletteOutput string

	// paletteFrom base palette for new.
	paletteFrom string

	// paletteOpen opens preview in browser.
	paletteOpen bool

	// palettePreviewAll shows all palettes in preview.
	palettePreviewAll bool

	// paletteCloneName is the name for the cloned palette.
	paletteCloneName string

	// paletteFetchName is the custom name for a fetched palette.
	paletteFetchName string

	// paletteFetchOutput is the output directory for a fetched palette.
	paletteFetchOutput string
)

func init() {
	rootCmd.AddCommand(paletteCmd)

	// List subcommand
	paletteCmd.AddCommand(paletteListCmd)
	paletteListCmd.Flags().StringVar(&paletteListVariant, "variant", "", "Filter by variant (light/dark)")
	paletteListCmd.Flags().BoolVar(&paletteJSON, "json", false, "Output as JSON")

	// Info subcommand
	paletteCmd.AddCommand(paletteInfoCmd)
	paletteInfoCmd.Flags().BoolVar(&paletteJSON, "json", false, "Output as JSON")

	// Check subcommand
	paletteCmd.AddCommand(paletteCheckCmd)
	paletteCheckCmd.Flags().BoolVar(&paletteStrict, "strict", false, "Include AAA level checks")
	paletteCheckCmd.Flags().BoolVar(&paletteAll, "all", false, "Check all palettes")
	paletteCheckCmd.Flags().BoolVar(&paletteJSON, "json", false, "Output as JSON")

	// Export subcommand
	paletteCmd.AddCommand(paletteExportCmd)
	paletteExportCmd.Flags().StringVarP(&paletteFormat, "format", "f", "css", "Export format (css, scss, json, tailwind)")
	paletteExportCmd.Flags().StringVarP(&paletteOutput, "output", "o", "", "Output file (default: stdout)")

	// Preview subcommand
	paletteCmd.AddCommand(palettePreviewCmd)
	palettePreviewCmd.Flags().StringVarP(&paletteOutput, "output", "o", "", "Output file (default: palette-preview.html)")
	palettePreviewCmd.Flags().BoolVar(&paletteOpen, "open", false, "Open preview in browser")
	palettePreviewCmd.Flags().BoolVar(&palettePreviewAll, "all", false, "Preview all available palettes")

	// New subcommand
	paletteCmd.AddCommand(paletteNewCmd)
	paletteNewCmd.Flags().StringVar(&paletteNewVariant, "variant", "dark", "Palette variant (light/dark)")
	paletteNewCmd.Flags().StringVar(&paletteFrom, "from", "", "Base palette to copy from")
	paletteNewCmd.Flags().StringVarP(&paletteOutput, "output", "o", "", "Output file (default: palettes/<name>.toml)")

	// Clone subcommand
	paletteCmd.AddCommand(paletteCloneCmd)
	paletteCloneCmd.Flags().StringVar(&paletteCloneName, "name", "", "Name for the cloned palette")

	// Fetch subcommand
	paletteCmd.AddCommand(paletteFetchCmd)
	paletteFetchCmd.Flags().StringVarP(&paletteFetchName, "name", "n", "", "Custom name for the palette")
	paletteFetchCmd.Flags().StringVarP(&paletteFetchOutput, "output", "o", "", "Output directory (default: ~/.config/markata-go/palettes/)")
}

// runPaletteListCommand lists available palettes.
func runPaletteListCommand(_ *cobra.Command, _ []string) error {
	loader := palettes.NewLoader()
	infos, err := loader.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover palettes: %w", err)
	}

	// Filter by variant if specified
	if paletteListVariant != "" {
		variant := palettes.Variant(paletteListVariant)
		filtered := make([]palettes.PaletteInfo, 0)
		for _, info := range infos {
			if info.Variant == variant {
				filtered = append(filtered, info)
			}
		}
		infos = filtered
	}

	if paletteJSON {
		data, err := json.MarshalIndent(infos, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Table output
	if len(infos) == 0 {
		fmt.Println("No palettes found.")
		return nil
	}

	fmt.Printf("%-25s %-8s %-10s %s\n", "NAME", "VARIANT", "SOURCE", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 70))
	for _, info := range infos {
		desc := info.Description
		if len(desc) > 25 {
			desc = desc[:22] + "..."
		}
		fmt.Printf("%-25s %-8s %-10s %s\n", info.Name, info.Variant, info.Source, desc)
	}

	return nil
}

// runPaletteInfoCommand shows palette details.
func runPaletteInfoCommand(_ *cobra.Command, args []string) error {
	name := args[0]

	loader := palettes.NewLoader()
	p, err := loader.Load(name)
	if err != nil {
		return fmt.Errorf("failed to load palette: %w", err)
	}

	if paletteJSON {
		data, err := p.ExportJSON(true)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Human-readable output
	fmt.Printf("Name:        %s\n", p.Name)
	fmt.Printf("Variant:     %s\n", p.Variant)
	if p.Author != "" {
		fmt.Printf("Author:      %s\n", p.Author)
	}
	if p.License != "" {
		fmt.Printf("License:     %s\n", p.License)
	}
	if p.Homepage != "" {
		fmt.Printf("Homepage:    %s\n", p.Homepage)
	}
	if p.Description != "" {
		fmt.Printf("Description: %s\n", p.Description)
	}
	fmt.Printf("Source:      %s\n", p.Source)
	if p.SourcePath != "" {
		fmt.Printf("Path:        %s\n", p.SourcePath)
	}
	fmt.Println()

	fmt.Printf("Colors:      %d raw, %d semantic, %d component\n",
		len(p.Colors), len(p.Semantic), len(p.Components))
	fmt.Println()

	// Show raw colors
	fmt.Println("Raw Colors:")
	for name, hex := range p.Colors {
		fmt.Printf("  %-15s %s\n", name, hex)
	}
	fmt.Println()

	// Show semantic colors
	fmt.Println("Semantic Colors:")
	for name, ref := range p.Semantic {
		hex := p.Resolve(name)
		fmt.Printf("  %-20s -> %-15s (%s)\n", name, ref, hex)
	}

	if len(p.Components) > 0 {
		fmt.Println()
		fmt.Println("Component Colors:")
		for name, ref := range p.Components {
			hex := p.Resolve(name)
			fmt.Printf("  %-25s -> %-15s (%s)\n", name, ref, hex)
		}
	}

	return nil
}

// runPaletteCheckCommand validates contrast ratios.
func runPaletteCheckCommand(_ *cobra.Command, args []string) error {
	loader := palettes.NewLoader()

	var palettesToCheck []*palettes.Palette

	if paletteAll {
		// Check all discovered palettes
		infos, err := loader.Discover()
		if err != nil {
			return fmt.Errorf("failed to discover palettes: %w", err)
		}
		for _, info := range infos {
			p, err := loader.Load(info.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load %s: %v\n", info.Name, err)
				continue
			}
			palettesToCheck = append(palettesToCheck, p)
		}
	} else {
		if len(args) == 0 {
			return fmt.Errorf("palette name required (or use --all)")
		}
		p, err := loader.Load(args[0])
		if err != nil {
			return fmt.Errorf("failed to load palette: %w", err)
		}
		palettesToCheck = append(palettesToCheck, p)
	}

	allPassed := true
	var allResults []map[string]interface{}

	for _, p := range palettesToCheck {
		var results []palettes.ContrastCheck
		if paletteStrict {
			results = p.CheckContrastStrict()
		} else {
			results = p.CheckContrast()
		}

		summary := palettes.SummarizeContrast(p.Name, results)
		if !summary.AllPassed {
			allPassed = false
		}

		if paletteJSON {
			allResults = append(allResults, map[string]interface{}{
				"palette": p.Name,
				"summary": summary,
				"results": results,
			})
		} else {
			fmt.Printf("\nChecking palette: %s\n", p.Name)
			fmt.Println(strings.Repeat("-", 50))

			for i := range results {
				fmt.Println(palettes.FormatContrastResult(results[i]))
			}

			fmt.Println()
			if summary.AllPassed {
				fmt.Printf("All %d checks passed!\n", summary.Passed)
			} else {
				fmt.Printf("Results: %d passed, %d failed, %d skipped\n",
					summary.Passed, summary.Failed, summary.Skipped)
			}
		}
	}

	if paletteJSON {
		data, err := json.MarshalIndent(allResults, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
	}

	if !allPassed {
		return fmt.Errorf("contrast checks failed")
	}

	return nil
}

// runPaletteExportCommand exports palette to different formats.
func runPaletteExportCommand(_ *cobra.Command, args []string) error {
	name := args[0]

	loader := palettes.NewLoader()
	p, err := loader.Load(name)
	if err != nil {
		return fmt.Errorf("failed to load palette: %w", err)
	}

	var output string
	switch paletteFormat {
	case paletteFormatCSS:
		output = p.GenerateCSS()
	case paletteFormatSCSS:
		output = p.GenerateSCSS()
	case paletteFormatJSON:
		data, err := p.ExportJSON(true)
		if err != nil {
			return fmt.Errorf("failed to generate JSON: %w", err)
		}
		output = string(data)
	case paletteFormatTailwind:
		output = p.GenerateTailwind()
	default:
		return fmt.Errorf("unknown format: %s (supported: css, scss, json, tailwind)", paletteFormat)
	}

	if paletteOutput != "" {
		if err := os.WriteFile(paletteOutput, []byte(output), 0o644); err != nil { //nolint:gosec // exported files should be readable
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("Exported to %s\n", paletteOutput)
	} else {
		fmt.Print(output)
	}

	return nil
}

// runPalettePreviewCommand generates HTML preview.
func runPalettePreviewCommand(_ *cobra.Command, args []string) error {
	loader := palettes.NewLoader()

	var palettesToPreview []*palettes.Palette

	if palettePreviewAll {
		// Load all discovered palettes
		infos, err := loader.Discover()
		if err != nil {
			return fmt.Errorf("failed to discover palettes: %w", err)
		}
		if len(infos) == 0 {
			return fmt.Errorf("no palettes found")
		}
		for _, info := range infos {
			p, err := loader.Load(info.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load %s: %v\n", info.Name, err)
				continue
			}
			palettesToPreview = append(palettesToPreview, p)
		}
		if len(palettesToPreview) == 0 {
			return fmt.Errorf("failed to load any palettes")
		}
	} else {
		if len(args) == 0 {
			return fmt.Errorf("palette name required (or use --all)")
		}
		name := args[0]
		p, err := loader.Load(name)
		if err != nil {
			return fmt.Errorf("failed to load palette: %w", err)
		}
		palettesToPreview = append(palettesToPreview, p)
	}

	var html string
	if len(palettesToPreview) == 1 {
		html = generatePreviewHTML(palettesToPreview[0])
	} else {
		html = generateAllPalettesPreviewHTML(palettesToPreview)
	}

	outputFile := paletteOutput
	if outputFile == "" {
		if palettePreviewAll {
			outputFile = "all-palettes-preview.html"
		} else {
			outputFile = "palette-preview.html"
		}
	}

	if err := os.WriteFile(outputFile, []byte(html), 0o644); err != nil { //nolint:gosec // preview files should be readable
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Preview generated: %s\n", outputFile)

	if paletteOpen {
		// Try to open in browser
		openBrowser(outputFile)
	}

	return nil
}

// runPaletteNewCommand creates a new palette.
func runPaletteNewCommand(_ *cobra.Command, args []string) error {
	name := args[0]

	var p *palettes.Palette

	if paletteFrom != "" {
		// Load base palette
		loader := palettes.NewLoader()
		base, err := loader.Load(paletteFrom)
		if err != nil {
			return fmt.Errorf("failed to load base palette: %w", err)
		}
		p = base.Clone()
		p.Name = name
		p.Author = ""
		p.Homepage = ""
		p.Description = fmt.Sprintf("Based on %s", paletteFrom)
	} else {
		// Create minimal palette
		p = palettes.NewPalette(name, palettes.Variant(paletteNewVariant))
		p.Description = "Custom color palette"

		// Add minimal default colors
		p.Colors[colorKeyText] = defaultTextColor
		p.Colors["background"] = defaultBgColor
		p.Colors["primary"] = "#7c3aed"
		p.Colors["secondary"] = "#06b6d4"
		p.Colors["accent"] = "#f59e0b"
		p.Colors["success"] = "#10b981"
		p.Colors["warning"] = "#f59e0b"
		p.Colors["error"] = "#ef4444"

		p.Semantic["text-primary"] = colorKeyText
		p.Semantic["bg-primary"] = "background"
		p.Semantic["accent"] = "primary"
		p.Semantic["link"] = "secondary"
	}

	p.Variant = palettes.Variant(paletteNewVariant)

	// Generate TOML
	toml := generatePaletteTOML(p)

	// Determine output file
	outputFile := paletteOutput
	if outputFile == "" {
		outputFile = fmt.Sprintf("palettes/%s.toml", normalizeFileName(name))
	}

	// Ensure directory exists
	if err := os.MkdirAll("palettes", 0o755); err != nil {
		return fmt.Errorf("failed to create palettes directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("file already exists: %s", outputFile)
	}

	if err := os.WriteFile(outputFile, []byte(toml), 0o644); err != nil { //nolint:gosec // palette files should be readable
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Created new palette: %s\n", outputFile)
	return nil
}

// runPaletteCloneCommand clones an existing palette.
func runPaletteCloneCommand(_ *cobra.Command, args []string) error {
	loader := palettes.NewLoader()

	// Get all available palettes for selection
	infos, err := loader.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover palettes: %w", err)
	}

	if len(infos) == 0 {
		return fmt.Errorf("no palettes available to clone")
	}

	var sourcePaletteName string

	if len(args) > 0 {
		// Palette name provided as argument
		sourcePaletteName = args[0]
	} else {
		// Show fuzzy picker
		idx, err := fuzzyfinder.Find(
			infos,
			func(i int) string {
				return infos[i].Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
				if i == -1 {
					return ""
				}
				return formatPalettePreview(infos[i])
			}),
		)
		if err != nil {
			if errors.Is(err, fuzzyfinder.ErrAbort) {
				return fmt.Errorf("selection canceled")
			}
			return fmt.Errorf("fuzzy finder error: %w", err)
		}
		sourcePaletteName = infos[idx].Name
	}

	// Load the source palette
	sourcePalette, err := loader.Load(sourcePaletteName)
	if err != nil {
		return fmt.Errorf("failed to load palette %q: %w", sourcePaletteName, err)
	}

	// Get the new palette name
	var newName string
	if paletteCloneName != "" {
		newName = paletteCloneName
	} else {
		// Prompt for the new name
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter name for the cloned palette [%s-custom]: ", normalizeFileName(sourcePalette.Name))
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		newName = strings.TrimSpace(input)
		if newName == "" {
			newName = normalizeFileName(sourcePalette.Name) + "-custom"
		}
	}

	// Validate the new name doesn't already exist
	normalizedNewName := normalizeFileName(newName)

	// Determine output path (user palettes directory)
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}
	palettesDir := filepath.Join(configDir, "markata-go", "palettes")
	outputFile := filepath.Join(palettesDir, normalizedNewName+".toml")

	// Check if file already exists
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("palette already exists: %s", outputFile)
	}

	// Ensure the palettes directory exists
	if err := os.MkdirAll(palettesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create palettes directory: %w", err)
	}

	// Clone the palette
	clonedPalette := sourcePalette.Clone()
	clonedPalette.Name = newName
	clonedPalette.Author = ""
	clonedPalette.Homepage = ""
	clonedPalette.Description = fmt.Sprintf("Cloned from %s", sourcePaletteName)

	// Generate TOML and write
	tomlContent := generatePaletteTOML(clonedPalette)
	if err := os.WriteFile(outputFile, []byte(tomlContent), 0o644); err != nil { //nolint:gosec // palette files should be readable
		return fmt.Errorf("failed to write palette file: %w", err)
	}

	fmt.Printf("Cloned palette %q to: %s\n", sourcePaletteName, outputFile)
	return nil
}

// formatPalettePreview formats a palette info for the fuzzy finder preview.
func formatPalettePreview(info palettes.PaletteInfo) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", info.Name))
	sb.WriteString(fmt.Sprintf("Variant: %s\n", info.Variant))
	if info.Author != "" {
		sb.WriteString(fmt.Sprintf("Author: %s\n", info.Author))
	}
	if info.Description != "" {
		sb.WriteString(fmt.Sprintf("Description: %s\n", info.Description))
	}
	sb.WriteString(fmt.Sprintf("Source: %s\n", info.Source))

	// Try to load the full palette for color preview
	loader := palettes.NewLoader()
	p, err := loader.Load(info.Name)
	if err == nil {
		sb.WriteString("\nColors:\n")
		count := 0
		for name, hex := range p.Colors {
			if count >= 8 {
				sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(p.Colors)-8))
				break
			}
			sb.WriteString(fmt.Sprintf("  %s: %s\n", name, hex))
			count++
		}
	}

	return sb.String()
}

// generatePreviewHTML generates an HTML preview page for a palette.
func generatePreviewHTML(p *palettes.Palette) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Palette Preview: `)
	sb.WriteString(p.Name)
	sb.WriteString(`</title>
  <style>
    * { box-sizing: border-box; }
    body { font-family: system-ui, sans-serif; margin: 0; padding: 20px; }
    h1 { margin-bottom: 5px; }
    .meta { color: #666; margin-bottom: 20px; }
    .section { margin-bottom: 30px; }
    .section h2 { border-bottom: 1px solid #ddd; padding-bottom: 5px; }
    .colors { display: flex; flex-wrap: wrap; gap: 10px; }
    .color-card {
      width: 150px;
      border: 1px solid #ddd;
      border-radius: 8px;
      overflow: hidden;
    }
    .color-swatch { height: 80px; }
    .color-info { padding: 8px; font-size: 12px; }
    .color-name { font-weight: bold; }
    .color-hex { font-family: monospace; color: #666; }
  </style>
</head>
<body>
  <h1>`)
	sb.WriteString(p.Name)
	sb.WriteString(`</h1>
  <p class="meta">`)
	sb.WriteString(string(p.Variant))
	if p.Author != "" {
		sb.WriteString(` &middot; by `)
		sb.WriteString(p.Author)
	}
	sb.WriteString(`</p>
`)

	// Raw colors section
	sb.WriteString(`  <div class="section">
    <h2>Raw Colors</h2>
    <div class="colors">
`)
	for name, hex := range p.Colors {
		sb.WriteString(fmt.Sprintf(`      <div class="color-card">
        <div class="color-swatch" style="background-color: %s;"></div>
        <div class="color-info">
          <div class="color-name">%s</div>
          <div class="color-hex">%s</div>
        </div>
      </div>
`, hex, name, hex))
	}
	sb.WriteString(`    </div>
  </div>
`)

	// Semantic colors section
	sb.WriteString(`  <div class="section">
    <h2>Semantic Colors</h2>
    <div class="colors">
`)
	for name := range p.Semantic {
		hex := p.Resolve(name)
		sb.WriteString(fmt.Sprintf(`      <div class="color-card">
        <div class="color-swatch" style="background-color: %s;"></div>
        <div class="color-info">
          <div class="color-name">%s</div>
          <div class="color-hex">%s</div>
        </div>
      </div>
`, hex, name, hex))
	}
	sb.WriteString(`    </div>
  </div>
`)

	// Contrast test section
	sb.WriteString(`  <div class="section">
    <h2>Contrast Preview</h2>
`)
	textPrimary := p.Resolve("text-primary")
	bgPrimary := p.Resolve("bg-primary")
	if textPrimary != "" && bgPrimary != "" {
		sb.WriteString(fmt.Sprintf(`    <div style="background-color: %s; color: %s; padding: 20px; border-radius: 8px;">
      <h3 style="margin-top: 0;">Sample Text</h3>
      <p>This is how text-primary looks on bg-primary.</p>
      <p><a href="#" style="color: %s;">This is a link</a></p>
    </div>
`, bgPrimary, textPrimary, p.Resolve("link")))
	}
	sb.WriteString(`  </div>
`)

	sb.WriteString(`</body>
</html>
`)

	return sb.String()
}

// generateAllPalettesPreviewHTML generates an HTML preview page for multiple palettes.
func generateAllPalettesPreviewHTML(paletteList []*palettes.Palette) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>All Palettes Preview</title>
  <style>
    * { box-sizing: border-box; }
    body { font-family: system-ui, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
    h1 { margin-bottom: 10px; }
    .summary { color: #666; margin-bottom: 20px; }
    .nav { position: sticky; top: 0; background: #fff; padding: 15px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
    .nav h2 { margin: 0 0 10px 0; font-size: 14px; text-transform: uppercase; color: #666; }
    .nav-links { display: flex; flex-wrap: wrap; gap: 8px; }
    .nav-links a { color: #0066cc; text-decoration: none; padding: 4px 8px; background: #f0f0f0; border-radius: 4px; font-size: 13px; }
    .nav-links a:hover { background: #e0e0e0; }
    .palette-card { background: #fff; border-radius: 12px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
    .palette-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; }
    .palette-header h2 { margin: 0; }
    .palette-meta { color: #666; font-size: 14px; }
    .variant-badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 12px; font-weight: bold; text-transform: uppercase; }
    .variant-dark { background: #333; color: #fff; }
    .variant-light { background: #f0f0f0; color: #333; }
    .colors { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 15px; }
    .color-swatch {
      width: 60px;
      height: 60px;
      border-radius: 8px;
      border: 1px solid rgba(0,0,0,0.1);
      position: relative;
      cursor: pointer;
    }
    .color-swatch:hover::after {
      content: attr(data-name) "\A" attr(data-hex);
      white-space: pre;
      position: absolute;
      bottom: 100%;
      left: 50%;
      transform: translateX(-50%);
      background: #333;
      color: #fff;
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 11px;
      font-family: monospace;
      z-index: 10;
      margin-bottom: 4px;
    }
    .contrast-preview {
      padding: 15px;
      border-radius: 8px;
      margin-top: 10px;
    }
    .contrast-preview h4 { margin: 0 0 8px 0; }
    .contrast-preview p { margin: 4px 0; }
    .contrast-preview a { text-decoration: underline; }
  </style>
</head>
<body>
  <h1>All Palettes Preview</h1>
`)

	sb.WriteString(fmt.Sprintf(`  <p class="summary">Showing %d palettes</p>
`, len(paletteList)))

	// Navigation
	sb.WriteString(`  <nav class="nav">
    <h2>Jump to Palette</h2>
    <div class="nav-links">
`)
	for _, p := range paletteList {
		sb.WriteString(fmt.Sprintf(`      <a href="#%s">%s</a>
`, normalizeFileName(p.Name), p.Name))
	}
	sb.WriteString(`    </div>
  </nav>

`)

	// Palette cards
	for _, p := range paletteList {
		variantClass := "variant-dark"
		if p.Variant == palettes.VariantLight {
			variantClass = "variant-light"
		}

		sb.WriteString(fmt.Sprintf(`  <div class="palette-card" id="%s">
    <div class="palette-header">
      <h2>%s</h2>
      <span class="variant-badge %s">%s</span>
    </div>
`, normalizeFileName(p.Name), p.Name, variantClass, p.Variant))

		if p.Author != "" || p.Description != "" {
			sb.WriteString(`    <p class="palette-meta">`)
			if p.Author != "" {
				sb.WriteString(fmt.Sprintf(`by %s`, p.Author))
			}
			if p.Author != "" && p.Description != "" {
				sb.WriteString(` &middot; `)
			}
			if p.Description != "" {
				sb.WriteString(p.Description)
			}
			sb.WriteString(`</p>
`)
		}

		// Color swatches
		sb.WriteString(`    <div class="colors">
`)
		for name, hex := range p.Colors {
			sb.WriteString(fmt.Sprintf(`      <div class="color-swatch" style="background-color: %s;" data-name="%s" data-hex="%s"></div>
`, hex, name, hex))
		}
		sb.WriteString(`    </div>
`)

		// Contrast preview
		textPrimary := p.Resolve("text-primary")
		bgPrimary := p.Resolve("bg-primary")
		if textPrimary != "" && bgPrimary != "" {
			linkColor := p.Resolve("link")
			if linkColor == "" {
				linkColor = textPrimary
			}
			sb.WriteString(fmt.Sprintf(`    <div class="contrast-preview" style="background-color: %s; color: %s;">
      <h4>Contrast Preview</h4>
      <p>Sample text on this background.</p>
      <p><a href="#" style="color: %s;">Sample link</a></p>
    </div>
`, bgPrimary, textPrimary, linkColor))
		}

		sb.WriteString(`  </div>

`)
	}

	sb.WriteString(`</body>
</html>
`)

	return sb.String()
}

// generatePaletteTOML generates TOML content for a palette.
func generatePaletteTOML(p *palettes.Palette) string {
	var sb strings.Builder

	sb.WriteString("# ")
	sb.WriteString(p.Name)
	sb.WriteString(" Color Palette\n\n")

	sb.WriteString("[palette]\n")
	sb.WriteString(fmt.Sprintf("name = %q\n", p.Name))
	sb.WriteString(fmt.Sprintf("variant = %q\n", p.Variant))
	if p.Author != "" {
		sb.WriteString(fmt.Sprintf("author = %q\n", p.Author))
	}
	if p.License != "" {
		sb.WriteString(fmt.Sprintf("license = %q\n", p.License))
	}
	if p.Homepage != "" {
		sb.WriteString(fmt.Sprintf("homepage = %q\n", p.Homepage))
	}
	if p.Description != "" {
		sb.WriteString(fmt.Sprintf("description = %q\n", p.Description))
	}

	sb.WriteString("\n# Raw Colors\n")
	sb.WriteString("[palette.colors]\n")
	for name, hex := range p.Colors {
		sb.WriteString(fmt.Sprintf("%s = %q\n", name, hex))
	}

	sb.WriteString("\n# Semantic Colors\n")
	sb.WriteString("[palette.semantic]\n")
	for name, ref := range p.Semantic {
		sb.WriteString(fmt.Sprintf("%s = %q\n", name, ref))
	}

	if len(p.Components) > 0 {
		sb.WriteString("\n# Component Colors\n")
		sb.WriteString("[palette.components]\n")
		for name, ref := range p.Components {
			sb.WriteString(fmt.Sprintf("%s = %q\n", name, ref))
		}
	}

	return sb.String()
}

// normalizeFileName converts a name to a file-friendly format.
func normalizeFileName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}

// openBrowser attempts to open a file in the default browser.
func openBrowser(path string) {
	// Platform-specific browser opening would go here
	// For now, just print a message
	fmt.Printf("Open %s in your browser to view the preview.\n", path)
}

// runPaletteFetchCommand fetches a palette from a Lospec URL.
func runPaletteFetchCommand(_ *cobra.Command, args []string) error {
	rawURL := args[0]

	// Validate and parse the URL
	normalizedURL, err := palettes.ParseLospecURL(rawURL)
	if err != nil {
		return fmt.Errorf("invalid Lospec URL: %w\nExpected format: https://lospec.com/palette-list/<name>.txt", err)
	}

	fmt.Printf("Fetching palette from: %s\n", normalizedURL)

	// Create client and fetch palette
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := palettes.NewLospecClient()
	p, err := client.FetchPalette(ctx, normalizedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch palette: %w", err)
	}

	// Override name if provided
	if paletteFetchName != "" {
		p.Name = paletteFetchName
	}

	// Determine output directory
	var outputDir string
	if paletteFetchOutput != "" {
		outputDir = paletteFetchOutput
	} else {
		outputDir, err = palettes.GetUserPalettesDir()
		if err != nil {
			return fmt.Errorf("failed to get user palettes directory: %w", err)
		}
	}

	// Generate output filename
	outputFile := filepath.Join(outputDir, normalizeFileName(p.Name)+".toml")

	// Check if file already exists
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("palette file already exists: %s\nUse a different name with --name or delete the existing file", outputFile)
	}

	// Save the palette
	if err := palettes.SavePaletteToFile(p, outputFile); err != nil {
		return fmt.Errorf("failed to save palette: %w", err)
	}

	fmt.Printf("Palette saved to: %s\n", outputFile)
	fmt.Printf("\nPalette details:\n")
	fmt.Printf("  Name:        %s\n", p.Name)
	fmt.Printf("  Variant:     %s\n", p.Variant)
	fmt.Printf("  Colors:      %d\n", len(p.Colors))
	fmt.Printf("  Source:      %s\n", p.Homepage)

	// Show semantic mappings
	if len(p.Semantic) > 0 {
		fmt.Printf("\nSemantic mappings:\n")
		for name, ref := range p.Semantic {
			fmt.Printf("  %-20s -> %s\n", name, ref)
		}
	}

	return nil
}
