# Asset-Management-System

Stand: 2026-02-14

## Übersicht

Im B2B E-Commerce sind **Dokumente genauso wichtig wie Bilder**. Ein Hydraulikzylinder braucht nicht nur ein Produktfoto, sondern auch ein Datenblatt, eine CAD-Zeichnung und ein CE-Zertifikat. Ein Industriereiniger ist ohne Sicherheitsdatenblatt schlicht nicht verkaufbar — das ist gesetzlich vorgeschrieben.

Gondolia braucht deshalb ein **generisches Asset-Management-System**, das sowohl Bilder als auch Dokumente verwaltet und diese flexibel mit Produkten, Kategorien und Herstellern verknüpft.

### Aktueller Ist-Zustand

Das heutige `ProductImage`-Struct ist minimal:

```go
type ProductImage struct {
    URL       string `json:"url"`
    AltText   string `json:"alt_text,omitempty"`
    SortOrder int    `json:"sort_order"`
    IsPrimary bool   `json:"is_primary"`
}
```

**Probleme:**
- Keine eigenständige Entität — Bilder sind direkt am Produkt embedded
- Kein Konzept für Dokumente
- Keine Wiederverwendung (ein Hersteller-Zertifikat müsste bei jedem Produkt separat hochgeladen werden)
- Keine Sprach-Zuordnung (DE-Datenblatt vs. FR-Datenblatt)
- Keine automatische Bildverarbeitung (Thumbnails, WebP-Konvertierung)
- Kein Storage-Abstraction-Layer

### Ziel-Architektur

```
Asset (generische Entität)
├── Image (Produktbilder, Kategorie-Bilder, Hersteller-Logos)
│   ├── Automatische Thumbnail-Generierung
│   ├── WebP/AVIF-Konvertierung
│   └── Responsive Varianten (Thumbnail, Medium, Full)
└── Document (PDFs, CAD-Dateien, Zertifikate)
    ├── Typ-Klassifizierung (Datenblatt, SDB, Zertifikat etc.)
    ├── Sprach-Zuordnung
    └── Versionierung
```

---

## 1. Produkt-Bilder (Images)

### Galerie & Varianten

Produkte können **mehrere Bilder** haben, sortierbar mit einem Hauptbild (Primary). Bei Variantenprodukten (vgl. [Produkttypen-Architektur](product-types.md)) existieren Bilder auf zwei Ebenen:

| Ebene | Beschreibung | Beispiel |
|-------|-------------|---------|
| **Master-Ebene** | Allgemeine Produktbilder, gelten für alle Varianten | Gesamtansicht eines Hydraulikzylinders |
| **Varianten-Ebene** | Spezifische Bilder pro Variante | Netzwerkkabel in Grau vs. Blau vs. Rot |

**Verhalten bei Variantenauswahl:**
- Variante hat eigene Bilder → diese anzeigen, Master-Bilder als Fallback in der Galerie ergänzen
- Variante hat keine eigenen Bilder → Master-Bilder anzeigen
- Bild-Wechsel bei Achsenauswahl (z.B. Farbe) → sofortiger Austausch ohne Page-Reload

### Bildformate & Responsive Varianten

Pro hochgeladenem Bild werden automatisch **mehrere Varianten** generiert:

| Variante | Max. Breite | Einsatz | Format |
|----------|-------------|---------|--------|
| `thumbnail` | 150px | Produktliste, Warenkorb, Varianten-Thumbnails | WebP + JPEG Fallback |
| `small` | 300px | Suchergebnisse, Related Products | WebP + JPEG Fallback |
| `medium` | 600px | Produktdetailseite (Standardansicht) | WebP + JPEG Fallback |
| `large` | 1200px | Zoom-Ansicht, Lightbox | WebP + AVIF + JPEG Fallback |
| `original` | unverändert | Download, Print | Originalformat |

**Format-Strategie:**
- **WebP** als primäres Auslieferungsformat (90%+ Browser-Support, 25–35% kleiner als JPEG)
- **AVIF** als progressives Format für `large`-Variante (50% kleiner als JPEG, wachsender Support)
- **JPEG** als Fallback für ältere Browser
- Content-Negotiation über `Accept`-Header oder `<picture>`-Element im Frontend

### Frontend-Anforderungen (Bilder)

| Feature | Beschreibung |
|---------|-------------|
| **Galerie** | Swipeable Bildergalerie, Thumbnails unterhalb/seitlich, aktives Bild hervorgehoben |
| **Zoom** | Hover-Zoom auf Desktop, Pinch-to-Zoom auf Mobile |
| **Lightbox** | Vollbild-Ansicht mit Navigation (Pfeiltasten, Swipe) |
| **Lazy Loading** | Bilder unterhalb des Viewports erst bei Annäherung laden (`loading="lazy"`, Intersection Observer) |
| **Alt-Text** | Mehrsprachiger Alt-Text pro Bild für Accessibility und SEO |
| **Placeholder** | Low-Quality Image Placeholder (LQIP) oder BlurHash während Bilder laden |

---

## 2. Produkt-Dokumente (Documents)

### Dokumenttypen im B2B-Kontext

