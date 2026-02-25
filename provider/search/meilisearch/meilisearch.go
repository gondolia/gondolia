// Package meilisearch provides a Meilisearch implementation of the Search provider.
package meilisearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/meilisearch/meilisearch-go"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/search"
)

func init() {
	provider.Register[search.SearchProvider]("search", "meilisearch",
		provider.Metadata{
			Name:        "meilisearch",
			DisplayName: "Meilisearch",
			Category:    "search",
			Version:     "1.0.0",
			Description: "Meilisearch search engine integration",
			ConfigSpec: []provider.ConfigField{
				{
					Key:         "host",
					Type:        "string",
					Required:    true,
					Description: "Meilisearch server host URL",
				},
				{
					Key:         "api_key",
					Type:        "secret",
					Required:    true,
					Description: "Meilisearch master/admin API key",
				},
			},
		},
		NewProvider,
	)
}

// Provider is a Meilisearch search provider.
type Provider struct {
	client meilisearch.ServiceManager
}

// NewProvider creates a new Meilisearch search provider.
func NewProvider(config map[string]any) (search.SearchProvider, error) {
	host, ok := config["host"].(string)
	if !ok || host == "" {
		return nil, fmt.Errorf("meilisearch: host is required")
	}

	apiKey, ok := config["api_key"].(string)
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("meilisearch: api_key is required")
	}

	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	return &Provider{
		client: client,
	}, nil
}

func (p *Provider) IndexDocuments(ctx context.Context, index string, documents []search.Document) (*search.TaskResult, error) {
	idx := p.client.Index(index)

	// Convert documents to []map[string]any
	docs := make([]map[string]any, len(documents))
	for i, doc := range documents {
		docs[i] = map[string]any(doc)
	}

	taskInfo, err := idx.AddDocumentsWithContext(ctx, docs, nil)
	if err != nil {
		return nil, fmt.Errorf("meilisearch: failed to index documents: %w", err)
	}

	return &search.TaskResult{
		TaskID: fmt.Sprintf("%d", taskInfo.TaskUID),
		Status: "enqueued",
	}, nil
}

func (p *Provider) DeleteDocuments(ctx context.Context, index string, ids []string) (*search.TaskResult, error) {
	idx := p.client.Index(index)

	taskInfo, err := idx.DeleteDocumentsWithContext(ctx, ids, nil)
	if err != nil {
		return nil, fmt.Errorf("meilisearch: failed to delete documents: %w", err)
	}

	return &search.TaskResult{
		TaskID: fmt.Sprintf("%d", taskInfo.TaskUID),
		Status: "enqueued",
	}, nil
}

