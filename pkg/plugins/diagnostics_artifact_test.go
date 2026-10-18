package plugins

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestDiagnosticsArtifactPlugin_Priority(t *testing.T) {
	plugin := NewDiagnosticsArtifactPlugin()
	if got := plugin.Priority(lifecycle.StageCleanup); got != lifecycle.PriorityLast+100 {
		t.Fatalf("Priority(StageCleanup) = %d, want %d", got, lifecycle.PriorityLast+100)
	}
	if got := plugin.Priority(lifecycle.StageWrite); got != lifecycle.PriorityDefault {
		t.Fatalf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestDiagnosticsArtifactPlugin_WritesAfterCompleteLifecycle(t *testing.T) {
	outputDir := t.TempDir()
	contentDir := t.TempDir()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
		Extra: map[string]interface{}{
			"markata_version": "test-version",
			"markata_commit":  "test-commit",
		},
	})
	m.RegisterPlugin(NewDiagnosticsArtifactPlugin())
	m.SetFiles([]string{"post.md"})
	m.ContentLedger().MarkLoaded("post.md")

	if err := m.Run(); err != nil {
		t.Fatalf("Manager.Run() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, diagnostics.DefaultArtifactPath))
	if err != nil {
		t.Fatalf("read diagnostics artifact: %v", err)
	}
	artifact, err := diagnostics.ParseArtifact(data)
	if err != nil {
		t.Fatalf("ParseArtifact() error = %v", err)
	}
	if artifact.Generator.Version != "test-version" || artifact.Generator.Commit != "test-commit" {
		t.Fatalf("generator = %#v", artifact.Generator)
	}
	if len(artifact.Entries) != 1 || artifact.Entries[0].Path != "post.md" {
		t.Fatalf("entries = %#v", artifact.Entries)
	}
}

func TestDiagnosticsArtifactPlugin_SkipsIncompleteLifecycle(t *testing.T) {
	outputDir := t.TempDir()
	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir

	if err := m.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("RunTo(StageWrite) error = %v", err)
	}
	if err := NewDiagnosticsArtifactPlugin().Cleanup(m); err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, diagnostics.DefaultArtifactPath)); !os.IsNotExist(err) {
		t.Fatalf("incomplete lifecycle created an artifact, stat error = %v", err)
	}
}

func TestDiagnosticsArtifactPlugin_PublishesBuildWithWarnings(t *testing.T) {
	outputDir := t.TempDir()
	destination := filepath.Join(outputDir, diagnostics.DefaultArtifactPath)

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	m.RegisterPlugins(&diagnosticsArtifactRenderFailure{}, NewDiagnosticsArtifactPlugin())
	if err := m.Run(); err != nil {
		t.Fatalf("Manager.Run() error = %v, want warning-only completion", err)
	}
	data, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read diagnostics artifact: %v", err)
	}
	if _, err := diagnostics.ParseArtifact(data); err != nil {
		t.Fatalf("ParseArtifact() error = %v", err)
	}
	if len(m.Warnings()) != 1 {
		t.Fatalf("warnings = %d, want 1", len(m.Warnings()))
	}
}

