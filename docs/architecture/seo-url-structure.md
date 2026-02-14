# SEO-freundliche URL-Struktur

Stand: 2026-02-14

## Übersicht

### Problem (Ist-Zustand)

Gondolia verwendet aktuell UUIDs in allen öffentlichen URLs:

```
/products/b0000000-0014-0000-0000-000000000001
/categories/a0000000-0000-0000-0000-000000000001
```

**Warum das inakzeptabel ist:**

| Problem | Auswirkung |
|---------|------------|
| Keine Keywords in URL | Google gewichtet URL-Keywords für Relevanz |
| Nicht merkbar | B2B-Kunden teilen Links per E-Mail, Chat, Telefon — UUID ist nicht kommunizierbar |
| Kein Vertrauen | Kryptische URLs wirken unseriös, CTR in Suchergebnissen sinkt |
| Keine Hierarchie sichtbar | Kategorien-Tiefe nicht erkennbar, Breadcrumbs nicht aus URL ableitbar |
| Analytics unlesbar | Traffic-Reports mit UUIDs sind nicht interpretierbar |

### Lösung

Alle öffentlichen URLs erhalten **menschenlesbare Slugs** mit Sprachwechsel-Unterstützung, hierarchischen Kategorie-Pfaden und durchgängigen Canonical/Redirect-Strategien.

---

## 1. URL-Schema

### 1.1 Grundstruktur

```
https://shop.example.com/{sprache}/{entitaet-pfad}/{slug}
```

Das Sprach-Präfix (`/de/`, `/fr/`, `/en/`) steht **immer** an erster Stelle. Die Entitäts-Pfadsegmente werden pro Sprache übersetzt.

### 1.2 Produkte

**Pattern:**

```
/{lang}/produkte/{product-slug}
/{lang}/products/{product-slug}        (EN)
/{lang}/produits/{product-slug}        (FR)
```

**Beispiele:**

| Produkt | URL (DE) |
|---------|----------|
| Hydraulikzylinder CDT3 50mm Hub 200mm | `/de/produkte/hydraulikzylinder-cdt3-50mm-hub-200mm` |
| Netzwerkkabel Cat6a 5m blau S/FTP | `/de/produkte/netzwerkkabel-cat6a-5m-blau-s-ftp` |
| Nitrilhandschuh puderfrei | `/de/produkte/nitrilhandschuh-puderfrei` |
| Spanplatte Eurospan 2800×2070 16mm Eiche | `/de/produkte/spanplatte-eurospan-2800x2070-16mm-eiche` |

**Entscheid: `/produkte/` statt `/p/`**

| Option | Pro | Contra |
|--------|-----|--------|
| `/produkte/{slug}` | SEO-Keyword, lesbar, professionell | Länger |
| `/p/{slug}` | Kurz | Kein SEO-Wert, kryptisch |
| `/p/{sku}` | Stabil, kurz | Kein SEO-Wert, SKU nicht immer lesbar |

**Empfehlung:** `/produkte/{slug}` für maximale SEO-Wirkung. SKU-basierte Kurz-URLs (`/p/HYD-ZYL-50-200`) als zusätzlicher Zugang mit 301-Redirect auf die kanonische Slug-URL.

### 1.3 Kategorien (hierarchisch)

**Pattern:**

```
/{lang}/kategorien/{level-0-slug}
/{lang}/kategorien/{level-0-slug}/{level-1-slug}
/{lang}/kategorien/{level-0-slug}/{level-1-slug}/{level-2-slug}
```

**Beispiele:**

```
/de/kategorien/industriebedarf
/de/kategorien/industriebedarf/hydraulik
/de/kategorien/industriebedarf/hydraulik/hydraulikzylinder
/de/kategorien/elektro-elektronik/kabel/netzwerkkabel
/de/kategorien/arbeitsschutz/handschutz/nitrilhandschuhe
```

**Maximale Tiefe:** 4 Ebenen. Tiefere Hierarchien werden in der URL abgeflacht (SEO-Best-Practice: Google bevorzugt kürzere URLs).

**Sprach-Beispiele:**

```
/de/kategorien/industriebedarf/hydraulik
/fr/categories/fournitures-industrielles/hydraulique
/en/categories/industrial-supplies/hydraulics
```

**Hinweis:** Die Slugs in der `categories`-Tabelle (siehe [Kategorie-Architektur](category-architecture.md)) werden pro Sprache gespeichert. Die URL-Pfad-Segmente werden aus den Slugs der gesamten Elternkette zusammengesetzt.

### 1.4 Hersteller

**Pattern:**

```
/{lang}/hersteller/{manufacturer-slug}
/{lang}/fabricants/{manufacturer-slug}    (FR)
/{lang}/manufacturers/{manufacturer-slug} (EN)
```

**Beispiele:**

```
/de/hersteller/nexoflux
/de/hersteller/bosch-rexroth
/de/hersteller/phoenix-contact
```

Hersteller-Slugs sind sprachunabhängig (Markennamen werden nicht übersetzt), aber das Pfadsegment (`hersteller`/`manufacturers`) ist sprachabhängig.

### 1.5 Kategorie + Produkt Kombination

Ein Produkt kann über seine Kategorie erreicht werden. Dies ist **kein kanonischer Pfad**, sondern ein Navigations-Kontext:

```
/de/kategorien/industriebedarf/hydraulik/hydraulikzylinder-cdt3-50mm-hub-200mm
```

**Verhalten:**
- Die URL funktioniert und zeigt das Produkt an
- `<link rel="canonical">` zeigt auf die kanonische Produkt-URL: `/de/produkte/hydraulikzylinder-cdt3-50mm-hub-200mm`
- Die Breadcrumb zeigt die Kategorie-Hierarchie des Kontext-Pfads

**Warum dieses Pattern?**
- Nutzer navigieren über Kategorien und erwarten, dass die URL den Pfad widerspiegelt
- Google sieht über Canonical die eindeutige URL — kein Duplicate Content
- Breadcrumbs sind kontextbezogen korrekt

### 1.6 Hersteller + Produkt Kombination

Analog zu Kategorien:

```
/de/hersteller/nexoflux/hydraulikzylinder-cdt3-50mm-hub-200mm
```

- Funktioniert, zeigt Produkt an
- Canonical → `/de/produkte/hydraulikzylinder-cdt3-50mm-hub-200mm`
- Breadcrumb: Home → Hersteller → NexoFlux → Hydraulikzylinder CDT3

