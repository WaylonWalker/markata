package plugins

import (
 "crypto/sha256"
 "fmt"
 "io/fs"
 "os"
 "path/filepath"
 "strings"

 "github.com/WaylonWalker/markata/pkg/lifecycle"
 "github.com/WaylonWalker/markata/pkg/templates"
)

// reconcileFinalCSSAssets runs after css_purge has rewritten output CSS.
// Templates are rendered earlier, so both the registered asset hashes and
// previously rendered HTML references must be updated together. Hashed copies
// are written by StaticAssetsPlugin.Cleanup after this reconciliation.
func reconcileFinalCSSAssets(m *lifecycle.Manager, outputFS *os.Root, outputDir string) error {
 type replacement struct{ from, to string }
 changes := make([]replacement, 0)
 finalHashes := make(map[string]string)
 for path, oldHash := range m.AssetHashes() {
  if !strings.EqualFold(filepath.Ext(path), ".css") {
   continue
  }
  relative, err := outputRelativePath(outputDir, filepath.Join(outputDir, path))
  if err != nil { return err }
  data, err := outputFS.ReadFile(relative)
  if os.IsNotExist(err) { continue }
  if err != nil { return fmt.Errorf("reading final stylesheet %s: %w", path, err) }
  hash := fmt.Sprintf("%x", sha256.Sum256(data))[:8]
  if hash == oldHash { continue }
  ext := filepath.Ext(path)
  base := strings.TrimSuffix(filepath.ToSlash(path), ext)
  changes = append(changes, replacement{
   from: base+"."+oldHash+ext,
   to: base+"."+hash+ext,
  })
  finalHashes[path]=hash
 }
 if len(changes)==0 { return nil }

 // Rewrite generated HTML references before publishing the final filenames.
 // Read/write through the descriptor-backed output root to prevent traversal.
 err := filepath.WalkDir(outputDir, func(path string, entry fs.DirEntry, walkErr error) error {
  if walkErr!=nil { return walkErr }
  if entry.IsDir() { return nil }
  if entry.Type()&os.ModeSymlink!=0 || !strings.EqualFold(filepath.Ext(path),".html") { return nil }
  relative, err := outputRelativePath(outputDir,path)
  if err!=nil { return err }
  data, err := outputFS.ReadFile(relative)
  if err!=nil { return err }
  before:=string(data)
  after:=before
  for _,change:=range changes { after=strings.ReplaceAll(after,change.from,change.to) }
  if after==before { return nil }
  return outputFS.WriteFile(relative,[]byte(after),0o644)
 })
 if err!=nil { return fmt.Errorf("updating HTML CSS URLs: %w",err) }
 for path,hash:=range finalHashes { m.SetAssetHash(path,hash) }
 templates.SetAssetHashes(finalHashes)
 return nil
}
