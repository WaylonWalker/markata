package workbench

import (
	"encoding/json"
	"net/http"
	"strings"
)

var pwaIconBytes = []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <rect width="512" height="512" rx="112" fill="#0f172a"/>
  <rect x="44" y="44" width="424" height="424" rx="84" fill="none" stroke="#f59e0b" stroke-width="20"/>
  <path d="M144 352V160h44l68 104 68-104h44v192h-44V230l-56 84h-24l-56-84v122z" fill="#f8fafc"/>
</svg>
`)

var pwaServiceWorkerBytes = []byte(`// Minimal Markata Editor PWA service worker.
// Keeps the app installable and enables OS Share Target while always
// preferring fresh network responses for workspace APIs and live previews.
self.addEventListener('install', (event) => {
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim());
});

self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') return;
  event.respondWith(fetch(event.request));
});
`)

type pwaIcon struct {
	Src     string `json:"src"`
	Sizes   string `json:"sizes"`
	Type    string `json:"type"`
	Purpose string `json:"purpose"`
}

type pwaShortcut struct {
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type pwaShareTargetParams struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type pwaShareTarget struct {
	Action string               `json:"action"`
	Method string               `json:"method"`
	Params pwaShareTargetParams `json:"params"`
}

type pwaManifest struct {
	Name            string         `json:"name"`
	ShortName       string         `json:"short_name"`
	Description     string         `json:"description"`
	StartURL        string         `json:"start_url"`
	Scope           string         `json:"scope"`
	Display         string         `json:"display"`
	BackgroundColor string         `json:"background_color"`
	ThemeColor      string         `json:"theme_color"`
	Icons           []pwaIcon      `json:"icons"`
	Shortcuts       []pwaShortcut  `json:"shortcuts"`
	ShareTarget     pwaShareTarget `json:"share_target"`
}

func (h *Handler) handlePWAManifest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	siteName := "Markata Editor"
	shortName := "Markata"
	if h.siteConfig != nil && strings.TrimSpace(h.siteConfig.Title) != "" {
		title := strings.TrimSpace(h.siteConfig.Title)
		siteName = title + " Editor"
		if len(title) <= 14 {
			shortName = title
		}
	}
	manifest := pwaManifest{
		Name:            siteName,
		ShortName:       shortName,
		Description:     "Mobile-first Markdown CMS and live editor for Markata sites",
		StartURL:        "/",
		Scope:           "/",
		Display:         "standalone",
		BackgroundColor: "#0f172a",
		ThemeColor:      "#0f172a",
		Icons: []pwaIcon{
			{
				Src:     "/pwa-icon.svg",
				Sizes:   "any",
				Type:    "image/svg+xml",
				Purpose: "any maskable",
			},
		},
		Shortcuts: []pwaShortcut{
			{
				Name:        "Today's Daily Note",
				ShortName:   "Daily",
				Description: "Open or create today's daily note",
				URL:         "/?flow=daily",
			},
			{
				Name:        "Quick Note",
				ShortName:   "Quick Note",
				Description: "Capture a quick thought or link",
				URL:         "/?flow=quick",
			},
			{
				Name:        "New Post",
				ShortName:   "New Post",
				Description: "Start a new post from a template",
				URL:         "/?flow=new",
			},
		},
		ShareTarget: pwaShareTarget{
			Action: "/",
			Method: "GET",
			Params: pwaShareTargetParams{
				Title: "share_title",
				Text:  "share_text",
				URL:   "share_url",
			},
		},
	}
	w.Header().Set("Content-Type", "application/manifest+json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	if err := json.NewEncoder(w).Encode(manifest); err != nil {
		return
	}
}

func (h *Handler) handlePWAServiceWorker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Service-Worker-Allowed", "/")
	if _, err := w.Write(pwaServiceWorkerBytes); err != nil {
		return
	}
}

func (h *Handler) handlePWAIcon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	if _, err := w.Write(pwaIconBytes); err != nil {
		return
	}
}
