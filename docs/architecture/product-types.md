# Produkttypen-Architektur

Stand: 2026-02-14

## Übersicht

Das aktuelle Produkt-Model in Gondolia ist ein **Flat-Product** — ein einfaches Produkt mit SKU, Name, Attributen und Preis. Für eine B2B-Plattform reicht das bei weitem nicht. Kunden erwarten:

- Produkte in verschiedenen Varianten (Grösse, Farbe, Ausführung)
- Massanfertigungen mit kundenspezifischen Parametern (Zuschnitte, Konfektionierung)
- Komplexe Konfigurationen mit abhängigen Optionen
- Fertige Pakete/Bundles mit Mengenrabatt
- Konfiguratoren für zusammengesetzte Produkte

Dieses Dokument definiert **5 Produkttypen** mit ihrem Datenmodell, API-Design, UX-Anforderungen und Preisbildung.

### Aktuelles Produkt-Model (Ist-Zustand)

```go
type Product struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    SKU         string
    Name        map[string]string    // i18n
    Description map[string]string    // i18n
    CategoryIDs []uuid.UUID
    Attributes  []ProductAttribute   // Key/Value Paare
    Status      ProductStatus
    Images      []ProductImage
}
```

**Problem:** Kein Konzept für Varianten, Parameter, Bundles oder Konfigurationen. Alles ist ein flaches Produkt mit einer SKU.

### Typ-Hierarchie

```
ProductType (Discriminator)
├── simple          → Heutiges Flat-Product (bleibt erhalten)
├── variant_master  → Master-Produkt mit Varianten-Achsen
├── parametric      → Produkt mit freien Parameter-Eingaben
├── complex         → Parametrisch mit abhängigen Regeln
├── bundle          → Paket aus mehreren Produkten
└── configurator    → Komponentenbasiert mit Abhängigkeiten
```

---

## 1. Variantenprodukte

### Beschreibung

Ein **Variantenprodukt** besteht aus einem Master-Produkt und einer Matrix von Varianten. Das Master-Produkt definiert die gemeinsamen Eigenschaften (Name, Beschreibung, Bilder), die Varianten definieren die konkreten Ausprägungen mit eigener SKU, eigenem Preis und eigenem Bestand.

**Abgrenzung:**
- Vs. Parametrisierbar: Varianten haben **feste, vordefinierte** Werte. Keine Freitexteingabe.
- Vs. Bundle: Eine Variante ist **ein** Produkt, kein Paket aus mehreren.
- Vs. Konfigurator: Varianten-Achsen sind **unabhängig** voneinander (keine Einschränkungen zwischen Achsen).

### Branchenübergreifende Beispiele

| Branche | Produkt | Achse 1 | Achse 2 | Achse 3 | Varianten |
|---------|---------|---------|---------|---------|-----------|
| Industriebedarf | Hydraulikzylinder Bosch Rexroth | Kolben-Ø (40, 50, 63, 80mm) | Hub (100, 200, 300, 500mm) | Betriebsdruck (160, 250 bar) | bis 32 |
| Elektro/Elektronik | Netzwerkkabel Cat6a | Länge (0.5, 1, 2, 3, 5, 10m) | Farbe (Grau, Blau, Gelb, Rot) | Schirmung (U/UTP, S/FTP) | bis 48 |
| Verpackung | Faltkarton FEFCO 0201 | Innenmaß (300×200×150, 400×300×200, 600×400×300mm) | Wellenstärke (B-Welle, BC-Welle, EB-Welle) | Farbe (Braun, Weiss) | bis 18 |
| Textil B2B | Poloshirt Hakro Performance | Grösse (XS, S, M, L, XL, XXL, 3XL) | Farbe (Schwarz, Weiss, Navy, Rot, Grün) | — | bis 35 |
| Chemie/Labor | Isopropanol technisch | Reinheit (99.5%, 99.8%, 99.9%) | Gebinde (1L, 2.5L, 5L, 10L, 25L) | — | bis 15 |
| Holz/Baustoffe | Spanplatte Egger Eurospan | Dimension (2800×2070, 2440×1220mm) | Dicke (8, 12, 16, 19, 25mm) | Dekor (Buche, Eiche, Weiss) | bis 30 |
| IT/Hardware | Dell PowerEdge R760 | CPU (Xeon Silver, Gold, Platinum) | RAM (64, 128, 256, 512GB) | Storage (2×960GB SSD, 4×1.92TB SSD) | bis 36 |

**Sparse Matrix:** Nicht alle Kombinationen müssen existieren. Beispiel: Hydraulikzylinder Ø 40mm bei 250 bar nur mit Hub 100 und 200mm verfügbar.

### Datenmodell

```sql
-- Varianten-Achsen (pro Master-Produkt)
CREATE TABLE variant_axes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,              -- z.B. "diameter", "stroke", "pressure"
    name JSONB NOT NULL,                    -- i18n: {"de": "Kolben-Ø", "fr": "Ø piston"}
    position INTEGER NOT NULL DEFAULT 0,    -- Reihenfolge der Achsen
    UNIQUE(product_id, code)
);

-- Achsen-Werte
CREATE TABLE variant_axis_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    axis_id UUID NOT NULL REFERENCES variant_axes(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,             -- z.B. "63mm", "250bar"
    label JSONB NOT NULL,                   -- i18n: {"de": "63 mm", "fr": "63 mm"}
    position INTEGER NOT NULL DEFAULT 0,    -- Sortierung
    UNIQUE(axis_id, code)
);

-- Varianten (konkrete Ausprägungen)
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL,
    axis_values JSONB NOT NULL,             -- {"diameter": "63mm", "stroke": "300mm", "pressure": "250bar"}
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, sku)
);

-- Index für schnelle Variantensuche
CREATE INDEX idx_variant_axis_values ON product_variants USING GIN (axis_values);
```

```go
// Domain-Model Erweiterung
type VariantAxis struct {
    ID       uuid.UUID         `json:"id"`
    Code     string            `json:"code"`
    Name     map[string]string `json:"name"`
    Position int               `json:"position"`
    Values   []AxisValue       `json:"values"`
}

type AxisValue struct {
    ID       uuid.UUID         `json:"id"`
    Code     string            `json:"code"`
    Label    map[string]string `json:"label"`
    Position int               `json:"position"`
}

type ProductVariant struct {
    ID         uuid.UUID         `json:"id"`
    ProductID  uuid.UUID         `json:"product_id"`
    SKU        string            `json:"sku"`
    AxisValues map[string]string `json:"axis_values"` // axis_code → value_code
    Status     ProductStatus     `json:"status"`
    // Preis und Bestand kommen aus Pricing/Inventory Service
}
```

### API-Design

```
# Varianten-Master abrufen (inkl. Achsen und verfügbare Werte)
GET /api/v1/products/:id
→ Response enthält: product + axes[] + variants[]

# Spezifische Variante selektieren
GET /api/v1/products/:id/variants?diameter=63mm&stroke=300mm&pressure=250bar
→ Response: variant mit SKU, Preis, Verfügbarkeit

# Verfügbare Werte für Achse (basierend auf bereits gewählten)
GET /api/v1/products/:id/variants/available?diameter=63mm
→ Response: welche Hübe und Betriebsdrücke für diesen Durchmesser existieren

# Bestellen
POST /api/v1/cart/items
{
    "variant_id": "uuid-der-variante",
    "quantity": 5
}
```

**Response-Beispiel (Industriebedarf — Hydraulikzylinder):**

```json
{
  "id": "master-uuid",
  "type": "variant_master",
  "sku": "REXROTH-CDT3",
  "name": {"de": "Bosch Rexroth Hydraulikzylinder CDT3"},
  "axes": [
    {
      "code": "diameter",
      "name": {"de": "Kolben-Ø"},
      "values": [
        {"code": "40mm", "label": {"de": "40 mm"}},
        {"code": "50mm", "label": {"de": "50 mm"}},
        {"code": "63mm", "label": {"de": "63 mm"}},
        {"code": "80mm", "label": {"de": "80 mm"}}
      ]
    },
    {
      "code": "stroke",
      "name": {"de": "Hub"},
      "values": [
        {"code": "100mm", "label": {"de": "100 mm"}},
        {"code": "200mm", "label": {"de": "200 mm"}},
        {"code": "300mm", "label": {"de": "300 mm"}},
        {"code": "500mm", "label": {"de": "500 mm"}}
      ]
    },
    {
      "code": "pressure",
      "name": {"de": "Betriebsdruck"},
      "values": [
        {"code": "160bar", "label": {"de": "160 bar"}},
        {"code": "250bar", "label": {"de": "250 bar"}}
      ]
    }
  ],
  "variants": [
    {
      "id": "variant-uuid-1",
      "sku": "REXROTH-CDT3-63-300-250",
      "axis_values": {"diameter": "63mm", "stroke": "300mm", "pressure": "250bar"},
      "price": {"net": 1245.00, "currency": "CHF"},
      "availability": {"in_stock": true, "quantity": 8}
    }
  ]
}
```

