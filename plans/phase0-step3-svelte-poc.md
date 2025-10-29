# Phase 0 - Step 3: Svelte Proof-of-Concept

**Phase:** 0 - Foundation & Proof-of-Concept
**Step:** 3 of 4
**Duration:** 3-4 days
**Dependencies:** Step 1 (Dev Environment Setup)

---

## Objectives

Build a minimal web-based GUI using SvelteKit + Go backend to validate the "browser-based UI" approach:

1. Prove SvelteKit static adapter can generate deployable files
2. Test Go serving embedded Svelte files via `embed` package
3. Validate HTTP API communication (Go backend ↔ Svelte frontend)
4. Implement day/night theme toggle with smooth transitions
5. Test cross-compilation to single Windows .exe
6. Assess visual capabilities (gradients, animations, modern design)
7. Demonstrate **why Svelte is superior for this project's vision**

---

## Why Svelte for This Project?

From `overview-plan.md`:

**Visual Requirements:**
- "Sleek, modern, gradient-heavy" design
- Large gradient banner (20% screen height) with smooth transitions
- Day/night mode with CSS variables and theme toggle
- Smooth animations (0.3s ease-in-out) on all state changes
- Custom components with sophisticated styling

**Why Svelte Excels:**
- ✅ CSS gradients and animations are trivial
- ✅ Theme toggle with CSS variables is built-in
- ✅ Hot reload for instant visual feedback
- ✅ Rich ecosystem (TailwindCSS, icon libraries, component libraries)
- ✅ Huge developer community (web skills are common)
- ✅ Professional, polished appearance out of the box
- ✅ Browser dev tools for debugging

**The Banner:** The centerpiece of the anti-impulsivity design requires smooth gradient transitions between RED/YELLOW/GREEN states. This is effortless in CSS but difficult in Fyne.

---

## Prerequisites

- Step 1 completed (Go 1.24+, Node.js 20+, backend tests passing)
- Node.js and npm installed and verified
- Backend compiles successfully

---

## Step-by-Step Instructions

### 1. Initialize SvelteKit Project

```bash
# Navigate to project root
cd /home/kali/fresh-start-trading-platform

# Create POC directory
mkdir -p poc/svelte-poc
cd poc/svelte-poc

# Initialize SvelteKit project
npm create svelte@latest .

# When prompted, select:
# - "Skeleton project" (minimal template)
# - TypeScript: Yes
# - ESLint: Yes
# - Prettier: Yes
# - Playwright: No (not needed for POC)
# - Vitest: No (not needed for POC)

# Install dependencies
npm install

# Verify development server works
npm run dev

# Expected: Server starts on http://localhost:5173/
# Open in browser to verify "Welcome to SvelteKit" page appears
```

---

### 2. Install and Configure Static Adapter

The static adapter compiles SvelteKit to plain HTML/CSS/JS files that can be served by any static file server (including Go's `embed`).

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc

# Install static adapter
npm install -D @sveltejs/adapter-static

# Verify installation
npm list @sveltejs/adapter-static
```

**Configure the adapter:**

Edit `svelte.config.js`:

```javascript
import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    // Consult https://kit.svelte.dev/docs/integrations#preprocessors
    // for more information about preprocessors
    preprocess: vitePreprocess(),

    kit: {
        // Use static adapter for generating static files
        adapter: adapter({
            pages: 'build',          // Output directory
            assets: 'build',         // Assets directory
            fallback: 'index.html',  // SPA fallback for client-side routing
            precompress: false,      // Don't generate .br/.gz files (not needed)
            strict: true             // Fail build if any page can't be prerendered
        })
    }
};

export default config;
```

**Configure prerendering:**

Create/edit `src/routes/+layout.ts`:

```typescript
// Enable prerendering for all routes
export const prerender = true;
// Use SPA mode (client-side routing)
export const ssr = false;
```

---

### 3. Install TailwindCSS

TailwindCSS makes gradient styling and responsive design trivial.

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc

# Install TailwindCSS and dependencies
npm install -D tailwindcss postcss autoprefixer

# Initialize TailwindCSS config
npx tailwindcss init -p

# This creates:
# - tailwind.config.js
# - postcss.config.js
```

