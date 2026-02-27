// Package opensearch provides an OpenSearch implementation of the Search provider.
package opensearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/search"
)

func init() {
	provider.Register[search.SearchProvider]("search", "opensearch",
		provider.Metadata{
			Name:        "opensearch",
			DisplayName: "OpenSearch",
			Category:    "search",
			Version:     "1.0.0",
			Description: "OpenSearch search engine integration",
			ConfigSpec: []provider.ConfigField{
				{
					Key:         "addresses",
					Type:        "array",
					Required:    true,
					Description: "OpenSearch server addresses",
				},
				{
					Key:         "username",
					Type:        "string",
					Required:    false,
					Description: "OpenSearch username",
				},
				{
					Key:         "password",
					Type:        "secret",
					Required:    false,
					Description: "OpenSearch password",
				},
				{
					Key:         "insecure_skip_verify",
					Type:        "bool",
					Required:    false,
					Description: "Skip TLS certificate verification",
				},
			},
		},
		NewProvider,
	)
}

// germanDecompoundWords is a curated word list for B2B/industrial product search.
// The dictionary_decompounder filter uses this to split German compound words
// (e.g., "Sicherheitshandschuhe" → "Sicherheit" + "Handschuhe").
var germanDecompoundWords = []string{
	// Safety / PPE
	"sicherheit", "schutz", "handschuh", "brille", "helm", "schuh", "stiefel",
	"weste", "maske", "gehör", "atem", "schnitt", "hitze", "kälte",
	// Tools & Hardware
	"werkzeug", "schrauben", "schlüssel", "bohrer", "säge", "hammer", "zange",
	"messer", "klinge", "schneider", "dreher", "schleifer", "fräser",
	// Materials
	"stahl", "holz", "kunststoff", "gummi", "metall", "aluminium", "kupfer",
	"eisen", "blech", "rohr", "platte", "folie", "gewebe", "faser",
	// Parts & Components
	"teil", "stück", "satz", "set", "band", "kette", "ring", "dichtung",
	"lager", "ventil", "pumpe", "motor", "antrieb", "getriebe", "welle",
	// Surfaces & Coating
	"lack", "farbe", "beschichtung", "oberfläche", "versiegelung", "grund",
	// Fasteners
	"schraube", "mutter", "bolzen", "nagel", "dübel", "niete", "klemme",
	"klammer", "haken", "öse", "clip",
	// Dimensions & Properties
	"lang", "breit", "hoch", "dick", "dünn", "schwer", "leicht", "groß", "klein",
	"rund", "flach", "eckig",
	// Electrical
	"kabel", "stecker", "dose", "schalter", "leitung", "lampe", "licht",
	"batterie", "trafo",
	// Containers & Packaging
	"box", "kasten", "behälter", "dose", "flasche", "kanister", "eimer",
	"tonne", "sack", "beutel", "karton", "palette",
	// Generic B2B
	"maschine", "anlage", "gerät", "apparat", "system", "einheit", "modul",
	"bau", "montage", "wartung", "reparatur", "ersatz", "zusatz", "zubehör",
	"industrie", "profi", "standard", "spezial", "universal",
	// Clothing / Workwear
	"hose", "jacke", "hemd", "overall", "mantel", "latex", "nitril",
	// Measurement
	"mess", "prüf", "anzeige", "sensor", "regler", "zähler",
	// Cleaning
	"reiniger", "reinigung", "pflege", "mittel", "lösung",
}

// OpenSearchConfig holds the OpenSearch configuration
type OpenSearchConfig struct {
	Addresses          []string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

// Provider is an OpenSearch search provider.
type Provider struct {
	client *opensearchapi.Client
}

// NewProvider creates a new OpenSearch search provider.
func NewProvider(config map[string]any) (search.SearchProvider, error) {
	// Parse addresses
	var addresses []string
	if addrAny, ok := config["addresses"]; ok {
		switch v := addrAny.(type) {
		case []string:
			addresses = v
		case []any:
			for _, addr := range v {
				if s, ok := addr.(string); ok {
					addresses = append(addresses, s)
				}
			}
		case string:
			addresses = []string{v}
		}
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("opensearch: at least one address is required")
	}

	username, _ := config["username"].(string)
	password, _ := config["password"].(string)
	insecureSkipVerify, _ := config["insecure_skip_verify"].(bool)

	// Create OpenSearch client config
	cfg := opensearchapi.Config{
		Client: opensearch.Config{
			Addresses: addresses,
			Username:  username,
			Password:  password,
		},
	}

	if insecureSkipVerify {
		cfg.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	client, err := opensearchapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("opensearch: failed to create client: %w", err)
	}

	return &Provider{
		client: client,
	}, nil
}

func (p *Provider) IndexDocuments(ctx context.Context, index string, documents []search.Document) (*search.TaskResult, error) {
	// OpenSearch bulk indexing
	var buf bytes.Buffer

	for _, doc := range documents {
		// Action line
		action := map[string]any{
			"index": map[string]any{
				"_index": index,
				"_id":    doc["id"],
			},
		}
		if err := json.NewEncoder(&buf).Encode(action); err != nil {
			return nil, fmt.Errorf("opensearch: failed to encode action: %w", err)
		}

		// Document line
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return nil, fmt.Errorf("opensearch: failed to encode document: %w", err)
		}
	}

	req := opensearchapi.BulkReq{
		Body: bytes.NewReader(buf.Bytes()),
	}

	resp, err := p.client.Bulk(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("opensearch: bulk indexing failed: %w", err)
	}

	if resp.Errors {
		return nil, fmt.Errorf("opensearch: bulk indexing returned errors")
	}

	// OpenSearch bulk is synchronous, so we return success immediately
	return &search.TaskResult{
		TaskID: "bulk-index",
		Status: "succeeded",
	}, nil
}