func (p *Provider) ConfigureIndex(ctx context.Context, index string, config search.IndexConfig) error {
	idx := p.client.Index(index)

	// Configure searchable attributes
	if len(config.SearchableAttributes) > 0 {
		if _, err := idx.UpdateSearchableAttributesWithContext(ctx, &config.SearchableAttributes); err != nil {
			return fmt.Errorf("meilisearch: failed to update searchable attributes: %w", err)
		}
	}

	// Configure filterable attributes
	if len(config.FilterableAttributes) > 0 {
		filterableAttrs := make([]any, len(config.FilterableAttributes))
		for i, attr := range config.FilterableAttributes {
			filterableAttrs[i] = attr
		}
		if _, err := idx.UpdateFilterableAttributesWithContext(ctx, &filterableAttrs); err != nil {
			return fmt.Errorf("meilisearch: failed to update filterable attributes: %w", err)
		}
	}

	// Configure sortable attributes
	if len(config.SortableAttributes) > 0 {
		if _, err := idx.UpdateSortableAttributesWithContext(ctx, &config.SortableAttributes); err != nil {
			return fmt.Errorf("meilisearch: failed to update sortable attributes: %w", err)
		}
	}

	// Configure synonyms
	if len(config.Synonyms) > 0 {
		if _, err := idx.UpdateSynonymsWithContext(ctx, &config.Synonyms); err != nil {
			return fmt.Errorf("meilisearch: failed to update synonyms: %w", err)
		}
	}

	// Configure stop words
	if len(config.StopWords) > 0 {
		if _, err := idx.UpdateStopWordsWithContext(ctx, &config.StopWords); err != nil {
			return fmt.Errorf("meilisearch: failed to update stop words: %w", err)
		}
	}

	// Configure typo tolerance
	if config.TypoTolerance != nil {
		typoTolerance := &meilisearch.TypoTolerance{
			Enabled: config.TypoTolerance.Enabled,
		}
		if config.TypoTolerance.MinWordSizeForTypos != nil {
			oneTypo := int64(config.TypoTolerance.MinWordSizeForTypos["oneTypo"])
			twoTypos := int64(config.TypoTolerance.MinWordSizeForTypos["twoTypos"])
			typoTolerance.MinWordSizeForTypos = meilisearch.MinWordSizeForTypos{
				OneTypo:  oneTypo,
				TwoTypos: twoTypos,
			}
		}
		if _, err := idx.UpdateTypoToleranceWithContext(ctx, typoTolerance); err != nil {
			return fmt.Errorf("meilisearch: failed to update typo tolerance: %w", err)
		}
	}

	return nil
}

func (p *Provider) Search(ctx context.Context, index string, query search.SearchQuery) (*search.SearchResult, error) {
	idx := p.client.Index(index)

	searchRequest := &meilisearch.SearchRequest{
		Query:  query.Query,
		Offset: int64(query.Offset),
		Limit:  int64(query.Limit),
	}

	// Build filter string from filters
	if len(query.Filters) > 0 {
		filterStr := buildMeilisearchFilter(query.Filters)
		searchRequest.Filter = filterStr
	}

	// Add facets
	if len(query.Facets) > 0 {
		searchRequest.Facets = query.Facets
	}

	// Add sort
	if len(query.Sort) > 0 {
		searchRequest.Sort = query.Sort
	}

	// Add highlight
	if len(query.Highlight) > 0 {
		searchRequest.AttributesToHighlight = query.Highlight
	}

	searchResp, err := idx.SearchWithContext(ctx, query.Query, searchRequest)
	if err != nil {
		return nil, fmt.Errorf("meilisearch: search failed: %w", err)
	}

	// Convert hits to documents
	hits := make([]search.Document, 0, len(searchResp.Hits))
	for _, hit := range searchResp.Hits {
		doc := make(map[string]any)
		for key, rawValue := range hit {
			var value any
			if err := json.Unmarshal(rawValue, &value); err == nil {
				doc[key] = value
			}
		}
		hits = append(hits, search.Document(doc))
	}

	// Convert facet distribution
	facets := make(map[string]map[string]int)
	if len(searchResp.FacetDistribution) > 0 {
		var facetDist map[string]map[string]int
		if err := json.Unmarshal(searchResp.FacetDistribution, &facetDist); err == nil {
			facets = facetDist
		}
	}

	return &search.SearchResult{
		Hits:             hits,
		TotalHits:        int(searchResp.EstimatedTotalHits),
		Facets:           facets,
		ProcessingTimeMs: int(searchResp.ProcessingTimeMs),
	}, nil
}

