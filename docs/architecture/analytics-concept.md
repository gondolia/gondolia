# Webshop Analytics Konzept

> **Quelle:** Anforderungskatalog Webshop-Analysetool (Michael Pletscher, 30.09.2025)

Dieses Dokument beschreibt die Umsetzung der Business-Anforderungen für Kundenverhalten-Analyse in Webshop V3.

---

## 1. Architektur-Übersicht

Die Analytics-Architektur baut auf der **Customer Journey Events**-Infrastruktur auf und erweitert diese um Aggregation, Reporting und Export-Funktionen.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ANALYTICS ARCHITEKTUR                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                        EVENT COLLECTION                                 │ │
│  │                                                                         │ │
│  │   Frontend          Backend Services         External                   │ │
│  │   ────────          ────────────────         ────────                   │ │
│  │   • Page Views      • Order Events           • SAP Events               │ │
│  │   • Clicks          • Cart Events            • Payment Events           │ │
│  │   • Scroll          • Search Events          • Shipping Events          │ │
│  │   • Forms           • Auth Events                                       │ │
│  │   • Errors          • Error Events                                      │ │
│  │                                                                         │ │
│  └─────────────────────────────┬──────────────────────────────────────────┘ │
│                                │                                             │
│                                ▼                                             │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                         KAFKA TOPICS                                    │ │
│  │                                                                         │ │
│  │   analytics.pageviews    analytics.interactions    analytics.funnels   │ │
│  │   analytics.searches     analytics.errors          analytics.sessions  │ │
│  │   customer.journey       order.events              cart.events         │ │
│  │                                                                         │ │
│  └─────────────────────────────┬──────────────────────────────────────────┘ │
│                                │                                             │
│                                ▼                                             │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                      ANALYTICS SERVICE                                  │ │
│  │                                                                         │ │
│  │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐                  │ │
│  │   │  Collector  │   │ Aggregator  │   │  Reporter   │                  │ │
│  │   │             │   │             │   │             │                  │ │
│  │   │ Raw Events  │──▶│ Metrics     │──▶│ Dashboards  │                  │ │
│  │   │ Validation  │   │ Rollups     │   │ Exports     │                  │ │
│  │   │ Enrichment  │   │ Funnels     │   │ API         │                  │ │
│  │   └─────────────┘   └─────────────┘   └─────────────┘                  │ │
│  │                                                                         │ │
│  └─────────────────────────────┬──────────────────────────────────────────┘ │
│                                │                                             │
│         ┌──────────────────────┼──────────────────────┐                     │
│         ▼                      ▼                      ▼                     │
│  ┌─────────────┐       ┌─────────────┐       ┌─────────────┐               │
│  │  TimescaleDB │       │   Grafana   │       │  Export API │               │
│  │  (Time-Series)│       │  Dashboards │       │  (PowerBI)  │               │
│  └─────────────┘       └─────────────┘       └─────────────┘               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Anforderungs-Mapping

### 2.1 Abdeckung der Business-Anforderungen

| Anforderung | Umsetzung in V3 | Priorität |
|-------------|-----------------|-----------|
| **Schnittstellen-Anbindung (PowerBI)** | REST API + Webhooks für Export | Hoch |
| **Bounce Rate** | Frontend Event Tracking + Session-Analyse | Hoch |
| **Exit-Rate pro Seite** | Page Exit Events | Hoch |
| **Warenkorb-Rate** | Funnel Event: `cart_add` | Hoch |
| **Checkout-Start-Rate** | Funnel Event: `checkout_start` | Hoch |
| **Abbruchrate im Checkout** | Step-Events: `checkout_step_*` | Hoch |
| **Click-Through-Rate** | Interaction Events | Mittel |
| **Scrolltiefe / Verweildauer** | Frontend Tracking | Mittel |
| **Suchfunktion-Nutzung** | Search Events | Hoch |
| **Null-Treffer-Suchen** | Search Events mit `results_count=0` | Hoch |
| **Fehler & Abbruchmeldungen** | Error Events (bereits in Journey) | Hoch |
| **Page Speed / Ladezeiten** | Performance Events (Web Vitals) | Mittel |
| **Mobile vs. Desktop Conversion** | Device-Segmentierung | Mittel |
| **Geräte-/Browser-Abbruchmuster** | Session + Device Tracking | Mittel |
| **Wiederkehrende Besucher** | Session-Linking über Customer ID | Mittel |
| **Warenkorbabbrecher** | Abandoned Cart Detection | Hoch |
| **Customer Lifetime Value (CLV)** | Aggregation aus Order Events | Mittel |
| **Warenkorb-Verhalten** | Cart Event Analyse | Hoch |
| **Suchanfragen (Interessenprofil)** | Search Query Aggregation | Mittel |
| **Klickpfade & Verweildauer** | Session Path Analysis | Mittel |
| **Login-/Account-Daten** | Auth Events + Profile Events | Niedrig |
| **Durchschn. Bestellwert (AOV)** | Order Event Aggregation | Hoch |
| **Bestellhäufigkeit** | Order Frequency Analysis | Mittel |
| **Retourenquote** | Return Events (via SAP) | Niedrig |
| **Bevorzugte Zahlungsarten** | Payment Event Analysis | Mittel |
| **Geräte- & Browserdaten** | User-Agent Tracking | Niedrig |
| **Standort (Geo-Analyse)** | IP-basierte Geo-Location | Niedrig |
| **Login-Zeiten** | Auth Events Timestamp Analysis | Niedrig |
| **Mein Konto-Bereich Nutzung** | Account Feature Events | Niedrig |
| **Dienstleistungen (Konfigurator)** | Service Usage Events | Mittel |

---

