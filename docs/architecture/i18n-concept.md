# Internationalisierung (i18n) Konzept

## Grundprinzip

> **Keine fremdsprachigen Fragmente sichtbar.**
>
> Wenn ein Benutzer Deutsch wählt, muss ALLES auf Deutsch sein.
> Kein englischer Fallback, keine technischen Begriffe, keine Platzhalter.

---

## 1. Unterstützte Sprachen

| Locale | Sprache | Region | Priorität |
|--------|---------|--------|-----------|
| `de-CH` | Deutsch | Schweiz | Primary |
| `fr-CH` | Französisch | Schweiz | Secondary |
| `it-CH` | Italienisch | Schweiz | Secondary |
| `en-GB` | Englisch | UK | Fallback (nur intern) |

### Sprachauswahl

```
┌─────────────────────────────────────────────────────────────────┐
│                    SPRACHAUSWAHL-HIERARCHIE                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. URL Parameter        /de/produkte, /fr/produits             │
│         ↓                                                        │
│  2. User Preference      Gespeichert im Profil                  │
│         ↓                                                        │
│  3. Browser Language     Accept-Language Header                  │
│         ↓                                                        │
│  4. Tenant Default       Konfiguriert pro Tenant                │
│         ↓                                                        │
│  5. System Default       de-CH                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Übersetzungsebenen

### 2.1 Übersicht

```
┌─────────────────────────────────────────────────────────────────┐
│                     ÜBERSETZUNGSEBENEN                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ STATISCHE TEXTE (UI)                                    │    │
│  │                                                         │    │
│  │ • Labels, Buttons, Menüs                               │    │
│  │ • Fehlermeldungen                                       │    │
│  │ • Tooltips, Hilfetexte                                 │    │
│  │                                                         │    │
│  │ Quelle: JSON-Dateien im Frontend                       │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ DYNAMISCHE INHALTE (Content)                            │    │
│  │                                                         │    │
│  │ • Produktnamen, Beschreibungen                         │    │
│  │ • Kategorien                                            │    │
│  │ • CMS-Texte, Banner                                    │    │
│  │                                                         │    │
│  │ Quelle: Datenbank (translations Tabellen)              │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ SYSTEM-TEXTE (Backend)                                  │    │
│  │                                                         │    │
│  │ • API-Fehlermeldungen                                  │    │
│  │ • E-Mail-Templates                                      │    │
│  │ • PDF-Dokumente                                         │    │
│  │                                                         │    │
│  │ Quelle: Backend Translation Files                      │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ ATTRIBUT-WERTE (PIM)                                    │    │
│  │                                                         │    │
│  │ • Select-Optionen (Farbe, Material)                    │    │
│  │ • Einheiten (Stück, m², kg)                            │    │
│  │ • Zertifikate, Eigenschaften                           │    │
│  │                                                         │    │
│  │ Quelle: Akeneo PIM                                     │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. Frontend i18n

### 3.1 Bibliothek: next-intl

```tsx
// next.config.js
const withNextIntl = require('next-intl/plugin')();

module.exports = withNextIntl({
  i18n: {
    locales: ['de-CH', 'fr-CH', 'it-CH'],
    defaultLocale: 'de-CH',
  },
});
```

### 3.2 Übersetzungsdateien

```
apps/shop/messages/
├── de-CH.json
├── fr-CH.json
└── it-CH.json

apps/admin/messages/
├── de-CH.json
├── fr-CH.json
└── it-CH.json

apps/support/messages/
├── de-CH.json
├── fr-CH.json
└── it-CH.json
```

### 3.3 Struktur der Übersetzungsdateien

