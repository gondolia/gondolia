# Kategorie-Architektur

## Problem (V2)

In V2 existieren Kategorien **nicht als eigenständige Entität**, sondern nur als:
- Facette in Algolia (`categories: ["flooring", "laminate"]`)
- Importierte Codes aus Akeneo PIM

**Nachteile:**
- ❌ Keine SEO-freundlichen URLs (`/produkte?filter=category:laminate`)
- ❌ Keine Kategorie-Landingpages
- ❌ Keine Meta-Descriptions für Kategorien
- ❌ Keine Breadcrumb-Navigation
- ❌ Keine kategorie-spezifischen Inhalte (Banner, Texte)
- ❌ Keine strukturierten Daten (Schema.org)
- ❌ Keine Kategorie-Hierarchie navigierbar

---

## Lösung (V3): Kategorien als First-Class Entity

### Datenmodell

```sql
-- Kategorie-Tabelle
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    -- Identifikation
    code VARCHAR(100) NOT NULL,           -- PIM-Code (z.B. "laminate_flooring")
    slug VARCHAR(255) NOT NULL,           -- URL-Slug (z.B. "laminat-bodenbelaege")

    -- Hierarchie
    parent_id UUID REFERENCES categories(id),
    path LTREE NOT NULL,                  -- Materialized Path für Performance
    level INTEGER NOT NULL DEFAULT 0,     -- Tiefe in Hierarchie
    position INTEGER NOT NULL DEFAULT 0,  -- Sortierung

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_visible BOOLEAN NOT NULL DEFAULT true,  -- In Navigation anzeigen

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(tenant_id, code),
    UNIQUE(tenant_id, slug)
);

-- Index für Hierarchie-Queries
CREATE INDEX idx_categories_path ON categories USING GIST (path);
CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_categories_tenant_active ON categories(tenant_id, is_active);

-- Kategorie-Übersetzungen (i18n)
CREATE TABLE category_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    language VARCHAR(5) NOT NULL,         -- de, en, fr, it

    -- SEO-relevante Felder
    name VARCHAR(255) NOT NULL,           -- "Laminat-Bodenbeläge"
    description TEXT,                      -- Kategorie-Beschreibung
    meta_title VARCHAR(70),               -- SEO Title (max 70 Zeichen)
    meta_description VARCHAR(160),        -- SEO Description (max 160 Zeichen)

    -- Content
    headline VARCHAR(255),                -- H1 auf Kategorie-Seite
    intro_text TEXT,                      -- Einleitungstext
    footer_text TEXT,                     -- Text unter Produkten

    UNIQUE(category_id, language)
);

-- Kategorie-Medien (Banner, Icons)
CREATE TABLE category_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,            -- banner, icon, thumbnail, hero

    url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    width INTEGER,
    height INTEGER,

    position INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_category_media_category ON category_media(category_id, type);
```

### Hierarchie-Beispiel

```
Master Catalog (Tenant: kurkj)
├── Bodenbeläge (level: 0, path: 'bodenbelaege')
│   ├── Laminat (level: 1, path: 'bodenbelaege.laminat')
│   │   ├── 7mm Laminat (level: 2, path: 'bodenbelaege.laminat.7mm')
│   │   ├── 8mm Laminat (level: 2, path: 'bodenbelaege.laminat.8mm')
│   │   └── 12mm Laminat (level: 2, path: 'bodenbelaege.laminat.12mm')
│   ├── Parkett (level: 1, path: 'bodenbelaege.parkett')
│   └── Vinyl (level: 1, path: 'bodenbelaege.vinyl')
├── Türen (level: 0, path: 'tueren')
│   ├── Innentüren (level: 1, path: 'tueren.innentueren')
│   └── Haustüren (level: 1, path: 'tueren.haustueren')
└── Werkzeuge (level: 0, path: 'werkzeuge')
```

### URL-Struktur (SEO-optimiert)

```
/kategorien/bodenbelaege                      → Hauptkategorie
/kategorien/bodenbelaege/laminat              → Unterkategorie
/kategorien/bodenbelaege/laminat/8mm-laminat  → Sub-Unterkategorie

Mit Sprache:
/de/kategorien/bodenbelaege
/fr/categories/revetements-de-sol
/en/categories/flooring
```

### API-Struktur