## 3. Event-Typen

### 3.1 Frontend Events (Browser)

```typescript
// Frontend Event Types
interface AnalyticsEvent {
  // Identifikation
  event_id: string;
  event_type: string;
  timestamp: string;

  // Session
  session_id: string;
  visitor_id: string;       // Persistent Cookie
  customer_id?: string;     // Falls eingeloggt

  // Context
  tenant_id: string;
  locale: string;

  // Page
  page_url: string;
  page_path: string;
  page_title: string;
  referrer?: string;

  // Device
  device_type: 'mobile' | 'tablet' | 'desktop';
  browser: string;
  browser_version: string;
  os: string;
  screen_width: number;
  screen_height: number;

  // Geo (Server-Side enriched)
  country_code?: string;
  region?: string;
  city?: string;

  // Event-specific Data
  properties: Record<string, any>;
}
```

### 3.2 Event-Katalog

```typescript
// === PAGE EVENTS ===

// Page View
{
  event_type: 'page_view',
  properties: {
    page_type: 'product' | 'category' | 'cart' | 'checkout' | 'home' | 'search' | 'account',
    category_id?: string,
    product_id?: string,
    search_query?: string,
  }
}

// Page Exit (via beforeunload)
{
  event_type: 'page_exit',
  properties: {
    time_on_page_ms: number,
    scroll_depth_percent: number,
    exit_type: 'navigation' | 'close' | 'back',
  }
}

// === INTERACTION EVENTS ===

// Click
{
  event_type: 'click',
  properties: {
    element_type: 'button' | 'link' | 'product_card' | 'banner' | 'menu',
    element_id?: string,
    element_text?: string,
    target_url?: string,
  }
}

// Scroll
{
  event_type: 'scroll',
  properties: {
    depth_percent: 25 | 50 | 75 | 100,
    page_height: number,
  }
}

// === SEARCH EVENTS ===

// Search Performed
{
  event_type: 'search',
  properties: {
    query: string,
    results_count: number,
    has_results: boolean,
    filters_applied: string[],
    sort_by?: string,
  }
}

// Search Result Click
{
  event_type: 'search_result_click',
  properties: {
    query: string,
    product_id: string,
    position: number,
  }
}

// === PRODUCT EVENTS ===

// Product View
{
  event_type: 'product_view',
  properties: {
    product_id: string,
    sku: string,
    name: string,
    category_id: string,
    category_name: string,
    price: number,
    currency: string,
    source: 'search' | 'category' | 'recommendation' | 'direct',
  }
}

// Product Impression (in List)
{
  event_type: 'product_impression',
  properties: {
    product_id: string,
    list_name: 'search_results' | 'category_list' | 'recommendations' | 'recently_viewed',
    position: number,
  }
}

// === CART EVENTS ===

// Add to Cart
{
  event_type: 'cart_add',
  properties: {
    product_id: string,
    sku: string,
    name: string,
    quantity: number,
    unit_price: number,
    currency: string,
    source: 'product_page' | 'quick_add' | 'search' | 'wishlist',
  }
}

// Remove from Cart
{
  event_type: 'cart_remove',
  properties: {
    product_id: string,
    quantity: number,
    reason?: 'user_action' | 'out_of_stock' | 'price_change',
  }
}

// Update Cart Quantity
{
  event_type: 'cart_update',
  properties: {
    product_id: string,
    old_quantity: number,
    new_quantity: number,
  }
}

// View Cart
{
  event_type: 'cart_view',
  properties: {
    item_count: number,
    cart_value: number,
    currency: string,
  }
}

// === CHECKOUT FUNNEL EVENTS ===

// Checkout Started
{
  event_type: 'checkout_start',
  properties: {
    cart_value: number,
    item_count: number,
    currency: string,
  }
}

// Checkout Step
{
  event_type: 'checkout_step',
  properties: {
    step: 'login' | 'shipping_address' | 'billing_address' | 'shipping_method' | 'payment_method' | 'review',
    step_number: number,
  }
}

// Checkout Step Completed
{
  event_type: 'checkout_step_complete',
  properties: {
    step: string,
    step_number: number,
    time_on_step_ms: number,
  }
}

// Checkout Abandoned
{
  event_type: 'checkout_abandon',
  properties: {
    last_step: string,
    cart_value: number,
    item_count: number,
    time_in_checkout_ms: number,
    reason?: 'navigation' | 'close' | 'error',
  }
}

// === ORDER EVENTS ===

// Order Placed
{
  event_type: 'order_placed',
  properties: {
    order_id: string,
    order_number: string,
    total: number,
    subtotal: number,
    shipping: number,
    tax: number,
    currency: string,
    item_count: number,
    payment_method: string,
    shipping_method: string,
    is_first_order: boolean,
  }
}

// === AUTH EVENTS ===

// Login
{
  event_type: 'login',
  properties: {
    method: 'password' | 'sso' | 'magic_link',
    success: boolean,
    error_code?: string,
  }
}

// Register
{
  event_type: 'register',
  properties: {
    method: 'form' | 'sso',
    success: boolean,
  }
}

// Logout
{
  event_type: 'logout',
  properties: {
    trigger: 'user' | 'session_timeout' | 'forced',
  }
}

// === ERROR EVENTS ===

// JavaScript Error
{
  event_type: 'js_error',
  properties: {
    message: string,
    stack?: string,
    filename: string,
    line: number,
    column: number,
  }
}

// API Error
{
  event_type: 'api_error',
  properties: {
    endpoint: string,
    status_code: number,
    error_code?: string,
    trace_id?: string,
  }
}

// === PERFORMANCE EVENTS ===

// Web Vitals
{
  event_type: 'web_vitals',
  properties: {
    lcp: number,  // Largest Contentful Paint
    fid: number,  // First Input Delay
    cls: number,  // Cumulative Layout Shift
    ttfb: number, // Time to First Byte
    fcp: number,  // First Contentful Paint
  }
}

// === SERVICE USAGE EVENTS ===

// Configurator
{
  event_type: 'configurator_use',
  properties: {
    configurator_type: 'hit' | 'worktop' | 'panel',
    action: 'start' | 'step' | 'complete' | 'abandon',
    step?: string,
    result_sku?: string,
  }
}

// Sample Request
{
  event_type: 'sample_request',
  properties: {
    product_id: string,
    sample_type: string,
  }
}
```

