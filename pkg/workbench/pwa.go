package workbench

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"strings"
	"sync"
)

const (
	pwaImageType        = "image/png"
	rasterLargeIconPath = "/pwa-icon-512.png"
	pwaSharePhotosField = "images"
	pwaGalleryFlagTrue  = "true"
)

var pwaIconBytes = []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <rect width="512" height="512" rx="112" fill="#0f172a"/>
  <rect x="44" y="44" width="424" height="424" rx="84" fill="none" stroke="#f59e0b" stroke-width="20"/>
  <path d="M144 352V160h44l68 104 68-104h44v192h-44V230l-56 84h-24l-56-84v122z" fill="#f8fafc"/>
</svg>
`)

// Raster icons back up the SVG to satisfy Chromium's installability checks.
var (
	pwaRasterOnce  sync.Once
	pwaRasterIcons map[int][]byte
)

type pwaPoint struct{ x, y float64 }

func pwaRounded(x, y, left, top, right, bottom, radius float64) bool {
	if x < left || x > right || y < top || y > bottom {
		return false
	}
	cx := math.Max(left+radius, math.Min(right-radius, x))
	cy := math.Max(top+radius, math.Min(bottom-radius, y))
	dx, dy := x-cx, y-cy
	return dx*dx+dy*dy <= radius*radius
}

func pwaPolygon(x, y float64, p []pwaPoint) bool {
	inside := false
	for i, j := 0, len(p)-1; i < len(p); j, i = i, i+1 {
		a, b := p[i], p[j]
		if (a.y > y) != (b.y > y) && x < (b.x-a.x)*(y-a.y)/(b.y-a.y)+a.x {
			inside = !inside
		}
	}
	return inside
}

func pwaRasterPNG(size int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	navy := color.RGBA{15, 23, 42, 255}
	gold := color.RGBA{245, 158, 11, 255}
	white := color.RGBA{248, 250, 252, 255}
	glyph := []pwaPoint{
		{144, 352}, {144, 160}, {188, 160}, {256, 264},
		{324, 160}, {368, 160}, {368, 352}, {324, 352},
		{324, 230}, {268, 314}, {244, 314}, {188, 230}, {188, 352},
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx := (float64(x) + .5) * 512 / float64(size)
			sy := (float64(y) + .5) * 512 / float64(size)
			c := navy
			if pwaRounded(sx, sy, 44, 44, 468, 468, 84) && !pwaRounded(sx, sy, 64, 64, 448, 448, 64) {
				c = gold
			}
			if pwaPolygon(sx, sy, glyph) {
				c = white
			}
			img.SetRGBA(x, y, c)
		}
	}
	var encoded bytes.Buffer
	if png.Encode(&encoded, img) != nil {
		return nil
	}
	return encoded.Bytes()
}

func (h *Handler) handlePWARasterIcon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	pwaRasterOnce.Do(func() {
		pwaRasterIcons = map[int][]byte{192: pwaRasterPNG(192), 512: pwaRasterPNG(512)}
	})
	size := 192
	if r.URL.Path == rasterLargeIconPath {
		size = 512
	}
	data := pwaRasterIcons[size]
	if len(data) == 0 {
		http.Error(w, "cannot generate PWA icon", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", pwaImageType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	if r.Method == http.MethodGet {
		if _, err := w.Write(data); err != nil {
			return
		}
	}
}

var pwaServiceWorkerBytes = []byte(`// Markata Web Share Target. POST bodies never reach the workspace server:
// stage images temporarily in this origin's IndexedDB and redirect to the
// editor, which asks the user to review/save using its existing guarded API.
const SHARE_DB = 'markata-share-captures-v1';
const SHARE_STORE = 'captures';
const MAX_FILES = 8;
const MAX_IMAGE = 15 * 1024 * 1024;
const MAX_TOTAL = 48 * 1024 * 1024;
const SHARE_TTL = 24 * 60 * 60 * 1000;

