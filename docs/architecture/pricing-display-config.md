# Preisanzeige-Konfiguration

Stand: 2026-02-14

## Übersicht

Die Preisanzeige in Gondolia ist **pro Tenant konfigurierbar**. Jeder Tenant entscheidet selbst, welche Preisinformationen in welchem Kontext (anonym, eingeloggt, kundenspezifisch) angezeigt werden. Die Konfiguration wird in der Datenbank gespeichert und über ein Admin-UI verwaltet.

### Warum tenant-spezifisch?

B2B-Shops haben fundamental unterschiedliche Preisstrategien:

| Tenant-Typ | Anonyme Besucher | Eingeloggte Kunden |
|------------|-------------------|---------------------|
| **Industriebedarf** (z.B. FlexoTech AG) | Keine Preise → «Bitte einloggen» | Kundenspezifisch aus ERP |
| **Laborbedarf** (z.B. LabSupply GmbH) | Listenpreise öffentlich | Rahmenvertrag-Preise |
| **Verpackung** (z.B. PackDirect AG) | Volle Staffelpreise öffentlich | Zusätzlich Kundenrabatt |

---

## 1. Preisanzeige-Modi

### 1.1 Ausgeloggter Zustand (anonym)

Der Tenant wählt **einen** der folgenden Modi:

| Modus | Code | Anzeige | SEO-Wirkung | Einsatz |
|-------|------|---------|-------------|---------|
| **Keine Preise** | `none` | «Preis auf Anfrage» oder «Bitte einloggen für Preise» | ❌ Kein Preis in Structured Data | Individuelle Preisgestaltung, Wettbewerbsschutz |
| **Listenpreise** | `list` | UVP/Bruttopreis ohne Rabatte | ✅ Preis in JSON-LD → Rich Snippets | Standardprodukte mit öffentlichem Listenpreis |
| **Ab-Preise** | `from` | «ab CHF X.XX» (niedrigster Staffelpreis) | ✅ `lowPrice` in AggregateOffer | Staffelpreise, Einstiegspreis kommunizieren |
| **Volle Preise** | `full` | Komplette Staffelpreis-Tabelle | ✅ Vollständige Preisinfo | Maximale Transparenz, Conversion-Optimierung |

**Verhalten bei `none`:**

```
┌─────────────────────────────────┐
│  Hydraulikzylinder CDT3 50mm   │
│                                 │
│  ┌───────────────────────────┐  │
│  │ 🔒 Preis auf Anfrage      │  │
│  │                           │  │
│  │ [Einloggen für Preise]    │  │
│  │ [Preis anfragen]          │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘
```

- **«Einloggen für Preise»** → Redirect zum Login
- **«Preis anfragen»** → Kontaktformular mit vorausgefülltem Produktnamen

**Verhalten bei `from`:**

```
┌─────────────────────────────────┐
│  Faltkarton 400×300×200mm       │
│                                 │
│  ab CHF 0.85 / Stück            │
│  Staffelpreise nach Login       │
│                                 │
│  [Einloggen für Ihren Preis]    │
└─────────────────────────────────┘
```

**Verhalten bei `full`:**

```
┌─────────────────────────────────┐
│  Faltkarton 400×300×200mm       │
│                                 │
│  Staffelpreise (netto):         │
│  ┌──────────┬──────────────┐    │
│  │ Menge    │ Preis/Stück  │    │
│  ├──────────┼──────────────┤    │
│  │ 1–49     │ CHF 1.20     │    │
│  │ 50–199   │ CHF 0.95     │    │
│  │ 200–499  │ CHF 0.88     │    │
│  │ ab 500   │ CHF 0.85     │    │
│  └──────────┴──────────────┘    │
│  zzgl. 8.1% MwSt.              │
└─────────────────────────────────┘
```

### 1.2 Eingeloggter Zustand (authentifiziert)

| Modus | Code | Anzeige | Preisquelle |
|-------|------|---------|-------------|
| **Listenpreise** | `list` | Katalogpreise ohne Kundenrabatte | Gondolia-DB (`tier_prices`) |
| **Kundenkonditionen** | `customer` | Individuelle Preise basierend auf Kundengruppe/Vertrag | Gondolia-DB (`customer_prices`) |
| **ERP-Livepreise** | `erp_live` | Echtzeit-Preis aus SAP/ERP | ERP-API mit Fallback |

**Kundenkonditionen — Priorisierung:**

```
1. Sonderkondition (kundenspezifisch, manuell gepflegt)
2. Rahmenvertrag-Preis (vertraglich vereinbart, Laufzeit)
3. Kundengruppen-Preis (Preisliste der Kundengruppe)
4. Staffelpreis (mengenabhängig, aus Katalog)
5. Listenpreis (Fallback)
```

Die erste zutreffende Kondition gewinnt. Ausnahme: Staffelrabatte werden **additiv** auf den konditionierten Preis angewendet, wenn `stack_volume_discounts = true` in der Tenant-Config.

### 1.3 Zusammenspiel anonym → eingeloggt

```
Anonymer Besucher                    Eingeloggter Kunde (Firma Müller AG)
─────────────────                    ────────────────────────────────────
"ab CHF 0.85"                        CHF 0.78 / Stück (Ihr Preis)
                                     ──────────────────────────────
                                     Listenpreis: CHF 1.20 (durchgestrichen)
                                     Ihr Rabatt: -35%
                                     
                                     Staffelpreise:
                                     1–49:    CHF 0.78
                                     50–199:  CHF 0.72
                                     200–499: CHF 0.68
                                     ab 500:  CHF 0.65
```

