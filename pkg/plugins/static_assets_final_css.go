package plugins

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/templates"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// A versioned asset is immutable. Previous releases may still refer to it,
// so purge/minification must only process the unversioned working stylesheet.
var versionedCSSFilename = regexp.MustCompile(`\.[0-9a-fA-F]{8}\.css$`)

func isVersionedCSSFile(path string) bool {
	return versionedCSSFilename.MatchString(filepath.Base(path))
}

type finalCSSReplacement struct {
	matcher     *regexp.Regexp
	destination string
}

// rewriteCSSLinkReferences edits local stylesheet URLs without modifying
// authored prose, code, comments, external links, or unrelated link elements.
func rewriteCSSLinkReferences(input []byte, relativeHTML string, changes []finalCSSReplacement) ([]byte, error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(input))
	var output bytes.Buffer
	for {
		kind := tokenizer.Next()
		if kind == html.ErrorToken {
			if errors.Is(tokenizer.Err(), io.EOF) {
				return output.Bytes(), nil
			}
			return nil, fmt.Errorf("tokenize HTML for CSS references: %w", tokenizer.Err())
		}
		raw := tokenizer.Raw()
		if kind != html.StartTagToken && kind != html.SelfClosingTagToken {
			output.Write(raw)
			continue
		}
		token := tokenizer.Token()
		if token.DataAtom != atom.Link || !isStylesheetLink(token.Attr) {
			output.Write(raw)
			continue
		}
		changed := false
		for i := range token.Attr {
			if token.Attr[i].Key != "href" {
				continue
			}
			updated := rewriteLocalCSSHref(token.Attr[i].Val, relativeHTML, changes)
			if updated != token.Attr[i].Val {
				token.Attr[i].Val = updated
				changed = true
			}
		}
		if changed {
			output.WriteString(token.String())
		} else {
			output.Write(raw)
		}
	}
}

func isStylesheetLink(attrs []html.Attribute) bool {
	var rel, as string
	for _, attr := range attrs {
		switch attr.Key {
		case "rel":
			rel = strings.ToLower(attr.Val)
		case "as":
			as = attr.Val
		}
	}
	for _, item := range strings.Fields(rel) {
		if item == cssStylesheetRel || (item == "preload" && strings.EqualFold(as, "style")) {
			return true
		}
	}
	return false
}

// Match the resolved output-relative path against a *complete* registered
// asset name. An external URL or a coincidental filename suffix is not ours.
func rewriteLocalCSSHref(href, relativeHTML string, changes []finalCSSReplacement) string {
	u, err := url.Parse(href)
	if err != nil || u.IsAbs() || u.Host != "" || u.Opaque != "" || u.Path == "" || strings.HasPrefix(href, "//") {
		return href
	}
	local := strings.TrimPrefix(u.Path, "/")
	if !strings.HasPrefix(u.Path, "/") {
		local = path.Join(path.Dir(filepath.ToSlash(relativeHTML)), u.Path)
	}
	if !fs.ValidPath(local) {
		return href
	}
	for _, change := range changes {
		updated := change.matcher.ReplaceAllString(local, change.destination)
		if updated == local {
			continue
		}
		oldName := path.Base(u.Path)
		u.Path = strings.TrimSuffix(u.Path, oldName) + path.Base(updated)
		u.RawPath = ""
		return u.String()
	}
	return href
}

// reconcileFinalCSSAssets reconciles final published CSS with already-rendered
// HTML after purge and minification, including HTML restored from build cache.
func reconcileFinalCSSAssets(m *lifecycle.Manager, outputFS *os.Root, outputDir string) error {
	changes, finalHashes, err := collectFinalCSSChanges(m.AssetHashes(), outputFS, outputDir)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		return nil
	}
	updatedHTMLHashes, err := rewriteOutputCSSLinks(outputFS, outputDir, changes)
	if err != nil {
		return err
	}
	if err := refreshImageLibraryCSSOwnership(GetBuildCache(m), outputDir, updatedHTMLHashes); err != nil {
		return fmt.Errorf("refreshing CSS-rewritten image library ownership: %w", err)
	}
	for asset, hash := range finalHashes {
		m.SetAssetHash(asset, hash)
	}
	templates.SetAssetHashes(finalHashes)
	return nil
}