### Frontend-UX-Anforderungen

1. **Achsen-Selektor:** Dropdowns oder Button-Gruppen für jede Achse. Bei ≤6 Werten: Buttons. Bei >6: Dropdown.
2. **Verfügbarkeitsanzeige:** Nicht existierende Kombinationen als ausgegraut/deaktiviert markieren. Nicht einfach verstecken — der Kunde muss sehen *warum* eine Kombination nicht geht.
3. **Live-Update:** Bei Auswahl einer Achse sofort Preis, Verfügbarkeit und Bild aktualisieren. Kein Page-Reload.
4. **Bild-Wechsel:** Varianten können eigene Bilder haben (z.B. verschiedene Farben bei Arbeitskleidung). Bei Auswahl Bild wechseln.
5. **Matrix-Ansicht (optional):** Für erfahrene B2B-Kunden: Tabellarische Ansicht aller Varianten mit Preis und «In den Warenkorb»-Button pro Zeile. Ermöglicht Schnellbestellung mehrerer Varianten.
6. **URL-Update:** Gewählte Variante in URL abbilden (`?variant=REXROTH-CDT3-63-300-250`) für Bookmarks und Teilen.

### Preisbildung

- Jede Variante hat einen **eigenen Preis** (aus Pricing Service / SAP-Konditionierung)
- Alternativ: Master-Preis + Aufschlag pro Achsenwert (z.B. Reinheit 99.9% +20%)
- B2B: Kundenspezifische Preise pro Variante möglich (SAP-Konditionen auf SKU-Ebene)
- Staffelpreise pro Variante: Ab 10 Stück -5%, ab 50 Stück -10%

---

## 2. Parametrisierbare Produkte

### Beschreibung

Ein **parametrisierbares Produkt** hat keine feste Variantenmatrix, sondern erlaubt dem Kunden die **freie Eingabe von Werten** innerhalb definierter Grenzen. Typisch für Zuschnitte, Konfektionierung, Meterware und konfigurierbare Einzelprodukte.

**Abgrenzung:**
- Vs. Varianten: Keine vordefinierte Matrix. Der Kunde gibt exakte Werte ein.
- Vs. Komplex: Parameter sind **unabhängig** voneinander. Keine Regeln zwischen Parametern.
- Vs. Konfigurator: Ein Produkt, nicht mehrere Komponenten.

### Branchenübergreifende Beispiele

| Branche | Produkt | Parameter 1 | Parameter 2 | Parameter 3 | Preisformel |
|---------|---------|-------------|-------------|-------------|-------------|
| Elektro | Kabel NYM-J 3×1.5mm² | Länge (1–500m, 1m Schritte) | — | — | Länge × Preis/m |
| Verpackung | Stretchfolie | Breite (100–500mm, 10mm) | Stärke (Dropdown: 12, 17, 20, 23µm) | — | Breite × Stärke-Faktor × Rollenpreis |
| Textil B2B | Stoff Meterware (z.B. Baumwoll-Twill) | Länge (0.5–50m, 0.1m) | Breite (Dropdown: 150, 160, 300cm) | — | Fläche × Preis/m² |
| Holz/Baustoffe | Plattenzuschnitt MDF | Länge (100–2800mm, 1mm) | Breite (100–2070mm, 1mm) | Dicke (Dropdown: 6, 12, 19, 25mm) | Fläche × Preis/m² + Zuschlag/Schnitt |
| Gastro | Edelstahl-Arbeitstisch nach Mass | Länge (600–3000mm, 50mm) | Tiefe (600–900mm, 50mm) | Höhe (Dropdown: 850, 900mm) | Fläche × Materialpreis + Fertigung |
| Chemie | Verdünnung/Mischung (z.B. Ethanol) | Volumen (5–1000L, 1L) | Konzentration (Dropdown: 70%, 80%, 96%) | — | Volumen × Konzentrations-Preis/L |
| IT/Hardware | Netzwerkkabel Meterware Cat7 | Länge (1–305m, 1m) | — | — | Länge × Preis/m |

### Datenmodell

```sql
-- Parameter-Definitionen (pro Produkt)
CREATE TABLE product_parameters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,              -- z.B. "length", "width", "concentration"
    name JSONB NOT NULL,                    -- i18n
    description JSONB,                      -- i18n, Hilfetext

    -- Typ & Einheit
    data_type VARCHAR(20) NOT NULL,         -- number, integer, select
    unit VARCHAR(20),                       -- mm, cm, m, m², kg, L, µm

    -- Grenzen (für number/integer)
    min_value DECIMAL(10,2),
    max_value DECIMAL(10,2),
    step DECIMAL(10,2),                     -- Schrittweite
    default_value DECIMAL(10,2),

    -- Optionen (für select)
    options JSONB,                          -- [{"code": "23um", "label": {"de": "23 µm"}, "value": 23}]

    -- Darstellung
    position INTEGER NOT NULL DEFAULT 0,
    is_required BOOLEAN NOT NULL DEFAULT true,

    UNIQUE(product_id, code)
);

-- Preisformeln
CREATE TABLE product_pricing_formulas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    formula TEXT NOT NULL,                  -- z.B. "(length * width / 1000000) * price_per_sqm + cut_surcharge"
    variables JSONB NOT NULL,               -- {"price_per_sqm": 45.00, "cut_surcharge": 8.50}
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ,
    UNIQUE(product_id, valid_from)
);
```

```go
type ProductParameter struct {
    ID           uuid.UUID         `json:"id"`
    Code         string            `json:"code"`
    Name         map[string]string `json:"name"`
    Description  map[string]string `json:"description,omitempty"`
    DataType     string            `json:"data_type"`     // "number", "integer", "select"
    Unit         string            `json:"unit,omitempty"` // "mm", "m²", "L"
    MinValue     *float64          `json:"min_value,omitempty"`
    MaxValue     *float64          `json:"max_value,omitempty"`
    Step         *float64          `json:"step,omitempty"`
    DefaultValue *float64          `json:"default_value,omitempty"`
    Options      []ParameterOption `json:"options,omitempty"` // für select
    Position     int               `json:"position"`
    IsRequired   bool              `json:"is_required"`
}

type ParameterOption struct {
    Code  string            `json:"code"`
    Label map[string]string `json:"label"`
    Value float64           `json:"value"`
}

type PricingFormula struct {
    Formula   string             `json:"formula"`
    Variables map[string]float64 `json:"variables"`
}
```

### API-Design

```
# Produkt mit Parametern abrufen
GET /api/v1/products/:id
→ Response enthält: product + parameters[] + pricing_formula

# Preis berechnen (Live-Preview)
POST /api/v1/products/:id/calculate-price
{
    "parameters": {
        "length": 50,
        "width": 300,
        "thickness": "20um"
    }
}
→ Response:
{
    "price": {
        "net": 187.50,
        "currency": "CHF",
        "breakdown": {
            "area_factor": 15.0,
            "thickness_factor": 1.25,
            "base_price": 150.00,
            "surcharge": 37.50,
            "total": 187.50
        }
    },
    "valid_for_seconds": 300
}

# Bestellen
POST /api/v1/cart/items
{
    "product_id": "uuid",
    "quantity": 1,
    "parameters": {
        "length": 50,
        "width": 300,
        "thickness": "20um"
    }
}
```

### Frontend-UX-Anforderungen

1. **Eingabefelder:** Numerische Inputs mit Min/Max-Validierung, Schrittweite, Einheit-Label.
2. **Live-Preisberechnung:** Bei jeder Parameteränderung Preis neu berechnen (debounced, ~300ms). Preis prominent anzeigen.
3. **Preisaufschlüsselung:** Transparente Darstellung: Menge × Einheitspreis + Zuschläge. B2B-Kunden wollen verstehen, wie der Preis zustande kommt.
4. **Validierung:** Sofortige Validierung im Frontend (Min/Max, Schrittweite). Fehler direkt am Feld anzeigen.
5. **Visualisierung (optional):** Massskizze die sich dynamisch mit den Eingaben ändert. Bei Kabeln z.B. Längenanzeige, bei Platten ein Rechteck mit Masslinie.
6. **Einheitenumrechnung:** Eingabe in mm, Anzeige auch in cm/m wo sinnvoll. Automatische Umrechnung.

### Preisbildung