| Dokumenttyp | Code | Beschreibung | Branchenbeispiel | Gesetzlich? |
|-------------|------|-------------|-------------------|-------------|
| **Datenblatt** | `datasheet` | Technische Spezifikationen, Leistungsdaten | Hydraulikzylinder: Kolbendurchmesser, Hübe, Druckbereiche | Nein |
| **Sicherheitsdatenblatt** | `safety_datasheet` | SDS/MSDS nach REACH-Verordnung (EG) 1907/2006 | Industriereiniger: Gefahrenstoffklassifizierung, Erste-Hilfe-Massnahmen | **Ja — gesetzlich vorgeschrieben!** |
| **Einbauanweisung** | `installation_guide` | Montage- und Installationsanleitung | Schaltschrank: Befestigungspunkte, Kabeleinführungen, Schutzabstände | Teilweise (CE) |
| **Bedienungsanleitung** | `operating_manual` | Betriebs- und Wartungsanleitung | Laborgeräte: Inbetriebnahme, Kalibrierung, Wartungsintervalle | Teilweise |
| **Zertifikat** | `certificate` | CE, TÜV, ISO, RoHS, REACH, UL, ATEX etc. | Sicherheitsschuh: EN ISO 20345:2022, CE-Kennzeichnung | Branchenabhängig |
| **Konformitätserklärung** | `declaration_of_conformity` | EU-Konformitätserklärung nach Maschinenrichtlinie | Frequenzumrichter: EU DoC gemäss 2014/35/EU | **Ja — gesetzlich vorgeschrieben** |
| **Prüfbericht** | `test_report` | Ergebnisse von Material-/Funktionsprüfungen | Schutzhandschuh: Abriebfestigkeit, Schnittfestigkeit nach EN 388 | Nein |
| **CAD-Zeichnung** | `cad_drawing` | 2D/3D-Konstruktionszeichnungen | Pneumatikzylinder: STEP-Datei für Integration in Kundenkonstruktion | Nein |
| **Produktbroschüre** | `brochure` | Marketing-/Informationsmaterial | Kochblock-Programm: Übersicht Modelle, Materialien, Referenzprojekte | Nein |
| **Pflegeanleitung** | `care_instruction` | Reinigungs- und Pflegehinweise | Arbeitskleidung: Waschtemperatur, Imprägnierungserneuerung | Nein |

### Dateiformate

| Format | MIME-Type | Einsatz |
|--------|-----------|---------|
| PDF | `application/pdf` | Datenblätter, SDB, Zertifikate, Anleitungen — **Standard** |
| DWG | `application/acad` | AutoCAD 2D-Zeichnungen |
| STEP | `application/step` | 3D-CAD-Modelle (ISO 10303) |
| IGES | `model/iges` | 3D-CAD-Modelle (älteres Format) |
| DXF | `application/dxf` | CAD-Austauschformat |

### Sprach-Zuordnung

Dokumente sind **sprachspezifisch**. Ein Sicherheitsdatenblatt muss in der Landessprache des Kunden vorliegen — das ist gesetzlich vorgeschrieben.

```
Produkt: Industriereiniger TechClean 5000
├── Sicherheitsdatenblatt DE → sdb_techclean5000_de.pdf
├── Sicherheitsdatenblatt FR → sdb_techclean5000_fr.pdf
├── Sicherheitsdatenblatt IT → sdb_techclean5000_it.pdf
├── Datenblatt DE → tds_techclean5000_de.pdf
├── Datenblatt EN → tds_techclean5000_en.pdf
└── CE-Zertifikat (sprachunabhängig) → ce_techclean5000.pdf
```

Im Frontend wird automatisch das Dokument in der Sprache des Benutzers angezeigt. Fallback: Deutsch → Englisch → erstbestes.

### Versionierung

Dokumente werden **versioniert**. Wenn ein neues Sicherheitsdatenblatt hochgeladen wird, ersetzt es die vorherige Version — die alte Version bleibt aber archiviert:

| Feld | Beschreibung |
|------|-------------|
| `version` | Versionsnummer (1, 2, 3, ...) — automatisch hochgezählt |
| `is_current` | Flag: Ist dies die aktuelle Version? |
| `superseded_by` | Verweis auf die neue Version |
| `valid_from` | Gültig ab (optional, z.B. bei regulatorischen Dokumenten) |

**Warum Versionierung wichtig ist:**
- SDB-Revisionen sind gesetzlich nachweispflichtig
- Kunden, die Version 3 des Datenblatts heruntergeladen haben, sollen nachvollziehbar den Stand kennen
- Bei Reklamationen muss die zum Kaufzeitpunkt gültige Dokumentversion abrufbar sein

---

## 3. Dokument-Viewer

### PDF-Viewer im Browser

Dokumente sollen **inline im Browser** anzeigbar sein — nicht nur als Download. B2B-Kunden wollen schnell in ein Datenblatt reinschauen, ohne es erst herunterladen zu müssen.

**Technologie-Entscheid:**

| Option | Vorteile | Nachteile | Empfehlung |
|--------|----------|-----------|------------|
| **pdf.js** (Mozilla) | Open Source, volle Kontrolle, Custom UI, kein Plugin nötig | Bundle-Grösse (~400KB gzipped), Rendering-Performance bei grossen PDFs | ✅ **Empfohlen** |
| **Browser-native** (`<embed>`, `<iframe>`) | Kein zusätzlicher Code, performant | Inkonsistentes UI über Browser, kein Custom Styling, Mobile-Support schlecht | Fallback |
| **Google Docs Viewer** | Einfach einzubinden | Datenschutz (Dokument wird an Google gesendet!), B2B-Kunden akzeptieren das nicht | ❌ Ausgeschlossen |

**Empfehlung:** **pdf.js** als primärer Viewer mit Browser-nativem Fallback. Konfigurierbar pro Tenant (manche Kunden wollen einfach nur Download).

### Viewer-Features

| Feature | Beschreibung |
|---------|-------------|
| **Inline-Anzeige** | PDF öffnet sich in einem Modal/Overlay auf der Produktseite |
| **Seiten-Navigation** | Blättern, Seitensprung, Seitenanzahl |
| **Zoom** | Rein-/Rauszoomen, Fit-to-Width, Fit-to-Page |
| **Download-Button** | Direkter Download des Originals |
| **Drucken** | Print-Funktion aus dem Viewer |
| **Mobile** | Touch-optimiert, Pinch-to-Zoom |
| **Fullscreen** | Vollbildmodus für grosse Dokumente |

### Preview-Thumbnails