---

## 2. Preisquellen

### 2.1 Architektur

```
                    ┌─────────────────────────────────────────────┐
                    │              Pricing Service                 │
                    │                                             │
   Anfrage ───────▶│  1. Tenant-Config laden                     │
                    │  2. Kontext bestimmen (anon/auth/customer)  │
                    │  3. Preisquelle auswählen                   │
                    │  4. Preis holen + transformieren             │
                    │  5. Anzeige-Regeln anwenden                 │
                    │                                             │
                    │         ┌──────────┐                        │
                    │         │  Cache    │◀─── Redis              │
                    │         └────┬─────┘                        │
                    │              │                               │
                    │    ┌─────────┼─────────┐                    │
                    │    ▼         ▼         ▼                    │
                    │  Katalog  Kunden-   ERP-Live                │
                    │  Preise   Preise    Preise                  │
                    │  (DB)     (DB)      (SAP API)               │
                    └─────────────────────────────────────────────┘
```

### 2.2 Katalogpreise (Gondolia-DB)

Basis-Preise direkt in Gondolia gepflegt. Staffelpreise wie bereits im Produkt-Modell vorgesehen (siehe [Produkttypen-Architektur](product-types.md), Abschnitt Preisbildung).

```sql
CREATE TABLE catalog_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),  -- NULL = Preis gilt für alle Varianten
    currency VARCHAR(3) NOT NULL DEFAULT 'CHF',
    
    -- Staffelpreise
    min_quantity INTEGER NOT NULL DEFAULT 1,
    price_net DECIMAL(12,4) NOT NULL,
    price_gross DECIMAL(12,4),                         -- Optional, berechnet sich aus net + VAT
    
    valid_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(tenant_id, product_id, variant_id, min_quantity, currency, valid_from)
);

CREATE INDEX idx_catalog_prices_lookup 
    ON catalog_prices(tenant_id, product_id, currency) 
    WHERE valid_to IS NULL OR valid_to > NOW();
```

### 2.3 Kundenspezifische Preise (Gondolia-DB)

```sql
CREATE TABLE customer_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Zuordnung: Kunde ODER Kundengruppe (nicht beides)
    customer_id UUID,                                  -- Spezifischer Kunde
    customer_group_id UUID,                            -- Kundengruppe / Preisliste
    
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    currency VARCHAR(3) NOT NULL DEFAULT 'CHF',
    
    -- Preistyp
    price_type VARCHAR(20) NOT NULL,                   -- 'fixed', 'discount_percent', 'discount_absolute'
    price_net DECIMAL(12,4),                           -- Bei fixed
    discount_percent DECIMAL(5,2),                     -- Bei discount_percent
    discount_absolute DECIMAL(12,4),                   -- Bei discount_absolute
    
    -- Staffel
    min_quantity INTEGER NOT NULL DEFAULT 1,
    
    -- Gültigkeit
    valid_from TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to TIMESTAMPTZ,
    
    -- Quelle
    source VARCHAR(20) NOT NULL DEFAULT 'manual',      -- 'manual', 'erp_import', 'contract'
    contract_reference VARCHAR(100),                    -- Rahmenvertrag-Nummer
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CHECK (customer_id IS NOT NULL OR customer_group_id IS NOT NULL)
);

CREATE INDEX idx_customer_prices_customer 
    ON customer_prices(tenant_id, customer_id, product_id, currency)
    WHERE valid_to IS NULL OR valid_to > NOW();

CREATE INDEX idx_customer_prices_group 
    ON customer_prices(tenant_id, customer_group_id, product_id, currency)
    WHERE valid_to IS NULL OR valid_to > NOW();
```

### 2.4 ERP-Import (Batch)

Regelmässiger Import von Preisen aus SAP/ERP:

```
SAP/ERP ──── (nächtlich, IDoc/RFC) ────▶ Import-Job ────▶ customer_prices (source='erp_import')
```

- **Frequenz:** Konfigurierbar pro Tenant (täglich, stündlich)
- **Format:** SAP-Konditionstabellen (KONV, KONH) → Mapping auf `customer_prices`
- **Delta-Import:** Nur geänderte Konditionen importieren (Änderungsdatum in SAP)
- **Validierung:** Negative Preise, fehlende Produkt-Zuordnung → Fehler-Log, nicht importieren
- **Monitoring:** Import-Dashboard mit Anzahl importierter/fehlerhafter Konditionen

### 2.5 ERP-Livepreise

Echtzeit-Abfrage beim Seitenaufruf oder Warenkorb-Aktualisierung:

```
Browser ──▶ Gondolia API ──▶ Pricing Service ──▶ SAP (RFC/REST)
                                    │
                                    ├── Cache Hit? → Cached Preis zurückgeben
                                    │
                                    └── Cache Miss → SAP anfragen
                                         │
                                         ├── SAP erreichbar → Preis cachen + zurückgeben
                                         │
                                         └── SAP nicht erreichbar → Fallback-Strategie
```

**Bulk-Abfrage (kritisch für Performance):**