- **Formelbasiert:** z.B. `length_m × price_per_meter` oder `(length_mm × width_mm / 1_000_000) × price_per_sqm + surcharge`
- **Variablen** kommen aus dem Pricing Service (mandantenspezifisch, kundenspezifisch)
- **Zuschläge:** Pro Schnitt, pro Bearbeitungsschritt, Express-Aufschlag
- **Mindestbestellwert:** Pro parametrischem Produkt konfigurierbar
- **Backend-Validierung:** Preis wird beim Hinzufügen zum Warenkorb und bei Bestellabschluss **nochmals** berechnet. Frontend-Preis ist nur Vorschau.

---

## 3. Komplexe Produkte (Abhängige Parameter)

### Beschreibung

Wie parametrisierbare Produkte, aber mit **Abhängigkeiten zwischen Parametern**. Die Wahl eines Parameters schränkt die erlaubten Werte anderer Parameter ein. Erfordert ein Regel-/Constraint-System.

**Abgrenzung:**
- Vs. Parametrisierbar: Parameter sind **nicht unabhängig**. Es gibt Regeln.
- Vs. Konfigurator: Immer noch **ein** Produkt mit Parametern, nicht mehrere Komponenten.
- Vs. Varianten: Keine feste Matrix, sondern Regeln die dynamisch auswerten.

### Branchenübergreifende Beispiele

**Beispiel 1: Schaltschrank-Konfiguration (Elektro)**

| Wenn | Dann |
|------|------|
| Schutzart = IP65 oder IP67 | Nur Edelstahl- oder Polyester-Gehäuse |
| Gehäusematerial = Stahlblech | Schutzart max. IP55 |
| Breite > 800mm | Nur Doppeltür verfügbar |
| Montageplatte = ja | Tiefe min. 210mm |
| Aussenaufstellung = ja | Schutzart min. IP55, Dachaufsatz erforderlich |

**Beispiel 2: Hydraulikschlauch konfektioniert (Industriebedarf)**

| Wenn | Dann |
|------|------|
| Betriebsdruck > 350 bar | Nur Schlauchtyp 2SN oder 4SP |
| Schlauchtyp = 4SP | Nennweite max. DN25 |
| Nennweite ≥ DN32 | Nur Pressfittinge (keine Schraubfittinge) |
| Medium = «Heissdampf» | Nur EPDM-Innenseele, max. Temp 210°C |
| Länge > 10m | Mindest-Nennweite DN12 |

**Beispiel 3: Server-Konfiguration (IT/Hardware)**

| Wenn | Dann |
|------|------|
| CPU = Xeon Platinum | RAM-Typ: nur DDR5 RDIMM |
| RAM > 512GB | Min. 2 CPUs erforderlich |
| GPU = NVIDIA A100 | Netzteil min. 1400W, Gehäuse = 2U oder 4U |
| Storage-Typ = NVMe | Max. 24 Laufwerke (statt 36 bei SAS) |
| Betriebssystem = VMware ESXi | Kompatibilitäts-Check auf RAM/NIC-Kombination |

**Beispiel 4: Industriereiniger-Mischung (Chemie)**

| Wenn | Dann |
|------|------|
| Konzentration > 30% | Gefahrgut-Versand, nur 25L/200L Gebinde |
| Gebinde = IBC 1000L | Nur Lieferung per Spedition |
| Anwendung = «Lebensmittelbereich» | Nur NSF-zugelassene Wirkstoffe |
| pH-Wert < 2 | Nur säurebeständige Gebinde (HDPE, Edelstahl) |

### Datenmodell

```sql
-- Regeln/Constraints zwischen Parametern
CREATE TABLE product_parameter_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(100),                      -- Beschreibung der Regel
    priority INTEGER NOT NULL DEFAULT 0,    -- Reihenfolge der Auswertung

    -- Bedingung (JSON-basiert)
    condition JSONB NOT NULL,
    -- Beispiel: {"parameter": "protection_class", "operator": "in", "value": ["IP65", "IP67"]}
    -- Beispiel: {"and": [{"parameter": "pressure", "operator": "gt", "value": 350}, {"parameter": "length", "operator": "gt", "value": 5}]}

    -- Auswirkung
    effect JSONB NOT NULL,
    -- Beispiel: {"parameter": "housing_material", "restrict_to": ["stainless_steel", "polyester"]}
    -- Beispiel: {"parameter": "nominal_width", "set_max": 25}
    -- Beispiel: {"parameter": "roof_attachment", "set_required": true}

    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_parameter_rules_product ON product_parameter_rules(product_id, is_active, priority);
```

```go
// Regel-Condition (verschachtelt möglich)
type RuleCondition struct {
    Parameter string         `json:"parameter,omitempty"`
    Operator  string         `json:"operator,omitempty"` // eq, neq, gt, gte, lt, lte, in
    Value     any            `json:"value,omitempty"`
    And       []RuleCondition `json:"and,omitempty"`
    Or        []RuleCondition `json:"or,omitempty"`
}

// Regel-Effekt
type RuleEffect struct {
    Parameter   string   `json:"parameter"`
    RestrictTo  []string `json:"restrict_to,omitempty"`  // Erlaubte Werte einschränken
    SetMin      *float64 `json:"set_min,omitempty"`
    SetMax      *float64 `json:"set_max,omitempty"`
    SetRequired *bool    `json:"set_required,omitempty"`
    Hide        *bool    `json:"hide,omitempty"`          // Parameter verstecken
}

type ParameterRule struct {
    ID        uuid.UUID     `json:"id"`
    Name      string        `json:"name"`
    Priority  int           `json:"priority"`
    Condition RuleCondition `json:"condition"`
    Effect    RuleEffect    `json:"effect"`
}
```

### API-Design

```
# Produkt abrufen (inkl. Parameter + Regeln)
GET /api/v1/products/:id
→ Response enthält: product + parameters[] + rules[]

# Regeln auswerten (bei Parameteränderung)
POST /api/v1/products/:id/evaluate-rules
{
    "parameters": {
        "protection_class": "IP65",
        "width": 1000,
        "depth": 300
    }
}
→ Response:
{
    "adjusted_parameters": {
        "housing_material": {
            "available_options": ["stainless_steel", "polyester"],
            "removed_options": ["steel_sheet"]
        },
        "door_type": {
            "available_options": ["double_door"],
            "removed_options": ["single_door"],
            "note": "Doppeltür erforderlich ab 800mm Breite"
        }
    },
    "violations": [],
    "price": { ... }
}

# Validierung + Preis (vor Warenkorb)
POST /api/v1/products/:id/validate-and-price
{
    "parameters": {
        "protection_class": "IP65",
        "housing_material": "stainless_steel",
        "width": 1000,
        "depth": 300,
        "door_type": "double_door"
    }
}
→ Response:
{
    "valid": true,
    "price": {"net": 2890.00, "currency": "CHF"},
    "violations": []
}
```

### Frontend-UX-Anforderungen

1. **Dynamische Formulare:** Bei Änderung eines Parameters sofort Regeln auswerten und abhängige Felder aktualisieren (Optionen einschränken, Min/Max anpassen, Felder ein-/ausblenden).
2. **Erklärungen:** Wenn eine Option deaktiviert wird, dem Kunden erklären *warum*. Tooltip: «Bei Schutzart IP65 ist Stahlblech-Gehäuse nicht verfügbar.»
3. **Reihenfolge:** Parameter in sinnvoller Reihenfolge anzeigen. Entscheidende Parameter zuerst (Schutzart → Gehäuse → Masse → Optionen).
4. **Validierung:** Doppelt — Frontend für sofortige UX, Backend für Sicherheit. Frontend-Regeln können vom Backend abweichen (Subset reicht).
5. **Regel-Engine im Frontend:** Regeln als JSON vom Backend laden und clientseitig auswerten. Kein API-Call bei jeder Parameteränderung (Performance). Nur für Preis API-Call.

### Preisbildung

- Wie parametrisierbare Produkte, aber Formel kann **regelabhängig** sein
- Beispiel: Schutzart IP67 = +25% auf Grundpreis, Edelstahl-Gehäuse = +40%
- Zuschläge abhängig von gewählten Optionen
- Backend-Validierung: Ungültige Kombinationen → Preis kann nicht berechnet werden → Fehler

### Regel-Engine

**Anforderungen:**

- Regeln als JSON gespeichert (keine Hardcodierung)
- Verschachtelung: AND/OR-Verknüpfungen
- Operatoren: eq, neq, gt, gte, lt, lte, in, not_in
- Effekte: restrict_to, set_min, set_max, hide, set_required
- Priorität: Bei Konflikten gewinnt die Regel mit höherer Priorität
- Performance: Alle Regeln eines Produkts werden beim Laden mitgeliefert und clientseitig ausgewertet
- **Identische Logik** in Frontend (TypeScript) und Backend (Go) — oder Backend als Single Source of Truth mit API-Call