### 1.7 Suche

**Pattern:**

```
/{lang}/suche?q={suchbegriff}
/{lang}/suche?q={suchbegriff}&kategorie={category-slug}&hersteller={manufacturer-slug}
/{lang}/suche?q={suchbegriff}&seite=2&sortierung=preis-aufsteigend
```

**Beispiele:**

```
/de/suche?q=hydraulikzylinder
/de/suche?q=hydraulik&kategorie=industriebedarf&hersteller=nexoflux
/de/suche?q=kabel+cat6a&seite=2&sortierung=preis-aufsteigend
/fr/recherche?q=hydraulique
/en/search?q=hydraulic+cylinder
```

**Query-Parameter:**

| Parameter | DE | FR | EN | Beschreibung |
|-----------|----|----|-----|-------------|
| `q` | `q` | `q` | `q` | Suchbegriff |
| `kategorie` | `kategorie` | `categorie` | `category` | Kategorie-Filter (Slug) |
| `hersteller` | `hersteller` | `fabricant` | `manufacturer` | Hersteller-Filter (Slug) |
| `seite` | `seite` | `page` | `page` | Seitennummer |
| `sortierung` | `sortierung` | `tri` | `sort` | Sortierung |

**SEO:** Suchseiten erhalten `<meta name="robots" content="noindex, follow">` — Suchresultate sollen nicht indexiert werden, aber die verlinkten Produkte schon.

### 1.8 Produkt-Varianten

Varianten (siehe [Produkttypen-Architektur](product-types.md)) werden über Query-Parameter abgebildet — **nicht** als eigene URL-Pfade.

**Pattern:**

```
/{lang}/produkte/{master-slug}?variante={variant-sku}
/{lang}/produkte/{master-slug}?durchmesser=63mm&hub=300mm&betriebsdruck=250bar
```

**Beispiele:**

```
/de/produkte/hydraulikzylinder-cdt3?variante=REXROTH-CDT3-63-300-250
/de/produkte/netzwerkkabel-cat6a?laenge=5m&farbe=blau&schirmung=s-ftp
/de/produkte/poloshirt-hakro-performance?groesse=l&farbe=navy
```

**Begründung für Query-Parameter statt eigene URLs:**
- Variantenprodukte haben bis zu 48 Kombinationen — 48 separate URLs würden Thin Content erzeugen
- Google erkennt Variantenprodukte besser als eine Seite mit strukturierten Daten
- Teilen: Die vollständige URL mit Query-Parametern ist teilbar und bookmarkbar
- Die Master-Seite rankt für den Oberbegriff, nicht jede Variante einzeln

**Canonical:** Varianten-URLs mit Query-Parametern haben als Canonical die Master-URL **ohne** Parameter:

```html
<link rel="canonical" href="https://shop.example.com/de/produkte/hydraulikzylinder-cdt3">
```

**Ausnahme:** Wenn Varianten eigenständige SEO-Relevanz haben (z.B. «Hydraulikzylinder 63mm» ist ein eigenständiger Suchbegriff), kann pro Variante ein eigenes Canonical gesetzt werden. Das ist eine Admin-Entscheidung pro Produkt.

### 1.9 Statische Seiten

```
/de/ueber-uns
/de/agb
/de/datenschutz
/de/impressum
/de/kontakt
/de/versand-und-lieferung
```

### 1.10 Zusammenfassung URL-Schema

| Entität | URL-Pattern | Beispiel |
|---------|-------------|----------|
| **Produkt** | `/{lang}/produkte/{slug}` | `/de/produkte/hydraulikzylinder-cdt3-50mm` |
| **Produkt (kurz)** | `/{lang}/p/{sku}` → 301 | `/de/p/HYD-ZYL-50` → 301 auf Slug-URL |
| **Kategorie** | `/{lang}/kategorien/{path}` | `/de/kategorien/industriebedarf/hydraulik` |
| **Hersteller** | `/{lang}/hersteller/{slug}` | `/de/hersteller/nexoflux` |
| **Kat. + Produkt** | `/{lang}/kategorien/{path}/{slug}` | `/de/kategorien/.../hydraulikzylinder-cdt3` |
| **Suche** | `/{lang}/suche?q=...` | `/de/suche?q=hydraulik` |
| **Variante** | `/{lang}/produkte/{slug}?variante=...` | `/de/produkte/...?variante=SKU` |

---

## 2. Slug-Generierung

### 2.1 Algorithmus

```go
// pkg/slug/slug.go
package slug

import (
    "regexp"
    "strings"
    "unicode"
)

var umlautReplacements = map[string]string{
    "ä": "ae", "ö": "oe", "ü": "ue",
    "Ä": "Ae", "Ö": "Oe", "Ü": "Ue",
    "ß": "ss",
    "é": "e", "è": "e", "ê": "e", "ë": "e",
    "à": "a", "â": "a",
    "ô": "o", "ò": "o",
    "î": "i", "ï": "i",
    "ù": "u", "û": "u",
    "ç": "c",
    "ñ": "n",
}

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)
var multiDash = regexp.MustCompile(`-+`)

// Generate erzeugt einen URL-Slug aus einem beliebigen String.
//
// Regeln:
//   1. Lowercase
//   2. Umlaute → Transliteration (ä→ae, ö→oe, ü→ue, ß→ss)
//   3. Sonderzeichen entfernen
//   4. Leerzeichen und Trennzeichen → Bindestrich
//   5. Mehrfache Bindestriche → einfacher Bindestrich
//   6. Führende/abschliessende Bindestriche entfernen
//   7. Max. 200 Zeichen (an Wortgrenze abschneiden)
func Generate(input string) string {
    s := strings.ToLower(input)

    for from, to := range umlautReplacements {
        s = strings.ReplaceAll(s, from, strings.ToLower(to))
    }

    s = nonAlphanumeric.ReplaceAllString(s, "-")
    s = multiDash.ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")

    // Max-Länge: an Wortgrenze kürzen
    if len(s) > 200 {
        s = s[:200]
        if idx := strings.LastIndex(s, "-"); idx > 100 {
            s = s[:idx]
        }
    }

    return s
}
```

### 2.2 Beispiele

