# Datahub Konzept

## 1. Vision

Der **Datahub** ist die zentrale Integrationsschicht für alle externen Schnittstellen. Statt für jede Integration eigenen Code zu schreiben, werden Verbindungen **visuell konfiguriert** - nicht programmiert.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              DATAHUB VISION                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  HEUTE (V2):                         MORGEN (V3 Datahub):                   │
│  ────────────                        ────────────────────                   │
│                                                                              │
│  Jede Integration = Code             Jede Integration = Visuelle Pipeline   │
│  ├── SAP: eigener Service            ├── Pipeline im Web UI designen       │
│  ├── Akeneo: eigene Jobs             ├── Transforms konfigurieren          │
│  ├── Import: eigene Parser           ├── Aktivieren                        │
│  └── Jeder Change = Deployment       └── Fertig (kein Deployment)          │
│                                                                              │
│  Probleme:                           Vorteile:                              │
│  • Viel Boilerplate                  • Schnelle Integration neuer Partner  │
│  • Fehleranfällig                    • IT konfiguriert selbst              │
│  • Schwer testbar                    • Einheitliches Monitoring            │
│  • Inkonsistentes Logging            • Automatische Retries & Error Hand.  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Technologie-Entscheidung: Apache Hop

### 2.1 Warum Apache Hop?

**Apache Hop** ist der direkte Nachfolger von Pentaho Kettle (PDI), entwickelt von den gleichen Machern. Es bietet:

| Anforderung | Apache Hop |
|-------------|------------|
| Web UI (nicht Desktop) | ✅ Hop Web |
| Visuelle Pipeline-Entwicklung | ✅ Drag & Drop |
| API Bereitstellung | ✅ Web Services |
| Komplexe Transformationen | ✅ >200 Transforms + JavaScript |
| Fixed-Width Parsing (Datanorm) | ✅ Text File Input |
| Kubernetes-ready | ✅ Container Images |
| Open Source | ✅ Apache 2.0 |

### 2.2 Kernkonzepte

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        APACHE HOP KERNKONZEPTE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         PIPELINES                                    │   │
│  │                                                                       │   │
│  │  Die Arbeitspferde: Daten lesen, transformieren, schreiben          │   │
│  │  Alle Transforms laufen PARALLEL                                     │   │
│  │                                                                       │   │
│  │  ┌─────┐    ┌─────┐    ┌─────┐    ┌─────┐    ┌─────┐              │   │
│  │  │Input│───▶│Trans│───▶│Trans│───▶│Trans│───▶│Output│              │   │
│  │  └─────┘    └─────┘    └─────┘    └─────┘    └─────┘              │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         WORKFLOWS                                    │   │
│  │                                                                       │   │
│  │  Orchestrierung: Pipelines steuern, Error Handling, Scheduling      │   │
│  │  Actions laufen SEQUENZIELL                                          │   │
│  │                                                                       │   │
│  │  ┌─────┐    ┌─────┐    ┌─────┐    ┌─────┐                          │   │
│  │  │Start│───▶│Check│─┬─▶│Pipe │───▶│Email│                          │   │
│  │  └─────┘    └─────┘ │  │line │    │     │                          │   │
│  │                     │  └─────┘    └─────┘                          │   │
│  │                     │                                                │   │
│  │                     └─▶[Error]───▶[Alert]                           │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         METADATA                                     │   │
│  │                                                                       │   │
│  │  Alles ist Metadata-getrieben:                                       │   │
│  │  • Connections (DB, REST, SFTP, S3, ...)                            │   │
│  │  • Web Services (API Endpoints)                                      │   │
│  │  • Environments (Dev, Staging, Prod)                                │   │
│  │  • Run Configurations                                                │   │
│  │                                                                       │   │
│  │  → Versionierbar in Git, deploybar via CI/CD                        │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Deployment-Komponenten

| Komponente | Beschreibung | Verwendung |
|------------|--------------|------------|
| **Hop Web** | Browser-basierte GUI | Pipeline-Entwicklung, Testing |
| **Hop Server** | Headless Runtime + REST API | Produktion, Web Services |
| **Hop Run** | CLI für Container/K8s | Batch Jobs, Scheduled Tasks |

---

