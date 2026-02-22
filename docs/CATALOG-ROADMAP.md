# Catalog Service — Projekt-Roadmap

Stand: 2026-02-18
Status-Legende: ✅ Done | 🔶 Teilweise | ❌ Offen | ⏳ Blocked

---

## Übersicht

Der Catalog Service braucht 5 Produkttypen (siehe `architecture/product-types.md`).
Wir arbeiten sie **sequentiell** ab — ein Typ wird fertig, bevor der nächste startet.
Zusätzlich gibt es übergreifende Verbesserungen (TODO-catalog-improvements.md).

**Reihenfolge:**
1. Varianten — Kernfunktion ✅ (PIM-Sync, Meilisearch, Migration-Tooling → Phase B)
2. Übergreifende Verbesserungen (Catalog Basics)
3. Parametrisierbare Produkte
4. Bundle-Produkte
5. Komplexe Produkte (abhängige Parameter)
6. Konfigurator-Produkte
7. Product Intelligence (Auto-Erkennung)
8. **Phase B: Infrastruktur** (PIM-Sync, Meilisearch-Index, Migration-Tooling — für alle Produkttypen gemeinsam)

---

## Meilenstein 1: Varianten — FERTIGSTELLEN

Referenz: `architecture/product-variants.md`

### Phase 1: DB + Domain + CRUD ✅
- [x] Migration 000001: `product_type`, `parent_id`, `variant_axes`, `variant_axis_values`
- [x] Migration 000002: Seed-Daten (IM3000 Motor 3 Achsen, SafeGrip Handschuhe 1 Achse, 7 Varianten)
- [x] Domain: `ProductType`, `VariantAxis`, `AxisOption`, `AxisValueEntry`, `ProductVariant`
- [x] Domain: Request-Structs (`CreateVariantParentRequest`, `CreateVariantRequest`)
- [x] Repository Interface: alle Varianten-Methoden (`ListVariants`, `FindVariantByAxisValues`, `GetAvailableAxisValues`, `Set/GetVariantAxes`, `Set/GetAxisValues`)
- [x] Postgres Repo: `variant_repo.go` (420 LOC)
- [x] Service: `variant_service.go` (Axis-Validation, Unique-Check, Vererbung)
- [x] Handler: CRUD Endpoints registriert
- [x] Routes: `GET/POST /:id/variants`, `GET /:id/variants/select`, `GET /:id/variants/available`

### Phase 2: Storefront-API 🔶
- [x] `GET /products/:id` → gibt `variant_axes` + `variants` zurück bei variant_parent
- [x] `GET /products/:id/variants/select?axis=value` → findet Variante
- [x] `GET /products/:id/variants/available?axis=value` → verfügbare Optionen
- [x] **Produktliste erweitern:** `price_range` (min/max) auf variant_parent berechnen
- [x] **Produktliste erweitern:** `variant_count` Feld
- [x] **Produktliste erweitern:** `variant_summary` (Achsenwerte-Zusammenfassung)
- [x] **Produktliste:** Varianten (product_type=variant) aus Listenansicht ausblenden (`ExcludeVariants: true`)
- [x] **Staffelpreise** in `/variants/select` Response (`tier_prices`, `VariantPrice` mit `TierPrice[]`)
- [x] **Variante direkt abrufbar:** `GET /products/:variant_id` → `parent_summary` + Bild-Vererbung + Name-Generierung

### Phase 3: PIM-Sync ⏭️ (verschoben → Phase B, nach allen Produkttypen)
- [ ] `PIMProvider` Interface erweitern: `FetchProductModels()`, `FetchVariantAxes()`, `SupportsVariants()`
- [ ] `pim.ProductModelPage`, `pim.ProductModel`, `pim.VariantAxisDefinition` Structs
- [ ] `SyncService.SyncProductModels()` — Provider-agnostischer Import
- [ ] PIM-Provider: Referenz-Implementierung
- [ ] Level-Flachdrückung (2 Level → 1 Level, PIM-abhängig)
- [ ] CSV/Excel-Provider: Minimal-Implementierung