```go
// Category Entity
type Category struct {
    ID          uuid.UUID            `json:"id"`
    TenantID    uuid.UUID            `json:"tenant_id"`
    Code        string               `json:"code"`
    Slug        string               `json:"slug"`
    ParentID    *uuid.UUID           `json:"parent_id,omitempty"`
    Path        string               `json:"path"`
    Level       int                  `json:"level"`
    Position    int                  `json:"position"`
    IsActive    bool                 `json:"is_active"`
    IsVisible   bool                 `json:"is_visible"`

    // Loaded via joins
    Translation *CategoryTranslation `json:"translation,omitempty"`
    Children    []*Category          `json:"children,omitempty"`
    Media       []*CategoryMedia     `json:"media,omitempty"`

    // Computed
    ProductCount int                 `json:"product_count,omitempty"`
    Breadcrumbs  []*BreadcrumbItem   `json:"breadcrumbs,omitempty"`
}

type CategoryTranslation struct {
    Language        string  `json:"language"`
    Name            string  `json:"name"`
    Description     *string `json:"description,omitempty"`
    MetaTitle       *string `json:"meta_title,omitempty"`
    MetaDescription *string `json:"meta_description,omitempty"`
    Headline        *string `json:"headline,omitempty"`
    IntroText       *string `json:"intro_text,omitempty"`
    FooterText      *string `json:"footer_text,omitempty"`
}

type BreadcrumbItem struct {
    Name string `json:"name"`
    Slug string `json:"slug"`
    URL  string `json:"url"`
}
```

### REST API Endpoints

```
GET  /api/v1/categories                    → Liste (mit Hierarchie-Option)
GET  /api/v1/categories/:slug              → Einzelne Kategorie
GET  /api/v1/categories/:slug/children     → Direkte Unterkategorien
GET  /api/v1/categories/:slug/products     → Produkte in Kategorie
GET  /api/v1/categories/:slug/breadcrumbs  → Breadcrumb-Pfad
GET  /api/v1/categories/tree               → Vollständiger Baum

# Admin
POST   /api/v1/admin/categories            → Kategorie erstellen
PUT    /api/v1/admin/categories/:id        → Kategorie aktualisieren
DELETE /api/v1/admin/categories/:id        → Kategorie löschen
PATCH  /api/v1/admin/categories/:id/move   → Position/Parent ändern
```

### Response-Beispiel

```json
GET /api/v1/categories/laminat?lang=de

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "laminate_flooring",
  "slug": "laminat",
  "level": 1,
  "is_active": true,

  "translation": {
    "language": "de",
    "name": "Laminat",
    "description": "Hochwertige Laminatböden für jeden Raum...",
    "meta_title": "Laminat kaufen | Große Auswahl | Shop Name",
    "meta_description": "Laminatböden in Premium-Qualität. ✓ Große Auswahl ✓ Schnelle Lieferung ✓ Top Preise. Jetzt entdecken!",
    "headline": "Laminatböden für Ihr Zuhause",
    "intro_text": "Entdecken Sie unsere große Auswahl an..."
  },

  "breadcrumbs": [
    { "name": "Home", "slug": "", "url": "/" },
    { "name": "Bodenbeläge", "slug": "bodenbelaege", "url": "/kategorien/bodenbelaege" },
    { "name": "Laminat", "slug": "laminat", "url": "/kategorien/bodenbelaege/laminat" }
  ],

  "children": [
    { "slug": "7mm-laminat", "name": "7mm Laminat", "product_count": 45 },
    { "slug": "8mm-laminat", "name": "8mm Laminat", "product_count": 128 },
    { "slug": "12mm-laminat", "name": "12mm Laminat", "product_count": 67 }
  ],

  "media": [
    {
      "type": "banner",
      "url": "https://cdn.example.com/categories/laminat-banner.jpg",
      "alt_text": "Laminatböden Kollektion",
      "width": 1920,
      "height": 400
    }
  ],

  "product_count": 240
}
```

---

## SEO-Optimierungen

### 1. URL-Struktur

```
✅ RICHTIG (V3):
/kategorien/bodenbelaege/laminat
/categories/flooring/laminate (EN)

❌ FALSCH (V2):
/produkte?filter[categories][0]=laminate_flooring
```

### 2. Meta-Tags (automatisch generiert)

```html
<!-- Kategorie-Seite -->
<title>Laminat kaufen | Große Auswahl | Shop Name</title>
<meta name="description" content="Laminatböden in Premium-Qualität. ✓ Große Auswahl ✓ Schnelle Lieferung ✓ Top Preise. Jetzt entdecken!">
<link rel="canonical" href="https://shop.example.com/kategorien/bodenbelaege/laminat">

<!-- Open Graph -->
<meta property="og:title" content="Laminat kaufen | Shop Name">
<meta property="og:description" content="Laminatböden in Premium-Qualität...">
<meta property="og:image" content="https://cdn.example.com/categories/laminat-og.jpg">
<meta property="og:type" content="website">

<!-- Hreflang für Mehrsprachigkeit -->
<link rel="alternate" hreflang="de" href="https://shop.example.com/de/kategorien/laminat">
<link rel="alternate" hreflang="fr" href="https://shop.example.com/fr/categories/stratifie">
<link rel="alternate" hreflang="en" href="https://shop.example.com/en/categories/laminate">
```

