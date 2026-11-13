package workbench

import (
	"net/http"
	"sync"
)

// deferredLanguage permits the editor shell and file endpoints to come online
// while the LSP backend warms up. Missing LSP must not prevent writing.
type deferredLanguage struct {
	mu      sync.RWMutex
	handler http.Handler
}

func (d *deferredLanguage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	handler := d.handler
	d.mu.RUnlock()
	if handler == nil {
		w.Header().Set("Retry-After", "1")
		http.Error(w, "language service initializing", http.StatusServiceUnavailable)
		return
	}
	handler.ServeHTTP(w, r)
}

func (d *deferredLanguage) ready(handler http.Handler) {
	d.mu.Lock()
	d.handler = handler
	d.mu.Unlock()
}