**Empfehlung:** Regeln als JSON an Frontend liefern, dort clientseitig auswerten. Bei Bestellung Backend-Validierung als Sicherheitsnetz. Vorteil: Kein API-Call bei jeder Parameteränderung. Nachteil: Regel-Engine muss in beiden Sprachen implementiert werden.

---

## 4. Bundle-Produkte (Packages)

### Beschreibung

Ein **Bundle** ist ein Paket aus mehreren Einzelprodukten, das als Einheit verkauft wird. Die enthaltenen Produkte sind eigenständige Katalogprodukte. Bundles bieten typischerweise einen Preisvorteil gegenüber Einzelkauf.

**Abgrenzung:**
- Vs. Konfigurator: Bundle hat **keine Abhängigkeiten** zwischen Komponenten. Es ist eine feste (oder teilweise variable) Zusammenstellung.
- Vs. Varianten: Bundle besteht aus **verschiedenen** Produkten, nicht verschiedenen Ausprägungen desselben Produkts.

### Branchenübergreifende Beispiele

| Branche | Bundle | Komponenten | Preis-Modell |
|---------|--------|-------------|--------------|
| Industriebedarf | Pneumatik-Starterpaket | 5× Zylinder DSBC Ø32, 10× Steckverbinder QS-6, 5m Schlauch PUN-6, 1× Wartungseinheit MSB4 | Summe - 12% |
| IT/Hardware | Arbeitsplatz-Komplett | 1× Monitor 27", 1× Dockingstation, 1× Tastatur, 1× Maus, 1× Headset | Festpreis (günstiger als Einzelteile) |
| Gastro/Hotel | Hotelzimmer-Grundausstattung | 2× Handtuch-Set, 1× Bettwäsche-Set, 1× Seifenspender, 1× Wasserkocher | Summe - 15% |
| Chemie/Labor | Labor-Starterset Analytik | 1× Becherglas-Set (6-teilig), 1× Pipetten-Set, 1× Schutzbrille, 1× Handschuh-Box (100 Stk) | Festpreis 89.00 CHF |
| Verpackung | Versandstation-Komplett | X Kartons (Menge variabel) + Füllmaterial + Klebeband + Etiketten | Summe berechnet |
| Elektro | LED-Hallenbeleuchung 500m² | 20× LED-Hallenstrahler 150W, 20× Aufhängung, 1× Steuermodul DALI, 500m Kabel | Projektpreis |
| Textil B2B | Musterkollektion Herbst/Winter | 12× verschiedene Stoffmuster (je 30×30cm) | Festpreis 0.00 CHF (Gratis-Muster) |

### Datenmodell

```sql
-- Bundle-Positionen
CREATE TABLE bundle_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bundle_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    item_product_id UUID NOT NULL REFERENCES products(id),
    quantity DECIMAL(10,2) NOT NULL DEFAULT 1,
    is_quantity_fixed BOOLEAN NOT NULL DEFAULT true,   -- Feste oder variable Menge
    min_quantity DECIMAL(10,2),                         -- Bei variabler Menge
    max_quantity DECIMAL(10,2),
    position INTEGER NOT NULL DEFAULT 0,
    is_optional BOOLEAN NOT NULL DEFAULT false,         -- Optionale Komponente

    UNIQUE(bundle_product_id, item_product_id)
);

-- Bundle-Preismodell
CREATE TABLE bundle_pricing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bundle_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    pricing_type VARCHAR(20) NOT NULL,       -- 'fixed', 'sum_discount_percent', 'sum_discount_absolute'
    fixed_price DECIMAL(10,2),               -- Bei pricing_type = 'fixed'
    discount_percent DECIMAL(5,2),           -- Bei pricing_type = 'sum_discount_percent'
    discount_absolute DECIMAL(10,2),         -- Bei pricing_type = 'sum_discount_absolute'
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ
);
```

```go
type BundleItem struct {
    ID              uuid.UUID `json:"id"`
    Product         *Product  `json:"product"`          // Das enthaltene Produkt
    Quantity        float64   `json:"quantity"`
    IsQuantityFixed bool      `json:"is_quantity_fixed"`
    MinQuantity     *float64  `json:"min_quantity,omitempty"`
    MaxQuantity     *float64  `json:"max_quantity,omitempty"`
    Position        int       `json:"position"`
    IsOptional      bool      `json:"is_optional"`
}

type BundlePricing struct {
    Type             string   `json:"type"` // "fixed", "sum_discount_percent", "sum_discount_absolute"
    FixedPrice       *float64 `json:"fixed_price,omitempty"`
    DiscountPercent  *float64 `json:"discount_percent,omitempty"`
    DiscountAbsolute *float64 `json:"discount_absolute,omitempty"`
}
```

### API-Design

```
# Bundle abrufen
GET /api/v1/products/:id
→ Response: product (type=bundle) + bundle_items[] + bundle_pricing

# Bundle-Preis berechnen (bei variablen Mengen)
POST /api/v1/products/:id/calculate-bundle-price
{
    "items": [
        {"product_id": "uuid-kartons", "quantity": 500},
        {"product_id": "uuid-fuellmaterial", "quantity": 20},
        {"product_id": "uuid-klebeband", "quantity": 12}
    ],
    "optional_items": ["uuid-etiketten"]
}
→ Response:
{
    "items": [
        {"product_id": "uuid-kartons", "quantity": 500, "unit_price": 1.20, "total": 600.00},
        {"product_id": "uuid-fuellmaterial", "quantity": 20, "unit_price": 8.50, "total": 170.00},
        {"product_id": "uuid-klebeband", "quantity": 12, "unit_price": 4.90, "total": 58.80}
    ],
    "subtotal": 828.80,
    "discount": -82.88,
    "total": {"net": 745.92, "currency": "CHF"},
    "savings_vs_individual": 82.88
}

# Bestellen
POST /api/v1/cart/items
{
    "product_id": "uuid-bundle",
    "bundle_items": [
        {"product_id": "uuid-kartons", "quantity": 500},
        {"product_id": "uuid-fuellmaterial", "quantity": 20},
        {"product_id": "uuid-klebeband", "quantity": 12}
    ]
}
```

### Frontend-UX-Anforderungen

1. **Bundle-Übersicht:** Alle enthaltenen Produkte mit Bild, Name, Menge, Einzelpreis auflisten. Ersparnis prominent anzeigen.
2. **Variable Mengen:** Bei nicht-festen Mengen: Mengen-Input pro Komponente. Live-Preisupdate.
3. **Optionale Komponenten:** Checkbox zum An-/Abwählen. Preis aktualisiert sich.
4. **Einzelprodukt-Links:** Jede Komponente verlinkt auf das Einzelprodukt (für Details).
5. **Vergleich:** «Als Einzelteile kaufen: CHF 828.80 — Als Bundle: CHF 745.92 — Sie sparen CHF 82.88»
6. **Verfügbarkeit:** Bundle nur bestellbar wenn **alle** Pflichtkomponenten verfügbar. Einzelverfügbarkeit pro Komponente anzeigen.

### Preisbildung

| Modell | Beschreibung | Beispiel |
|--------|--------------|---------|
| **Festpreis** | Bundle hat einen fixen Preis, unabhängig von Einzelpreisen | Musterkollektion = 0.00 CHF |
| **Summe - Prozent** | Summe der Einzelpreise minus X% | Pneumatik-Starterpaket = Summe - 12% |
| **Summe - Absolut** | Summe der Einzelpreise minus fester Betrag | Bundle = Summe - 50.00 CHF |
| **Kundenspezifisch** | SAP-Konditionen auf Bundle-SKU | Enterprise-Kunden: eigener Bundle-Preis |

---

## 5. Konfigurator-Produkte

### Beschreibung

Ein **Konfigurator-Produkt** besteht aus wählbaren **Komponenten-Gruppen** (Steps), wobei die Auswahl in einem Step die verfügbaren Optionen in nachfolgenden Steps einschränkt. Das Ergebnis ist ein kundenindividuell zusammengestelltes Produkt.

**Abgrenzung:**
- Vs. Bundle: Bundle hat **keine Abhängigkeiten** zwischen Komponenten. Konfigurator hat kaskadierene Einschränkungen.
- Vs. Komplex: Komplexe Produkte sind **ein** Produkt mit abhängigen Parametern. Konfigurator besteht aus **mehreren wählbaren Komponenten-Produkten**.
- Vs. Varianten: Varianten haben feste Achsen. Konfigurator hat Steps mit Abhängigkeiten.

