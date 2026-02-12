# Search Engine Evaluation: Algolia vs Alternativen

## Ausgangssituation

**Aktuell (V2):** Algolia (SaaS)
**Frage:** Können wir auf eine self-hosted Lösung wechseln?

---

## Kundenfeedback: Die Suche funktioniert!

> **Wichtig:** Die aktuelle Suche wird von Kunden **positiv hervorgehoben**. Produkte lassen sich sehr gut finden.

### Was Kunden schätzen

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  POSITIVES KUNDENFEEDBACK ZUR SUCHE                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ✅ "Produkte lassen sich sehr gut finden"                                  │
│  ✅ Kategorieübergreifende Suche (wichtig für B2B!)                         │
│  ✅ Schnelle Ergebnisse                                                      │
│  ✅ Relevante Treffer                                                        │
│  ✅ Typo-Toleranz funktioniert                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### B2B-Suchverhalten

B2B-Kunden suchen **anders** als B2C:

| Aspekt | B2C | B2B |
|--------|-----|-----|
| **Suchstrategie** | Browsen in Kategorien | Direkte Suche nach Artikelnummer/Name |
| **Kategorie-Bindung** | Meist innerhalb einer Kategorie | **Kategorieübergreifend** |
| **Suchanfragen** | "Laminat Eiche" | "Swiss Krono D4152" (exakte SKU) |
| **Erwartung** | Inspiration | Schnelles Finden |

### Design-Prinzip: 1:1 Übernahme

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  ⚠️  WICHTIG: KEIN REDESIGN DER SUCHE!                                      │
│  ════════════════════════════════════                                        │
│                                                                              │
│  Die Suche funktioniert. Kunden sind zufrieden.                             │
│  → Das Design und UX soll 1:1 übernommen werden.                            │
│                                                                              │
│  Was NICHT ändern:                                                          │
│  • Such-UI (Autocomplete, Dropdown, Layout)                                 │
│  • Ergebnis-Darstellung (Kacheln, Liste)                                    │
│  • Filter-Sidebar (Facetten, Preis-Slider)                                  │
│  • Sortier-Optionen                                                          │
│  • Kategorieübergreifende Suche                                             │
│                                                                              │
│  Was ÄNDERN (Backend):                                                       │
│  • Engine: Algolia → Meilisearch                                            │
│  • Hosting: SaaS → Self-Hosted                                              │
│  • Kosten: $800+/Monat → ~$50/Monat                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Zu erhaltende Features

| Feature | Beschreibung | Meilisearch |
|---------|--------------|-------------|
| **Instant Search** | Ergebnisse beim Tippen | ✅ <50ms |
| **Kategorieübergreifend** | Suche über alle Produkte | ✅ Multi-Index Query |
| **Facetten** | Filter nach Attributen | ✅ Faceted Search |
| **Typo-Toleranz** | "Lamiant" → "Laminat" | ✅ Auto |
| **Highlighting** | Suchbegriff hervorheben | ✅ Eingebaut |
| **Autocomplete** | Vorschläge beim Tippen | ✅ Prefix Search |
| **SKU-Suche** | Exakte Artikelnummer | ✅ Exact Match Boost |
| **Synonyme** | Support-gepflegte Liste | ✅ API + Admin UI |

### Kundenspezifische Suche (Self-Hosted Vorteil)

Mit eigenem Suchindex können wir **kundeneigene Daten** suchbar machen - bei Algolia wäre das datenschutzrechtlich problematisch und teuer.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  KUNDENSPEZIFISCHE SUCHE                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  "Zeig mir meine letzten Bestellungen mit Laminat"                         │
│  "Was habe ich letztes Jahr für Projekt Müller bestellt?"                  │
│  "Meine häufigsten Produkte"                                                │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  SUCHBARE KUNDENDATEN:                                                      │
│                                                                              │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐ │
│  │  Bestellhistorie    │  │  Eigene Preise      │  │  Favoriten/Listen   │ │
│  │                     │  │                     │  │                     │ │
│  │  • Bestellnummer    │  │  • Kundenpreise     │  │  • Merkliste        │ │
│  │  • Bestelldatum     │  │  • Staffelpreise    │  │  • Projektlisten    │ │
│  │  • Produkte         │  │  • Rahmenverträge   │  │  • "Oft bestellt"   │ │
│  │  • Projektreferenz  │  │                     │  │                     │ │
│  │  • Lieferadresse    │  │                     │  │                     │ │
│  └─────────────────────┘  └─────────────────────┘  └─────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Use Cases:**

