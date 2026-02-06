// Package noop provides a no-op implementation of the Storage provider.
package noop

import (
	"context"
	"io"
	"time"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/storage"
)

func init() {
	provider.Register[storage.StorageProvider]("storage", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op Storage Provider",
			Category:    "storage",
			Version:     "1.0.0",
			Description: "A no-operation storage provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op Storage provider.
type Provider struct{}

// NewProvider creates a new no-op Storage provider.
func NewProvider(config map[string]any) (storage.StorageProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) Upload(ctx context.Context, path string, reader io.Reader, opts storage.UploadOptions) (*storage.FileInfo, error) {
	return &storage.FileInfo{
		Path:         path,
		Size:         0,
		ContentType:  opts.ContentType,
		LastModified: time.Now(),
		Metadata:     opts.Metadata,
		ETag:         "noop-etag",
	}, nil
}

func (p *Provider) Download(ctx context.Context, path string) (io.ReadCloser, *storage.FileInfo, error) {
	return io.NopCloser(nil), &storage.FileInfo{
		Path:         path,
		Size:         0,
		ContentType:  "application/octet-stream",
		LastModified: time.Now(),
		Metadata:     make(map[string]string),
		ETag:         "noop-etag",
	}, nil
}

func (p *Provider) Delete(ctx context.Context, path string) error {
	return nil
}

func (p *Provider) Exists(ctx context.Context, path string) (bool, error) {
	return false, nil
}

func (p *Provider) GetSignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	return "https://example.com/noop/" + path, nil
}

func (p *Provider) List(ctx context.Context, prefix string, opts storage.ListOptions) ([]storage.FileInfo, error) {
	return []storage.FileInfo{}, nil
}

func (p *Provider) Metadata() storage.Metadata {
	return storage.Metadata{
		Name:   "noop",
		Region: "noop",
		Bucket: "noop",
	}
}
