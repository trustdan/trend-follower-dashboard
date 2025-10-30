# Phase 0 - Step 2: Fyne Proof-of-Concept

**Phase:** 0 - Foundation & Proof-of-Concept
**Step:** 2 of 4
**Duration:** 2-3 days
**Dependencies:** Step 1 (Dev Environment Setup)

---

## Objectives

Build a minimal desktop GUI application using Fyne to validate the "pure Go desktop app" approach:

1. Prove Go GUI can work for this project
2. Test direct in-process calls to backend (no HTTP, no CLI spawning)
3. Validate cross-compilation to Windows .exe
4. Test packaging and deployment
5. Assess Fyne's suitability for the full application
6. Provide a fallback option if Svelte approach fails

---

## Why Fyne?

**Fyne** is a pure Go GUI toolkit with:
- Cross-platform support (Windows, Linux, macOS)
- Material Design look and feel
- No external dependencies (no DLLs, no system libraries)
- Simple API for basic UIs
- Active community and good documentation
- Built-in packaging tools for .exe, .app, .deb

**Website:** https://fyne.io/

---

## Prerequisites

- Step 1 completed (Go 1.24+, backend tests passing)
- Backend compiles successfully
- Cross-compilation to Windows verified

---

## Step-by-Step Instructions

### 1. Install Fyne

```bash
# Navigate to project root
cd /home/kali/fresh-start-trading-platform

# Create a POC directory
mkdir -p poc/fyne-poc
cd poc/fyne-poc

# Initialize a Go module for the POC
go mod init github.com/fresh-start-trading-platform/fyne-poc

# Install Fyne v2
go get fyne.io/fyne/v2

# Verify installation
go list -m fyne.io/fyne/v2
# Expected: fyne.io/fyne/v2 v2.4.5 (or similar)
```

**System Dependencies (Linux):**

Fyne requires some system libraries on Linux:

```bash
# For Debian/Ubuntu/Kali
sudo apt-get install -y \
    gcc \
    libgl1-mesa-dev \
    xorg-dev

# Verify gcc is installed
gcc --version
```

---

### 2. Create Minimal Fyne Application

Create `poc/fyne-poc/main.go`:

```go
package main

import (
    "fmt"
    "log"

    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func main() {
    // Create a new Fyne application
    myApp := app.New()
    myWindow := myApp.NewWindow("TF-Engine Settings")
    myWindow.Resize(fyne.NewSize(600, 400))

    // Create UI components
    title := widget.NewLabel("TF-Engine Settings")
    title.TextStyle.Bold = true

    equityLabel := widget.NewLabel("Equity: $0")
    riskLabel := widget.NewLabel("Risk %: 0.00%")
    portfolioCapLabel := widget.NewLabel("Portfolio Cap: 0.00%")
    bucketCapLabel := widget.NewLabel("Bucket Cap: 0.00%")

    refreshBtn := widget.NewButton("Refresh", func() {
        equityLabel.SetText("Equity: $100,000")
        riskLabel.SetText("Risk %: 0.75%")
        portfolioCapLabel.SetText("Portfolio Cap: 4.00%")
        bucketCapLabel.SetText("Bucket Cap: 1.50%")
        log.Println("Data refreshed")
    })

    updateBtn := widget.NewButton("Update", func() {
        log.Println("Update clicked (not implemented in POC)")
    })

    statusLabel := widget.NewLabel("Status: Ready")

    // Layout
    content := container.NewVBox(
        title,
        widget.NewSeparator(),
        equityLabel,
        riskLabel,
        portfolioCapLabel,
        bucketCapLabel,
        widget.NewSeparator(),
        container.NewHBox(refreshBtn, updateBtn),
        widget.NewSeparator(),
        statusLabel,
    )

    myWindow.SetContent(content)
    myWindow.ShowAndRun()
}
```

---

### 3. Test the Fyne Application (Linux)

```bash
cd /home/kali/fresh-start-trading-platform/poc/fyne-poc

# Run the application
go run main.go

# Expected: A GUI window should appear with:
# - Title: "TF-Engine Settings"
# - Labels for equity, risk %, portfolio cap, bucket cap
# - Refresh and Update buttons
# - Status label

# Test the "Refresh" button
# Expected: Labels update to show hardcoded values
```

**If the window doesn't appear:**
- Check terminal for error messages
- Verify system dependencies are installed
- Try setting `DISPLAY` environment variable if using WSL2:
  ```bash
  # For WSL2 with VcXsrv or X410
  export DISPLAY=:0
  ```