### Branchenübergreifende Beispiele

**Beispiel 1: Grossküchen-Kochblock (Gastro)**

```
Step 1: Grundmodul
  → Kochfeld-Zeile 4er, Kochfeld-Zeile 6er, Kochfeld-Zeile 8er

Step 2: Kochstellen-Typ (abhängig von Grundmodul)
  → 4er: Gas, Induktion
  → 6er: Gas, Induktion, Kombi (Gas+Induktion)
  → 8er: nur Kombi

Step 3: Unterbau (abhängig von Kochstellen-Typ)
  → Gas: offener Unterbau, Backofen GN2/1
  → Induktion: Kühlunterbau, Schubladenblock, Backofen GN2/1
  → Kombi: Backofen GN2/1, Neutralunterbau

Step 4: Material/Oberfläche
  → Edelstahl AISI 304, Edelstahl AISI 316 (Seeluft)

Step 5: Optionen
  → Spritzschutzrückwand, Ablage, Warmhaltebrücke

Step 6: Anschlüsse (abhängig von Kochstellen-Typ)
  → Gas: Gasanschluss R½", R¾"
  → Induktion: Starkstrom 32A, 63A
```

**Beispiel 2: 19"-Rack-Konfiguration (IT/Hardware)**

```
Step 1: Gehäuse → 24HE, 42HE, 47HE
Step 2: Tiefe (abhängig von Gehäuse) → 600mm, 800mm, 1000mm, 1200mm
Step 3: Klimatisierung → Passiv, Dachlüfter, Seitenkühlung (abhängig von HE + Tiefe)
Step 4: Stromversorgung → 1× PDU, 2× PDU redundant (abhängig von HE)
Step 5: Kabelmanagement → Standard, Premium (abhängig von Tiefe)
Step 6: Sicherheit → Schloss, Codelock, Biometrisch
Step 7: Zubehör → Fachböden, Blindplatten, Kabeldurchführungen (optional)
```

**Beispiel 3: Berufsbekleidungs-Set (Textil B2B)**

```
Step 1: Branche → Handwerk, Industrie, Medizin, Gastronomie
Step 2: Oberteil (abhängig von Branche)
  → Handwerk: Arbeitsjacke, Softshelljacke, Fleece
  → Industrie: Bundjacke, Warnschutzjacke
  → Medizin: Kasack, Laborkittel
  → Gastronomie: Kochjacke, Serviceweste

Step 3: Hose (abhängig von Branche)
  → Handwerk: Bundhose, Latzhose
  → Industrie: Bundhose, Multinormhose
  → Medizin: Schlupfhose
  → Gastronomie: Kochhose, Servicehose

Step 4: Grösse (parametrisierbar) → Konfektionsgrösse oder Körpermasse
Step 5: Veredelung → Logo-Stick, Namens-Stick, Bedruckung (optional)
Step 6: Schuhe (optional, abhängig von Branche)
  → Handwerk/Industrie: Sicherheitsschuhe S1/S2/S3
  → Medizin: Berufsschuhe OB
```

### Datenmodell

```sql
-- Konfigurator-Steps
CREATE TABLE configurator_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name JSONB NOT NULL,                    -- i18n
    description JSONB,                      -- i18n
    position INTEGER NOT NULL DEFAULT 0,
    is_required BOOLEAN NOT NULL DEFAULT true,
    step_type VARCHAR(20) NOT NULL,         -- 'select_one', 'select_many', 'parametric'

    UNIQUE(product_id, code)
);

-- Optionen pro Step
CREATE TABLE configurator_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    step_id UUID NOT NULL REFERENCES configurator_steps(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    name JSONB NOT NULL,                    -- i18n
    description JSONB,
    product_id UUID REFERENCES products(id), -- Optionale Verknüpfung mit Katalogprodukt
    price_modifier JSONB,                   -- {"type": "absolute", "value": 45.00} oder {"type": "percent", "value": 10}
    image_url VARCHAR(500),
    position INTEGER NOT NULL DEFAULT 0,

    UNIQUE(step_id, code)
);

-- Abhängigkeiten zwischen Steps
CREATE TABLE configurator_dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    -- Wenn in Step X Option Y gewählt wurde...
    source_step_id UUID NOT NULL REFERENCES configurator_steps(id),
    source_option_codes JSONB NOT NULL,     -- ["gas", "kombi"]

    -- ...dann sind in Step Z nur diese Optionen verfügbar
    target_step_id UUID NOT NULL REFERENCES configurator_steps(id),
    allowed_option_codes JSONB NOT NULL,    -- ["open_base", "oven_gn21"]

    priority INTEGER NOT NULL DEFAULT 0
);
```

```go
type ConfiguratorStep struct {
    ID          uuid.UUID            `json:"id"`
    Code        string               `json:"code"`
    Name        map[string]string    `json:"name"`
    Description map[string]string    `json:"description,omitempty"`
    Position    int                  `json:"position"`
    IsRequired  bool                 `json:"is_required"`
    StepType    string               `json:"step_type"` // "select_one", "select_many", "parametric"
    Options     []ConfiguratorOption `json:"options"`
}

type ConfiguratorOption struct {
    ID            uuid.UUID         `json:"id"`
    Code          string            `json:"code"`
    Name          map[string]string `json:"name"`
    ProductID     *uuid.UUID        `json:"product_id,omitempty"`
    PriceModifier *PriceModifier    `json:"price_modifier,omitempty"`
    ImageURL      string            `json:"image_url,omitempty"`
    Position      int               `json:"position"`
    IsAvailable   bool              `json:"is_available"` // Dynamisch berechnet
}

type PriceModifier struct {
    Type  string  `json:"type"`  // "absolute", "percent"
    Value float64 `json:"value"`
}

type ConfiguratorDependency struct {
    SourceStepCode      string   `json:"source_step_code"`
    SourceOptionCodes   []string `json:"source_option_codes"`
    TargetStepCode      string   `json:"target_step_code"`
    AllowedOptionCodes  []string `json:"allowed_option_codes"`
}
```

### API-Design

```
# Konfigurator laden
GET /api/v1/products/:id
→ Response: product (type=configurator) + steps[] + dependencies[]

# Verfügbare Optionen für nächsten Step (basierend auf bisheriger Auswahl)
POST /api/v1/products/:id/configure
{
    "selections": {
        "base_module": "line_6",
        "cooking_type": "induction"
    }
}
→ Response:
{
    "steps": [
        {"code": "base_module", "selected": "line_6", "locked": true},
        {"code": "cooking_type", "selected": "induction", "locked": true},
        {
            "code": "undercounter",
            "available_options": [
                {"code": "cooling_unit", "name": {"de": "Kühlunterbau"}, "price_modifier": {"type": "absolute", "value": 2800}},
                {"code": "drawer_block", "name": {"de": "Schubladenblock"}, "price_modifier": {"type": "absolute", "value": 1950}},
                {"code": "oven_gn21", "name": {"de": "Backofen GN 2/1"}, "price_modifier": {"type": "absolute", "value": 4200}}
            ]
        },
        {"code": "material", "awaiting_selection": true},
        {"code": "options", "awaiting_selection": true},
        {"code": "connections", "awaiting_selection": true}
    ],
    "current_price": {"net": 12500.00, "currency": "CHF"},
    "is_complete": false
}

# Konfiguration abschliessen + bestellen
POST /api/v1/cart/items
{
    "product_id": "uuid-konfigurator",
    "configuration": {
        "base_module": "line_6",
        "cooking_type": "induction",
        "undercounter": "cooling_unit",
        "material": "aisi_304",
        "options": ["splash_guard", "warming_bridge"],
        "connections": "power_63a"
    }
}
```

### Frontend-UX-Anforderungen

1. **Wizard/Stepper:** Step-by-Step Führung durch die Konfiguration. Fortschrittsanzeige (Step 2 von 6).
2. **Visuelle Vorschau:** Konfiguration als Bild/3D-Ansicht die sich mit jeder Auswahl aktualisiert. Bei Grossküche: Kochblock-Visualisierung mit gewählten Modulen.
3. **Zurück-Navigation:** Jederzeit zu früheren Steps zurück. Bei Änderung: nachfolgende Auswahlen zurücksetzen mit Hinweis.
4. **Nicht-verfügbare Optionen:** Ausgegraut mit Erklärung anzeigen, nicht verstecken. «Kühlunterbau ist bei Gas-Kochfeld nicht verfügbar.»
5. **Live-Preis:** Gesamtpreis nach jeder Auswahl aktualisieren. Aufschlüsselung zeigen.
6. **Zusammenfassung:** Vor dem Warenkorb: komplette Zusammenfassung aller gewählten Optionen mit Preisdetails.
7. **Konfiguration speichern/teilen:** B2B-Kunden wollen Konfigurationen speichern, teilen, als Angebot anfordern.

