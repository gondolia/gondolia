# KI-Assistenten Konzept

## 1. Ausgangslage

### 1.1 Anforderungen

Es besteht Interesse an KI-Assistenten von:
- **Kunden** (B2B): Schnellere Produktfindung, technische Beratung
- **Interner Verkauf**: Kundenhistorie, Produktwissen, Beratungsunterstützung

### 1.2 Kritische Randbedingung

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  ⚠️  DOMAIN-WISSEN BLEIBT INTERN                                            │
│  ════════════════════════════════                                            │
│                                                                              │
│  Folgende Daten dürfen NICHT an externe KI-Dienste:                        │
│                                                                              │
│  • Produktdaten (Beschreibungen, techn. Daten)                             │
│  • Preise und Margen                                                        │
│  • Kundeninformationen                                                      │
│  • Bestellhistorie                                                          │
│  • Lieferanten-Konditionen                                                  │
│  • Internes Fachwissen                                                      │
│                                                                              │
│  Grund: Wettbewerbsvorteil + Datenschutz                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Use Cases

### 2.1 Übersicht

| # | Use Case | Zielgruppe | Priorität | Komplexität |
|---|----------|------------|-----------|-------------|
| 1 | Produktsuche (Natural Language) | Kunden | Hoch | Mittel |
| 2 | Technische Beratung | Kunden | Hoch | Mittel |
| 3 | Chat-Support | Kunden | Mittel | Mittel |
| 4 | Konfigurations-Hilfe | Kunden | Mittel | Hoch |
| 5 | Sales Assistant (intern) | Verkauf | Hoch | Mittel |

### 2.2 Use Case 1: Produktsuche (Natural Language)

**Szenario:**
```
Kunde: "Ich suche einen Bodenbelag für mein Badezimmer,
        sollte wasserfest sein und nicht zu teuer"

Assistent: "Für Badezimmer empfehle ich wasserfeste Optionen:

        1. Vinyl-Designboden 'AquaProtect' (CHF 34.90/m²)
           → 100% wasserfest, einfache Klick-Montage

        2. Feinsteinzeug 'Stone Grey' (CHF 28.50/m²)
           → Fliesen-Optik, rutschfest

        Soll ich nach einer bestimmten Optik filtern?"
```

**Technisch:**
- Semantische Suche über Produktkatalog
- Attribut-Extraktion aus Freitext (wasserfest → Attribut-Filter)
- Preis-Bewusstsein ("nicht zu teuer" → Budget-Segment)

### 2.3 Use Case 2: Technische Beratung

**Szenario:**
```
Kunde: "Kann ich das Parkett 'Eiche Rustikal' auf
        Fussbodenheizung verlegen?"

Assistent: "Ja, das Parkett 'Eiche Rustikal' ist für
        Fussbodenheizung geeignet, mit Einschränkungen:

        ✅ Warmwasser-Fussbodenheizung: Ja
        ⚠️ Max. Oberflächentemperatur: 27°C
        ⚠️ Elektrische FBH: Nicht empfohlen

        Wichtig: Verwenden Sie eine geeignete Dämmunterlage.
        → Empfehlung: 'ThermoBase Plus' (siehe Zubehör)

        Quelle: Produktdatenblatt Seite 3"
```

**Technisch:**
- RAG über Produktdatenblätter (PDFs)
- Technische Spezifikationen aus PIM
- Verknüpfung zu Zubehör-Empfehlungen

### 2.4 Use Case 3: Chat-Support

**Szenario:**
```
Kunde: "Wo ist meine Bestellung 12847?"

Assistent: "Ihre Bestellung #12847 vom 15.01.2025:

        Status: In Zustellung
        Lieferung: Heute, 14:00-18:00 Uhr
        Spediteur: Planzer (Sendungsnr. PLZ-123456)

        → Sendungsverfolgung öffnen

        Kann ich sonst noch helfen?"
```