```go
// NICHT: 50 einzelne SAP-Calls für 50 Warenkorb-Positionen
// SONDERN: Ein Bulk-Call mit allen Positionen

type ERPPriceRequest struct {
    TenantID   uuid.UUID           `json:"tenant_id"`
    CustomerID uuid.UUID           `json:"customer_id"`
    Items      []ERPPriceLineItem  `json:"items"`
}

type ERPPriceLineItem struct {
    ProductSKU string  `json:"product_sku"`
    VariantSKU string  `json:"variant_sku,omitempty"`
    Quantity   float64 `json:"quantity"`
}

type ERPPriceResponse struct {
    Items []ERPPriceResult `json:"items"`
    Timestamp time.Time    `json:"timestamp"`
}

type ERPPriceResult struct {
    ProductSKU   string   `json:"product_sku"`
    PriceNet     float64  `json:"price_net"`
    Currency     string   `json:"currency"`
    ConditionType string  `json:"condition_type"` // z.B. "PR00", "ZK01"
    ValidUntil   *time.Time `json:"valid_until,omitempty"`
}
```

**SAP-Integration:** Nutzung von BAPI_PRICING_SIMULATE oder custom RFC für Bulk-Preisermittlung. Detaillierte SAP-Integration siehe [SAP-Integration](sap-integration.md).

**Fallback-Strategie:**

| Strategie | Code | Verhalten | Einsatz |
|-----------|------|-----------|---------|
| **Verstecken** | `hide` | Kein Preis anzeigen, «Preis wird ermittelt» | Wenn Preis ohne ERP nicht sinnvoll ist |
| **Cached** | `cached` | Letzten bekannten Preis anzeigen (mit Hinweis «Preis evtl. nicht aktuell») | Standard-Empfehlung |
| **Listenpreis** | `list_price` | Katalog-Listenpreis als Fallback | Wenn Listenpreise gepflegt sind |

**Caching-Strategie:**

```
Cache-Key:  price:{tenant_id}:{customer_id}:{product_sku}:{variant_sku}:{quantity}
TTL:        Konfigurierbar pro Tenant (Default: 300 Sekunden)
Storage:    Redis
```

| Szenario | Cache-Verhalten |
|----------|-----------------|
| Produktseite aufrufen | Cache prüfen → Hit: sofort liefern. Miss: SAP anfragen, cachen |
| Warenkorb öffnen | Bulk-Anfrage für alle Positionen, Cache pro Position prüfen |
| Menge ändern | Neuer Cache-Key (Menge ist Teil des Keys), SAP anfragen |
| Bestellung abschliessen | **Immer** frischen Preis aus SAP holen (kein Cache) |

**Cache-Invalidierung:**

- TTL-basiert (automatisch nach Ablauf)
- Event-basiert: SAP-Konditionsänderung → Invalidierung via Message Queue
- Manuell: Admin kann Cache pro Tenant/Kunde flushen

### 2.6 Hybride Preisquelle

Kombination aus Katalogpreis und ERP:

```
1. ERP-Preis anfragen
2. ERP erreichbar?
   JA  → ERP-Preis verwenden, cachen
   NEIN → Fallback-Strategie anwenden:
          - 'cached': Letzten ERP-Preis aus Cache
          - 'list_price': Katalogpreis aus DB
          - 'hide': Kein Preis anzeigen
```

---

## 3. Tenant-Konfiguration

### 3.1 Datenmodell

```sql
CREATE TABLE tenant_pricing_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL UNIQUE,
    
    -- Preisanzeige-Modi
    anonymous_price_display VARCHAR(10) NOT NULL DEFAULT 'none',
        -- 'none': Keine Preise (Preis auf Anfrage)
        -- 'list': Listenpreise (UVP)
        -- 'from': Ab-Preise (niedrigster Staffelpreis)
        -- 'full': Volle Staffelpreise
    
    authenticated_price_display VARCHAR(10) NOT NULL DEFAULT 'list',
        -- 'list': Listenpreise ohne Kundenrabatte
        -- 'customer': Kundenspezifische Konditionen
        -- 'erp_live': Echtzeit aus ERP
    
    -- ERP-Preisquelle
    erp_price_source VARCHAR(15) NOT NULL DEFAULT 'none',
        -- 'none': Keine ERP-Preise
        -- 'batch_import': Regelmässiger Import
        -- 'live_query': Echtzeit-Abfrage
    
    erp_fallback_strategy VARCHAR(15) NOT NULL DEFAULT 'cached',
        -- 'hide': Kein Preis anzeigen
        -- 'cached': Letzten bekannten Preis
        -- 'list_price': Katalog-Listenpreis
    
    -- Caching
    price_cache_ttl_seconds INTEGER NOT NULL DEFAULT 300,
    
    -- Anzeige-Optionen
    show_discount_percentage BOOLEAN NOT NULL DEFAULT false,
        -- Zeigt "-15%" neben dem konditionierten Preis
    show_list_price_strikethrough BOOLEAN NOT NULL DEFAULT false,
        -- Zeigt durchgestrichenen Listenpreis neben dem Kundenpreis
    show_volume_discount_table BOOLEAN NOT NULL DEFAULT true,
        -- Staffelpreis-Tabelle auf Produktseite
    stack_volume_discounts BOOLEAN NOT NULL DEFAULT false,
        -- Staffelrabatte additiv auf Kundenkonditionen
    
    -- MwSt
    price_includes_vat BOOLEAN NOT NULL DEFAULT false,
        -- true: Preise inkl. MwSt (B2C-artig)
        -- false: Preise exkl. MwSt (B2B-Standard)
    vat_rate DECIMAL(5,2) NOT NULL DEFAULT 8.1,
        -- Standard-MwSt-Satz (8.1% CH, 19.0% DE, 20.0% AT)
    vat_display_hint VARCHAR(20) NOT NULL DEFAULT 'net',
        -- 'net': "zzgl. 8.1% MwSt."
        -- 'gross': "inkl. 8.1% MwSt."
        -- 'both': "CHF 100.00 netto (CHF 108.10 brutto)"
    
    -- Texte (i18n)
    anonymous_no_price_text JSONB NOT NULL DEFAULT '{"de": "Preis auf Anfrage", "fr": "Prix sur demande", "en": "Price on request"}',
    anonymous_login_cta_text JSONB NOT NULL DEFAULT '{"de": "Einloggen für Preise", "fr": "Connectez-vous pour les prix", "en": "Login for prices"}',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT chk_anonymous_display CHECK (anonymous_price_display IN ('none', 'list', 'from', 'full')),
    CONSTRAINT chk_authenticated_display CHECK (authenticated_price_display IN ('list', 'customer', 'erp_live')),
    CONSTRAINT chk_erp_source CHECK (erp_price_source IN ('none', 'batch_import', 'live_query')),
    CONSTRAINT chk_erp_fallback CHECK (erp_fallback_strategy IN ('hide', 'cached', 'list_price')),
    CONSTRAINT chk_vat_hint CHECK (vat_display_hint IN ('net', 'gross', 'both'))
);
```

