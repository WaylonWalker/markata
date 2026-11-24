package services

import (
	"context"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/listcache"
)

// buildService implements BuildService using lifecycle.Manager.
type buildService struct {
	manager     *lifecycle.Manager
	subscribers []chan BuildEvent
	mu          sync.RWMutex
}

// newBuildService creates a new BuildService.
func newBuildService(m *lifecycle.Manager) BuildService {
	return &buildService{
		manager:     m,
		subscribers: make([]chan BuildEvent, 0),
	}
}

// Build runs the build process.
func (s *buildService) Build(_ context.Context, opts BuildOptions) (*BuildResult, error) {
	start := time.Now()

	s.emit(BuildEvent{
		Type:    BuildEventStart,
		Message: "Starting build",
	})

	if opts.Concurrency > 0 {
		s.manager.SetConcurrency(opts.Concurrency)
	}

	err := s.manager.Run()

	result := &BuildResult{
		Success:        err == nil,
		Duration:       time.Since(start),
		PostsProcessed: len(s.manager.Posts()),
	}

	if err != nil {
		result.Errors = []error{err}
		s.emit(BuildEvent{
			Type:    BuildEventError,
			Message: err.Error(),
			Error:   err,
		})
	} else {
		s.emit(BuildEvent{
			Type:     BuildEventComplete,
			Message:  "Build complete",
			Progress: 100,
		})
	}

	// Collect warnings
	for _, w := range s.manager.Warnings() {
		result.Warnings = append(result.Warnings, w.Error())
	}

	return result, err
}

// LoadOnly runs only the load stage (for TUI browsing without full build).
func (s *buildService) LoadOnly(_ context.Context) error {
	return s.manager.RunTo(lifecycle.StageLoad)
}

// LoadForTUI runs through Collect stage for TUI browsing.
// This includes Transform (for stats, auto-titles) and Collect (for feeds).
func (s *buildService) LoadForTUI(ctx context.Context) error {
	if opts, ok := listcache.OptionsFromManager(s.manager); ok {
		return listcache.LoadOrRefresh(ctx, s.manager, opts)
	}
	return s.manager.RunTo(lifecycle.StageCollect)
}

// Subscribe returns a channel for build progress events.
func (s *buildService) Subscribe() <-chan BuildEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan BuildEvent, 10)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *buildService) emit(event BuildEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.subscribers {
		select {
		case ch <- event:
		default:
			// Don't block if channel is full
		}
	}
}
