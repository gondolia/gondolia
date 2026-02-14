# TODO: Catalog Service — Technische Verbesserungen

Stand: 2026-02-14

## 1. Backend: `product_count` im Category-Model

**Problem:** Das Backend liefert kein `product_count` Feld bei Kategorien. Das Frontend muss pro Kategorie einen separaten API-Call machen, nur um die Anzahl zu ermitteln (N+1 Problem).

**Aktueller Workaround:** Frontend ruft `GET /categories/:id/products?limit=1` pro Kategorie auf und nutzt `total` aus der Response.

**Lösung:**
- `product_count` Feld in `domain.Category` hinzufügen
- Im Repository per SQL Subquery berechnen: `SELECT c.*, (SELECT COUNT(*) FROM products WHERE category_id = c.id) AS product_count`
- Für Elternkategorien: rekursive Zählung über alle Kind-Kategorien (CTE oder materialized count)

**Impact:** Reduziert API-Calls auf der Kategorieübersicht von N+1 auf 1.

---

## 2. Backend: `include_children` Parameter für Kategorie-Produkte

**Problem:** `GET /categories/:id/products` liefert nur Produkte die *direkt* dieser Kategorie zugeordnet sind. Bei Elternkategorien (z.B. "Holz & Holzwerkstoffe") kommen 0 Produkte zurück, obwohl die Unterkategorien Produkte haben.

**Aktueller Workaround:** Frontend sammelt rekursiv alle Kind-Kategorie-IDs (`collectCategoryIds()`), macht pro Kind-Kategorie einen API-Call, merged und dedupliziert client-seitig.

**Lösung:**
- Neuer Query-Parameter: `GET /categories/:id/products?include_children=true`
- Backend sammelt Kind-Kategorie-IDs per CTE und filtert Produkte mit `WHERE category_id IN (...)`
- Serverseitige Pagination bleibt erhalten

**Impact:** Reduziert API-Calls bei Elternkategorien von N auf 1, ermöglicht echte serverseitige Pagination.

---

## 3. Frontend: Serverseitige Pagination für Kategorie-Produkte

**Problem:** Bei Elternkategorien werden aktuell *alle* Produkte geladen (`limit: 100` pro Unterkategorie) und client-seitig in 12er-Seiten geschnitten. Skaliert nicht bei grossen Katalogen.

**Abhängigkeit:** Erfordert #2 (`include_children` Parameter).

**Lösung:**
- Nach Implementierung von `include_children`: Frontend nutzt serverseitige Pagination wie bei Leaf-Kategorien
- `limit` und `offset` werden direkt an Backend durchgereicht
- Kein client-seitiges Merging/Deduplication mehr nötig

**Impact:** Konstanter Memory-Verbrauch im Browser, schnellere Ladezeiten bei grossen Kategorien.

---

## 4. Frontend: Lazy Loading / Infinite Scroll (Nice-to-have)

**Problem:** Alle Produkte einer Seite werden sofort beim Seitenaufruf geladen. Kein progressives Laden.

**Lösung (nach #3):**
- Option A: Infinite Scroll mit Intersection Observer
- Option B: "Mehr laden" Button
- In beiden Fällen: nächste Seite erst bei Bedarf vom Backend holen

**Priorität:** Niedrig — erst relevant bei >50 Produkten pro Kategorie.

---

## Prioritäten

| # | Verbesserung | Aufwand | Impact | Priorität |
|---|---|---|---|---|
| 1 | `product_count` im Backend | Klein | Hoch (N+1 eliminieren) | **Hoch** |
| 2 | `include_children` Parameter | Mittel | Hoch (echte Pagination) | **Hoch** |
| 3 | Serverseitige Pagination | Klein | Mittel (nach #2 trivial) | Mittel |
| 4 | Lazy Loading | Mittel | Niedrig (UX) | Niedrig |
