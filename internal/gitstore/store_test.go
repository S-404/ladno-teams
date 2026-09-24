package gitstore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
)

func TestInitAndOpen(t *testing.T) {
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	guid := uuid.New()
	if err := s.Init(guid, "Demo", "v1"); err != nil {
		t.Fatal(err)
	}
	path := s.Path(guid)
	if _, err := os.Stat(filepath.Join(path, "HEAD")); err != nil {
		t.Fatalf("bare HEAD missing: %v", err)
	}
	repo, err := git.PlainOpen(path)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	if ref.Name().Short() != DefaultBranch {
		t.Fatalf("branch=%s", ref.Name().Short())
	}
	// idempotent
	if err := s.Init(guid, "Demo", "v1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(guid); err != nil {
		t.Fatal(err)
	}
	if s.Exists(guid) {
		t.Fatal("expected deleted")
	}
}