## 3. Architektur-Übersicht

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        DATAHUB ARCHITEKTUR (V3)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  EXTERNE SYSTEME                                                            │
│  ───────────────                                                            │
│                                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │Konfigurator │  │  SFTP/S3    │  │ Kunden-ERP  │  │   Akeneo    │       │
│  │  (REST)     │  │ (Dateien)   │  │  (REST)     │  │  (Webhook)  │       │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │
│         │                │                │                │               │
│         │ POST           │ FILE           │ POST           │ POST          │
│         │                │                │                │               │
│         └────────────────┴────────────────┴────────────────┘               │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         APACHE HOP                                   │   │
│  │  ┌─────────────────────────────────────────────────────────────┐   │   │
│  │  │                      HOP WEB                                 │   │   │
│  │  │                                                              │   │   │
│  │  │  Browser-basierte Entwicklungsumgebung                      │   │   │
│  │  │  • Pipeline Designer (Drag & Drop)                          │   │   │
│  │  │  • Workflow Designer                                         │   │   │
│  │  │  • Metadata Management                                       │   │   │
│  │  │  • Testing & Debugging                                       │   │   │
│  │  │                                                              │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐   │   │
│  │  │                     HOP SERVER                               │   │   │
│  │  │                                                              │   │   │
│  │  │  Runtime für Produktion                                      │   │   │
│  │  │  • Web Services (REST Endpoints)                            │   │   │
│  │  │  • Workflow Execution                                        │   │   │
│  │  │  • Scheduling                                                │   │   │
│  │  │  • Monitoring & Logging                                      │   │   │
│  │  │                                                              │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  │                                                                      │   │
│  └──────────────────────────────┬──────────────────────────────────────┘   │
│                                 │                                           │
│                                 │ REST API Calls                            │
│                                 ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                       V3 MICROSERVICES                               │   │
│  │                                                                       │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │   Order     │  │  Catalog    │  │  Inventory  │                  │   │
│  │  │   Service   │  │   Service   │  │   Service   │                  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                  │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Verfügbare Transforms

Apache Hop bietet über **200 vorgefertigte Transforms**:

### 4.1 Input/Output

| Transform | Beschreibung | Use Case |
|-----------|--------------|----------|
| **Text File Input** | CSV, Fixed-Width, Delimited | Datanorm, ComNorm |
| **CSV File Input** | Optimiert für CSV | Standard-Imports |
| **JSON Input** | JSON Dateien/Streams | Konfigurator-Aufträge |
| **XML Input** | XML Dateien | Legacy-Formate |
| **REST Client** | HTTP API Calls | Shop-Services aufrufen |
| **Table Input** | SQL Queries | Datenbank-Abfragen |
| **S3 File Input** | AWS S3 Dateien | Cloud-Imports |
| **Kafka Consumer** | Event Streams | Real-time Events |

### 4.2 Transformation

| Transform | Beschreibung | Use Case |
|-----------|--------------|----------|
| **JavaScript** | Komplexe Logik | Berechnungen, Validierung |
| **Value Mapper** | Wert-Übersetzungen | Code-Mappings |
| **Filter Rows** | Bedingtes Filtern | Nur bestimmte Sätze |
| **Switch/Case** | Verzweigung | Nach Satzart verarbeiten |
| **Calculator** | Berechnungen | Preise, Mengen |
| **String Operations** | Text-Manipulation | Formatierung |
| **Regex Evaluation** | Pattern Matching | Parsing |
| **Data Validator** | Schema-Validierung | Input-Prüfung |

### 4.3 Output

| Transform | Beschreibung | Use Case |
|-----------|--------------|----------|
| **REST Client** | HTTP POST/PUT | An Shop-Services senden |
| **JSON Output** | JSON generieren | API Responses |
| **Table Output** | DB Insert/Update | Direktes Speichern |
| **Text File Output** | CSV/Fixed-Width | Export-Dateien |
| **Kafka Producer** | Event Publishing | Async Events |

---

## 5. Use Cases im Detail

### 5.1 Datanorm Produkt-Import

