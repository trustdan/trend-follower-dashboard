# TF-Engine GUI Build Instructions

Quick reference for rebuilding the TF-Engine GUI application.

---

## Quick Start

### Windows (PowerShell)

**Using the build script (easiest):**
```powershell
cd ui
.\build.bat           # Creates tf-gui.exe
.\build.bat v10       # Creates tf-gui-v10.exe
```

**Manual build:**
```powershell
cd ui
go build -o tf-gui.exe *.go
```

### Linux / macOS (Bash)

**Using the build script:**
```bash
cd ui
./build.sh           # Creates tf-gui.exe
./build.sh v10       # Creates tf-gui-v10.exe
```

**Manual build:**
```bash
cd ui
go build -o tf-gui *.go
```

---

## Build Scripts

### build.bat (Windows)

**Features:**
- Builds with version number or default name
- Shows file size after build
- Color-coded success/failure messages

**Usage:**
```powershell
# Default build
.\build.bat

# Versioned build
.\build.bat v10
.\build.bat beta
.\build.bat test
```

### build.sh (Linux/macOS)

**Features:**
- Same as build.bat but for Unix systems
- Requires execute permission

**Setup:**
```bash
chmod +x build.sh
```

**Usage:**
```bash
# Default build
./build.sh

# Versioned build
./build.sh v10
./build.sh beta
./build.sh test
```

---

## Manual Build Commands

### Build Current Version
```powershell
cd ui
go build -o tf-gui-v9.exe *.go
```

### Build Next Version
```powershell
cd ui
go build -o tf-gui-v10.exe *.go
```

### Build for Testing
```powershell
cd ui
go build -o tf-gui-test.exe *.go
```

### Build with Specific Files (if wildcards don't work)
```powershell
cd ui
go build -o tf-gui.exe main.go dashboard.go checklist.go position_sizing.go heat_check.go trade_entry.go scanner.go calendar.go theme.go widgets.go utils.go keybindings.go keybindings_v2.go
```

---

## Cross-Platform Building

### Build Windows Binary from Linux/macOS
```bash
cd ui
GOOS=windows GOARCH=amd64 go build -o tf-gui.exe *.go
```

### Build Linux Binary from Windows
```powershell
cd ui
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o tf-gui *.go
```

### Build macOS Binary from Windows
```powershell
cd ui
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o tf-gui-mac *.go
```

---

## Troubleshooting

### Error: "go: cannot find main module"
**Solution:** Make sure you're in the `ui/` directory
```powershell
cd ui
pwd  # Should show: .../trend-follower-dashboard/ui
```

### Error: "package X is not in GOROOT"
**Solution:** Install missing dependencies
```powershell
cd ui
go mod tidy
go mod download
```

### Error: "cannot find package"
**Solution:** Ensure all .go files are present
```powershell
ls *.go
# Should show:
# - main.go
# - dashboard.go
# - checklist.go
# - position_sizing.go
# - heat_check.go
# - trade_entry.go
# - scanner.go
# - calendar.go
# - theme.go
# - widgets.go
# - utils.go
# - keybindings.go
# - keybindings_v2.go
```

### Error: Build succeeds but binary doesn't run
**Solution:** Check for missing assets
```powershell
# Ensure database directory is accessible
cd ..
ls backend/internal/storage/
```

---

## Build Optimization

### Smaller Binary Size
```powershell
go build -ldflags="-s -w" -o tf-gui.exe *.go
```
- `-s` — Strip debug symbols
- `-w` — Strip DWARF debugging info
- Can reduce size by 20-30%

### Faster Builds (Development)
```powershell
go build -o tf-gui.exe *.go
# No optimization flags = faster compile
```

### Release Build (Optimized)
```powershell
go build -ldflags="-s -w" -o tf-gui.exe *.go
# Use this for distribution
```

---

## After Building

### Test the Binary
```powershell
# Run it
.\tf-gui.exe

# Check version
# (No version flag yet - check window title or logs)

# Check log output
notepad tf-gui.log
```

### Clean Up Old Builds
```powershell
# Remove old versions
rm tf-gui-v1.exe
rm tf-gui-v2.exe
# ... etc

# Or remove all versioned builds
rm tf-gui-v*.exe

# Keep only latest
rm tf-gui.exe
.\build.bat v9
```

---

## Development Workflow

### Typical Development Cycle

1. **Make code changes**
   ```powershell
   code main.go  # or any other .go file
   ```

2. **Build**
   ```powershell
   .\build.bat test
   ```

3. **Test**
   ```powershell
   .\tf-gui-test.exe
   ```

4. **Check logs**
   ```powershell
   notepad tf-gui.log
   ```

5. **Repeat** until satisfied

6. **Create release build**
   ```powershell
   .\build.bat v10
   ```

---

## File Sizes

Typical build sizes:
- **Debug build:** ~50-60 MB
- **Release build (-ldflags):** ~35-45 MB
- **Compressed (UPX):** ~15-20 MB (optional)

---

## Quick Reference

| Task | Command |
|------|---------|
| Default build | `.\build.bat` or `./build.sh` |
| Versioned build | `.\build.bat v10` |
| Manual build | `go build -o tf-gui.exe *.go` |
| Optimized build | `go build -ldflags="-s -w" -o tf-gui.exe *.go` |
| Cross-compile | `GOOS=windows GOARCH=amd64 go build` |
| Clean | `rm tf-gui*.exe` |
| Test | `.\tf-gui.exe` |
| Check logs | `notepad tf-gui.log` |

---

## Version History

| Version | Date | Key Features |
|---------|------|--------------|
| v1-v6 | 2025-10-30 | Button text fixes |
| v7 | 2025-10-30 | All buttons fixed |
| v8 | 2025-10-30 | 4 major features |
| v8-fixed | 2025-10-30 | Dialog sizing fixes |
| v9 | 2025-10-30 | VIM mode + toggle |
| v10+ | TBD | Future features |

---

## Need Help?

- **Build errors:** Check `go.mod` and `go.sum` are present
- **Runtime errors:** Check `tf-gui.log`
- **Missing features:** Check version number in code
- **Database errors:** Ensure `trading.db` is accessible

**Location:** All builds should be in the `ui/` directory alongside source files.
