package gitstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/uuid"
)

const DefaultBranch = "main"

// Store manages bare git repositories for workspaces on disk.
type Store struct {
	root string
}

func New(root string) (*Store, error) {
	if root == "" {
		root = "./data/git-repos"
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("gitstore mkdir: %w", err)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &Store{root: abs}, nil
}

func (s *Store) Root() string { return s.root }

func (s *Store) Path(workspaceGuid uuid.UUID) string {
	return filepath.Join(s.root, workspaceGuid.String()+".git")
}

// Exists reports whether the bare repo directory is present.
func (s *Store) Exists(workspaceGuid uuid.UUID) bool {
	st, err := os.Stat(s.Path(workspaceGuid))
	return err == nil && st.IsDir()
}

// Init creates a bare repo with an initial commit on main (seed workspace tree).
func (s *Store) Init(workspaceGuid uuid.UUID, name, version string) error {
	path := s.Path(workspaceGuid)
	if s.Exists(workspaceGuid) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	bare, err := git.PlainInit(path, true)
	if err != nil {
		return fmt.Errorf("plain init bare: %w", err)
	}

	if err := seedBareRepo(bare, name, version); err != nil {
		_ = os.RemoveAll(path)
		return err
	}
	return nil
}

// Delete removes the bare repository directory.
func (s *Store) Delete(workspaceGuid uuid.UUID) error {
	path := s.Path(workspaceGuid)
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove bare repo: %w", err)
	}
	return nil
}

// Ensure creates the bare repo if missing (for workspaces created before gitstore).
func (s *Store) Ensure(workspaceGuid uuid.UUID, name, version string) error {
	if s.Exists(workspaceGuid) {
		return nil
	}
	return s.Init(workspaceGuid, name, version)
}

type seedWorkspace struct {
	Name               string `json:"name"`
	Version            string `json:"version,omitempty"`
	ConnectionConfig   string `json:"connectionConfig,omitempty"`
	FolderNestingLimit *int   `json:"folderNestingLimit,omitempty"`
}

func seedBareRepo(bare *git.Repository, name, version string) error {
	storer := memory.NewStorage()
	fs := memfs.New()
	r, err := git.Init(storer, fs)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	meta, err := json.MarshalIndent(seedWorkspace{Name: name, Version: version}, "", "  ")
	if err != nil {
		return err
	}
	f, err := fs.Create("workspace.json")
	if err != nil {
		return err
	}
	if _, err := f.Write(append(meta, '\n')); err != nil {
		_ = f.Close()
		return err
	}
	_ = f.Close()

	if err := fs.MkdirAll("collections", 0o755); err != nil {
		return err
	}
	if err := fs.MkdirAll("envs", 0o755); err != nil {
		return err
	}
	// keep dirs via .gitkeep
	for _, p := range []string{"collections/.gitkeep", "envs/.gitkeep"} {
		kf, err := fs.Create(p)
		if err != nil {
			return err
		}
		_ = kf.Close()
	}

	if _, err := w.Add("."); err != nil {
		return err
	}
	hash, err := w.Commit("initial workspace", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "ladno-teams",
			Email: "teams@ladno.local",
			When:  time.Now().UTC(),
		},
	})
	if err != nil {
		return err
	}

	ref := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+DefaultBranch), hash)
	if err := bare.Storer.SetReference(ref); err != nil {
		return err
	}
	head := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName("refs/heads/"+DefaultBranch))
	if err := bare.Storer.SetReference(head); err != nil {
		return err
	}

	// Copy objects from memory storer into bare.
	return copyObjects(storer, bare)
}

func copyObjects(from *memory.Storage, to *git.Repository) error {
	iter, err := from.IterEncodedObjects(plumbing.AnyObject)
	if err != nil {
		return err
	}
	defer iter.Close()
	return iter.ForEach(func(obj plumbing.EncodedObject) error {
		_, err := to.Storer.SetEncodedObject(obj)
		return err
	})
}