Für jedes PDF-Dokument wird beim Upload automatisch ein **Preview-Thumbnail** der ersten Seite generiert. Dieses wird in der Dokumentliste auf der Produktseite angezeigt.

---

## 4. Asset-Entität / Datenmodell

### Generische Asset-Entität

Ein **Asset** ist die zentrale Entität für alle Dateien im System — Bilder und Dokumente gleichermassen.

### SQL-Schema

```sql
-- Migration: 000005_create_assets.up.sql

-- ============================================
-- Assets (Haupttabelle)
-- ============================================
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,

    -- Typ & Klassifizierung
    asset_type VARCHAR(20) NOT NULL,          -- 'image', 'document'
    document_type VARCHAR(50),                -- Nur bei asset_type='document': 'datasheet', 'safety_datasheet', 'certificate' etc.
    mime_type VARCHAR(100) NOT NULL,           -- 'image/jpeg', 'application/pdf', 'application/step'
    
    -- Datei-Informationen
    filename VARCHAR(500) NOT NULL,           -- Originaler Dateiname
    file_size BIGINT NOT NULL,                -- Dateigrösse in Bytes
    storage_path VARCHAR(1000) NOT NULL,      -- Pfad im Object Storage (S3-Key)
    
    -- Metadaten
    title JSONB,                              -- i18n: {"de": "Datenblatt Hydraulikzylinder CDT3", "fr": "Fiche technique..."}
    alt_text JSONB,                           -- i18n: Nur für Bilder relevant (Accessibility/SEO)
    description JSONB,                        -- i18n: Optionale Beschreibung
    language VARCHAR(5),                      -- ISO 639-1: 'de', 'fr', 'en' — Dokumentsprache (NULL = sprachunabhängig)
    
    -- Bild-spezifisch
    width INTEGER,                            -- Bildbreite in Pixel (NULL bei Dokumenten)
    height INTEGER,                           -- Bildhöhe in Pixel
    
    -- Versionierung (nur Dokumente)
    version INTEGER NOT NULL DEFAULT 1,
    is_current BOOLEAN NOT NULL DEFAULT true,
    superseded_by UUID REFERENCES assets(id),
    valid_from TIMESTAMPTZ,
    
    -- Verarbeitungsstatus
    processing_status VARCHAR(20) NOT NULL DEFAULT 'pending',  -- 'pending', 'processing', 'ready', 'failed'
    
    -- Externe Systeme
    pim_identifier VARCHAR(255),
    erp_identifier VARCHAR(255),
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indizes
CREATE INDEX idx_assets_tenant_type ON assets(tenant_id, asset_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_tenant_doctype ON assets(tenant_id, document_type) WHERE document_type IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_assets_processing ON assets(processing_status) WHERE processing_status != 'ready';
CREATE INDEX idx_assets_pim ON assets(tenant_id, pim_identifier) WHERE pim_identifier IS NOT NULL;
CREATE INDEX idx_assets_version ON assets(superseded_by) WHERE superseded_by IS NOT NULL;

-- ============================================
-- Asset-Varianten (generierte Bildformate)
-- ============================================
CREATE TABLE asset_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    variant_key VARCHAR(50) NOT NULL,         -- 'thumbnail', 'small', 'medium', 'large'
    format VARCHAR(10) NOT NULL,              -- 'webp', 'avif', 'jpeg'
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    file_size BIGINT NOT NULL,
    storage_path VARCHAR(1000) NOT NULL,
    
    UNIQUE(asset_id, variant_key, format)
);

CREATE INDEX idx_asset_variants_asset ON asset_variants(asset_id);

-- ============================================
-- Asset-Zuordnungen (M:N Beziehungen)
-- ============================================
CREATE TABLE asset_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    
    -- Polymorphe Zuordnung
    entity_type VARCHAR(30) NOT NULL,         -- 'product', 'product_variant', 'category', 'manufacturer'
    entity_id UUID NOT NULL,
    
    -- Kontext
    role VARCHAR(30) NOT NULL DEFAULT 'gallery',  -- 'gallery', 'primary', 'logo', 'icon', 'document'
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(asset_id, entity_type, entity_id, role)
);

CREATE INDEX idx_asset_assignments_entity ON asset_assignments(entity_type, entity_id);
CREATE INDEX idx_asset_assignments_asset ON asset_assignments(asset_id);
```

```sql
-- Migration: 000005_create_assets.down.sql

DROP TABLE IF EXISTS asset_assignments;
DROP TABLE IF EXISTS asset_variants;
DROP TABLE IF EXISTS assets;
```

### Beziehungsmodell

```
Asset M ──── N Product          (über asset_assignments)
Asset M ──── N ProductVariant   (über asset_assignments)
Asset M ──── N Category         (über asset_assignments)
Asset M ──── N Manufacturer     (über asset_assignments)
```

**M:N-Beziehung ist entscheidend.** Beispiele:

| Szenario | Beschreibung |
|----------|-------------|
| Hersteller-Zertifikat | Ein ISO 9001-Zertifikat von Festo hängt an **allen** Festo-Produkten — und am Hersteller selbst |
| Kategorie-Bild | Ein Kategoriebild «Hydraulik» wird auch als Fallback-Bild für Produkte ohne eigenes Bild verwendet |
| Gemeinsames Datenblatt | Ein Datenblatt «Pneumatikzylinder-Serie DSBC» gilt für 20 Varianten |
| Hersteller-Logo | Das Festo-Logo hängt am Manufacturer, wird aber auf der Produktseite angezeigt |

### Go Domain Structs

