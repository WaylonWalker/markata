package buildcache

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/encryption"
)

func TestSourceEncryptionCache_SaveLoadAndValidate(t *testing.T) {
	cacheDir := t.TempDir()
	cache, err := LoadSourceEncryptionCache(cacheDir)
	if err != nil {
		t.Fatalf("load cache: %v", err)
	}
	password := "h7Qm!2Vx9#Lp4@Td" //nolint:gosec // test password
	body := "secret body\n"
	encrypted, err := encryption.EncryptSourceMarkdown(body, "default", password)
	if err != nil {
		t.Fatalf("encrypt body: %v", err)
	}
	cache.Put("post.md", "default", encrypted)
	if err := cache.Save(); err != nil {
		t.Fatalf("save cache: %v", err)
	}

	reloaded, err := LoadSourceEncryptionCache(cacheDir)
	if err != nil {
		t.Fatalf("reload persisted cache: %v", err)
	}
	if got, ok := reloaded.Get("post.md", body, "default", password); !ok || got != encrypted {
		t.Fatal("persisted cache did not validate unchanged ciphertext")
	}
	if _, ok := reloaded.Get("post.md", "changed body\n", "default", password); ok {
		t.Fatal("cache accepted changed plaintext")
	}
	if _, ok := reloaded.Get("post.md", body, "default", "wrong-password"); ok {
		t.Fatal("cache accepted the wrong password")
	}

	reloaded.Entries["post.md"] = SourceEncryptionEntry{KeyName: "default", EncryptedBody: "tampered"}
	if _, ok := reloaded.Get("post.md", body, "default", password); ok {
		t.Fatal("cache accepted tampered ciphertext")
	}
}