### Preisbildung

- **Basispreis** des Konfigurators (z.B. Grundpreis Kochfeld-Zeile)
- **+ Preismodifier** pro gewählter Option (absolut oder prozentual)
- **+ Parametrischer Anteil** (z.B. Konfektionsgrösse bei Bekleidungs-Konfigurator)
- **= Gesamtpreis** der Konfiguration
- Kundenspezifische Rabatte auf den Gesamtpreis möglich

---

## 6. Automatische Varianten-Erkennung und Produktmanagement-Assistenz

### Motivation

Das System soll Produktmanager aktiv dabei unterstützen, aus einfachen (flachen) Produkten Variantenprodukte zu erstellen. Besonders relevant bei:

- **Migration aus Altsystemen** (SAP, CSV-Import): Historisch gewachsene Kataloge, in denen alles als Einzelprodukt mit eigener SKU angelegt wurde — obwohl es sich faktisch um Varianten eines Master-Produkts handelt
- **Kataloge mit tausenden Produkten**: Manuelle Analyse, welche Produkte zusammengehören, ist ab wenigen hundert Produkten nicht mehr wirtschaftlich

**Mehrwert für den Produktmanager:** Statt manuell 5'000 Produkte durchzugehen und Gruppen zu bilden, bekommt er eine kuratierte Liste von Vorschlägen: «Diese 12 Kabel könnten ein Variantenprodukt mit den Achsen Aderzahl × Querschnitt sein.» Er prüft, bestätigt oder verwirft — und spart Wochen an Arbeit.

### 6.1 Ähnlichkeitserkennung / Clustering

Das System analysiert den bestehenden Produktkatalog und erkennt automatisch Gruppen von Produkten, die potenzielle Varianten eines gemeinsamen Master-Produkts sind.

**Erkennungskriterien:**

| Kriterium | Gewichtung | Beispiel |
|-----------|------------|---------|
| **Namensähnlichkeit** | Hoch | «NYM-J 3x1.5mm² 50m», «NYM-J 3x2.5mm² 50m», «NYM-J 5x1.5mm² 50m» |
| **Gleiche Kategorie** | Mittel | Alle in Kategorie «Installationskabel» |
| **Ähnliche Attribute** | Mittel | Gleicher Hersteller, gleiche Materialart, gleiche Zertifizierung |
| **SKU-Muster** | Hoch | SKUs wie `NYM-3x15-50`, `NYM-3x25-50`, `NYM-5x15-50` — gemeinsamer Präfix, variierende Suffixe |
| **Preisstruktur** | Niedrig | Ähnlicher Preisbereich (nicht 10× teurer/günstiger) |

**Algorithmus-Optionen:**

| Ansatz | Vorteile | Nachteile | Empfehlung |
|--------|----------|-----------|------------|
| **Regelbasiert** (String-Matching, SKU-Pattern) | Deterministisch, erklärbar, schnell | Starr, erkennt nur offensichtliche Muster | Phase 1 — Sofort umsetzbar |
| **String-Similarity** (Levenshtein, Jaro-Winkler, n-gram) | Findet auch unsaubere Benennungen | Viele False Positives bei kurzen Namen | Kombination mit regelbasiert |
| **NLP/Embeddings** (Sentence-Transformers, OpenAI Embeddings) | Erkennt semantische Ähnlichkeit («Schraubendreher» ≈ «Schraubenzieher») | Overhead, Black Box, Kosten | Phase 2 — Für feinere Analyse |
| **Kombination** | Bestes Ergebnis | Komplexer zu implementieren | Langfristig empfohlen |

**Ergebnis:** Gruppierungsvorschläge mit Konfidenz-Score:

```json
{
  "cluster_id": "uuid",
  "confidence": 0.92,
  "suggested_master_name": "NYM-J Installationskabel",
  "products": [
    {"id": "uuid-1", "name": "NYM-J 3x1.5mm² 50m", "sku": "NYM-3x15-50"},
    {"id": "uuid-2", "name": "NYM-J 3x2.5mm² 50m", "sku": "NYM-3x25-50"},
    {"id": "uuid-3", "name": "NYM-J 5x1.5mm² 50m", "sku": "NYM-5x15-50"},
    {"id": "uuid-4", "name": "NYM-J 5x2.5mm² 50m", "sku": "NYM-5x25-50"},
    {"id": "uuid-5", "name": "NYM-J 3x1.5mm² 100m", "sku": "NYM-3x15-100"},
    {"id": "uuid-6", "name": "NYM-J 3x2.5mm² 100m", "sku": "NYM-3x25-100"}
  ],
  "suggested_axes": [
    {"code": "wire_count", "name": "Aderzahl", "values": ["3", "5"]},
    {"code": "cross_section", "name": "Querschnitt", "values": ["1.5mm²", "2.5mm²"]},
    {"code": "length", "name": "Länge", "values": ["50m", "100m"]}
  ]
}
```

### 6.2 Vorschlagsgenerierung

Für jede erkannte Gruppe generiert das System einen konkreten Umwandlungsvorschlag:

**Master-Produkt-Bestimmung:**
- Das Produkt mit der vollständigsten Beschreibung, den meisten Attributen und den besten Bildern wird als Master vorgeschlagen
- Alternativ: ein neues Master-Produkt aus den gemeinsamen Attributen aller Gruppenmitglieder erzeugen

**Achsen-Erkennung:**
- System analysiert die Unterschiede zwischen den Produkten und leitet daraus Variantenachsen ab
- Beispiel: Aus «Nitrilhandschuh S», «Nitrilhandschuh M», «Nitrilhandschuh L» ergibt sich Achse «Grösse» mit Werten S, M, L

**Attribut-Aufteilung:**

| Attribut-Typ | Zuordnung | Beispiel |
|-------------|-----------|---------|
| Identisch bei allen Produkten | → Master-Attribut | Hersteller, Material, Farbe (wenn gleich) |
| Variiert zwischen Produkten | → Variantenkoordinate (Achse) | Grösse, Durchmesser, Länge |
| Variiert, aber keine sinnvolle Achse | → Bleibt an der einzelnen Variante | Individuelle Beschreibungstexte |

**Review-Vorschlag (was der Produktmanager sieht):**

```
┌─────────────────────────────────────────────────────────────┐
│ Vorschlag #14: Nitrilhandschuh (Konfidenz: 95%)            │
│                                                             │
│ Master-Produkt: "Nitrilhandschuh puderfrei, blau"          │
│ Achsen: Grösse (S, M, L, XL)                              │
│                                                             │
│ Gemeinsame Attribute:                                       │
│   Hersteller: SafeGrip    Material: Nitril                 │
│   Farbe: Blau              Normen: EN 374, EN 455          │
│                                                             │
│ Varianten:                                                  │
│   S  → SKU: SG-NIT-S   │ Preis: 12.90                     │
│   M  → SKU: SG-NIT-M   │ Preis: 12.90                     │
│   L  → SKU: SG-NIT-L   │ Preis: 12.90                     │
│   XL → SKU: SG-NIT-XL  │ Preis: 13.50                     │
│                                                             │
│ [✓ Bestätigen]  [✎ Anpassen]  [✗ Verwerfen]               │
└─────────────────────────────────────────────────────────────┘
```

### 6.3 Workflow / UX

**Dashboard-Integration:**

- Widget auf dem Produktmanager-Dashboard: «**47 potenzielle Variantengruppen** erkannt — Jetzt prüfen»
- Fortschrittsanzeige: «12 von 47 Gruppen bearbeitet»
- Filter: Nach Konfidenz, Kategorie, Anzahl Produkte pro Gruppe

**Detailansicht pro Gruppe:**

1. **Gruppierung anzeigen:** Alle erkannten Produkte als Karten oder Liste
2. **Drag & Drop:** Produkte zwischen Gruppen verschieben, aus Gruppen entfernen, neue hinzufügen
3. **Achsen-Editor:** Vorgeschlagene Achsen anpassen, umbenennen, zusammenlegen oder neue hinzufügen
4. **Vorschau:** «So würde das Variantenprodukt im Shop aussehen» — inklusive Achsen-Selektor, Preismatrix
5. **Vergleichsansicht:** Attribute aller Produkte nebeneinander, Unterschiede farblich markiert

**Batch-Aktionen:**

- Mehrere Gruppen auf einmal bestätigen (bei hoher Konfidenz)
- Filter: «Alle Gruppen mit Konfidenz > 90% bestätigen»
- Fortschrittsbalken bei Massenkonvertierung

