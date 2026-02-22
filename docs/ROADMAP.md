# Gondolia — Roadmap

Stand: 2026-02-22

---

## Phase 1: Foundation ✅

### Identity Service ✅
- [x] User Management (CRUD, Rollen, Permissions)
- [x] Company Management (Multi-Tenancy)
- [x] JWT Auth (Access + Refresh Token, HttpOnly Cookie)
- [x] Login/Logout Frontend

### Catalog Service — Basis ✅
- [x] Produkte (i18n Name, Description, Attributes, Images)
- [x] Kategorien (hierarchisch, parent_id, product_count mit CTE)
- [x] Kategorien: Bilder, Beschreibung, SEO
- [x] Preise (B2B-Staffelpreise, Währung)
- [x] Multi-Tenancy (Row-Level, X-Tenant-ID)
- [x] Suche (DB-basiert, Volltext)
- [x] Attribut-Übersetzungen (i18n Labels + Einheiten)

### Frontend — Basis ✅
- [x] Next.js 14 mit Tailwind CSS
- [x] Login-Flow
- [x] Produktliste mit Infinite Scroll
- [x] Produkttyp-Filter (Chips: Alle/Einfach/Varianten/Parametrisch/Bundles)
- [x] Kategorie-Sidebar mit Hierarchie
- [x] Kategorie-Produkte mit include_children
- [x] Breadcrumbs mit vollständiger Hierarchie
- [x] Produktdetailseite mit übersetzten Attributen

### Infrastruktur ✅
- [x] Docker Compose mit Traefik (Production-naher Dev-Setup)
- [x] PostgreSQL (Database per Service)
- [x] Migrations-System
- [x] 4-Layer-Architektur (Handler → Service → Repository → Domain)

---

## Phase 2: Produkttypen ✅ (4/6)

### simple ✅
- [x] Flat Product mit SKU, Preis, Attributen

### variant_parent + variant ✅
- [x] Parent mit Varianten-Achsen (Farbe, Grösse, etc.)
- [x] Kind-Produkte mit eigener SKU/Preis/Bestand
- [x] Variant Selector (Pill-Buttons, Verfügbarkeitslogik)
- [x] Matrix View (Order-Mode, Compact Cards)
- [x] Tier Pricing pro Variante (Lazy-Loading)
- [x] 74 variant_parents + 199 variants (Seed-Daten)

### parametric ✅
- [x] Select-Achsen (Dropdown/Buttons) → SKU-Auflösung
- [x] Range-Achsen (numerische Eingabe) → Mengenberechnung
- [x] Formeln: per_m2, per_running_meter, per_unit, fixed
- [x] SKU-Mapping (Kombination → konkrete Artikelnummer)
- [x] Live-Preisberechnung (Debounced API)
- [x] 3 Seed-Produkte (Stahlblech, Hydraulikschlauch, Alu-Platte)

### bundle ✅
- [x] Fixed Bundles (feste Zusammenstellung, berechneter Preis)
- [x] Configurable Bundles (wählbare Mengen, min/max)
- [x] Parametrische Komponenten in Bundles (Select + Range Achsen)
- [x] Bundle-Preisberechnung (computed + fixed mode)
- [x] Embedded ParametricConfigurator (ohne doppelten Warenkorb)
- [x] 3 Seed-Bundles mit SVG-Bildern

### complex ❌ (Phase 5)
- [ ] Abhängige Parameter mit Constraint-System
- [ ] Regeln: "Wenn Material=Edelstahl, dann max. Dicke=50mm"
- [ ] Constraint Engine (siehe Konfigurator-Showcase)
- Referenz: `docs/architecture/product-types.md` §3

### configurator ❌ (Phase 5)
- [ ] Step-basierte Konfiguration mit Abhängigkeiten
- [ ] Komponenten-Gruppen mit kaskadierende Einschränkungen
- [ ] Konfigurator-Showcase als Proof of Concept vorhanden
- Referenz: `docs/architecture/product-types.md` §5

---

## Phase 3: Commerce 🔜 (nächster Schritt)

### Cart Service
- [ ] Warenkorb (CRUD, Session-basiert + User-basiert)
- [ ] Produkte hinzufügen (alle 4 Typen: simple, variant, parametric, bundle)
- [ ] Mengen ändern, Positionen entfernen
- [ ] Warenkorb-Zusammenfassung (Zwischensumme, MwSt)
- [ ] Warenkorb persistent über Sessions
- [ ] Mini-Cart im Header