| Feature | Beschreibung | Nutzen |
|---------|--------------|--------|
| **Bestellsuche** | "Bestellung vom März" / "Lieferschein 12345" | Schneller Self-Service |
| **Nachbestellung** | Produkte aus alter Bestellung anzeigen | 1-Click Reorder |
| **Projektsuche** | "Was habe ich für Baustelle X bestellt?" | B2B-Workflow |
| **Personalisierte Ergebnisse** | Bereits gekaufte Produkte höher ranken | Bessere Relevanz |
| **Preissuche** | "Zeig mir Produkte unter meinem Rahmenvertrag" | B2B-Komfort |

**Index-Struktur:**

```
Indices:
  {tenant}_products_{language}           # Produktkatalog (alle Kunden)
  {tenant}_customer_{customer_id}        # Kundenspezifischer Index
```

**Kundenspezifisches Dokument:**

```json
{
  "id": "order-12345-line-1",
  "type": "order_item",
  "customer_id": "cust-789",

  // Bestelldaten
  "order_id": "12345",
  "order_date": "2024-11-15",
  "order_reference": "Projekt Müller Umbau",

  // Produktdaten (zum Zeitpunkt der Bestellung)
  "sku": "LAM-OAK-001",
  "product_name": "Laminat Eiche Natur 8mm",
  "quantity": 45,
  "unit": "m²",
  "price_paid": 34.90,

  // Für Suche
  "searchable_text": "Laminat Eiche Natur Projekt Müller Umbau November 2024"
}
```

**Sicherheit:**

```go
// Tenant-Token für Kundenspezifischen Index
func (s *SearchService) GetCustomerSearchKey(ctx context.Context, customerID string) (string, error) {
    tenantID := auth.TenantFromContext(ctx)

    // Meilisearch Tenant Token - kann NUR diesen Index durchsuchen
    token, err := s.meili.GenerateTenantToken(
        s.config.SearchKeyUID,
        map[string]interface{}{
            "filter": fmt.Sprintf("customer_id = %s", customerID),
        },
        &meilisearch.TenantTokenOptions{
            ExpiresAt: time.Now().Add(24 * time.Hour),
        },
    )
    return token, err
}
```

**Warum bei Algolia schwierig:**

| Aspekt | Algolia | Meilisearch (Self-Hosted) |
|--------|---------|---------------------------|
| **Datenschutz** | Kundendaten bei US-Firma | Daten bleiben bei uns |
| **Kosten** | Pro Record (teuer!) | Flat (eigene Infra) |
| **Compliance** | Auftragsverarbeitung nötig | Intern = kein Problem |
| **Flexibilität** | Begrenzte Index-Struktur | Volle Kontrolle |

### Synonym-Management (Support-Anforderung)

Der Support pflegt aktuell Synonyme direkt in Algolia. Diese Möglichkeit muss erhalten bleiben.

**Aktuelle Synonyme (Beispiel aus V2):**

```csv
vollkernplatte, kompaktplatte
swisspearl, eternit
oel, öl
Türrahmenprofil, Kanteln
fuma, tischlerplatte
Schichtstoffplatte, HPL, CPL, Kunstharzplatte, Schichtstoff, Kunstharz
Brandschutz, Feuer, feuerhemmend, hitzebeständig, Hitzeschutz
Multiplex, Sperrholz
```