### Phase 4: Frontend ✅ (Matrix-Ansicht separat)
- [x] `VariantSelector.tsx` (317 LOC) — Buttons/Dropdowns, abhängige Filterung
- [x] Preis/SKU/Bild-Switching bei Variantenwahl
- [x] i18n Label-Formatierung (formatOptionLabel, formatAxisLabel)
- [x] **Dropdown-Fallback** bei >6 Optionen pro Achse (korrigiert von >8)
- [x] **URL-Update:** `?variant=SKU` für Bookmarks/Teilen (war schon implementiert)
- [x] **Matrix-Ansicht** (B2B-Schnellbestellung): Tabelle aller Varianten mit Mengen-Input ← Phase 4b
- [x] **Preisbereich** in Produktliste: «ab CHF 485.00» (war schon in ProductCard)
- [x] **Varianten-Zusammenfassung** in Produktkarte: Tags mit Achsenwerten

### Phase 5: Suchindex (Meilisearch) ⏭️ (verschoben → Phase B, nach allen Produkttypen)
- [ ] Parent im Index mit allen Varianten-SKUs (für SKU-Suche)
- [ ] Achsenwerte als filterbare Facetten
- [ ] Preisbereich (min/max) im Index
- [ ] Aggregierte Verfügbarkeit
- [ ] Facetten-Counts in Suchresultaten

### Phase 6: Migration-Tooling ⏭️ (verschoben → Phase B, nach allen Produkttypen)
- [ ] `MigrationService.MigrateToVariants()` — Einzelprodukte → Variantenprodukt
- [ ] Rückwärtskompatibilität (URLs, Warenkorb, Preise)
- [ ] Audit-Log für Konvertierungen

### Design-Entscheide ✅ (dokumentiert in product-variants.md, Abschnitt 11)
- [x] **Max. 4 Achsen** — Validierung in CreateVariantParent() einbauen
- [x] **Soft-Limit 200 Varianten** — Warnung im Admin-UI, kein Hard-Limit
- [x] **Variante direkt abrufbar** — mit parent_id + parent_summary, Frontend redirectet auf Parent-PDP
- [x] **Bild-Override ersetzt** — variant.Images > 0 → nur Varianten-Bilder, sonst Parent
- [x] **Name auto-generiert** — «Parent — Achse1, Achse2» mit manuellem Override
- [x] **Fertige Labels** — Achsen-Optionen enthalten Einheit im Label, kein separates System
- [x] **Cross-Sell erst mit Warenkorb** — PDP-Achsenselektor ist bereits Cross-Sell

### Noch umzusetzen aus den Entscheiden
- [x] Validierung: max. 4 Achsen in VariantService (`ErrTooManyAxes`)
- [ ] Validierung: Soft-Limit 200 Warnung (Admin-UI, niedrige Prio)
- [ ] `GET /products/:variant_id` → parent_summary Feld + Redirect-Info
- [x] Bild-Vererbungslogik: `GetEffectiveImages()` in domain + Repo
- [x] Name-Generierung: `GenerateVariantName()` in domain + SelectVariant

---

## Meilenstein 2: Catalog Basics — Übergreifende Verbesserungen

Referenz: `architecture/TODO-catalog-improvements.md`

### Backend
- [x] `product_count` in Category-Response (eliminiert N+1) — Bugfix: LATERAL join für korrekte per-category counts
- [x] `include_children=true` Parameter für Kategorie-Produkte (rekursive CTE) — war schon implementiert
- [ ] Auto-Migration Setup (golang-migrate + `schema_migrations` Tabelle)

### Frontend
- [ ] Serverseitige Pagination für Kategorie-Produkte (nach include_children)
- [ ] Breadcrumbs mit vollständiger Kategorie-Hierarchie
- [ ] Attribut-Namen i18n auf Produktdetailseite (✅ Backend: attribute_translations existiert)
- [ ] Lazy Loading / Infinite Scroll (niedrige Prio)

