# Gondolia

**Open Source B2B E-Commerce Platform** built with Go microservices and Next.js.

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-Enabled-blue)](https://opentelemetry.io/)

---

## 🚀 Overview

Gondolia is a modern, cloud-native B2B e-commerce platform designed for enterprise use cases:

- **Multi-Tenant Architecture**: Isolated data, shared infrastructure
- **Provider Pattern**: Pluggable integrations for ERP, PIM, Payment, Search, and more
- **Microservices**: Clean 4-layer architecture (Handler → Service → Repository → Models)
- **Observability**: OpenTelemetry tracing, structured logging, Prometheus metrics
- **Production-Ready**: Kubernetes deployment, health checks, graceful shutdown

### Key Features

- 🔐 **Identity & Access**: JWT authentication, role-based access control, SSO ready
- 📦 **Provider System**: Swap ERP (SAP, Dynamics), PIM (Akeneo, Pimcore), Payment (Saferpay, Stripe), and more via interfaces
- 🌍 **Multi-Language**: i18n-ready frontend and backend
- 📊 **B2B Features**: Company management, hierarchical users, contract pricing, order approval workflows
- 🔍 **Full-Text Search**: Pluggable search providers (Meilisearch, Algolia, Elasticsearch)
- 🧪 **Testable**: Interface-based design, dependency injection, mocks included

---

## 📦 Architecture

### Services

| Service | Status | Description |
|---------|--------|-------------|
| **identity** | ✅ Complete | User, company, role management with JWT auth |
| **catalog** | ⏳ Planned | Product catalog, categories, pricing |
| **cart** | ⏳ Planned | Shopping cart with multi-currency support |
| **order** | ⏳ Planned | Order management with ERP integration |
| **inventory** | ⏳ Planned | Stock management and availability |
| **shipping** | ⏳ Planned | Shipping calculation and fulfillment |
| **notification** | ⏳ Planned | Email, SMS, push notifications |

### Provider Pattern

Gondolia uses a **provider pattern** for external integrations. Providers are registered via `init()` and discovered through a type-safe global registry.

**Available Provider Interfaces:**

- `ERPProvider` — Integrate with SAP, Microsoft Dynamics, etc.
- `PIMProvider` — Sync product data from Akeneo, Pimcore, etc.
- `PaymentProvider` — Accept payments via Saferpay, Stripe, Adyen, etc.
- `SearchProvider` — Full-text search with Meilisearch, Algolia, Elasticsearch, etc.
- `AuthProvider` — SSO with Azure AD, Keycloak, Auth0, etc.
- `StorageProvider` — Store files in S3, Azure Blob, MinIO, etc.
- `FulfillmentProvider` — Logistics integration (DHL, UPS, etc.)
- `NotificationProvider` — Send emails, SMS, push notifications
- `CRMProvider` — Sync contacts with Dynamics CRM, Salesforce, etc.
- `TaxProvider` — Calculate taxes (Avalara, TaxJar, etc.)

See [`docs/architecture/adr-001-provider-pattern.md`](./docs/architecture/adr-001-provider-pattern.md) for the full specification.

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.23+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15
- **Cache**: Redis
- **Message Queue**: Kafka (planned)
- **Tracing**: OpenTelemetry + Jaeger
- **Logging**: Zap (structured JSON logs)

### Frontend
- **Framework**: Next.js 14 (App Router)
- **UI**: React 18 + Tailwind CSS
- **State**: Zustand
- **API Client**: Fetch API with OpenAPI types

### Infrastructure
- **Orchestration**: Kubernetes (Helm charts included)
- **Dev Environment**: K3d (local Kubernetes)
- **CI/CD**: GitHub Actions (coming soon)
- **Registry**: `ghcr.io/gondolia/gondolia`

---

## 📖 Getting Started

### Prerequisites

- **Go 1.23+**
- **Node.js 20+**
- **Docker & Docker Compose**
- **K3d** (for local Kubernetes development)

### Quick Start

1. **Clone the repository**

```bash
git clone https://github.com/gondolia/gondolia.git
cd gondolia
```

2. **Start local infrastructure**

```bash
# Start PostgreSQL, Redis, Jaeger
docker-compose up -d
```

3. **Run Identity Service**

```bash
cd services/identity
make run
```

4. **Run Frontend**

```bash
cd frontend
npm install
npm run dev
```

5. **Access the application**

- **Frontend**: http://localhost:3000
- **Identity API**: http://localhost:8080
- **Jaeger UI**: http://localhost:16686

---

## 📚 Documentation

- [Architecture Overview](./docs/architecture/overview.md)
- [Service Layer Structure](./docs/architecture/service-layer-structure.md)
- [Provider Pattern (ADR-001)](./docs/architecture/adr-001-provider-pattern.md)
- [B2B Self-Service Principles](./docs/architecture/b2b-self-service-principles.md)
- [i18n Concept](./docs/architecture/i18n-concept.md)
- [Search Engine Evaluation](./docs/architecture/search-engine-evaluation.md)

---

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guide](./CONTRIBUTING.md) first.

### Contributor License Agreement (CLA)

Before we can merge your pull request, you need to sign our [Contributor License Agreement](./CONTRIBUTING.md#contributor-license-agreement-cla). This ensures that the project can remain open source under AGPL v3 while allowing commercial use for contributors.

---

## 📜 License

Gondolia is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

This means:
- ✅ You can use, modify, and distribute this software freely
- ✅ You can use it in commercial projects
- ⚠️ If you modify the software and run it as a network service, you **must** release your changes under AGPL-3.0
- ⚠️ If you integrate Gondolia into a SaaS product, you **must** release your modifications

See [LICENSE](./LICENSE) for the full text.

**Why AGPL?** We chose AGPL to ensure that improvements to Gondolia remain open source, even when the software is used as a service. This prevents "cloud lock-in" of open source improvements.

---

## 🌟 Roadmap

- [x] Identity Service (authentication, authorization)
- [x] Provider Pattern (pluggable integrations)
- [ ] Catalog Service (products, categories, pricing)
- [ ] Cart Service (shopping cart with promotions)
- [ ] Order Service (order processing with ERP integration)
- [ ] Inventory Service (stock management)
- [ ] Shipping Service (fulfillment integration)
- [ ] Notification Service (email, SMS, push)
- [ ] Admin Dashboard (management UI)
- [ ] GraphQL Gateway (unified API)

---

## 🙏 Acknowledgments

Gondolia is inspired by modern e-commerce platforms like:
- [Medusa.js](https://medusajs.com/) — Provider pattern inspiration
- [Vendure](https://www.vendure.io/) — Plugin architecture insights
- [Saleor](https://saleor.io/) — GraphQL-first approach

---

## 📞 Support & Community

- **Issues**: [GitHub Issues](https://github.com/gondolia/gondolia/issues)
- **Discussions**: [GitHub Discussions](https://github.com/gondolia/gondolia/discussions)
- **Discord**: Coming soon

---

**Made with ❤️ by the Gondolia Community**