**Technisch:**
- Integration mit Order Service
- Tracking-Informationen von Spediteuren
- Eskalation an menschlichen Support bei Problemen

### 2.5 Use Case 4: Konfigurations-Hilfe

**Szenario:**
```
Kunde: "Ich möchte mein Wohnzimmer mit Laminat auslegen.
        Der Raum ist 5.2m x 4.8m."

Assistent: "Für Ihr Wohnzimmer (24.96 m²) berechne ich:

        Benötigte Menge (inkl. 10% Verschnitt): 27.5 m²

        Passende Produkte:
        1. Swiss Krono Eiche Natur - 11 Pakete (CHF 538.-)
        2. Parador Classic - 10 Pakete (CHF 489.-)

        Zusätzlich empfohlen:
        • Trittschalldämmung: 28 m² (CHF 84.-)
        • Sockelleisten: 20 lfm (CHF 60.-)
        • Keile-Set für Dehnungsfuge (CHF 12.90)

        Soll ich einen Warenkorb erstellen?"
```

**Technisch:**
- Flächenberechnung + Verschnitt
- Paket-Rundung basierend auf Produktdaten
- Cross-Selling Logik (Zubehör)
- Warenkorb-Integration

### 2.6 Use Case 5: Sales Assistant (Intern)

**Szenario:**
```
Verkäufer: "Zeig mir die Bestellhistorie von Schreinerei Müller"

Assistent: "Schreinerei Müller AG (Kunde seit 2018):

        Letzte 12 Monate:
        • 47 Bestellungen, Gesamtwert CHF 128'450.-
        • Hauptprodukte: Multiplex (45%), Schichtstoff (30%)
        • Ø Bestellwert: CHF 2'733.-

        Auffälligkeiten:
        ⚠️ Bestellvolumen -15% vs. Vorjahr
        ⚠️ Keine Bestellung seit 23 Tagen (sonst wöchentlich)

        Offene Angebote: 2 (CHF 8'400.-)
        Offene Reklamationen: 0

        Kontakt: Hans Müller, 079 123 45 67"
```

**Technisch:**
- Zugriff auf CRM/Kundendaten
- Bestellanalyse und Trends
- Anomalie-Erkennung
- NUR intern, mit Berechtigungen

---

## 3. Technische Architektur

### 3.1 Optionen für LLM

| Option | Datenschutz | Qualität | Kosten | Aufwand |
|--------|-------------|----------|--------|---------|
| **Self-Hosted (Llama 3, Mistral)** | ✅ Lokal | ⚠️ Gut | Mittel | Hoch |
| **Private Cloud (Azure OpenAI)** | ✅ Vertrag | ✅ Sehr gut | Hoch | Mittel |
| **Anthropic API (Claude)** | ⚠️ Vertrag | ✅ Sehr gut | Mittel | Niedrig |
| **Local + Cloud Hybrid** | ✅ Flexibel | ✅ Sehr gut | Mittel | Mittel |