| Eingabe | Slug |
|---------|------|
| Hydraulikzylinder CDT3 50mm Hub 200mm | `hydraulikzylinder-cdt3-50mm-hub-200mm` |
| Bosch Rexroth | `bosch-rexroth` |
| Türen & Zargen | `tueren-zargen` |
| 3M™ Schutzbrille SecureFit™ 400 | `3m-schutzbrille-securefit-400` |
| Ø 63mm Edelstahl-Rohr | `63mm-edelstahl-rohr` |
| Faltkarton 300×200×150mm braun | `faltkarton-300x200x150mm-braun` |
| LED-Hallenstrahler 150W IP65 | `led-hallenstrahler-150w-ip65` |

### 2.3 Eindeutigkeit

Slugs müssen pro Entitäts-Typ, Tenant und Sprache eindeutig sein.

**Duplikat-Auflösung:**

```go
func (s *service) EnsureUniqueSlug(ctx context.Context, tenantID uuid.UUID, baseSlug string, entityType string, lang string) (string, error) {
    slug := baseSlug
    suffix := 1

    for {
        exists, err := s.repo.SlugExists(ctx, tenantID, slug, entityType, lang)
        if err != nil {
            return "", err
        }
        if !exists {
            return slug, nil
        }
        suffix++
        slug = fmt.Sprintf("%s-%d", baseSlug, suffix)
    }
}
```

**Beispiel:** Zwei Produkte heissen «Hydraulikzylinder 50mm»:
- Erstes Produkt: `hydraulikzylinder-50mm`
- Zweites Produkt: `hydraulikzylinder-50mm-2`

### 2.4 Slug-Feld in Entitäten

Alle öffentlich zugänglichen Entitäten erhalten ein Slug-Feld:

| Entität | Slug-Typ | Beispiel |
|---------|----------|---------|
| **Product** | `map[string]string` (i18n) | `{"de": "hydraulikzylinder-cdt3", "fr": "verin-hydraulique-cdt3"}` |
| **Category** | `string` (pro Sprache in `category_translations`) | Bereits vorhanden (siehe Kategorie-Architektur) |
| **Manufacturer** | `map[string]string` (i18n) | Bereits vorhanden (siehe Hersteller-Entität) |

**Product-Model Erweiterung:**

```go
type Product struct {
    // ... bestehende Felder ...

    Slug map[string]string `json:"slug" db:"slug"` // NEU: {"de": "...", "fr": "...", "en": "..."}
}
```

```sql
-- Migration: add slug to products
ALTER TABLE products ADD COLUMN slug JSONB NOT NULL DEFAULT '{}';
CREATE UNIQUE INDEX idx_products_tenant_slug ON products
    USING GIN (tenant_id, slug) WHERE deleted_at IS NULL;
```

### 2.5 Manuelle Überschreibung

Admins können Slugs manuell überschreiben. Anwendungsfälle:
- Produktname zu lang → kürzerer Slug gewünscht
- Marketing-Kampagne mit spezifischer URL
- SEO-Optimierung mit gezielten Keywords

**Admin-API:**

```json
PUT /api/v1/admin/products/:id
{
    "slug": {
        "de": "hyd-zylinder-cdt3-50",
        "fr": "verin-hyd-cdt3-50"
    }
}
```

**Regeln bei manueller Änderung:**
- Eindeutigkeits-Check wird durchgeführt
- Alter Slug wird in die Redirect-Tabelle eingetragen (→ Abschnitt 5.4)
- `slug_is_manual` Flag wird gesetzt, damit automatische Regeneration bei Namensänderung den manuellen Slug nicht überschreibt

### 2.6 Slug-Änderung bei Namensänderung

| Szenario | Verhalten |
|----------|-----------|
| Produktname ändert sich, Slug ist automatisch generiert | Neuer Slug wird generiert, alter Slug → Redirect-Tabelle |
| Produktname ändert sich, Slug ist manuell gesetzt | Slug bleibt unverändert (manuell hat Vorrang) |
| Admin ändert Slug manuell | Alter Slug → Redirect-Tabelle, neuer Slug aktiv |

---

## 3. Sprach-URLs

### 3.1 Strategie: Pfad-Präfix

**Entscheid: Pfad-Präfix** (`/de/...`, `/fr/...`) statt Subdomain (`de.shop.example.com`).

| Option | Pro | Contra |
|--------|-----|--------|
| **Pfad-Präfix** `/de/...` | Ein Domain, einfaches Hosting, Link-Juice konsolidiert, einfaches Deployment | URL etwas länger |
| **Subdomain** `de.shop...` | Saubere Trennung | Separate DNS, SSL-Certs, Google wertet als separate Sites, Link-Juice verteilt |
| **TLD** `.de` / `.ch` / `.fr` | Geo-Targeting | Komplett separate Infrastruktur, teuer, für B2B überdimensioniert |

**Empfehlung:** Pfad-Präfix. Industriestandard für mehrsprachige B2B-Shops. Google empfiehlt explizit diesen Ansatz.

### 3.2 Sprach-Routing

```
https://shop.example.com/de/produkte/hydraulikzylinder-cdt3
https://shop.example.com/fr/produits/verin-hydraulique-cdt3
https://shop.example.com/en/products/hydraulic-cylinder-cdt3
https://shop.example.com/it/prodotti/cilindro-idraulico-cdt3
```

**Pfadsegment-Übersetzungen:**

| Segment | DE | FR | EN | IT |
|---------|----|----|----|----|
| Produkte | `produkte` | `produits` | `products` | `prodotti` |
| Kategorien | `kategorien` | `categories` | `categories` | `categorie` |
| Hersteller | `hersteller` | `fabricants` | `manufacturers` | `produttori` |
| Suche | `suche` | `recherche` | `search` | `ricerca` |

Diese Übersetzungen werden als **Konfiguration** im Frontend hinterlegt, nicht in der Datenbank.

### 3.3 hreflang Tags

Jede Seite liefert hreflang-Tags für alle verfügbaren Sprachen:

```html
<link rel="alternate" hreflang="de" href="https://shop.example.com/de/produkte/hydraulikzylinder-cdt3" />
<link rel="alternate" hreflang="fr" href="https://shop.example.com/fr/produits/verin-hydraulique-cdt3" />
<link rel="alternate" hreflang="en" href="https://shop.example.com/en/products/hydraulic-cylinder-cdt3" />
<link rel="alternate" hreflang="x-default" href="https://shop.example.com/de/produkte/hydraulikzylinder-cdt3" />
```

**`x-default`** zeigt auf die Standard-Sprache des Tenants (konfigurierbar, typischerweise `de`).

**Implementierung:** Die API liefert pro Entität die Slugs aller Sprachen. Das Frontend baut daraus die hreflang-Links:

```json
// API Response
{
    "slug": {
        "de": "hydraulikzylinder-cdt3",
        "fr": "verin-hydraulique-cdt3",
        "en": "hydraulic-cylinder-cdt3"
    }
}
```

### 3.4 Canonical URLs

Jede Seite hat **genau eine** Canonical URL — immer die kanonische URL in der aktuellen Sprache:

```html
<!-- Auf /de/produkte/hydraulikzylinder-cdt3 -->
<link rel="canonical" href="https://shop.example.com/de/produkte/hydraulikzylinder-cdt3" />

<!-- Auf /de/kategorien/industriebedarf/hydraulik/hydraulikzylinder-cdt3 (Kategorie-Kontext) -->
<link rel="canonical" href="https://shop.example.com/de/produkte/hydraulikzylinder-cdt3" />
```

### 3.5 Redirect bei Sprachwechsel

Wenn ein Nutzer die Sprache wechselt, wird er auf die **äquivalente Seite** in der neuen Sprache weitergeleitet:

```
Nutzer ist auf:  /de/produkte/hydraulikzylinder-cdt3
Wechselt zu FR:  → 302 Redirect → /fr/produits/verin-hydraulique-cdt3
```

**302 (nicht 301)** — weil der Redirect nutzerabhängig ist, nicht permanent.

**Fallback:** Wenn eine Übersetzung fehlt (Slug existiert nicht in Zielsprache):
1. Fallback auf Default-Sprache des Tenants
2. Wenn auch das fehlt: 404

### 3.6 Root-URL und Spracherkennung

```
https://shop.example.com/
```

**Verhalten:**
1. `Accept-Language` Header auswerten
2. Falls Sprache unterstützt → `302 Redirect` auf `/{lang}/`
3. Falls nicht → Redirect auf Default-Sprache (`/de/`)
4. Sprach-Präferenz in Cookie speichern für Folgebesuche

**Wichtig:** Kein `301` für die Root-URL — die Weiterleitung ist nutzerabhängig.

---

## 4. Breadcrumbs

### 4.1 Breadcrumb-Logik

Breadcrumbs basieren auf dem **Navigationskontext** — nicht auf einer festen Zuordnung:

| URL | Breadcrumb |
|-----|------------|
| `/de/produkte/hydraulikzylinder-cdt3` | Home → Produkte → Hydraulikzylinder CDT3 |
| `/de/kategorien/industriebedarf/hydraulik/hydraulikzylinder-cdt3` | Home → Industriebedarf → Hydraulik → Hydraulikzylinder CDT3 |
| `/de/hersteller/nexoflux/hydraulikzylinder-cdt3` | Home → Hersteller → NexoFlux → Hydraulikzylinder CDT3 |
| `/de/kategorien/industriebedarf/hydraulik` | Home → Industriebedarf → Hydraulik |

### 4.2 Produkt gehört zu mehreren Kategorien

Ein Produkt kann mehreren Kategorien zugeordnet sein (z.B. «Hydraulikzylinder» in «Industriebedarf → Hydraulik» und «Antriebstechnik → Linearantriebe»).

**Strategie: Primärkategorie**

Jedes Produkt hat eine **Primärkategorie** (`primary_category_id`). Diese bestimmt:
- Den Default-Breadcrumb auf `/de/produkte/{slug}`
- Die Kategorie in strukturierten Daten
- Den Pfad in der Sitemap

```go
type Product struct {
    // ... bestehende Felder ...
    PrimaryCategoryID *uuid.UUID `json:"primary_category_id" db:"primary_category_id"`
}
```

**Wenn über Kategorie navigiert:** Breadcrumb zeigt den tatsächlichen Navigationskontext (nicht die Primärkategorie).

**Admin-Konfiguration:** Primärkategorie wird beim Erstellen automatisch gesetzt (erste zugewiesene Kategorie) und kann manuell geändert werden.

### 4.3 JSON-LD Strukturierte Daten

```html
<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    "itemListElement": [
        {
            "@type": "ListItem",
            "position": 1,
            "name": "Home",
            "item": "https://shop.example.com/de/"
        },
        {
            "@type": "ListItem",
            "position": 2,
            "name": "Industriebedarf",
            "item": "https://shop.example.com/de/kategorien/industriebedarf"
        },
        {
            "@type": "ListItem",
            "position": 3,
            "name": "Hydraulik",
            "item": "https://shop.example.com/de/kategorien/industriebedarf/hydraulik"
        },
        {
            "@type": "ListItem",
            "position": 4,
            "name": "Hydraulikzylinder CDT3",
            "item": "https://shop.example.com/de/produkte/hydraulikzylinder-cdt3"
        }
    ]
}
</script>
```

### 4.4 API-Unterstützung

Die Breadcrumb-Daten werden vom Backend mitgeliefert. Kategorie-API liefert bereits `breadcrumbs` (siehe [Kategorie-Architektur](category-architecture.md)). Produkt-API wird erweitert:

```json
GET /api/v1/products/by-slug/hydraulikzylinder-cdt3?lang=de

{
    "id": "...",
    "name": {"de": "Hydraulikzylinder CDT3"},
    "slug": {"de": "hydraulikzylinder-cdt3"},
    "breadcrumbs": [
        {"name": "Home", "slug": "", "url": "/de/"},
        {"name": "Industriebedarf", "slug": "industriebedarf", "url": "/de/kategorien/industriebedarf"},
        {"name": "Hydraulik", "slug": "hydraulik", "url": "/de/kategorien/industriebedarf/hydraulik"},
        {"name": "Hydraulikzylinder CDT3", "slug": "hydraulikzylinder-cdt3", "url": "/de/produkte/hydraulikzylinder-cdt3"}
    ]
}
```

---

## 5. Technische Umsetzung

### 5.1 Backend: Slug-Felder in Domain-Models

Alle öffentlich sichtbaren Entitäten erhalten ein `slug`-Feld (soweit nicht bereits vorhanden):

| Entität | Slug-Feld | Status |
|---------|-----------|--------|
| Category | `slug` (in `category_translations`) | ✅ Bereits vorhanden |
| Manufacturer | `slug JSONB` | ✅ Bereits vorhanden |
| Product | `slug JSONB` | ❌ **Neu** — muss ergänzt werden |

**Migration:**

