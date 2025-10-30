# Development Environment

**Last Updated:** 2025-10-29

## System Information

- **OS:** Kali GNU/Linux Rolling 2025.2
- **Kernel:** Linux 5.15.167.4-microsoft-standard-WSL2
- **User:** root
- **Home:** /root

## Installed Tools

### Go

- **Version:** go1.24.2 linux/amd64
- **Location:** /usr/local/go
- **GOPATH:** /root/go
- **Verification:** `go version`

### Node.js

- **Version:** v20.19.0
- **NPM Version:** 9.2.0
- **Location:** /usr/bin/node
- **Verification:** `node --version`

### Backend Status

- **Tests Passing:** ✓ All tests pass (96.9% coverage on domain logic)
- **Binary Builds:** ✓ Linux binary builds successfully (18MB)
- **Cross-Compile:** ✓ Windows .exe builds successfully (11MB PE32+)

### Directories

- **Project Root:** /home/kali/fresh-start-trading-platform
- **Backend:** /home/kali/fresh-start-trading-platform/backend
- **Logs:** /home/kali/fresh-start-trading-platform/logs
- **Plans:** /home/kali/fresh-start-trading-platform/plans
- **Docs:** /home/kali/fresh-start-trading-platform/docs

## Quick Commands

### Backend Development

**Linux/macOS (bash):**
```bash
# Navigate to backend
cd /home/kali/fresh-start-trading-platform/backend

# Build Linux binary
go build -o tf-engine cmd/tf-engine/main.go

# Build Windows binary
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o tf-engine.exe cmd/tf-engine/main.go

# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific domain tests
go test ./internal/domain/... -v

# Run storage tests
go test ./internal/storage/... -v
```

**Windows (PowerShell):**
```powershell
# Navigate to backend
cd C:\Users\Dan\trend-follower-dashboard\backend

# Build Windows binary
go build -o tf-engine.exe cmd/tf-engine/main.go

# Build for Linux (cross-compile)
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -trimpath -ldflags "-s -w" -o tf-engine cmd/tf-engine/main.go

# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific domain tests
go test ./internal/domain/... -v

# Run storage tests
go test ./internal/storage/... -v
```

### Database Operations

```bash
# Navigate to backend
cd /home/kali/fresh-start-trading-platform/backend

# Initialize database
./tf-engine init

# Configure settings
./tf-engine set-setting --key equity --value 100000
./tf-engine set-setting --key risk_pct --value 0.75
./tf-engine set-setting --key portfolio_heat_cap --value 4.0
./tf-engine set-setting --key bucket_heat_cap --value 1.5

# Get all settings
./tf-engine get-settings
```

### Cross-Compilation

**From Linux/macOS (bash):**
```bash
# Windows (amd64)
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o tf-engine.exe cmd/tf-engine/main.go

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o tf-engine cmd/tf-engine/main.go

# macOS (amd64)
GOOS=darwin GOARCH=amd64 go build -o tf-engine-mac cmd/tf-engine/main.go

# macOS (arm64 / M1/M2)
GOOS=darwin GOARCH=arm64 go build -o tf-engine-mac-arm cmd/tf-engine/main.go
```

**From Windows (PowerShell):**
```powershell
# Windows (native - no env vars needed)
go build -trimpath -ldflags "-s -w" -o tf-engine.exe cmd/tf-engine/main.go

# Linux (amd64)
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o tf-engine cmd/tf-engine/main.go

# macOS (amd64)
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o tf-engine-mac cmd/tf-engine/main.go

# macOS (arm64 / M1/M2)
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o tf-engine-mac-arm cmd/tf-engine/main.go
```

## Testing Strategy

### Run Full Test Suite

```bash
cd /home/kali/fresh-start-trading-platform/backend
go test ./... -v
```

### Run Specific Tests