---

## 4. Frontend Tracking Implementation

### 4.1 Analytics SDK

```typescript
// lib/analytics/tracker.ts

interface TrackerConfig {
  endpoint: string;
  tenantId: string;
  batchSize: number;
  flushInterval: number;
  debug: boolean;
}

class AnalyticsTracker {
  private config: TrackerConfig;
  private queue: AnalyticsEvent[] = [];
  private sessionId: string;
  private visitorId: string;
  private customerId?: string;

  constructor(config: TrackerConfig) {
    this.config = config;
    this.sessionId = this.getOrCreateSessionId();
    this.visitorId = this.getOrCreateVisitorId();

    // Flush on interval
    setInterval(() => this.flush(), config.flushInterval);

    // Flush before page unload
    window.addEventListener('beforeunload', () => this.flush());

    // Auto-track page views
    this.trackPageView();

    // Auto-track Web Vitals
    this.trackWebVitals();

    // Auto-track errors
    this.setupErrorTracking();
  }

  // Identify logged-in user
  identify(customerId: string) {
    this.customerId = customerId;
  }

  // Track generic event
  track(eventType: string, properties: Record<string, any> = {}) {
    const event: AnalyticsEvent = {
      event_id: crypto.randomUUID(),
      event_type: eventType,
      timestamp: new Date().toISOString(),
      session_id: this.sessionId,
      visitor_id: this.visitorId,
      customer_id: this.customerId,
      tenant_id: this.config.tenantId,
      locale: document.documentElement.lang || 'de-CH',
      page_url: window.location.href,
      page_path: window.location.pathname,
      page_title: document.title,
      referrer: document.referrer,
      device_type: this.getDeviceType(),
      browser: this.getBrowser(),
      browser_version: this.getBrowserVersion(),
      os: this.getOS(),
      screen_width: window.screen.width,
      screen_height: window.screen.height,
      properties,
    };

    this.queue.push(event);

    if (this.queue.length >= this.config.batchSize) {
      this.flush();
    }
  }

  // Convenience methods
  trackPageView(properties?: Record<string, any>) {
    this.track('page_view', {
      page_type: this.detectPageType(),
      ...properties,
    });
  }

  trackClick(element: HTMLElement, properties?: Record<string, any>) {
    this.track('click', {
      element_type: this.getElementType(element),
      element_id: element.id,
      element_text: element.textContent?.slice(0, 100),
      ...properties,
    });
  }

  trackSearch(query: string, resultsCount: number, filters?: string[]) {
    this.track('search', {
      query,
      results_count: resultsCount,
      has_results: resultsCount > 0,
      filters_applied: filters || [],
    });
  }

  trackAddToCart(product: ProductInfo, quantity: number, source: string) {
    this.track('cart_add', {
      product_id: product.id,
      sku: product.sku,
      name: product.name,
      quantity,
      unit_price: product.price,
      currency: product.currency,
      source,
    });
  }

  trackCheckoutStep(step: string, stepNumber: number) {
    this.track('checkout_step', { step, step_number: stepNumber });
  }

  trackOrder(order: OrderInfo) {
    this.track('order_placed', {
      order_id: order.id,
      order_number: order.number,
      total: order.total,
      subtotal: order.subtotal,
      shipping: order.shipping,
      tax: order.tax,
      currency: order.currency,
      item_count: order.items.length,
      payment_method: order.paymentMethod,
      shipping_method: order.shippingMethod,
      is_first_order: order.isFirstOrder,
    });
  }

  // Scroll tracking
  setupScrollTracking() {
    let maxScroll = 0;
    const thresholds = [25, 50, 75, 100];
    const tracked = new Set<number>();

    window.addEventListener('scroll', () => {
      const scrollPercent = Math.round(
        (window.scrollY / (document.body.scrollHeight - window.innerHeight)) * 100
      );

      if (scrollPercent > maxScroll) {
        maxScroll = scrollPercent;

        for (const threshold of thresholds) {
          if (maxScroll >= threshold && !tracked.has(threshold)) {
            tracked.add(threshold);
            this.track('scroll', {
              depth_percent: threshold,
              page_height: document.body.scrollHeight,
            });
          }
        }
      }
    });
  }

  // Web Vitals tracking
  private trackWebVitals() {
    if (typeof window !== 'undefined' && 'PerformanceObserver' in window) {
      // LCP
      new PerformanceObserver((list) => {
        const entries = list.getEntries();
        const lastEntry = entries[entries.length - 1];
        this.track('web_vitals', { lcp: lastEntry.startTime });
      }).observe({ entryTypes: ['largest-contentful-paint'] });

      // FID
      new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry: any) => {
          this.track('web_vitals', { fid: entry.processingStart - entry.startTime });
        });
      }).observe({ entryTypes: ['first-input'] });

      // CLS
      let clsValue = 0;
      new PerformanceObserver((list) => {
        for (const entry of list.getEntries() as any[]) {
          if (!entry.hadRecentInput) {
            clsValue += entry.value;
          }
        }
        this.track('web_vitals', { cls: clsValue });
      }).observe({ entryTypes: ['layout-shift'] });
    }
  }

  // Error tracking
  private setupErrorTracking() {
    window.addEventListener('error', (event) => {
      this.track('js_error', {
        message: event.message,
        filename: event.filename,
        line: event.lineno,
        column: event.colno,
      });
    });

    window.addEventListener('unhandledrejection', (event) => {
      this.track('js_error', {
        message: event.reason?.message || 'Unhandled Promise Rejection',
        stack: event.reason?.stack,
      });
    });
  }

  // Send events to backend
  private async flush() {
    if (this.queue.length === 0) return;

    const events = [...this.queue];
    this.queue = [];

    try {
      await fetch(this.config.endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ events }),
        keepalive: true, // Wichtig für beforeunload
      });
    } catch (error) {
      // Re-queue failed events
      this.queue.unshift(...events);
      console.error('Analytics flush failed:', error);
    }
  }

  // Helper methods
  private getOrCreateSessionId(): string {
    let sessionId = sessionStorage.getItem('analytics_session_id');
    if (!sessionId) {
      sessionId = crypto.randomUUID();
      sessionStorage.setItem('analytics_session_id', sessionId);
    }
    return sessionId;
  }

  private getOrCreateVisitorId(): string {
    let visitorId = localStorage.getItem('analytics_visitor_id');
    if (!visitorId) {
      visitorId = crypto.randomUUID();
      localStorage.setItem('analytics_visitor_id', visitorId);
    }
    return visitorId;
  }

  private getDeviceType(): 'mobile' | 'tablet' | 'desktop' {
    const width = window.innerWidth;
    if (width < 768) return 'mobile';
    if (width < 1024) return 'tablet';
    return 'desktop';
  }

  private detectPageType(): string {
    const path = window.location.pathname;
    if (path === '/' || path === '') return 'home';
    if (path.includes('/product/')) return 'product';
    if (path.includes('/category/')) return 'category';
    if (path.includes('/cart')) return 'cart';
    if (path.includes('/checkout')) return 'checkout';
    if (path.includes('/search')) return 'search';
    if (path.includes('/account')) return 'account';
    return 'other';
  }

  // ... weitere Helper-Methoden
}

export const analytics = new AnalyticsTracker({
  endpoint: process.env.NEXT_PUBLIC_ANALYTICS_ENDPOINT!,
  tenantId: process.env.NEXT_PUBLIC_TENANT_ID!,
  batchSize: 10,
  flushInterval: 5000,
  debug: process.env.NODE_ENV === 'development',
});
```