func (p *Provider) DeleteDocuments(ctx context.Context, index string, ids []string) (*search.TaskResult, error) {
	var buf bytes.Buffer

	for _, id := range ids {
		// Action line
		action := map[string]any{
			"delete": map[string]any{
				"_index": index,
				"_id":    id,
			},
		}
		if err := json.NewEncoder(&buf).Encode(action); err != nil {
			return nil, fmt.Errorf("opensearch: failed to encode delete action: %w", err)
		}
	}

	req := opensearchapi.BulkReq{
		Body: bytes.NewReader(buf.Bytes()),
	}

	resp, err := p.client.Bulk(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("opensearch: delete failed: %w", err)
	}

	if resp.Errors {
		return nil, fmt.Errorf("opensearch: delete returned errors")
	}

	return &search.TaskResult{
		TaskID: "bulk-delete",
		Status: "succeeded",
	}, nil
}

func (p *Provider) ConfigureIndex(ctx context.Context, index string, config search.IndexConfig) error {
	// Create index with mappings and settings
	mappings := map[string]any{
		"properties": map[string]any{
			"name_de": map[string]any{
				"type":     "text",
				"analyzer": "german",
				"fields": map[string]any{
					"prefix": map[string]any{"type": "text", "analyzer": "autocomplete", "search_analyzer": "autocomplete_search"},
				},
			},
			"name_en": map[string]any{
				"type":     "text",
				"analyzer": "english",
				"fields": map[string]any{
					"prefix": map[string]any{"type": "text", "analyzer": "autocomplete", "search_analyzer": "autocomplete_search"},
				},
			},
			"name_fr": map[string]any{
				"type":     "text",
				"analyzer": "french",
				"fields": map[string]any{
					"prefix": map[string]any{"type": "text", "analyzer": "autocomplete", "search_analyzer": "autocomplete_search"},
				},
			},
			"name_it": map[string]any{
				"type":     "text",
				"analyzer": "italian",
				"fields": map[string]any{
					"prefix": map[string]any{"type": "text", "analyzer": "autocomplete", "search_analyzer": "autocomplete_search"},
				},
			},
			"description_de": map[string]any{"type": "text", "analyzer": "german"},
			"description_en": map[string]any{"type": "text", "analyzer": "english"},
			"sku": map[string]any{
				"type": "keyword",
				"fields": map[string]any{
					"search": map[string]any{"type": "text"},
				},
			},
			"product_type": map[string]any{"type": "keyword"},
			"status":       map[string]any{"type": "keyword"},
			"tenant_id":    map[string]any{"type": "keyword"},
			"category_ids": map[string]any{"type": "keyword"},
			"created_at":   map[string]any{"type": "date", "format": "epoch_second"},
			"updated_at":   map[string]any{"type": "date", "format": "epoch_second"},
		},
	}

	settings := map[string]any{
		"analysis": map[string]any{
			"filter": map[string]any{
				"german_decompounder": map[string]any{
					"type":             "dictionary_decompounder",
					"word_list":        germanDecompoundWords,
					"min_word_size":    5,
					"min_subword_size": 3,
					"max_subword_size": 15,
				},
				"german_stemmer": map[string]any{
					"type":     "stemmer",
					"language": "light_german",
				},
				"german_normalization": map[string]any{
					"type": "german_normalization",
				},
				"french_stemmer": map[string]any{
					"type":     "stemmer",
					"language": "light_french",
				},
				"french_elision": map[string]any{
					"type":     "elision",
					"articles": []string{"l", "m", "t", "qu", "n", "s", "j", "d", "c"},
				},
				"italian_stemmer": map[string]any{
					"type":     "stemmer",
					"language": "light_italian",
				},
				"italian_elision": map[string]any{
					"type":     "elision",
					"articles": []string{"c", "l", "all", "dall", "dell", "nell", "sull", "coll", "pell", "gl", "agl", "dagl", "degl", "negl", "sugl", "un", "m", "t", "s", "v", "d"},
				},
				"english_stemmer": map[string]any{
					"type":     "stemmer",
					"language": "light_english",
				},
				"autocomplete_filter": map[string]any{
					"type":     "edge_ngram",
					"min_gram": 2,
					"max_gram": 20,
				},
			},
			"analyzer": map[string]any{
				"autocomplete": map[string]any{
					"type":      "custom",
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "autocomplete_filter"},
				},
				"autocomplete_search": map[string]any{
					"type":      "custom",
					"tokenizer": "standard",
					"filter":    []string{"lowercase"},
				},
				"german": map[string]any{
					"type":      "custom",
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "german_decompounder", "german_normalization", "german_stemmer"},
				},
				"french": map[string]any{
					"type":      "custom",
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "french_elision", "french_stemmer"},
				},
				"italian": map[string]any{
					"type":      "custom",
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "italian_elision", "italian_stemmer"},
				},
			},
		},
	}

	body := map[string]any{
		"settings": settings,
		"mappings": mappings,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("opensearch: failed to marshal index config: %w", err)
	}

	req := opensearchapi.IndicesCreateReq{
		Index: index,
		Body:  bytes.NewReader(bodyBytes),
	}

	_, err = p.client.Indices.Create(ctx, req)
	if err != nil {
		// Ignore "already exists" errors
		if !strings.Contains(err.Error(), "resource_already_exists") {
			return fmt.Errorf("opensearch: failed to configure index: %w", err)
		}
	}

	return nil
}

