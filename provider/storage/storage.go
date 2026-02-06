// Package storage defines the interface for object storage integrations.
package storage

import (
	"context"
	"io"
	"time"
)

// StorageProvider abstracts object storage (S3, Azure Blob, MinIO, local filesystem).
type StorageProvider interface {
	// Upload uploads a file.
	Upload(ctx context.Context, path string, reader io.Reader, opts UploadOptions) (*FileInfo, error)

	// Download downloads a file.
	Download(ctx context.Context, path string) (io.ReadCloser, *FileInfo, error)

	// Delete deletes a file.
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists.
	Exists(ctx context.Context, path string) (bool, error)

	// GetSignedURL returns a temporary URL for direct access.
	GetSignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)

	// List lists files in a path.
	List(ctx context.Context, prefix string, opts ListOptions) ([]FileInfo, error)

	// Metadata returns provider information.
	Metadata() Metadata
}

// UploadOptions contains options for uploading files.
type UploadOptions struct {
	ContentType string
	Metadata    map[string]string
	ACL         string // "private", "public-read"
}

// ListOptions contains options for listing files.
type ListOptions struct {
	MaxKeys int
	Cursor  string
}

// FileInfo represents information about a stored file.
type FileInfo struct {
	Path         string
	Size         int64
	ContentType  string
	LastModified time.Time
	Metadata     map[string]string
	ETag         string
}

// Metadata provides information about the storage provider.
type Metadata struct {
	Name   string
	Region string
	Bucket string
}