---

### 4. Integrate with Backend Domain Logic

Now modify the POC to call real backend functions.

**Update `go.mod` to reference the backend:**

```bash
cd /home/kali/fresh-start-trading-platform/poc/fyne-poc

# Add replace directive to use local backend
# Edit go.mod and add:
```

Edit `poc/fyne-poc/go.mod`:
```go
module github.com/fresh-start-trading-platform/fyne-poc

go 1.24

require fyne.io/fyne/v2 v2.4.5

// Add this line to reference local backend
replace github.com/fresh-start-trading-platform/backend => ../../backend
```

**Update `main.go` to use backend:**

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"

    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    // Import backend packages
    "github.com/fresh-start-trading-platform/backend/internal/storage"
)

func main() {
    // Initialize database
    dbPath := filepath.Join("../../", "trading.db")
    db, err := storage.NewDB(dbPath)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    // Create Fyne application
    myApp := app.New()
    myWindow := myApp.NewWindow("TF-Engine Settings")
    myWindow.Resize(fyne.NewSize(600, 400))

    // Create UI components
    title := widget.NewLabel("TF-Engine Settings")
    title.TextStyle.Bold = true

    equityLabel := widget.NewLabel("Equity: $0")
    riskLabel := widget.NewLabel("Risk %: 0.00%")
    portfolioCapLabel := widget.NewLabel("Portfolio Cap: 0.00%")
    bucketCapLabel := widget.NewLabel("Bucket Cap: 0.00%")

    refreshBtn := widget.NewButton("Refresh", func() {
        // Call backend to get settings
        settings, err := db.GetSettings()
        if err != nil {
            log.Printf("Error getting settings: %v", err)
            return
        }

        // Update UI with real data
        equityLabel.SetText(fmt.Sprintf("Equity: $%.2f", settings.Equity))
        riskLabel.SetText(fmt.Sprintf("Risk %%: %.2f%%", settings.RiskPct))
        portfolioCapLabel.SetText(fmt.Sprintf("Portfolio Cap: %.2f%%", settings.PortfolioCap))
        bucketCapLabel.SetText(fmt.Sprintf("Bucket Cap: %.2f%%", settings.BucketCap))

        log.Println("Data refreshed from database")
    })

    updateBtn := widget.NewButton("Update", func() {
        // For POC, just log; full implementation would save to DB
        log.Println("Update clicked (would save to database)")
    })

    statusLabel := widget.NewLabel("Status: Ready")

    // Layout
    content := container.NewVBox(
        title,
        widget.NewSeparator(),
        equityLabel,
        riskLabel,
        portfolioCapLabel,
        bucketCapLabel,
        widget.NewSeparator(),
        container.NewHBox(refreshBtn, updateBtn),
        widget.NewSeparator(),
        statusLabel,
    )

    myWindow.SetContent(content)
    myWindow.ShowAndRun()
}
```

**Run the integrated version:**

```bash
# Ensure database exists
cd /home/kali/fresh-start-trading-platform
ls -la trading.db

# If database doesn't exist, initialize it
cd backend
./tf-engine init
cd ..

# Run the POC
cd poc/fyne-poc
go run main.go

# Click "Refresh" button
# Expected: Labels should show real values from trading.db
```

---

### 5. Cross-Compile Fyne to Windows

Fyne supports cross-compilation with some additional steps.

**Install fyne-cross (cross-compilation tool):**

```bash
# Install fyne-cross
go install github.com/fyne-io/fyne-cross@latest

# Verify installation
fyne-cross version
```

**Cross-compile to Windows:**

```bash
cd /home/kali/fresh-start-trading-platform/poc/fyne-poc

# Build Windows executable using fyne-cross
# This uses Docker to cross-compile, ensuring all dependencies are correct
fyne-cross windows -arch=amd64 -app-id=com.tf-engine.settings

# Alternative: Manual cross-compilation (may not work without proper setup)
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o fyne-poc.exe .

# Output should be in:
# fyne-cross/bin/windows-amd64/fyne-poc.exe
```

**Copy to Windows for testing:**

```bash
# Create export directory
mkdir -p /home/kali/fresh-start-trading-platform/poc/fyne-poc/dist

# Copy Windows binary
cp fyne-cross/bin/windows-amd64/fyne-poc.exe dist/

# Copy database for testing
cp /home/kali/fresh-start-trading-platform/trading.db dist/

# Create export zip
cd dist
zip fyne-poc-windows.zip fyne-poc.exe trading.db
cd ..