**Anforderung: Admin UI für Synonym-Pflege**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  ADMIN PORTAL: SYNONYM-VERWALTUNG                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Suche: [___________________] [+ Neues Synonym]                             │
│                                                                              │
│  │ Synonymgruppe                              │ Tenant  │ Aktionen │        │
│  ├────────────────────────────────────────────┼─────────┼──────────┤        │
│  │ vollkernplatte, kompaktplatte              │ kurkj   │ ✏️ 🗑️    │        │
│  │ swisspearl, eternit                        │ kurkj   │ ✏️ 🗑️    │        │
│  │ Schichtstoffplatte, HPL, CPL, Kunstharz... │ kurkj   │ ✏️ 🗑️    │        │
│  │ Brandschutz, Feuer, feuerhemmend, ...      │ alle    │ ✏️ 🗑️    │        │
│  │                                                                          │
│                                                                              │
│  [CSV Import]  [CSV Export]                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Meilisearch Synonym-API:**

```go
// Synonyme setzen
synonyms := map[string][]string{
    "vollkernplatte": {"kompaktplatte"},
    "kompaktplatte":  {"vollkernplatte"},
    "swisspearl":     {"eternit"},
    "eternit":        {"swisspearl"},
    "hpl":            {"schichtstoffplatte", "cpl", "kunstharzplatte"},
    // Bidirektional für alle Varianten
}

index.UpdateSynonyms(synonyms)
```

**Features für Admin UI:**

| Feature | Beschreibung |
|---------|--------------|
| **CRUD** | Synonymgruppen erstellen, bearbeiten, löschen |
| **Tenant-spezifisch** | Synonyme pro Mandant oder global |
| **CSV Import/Export** | Migration von Algolia, Backup |
| **Vorschau** | Testen wie Suche mit Synonym funktioniert |
| **Audit Log** | Wer hat wann was geändert |

---

## Kandidaten-Übersicht

| Engine | Typ | Sprache | Lizenz | Self-Hosted |
|--------|-----|---------|--------|-------------|
| **Algolia** | SaaS only | - | Proprietär | Nein |
| **Elasticsearch** | Open Source | Java | Open Source | Ja |
| **Meilisearch** | Open Source | Rust | MIT | Ja |
| **Typesense** | Open Source | C++ | GPL | Ja |
| **OpenSearch** | Open Source | Java | Apache 2.0 | Ja |

---

## Detaillierter Vergleich

### 1. Algolia (Aktuell)

**Vorteile:**
- Extrem schnelle Suche (<50ms)
- Zero-Config Typo-Toleranz
- Eingebaute Personalisierung & Recommendations
- A/B Testing für Suchergebnisse
- Umfangreiche Analytics
- Auto-Scaling

**Nachteile:**
- **Kosten:** Sehr teuer bei Skalierung (pro 1000 Requests + Records)
- **Vendor Lock-in:** Kein Self-Hosting möglich
- **Datenhoheit:** Daten liegen bei Algolia
- **Keine Custom Ranking:** Begrenzte Anpassungsmöglichkeiten

**Kosten-Beispiel:**
```
100.000 Records × 10 Tenants = 1.000.000 Records
500.000 Searches/Monat
→ Ca. $500-1000+/Monat
```

---

### 2. Meilisearch ⭐ Empfehlung

**Vorteile:**
- **Extrem schnell:** <50ms Antwortzeit
- **Rust-basiert:** Effizient, geringer Speicherverbrauch
- **Einfache API:** REST JSON, ähnlich wie Algolia
- **Typo-Toleranz:** Out-of-the-box
- **Faceted Search:** Vollständig unterstützt
- **Hybrid Search:** Keyword + Vektor/Semantic (neu 2024)
- **Multi-Tenancy:** Tenant-Token für Index-Isolation
- **Kubernetes-ready:** Offizielles Helm Chart
- **MIT Lizenz:** Keine Einschränkungen

**Nachteile:**
- Keine eingebaute Personalisierung
- Keine A/B Testing Features
- Noch relativ jung (aber stabil)

**Features für Webshop:**
```
✅ Faceted Search (Farbe, Größe, Preis, etc.)
✅ Typo-Toleranz
✅ Synonyme
✅ Stop Words
✅ Ranking Rules (custom sortierbar)
✅ Filtering
✅ Geo Search
✅ Multi-Index
✅ API Key Management (per Tenant)
✅ Instant Search (<50ms)
```