func TestDiagnosticsArtifactPlugin_PreservesArtifactAfterFailedBuild(t *testing.T) {
	outputDir := t.TempDir()
	destination := filepath.Join(outputDir, diagnostics.DefaultArtifactPath)
	if err := writeDiagnosticsArtifact(destination, []byte("previous artifact\n")); err != nil {
		t.Fatalf("write previous artifact: %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	m.RegisterPlugins(&diagnosticsArtifactCriticalRenderFailure{}, NewDiagnosticsArtifactPlugin())
	if err := m.Run(); err == nil {
		t.Fatal("Manager.Run() succeeded after critical render failure")
	}
	data, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read previous artifact: %v", err)
	}
	if string(data) != "previous artifact\n" {
		t.Fatalf("artifact changed after failed build: %q", data)
	}
}

func TestDiagnosticsArtifactPlugin_SkipsNonProductionBuildModes(t *testing.T) {
	for _, mode := range []string{"fast_mode", "incremental_mode"} {
		t.Run(mode, func(t *testing.T) {
			outputDir := t.TempDir()
			m := lifecycle.NewManager()
			m.Config().OutputDir = outputDir
			m.Config().Extra[mode] = true
			m.RegisterPlugin(NewDiagnosticsArtifactPlugin())

			if err := m.Run(); err != nil {
				t.Fatalf("Manager.Run() error = %v", err)
			}
			if _, err := os.Stat(filepath.Join(outputDir, diagnostics.DefaultArtifactPath)); !os.IsNotExist(err) {
				t.Fatalf("%s build created an artifact, stat error = %v", mode, err)
			}
		})
	}
}

func TestDiagnosticsArtifactPlugin_WriteFailureIsCritical(t *testing.T) {
	outputDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(outputDir, ".markata"), []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("create conflicting output path: %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	m.RegisterPlugin(NewDiagnosticsArtifactPlugin())
	err := m.Run()
	if err == nil {
		t.Fatal("Manager.Run() succeeded after artifact write failure")
	}
	var hookErrors *lifecycle.HookErrors
	if !errors.As(err, &hookErrors) || !hookErrors.HasCritical() {
		t.Fatalf("error = %v, want critical hook error", err)
	}
	if m.HasRun(lifecycle.StageCleanup) {
		t.Fatal("cleanup stage was marked complete after critical artifact failure")
	}
}

func TestWriteDiagnosticsArtifact_ReplacesExistingFile(t *testing.T) {
	destination := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
	if err := writeDiagnosticsArtifact(destination, []byte("old\n")); err != nil {
		t.Fatalf("write old artifact: %v", err)
	}
	if err := writeDiagnosticsArtifact(destination, []byte("new\n")); err != nil {
		t.Fatalf("write new artifact: %v", err)
	}
	data, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if string(data) != "new\n" {
		t.Fatalf("artifact = %q, want new content", data)
	}
}

func TestWriteDiagnosticsArtifact_RemovesTemporaryFileAfterReplacementFailure(t *testing.T) {
	directory := t.TempDir()
	destination := filepath.Join(directory, diagnostics.DefaultArtifactPath)
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatalf("create conflicting artifact directory: %v", err)
	}

	if err := writeDiagnosticsArtifact(destination, []byte("artifact")); err == nil {
		t.Fatal("writeDiagnosticsArtifact() succeeded with a directory destination")
	}
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(destination), ".diagnostics-*.tmp"))
	if err != nil {
		t.Fatalf("find temporary artifacts: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary artifacts remain after failed replacement: %v", matches)
	}
}

func TestWriteDiagnosticsArtifactWith_PreservesPreviousOnFailure(t *testing.T) {
	sentinel := errors.New("injected writer failure")
	snapshot := diagnostics.ContentLedgerSnapshot{
		Entries: []diagnostics.ContentDisposition{{Path: "post.md"}},
	}
	for _, test := range []struct {
		name  string
		write func(io.Writer) error
		want  error
	}{
		{
			name: "serialization",
			write: func(writer io.Writer) error {
				return diagnostics.WriteArtifact(writer, snapshot, diagnostics.ArtifactBuildInfo{
					BuiltAt: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
				})
			},
		},
		{
			name: "partial serialization",
			write: func(writer io.Writer) error {
				if _, err := io.WriteString(writer, "partial JSON"); err != nil {
					return err
				}
				return sentinel
			},
			want: sentinel,
		},
		{
			name: "stream writer error",
			write: func(writer io.Writer) error {
				return diagnostics.WriteArtifact(&diagnosticsArtifactFailWriter{writer: writer, err: sentinel},
					snapshot, diagnostics.ArtifactBuildInfo{BuiltAt: time.Unix(1, 0)})
			},
			want: sentinel,
		},
		{
			name: "stream short write",
			write: func(writer io.Writer) error {
				return diagnostics.WriteArtifact(&diagnosticsArtifactFailWriter{writer: writer},
					snapshot, diagnostics.ArtifactBuildInfo{BuiltAt: time.Unix(1, 0)})
			},
			want: io.ErrShortWrite,
		},
		{
			name: "buffered flush short write",
			write: func(writer io.Writer) error {
				buffered, ok := writer.(*bufio.Writer)
				if !ok {
					return errors.New("atomic artifact writer must provide a buffered writer")
				}
				buffered.Reset(&diagnosticsArtifactFailWriter{failAt: 1})
				_, err := buffered.WriteString("buffered JSON")
				return err
			},
			want: io.ErrShortWrite,
		},
		{
			name: "buffered flush failure",
			write: func(writer io.Writer) error {
				// Inject a failing sink beneath the buffer: serialization succeeds
				// and only the atomic helper's final flush reports the failure.
				buffered, ok := writer.(*bufio.Writer)
				if !ok {
					return errors.New("atomic artifact writer must provide a buffered writer")
				}
				buffered.Reset(&diagnosticsArtifactFailWriter{err: sentinel, failAt: 1})
				_, err := buffered.WriteString("buffered JSON")
				return err
			},
			want: sentinel,
		},
	} {
		for _, existing := range []bool{false, true} {
			t.Run(test.name+fmtArtifactExisting(existing), func(t *testing.T) {
				destination := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
				if existing {
					if err := writeDiagnosticsArtifact(destination, []byte("previous artifact")); err != nil {
						t.Fatal(err)
					}
				}
				err := writeDiagnosticsArtifactWith(destination, test.write)
				if err == nil || (test.want != nil && !errors.Is(err, test.want)) {
					t.Fatalf("error = %v, want failure %v", err, test.want)
				}
				data, readErr := os.ReadFile(destination)
				if existing {
					if readErr != nil || string(data) != "previous artifact" {
						t.Fatalf("previous artifact changed: %q, %v", data, readErr)
					}
				} else if !os.IsNotExist(readErr) {
					t.Fatalf("failed publication created destination: %v", readErr)
				}
				matches, err := filepath.Glob(filepath.Join(filepath.Dir(destination), ".diagnostics-*.tmp"))
				if err != nil || len(matches) != 0 {
					t.Fatalf("temporary artifacts remain: %v, %v", matches, err)
				}
			})
		}
	}
}