func collectFinalCSSChanges(hashes map[string]string, outputFS *os.Root, outputDir string) ([]finalCSSReplacement, map[string]string, error) {
	paths := make([]string, 0, len(hashes))
	for asset := range hashes {
		if strings.EqualFold(filepath.Ext(asset), ".css") {
			paths = append(paths, asset)
		}
	}
	sort.Strings(paths)
	changes := make([]finalCSSReplacement, 0, len(paths))
	finalHashes := make(map[string]string)
	for _, asset := range paths {
		relative, err := outputRelativePath(outputDir, filepath.Join(outputDir, asset))
		if err != nil {
			return nil, nil, err
		}
		data, err := outputFS.ReadFile(relative)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("reading final stylesheet %s: %w", asset, err)
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(data))[:8]
		ext := filepath.Ext(asset)
		base := strings.TrimSuffix(filepath.ToSlash(asset), ext)
		matcher, err := regexp.Compile("^" + regexp.QuoteMeta(base) + `\.[0-9a-fA-F]{8}` + regexp.QuoteMeta(ext) + "$")
		if err != nil {
			return nil, nil, err
		}
		changes = append(changes, finalCSSReplacement{matcher: matcher, destination: base + "." + hash + ext})
		if hash != hashes[asset] {
			finalHashes[asset] = hash
		}
	}
	return changes, finalHashes, nil
}

// Read/write only output-root-confined regular HTML files. Preserve output
// unchanged when no link needs rewriting, including warm builds.
func rewriteOutputCSSLinks(outputFS *os.Root, outputDir string, changes []finalCSSReplacement) (map[string]string, error) {
	updatedHTMLHashes := make(map[string]string)
	err := filepath.WalkDir(outputDir, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() || !strings.EqualFold(filepath.Ext(filePath), ".html") {
			return nil
		}
		relative, err := outputRelativePath(outputDir, filePath)
		if err != nil {
			return err
		}
		data, err := outputFS.ReadFile(relative)
		if err != nil {
			return err
		}
		after, err := rewriteCSSLinkReferences(data, relative, changes)
		if err != nil {
			return fmt.Errorf("reconcile CSS links in %s: %w", relative, err)
		}
		if bytes.Equal(data, after) {
			return nil
		}
		if err := outputFS.WriteFile(relative, after, 0o644); err != nil {
			return err
		}
		updatedHTMLHashes[filepath.Clean(relative)] = buildcache.ContentHash(string(after))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("updating HTML CSS URLs: %w", err)
	}
	return updatedHTMLHashes, nil
}

// refreshImageLibraryCSSOwnership keeps the image library's ownership hashes
// aligned with the final HTML after CSS links are rewritten in Cleanup. The
// images plugin registers its hashes during Write, so otherwise an unchanged
// warm build sees its own output as externally modified. Never adopt files
// that the image library did not already own.
func refreshImageLibraryCSSOwnership(cache *buildcache.Cache, outputDir string, updated map[string]string) error {
	if cache == nil || len(updated) == 0 {
		return nil
	}
	contentRoot, recordedOutput, owned := cache.GetImageLibraryOutputs()
	absoluteOutput, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("resolve CSS output directory for ownership: %w", err)
	}
	if recordedOutput == "" || filepath.Clean(recordedOutput) != filepath.Clean(absoluteOutput) {
		return nil
	}
	changed := false
	for relative, hash := range updated {
		if _, ok := owned[filepath.Clean(relative)]; ok {
			owned[filepath.Clean(relative)] = hash
			changed = true
		}
	}
	if !changed {
		return nil
	}
	absolute := make(map[string]string, len(owned))
	for relative, hash := range owned {
		absolute[filepath.Join(absoluteOutput, relative)] = hash
	}
	cache.SetImageLibraryOutputs(contentRoot, absoluteOutput, absolute)
	// The build-cache plugin saves earlier in Cleanup. Persist this post-write
	// correction so a subsequent process also recognizes Markata-owned output.
	return cache.Save()
}