**Kubernetes Deployment:**
```yaml
# Helm Chart verfügbar
helm repo add meilisearch https://meilisearch.github.io/meilisearch-kubernetes
helm install meilisearch meilisearch/meilisearch \
  --set environment.MEILI_MASTER_KEY=xxx \
  --set persistence.enabled=true \
  --set persistence.size=10Gi
```

**Ressourcen:**
```yaml
# Empfohlen für 1M+ Dokumente
resources:
  requests:
    memory: "2Gi"
    cpu: "1000m"
  limits:
    memory: "4Gi"
    cpu: "2000m"
```

---

### 3. Typesense

**Vorteile:**
- Sehr schnell (<50ms)
- C++ basiert, effizient
- Einfache API
- Dynamic Sorting (ohne separate Indizes)
- Vector Search Support
- Single Binary (einfaches Deployment)

**Nachteile:**
- RAM-basierter Index (teurer bei großen Datenmengen)
- GPL Lizenz (Copyleft)
- Weniger Dokumentation als Meilisearch

**Kosten Self-Hosted:**
```
$20/Monat VPS → 7.6x günstiger als Algolia
```

---

### 4. Elasticsearch / OpenSearch

**Vorteile:**
- Extrem mächtig für komplexe Queries
- Analytics & Aggregationen
- NLP & Machine Learning
- Bewährt in Enterprise-Umgebungen
- Riesiges Ökosystem

**Nachteile:**
- **Komplex:** Hoher Konfigurationsaufwand
- **Ressourcenhungrig:** JVM-basiert, braucht viel RAM
- **Langsamer Setup:** Wochen statt Tage
- **Overhead:** Für reine Produktsuche überdimensioniert

**Ressourcen:**
```yaml
# Minimum für Production
resources:
  requests:
    memory: "4Gi"  # Mindestens!
    cpu: "2000m"
  limits:
    memory: "8Gi"
    cpu: "4000m"
```

---

## Feature-Vergleich für Webshop

| Feature | Algolia | Meilisearch | Typesense | Elasticsearch |
|---------|---------|-------------|-----------|---------------|
| **Instant Search** | ✅ <50ms | ✅ <50ms | ✅ <50ms | ⚠️ 100-500ms |
| **Typo-Toleranz** | ✅ Auto | ✅ Auto | ✅ Auto | ⚠️ Konfiguration |
| **Faceted Search** | ✅ | ✅ | ✅ | ✅ |
| **Synonyme** | ✅ | ✅ | ✅ | ✅ |
| **Geo Search** | ✅ | ✅ | ✅ | ✅ |
| **Multi-Tenancy** | ✅ | ✅ | ✅ | ✅ |
| **Highlighting** | ✅ | ✅ | ✅ | ✅ |
| **Custom Ranking** | ⚠️ | ✅ | ✅ | ✅ |
| **Vector/Semantic** | ✅ | ✅ (neu) | ✅ | ✅ |
| **Self-Hosted** | ❌ | ✅ | ✅ | ✅ |
| **Setup-Zeit** | Minuten | Stunden | Stunden | Wochen |
| **RAM-Bedarf** | - | Niedrig | Mittel | Hoch |

---

## Kosten-Vergleich (1M Records, 500k Searches/Monat)

| Lösung | Monatliche Kosten | Anmerkung |
|--------|-------------------|-----------|
| **Algolia** | ~$800-1500 | Pro Record + Search |
| **Meilisearch Cloud** | ~$100-200 | Managed |
| **Meilisearch Self-Hosted** | ~$50-100 | 2x 4GB VPS/K8s Nodes |
| **Typesense Self-Hosted** | ~$40-80 | 2x 4GB VPS |
| **Elasticsearch** | ~$200-400 | 3x 8GB Nodes minimum |

---

## Empfehlung: Meilisearch

### Gründe

1. **Feature-Parität mit Algolia**
   - Gleiche Kernfunktionen für Produktsuche
   - Ähnliche API-Struktur (Migration einfacher)