**Szenario:** Lieferant liefert Produktdaten im Datanorm-Format (Fixed-Width, mehrere Satzarten).

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  PIPELINE: datanorm-product-import                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐                                                           │
│  │  Text File   │  Datanorm-Datei (Fixed-Width)                            │
│  │    Input     │  • Satzart Position 1-1                                  │
│  │              │  • Artikelnummer Position 2-12                           │
│  │              │  • etc.                                                   │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │   Switch /   │  Nach Satzart verzweigen:                                │
│  │    Case      │  • "A" → Artikelstamm                                    │
│  │              │  • "B" → Preisdaten                                      │
│  │              │  • "P" → Warengruppe                                     │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│    ┌────┴────┬────────────┐                                                 │
│    ▼         ▼            ▼                                                 │
│  ┌─────┐  ┌─────┐     ┌─────┐                                              │
│  │ "A" │  │ "B" │     │ "P" │                                              │
│  │Parse│  │Parse│     │Parse│                                              │
│  └──┬──┘  └──┬──┘     └──┬──┘                                              │
│     │        │           │                                                  │
│     └────────┴───────────┘                                                  │
│              │                                                               │
│              ▼                                                               │
│  ┌──────────────┐                                                           │
│  │  JavaScript  │  Komplexe Logik:                                         │
│  │  Transform   │  • Preis-Berechnung (Staffeln)                          │
│  │              │  • EAN-Validierung                                       │
│  │              │  • Kategorie-Mapping                                     │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │    Data      │  Pflichtfelder prüfen                                    │
│  │  Validator   │  • SKU vorhanden?                                        │
│  │              │  • Preis > 0?                                            │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │    REST      │  POST /api/v1/products/import                            │
│  │   Client     │  → Catalog Service                                       │
│  │              │                                                           │
│  └──────────────┘                                                           │
│                                                                              │
│  Trigger: S3 Bucket Watch oder Scheduled (täglich 02:00)                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Konfigurator Auftrags-Import (API Endpoint)

**Szenario:** Externer Konfigurator sendet Aufträge per REST API.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  WEB SERVICE: konfigurator-order-import                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Endpoint: POST /hop/webService/?service=konfigurator-order                 │
│                                                                              │
│  Request Body:                                                              │
│  {                                                                           │
│    "konfigurator_id": "XYZ-123",                                           │
│    "customer_number": "CUST-456",                                          │
│    "items": [                                                                │
│      { "sku": "PROD-A", "quantity": 2, "config": {...} },                  │
│      { "sku": "PROD-B", "quantity": 1, "config": {...} }                   │
│    ],                                                                        │
│    "delivery_address": {...}                                                │
│  }                                                                           │
│                                                                              │
│  ┌──────────────┐                                                           │
│  │ Get Request  │  JSON Body aus POST Request                              │
│  │    Body      │                                                           │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │    JSON      │  Request parsen                                          │
│  │   Input      │                                                           │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │    Data      │  • customer_number vorhanden?                            │
│  │  Validator   │  • items nicht leer?                                     │
│  │              │  • SKUs gültig?                                          │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │  JavaScript  │  • Config in Produktoptionen umwandeln                   │
│  │  Transform   │  • Preise aus Catalog Service holen                      │
│  │              │  • Order-Struktur aufbauen                               │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │    REST      │  POST /api/v1/orders                                     │
│  │   Client     │  → Order Service                                         │
│  └──────┬───────┘                                                           │
│         │                                                                    │
│         ▼                                                                    │
│  ┌──────────────┐                                                           │
│  │    JSON      │  Response für Konfigurator                               │
│  │   Output     │  { "order_id": "...", "status": "created" }             │
│  └──────────────┘                                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 ComNorm Bestell-Import

**Szenario:** B2B-Kunde sendet Bestellungen im ComNorm-Format.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  WORKFLOW: comnorm-order-import                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐                 │
│  │  Start  │───▶│  Check  │───▶│ Pipeline│───▶│  Move   │                 │
│  │         │    │  Files  │    │  Run    │    │  File   │                 │
│  └─────────┘    └────┬────┘    └────┬────┘    └────┬────┘                 │
│                      │              │              │                        │
│                      │              │              ▼                        │
│                      │              │         ┌─────────┐                  │
│                      │              │         │ Success │                  │
│                      │              │         │  Email  │                  │
│                      │              │         └─────────┘                  │
│                      │              │                                       │
│                      │              └─────[Error]─────┐                    │
│                      │                                │                    │
│                      │                                ▼                    │
│                      │                          ┌─────────┐               │
│                      │                          │  Error  │               │
│                      │                          │  Alert  │               │
│                      │                          └─────────┘               │
│                      │                                                     │
│                      └─────[No Files]────▶ [End]                          │
│                                                                              │
│  Pipeline "comnorm-order-parse":                                            │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐               │
│  │  Text    │──▶│  Parse   │──▶│ Validate │──▶│  REST    │               │
│  │  Input   │   │ ComNorm  │   │  Order   │   │  Client  │               │
│  └──────────┘   └──────────┘   └──────────┘   └──────────┘               │
│                                                                              │
│  Trigger: SFTP Watch oder Schedule (alle 15 Min)                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Kubernetes Deployment