### 3.2 Go Domain-Model

```go
type AnonymousPriceDisplay string
const (
    AnonymousPriceNone AnonymousPriceDisplay = "none"
    AnonymousPriceList AnonymousPriceDisplay = "list"
    AnonymousPriceFrom AnonymousPriceDisplay = "from"
    AnonymousPriceFull AnonymousPriceDisplay = "full"
)

type AuthenticatedPriceDisplay string
const (
    AuthenticatedPriceList     AuthenticatedPriceDisplay = "list"
    AuthenticatedPriceCustomer AuthenticatedPriceDisplay = "customer"
    AuthenticatedPriceERPLive  AuthenticatedPriceDisplay = "erp_live"
)

type ERPPriceSource string
const (
    ERPSourceNone        ERPPriceSource = "none"
    ERPSourceBatchImport ERPPriceSource = "batch_import"
    ERPSourceLiveQuery   ERPPriceSource = "live_query"
)

type ERPFallbackStrategy string
const (
    ERPFallbackHide      ERPFallbackStrategy = "hide"
    ERPFallbackCached    ERPFallbackStrategy = "cached"
    ERPFallbackListPrice ERPFallbackStrategy = "list_price"
)

type TenantPricingConfig struct {
    ID                        uuid.UUID                 `json:"id" db:"id"`
    TenantID                  uuid.UUID                 `json:"tenant_id" db:"tenant_id"`
    
    AnonymousPriceDisplay     AnonymousPriceDisplay     `json:"anonymous_price_display" db:"anonymous_price_display"`
    AuthenticatedPriceDisplay AuthenticatedPriceDisplay  `json:"authenticated_price_display" db:"authenticated_price_display"`
    
    ERPPriceSource            ERPPriceSource            `json:"erp_price_source" db:"erp_price_source"`
    ERPFallbackStrategy       ERPFallbackStrategy       `json:"erp_fallback_strategy" db:"erp_fallback_strategy"`
    
    PriceCacheTTLSeconds      int                       `json:"price_cache_ttl_seconds" db:"price_cache_ttl_seconds"`
    
    ShowDiscountPercentage    bool                      `json:"show_discount_percentage" db:"show_discount_percentage"`
    ShowListPriceStrikethrough bool                     `json:"show_list_price_strikethrough" db:"show_list_price_strikethrough"`
    ShowVolumeDiscountTable   bool                      `json:"show_volume_discount_table" db:"show_volume_discount_table"`
    StackVolumeDiscounts      bool                      `json:"stack_volume_discounts" db:"stack_volume_discounts"`
    
    PriceIncludesVAT          bool                      `json:"price_includes_vat" db:"price_includes_vat"`
    VATRate                   float64                   `json:"vat_rate" db:"vat_rate"`
    VATDisplayHint            string                    `json:"vat_display_hint" db:"vat_display_hint"`
    
    AnonymousNoPriceText      map[string]string         `json:"anonymous_no_price_text" db:"anonymous_no_price_text"`
    AnonymousLoginCTAText     map[string]string         `json:"anonymous_login_cta_text" db:"anonymous_login_cta_text"`
}
```

### 3.3 Validierungsregeln

| Regel | Beschreibung |
|-------|-------------|
| `erp_live` erfordert `erp_price_source != 'none'` | ERP-Livepreise brauchen eine konfigurierte ERP-Quelle |
| `erp_fallback_strategy` nur relevant wenn `erp_price_source = 'live_query'` | Fallback nur bei Live-Abfragen |
| `show_discount_percentage` nur sinnvoll bei `authenticated_price_display = 'customer'` | Rabatt-Anzeige braucht Kundenkonditionen |
| `show_list_price_strikethrough` nur bei `customer` oder `erp_live` | Durchgestrichener Preis nur sinnvoll wenn anderer Preis daneben steht |
| `vat_rate > 0` | MwSt-Satz muss positiv sein |

### 3.4 Tenant-Beispiele

**FlexoTech AG (Industriebedarf):**

