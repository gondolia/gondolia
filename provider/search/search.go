// Package search defines the interface for search engine integrations.
package search

import "context"

// SearchProvider abstracts search engines (Meilisearch, Algolia, Elasticsearch, etc.).
type SearchProvider interface {
	// --- Indexing ---

	// IndexDocuments indexes documents into an index.
	IndexDocuments(ctx context.Context, index string, documents []Document) (*TaskResult, error)

	// DeleteDocuments deletes documents from an index.
	DeleteDocuments(ctx context.Context, index string, ids []string) (*TaskResult, error)

	// ConfigureIndex configures index settings.
	ConfigureIndex(ctx context.Context, index string, config IndexConfig) error

	// --- Search ---

	// Search executes a search query.
	Search(ctx context.Context, index string, query SearchQuery) (*SearchResult, error)

	// --- Management ---

	// CreateIndex creates a new index.
	CreateIndex(ctx context.Context, index string, primaryKey string) error

	// DeleteIndex deletes an index.
	DeleteIndex(ctx context.Context, index string) error

	// GetTaskStatus checks the status of an asynchronous task.
	GetTaskStatus(ctx context.Context, taskID string) (*TaskResult, error)

	// Health checks the availability of the search engine.
	Health(ctx context.Context) error

	// Metadata returns provider information.
	Metadata() Metadata
}

// Document represents a document to be indexed.
type Document map[string]any

// SearchQuery represents a search request.
type SearchQuery struct {
	Query     string
	Filters   []Filter
	Facets    []string
	Sort      []string
	Offset    int
	Limit     int
	Highlight []string
}

// Filter represents a search filter.
type Filter struct {
	Field    string
	Operator string // "=", "!=", ">", "<", ">=", "<=", "IN", "NOT IN"
	Value    any
}

// SearchResult represents search results.
type SearchResult struct {
	Hits             []Document
	TotalHits        int
	Facets           map[string]map[string]int
	ProcessingTimeMs int
}

// IndexConfig represents index configuration.
type IndexConfig struct {
	SearchableAttributes []string
	FilterableAttributes []string
	SortableAttributes   []string
	Synonyms             map[string][]string
	StopWords            []string
	TypoTolerance        *TypoTolerance
}

// TypoTolerance configures typo tolerance settings.
type TypoTolerance struct {
	Enabled             bool
	MinWordSizeForTypos map[string]int
}

// TaskResult represents the status of an asynchronous operation.
type TaskResult struct {
	TaskID string
	Status string // "enqueued", "processing", "succeeded", "failed"
	Error  string
}

// Metadata provides information about the search provider.
type Metadata struct {
	Name     string
	Version  string
	Features []string // e.g. ["facets", "typo-tolerance", "synonyms", "geo-search"]
}