### 3.2 Empfehlung: Hybrid-Ansatz

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         KI-ARCHITEKTUR (HYBRID)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                        ┌─────────────────────┐                              │
│                        │    AI Gateway       │                              │
│                        │    (Go Service)     │                              │
│                        └──────────┬──────────┘                              │
│                                   │                                          │
│              ┌────────────────────┼────────────────────┐                    │
│              │                    │                    │                    │
│              ▼                    ▼                    ▼                    │
│  ┌───────────────────┐ ┌───────────────────┐ ┌───────────────────┐        │
│  │   Self-Hosted     │ │  Vector Store     │ │  External LLM     │        │
│  │   LLM (Ollama)    │ │  (Milvus/Qdrant)  │ │  (falls nötig)    │        │
│  │                   │ │                   │ │                   │        │
│  │  • Llama 3 70B    │ │  • Produkt-Emb.   │ │  • Azure OpenAI   │        │
│  │  • Mistral        │ │  • Doku-Emb.      │ │  • Nur anonyme    │        │
│  │  • Für sensible   │ │  • FAQ-Emb.       │ │    Anfragen       │        │
│  │    Anfragen       │ │                   │ │                   │        │
│  └───────────────────┘ └───────────────────┘ └───────────────────┘        │
│              │                    │                    │                    │
│              └────────────────────┼────────────────────┘                    │
│                                   │                                          │
│                        ┌──────────▼──────────┐                              │
│                        │   RAG Pipeline      │                              │
│                        │                     │                              │
│                        │ 1. Query verstehen  │                              │
│                        │ 2. Relevante Docs   │                              │
│                        │ 3. Context bauen    │                              │
│                        │ 4. LLM Antwort      │                              │
│                        │ 5. Quellen angeben  │                              │
│                        └─────────────────────┘                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Datenfluss-Regeln

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  DATENKLASSIFIZIERUNG                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  🔴 STRENG INTERN (nur Self-Hosted LLM)                                     │
│  ────────────────────────────────────────                                   │
│  • Kundendaten (Namen, Adressen, Bestellungen)                             │
│  • Preise und Margen                                                        │
│  • Lieferanten-Konditionen                                                  │
│  • Interne Verkaufszahlen                                                   │
│                                                                              │
│  🟡 VERTRAULICH (Private Cloud mit Vertrag OK)                              │
│  ────────────────────────────────────────────                               │
│  • Produktbeschreibungen                                                    │
│  • Technische Datenblätter                                                  │
│  • Allgemeine FAQs                                                          │
│                                                                              │
│  🟢 ÖFFENTLICH (Externe API OK)                                             │
│  ────────────────────────────────                                           │
│  • Allgemeine Fragen ohne Kontext                                          │
│  • Anonymisierte Anfragen                                                   │
│  • Öffentlich verfügbare Infos                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.4 RAG-Pipeline Detail

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  RAG (Retrieval Augmented Generation)                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  INDEXIERUNG (offline):                                                     │
│  ──────────────────────                                                     │
│                                                                              │
│  Produktdaten ──┐                                                           │
│  Datenblätter ──┼──▶ Chunking ──▶ Embedding ──▶ Vector Store               │
│  FAQs ──────────┤      │              │           (Milvus)                  │
│  Support-Docs ──┘      │              │                                     │
│                        │              │                                     │
│                   Split in        Llama/                                    │
│                   ~500 Token     Mistral                                    │
│                   Chunks         Embeddings                                 │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  ANFRAGE (online):                                                          │
│  ─────────────────                                                          │
│                                                                              │
│  User Query                                                                 │
│      │                                                                       │
│      ▼                                                                       │
│  ┌─────────────────┐                                                        │
│  │ Query Embedding │                                                        │
│  └────────┬────────┘                                                        │
│           │                                                                  │
│           ▼                                                                  │
│  ┌─────────────────┐     ┌─────────────────┐                               │
│  │ Similarity      │────▶│ Top-K Chunks    │                               │
│  │ Search          │     │ (k=5)           │                               │
│  └─────────────────┘     └────────┬────────┘                               │
│                                   │                                         │
│                                   ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ PROMPT                                                               │   │
│  │                                                                       │   │
│  │ System: Du bist ein Produktberater für Bodenbeläge...               │   │
│  │                                                                       │   │
│  │ Context:                                                             │   │
│  │ [Chunk 1: Vinyl AquaProtect ist 100% wasserfest...]                 │   │
│  │ [Chunk 2: Für Badezimmer empfehlen wir...]                          │   │
│  │ [Chunk 3: Preisliste Vinyl: AquaProtect CHF 34.90...]               │   │
│  │                                                                       │   │
│  │ User: Ich suche wasserfesten Bodenbelag fürs Bad                    │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                   │                                         │
│                                   ▼                                         │
│                          ┌─────────────────┐                               │
│                          │      LLM        │                               │
│                          │   (Llama 3)     │                               │
│                          └────────┬────────┘                               │
│                                   │                                         │
│                                   ▼                                         │
│                            Antwort + Quellen                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Infrastruktur