func (p *Provider) Search(ctx context.Context, index string, query search.SearchQuery) (*search.SearchResult, error) {
	// Build multi_match query across all language fields
	var queryObj map[string]any

	if query.Query != "" {
		// Combined query strategy:
		// 1. Prefix match on autocomplete subfields (highest boost — for suggestions/typeahead)
		// 2. Analyzed match on language fields with decompounder (for full-word search)
		// 3. SKU search
		// No fuzziness to avoid false positives (e.g., "Leder" → "LED")
		queryObj = map[string]any{
			"bool": map[string]any{
				"should": []map[string]any{
					// Prefix matching (autocomplete) — highest priority
					{"multi_match": map[string]any{
						"query":  query.Query,
						"fields": []string{"name_de.prefix^10", "name_en.prefix^5", "name_fr.prefix^3", "name_it.prefix^3"},
						"type":   "best_fields",
					}},
					// Analyzed match (decompounder, stemming) — for complete words, with typo tolerance
					{"multi_match": map[string]any{
						"query":     query.Query,
						"fields":    []string{"name_de^3", "name_en^2", "name_fr", "name_it", "description_de", "description_en"},
						"type":      "best_fields",
						"fuzziness": "AUTO:4,7",
					}},
					// SKU match
					{"match": map[string]any{
						"sku.search": map[string]any{"query": query.Query, "boost": 2},
					}},
				},
				"minimum_should_match": 1,
			},
		}
	} else {
		// Match all if no query
		queryObj = map[string]any{
			"match_all": map[string]any{},
		}
	}

	// Build filters
	var filterQueries []map[string]any
	for _, f := range query.Filters {
		filterQueries = append(filterQueries, buildOpenSearchFilter(f))
	}

	// Combine query and filters
	boolQuery := map[string]any{
		"must": queryObj,
	}
	if len(filterQueries) > 0 {
		boolQuery["filter"] = filterQueries
	}

	searchBody := map[string]any{
		"query": map[string]any{
			"bool": boolQuery,
		},
		"from": query.Offset,
		"size": query.Limit,
	}

	// Add sort if specified
	if len(query.Sort) > 0 {
		searchBody["sort"] = query.Sort
	}

	// Add aggregations for faceted search
	searchBody["aggs"] = map[string]any{
		"categories": map[string]any{
			"terms": map[string]any{
				"field": "category_ids",
				"size":  50,
			},
		},
		"product_types": map[string]any{
			"terms": map[string]any{
				"field": "product_type",
				"size":  10,
			},
		},
	}

	bodyBytes, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("opensearch: failed to marshal search query: %w", err)
	}

	req := opensearchapi.SearchReq{
		Indices: []string{index},
		Body:    bytes.NewReader(bodyBytes),
	}

	resp, err := p.client.Search(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("opensearch: search failed: %w", err)
	}

	// Convert hits to documents
	hits := make([]search.Document, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var doc map[string]any
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			continue
		}
		hits = append(hits, search.Document(doc))
	}

	// Parse aggregations into facets
	facets := make(map[string]map[string]int)
	if len(resp.Aggregations) > 0 {
		var aggs map[string]json.RawMessage
		if err := json.Unmarshal(resp.Aggregations, &aggs); err == nil {
			for aggName, rawAgg := range aggs {
				var agg struct {
					Buckets []struct {
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				}
				if err := json.Unmarshal(rawAgg, &agg); err == nil && len(agg.Buckets) > 0 {
					facets[aggName] = make(map[string]int)
					for _, bucket := range agg.Buckets {
						facets[aggName][bucket.Key] = bucket.DocCount
					}
				}
			}
		}
	}

	return &search.SearchResult{
		Hits:             hits,
		TotalHits:        int(resp.Hits.Total.Value),
		ProcessingTimeMs: resp.Took,
		Facets:           facets,
	}, nil
}

