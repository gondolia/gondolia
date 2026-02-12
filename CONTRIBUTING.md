# Contributing to Gondolia

Thank you for your interest in contributing to Gondolia! We welcome contributions from the community.

## 📜 Contributor License Agreement (CLA)

**Before we can accept your contributions, you need to sign our Contributor License Agreement (CLA).**

### Why a CLA?

The CLA ensures that:

1. **The project remains open source** under AGPL v3
2. **You retain copyright** to your contributions
3. **You grant Gondolia the right** to use and distribute your contributions
4. **Commercial use is protected** — companies can build on Gondolia without legal uncertainty
5. **Patent protection** — contributors grant a patent license to prevent patent trolling

### How to Sign the CLA

We use [CLA Assistant](https://cla-assistant.io/) for automated CLA signing:

1. Open a pull request
2. The CLA bot will comment with a link to sign
3. Sign electronically (takes ~1 minute)
4. Your PR will be unblocked

**One-time process** — you only need to sign once for all future contributions.

### CLA Summary

By signing, you confirm that:

- ✅ You own the copyright to your contribution OR have permission to contribute it
- ✅ You grant Gondolia a perpetual, worldwide, non-exclusive, royalty-free license to use your contribution under AGPL v3
- ✅ You grant a patent license for any patents covering your contribution
- ✅ Your contribution does not violate any third-party rights

**Full CLA text**: [CLA.md](./CLA.md) *(coming soon)*

---

## 🛠️ Development Setup

### Prerequisites

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose
- K3d (for local Kubernetes)
- PostgreSQL 15 (via Docker)
- Redis (via Docker)

### Local Setup

1. **Clone the repository**

```bash
git clone https://github.com/gondolia/gondolia.git
cd gondolia
```

2. **Install Go dependencies**

```bash
go mod download
```

3. **Start infrastructure**

```bash
docker-compose up -d
```

4. **Run tests**

```bash
make test
```

5. **Run a service locally**

```bash
cd services/identity
make run
```

6. **Run frontend**

```bash
cd frontend
npm install
npm run dev
```

---

## 🧪 Testing

We use standard Go testing:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

**Testing Guidelines:**

- ✅ Write unit tests for all service logic
- ✅ Use interface mocks for repositories
- ✅ Test error paths, not just happy paths
- ✅ Aim for >80% coverage on critical paths

---

## 📝 Code Style

### Go Code Style

We follow standard Go conventions:

- Use `gofmt` (already handled by `go fmt`)
- Use `golangci-lint` for linting
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Follow [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

**Run linter:**

```bash
golangci-lint run
```

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring (no behavior change)
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (dependencies, build scripts)
- `perf`: Performance improvements

**Examples:**

```
feat(identity): add SSO provider interface

Implement AuthProvider interface for pluggable SSO integrations.
Includes Azure AD example implementation.

Closes #42
```

```
fix(catalog): handle null prices gracefully

Previously crashed on products without prices.
Now returns 0 with a warning log.

Fixes #123
```

---

## 🔀 Pull Request Process

1. **Fork the repository** and create a feature branch

```bash
git checkout -b feat/my-awesome-feature
```

2. **Make your changes**

- Write clean, idiomatic Go code
- Add tests for new functionality
- Update documentation if needed
- Run `go fmt` and `golangci-lint`

3. **Commit with conventional commits**

```bash
git commit -m "feat(catalog): add bulk import API"
```

4. **Push to your fork**

```bash
git push origin feat/my-awesome-feature
```

5. **Open a Pull Request**

- Describe what your PR does
- Reference related issues (e.g., "Closes #42")
- Sign the CLA when prompted

6. **Code Review**

- Maintainers will review your PR
- Address feedback in new commits
- Once approved, we'll merge it!

---

## 📦 Writing a Provider

Providers are the core extensibility mechanism in Gondolia. Here's how to write one:

### 1. Implement the Interface

```go
// github.com/yourname/gondolia-stripe/provider.go
package stripe

import (
    "context"
    "github.com/gondolia/gondolia/provider"
    "github.com/gondolia/gondolia/provider/payment"
)

type Provider struct {
    secretKey string
}

func NewProvider(config map[string]any) (payment.PaymentProvider, error) {
    secretKey, _ := config["secret_key"].(string)
    return &Provider{secretKey: secretKey}, nil
}

func (p *Provider) Initialize(ctx context.Context, req payment.InitializeRequest) (*payment.PaymentSession, error) {
    // Implementation
}

// ... implement other methods
```

### 2. Register via `init()`

```go
func init() {
    provider.Register[payment.PaymentProvider]("payment", "stripe",
        provider.Metadata{
            Name:        "stripe",
            DisplayName: "Stripe",
            Category:    "payment",
            Version:     "1.0.0",
            Description: "Stripe payment gateway",
            ConfigSpec: []provider.ConfigField{
                {Key: "secret_key", Type: "secret", Required: true},
            },
        },
        NewProvider,
    )
}
```

### 3. Publish as a Go Module

```bash
git tag v1.0.0
git push origin v1.0.0
```

### 4. Users Import Your Provider

```go
import _ "github.com/yourname/gondolia-stripe"
```

**See [`docs/architecture/adr-001-provider-pattern.md`](./docs/architecture/adr-001-provider-pattern.md) for the full guide.**

---

## 🐛 Reporting Bugs

Found a bug? Please [open an issue](https://github.com/gondolia/gondolia/issues/new) with:

- **Description**: What went wrong?
- **Steps to reproduce**: How can we trigger the bug?
- **Expected behavior**: What should have happened?
- **Actual behavior**: What actually happened?
- **Environment**: OS, Go version, service version
- **Logs**: Relevant error messages or stack traces

---

## 💡 Requesting Features

Have an idea? [Open a discussion](https://github.com/gondolia/gondolia/discussions/new) first!

We prefer discussing features before implementation to ensure they align with the project's direction.

---

## 🙏 Thank You!

Every contribution makes Gondolia better. Whether it's code, documentation, bug reports, or feedback — thank you for being part of the community! ❤️

---

## 📞 Questions?

- **GitHub Discussions**: [github.com/gondolia/gondolia/discussions](https://github.com/gondolia/gondolia/discussions)
- **Discord**: Coming soon

---

**Happy coding!** 🚀