### 4.2 React Integration

```tsx
// hooks/useAnalytics.ts
import { analytics } from '@/lib/analytics/tracker';
import { useEffect } from 'react';
import { usePathname, useSearchParams } from 'next/navigation';

export function usePageTracking() {
  const pathname = usePathname();
  const searchParams = useSearchParams();

  useEffect(() => {
    analytics.trackPageView();
  }, [pathname, searchParams]);
}

export function useScrollTracking() {
  useEffect(() => {
    analytics.setupScrollTracking();
  }, []);
}

// components/AnalyticsProvider.tsx
export function AnalyticsProvider({ children, customerId }: {
  children: React.ReactNode;
  customerId?: string;
}) {
  usePageTracking();
  useScrollTracking();

  useEffect(() => {
    if (customerId) {
      analytics.identify(customerId);
    }
  }, [customerId]);

  return <>{children}</>;
}
```

---

## 5. Backend Analytics Service

### 5.1 Service-Struktur

```
services/analytics/
├── cmd/
│   └── main.go
├── internal/
│   ├── collector/           # Event Collection
│   │   ├── handler.go       # HTTP Endpoint für Events
│   │   ├── validator.go     # Event Validation
│   │   ├── enricher.go      # Geo-IP, User-Agent Parsing
│   │   └── publisher.go     # Kafka Publisher
│   │
│   ├── aggregator/          # Metrics Aggregation
│   │   ├── consumer.go      # Kafka Consumer
│   │   ├── funnel.go        # Conversion Funnel
│   │   ├── session.go       # Session Analysis
│   │   ├── rollup.go        # Time-based Rollups
│   │   └── scheduler.go     # Aggregation Jobs
│   │
│   ├── reporter/            # Reporting & Export
│   │   ├── handler.go       # REST API
│   │   ├── dashboards.go    # Grafana Provisioning
│   │   ├── export.go        # PowerBI Export
│   │   └── alerts.go        # Threshold Alerts
│   │
│   ├── models/
│   │   ├── events.go
│   │   ├── metrics.go
│   │   └── funnels.go
│   │
│   └── repository/
│       ├── timescale.go     # TimescaleDB
│       └── redis.go         # Real-time Counters
│
├── migrations/
│   ├── 000001_events.up.sql
│   ├── 000002_metrics.up.sql
│   └── 000003_funnels.up.sql
│
└── api/
    └── proto/
        └── analytics.proto
```

### 5.2 Datenbank-Schema (TimescaleDB)