```json
{
    "anonymous_price_display": "none",
    "authenticated_price_display": "erp_live",
    "erp_price_source": "live_query",
    "erp_fallback_strategy": "cached",
    "price_cache_ttl_seconds": 600,
    "show_discount_percentage": false,
    "show_list_price_strikethrough": false,
    "show_volume_discount_table": false,
    "price_includes_vat": false,
    "vat_rate": 8.1,
    "vat_display_hint": "net"
}
```

**LabSupply GmbH (Laborbedarf):**

```json
{
    "anonymous_price_display": "list",
    "authenticated_price_display": "customer",
    "erp_price_source": "batch_import",
    "erp_fallback_strategy": "list_price",
    "price_cache_ttl_seconds": 300,
    "show_discount_percentage": true,
    "show_list_price_strikethrough": true,
    "show_volume_discount_table": true,
    "price_includes_vat": false,
    "vat_rate": 19.0,
    "vat_display_hint": "net"
}
```

**PackDirect AG (Verpackung):**

```json
{
    "anonymous_price_display": "full",
    "authenticated_price_display": "customer",
    "erp_price_source": "none",
    "erp_fallback_strategy": "list_price",
    "price_cache_ttl_seconds": 300,
    "show_discount_percentage": true,
    "show_list_price_strikethrough": true,
    "show_volume_discount_table": true,
    "stack_volume_discounts": true,
    "price_includes_vat": false,
    "vat_rate": 8.1,
    "vat_display_hint": "net"
}
```

---

## 4. API-Design

### 4.1 Preis-Kontext im Product-Response

Die Product-API liefert **kontextabhängige** Preisinformationen. Der Kontext wird aus dem Auth-Token und der Tenant-Config abgeleitet:

```
GET /api/v1/products/:id
Authorization: Bearer <token>   ← Optional
X-Tenant-ID: <tenant-id>
```

**Response-Preis-Objekt (variiert je nach Kontext):**

**Anonym, Tenant-Config `none`:**

```json
{
    "id": "...",
    "name": {"de": "Hydraulikzylinder CDT3 50mm"},
    "price": {
        "display_mode": "none",
        "message": {"de": "Preis auf Anfrage"},
        "login_cta": {"de": "Einloggen für Preise"},
        "login_url": "/login?redirect=/de/produkte/hydraulikzylinder-cdt3"
    }
}
```

**Anonym, Tenant-Config `from`:**

```json
{
    "price": {
        "display_mode": "from",
        "from_price": {
            "net": 385.00,
            "currency": "CHF"
        },
        "vat_hint": "zzgl. 8.1% MwSt.",
        "login_cta": {"de": "Einloggen für Ihren Preis"}
    }
}
```

**Anonym, Tenant-Config `full`:**

```json
{
    "price": {
        "display_mode": "full",
        "tiers": [
            {"min_quantity": 1, "price_net": 1.20, "currency": "CHF"},
            {"min_quantity": 50, "price_net": 0.95, "currency": "CHF"},
            {"min_quantity": 200, "price_net": 0.88, "currency": "CHF"},
            {"min_quantity": 500, "price_net": 0.85, "currency": "CHF"}
        ],
        "vat_hint": "zzgl. 8.1% MwSt."
    }
}
```

**Eingeloggt, Kundenkonditionen:**

```json
{
    "price": {
        "display_mode": "customer",
        "customer_price": {
            "net": 0.78,
            "currency": "CHF"
        },
        "list_price": {
            "net": 1.20,
            "currency": "CHF",
            "strikethrough": true
        },
        "discount": {
            "percent": 35.0,
            "show": true
        },
        "tiers": [
            {"min_quantity": 1, "price_net": 0.78, "currency": "CHF"},
            {"min_quantity": 50, "price_net": 0.72, "currency": "CHF"},
            {"min_quantity": 200, "price_net": 0.68, "currency": "CHF"},
            {"min_quantity": 500, "price_net": 0.65, "currency": "CHF"}
        ],
        "source": "contract",
        "contract_reference": "RV-2025-0847",
        "vat_hint": "zzgl. 8.1% MwSt."
    }
}
```

**Eingeloggt, ERP-Live:**

```json
{
    "price": {
        "display_mode": "erp_live",
        "customer_price": {
            "net": 412.50,
            "currency": "CHF"
        },
        "source": "erp_live",
        "cached": false,
        "valid_until": "2026-02-14T21:52:00Z",
        "vat_hint": "zzgl. 8.1% MwSt."
    }
}
```

**ERP nicht erreichbar, Fallback `cached`:**

```json
{
    "price": {
        "display_mode": "erp_live",
        "customer_price": {
            "net": 412.50,
            "currency": "CHF"
        },
        "source": "erp_cached",
        "cached": true,
        "cached_at": "2026-02-14T19:30:00Z",
        "notice": {"de": "Preis basiert auf letzter Aktualisierung. Endpreis bei Bestellung."},
        "vat_hint": "zzgl. 8.1% MwSt."
    }
}
```

### 4.2 Dedizierter Preis-Endpoint (ERP-Live, Mengenabfrage)

Für ERP-Live-Preise mit spezifischer Menge:

```
GET /api/v1/products/:id/price?quantity=50
Authorization: Bearer <token>
X-Tenant-ID: <tenant-id>
```

Response:

```json
{
    "product_id": "...",
    "quantity": 50,
    "unit_price": {
        "net": 0.72,
        "gross": 0.78,
        "currency": "CHF"
    },
    "total_price": {
        "net": 36.00,
        "gross": 38.92,
        "currency": "CHF"
    },
    "source": "erp_live",
    "valid_for_seconds": 300
}
```

