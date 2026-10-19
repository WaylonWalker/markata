package templates

import (
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
	htmlparser "golang.org/x/net/html"
)

func editorialFontContext(picker bool, presentation string, populated bool) Context {
	config := models.NewConfig()
	config.Theme.Switcher.Enabled = &picker
	config.Extra = map[string]any{
		"fontpack_css":                    true,
		"fontpack_preload_urls":           []string{"/body.woff2"},
		"fontpack_editorial_preload_urls": []string{"/body.woff2", "/code.woff2"},
	}
	if picker {
		config.Extra["fontpack_preloads_json"] = `{"handwritten":["/body.woff2"],"technical":["/technical-body.woff2"],"system":[]}`
		config.Extra["fontpack_editorial_preloads_json"] = `{"handwritten":["/body.woff2","/code.woff2"],"technical":["/technical-body.woff2","/technical-code.woff2"],"system":[]}`
	}
	ctx := NewContext(nil, "", config)
	ctx.Set("feed_presentation", presentation)
	ctx.Feed = &models.FeedConfig{Views: []string{models.FeedViewDefault, models.FeedViewCalendar}}
	sections := []int{}
	if populated {
		sections = append(sections, 1)
	}
	ctx.Set("feed_sections", sections)
	return ctx
}

func TestEditorialFontPreloadStaticLayouts(t *testing.T) {
	for _, dir := range []string{"../../templates", "../themes/default/templates"} {
		engine, err := NewEngine(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, tc := range []struct {
			name, presentation string
			populated          bool
			want               []string
		}{
			{"nonempty editorial", "editorial", true, []string{"/body.woff2", "/code.woff2"}},
			{"empty editorial", "editorial", false, []string{"/body.woff2"}},
			{"classic feed", "classic", true, []string{"/body.woff2"}},
			{"ordinary page", "", false, []string{"/body.woff2"}},
		} {
			t.Run(dir+"/"+tc.name, func(t *testing.T) {
				rendered, err := engine.Render("base.html", editorialFontContext(false, tc.presentation, tc.populated))
				if err != nil {
					t.Fatal(err)
				}
				var got []string
				tokens := htmlparser.NewTokenizer(strings.NewReader(rendered))
				for {
					tt := tokens.Next()
					if tt == htmlparser.ErrorToken {
						if !errors.Is(tokens.Err(), io.EOF) {
							t.Fatal(tokens.Err())
						}
						break
					}
					if tt != htmlparser.StartTagToken && tt != htmlparser.SelfClosingTagToken {
						continue
					}
					token := tokens.Token()
					if token.Data != "link" {
						continue
					}
					attrs := map[string]string{}
					for _, a := range token.Attr {
						attrs[a.Key] = a.Val
					}
					if attrs["rel"] != "preload" || attrs["as"] != "font" {
						continue
					}
					if _, ok := attrs["crossorigin"]; !ok || attrs["type"] != "font/woff2" {
						t.Fatalf("unusable font preload: %v", attrs)
					}
					got = append(got, attrs["href"])
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("preloads = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

func TestEditorialFontBootstrapHonorsStoredPack(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node.js unavailable for bootstrap execution")
	}
	for _, dir := range []string{"../../templates", "../themes/default/templates"} {
		engine, err := NewEngine(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, tc := range []struct {
			name, presentation, stored string
			populated                  bool
			want                       []string
		}{
			{"default editorial", "editorial", "", true, []string{"/body.woff2", "/code.woff2"}},
			{"stored editorial", "editorial", "technical", true, []string{"/technical-body.woff2", "/technical-code.woff2"}},
			{"system editorial", "editorial", "system", true, []string{}},
			{"classic stored pack", "classic", "technical", true, []string{"/technical-body.woff2"}},
			{"empty stored editorial", "editorial", "technical", false, []string{"/technical-body.woff2"}},
		} {
			t.Run(dir+"/"+tc.name, func(t *testing.T) {
				rendered, err := engine.Render("partials/theme-bootstrap.html", editorialFontContext(true, tc.presentation, tc.populated))
				if err != nil {
					t.Fatal(err)
				}
				var scripts []string
				for _, match := range regexp.MustCompile(`(?s)<script>(.*?)</script>`).FindAllStringSubmatch(rendered, -1) {
					scripts = append(scripts, match[1])
				}
				input, err := json.Marshal(map[string]any{"scripts": scripts, "stored": tc.stored, "calendar": tc.presentation == "editorial" && tc.populated})
				if err != nil {
					t.Fatal(err)
				}
				// Execute the rendered bootstrap, including its existing picker restoration.
				cmd := exec.Command(node, "-e", `
const fs=require('node:fs'),vm=require('node:vm'),input=JSON.parse(fs.readFileSync(0,'utf8')),links=[],classes=new Set();
const sandbox={window:{},localStorage:{getItem:key=>key==='theme-fontpack'?input.stored:null},document:{documentElement:{dataset:{fontpack:'handwritten'},classList:{toggle(){},add:name=>classes.add(name)}},createElement:()=>({}),head:{appendChild:link=>links.push(link)}}};
for(const script of input.scripts)vm.runInNewContext(script,sandbox);
for(const link of links)if(link.rel!=='preload'||link.as!=='font'||link.type!=='font/woff2'||link.crossOrigin!=='anonymous')throw Error('Invalid preload');
if(classes.has('editorial-calendar-js')!==input.calendar)throw Error('Wrong calendar enhancement marker');
process.stdout.write(JSON.stringify(links.map(link=>link.href)));`)
				cmd.Stdin = strings.NewReader(string(input))
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("bootstrap execution: %v: %s", err, output)
				}
				var got []string
				if err := json.Unmarshal(output, &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("preloads = %v, want %v", got, tc.want)
				}
			})
		}
	}
}

func TestEditorialCalendarMarkerRequiresCalendarView(t *testing.T) {
	for _, dir := range []string{"../../templates", "../themes/default/templates"} {
		engine, err := NewEngine(dir)
		if err != nil {
			t.Fatal(err)
		}
		ctx := editorialFontContext(true, "editorial", true)
		ctx.Feed.Views = []string{models.FeedViewDefault}
		rendered, err := engine.Render("partials/theme-bootstrap.html", ctx)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(rendered, "classList.add('editorial-calendar-js')") {
			t.Fatal("non-calendar feed got a calendar marker")
		}
	}
}