```sql
-- TimescaleDB Extension aktivieren
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- ============================================
-- RAW EVENTS (Hypertable)
-- ============================================

CREATE TABLE analytics_events (
    id UUID DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,

    -- Event Identifikation
    event_type VARCHAR(50) NOT NULL,
    event_id UUID NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,

    -- Session & User
    session_id UUID NOT NULL,
    visitor_id UUID NOT NULL,
    customer_id UUID,
    company_id UUID,

    -- Page Context
    page_url TEXT,
    page_path VARCHAR(500),
    page_type VARCHAR(50),
    referrer TEXT,

    -- Device Info
    device_type VARCHAR(20),
    browser VARCHAR(50),
    browser_version VARCHAR(20),
    os VARCHAR(50),
    screen_width INT,
    screen_height INT,

    -- Geo (enriched)
    country_code CHAR(2),
    region VARCHAR(100),
    city VARCHAR(100),

    -- Event Properties (flexible)
    properties JSONB,

    -- Metadata
    received_at TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (tenant_id, timestamp, id)
);

-- Convert to Hypertable (partitioned by time)
SELECT create_hypertable('analytics_events', 'timestamp',
    partitioning_column => 'tenant_id',
    number_partitions => 4
);

-- Indexes
CREATE INDEX idx_events_session ON analytics_events (session_id, timestamp DESC);
CREATE INDEX idx_events_visitor ON analytics_events (visitor_id, timestamp DESC);
CREATE INDEX idx_events_customer ON analytics_events (customer_id, timestamp DESC) WHERE customer_id IS NOT NULL;
CREATE INDEX idx_events_type ON analytics_events (tenant_id, event_type, timestamp DESC);
CREATE INDEX idx_events_page_type ON analytics_events (tenant_id, page_type, timestamp DESC);

-- Retention Policy (90 Tage für Raw Events)
SELECT add_retention_policy('analytics_events', INTERVAL '90 days');

-- ============================================
-- AGGREGATED METRICS (Continuous Aggregates)
-- ============================================

-- Hourly Page Views
CREATE MATERIALIZED VIEW analytics_pageviews_hourly
WITH (timescaledb.continuous) AS
SELECT
    tenant_id,
    time_bucket('1 hour', timestamp) AS bucket,
    page_type,
    device_type,
    COUNT(*) AS views,
    COUNT(DISTINCT session_id) AS sessions,
    COUNT(DISTINCT visitor_id) AS visitors,
    COUNT(DISTINCT customer_id) AS customers
FROM analytics_events
WHERE event_type = 'page_view'
GROUP BY tenant_id, bucket, page_type, device_type;

-- Refresh Policy
SELECT add_continuous_aggregate_policy('analytics_pageviews_hourly',
    start_offset => INTERVAL '3 hours',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour'
);

-- Daily Conversion Funnel
CREATE MATERIALIZED VIEW analytics_funnel_daily
WITH (timescaledb.continuous) AS
SELECT
    tenant_id,
    time_bucket('1 day', timestamp) AS bucket,
    device_type,

    -- Funnel Steps
    COUNT(DISTINCT CASE WHEN event_type = 'page_view' THEN session_id END) AS sessions,
    COUNT(DISTINCT CASE WHEN event_type = 'product_view' THEN session_id END) AS product_views,
    COUNT(DISTINCT CASE WHEN event_type = 'cart_add' THEN session_id END) AS cart_adds,
    COUNT(DISTINCT CASE WHEN event_type = 'cart_view' THEN session_id END) AS cart_views,
    COUNT(DISTINCT CASE WHEN event_type = 'checkout_start' THEN session_id END) AS checkout_starts,
    COUNT(DISTINCT CASE WHEN event_type = 'order_placed' THEN session_id END) AS orders,

    -- Revenue
    SUM(CASE WHEN event_type = 'order_placed'
        THEN (properties->>'total')::DECIMAL END) AS revenue

FROM analytics_events
GROUP BY tenant_id, bucket, device_type;

-- Search Analytics
CREATE MATERIALIZED VIEW analytics_searches_daily
WITH (timescaledb.continuous) AS
SELECT
    tenant_id,
    time_bucket('1 day', timestamp) AS bucket,
    properties->>'query' AS query,
    COUNT(*) AS search_count,
    AVG((properties->>'results_count')::INT) AS avg_results,
    SUM(CASE WHEN (properties->>'results_count')::INT = 0 THEN 1 ELSE 0 END) AS zero_results_count,
    COUNT(DISTINCT session_id) AS unique_searchers
FROM analytics_events
WHERE event_type = 'search'
GROUP BY tenant_id, bucket, properties->>'query';

-- ============================================
-- SESSION ANALYSIS
-- ============================================

CREATE TABLE analytics_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    session_id UUID NOT NULL,
    visitor_id UUID NOT NULL,
    customer_id UUID,

    -- Timing
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    duration_seconds INT,

    -- Device
    device_type VARCHAR(20),
    browser VARCHAR(50),
    os VARCHAR(50),
    country_code CHAR(2),

    -- Behavior
    page_views INT DEFAULT 0,
    unique_pages INT DEFAULT 0,
    events_count INT DEFAULT 0,

    -- Entry/Exit
    entry_page VARCHAR(500),
    exit_page VARCHAR(500),

    -- Conversion
    is_bounce BOOLEAN DEFAULT false,  -- Nur 1 Page View
    added_to_cart BOOLEAN DEFAULT false,
    started_checkout BOOLEAN DEFAULT false,
    completed_order BOOLEAN DEFAULT false,
    order_value DECIMAL(12,2),

    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(tenant_id, session_id)
);

CREATE INDEX idx_sessions_tenant_date ON analytics_sessions (tenant_id, started_at DESC);
CREATE INDEX idx_sessions_visitor ON analytics_sessions (visitor_id, started_at DESC);
CREATE INDEX idx_sessions_bounce ON analytics_sessions (tenant_id, is_bounce, started_at DESC);

-- ============================================
-- CUSTOMER METRICS (CLV, etc.)
-- ============================================

CREATE TABLE analytics_customer_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,

    -- Lifetime Metrics
    first_seen_at TIMESTAMPTZ,
    last_seen_at TIMESTAMPTZ,
    total_sessions INT DEFAULT 0,
    total_page_views INT DEFAULT 0,

    -- Purchase Metrics
    first_order_at TIMESTAMPTZ,
    last_order_at TIMESTAMPTZ,
    total_orders INT DEFAULT 0,
    total_revenue DECIMAL(12,2) DEFAULT 0,
    avg_order_value DECIMAL(12,2),

    -- CLV (calculated)
    customer_lifetime_value DECIMAL(12,2),
    predicted_next_order_date DATE,
    churn_risk_score DECIMAL(3,2),  -- 0.00 - 1.00

    -- Preferences
    preferred_device VARCHAR(20),
    preferred_payment_method VARCHAR(50),
    preferred_categories JSONB,  -- Array of category_ids

    -- Calculated at
    calculated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(tenant_id, customer_id)
);

-- ============================================
-- ABANDONED CARTS
-- ============================================

CREATE TABLE analytics_abandoned_carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    session_id UUID NOT NULL,
    customer_id UUID,

    -- Cart Info
    cart_value DECIMAL(12,2),
    item_count INT,
    items JSONB,  -- Array of {product_id, sku, name, quantity, price}

    -- Abandonment
    abandoned_at TIMESTAMPTZ NOT NULL,
    abandoned_at_step VARCHAR(50),  -- 'cart', 'checkout_login', 'checkout_shipping', etc.

    -- Recovery
    recovery_email_sent_at TIMESTAMPTZ,
    recovered_at TIMESTAMPTZ,
    recovered_order_id UUID,

    -- Status
    status VARCHAR(20) DEFAULT 'abandoned',  -- 'abandoned', 'email_sent', 'recovered', 'expired'

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_abandoned_tenant ON analytics_abandoned_carts (tenant_id, abandoned_at DESC);
CREATE INDEX idx_abandoned_customer ON analytics_abandoned_carts (customer_id, abandoned_at DESC) WHERE customer_id IS NOT NULL;
CREATE INDEX idx_abandoned_status ON analytics_abandoned_carts (tenant_id, status);
```