```go
// internal/asset/models.go
package asset

import (
    "time"
    "github.com/google/uuid"
)

// AssetType unterscheidet die Hauptkategorien
type AssetType string

const (
    AssetTypeImage    AssetType = "image"
    AssetTypeDocument AssetType = "document"
)

// DocumentType klassifiziert Dokumente
type DocumentType string

const (
    DocumentTypeDatasheet              DocumentType = "datasheet"
    DocumentTypeSafetyDatasheet        DocumentType = "safety_datasheet"
    DocumentTypeInstallationGuide      DocumentType = "installation_guide"
    DocumentTypeOperatingManual        DocumentType = "operating_manual"
    DocumentTypeCertificate            DocumentType = "certificate"
    DocumentTypeDeclarationOfConformity DocumentType = "declaration_of_conformity"
    DocumentTypeTestReport             DocumentType = "test_report"
    DocumentTypeCADDrawing             DocumentType = "cad_drawing"
    DocumentTypeBrochure               DocumentType = "brochure"
    DocumentTypeCareInstruction        DocumentType = "care_instruction"
)

// ProcessingStatus zeigt den Verarbeitungsstand
type ProcessingStatus string

const (
    ProcessingStatusPending    ProcessingStatus = "pending"
    ProcessingStatusProcessing ProcessingStatus = "processing"
    ProcessingStatusReady      ProcessingStatus = "ready"
    ProcessingStatusFailed     ProcessingStatus = "failed"
)

// EntityType für polymorphe Zuordnungen
type EntityType string

const (
    EntityTypeProduct        EntityType = "product"
    EntityTypeProductVariant EntityType = "product_variant"
    EntityTypeCategory       EntityType = "category"
    EntityTypeManufacturer   EntityType = "manufacturer"
)

// AssignmentRole definiert den Kontext einer Zuordnung
type AssignmentRole string

const (
    AssignmentRoleGallery  AssignmentRole = "gallery"
    AssignmentRolePrimary  AssignmentRole = "primary"
    AssignmentRoleLogo     AssignmentRole = "logo"
    AssignmentRoleIcon     AssignmentRole = "icon"
    AssignmentRoleDocument AssignmentRole = "document"
)

// Asset ist die zentrale Entität
type Asset struct {
    ID               uuid.UUID         `json:"id" db:"id"`
    TenantID         uuid.UUID         `json:"tenant_id" db:"tenant_id"`
    AssetType        AssetType         `json:"asset_type" db:"asset_type"`
    DocumentType     *DocumentType     `json:"document_type,omitempty" db:"document_type"`
    MimeType         string            `json:"mime_type" db:"mime_type"`
    Filename         string            `json:"filename" db:"filename"`
    FileSize         int64             `json:"file_size" db:"file_size"`
    StoragePath      string            `json:"-" db:"storage_path"` // Nie an Frontend ausliefern
    Title            map[string]string `json:"title,omitempty" db:"title"`
    AltText          map[string]string `json:"alt_text,omitempty" db:"alt_text"`
    Description      map[string]string `json:"description,omitempty" db:"description"`
    Language         *string           `json:"language,omitempty" db:"language"`
    Width            *int              `json:"width,omitempty" db:"width"`
    Height           *int              `json:"height,omitempty" db:"height"`
    Version          int               `json:"version" db:"version"`
    IsCurrent        bool              `json:"is_current" db:"is_current"`
    SupersededBy     *uuid.UUID        `json:"superseded_by,omitempty" db:"superseded_by"`
    ValidFrom        *time.Time        `json:"valid_from,omitempty" db:"valid_from"`
    ProcessingStatus ProcessingStatus  `json:"processing_status" db:"processing_status"`
    PIMIdentifier    *string           `json:"pim_identifier,omitempty" db:"pim_identifier"`
    ERPIdentifier    *string           `json:"erp_identifier,omitempty" db:"erp_identifier"`
    CreatedAt        time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
    DeletedAt        *time.Time        `json:"-" db:"deleted_at"`

    // Berechnete Felder
    URLs     *AssetURLs     `json:"urls,omitempty" db:"-"`
    Variants []AssetVariant `json:"variants,omitempty" db:"-"`
}

// AssetURLs enthält die öffentlichen URLs (über CDN)
type AssetURLs struct {
    Original  string `json:"original"`
    Thumbnail string `json:"thumbnail,omitempty"` // Nur bei Bildern
    Small     string `json:"small,omitempty"`
    Medium    string `json:"medium,omitempty"`
    Large     string `json:"large,omitempty"`
}

// AssetVariant ist eine generierte Bildvariante
type AssetVariant struct {
    ID          uuid.UUID `json:"id" db:"id"`
    AssetID     uuid.UUID `json:"asset_id" db:"asset_id"`
    VariantKey  string    `json:"variant_key" db:"variant_key"`
    Format      string    `json:"format" db:"format"`
    Width       int       `json:"width" db:"width"`
    Height      int       `json:"height" db:"height"`
    FileSize    int64     `json:"file_size" db:"file_size"`
    StoragePath string    `json:"-" db:"storage_path"`
}

// AssetAssignment verknüpft ein Asset mit einer Entität
type AssetAssignment struct {
    ID         uuid.UUID      `json:"id" db:"id"`
    AssetID    uuid.UUID      `json:"asset_id" db:"asset_id"`
    EntityType EntityType     `json:"entity_type" db:"entity_type"`
    EntityID   uuid.UUID      `json:"entity_id" db:"entity_id"`
    Role       AssignmentRole `json:"role" db:"role"`
    SortOrder  int            `json:"sort_order" db:"sort_order"`
    IsPrimary  bool           `json:"is_primary" db:"is_primary"`
    CreatedAt  time.Time      `json:"created_at" db:"created_at"`
}
```

### Erweiterung des Product-Structs

Das bisherige `Images []ProductImage` wird durch die Asset-Beziehung ersetzt:

```go
// internal/domain/product.go — Erweiterung

type Product struct {
    // ... bestehende Felder ...

    // NEU: Assets statt Images
    // Images []ProductImage  ← ENTFÄLLT (Migration)
    Assets []AssetAssignment `json:"assets,omitempty" db:"-"` // Lazy-loaded über asset_assignments
}
```