**Configure TailwindCSS:**

Edit `tailwind.config.js`:

```javascript
/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    theme: {
        extend: {
            // Custom colors for day/night mode
            colors: {
                // Banner gradients
                'banner-red-start': '#DC2626',
                'banner-red-end': '#991B1B',
                'banner-yellow-start': '#F59E0B',
                'banner-yellow-end': '#FBBF24',
                'banner-green-start': '#10B981',
                'banner-green-end': '#059669',
            },
        },
    },
    plugins: [],
    // Enable dark mode via class
    darkMode: 'class',
};
```

**Add TailwindCSS to your styles:**

Create `src/app.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

/* CSS Variables for theming */
:root {
    /* Day mode (light theme) */
    --bg-primary: #FFFFFF;
    --bg-secondary: #F9FAFB;
    --bg-tertiary: #F3F4F6;
    --text-primary: #111827;
    --text-secondary: #6B7280;
    --text-tertiary: #9CA3AF;
    --border-color: #E5E7EB;
    --border-focus: #3B82F6;
}

:root.dark {
    /* Night mode (dark theme) */
    --bg-primary: #0F172A;
    --bg-secondary: #1E293B;
    --bg-tertiary: #334155;
    --text-primary: #F1F5F9;
    --text-secondary: #CBD5E1;
    --text-tertiary: #94A3B8;
    --border-color: #334155;
    --border-focus: #60A5FA;
}

/* Smooth transitions for theme changes */
* {
    transition: background-color 0.3s ease-in-out, color 0.3s ease-in-out, border-color 0.3s ease-in-out;
}

/* Banner gradient classes */
.banner-red {
    background: linear-gradient(135deg, #DC2626 0%, #991B1B 100%);
}

.banner-yellow {
    background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
}

.banner-green {
    background: linear-gradient(135deg, #10B981 0%, #059669 100%);
}

/* Banner animation */
.banner {
    transition: background 0.3s ease-in-out;
    animation: pulse 0.5s ease-in-out;
}

@keyframes pulse {
    0% {
        transform: scale(1);
    }
    50% {
        transform: scale(1.02);
    }
    100% {
        transform: scale(1);
    }
}
```

**Import CSS in layout:**

Edit `src/routes/+layout.svelte`:

```svelte
<script lang="ts">
    import '../app.css';
</script>

<slot />
```

---

### 4. Create Settings Page with Theme Toggle

Create `src/routes/+page.svelte`:

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import { writable } from 'svelte/store';

    // Theme store
    const isDark = writable(false);

    // Settings data
    let settings = {
        equity: 0,
        riskPct: 0,
        portfolioCap: 0,
        bucketCap: 0,
    };

    let loading = false;
    let error = '';

    // Load theme preference from localStorage
    onMount(() => {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === 'dark') {
            isDark.set(true);
            document.documentElement.classList.add('dark');
        }
    });

    // Toggle theme
    function toggleTheme() {
        isDark.update(v => {
            const newValue = !v;
            if (newValue) {
                document.documentElement.classList.add('dark');
                localStorage.setItem('theme', 'dark');
            } else {
                document.documentElement.classList.remove('dark');
                localStorage.setItem('theme', 'light');
            }
            return newValue;
        });
    }

    // Fetch settings from API
    async function fetchSettings() {
        loading = true;
        error = '';

        try {
            const response = await fetch('/api/settings');
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            const data = await response.json();
            settings = data;
        } catch (err) {
            error = err instanceof Error ? err.message : 'Failed to fetch settings';
            console.error('Error fetching settings:', err);
        } finally {
            loading = false;
        }
    }

    // Initialize on mount
    onMount(() => {
        fetchSettings();
    });
</script>

