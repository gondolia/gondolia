# Catalog Frontend Implementation - Gondolia

## Implementierte Features

### 1. **Types** (`src/types/catalog.ts`)
- ✅ Product, Category, PriceScale, Manufacturer Types
- ✅ API Response Types (snake_case mapping)
- ✅ PaginatedResponse Type
- ✅ ProductSearchParams Type
- ✅ Mapper Functions (API → Frontend)

### 2. **API Client Erweiterung** (`src/lib/api/client.ts`)
Neue Methoden:
- `getProducts(params)` - Produktliste mit Pagination & Filter
- `getProduct(id)` - Produktdetail
- `searchProducts(query)` - Volltextsuche
- `getCategories()` - Alle Kategorien (hierarchisch)
- `getCategory(id)` - Kategorie-Detail
- `getCategoryProducts(categoryId, params)` - Produkte einer Kategorie
- `getProductPrices(productId)` - Staffelpreise
- `getManufacturers()` - Hersteller (optional)

### 3. **UI Komponenten** (`src/components/catalog/`)
- ✅ `ProductCard.tsx` - Produktkarte mit Bild, Preis, Lagerbestand
- ✅ `CategorySidebar.tsx` - Hierarchische Kategorie-Navigation
- ✅ `Pagination.tsx` - Pagination mit Seitennummern

### 4. **Pages**

#### Produktliste (`app/(auth)/products/page.tsx`)
- ✅ Grid-Layout mit Produktkarten
- ✅ Kategorie-Filter Sidebar (hierarchisch)
- ✅ Suchfunktion
- ✅ Pagination
- ✅ Loading States
- ✅ Error Handling
- ✅ Responsive Design

#### Produktdetail (`app/(auth)/products/[id]/page.tsx`)
- ✅ Produktbild (Placeholder wenn kein Bild)
- ✅ Produktinformationen (Name, SKU, Beschreibung)
- ✅ Preis & Lagerbestand
- ✅ Staffelpreise-Tabelle mit Rabatten
- ✅ Mengenauswahl
- ✅ Gesamtpreis-Berechnung
- ✅ Breadcrumb Navigation
- ✅ Technische Daten (Attributes)
- ✅ Responsive Design

#### Kategorieübersicht (`app/(auth)/categories/page.tsx`)
- ✅ Grid mit Hauptkategorien
- ✅ Produktanzahl pro Kategorie
- ✅ Anzahl Unterkategorien
- ✅ Kategoriebilder (optional)

#### Kategoriedetail (`app/(auth)/categories/[id]/page.tsx`)
- ✅ Kategoriebeschreibung
- ✅ Unterkategorien anzeigen
- ✅ Produktliste der Kategorie
- ✅ Pagination
- ✅ Breadcrumb

### 5. **Navigation** (`src/components/layout/Header.tsx`)
- ✅ Links zu Produkten & Kategorien
- ✅ Active State Highlighting
- ✅ Responsive Navigation

## Stil & Features

- ✅ **Tailwind CSS** - Alle Komponenten nutzen Tailwind
- ✅ **Dark Mode** - Vollständige Dark Mode Unterstützung
- ✅ **Responsive** - Mobile-first Design
- ✅ **Deutsche Texte** - UI-Texte auf Deutsch
- ✅ **Loading States** - Spinner für asynchrone Aktionen
- ✅ **Error Handling** - Fehler werden angezeigt
- ✅ **Konsistenter Stil** - Gleicher Stil wie Auth-Pages

## API Endpoints (Backend)

Die Frontend-Pages erwarten folgende Backend-Endpoints:

### Produkte
```
GET  /api/v1/catalog/products              # Liste mit Pagination
GET  /api/v1/catalog/products/:id          # Detail
GET  /api/v1/catalog/products/search?q=    # Suche
```

### Kategorien
```
GET  /api/v1/catalog/categories            # Alle (hierarchisch)
GET  /api/v1/catalog/categories/:id        # Detail
GET  /api/v1/catalog/categories/:id/products  # Produkte
```

### Preise
```
GET  /api/v1/catalog/prices/product/:id    # Staffelpreise
```

## Query Parameters

### Produktliste
- `q` - Suchquery
- `category` - Kategorie-ID (für Sidebar)
- `page` - Seitennummer (default: 1)
- `limit` - Einträge pro Seite (default: 12)

### Backend Query Params (API)
- `category_id` - Kategorie-Filter
- `manufacturer_id` - Hersteller-Filter
- `min_price` - Mindestpreis
- `max_price` - Maximalpreis
- `sort_by` - Sortierung (name, price, created_at)
- `sort_order` - Reihenfolge (asc, desc)

## Verwendung

### Produktliste aufrufen
```
/products                    # Alle Produkte
/products?q=hammer           # Suche nach "hammer"
/products?category=abc123    # Produkte der Kategorie
/products?page=2             # Seite 2
```

### Produktdetail
```
/products/{product-id}
```

### Kategorien
```
/categories                  # Übersicht
/categories/{category-id}    # Detail mit Produkten
```

## Fehlende Features (für später)

- [ ] Warenkorb-Funktionalität
- [ ] Favoriten / Merkliste
- [ ] Produktvergleich
- [ ] Erweiterte Filter (Preis-Range, Hersteller)
- [ ] Sortierung in UI
- [ ] Produktbewertungen
- [ ] Bilder-Galerie (Mehrere Bilder pro Produkt)
- [ ] PDF-Export (Produktdatenblatt)

## Dependencies

Keine neuen Dependencies hinzugefügt! Alle Features nutzen:
- React 18
- Next.js 14
- Tailwind CSS
- Bestehende UI-Komponenten (Panel, Button, Input)
- `class-variance-authority` (bereits installiert)

## Testing

Zum Testen:
1. Backend starten (Catalog Service muss laufen)
2. Frontend starten: `npm run dev`
3. In Browser: http://localhost:3000/products

## Anmerkungen

- Alle API-Calls gehen über den zentralen `apiClient`
- Auto-Refresh Token funktioniert auch für Catalog-Calls
- X-Tenant-ID Header wird automatisch gesetzt
- API Response Mapping (snake_case → camelCase) erfolgt automatisch
- Pagination ist URL-basiert (bookmarkable)
- Kategorie-Sidebar ist collapsible/expandable