func fmtArtifactExisting(existing bool) string {
	if existing {
		return "/existing"
	}
	return "/first"
}

// Accept the header through the real temporary writer, then fail an entry write.
type diagnosticsArtifactFailWriter struct {
	writer io.Writer
	err    error
	calls  int
	failAt int
}

func (w *diagnosticsArtifactFailWriter) Write(data []byte) (int, error) {
	w.calls++
	failAt := w.failAt
	if failAt == 0 {
		failAt = 2
	}
	if w.calls == failAt {
		if w.err != nil {
			return 0, w.err
		}
		return len(data) - 1, nil
	}
	return w.writer.Write(data)
}

func TestWriteDiagnosticsArtifactWith_StreamingByteParity(t *testing.T) {
	snapshot := diagnostics.ContentLedgerSnapshot{
		Entries: []diagnostics.ContentDisposition{
			{Path: "z.md", Feeds: []diagnostics.ContentFeedDisposition{
				{Feed: "z", Included: false, Reasons: []string{diagnostics.ReasonContentFiltered}},
				{Feed: "a", Included: true},
			}},
			{Path: "./a.md"},
		},
	}
	info := diagnostics.ArtifactBuildInfo{BuiltAt: time.Unix(1, 0), Executor: diagnostics.ArtifactExecutorLegacy}
	want, err := diagnostics.MarshalArtifact(snapshot, info)
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
	if err := writeDiagnosticsArtifactWith(destination, func(writer io.Writer) error {
		return diagnostics.WriteArtifact(writer, snapshot, info)
	}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(destination)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("published bytes differ: error = %v\ngot:\n%s\nwant:\n%s", err, got, want)
	}
}

