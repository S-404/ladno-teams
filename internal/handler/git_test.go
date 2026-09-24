package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ladno-teams/internal/entity/dto"
	"ladno-teams/internal/exception"
	"ladno-teams/internal/gitstore"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/google/uuid"
)

type fakeGitWorkspace struct{}

func (fakeGitWorkspace) AuthorizeGitRead(userGuid, workspaceGuid uuid.UUID) *exception.ApiError {
	return nil
}
func (fakeGitWorkspace) AuthorizeGitWrite(userGuid, workspaceGuid uuid.UUID) *exception.ApiError {
	return nil
}
func (fakeGitWorkspace) EnsureGitRepo(workspaceGuid uuid.UUID) error { return nil }

func TestInfoRefsAdvertise(t *testing.T) {
	dir := t.TempDir()
	store, err := gitstore.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	guid := uuid.New()
	if err := store.Init(guid, "Demo", "v1"); err != nil {
		t.Fatal(err)
	}

	body, err := advertiseRefs(store.Path(guid), transport.UploadPackServiceName)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) < 20 {
		t.Fatalf("short body: %q", body)
	}
	if !containsService(body, "git-upload-pack") {
		t.Fatalf("missing service line: %q", body[:min(80, len(body))])
	}

	gin.SetMode(gin.TestMode)
	h := NewGitHandler(fakeGitWorkspace{}, store)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxUser, dto.UserDto{Guid: uuid.New(), Login: "t"})
	})
	r.GET("/git/workspaces/:guid/info/refs", h.InfoRefs)
	r.POST("/git/workspaces/:guid/git-upload-pack", h.UploadPack)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/git/workspaces/"+guid.String()+"/info/refs?service=git-upload-pack", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/x-git-upload-pack-advertisement" {
		t.Fatalf("content-type=%q", ct)
	}
}

func TestGitSubcommand(t *testing.T) {
	if got := gitSubcommand(transport.UploadPackServiceName); got != "upload-pack" {
		t.Fatalf("got %q", got)
	}
	if got := gitSubcommand(transport.ReceivePackServiceName); got != "receive-pack" {
		t.Fatalf("got %q", got)
	}
}

func TestCloneRoundTripViaHTTP(t *testing.T) {
	dir := t.TempDir()
	store, err := gitstore.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	guid := uuid.New()
	if err := store.Init(guid, "Demo", "v1"); err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	h := NewGitHandler(fakeGitWorkspace{}, store)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxUser, dto.UserDto{Guid: uuid.New(), Login: "t"})
	})
	r.GET("/git/workspaces/:guid/info/refs", h.InfoRefs)
	r.POST("/git/workspaces/:guid/git-upload-pack", h.UploadPack)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	cloneDir := t.TempDir() + "/clone"
	remote := srv.URL + "/git/workspaces/" + guid.String()
	_, err = git.PlainClone(cloneDir, false, &git.CloneOptions{URL: remote})
	if err != nil {
		t.Fatalf("clone: %v", err)
	}
	repo, err := git.PlainOpen(cloneDir)
	if err != nil {
		t.Fatal(err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	if head.Hash().IsZero() {
		t.Fatal("empty head")
	}
}

func containsService(b []byte, service string) bool {
	return bytesContains(b, []byte("# service="+service))
}

func bytesContains(b, sub []byte) bool {
	for i := 0; i+len(sub) <= len(b); i++ {
		ok := true
		for j := range sub {
			if b[i+j] != sub[j] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
