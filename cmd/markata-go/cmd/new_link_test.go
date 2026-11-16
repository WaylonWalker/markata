package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func linkTestCommand(t *testing.T) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	previousCommand, previousNoInput := currentCmd, noInput
	previousSiteSelected, previousSiteDir := siteDirSelected, activeSiteDir
	resetNewCommandFlags()
	siteDirSelected = false
	t.Cleanup(func() {
		resetNewCommandFlags()
		currentCmd, noInput = previousCommand, previousNoInput
		siteDirSelected, activeSiteDir = previousSiteSelected, previousSiteDir
	})
	t.Chdir(t.TempDir())
	cmd := &cobra.Command{Use: "markata new"}
	cmd.Flags().StringVarP(&newTemplate, "template", "t", "post", "")
	cmd.Flags().StringVar(&newDir, "dir", "", "")
	cmd.Flags().StringVar(&newTags, "tags", "", "")
	cmd.Flags().BoolVar(&newDraft, "draft", false, "")
	cmd.Flags().String("title", "", "")
	cmd.SetOut(&bytes.Buffer{})
	stderr := &bytes.Buffer{}
	cmd.SetErr(stderr)
	cmd.SetIn(strings.NewReader(""))
	return cmd, stderr
}

func readLinkFrontmatter(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		t.Fatalf("missing frontmatter: %s", content)
	}
	var fm map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		t.Fatal(err)
	}
	return fm
}

func TestRunNewCommand_LinkURLMetadata(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	noInput = true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/original" {
			http.Redirect(w, r, "/articles/final", http.StatusFound)
			return
		}
		fmt.Fprint(w, `<title>Wrong HTML title</title>
<meta content="A link &amp; a title" property="og:title">
<meta property=og:description content="A useful description">
<meta property="og:image" content="../preview.png">
<meta property="og:url" content="https://canonical.example/not-the-entered-url">`)
	}))
	t.Cleanup(server.Close)
	if err := cmd.ParseFlags([]string{"--template", "bookmark", "--tags", "go,links", "--draft"}); err != nil {
		t.Fatal(err)
	}
	destination := server.URL + "/original"
	if err := runNewCommand(cmd, []string{destination}); err != nil {
		t.Fatal(err)
	}
	fm := readLinkFrontmatter(t, filepath.Join("pages", "link", "a-link-a-title.md"))
	for key, want := range map[string]interface{}{
		"title": "A link & a title", "description": "A useful description",
		"image": server.URL + "/preview.png", "url": destination,
		"template": "bookmark", "draft": true,
	} {
		if fm[key] != want {
			t.Errorf("%s = %v, want %v", key, fm[key], want)
		}
	}
	if tags, ok := fm["tags"].([]interface{}); !ok || len(tags) != 2 {
		t.Errorf("tags = %#v", fm["tags"])
	}
}

func TestRunNewCommand_LinkLowConfidence(t *testing.T) {
	for _, tt := range []struct {
		name       string
		page       string
		status     int
		path       string
		suggestion string
	}{
		{"html title", "<title>HTML suggestion</title>", 200, "/article", "HTML suggestion"},
		{"twitter title", `<meta name="twitter:title" content="Twitter suggestion">`, 200, "/article", "Twitter suggestion"},
		{"heading", "<h1>A <em>useful</em> heading</h1>", 200, "/article", "A useful heading"},
		{"slug", "", 200, "/my-useful-article.html?secret=no", "my useful article"},
		{"fetch failure", "", 503, "/another-article", "another article"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd, stderr := linkTestCommand(t)
			noInput = true
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.page)
			}))
			t.Cleanup(server.Close)
			if err := cmd.ParseFlags([]string{"--template", "link", "--dir", "bookmarks", "--draft", "--tags", "keep"}); err != nil {
				t.Fatal(err)
			}
			destination := server.URL + tt.path
			err := runNewCommand(cmd, []string{destination})
			if err == nil {
				t.Fatal("expected a low-confidence title error, not a file")
			}
			for _, want := range []string{"--title", tt.suggestion, destination, "--dir", "bookmarks", "--draft", "--tags", "keep"} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("recovery error missing %q: %v", want, err)
				}
			}
			if _, err := os.Stat("bookmarks"); !os.IsNotExist(err) {
				t.Errorf("low-confidence creation wrote output: %v", err)
			}
			if tt.status != 200 && !strings.Contains(stderr.String(), "503") {
				t.Errorf("fetch failure was not explained: %s", stderr)
			}
		})
	}
}

