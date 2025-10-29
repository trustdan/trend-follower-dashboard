# Phase 0 - Step 1: Development Environment Setup

**Phase:** 0 - Foundation & Proof-of-Concept
**Step:** 1 of 4
**Duration:** 1 day
**Dependencies:** None (starting point)

---

## Objectives

Set up a complete development environment on Linux (WSL2/Kali) that enables:
1. Go backend development and testing
2. Cross-compilation to Windows .exe
3. Node.js/Svelte frontend development
4. Comprehensive logging infrastructure
5. Reproducible builds

---

## Prerequisites

- Linux environment (WSL2/Kali or native)
- Internet access for downloading tools
- Windows machine available for testing cross-compiled binaries

---

## Step-by-Step Instructions

### 1. Verify Linux Environment

```bash
# Check Linux distribution
cat /etc/os-release

# Check if running in WSL
uname -a | grep -i microsoft

# Note your username and home directory
whoami
echo $HOME
```

**Expected outcome:**
- Kali Linux or similar Debian-based distribution
- WSL2 kernel detected (if using WSL)
- Home directory path noted for future reference

---

### 2. Install Go 1.24+

```bash
# Check if Go is already installed
go version

# If Go is not installed or version < 1.24:
# Download Go 1.24.2 (or latest 1.24.x)
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz

# Remove old Go installation (if exists)
sudo rm -rf /usr/local/go

# Extract new Go
sudo tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz

# Add Go to PATH (if not already in ~/.bashrc or ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
# Expected: go version go1.24.2 linux/amd64
```

**Verification:**
```bash
go version
# Should output: go version go1.24.2 linux/amd64 (or similar)

go env GOPATH
# Should output: /home/yourusername/go

go env GOOS
# Should output: linux

go env GOARCH
# Should output: amd64
```

---

### 3. Verify Backend Compiles and Tests Pass

```bash
# Navigate to backend directory
cd /home/kali/fresh-start-trading-platform/backend

# Verify go.mod exists
ls -la go.mod

# Download dependencies
go mod download

# Run all tests (comprehensive test suite)
go test ./... -v

# Expected: All tests pass
# Look for: PASS indicators, no FAIL

# Build the backend binary
go build -o tf-engine cmd/tf-engine/main.go

# Verify binary was created
ls -lh tf-engine

# Run a simple command to verify it works
./tf-engine --help

# Expected: Help text with available commands
```

**If tests fail:**
- Review error messages carefully
- Check if SQLite dependencies are missing (rarely needed for pure Go)
- Ensure Go version is 1.24+
- Report issues in documentation

---

### 4. Test Cross-Compilation to Windows

This is critical for the entire project workflow.

```bash
# Navigate to backend
cd /home/kali/fresh-start-trading-platform/backend

# Cross-compile to Windows .exe (pure Go, no cgo)
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" -o tf-engine.exe cmd/tf-engine/main.go

# Verify Windows binary was created
ls -lh tf-engine.exe
file tf-engine.exe

# Expected file output:
# tf-engine.exe: PE32+ executable (console) x86-64 (stripped to external PDB), for MS Windows
```

**Test on Windows (if available now):**
1. Copy `tf-engine.exe` to Windows machine
2. Open PowerShell or CMD in the directory
3. Run: `.\tf-engine.exe --help`
4. Should display help text without errors

**If Windows testing not available yet:**
- Document that cross-compilation succeeded
- Plan Windows testing for later in Phase 0

---

### 5. Install Node.js 20+

```bash
# Check if Node.js is installed
node --version
npm --version

# If Node.js is not installed or version < 20:
# Option A: Using NodeSource repository (recommended)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Option B: Using nvm (Node Version Manager)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 20
nvm use 20

# Verify installation
node --version
# Expected: v20.x.x

npm --version
# Expected: 10.x.x or similar
```

**Verification:**
```bash
node --version
# Should output: v20.11.0 (or similar 20.x.x)

npm --version
# Should output: 10.2.4 (or similar)

# Test npm works
npm --help
```

---

### 6. Set Up Logging Infrastructure

```bash
# Navigate to project root
cd /home/kali/fresh-start-trading-platform

# Create logs directory
mkdir -p logs

# Verify directory creation
ls -ld logs

# Create a .gitignore for logs (even though we're not using Git in Linux)
# This prepares for Windows Git repo
echo "logs/*.log" >> .gitignore
echo "logs/*.txt" >> .gitignore

# Test write permissions
touch logs/test.log
rm logs/test.log

# Expected: No permission errors
```

