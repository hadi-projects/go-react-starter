package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// LocalDriver stores files on the local filesystem.
// BasePath is the absolute or relative root directory, e.g. "./storage/uploads".
type LocalDriver struct {
	BasePath string
}

func NewLocalDriver(basePath string) *LocalDriver {
	os.MkdirAll(basePath, 0755)
	return &LocalDriver{BasePath: basePath}
}

// Save writes r to BasePath/key, creating parent directories as needed.
func (d *LocalDriver) Save(_ context.Context, key string, r io.Reader) error {
	fullPath := filepath.Join(d.BasePath, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

// Get opens a file for reading.
func (d *LocalDriver) Get(_ context.Context, key string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(d.BasePath, key))
}

// Delete removes the file at BasePath/key.
func (d *LocalDriver) Delete(_ context.Context, key string) error {
	err := os.Remove(filepath.Join(d.BasePath, key))
	if os.IsNotExist(err) {
		return nil // idempotent
	}
	return err
}

// Exists reports whether the file exists.
func (d *LocalDriver) Exists(_ context.Context, key string) bool {
	_, err := os.Stat(filepath.Join(d.BasePath, key))
	return err == nil
}