func (p *Provider) CreateIndex(ctx context.Context, index string, primaryKey string) error {
	// Simple index creation (full config happens in ConfigureIndex)
	req := opensearchapi.IndicesCreateReq{
		Index: index,
	}

	_, err := p.client.Indices.Create(ctx, req)
	if err != nil {
		// Ignore "already exists" errors
		if !strings.Contains(err.Error(), "resource_already_exists") {
			return fmt.Errorf("opensearch: failed to create index: %w", err)
		}
	}

	return nil
}

func (p *Provider) DeleteIndex(ctx context.Context, index string) error {
	req := opensearchapi.IndicesDeleteReq{
		Indices: []string{index},
	}

	_, err := p.client.Indices.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("opensearch: failed to delete index: %w", err)
	}

	return nil
}

func (p *Provider) GetTaskStatus(ctx context.Context, taskID string) (*search.TaskResult, error) {
	// OpenSearch operations are mostly synchronous, so tasks are completed immediately
	return &search.TaskResult{
		TaskID: taskID,
		Status: "succeeded",
	}, nil
}

func (p *Provider) Health(ctx context.Context) error {
	resp, err := p.client.Cluster.Health(ctx, nil)
	if err != nil {
		return fmt.Errorf("opensearch: health check failed: %w", err)
	}

	if resp.Status != "green" && resp.Status != "yellow" {
		return fmt.Errorf("opensearch: cluster status is %s", resp.Status)
	}

	return nil
}

func (p *Provider) Metadata() search.Metadata {
	return search.Metadata{
		Name:    "opensearch",
		Version: "2.18.0",
		Features: []string{
			"multilingual",
			"fuzzy-search",
			"facets",
			"filtering",
			"sorting",
		},
	}
}

// buildOpenSearchFilter converts search.Filter to OpenSearch query DSL
func buildOpenSearchFilter(f search.Filter) map[string]any {
	switch f.Operator {
	case "=":
		// Use "terms" for slice values (e.g., category hierarchy), "term" for single values
		switch v := f.Value.(type) {
		case []string:
			return map[string]any{
				"terms": map[string]any{
					f.Field: v,
				},
			}
		default:
			return map[string]any{
				"term": map[string]any{
					f.Field: f.Value,
				},
			}
		}
	case "!=":
		return map[string]any{
			"bool": map[string]any{
				"must_not": map[string]any{
					"term": map[string]any{
						f.Field: f.Value,
					},
				},
			},
		}
	case ">":
		return map[string]any{
			"range": map[string]any{
				f.Field: map[string]any{
					"gt": f.Value,
				},
			},
		}
	case "<":
		return map[string]any{
			"range": map[string]any{
				f.Field: map[string]any{
					"lt": f.Value,
				},
			},
		}
	case ">=":
		return map[string]any{
			"range": map[string]any{
				f.Field: map[string]any{
					"gte": f.Value,
				},
			},
		}
	case "<=":
		return map[string]any{
			"range": map[string]any{
				f.Field: map[string]any{
					"lte": f.Value,
				},
			},
		}
	case "IN":
		return map[string]any{
			"terms": map[string]any{
				f.Field: f.Value,
			},
		}
	case "NOT IN":
		return map[string]any{
			"bool": map[string]any{
				"must_not": map[string]any{
					"terms": map[string]any{
						f.Field: f.Value,
					},
				},
			},
		}
	default:
		return map[string]any{
			"match_all": map[string]any{},
		}
	}
}