### 5.3 Aggregation Jobs

```go
// internal/aggregator/funnel.go
package aggregator

import (
    "context"
    "time"
)

type FunnelAggregator struct {
    db     *timescale.Client
    redis  *redis.Client
    logger *zap.Logger
}

// Real-time Funnel Update (für Live-Dashboard)
func (a *FunnelAggregator) UpdateRealtimeFunnel(ctx context.Context, event *AnalyticsEvent) error {
    key := fmt.Sprintf("funnel:%s:%s", event.TenantID, time.Now().Format("2006-01-02"))

    switch event.EventType {
    case "page_view":
        a.redis.PFAdd(ctx, key+":sessions", event.SessionID)
    case "product_view":
        a.redis.PFAdd(ctx, key+":product_views", event.SessionID)
    case "cart_add":
        a.redis.PFAdd(ctx, key+":cart_adds", event.SessionID)
    case "checkout_start":
        a.redis.PFAdd(ctx, key+":checkout_starts", event.SessionID)
    case "order_placed":
        a.redis.PFAdd(ctx, key+":orders", event.SessionID)
        a.redis.IncrByFloat(ctx, key+":revenue", event.Properties["total"].(float64))
    }

    return nil
}

// Get Real-time Funnel
func (a *FunnelAggregator) GetRealtimeFunnel(ctx context.Context, tenantID string) (*FunnelMetrics, error) {
    key := fmt.Sprintf("funnel:%s:%s", tenantID, time.Now().Format("2006-01-02"))

    sessions, _ := a.redis.PFCount(ctx, key+":sessions").Result()
    productViews, _ := a.redis.PFCount(ctx, key+":product_views").Result()
    cartAdds, _ := a.redis.PFCount(ctx, key+":cart_adds").Result()
    checkoutStarts, _ := a.redis.PFCount(ctx, key+":checkout_starts").Result()
    orders, _ := a.redis.PFCount(ctx, key+":orders").Result()
    revenue, _ := a.redis.Get(ctx, key+":revenue").Float64()

    return &FunnelMetrics{
        Date:           time.Now().Format("2006-01-02"),
        Sessions:       sessions,
        ProductViews:   productViews,
        CartAdds:       cartAdds,
        CheckoutStarts: checkoutStarts,
        Orders:         orders,
        Revenue:        revenue,
        ConversionRate: float64(orders) / float64(sessions) * 100,
    }, nil
}

// Session Builder (aus Events)
func (a *FunnelAggregator) BuildSession(ctx context.Context, sessionID string) (*Session, error) {
    events, err := a.db.GetEventsBySession(ctx, sessionID)
    if err != nil {
        return nil, err
    }

    if len(events) == 0 {
        return nil, ErrSessionNotFound
    }

    session := &Session{
        SessionID:  sessionID,
        TenantID:   events[0].TenantID,
        VisitorID:  events[0].VisitorID,
        CustomerID: events[0].CustomerID,
        StartedAt:  events[0].Timestamp,
        EndedAt:    events[len(events)-1].Timestamp,
        DeviceType: events[0].DeviceType,
        Browser:    events[0].Browser,
        OS:         events[0].OS,
        Country:    events[0].CountryCode,
    }

    uniquePages := make(map[string]bool)

    for _, e := range events {
        session.EventsCount++

        switch e.EventType {
        case "page_view":
            session.PageViews++
            uniquePages[e.PagePath] = true
            if session.EntryPage == "" {
                session.EntryPage = e.PagePath
            }
            session.ExitPage = e.PagePath

        case "cart_add":
            session.AddedToCart = true

        case "checkout_start":
            session.StartedCheckout = true

        case "order_placed":
            session.CompletedOrder = true
            session.OrderValue = e.Properties["total"].(float64)
        }
    }

    session.UniquePages = len(uniquePages)
    session.IsBounce = session.PageViews == 1
    session.Duration = int(session.EndedAt.Sub(session.StartedAt).Seconds())

    return session, nil
}
```

