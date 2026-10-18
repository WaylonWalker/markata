package builderadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/WaylonWalker/markata/pkg/servecontrol"
)

type adminStateResponse struct {
	State         State                 `json:"state"`
	Revision      RevisionStatus        `json:"revision"`
	Releases      []ReleaseView         `json:"releases"`
	CurrentID     string                `json:"current_release_id,omitempty"`
	CurrentPath   string                `json:"current_release_path,omitempty"`
	Config        Config                `json:"config"`
	RefreshTasks  []string              `json:"refresh_tasks"`
	ControlCenter servecontrol.Snapshot `json:"control_center"`
}

type apiErrorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiOperationResponse struct {
	Operation QueuedOperation `json:"operation"`
	Status    string          `json:"status"`
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apiErrorResponse{Error: apiError{Code: code, Message: message}})
}

func writeAPIJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func apiMethodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed.")
}

func (s *Service) apiActorID(r *http.Request) (string, error) {
	profile, err := operatorProfileFromHeaders(r.Header, s.cfg.AuthHeaders)
	if err != nil {
		return "", err
	}
	return profile.UserID, nil
}

func (s *Service) handleAPIState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiMethodNotAllowed(w, http.MethodGet)
		return
	}
	writeAPIJSON(w, http.StatusOK, s.stateResponse())
}

func (s *Service) handleAPIBuildDetail(w http.ResponseWriter, r *http.Request) {
	s.handleAPIOperationDetail(w, r, "/builds/", "build")
}

func (s *Service) handleAPIRefreshDetail(w http.ResponseWriter, r *http.Request) {
	s.handleAPIOperationDetail(w, r, "/refreshes/", "refresh")
}

func (s *Service) handleAPIOperationDetail(w http.ResponseWriter, r *http.Request, prefix, kind string) {
	if r.Method != http.MethodGet {
		apiMethodNotAllowed(w, http.MethodGet)
		return
	}
	resourceName := "Build"
	if kind == "refresh" {
		resourceName = "Refresh"
	}
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, prefix), "/")
	if id == "" || strings.Contains(id, "/") {
		writeAPIError(w, http.StatusNotFound, "not_found", resourceName+" not found.")
		return
	}
	state := s.viewState()
	if kind == "build" {
		for _, build := range state.Builds {
			if build.ID == id {
				writeAPIJSON(w, http.StatusOK, build)
				return
			}
		}
	} else {
		for _, refresh := range state.Refresh {
			if refresh.ID == id {
				writeAPIJSON(w, http.StatusOK, refresh)
				return
			}
		}
	}
	if state.Running != nil && state.Running.ID == id && state.Running.Kind == kind {
		writeAPIJSON(w, http.StatusOK, map[string]any{"operation": state.Running, "status": "running"})
		return
	}
	for _, queued := range state.Queue {
		if queued.ID == id && queued.Kind == kind {
			writeAPIJSON(w, http.StatusOK, apiOperationResponse{Operation: queued, Status: "queued"})
			return
		}
	}
	writeAPIError(w, http.StatusNotFound, "not_found", resourceName+" not found.")
}

func (s *Service) handleAPIBuilds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiMethodNotAllowed(w, http.MethodPost)
		return
	}
	if s.handleStandbyMutation(w, r) {
		return
	}
	if err := decodeEmptyAPIRequest(w, r); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	actorID, err := s.apiActorID(r)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, "authentication_required", "A trusted operator identity is required.")
		return
	}
	operation, err := s.enqueueBuildAs("agent-api", "Builder-admin API build", nil, actorID)
	if err != nil {
		writeAPIQueueError(w, err)
		return
	}
	writeAPIJSON(w, http.StatusAccepted, apiOperationResponse{Operation: operation, Status: "queued"})
}

func (s *Service) handleAPIRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiMethodNotAllowed(w, http.MethodPost)
		return
	}
	if s.handleStandbyMutation(w, r) {
		return
	}
	if err := decodeEmptyAPIRequest(w, r); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	taskName := strings.Trim(strings.TrimPrefix(r.URL.Path, "/refresh/"), "/")
	if taskName == "" || strings.Contains(taskName, "/") {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "A refresh task name is required.")
		return
	}
	actorID, err := s.apiActorID(r)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, "authentication_required", "A trusted operator identity is required.")
		return
	}
	operation, err := s.enqueueRefreshAs(taskName, "agent-api", "Builder-admin API refresh", actorID)
	if err != nil {
		if errors.Is(err, ErrBuildQueueFull) {
			writeAPIQueueError(w, err)
			return
		}
		writeAPIError(w, http.StatusNotFound, "refresh_task_not_found", err.Error())
		return
	}
	writeAPIJSON(w, http.StatusAccepted, apiOperationResponse{Operation: operation, Status: "queued"})
}

func (s *Service) handleAPIReleaseAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apiMethodNotAllowed(w, http.MethodPost)
		return
	}
	if s.handleStandbyMutation(w, r) {
		return
	}
	if err := decodeEmptyAPIRequest(w, r); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	trimmed := strings.Trim(strings.TrimPrefix(r.URL.Path, "/releases/"), "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] != "rollback" {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Expected /releases/{release-id}/rollback.")
		return
	}
	releaseID := parts[0]
	found := false
	for _, release := range s.discoverReleases() {
		if release.ID == releaseID {
			found = true
			break
		}
	}
	if !found {
		writeAPIError(w, http.StatusNotFound, "release_not_found", fmt.Sprintf("Release %q was not found.", releaseID))
		return
	}
	actorID, err := s.apiActorID(r)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, "authentication_required", "A trusted operator identity is required.")
		return
	}
	operation, err := s.enqueueRollbackAs(releaseID, "agent-api", "Builder-admin API rollback", actorID)
	if err != nil {
		if errors.Is(err, ErrBuildQueueFull) {
			writeAPIQueueError(w, err)
			return
		}
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	writeAPIJSON(w, http.StatusAccepted, apiOperationResponse{Operation: operation, Status: "queued"})
}

func (s *Service) handleAPILogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiMethodNotAllowed(w, http.MethodGet)
		return
	}
	name := strings.Trim(strings.TrimPrefix(r.URL.Path, "/logs/"), "/")
	if name == "" || strings.Contains(name, "/") {
		writeAPIError(w, http.StatusNotFound, "not_found", "Log not found.")
		return
	}
	s.handleLogs(w, r)
}

func decodeEmptyAPIRequest(w http.ResponseWriter, r *http.Request) error {
	if r.Body == nil || r.ContentLength == 0 {
		return nil
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request map[string]json.RawMessage
	if err := decoder.Decode(&request); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("request body must be an empty JSON object: %w", err)
	} else if errors.Is(err, io.EOF) {
		return nil
	}
	if request == nil || len(request) != 0 {
		return fmt.Errorf("request body must be an empty JSON object")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("request body must contain one JSON object")
	}
	return nil
}

func writeAPIQueueError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrBuildQueueFull) {
		w.Header().Set("Retry-After", "30")
		writeAPIError(w, http.StatusServiceUnavailable, "queue_full", "Builder-admin queue is full; retry later.")
		return
	}
	writeAPIError(w, http.StatusInternalServerError, "operation_failed", err.Error())
}