```json
// de-CH.json
{
  "common": {
    "save": "Speichern",
    "cancel": "Abbrechen",
    "delete": "Löschen",
    "edit": "Bearbeiten",
    "search": "Suchen",
    "loading": "Wird geladen...",
    "error": "Fehler",
    "success": "Erfolgreich"
  },
  "navigation": {
    "home": "Startseite",
    "products": "Produkte",
    "categories": "Kategorien",
    "cart": "Warenkorb",
    "account": "Mein Konto",
    "orders": "Bestellungen",
    "logout": "Abmelden"
  },
  "product": {
    "addToCart": "In den Warenkorb",
    "outOfStock": "Nicht verfügbar",
    "inStock": "Auf Lager",
    "price": "Preis",
    "quantity": "Menge",
    "unit": "Einheit",
    "description": "Beschreibung",
    "specifications": "Spezifikationen",
    "documents": "Dokumente",
    "relatedProducts": "Ähnliche Produkte"
  },
  "cart": {
    "title": "Warenkorb",
    "empty": "Ihr Warenkorb ist leer",
    "subtotal": "Zwischensumme",
    "shipping": "Versand",
    "tax": "MwSt.",
    "total": "Gesamtsumme",
    "checkout": "Zur Kasse",
    "continueShopping": "Weiter einkaufen"
  },
  "checkout": {
    "title": "Kasse",
    "billingAddress": "Rechnungsadresse",
    "shippingAddress": "Lieferadresse",
    "sameAsBilling": "Gleich wie Rechnungsadresse",
    "paymentMethod": "Zahlungsart",
    "shippingMethod": "Versandart",
    "placeOrder": "Bestellung aufgeben",
    "orderConfirmation": "Bestellbestätigung"
  },
  "errors": {
    "required": "Dieses Feld ist erforderlich",
    "invalidEmail": "Ungültige E-Mail-Adresse",
    "minLength": "Mindestens {min} Zeichen erforderlich",
    "maxLength": "Maximal {max} Zeichen erlaubt",
    "invalidPhone": "Ungültige Telefonnummer",
    "networkError": "Verbindungsfehler. Bitte versuchen Sie es erneut.",
    "serverError": "Ein Fehler ist aufgetreten. Bitte versuchen Sie es später erneut.",
    "notFound": "Seite nicht gefunden",
    "unauthorized": "Zugriff verweigert"
  },
  "units": {
    "piece": "Stück",
    "pieces": "Stück",
    "squareMeter": "m²",
    "linearMeter": "lfm",
    "kilogram": "kg",
    "liter": "l",
    "pack": "Paket",
    "pallet": "Palette"
  },
  "dateTime": {
    "today": "Heute",
    "yesterday": "Gestern",
    "tomorrow": "Morgen",
    "daysAgo": "vor {count} Tagen",
    "hoursAgo": "vor {count} Stunden",
    "minutesAgo": "vor {count} Minuten",
    "justNow": "Gerade eben"
  }
}
```

### 3.4 Verwendung in Komponenten

```tsx
// components/ProductCard.tsx
import { useTranslations } from 'next-intl';

export function ProductCard({ product }: { product: Product }) {
  const t = useTranslations('product');
  const tUnits = useTranslations('units');

  return (
    <Card>
      <CardContent>
        <h3>{product.name}</h3>
        <p>{formatPrice(product.price)} / {tUnits(product.unit)}</p>

        {product.inStock ? (
          <Badge variant="success">{t('inStock')}</Badge>
        ) : (
          <Badge variant="destructive">{t('outOfStock')}</Badge>
        )}

        <Button>{t('addToCart')}</Button>
      </CardContent>
    </Card>
  );
}
```

### 3.5 Pluralisierung

```json
// de-CH.json
{
  "cart": {
    "itemCount": "{count, plural, =0 {Keine Artikel} =1 {1 Artikel} other {# Artikel}}"
  }
}
```

```tsx
const t = useTranslations('cart');
t('itemCount', { count: items.length }); // "5 Artikel"
```

### 3.6 Datum & Zahlen

