# B2B Self-Service Prinzipien

> **Kernthese:** B2B-Unternehmen haben kein "Prozessproblem" - sie haben ein "Posteingang-Problem".

---

## Das Problem: Der typische B2B-Bestellprozess

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AKTUELLER ZUSTAND (Chaos)                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Kunde ──Mail/Telefon──▶ Vertrieb ◀──Ping-Pong──▶ Betrieb/Produktion  │
│     │                        │                           │              │
│     │                        ▼                           │              │
│     │                   ┌─────────┐                      │              │
│     │                   │ Offerte │ ◀── Stunden/Tage ────┘              │
│     │                   └────┬────┘                                     │
│     │                        │                                          │
│     │                        ▼                                          │
│     │                   Bestellung                                      │
│     │                        │                                          │
│     │                        ▼                                          │
│     └─────────────────▶ SPIEL VON VORNE ◀───────────────────────────────┤
│                                                                         │
│   Symptome:                                                             │
│   • "Lautstärke gewinnt" - wer am meisten drängelt, wird bedient       │
│   • Alle beschäftigt, trotzdem langsam                                 │
│   • Durchlaufzeit hoch, Nerven runter                                  │
│   • Umsatz bleibt liegen                                               │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Die Lösung: Self-Service Portal

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      ZIEL-ZUSTAND (Self-Service)                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Kunde ──Self-Service──▶ Webshop ──Direkt──▶ SAP/ERP                  │
│     │                        │                    │                     │
│     │    ┌───────────────────┼────────────────────┼──────────┐         │
│     │    │                   │                    │          │         │
│     │    ▼                   ▼                    ▼          ▼         │
│     │  Produkt            Preis in            Bestellung   Lieferung   │
│     │  auswählen          Sekunden            direkt       tracken     │
│     │                                         im ERP                    │
│     │                                                                   │
│     └───────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Ergebnis:                                                             │
│   • Kunden bestellen selbst - ohne Vertriebskontakt                    │
│   • Preise/Offerten in Sekunden statt Stunden/Tagen                    │
│   • Bestellung geht direkt an ERP/Partner                              │
│   • Rückfragen und Ping-Pong gehen massiv runter                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Design-Prinzipien

### 1. Self-Service First

**Leitfrage:** Kann der Kunde das selbst erledigen, ohne jemanden anzurufen?

| Aktion | Ohne Self-Service | Mit Self-Service |
|--------|-------------------|------------------|
| Preis anfragen | Mail → Warten → Antwort | Sofort im Shop sichtbar |
| Verfügbarkeit prüfen | Anruf → Rückruf | Echtzeit-Bestand |
| Bestellung aufgeben | Mail → Bestätigung → Nachfragen | Click → Fertig |
| Lieferstatus | Anruf → "Ich frag nach" | Tracking im Portal |
| Rechnung finden | Mail → Warten | Download im Portal |

**Regel:** Wenn der Kunde für eine Standard-Aktion den Hörer abnehmen muss, ist das ein Bug.

---

### 2. Die 80/20-Regel für Automatisierung

**Leitfrage:** "Von 100 Bestellungen, wie oft kommt dieser Fall vor?"

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PRIORISIERUNG NACH HÄUFIGKEIT                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   80% der Fälle          │  15% der Fälle        │  5% der Fälle       │
│   ────────────────────   │  ─────────────────    │  ─────────────────  │
│                          │                        │                     │
│   AUTOMATISIEREN         │  PHASE 2               │  MANUELL LASSEN     │
│   (Sofort)               │  (Später)              │  (Nie)              │
│                          │                        │                     │
│   • Standard-Produkte    │  • Sonderanfertigungen │  • Einmalige        │
│   • Standard-Mengen      │  • Spezialpreise       │    Spezialfälle     │
│   • Standard-Lieferung   │  • Projektgeschäft     │  • Exceptions       │
│   • Bekannte Kunden      │  • Neukunden-Onboard.  │  • "Das eine Mal"   │
│                          │                        │                     │
└─────────────────────────────────────────────────────────────────────────┘
```

**Regel:** Nicht jeden Sonderfall automatisieren. Wenn etwas selten vorkommt: Phase 2. Oder 3. Oder gar nie.

---

### 3. Durchlaufzeit als KPI

**Leitfrage:** Wo bleibt der Prozess hängen?

```
Anfrage ──────▶ Offerte ──────▶ Auftrag ──────▶ Lieferung
   │              │               │               │
   ▼              ▼               ▼               ▼
┌──────┐      ┌──────┐        ┌──────┐        ┌──────┐
│ ZIEL │      │ ZIEL │        │ ZIEL │        │ ZIEL │
│  0s  │      │ <60s │        │ <5m  │        │Track │
│      │      │      │        │      │        │ able │
└──────┘      └──────┘        └──────┘        └──────┘