**Create logging configuration:**

Create `backend/internal/logx/config.go`:
```go
package logx

import (
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

// LogLevel defines logging verbosity
type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
)

// Config holds logger configuration
type Config struct {
    Level      LogLevel
    LogToFile  bool
    LogToConsole bool
    LogDir     string
    FileName   string
}

// DefaultConfig returns default logging configuration
func DefaultConfig() *Config {
    return &Config{
        Level:        INFO,
        LogToFile:    true,
        LogToConsole: true,
        LogDir:       "logs",
        FileName:     "tf-engine.log",
    }
}

// NewLogger creates a logger with the given configuration
func NewLogger(cfg *Config) (*log.Logger, error) {
    var writers []io.Writer

    if cfg.LogToConsole {
        writers = append(writers, os.Stdout)
    }

    if cfg.LogToFile {
        if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
            return nil, err
        }

        logPath := filepath.Join(cfg.LogDir, cfg.FileName)
        f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return nil, err
        }

        writers = append(writers, f)
    }

    return log.New(io.MultiWriter(writers...), "", log.LstdFlags|log.Lshortfile), nil
}
```

**Verify logging directory structure:**
```bash
tree logs/
# Expected: logs/ directory exists (empty for now)
```

---

### 7. Configure VSCode/Cursor Workspace

**Install VSCode or Cursor (if not already installed):**
```bash
# For Cursor (recommended for AI-assisted coding)
# Download from: https://cursor.sh/

# For VSCode
# Download from: https://code.visualstudio.com/
```

**Install Go Extension:**
1. Open VSCode/Cursor
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Go" (by Go Team at Google)
4. Click Install

**Install Svelte Extension:**
1. Search for "Svelte for VS Code" (by Svelte)
2. Click Install

**Configure workspace settings:**

Create `.vscode/settings.json`:
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "goimports",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  },
  "go.testFlags": ["-v"],
  "go.testTimeout": "30s",
  "svelte.enable-ts-plugin": true,
  "files.eol": "\n",
  "files.insertFinalNewline": true,
  "files.trimTrailingWhitespace": true
}
```

Create `.vscode/launch.json` for debugging:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch tf-engine",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/cmd/tf-engine",
      "args": []
    },
    {
      "name": "Test Current File",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${file}"
    }
  ]
}
```

---

### 8. Verify Backend Test Suite

Run the comprehensive backend test suite to ensure everything works:

```bash
cd /home/kali/fresh-start-trading-platform/backend

# Run all tests with verbose output
go test ./... -v

# Run tests with coverage report
go test ./... -cover

# Run specific domain tests
go test ./internal/domain/... -v

# Run storage tests
go test ./internal/storage/... -v

# Expected: All tests PASS
```

**Critical tests to verify:**
- Position sizing calculations (Van Tharp method)
- Checklist evaluation (RED/YELLOW/GREEN logic)
- Heat management (portfolio and bucket caps)
- 5 gates validation
- Database operations (CRUD for positions, decisions, candidates)

**If any tests fail:**
- Document which tests failed
- Check error messages
- Verify SQLite is working (should be built into Go)
- Do NOT proceed until tests pass

---

### 9. Document Environment Setup

Create `docs/dev-environment.md`:

```markdown
# Development Environment

**Last Updated:** 2025-10-29

## System Information

- **OS:** Kali Linux on WSL2 (or native Linux)
- **Kernel:** Linux 5.15.167.4-microsoft-standard-WSL2
- **User:** [your username]
- **Home:** /home/[your username]

## Installed Tools

### Go
- **Version:** go1.24.2 linux/amd64
- **Location:** /usr/local/go
- **GOPATH:** /home/[username]/go
- **Verification:** `go version`

### Node.js
- **Version:** v20.11.0
- **NPM Version:** 10.2.4
- **Location:** /usr/bin/node
- **Verification:** `node --version`

### Backend Status
- **Tests Passing:** ✓ All tests pass
- **Binary Builds:** ✓ Linux binary builds successfully
- **Cross-Compile:** ✓ Windows .exe builds successfully

### Directories
- **Project Root:** /home/kali/fresh-start-trading-platform
- **Backend:** /home/kali/fresh-start-trading-platform/backend
- **Logs:** /home/kali/fresh-start-trading-platform/logs
- **Plans:** /home/kali/fresh-start-trading-platform/plans

## Quick Commands

### Backend Development
```bash
# Build Linux binary
cd backend && go build -o tf-engine cmd/tf-engine/main.go