```tsx
import { useFormatter } from 'next-intl';

function PriceDisplay({ price }: { price: number }) {
  const format = useFormatter();

  return (
    <span>
      {format.number(price, {
        style: 'currency',
        currency: 'CHF'
      })}
    </span>
  );
  // Output: CHF 1'234.50 (de-CH)
  // Output: CHF 1 234,50 (fr-CH)
}

function DateDisplay({ date }: { date: Date }) {
  const format = useFormatter();

  return (
    <span>
      {format.dateTime(date, {
        day: 'numeric',
        month: 'long',
        year: 'numeric'
      })}
    </span>
  );
  // Output: 15. Januar 2024 (de-CH)
  // Output: 15 janvier 2024 (fr-CH)
}
```

---

## 4. Backend i18n

### 4.1 API-Fehlermeldungen

```go
// pkg/i18n/messages.go
package i18n

type MessageKey string

const (
    ErrNotFound         MessageKey = "error.not_found"
    ErrUnauthorized     MessageKey = "error.unauthorized"
    ErrValidation       MessageKey = "error.validation"
    ErrInternalServer   MessageKey = "error.internal_server"
    ErrInvalidInput     MessageKey = "error.invalid_input"
    ErrOrderFailed      MessageKey = "error.order_failed"
    ErrPaymentFailed    MessageKey = "error.payment_failed"
    ErrStockInsufficient MessageKey = "error.stock_insufficient"
)

var messages = map[string]map[MessageKey]string{
    "de-CH": {
        ErrNotFound:         "Die angeforderte Ressource wurde nicht gefunden.",
        ErrUnauthorized:     "Sie sind nicht berechtigt, diese Aktion durchzuführen.",
        ErrValidation:       "Die eingegebenen Daten sind ungültig.",
        ErrInternalServer:   "Ein interner Fehler ist aufgetreten. Bitte versuchen Sie es später erneut.",
        ErrInvalidInput:     "Ungültige Eingabe: {field}",
        ErrOrderFailed:      "Die Bestellung konnte nicht aufgegeben werden.",
        ErrPaymentFailed:    "Die Zahlung ist fehlgeschlagen.",
        ErrStockInsufficient: "Das Produkt \"{product}\" ist nicht in ausreichender Menge verfügbar.",
    },
    "fr-CH": {
        ErrNotFound:         "La ressource demandée n'a pas été trouvée.",
        ErrUnauthorized:     "Vous n'êtes pas autorisé à effectuer cette action.",
        ErrValidation:       "Les données saisies ne sont pas valides.",
        ErrInternalServer:   "Une erreur interne s'est produite. Veuillez réessayer plus tard.",
        ErrInvalidInput:     "Entrée invalide: {field}",
        ErrOrderFailed:      "La commande n'a pas pu être passée.",
        ErrPaymentFailed:    "Le paiement a échoué.",
        ErrStockInsufficient: "Le produit \"{product}\" n'est pas disponible en quantité suffisante.",
    },
    "it-CH": {
        ErrNotFound:         "La risorsa richiesta non è stata trovata.",
        ErrUnauthorized:     "Non sei autorizzato a eseguire questa azione.",
        ErrValidation:       "I dati inseriti non sono validi.",
        ErrInternalServer:   "Si è verificato un errore interno. Riprova più tardi.",
        ErrInvalidInput:     "Input non valido: {field}",
        ErrOrderFailed:      "Non è stato possibile effettuare l'ordine.",
        ErrPaymentFailed:    "Il pagamento non è riuscito.",
        ErrStockInsufficient: "Il prodotto \"{product}\" non è disponibile in quantità sufficiente.",
    },
}

func Translate(locale string, key MessageKey, params map[string]string) string {
    msg, ok := messages[locale][key]
    if !ok {
        // Fallback zu de-CH, NIEMALS zu en-GB für User-facing Messages
        msg = messages["de-CH"][key]
    }

    // Parameter ersetzen
    for k, v := range params {
        msg = strings.ReplaceAll(msg, "{"+k+"}", v)
    }

    return msg
}
```