```sql
-- 000010_add_product_slugs.up.sql
ALTER TABLE products ADD COLUMN slug JSONB NOT NULL DEFAULT '{}';
ALTER TABLE products ADD COLUMN slug_is_manual BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE products ADD COLUMN primary_category_id UUID REFERENCES categories(id);

-- Unique Index pro Tenant und Sprache
-- (GIN-Index auf JSONB für Slug-Lookup)
CREATE INDEX idx_products_slug ON products USING GIN (slug jsonb_path_ops)
    WHERE deleted_at IS NULL;

-- Initiale Slug-Generierung aus bestehenden Produktnamen
-- (wird per Migration-Script durchgeführt)
```

### 5.2 Backend: Lookup by Slug

Neue API-Endpoints für Slug-basierte Lookups:

```go
// Repository Interface
type ProductRepository interface {
    // Bestehend
    GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Product, error)

    // Neu
    GetBySlug(ctx context.Context, tenantID uuid.UUID, slug string, lang string) (*Product, error)
}
```

```sql
-- Query: Produkt per Slug finden
SELECT * FROM products
WHERE tenant_id = $1
  AND slug->>$2 = $3    -- $2 = Sprache, $3 = Slug-Wert
  AND deleted_at IS NULL;
```

**API-Endpoints:**

```
GET /api/v1/products/by-slug/:slug?lang=de       → Produkt per Slug
GET /api/v1/categories/by-slug/:slug?lang=de      → Kategorie per Slug (bereits vorhanden)
GET /api/v1/manufacturers/by-slug/:slug?lang=de   → Hersteller per Slug (bereits vorhanden)
```

### 5.3 Routing: Slug-Auflösung

```
Browser                  Next.js Frontend           Gondolia API            PostgreSQL
   │                          │                          │                       │
   │ GET /de/produkte/xyz     │                          │                       │
   │ ─────────────────────▶   │                          │                       │
   │                          │ GET /api/v1/products/    │                       │
   │                          │   by-slug/xyz?lang=de    │                       │
   │                          │ ─────────────────────▶   │                       │
   │                          │                          │ SELECT ... WHERE      │
   │                          │                          │   slug->>'de' = 'xyz' │
   │                          │                          │ ─────────────────────▶│
   │                          │                          │                       │
   │                          │                          │ ◀─────────────────────│
   │                          │ ◀─────────────────────   │                       │
   │ ◀─────────────────────   │                          │                       │
```

**Caching-Strategie:**

| Layer | Cache | TTL | Invalidierung |
|-------|-------|-----|---------------|
| **CDN/Edge** | Seiten-Cache (Vercel/Cloudflare) | 60s (ISR) | On-Demand Revalidation bei Produktänderung |
| **Redis** | Slug → UUID Mapping | 1h | Event-basiert (`product.updated`, `product.deleted`) |
| **DB** | GIN-Index auf `slug` JSONB | — | — |

**Redis Slug-Cache:**

```
slug:product:de:hydraulikzylinder-cdt3 → b0000000-0014-0000-0000-000000000001
slug:category:de:industriebedarf       → a0000000-0000-0000-0000-000000000001
slug:manufacturer:de:nexoflux          → c0000000-0001-0000-0000-000000000001
```

```go
func (s *service) ResolveSlug(ctx context.Context, tenantID uuid.UUID, entityType, lang, slug string) (uuid.UUID, error) {
    // 1. Redis Cache
    cacheKey := fmt.Sprintf("slug:%s:%s:%s:%s", tenantID, entityType, lang, slug)
    if id, err := s.cache.Get(ctx, cacheKey); err == nil {
        return uuid.Parse(id)
    }

    // 2. DB Lookup
    id, err := s.repo.GetIDBySlug(ctx, tenantID, entityType, lang, slug)
    if err != nil {
        return uuid.Nil, err
    }

    // 3. Cache setzen
    s.cache.Set(ctx, cacheKey, id.String(), 1*time.Hour)
    return id, nil
}
```

### 5.4 Redirects: Alte URLs → Neue URLs

#### UUID-URLs → Slug-URLs (Migration)

Alle alten UUID-URLs müssen per **301 Redirect** auf die neuen Slug-URLs weiterleiten:

```
/products/b0000000-0014-0000-0000-000000000001
→ 301 → /de/produkte/hydraulikzylinder-cdt3

/categories/a0000000-0000-0000-0000-000000000001
→ 301 → /de/kategorien/industriebedarf/hydraulik
```

#### Redirect-Tabelle

Für Slug-Änderungen (Umbenennungen) und alte URLs:

```sql
CREATE TABLE url_redirects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    source_path VARCHAR(500) NOT NULL,       -- Alter Pfad (ohne Domain)
    target_path VARCHAR(500) NOT NULL,       -- Neuer Pfad
    status_code SMALLINT NOT NULL DEFAULT 301,
    entity_type VARCHAR(50),                 -- 'product', 'category', 'manufacturer'
    entity_id UUID,                          -- Zugehörige Entität
    hits INTEGER NOT NULL DEFAULT 0,         -- Zähler für Monitoring
    last_hit_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,                  -- Optional: Redirect verfällt

    UNIQUE(tenant_id, source_path)
);

CREATE INDEX idx_redirects_lookup ON url_redirects(tenant_id, source_path)
    WHERE expires_at IS NULL OR expires_at > NOW();
```

#### Redirect-Ketten vermeiden

Wenn Slug A → Slug B → Slug C umbenannt wird, muss die Redirect-Tabelle **direkt** von A → C zeigen (nicht A → B → C). Das wird beim Eintragen neuer Redirects automatisch aufgelöst:

```go
func (s *service) CreateRedirect(ctx context.Context, tenantID uuid.UUID, source, target string) error {
    // Prüfen ob 'source' bereits ein Redirect-Ziel ist
    // Falls ja: alle Redirects die auf 'source' zeigen, auf 'target' umbiegen
    existingRedirects, _ := s.repo.GetRedirectsByTarget(ctx, tenantID, source)
    for _, r := range existingRedirects {
        r.TargetPath = target
        s.repo.Update(ctx, r)
    }

    return s.repo.Create(ctx, &URLRedirect{
        TenantID:   tenantID,
        SourcePath: source,
        TargetPath: target,
        StatusCode: 301,
    })
}
```

#### Initiale UUID-Redirect-Migration