### 6.1 Architektur

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  KUBERNETES NAMESPACE: datahub                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Ingress                                                             │   │
│  │  ────────                                                            │   │
│  │  datahub.shop.ch/web    → hop-web:8080                              │   │
│  │  datahub.shop.ch/api    → hop-server:8182                           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────────────┐     ┌──────────────────────┐                    │
│  │  Deployment:         │     │  Deployment:         │                    │
│  │  hop-web             │     │  hop-server          │                    │
│  │                      │     │                      │                    │
│  │  Image:              │     │  Image:              │                    │
│  │  apache/hop-web:2.8  │     │  apache/hop:2.8      │                    │
│  │                      │     │                      │                    │
│  │  Replicas: 1         │     │  Replicas: 2-5 (HPA)│                    │
│  │  (Development)       │     │  (Production)        │                    │
│  │                      │     │                      │                    │
│  │  Mounts:             │     │  Mounts:             │                    │
│  │  - /hop/projects     │     │  - /hop/projects     │                    │
│  │  - /hop/config       │     │  - /hop/config       │                    │
│  └──────────┬───────────┘     └──────────┬───────────┘                    │
│             │                            │                                 │
│             └────────────┬───────────────┘                                 │
│                          │                                                  │
│                          ▼                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  PersistentVolumeClaim: hop-projects                                │   │
│  │                                                                       │   │
│  │  Enthält:                                                            │   │
│  │  - Pipelines (.hpl)                                                  │   │
│  │  - Workflows (.hwf)                                                  │   │
│  │  - Metadata (connections, web services, etc.)                       │   │
│  │                                                                       │   │
│  │  Synchronisiert via Git oder Volume Mount                           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────────────┐     ┌──────────────────────┐                    │
│  │  StatefulSet:        │     │  ConfigMap:          │                    │
│  │  hop-postgres        │     │  hop-config          │                    │
│  │                      │     │                      │                    │
│  │  Hop Audit DB        │     │  - hop-server.xml    │                    │
│  │  (optional)          │     │  - environments      │                    │
│  └──────────────────────┘     └──────────────────────┘                    │
│                                                                              │
│  ┌──────────────────────┐                                                  │
│  │  Secret:             │                                                  │
│  │  hop-credentials     │                                                  │
│  │                      │                                                  │
│  │  - DB Passwords      │                                                  │
│  │  - API Keys          │                                                  │
│  │  - SFTP Keys         │                                                  │
│  └──────────────────────┘                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Kubernetes Manifests

```yaml
# deployment-hop-server.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hop-server
  namespace: datahub
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hop-server
  template:
    metadata:
      labels:
        app: hop-server
    spec:
      containers:
      - name: hop-server
        image: apache/hop:2.8.0
        command: ["/opt/hop/hop-server.sh"]
        args: ["/hop/config/hop-server.xml"]
        ports:
        - containerPort: 8182
        env:
        - name: HOP_PROJECT_FOLDER
          value: "/hop/projects/datahub"
        - name: HOP_ENVIRONMENT
          value: "production"
        - name: HOP_LOG_LEVEL
          value: "Basic"
        volumeMounts:
        - name: projects
          mountPath: /hop/projects
        - name: config
          mountPath: /hop/config
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /hop/status
            port: 8182
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /hop/status
            port: 8182
          initialDelaySeconds: 10
          periodSeconds: 5
      volumes:
      - name: projects
        persistentVolumeClaim:
          claimName: hop-projects
      - name: config
        configMap:
          name: hop-config
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hop-server-hpa
  namespace: datahub
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: hop-server
  minReplicas: 2
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

```yaml
# hop-server-config.xml (ConfigMap)
<hop-server-config>
  <webserver>
    <hostname>0.0.0.0</hostname>
    <port>8182</port>
    <shutdownInterface>true</shutdownInterface>
  </webserver>

  <metadata_folder>/hop/projects/datahub/metadata</metadata_folder>

  <max_log_lines>10000</max_log_lines>
  <max_log_timeout_minutes>1440</max_log_timeout_minutes>
  <object_timeout_minutes>1440</object_timeout_minutes>