### 4.2 API Response mit Locale

```go
// Handler extrahiert Locale aus Request
func (h *Handler) GetProduct(c *gin.Context) {
    locale := c.GetHeader("Accept-Language")
    if locale == "" {
        locale = "de-CH"
    }

    product, err := h.service.GetByID(ctx, id, locale)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{
                "code":    "NOT_FOUND",
                "message": i18n.Translate(locale, i18n.ErrNotFound, nil),
            },
        })
        return
    }

    c.JSON(http.StatusOK, product)
}
```

### 4.3 Lokalisierte Datenbank-Abfragen

```go
// Repository gibt nur die angeforderte Sprache zurück
func (r *repository) GetByID(ctx context.Context, id, locale string) (*Product, error) {
    var product Product

    err := r.db.GetContext(ctx, &product, `
        SELECT
            p.id,
            p.sku,
            p.price,
            pt.name,
            pt.description,
            pt.meta_title,
            pt.meta_description
        FROM products p
        LEFT JOIN product_translations pt ON p.id = pt.product_id
            AND pt.locale = $2
        WHERE p.id = $1
    `, id, locale)

    return &product, err
}
```

---

## 5. Datenbank-Schema

### 5.1 Translations-Tabellen Pattern

```sql
-- Produkte: Stammdaten
CREATE TABLE products (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    sku VARCHAR(50) NOT NULL,
    price DECIMAL(12, 5) NOT NULL,
    -- ... nicht-lokalisierte Felder
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Produkte: Übersetzungen
CREATE TABLE product_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    locale VARCHAR(5) NOT NULL,  -- 'de-CH', 'fr-CH', 'it-CH'

    -- Lokalisierte Felder
    name VARCHAR(255) NOT NULL,
    description TEXT,
    short_description TEXT,
    meta_title VARCHAR(70),
    meta_description VARCHAR(160),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(product_id, locale)
);

CREATE INDEX idx_product_translations_locale ON product_translations(product_id, locale);
```

### 5.2 Alle Entitäten mit Übersetzungen

```sql
-- Kategorien
CREATE TABLE category_translations (
    category_id UUID NOT NULL REFERENCES categories(id),
    locale VARCHAR(5) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    meta_title VARCHAR(70),
    meta_description VARCHAR(160),
    UNIQUE(category_id, locale)
);

-- Attribut-Optionen (z.B. Farben, Materialien)
CREATE TABLE attribute_option_translations (
    option_id UUID NOT NULL REFERENCES attribute_options(id),
    locale VARCHAR(5) NOT NULL,
    label VARCHAR(255) NOT NULL,
    UNIQUE(option_id, locale)
);

-- CMS Textblöcke
CREATE TABLE textblock_translations (
    textblock_id UUID NOT NULL REFERENCES textblocks(id),
    locale VARCHAR(5) NOT NULL,
    title VARCHAR(255),
    content TEXT NOT NULL,
    UNIQUE(textblock_id, locale)
);

-- E-Mail Templates
CREATE TABLE email_template_translations (
    template_id UUID NOT NULL REFERENCES email_templates(id),
    locale VARCHAR(5) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body_html TEXT NOT NULL,
    body_text TEXT,
    UNIQUE(template_id, locale)
);
```

---

## 6. Fallback-Strategie

### 6.1 Strikte Regel: Kein sichtbarer Fallback

```
❌ FALSCH: Englischen Text anzeigen wenn Deutsch fehlt
   "Color: Red" (statt "Farbe: Rot")

❌ FALSCH: Technische Codes anzeigen
   "STATUS_PENDING" (statt "In Bearbeitung")

❌ FALSCH: Leere Felder anzeigen
   "Beschreibung: " (ohne Text)

✅ RICHTIG: Immer vollständige Übersetzung
   "Farbe: Rot"
   "Status: In Bearbeitung"
   Feld komplett ausblenden wenn kein Inhalt
```

