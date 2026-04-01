package storage

import (
	"context"
	"io"
)

// Driver is the abstraction for file storage backends.
// The default implementation is LocalDriver (filesystem).
// Can be swapped for S3/MinIO without changing service code.
type Driver interface {
	// Save stores data from r using key as the file path identifier.
	Save(ctx context.Context, key string, r io.Reader) error
	// Get opens a file for reading. Caller must Close() the result.
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete removes the file permanently.
	Delete(ctx context.Context, key string) error
	// Exists reports whether a file with the given key exists.
	Exists(ctx context.Context, key string) bool
}