2. **Kubernetes-Native**
   - Offizielles Helm Chart
   - StatefulSet mit Persistence
   - Passt zu unserer K8s-only Strategie

3. **Kostenersparnis**
   - 80-90% günstiger als Algolia
   - Keine Vendor Lock-in Kosten

4. **Datenhoheit**
   - Daten bleiben in unserer Infrastruktur
   - DSGVO-konform ohne Auftragsverarbeitung

5. **Zukunftssicher**
   - Aktive Entwicklung (Rust, modern)
   - Hybrid Search für AI-Features vorbereitet

6. **Einfache Migration**
   - REST API ähnlich zu Algolia
   - Faceted Search funktioniert identisch

### Risiken

- Jüngeres Projekt (aber stabil, MIT-lizenziert)
- Keine eingebaute Personalisierung (kann im Shop-Service gebaut werden)
- Kein A/B Testing (brauchen wir das wirklich?)

---

## Migrations-Strategie

### Phase 1: Abstraktion (V3 von Anfang an)

```go
// Abstrakte Search-Interface
type SearchEngine interface {
    Index(ctx context.Context, tenant, indexType string, docs []Document) error
    Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
    Delete(ctx context.Context, tenant, indexType string, ids []string) error
    Configure(ctx context.Context, tenant, indexType string, config *IndexConfig) error
}

// Implementierungen
type AlgoliaEngine struct { ... }     // Für Migration/Fallback
type MeilisearchEngine struct { ... } // Neue Lösung
```

### Phase 2: Parallel-Betrieb (Optional)

```
┌─────────────┐     ┌─────────────┐
│   Search    │────▶│   Algolia   │  (Lesen)
│   Service   │     └─────────────┘
│             │
│             │────▶│ Meilisearch │  (Schreiben + Lesen für Tests)
└─────────────┘     └─────────────┘
```

### Phase 3: Umstellung

```
Feature Flag: USE_MEILISEARCH=true
→ Vollständiger Wechsel auf Meilisearch
→ Algolia kündigen
```

---

## Meilisearch Integration für V3

### Index-Struktur

```
Indices pro Tenant und Sprache:
  {tenant}_products_{language}     # z.B. kurkj_products_de
  {tenant}_categories_{language}
```

### Dokument-Schema

```json
{
  "id": "ABC-123",
  "sku": "ABC-123",
  "name": "Laminat Eiche Natur",
  "description": "Hochwertiges Laminat...",
  "price": 45.90,
  "categories": ["bodenbelaege", "laminat"],
  "brand": "Swiss Krono",

  "attributes": {
    "thickness": "8mm",
    "width": "193mm",
    "length": "1380mm",
    "color": "Eiche Natur",
    "surface": "Matt",
    "wood_type": "Eiche"
  },

  "in_stock": true,
  "stock_quantity": 1250,

  "image_url": "https://cdn.example.com/products/abc-123.jpg",
  "tenant": "kurkj",

  "_geo": {
    "lat": 47.05,
    "lng": 7.45
  }
}
```

### Index-Konfiguration

```go
config := &meilisearch.IndexConfig{
    PrimaryKey: "id",

    SearchableAttributes: []string{
        "name",
        "description",
        "sku",
        "brand",
        "attributes.color",
    },

    FilterableAttributes: []string{
        "categories",
        "brand",
        "price",
        "in_stock",
        "attributes.thickness",
        "attributes.color",
        "attributes.surface",
        "attributes.wood_type",
        "tenant",
    },

    SortableAttributes: []string{
        "price",
        "name",
        "stock_quantity",
    },

    RankingRules: []string{
        "words",
        "typo",
        "proximity",
        "attribute",
        "sort",
        "exactness",
    },

    TypoTolerance: &meilisearch.TypoTolerance{
        Enabled: true,
        MinWordSizeForTypos: meilisearch.MinWordSize{
            OneTypo:  4,
            TwoTypos: 8,
        },
    },

    Synonyms: map[string][]string{
        "laminat":  {"laminatboden", "klicklaminat"},
        "parkett":  {"parkettboden", "echtholzparkett"},
        "eiche":    {"oak", "eichenholz"},
    },
}
```