### 6.2 Interner Fallback (nur für Entwickler/Logs)

```go
// Nur für interne Logs, niemals für User-facing Content
func TranslateInternal(locale string, key MessageKey) string {
    if msg, ok := messages[locale][key]; ok {
        return msg
    }
    if msg, ok := messages["de-CH"][key]; ok {
        return msg
    }
    // Nur im Log, niemals zum User
    return string(key)
}
```

### 6.3 Validierung bei Import (Akeneo)

```go
// Beim Import: Fehlende Übersetzungen blockieren
func ValidateProductTranslations(product *AkeneoProduct) error {
    requiredLocales := []string{"de-CH", "fr-CH", "it-CH"}
    requiredFields := []string{"name", "description"}

    for _, locale := range requiredLocales {
        for _, field := range requiredFields {
            if product.GetTranslation(locale, field) == "" {
                return fmt.Errorf(
                    "fehlende Übersetzung: Produkt %s, Feld %s, Sprache %s",
                    product.SKU, field, locale,
                )
            }
        }
    }
    return nil
}
```

### 6.4 Fehlende Übersetzungen melden

```go
// Monitoring: Fehlende Übersetzungen tracken
type MissingTranslation struct {
    EntityType string    // "product", "category"
    EntityID   string
    Locale     string
    Field      string
    DetectedAt time.Time
}

func (s *TranslationService) ReportMissing(ctx context.Context, mt *MissingTranslation) {
    // 1. In Datenbank speichern für Admin-Übersicht
    s.repo.SaveMissingTranslation(ctx, mt)

    // 2. Alert für Support/Admin
    s.alertService.CreateAlert(ctx, &Alert{
        Type:     "missing_translation",
        Severity: "warning",
        Title:    fmt.Sprintf("Fehlende Übersetzung: %s", mt.Field),
        Metadata: map[string]interface{}{
            "entity_type": mt.EntityType,
            "entity_id":   mt.EntityID,
            "locale":      mt.Locale,
            "field":       mt.Field,
        },
    })
}
```

---

## 7. Einheiten & Formatierung

### 7.1 Einheiten-Übersetzungen

```json
// de-CH.json
{
  "units": {
    "STK": "Stück",
    "M2": "m²",
    "LFM": "lfm",
    "KG": "kg",
    "L": "l",
    "PAK": "Paket",
    "PAL": "Palette",
    "KRT": "Karton",
    "ROL": "Rolle"
  }
}
```

```json
// fr-CH.json
{
  "units": {
    "STK": "pièce",
    "M2": "m²",
    "LFM": "ml",
    "KG": "kg",
    "L": "l",
    "PAK": "paquet",
    "PAL": "palette",
    "KRT": "carton",
    "ROL": "rouleau"
  }
}
```

### 7.2 Zahlenformatierung

| Locale | Beispiel | Tausender | Dezimal |
|--------|----------|-----------|---------|
| de-CH | 1'234.50 | ' | . |
| fr-CH | 1 234,50 | (space) | , |
| it-CH | 1'234.50 | ' | . |

```tsx
// Konsistente Formatierung
function formatNumber(value: number, locale: string): string {
  return new Intl.NumberFormat(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value);
}

// CHF immer mit CHF prefix, nicht Fr.
function formatPrice(value: number, locale: string): string {
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: 'CHF',
  }).format(value);
}
```

### 7.3 Datumsformatierung

| Locale | Kurz | Lang |
|--------|------|------|
| de-CH | 15.01.2024 | 15. Januar 2024 |
| fr-CH | 15.01.2024 | 15 janvier 2024 |
| it-CH | 15.01.2024 | 15 gennaio 2024 |

---

## 8. Admin & Support Portal

### 8.1 Sprache des Portals

```tsx
// Admin/Support Portal: Sprache aus User-Profil
function AdminLayout({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const locale = user?.preferredLocale || 'de-CH';

  return (
    <NextIntlClientProvider locale={locale} messages={messages[locale]}>
      {children}
    </NextIntlClientProvider>
  );
}
```