**Migration:** Die bestehenden `ProductImage`-Einträge werden in `assets` + `asset_assignments` migriert. Eine Datenmigration überführt URL → Storage-Path, AltText → alt_text, SortOrder → sort_order, IsPrimary → is_primary/role.

---

## 5. Storage

### Storage-Architektur

```
┌─────────────┐     ┌──────────────┐     ┌───────────┐     ┌─────────┐
│ Admin-Upload │────▶│ Asset Service │────▶│ Object    │────▶│ CDN     │
│ PIM-Import   │     │ (Processing) │     │ Storage   │     │ (Edge)  │
│ ERP-Import   │     └──────────────┘     │ (S3/MinIO)│     └─────────┘
└─────────────┘            │              └───────────┘          │
                           │                                     │
                    ┌──────▼──────┐                        ┌─────▼─────┐
                    │ Processing  │                        │ Frontend  │
                    │ Pipeline    │                        │ (Browser) │
                    └─────────────┘                        └───────────┘
```

### Object Storage

| Option | Einsatz | Konfiguration |
|--------|---------|--------------|
| **MinIO** | Self-Hosted / On-Premise | S3-kompatible API, Docker-Deployment |
| **AWS S3** | Cloud / AWS | Native S3 |
| **Azure Blob** | Cloud / Azure | S3-Kompatibilitätslayer oder native SDK |

**Empfehlung:** S3-kompatible Abstraction (`github.com/aws/aws-sdk-go-v2`), damit MinIO und AWS S3 austauschbar sind. Kein Vendor Lock-in.

### Bucket-Struktur

```
gondolia-assets/
  {tenant_id}/
    images/
      {asset_id}/
        original.jpg
        thumbnail.webp
        thumbnail.jpg
        small.webp
        small.jpg
        medium.webp
        medium.jpg
        large.webp
        large.avif
        large.jpg
    documents/
      {asset_id}/
        original.pdf
        preview.jpg          # Thumbnail der ersten Seite
```

### CDN-Anbindung

- **CloudFront** (AWS) oder **Cloudflare** als CDN
- Assets werden über CDN-URLs ausgeliefert, nie direkt vom Object Storage
- **Signed URLs** mit Ablaufzeit für nicht-öffentliche Dokumente (z.B. kundenspezifische Preislisten)
- **Cache-Invalidation** bei Asset-Update über CDN-API
- URL-Pattern: `https://assets.{tenant-domain}/images/{asset_id}/medium.webp`

### Upload-Pipeline

```
Upload → Validierung → Virus-Scan → Speicherung → Processing → Ready
```

| Schritt | Beschreibung | Technologie |
|---------|-------------|-------------|
| **Validierung** | MIME-Type prüfen, Dateigrösse prüfen, Dateiendung prüfen | Go (net/http.DetectContentType) |
| **Virus-Scan** | Datei auf Malware scannen | ClamAV (Open Source) via clamd-Socket |
| **Speicherung** | Original in Object Storage ablegen | S3 PutObject |
| **Processing (Bilder)** | Thumbnails generieren, WebP/AVIF konvertieren | libvips (via govips oder Kommandozeile) |
| **Processing (PDFs)** | Preview-Thumbnail der ersten Seite | ImageMagick / Ghostscript |
| **Ready** | `processing_status` auf `ready` setzen, Event publizieren | — |

**Asynchrone Verarbeitung:** Upload liefert sofort eine Response mit `processing_status: "pending"`. Die Bildverarbeitung läuft asynchron über eine Message Queue (NATS/RabbitMQ). Das Frontend zeigt einen Platzhalter bis `processing_status: "ready"`.

### Grössen-Limits

| Asset-Typ | Max. Dateigrösse | Erlaubte MIME-Types |
|-----------|-----------------|---------------------|
| Image | 20 MB | `image/jpeg`, `image/png`, `image/webp`, `image/svg+xml` |
| Document (PDF) | 50 MB | `application/pdf` |
| Document (CAD) | 100 MB | `application/step`, `application/acad`, `model/iges`, `application/dxf` |

---

## 6. Frontend-UX

### Produktdetailseite

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  ┌─────────────────────────┐   Hydraulikzylinder CDT3           │
│  │                         │   Hersteller: [Logo] Bosch Rexroth │
│  │      [Hauptbild]        │                                    │
│  │      (Medium/Large)     │   Kolben-Ø: [40] [50] [63] [80]mm │
│  │                         │   Hub:      [100] [200] [300] 500  │
│  │                         │   Druck:    [160] [250] bar        │
│  ├──┬──┬──┬──┬─────────────┤                                    │
│  │T1│T2│T3│T4│             │   CHF 1'245.00 (netto)             │
│  └──┴──┴──┴──┴─────────────┘   [In den Warenkorb]              │
│                                                                 │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                                 │
│  [Beschreibung] [Technische Daten] [📄 Dokumente (5)]          │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ 📄 Dokumente                               [⬇ Alle ZIP] │   │
│  │                                                          │   │
│  │ ┌────┐  Datenblatt CDT3              PDF · 2.4 MB · DE  │   │
│  │ │prev│  Technische Spezifikationen   [👁 Ansehen] [⬇]   │   │
│  │ └────┘                                                   │   │
│  │ ┌────┐  Sicherheitsdatenblatt        PDF · 1.8 MB · DE  │   │
│  │ │prev│  SDS gemäss REACH 1907/2006   [👁 Ansehen] [⬇]   │   │
│  │ └────┘                                                   │   │
│  │ ┌────┐  CE-Konformitätserklärung     PDF · 0.3 MB       │   │
│  │ │prev│  EU DoC 2006/42/EG            [👁 Ansehen] [⬇]   │   │
│  │ └────┘                                                   │   │
│  │ ┌────┐  CAD-Modell CDT3-63           STEP · 8.2 MB      │   │
│  │ │ 3D │  3D-Modell für Konstruktion   [⬇ Download]       │   │
│  │ └────┘                                                   │   │
│  │ ┌────┐  ISO 9001 Zertifikat          PDF · 0.5 MB       │   │
│  │ │prev│  Gültig bis 2027-03           [👁 Ansehen] [⬇]   │   │
│  │ └────┘                                                   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Dokument-Liste