# Note the path for Windows access
echo "Windows path: \\\\wsl\$\\Kali\\home\\kali\\fresh-start-trading-platform\\poc\\fyne-poc\\dist"
```

**Test on Windows:**
1. From Windows Explorer, navigate to `\\wsl$\Kali\home\kali\fresh-start-trading-platform\poc\fyne-poc\dist`
2. Copy `fyne-poc.exe` and `trading.db` to a Windows directory
3. Double-click `fyne-poc.exe`
4. Window should appear with TF-Engine Settings
5. Click "Refresh" - should show settings from database

---

### 6. Test Fyne Packaging

Fyne includes built-in packaging tools:

```bash
cd /home/kali/fresh-start-trading-platform/poc/fyne-poc

# Package as application
fyne package -os linux -icon ../../docs/icon.png

# This creates a .tar.xz for Linux distribution

# For Windows (if fyne-cross worked):
fyne package -os windows -icon ../../docs/icon.png

# For macOS:
fyne package -os darwin -icon ../../docs/icon.png
```

---

## Pros and Cons Analysis

Document the findings for decision-making in Step 4.

### Pros (Why Choose Fyne)

**Technical:**
- ✅ Pure Go (same language as backend, no context switching)
- ✅ Single binary output (no external dependencies)
- ✅ Direct function calls (no HTTP overhead, no IPC complexity)
- ✅ Fast performance (compiled, native code)
- ✅ Good cross-platform support

**Development:**
- ✅ Simple API for basic UIs
- ✅ Built-in widgets (buttons, labels, forms, etc.)
- ✅ Material Design look is modern and clean
- ✅ Built-in packaging tools

**Deployment:**
- ✅ Easy distribution (single .exe file)
- ✅ No installation required (portable)
- ✅ Small binary size (~10-20 MB)

### Cons (Why Fyne May Not Be Ideal)

**UI Limitations:**
- ❌ **Limited styling flexibility** - Hard to achieve custom gradients and sophisticated designs
- ❌ **Fixed Material Design** - Can't easily create the sleek, gradient-heavy UI specified in overview-plan
- ❌ **Day/Night mode** - Possible but not as elegant as CSS-based themes
- ❌ Complex animations difficult - No native support for smooth gradient transitions
- ❌ Layout control less precise than CSS Grid/Flexbox

**Development Speed:**
- ❌ Slower iteration - Need to recompile to see changes (vs hot reload in Svelte)
- ❌ No browser dev tools - Harder to debug UI issues
- ❌ Limited component ecosystem - Must build everything custom

**Visual Appeal:**
- ❌ **Cannot easily achieve the "sleek, modern, gradient-heavy" design** specified in overview-plan
- ❌ Banner gradients (RED → crimson, etc.) difficult to implement smoothly
- ❌ Sophisticated micro-interactions harder to create
- ❌ Less "polished" look compared to modern web UIs

**Specific to This Project:**
- ❌ The large gradient banner (20% screen height, smooth transitions) is a core requirement
- ❌ Day/night mode with smooth theme toggle is harder to implement elegantly
- ❌ TradingView integration (opening browser) works but feels disconnected

### Verdict for This Project

**Fyne is a solid fallback** but **not ideal for the TF-Engine vision** because:

1. The overview-plan specifies a **"sleek, modern, gradient-heavy"** UI
2. The banner component is central to the anti-impulsivity design
3. Day/night mode with smooth transitions is a core feature
4. Visual appeal is important for daily use and discipline enforcement

**Recommendation:** Proceed with Svelte POC (Step 3) to validate if it provides superior visual capabilities while maintaining reasonable complexity.

---

## Documentation Requirements

Create `poc/fyne-poc/README.md`:

```markdown
# Fyne Proof-of-Concept

**Status:** ✅ Completed
**Date:** 2025-10-29

## What This Proves

- Fyne can build a Go desktop GUI
- Backend functions can be called directly (in-process)
- Cross-compilation to Windows works
- Single binary distribution is possible

## What We Learned

### Pros
- Simple API, easy to get started
- Pure Go, no external dependencies
- Fast compilation and execution
- Good for basic UIs

### Cons
- Limited styling flexibility
- Hard to achieve sophisticated gradients and animations
- Day/night mode less elegant than CSS
- Material Design is somewhat rigid

## Decision Impact

Fyne is a **solid fallback** if Svelte proves too complex.

For this project's vision (sleek, gradient-heavy, modern UI with smooth transitions), **Svelte is better aligned** with requirements.