---

## 6. Export API (PowerBI)

### 6.1 REST API Endpoints

```go
// internal/reporter/handler.go
package reporter

import (
    "github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
    // Real-time Metrics
    r.GET("/metrics/realtime", h.GetRealtimeMetrics)
    r.GET("/metrics/funnel", h.GetFunnelMetrics)

    // Historical Data (for PowerBI)
    r.GET("/export/pageviews", h.ExportPageviews)
    r.GET("/export/sessions", h.ExportSessions)
    r.GET("/export/funnels", h.ExportFunnels)
    r.GET("/export/searches", h.ExportSearches)
    r.GET("/export/customers", h.ExportCustomerMetrics)
    r.GET("/export/abandoned-carts", h.ExportAbandonedCarts)

    // Aggregated Reports
    r.GET("/reports/conversion", h.GetConversionReport)
    r.GET("/reports/bounce-rate", h.GetBounceRateReport)
    r.GET("/reports/search-analytics", h.GetSearchAnalyticsReport)
    r.GET("/reports/device-breakdown", h.GetDeviceBreakdownReport)
}

// Export Pageviews (PowerBI compatible)
// GET /api/analytics/export/pageviews?from=2025-01-01&to=2025-01-31&granularity=daily
func (h *Handler) ExportPageviews(c *gin.Context) {
    tenantID := middleware.GetTenantID(c)
    from := c.Query("from")
    to := c.Query("to")
    granularity := c.DefaultQuery("granularity", "daily") // hourly, daily, weekly

    data, err := h.repo.GetPageviewsExport(c.Request.Context(), tenantID, from, to, granularity)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // Support CSV and JSON
    format := c.DefaultQuery("format", "json")
    if format == "csv" {
        c.Header("Content-Type", "text/csv")
        c.Header("Content-Disposition", "attachment; filename=pageviews.csv")
        h.writeCSV(c.Writer, data)
        return
    }

    c.JSON(200, gin.H{
        "data": data,
        "meta": gin.H{
            "from":        from,
            "to":          to,
            "granularity": granularity,
            "total_rows":  len(data),
        },
    })
}

// Response Structure für PowerBI
type PageviewsExport struct {
    Date        string  `json:"date" csv:"Datum"`
    PageType    string  `json:"page_type" csv:"Seitentyp"`
    DeviceType  string  `json:"device_type" csv:"Gerät"`
    Views       int64   `json:"views" csv:"Seitenaufrufe"`
    Sessions    int64   `json:"sessions" csv:"Sessions"`
    Visitors    int64   `json:"visitors" csv:"Besucher"`
    Customers   int64   `json:"customers" csv:"Kunden"`
    BounceRate  float64 `json:"bounce_rate" csv:"Absprungrate"`
    AvgDuration float64 `json:"avg_duration" csv:"Ø Verweildauer"`
}

type FunnelExport struct {
    Date              string  `json:"date" csv:"Datum"`
    DeviceType        string  `json:"device_type" csv:"Gerät"`
    Sessions          int64   `json:"sessions" csv:"Sessions"`
    ProductViews      int64   `json:"product_views" csv:"Produktansichten"`
    CartAdds          int64   `json:"cart_adds" csv:"Warenkorb"`
    CheckoutStarts    int64   `json:"checkout_starts" csv:"Checkout gestartet"`
    Orders            int64   `json:"orders" csv:"Bestellungen"`
    Revenue           float64 `json:"revenue" csv:"Umsatz"`
    ConversionRate    float64 `json:"conversion_rate" csv:"Conversion Rate %"`
    CartAddRate       float64 `json:"cart_add_rate" csv:"Warenkorb-Rate %"`
    CheckoutRate      float64 `json:"checkout_rate" csv:"Checkout-Rate %"`
}
```

### 6.2 PowerBI Integration

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         POWERBI INTEGRATION                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Option 1: REST API (Empfohlen)                                             │
│  ──────────────────────────────                                             │
│  PowerBI Desktop → Get Data → Web → API URL                                │
│                                                                              │
│  URL: https://api.shop.ch/analytics/export/funnels                          │
│       ?from={StartDate}&to={EndDate}&format=json                            │
│                                                                              │
│  • Parametrisierbar mit PowerBI Parameters                                  │
│  • OAuth2 Authentication                                                     │
│  • Scheduled Refresh möglich                                                 │
│                                                                              │
│  Option 2: Direct DB Connection                                              │
│  ──────────────────────────────                                             │
│  PowerBI Gateway → PostgreSQL → TimescaleDB                                 │
│                                                                              │
│  • Read-Replica verwenden                                                    │
│  • Views für Reporting erstellen                                             │
│  • Nur für große Datenmengen                                                 │
│                                                                              │
│  Option 3: Export to Data Lake                                               │
│  ─────────────────────────────                                              │
│  Scheduled Job → Parquet Files → Azure Blob → PowerBI                       │
│                                                                              │
│  • Für historische Analysen                                                  │
│  • Große Datenmengen                                                         │
│  • ML/AI Integrationen                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Grafana Dashboards

### 7.1 Dashboard Provisioning