| Element | Beschreibung |
|---------|-------------|
| **Preview-Thumbnail** | Automatisch generiertes Vorschaubild der ersten PDF-Seite |
| **Dokumenttyp-Icon** | Visuelles Icon je nach Typ (📊 Datenblatt, ⚠️ SDB, 🏅 Zertifikat, 📐 CAD) |
| **Titel** | Dokumenttitel (mehrsprachig) |
| **Metadaten** | Format, Dateigrösse, Sprache |
| **Aktionen** | «Ansehen» (öffnet PDF-Viewer), «Download» (direkter Download) |
| **ZIP-Download** | Button «Alle Dokumente herunterladen» — generiert on-the-fly ein ZIP-Archiv |

### ZIP-Download

```
POST /api/v1/products/:id/documents/download-all
→ Response: application/zip (Streaming)
```

Das ZIP wird **on-the-fly** generiert (Streaming, kein Zwischenspeichern), um Speicher zu sparen. Enthält alle aktuellen Dokumente des Produkts in der Sprache des Benutzers.

---

## 7. API-Design

### Public API (Storefront)

#### Assets eines Produkts abrufen

```
GET /api/v1/products/:id/assets?type=image
GET /api/v1/products/:id/assets?type=document&language=de
GET /api/v1/products/:id/assets?document_type=safety_datasheet
```

**Response (Bilder):**

```json
{
    "data": [
        {
            "id": "asset-uuid-1",
            "asset_type": "image",
            "title": {"de": "Hydraulikzylinder CDT3 — Frontalansicht"},
            "alt_text": {"de": "Bosch Rexroth Hydraulikzylinder CDT3, Kolben-Ø 63mm"},
            "is_primary": true,
            "sort_order": 0,
            "urls": {
                "thumbnail": "https://assets.example.com/img/uuid-1/thumbnail.webp",
                "small": "https://assets.example.com/img/uuid-1/small.webp",
                "medium": "https://assets.example.com/img/uuid-1/medium.webp",
                "large": "https://assets.example.com/img/uuid-1/large.webp",
                "original": "https://assets.example.com/img/uuid-1/original.jpg"
            },
            "width": 2400,
            "height": 1600
        }
    ]
}
```

**Response (Dokumente):**

```json
{
    "data": [
        {
            "id": "asset-uuid-10",
            "asset_type": "document",
            "document_type": "safety_datasheet",
            "title": {"de": "Sicherheitsdatenblatt TechClean 5000"},
            "filename": "sdb_techclean5000_de.pdf",
            "mime_type": "application/pdf",
            "file_size": 1887436,
            "language": "de",
            "version": 3,
            "valid_from": "2025-11-01T00:00:00Z",
            "urls": {
                "original": "https://assets.example.com/doc/uuid-10/original.pdf",
                "thumbnail": "https://assets.example.com/doc/uuid-10/preview.jpg"
            }
        }
    ]
}
```

#### Dokumente ZIP-Download

```
POST /api/v1/products/:id/documents/download-all?language=de
→ Content-Type: application/zip
→ Content-Disposition: attachment; filename="hydraulikzylinder-cdt3-dokumente.zip"
```

### Admin API

#### Asset-Upload

```
POST /api/v1/admin/assets
Content-Type: multipart/form-data

file: <binary>
tenant_id: uuid
asset_type: image | document
document_type: datasheet (optional)
language: de (optional)
title: {"de": "Datenblatt CDT3"} (optional)
alt_text: {"de": "..."} (optional)
```

**Response:**

```json
{
    "id": "new-asset-uuid",
    "processing_status": "pending",
    "filename": "datenblatt_cdt3.pdf",
    "file_size": 2456789,
    "mime_type": "application/pdf"
}
```

#### Bulk-Upload

```
POST /api/v1/admin/assets/bulk
Content-Type: multipart/form-data

files[]: <binary> (mehrere Dateien)
metadata: [
    {"filename": "bild1.jpg", "asset_type": "image", "title": {"de": "Frontansicht"}},
    {"filename": "bild2.jpg", "asset_type": "image", "title": {"de": "Seitenansicht"}},
    {"filename": "datenblatt.pdf", "asset_type": "document", "document_type": "datasheet", "language": "de"}
]
```

**Response:**

```json
{
    "uploaded": 3,
    "failed": 0,
    "assets": [
        {"id": "uuid-1", "filename": "bild1.jpg", "processing_status": "pending"},
        {"id": "uuid-2", "filename": "bild2.jpg", "processing_status": "pending"},
        {"id": "uuid-3", "filename": "datenblatt.pdf", "processing_status": "pending"}
    ]
}
```

#### Asset-Metadaten aktualisieren

```
PUT /api/v1/admin/assets/:id
{
    "title": {"de": "Neuer Titel", "fr": "Nouveau titre"},
    "alt_text": {"de": "Beschreibender Alt-Text"},
    "document_type": "certificate",
    "language": "de"
}
```

#### Asset-Zuordnung

```
# Asset einem Produkt zuordnen
POST /api/v1/admin/assets/:asset_id/assign
{
    "entity_type": "product",
    "entity_id": "product-uuid",
    "role": "gallery",
    "sort_order": 0,
    "is_primary": false
}

# Zuordnung entfernen (Asset bleibt erhalten)
DELETE /api/v1/admin/assets/:asset_id/assign/:assignment_id

# Zuordnungen eines Produkts sortieren
PUT /api/v1/admin/products/:id/assets/reorder
{
    "assignments": [
        {"assignment_id": "uuid-1", "sort_order": 0, "is_primary": true},
        {"assignment_id": "uuid-2", "sort_order": 1, "is_primary": false},
        {"assignment_id": "uuid-3", "sort_order": 2, "is_primary": false}
    ]
}
```

