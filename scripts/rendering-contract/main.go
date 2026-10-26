// Command generate-rendering-contract projects the language-neutral contract to
// the Go and browser consumers. Run it from the markata-go repository root.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"
)

func main() {
	check := flag.Bool("check", false, "check generated files without writing")
	resolveFixtures := flag.Bool("resolve-fixtures", false, "write resolved shared render-plan fixtures as JSON")
	derivePalettes := flag.Bool("derive-palettes", false, "recompute roles of derived counterpart palettes in the source contract from pkg/palettes")
	flag.Parse()
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	source := filepath.Join(root, "spec", "rendering-contract", "contract-v1.json")
	raw, err := os.ReadFile(source)
	if err != nil {
		panic(err)
	}
	if *derivePalettes {
		raw = syncDerivedPalettes(raw)
		//nolint:gosec // G306: this is a checked-in generated source artifact.
		if err := os.WriteFile(source, raw, 0o644); err != nil {
			panic(err)
		}
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		panic(err)
	}
	if *resolveFixtures {
		fixturesPath := filepath.Join(root, "spec", "rendering-contract", "render-plan-fixtures.json")
		fixtureRaw, err := os.ReadFile(fixturesPath)
		if err != nil {
			panic(err)
		}
		var fixtures renderingcontract.RenderPlanFixtures
		if err := json.Unmarshal(fixtureRaw, &fixtures); err != nil {
			panic(err)
		}
		results := make([]map[string]any, 0, len(fixtures.Fixtures))
		for i := range fixtures.Fixtures {
			fixture := &fixtures.Fixtures[i]
			state, plan, _, err := renderingcontract.ResolveFixture(*fixture)
			if err != nil {
				panic(err)
			}
			results = append(results, map[string]any{"id": fixture.ID, "normalized": state, "plan": plan})
		}
		encoded, err := json.Marshal(results)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(encoded))
		return
	}
	pretty, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	pretty = append(pretty, '\n')
	jsonTargets := []string{filepath.Join(root, "pkg", "renderingcontract", "contract-v1.json")}
	jsTargets := map[string][]byte{
		filepath.Join(root, "..", "md.waylonwalker.com", "src", "rendering-contract.generated.js"):     append([]byte("// Generated. Do not edit.\nexport const RENDERING_CONTRACT = Object.freeze("), append(pretty, []byte(");\nexport const CONTRACT_PALETTES = RENDERING_CONTRACT.palettes;\nexport const CONTRACT_ENUMS = RENDERING_CONTRACT.enums;\n")...)...),
		filepath.Join(root, "..", "themes.waylonwalker.com", "src", "rendering-contract.generated.js"): append([]byte("// Generated. Do not edit.\nwindow.renderingContract = "), append(pretty, []byte(";\n")...)...),
	}
	for _, target := range jsonTargets {
		if err := project(target, pretty, *check); err != nil {
			panic(err)
		}
	}
	for target, data := range jsTargets {
		if *check {
			if _, err := os.Stat(target); os.IsNotExist(err) {
				// The browser projections live in sibling repositories during the
				// coordinated workspace, but are not present in a standalone Go
				// checkout. CI still verifies the Go projection here; the full
				// workspace check verifies all three when siblings are available.
				fmt.Printf("skipping absent sibling projection %s\n", target)
				continue
			}
		}
		if err := project(target, data, *check); err != nil {
			panic(err)
		}
	}
	if *check {
		fmt.Println("rendering contract generated artifacts are current")
	} else {
		fmt.Println("rendering contract generated artifacts updated")
	}
}

func project(path string, expected []byte, check bool) error {
	actual, err := os.ReadFile(path)
	if check {
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if !bytes.Equal(actual, expected) {
			return fmt.Errorf("%s is stale", path)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	//nolint:gosec // Generated static assets must be readable by the web server.
	return os.WriteFile(path, expected, 0o644)
}

// syncDerivedPalettes rewrites the roles of every contract palette that has no
// hand-authored source (an "<id>-light"/"<id>-dark" counterpart of a shipped
// palette) using palettes.DeriveCounterpart, so the contract and the palette
// loader agree on derived colors. Edits are applied textually so the source
// document keeps its key order and formatting.
func syncDerivedPalettes(raw []byte) []byte {
	var doc struct {
		Palettes []struct {
			ID string `json:"id"`
		} `json:"palettes"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		panic(err)
	}
	loader := palettes.NewLoaderWithPaths(nil)
	text := string(raw)
	updated := 0
	for _, entry := range doc.Palettes {
		if palettes.HasBuiltin(entry.ID) {
			continue
		}
		base, variant, ok := palettes.IsDerivedCounterpart(entry.ID)
		if !ok || !palettes.HasBuiltin(base) {
			continue
		}
		p, err := loader.Load(entry.ID)
		if err != nil || p.Variant != variant {
			continue
		}
		idMatch := regexp.MustCompile(`"id":\s*` + regexp.QuoteMeta(fmt.Sprintf("%q", entry.ID))).FindStringIndex(text)
		if idMatch == nil {
			continue
		}
		idPos := idMatch[0]
		rolesMatch := regexp.MustCompile(`"roles":\s*\{`).FindStringIndex(text[idPos:])
		if rolesMatch == nil {
			continue
		}
		rolesStart := idPos + rolesMatch[0]
		rolesEnd := rolesStart + strings.Index(text[rolesStart:], "}")
		block := text[rolesStart:rolesEnd]
		for role, semantic := range map[string]string{"accent": "accent", "background": "bg-primary", "ink": "text-primary", "surface": "bg-surface"} {
			re := regexp.MustCompile(fmt.Sprintf(`(%q:\s*)"#[0-9a-fA-F]{6}"`, role))
			block = re.ReplaceAllString(block, `${1}`+fmt.Sprintf("%q", p.Resolve(semantic)))
		}
		text = text[:rolesStart] + block + text[rolesEnd:]
		updated++
	}
	fmt.Printf("synced %d derived contract palettes\n", updated)
	return []byte(text)
}