### Service Implementation

```go
// services/search/internal/meilisearch/client.go
package meilisearch

type Client struct {
    client *meilisearch.Client
    config *Config
}

func (c *Client) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    index := c.client.Index(c.indexName(req.Tenant, req.IndexType, req.Language))

    searchReq := &meilisearch.SearchRequest{
        Query:                 req.Query,
        Limit:                 int64(req.Limit),
        Offset:                int64(req.Offset),
        Filter:                c.buildFilters(req.Filters),
        Facets:                req.Facets,
        AttributesToHighlight: []string{"name", "description"},
        HighlightPreTag:       "<mark>",
        HighlightPostTag:      "</mark>",
    }

    if req.Sort != "" {
        searchReq.Sort = []string{req.Sort}
    }

    result, err := index.Search(req.Query, searchReq)
    if err != nil {
        return nil, fmt.Errorf("meilisearch search: %w", err)
    }

    return c.mapResponse(result), nil
}

func (c *Client) Index(ctx context.Context, tenant, indexType, language string, docs []Document) error {
    index := c.client.Index(c.indexName(tenant, indexType, language))

    task, err := index.AddDocuments(docs, "id")
    if err != nil {
        return fmt.Errorf("add documents: %w", err)
    }

    // Optional: Auf Completion warten
    if c.config.WaitForIndexing {
        _, err = c.client.WaitForTask(task.TaskUID)
        if err != nil {
            return fmt.Errorf("wait for task: %w", err)
        }
    }

    return nil
}
```

### Kubernetes Deployment

```yaml
# infrastructure/kubernetes/base/meilisearch/deployment.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: meilisearch
spec:
  serviceName: meilisearch
  replicas: 1
  selector:
    matchLabels:
      app: meilisearch
  template:
    metadata:
      labels:
        app: meilisearch
    spec:
      containers:
      - name: meilisearch
        image: getmeili/meilisearch:v1.6
        ports:
        - containerPort: 7700
        env:
        - name: MEILI_MASTER_KEY
          valueFrom:
            secretKeyRef:
              name: meilisearch-secrets
              key: master-key
        - name: MEILI_ENV
          value: "production"
        - name: MEILI_DB_PATH
          value: "/meili_data"
        resources:
          requests:
            memory: "2Gi"
            cpu: "500m"
          limits:
            memory: "4Gi"
            cpu: "2000m"
        volumeMounts:
        - name: data
          mountPath: /meili_data
        livenessProbe:
          httpGet:
            path: /health
            port: 7700
          initialDelaySeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 7700
          initialDelaySeconds: 5
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 20Gi
---
apiVersion: v1
kind: Service
metadata:
  name: meilisearch
spec:
  ports:
  - port: 7700
  selector:
    app: meilisearch
```

---

## Fazit

| Kriterium | Algolia | Meilisearch | Entscheidung |
|-----------|---------|-------------|--------------|
| Features für Webshop | ✅ | ✅ | Gleich |
| Kosten | ❌ Teuer | ✅ 80-90% günstiger | **Meilisearch** |
| Self-Hosted | ❌ | ✅ | **Meilisearch** |
| Datenhoheit | ❌ | ✅ | **Meilisearch** |
| K8s Integration | ⚠️ SaaS | ✅ Helm Chart | **Meilisearch** |
| Setup-Komplexität | ✅ Einfach | ✅ Einfach | Gleich |
| Risiko | ✅ Etabliert | ⚠️ Jünger | Algolia |

**Empfehlung:** Meilisearch für V3, mit Abstraktionsschicht für Flexibilität.

---

## Quellen

- [Typesense vs Algolia vs Elasticsearch vs Meilisearch](https://typesense.org/typesense-vs-algolia-vs-elasticsearch-vs-meilisearch/)
- [Meilisearch: Algolia Alternatives](https://www.meilisearch.com/blog/algolia-alternatives)
- [Meilisearch Documentation](https://docs.meilisearch.com/)