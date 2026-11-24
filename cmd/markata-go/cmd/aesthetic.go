package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// TODO: Import the aesthetic package when it's implemented.
// import "github.com/WaylonWalker/markata/pkg/aesthetics"

// tokenNone is the constant value for "none" tokens.
const tokenNone = "none"

// Aesthetic represents a design aesthetic with its tokens and description.
// TODO: Move to pkg/aesthetics package when implemented.
type Aesthetic struct {
	Name        string
	Description string
	Tokens      AestheticTokens
}

// AestheticTokens contains the design tokens for an aesthetic.
type AestheticTokens struct {
	RadiusSm   string
	RadiusMd   string
	RadiusLg   string
	RadiusFull string
	Spacing    string
	Border     string
	Shadow     string
}

// builtinAesthetics contains the predefined aesthetics.
// TODO: Load from pkg/aesthetics package when implemented.
var builtinAesthetics = []Aesthetic{
	{
		Name:        "brutal",
		Description: "Brutalist design: harsh, uncompromising, raw",
		Tokens: AestheticTokens{
			RadiusSm:   "0",
			RadiusMd:   "0",
			RadiusLg:   "0",
			RadiusFull: "0",
			Spacing:    "0.75x scale",
			Border:     "3px solid",
			Shadow:     tokenNone,
		},
	},
	{
		Name:        "precision",
		Description: "Technical/engineering: clean, exact, minimal",
		Tokens: AestheticTokens{
			RadiusSm:   "2px",
			RadiusMd:   "2px",
			RadiusLg:   "4px",
			RadiusFull: "4px",
			Spacing:    "1x scale",
			Border:     "1px solid",
			Shadow:     tokenNone,
		},
	},
	{
		Name:        "balanced",
		Description: "Default harmonious: comfortable, balanced",
		Tokens: AestheticTokens{
			RadiusSm:   "0.25rem",
			RadiusMd:   "0.375rem",
			RadiusLg:   "0.5rem",
			RadiusFull: "9999px",
			Spacing:    "1x scale",
			Border:     "1px solid",
			Shadow:     "0 1px 3px rgba(0,0,0,0.1)",
		},
	},
	{
		Name:        "elevated",
		Description: "Layered/premium: depth, floating cards",
		Tokens: AestheticTokens{
			RadiusSm:   "0.5rem",
			RadiusMd:   "0.75rem",
			RadiusLg:   "1rem",
			RadiusFull: "9999px",
			Spacing:    "1.25x scale",
			Border:     tokenNone,
			Shadow:     "0 4px 12px rgba(0,0,0,0.15)",
		},
	},
	{
		Name:        "minimal",
		Description: "Maximum whitespace: sparse, intentional",
		Tokens: AestheticTokens{
			RadiusSm:   "0",
			RadiusMd:   "0",
			RadiusLg:   "0",
			RadiusFull: "0",
			Spacing:    "1.5x scale",
			Border:     tokenNone,
			Shadow:     tokenNone,
		},
	},
}

// aestheticCmd represents the aesthetic command group.
var aestheticCmd = &cobra.Command{
	Use:   "aesthetic",
	Short: "Design aesthetic commands",
	Long: `Commands for managing design aesthetics.

Aesthetics define the visual style tokens for your site, including
border radius, spacing scale, border styles, and shadow effects.

Subcommands:
  list  - List available aesthetics
  show  - Show details of a specific aesthetic`,
}

// aestheticListCmd lists available aesthetics.
var aestheticListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available aesthetics",
	Long: `List all available design aesthetics with descriptions.

Aesthetics control the visual "feel" of your site through design tokens
like border radius, spacing, and shadow intensity.

Example usage:
  markata-go aesthetic list`,
	RunE: runAestheticListCommand,
}

// aestheticShowCmd shows details of a specific aesthetic.
var aestheticShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show aesthetic details",
	Long: `Display detailed information about a specific aesthetic.

Shows all design tokens and a CSS preview of the generated custom properties.

Example usage:
  markata-go aesthetic show brutal
  markata-go aesthetic show balanced`,
	Args: cobra.ExactArgs(1),
	RunE: runAestheticShowCommand,
}

func init() {
	rootCmd.AddCommand(aestheticCmd)

	// List subcommand
	aestheticCmd.AddCommand(aestheticListCmd)

	// Show subcommand
	aestheticCmd.AddCommand(aestheticShowCmd)
}

// runAestheticListCommand lists available aesthetics.
func runAestheticListCommand(_ *cobra.Command, _ []string) error {
	// TODO: Load aesthetics from pkg/aesthetics package when implemented.
	// loader := aesthetics.NewLoader()
	// aestheticList, err := loader.Discover()
	// if err != nil {
	//     return fmt.Errorf("failed to discover aesthetics: %w", err)
	// }

	aestheticList := builtinAesthetics

	if len(aestheticList) == 0 {
		fmt.Println("No aesthetics found.")
		return nil
	}

	fmt.Println("Available aesthetics:")
	for i := range aestheticList {
		fmt.Printf("  %-10s - %s\n", aestheticList[i].Name, aestheticList[i].Description)
	}

	return nil
}

// runAestheticShowCommand shows details of a specific aesthetic.
func runAestheticShowCommand(_ *cobra.Command, args []string) error {
	name := args[0]

	// TODO: Load aesthetic from pkg/aesthetics package when implemented.
	// loader := aesthetics.NewLoader()
	// aesthetic, err := loader.Load(name)
	// if err != nil {
	//     return fmt.Errorf("failed to load aesthetic: %w", err)
	// }

	var aesthetic *Aesthetic
	for i := range builtinAesthetics {
		if strings.EqualFold(builtinAesthetics[i].Name, name) {
			aesthetic = &builtinAesthetics[i]
			break
		}
	}

	if aesthetic == nil {
		return fmt.Errorf("aesthetic not found: %s", name)
	}

	fmt.Printf("Aesthetic: %s\n", aesthetic.Name)
	fmt.Printf("Description: %s\n", aesthetic.Description)
	fmt.Println()

	fmt.Println("Tokens:")
	fmt.Printf("  radius:  %s (sm), %s (md), %s (lg)\n",
		aesthetic.Tokens.RadiusSm, aesthetic.Tokens.RadiusMd, aesthetic.Tokens.RadiusLg)
	fmt.Printf("  spacing: %s\n", aesthetic.Tokens.Spacing)
	fmt.Printf("  border:  %s\n", aesthetic.Tokens.Border)
	fmt.Printf("  shadow:  %s\n", aesthetic.Tokens.Shadow)
	fmt.Println()

	fmt.Println("CSS Preview:")
	fmt.Printf("  --radius-sm: %s;\n", aesthetic.Tokens.RadiusSm)
	fmt.Printf("  --radius-md: %s;\n", aesthetic.Tokens.RadiusMd)
	fmt.Printf("  --radius-lg: %s;\n", aesthetic.Tokens.RadiusLg)
	fmt.Printf("  --radius-full: %s;\n", aesthetic.Tokens.RadiusFull)
	if aesthetic.Tokens.Shadow != tokenNone {
		fmt.Printf("  --shadow: %s;\n", aesthetic.Tokens.Shadow)
	} else {
		fmt.Println("  --shadow: none;")
	}
	if aesthetic.Tokens.Border != tokenNone {
		fmt.Printf("  --border: %s;\n", aesthetic.Tokens.Border)
	} else {
		fmt.Println("  --border: none;")
	}

	return nil
}