function openShareDatabase() {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(SHARE_DB, 1);
    req.onupgradeneeded = () => req.result.createObjectStore(SHARE_STORE, { keyPath: 'id' });
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

async function stageShare(capture) {
  const db = await openShareDatabase();
  try {
    await new Promise((resolve, reject) => {
      const tx = db.transaction(SHARE_STORE, 'readwrite');
      const store = tx.objectStore(SHARE_STORE);
      // Cleanup expired captures during every share; canceled previews
      // should not keep private photos in browser storage forever.
      const cursor = store.openCursor();
      cursor.onsuccess = () => {
        const item = cursor.result;
        if (!item) return;
        if (!item.value?.createdAt || Date.now() - item.value.createdAt > SHARE_TTL) {
          item.delete();
        }
        item.continue();
      };
      store.put(capture);
      tx.oncomplete = resolve;
      tx.onerror = () => reject(tx.error);
      tx.onabort = () => reject(tx.error || Error('Cannot stage shared photo'));
    });
  } finally {
    db.close();
  }
}

self.addEventListener('install', () => self.skipWaiting());
self.addEventListener('activate', (event) => event.waitUntil(self.clients.claim()));

self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url);
  if (url.origin === self.location.origin && url.pathname === '/share-target'
      && event.request.method === 'POST') {
    event.respondWith((async () => {
      try {
        const form = await event.request.formData();
        const photos = form.getAll('images').filter((item) => item instanceof File && item.size);
        const mediaType = (photo) => {
          const suffix = String(photo.name || '').split('.').pop().toLowerCase();
          const types = {png: 'image/png', jpg: 'image/jpeg', jpeg: 'image/jpeg',
            gif: 'image/gif', webp: 'image/webp'};
          return photo.type || types[suffix] || '';
        };
        if (photos.length > MAX_FILES) throw Error('Share a maximum of 8 photos at once');
        const size = photos.reduce((sum, photo) => sum + photo.size, 0);
        if (size > MAX_TOTAL || photos.some((photo) => photo.size > MAX_IMAGE)) {
          throw Error('Each photo must be 15 MiB or less; 48 MiB total');
        }
        if (photos.some((photo) => !['image/png','image/jpeg','image/gif','image/webp'].includes(mediaType(photo)))) {
          throw Error('Only PNG, JPEG, GIF, and WebP photos are supported');
        }
        const id = crypto.randomUUID();
        await stageShare({
          id, createdAt: Date.now(),
          title: String(form.get('share_title') || '').slice(0, 2048),
          text: String(form.get('share_text') || '').slice(0, 200000),
          url: String(form.get('share_url') || '').slice(0, 4096),
          files: photos,
        });
        if (event.request.mode === 'navigate') {
          return Response.redirect(new URL('/?share_id=' + encodeURIComponent(id), self.location.origin).href, 303);
        }
        return new Response(JSON.stringify({ share_id: id }), {
          status: 201, headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' },
        });
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Could not receive shared photos';
        if (event.request.mode === 'navigate') {
          return Response.redirect(new URL('/?share_error=' + encodeURIComponent(message), self.location.origin).href, 303);
        }
        return new Response(JSON.stringify({ error: message }), {
          status: 400, headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' },
        });
      }
    })());
    return;
  }
  // Non-share requests go directly to the network; no GET interception.
  // This keeps editor APIs responsive and makes browser request routing reliable.
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

type pwaShareTargetFile struct {
	Name   string   `json:"name"`
	Accept []string `json:"accept"`
}

type pwaShareTarget struct {
	Action  string               `json:"action"`
	Method  string               `json:"method"`
	Enctype string               `json:"enctype"`
	Params  pwaShareTargetParams `json:"params"`
	Files   []pwaShareTargetFile `json:"files"`
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
			{Src: "/pwa-icon-192.png", Sizes: "192x192", Type: pwaImageType, Purpose: "any"},
			{Src: rasterLargeIconPath, Sizes: "512x512", Type: pwaImageType, Purpose: "any maskable"},
			{Src: "/pwa-icon.svg", Sizes: "any", Type: "image/svg+xml", Purpose: "any maskable"},
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
			Action:  "/share-target",
			Method:  "POST",
			Enctype: "multipart/form-data",
			Files: []pwaShareTargetFile{{
				Name:   pwaSharePhotosField,
				Accept: []string{"image/png", "image/jpeg", "image/gif", "image/webp"},
			}},
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
