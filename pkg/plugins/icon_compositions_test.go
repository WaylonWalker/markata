package plugins

import (
	"strings"
	"testing"
)

func TestIconCompositionsUseExistingShortcodesAndNestedContainers(t *testing.T) {
	p := NewIconsPlugin()
	svg, ok := prepareIconSVG(testIconSVG)
	if !ok {
		t.Fatal("test SVG must pass normal icon preparation")
	}
	for _, name := range []string{"rocket", "server", "database"} {
		icon := iconAsset{name: "lucide/" + name, svg: svg}
		p.icons["lucide-"+name] = icon
		p.icons["lucide/"+name] = icon
	}

	source := `::: md-feature

:::: md-feature-icon
:lucide-rocket:

::::

:::: md-feature-body
### Verify production
Confirm the live site.

::::

:::

::: md-columns

:::: md-card {data-align="center"}

::::: md-icon {data-size="lg"}
:lucide-server:

:::::

### Falcon2
K3s node

::::

:::: md-card {data-align="center"}

::::: md-icon {data-size="lg"}
:lucide-database:

:::::

### Falcon3
NAS + ZFS

::::

:::`
	expanded := p.processContent(source)
	if got := strings.Count(expanded, `data-icon="lucide/`); got != 3 {
		t.Fatalf("expected three existing icon shortcodes to be expanded, got %d:\n%s", got, expanded)
	}
	if !strings.Contains(expanded, `aria-hidden="true"`) || !strings.Contains(expanded, `width="1em" height="1em"`) {
		t.Fatalf("icons should remain decorative and em-sized: %s", expanded)
	}
	// Goldmark's container parser must preserve the new class/attribute
	// contracts for regular Markdown and keep all the nested closers matched.
	rendered := renderContainerMarkdown(t, source)
	assertContainerHTMLInOrder(t, rendered,
		`<div class="md-feature">`,
		`<div class="md-feature-icon">`,
		`<p>:lucide-rocket:</p>`,
		`</div>`,
		`<div class="md-feature-body">`,
		`<h3>Verify production</h3>`,
		`</div>`,
		`</div>`,
		`<div class="md-columns">`,
		`<div class="md-card" data-align="center">`,
		`<div class="md-icon" data-size="lg">`,
		`<p>:lucide-server:</p>`,
		`</div>`,
	)
}

func TestIconsInFencedRichMarkdownExamplesStayLiteral(t *testing.T) {
	p := NewIconsPlugin()
	icon := iconAsset{name:"lucide/rocket",svg:testIconSVG}
	p.icons["lucide-rocket"] = icon
	input := "Actual :lucide-rocket:\n\n````markdown\n::: md-feature\n:lucide-rocket:\n:::\n````\n"
	output := p.processContent(input)
	if got := strings.Count(output, `data-icon="lucide/rocket"`); got != 1 {
		t.Fatalf("icon examples in code fences must remain literal, got %s",output)
	}
	if !strings.Contains(output, "::: md-feature\n:lucide-rocket:\n:::") {
		t.Fatalf("fenced example was modified: %s",output)
	}
}
