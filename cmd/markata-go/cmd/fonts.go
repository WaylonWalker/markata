package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/fontpacks"
	"github.com/spf13/cobra"
)

var fontsCmd = &cobra.Command{Use: "fonts", Short: "Inspect and maintain Markata font packs"}
var fontsListCmd = &cobra.Command{Use: "list", Short: "List available font sources", RunE: runFontsList}
var fontsPacksCmd = &cobra.Command{Use: "packs", Short: "List available font packs", RunE: runFontsPacks}
var fontsShowCmd = &cobra.Command{Use: "show <font-or-pack>", Args: cobra.ExactArgs(1), RunE: runFontsShow}
var fontsDoctorCmd = &cobra.Command{Use: "doctor", Short: "Check optional FontTools support", RunE: runFontsDoctor}
var fontsVerifyCmd = &cobra.Command{Use: "verify [pack]", Short: "Validate the local font catalog", Args: cobra.MaximumNArgs(1), RunE: runFontsVerify}
var fontsReportCmd = &cobra.Command{Use: "report", Short: "Report selected pack size", RunE: runFontsReport}
var fontsLicensesCmd = &cobra.Command{Use: "licenses", Short: "List catalog license records", RunE: runFontsLicenses}

const defaultFontpack = "system"

func init() {
	rootCmd.AddCommand(fontsCmd)
	fontsCmd.AddCommand(fontsListCmd, fontsPacksCmd, fontsShowCmd, fontsDoctorCmd, fontsVerifyCmd, fontsReportCmd, fontsLicensesCmd)
	fontsCmd.AddCommand(unsupportedFontsCommand("vendor"), unsupportedFontsCommand("rebuild"), unsupportedFontsCommand("add"), unsupportedFontsCommand("remove"))
}
func loadFontSource() (*fontpacks.CatalogSource, error) { return fontpacks.BuiltinSource() }
func loadConfiguredFontSource() (*fontpacks.CatalogSource, error) {
	path, err := config.Discover()
	if err != nil {
		return loadFontSource()
	}
	cfg, err := config.Load(path)
	if err != nil || cfg.FontpacksFile == "" {
		return loadFontSource()
	}
	catalogPath := cfg.FontpacksFile
	if !filepath.IsAbs(catalogPath) {
		catalogPath = filepath.Join(filepath.Dir(path), catalogPath)
	}
	return fontpacks.LoadSource(catalogPath)
}
func loadFontCatalog() (*fontpacks.Catalog, error) {
	s, err := loadConfiguredFontSource()
	if err != nil {
		return nil, err
	}
	return s.Catalog, nil
}