func TestWriteDiagnosticsArtifactMeasured_Phases(t *testing.T) {
	for _, size := range []int{128, 32 * 1024, 128 * 1024} {
		t.Run(fmtArtifactSize(size), func(t *testing.T) {
			data := bytes.Repeat([]byte("x"), size)
			destination := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
			var phases diagnosticsArtifactPhases
			err := writeDiagnosticsArtifactMeasured(destination, func(writer io.Writer) error {
				buffered, ok := writer.(*bufio.Writer)
				if !ok || buffered.Size() != diagnosticsArtifactBufferSize {
					t.Fatal("publication must use the fixed bounded buffer")
				}
				_, err := writer.Write(data)
				return err
			}, &phases)
			if err != nil {
				t.Fatal(err)
			}
			if phases.published != 1 || phases.fileWriteBytes != int64(size) || phases.fileWriteCalls < 1 {
				t.Fatalf("publication counts = %#v", phases)
			}
			if phases.serialization < 0 || phases.flush < 0 ||
				phases.serializationWall+phases.flushWall != phases.serialization+phases.flush+phases.fileWrite {
				t.Fatalf("file writes counted twice or omitted: %#v", phases)
			}
			// Fast operations can complete within a single clock tick on Windows.
			if phases.fileWrite < 0 || phases.sync < 0 || phases.close < 0 || phases.replace < 0 {
				t.Fatalf("negative operation duration: %#v", phases)
			}
			if size < diagnosticsArtifactBufferSize && phases.serialization != phases.serializationWall {
				t.Fatal("small buffered write counted as underlying file write")
			}
			if size > diagnosticsArtifactBufferSize && phases.fileWrite > 0 && phases.serialization == phases.serializationWall {
				t.Fatal("large direct file write not separated from serialization")
			}
			got, err := os.ReadFile(destination)
			if err != nil || !bytes.Equal(got, data) {
				t.Fatal("measured publication changed bytes")
			}
			info, err := os.Stat(destination)
			if err != nil {
				t.Fatal(err)
			}
			// Windows only represents the owner write bit in Chmod permissions.
			if runtime.GOOS == "windows" {
				if info.Mode().Perm()&0o200 == 0 {
					t.Fatal("measured publication is not writable")
				}
			} else if info.Mode().Perm() != 0o644 {
				t.Fatal("measured publication changed permissions")
			}
		})
	}
}

func fmtArtifactSize(size int) string {
	if size < diagnosticsArtifactBufferSize {
		return "buffered"
	}
	return "direct"
}

type diagnosticsArtifactPartialSink struct {
	err error
}

func (w diagnosticsArtifactPartialSink) Write(data []byte) (int, error) {
	return len(data) / 2, w.err
}

func TestDiagnosticsArtifactMeasuredWriter_PartialWrites(t *testing.T) {
	sentinel := errors.New("underlying write failed")
	for _, failure := range []error{nil, sentinel} {
		var phases diagnosticsArtifactPhases
		writer := diagnosticsArtifactMeasuredWriter{
			writer: diagnosticsArtifactPartialSink{err: failure}, phases: &phases,
		}
		written, err := writer.Write([]byte("123456"))
		want := failure
		if want == nil {
			want = io.ErrShortWrite
		}
		if written != 3 || !errors.Is(err, want) || phases.fileWriteCalls != 1 ||
			phases.fileWriteBytes != 3 || phases.fileWrite < 0 {
			t.Fatalf("partial write: %d, %v, %#v", written, err, phases)
		}
	}
}

func TestWriteDiagnosticsArtifactMeasured_FailurePhases(t *testing.T) {
	destination := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
	sentinel := errors.New("serialization failed")
	var phases diagnosticsArtifactPhases
	err := writeDiagnosticsArtifactMeasured(destination, func(writer io.Writer) error {
		if _, err := writer.Write(bytes.Repeat([]byte("x"), 2*diagnosticsArtifactBufferSize)); err != nil {
			return err
		}
		return sentinel
	}, &phases)
	if !errors.Is(err, sentinel) || phases.published != 0 || phases.close < 0 ||
		phases.sync != 0 || phases.replace != 0 || phases.flushWall != 0 ||
		phases.fileWriteBytes != 2*diagnosticsArtifactBufferSize {
		t.Fatalf("failure phases = %#v, error %v", phases, err)
	}
}

type diagnosticsArtifactRenderFailure struct{}

func (p *diagnosticsArtifactRenderFailure) Name() string {
	return "diagnostics_artifact_test_failure"
}

func (p *diagnosticsArtifactRenderFailure) Render(_ *lifecycle.Manager) error {
	return errors.New("expected render warning")
}

var _ lifecycle.RenderPlugin = (*diagnosticsArtifactRenderFailure)(nil)

type diagnosticsArtifactCriticalRenderFailure struct{}

func (p *diagnosticsArtifactCriticalRenderFailure) Name() string {
	return "diagnostics_artifact_test_critical_failure"
}

func (p *diagnosticsArtifactCriticalRenderFailure) Render(_ *lifecycle.Manager) error {
	return diagnosticsArtifactCriticalError{}
}

type diagnosticsArtifactCriticalError struct{}

func (diagnosticsArtifactCriticalError) Error() string {
	return "expected critical render failure"
}

func (diagnosticsArtifactCriticalError) IsCritical() bool {
	return true
}

var _ lifecycle.RenderPlugin = (*diagnosticsArtifactCriticalRenderFailure)(nil)
var _ lifecycle.CriticalError = diagnosticsArtifactCriticalError{}