### Order Service
- [ ] Bestellung aus Warenkorb erstellen
- [ ] Bestellstatus-Workflow (draft → confirmed → processing → shipped → delivered)
- [ ] Bestellübersicht (Liste + Detail)
- [ ] Bestellhistorie pro Kunde
- [ ] PDF-Generierung (Auftragsbestätigung)

### Checkout
- [ ] Adresseingabe (Liefer- + Rechnungsadresse)
- [ ] Versandart-Auswahl
- [ ] Bestellzusammenfassung + Bestätigung
- [ ] Gastbestellung vs. eingeloggt

---

## Phase 4: B2B Essentials

### Inventory Service
- [ ] Lagerbestand pro SKU/Werk
- [ ] Verfügbarkeitsanzeige im Frontend
- [ ] Bestandsreservierung bei Bestellung
- [ ] Multi-Lager (Plant/Zone)
- Referenz: `docs/architecture/overview.md`

### Shipping Service
- [ ] Adress-Validierung
- [ ] PLZ-Zonen + Versandkostenberechnung
- [ ] Carrier-Integration (vorbereitet via Provider Pattern)
- Referenz: `docs/architecture/adr-001-provider-pattern.md`

### Pricing Erweitert
- [ ] Kundenspezifische Preise (pro Company/Gruppe)
- [ ] Rabattregeln
- [ ] Preisanzeige-Konfiguration (Netto/Brutto, MwSt)
- Referenz: `docs/architecture/pricing-display-config.md`

### Search
- [ ] Meilisearch Integration (empfohlen in Evaluation)
- [ ] Facetten-Filter (Preis, Attribute, Hersteller)
- [ ] Suggest/Autocomplete
- Referenz: `docs/architecture/search-engine-evaluation.md`

---

## Phase 5: Advanced Features

### Komplexe Produkttypen
- [ ] Complex Products (abhängige Parameter + Constraint Engine)
- [ ] Configurator Products (Step-basiert, Komponenten-Auswahl)
- Referenz: `docs/architecture/product-types.md` §3, §5

### Asset Management
- [ ] Zentrale Medienverwaltung
- [ ] Bild-Optimierung (Resize, WebP)
- [ ] Storage Provider (Azure Blob, S3, MinIO)
- Referenz: `docs/architecture/asset-management.md`

### PIM Integration
- [ ] Akeneo/Pimcore Sync
- [ ] Automatischer Produkt-Import
- [ ] Attribut-Mapping
- Referenz: `docs/architecture/pim-attribute-system.md`

---

## Phase 6: Enterprise & Integration

### SAP Integration
- [ ] Kunden-Sync (BAPI)
- [ ] Bestell-Sync
- [ ] Preis/Bestand-Abgleich
- [ ] gondolia-kuratle Private Repo
- Referenz: `docs/architecture/adr-001-provider-pattern.md`

### Weitere Services
- [ ] Notification Service (E-Mail, Push)
- [ ] CMS (Headless, Payload)
- [ ] Analytics (Conversion Funnel, PowerBI)
- [ ] Support Portal (Customer Journey Tracking)
- [ ] KI-Assistent (RAG, Produktberatung)
- [ ] Datahub (visuelle Integration-Konfiguration)

### Multi-Shop
- [ ] Theming-System
- [ ] Multi-Instance Orchestration
- [ ] SEO-freundliche URLs
- Referenz: `docs/architecture/theming-multi-shop.md`, `seo-url-structure.md`

---

## Architektur-Entscheidungen

| Entscheidung | Status | Dokument |
|---|---|---|
| Provider Pattern | ✅ Akzeptiert | `adr-001-provider-pattern.md` |
| 4-Layer Service-Architektur | ✅ Umgesetzt | `service-layer-structure.md` |
| Database per Service | ✅ Umgesetzt | `overview.md` |
| Meilisearch für Search | 📋 Geplant | `search-engine-evaluation.md` |
| Event-Driven (Kafka) | 📋 Geplant | `overview.md` |
| Headless CMS (Payload) | 📋 Geplant | `cms-concept.md` |

---

## Service-Verzeichnisse

| Service | Code | Status |
|---|---|---|
| identity | `services/identity/` | ✅ Produktiv |
| catalog | `services/catalog/` | ✅ Produktiv |
| cart | `services/cart/` | 📁 Platzhalter |
| order | `services/order/` | 📁 Platzhalter |
| inventory | `services/inventory/` | 📁 Platzhalter |
| shipping | `services/shipping/` | 📁 Platzhalter |
| notification | `services/notification/` | 📁 Platzhalter |
| gateway | `services/gateway/` | 📁 Platzhalter |