func TestRunNewCommand_LinkTitleOverride(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	noInput = true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `<meta property="og:title" content="Fetched title"><meta property="og:image" content="/image.png">`)
	}))
	t.Cleanup(server.Close)
	if err := cmd.ParseFlags([]string{"--template", "link", "--title", "User title"}); err != nil {
		t.Fatal(err)
	}
	if err := runNewCommand(cmd, []string{server.URL}); err != nil {
		t.Fatal(err)
	}
	fm := readLinkFrontmatter(t, filepath.Join("pages", "link", "user-title.md"))
	if fm["title"] != "User title" || fm["image"] != server.URL+"/image.png" {
		t.Fatalf("explicit title or metadata lost: %#v", fm)
	}
}

func TestRunInteractiveMode_LinkURLBeforeTitle(t *testing.T) {
	cmd, stderr := linkTestCommand(t)
	noInput = false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `<meta property="og:title" content="Fetched title">`)
	}))
	t.Cleanup(server.Close)
	if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
		t.Fatal(err)
	}
	cmd.SetIn(strings.NewReader(server.URL + "/article\n\n"))
	currentCmd = cmd
	input, err := runInteractiveMode(cmd, builtinTemplates())
	if err != nil {
		t.Fatal(err)
	}
	if input.title != "Fetched title" || input.template.Frontmatter["url"] != server.URL+"/article" {
		t.Fatalf("URL-first defaults not applied: %#v", input)
	}
	if strings.Contains(stderr.String(), "Title:") {
		t.Fatalf("trusted metadata unnecessarily asked for a title: %s", stderr)
	}
}

func TestPromptNewLinkTitle_Confirmation(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		want  string
	}{
		{"accept suggestion", "\n\n", "a useful link"},
		{"replace suggestion", "My chosen title\n\n", "My chosen title"},
		{"EOF is not confirmation", "", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd, _ := linkTestCommand(t)
			currentCmd = cmd
			title, err := promptNewLinkTitle(bufio.NewReader(strings.NewReader(tt.input)), "a useful link")
			if tt.want == "" {
				if err == nil {
					t.Fatal("EOF silently accepted an inferred title")
				}
				return
			}
			if err != nil || title != tt.want {
				t.Fatalf("title = %q, error = %v", title, err)
			}
		})
	}
}

func TestRunNewCommand_LinkLowConfidenceNoTerminal(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	noInput, newPlain = false, true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "<title>A suggested title</title>")
	}))
	t.Cleanup(server.Close)
	if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
		t.Fatal(err)
	}
	cmd.SetIn(strings.NewReader("\n"))
	err := runNewCommand(cmd, []string{server.URL})
	if err == nil || !strings.Contains(err.Error(), "--title=") {
		t.Fatalf("non-terminal input silently accepted an inferred title: %v", err)
	}
	if _, err := os.Stat("pages"); !os.IsNotExist(err) {
		t.Fatalf("unattended creation wrote output: %v", err)
	}
}

func TestRunNewCommand_LinkOverrideAfterFetchFailure(t *testing.T) {
	cmd, stderr := linkTestCommand(t)
	noInput = true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(server.Close)
	if err := cmd.ParseFlags([]string{"--template", "link", "--title", "Explicit title"}); err != nil {
		t.Fatal(err)
	}
	if err := runNewCommand(cmd, []string{server.URL}); err != nil {
		t.Fatal(err)
	}
	fm := readLinkFrontmatter(t, filepath.Join("pages", "link", "explicit-title.md"))
	if fm["title"] != "Explicit title" || !strings.Contains(stderr.String(), "403") {
		t.Fatalf("override/fetch warning lost: %#v; %s", fm, stderr)
	}
}

