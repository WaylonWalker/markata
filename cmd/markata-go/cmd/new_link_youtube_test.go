package cmd

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

func stubNewLinkTransport(t *testing.T, handler func(*http.Request) (int, string)) {
	t.Helper()
	previous := http.DefaultTransport
	http.DefaultTransport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		status, body := handler(request)
		return &http.Response{
			StatusCode: status,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    request,
		}, nil
	})
	t.Cleanup(func() { http.DefaultTransport = previous })
}

func TestRunNewCommand_LinkYouTubeOEmbed(t *testing.T) {
	const videoTitle = `"Nobody's Going to Write the Next React" ft. Armin Ronacher & Mario Zechner - Ep 08`
	const thumbnail = "https://i.ytimg.com/vi/kODgReWwrWM/hqdefault.jpg"
	for _, tt := range []struct {
		name     string
		url      string
		override string
	}{
		{"watch", "https://www.youtube.com/watch?v=kODgReWwrWM", ""},
		{"short URL", "https://youtu.be/kODgReWwrWM?t=60", ""},
		{"mobile", "https://m.youtube.com/watch?v=kODgReWwrWM", ""},
		{"shorts", "https://www.youtube.com/shorts/kODgReWwrWM", ""},
		{"explicit title", "https://www.youtube.com/watch?v=kODgReWwrWM", "My chosen title"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd, stderr := linkTestCommand(t)
			noInput = true
			var requests []string
			stubNewLinkTransport(t, func(request *http.Request) (int, string) {
				requests = append(requests, request.URL.String())
				if request.URL.Host == "www.youtube.com" && request.URL.Path == "/oembed" {
					if request.URL.Query().Get("url") != tt.url || request.URL.Query().Get("format") != "json" {
						t.Errorf("oEmbed query did not preserve the destination: %s", request.URL)
					}
					if request.Header.Get("Accept") != "application/json" {
						t.Errorf("oEmbed Accept = %q", request.Header.Get("Accept"))
					}
					return http.StatusOK, `{"title":"\u0022Nobody\u0027s Going to Write the Next React\u0022 ft. Armin Ronacher \u0026 Mario Zechner - Ep 08","type":"video","thumbnail_url":"` + thumbnail + `"}`
				}
				return http.StatusOK, "<html>" + strings.Repeat("x", 1024*1024+1) + "</html>"
			})
			flags := []string{"--template", "link"}
			wantTitle := videoTitle
			if tt.override != "" {
				flags = append(flags, "--title", tt.override)
				wantTitle = tt.override
			}
			if err := cmd.ParseFlags(flags); err != nil {
				t.Fatal(err)
			}
			if err := runNewCommand(cmd, []string{tt.url}); err != nil {
				t.Fatalf("YouTube link creation failed: %v; warnings: %s", err, stderr)
			}
			fm := readLinkFrontmatter(t, filepath.Join("pages", "link", generateSlug(wantTitle)+".md"))
			if fm["title"] != wantTitle || fm["image"] != thumbnail || fm["url"] != tt.url {
				t.Fatalf("YouTube defaults or overrides lost: %#v", fm)
			}
			if len(requests) != 1 || !strings.Contains(requests[0], "/oembed?") {
				t.Errorf("downloaded a watch page instead of only oEmbed: %v", requests)
			}
			if stderr.Len() != 0 {
				t.Errorf("successful oEmbed fetch warned: %s", stderr)
			}
		})
	}
}

func TestRunNewCommand_LinkYouTubeOEmbedFallback(t *testing.T) {
	cmd, stderr := linkTestCommand(t)
	noInput = true
	var paths []string
	stubNewLinkTransport(t, func(request *http.Request) (int, string) {
		paths = append(paths, request.URL.Path)
		if request.URL.Path == "/oembed" {
			return http.StatusServiceUnavailable, "Unavailable"
		}
		return http.StatusOK, `<meta property="og:title" content="HTML fallback title">`
	})
	if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
		t.Fatal(err)
	}
	destination := "https://www.youtube.com/watch?v=kODgReWwrWM"
	if err := runNewCommand(cmd, []string{destination}); err != nil {
		t.Fatal(err)
	}
	fm := readLinkFrontmatter(t, filepath.Join("pages", "link", "html-fallback-title.md"))
	if fm["url"] != destination || len(paths) != 2 || paths[0] != "/oembed" {
		t.Fatalf("fallback did not preserve destination or provider priority: %#v, %v", fm, paths)
	}
	if !strings.Contains(stderr.String(), "503") {
		t.Fatalf("provider failure was hidden: %s", stderr)
	}
}

func TestRunNewCommand_LinkYouTubeUnavailableMetadata(t *testing.T) {
	for _, response := range []struct {
		name   string
		status int
		body   string
	}{
		{"missing title", http.StatusOK, `{"thumbnail_url":"https://i.ytimg.com/vi/kODgReWwrWM/hqdefault.jpg"}`},
		{"invalid JSON", http.StatusOK, "not JSON"},
		{"oversized response", http.StatusOK, strings.Repeat("x", 1024*1024+1)},
		{"private video", http.StatusUnauthorized, "Unauthorized"},
	} {
		t.Run(response.name, func(t *testing.T) {
			cmd, stderr := linkTestCommand(t)
			noInput = true
			stubNewLinkTransport(t, func(request *http.Request) (int, string) {
				if request.URL.Path == "/oembed" {
					return response.status, response.body
				}
				return http.StatusOK, strings.Repeat("x", 1024*1024+1)
			})
			if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
				t.Fatal(err)
			}
			err := runNewCommand(cmd, []string{"https://www.youtube.com/watch?v=kODgReWwrWM"})
			if err == nil || !strings.Contains(err.Error(), "--title='YouTube video kODgReWwrWM'") {
				t.Fatalf("unavailable metadata needs a useful, untrusted suggestion: %v", err)
			}
			if !strings.Contains(stderr.String(), "YouTube") || !strings.Contains(stderr.String(), "1 MiB") {
				t.Errorf("provider/page errors were not both explained: %s", stderr)
			}
		})
	}
}

func TestRunNewCommand_LinkYouTubeLookalikeUsesHTML(t *testing.T) {
	cmd, _ := linkTestCommand(t)
	noInput = true
	stubNewLinkTransport(t, func(request *http.Request) (int, string) {
		if request.URL.Host != "youtube.com.example.org" {
			t.Errorf("lookalike hostname used a provider endpoint: %s", request.URL)
		}
		return http.StatusOK, `<meta property="og:title" content="An ordinary page">`
	})
	if err := cmd.ParseFlags([]string{"--template", "link"}); err != nil {
		t.Fatal(err)
	}
	if err := runNewCommand(cmd, []string{"https://youtube.com.example.org/watch?v=kODgReWwrWM"}); err != nil {
		t.Fatal(err)
	}
	readLinkFrontmatter(t, filepath.Join("pages", "link", "an-ordinary-page.md"))
}