### 4.3 Bulk-Preis-Endpoint (Warenkorb)

Für Warenkorb-Szenarien mit mehreren Produkten:

```
POST /api/v1/prices/bulk
Authorization: Bearer <token>
X-Tenant-ID: <tenant-id>

{
    "items": [
        {"product_id": "uuid-1", "variant_id": "uuid-v1", "quantity": 50},
        {"product_id": "uuid-2", "quantity": 10},
        {"product_id": "uuid-3", "variant_id": "uuid-v3", "quantity": 200}
    ]
}
```

Response:

```json
{
    "items": [
        {
            "product_id": "uuid-1",
            "variant_id": "uuid-v1",
            "quantity": 50,
            "unit_price_net": 0.72,
            "total_price_net": 36.00,
            "currency": "CHF",
            "source": "erp_live"
        },
        {
            "product_id": "uuid-2",
            "quantity": 10,
            "unit_price_net": 45.00,
            "total_price_net": 450.00,
            "currency": "CHF",
            "source": "customer_price"
        },
        {
            "product_id": "uuid-3",
            "variant_id": "uuid-v3",
            "quantity": 200,
            "unit_price_net": 0.68,
            "total_price_net": 136.00,
            "currency": "CHF",
            "source": "catalog"
        }
    ],
    "subtotal_net": 622.00,
    "currency": "CHF"
}
```

Intern löst dieser Endpoint **einen** Bulk-Call zum ERP aus (nicht N einzelne Calls).

### 4.4 Admin-API für Tenant-Config

```
# Config lesen
GET /api/v1/admin/tenants/:id/pricing-config
→ TenantPricingConfig

# Config aktualisieren
PUT /api/v1/admin/tenants/:id/pricing-config
{
    "anonymous_price_display": "list",
    "authenticated_price_display": "customer",
    ...
}

# Config validieren (Dry Run)
POST /api/v1/admin/tenants/:id/pricing-config/validate
{
    "anonymous_price_display": "list",
    "authenticated_price_display": "erp_live",
    "erp_price_source": "none"
}
→ {"valid": false, "errors": ["erp_live erfordert erp_price_source != 'none'"]}
```

---

## 5. Frontend

### 5.1 Preis-Komponente

Die Frontend-Preis-Komponente rendert sich vollständig basierend auf dem `price`-Objekt aus der API. Keine Business-Logik im Frontend — die API entscheidet, was angezeigt wird.

```typescript
// Vereinfachte Komponentenstruktur
interface PriceDisplayProps {
    price: ProductPrice;
    locale: string;
}

function PriceDisplay({ price, locale }: PriceDisplayProps) {
    switch (price.display_mode) {
        case 'none':
            return <PriceOnRequest message={price.message} loginCta={price.login_cta} />;
        case 'list':
            return <ListPrice price={price.list_price} vatHint={price.vat_hint} />;
        case 'from':
            return <FromPrice price={price.from_price} vatHint={price.vat_hint} loginCta={price.login_cta} />;
        case 'full':
            return <TierPriceTable tiers={price.tiers} vatHint={price.vat_hint} />;
        case 'customer':
            return <CustomerPrice price={price} />;
        case 'erp_live':
            return <ERPPrice price={price} />;
    }
}
```

### 5.2 Darstellungsvarianten

**«Preis auf Anfrage» (mode: none):**
- Text: Konfigurierbar pro Tenant (i18n)
- CTA-Buttons: «Einloggen für Preise» (→ Login) und/oder «Preis anfragen» (→ Kontaktformular)
- Kontaktformular: Produkt-ID und -Name vorausgefüllt

**Durchgestrichener Listenpreis + Kundenpreis:**

```
CHF 0.78 / Stück
CHF 1.20  ← durchgestrichen
-35%      ← optional (show_discount_percentage)
```

**Staffelpreis-Tabelle:**

```
┌──────────┬──────────────┐
│ Menge    │ Preis/Stück  │
├──────────┼──────────────┤
│ 1–49     │ CHF 0.78  ◀  │  ← aktuelle Auswahl markiert
│ 50–199   │ CHF 0.72     │
│ 200–499  │ CHF 0.68     │
│ ab 500   │ CHF 0.65     │
└──────────┴──────────────┘
```

**MwSt-Hinweis:**

| Config | Anzeige |
|--------|---------|
| `vat_display_hint: "net"` | «zzgl. 8.1% MwSt.» |
| `vat_display_hint: "gross"` | «inkl. 8.1% MwSt.» |
| `vat_display_hint: "both"` | «CHF 100.00 netto (CHF 108.10 brutto)» |

**ERP-Cached-Hinweis:**

```
CHF 412.50
ℹ️ Preis basiert auf letzter Aktualisierung. Endpreis bei Bestellung.
```

### 5.3 SEO: Structured Data

Die Preis-Ausgabe in JSON-LD richtet sich nach der Tenant-Config für anonyme Besucher (Google ist ein anonymer Besucher):

| Config | JSON-LD |
|--------|---------|
| `anonymous_price_display: "none"` | Kein `offers`-Block in Product-Schema |
| `anonymous_price_display: "list"` | `"price": "1.20"` (Listenpreis) |
| `anonymous_price_display: "from"` | `"lowPrice": "0.85"` (AggregateOffer) |
| `anonymous_price_display: "full"` | `"lowPrice": "0.85", "highPrice": "1.20"` (AggregateOffer) |

