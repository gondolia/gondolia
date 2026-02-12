# CMS Konzept für Webshop V3

## 1. Hintergrund

### 1.1 Die Anfrage von Marketing (2023)

Vor zwei Jahren kam vom Marketing die Anfrage, die **Typo3-Seiten auf Statamic zu konsolidieren**. Das Ziel war nachvollziehbar:

- **Ein System** statt zwei verschiedene
- **Einheitliches Design** über alle Seiten
- **Eine Domain** für besseres SEO
- **Weniger Pflegeaufwand** für das Marketing-Team

### 1.2 Warum wir damals abgelehnt haben

Die Übernahme der Typo3-Inhalte in Statamic war aus technischer Sicht nicht sinnvoll:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PROBLEME MIT STATAMIC (V2)                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. FILE-BASED STORAGE                                                       │
│     ─────────────────────                                                   │
│     Statamic speichert alles in YAML/Markdown-Dateien.                      │
│     Bei Bulk-Operationen (viele Seiten gleichzeitig) gab es                 │
│     regelmässig Probleme: Race Conditions, Git-Konflikte,                   │
│     korrupte Dateien.                                                       │
│                                                                              │
│     → Mehr Typo3-Content = mehr Probleme                                    │
│                                                                              │
│  2. MONOLITH-KOPPLUNG                                                        │
│     ────────────────────                                                    │
│     Statamic ist IN den Laravel-Shop eingewoben.                            │
│     Jedes Shop-Update birgt das Risiko, Statamic zu brechen.               │
│     Jedes Statamic-Update muss mit Shop-Releases koordiniert werden.       │
│                                                                              │
│     → Mehr Content = mehr Risiko bei Updates                                │
│                                                                              │
│  3. ARBEITSAUFWAND SHOP-TEAM                                                 │
│     ────────────────────────────                                            │
│     Migration, Templating, Testing - alles hätte das Shop-Team             │
│     machen müssen, während gleichzeitig der Shop weiterlaufen muss.        │
│                                                                              │
│     → Ressourcen-Konflikt mit Shop-Entwicklung                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Fazit damals:** Die Konsolidierung wäre technisch möglich gewesen, hätte aber die bestehenden Statamic-Probleme verschärft und erhebliche Ressourcen gebunden - ohne die grundlegenden Architektur-Probleme zu lösen.

---

## 2. Die Chance mit V3

### 2.1 Warum jetzt der richtige Zeitpunkt ist

Mit Webshop V3 bauen wir die gesamte Architektur neu auf. Das ist die Gelegenheit, die **Marketing-Anfrage von 2023 richtig umzusetzen** - mit einer Lösung, die beide Probleme löst:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  V3 = NEUSTART                                                               │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ✓ Kein Monolith mehr      → CMS kann unabhängig laufen            │   │
│  │  ✓ Kubernetes-Infrastruktur → Separate Deployments möglich          │   │
│  │  ✓ API-First Architektur    → CMS als Content-API                   │   │
│  │  ✓ PostgreSQL überall       → Echte Datenbank statt Flat-Files      │   │
│  │  ✓ Next.js Frontend         → Moderne CMS-Integration               │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  → Wir können jetzt machen, was 2023 nicht sinnvoll war.                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Was Marketing bekommt