### 4.1 Self-Hosted LLM mit Ollama

```yaml
# Kubernetes Deployment für Ollama
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ollama
  namespace: ai
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ollama
  template:
    metadata:
      labels:
        app: ollama
    spec:
      containers:
      - name: ollama
        image: ollama/ollama:latest
        ports:
        - containerPort: 11434
        resources:
          requests:
            memory: "16Gi"
            cpu: "4000m"
            nvidia.com/gpu: 1  # Falls GPU verfügbar
          limits:
            memory: "32Gi"
            cpu: "8000m"
            nvidia.com/gpu: 1
        volumeMounts:
        - name: models
          mountPath: /root/.ollama
      volumes:
      - name: models
        persistentVolumeClaim:
          claimName: ollama-models
---
# Models vorladen
apiVersion: batch/v1
kind: Job
metadata:
  name: ollama-pull-models
spec:
  template:
    spec:
      containers:
      - name: pull
        image: curlimages/curl
        command:
        - sh
        - -c
        - |
          curl -X POST http://ollama:11434/api/pull -d '{"name": "llama3:70b"}'
          curl -X POST http://ollama:11434/api/pull -d '{"name": "mistral:7b"}'
      restartPolicy: OnFailure
```

### 4.2 Vector Store (Milvus)

```yaml
# Milvus für Embeddings
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: milvus
  namespace: ai
spec:
  serviceName: milvus
  replicas: 1
  selector:
    matchLabels:
      app: milvus
  template:
    metadata:
      labels:
        app: milvus
    spec:
      containers:
      - name: milvus
        image: milvusdb/milvus:v2.3-latest
        ports:
        - containerPort: 19530
        - containerPort: 9091
        env:
        - name: ETCD_ENDPOINTS
          value: "etcd:2379"
        - name: MINIO_ADDRESS
          value: "minio:9000"
        resources:
          requests:
            memory: "4Gi"
            cpu: "2000m"
          limits:
            memory: "8Gi"
            cpu: "4000m"
```

### 4.3 AI Gateway Service

```go
// services/ai-gateway/internal/gateway/service.go
package gateway

type Service struct {
    ollama      *ollama.Client
    vectorStore *milvus.Client
    classifier  *DataClassifier
}

// Chat verarbeitet eine Anfrage mit der richtigen LLM-Auswahl
func (s *Service) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
    // 1. Datenklassifizierung
    classification := s.classifier.Classify(req)

    // 2. LLM auswählen basierend auf Sensitivität
    var llm LLMClient
    switch classification.Level {
    case DataLevelInternal:
        llm = s.ollama // Nur lokal
    case DataLevelConfidential:
        llm = s.azureOpenAI // Private Cloud
    case DataLevelPublic:
        llm = s.selectBestAvailable()
    }

    // 3. RAG: Relevante Dokumente finden
    docs, err := s.vectorStore.SimilaritySearch(ctx, req.Query, 5)
    if err != nil {
        return nil, err
    }

    // 4. Prompt bauen
    prompt := s.buildPrompt(req, docs, classification)

    // 5. LLM aufrufen
    response, err := llm.Complete(ctx, prompt)
    if err != nil {
        return nil, err
    }

    // 6. Quellen anhängen
    return &ChatResponse{
        Answer:  response.Text,
        Sources: s.extractSources(docs),
        Model:   llm.Name(),
    }, nil
}

// DataClassifier bestimmt die Sensitivität einer Anfrage
type DataClassifier struct {
    patterns []ClassificationRule
}

func (c *DataClassifier) Classify(req *ChatRequest) *Classification {
    // Prüfe auf sensible Inhalte
    if containsCustomerData(req) || containsPricing(req) {
        return &Classification{Level: DataLevelInternal}
    }
    if containsProductData(req) {
        return &Classification{Level: DataLevelConfidential}
    }
    return &Classification{Level: DataLevelPublic}
}
```

---

## 5. UI-Integration