**Rückgängig-Funktion:**

- Variantenprodukt jederzeit wieder in Einzelprodukte auflösen
- Audit-Log: Wer hat wann welche Gruppe konvertiert
- Soft-Delete der ursprünglichen Einzelprodukte (30 Tage wiederherstellbar)

### 6.4 Weitere Assistenz-Features

Über die Varianten-Erkennung hinaus kann das System den Produktmanager bei weiteren Aufgaben unterstützen:

**Duplikat-Erkennung:**
- Fast identische Produkte erkennen (>95% Attributübereinstimmung, ähnlicher Name)
- Vorschlag: Zusammenführen oder eines deaktivieren
- Beispiel: «Hydraulikzylinder 50mm» und «Hydraulik-Zylinder 50 mm» — wahrscheinlich dasselbe Produkt

**Attribut-Normalisierung:**
- Inkonsistente Schreibweisen erkennen und vereinheitlichen
- «50 mm» vs. «50mm» vs. «50 Millimeter» → normalisiert zu «50 mm»
- «Ø 63» vs. «DN63» vs. «63mm Durchmesser» → einheitliches Format
- Vorschlag mit Vorher/Nachher-Vergleich, Produktmanager bestätigt

**Kategorie-Vorschläge:**
- Unkategorisierte Produkte automatisch einer Kategorie zuordnen (basierend auf Namen, Attributen, Ähnlichkeit zu kategorisierten Produkten)
- Falsch kategorisierte Produkte erkennen («Schraubendreher» in Kategorie «Kabel»)

**Preisanomalien:**
- Produkte identifizieren, deren Preis signifikant vom Durchschnitt ähnlicher Produkte abweicht
- Beispiel: «Nitrilhandschuh L kostet 129.00 CHF statt 12.90 CHF — vermutlich Tippfehler»
- Dashboard-Warnung mit Ein-Klick-Korrektur

**Lücken in der Variantenmatrix:**
- Bei bestehenden Variantenprodukten fehlende Kombinationen erkennen
- Beispiel: «Poloshirt Hakro — Grösse M fehlt in Farbe Rot. Alle anderen Grössen sind vorhanden.»
- Vorschlag: Fehlende Variante anlegen oder bewusst als nicht verfügbar markieren

### 6.5 Technische Umsetzung

**Architektur-Entscheid: Eigener Service vs. Teil des Catalog Service**

| Option | Vorteile | Nachteile |
|--------|----------|-----------|
| **Eigener Microservice** («Product Intelligence Service») | Unabhängig skalierbar, eigene Deployment-Zyklen, kann ML-Modelle nutzen ohne Catalog zu belasten | Zusätzlicher Service zu betreiben, Netzwerk-Overhead |
| **Teil vom Catalog Service** | Kein zusätzlicher Service, direkter DB-Zugriff | Kopplung, Catalog wird komplexer, schwerer skalierbar bei ML-Last |

**Empfehlung:** Eigener Service «Product Intelligence Service», der den Catalog Service via API konsumiert. Gründe: ML-Workloads haben andere Ressourcen-Anforderungen (CPU/RAM-intensiv), unabhängige Release-Zyklen, saubere Trennung.

**Batch-Job vs. Echtzeit:**

| Modus | Einsatz |
|-------|---------|
| **Batch-Job** (geplant, z.B. nächtlich) | Initiale Analyse des gesamten Katalogs, Re-Analyse nach grossen Imports |
| **Event-basiert** (bei Produktänderung) | Neue Produkte prüfen, ob sie in eine bestehende Gruppe passen |
| **On-Demand** (manuell ausgelöst) | Produktmanager startet Analyse für bestimmte Kategorie oder Import-Batch |

**ML/AI vs. Regelbasiert:**

| Ansatz | Stärken | Schwächen | Empfehlung |
|--------|---------|-----------|------------|
| **Regelbasiert** | Deterministisch, erklärbar, kein Training nötig, schnell | Erkennt nur explizit programmierte Muster | Phase 1: SKU-Pattern, Name-Prefix, Kategorie-Match |
| **ML/AI** | Erkennt subtile Muster, lernt aus Feedback | Black Box, benötigt Trainingsdaten, Infrastruktur | Phase 2: Embeddings für semantische Ähnlichkeit |
| **Hybrid** | Bestes aus beiden Welten | Komplexer | Langfristig: Regeln als Vorfilter, ML für Feinanalyse |

**Performance bei 100k+ Produkten:**

- Vorfilterung über Kategorie und SKU-Muster (reduziert Vergleiche drastisch)
- Paarweise Vergleiche nur innerhalb einer Kategorie (nicht über den gesamten Katalog)
- Batch-Verarbeitung mit Fortschrittsanzeige
- Ergebnisse cachen — nur Delta bei neuen/geänderten Produkten neu analysieren
- Bei Embeddings: Vektor-Datenbank (pgvector oder dediziert) für effiziente Similarity-Search

### Branchenbeispiele

**Elektro — Installationskabel:**

```
Ausgangslage (Einzelprodukte):
  "NYM-J 3x1.5mm² 50m"   SKU: NYM-3x15-50
  "NYM-J 3x2.5mm² 50m"   SKU: NYM-3x25-50
  "NYM-J 5x1.5mm² 50m"   SKU: NYM-5x15-50
  "NYM-J 5x2.5mm² 50m"   SKU: NYM-5x25-50
  "NYM-J 3x1.5mm² 100m"  SKU: NYM-3x15-100
  "NYM-J 3x2.5mm² 100m"  SKU: NYM-3x25-100

Erkannt als Variantenprodukt:
  Master: "NYM-J Installationskabel"
  Achse 1: Aderzahl × Querschnitt → 3x1.5mm², 3x2.5mm², 5x1.5mm², 5x2.5mm²
  Achse 2: Ringlänge → 50m, 100m
```

**Arbeitsschutz — Handschuhe:**

```
Ausgangslage:
  "Nitrilhandschuh puderfrei S"  SKU: SG-NIT-S
  "Nitrilhandschuh puderfrei M"  SKU: SG-NIT-M
  "Nitrilhandschuh puderfrei L"  SKU: SG-NIT-L
  "Nitrilhandschuh puderfrei XL" SKU: SG-NIT-XL

Erkannt als Variantenprodukt:
  Master: "Nitrilhandschuh puderfrei"
  Achse: Grösse → S, M, L, XL
```

**Verpackung — Faltkartons:**

```
Ausgangslage:
  "Faltkarton 300x200x150 braun"  SKU: FK-300200150-BR
  "Faltkarton 300x200x150 weiss"  SKU: FK-300200150-WS
  "Faltkarton 400x300x200 braun"  SKU: FK-400300200-BR
  "Faltkarton 400x300x200 weiss"  SKU: FK-400300200-WS

Erkannt als Variantenprodukt:
  Master: "Faltkarton FEFCO 0201"
  Achse 1: Dimension → 300×200×150mm, 400×300×200mm
  Achse 2: Farbe → Braun, Weiss
```

---

## Übergreifende Themen

### Typ-Interaktionen

Die Produkttypen können **kombiniert** auftreten:

| Kombination | Beispiel | Umsetzung |
|-------------|----------|-----------|
| Bundle aus Varianten | Pneumatik-Set mit Zylindern in verschiedenen Hublängen | Bundle-Item verweist auf Variant (nicht Master) |
| Bundle mit parametrischem Produkt | Kabel-Set: 3 Kabeltypen mit individuellen Längen | Bundle-Item hat zusätzlich `parameters` |
| Konfigurator mit Varianten-Steps | Rack-Konfigurator: PDU als Variantenprodukt (16A/32A) | ConfiguratorOption verweist auf VariantMaster |
| Konfigurator mit parametrischem Step | Berufsbekleidung: Grösse als parametrischer Step | Step mit `step_type: "parametric"` |

**Technische Konsequenz:** Das Warenkorb-Item muss generisch genug sein, um alle Typen abzubilden.

### Warenkorb-Auswirkungen

```go
type CartItem struct {
    ID          uuid.UUID              `json:"id"`
    ProductID   uuid.UUID              `json:"product_id"`
    ProductType string                 `json:"product_type"` // simple, variant, parametric, complex, bundle, configurator
    Quantity    float64                `json:"quantity"`

    // Typ-spezifische Daten
    VariantID      *uuid.UUID         `json:"variant_id,omitempty"`       // Für variant
    Parameters     map[string]any     `json:"parameters,omitempty"`       // Für parametric/complex
    BundleItems    []CartBundleItem   `json:"bundle_items,omitempty"`     // Für bundle
    Configuration  map[string]any     `json:"configuration,omitempty"`    // Für configurator

    // Berechnete Werte
    UnitPrice   float64               `json:"unit_price"`
    TotalPrice  float64               `json:"total_price"`
    Description string                `json:"description"` // Menschenlesbare Zusammenfassung

    // Validierung
    IsValid     bool                  `json:"is_valid"`
    Errors      []string              `json:"errors,omitempty"`
}
```

