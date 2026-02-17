# Contributing to Ahrefs CLI

First off, thank you for considering contributing to Ahrefs CLI! 🎉

This project was built with AI agents and the open-source community in mind. Your contributions help make the Ahrefs API more accessible to everyone.

## 🚀 Quick Start for Contributors

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/ahrefs-cli`
3. Create a branch: `git checkout -b feature/amazing-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit: `git commit -m 'Add amazing feature'`
7. Push: `git push origin feature/amazing-feature`
8. Open a Pull Request

## 🎯 What We're Looking For

### High Priority
- **New API endpoints** - Help us expand coverage of the Ahrefs API v3
- **Bug fixes** - Found a bug? Fix it and send a PR!
- **Tests** - More test coverage is always welcome
- **Documentation** - Improve examples, fix typos, add clarity

### Medium Priority
- **Performance improvements** - Make it faster!
- **Code quality** - Refactoring, better error handling, etc.
- **CI/CD** - Help us automate builds and releases

### Nice to Have
- **Shell completions** - bash, zsh, fish support
- **Pre-built binaries** - Multi-platform releases
- **Docker image** - Containerized version

## 🏗️ Development Setup

### Requirements
- Go 1.25 or higher
- Make (optional, but recommended)

### Setup Steps

```bash
# Clone the repo
git clone https://github.com/amine/ahrefs-cli
cd ahrefs-cli

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Try it out
./ahrefs --help
```

## ✅ Code Standards

### Go Code Style
- Follow standard Go conventions (use `gofmt` or `make fmt`)
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions small and focused

### Testing
- Write tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Test both success and error cases

### Commits
- Use clear, descriptive commit messages
- Start with a verb: "Add", "Fix", "Update", "Refactor"
- Reference issue numbers when applicable

**Good:**
```
Add refdomains endpoint to Site Explorer
Fix retry logic for 429 rate limit errors
Update README with new usage examples
```

**Bad:**
```
changes
fix stuff
wip
```

## 📝 Adding a New Endpoint

Here's the typical workflow for adding a new API endpoint:

### 1. Add Response Models

Create structs in `pkg/models/`:

```go
// pkg/models/siteexplorer.go
type RefDomainsResponse struct {
    RefDomains []RefDomain `json:"refdomains"`
}

type RefDomain struct {
    Domain       string  `json:"domain"`
    DomainRating float64 `json:"domain_rating"`
    Backlinks    int     `json:"backlinks"`
}
```

### 2. Create Command

Add command to appropriate category in `cmd/`:

```go
// cmd/siteexplorer/refdomains.go
func newRefDomainsCmd() *cobra.Command {
    var target, mode string
    var limit int

    cmd := &cobra.Command{
        Use:   "refdomains",
        Short: "List referring domains",
        Long:  "Get a list of domains that link to the target.",
        Example: `  # Get referring domains
  ahrefs site-explorer refdomains --target example.com`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return runRefDomains(target, mode, limit)
        },
    }

    cmd.Flags().StringVar(&target, "target", "", "Target domain (required)")
    cmd.Flags().StringVar(&mode, "mode", "domain", "Mode: exact, domain, prefix, subdomains")
    cmd.Flags().IntVar(&limit, "limit", 100, "Max results")

    cmd.MarkFlagRequired("target")

    return cmd
}
```

### 3. Implement Handler

```go
func runRefDomains(target, mode string, limit int) error {
    flags := cmd.GetGlobalFlags()

    // Get API key
    apiKey := flags.APIKey
    if apiKey == "" {
        apiKey = config.GetAPIKey()
    }

    // Create client
    c := client.NewClient(client.Config{APIKey: apiKey})

    // Build params
    params := url.Values{}
    params.Set("target", target)
    params.Set("mode", mode)
    params.Set("limit", fmt.Sprintf("%d", limit))

    // Dry run support
    if flags.DryRun {
        fmt.Printf("✓ Valid request. Would call: GET %s/site-explorer/refdomains?%s\n",
            client.BaseURL, params.Encode())
        return nil
    }

    // Make request
    resp, err := c.Get(context.Background(), "/site-explorer/refdomains", params)
    if err != nil {
        w, _ := output.NewWriter(flags.OutputFormat, flags.OutputFile)
        w.WriteError(err)
        return err
    }

    // Parse response
    var result models.RefDomainsResponse
    if err := json.Unmarshal(resp.Body, &result); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }

    // Output
    w, err := output.NewWriter(flags.OutputFormat, flags.OutputFile)
    if err != nil {
        return err
    }
    defer w.Close()

    return w.WriteSuccess(result, &resp.Meta)
}
```

### 4. Wire It Up

Add to parent command in the same file:

```go
func NewSiteExplorerCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "site-explorer",
        Short: "Site Explorer API endpoints",
    }

    cmd.AddCommand(newDomainRatingCmd())
    cmd.AddCommand(newBacklinksCmd())
    cmd.AddCommand(newRefDomainsCmd()) // <-- Add this

    return cmd
}
```

### 5. Write Tests

Add tests to `pkg/client/client_test.go` or create endpoint-specific tests:

```go
func TestRefDomains(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"refdomains":[{"domain":"example.com","domain_rating":50}]}`))
    }))
    defer server.Close()

    c := client.NewClient(client.Config{
        APIKey:  "test",
        BaseURL: server.URL,
    })

    resp, err := c.Get(context.Background(), "/site-explorer/refdomains", url.Values{
        "target": []string{"test.com"},
    })

    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if resp.StatusCode != 200 {
        t.Errorf("Expected 200, got %d", resp.StatusCode)
    }
}
```

### 6. Update Documentation

- Add endpoint to README.md
- Add usage example
- Update roadmap status

## 🐛 Reporting Bugs

**Before submitting a bug report:**
- Check existing issues to avoid duplicates
- Verify you're using the latest version
- Test with free test queries (`ahrefs.com`, `wordcount.com`) if possible

**When submitting a bug report, include:**
- Go version (`go version`)
- CLI version (`ahrefs --version`)
- Command you ran (sanitize your API key!)
- Expected behavior
- Actual behavior
- Error messages or output

## 💡 Feature Requests

We love feature requests! Open an issue with:
- Clear description of the feature
- Why it's useful
- Example usage
- Relevant API documentation links

## 📜 Code of Conduct

- Be respectful and inclusive
- Assume good intentions
- Provide constructive feedback
- Help others when you can
- Keep discussions on-topic

## 🙏 Recognition

Contributors will be recognized in:
- Release notes
- Contributors section (coming soon)
- Commit history

## 📞 Questions?

- Open a [Discussion](https://github.com/amine/ahrefs-cli/discussions)
- Check existing [Issues](https://github.com/amine/ahrefs-cli/issues)
- Review the [README](README.md)

---

**Thank you for contributing!** 🎉

Every contribution, no matter how small, makes this project better for the community.