### 5.1 Kunden-Chat Widget

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  SHOP FRONTEND                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────┐                          │
│  │ 🔍 Suche: [wasserfester Boden Bad        ] 🔎│                          │
│  └──────────────────────────────────────────────┘                          │
│                                                                              │
│  [Kategorien] [Angebote] [Neu] [...]           [🛒 Warenkorb] [👤 Login]   │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│                                                                              │
│                    ... Shop Content ...                                      │
│                                                                              │
│                                                                              │
│                                                                              │
│                                              ┌──────────────────────────┐   │
│                                              │ 💬 Produktberater        │   │
│                                              │                          │   │
│                                              │ Wie kann ich helfen?     │   │
│                                              │                          │   │
│                                              │ ○ Produkt finden         │   │
│                                              │ ○ Technische Frage       │   │
│                                              │ ○ Bestellung verfolgen   │   │
│                                              │                          │   │
│                                              │ [___________________]    │   │
│                                              │ [Fragen Sie mich...]  ➤  │   │
│                                              └──────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Interner Sales Assistant

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  ADMIN PORTAL - SALES ASSISTANT                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  [Dashboard] [Kunden] [Bestellungen] [🤖 AI Assistant] [...]               │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                       │   │
│  │  🤖 Sales Assistant                                      [Historie]  │   │
│  │  ─────────────────                                                    │   │
│  │                                                                       │   │
│  │  Du: Zeig mir Kunden die letztes Jahr viel bestellt haben            │   │
│  │      aber dieses Jahr weniger                                         │   │
│  │                                                                       │   │
│  │  ────────────────────────────────────────────────────────────────    │   │
│  │                                                                       │   │
│  │  🤖 Hier sind Kunden mit Umsatzrückgang (>20%):                      │   │
│  │                                                                       │   │
│  │  │ Kunde                  │ 2024      │ 2025 YTD  │ Diff    │        │   │
│  │  ├────────────────────────┼───────────┼───────────┼─────────┤        │   │
│  │  │ Schreinerei Müller AG  │ CHF 152k  │ CHF 28k   │ -45%    │        │   │
│  │  │ Holzbau Weber          │ CHF 89k   │ CHF 18k   │ -38%    │        │   │
│  │  │ Parkett Plus GmbH      │ CHF 67k   │ CHF 15k   │ -32%    │        │   │
│  │                                                                       │   │
│  │  Empfehlung: Diese Kunden könnten von einer Kontaktaufnahme          │   │
│  │  profitieren. Soll ich Gesprächspunkte vorbereiten?                  │   │
│  │                                                                       │   │
│  │  [Ja, für Schreinerei Müller] [Alle exportieren] [Neue Frage]        │   │
│  │                                                                       │   │
│  │  ────────────────────────────────────────────────────────────────    │   │
│  │                                                                       │   │
│  │  [___________________________________________________] [Fragen ➤]   │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Datenschutz & Compliance

### 6.1 Vertragsanforderungen (falls externe API)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  ANFORDERUNGEN AN EXTERNE LLM-ANBIETER                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Muss vertraglich garantiert sein:                                          │
│                                                                              │
│  ✅ Kein Training auf unseren Daten                                         │
│  ✅ Daten werden nicht gespeichert (oder max. 30 Tage)                      │
│  ✅ Auftragsverarbeitungsvertrag (AVV) nach DSGVO                           │
│  ✅ Datenverarbeitung in EU/CH                                              │
│  ✅ SOC 2 Type II Zertifizierung                                            │
│  ✅ Recht auf Löschung                                                      │
│                                                                              │
│  Anbieter die das erfüllen:                                                 │
│  • Azure OpenAI (mit Enterprise Agreement)                                  │
│  • Anthropic Claude (mit Business Agreement)                                │
│  • Google Vertex AI (mit entsprechendem Vertrag)                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Audit-Logging