**Was steht im Warenkorb?**

| Typ | Warenkorb-Darstellung |
|-----|-----------------------|
| Simple | Produkt + Menge |
| Variant | Master-Name + gewählte Achsenwerte + Menge. Z.B. «Hydraulikzylinder CDT3 — Ø63mm, Hub 300mm, 250 bar — 2 Stk» |
| Parametric | Produkt + eingegebene Parameter. Z.B. «Kabel NYM-J 3×1.5mm² — 45m — 1 Stk» |
| Complex | Wie Parametric, mit allen gewählten Optionen |
| Bundle | Bundle-Name + aufgeklappte Einzelpositionen. Z.B. «Pneumatik-Starterpaket: 5× Zylinder, 10× Steckverbinder, 5m Schlauch, 1× Wartungseinheit» |
| Configurator | Konfigurator-Name + Zusammenfassung aller Steps. Z.B. «Kochblock — 6er-Zeile, Induktion, Kühlunterbau, Edelstahl 304, Spritzschutz» |

### SAP/ERP-Integration

**Herausforderung:** SAP kennt diese Produkttypen nicht nativ. Bestellungen müssen in SAP-kompatible Strukturen übersetzt werden.

| Produkttyp | SAP-Übermittlung |
|------------|------------------|
| Simple | 1 Position: Material + Menge |
| Variant | 1 Position: Varianten-SKU + Menge (jede Variante hat eigene SAP-Materialnummer) |
| Parametric | 1 Position: Material + Menge + Konfigurationsmerkmale (CUOBJ / Merkmalsbewertung) |
| Complex | Wie Parametric (SAP-Variantenkonfiguration) |
| Bundle | N Positionen: Jede Komponente als eigene Auftragsposition, verknüpft über Bundle-Referenz |
| Configurator | Stückliste (BOM) in SAP oder N Positionen mit Konfigurationsmerkmalen |

**Offene Fragen:**
- [ ] Unterstützt das vorhandene SAP-System Variantenkonfiguration (LO-VC)?
- [ ] Gibt es SAP-Materialnummern für jede Variante oder nur für Master-Produkte?
- [ ] Wie werden parametrische Produkte heute in SAP abgebildet?
- [ ] Werden Bundles als Stücklisten (BOM) oder als Einzelpositionen übertragen?

### Suchbarkeit (Meilisearch)

| Produkttyp | Index-Strategie |
|------------|-----------------|
| Simple | Standard-Index |
| Variant Master | Master + alle Varianten-Werte als durchsuchbare Attribute. Filter auf Achsenwerte. |
| Variant (einzeln) | Nicht separat indexiert — nur über Master findbar |
| Parametric | Master indexiert mit Parameter-Bereichen als Attribute (z.B. «Kabel ab 1m bis 500m konfigurierbar») |
| Complex | Wie Parametric |
| Bundle | Bundle-Name + Namen aller enthaltenen Produkte indexiert. Filter: «Bundles verfügbar» |
| Configurator | Konfigurator-Name + alle möglichen Optionen als durchsuchbare Attribute |

**Meilisearch-Dokument für Variantenprodukt:**

```json
{
    "id": "master-uuid",
    "type": "variant_master",
    "name": "Bosch Rexroth Hydraulikzylinder CDT3",
    "variant_values": {
        "diameter": ["40mm", "50mm", "63mm", "80mm"],
        "stroke": ["100mm", "200mm", "300mm", "500mm"],
        "pressure": ["160bar", "250bar"]
    },
    "searchable_text": "Bosch Rexroth Hydraulikzylinder CDT3 40mm 50mm 63mm 80mm 100mm 200mm 300mm 500mm 160bar 250bar",
    "price_range": {"min": 485.00, "max": 2890.00}
}
```

### Priorisierung / Implementierungsreihenfolge

| Phase | Produkttyp | Begründung | Geschätzter Aufwand |
|-------|------------|------------|---------------------|
| **Phase 1** | **Variantenprodukte** | Häufigster Typ im B2B. Deckt 60-70% der Produkte ab (Industriekomponenten, Elektronik, Verpackung, Textil). Grundlage für alles weitere. | 3–4 Wochen |
| **Phase 2** | **Parametrisierbare Produkte** | Meterware, Zuschnitte und Konfektionierung sind Kerngeschäft vieler B2B-Branchen. Differenzierungsmerkmal gegenüber reinen Online-Shops. | 3–4 Wochen |
| **Phase 3** | **Bundle-Produkte** | Relativ einfach zu implementieren. Sofortiger Business-Value durch Upselling (Startersets, Komplett-Pakete). | 2–3 Wochen |
| **Phase 4** | **Komplexe Produkte** | Erweiterung von Phase 2. Regel-Engine ist der Hauptaufwand. Relevant für Schaltschränke, Hydraulik, Server-Konfigurationen. | 3–4 Wochen |
| **Phase 5** | **Konfigurator-Produkte** | Höchste Komplexität. Erst wenn die Basis steht. Kann mandantenspezifisch eingeführt werden (Grossküche, IT-Rack, Bekleidung). | 4–6 Wochen |

**Empfehlung:** Phase 1 und 2 parallel starten, da unterschiedliche Domänen. Phase 3 kann teilweise parallel zu Phase 2.

### Datenmodell-Erweiterung am Product

Das bestehende `Product`-Struct wird um ein `Type`-Feld erweitert:

```go
type ProductType string

const (
    ProductTypeSimple       ProductType = "simple"
    ProductTypeVariantMaster ProductType = "variant_master"
    ProductTypeParametric   ProductType = "parametric"
    ProductTypeComplex      ProductType = "complex"
    ProductTypeBundle       ProductType = "bundle"
    ProductTypeConfigurator ProductType = "configurator"
)

type Product struct {
    // ... bestehende Felder ...
    Type ProductType `json:"type" db:"type"`

    // Typ-spezifische Daten (lazy-loaded, nicht in DB-Tabelle products)
    Axes          []VariantAxis        `json:"axes,omitempty"`          // variant_master
    Variants      []ProductVariant     `json:"variants,omitempty"`      // variant_master
    Parameters    []ProductParameter   `json:"parameters,omitempty"`    // parametric, complex
    Rules         []ParameterRule      `json:"rules,omitempty"`         // complex
    BundleItems   []BundleItem         `json:"bundle_items,omitempty"`  // bundle
    BundlePricing *BundlePricing       `json:"bundle_pricing,omitempty"` // bundle
    Steps         []ConfiguratorStep   `json:"steps,omitempty"`         // configurator
    Dependencies  []ConfiguratorDependency `json:"dependencies,omitempty"` // configurator
}
```

### Migration

```sql
-- Schritt 1: Type-Feld hinzufügen
ALTER TABLE products ADD COLUMN type VARCHAR(20) NOT NULL DEFAULT 'simple';

-- Schritt 2: Neue Tabellen erstellen (siehe Datenmodelle oben)
-- variant_axes, variant_axis_values, product_variants
-- product_parameters, product_pricing_formulas
-- product_parameter_rules
-- bundle_items, bundle_pricing
-- configurator_steps, configurator_options, configurator_dependencies

-- Schritt 3: Bestehende Produkte bleiben type='simple'
```

---

## Offene Fragen

1. **Akeneo-Mapping:** Wie werden die verschiedenen Produkttypen in Akeneo abgebildet? Gibt es dort bereits Varianten-Support?
2. **SAP-Variantenkonfiguration:** Welche SAP-Module sind für Variantenkonfiguration im Einsatz (LO-VC, IPC)?
3. **Preise in SAP:** Werden Variantenpreise pro SKU in SAP gepflegt oder gibt es Preisregeln?
4. **Regel-Komplexität:** Wie komplex sind die realen Regeln bei komplexen Produkten? Reicht ein einfaches Condition/Effect-System oder braucht es einen vollwertigen Regelinterpreter?
5. **Konfigurator-Persistenz:** Sollen angefangene Konfigurationen serverseitig gespeichert werden (Resume nach Session-Verlust)?
6. **Masseinheiten:** Welche Einheiten sind im Einsatz? (mm, cm, m, m², m³, kg, L, Stk, Lfm, Packung, bar, µm) — Einheitensystem standardisieren.
7. **Performance:** Bei grossen Variantenmatrizen (>100 Varianten): Alle Varianten im initialen Response liefern oder lazy-load?
