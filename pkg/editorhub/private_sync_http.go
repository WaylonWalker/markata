package editorhub

import (
	"errors"
	"net/http"

	"github.com/WaylonWalker/markata/pkg/editorsync"
)

// NewPrivateSyncHandler composes only private-draft site discovery, writes,
// snapshots, incremental reads and lazy source hydration. It is unmounted.
// Never mount trusted-local workbench, LSP, Git or Builder Admin behind it.
//
// The caller must supply a trusted Authenticator and verify the real
// ForwardAuth/proxy boundary before mounting this behind a hosted origin.
func NewPrivateSyncHandler(
	registry *Registry,
	authenticate Authenticator,
	engine *editorsync.Engine,
	reader *editorsync.PrivateReader,
	allowedOrigin string,
) (http.Handler, error) {
	if registry == nil || authenticate == nil || engine == nil || reader == nil {
		return nil, errors.New("private sync requires registry, verified authenticator, writer and reader")
	}
	discovery, err := NewHandler(registry, authenticate)
	if err != nil {
		return nil, err
	}
	writes, err := NewPrivateMutationHandler(registry, authenticate, engine, allowedOrigin)
	if err != nil {
		return nil, err
	}
	reads, err := NewPrivateReadHandler(registry, authenticate, reader, allowedOrigin)
	if err != nil {
		return nil, err
	}

	// Each underlying route enforces site and subject permissions again.
	// No catch-all or fallback to an operator interface is registered.
	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/sites", discovery)
	mux.Handle("POST /api/v1/sites/{site}/mutations", writes)
	mux.Handle("GET /api/v1/sites/{site}/conflicts/{mutation}", writes)
	mux.Handle("POST /api/v1/sites/{site}/snapshots", reads)
	mux.Handle("GET /api/v1/sites/{site}/documents", reads)
	mux.Handle("GET /api/v1/sites/{site}/documents/{document}", reads)
	mux.Handle("GET /api/v1/sites/{site}/changes", reads)
	return mux, nil
}