#### Neue Dokumentversion hochladen

```
POST /api/v1/admin/assets/:id/new-version
Content-Type: multipart/form-data

file: <binary>
valid_from: 2026-03-01 (optional)
```

Erstellt ein neues Asset, setzt `superseded_by` beim alten, übernimmt alle Zuordnungen.

### Events

| Event | Trigger | Konsumenten |
|-------|---------|-------------|
| `asset.uploaded` | Neues Asset hochgeladen | Processing Pipeline |
| `asset.processing.complete` | Verarbeitung abgeschlossen | Frontend (Websocket/Polling), Meilisearch |
| `asset.processing.failed` | Verarbeitung fehlgeschlagen | Admin-Notification |
| `asset.assigned` | Asset einer Entität zugeordnet | Meilisearch (Produkt-Update), Cache-Invalidation |
| `asset.unassigned` | Zuordnung entfernt | Meilisearch, Cache-Invalidation |
| `asset.deleted` | Soft-Delete | Meilisearch, CDN-Invalidation |
| `asset.version.created` | Neue Dokumentversion | Benachrichtigung an Kunden (optional) |

---

## 8. PIM/ERP-Integration

### Akeneo PIM

Akeneo unterscheidet zwischen **Product Media Files** (Bilder) und **Asset Manager** (Dokumente, Videos, Links). Beide werden in Gondolia als Assets importiert.

**Mapping Akeneo → Gondolia:**

| Akeneo | Gondolia |
|--------|---------|
| Product Media File (image) | Asset (type=image) + Assignment (entity=product) |
| Asset Manager Asset (PDF) | Asset (type=document) + Assignment (entity=product) |
| Reference Entity Image (z.B. Hersteller-Logo) | Asset (type=image) + Assignment (entity=manufacturer) |
| Asset Code | `pim_identifier` |
| Asset Locale | `language` |
| Asset Categories | `document_type` (Mapping-Tabelle) |

**Sync-Ablauf:**

1. Akeneo-Event oder geplanter Job triggert Import
2. Asset-Binärdatei wird aus Akeneo heruntergeladen
3. Asset wird in Gondolia Object Storage hochgeladen
4. Asset-Metadaten werden erstellt/aktualisiert (`pim_identifier` als Matching-Key)
5. Zuordnungen (Product → Asset) werden aus Akeneo-Produktdaten erstellt
6. Processing-Pipeline generiert Varianten (Thumbnails, WebP etc.)

**Wichtig:** Akeneo ist das führende System für Assets. Manuelle Änderungen in Gondolia werden beim nächsten Sync **überschrieben** — ausser sie betreffen Felder die Akeneo nicht liefert (z.B. `document_type`-Klassifizierung).

### SAP Document Management System (DMS)

SAP DMS speichert Dokumente als «Document Info Records» (DIR) mit Verknüpfung zu Materialstammdaten.

**Mapping SAP DMS → Gondolia:**

| SAP DMS | Gondolia |
|---------|---------|
| Document Info Record (DIR) | Asset (type=document) |
| DIR Nummer + Version | `erp_identifier` |
| Dokumentart (z.B. ZCE = Zertifikat) | `document_type` |
| Originaldatei | Binärdaten → Object Storage |
| Verknüpfung DIR ↔ Material | Asset Assignment (entity=product, via SAP-Materialnummer → Gondolia-Produkt) |

**Sync-Richtung:** SAP → Gondolia (unidirektional). Dokumente werden aus SAP importiert und in Gondolia gespeichert. Gondolia liest, SAP schreibt.

**Trigger:** IDoc `DOCUMENT_CREATE_DOC` / `DOCUMENT_CHANGE_DOC` oder RFC-Abruf über `BAPI_DOCUMENT_GETDETAIL2`.

---

## 9. Service-Architektur

Das Asset-Management wird als **Domain innerhalb des Catalog-Service** implementiert — analog zum Hersteller (vgl. [Hersteller-Entität](manufacturer-entity.md)).

**Begründung:** Assets sind primär Anhängsel von Katalog-Entitäten (Produkte, Kategorien, Hersteller). Ein eigener Microservice wäre Overengineering — die Verarbeitungslogik (Bildkonvertierung etc.) läuft ohnehin asynchron über Worker.

### Verzeichnisstruktur

```
services/catalog/
  internal/
    asset/
      handler.go           # HTTP Handler (Public + Admin)
      service.go            # Business Logic
      repository.go         # DB-Zugriff
      models.go             # Structs, DTOs, Enums
      errors.go             # Domain-spezifische Fehler
      storage.go            # S3-Abstraction (Interface)
      storage_s3.go         # S3/MinIO-Implementierung
      processing.go         # Bildverarbeitung, Thumbnail-Generierung
      processing_worker.go  # Async Worker (konsumiert Queue)
      asset_test.go         # Unit Tests
    product/
      ...                   # Erweitert: Images → Assets
    manufacturer/
      ...                   # logo_url → Asset-Beziehung
    domain/
      product.go            # Images-Feld entfällt, Assets-Beziehung neu
```

### Storage Interface

```go
// internal/asset/storage.go

type StorageProvider interface {
    Upload(ctx context.Context, key string, data io.Reader, contentType string) error
    Download(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    GetURL(ctx context.Context, key string) (string, error)
    GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}
```

---

## 10. Branchenbeispiele

### Industriebedarf — Hydraulikzylinder (Bosch Rexroth CDT3)

