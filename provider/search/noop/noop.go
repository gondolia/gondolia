// Package noop provides a no-op implementation of the Search provider.
package noop

import (
	"context"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/search"
)

func init() {
	provider.Register[search.SearchProvider]("search", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op Search Provider",
			Category:    "search",
			Version:     "1.0.0",
			Description: "A no-operation search provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op Search provider.
type Provider struct{}

// NewProvider creates a new no-op Search provider.
func NewProvider(config map[string]any) (search.SearchProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) IndexDocuments(ctx context.Context, index string, documents []search.Document) (*search.TaskResult, error) {
	return &search.TaskResult{
		TaskID: "noop-task-001",
		Status: "succeeded",
		Error:  "",
	}, nil
}

func (p *Provider) DeleteDocuments(ctx context.Context, index string, ids []string) (*search.TaskResult, error) {
	return &search.TaskResult{
		TaskID: "noop-task-002",
		Status: "succeeded",
		Error:  "",
	}, nil
}

func (p *Provider) ConfigureIndex(ctx context.Context, index string, config search.IndexConfig) error {
	return nil
}

func (p *Provider) Search(ctx context.Context, index string, query search.SearchQuery) (*search.SearchResult, error) {
	return &search.SearchResult{
		Hits:             []search.Document{},
		TotalHits:        0,
		Facets:           make(map[string]map[string]int),
		ProcessingTimeMs: 0,
	}, nil
}

func (p *Provider) CreateIndex(ctx context.Context, index string, primaryKey string) error {
	return nil
}

func (p *Provider) DeleteIndex(ctx context.Context, index string) error {
	return nil
}

func (p *Provider) GetTaskStatus(ctx context.Context, taskID string) (*search.TaskResult, error) {
	return &search.TaskResult{
		TaskID: taskID,
		Status: "succeeded",
		Error:  "",
	}, nil
}

func (p *Provider) Health(ctx context.Context) error {
	return nil
}

func (p *Provider) Metadata() search.Metadata {
	return search.Metadata{
		Name:     "noop",
		Version:  "1.0.0",
		Features: []string{},
	}
}
