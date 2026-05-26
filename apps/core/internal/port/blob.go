package port

import (
	"context"
	"io"
	"time"
)

// BlobStore is an abstract object storage interface.
type BlobStore interface {
	Put(ctx context.Context, key string, data io.Reader, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	SignedURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}