### 8.2 Mehrsprachige Daten bearbeiten

```tsx
// Produkt-Editor mit Tabs pro Sprache
function ProductEditor({ product }: { product: Product }) {
  const locales = ['de-CH', 'fr-CH', 'it-CH'];
  const [activeLocale, setActiveLocale] = useState('de-CH');

  return (
    <div>
      <Tabs value={activeLocale} onValueChange={setActiveLocale}>
        <TabsList>
          {locales.map(locale => (
            <TabsTrigger key={locale} value={locale}>
              {localeNames[locale]}
              {!product.translations[locale]?.name && (
                <Badge variant="destructive" className="ml-2">!</Badge>
              )}
            </TabsTrigger>
          ))}
        </TabsList>

        {locales.map(locale => (
          <TabsContent key={locale} value={locale}>
            <TranslationForm
              locale={locale}
              data={product.translations[locale]}
            />
          </TabsContent>
        ))}
      </Tabs>
    </div>
  );
}
```

### 8.3 Übersetzungs-Status Übersicht

```tsx
// Admin: Übersicht fehlender Übersetzungen
function TranslationStatus() {
  const { data: stats } = useTranslationStats();

  return (
    <Card>
      <CardHeader>
        <CardTitle>Übersetzungs-Status</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Bereich</TableHead>
              <TableHead>DE</TableHead>
              <TableHead>FR</TableHead>
              <TableHead>IT</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>Produkte</TableCell>
              <TableCell><CompletionBadge value={stats.products.de} /></TableCell>
              <TableCell><CompletionBadge value={stats.products.fr} /></TableCell>
              <TableCell><CompletionBadge value={stats.products.it} /></TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Kategorien</TableCell>
              <TableCell><CompletionBadge value={stats.categories.de} /></TableCell>
              <TableCell><CompletionBadge value={stats.categories.fr} /></TableCell>
              <TableCell><CompletionBadge value={stats.categories.it} /></TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
```

---

## 9. E-Mail & Dokumente

### 9.1 E-Mail Templates

```go
// E-Mail in Kundensprache versenden
func (s *NotificationService) SendOrderConfirmation(ctx context.Context, order *Order) error {
    // Sprache des Kunden verwenden
    locale := order.Customer.PreferredLocale
    if locale == "" {
        locale = "de-CH"
    }

    template, err := s.templateRepo.GetByKey(ctx, "order_confirmation", locale)
    if err != nil {
        return err
    }

    // Template rendern
    body, err := s.renderTemplate(template, map[string]interface{}{
        "order":    order,
        "customer": order.Customer,
    })

    return s.emailSender.Send(ctx, &Email{
        To:      order.Customer.Email,
        Subject: template.Subject,
        Body:    body,
    })
}
```

### 9.2 PDF-Generierung

```go
// Auftragsbestätigung in Kundensprache
func (s *PDFService) GenerateOrderConfirmation(ctx context.Context, order *Order) ([]byte, error) {
    locale := order.Customer.PreferredLocale

    // Lokalisierte Labels laden
    labels := s.loadLabels(locale)

    // PDF generieren
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()

    // Header mit lokalisierten Texten
    pdf.Cell(0, 10, labels["order_confirmation"])
    pdf.Cell(0, 10, fmt.Sprintf("%s: %s", labels["order_number"], order.Number))
    pdf.Cell(0, 10, fmt.Sprintf("%s: %s", labels["date"], formatDate(order.CreatedAt, locale)))

    // ... weitere lokalisierte Inhalte

    return pdf.Output()
}
```

---

## 10. SEO & URLs

### 10.1 Lokalisierte URLs

```
Deutsch:   /de/bodenbelaege/laminat/eiche-natur
Französisch: /fr/revetements-de-sol/stratifie/chene-naturel
Italienisch: /it/rivestimenti/laminato/rovere-naturale
```