### 3. Strukturierte Daten (Schema.org)

```json
{
  "@context": "https://schema.org",
  "@type": "CollectionPage",
  "name": "Laminat",
  "description": "Hochwertige Laminatböden für jeden Raum...",
  "url": "https://shop.example.com/kategorien/bodenbelaege/laminat",
  "image": "https://cdn.example.com/categories/laminat-banner.jpg",

  "breadcrumb": {
    "@type": "BreadcrumbList",
    "itemListElement": [
      {
        "@type": "ListItem",
        "position": 1,
        "name": "Home",
        "item": "https://shop.example.com/"
      },
      {
        "@type": "ListItem",
        "position": 2,
        "name": "Bodenbeläge",
        "item": "https://shop.example.com/kategorien/bodenbelaege"
      },
      {
        "@type": "ListItem",
        "position": 3,
        "name": "Laminat",
        "item": "https://shop.example.com/kategorien/bodenbelaege/laminat"
      }
    ]
  },

  "mainEntity": {
    "@type": "ItemList",
    "numberOfItems": 240,
    "itemListElement": [
      {
        "@type": "Product",
        "position": 1,
        "name": "Swiss Krono Laminat Eiche Natur 8mm",
        "url": "https://shop.example.com/produkte/swiss-krono-eiche-natur-8mm"
      }
    ]
  }
}
```

### 4. Sitemap-Integration

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <!-- Kategorien mit Priorität -->
  <url>
    <loc>https://shop.example.com/kategorien/bodenbelaege</loc>
    <lastmod>2024-01-15</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>https://shop.example.com/kategorien/bodenbelaege/laminat</loc>
    <lastmod>2024-01-15</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.7</priority>
  </url>
</urlset>
```

---

## Synchronisation mit Akeneo PIM

### Import-Flow

```
Akeneo PIM                     Catalog Service                    Database
     │                               │                                │
     │  GET /categories              │                                │
     │  ─────────────────────────▶   │                                │
     │                               │                                │
     │  { categories: [...] }        │                                │
     │  ◀─────────────────────────   │                                │
     │                               │                                │
     │                               │  Upsert categories             │
     │                               │  ─────────────────────────────▶│
     │                               │                                │
     │                               │  Generate slugs (if missing)   │
     │                               │  Build hierarchy (path, level) │
     │                               │  ─────────────────────────────▶│
```

### Slug-Generierung

```go
// Automatische Slug-Generierung aus PIM-Name
func generateSlug(name, language string) string {
    slug := strings.ToLower(name)

    // Umlaute ersetzen
    replacements := map[string]string{
        "ä": "ae", "ö": "oe", "ü": "ue",
        "ß": "ss", "é": "e", "è": "e",
        "à": "a", "ô": "o", "î": "i",
    }
    for from, to := range replacements {
        slug = strings.ReplaceAll(slug, from, to)
    }

    // Nicht-alphanumerische Zeichen durch Bindestrich
    slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")

    // Mehrfache Bindestriche entfernen
    slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

    // Trim
    slug = strings.Trim(slug, "-")

    return slug
}

// Beispiele:
// "Laminat-Bodenbeläge" → "laminat-bodenbelaege"
// "Türen & Zargen"      → "tueren-zargen"
// "7mm Laminat"         → "7mm-laminat"
```

### Mapping Akeneo → Shop

```go
type AkeneoCategoryMapper struct {
    slugGenerator SlugGenerator
    translations  map[string]string // code → custom slug (override)
}