```yaml
# infrastructure/kubernetes/base/grafana/dashboards/analytics-overview.json
{
  "title": "Webshop Analytics Overview",
  "panels": [
    {
      "title": "Conversion Funnel (Live)",
      "type": "barchart",
      "gridPos": { "x": 0, "y": 0, "w": 12, "h": 8 },
      "targets": [
        {
          "datasource": "TimescaleDB",
          "rawSql": "SELECT ... FROM analytics_funnel_daily WHERE bucket >= NOW() - INTERVAL '7 days'"
        }
      ]
    },
    {
      "title": "Bounce Rate by Page Type",
      "type": "piechart",
      "gridPos": { "x": 12, "y": 0, "w": 6, "h": 8 }
    },
    {
      "title": "Sessions by Device",
      "type": "stat",
      "gridPos": { "x": 18, "y": 0, "w": 6, "h": 8 }
    },
    {
      "title": "Zero-Result Searches",
      "type": "table",
      "gridPos": { "x": 0, "y": 8, "w": 12, "h": 6 }
    },
    {
      "title": "Abandoned Carts",
      "type": "timeseries",
      "gridPos": { "x": 12, "y": 8, "w": 12, "h": 6 }
    }
  ]
}
```

### 7.2 Alert Rules

```yaml
# Alerting für kritische Metriken
groups:
  - name: analytics-alerts
    rules:
      - alert: HighBounceRate
        expr: analytics_bounce_rate > 0.7
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "Bounce Rate über 70%"

      - alert: ConversionDrop
        expr: (analytics_conversion_rate - analytics_conversion_rate offset 1d) / analytics_conversion_rate offset 1d < -0.3
        for: 2h
        labels:
          severity: critical
        annotations:
          summary: "Conversion Rate um 30% gesunken"

      - alert: HighCheckoutAbandonment
        expr: analytics_checkout_abandonment_rate > 0.8
        for: 30m
        labels:
          severity: warning
        annotations:
          summary: "Checkout-Abbruchrate über 80%"
```

---

## 8. Datenschutz (DSGVO)

### 8.1 Anonymisierung

```go
// Visitor ID ist pseudonym, nicht personenbezogen
// Customer ID nur bei Login verknüpft

// IP-Adressen werden NICHT gespeichert
// Nur Geo-Location (Land/Region) aus IP abgeleitet

// Retention Policies
const (
    RawEventsRetention = 90 * 24 * time.Hour  // 90 Tage
    AggregatedRetention = 2 * 365 * 24 * time.Hour // 2 Jahre
    CustomerMetricsRetention = 3 * 365 * 24 * time.Hour // 3 Jahre (nach letzter Aktivität)
)
```

### 8.2 Consent Management

```typescript
// Frontend: Nur tracken wenn Consent gegeben
if (hasAnalyticsConsent()) {
  analytics.track('page_view', ...);
}

// Cookie-Banner Integration
function initAnalytics(consent: ConsentSettings) {
  if (consent.analytics) {
    analytics.enable();
  } else {
    analytics.disable();
  }
}
```

---

## 9. Umsetzungs-Roadmap

### Phase 1: Foundation (MVP)

| Feature | Priorität | Aufwand |
|---------|-----------|---------|
| Frontend SDK (Page Views, Clicks) | Hoch | M |
| Backend Collector Endpoint | Hoch | M |
| Kafka Topics Setup | Hoch | S |
| TimescaleDB Schema | Hoch | M |
| Basic Conversion Funnel | Hoch | M |
| Grafana Dashboard (Basic) | Hoch | S |

### Phase 2: Core Analytics

| Feature | Priorität | Aufwand |
|---------|-----------|---------|
| Search Analytics | Hoch | M |
| Cart Event Tracking | Hoch | M |
| Checkout Funnel (per Step) | Hoch | L |
| Session Analysis | Mittel | L |
| Abandoned Cart Detection | Hoch | M |
| PowerBI Export API | Hoch | M |

### Phase 3: Advanced

| Feature | Priorität | Aufwand |
|---------|-----------|---------|
| Customer Lifetime Value | Mittel | L |
| Scroll/Engagement Tracking | Mittel | M |
| Web Vitals | Mittel | S |
| Device/Browser Analysis | Niedrig | M |
| Geo-Analysis | Niedrig | M |
| Configurator Usage | Mittel | M |

### Aufwands-Legende
- **S** = Small (1-2 Tage)
- **M** = Medium (3-5 Tage)
- **L** = Large (1-2 Wochen)

---

## 10. Zusammenfassung

### Abdeckung der Anforderungen

| Kategorie | Abdeckung | Details |
|-----------|-----------|---------|
| **Schnittstellen (PowerBI)** | ✅ 100% | REST API + CSV Export |
| **Absprungraten** | ✅ 100% | Bounce Rate, Exit Rate |
| **Conversion Funnel** | ✅ 100% | Alle Steps tracked |
| **Interaktionsraten** | ✅ 100% | CTR, Scroll, Search |
| **Technische KPIs** | ✅ 100% | Errors, Web Vitals, Device |
| **Wiederkehrer** | ✅ 100% | Session Linking |
| **On-Site Tracking** | ✅ 100% | CLV, Cart, Search Profile |
| **Transaktionsdaten** | ✅ 100% | AOV, Frequency, Payment |
| **Nutzungsdaten** | ⚠️ 80% | Account, Configurator |

### Architektur-Vorteile

- **Event-Driven**: Gleiche Infrastruktur wie Customer Journey (Support Portal)
- **Skalierbar**: TimescaleDB für Time-Series, Kafka für Throughput
- **Real-time + Historical**: Live-Dashboard + PowerBI für Analysen
- **DSGVO-konform**: Pseudonymisierung, Retention Policies, Consent

### Nächste Schritte

1. [ ] Analytics Service im Monorepo erstellen
2. [ ] Frontend SDK implementieren
3. [ ] TimescaleDB zu Kubernetes hinzufügen
4. [ ] Grafana Dashboards provisionieren
5. [ ] PowerBI Export testen