### 10.2 Hreflang Tags

```tsx
// Automatische hreflang Tags
function ProductPage({ product }: { product: Product }) {
  const locales = ['de-CH', 'fr-CH', 'it-CH'];

  return (
    <>
      <Head>
        {locales.map(locale => (
          <link
            key={locale}
            rel="alternate"
            hrefLang={locale}
            href={`https://shop.example.com/${locale}/${product.slugs[locale]}`}
          />
        ))}
        <link
          rel="alternate"
          hrefLang="x-default"
          href={`https://shop.example.com/de-CH/${product.slugs['de-CH']}`}
        />
      </Head>
      {/* ... */}
    </>
  );
}
```

### 10.3 Sitemap mit Sprachen

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml">
  <url>
    <loc>https://shop.example.com/de-CH/produkte/eiche-laminat</loc>
    <xhtml:link rel="alternate" hreflang="de-CH"
                href="https://shop.example.com/de-CH/produkte/eiche-laminat"/>
    <xhtml:link rel="alternate" hreflang="fr-CH"
                href="https://shop.example.com/fr-CH/produits/stratifie-chene"/>
    <xhtml:link rel="alternate" hreflang="it-CH"
                href="https://shop.example.com/it-CH/prodotti/laminato-rovere"/>
  </url>
</urlset>
```

---

## 11. Qualitätssicherung

### 11.1 Automatische Prüfungen

```typescript
// CI/CD: Übersetzungs-Vollständigkeit prüfen
describe('Translation completeness', () => {
  const baseLocale = 'de-CH';
  const locales = ['fr-CH', 'it-CH'];

  test('all locales have same keys as base', () => {
    const baseKeys = getAllKeys(messages[baseLocale]);

    for (const locale of locales) {
      const localeKeys = getAllKeys(messages[locale]);
      const missing = baseKeys.filter(k => !localeKeys.includes(k));

      expect(missing).toEqual([]);
    }
  });

  test('no empty translations', () => {
    for (const locale of [baseLocale, ...locales]) {
      const emptyKeys = findEmptyValues(messages[locale]);
      expect(emptyKeys).toEqual([]);
    }
  });
});
```

### 11.2 Lint-Regeln

```typescript
// ESLint Plugin für i18n
// eslint-plugin-i18n-checker

// ❌ Nicht erlaubt: Hardcoded Strings
<Button>Save</Button>

// ✅ Richtig: Übersetzungsfunktion
<Button>{t('common.save')}</Button>
```

### 11.3 Übersetzungs-Review Workflow

```
1. Entwickler fügt neuen Text hinzu (de-CH)
2. CI markiert fehlende Übersetzungen
3. Übersetzer ergänzen fr-CH, it-CH
4. Review durch Native Speaker
5. Merge nur wenn alle Sprachen vollständig
```

---

## 12. Zusammenfassung

### Grundregeln

1. **Keine englischen Fallbacks** für User-facing Content
2. **Alle Pflichtfelder** müssen in allen 3 Sprachen vorhanden sein
3. **Einheiten, Währungen, Daten** werden korrekt lokalisiert
4. **Fehlende Übersetzungen** werden gemeldet und blockieren ggf. den Import
5. **Admin/Support** kann Übersetzungsstatus einsehen

### Checkliste für neue Features

```
□ UI-Texte in allen 3 JSON-Dateien hinzugefügt
□ API-Fehlermeldungen in allen 3 Sprachen
□ Datenbank-Schema mit _translations Tabelle
□ Tests für Übersetzungs-Vollständigkeit
□ Dokumentation aktualisiert
```

### Tech Stack

| Bereich | Lösung |
|---------|--------|
| Frontend | next-intl |
| Backend | Custom i18n Package |
| Datenbank | Translations-Tabellen |
| Validierung | CI/CD Tests |
| Monitoring | Missing Translation Alerts |
