package builderadmin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
)

type publicationRecord struct {
	Source  fileidentity.Identity
	Release fileidentity.Identity
	Digest  string
}

type publicationManifest struct {
	ReleaseID string
	Workspace string
	Committed time.Time
	Files     map[string]publicationRecord
}

type publicationEnvelope struct {
	Version  int
	Payload  json.RawMessage
	Checksum string
}

func loadPublicationManifest(workspace, baseline string) publicationManifest {
	empty := publicationManifest{}
	file, err := os.Open(workspace + ".publication.json")
	if err != nil {
		return empty
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, (32<<20)+1))
	if err != nil || len(data) > 32<<20 {
		return empty
	}
	var envelope publicationEnvelope
	if json.Unmarshal(data, &envelope) != nil || envelope.Version != 1 {
		return empty
	}
	digest := sha256.Sum256(envelope.Payload)
	if hex.EncodeToString(digest[:]) != envelope.Checksum {
		return empty
	}
	var manifest publicationManifest
	absolute, err := filepath.Abs(workspace)
	if err != nil || json.Unmarshal(envelope.Payload, &manifest) != nil || manifest.Workspace != absolute || manifest.ReleaseID != baseline || manifest.Committed.IsZero() {
		return empty
	}
	return manifest
}

func savePublicationManifest(workspace, release string, files map[string]publicationRecord) bool {
	absolute, err := filepath.Abs(workspace)
	if err != nil {
		return false
	}
	payload, err := json.Marshal(publicationManifest{ReleaseID: release, Workspace: absolute, Committed: time.Now(), Files: files})
	if err != nil {
		return false
	}
	digest := sha256.Sum256(payload)
	data, err := json.Marshal(publicationEnvelope{Version: 1, Payload: payload, Checksum: hex.EncodeToString(digest[:])})
	if err != nil || len(data) > 32<<20 {
		return false
	}
	file, err := os.CreateTemp(filepath.Dir(workspace), ".publication-*")
	if err != nil {
		return false
	}
	defer os.Remove(file.Name())
	_, writeErr := file.Write(data)
	syncErr := file.Sync()
	closeErr := file.Close()
	if writeErr != nil || syncErr != nil || closeErr != nil {
		return false
	}
	return os.Rename(file.Name(), workspace+".publication.json") == nil
}

func cachedSourceDigest(record publicationRecord, identity fileidentity.Identity, committed, now time.Time) bool {
	return record.Source == identity && !now.Before(committed) && identity.StableBefore(committed) && validPublicationDigest(record.Digest)
}

func validPublicationDigest(digest string) bool {
	decoded, err := hex.DecodeString(digest)
	return err == nil && len(decoded) == sha256.Size
}

func hashPublicationFile(file *os.File, buffer []byte) (string, int64, error) {
	hash := sha256.New()
	count, err := io.CopyBuffer(hash, file, buffer)
	return hex.EncodeToString(hash.Sum(nil)), count, err
}

func (p *workspaceDeltaPublisher) sourceDigest(path string, info fs.FileInfo, record publicationRecord) (string, error) {
	identity := fileidentity.Read(info)
	if cachedSourceDigest(record, identity, p.previous.Committed, time.Now()) {
		p.stats.CachedSourceHashes++
		return record.Digest, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	digest, count, err := hashPublicationFile(file, p.sourceBuffer)
	p.stats.ComparedBytes += count
	if err != nil {
		return "", err
	}
	after, err := file.Stat()
	if err != nil {
		return "", err
	}
	if identity.Inode != 0 && fileidentity.Read(after) != identity {
		return "", fmt.Errorf("workspace file changed during publication: %s", path)
	}
	return digest, nil
}

func (p *workspaceDeltaPublisher) baselineDigest(old string, source fs.FileInfo, record publicationRecord) string {
	info, err := p.root.Lstat(old)
	if err != nil || !info.Mode().IsRegular() || info.Size() != source.Size() || info.Mode().Perm() != source.Mode().Perm() || os.SameFile(info, source) {
		return ""
	}
	identity := fileidentity.Read(info)
	identity.Changed = 0
	if identity.Inode != 0 && identity == record.Release && validPublicationDigest(record.Digest) {
		p.stats.CachedReleaseHashes++
		return record.Digest
	}
	file, err := p.root.Open(old)
	if err != nil {
		return ""
	}
	defer file.Close()
	digest, count, err := hashPublicationFile(file, p.baselineBuffer)
	p.stats.ComparedBytes += count
	if err != nil {
		return ""
	}
	return digest
}

func (p *workspaceDeltaPublisher) copyRegular(path, target, next, rel string, info fs.FileInfo) error {
	record := p.previous.Files[rel]
	digest, err := p.sourceDigest(path, info, record)
	if err != nil {
		return err
	}
	old := filepath.Join(p.baselineID, rel)
	linked := p.baselineID != "" && digest == p.baselineDigest(old, info, record) && p.link(p.root, old, next) == nil
	if linked {
		p.stats.LinkedFiles++
	} else {
		if err := copyWorkspaceFile(path, target, info.Mode().Perm(), info.ModTime()); err != nil {
			return err
		}
		p.stats.CopiedFiles++
		p.stats.CopiedBytes += info.Size()
	}
	final, err := p.root.Lstat(next)
	if err != nil {
		return err
	}
	identity := fileidentity.Read(final)
	identity.Changed = 0
	p.records[rel] = publicationRecord{Source: fileidentity.Read(info), Release: identity, Digest: digest}
	return nil
}