```go
// Einmaliger Migrationsjob: UUID-URLs → Slug-URLs
func MigrateUUIDRedirects(ctx context.Context, tenantID uuid.UUID) error {
    products, _ := productRepo.GetAll(ctx, tenantID)
    for _, p := range products {
        for lang, slug := range p.Slug {
            source := fmt.Sprintf("/products/%s", p.ID)
            target := fmt.Sprintf("/%s/produkte/%s", lang, slug) // TODO: Pfadsegment pro Sprache
            redirectRepo.Create(ctx, &URLRedirect{
                TenantID:   tenantID,
                SourcePath: source,
                TargetPath: target,
                StatusCode: 301,
                EntityType: "product",
                EntityID:   p.ID,
            })
        }
    }
    // Analog für Categories und Manufacturers
    return nil
}
```

### 5.5 Frontend: Next.js Dynamic Routes

```
app/
├── [lang]/
│   ├── layout.tsx
│   ├── page.tsx                              → Startseite
│   ├── produkte/                             → wird per rewrites auf [lang]/[...productPath] gemappt
│   ├── kategorien/
│   │   └── [...categoryPath]/
│   │       └── page.tsx                      → Kategorie oder Kategorie+Produkt
│   ├── hersteller/
│   │   └── [manufacturerSlug]/
│   │       └── page.tsx                      → Hersteller-Detail
│   ├── suche/
│   │   └── page.tsx                          → Suchseite
│   └── p/
│       └── [sku]/
│           └── page.tsx                      → SKU-Redirect
```

**Next.js Middleware für Sprach-Routing:**

```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const SUPPORTED_LANGS = ['de', 'fr', 'en', 'it'];
const DEFAULT_LANG = 'de';

export function middleware(request: NextRequest) {
    const pathname = request.nextUrl.pathname;

    // Root → Sprachwahl
    if (pathname === '/') {
        const preferredLang = getPreferredLanguage(request);
        return NextResponse.redirect(new URL(`/${preferredLang}/`, request.url), 302);
    }

    // UUID-Redirect prüfen (Legacy-URLs)
    if (pathname.match(/\/(?:products|categories)\/[0-9a-f-]{36}/)) {
        // API-Call: Redirect-Tabelle abfragen
        // → 301 auf Slug-URL
    }

    // Sprach-Präfix prüfen
    const langMatch = pathname.match(/^\/([a-z]{2})\//);
    if (!langMatch || !SUPPORTED_LANGS.includes(langMatch[1])) {
        return NextResponse.redirect(new URL(`/${DEFAULT_LANG}${pathname}`, request.url), 302);
    }
}
```

### 5.6 Canonical URLs

**Regeln:**

| Seite | Canonical |
|-------|-----------|
| `/de/produkte/{slug}` | Sich selbst |
| `/de/kategorien/.../` | Sich selbst |
| `/de/kategorien/.../produkte/{slug}` | `/de/produkte/{slug}` |
| `/de/hersteller/.../produkte/{slug}` | `/de/produkte/{slug}` |
| `/de/produkte/{slug}?variante=...` | `/de/produkte/{slug}` (ohne Query-Parameter) |
| `/de/produkte/{slug}?seite=2` | `/de/produkte/{slug}` (Seite 1) |
| `/de/kategorien/.../seite=2` | `/de/kategorien/.../` (Seite 1) |

**Implementierung im Frontend:**

```typescript
// components/SEOHead.tsx
function getCanonicalUrl(pathname: string, entityType: string, slug: string, lang: string): string {
    const base = `https://shop.example.com`;
    const pathSegments: Record<string, Record<string, string>> = {
        product: { de: 'produkte', fr: 'produits', en: 'products' },
        category: { de: 'kategorien', fr: 'categories', en: 'categories' },
        manufacturer: { de: 'hersteller', fr: 'fabricants', en: 'manufacturers' },
    };

    return `${base}/${lang}/${pathSegments[entityType][lang]}/${slug}`;
}
```

---

## 6. Weitere SEO-Themen

### 6.1 Meta-Tags

Jede Seite liefert individuelle Meta-Tags. Die Daten kommen aus den Entitäten (Felder `meta_title`, `meta_description`) mit Fallback-Generierung:

**Produkt-Seite:**

```html
<title>Hydraulikzylinder CDT3 50mm Hub 200mm | NexoFlux | Shop Name</title>
<meta name="description" content="NexoFlux Hydraulikzylinder CDT3 mit 50mm Kolben-Ø und 200mm Hub. Betriebsdruck 160 bar. ✓ B2B-Preise ✓ Schnelle Lieferung ✓ Technische Beratung.">
```

**Fallback-Template (wenn kein manueller Meta-Title gepflegt):**

```
Produkt: "{product_name} | {manufacturer_name} | {shop_name}"
Kategorie: "{category_name} kaufen | {shop_name}"
Hersteller: "{manufacturer_name} Produkte | {shop_name}"
Suche: "Suche: {query} | {shop_name}"
```

**Meta-Description Fallback:** Ersten 155 Zeichen der Produktbeschreibung, sauber am Satzende oder Wort abgeschnitten.

### 6.2 Open Graph / Twitter Cards

```html
<!-- Produkt-Seite -->
<meta property="og:type" content="product" />
<meta property="og:title" content="Hydraulikzylinder CDT3 50mm Hub 200mm" />
<meta property="og:description" content="NexoFlux Hydraulikzylinder CDT3..." />
<meta property="og:image" content="https://cdn.example.com/products/hyd-zyl-cdt3-og.jpg" />
<meta property="og:url" content="https://shop.example.com/de/produkte/hydraulikzylinder-cdt3" />
<meta property="og:site_name" content="Shop Name" />
<meta property="og:locale" content="de_CH" />
<meta property="og:locale:alternate" content="fr_CH" />
<meta property="og:locale:alternate" content="en_US" />

<meta property="product:price:amount" content="485.00" />
<meta property="product:price:currency" content="CHF" />

<!-- Twitter Cards -->
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="Hydraulikzylinder CDT3 50mm Hub 200mm" />
<meta name="twitter:description" content="NexoFlux Hydraulikzylinder CDT3..." />
<meta name="twitter:image" content="https://cdn.example.com/products/hyd-zyl-cdt3-og.jpg" />
```

**OG-Image Generierung:** Für Produkte ohne eigenes OG-Bild wird automatisch ein Bild generiert (Produktbild auf Template mit Logo und Preis). Für Kategorien: Banner-Bild oder generiertes Bild.

### 6.3 Sitemap.xml

**Sitemap-Index:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <sitemap>
        <loc>https://shop.example.com/sitemap-categories-de.xml</loc>
        <lastmod>2026-02-14</lastmod>
    </sitemap>
    <sitemap>
        <loc>https://shop.example.com/sitemap-products-de-1.xml</loc>
        <lastmod>2026-02-14</lastmod>
    </sitemap>
    <sitemap>
        <loc>https://shop.example.com/sitemap-products-de-2.xml</loc>
        <lastmod>2026-02-14</lastmod>
    </sitemap>
    <sitemap>
        <loc>https://shop.example.com/sitemap-manufacturers-de.xml</loc>
        <lastmod>2026-02-14</lastmod>
    </sitemap>
    <!-- Analog für FR, EN, IT -->
</sitemapindex>
```

