package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/WaylonWalker/markata/pkg/config"
)

func configFilesHash(cfgPath string, mergeFiles []string) (string, error) {
	paths := resolveConfigPaths(cfgPath, mergeFiles)
	if len(paths) == 0 {
		return "", nil
	}

	h := sha256.New()
	for _, path := range paths {
		if err := hashFile(h, path); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func resolveConfigPaths(cfgPath string, mergeFiles []string) []string {
	paths := make([]string, 0, 1+len(mergeFiles))
	basePath := cfgPath
	if basePath == "" {
		discovered, err := config.Discover()
		if err == nil {
			basePath = discovered
		}
	}
	if basePath != "" {
		paths = append(paths, basePath)
	}
	paths = append(paths, mergeFiles...)
	return paths
}

func hashFile(w io.Writer, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %s: %w", path, err)
	}
	if _, err := io.WriteString(w, path); err != nil {
		return err
	}
	if _, err := w.Write([]byte{0}); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	if _, err := w.Write([]byte{0}); err != nil {
		return err
	}
	return nil
}
