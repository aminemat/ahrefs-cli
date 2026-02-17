# Ahrefs CLI

> AI agent-first command-line interface for the Ahrefs API v3.

[![Tests](https://github.com/aminemat/ahrefs-cli/actions/workflows/test.yml/badge.svg)](https://github.com/aminemat/ahrefs-cli/actions/workflows/test.yml)
[![golangci-lint](https://github.com/aminemat/ahrefs-cli/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/aminemat/ahrefs-cli/actions/workflows/golangci-lint.yml)
[![codecov](https://codecov.io/gh/aminemat/ahrefs-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/aminemat/ahrefs-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/aminemat/ahrefs-cli)](https://goreportcard.com/report/github.com/aminemat/ahrefs-cli)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

---

## 🎯 Project Goals

This CLI was built with **one primary mission**: make the Ahrefs API maximally discoverable and usable for AI coding agents.

**Why AI agents?** Because if an AI can use your CLI without reading docs, humans definitely can too.

### Core Principles

1. **Discoverability First** - Every command is self-documenting
2. **Machine-Readable** - All metadata available as structured JSON
3. **Validation Before Execution** - Dry-run mode for every command
4. **Minimal Dependencies** - Only Cobra + Go stdlib
5. **Structured Errors** - Rich error messages with suggestions

---

## ⚡ Quick Start

### Installation

```bash
# Install with Go
go install github.com/aminemat/ahrefs-cli@latest

# Or clone and build
git clone https://github.com/aminemat/ahrefs-cli
cd ahrefs-cli
make build
```

### Set Your API Key

```bash
# Save to config file (~/.ahrefsrc)
ahrefs config set-key YOUR_API_KEY_HERE

# Or use environment variable
export AHREFS_API_KEY=YOUR_API_KEY_HERE
```

### Your First Query

```bash
# Get domain rating for ahrefs.com (free test query)
ahrefs site-explorer domain-rating --target ahrefs.com --date 2024-01-01

# Output
{
  "status": "success",
  "data": {
    "domain_rating": {
      "domain_rating": 91
    }
  },
  "meta": {
    "response_time_ms": 472
  }
}
```

---

## 🚀 Current Status

### ✅ Implemented & Tested

**Foundation:**
- ✅ HTTP client with Bearer auth
- ✅ Automatic retries with exponential backoff
- ✅ Rate limiting support
- ✅ Config management (`~/.ahrefsrc`)
- ✅ Multiple output formats (JSON, YAML, CSV, Table)
- ✅ **87.7% test coverage** on HTTP client

**Site Explorer Endpoints:**
- ✅ `domain-rating` - Get domain rating
- ✅ `backlinks-stats` - Get backlink statistics
- ✅ `backlinks` - List backlinks (partial)

**AI Agent Features:**
- ✅ `--list-commands` - Full command tree as JSON
- ✅ `--dry-run` - Validate requests without executing
- ✅ `--verbose` - Debug mode
- ✅ Structured error responses with suggestions

### 🔨 In Progress

**Site Explorer (Remaining):**
- ⏳ `refdomains` - Referring domains
- ⏳ `anchors` - Anchor text distribution
- ⏳ `organic-keywords` - Organic keyword rankings
- ⏳ `top-pages` - Top-performing pages
- ⏳ `outlinks-stats` - Outbound link stats

**Other API Categories:**
- ⏳ Keywords Explorer
- ⏳ SERP Overview
- ⏳ Rank Tracker
- ⏳ Site Audit
- ⏳ Brand Radar

**Advanced Features:**
- ⏳ `--describe` flag (JSON schema for each endpoint)
- ⏳ `--list-fields` flag (available fields per endpoint)
- ⏳ Shell completions (bash/zsh/fish)
- ⏳ Pre-built binaries via GitHub Actions
- ⏳ Docker image

---

## 📖 Usage Guide

### For AI Coding Agents

**Step 1: Discover Available Commands**
```bash
ahrefs --list-commands
# Returns complete command tree with all flags and examples as JSON
```

**Step 2: Validate Before Execution**
```bash
ahrefs site-explorer domain-rating \
  --target ahrefs.com \
  --date 2024-01-01 \
  --dry-run

# Output: ✓ Valid request. Would call: GET https://api.ahrefs.com/v3/...
```

**Step 3: Execute & Parse Structured Output**
```bash
ahrefs site-explorer domain-rating \
  --target ahrefs.com \
  --date 2024-01-01 \
  --format json

# Always returns: {"status":"success|error", "data":{...}, "meta":{...}}
```

**Step 4: Handle Errors Programmatically**
```json
{
  "status": "error",
  "error": {
    "code": "AUTH_ERROR",
    "message": "Invalid API key",
    "suggestion": "Run 'ahrefs config set-key <your-key>' to configure",
    "docs_url": "https://docs.ahrefs.com/..."
  }
}
```

### For Humans

```bash
# Get help at any level
ahrefs --help
ahrefs site-explorer --help
ahrefs site-explorer domain-rating --help

# Each command includes detailed examples
ahrefs site-explorer backlinks --help

# Switch output formats
ahrefs site-explorer domain-rating --target ahrefs.com --date 2024-01-01 --format table

# Save output to file
ahrefs site-explorer domain-rating --target ahrefs.com --date 2024-01-01 -o output.json

# Use verbose mode for debugging
ahrefs site-explorer domain-rating --target ahrefs.com --date 2024-01-01 --verbose
```

---

## 🧪 Testing with Free Queries

Ahrefs provides **free test queries** for `ahrefs.com` and `wordcount.com`:

```bash
# These are FREE (no API units consumed)
ahrefs site-explorer domain-rating --target ahrefs.com --date 2024-01-01
ahrefs site-explorer domain-rating --target wordcount.com --date 2024-01-01

# Results are capped at 100 rows
ahrefs site-explorer backlinks-stats --target ahrefs.com --date 2024-01-01
```

**Note:** Full API access requires an Ahrefs Enterprise plan.

---

## 🏗️ Architecture

```
ahrefs-cli/
├── cmd/                      # Command implementations
│   ├── root.go              # Root command + --list-commands
│   ├── config/              # Config management
│   └── siteexplorer/        # Site Explorer endpoints
├── pkg/
│   ├── client/              # HTTP client (87.7% test coverage!)
│   │   ├── client.go
│   │   └── client_test.go
│   ├── models/              # API response structs
│   ├── output/              # Multi-format output (JSON/YAML/CSV/Table)
│   ├── schema/              # JSON schema generator (planned)
│   └── validator/           # Request validation (planned)
├── internal/
│   └── config/              # Config file I/O (~/.ahrefsrc)
├── main.go
├── Makefile                 # make build, test, install
└── README.md
```

**Design Decisions:**
- **Minimal deps:** Only Cobra for CLI framework
- **Stdlib only:** HTTP, JSON, CSV, testing - all native Go
- **No external test deps:** Pure `testing` package with `httptest`
- **Structured errors:** Always return JSON-parseable errors

---

## 🛠️ Development

### Requirements
- Go 1.25+
- Make (optional, for convenience targets)

### Build & Test

```bash
# Build
make build

# Run tests
make test

# Get test coverage
make test-coverage

# Format code
make fmt

# Run all checks
make all
```

### Adding New Endpoints

1. Add model to `pkg/models/`
2. Create command in `cmd/<category>/`
3. Wire up in `main.go`
4. Add tests
5. Update README

**Example:** See `cmd/siteexplorer/siteexplorer.go:newDomainRatingCmd()`

---

## 📊 Project Stats

| Metric | Value |
|--------|-------|
| Go Version | 1.25+ |
| Dependencies | 2 (cobra, pflag) |
| Test Coverage | 87.7% (client) |
| Lines of Code | ~1,500 |
| Endpoints | 3 (more coming!) |
| Output Formats | 4 (JSON, YAML, CSV, Table) |

---

## 🗺️ Roadmap

### v0.2.0 - Site Explorer Complete
- [ ] All remaining Site Explorer endpoints
- [ ] Field selection support (`--select`)
- [ ] Filtering support (`--where`)
- [ ] Pagination support (`--limit`, `--offset`)

### v0.3.0 - Keywords Explorer
- [ ] Keywords overview
- [ ] Volume history
- [ ] Related terms
- [ ] Volume by country

### v0.4.0 - Advanced Features
- [ ] SERP Overview endpoints
- [ ] Rank Tracker endpoints
- [ ] Site Audit endpoints
- [ ] Brand Radar endpoints

### v1.0.0 - Production Ready
- [ ] Full test coverage (>90%)
- [ ] Shell completions
- [ ] Pre-built binaries for all platforms
- [ ] Docker image
- [ ] Comprehensive documentation
- [ ] `--describe` flag (JSON schema introspection)
- [ ] `--list-fields` flag per endpoint

---

## 🤝 Contributing

Contributions are **very welcome**! This is an open-source project built for the community.

### How to Contribute

1. **Fork the repo**
2. **Create a feature branch** (`git checkout -b feature/amazing-endpoint`)
3. **Write tests** for new code
4. **Run tests** (`make test`)
5. **Commit changes** (`git commit -m 'Add amazing endpoint'`)
6. **Push to branch** (`git push origin feature/amazing-endpoint`)
7. **Open a Pull Request**

### What to Contribute

- 🆕 New API endpoints
- 🐛 Bug fixes
- 📝 Documentation improvements
- ✅ More tests
- 🎨 Code quality improvements
- 💡 Feature ideas (open an issue first!)

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **Ahrefs** for providing the API
- **Cobra** for the excellent CLI framework
- **Go team** for the amazing standard library

---

## 🔗 Links

- **Ahrefs API Docs:** https://docs.ahrefs.com/docs/api/reference/introduction
- **API v3 Overview:** https://help.ahrefs.com/en/articles/6559232-about-api-v3-for-enterprise-plan
- **Free Test Queries:** https://docs.ahrefs.com/docs/api/reference/free-test-queries

---

## 💬 Support

- 🐛 **Bug reports:** [Open an issue](https://github.com/aminemat/ahrefs-cli/issues)
- 💡 **Feature requests:** [Open an issue](https://github.com/aminemat/ahrefs-cli/issues)
- 💬 **Questions:** [Start a discussion](https://github.com/aminemat/ahrefs-cli/discussions)

---

<div align="center">

**Built with ❤️ for AI agents and humans alike**

Made by [amine](https://github.com/amine) • [Star this repo](https://github.com/aminemat/ahrefs-cli) if you find it useful!

</div>
