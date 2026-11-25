package builderadmin

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func workspaceReuseMarkerPath(workDir string) string {
	return workDir + ".base-release"
}

// reusableWorkspaceMatches checks retained workspace proof without consuming it.
func reusableWorkspaceMatches(workDir, currentReleaseID string) (bool, error) {
	if strings.TrimSpace(currentReleaseID) == "" {
		return false, nil
	}
	info, err := os.Stat(workDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("inspect reusable workspace: %w", err)
	}
	if !info.IsDir() {
		return false, nil
	}
	data, err := os.ReadFile(workspaceReuseMarkerPath(workDir))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("read reusable workspace marker: %w", err)
	}
	return strings.TrimSpace(string(data)) == currentReleaseID, nil
}

// claimReusableWorkspace returns true only when the existing workspace was
// produced by the currently-live release. The marker is consumed before the
// caller mutates the workspace so a failed or interrupted build can never be
// mistaken for a clean reusable base on the next run.
func claimReusableWorkspace(workDir, currentReleaseID string) (bool, error) {
	if strings.TrimSpace(currentReleaseID) == "" {
		return false, nil
	}
	info, err := os.Stat(workDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("inspect reusable workspace: %w", err)
	}
	if !info.IsDir() {
		return false, nil
	}

	markerPath := workspaceReuseMarkerPath(workDir)
	data, err := os.ReadFile(markerPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("read reusable workspace marker: %w", err)
	}
	if strings.TrimSpace(string(data)) != currentReleaseID {
		return false, nil
	}
	if err := os.Remove(markerPath); err != nil {
		return false, fmt.Errorf("claim reusable workspace: %w", err)
	}
	return true, nil
}

// markReusableWorkspace records that the retained node-local workspace is an
// exact successful build of releaseID. Staged incremental publication keeps
// this workspace independent of every immutable release.
func markReusableWorkspace(workDir, releaseID string) error {
	markerPath := workspaceReuseMarkerPath(workDir)
	info, err := os.Stat(workDir)
	if err != nil {
		if os.IsNotExist(err) {
			return removeWorkspaceReuseMarker(markerPath)
		}
		return fmt.Errorf("inspect promoted workspace: %w", err)
	}
	if !info.IsDir() || strings.TrimSpace(releaseID) == "" {
		return removeWorkspaceReuseMarker(markerPath)
	}

	tmp := markerPath + ".tmp"
	if err := os.WriteFile(tmp, []byte(releaseID+"\n"), 0o600); err != nil {
		return fmt.Errorf("write reusable workspace marker: %w", err)
	}
	if err := os.Rename(tmp, markerPath); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("install reusable workspace marker: %w", err)
	}
	return nil
}

func invalidateReusableWorkspace(workDir string) error {
	return removeWorkspaceReuseMarker(workspaceReuseMarkerPath(workDir))
}

func removeWorkspaceReuseMarker(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