| Wunsch von 2023 | Umsetzung in V3 |
|-----------------|-----------------|
| Ein System statt zwei | ✅ Ein Headless CMS für alles |
| Einheitliches Design | ✅ Gleiche Komponenten wie Shop |
| Eine Domain | ✅ shop.ch/* für alle Seiten |
| Weniger Pflegeaufwand | ✅ Ein Login, eine Oberfläche |
| Besseres SEO | ✅ Eine Sitemap, keine Fragmentierung |

**Zusätzlich** (was 2023 nicht möglich gewesen wäre):

- **Preview Mode** - Drafts im echten Shop-Design sehen
- **Scheduled Publishing** - Inhalte zeitgesteuert veröffentlichen
- **Multi-Channel** - Gleicher Content für App, Displays, Newsletter
- **Block-basierter Editor** - Flexibles Seitenlayout ohne Entwickler

---

## 3. Aktuelle Situation (Ist-Zustand)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          AKTUELLER ZUSTAND                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────┐    ┌───────────────────────────────┐    │
│  │     WEBSHOP (Laravel)          │    │         TYPO3                  │    │
│  │                                │    │                                │    │
│  │  ┌──────────────────────────┐ │    │  • Marketing-Seiten           │    │
│  │  │      STATAMIC            │ │    │  • Landingpages               │    │
│  │  │                          │ │    │  • News/Blog                  │    │
│  │  │  • Content-Seiten        │ │    │  • Kampagnen                  │    │
│  │  │  • Produktinfos          │ │    │                                │    │
│  │  │  • Banner                │ │    │  Eigene Domain                │    │
│  │  │  • Textblöcke            │ │    │  Separates Hosting            │    │
│  │  │                          │ │    │  Eigenes Design               │    │
│  │  │  File-Based Storage      │ │    │                                │    │
│  │  └──────────────────────────┘ │    └───────────────────────────────┘    │
│  │                                │                                         │
│  └───────────────────────────────┘                                         │
│                                                                              │
│  Probleme Marketing:              Probleme IT:                              │
│  • Doppelte Pflege               • Update-Konflikte (Statamic)            │
│  • Inkonsistentes Design         • Bulk-Operationen scheitern             │
│  • SEO-Fragmentierung            • Tight Coupling an Monolith             │
│  • Zwei Logins, zwei Systeme     • Wartung von zwei Systemen              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Ziel-Architektur

### 4.1 Headless CMS als Unified Content Hub

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          V3 CONTENT ARCHITEKTUR                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                        ┌─────────────────────────┐                          │
│                        │     HEADLESS CMS        │                          │
│                        │                         │                          │
│                        │  • Content-Seiten       │                          │
│                        │  • Landingpages         │                          │
│                        │  • Banner/Teaser        │                          │
│                        │  • News/Blog            │                          │
│                        │  • Produkttexte         │                          │
│                        │  • FAQ                  │                          │
│                        │  • Rechtliches          │                          │
│                        │                         │                          │
│                        │  PostgreSQL Backend     │                          │
│                        │  Media in S3/Blob       │                          │
│                        │  GraphQL + REST API     │                          │
│                        └───────────┬─────────────┘                          │
│                                    │                                         │
│                                    │ API                                     │
│                                    │                                         │
│         ┌──────────────────────────┼──────────────────────────┐             │
│         │                          │                          │             │
│         ▼                          ▼                          ▼             │
│  ┌─────────────┐           ┌─────────────┐           ┌─────────────┐       │
│  │   WEBSHOP   │           │  MARKETING  │           │   WEITERE   │       │
│  │  (Next.js)  │           │   PAGES     │           │   KANÄLE    │       │
│  │             │           │  (Next.js)  │           │             │       │
│  │ shop.ch     │           │ shop.ch/*   │           │ • App       │       │
│  │             │           │             │           │ • Displays  │       │
│  │ Produktseiten│           │ /magazin    │           │ • Emails    │       │
│  │ Checkout    │           │ /ratgeber   │           │             │       │
│  │ Account     │           │ /aktionen   │           │             │       │
│  └─────────────┘           └─────────────┘           └─────────────┘       │
│                                                                              │
│         EINE Domain         │         EINE Codebasis        │               │
│         EINE Sitemap        │         EIN Design-System     │               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Vorteile der neuen Architektur

| Aspekt | Vorher | Nachher |
|--------|--------|---------|
| **Systeme** | 2 (Statamic + Typo3) | 1 Headless CMS |
| **Domains** | 2 (shop + marketing) | 1 unified |
| **Deployments** | Gekoppelt | Unabhängig |
| **Content-API** | Keine | GraphQL + REST |
| **Datenbank** | File-based | PostgreSQL |
| **Media Storage** | Lokal | S3/Azure Blob |
| **Skalierung** | Problematisch | Horizontal |
| **Multi-Channel** | Nicht möglich | API-first |

---

## 5. CMS-Evaluation

### 5.1 Anforderungen

| Anforderung | Priorität | Beschreibung |
|-------------|-----------|--------------|
| **Headless/API-first** | Must | GraphQL oder REST API |
| **Self-hosted möglich** | Must | Daten in eigener Infrastruktur |
| **PostgreSQL Support** | Must | Keine File-based DB |
| **Media Library** | Must | S3/Azure Blob Integration |
| **Mehrsprachigkeit** | Must | de-CH, fr-CH, it-CH |
| **Visual Editor** | Should | WYSIWYG für Marketing |
| **Webhooks** | Should | Cache Invalidation |
| **Roles & Permissions** | Should | Marketing vs. Admin |
| **Preview Mode** | Should | Draft-Vorschau im Frontend |
| **Versionierung** | Nice | Content History |
| **Workflows** | Nice | Approval Prozesse |

### 5.2 Kandidaten-Vergleich

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          CMS KANDIDATEN                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │    STRAPI       │  │    PAYLOAD      │  │   DIRECTUS      │             │
│  │                 │  │                 │  │                 │             │
│  │  Node.js        │  │  Node.js        │  │  Node.js        │             │
│  │  Open Source    │  │  Open Source    │  │  Open Source    │             │
│  │  PostgreSQL ✓   │  │  PostgreSQL ✓   │  │  PostgreSQL ✓   │             │
│  │  REST + GraphQL │  │  REST + GraphQL │  │  REST + GraphQL │             │
│  │  Self-hosted ✓  │  │  Self-hosted ✓  │  │  Self-hosted ✓  │             │
│  │                 │  │                 │  │                 │             │
│  │  ⭐ Marktführer │  │  ⭐ TypeScript  │  │  ⭐ SQL-first   │             │
│  │  ⭐ Große       │  │    native       │  │  ⭐ Bestehendes │             │
│  │    Community    │  │  ⭐ Code-first  │  │    Schema       │             │
│  │  ⭐ Plugin-     │  │  ⭐ Next.js     │  │    nutzbar      │             │
│  │    Ökosystem    │  │    Integration  │  │                 │             │
│  │                 │  │                 │  │                 │             │
│  │  ⚠️ v5 Breaking │  │  ⚠️ Jüngeres   │  │  ⚠️ Weniger    │             │
│  │    Changes      │  │    Projekt      │  │    Flexibilität │             │
│  │                 │  │                 │  │                 │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  CONTENTFUL     │  │    SANITY       │  │   STORYBLOK     │             │
│  │     (SaaS)      │  │     (SaaS)      │  │     (SaaS)      │             │
│  │                 │  │                 │  │                 │             │
│  │  ⭐ Enterprise- │  │  ⭐ Real-time   │  │  ⭐ Visual      │             │
│  │    grade        │  │    Collab       │  │    Editor       │             │
│  │  ⭐ Keine       │  │  ⭐ GROQ Query  │  │  ⭐ Component-  │             │
│  │    Wartung      │  │  ⭐ Portable    │  │    based        │             │
│  │                 │  │    Text         │  │                 │             │
│  │  ⚠️ $$$        │  │  ⚠️ $$         │  │  ⚠️ $$         │             │
│  │  ⚠️ Vendor     │  │  ⚠️ Learning   │  │  ⚠️ Vendor     │             │
│  │    Lock-in      │  │    Curve        │  │    Lock-in      │             │
│  │                 │  │                 │  │                 │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Detaillierter Vergleich

| Kriterium | Strapi | Payload | Directus | Sanity | Storyblok |
|-----------|--------|---------|----------|--------|-----------|
| **Lizenz** | MIT | MIT | GPL/BSL | Proprietary | Proprietary |
| **Self-hosted** | ✅ | ✅ | ✅ | ❌ | ❌ |
| **PostgreSQL** | ✅ | ✅ | ✅ | N/A | N/A |
| **GraphQL** | ✅ | ✅ | ✅ | ✅ (GROQ) | ✅ |
| **i18n** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Visual Editor** | Plugin | ✅ | ❌ | ✅ | ✅✅ |
| **TypeScript** | Partial | ✅ Native | Partial | ✅ | N/A |
| **Media Library** | ✅ S3 | ✅ S3 | ✅ S3 | ✅ | ✅ |
| **Webhooks** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Preview Mode** | Manual | ✅ Native | Manual | ✅ | ✅ |
| **Community** | Sehr groß | Wachsend | Mittel | Groß | Mittel |
| **Kosten (Self)** | $0 | $0 | $0 | - | - |
| **Kosten (Cloud)** | Ab $99/m | Ab $199/m | Ab $99/m | Ab $99/m | Ab $99/m |

### 5.4 Empfehlung

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  EMPFEHLUNG: PAYLOAD CMS                                                    │
│  ═══════════════════════                                                    │
│                                                                              │
│  Gründe:                                                                     │
│                                                                              │
│  1. TypeScript Native                                                        │
│     → Passt zu Next.js Stack                                                │
│     → Type-safe Content Queries                                             │
│     → Bessere DX                                                            │
│                                                                              │
│  2. Code-First Schema                                                        │
│     → Schema als Code (versioniert)                                         │
│     → Keine UI-Klickerei für Änderungen                                    │
│     → CI/CD für Content-Typen                                              │
│                                                                              │
│  3. Native Next.js Integration                                              │
│     → Kann IN Next.js App laufen                                           │
│     → Shared Types                                                          │
│     → Preview Mode out-of-box                                               │
│                                                                              │
│  4. PostgreSQL + S3                                                          │
│     → Gleiche Infrastruktur wie Shop                                        │
│     → Keine neuen Systeme                                                   │
│                                                                              │
│  5. Visual Editor (Lexical)                                                  │
│     → Rich Text mit Blöcken                                                 │
│     → Marketing-freundlich                                                  │
│                                                                              │
│  Alternative: Strapi                                                        │
│  → Falls größere Plugin-Auswahl benötigt                                   │
│  → Falls Team bereits Strapi-Erfahrung hat                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Payload CMS Architektur

### 6.1 Deployment-Optionen

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       DEPLOYMENT OPTIONEN                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  OPTION A: Standalone (Empfohlen)                                           │
│  ─────────────────────────────────                                          │
│                                                                              │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                   │
│  │  Payload    │     │   Next.js   │     │  Go Backend │                   │
│  │    CMS      │     │   Frontend  │     │  Services   │                   │
│  │             │     │             │     │             │                   │
│  │ cms.shop.ch │     │  shop.ch    │     │ api.shop.ch │                   │
│  │             │◀────│             │────▶│             │                   │
│  │ Port 3001   │ API │ Port 3000   │ API │ Port 8080   │                   │
│  └──────┬──────┘     └─────────────┘     └─────────────┘                   │
│         │                                                                    │
│         ▼                                                                    │
│  ┌─────────────┐     ┌─────────────┐                                       │
│  │ PostgreSQL  │     │  S3 / Blob  │                                       │
│  │  (shared)   │     │   (Media)   │                                       │
│  └─────────────┘     └─────────────┘                                       │
│                                                                              │
│  Vorteile:                                                                   │
│  • Unabhängige Deployments                                                  │
│  • CMS kann neugestartet werden ohne Shop-Impact                           │
│  • Eigene Skalierung                                                        │
│                                                                              │
│  ────────────────────────────────────────────────────────────────────────   │
│                                                                              │
│  OPTION B: Embedded in Next.js                                              │
│  ─────────────────────────────────                                          │
│                                                                              │
│  ┌─────────────────────────────────┐     ┌─────────────┐                   │
│  │         Next.js App             │     │  Go Backend │                   │
│  │                                 │     │  Services   │                   │
│  │  ┌───────────┐  ┌───────────┐  │     │             │                   │
│  │  │  Payload  │  │   Shop    │  │────▶│ api.shop.ch │                   │
│  │  │  /admin   │  │   Pages   │  │     │             │                   │
│  │  └───────────┘  └───────────┘  │     └─────────────┘                   │
│  │                                 │                                        │
│  │  shop.ch                        │                                        │
│  └─────────────────────────────────┘                                       │
│                                                                              │
│  Vorteile:                                                                   │
│  • Ein Deployment                                                            │
│  • Shared Types                                                              │
│  • Einfacher Preview Mode                                                   │
│                                                                              │
│  Nachteile:                                                                  │
│  • Gekoppelte Releases                                                       │
│  • Größeres Bundle                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Empfehlung:** Option A (Standalone) - vermeidet die Probleme des aktuellen Setups.

### 6.2 Content-Modell

```typescript
// payload.config.ts
import { buildConfig } from 'payload/config';
import { postgresAdapter } from '@payloadcms/db-postgres';
import { s3Adapter } from '@payloadcms/plugin-cloud-storage/s3';
import { lexicalEditor } from '@payloadcms/richtext-lexical';

export default buildConfig({
  admin: {
    user: 'users',
  },

  db: postgresAdapter({
    pool: {
      connectionString: process.env.DATABASE_URL,
    },
  }),

  editor: lexicalEditor({}),

  collections: [
    // === PAGES ===
    {
      slug: 'pages',
      admin: {
        useAsTitle: 'title',
        defaultColumns: ['title', 'slug', 'status', 'updatedAt'],
      },
      versions: {
        drafts: true,
      },
      fields: [
        {
          name: 'title',
          type: 'text',
          required: true,
          localized: true,
        },
        {
          name: 'slug',
          type: 'text',
          required: true,
          unique: true,
          admin: {
            position: 'sidebar',
          },
        },
        {
          name: 'template',
          type: 'select',
          options: [
            { label: 'Standard', value: 'standard' },
            { label: 'Landing Page', value: 'landing' },
            { label: 'Ratgeber', value: 'guide' },
            { label: 'Aktion', value: 'campaign' },
          ],
          defaultValue: 'standard',
        },
        {
          name: 'content',
          type: 'blocks',
          localized: true,
          blocks: [
            HeroBlock,
            TextBlock,
            ImageBlock,
            ProductGridBlock,
            CTABlock,
            FAQBlock,
            VideoBlock,
            TestimonialBlock,
          ],
        },
        {
          name: 'seo',
          type: 'group',
          fields: [
            { name: 'metaTitle', type: 'text', localized: true },
            { name: 'metaDescription', type: 'textarea', localized: true },
            { name: 'ogImage', type: 'upload', relationTo: 'media' },
          ],
        },
      ],
    },

    // === BLOG/NEWS ===
    {
      slug: 'posts',
      admin: {
        useAsTitle: 'title',
      },
      versions: {
        drafts: true,
      },
      fields: [
        { name: 'title', type: 'text', required: true, localized: true },
        { name: 'slug', type: 'text', required: true, unique: true },
        { name: 'excerpt', type: 'textarea', localized: true },
        { name: 'featuredImage', type: 'upload', relationTo: 'media' },
        { name: 'content', type: 'richText', localized: true },
        { name: 'category', type: 'relationship', relationTo: 'post-categories' },
        { name: 'author', type: 'relationship', relationTo: 'users' },
        { name: 'publishedAt', type: 'date' },
      ],
    },

    // === PRODUCT CONTENT ===
    {
      slug: 'product-content',
      admin: {
        useAsTitle: 'sku',
        description: 'Zusätzlicher Content für Produkte (über PIM hinaus)',
      },
      fields: [
        { name: 'sku', type: 'text', required: true, unique: true },
        { name: 'additionalDescription', type: 'richText', localized: true },
        { name: 'applicationTips', type: 'richText', localized: true },
        { name: 'videos', type: 'array', fields: [
          { name: 'title', type: 'text', localized: true },
          { name: 'youtubeId', type: 'text' },
        ]},
        { name: 'downloads', type: 'array', fields: [
          { name: 'title', type: 'text', localized: true },
          { name: 'file', type: 'upload', relationTo: 'media' },
        ]},
      ],
    },

    // === BANNERS ===
    {
      slug: 'banners',
      admin: {
        useAsTitle: 'name',
      },
      fields: [
        { name: 'name', type: 'text', required: true },
        { name: 'placement', type: 'select', options: [
          { label: 'Homepage Hero', value: 'homepage_hero' },
          { label: 'Homepage Secondary', value: 'homepage_secondary' },
          { label: 'Category Header', value: 'category_header' },
          { label: 'Cart Upsell', value: 'cart_upsell' },
          { label: 'Checkout', value: 'checkout' },
        ]},
        { name: 'image', type: 'upload', relationTo: 'media', required: true },
        { name: 'imageMobile', type: 'upload', relationTo: 'media' },
        { name: 'title', type: 'text', localized: true },
        { name: 'subtitle', type: 'text', localized: true },
        { name: 'ctaText', type: 'text', localized: true },
        { name: 'ctaLink', type: 'text' },
        { name: 'startDate', type: 'date' },
        { name: 'endDate', type: 'date' },
        { name: 'active', type: 'checkbox', defaultValue: true },
      ],
    },

    // === FAQ ===
    {
      slug: 'faqs',
      admin: {
        useAsTitle: 'question',
      },
      fields: [
        { name: 'question', type: 'text', required: true, localized: true },
        { name: 'answer', type: 'richText', required: true, localized: true },
        { name: 'category', type: 'relationship', relationTo: 'faq-categories' },
        { name: 'order', type: 'number' },
      ],
    },

    // === NAVIGATION ===
    {
      slug: 'navigation',
      admin: {
        useAsTitle: 'name',
      },
      fields: [
        { name: 'name', type: 'text', required: true },
        { name: 'location', type: 'select', options: [
          { label: 'Header Main', value: 'header_main' },
          { label: 'Header Secondary', value: 'header_secondary' },
          { label: 'Footer', value: 'footer' },
          { label: 'Footer Legal', value: 'footer_legal' },
        ]},
        { name: 'items', type: 'array', fields: [
          { name: 'label', type: 'text', required: true, localized: true },
          { name: 'type', type: 'select', options: [
            { label: 'Internal Page', value: 'page' },
            { label: 'Category', value: 'category' },
            { label: 'External URL', value: 'external' },
          ]},
          { name: 'page', type: 'relationship', relationTo: 'pages',
            admin: { condition: (_, siblingData) => siblingData?.type === 'page' }},
          { name: 'categoryId', type: 'text',
            admin: { condition: (_, siblingData) => siblingData?.type === 'category' }},
          { name: 'url', type: 'text',
            admin: { condition: (_, siblingData) => siblingData?.type === 'external' }},
          { name: 'children', type: 'array', fields: [
            { name: 'label', type: 'text', localized: true },
            { name: 'url', type: 'text' },
          ]},
        ]},
      ],
    },

    // === MEDIA ===
    {
      slug: 'media',
      upload: {
        staticDir: 'media',
        imageSizes: [
          { name: 'thumbnail', width: 150, height: 150, position: 'centre' },
          { name: 'card', width: 400, height: 300, position: 'centre' },
          { name: 'hero', width: 1920, height: 600, position: 'centre' },
        ],
        adminThumbnail: 'thumbnail',
        mimeTypes: ['image/*', 'application/pdf'],
      },
      fields: [
        { name: 'alt', type: 'text', required: true, localized: true },
        { name: 'caption', type: 'text', localized: true },
      ],
    },

    // === USERS ===
    {
      slug: 'users',
      auth: true,
      admin: {
        useAsTitle: 'email',
      },
      fields: [
        { name: 'name', type: 'text' },
        { name: 'role', type: 'select', options: [
          { label: 'Admin', value: 'admin' },
          { label: 'Editor', value: 'editor' },
          { label: 'Marketing', value: 'marketing' },
        ]},
      ],
    },
  ],

  localization: {
    locales: ['de-CH', 'fr-CH', 'it-CH'],
    defaultLocale: 'de-CH',
    fallback: true,
  },

  plugins: [
    s3Adapter({
      collections: { media: true },
      bucket: process.env.S3_BUCKET!,
      config: {
        endpoint: process.env.S3_ENDPOINT,
        credentials: {
          accessKeyId: process.env.S3_ACCESS_KEY!,
          secretAccessKey: process.env.S3_SECRET_KEY!,
        },
      },
    }),
  ],
});
```

### 6.3 Content Blocks

```typescript
// blocks/HeroBlock.ts
import { Block } from 'payload/types';

export const HeroBlock: Block = {
  slug: 'hero',
  labels: {
    singular: 'Hero Banner',
    plural: 'Hero Banners',
  },
  fields: [
    { name: 'image', type: 'upload', relationTo: 'media', required: true },
    { name: 'imageMobile', type: 'upload', relationTo: 'media' },
    { name: 'title', type: 'text', required: true, localized: true },
    { name: 'subtitle', type: 'text', localized: true },
    { name: 'alignment', type: 'select', options: ['left', 'center', 'right'], defaultValue: 'center' },
    { name: 'overlay', type: 'checkbox', defaultValue: true },
    { name: 'cta', type: 'group', fields: [
      { name: 'text', type: 'text', localized: true },
      { name: 'link', type: 'text' },
      { name: 'style', type: 'select', options: ['primary', 'secondary', 'outline'] },
    ]},
  ],
};

// blocks/ProductGridBlock.ts
export const ProductGridBlock: Block = {
  slug: 'productGrid',
  labels: {
    singular: 'Produkt-Grid',
    plural: 'Produkt-Grids',
  },
  fields: [
    { name: 'title', type: 'text', localized: true },
    { name: 'source', type: 'select', options: [
      { label: 'Manuelle Auswahl', value: 'manual' },
      { label: 'Kategorie', value: 'category' },
      { label: 'Bestseller', value: 'bestsellers' },
      { label: 'Neu', value: 'new' },
    ]},
    { name: 'skus', type: 'array', fields: [
      { name: 'sku', type: 'text' },
    ], admin: { condition: (_, siblingData) => siblingData?.source === 'manual' }},
    { name: 'categoryId', type: 'text',
      admin: { condition: (_, siblingData) => siblingData?.source === 'category' }},
    { name: 'limit', type: 'number', defaultValue: 4 },
    { name: 'columns', type: 'select', options: ['2', '3', '4'], defaultValue: '4' },
  ],
};

// blocks/FAQBlock.ts
export const FAQBlock: Block = {
  slug: 'faq',
  labels: {
    singular: 'FAQ Sektion',
    plural: 'FAQ Sektionen',
  },
  fields: [
    { name: 'title', type: 'text', localized: true },
    { name: 'category', type: 'relationship', relationTo: 'faq-categories' },
    { name: 'items', type: 'relationship', relationTo: 'faqs', hasMany: true },
  ],
};
```

---

## 7. Frontend Integration

### 7.1 Content Fetching

```typescript
// lib/cms/client.ts
import { getPayloadClient } from 'payload';

// Server-side: Direct DB access (schneller)
export async function getPage(slug: string, locale: string) {
  const payload = await getPayloadClient();

  const result = await payload.find({
    collection: 'pages',
    where: {
      slug: { equals: slug },
      _status: { equals: 'published' },
    },
    locale,
    depth: 2,
  });

  return result.docs[0] || null;
}

// Client-side: REST API
export async function getPageAPI(slug: string, locale: string) {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_CMS_URL}/api/pages?where[slug][equals]=${slug}&locale=${locale}`,
    { next: { revalidate: 60 } } // ISR: 60 Sekunden
  );
  const data = await res.json();
  return data.docs[0] || null;
}

// GraphQL Query
export const GET_PAGE = `
  query GetPage($slug: String!, $locale: LocaleInputType!) {
    Pages(where: { slug: { equals: $slug } }, locale: $locale, limit: 1) {
      docs {
        id
        title
        slug
        template
        content
        seo {
          metaTitle
          metaDescription
          ogImage {
            url
          }
        }
      }
    }
  }
`;
```

### 7.2 Page Rendering

```tsx
// app/[locale]/[...slug]/page.tsx
import { getPage } from '@/lib/cms/client';
import { notFound } from 'next/navigation';
import { BlockRenderer } from '@/components/cms/BlockRenderer';

interface Props {
  params: {
    locale: string;
    slug: string[];
  };
}

export async function generateMetadata({ params }: Props) {
  const page = await getPage(params.slug.join('/'), params.locale);
  if (!page) return {};

  return {
    title: page.seo?.metaTitle || page.title,
    description: page.seo?.metaDescription,
    openGraph: {
      images: page.seo?.ogImage ? [page.seo.ogImage.url] : [],
    },
  };
}

export default async function CMSPage({ params }: Props) {
  const page = await getPage(params.slug.join('/'), params.locale);

  if (!page) {
    notFound();
  }

  return (
    <main>
      {page.content?.map((block, index) => (
        <BlockRenderer key={index} block={block} />
      ))}
    </main>
  );
}

// components/cms/BlockRenderer.tsx
import { HeroBlock } from './blocks/HeroBlock';
import { TextBlock } from './blocks/TextBlock';
import { ProductGridBlock } from './blocks/ProductGridBlock';
import { FAQBlock } from './blocks/FAQBlock';

const blockComponents = {
  hero: HeroBlock,
  text: TextBlock,
  productGrid: ProductGridBlock,
  faq: FAQBlock,
  // ... weitere Blocks
};

export function BlockRenderer({ block }: { block: any }) {
  const Component = blockComponents[block.blockType];

  if (!Component) {
    console.warn(`Unknown block type: ${block.blockType}`);
    return null;
  }

  return <Component {...block} />;
}

// components/cms/blocks/HeroBlock.tsx
import Image from 'next/image';
import Link from 'next/link';

interface HeroBlockProps {
  image: { url: string; alt: string };
  imageMobile?: { url: string; alt: string };
  title: string;
  subtitle?: string;
  alignment: 'left' | 'center' | 'right';
  overlay: boolean;
  cta?: {
    text: string;
    link: string;
    style: 'primary' | 'secondary' | 'outline';
  };
}

export function HeroBlock({ image, imageMobile, title, subtitle, alignment, overlay, cta }: HeroBlockProps) {
  return (
    <section className="relative h-[60vh] min-h-[400px]">
      <picture>
        {imageMobile && (
          <source media="(max-width: 768px)" srcSet={imageMobile.url} />
        )}
        <Image
          src={image.url}
          alt={image.alt}
          fill
          className="object-cover"
          priority
        />
      </picture>

      {overlay && <div className="absolute inset-0 bg-black/40" />}

      <div className={`relative z-10 flex h-full items-center ${
        alignment === 'left' ? 'justify-start' :
        alignment === 'right' ? 'justify-end' : 'justify-center'
      }`}>
        <div className={`max-w-2xl px-6 text-white ${
          alignment === 'center' ? 'text-center' : ''
        }`}>
          <h1 className="text-4xl font-bold md:text-5xl">{title}</h1>
          {subtitle && <p className="mt-4 text-xl">{subtitle}</p>}
          {cta?.text && (
            <Link
              href={cta.link}
              className={`mt-6 inline-block rounded px-6 py-3 font-semibold ${
                cta.style === 'primary' ? 'bg-primary text-white' :
                cta.style === 'secondary' ? 'bg-white text-black' :
                'border-2 border-white text-white'
              }`}
            >
              {cta.text}
            </Link>
          )}
        </div>
      </div>
    </section>
  );
}
```

### 7.3 Preview Mode

```typescript
// app/api/preview/route.ts
import { draftMode } from 'next/headers';
import { redirect } from 'next/navigation';

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const secret = searchParams.get('secret');
  const slug = searchParams.get('slug');

  // Validate secret
  if (secret !== process.env.PREVIEW_SECRET) {
    return new Response('Invalid token', { status: 401 });
  }

  // Enable draft mode
  draftMode().enable();

  // Redirect to the page
  redirect(slug || '/');
}

// In Page-Komponente
import { draftMode } from 'next/headers';

export default async function CMSPage({ params }: Props) {
  const { isEnabled: isDraft } = draftMode();

  const page = await getPage(params.slug.join('/'), params.locale, {
    draft: isDraft, // Auch Drafts laden
  });

  // ...
}
```

---

## 8. Cache & Performance

### 8.1 Caching-Strategie

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          CACHING STRATEGIE                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Payload CMS                    Next.js                    CDN              │
│  ──────────                     ───────                    ───              │
│                                                                              │
│  ┌─────────────┐  Webhook   ┌─────────────┐           ┌─────────────┐      │
│  │             │───────────▶│             │──────────▶│             │      │
│  │  Publish    │            │  Revalidate │           │  Purge      │      │
│  │  Content    │            │  Path       │           │  Cache      │      │
│  │             │            │             │           │             │      │
│  └─────────────┘            └─────────────┘           └─────────────┘      │
│                                                                              │
│  Caching Levels:                                                            │
│                                                                              │
│  1. PostgreSQL Query Cache (Payload)                                        │
│  2. Redis (API Response Cache)                                              │
│  3. Next.js ISR (Incremental Static Regeneration)                          │
│  4. CDN Edge Cache (Cloudflare/Vercel)                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Webhook für Cache Invalidation

```typescript
// Payload: Webhook bei Publish
// payload.config.ts
{
  collections: [
    {
      slug: 'pages',
      hooks: {
        afterChange: [
          async ({ doc, operation }) => {
            if (operation === 'update' && doc._status === 'published') {
              // Next.js Revalidation
              await fetch(`${process.env.FRONTEND_URL}/api/revalidate`, {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json',
                  'x-revalidate-token': process.env.REVALIDATE_TOKEN!,
                },
                body: JSON.stringify({
                  type: 'page',
                  slug: doc.slug,
                }),
              });
            }
          },
        ],
      },
    },
  ],
}

// Next.js: Revalidation Endpoint
// app/api/revalidate/route.ts
import { revalidatePath, revalidateTag } from 'next/cache';

export async function POST(request: Request) {
  const token = request.headers.get('x-revalidate-token');

  if (token !== process.env.REVALIDATE_TOKEN) {
    return new Response('Unauthorized', { status: 401 });
  }

  const { type, slug } = await request.json();

  if (type === 'page') {
    // Revalidate specific page
    revalidatePath(`/de-CH/${slug}`);
    revalidatePath(`/fr-CH/${slug}`);
    revalidatePath(`/it-CH/${slug}`);
  } else if (type === 'navigation') {
    // Revalidate all pages (navigation changed)
    revalidateTag('navigation');
  } else if (type === 'banner') {
    // Revalidate homepage
    revalidatePath('/de-CH');
    revalidatePath('/fr-CH');
    revalidatePath('/it-CH');
  }

  return new Response('OK');
}
```

---

## 9. Migration von Statamic + Typo3

### 9.1 Migrations-Plan

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          MIGRATIONS-PHASEN                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Phase 1: Setup (2 Wochen)                                                  │
│  ──────────────────────────                                                 │
│  • Payload CMS aufsetzen (Kubernetes)                                       │
│  • PostgreSQL Schema migrieren                                              │
│  • S3 Bucket für Media konfigurieren                                        │
│  • Content-Types definieren                                                 │
│  • Admin-User anlegen                                                       │
│                                                                              │
│  Phase 2: Content-Migration (3 Wochen)                                      │
│  ─────────────────────────────────────                                      │
│  • Export-Scripts für Statamic schreiben                                    │
│  • Export-Scripts für Typo3 schreiben                                       │
│  • Media-Migration (Bilder, PDFs)                                           │
│  • Content importieren und verifizieren                                     │
│  • URL-Redirects vorbereiten                                                │
│                                                                              │
│  Phase 3: Frontend-Integration (2 Wochen)                                   │
│  ────────────────────────────────────────                                   │
│  • Block-Komponenten in Next.js erstellen                                   │
│  • Page-Rendering implementieren                                            │
│  • Navigation integrieren                                                   │
│  • Preview Mode einrichten                                                  │
│                                                                              │
│  Phase 4: Parallel-Betrieb (2 Wochen)                                       │
│  ────────────────────────────────────                                       │
│  • Beide Systeme parallel laufen                                            │
│  • Marketing testet neues System                                            │
│  • Bug-Fixing                                                               │
│                                                                              │
│  Phase 5: Cutover (1 Woche)                                                 │
│  ──────────────────────────                                                 │
│  • DNS-Umstellung für Typo3-Seiten                                          │
│  • Redirects aktivieren                                                     │
│  • Statamic deaktivieren                                                    │
│  • Monitoring                                                               │
│                                                                              │
│  Total: ~10 Wochen                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 Export-Scripts

```typescript
// scripts/migrate-statamic.ts
import fs from 'fs';
import path from 'path';
import yaml from 'js-yaml';

interface StatamicEntry {
  id: string;
  title: string;
  content: string;
  // ...
}

async function migrateStatamicContent() {
  const contentDir = '/path/to/statamic/content/collections/pages';
  const files = fs.readdirSync(contentDir);

  const pages = files
    .filter(f => f.endsWith('.md'))
    .map(file => {
      const content = fs.readFileSync(path.join(contentDir, file), 'utf-8');
      const [_, frontmatter, body] = content.split('---');
      const meta = yaml.load(frontmatter) as StatamicEntry;

      return {
        title: meta.title,
        slug: file.replace('.md', ''),
        content: convertToPayloadBlocks(body, meta),
        // ...
      };
    });

  // Import to Payload
  for (const page of pages) {
    await payload.create({
      collection: 'pages',
      data: page,
    });
    console.log(`Migrated: ${page.slug}`);
  }
}

function convertToPayloadBlocks(markdown: string, meta: any): any[] {
  // Convert Statamic Bard/Markdown to Payload Blocks
  // ...
}
```

---

## 10. Zusammenfassung

### 10.1 Vorteile des neuen Systems

| Bereich | Verbesserung |
|---------|--------------|
| **Wartung** | Ein System statt zwei (Statamic + Typo3) |
| **Deployments** | CMS unabhängig vom Shop |
| **Datenbank** | PostgreSQL statt File-based |
| **Skalierung** | Horizontal skalierbar, kein Bulk-Problem |
| **Multi-Channel** | API-first für App, Displays, etc. |
| **SEO** | Eine Domain, eine Sitemap |
| **DX** | TypeScript, Code-first Schema |
| **UX** | Einheitliches Design über alle Seiten |

### 10.2 Empfohlene Lösung

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  EMPFEHLUNG                                                                  │
│  ══════════                                                                  │
│                                                                              │
│  CMS:        Payload CMS (self-hosted)                                      │
│  Deployment: Standalone in Kubernetes                                        │
│  Datenbank:  PostgreSQL (shared mit Shop)                                   │
│  Media:      Azure Blob / S3                                                │
│  Sprachen:   de-CH, fr-CH, it-CH                                            │
│                                                                              │
│  Ersetzt:    Statamic + Typo3                                               │
│  Integriert: Next.js Frontend (Shop + Marketing Pages)                      │
│                                                                              │
│  Aufwand:    ~10 Wochen (Setup + Migration + Integration)                   │
│  Kosten:     $0 Lizenz, nur Hosting                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.3 Nächste Schritte

1. [ ] Stakeholder-Präsentation für Marketing
2. [ ] Payload CMS PoC aufsetzen
3. [ ] Content-Audit (Statamic + Typo3)
4. [ ] Block-Library definieren
5. [ ] Migrations-Scripts entwickeln
