package plugins

import (
	"strings"
	"testing"
)

func TestNativeMarkdownTabsRenderAccessibleMarkup(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-tabs {#install data-label="Installation methods" style="pills"}

:::: md-tab {label="curl & shell"}
~~~bash
curl -fsSL example.com/install.sh | bash
~~~
::::

:::: md-tab {label="mise"}
Visit **mise**.
::::

:::`)
	assertContainerHTMLInOrder(t, got,
		`<div class="md-tabs" data-md-tabs data-style="pills" id="install">`,
		`role="tablist" aria-label="Installation methods" hidden`,
		`role="tab" id="md-tab-1-0-tab" aria-controls="md-tab-1-0-panel" aria-selected="true" tabindex="0">curl &amp; shell</button>`,
		`role="tab" id="md-tab-1-1-tab" aria-controls="md-tab-1-1-panel" aria-selected="false" tabindex="-1">mise</button>`,
		`<section class="md-tab" role="tabpanel" id="md-tab-1-0-panel" aria-labelledby="md-tab-1-0-tab"><h3 class="md-tab-fallback-label">curl &amp; shell</h3>`,
		"curl -fsSL example.com/install.sh | bash",
		`<section class="md-tab" role="tabpanel" id="md-tab-1-1-panel" aria-labelledby="md-tab-1-1-tab"><h3 class="md-tab-fallback-label">mise</h3>`,
		"<strong>mise</strong>", "</section>", "</div></div>",
	)
	if strings.Contains(got, "<wa-tab-group") || strings.Contains(got, "<script") {
		t.Fatalf("native tabs must not depend on Web Awesome or inline JS: %s", got)
	}
}

func TestNativeMarkdownTabsMultipleGroupsAndFallback(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-tabs

:::: md-tab {label="A"}
First
::::
:::

::: md-tabs {data-style="underline"}

:::: md-tab
Second
::::
:::`)
	for _, want := range []string{
		`id="md-tab-1-0-tab"`, `id="md-tab-2-0-tab"`,
		`data-style="underline"`, `<h3 class="md-tab-fallback-label">Tab 1</h3>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q:\n%s", want, got)
		}
	}
}

func TestNativeMarkdownTabsInvalidGroupsRemainReadable(t *testing.T) {
	for _, src := range []string{
		"::: md-tabs\n\nordinary body\n\n:::",
		"::: md-tabs\n\n:::: md-tab\nOne\n::::\n\nother body\n\n:::",
	} {
		got := renderContainerMarkdown(t, src)
		if strings.Contains(got, "data-md-tabs") || !strings.Contains(got, "body") {
			t.Fatalf("invalid group should remain normal Markdown: %s", got)
		}
	}
}

func TestNativeMarkdownTabsEscapeLabels(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-tabs {data-label="Install & test"}

:::: md-tab {label="<script>alert(1)</script>"}
Hello.
::::
:::`)
	if strings.Contains(got, "<script>") || !strings.Contains(got, "&lt;script&gt;") || !strings.Contains(got, "Install &amp; test") {
		t.Fatalf("unsafe tab label: %s", got)
	}
}