func (m *AkeneoCategoryMapper) Map(akeneo *AkeneoCategory, tenant *Tenant) (*Category, error) {
    category := &Category{
        TenantID: tenant.ID,
        Code:     akeneo.Code,
        IsActive: true,
    }

    // Custom Slug oder generiert
    if customSlug, ok := m.translations[akeneo.Code]; ok {
        category.Slug = customSlug
    } else {
        // Slug aus deutschem Namen generieren (Fallback)
        if name := akeneo.Labels["de_CH"]; name != "" {
            category.Slug = m.slugGenerator.Generate(name)
        } else {
            category.Slug = m.slugGenerator.Generate(akeneo.Code)
        }
    }

    // Parent finden
    if akeneo.Parent != "" {
        // Parent-Kategorie muss vorher importiert sein
        parent, err := m.repo.FindByCode(tenant.ID, akeneo.Parent)
        if err != nil {
            return nil, fmt.Errorf("parent not found: %s", akeneo.Parent)
        }
        category.ParentID = &parent.ID
        category.Path = parent.Path + "." + category.Slug
        category.Level = parent.Level + 1
    } else {
        category.Path = category.Slug
        category.Level = 0
    }

    // Übersetzungen mappen
    category.Translations = make([]*CategoryTranslation, 0)
    for locale, name := range akeneo.Labels {
        lang := m.mapLocale(locale) // "de_CH" → "de"
        category.Translations = append(category.Translations, &CategoryTranslation{
            Language: lang,
            Name:     name,
        })
    }

    return category, nil
}
```

---

## Integration mit Meilisearch

### Produkt-Dokument erweitert

```json
{
  "id": "ABC-123",
  "sku": "ABC-123",
  "name": "Swiss Krono Laminat Eiche Natur 8mm",

  "categories": [
    {
      "code": "laminate_flooring",
      "slug": "laminat",
      "name": "Laminat",
      "path": "bodenbelaege.laminat",
      "level": 1
    },
    {
      "code": "8mm_laminate",
      "slug": "8mm-laminat",
      "name": "8mm Laminat",
      "path": "bodenbelaege.laminat.8mm-laminat",
      "level": 2
    }
  ],

  "category_codes": ["laminate_flooring", "8mm_laminate"],
  "category_slugs": ["laminat", "8mm-laminat"],
  "category_paths": ["bodenbelaege.laminat", "bodenbelaege.laminat.8mm-laminat"],

  "brand": "Swiss Krono",
  "price": 45.90
}
```

### Meilisearch Filter

```go
// Produkte in Kategorie und allen Unterkategorien
filter := fmt.Sprintf("category_paths CONTAINS '%s'", category.Path)

// Nur direkte Kategorie
filter := fmt.Sprintf("category_slugs = '%s'", category.Slug)
```

### Facetten mit Kategorie-Hierarchie

```json
{
  "facets": ["category_paths"],
  "filter": "category_paths CONTAINS 'bodenbelaege'"
}

// Response
{
  "facetDistribution": {
    "category_paths": {
      "bodenbelaege.laminat": 240,
      "bodenbelaege.laminat.7mm-laminat": 45,
      "bodenbelaege.laminat.8mm-laminat": 128,
      "bodenbelaege.laminat.12mm-laminat": 67,
      "bodenbelaege.parkett": 180,
      "bodenbelaege.vinyl": 95
    }
  }
}
```

---

## Service-Zuordnung

```
services/
├── catalog/                    # Hauptservice für Kategorien + Produkte
│   ├── internal/
│   │   ├── category/          # Kategorie-Domain
│   │   │   ├── handler.go     # REST API
│   │   │   ├── service.go     # Business Logic
│   │   │   ├── repository.go  # DB Access
│   │   │   ├── models.go      # Entities
│   │   │   └── mapper.go      # Akeneo Mapping
│   │   ├── product/           # Produkt-Domain
│   │   └── search/            # Meilisearch Integration
```

---

## Migration von V2

### Schritt 1: Kategorie-Struktur aus PIM extrahieren

```sql
-- Kategorien aus Akeneo importieren
INSERT INTO categories (tenant_id, code, slug, path, level)
SELECT
    t.id,
    ac.code,
    generate_slug(ac.label_de),
    build_path(ac.code, ac.parent_code),
    count_parents(ac.code)
FROM akeneo_categories ac
CROSS JOIN tenants t
WHERE t.code = 'kurkj';
```

### Schritt 2: Übersetzungen migrieren

```sql
-- Übersetzungen aus PIM Labels
INSERT INTO category_translations (category_id, language, name)
SELECT
    c.id,
    'de',
    ac.label_de
FROM categories c
JOIN akeneo_categories ac ON c.code = ac.code;
```

### Schritt 3: SEO-Felder manuell pflegen

Nach dem initialen Import müssen SEO-relevante Felder manuell gepflegt werden:
- `meta_title`
- `meta_description`
- `headline`
- `intro_text`
- `footer_text`
- Kategorie-Banner

---

## Zusammenfassung

| Aspekt | V2 (Algolia Facette) | V3 (First-Class Entity) |
|--------|----------------------|-------------------------|
| **URL** | Query-Parameter | SEO-freundlich |
| **SEO** | Keine Meta-Tags | Vollständige SEO-Kontrolle |
| **Navigation** | Facetten-Filter | Breadcrumbs + Hierarchie |
| **Content** | Keiner | Beschreibungen, Banner |
| **Structured Data** | Keine | Schema.org |
| **Sitemap** | Nur Produkte | Kategorien + Produkte |
| **Multi-Language** | Nur Facetten-Labels | Vollständige Übersetzungen |
| **PIM-Sync** | Nur Codes | Entities mit Mapping |