Siehe auch [SEO-URL-Struktur](seo-url-structure.md), Abschnitt 7.2: Preisanzeige: SEO vs. B2B-Praxis.

**Wichtig:** JSON-LD enthält **immer** den öffentlichen Preis (Listenpreis), nie kundenspezifische Konditionen.

---

## 6. Sicherheit

### 6.1 Preis-Isolation

Kundenspezifische Preise dürfen **unter keinen Umständen** an andere Kunden oder anonyme Besucher leaken.

**Massnahmen:**

| Ebene | Massnahme |
|-------|-----------|
| **API** | Pricing Service prüft Auth-Token → Tenant + Customer. Preis-Response wird serverseitig zusammengebaut. Kein Frontend-Filter. |
| **Cache** | Cache-Keys enthalten immer `tenant_id` + `customer_id`. Kein Shared Cache für verschiedene Kunden. |
| **Logs** | Kundenspezifische Preise werden in Logs **nicht** mit Klartext-Kundennamen geloggt. Nur UUIDs. |
| **CDN** | Produktseiten mit ERP-Livepreisen: `Cache-Control: private, no-store`. Nur anonyme Preise dürfen CDN-gecacht werden. |
| **Admin** | Admin sieht alle Preise. Audit-Log bei Preisänderungen. |

### 6.2 Cache-Key-Schema

```
# Anonymer Besucher (CDN-cacheable)
price:anon:{tenant_id}:{product_id}:{variant_id}

# Eingeloggter Kunde (NICHT CDN-cacheable)
price:auth:{tenant_id}:{customer_id}:{product_id}:{variant_id}:{quantity}
```

**Regel:** Anonyme Preise → Redis + CDN. Authentifizierte Preise → nur Redis, `Cache-Control: private`.

### 6.3 Rate Limiting

ERP-Live-Preisabfragen sind teuer (SAP-Last). Rate Limiting pro Kunde:

| Endpoint | Limit |
|----------|-------|
| `GET /products/:id/price` | 60 req/min pro Kunde |
| `POST /prices/bulk` | 10 req/min pro Kunde, max. 100 Items |

---

## 7. Pricing Service — Interner Ablauf

### 7.1 Preis-Auflösung (Entscheidungsbaum)

```
Anfrage eingehend
│
├── Auth-Token vorhanden?
│   │
│   ├── NEIN (anonym)
│   │   │
│   │   └── Tenant-Config: anonymous_price_display
│   │       ├── 'none'  → {display_mode: "none", message: ...}
│   │       ├── 'list'  → Katalogpreis laden (Menge 1)
│   │       ├── 'from'  → Niedrigsten Staffelpreis laden
│   │       └── 'full'  → Alle Staffelpreise laden
│   │
│   └── JA (authentifiziert)
│       │
│       └── Tenant-Config: authenticated_price_display
│           ├── 'list'     → Katalogpreis laden
│           ├── 'customer' → Kundenpreis-Kaskade (siehe 1.2)
│           └── 'erp_live' → ERP-Preis anfragen
│                             ├── Erfolg → Preis zurückgeben
│                             └── Fehler → Fallback-Strategie
│                                 ├── 'hide'       → {display_mode: "erp_live", notice: "..."}
│                                 ├── 'cached'     → Letzten Cache-Wert
│                                 └── 'list_price' → Katalogpreis
│
└── Preis-Objekt zusammenbauen
    ├── Anzeige-Optionen anwenden (Strikethrough, Discount, VAT)
    └── Response zurückgeben
```

### 7.2 Service-Interface

```go
type PricingService interface {
    // Preis für ein Produkt im aktuellen Kontext
    GetProductPrice(ctx context.Context, req ProductPriceRequest) (*ProductPriceResponse, error)
    
    // Bulk-Preise (Warenkorb)
    GetBulkPrices(ctx context.Context, req BulkPriceRequest) (*BulkPriceResponse, error)
    
    // Preis für spezifische Menge (ERP-Live)
    GetQuantityPrice(ctx context.Context, req QuantityPriceRequest) (*QuantityPriceResponse, error)
    
    // Tenant-Config
    GetPricingConfig(ctx context.Context, tenantID uuid.UUID) (*TenantPricingConfig, error)
    UpdatePricingConfig(ctx context.Context, config *TenantPricingConfig) error
}

type ProductPriceRequest struct {
    TenantID   uuid.UUID  `json:"tenant_id"`
    ProductID  uuid.UUID  `json:"product_id"`
    VariantID  *uuid.UUID `json:"variant_id,omitempty"`
    CustomerID *uuid.UUID `json:"customer_id,omitempty"` // nil = anonym
    Quantity   int        `json:"quantity,omitempty"`     // Default: 1
}
```

---

## 8. Admin-UI

### 8.1 Tenant-Pricing-Konfiguration

Das Admin-UI bietet ein dediziertes Panel für die Preisanzeige-Konfiguration pro Tenant:

**Sektion 1: Anonyme Besucher**
- Radio-Buttons: Keine Preise / Listenpreise / Ab-Preise / Volle Preise
- Vorschau: Live-Preview wie ein Produkt für anonyme Besucher aussieht
- Texte editierbar: «Preis auf Anfrage»-Text, Login-CTA-Text (i18n)

