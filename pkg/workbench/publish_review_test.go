package workbench

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
	"github.com/WaylonWalker/markata/pkg/models"
)

type publishWorkspaceFixture struct{ root string }

func (f publishWorkspaceFixture) Config() (*models.Config, error) { return nil, nil }
func (f publishWorkspaceFixture) Paths() ([]editorworkspace.File, error) { return nil, nil }
func (f publishWorkspaceFixture) Files() ([]editorworkspace.File, error) { return nil, nil }
func (f publishWorkspaceFixture) Create(string, string) error { return nil }
func (f publishWorkspaceFixture) Write(string, string) error { return nil }
func (f publishWorkspaceFixture) Read(path string) (string, error) {
	body, err := os.ReadFile(filepath.Join(f.root, filepath.FromSlash(path)))
	return string(body), err
}

func gitForReviewTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %s: %v", args, output, err)
	}
}

func TestPublishReviewUsesExactCommittedSourceNotWorkingIndex(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil { t.Skip("git not installed") }
	dir := t.TempDir()
	gitForReviewTest(t,dir,"init","-q")
	gitForReviewTest(t,dir,"config","user.name","Review Test")
	gitForReviewTest(t,dir,"config","user.email","review@example.com")
	if err := os.MkdirAll(filepath.Join(dir,"posts"),0755); err != nil { t.Fatal(err) }
	path := filepath.Join(dir,"posts","example.md")
	before := "---\npublished: false\n---\n\n# Original\n"
	after := "---\npublished: true\n---\n\n# Changed\n"
	if err := os.WriteFile(path,[]byte(before),0644); err != nil { t.Fatal(err) }
	gitForReviewTest(t,dir,"add","posts/example.md")
	gitForReviewTest(t,dir,"commit","-qm","baseline")
	if err := os.WriteFile(path,[]byte(after),0644); err != nil { t.Fatal(err) }
	handler := &Handler{root:dir,workspace:publishWorkspaceFixture{root:dir}}
	request := httptest.NewRequest(http.MethodGet,"/workspace/publish/review?path=posts%2Fexample.md",nil)
	recorder := httptest.NewRecorder()
	handler.servePublishReview(recorder,request)
	if recorder.Code != http.StatusOK { t.Fatalf("status %d: %s",recorder.Code,recorder.Body.String()) }
	var got publishReview
	if err := json.Unmarshal(recorder.Body.Bytes(),&got); err != nil {t.Fatal(err)}
	if !got.BaselineExists || !got.Changed || got.Baseline!=before || got.Current!=after {
		t.Fatalf("wrong review snapshot: %+v",got)
	}
	if got.CanPublish || got.BaselineKind!="last-committed" {
		t.Fatalf("must not advertise a deployment: %+v",got)
	}
	if got.Head=="" || strings.Contains(got.Head,"\n") { t.Fatalf("bad HEAD: %q",got.Head) }
}

func TestPublishReviewRejectsUnsafePaths(t *testing.T) {
	dir:=t.TempDir()
	handler:=&Handler{root:dir,workspace:publishWorkspaceFixture{root:dir}}
	for _,path:=range []string{"../escape.md","/tmp/elsewhere.md","a.txt"} {
		recorder:=httptest.NewRecorder()
		handler.servePublishReview(recorder,httptest.NewRequest(http.MethodGet,
			"/workspace/publish/review?path="+path,nil))
		if recorder.Code!=http.StatusBadRequest {
			t.Errorf("%q returned %d, expected 400",path,recorder.Code)
		}
	}
}

func TestPublishReviewUntrackedDocument(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil { t.Skip("git not installed") }
	dir := t.TempDir()
	gitForReviewTest(t,dir,"init","-q")
	gitForReviewTest(t,dir,"config","user.name","Review Test")
	gitForReviewTest(t,dir,"config","user.email","review@example.com")
	if err:=os.WriteFile(filepath.Join(dir,"README.md"),[]byte("base\n"),0644); err!=nil{t.Fatal(err)}
	gitForReviewTest(t,dir,"add","README.md")
	gitForReviewTest(t,dir,"commit","-qm","baseline")
	if err:=os.WriteFile(filepath.Join(dir,"new.md"),[]byte("# New\n"),0644); err!=nil{t.Fatal(err)}
	handler:=&Handler{root:dir,workspace:publishWorkspaceFixture{root:dir}}
	recorder:=httptest.NewRecorder()
	handler.servePublishReview(recorder,httptest.NewRequest(http.MethodGet,"/workspace/publish/review?path=new.md",nil))
	if recorder.Code!=http.StatusOK {t.Fatalf("status %d: %s",recorder.Code,recorder.Body.String())}
	var got publishReview
	if err:=json.Unmarshal(recorder.Body.Bytes(),&got);err!=nil{t.Fatal(err)}
	if got.BaselineExists || !got.Changed || got.Baseline!="" || got.Current!="# New\n" {
		t.Fatalf("incorrect new-document review: %+v",got)
	}
}
