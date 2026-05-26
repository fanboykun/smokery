// Package fs provides a filesystem-backed BlobStore. The CLI uses this to
// write artifacts (reports, JSON dumps) to a local directory.
package fs

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fanboykun/smokery/apps/core/internal/port"
)

// BlobStore stores blobs as files under a root directory. Keys are interpreted
// as relative paths beneath the root.
type BlobStore struct {
	root string
}

// New creates a BlobStore rooted at dir. The directory is created if missing.
func New(dir string) (*BlobStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &BlobStore{root: dir}, nil
}

var _ port.BlobStore = (*BlobStore)(nil)

func (s *BlobStore) path(key string) string {
	return filepath.Join(s.root, filepath.Clean("/"+key))
}

func (s *BlobStore) Put(ctx context.Context, key string, data io.Reader, contentType string) error {
	p := s.path(key)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, data)
	return err
}

func (s *BlobStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	return os.Open(s.path(key))
}

func (s *BlobStore) Delete(ctx context.Context, key string) error {
	err := os.Remove(s.path(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// SignedURL returns a file:// URL. The CLI doesn't need real signed URLs.
func (s *BlobStore) SignedURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	abs, err := filepath.Abs(s.path(key))
	if err != nil {
		return "", err
	}
	return "file://" + abs, nil
}