func (p *Provider) CreateIndex(ctx context.Context, index string, primaryKey string) error {
	task, err := p.client.CreateIndexWithContext(ctx, &meilisearch.IndexConfig{
		Uid:        index,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		return fmt.Errorf("meilisearch: failed to create index: %w", err)
	}

	// Wait for task to complete (use 0 for default timeout)
	_, err = p.client.WaitForTaskWithContext(ctx, task.TaskUID, 0)
	if err != nil {
		return fmt.Errorf("meilisearch: failed to wait for index creation: %w", err)
	}

	return nil
}

func (p *Provider) DeleteIndex(ctx context.Context, index string) error {
	task, err := p.client.DeleteIndexWithContext(ctx, index)
	if err != nil {
		return fmt.Errorf("meilisearch: failed to delete index: %w", err)
	}

	// Wait for task to complete (use 0 for default timeout)
	_, err = p.client.WaitForTaskWithContext(ctx, task.TaskUID, 0)
	if err != nil {
		return fmt.Errorf("meilisearch: failed to wait for index deletion: %w", err)
	}

	return nil
}

func (p *Provider) GetTaskStatus(ctx context.Context, taskID string) (*search.TaskResult, error) {
	var uid int64
	if _, err := fmt.Sscanf(taskID, "%d", &uid); err != nil {
		return nil, fmt.Errorf("meilisearch: invalid task ID: %w", err)
	}

	task, err := p.client.GetTaskWithContext(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("meilisearch: failed to get task status: %w", err)
	}

	status := strings.ToLower(string(task.Status))
	result := &search.TaskResult{
		TaskID: taskID,
		Status: status,
	}

	if task.Error.Message != "" {
		result.Error = task.Error.Message
	}

	return result, nil
}

func (p *Provider) Health(ctx context.Context) error {
	health, err := p.client.HealthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("meilisearch: health check failed: %w", err)
	}
	if health.Status != "available" {
		return fmt.Errorf("meilisearch: server is not available (status: %s)", health.Status)
	}
	return nil
}

func (p *Provider) Metadata() search.Metadata {
	return search.Metadata{
		Name:    "meilisearch",
		Version: "1.12",
		Features: []string{
			"facets",
			"typo-tolerance",
			"synonyms",
			"stop-words",
			"highlighting",
			"filtering",
			"sorting",
		},
	}
}

// buildMeilisearchFilter converts search.Filter to Meilisearch filter syntax
func buildMeilisearchFilter(filters []search.Filter) string {
	if len(filters) == 0 {
		return ""
	}

	var parts []string
	for _, f := range filters {
		filterStr := buildSingleFilter(f)
		if filterStr != "" {
			parts = append(parts, filterStr)
		}
	}

	// Combine filters with AND
	return strings.Join(parts, " AND ")
}

func buildSingleFilter(f search.Filter) string {
	switch f.Operator {
	case "=":
		return fmt.Sprintf("%s = %v", f.Field, formatFilterValue(f.Value))
	case "!=":
		return fmt.Sprintf("%s != %v", f.Field, formatFilterValue(f.Value))
	case ">":
		return fmt.Sprintf("%s > %v", f.Field, f.Value)
	case "<":
		return fmt.Sprintf("%s < %v", f.Field, f.Value)
	case ">=":
		return fmt.Sprintf("%s >= %v", f.Field, f.Value)
	case "<=":
		return fmt.Sprintf("%s <= %v", f.Field, f.Value)
	case "IN":
		if values, ok := f.Value.([]any); ok {
			var formatted []string
			for _, v := range values {
				formatted = append(formatted, fmt.Sprintf("%v", formatFilterValue(v)))
			}
			return fmt.Sprintf("%s IN [%s]", f.Field, strings.Join(formatted, ", "))
		}
	case "NOT IN":
		if values, ok := f.Value.([]any); ok {
			var formatted []string
			for _, v := range values {
				formatted = append(formatted, fmt.Sprintf("%v", formatFilterValue(v)))
			}
			return fmt.Sprintf("%s NOT IN [%s]", f.Field, strings.Join(formatted, ", "))
		}
	}
	return ""
}

func formatFilterValue(value any) string {
	switch v := value.(type) {
	case string:
		// Escape quotes and wrap in quotes
		escaped := strings.ReplaceAll(v, `"`, `\"`)
		return fmt.Sprintf(`"%s"`, escaped)
	default:
		return fmt.Sprintf("%v", v)
	}
}
