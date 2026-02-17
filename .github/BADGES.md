# Badges Setup Guide

## Current Badges

The README includes the following badges (all FREE):

### 1. **Tests** (GitHub Actions)
```markdown
[![Tests](https://github.com/amine/ahrefs-cli/actions/workflows/test.yml/badge.svg)](https://github.com/amine/ahrefs-cli/actions/workflows/test.yml)
```
- **What it shows:** Pass/fail status of test suite
- **Setup:** Automatic (already configured in `.github/workflows/test.yml`)
- **Updates:** Every push to main/master and every PR

### 2. **golangci-lint** (GitHub Actions)
```markdown
[![golangci-lint](https://github.com/amine/ahrefs-cli/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/amine/ahrefs-cli/actions/workflows/golangci-lint.yml)
```
- **What it shows:** Code quality/linting status
- **Setup:** Automatic (already configured in `.github/workflows/golangci-lint.yml`)
- **Updates:** Every push to main/master and every PR

### 3. **Codecov** (Test Coverage)
```markdown
[![codecov](https://codecov.io/gh/amine/ahrefs-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/amine/ahrefs-cli)
```
- **What it shows:** Test coverage percentage with visual graph
- **Setup Required:**
  1. Go to https://codecov.io/
  2. Sign in with GitHub
  3. Enable `ahrefs-cli` repository
  4. No token needed for public repos!
- **Updates:** Automatic via GitHub Actions after every test run

### 4. **Go Report Card**
```markdown
[![Go Report Card](https://goreportcard.com/badge/github.com/amine/ahrefs-cli)](https://goreportcard.com/report/github.com/amine/ahrefs-cli)
```
- **What it shows:** Code quality grade (A+ to F) based on multiple metrics
- **Setup:**
  1. Go to https://goreportcard.com/
  2. Enter your repo URL: `github.com/amine/ahrefs-cli`
  3. Click "Generate Report"
  4. Badge updates automatically every 24 hours
- **Checks:** gofmt, go vet, gocyclo, golint, ineffassign, license, misspell

### 5. **Go Version**
```markdown
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
```
- **What it shows:** Minimum Go version required
- **Setup:** Already added (static badge)
- **Updates:** Manual (update when min Go version changes)

### 6. **License**
```markdown
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
```
- **What it shows:** Project license
- **Setup:** Already added (static badge)
- **Updates:** Manual (if license changes)

### 7. **PRs Welcome**
```markdown
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
```
- **What it shows:** Indicates the project welcomes contributions
- **Setup:** Already added (static badge)
- **Updates:** None needed

## Post-Push Setup Steps

### 1. Enable Codecov (2 minutes)
After your first push:
1. Visit https://codecov.io/
2. Click "Sign in with GitHub"
3. Authorize Codecov
4. Click "Add new repository"
5. Find and enable `ahrefs-cli`
6. **Done!** Badge will show coverage after first test run

### 2. Generate Go Report Card (30 seconds)
After your first push:
1. Visit https://goreportcard.com/
2. Enter: `github.com/amine/ahrefs-cli`
3. Click "Generate Report"
4. Wait ~30 seconds for first analysis
5. **Done!** Badge will show your grade

### 3. Wait for GitHub Actions (automatic)
After your first push:
1. GitHub Actions will automatically run
2. Visit https://github.com/amine/ahrefs-cli/actions
3. Watch tests and linting complete
4. **Done!** Badges will update automatically

## What Happens After Push

### Immediately:
- ✅ GitHub Actions trigger (test + lint workflows)
- ✅ Tests run with coverage collection
- ✅ Coverage uploaded to Codecov

### Within 1 minute:
- ✅ Test badge updates (pass/fail)
- ✅ Lint badge updates (pass/fail)
- ✅ Codecov badge shows coverage %

### Within 24 hours:
- ✅ Go Report Card updates automatically
- ✅ Shows code quality grade

## Badge Monitoring

All badges are public and visible on your README. They'll automatically update:

- **Tests/Lint:** Every push, every PR
- **Coverage:** Every push, every PR
- **Go Report Card:** Daily (or on-demand)

## Optional: Additional Badges

### GitHub Stars
```markdown
[![GitHub Stars](https://img.shields.io/github/stars/amine/ahrefs-cli?style=social)](https://github.com/amine/ahrefs-cli)
```

### GitHub Issues
```markdown
[![GitHub Issues](https://img.shields.io/github/issues/amine/ahrefs-cli)](https://github.com/amine/ahrefs-cli/issues)
```

### Last Commit
```markdown
[![Last Commit](https://img.shields.io/github/last-commit/amine/ahrefs-cli)](https://github.com/amine/ahrefs-cli/commits)
```

### Release Version
```markdown
[![Release](https://img.shields.io/github/v/release/amine/ahrefs-cli)](https://github.com/amine/ahrefs-cli/releases)
```

## Troubleshooting

### Badge shows "unknown" or "invalid"
- Wait 1-2 minutes for first run to complete
- Check GitHub Actions are enabled
- Verify branch name (main vs master)

### Codecov badge not showing
1. Check https://app.codecov.io/gh/amine/ahrefs-cli
2. Verify repository is enabled
3. Check GitHub Actions logs for upload errors
4. For public repos, no token needed!

### Go Report Card not updating
- Visit https://goreportcard.com/report/github.com/amine/ahrefs-cli
- Click "Refresh" to force update
- Card updates daily automatically

## Cost

**All badges are 100% FREE for public repositories!**

- GitHub Actions: 2000 minutes/month free
- Codecov: Unlimited for open source
- Go Report Card: Free forever
- Shields.io badges: Free forever

## Summary

After pushing to GitHub:
1. ✅ 4 badges work immediately (Tests, Lint, Go Version, License, PRs Welcome)
2. 🔧 Enable Codecov (2 min)
3. 🔧 Generate Go Report Card (30 sec)
4. ✅ All badges auto-update forever

**Total setup time: ~3 minutes** 🚀