func configuredFontpack() string {
	path, err := config.Discover()
	if err != nil {
		return defaultFontpack
	}
	cfg, err := config.Load(path)
	if err != nil || cfg.Fontpack == "" {
		return defaultFontpack
	}
	return cfg.Fontpack
}
func runFontsList(*cobra.Command, []string) error {
	c, err := loadFontCatalog()
	if err != nil {
		return err
	}
	for _, id := range fontpacks.SortedKeys(c.FontSources) {
		outlnf("%s\t%s", id, c.FontSources[id].Family)
	}
	return nil
}
func runFontsPacks(*cobra.Command, []string) error {
	c, err := loadFontCatalog()
	if err != nil {
		return err
	}
	for _, id := range fontpacks.SortedKeys(c.FontPacks) {
		p := c.FontPacks[id]
		outlnf("%s\t%s\t%s", id, p.Performance.Class, p.Description)
	}
	return nil
}
func runFontsShow(_ *cobra.Command, args []string) error {
	c, err := loadFontCatalog()
	if err != nil {
		return err
	}
	name := args[0]
	if p, ok := c.FontPacks[name]; ok {
		outlnf("Pack: %s\nPerformance: %s\n%s", name, p.Performance.Class, p.Description)
		return nil
	}
	if s, ok := c.FontSources[name]; ok {
		outlnf("Source: %s\nProvider: %s\nFamily: %s", name, s.Provider, s.Family)
		return nil
	}
	_, _, err = c.ResolvePack(name)
	return err
}
func runFontsDoctor(*cobra.Command, []string) error {
	py, pyErr := exec.LookPath("python")
	subset, subsetErr := exec.LookPath("pyftsubset")
	outln("FontTools doctor")
	if pyErr != nil {
		outln("  python: missing")
	} else {
		outln("  python: " + py)
	}
	if subsetErr != nil {
		outln("  pyftsubset: missing")
	} else {
		outln("  pyftsubset: " + subset)
	}
	if pyErr != nil || subsetErr != nil {
		outln("Custom font processing requires FontTools with WOFF2/Brotli support.\nRecommended:\n  uv tool install \"fonttools[woff]\"\nAlternative:\n  python -m pip install \"fonttools[woff]\"")
	}
	return nil
}
func runFontsVerify(_ *cobra.Command, args []string) error {
	selected := ""
	if len(args) == 1 {
		selected = args[0]
	}
	var verifySource *fontpacks.CatalogSource
	var err error
	if selected == "" {
		// Bare verification is an integrity check of the shipped catalog, not
		// a check of the zero-download default pack.
		verifySource, err = fontpacks.BuiltinSource()
		if err != nil {
			return err
		}
	} else {
		verifySource, err = loadConfiguredFontSource()
		if err != nil {
			return err
		}
	}
	if err := fontpacks.VerifySource(verifySource, selected); err != nil {
		return err
	}
	if selected == "" {
		packs, sources := fontpacks.BundledScope(verifySource.Catalog)
		outlnf("Verified %d bundled font sources across %d built-in packs.", sources, packs)
	} else {
		resolved, pack, resolveErr := verifySource.Catalog.ResolvePack(selected)
		if resolveErr != nil {
			return resolveErr
		}
		_, sources := fontpacks.BundledScopeForPack(pack)
		outlnf("Verified fontpack %q: %d font sources.", resolved, sources)
	}
	return nil
}
func runFontsReport(*cobra.Command, []string) error {
	s, err := loadConfiguredFontSource()
	if err != nil {
		return err
	}
	c := s.Catalog
	name := configuredFontpack()
	resolvedName, p, err := c.ResolvePack(name)
	if err != nil {
		return err
	}
	families := map[string]bool{}
	for _, role := range p.Roles {
		if role.Source != "" {
			families[role.Source] = true
		}
	}
	files, bytes := 0, int64(0)
	tierFiles := map[string]int{}
	tierBytes := map[string]int64{}
	outputDir := defaultOutputDir
	if path, discoverErr := config.Discover(); discoverErr == nil {
		if cfg, loadErr := config.Load(path); loadErr == nil {
			outputDir = cfg.OutputDir
			if !filepath.IsAbs(outputDir) {
				outputDir = filepath.Join(filepath.Dir(path), outputDir)
			}
		}
	}
	if names, readErr := fontpacks.ManagedFontFiles(outputDir); readErr == nil {
		files = len(names)
		for _, name := range names {
			if info, infoErr := os.Stat(filepath.Join(outputDir, "assets", "fonts", name)); infoErr == nil {
				bytes += info.Size()
				for _, tier := range []string{"display-core", "prose-core", "code-core", "latin-ext", "full"} {
					if strings.HasSuffix(name, "-"+tier+".woff2") {
						tierFiles[tier]++
						tierBytes[tier] += info.Size()
						break
					}
				}
			}
		}
	}
	outlnf("Pack: %s\nPerformance class: %s\nFamilies: %d\nFiles emitted: %d\n%s", resolvedName, p.Performance.Class, len(families), files, fontReportAssetBytesLine(bytes))
	for _, tier := range []string{"display-core", "prose-core", "code-core", "latin-ext", "full"} {
		if tierFiles[tier] > 0 {
			outlnf("  %s: %d files, %d bytes", tier, tierFiles[tier], tierBytes[tier])
		}
	}
	return nil
}

func fontReportAssetBytesLine(bytes int64) string {
	return fmt.Sprintf("Emitted font asset bytes: %d", bytes)
}

func runFontsLicenses(*cobra.Command, []string) error {
	c, err := loadFontCatalog()
	if err != nil {
		return err
	}
	outln("Licenses are recorded in family manifests and THIRD_PARTY_FONTS.md.")
	outln(strings.Join(fontpacks.SortedKeys(c.FontSources), ", "))
	return nil
}

func unsupportedFontsCommand(name string) *cobra.Command {
	return &cobra.Command{Use: name, RunE: func(*cobra.Command, []string) error {
		return fmt.Errorf("%s is reserved for the catalog maintenance workflow; use markata fonts doctor before installing FontTools", name)
	}}
}