```bash
# Test position sizing
go test ./internal/domain/ -v -run TestCalculateSize

# Test checklist validation
go test ./internal/domain/ -v -run TestEvaluateChecklist

# Test heat management
go test ./internal/domain/ -v -run TestCalculateHeat

# Test 5 gates
go test ./internal/domain/ -v -run TestValidateHardGates

# Test database operations
go test ./internal/storage/ -v
```

### Coverage Report

```bash
# Overall coverage
go test ./... -cover

# Detailed coverage HTML report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## IDE Configuration

### VSCode/Cursor Extensions

**Required:**
- Go (by Go Team at Google)
- Svelte for VS Code (by Svelte)

**Recommended:**
- GitLens
- Error Lens
- Better Comments

### Workspace Settings

Configuration files created at `.vscode/`:
- `settings.json` - Go formatting, linting, and editor settings
- `launch.json` - Debug configurations for tf-engine and tests

### Debugging

**Launch tf-engine in debug mode:**
1. Open VSCode/Cursor
2. Press F5 or select "Launch tf-engine" from Run menu
3. Set breakpoints in code
4. Step through execution

**Debug tests:**
1. Open test file
2. Press F5 or select "Test Current File"
3. Breakpoints in test functions will be hit

## Known Issues

None at this time. All setup steps completed successfully.

## Environment Variables

### Optional Environment Variables

```bash
# Enable debug logging to stderr (in addition to file)
export TF_DEBUG=1

# Set log level (debug, info, warn, error)
export LOG_LEVEL=debug

# Custom database location
export DB_PATH=/path/to/custom/trading.db
```

## Verification Checklist

✅ Go 1.24.2 installed and `go version` works
✅ Backend compiles successfully (`go build` produces binary)
✅ All backend tests pass (`go test ./... -v`)
✅ Cross-compilation to Windows .exe succeeds
✅ Windows binary is PE32+ executable (verified with `file` command)
✅ Node.js 20.19.0 installed and `node --version` works
✅ NPM 9.2.0 installed and `npm --version` works
✅ Logs directory created (`logs/`)
✅ VSCode/Cursor configured with settings.json and launch.json
✅ No critical errors or blockers

## Test Coverage Summary

As of 2025-10-29:

| Package | Coverage | Notes |
|---------|----------|-------|
| internal/domain | 96.9% | Core business logic (position sizing, checklist, heat, gates) |
| internal/storage | 77.1% | Database operations (SQLite CRUD) |
| internal/logx | 73.3% | Logging infrastructure |
| internal/scrape | 42.1% | FINVIZ web scraping |
| cmd/tf-engine | 0.0% | No unit tests (integration tested) |
| internal/cli | 0.0% | No unit tests (integration tested) |
| internal/server | 0.0% | No unit tests (legacy, minimal use) |

**Critical tests verified:**
- ✅ Position sizing calculations (Van Tharp method)
- ✅ Checklist evaluation (RED/YELLOW/GREEN logic)
- ✅ Heat management (portfolio and bucket caps)
- ✅ 5 gates validation
- ✅ Database operations (CRUD for positions, decisions, candidates)
- ✅ Cooldown management
- ✅ Impulse brake timer

## Next Steps

✅ **Phase 0 Step 1: Development Environment Setup** - COMPLETE

**Ready to proceed to:**
- Phase 0 Step 2: Fyne Proof-of-Concept (GUI framework evaluation)

See: `plans/phase0-step2-fyne-poc.md`

## References

- [Go Installation](https://go.dev/doc/install)
- [Go Cross-Compilation](https://go.dev/wiki/WindowsCrossCompiling)
- [Node.js Downloads](https://nodejs.org/)
- [VSCode Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
- [Project Overview](../README.md)
- [CLAUDE.md - Development Guidelines](../CLAUDE.md)
- [Anti-Impulsivity Design](./anti-impulsivity.md)
- [Development Philosophy](./dev/DEVELOPMENT_PHILOSOPHY.md)

---

**Status:** ✅ Complete
**Date:** 2025-10-29
**Duration:** ~1.5 hours (including documentation)