func TestRunNewCommand_LinkInvalidURL(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
		t.Fatal(err)
	}
	err := runNewCommand(cmd, []string{"https://"})
	if err == nil || !strings.Contains(err.Error(), "absolute HTTP(S) URL") {
		t.Fatalf("invalid URL did not fail clearly: %v", err)
	}
}

func TestRunNewCommand_LinkPreservesTitleArguments(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	noInput = true
	if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
		t.Fatal(err)
	}
	if err := runNewCommand(cmd, []string{"My existing title"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join("pages", "link", "my-existing-title.md")); err != nil {
		t.Fatal(err)
	}
}

func TestNewLinkMetadata_UntrustedHostAndURLTitles(t *testing.T) {
	for _, title := range []string{"example.com", "www.example.com/", "https://example.com/article", "   "} {
		if usableNewLinkTitle(title, "https://example.com/article") {
			t.Errorf("trusted non-title %q", title)
		}
	}
}

func TestResolveNewLinkMetadata_Cancellation(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cmd.SetContext(ctx)
	_, err := resolveNewLinkMetadata(cmd, "https://example.com/article")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation became a metadata fallback: %v", err)
	}
}

func TestApplyNewWizardState_LinkDefaultsAndOverrides(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	if err := cmd.ParseFlags([]string{"--template", "bookmark", "--tags=", "--draft=false", "--dir", "chosen"}); err != nil {
		t.Fatal(err)
	}
	currentCmd = cmd
	templates := builtinTemplates()
	templates["link"].Frontmatter["private"] = true
	templates["link"].Frontmatter["draft"] = true
	templates["link"].Frontmatter["tags"] = []string{"starter"}
	state := &newWizardState{
		Template: tmplLink, LinkURL: "https://example.com/article",
		PrivateChoice: "inherit",
	}
	initializeNewLinkWizard(cmd, state, templates["link"], newLinkMetadata{
		Title: "Fetched title", Description: "Fetched description", TrustedTitle: true,
	})
	if err := applyNewWizardState(state, templates); err != nil {
		t.Fatal(err)
	}
	fm := readLinkFrontmatter(t, filepath.Join("chosen", "fetched-title.md"))
	if fm["private"] != true || fm["draft"] != false || fm["template"] != "bookmark" {
		t.Fatalf("starter defaults or explicit flags lost: %#v", fm)
	}
	if tags, ok := fm["tags"].([]interface{}); !ok || len(tags) != 0 {
		t.Fatalf("explicit empty tags lost: %#v", fm["tags"])
	}
	if templates["link"].Frontmatter["url"] != "" {
		t.Fatal("wizard mutated the reusable starter")
	}
}

func TestRunInteractiveMode_LinkPreservesStarterDraft(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	noInput = false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `<meta property="og:title" content="Fetched title">`)
	}))
	t.Cleanup(server.Close)
	if err := os.Mkdir("content-templates", 0o755); err != nil {
		t.Fatal(err)
	}
	starter := "---\n_directory: bookmarks\ndraft: true\nprivate: true\ntags: [keep]\n---\n\nMy commentary"
	if err := os.WriteFile(filepath.Join("content-templates", "link.md"), []byte(starter), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
		t.Fatal(err)
	}
	cmd.SetIn(strings.NewReader(server.URL + "\n\n"))
	if err := runNewCommand(cmd, nil); err != nil {
		t.Fatal(err)
	}
	fm := readLinkFrontmatter(t, filepath.Join("bookmarks", "fetched-title.md"))
	if fm["draft"] != true || fm["private"] != true {
		t.Fatalf("fast creation replaced starter publication/privacy defaults: %#v", fm)
	}
}