# Build Windows binary
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o tf-engine.exe cmd/tf-engine/main.go

# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover
```

### Cross-Compilation
```bash
# Windows (amd64)
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o app.exe ./cmd/app

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

# macOS (amd64)
GOOS=darwin GOARCH=amd64 go build -o app-mac ./cmd/app

# macOS (arm64 / M1/M2)
GOOS=darwin GOARCH=arm64 go build -o app-mac-arm ./cmd/app
```

## Known Issues

[Document any issues encountered during setup]

## Next Steps

- Proceed to Phase 0 Step 2: Fyne POC
```

---

## Verification Checklist

Before proceeding to Step 2, verify:

- [ ] Go 1.24+ installed and `go version` works
- [ ] Backend compiles successfully (`go build` produces binary)
- [ ] All backend tests pass (`go test ./... -v`)
- [ ] Cross-compilation to Windows .exe succeeds
- [ ] Windows binary is PE32+ executable (check with `file` command)
- [ ] Node.js 20+ installed and `node --version` works
- [ ] NPM installed and `npm --version` works
- [ ] Logs directory created (`logs/`)
- [ ] VSCode/Cursor configured with Go and Svelte extensions
- [ ] `docs/dev-environment.md` created with current setup details
- [ ] No critical errors or blockers

---

## Expected Outputs

After completing this step, you should have:

1. **Go Environment:**
   - `go version` outputs go1.24.2 (or newer)
   - `backend/tf-engine` (Linux binary)
   - `backend/tf-engine.exe` (Windows binary)
   - All tests passing

2. **Node Environment:**
   - `node --version` outputs v20.x.x
   - `npm --version` outputs 10.x.x

3. **Project Structure:**
   - `logs/` directory created
   - `.vscode/settings.json` configured
   - `.vscode/launch.json` for debugging

4. **Documentation:**
   - `docs/dev-environment.md` with current setup details

---

## Troubleshooting

### Go Installation Issues

**Problem:** `go version` not found after installation
**Solution:**
```bash
# Ensure PATH is set
echo $PATH | grep "/usr/local/go/bin"

# If not present, add to ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**Problem:** Tests fail with "cannot find package"
**Solution:**
```bash
cd backend
go mod download
go mod tidy
go test ./... -v
```

### Cross-Compilation Issues

**Problem:** Windows binary builds but doesn't run on Windows
**Solution:**
- Ensure `CGO_ENABLED=0` is set (pure Go only)
- Use `-trimpath -ldflags "-s -w"` flags
- Verify binary with `file tf-engine.exe` shows "PE32+"
- Test on actual Windows machine (not Wine)

**Problem:** "cgo" errors during cross-compilation
**Solution:**
- Set `CGO_ENABLED=0` explicitly
- Remove any dependencies that require cgo
- Check `go.mod` for cgo-dependent packages

### Node.js Issues

**Problem:** Node.js version is too old
**Solution:**
```bash
# Use nvm to install specific version
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 20
nvm use 20
```

---

## Time Estimate

- **Go Installation:** 15-30 minutes
- **Backend Verification:** 15 minutes
- **Cross-Compilation Testing:** 15 minutes
- **Node.js Installation:** 10-15 minutes
- **Logging Setup:** 10 minutes
- **VSCode Configuration:** 10 minutes
- **Documentation:** 15 minutes

**Total:** ~1.5-2 hours (allowing for troubleshooting)

---

## References

- [Go Installation](https://go.dev/doc/install)
- [Go Cross-Compilation](https://go.dev/wiki/WindowsCrossCompiling)
- [Node.js Downloads](https://nodejs.org/)
- [VSCode Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
- [overview-plan.md - System Architecture](../plans/overview-plan.md#system-architecture)
- [overview-plan.md - Technology Stack](../plans/overview-plan.md#technology-stack)
- [RULES.md - Section 5.3](../1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md)

---

## Next Step

Proceed to: **[Phase 0 - Step 2: Fyne Proof-of-Concept](phase0-step2-fyne-poc.md)**

---

**Status:** 📋 Ready for Execution
**Last Updated:** 2025-10-29