<div class="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)] p-8">
    <!-- Header with Theme Toggle -->
    <header class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold">TF-Engine Settings</h1>

        <button
            on:click={toggleTheme}
            class="p-2 rounded-lg bg-[var(--bg-secondary)] hover:bg-[var(--bg-tertiary)] border border-[var(--border-color)]"
            aria-label="Toggle theme"
        >
            {#if $isDark}
                <!-- Sun icon (day mode) -->
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
            {:else}
                <!-- Moon icon (night mode) -->
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                </svg>
            {/if}
        </button>
    </header>

    <!-- Main Content Card -->
    <div class="max-w-2xl mx-auto bg-[var(--bg-secondary)] rounded-xl shadow-lg p-6 border border-[var(--border-color)]">
        <h2 class="text-2xl font-semibold mb-6">Account Settings</h2>

        {#if loading}
            <div class="text-center py-8">
                <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--border-focus)]"></div>
                <p class="mt-2 text-[var(--text-secondary)]">Loading...</p>
            </div>
        {:else if error}
            <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-4">
                <p class="text-red-800 dark:text-red-200">{error}</p>
            </div>
        {/if}

        <!-- Settings Display -->
        <div class="space-y-4">
            <div class="flex justify-between items-center py-3 border-b border-[var(--border-color)]">
                <span class="text-[var(--text-secondary)]">Equity:</span>
                <span class="font-semibold text-lg">${settings.equity.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
            </div>

            <div class="flex justify-between items-center py-3 border-b border-[var(--border-color)]">
                <span class="text-[var(--text-secondary)]">Risk % per Unit:</span>
                <span class="font-semibold text-lg">{settings.riskPct.toFixed(2)}%</span>
            </div>

            <div class="flex justify-between items-center py-3 border-b border-[var(--border-color)]">
                <span class="text-[var(--text-secondary)]">Portfolio Heat Cap:</span>
                <span class="font-semibold text-lg">{settings.portfolioCap.toFixed(2)}%</span>
            </div>

            <div class="flex justify-between items-center py-3">
                <span class="text-[var(--text-secondary)]">Sector Bucket Cap:</span>
                <span class="font-semibold text-lg">{settings.bucketCap.toFixed(2)}%</span>
            </div>
        </div>

        <!-- Action Buttons -->
        <div class="mt-6 flex gap-4">
            <button
                on:click={fetchSettings}
                disabled={loading}
                class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
                Refresh
            </button>

            <button
                disabled
                class="flex-1 px-4 py-2 bg-gray-300 dark:bg-gray-700 text-gray-500 dark:text-gray-400 font-medium rounded-lg cursor-not-allowed"
            >
                Update (Not Implemented)
            </button>
        </div>

        <!-- Status -->
        <p class="mt-4 text-sm text-[var(--text-tertiary)] text-center">
            Status: Ready • API: {loading ? 'Loading...' : error ? 'Error' : 'Connected'}
        </p>
    </div>

    <!-- Demo Banner (shows gradient transitions) -->
    <div class="max-w-2xl mx-auto mt-8">
        <h3 class="text-xl font-semibold mb-4">Banner Demo (Anti-Impulsivity Core)</h3>
        <div class="banner banner-green rounded-xl p-8 text-white text-center shadow-lg">
            <p class="text-3xl font-bold mb-2">✓ OK TO TRADE ✓</p>
            <p class="text-lg opacity-90">All gates pass • Quality score met</p>
        </div>
    </div>
</div>
```

**Test the Svelte app:**

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc

# Run development server with hot reload
npm run dev

# Open browser to: http://localhost:5173/

# Expected:
# - Page displays "TF-Engine Settings"
# - Theme toggle button in top-right (sun/moon icon)
# - Settings card with placeholder data
# - "Refresh" and "Update" buttons
# - Green banner demo at bottom

# Test theme toggle:
# - Click sun/moon icon
# - Page should smoothly transition to dark mode
# - Theme preference should persist on reload
```

---

### 5. Build Static Files

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc

# Build for production
npm run build

# Expected output: Static files in build/ directory
ls -la build/

# Should contain:
# - index.html
# - _app/ (directory with CSS and JS chunks)
# - favicon.png
# - etc.
```

---

### 6. Create Go HTTP Server to Serve Svelte

Now create a Go backend that serves the Svelte static files and provides an API endpoint.

**Create Go module:**

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc

# Create Go server directory
mkdir go-server
cd go-server

# Initialize Go module
go mod init github.com/fresh-start-trading-platform/svelte-poc-server
```

**Create directory for embedded files:**

```bash
mkdir -p webui/dist
```

**Copy Svelte build output:**

```bash
cp -r ../build/* webui/dist/
```

**Create embedding package:**

Create `go-server/webui/embed.go`:

```go
// Package webui exposes an fs.FS with the compiled Svelte app.
package webui

import (
    "embed"
    "io/fs"
)

//go:embed all:dist
var dist embed.FS

// Sub returns an fs.FS rooted at the "dist" folder.
func Sub() (fs.FS, error) {
    return fs.Sub(dist, "dist")
}
```

**Create main server:**

Create `go-server/main.go`:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/fresh-start-trading-platform/svelte-poc-server/webui"
    "github.com/fresh-start-trading-platform/backend/internal/storage"
)