</hop-server-config>
```

---

## 7. Environment-Konzept

Apache Hop unterstützt die Trennung von Code und Konfiguration:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  PROJECTS & ENVIRONMENTS                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  PROJECT: datahub                                                           │
│  ────────────────                                                           │
│  Enthält die Pipelines und Workflows (WAS wird verarbeitet)                │
│                                                                              │
│  datahub/                                                                   │
│  ├── pipelines/                                                             │
│  │   ├── datanorm-import.hpl                                               │
│  │   ├── comnorm-order.hpl                                                 │
│  │   └── konfigurator-order.hpl                                            │
│  ├── workflows/                                                             │
│  │   ├── daily-import.hwf                                                  │
│  │   └── order-processing.hwf                                              │
│  └── metadata/                                                              │
│      ├── connections/                                                       │
│      └── web-services/                                                      │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                              │
│  ENVIRONMENTS (WO wird verarbeitet)                                         │
│  ──────────────────────────────────                                         │
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐            │
│  │   Development   │  │    Staging      │  │   Production    │            │
│  │                 │  │                 │  │                 │            │
│  │ SHOP_API_URL=   │  │ SHOP_API_URL=   │  │ SHOP_API_URL=   │            │
│  │ http://localhost│  │ https://staging │  │ https://api     │            │
│  │                 │  │                 │  │                 │            │
│  │ SFTP_HOST=      │  │ SFTP_HOST=      │  │ SFTP_HOST=      │            │
│  │ localhost       │  │ sftp-staging    │  │ sftp.partner.ch │            │
│  │                 │  │                 │  │                 │            │
│  │ DB_HOST=        │  │ DB_HOST=        │  │ DB_HOST=        │            │
│  │ localhost:5432  │  │ db-staging      │  │ db-prod-cluster │            │
│  │                 │  │                 │  │                 │            │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘            │
│                                                                              │
│  → Gleiche Pipeline läuft in allen Environments                            │
│  → Nur Variablen ändern sich                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Monitoring & Observability

### 8.1 Hop Server Status

Hop Server bietet eine eingebaute Status-Seite:

```
http://hop-server:8182/hop/status

┌─────────────────────────────────────────────────────────────────────────────┐
│  HOP SERVER STATUS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Server: hop-server-5d7b9c4f8-abc12                                         │
│  Uptime: 3 days, 14 hours                                                   │
│  Memory: 512MB / 2048MB                                                     │
│                                                                              │
│  ACTIVE PIPELINES                                                           │
│  ─────────────────                                                          │
│  │ Pipeline              │ Status    │ Started     │ Rows      │          │
│  ├───────────────────────┼───────────┼─────────────┼───────────┤          │
│  │ datanorm-import       │ Running   │ 14:32:05    │ 12,847    │          │
│  │ order-processing      │ Idle      │ -           │ -         │          │
│  │                                                                          │
│  RECENT EXECUTIONS                                                          │
│  ─────────────────                                                          │
│  │ Pipeline              │ Status    │ Duration    │ Rows      │          │
│  ├───────────────────────┼───────────┼─────────────┼───────────┤          │
│  │ konfigurator-order    │ ✅ Success│ 1.2s        │ 1         │          │
│  │ datanorm-import       │ ✅ Success│ 45.3s       │ 8,234     │          │
│  │ comnorm-order         │ ❌ Failed │ 0.8s        │ 0         │          │
│  │                                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Integration mit V3 Observability

```yaml
# Hop Logs → Loki (via stdout JSON)
env:
- name: HOP_LOG_LEVEL
  value: "Basic"
- name: HOP_OPTIONS
  value: "-Dhop.log.stdout.layout=JSON"

# Metrics → Prometheus (Custom Exporter oder Log-Parsing)
# Traces → Jaeger (über REST Client Instrumentation)
```

### 8.3 Alerting

```yaml
# Alert: Pipeline Failures
groups:
- name: datahub
  rules:
  - alert: HopPipelineFailed
    expr: |
      increase(hop_pipeline_failures_total[5m]) > 0
    labels:
      severity: warning
    annotations:
      summary: "Pipeline {{ $labels.pipeline }} failed"

  - alert: HopServerDown
    expr: up{job="hop-server"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Hop Server is down"
```

---

## 9. Sicherheit

### 9.1 Authentifizierung

```xml
<!-- Hop Server mit Basic Auth -->
<hop-server-config>
  <webserver>
    <hostname>0.0.0.0</hostname>
    <port>8182</port>
  </webserver>

  <security>
    <username>admin</username>
    <password>${HOP_ADMIN_PASSWORD}</password>
  </security>
</hop-server-config>
```