## How to Run

### Linux
```bash
cd poc/fyne-poc
go run main.go
```

### Windows
1. Cross-compile: `fyne-cross windows -arch=amd64`
2. Copy `fyne-cross/bin/windows-amd64/fyne-poc.exe` to Windows
3. Copy `trading.db` to same directory
4. Run `fyne-poc.exe`

## Files Created

- `main.go` - Minimal Fyne app with backend integration
- `go.mod` - Dependencies
- `go.sum` - Dependency checksums
- `README.md` - This file
```

---

## Verification Checklist

Before proceeding to Step 3, verify:

- [ ] Fyne installed successfully
- [ ] Minimal Fyne app runs on Linux and displays window
- [ ] Backend integration works (can read from trading.db)
- [ ] "Refresh" button updates UI with real database values
- [ ] Cross-compilation to Windows .exe succeeds
- [ ] Windows .exe tested (on Windows machine if available)
- [ ] Pros/cons analysis documented
- [ ] `poc/fyne-poc/README.md` created

---

## Expected Outputs

After completing this step, you should have:

1. **Working Fyne POC:**
   - `poc/fyne-poc/main.go` with backend integration
   - Application runs on Linux
   - Application displays settings from database

2. **Windows Executable:**
   - `poc/fyne-poc/dist/fyne-poc.exe`
   - Verified to run on Windows (if tested)

3. **Documentation:**
   - `poc/fyne-poc/README.md` with findings
   - Pros/cons analysis documented

4. **Decision Input:**
   - Clear understanding of Fyne's capabilities and limitations
   - Evidence for Step 4 decision (Fyne vs Svelte)

---

## Troubleshooting

### Fyne Installation Issues

**Problem:** `go get fyne.io/fyne/v2` fails
**Solution:**
```bash
# Ensure Go is up to date
go version

# Clear module cache
go clean -modcache

# Try again
go get -u fyne.io/fyne/v2
```

**Problem:** System dependencies missing (Linux)
**Solution:**
```bash
# Install all required libraries
sudo apt-get update
sudo apt-get install -y gcc libgl1-mesa-dev xorg-dev

# Verify gcc
gcc --version
```

### Window Not Appearing

**Problem:** `go run main.go` completes but no window appears
**Solution:**
```bash
# Check for X server (if using WSL2)
echo $DISPLAY

# Install VcXsrv or X410 on Windows
# Start X server
# Set DISPLAY
export DISPLAY=:0

# Try again
go run main.go
```

### Cross-Compilation Issues

**Problem:** fyne-cross fails with Docker errors
**Solution:**
```bash
# Install Docker if not present
sudo apt-get install docker.io

# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in, then try again
```

**Problem:** Manual cross-compilation fails
**Solution:**
```bash
# Use fyne-cross instead (recommended)
fyne-cross windows -arch=amd64

# Or install MinGW toolchain
sudo apt-get install gcc-mingw-w64
```

### Backend Integration Issues

**Problem:** Cannot find database
**Solution:**
```bash
# Ensure trading.db exists
ls -la ../../trading.db

# Initialize if missing
cd ../../backend
./tf-engine init
cd -
```

**Problem:** Import cycle or module resolution errors
**Solution:**
```bash
# Verify go.mod has replace directive
cat go.mod | grep replace

# Run go mod tidy
go mod tidy
```

---

## Time Estimate

- **Fyne Installation:** 15-30 minutes
- **Minimal App Creation:** 30-45 minutes
- **Backend Integration:** 45-60 minutes
- **Testing:** 30 minutes
- **Cross-Compilation:** 30-45 minutes
- **Windows Testing:** 30 minutes (if available)
- **Documentation:** 30 minutes

**Total:** ~4-6 hours (1 day with breaks and troubleshooting)

---

## References

- [Fyne Documentation](https://developer.fyne.io/)
- [Fyne Getting Started](https://developer.fyne.io/started/)
- [Fyne Cross-Compilation](https://developer.fyne.io/started/cross-compiling)
- [fyne-cross GitHub](https://github.com/fyne-io/fyne-cross)
- [overview-plan.md - Proof-of-Concept Approach](../plans/overview-plan.md#proof-of-concept-approach)
- [overview-plan.md - Technology Recommendation](../plans/overview-plan.md#gui-implementation-guidance)

---

## Next Step

Proceed to: **[Phase 0 - Step 3: Svelte Proof-of-Concept](phase0-step3-svelte-poc.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
**Last Updated:** 2025-10-29