**Produkt-Sitemap (max. 50'000 URLs pro Datei):**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml">
    <url>
        <loc>https://shop.example.com/de/produkte/hydraulikzylinder-cdt3</loc>
        <lastmod>2026-02-10</lastmod>
        <changefreq>weekly</changefreq>
        <priority>0.8</priority>
        <xhtml:link rel="alternate" hreflang="de" href="https://shop.example.com/de/produkte/hydraulikzylinder-cdt3" />
        <xhtml:link rel="alternate" hreflang="fr" href="https://shop.example.com/fr/produits/verin-hydraulique-cdt3" />
        <xhtml:link rel="alternate" hreflang="en" href="https://shop.example.com/en/products/hydraulic-cylinder-cdt3" />
    </url>
</urlset>
```

**Prioritäten:**

| Seiten-Typ | Priority | changefreq |
|------------|----------|------------|
| Startseite | 1.0 | daily |
| Hauptkategorien (Level 0) | 0.9 | weekly |
| Unterkategorien (Level 1+) | 0.7–0.8 | weekly |
| Produkte | 0.8 | weekly |
| Hersteller | 0.6 | monthly |
| Statische Seiten | 0.3 | monthly |

**Generierung:** Sitemaps werden als **statische Dateien** generiert (Cronjob, 1× pro Nacht) und über CDN ausgeliefert. Bei Änderungen (Produkt erstellt/gelöscht) wird die Regenerierung getriggert (Event-basiert, max. alle 15 Minuten).

### 6.4 robots.txt

```
# robots.txt
User-agent: *

# Erlaubt
Allow: /de/
Allow: /fr/
Allow: /en/
Allow: /it/

# Nicht indexieren
Disallow: /api/
Disallow: /admin/
Disallow: /warenkorb/
Disallow: /konto/
Disallow: /kasse/
Disallow: /*/suche
Disallow: /*?seite=
Disallow: /*?sortierung=
Disallow: /*?variante=

# Sitemaps
Sitemap: https://shop.example.com/sitemap-index.xml

# Crawl-Rate
Crawl-delay: 1
```

**Hinweis:** `Disallow: /*?variante=` verhindert die Indexierung einzelner Varianten-URLs. Die Master-Produktseite wird indexiert, Varianten nicht.

### 6.5 Structured Data (JSON-LD)

#### Product

```html
<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "Product",
    "name": "Hydraulikzylinder CDT3 50mm Hub 200mm",
    "description": "NexoFlux Hydraulikzylinder CDT3 mit 50mm Kolben-Ø und 200mm Hub...",
    "image": [
        "https://cdn.example.com/products/hyd-zyl-cdt3-1.jpg",
        "https://cdn.example.com/products/hyd-zyl-cdt3-2.jpg"
    ],
    "sku": "NXF-CDT3-50-200",
    "mpn": "CDT3-50-200-160",
    "brand": {
        "@type": "Brand",
        "name": "NexoFlux"
    },
    "manufacturer": {
        "@type": "Organization",
        "name": "NexoFlux",
        "url": "https://shop.example.com/de/hersteller/nexoflux"
    },
    "category": "Industriebedarf > Hydraulik > Hydraulikzylinder",
    "url": "https://shop.example.com/de/produkte/hydraulikzylinder-cdt3-50mm-hub-200mm",
    "offers": {
        "@type": "AggregateOffer",
        "priceCurrency": "CHF",
        "lowPrice": "385.00",
        "highPrice": "485.00",
        "offerCount": "4",
        "availability": "https://schema.org/InStock"
    },
    "additionalProperty": [
        {
            "@type": "PropertyValue",
            "name": "Kolben-Ø",
            "value": "50 mm"
        },
        {
            "@type": "PropertyValue",
            "name": "Hub",
            "value": "200 mm"
        },
        {
            "@type": "PropertyValue",
            "name": "Betriebsdruck",
            "value": "160 bar"
        }
    ]
}
</script>
```

#### BreadcrumbList

Siehe Abschnitt 4.3.

#### Organization (Startseite)

```html
<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "Organization",
    "name": "Shop Name GmbH",
    "url": "https://shop.example.com",
    "logo": "https://shop.example.com/logo.svg",
    "contactPoint": {
        "@type": "ContactPoint",
        "telephone": "+41-XX-XXX-XX-XX",
        "contactType": "customer service",
        "availableLanguage": ["German", "French", "English"]
    },
    "sameAs": [
        "https://www.linkedin.com/company/shopname"
    ]
}
</script>
```

### 6.6 Core Web Vitals

| Metrik | Ziel | Massnahmen |
|--------|------|------------|
| **LCP** (Largest Contentful Paint) | < 2.5s | Bilder: WebP/AVIF, CDN, `priority`-Attribut für Hero-Bilder, ISR für Produktseiten |
| **FID/INP** (Interaction to Next Paint) | < 200ms | Varianten-Selektor clientseitig, keine API-Calls bei Achsen-Auswahl |
| **CLS** (Cumulative Layout Shift) | < 0.1 | Bild-Dimensionen immer angeben, Skeleton-Loading für Preise, Fonts per `font-display: swap` |

**Architektur-Entscheide für Performance:**

- **ISR (Incremental Static Regeneration):** Produktseiten werden statisch generiert und bei Änderungen revalidiert. Kein SSR bei jedem Request.
- **Edge-Caching:** Slug-Auflösung wird am Edge gecacht (Vercel Edge Middleware oder Cloudflare Workers).
- **Lazy Loading:** Bilder unterhalb des Folds, Herstellerbeschreibung, Related Products — alles lazy.
- **Preloading:** Kritische API-Aufrufe per `<link rel="preload">` vorladen.

---

## 7. B2B-spezifische SEO-Aspekte

### 7.1 Öffentlich vs. Geschützt

In B2B-Shops ist nicht alles öffentlich. Die Abgrenzung ist SEO-kritisch:

| Bereich | Sichtbarkeit | SEO-Strategie |
|---------|-------------|---------------|
| **Produktkatalog** (Name, Bild, Beschreibung, Spezifikationen) | ✅ Öffentlich | Voll indexierbar, Meta-Tags, Structured Data |
| **Preise** | ⚠️ Konfigurierbar pro Tenant | Siehe 7.2 |
| **Verfügbarkeit/Bestand** | ⚠️ Konfigurierbar | Generisch «Auf Anfrage» oder echte Zahlen |
| **Warenkorb/Bestellung** | 🔒 Login-Pflicht | `noindex`, nicht in Sitemap |
| **Kundenkonto** | 🔒 Login-Pflicht | `noindex`, nicht in Sitemap |
| **Kundenspezifische Preise** | 🔒 Login-Pflicht | Nur nach Login sichtbar |

**robots.txt und Meta-Tags** stellen sicher, dass geschützte Bereiche nicht indexiert werden.

### 7.2 Preisanzeige: SEO vs. B2B-Praxis

| Strategie | SEO-Wirkung | B2B-Praxis | Empfehlung |
|-----------|-------------|------------|------------|
| **Preise öffentlich sichtbar** | ✅ Google kann `Product`-Schema mit Preis anzeigen → Rich Snippets, höhere CTR | ⚠️ Wettbewerber sehen Preise, Kunden erwarten individuelle Konditionen | Für Standardprodukte mit Listenpreis |
| **Preise nur nach Login** | ❌ Kein Preis in Structured Data, keine Rich Snippets | ✅ B2B-Standard für individuelle Preisgestaltung | Für kundenspezifische Preise |
| **Listenpreis öffentlich + individuelle Preise nach Login** | ✅ SEO-Vorteil durch Preis in Structured Data | ✅ Best of both worlds | **Empfohlen** |

**Empfohlene Implementierung:**

```
Nicht eingeloggt → Listenpreis (UVP) anzeigen + "Ihr Preis nach Login"
Eingeloggt       → Kundenspezifischer Preis (aus SAP-Konditionen)
```

**Structured Data:** Immer den Listenpreis im JSON-LD ausgeben (auch wenn eingeloggter Nutzer einen anderen Preis sieht). Google sieht nur den öffentlichen Preis.

### 7.3 Google Merchant Center / B2B-Marktplätze

Gondolia sollte einen **Produkt-Feed** generieren können:

| Plattform | Feed-Format | Relevanz |
|-----------|-------------|----------|
| Google Merchant Center | XML / Content API | Hoch — Google Shopping Anzeigen |
| Mercateo | BMEcat/OpenTrans | Mittel — B2B-Marktplatz DACH |
| Amazon Business | Amazon Flat File | Mittel — wachsender B2B-Kanal |
| Wer liefert was (wlw) | CSV / API | Mittel — B2B-Verzeichnis DACH |

**Feed-Generierung:**

```
GET /api/v1/feeds/google-merchant?tenant_id=...&lang=de
→ XML-Feed mit allen aktiven Produkten, Listenpreisen, Bildern, Kategorien
```

**Felder für Google Merchant:**

| Google-Feld | Quelle in Gondolia |
|-------------|-------------------|
| `id` | Product UUID |
| `title` | `product.name[lang]` |
| `description` | `product.description[lang]` |
| `link` | Kanonische Produkt-URL |
| `image_link` | Erstes Produktbild |
| `price` | Listenpreis (öffentlich) |
| `brand` | `manufacturer.name` |
| `gtin` / `mpn` | Aus Produktattributen |
| `google_product_category` | Mapping Gondolia-Kategorie → Google Taxonomy |
| `availability` | `in_stock` / `out_of_stock` / `preorder` |
| `condition` | `new` (B2B typischerweise immer) |

---

## Priorisierung

| Phase | Aufgabe | Aufwand | Abhängigkeiten |
|-------|---------|---------|----------------|
| **Phase 1** | Slug-Feld in Products, Slug-Generierung, `by-slug` API-Endpoints | 3–4 Tage | — |
| **Phase 2** | URL-Redirect-Tabelle, UUID→Slug Migration, 301 Redirects | 2–3 Tage | Phase 1 |
| **Phase 3** | Next.js Routing mit Slugs, Sprach-Routing, Middleware | 3–4 Tage | Phase 1 |
| **Phase 4** | Canonical URLs, hreflang Tags, Meta-Tags | 2–3 Tage | Phase 3 |
| **Phase 5** | JSON-LD Structured Data (Product, Breadcrumb, Organization) | 2–3 Tage | Phase 3 |
| **Phase 6** | Sitemap-Generierung, robots.txt | 1–2 Tage | Phase 3 |
| **Phase 7** | Open Graph / Twitter Cards, OG-Image Generierung | 1–2 Tage | Phase 3 |
| **Phase 8** | Redis Slug-Cache, Performance-Optimierung | 2–3 Tage | Phase 1 |
| **Phase 9** | Google Merchant Feed, B2B-Marktplatz-Feeds | 3–5 Tage | Phase 4, 5 |

**Gesamtaufwand:** ~4–6 Wochen

---

## Offene Fragen

1. **Sprachen pro Tenant:** Welche Sprachen werden initial unterstützt? (Empfehlung: DE, FR, EN — IT optional als Phase 2)
2. **Preisanzeige-Policy:** Pro Tenant konfigurierbar? Oder globale Entscheidung? (Empfehlung: Tenant-Konfiguration)
3. **Slug-Sprache:** Werden Produkt-Slugs in alle Sprachen übersetzt oder reicht ein Slug pro Produkt? (Empfehlung: Pro Sprache, aber Akeneo/PIM liefert die Übersetzungen)
4. **Kategorie-Tiefe:** Aktuell 4 Ebenen empfohlen — gibt es Tenants die tiefere Hierarchien brauchen?
5. **SKU-Kurz-URLs:** Sollen `/p/{sku}`-URLs implementiert werden? (Empfehlung: Ja, als Convenience-Feature für B2B-Kunden die SKUs kennen)
6. **Google Merchant:** Haben aktuelle Tenants bereits ein Google Merchant Center Konto? Welche Feed-Formate werden gebraucht?
7. **Bestehende externe Links:** Gibt es bekannte Backlinks oder Bookmarks auf aktuelle UUID-URLs die redirected werden müssen? (Monitoring nach Go-Live wichtig)
