package editorworkspace

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspace_ReadImageConfinesStaticFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "static", "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	want := []byte("image bytes")
	if err := os.WriteFile(filepath.Join(root, "static", "images", "photo.jpg"), want, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "private.jpg"), want, 0o600); err != nil {
		t.Fatal(err)
	}
	workspace, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := workspace.ReadImage("images/photo.jpg")
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("ReadImage() = %q, %v", got, err)
	}
	for _, path := range []string{"../private.jpg", "images/photo.md", "/images/photo.jpg"} {
		if _, err := workspace.ReadImage(path); err == nil {
			t.Errorf("ReadImage(%q) accepted an unsafe path", path)
		}
	}
}

func TestWorkspace_ConfiguredDiscoveryAndSource(t *testing.T) {
	root := t.TempDir()
	write := func(path, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, path), []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("markata-go.toml", "[markata-go.glob]\npatterns = ['*.md']\nuse_gitignore = true\n")
	write(".gitignore", "ignored.md\n")
	source := "---\ntitle: From loader\ntags: [test]\n---\n# Original\n"
	write("post.md", source)
	write("ignored.md", "# Hidden")
	write("secret.txt", "private")
	workspace, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	files, err := workspace.Files()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Path != "post.md" || files[0].Title == nil || *files[0].Title != "From loader" {
		t.Fatalf("unexpected files: %+v", files)
	}
	got, err := workspace.Read("post.md")
	if err != nil || got != source {
		t.Fatalf("source = %q, error = %v", got, err)
	}
	for _, path := range []string{"../outside.md", "ignored.md", "secret.txt", filepath.Join(root, "post.md")} {
		if _, err := workspace.Read(path); err == nil {
			t.Fatalf("read allowed %q", path)
		}
	}
}

func TestWorkspace_RejectsSymlinkEscape(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "markata-go.toml"), []byte("[markata-go.glob]\npatterns = ['*.md']\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(outside, "secret.md")
	if err := os.WriteFile(target, []byte("# Outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, "escape.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	workspace, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	files, err := workspace.Files()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("escaped file listed: %+v", files)
	}
	if _, err := workspace.Read("escape.md"); err == nil {
		t.Fatal("escaped file readable")
	}
}

func TestWorkspace_PathsDoesNotParseBodies(t *testing.T) {
	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_EDITORPERFORMANCEFIXTUREMISSING", "")
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "markata-go.toml"), []byte("[markata-go.glob]\npatterns=['*.md']\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "broken.md"), []byte("<!-- markata-encrypted-source:v1 key=editorperformancefixturemissing -->\nYWJj\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	workspace, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	paths, err := workspace.Paths()
	if err != nil || len(paths) != 1 || paths[0].Path != "broken.md" {
		t.Fatalf("paths=%v error=%v", paths, err)
	}
	if paths[0].Title != nil || len(paths[0].Tags) != 0 {
		t.Fatal("path discovery returned parsed metadata")
	}
	if _, err := workspace.Files(); err == nil {
		t.Fatal("fixture must fail metadata parsing")
	}
}