MESSEN:
• Zeit von Anfrage bis Offerte
• Zeit von Offerte bis Bestellung (Conversion)
• Zeit von Bestellung bis Auftragsbestätigung
• Zeit von Bestellung bis Lieferung
```

**Messwerte für Webshop V3:**

| Metrik | IST (V2/manuell) | SOLL (V3) |
|--------|------------------|-----------|
| Preis/Offerte | Stunden-Tage | < 1 Sekunde |
| Bestellung aufgeben | 10-30 Min | < 3 Min |
| Auftragsbestätigung | Stunden | Sofort |
| Lieferstatus | Anruf nötig | Self-Service |

---

### 4. Direktintegration statt Medienbrüche

**Leitfrage:** Wie viele Systeme/Menschen muss die Information durchlaufen?

```
SCHLECHT (Medienbrüche):
┌────────┐   Mail   ┌──────────┐  Excel  ┌─────────┐  Eingabe  ┌─────┐
│ Kunde  │ ───────▶ │ Vertrieb │ ──────▶ │ Betrieb │ ────────▶ │ SAP │
└────────┘          └──────────┘         └─────────┘           └─────┘
                          ▲                    │
                          └────── Rückfragen ──┘

GUT (Direktintegration):
┌────────┐          ┌─────────┐           ┌─────┐
│ Kunde  │ ───────▶ │ Webshop │ ────────▶ │ SAP │
└────────┘          └─────────┘           └─────┘
```

**Für Webshop V3:**

| Prozess | Direktintegration |
|---------|-------------------|
| Preisanfrage | Webshop → SAP (Echtzeit) |
| Bestellung | Webshop → SAP Order |
| Verfügbarkeit | Webshop → SAP Bestand |
| Lieferstatus | SAP → Webshop → Kunde |
| Rechnung | SAP → Webshop Download |

---

## Anforderungen für Webshop V3

### Must-Have (80% der Fälle)

1. **Produkt-Konfigurator**
   - Kunde klickt sich Standard-Aufträge selbst zusammen
   - Keine Freitext-Felder für Standardfälle
   - Geführte Auswahl statt offene Anfrage

2. **Echtzeit-Preise**
   - Kundenspezifische Preise aus SAP
   - Staffelpreise sofort sichtbar
   - Keine "Preis auf Anfrage" für Standardprodukte

3. **Sofort-Bestellung**
   - Warenkorb → Checkout → Fertig
   - Bestellung geht direkt an SAP
   - Auftragsbestätigung sofort (nicht nach Stunden)

4. **Self-Service Portal**
   - Bestellhistorie
   - Lieferstatus-Tracking
   - Rechnungs-Download
   - Reorder-Funktion

### Nice-to-Have (Phase 2)

1. **Angebotsanfrage für Sonderfälle**
   - Strukturiertes Formular (kein Freitext-Mail)
   - Automatische Weiterleitung an richtigen Ansprechpartner
   - Status-Tracking für den Kunden

2. **Projektgeschäft**
   - Mehrere Liefertermine
   - Teillieferungen
   - Projekt-Preise

3. **Sonderanfertigungen**
   - Konfigurator für Varianten
   - Automatische Machbarkeits-Prüfung
   - Lieferzeit-Kalkulation

---

## Anti-Patterns (was wir NICHT machen)

| Anti-Pattern | Problem | Stattdessen |
|--------------|---------|-------------|
| "Preis auf Anfrage" | Erzeugt Mail/Anruf | Echtzeit-Preis aus SAP |
| "Rufen Sie uns an" | Medienbruch | Self-Service |
| "Wir melden uns" | Wartezeit, Unsicherheit | Sofort-Bestätigung |
| Alles automatisieren | Overengineering | 80/20 Fokus |
| PDF-Formulare | Keine Integration | Online-Formular → SAP |

---

## Metriken & KPIs

### Business-Metriken

| Metrik | Beschreibung | Ziel |
|--------|--------------|------|
| **Self-Service-Quote** | % Bestellungen ohne Vertriebskontakt | > 80% |
| **Durchlaufzeit** | Anfrage → Lieferung | -50% vs. V2 |
| **Conversion Rate** | Warenkorb → Bestellung | > 70% |
| **Reorder-Rate** | Wiederkehrende Bestellungen | +20% |

### Technische Metriken

| Metrik | Beschreibung | Ziel |
|--------|--------------|------|
| **Preis-Response-Zeit** | Zeit für SAP-Preisabfrage | < 500ms |
| **Order-Processing** | Zeit für Bestellübertragung | < 5s |
| **Verfügbarkeits-Sync** | Aktualität der Bestandsdaten | < 15 Min |

---

## Zusammenfassung

> **Der Webshop soll Kunden ermöglichen zu bestellen, ohne dass jemand angerufen werden muss.**

| Prinzip | Umsetzung |
|---------|-----------|
| Self-Service First | Alles was der Kunde selbst kann, soll er selbst können |
| 80/20 Fokus | Standardfälle automatisieren, Sonderfälle bewusst manuell |
| Durchlaufzeit messen | Wo bleibt's hängen? → Dort optimieren |
| Direktintegration | Webshop ↔ SAP ohne Medienbrüche |

**Das Ziel ist nicht "ein Portal haben" - das Ziel ist: Bestellen. Fertig. Weiterarbeiten.**