| Asset | Typ | Format | Sprachen |
|-------|-----|--------|----------|
| Produktfoto Frontal | Image (Primary) | JPEG → WebP/AVIF | — |
| Produktfoto Schnittzeichnung | Image (Gallery) | JPEG → WebP/AVIF | — |
| Technisches Datenblatt | Document (datasheet) | PDF, 12 Seiten | DE, EN, FR |
| CE-Konformitätserklärung | Document (declaration_of_conformity) | PDF, 2 Seiten | DE, EN |
| CAD-Modell 3D | Document (cad_drawing) | STEP | — |
| CAD-Zeichnung 2D | Document (cad_drawing) | DWG | — |
| ISO 9001 Zertifikat | Document (certificate) | PDF | — (sprachunabhängig) |

### Chemie — Industriereiniger (TechClean 5000)

| Asset | Typ | Format | Sprachen |
|-------|-----|--------|----------|
| Produktfoto Gebinde | Image (Primary) | JPEG → WebP/AVIF | — |
| **Sicherheitsdatenblatt** ⚠️ | Document (safety_datasheet) | PDF, 16 Seiten | **DE, FR, IT** (gesetzlich!) |
| Technisches Datenblatt | Document (datasheet) | PDF | DE, EN |
| Analysezertifikat Charge | Document (certificate) | PDF | DE |
| Anwendungshinweis | Document (operating_manual) | PDF | DE, FR |

### Elektro — Frequenzumrichter (Siemens SINAMICS G120)

| Asset | Typ | Format | Sprachen |
|-------|-----|--------|----------|
| Produktfoto | Image (Primary) | JPEG → WebP/AVIF | — |
| Betriebsanleitung | Document (operating_manual) | PDF, 340 Seiten | DE, EN, FR |
| Einbauanweisung | Document (installation_guide) | PDF, 24 Seiten | DE, EN |
| EU-Konformitätserklärung | Document (declaration_of_conformity) | PDF | DE, EN |
| Schaltplan Standardanschluss | Document (cad_drawing) | PDF/DWG | — |
| RoHS-Zertifikat | Document (certificate) | PDF | EN |

### Arbeitsschutz — Sicherheitsschuh (Uvex 1 x-tended)

| Asset | Typ | Format | Sprachen |
|-------|-----|--------|----------|
| Produktfoto seitlich | Image (Primary) | JPEG → WebP/AVIF | — |
| Produktfoto Sohle | Image (Gallery) | JPEG → WebP/AVIF | — |
| Produktfoto getragen | Image (Gallery) | JPEG → WebP/AVIF | — |
| Prüfbericht EN ISO 20345 | Document (test_report) | PDF | DE |
| CE-Zertifikat | Document (certificate) | PDF | — |
| Pflegeanleitung | Document (care_instruction) | PDF | DE, FR, IT |
| Grössentabelle | Document (datasheet) | PDF | DE |

---

## 11. Priorisierung

| Phase | Aufgabe | Aufwand |
|-------|---------|---------|
| **Phase 1** | DB-Migration, Asset-Entität, Repository, Storage-Interface (S3/MinIO), Upload-Endpoint (Admin) | 4–5 Tage |
| **Phase 2** | Processing-Pipeline (Thumbnail-Generierung, WebP-Konvertierung), Async Worker | 3–4 Tage |
| **Phase 3** | Asset-Zuordnungen (M:N), Product/Category/Manufacturer-Integration, Public API | 3–4 Tage |
| **Phase 4** | Migration `ProductImage` → Assets, Backward-Compatibility | 2–3 Tage |
| **Phase 5** | Frontend: Bildergalerie, Dokument-Tab, PDF-Viewer (pdf.js) | 5–7 Tage |
| **Phase 6** | Bulk-Upload, ZIP-Download, Versionierung | 2–3 Tage |
| **Phase 7** | PIM-Sync (Akeneo), SAP DMS-Anbindung | 3–5 Tage |
| **Phase 8** | CDN-Integration, Virus-Scan (ClamAV), AVIF-Support | 2–3 Tage |

**Gesamtaufwand:** ~5–7 Wochen

---

## Offene Fragen

1. **Storage-Entscheid:** MinIO (Self-Hosted) oder AWS S3? Oder beides als Option (pro Tenant konfigurierbar)?
2. **CDN:** CloudFront, Cloudflare oder Bunny CDN? Gibt es bereits eine CDN-Infrastruktur?
3. **Virus-Scan:** Ist ClamAV akzeptabel oder wird ein kommerzieller Scanner benötigt? SLA für Scan-Dauer?
4. **Akeneo Asset Manager:** Wird der Akeneo Asset Manager genutzt oder nur Product Media Files? Welche Dokumenttypen sind in Akeneo klassifiziert?
5. **SAP DMS:** Wird SAP DMS aktiv genutzt? Welche Dokumentarten (DIR-Typen) sind relevant?
6. **Signed URLs:** Sollen bestimmte Dokumente nur authentifizierten Benutzern zugänglich sein (z.B. kundenspezifische Preislisten)? Oder sind alle Assets öffentlich?
7. **Hersteller-Logo Migration:** Das aktuelle `logo_url`-Feld im Manufacturer-Struct (vgl. [Hersteller-Entität](manufacturer-entity.md)) — soll es durch eine Asset-Beziehung ersetzt werden oder als Convenience-Feld bestehen bleiben?
8. **CAD-Viewer:** Sollen CAD-Dateien (STEP, DWG) im Browser anzeigbar sein oder reicht Download? Ein 3D-Viewer (z.B. three.js mit STEP-Parser) wäre möglich, aber aufwändig.
9. **Wasserzeichen:** Sollen Bilder im Frontend mit einem Wasserzeichen versehen werden können (Tenant-konfigurierbar)?
10. **Retention Policy:** Wie lange werden alte Dokumentversionen und gelöschte Assets aufbewahrt? (Empfehlung: 1 Jahr für gelöschte Assets, unbegrenzt für Dokumentversionen wegen gesetzlicher Nachweispflicht)