func main() {
    // Initialize database
    dbPath := filepath.Join("../../..", "trading.db")
    db, err := storage.NewDB(dbPath)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    // Get embedded Svelte files
    sfs, err := webui.Sub()
    if err != nil {
        log.Fatalf("Failed to get embedded files: %v", err)
    }

    // Create HTTP server
    mux := http.NewServeMux()

    // API endpoint: GET /api/settings
    mux.HandleFunc("/api/settings", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        settings, err := db.GetSettings()
        if err != nil {
            log.Printf("Error getting settings: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(settings)
    })

    // Serve static files (SPA fallback)
    mux.Handle("/", http.FileServer(http.FS(sfs)))

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server on http://localhost:%s", port)
    log.Printf("Press Ctrl+C to stop")

    if err := http.ListenAndServe(":"+port, mux); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

**Update go.mod to reference backend:**

Edit `go-server/go.mod`:

```go
module github.com/fresh-start-trading-platform/svelte-poc-server

go 1.24

// Reference local backend
replace github.com/fresh-start-trading-platform/backend => ../../../backend
```

**Download dependencies:**

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc/go-server

go mod tidy
```

---

### 7. Test Go Server with Embedded Svelte

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc/go-server

# Run the Go server
go run main.go

# Expected output:
# Starting server on http://localhost:8080
# Press Ctrl+C to stop

# Open browser to: http://localhost:8080/

# Expected:
# - Svelte app loads
# - Click "Refresh" button
# - Settings should load from database via API
# - Theme toggle works
# - No console errors
```

---

### 8. Build Windows Executable with Embedded UI

This is the key test: can we create a single .exe that contains everything?

**Create sync script:**

Create `poc/svelte-poc/sync-ui-to-go.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

echo "[1/3] Building Svelte UI..."
npm ci
npm run build

echo "[2/3] Staging UI for Go embed..."
rm -rf go-server/webui/dist
mkdir -p go-server/webui/dist
cp -R build/* go-server/webui/dist/

echo "[3/3] Done. Staged at go-server/webui/dist/"
```

Make executable:

```bash
chmod +x poc/svelte-poc/sync-ui-to-go.sh
```

**Build Windows executable:**

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc

# Sync UI to Go
./sync-ui-to-go.sh

# Build Windows binary (pure Go, no cgo)
cd go-server
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" -o svelte-poc.exe .

# Verify binary
ls -lh svelte-poc.exe
file svelte-poc.exe

# Expected:
# svelte-poc.exe: PE32+ executable (console) x86-64 (stripped to external PDB), for MS Windows
```

**Create export package:**

```bash
cd /home/kali/fresh-start-trading-platform/poc/svelte-poc

mkdir -p dist
cp go-server/svelte-poc.exe dist/
cp ../../trading.db dist/

# Create README for Windows testing
cat > dist/README.txt << 'EOF'
TF-Engine Svelte POC - Windows Testing
=======================================

1. Ensure trading.db is in the same directory as svelte-poc.exe
2. Double-click svelte-poc.exe to start the server
3. Server will start on http://localhost:8080
4. Open your browser to http://localhost:8080
5. Test the theme toggle (sun/moon icon)
6. Test the Refresh button (loads data from database)

To stop: Close the console window or press Ctrl+C

If port 8080 is in use:
  set PORT=8090
  svelte-poc.exe
EOF

# Create zip for Windows
cd dist
zip svelte-poc-windows.zip svelte-poc.exe trading.db README.txt
cd ..

echo "Windows package ready: poc/svelte-poc/dist/svelte-poc-windows.zip"
echo "Windows path: \\\\wsl\$\\Kali\\home\\kali\\fresh-start-trading-platform\\poc\\svelte-poc\\dist"
```

---

### 9. Test on Windows

**Steps:**
1. From Windows Explorer, navigate to: `\\wsl$\Kali\home\kali\fresh-start-trading-platform\poc\svelte-poc\dist`
2. Copy `svelte-poc-windows.zip` to a Windows directory
3. Extract the zip
4. Double-click `svelte-poc.exe`
5. Console window should appear with: "Starting server on http://localhost:8080"
6. Open browser to `http://localhost:8080`
7. Verify:
   - Svelte app loads correctly
   - Theme toggle works (sun/moon icon switches, UI transitions smoothly)
   - Click "Refresh" - settings should load from `trading.db`
   - Green banner demo displays with gradient
   - No errors in browser console
8. Test theme persistence:
   - Toggle to dark mode
   - Refresh the page
   - Should remain in dark mode
9. Close browser, stop server (Ctrl+C or close console)

---

## Pros and Cons Analysis

Document findings for Step 4 decision.

### Pros (Why Choose Svelte)

**Visual Capabilities (Critical for This Project):**
- ✅ **Perfect for gradient-heavy design** - CSS gradients are effortless
- ✅ **Smooth animations** - Built-in transitions and CSS animations
- ✅ **Day/night mode** - Elegant with CSS variables and class toggle
- ✅ **Banner component** - Can achieve the exact RED/YELLOW/GREEN gradient vision
- ✅ **Modern, polished look** - Matches the "sleek" requirement
- ✅ **TailwindCSS integration** - Rapid UI development with utility classes

**Development Experience:**
- ✅ **Hot reload** - Instant feedback, no recompilation
- ✅ **Browser dev tools** - Inspect, debug, test easily
- ✅ **TypeScript support** - Type safety in frontend
- ✅ **Rich ecosystem** - Component libraries, icons, charts
- ✅ **Common skills** - Web development knowledge transfers

**Technical:**
- ✅ **Single binary** - Go can embed Svelte static files
- ✅ **Fast performance** - Svelte compiles to vanilla JS, no runtime
- ✅ **API communication** - Standard HTTP/JSON, well understood
- ✅ **Static adapter** - No server-side rendering complexity
- ✅ **Cross-platform** - Browser is the same everywhere

### Cons (Challenges)

**Complexity:**
- ⚠️ **Two-language stack** - Go backend + TypeScript frontend
- ⚠️ **Build pipeline** - Must build Svelte, then embed in Go
- ⚠️ **More moving parts** - npm, SvelteKit, Go, HTTP API

**Debugging:**
- ⚠️ **API layer** - HTTP communication adds complexity vs direct calls
- ⚠️ **State management** - Must sync state between frontend and backend
- ⚠️ **CORS issues** - (Avoided by serving from same origin)

**Distribution:**
- ⚠️ **Larger binary** - Embedded HTML/CSS/JS increases size (~5-10 MB more)
- ⚠️ **Browser required** - Must open browser (but auto-open is easy)

### Verdict for This Project

**Svelte is the CLEAR WINNER** for TF-Engine because:

1. **The banner is the core of the design** - Gradients and smooth transitions are essential
2. **Visual appeal matters for daily use** - Discipline is easier with a pleasant interface
3. **Day/night mode is a requirement** - CSS makes this elegant
4. **Development speed** - Hot reload means faster iteration
5. **Professional appearance** - Matches modern SaaS app expectations

**The added complexity is justified** by the superior visual capabilities and developer experience.

---

## Verification Checklist

Before proceeding to Step 4, verify:

- [ ] SvelteKit project initialized successfully
- [ ] Static adapter configured and builds to `build/` directory
- [ ] TailwindCSS installed and configured
- [ ] Settings page created with theme toggle
- [ ] Theme toggle works and persists to localStorage
- [ ] Day/night mode transitions smoothly (0.3s)
- [ ] Go server created and serves embedded Svelte files
- [ ] API endpoint `/api/settings` returns data from database
- [ ] "Refresh" button calls API and updates UI
- [ ] Windows .exe builds successfully (with embedded UI)
- [ ] Windows .exe tested (on Windows machine if available)
- [ ] Single binary contains both backend and frontend
- [ ] Pros/cons analysis documented

---

## Expected Outputs

After completing this step, you should have:

1. **Working Svelte POC:**
   - `poc/svelte-poc/` - SvelteKit project
   - `poc/svelte-poc/build/` - Static build output
   - `poc/svelte-poc/go-server/` - Go HTTP server
   - Application runs on Linux (dev server and Go server)

2. **Windows Executable:**
   - `poc/svelte-poc/dist/svelte-poc.exe` - Single binary with embedded UI
   - `poc/svelte-poc/dist/svelte-poc-windows.zip` - Distribution package

3. **Documentation:**
   - `poc/svelte-poc/README.md` with findings
   - Pros/cons analysis documented
   - Screenshots (optional but recommended)

4. **Decision Input:**
   - Clear evidence that Svelte can achieve the visual requirements
   - Demonstration of theme toggle and smooth transitions
   - Proof that single binary deployment works

---

## Troubleshooting

### SvelteKit Issues

**Problem:** `npm create svelte` fails
**Solution:**
```bash
# Clear npm cache
npm cache clean --force

# Update npm
npm install -g npm@latest

# Try again
npm create svelte@latest .
```

**Problem:** Build fails with "Cannot find module"
**Solution:**
```bash
# Reinstall dependencies
rm -rf node_modules
npm install
```

### Static Adapter Issues

**Problem:** Build fails with "Not found: /..."
**Solution:**
- Ensure `export const prerender = true` is in `src/routes/+layout.ts`
- Check `svelte.config.js` has `fallback: 'index.html'`

**Problem:** API calls don't work
**Solution:**
- Use relative URLs (`/api/settings`, not `http://localhost:8080/api/settings`)
- Ensure Go server is handling `/api/*` before static files

### Go Server Issues

**Problem:** Embedded files not found
**Solution:**
```bash
# Ensure dist directory has files
ls -la go-server/webui/dist/

# Rebuild if empty
cd poc/svelte-poc
./sync-ui-to-go.sh
```

**Problem:** Database connection fails
**Solution:**
```bash
# Check database path
ls -la ../../../trading.db

# Adjust path in main.go if needed
```

### Windows Testing Issues

**Problem:** .exe won't run on Windows
**Solution:**
- Ensure `CGO_ENABLED=0` was set during build
- Check if Windows Defender blocked the .exe
- Run from CMD/PowerShell to see error messages

**Problem:** API calls fail on Windows
**Solution:**
- Ensure `trading.db` is in same directory as .exe
- Check console output for error messages
- Verify database path logic in Go code

---

## Time Estimate

- **SvelteKit Setup:** 30-45 minutes
- **Static Adapter Configuration:** 30 minutes
- **TailwindCSS Setup:** 30 minutes
- **Settings Page Creation:** 2-3 hours
- **Theme Toggle Implementation:** 1-2 hours
- **Go Server Creation:** 1-2 hours
- **Embedding & Integration:** 1-2 hours
- **Windows Build & Testing:** 1-2 hours
- **Documentation:** 1 hour

**Total:** ~10-14 hours (2-3 days with breaks)

---

## References

- [SvelteKit Documentation](https://kit.svelte.dev/docs)
- [Static Adapter](https://kit.svelte.dev/docs/adapter-static)
- [TailwindCSS](https://tailwindcss.com/docs)
- [Go embed Package](https://pkg.go.dev/embed)
- [overview-plan.md - Visual Design Philosophy](../plans/overview-plan.md#visual-design-philosophy)
- [overview-plan.md - Why Svelte](../plans/overview-plan.md#proof-of-concept-approach)
- [RULES.md - Section 5](../1._RULES.md—Operating_Rules_for_This_Project-(Claude_Code).md)

---

## Next Step

Proceed to: **[Phase 0 - Step 4: Technology Decision & Build Pipeline](phase0-step4-decision-pipeline.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
**Last Updated:** 2025-10-29