Für Produktion: Ingress mit OAuth2 Proxy oder mTLS.

### 9.2 Secrets Management

```yaml
# Credentials als Kubernetes Secrets
apiVersion: v1
kind: Secret
metadata:
  name: hop-credentials
  namespace: datahub
type: Opaque
stringData:
  SFTP_PASSWORD: "xxx"
  DB_PASSWORD: "xxx"
  SHOP_API_KEY: "xxx"

# In Hop als Environment Variables verwenden
# ${SFTP_PASSWORD} in Connection Metadata
```

### 9.3 Network Policies

```yaml
# Hop Server darf nur Shop-Services erreichen
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: hop-server-egress
  namespace: datahub
spec:
  podSelector:
    matchLabels:
      app: hop-server
  policyTypes:
  - Egress
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: webshop
    ports:
    - port: 8080
  - to:
    - namespaceSelector:
        matchLabels:
          name: datahub
    ports:
    - port: 5432  # PostgreSQL
```

---

## 10. Migration & Rollout

### 10.1 Phasen

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  ROLLOUT PHASEN                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  PHASE 1: Setup & PoC                                                       │
│  ─────────────────────                                                      │
│  • Hop Web + Server in Kubernetes deployen                                  │
│  • Erste Pipeline: Datanorm Import (ein Lieferant)                         │
│  • Connection zu Catalog Service testen                                     │
│  • Team-Schulung Hop Basics                                                 │
│                                                                              │
│  PHASE 2: Weitere Imports                                                   │
│  ─────────────────────────                                                  │
│  • ComNorm Bestell-Import                                                   │
│  • Weitere Datanorm-Lieferanten                                            │
│  • Monitoring & Alerting einrichten                                        │
│                                                                              │
│  PHASE 3: API Endpoints                                                     │
│  ───────────────────────                                                    │
│  • Konfigurator Auftrags-Import als Web Service                            │
│  • Kunden-ERP Integrationen                                                 │
│  • Dokumentation für Partner                                                │
│                                                                              │
│  PHASE 4: Erweiterte Integrationen                                         │
│  ─────────────────────────────────                                          │
│  • SAP Outbound über Hop (optional)                                        │
│  • Akeneo Webhook Processing                                                │
│  • Batch-Exports für B2B-Kunden                                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.2 PoC Scope

Für den ersten PoC empfehle ich:

1. **Use Case:** Datanorm Import (ein Lieferant)
2. **Komponenten:** Hop Web + Hop Server in Docker/K8s
3. **Ziel:** Datei parsen → REST Call an Catalog Service
4. **Dauer:** 1-2 Wochen

---

## 11. Zusammenfassung

### 11.1 Empfohlene Lösung

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  DATAHUB STACK                                                              │
│  ═════════════                                                              │
│                                                                              │
│  Tool:           Apache Hop                                                 │
│  Version:        2.8.x (aktuell)                                           │
│  Lizenz:         Apache 2.0 (Open Source)                                  │
│                                                                              │
│  Komponenten:                                                               │
│  • Hop Web       → Pipeline-Entwicklung (Browser)                          │
│  • Hop Server    → Runtime + Web Services                                  │
│  • Hop Run       → CLI für Batch Jobs                                      │
│                                                                              │
│  Deployment:     Kubernetes (Namespace: datahub)                           │
│  Storage:        PVC für Projects, Git für Versionierung                   │
│                                                                              │
│  Unterstützte Formate:                                                      │
│  • Datanorm (Fixed-Width)    ✅                                            │
│  • ComNorm                   ✅                                            │
│  • JSON/XML                  ✅                                            │
│  • CSV                       ✅                                            │
│  • REST APIs                 ✅                                            │
│                                                                              │
│  Vorteile:                                                                  │
│  • Pentaho-Erfahrung übertragbar                                           │
│  • Web UI (nicht Desktop)                                                   │
│  • API Bereitstellung möglich                                              │
│  • >200 vorgefertigte Transforms                                           │
│  • JavaScript für komplexe Logik                                           │
│  • Kubernetes-ready                                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 11.2 Nächste Schritte

1. [ ] Apache Hop lokal testen (Docker)
2. [ ] PoC: Datanorm Import Pipeline erstellen
3. [ ] Kubernetes Deployment aufsetzen
4. [ ] Team-Schulung durchführen
5. [ ] Erste Produktiv-Integration migrieren