```go
// Jede KI-Anfrage wird geloggt
type AIAuditLog struct {
    ID              string    `json:"id"`
    Timestamp       time.Time `json:"timestamp"`
    UserID          string    `json:"user_id"`
    UserType        string    `json:"user_type"` // customer, sales, support
    Query           string    `json:"query"`     // Hash oder anonymisiert
    DataClassLevel  string    `json:"data_class"`
    LLMUsed         string    `json:"llm_used"`
    ResponseTime    int       `json:"response_time_ms"`
    TokensUsed      int       `json:"tokens_used"`
    SourcesUsed     []string  `json:"sources_used"`
}
```

---

## 7. Rollout-Phasen

### 7.1 Phasen-Plan

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  ROLLOUT PHASEN                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  PHASE 1: Infrastruktur                                                     │
│  ─────────────────────────                                                  │
│  • Ollama + Milvus in Kubernetes deployen                                   │
│  • AI Gateway Service erstellen                                             │
│  • Produkt-Embeddings generieren                                            │
│  • Interne Tests                                                            │
│                                                                              │
│  PHASE 2: Interner Sales Assistant (Pilot)                                  │
│  ─────────────────────────────────────────                                  │
│  • Use Case 5 implementieren                                                │
│  • Nur für ausgewählte Verkäufer                                           │
│  • Feedback sammeln, iterieren                                              │
│  • Kein Kundenrisiko                                                        │
│                                                                              │
│  PHASE 3: Technische Beratung (Kunden)                                     │
│  ─────────────────────────────────────                                      │
│  • Use Case 2 implementieren                                                │
│  • RAG über Produktdatenblätter                                            │
│  • Beta mit ausgewählten Kunden                                            │
│  • Prominente "Beta"-Kennzeichnung                                         │
│                                                                              │
│  PHASE 4: Produktsuche Natural Language                                     │
│  ─────────────────────────────────────                                      │
│  • Use Case 1 in Suche integrieren                                         │
│  • A/B Test: Klassisch vs. KI-unterstützt                                  │
│  • Metriken: Conversion, Findability                                       │
│                                                                              │
│  PHASE 5: Vollständiger Chat-Support                                       │
│  ───────────────────────────────────                                        │
│  • Use Cases 3 + 4                                                          │
│  • Eskalation an menschlichen Support                                      │
│  • 24/7 Verfügbarkeit für Basis-Anfragen                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Empfehlung: Start mit Phase 1-2

**Warum intern zuerst?**
- Kein Kundenrisiko bei Fehlern
- Schnelles Feedback von eigenem Team
- Datenklassifizierung kann getestet werden
- Verkäufer können "trainieren" was gute Antworten sind

---

## 8. Zusammenfassung

### 8.1 Empfohlene Lösung

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  KI-ASSISTENTEN STACK                                                       │
│  ════════════════════                                                       │
│                                                                              │
│  LLM (primär):     Ollama mit Llama 3 / Mistral (self-hosted)              │
│  LLM (optional):   Azure OpenAI (für weniger sensible Anfragen)            │
│  Vector Store:     Milvus (self-hosted)                                     │
│  Embeddings:       Llama 3 oder sentence-transformers                       │
│  Gateway:          Go Service im V3 Stack                                   │
│                                                                              │
│  Prioritäre Use Cases:                                                      │
│  1. Sales Assistant (intern) ← Start hier                                  │
│  2. Technische Beratung (Kunden)                                           │
│  3. Produktsuche Natural Language                                          │
│                                                                              │
│  Datenschutz:                                                               │
│  • Sensible Daten nur über Self-Hosted LLM                                 │
│  • Audit-Logging für alle Anfragen                                         │
│  • Datenklassifizierung automatisch                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Nächste Schritte

1. [ ] Hardware-Anforderungen für LLM klären (GPU?)
2. [ ] Ollama PoC lokal aufsetzen
3. [ ] Produkt-Embeddings mit Testdaten generieren
4. [ ] Sales Assistant Prototyp für einen Use Case
5. [ ] Feedback von Verkaufsteam einholen
