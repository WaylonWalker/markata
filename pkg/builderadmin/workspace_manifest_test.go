package builderadmin

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
)

func TestPublicationManifest_WarmAndRestoredModificationTime(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux change-time cache")
	}
	root := t.TempDir()
	work, releases := filepath.Join(root, "work"), filepath.Join(root, "releases")
	writeDeltaFile(t, work, "post.html", "before")
	time.Sleep(1100 * time.Millisecond) // Metadata must precede manifest by a full timestamp tick.
	first, stats, err := stageIncrementalWorkspace(work, releases, "first", "")
	if err != nil || !stats.ManifestSaved {
		t.Fatalf("prime: %+v %v", stats, err)
	}
	original, err := os.Stat(filepath.Join(work, "post.html"))
	if err != nil {
		t.Fatal(err)
	}
	_, stats, err = stageIncrementalWorkspace(work, releases, "second", "first")
	if err != nil || stats.LinkedFiles != 1 || stats.ComparedBytes != 0 || stats.CachedSourceHashes != 1 || stats.CachedReleaseHashes != 1 {
		t.Fatalf("warm: %+v %v", stats, err)
	}
	writeDeltaFile(t, work, "post.html", "after!")
	if err := os.Chtimes(filepath.Join(work, "post.html"), original.ModTime(), original.ModTime()); err != nil {
		t.Fatal(err)
	}
	third, stats, err := stageIncrementalWorkspace(work, releases, "third", "second")
	if err != nil || stats.CopiedFiles != 1 || stats.CachedSourceHashes != 0 || stats.ComparedBytes != 6 {
		t.Fatalf("edited: %+v %v", stats, err)
	}
	assertWorkspaceContent(t, filepath.Join(first, "post.html"), "before")
	assertWorkspaceContent(t, filepath.Join(third, "post.html"), "after!")
}

func TestPublicationManifest_InvalidCacheFallsBack(t *testing.T) {
	for _, scenario := range []string{"corrupt", "wrong-release", "wrong-workspace", "write-failure"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			work, releases := filepath.Join(root, "work"), filepath.Join(root, "releases")
			writeDeltaFile(t, work, "post.html", "content")
			if _, _, err := stageIncrementalWorkspace(work, releases, "first", ""); err != nil {
				t.Fatal(err)
			}
			manifestPath := work + ".publication.json"
			switch scenario {
			case "corrupt":
				if err := os.WriteFile(manifestPath, []byte(`{"Version":1,"Payload":{},"Checksum":"bad"}`), 0o600); err != nil {
					t.Fatal(err)
				}
			case "wrong-release":
				if !savePublicationManifest(work, "other", nil) {
					t.Fatal("save manifest")
				}
			case "wrong-workspace":
				other := filepath.Join(root, "other")
				if !savePublicationManifest(other, "first", nil) {
					t.Fatal("save manifest")
				}
				if err := os.Rename(other+".publication.json", manifestPath); err != nil {
					t.Fatal(err)
				}
			case "write-failure":
				if err := os.Remove(manifestPath); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(manifestPath, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			published, stats, err := stageIncrementalWorkspace(work, releases, "second", "first")
			if err != nil || stats.LinkedFiles != 1 || stats.ComparedBytes != 14 || stats.CachedSourceHashes != 0 || stats.CachedReleaseHashes != 0 {
				t.Fatalf("fallback: %+v %v", stats, err)
			}
			if scenario == "write-failure" && stats.ManifestSaved {
				t.Fatal("reported successful manifest persistence")
			}
			assertWorkspaceContent(t, filepath.Join(published, "post.html"), "content")
			if _, err := os.Stat(filepath.Join(published, ".publication.json")); !os.IsNotExist(err) {
				t.Fatal("private manifest published")
			}
		})
	}
}

func TestPublicationManifest_TimestampUncertainty(t *testing.T) {
	sum := sha256.Sum256([]byte("content"))
	committed := time.Now()
	identity := fileidentity.Identity{Inode: 1, Modified: committed.Add(-2 * time.Second).UnixNano(), Changed: committed.Add(-2 * time.Second).UnixNano()}
	record := publicationRecord{Source: identity, Digest: hex.EncodeToString(sum[:])}
	if !cachedSourceDigest(record, identity, committed, committed) {
		t.Fatal("old unchanged identity rejected")
	}
	if cachedSourceDigest(record, identity, committed, committed.Add(-time.Second)) {
		t.Fatal("backwards clock accepted")
	}
	identity.Changed = committed.UnixNano()
	record.Source = identity
	if cachedSourceDigest(record, identity, committed, committed) {
		t.Fatal("uncertain change time accepted")
	}
}