---

## Meilenstein 3: Parametrisierbare Produkte

Referenz: `architecture/product-types.md`, Abschnitt 2

### DB + Domain
- [ ] `ProductType` erweitern: `parametric`
- [ ] CHECK Constraint erweitern
- [ ] Tabelle `product_parameters` (code, data_type, unit, min/max/step, options)
- [ ] Tabelle `product_pricing_formulas` (Formel + Variablen)
- [ ] Domain-Structs: `ProductParameter`, `ParameterOption`, `PricingFormula`

### Backend
- [ ] Repository: Parameter-CRUD
- [ ] Service: Preis-Berechnung (Formel-Evaluator, z.B. `expr` Library)
- [ ] Handler: `POST /products/:id/calculate-price`
- [ ] Validierung: Min/Max/Step Server-seitig
- [ ] Seed-Daten: 2–3 Beispielprodukte (Kabel auf Mass, Plattenzuschnitt)

### Frontend
- [ ] Parametrische Eingabefelder (Slider/Number mit Einheit)
- [ ] Live-Preisberechnung (debounced)
- [ ] Preisaufschlüsselung anzeigen

### Warenkorb-Integration
- [ ] `CartItem.Parameters` Feld
- [ ] Parameter-Validierung bei Warenkorb-Hinzufügung
- [ ] Preis-Neuberechnung bei Checkout

---

## Meilenstein 4: Bundle-Produkte

Referenz: `architecture/product-types.md`, Abschnitt 4

### DB + Domain
- [ ] `ProductType` erweitern: `bundle`
- [ ] Tabelle `bundle_items` (item_product_id, quantity, fixed/variable, optional)
- [ ] Tabelle `bundle_pricing` (fixed, sum_discount_percent, sum_discount_absolute)
- [ ] Domain-Structs: `BundleItem`, `BundlePricing`

### Backend
- [ ] Repository: Bundle-CRUD
- [ ] Service: Bundle-Preisberechnung
- [ ] Handler: `POST /products/:id/calculate-bundle-price`
- [ ] Verfügbarkeitsprüfung (alle Pflichtkomponenten auf Lager?)
- [ ] Seed-Daten: 2 Beispiel-Bundles

### Frontend
- [ ] Bundle-Übersicht (Komponenten-Liste mit Bildern/Preisen)
- [ ] Variable Mengen-Inputs
- [ ] Optionale Komponenten (Checkbox)
- [ ] Ersparnis-Anzeige vs. Einzelkauf

### Warenkorb-Integration
- [ ] `CartItem.BundleItems` Feld
- [ ] Bundle-Auflösung (N Positionen) für ERP-Systeme

---

## Meilenstein 5: Komplexe Produkte (Abhängige Parameter)

Referenz: `architecture/product-types.md`, Abschnitt 3

### DB + Domain
- [ ] `ProductType` erweitern: `complex`
- [ ] Tabelle `product_parameter_rules` (condition JSON, effect JSON, priority)
- [ ] Domain-Structs: `RuleCondition`, `RuleEffect`, `ParameterRule`

### Backend
- [ ] Regel-Engine (JSON-basiert, AND/OR, Operatoren: eq/gt/in/...)
- [ ] Handler: `POST /products/:id/evaluate-rules`
- [ ] Handler: `POST /products/:id/validate-and-price`
- [ ] Seed-Daten: Schaltschrank-Konfigurator oder Hydraulikschlauch

### Frontend
- [ ] Dynamische Formulare (Felder ein-/ausblenden, Optionen einschränken)
- [ ] Erklärungen bei deaktivierten Optionen (Tooltip)
- [ ] Sinnvolle Reihenfolge (entscheidende Parameter zuerst)

### Entscheid
- [ ] Regel-Engine nur Backend (API-Call) oder auch Frontend (JSON-Regeln clientseitig)?

---