**Sektion 2: Eingeloggte Kunden**
- Radio-Buttons: Listenpreise / Kundenkonditionen / ERP-Livepreise
- Bei ERP-Live: Zusatzoptionen (Fallback, Cache-TTL)
- Vorschau: Live-Preview mit Beispiel-Kundenpreis

**Sektion 3: Anzeige-Optionen**
- Toggles: Rabatt-Prozent anzeigen, Listenpreis durchgestrichen, Staffeltabelle
- MwSt-Einstellungen: Netto/Brutto, Satz, Hinweis-Modus

**Sektion 4: ERP-Einstellungen** (nur wenn ERP-Quelle != none)
- Quelle: Batch-Import / Live-Query
- Fallback-Strategie
- Cache-TTL (Slider: 60s – 3600s)
- Test-Button: «ERP-Verbindung testen» → Probe-Preisabfrage

---

## 9. Migration / Erweiterung Produkt-Model

Das bestehende `Product`-Struct (siehe `product.go`) enthält aktuell **keine Preisinformationen**. Preise werden bewusst in separaten Tabellen geführt (siehe Abschnitt 2.2, 2.3). Das Product-Model wird **nicht** um Preisfelder erweitert — stattdessen liefert der Pricing Service die Preise kontextabhängig.

**Grund:** Ein Produkt hat keinen «einen Preis». Der Preis hängt ab von:
- Tenant-Config
- Auth-Kontext (anonym vs. Kunde)
- Kundengruppe / Vertrag
- Menge
- Zeitpunkt (Gültigkeit, ERP-Kurs)

Das Product-Domain-Model bleibt schlank. Preise werden **on-demand** vom Pricing Service ermittelt.

---

## 10. Performance-Überlegungen

| Szenario | Latenz-Ziel | Strategie |
|----------|-------------|-----------|
| Anonymer Besucher, Katalogpreise | < 50ms | Redis-Cache, CDN |
| Eingeloggt, Kundenpreise (DB) | < 100ms | Redis-Cache pro Customer, DB-Index |
| Eingeloggt, ERP-Live (Cache Hit) | < 50ms | Redis-Cache mit TTL |
| Eingeloggt, ERP-Live (Cache Miss) | < 2s | SAP-Call, Timeout 3s, Fallback |
| Warenkorb Bulk (50 Positionen) | < 3s | Ein SAP-Bulk-Call, paralleles Cache-Lookup |
| Produktliste (20 Produkte, anonym) | < 100ms | Batch-Laden aus Redis/DB |
| Produktliste (20 Produkte, eingeloggt) | < 500ms | Batch-Laden, ggf. ERP-Bulk |

**Circuit Breaker für ERP:**

```go
// ERP-Aufrufe mit Circuit Breaker absichern
// Nach 5 Fehlern in 30s → Circuit Open → Fallback für 60s
erpBreaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "erp-pricing",
    MaxRequests: 3,              // Requests im Half-Open State
    Interval:    30 * time.Second,
    Timeout:     60 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures >= 5
    },
})
```

---

## Priorisierung

| Phase | Aufgabe | Aufwand | Abhängigkeiten |
|-------|---------|---------|----------------|
| **Phase 1** | Tenant-Pricing-Config Tabelle + Admin-API + Validierung | 2–3 Tage | Identity Service (Tenant-Model) |
| **Phase 2** | Katalogpreise (DB-Tabelle, CRUD-API) | 2–3 Tage | Catalog Service |
| **Phase 3** | Pricing Service: anonyme Preisauflösung (none/list/from/full) | 3–4 Tage | Phase 1, 2 |
| **Phase 4** | Kundenspezifische Preise (DB, Kaskade, Import) | 3–4 Tage | Phase 3 |
| **Phase 5** | Frontend: Preis-Komponente (alle Display-Modi) | 3–4 Tage | Phase 3 |
| **Phase 6** | ERP-Live-Integration (SAP Bulk-Call, Caching, Circuit Breaker) | 5–7 Tage | Phase 3, SAP-Integration |
| **Phase 7** | Admin-UI: Tenant-Pricing-Config Panel | 2–3 Tage | Phase 1, 5 |
| **Phase 8** | SEO: JSON-LD Preise basierend auf Config | 1–2 Tage | Phase 5 |

**Gesamtaufwand:** ~4–6 Wochen

---

## Offene Fragen

1. **Währungen:** Unterstützt Gondolia mehrere Währungen pro Tenant? (CHF + EUR für CH/DE-übergreifende Tenants?) → Auswirkung auf Cache-Keys und Preistabellen.
2. **Preishistorie:** Sollen historische Preise gespeichert werden? (Für Nachvollziehbarkeit, «Preis am Bestelldatum» bei Reklamationen)
3. **Rundung:** Schweizer Rappen-Rundung (auf 0.05 CHF runden) vs. Cent-genaue Preise in EUR. Pro Tenant konfigurierbar?
4. **Mindestbestellwert:** Gehört zur Pricing-Config oder zum Order-Service?
5. **Preisanzeige in Suchergebnissen (Meilisearch):** Sollen Preise im Suchindex sein? (Schnellere Anzeige, aber komplexere Indexierung pro Customer)
6. **Mengeneinheiten:** Preis pro Stück, pro Meter, pro kg — wie wird die Mengeneinheit in der Preisanzeige berücksichtigt?
7. **SAP-Konditionsarten:** Welche SAP-Konditionsarten (KONV-Tabelle) müssen unterstützt werden? (PR00, ZK01, ZK02, ...)