## Meilenstein 6: Konfigurator-Produkte

Referenz: `architecture/product-types.md`, Abschnitt 5

### DB + Domain
- [ ] `ProductType` erweitern: `configurator`
- [ ] Tabelle `configurator_steps` (code, step_type, position, required)
- [ ] Tabelle `configurator_options` (product_id ref, price_modifier)
- [ ] Tabelle `configurator_dependencies` (source_step → target_step, allowed_options)
- [ ] Domain-Structs: `ConfiguratorStep`, `ConfiguratorOption`, `ConfiguratorDependency`, `PriceModifier`

### Backend
- [ ] Repository: Step/Option/Dependency CRUD
- [ ] Service: Dependency-Auswertung, Preisberechnung
- [ ] Handler: `POST /products/:id/configure`
- [ ] Seed-Daten: Grossküchen-Kochblock oder 19"-Rack

### Frontend
- [ ] Wizard/Stepper UI (Step-by-Step)
- [ ] Abhängige Optionen (dynamisch filtern)
- [ ] Zusammenfassung vor Warenkorb
- [ ] Konfiguration speichern/teilen (B2B)

### Warenkorb-Integration
- [ ] `CartItem.Configuration` Feld
- [ ] Konfiguration → ERP-Stückliste (BOM) oder Einzelpositionen

---

## Meilenstein 7: Product Intelligence

Referenz: `architecture/product-types.md`, Abschnitt 6

### Phase 1: Regelbasiert
- [ ] Ähnlichkeitserkennung: SKU-Pattern, Name-Prefix, Kategorie-Match
- [ ] Clustering: Produktgruppen mit Konfidenz-Score
- [ ] Achsen-Erkennung aus Unterschieden
- [ ] Review-UI: Vorschlag bestätigen/anpassen/verwerfen

### Phase 2: ML/AI (optional)
- [ ] Embeddings für semantische Ähnlichkeit
- [ ] pgvector oder dedizierter Vektor-Index
- [ ] Feedback-Loop (bestätigte Vorschläge als Trainingsdaten)

### Weitere Assistenz
- [ ] Duplikat-Erkennung
- [ ] Attribut-Normalisierung
- [ ] Kategorie-Vorschläge
- [ ] Preisanomalien-Erkennung
- [ ] Lücken in Variantenmatrix erkennen

---

## Übergreifend: Warenkorb

Wird benötigt ab Meilenstein 3 (parametrisierbar), aber Design jetzt festlegen:

- [ ] Generisches `CartItem`-Struct (product_type-spezifische Felder)
- [ ] Cart-Service als eigener Microservice oder Teil von Catalog?
- [ ] ERP-Mapping pro Produkttyp klären

---

## Übergreifend: ERP-Integration

Offene Fragen:
- [ ] ERP-Variantenkonfiguration Support?
- [ ] Materialnummern pro Variante oder nur Master?
- [ ] Parametrische Produkte im ERP (Merkmalsbewertung)?
- [ ] Bundles als Stückliste (BOM) oder Einzelpositionen?

---

## Aktueller Code-Stand (Referenz)

| Datei | LOC | Status |
|-------|-----|--------|
| `domain/variant.go` | 100 | ✅ |
| `domain/product.go` | ~80 | ✅ (nur simple/variant_parent/variant) |
| `repository/interfaces.go` | 74 | ✅ |
| `repository/postgres/variant_repo.go` | 420 | ✅ |
| `repository/postgres/product_repo.go` | 373 | ✅ |
| `service/variant_service.go` | 244 | ✅ |
| `service/product_service.go` | 220 | ✅ |
| `service/price_service.go` | 151 | ✅ |
| `service/search_service.go` | 70 | Stub |
| `service/sync_service.go` | 240 | Stub (PIM-Sync) |
| `handler/variant.go` | 259 | ✅ |
| `frontend/VariantSelector.tsx` | 317 | ✅ |
| Migrations | 6 | 000001–000006 |
