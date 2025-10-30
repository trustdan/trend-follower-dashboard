# Phase 1 - Step 6: Application Layout & Navigation

**Phase:** 1 - Dashboard & FINVIZ Scanner
**Step:** 6 of 9 (overall), 2 of 5 (Phase 1)
**Duration:** 2 days
**Dependencies:** Step 5 (Backend API Foundation)

---

## Objectives

Build the application shell that houses all screens:

1. Create main layout component (`App.svelte` / `+layout.svelte`)
2. Implement file-based routing (SvelteKit routes)
3. Create header component with app title, theme toggle, settings icon
4. Create navigation component (sidebar with links to all main screens)
5. Implement theme toggle functionality (day/night mode)
6. Ensure theme preference persists to localStorage
7. Create placeholder routes for all main screens
8. Test navigation between screens with smooth transitions
9. Add frontend logging utility

---

## Prerequisites

- Step 5 completed (Backend API functional)
- Production UI structure created (`ui/src/lib/...`)
- TailwindCSS configured with custom theme
- SvelteKit with static adapter configured

---

## Step-by-Step Instructions

### 1. Create Theme Store

Create `ui/src/lib/stores/theme.ts`:

```typescript
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Theme type
export type Theme = 'light' | 'dark';

// Create writable store with initial value from localStorage or system preference
function createThemeStore() {
    // Get initial theme
    const getInitialTheme = (): Theme => {
        if (!browser) return 'light';

        // Check localStorage first
        const stored = localStorage.getItem('theme') as Theme | null;
        if (stored) return stored;

        // Check system preference
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            return 'dark';
        }

        return 'light';
    };

    const { subscribe, set, update } = writable<Theme>(getInitialTheme());

    return {
        subscribe,
        toggle: () => {
            update(current => {
                const newTheme: Theme = current === 'light' ? 'dark' : 'light';

                if (browser) {
                    // Update localStorage
                    localStorage.setItem('theme', newTheme);

                    // Update document class
                    if (newTheme === 'dark') {
                        document.documentElement.classList.add('dark');
                    } else {
                        document.documentElement.classList.remove('dark');
                    }
                }

                return newTheme;
            });
        },
        set: (theme: Theme) => {
            set(theme);

            if (browser) {
                localStorage.setItem('theme', theme);

                if (theme === 'dark') {
                    document.documentElement.classList.add('dark');
                } else {
                    document.documentElement.classList.remove('dark');
                }
            }
        }
    };
}

export const theme = createThemeStore();
```

---

### 2. Create Logger Utility

Create `ui/src/lib/utils/logger.ts`:

```typescript
// Logger utility for frontend
export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

interface LogEntry {
    timestamp: string;
    level: LogLevel;
    message: string;
    data?: any;
}

class Logger {
    private logs: LogEntry[] = [];
    private maxLogs = 1000; // Keep last 1000 logs

    private log(level: LogLevel, message: string, data?: any) {
        const entry: LogEntry = {
            timestamp: new Date().toISOString(),
            level,
            message,
            data
        };

        this.logs.push(entry);

        // Keep only maxLogs entries
        if (this.logs.length > this.maxLogs) {
            this.logs.shift();
        }

        // Console output with color coding
        const styles: Record<LogLevel, string> = {
            debug: 'color: gray',
            info: 'color: blue',
            warn: 'color: orange',
            error: 'color: red'
        };

        console.log(
            `%c[${entry.timestamp}] [${level.toUpperCase()}] ${message}`,
            styles[level],
            data || ''
        );
    }

    debug(message: string, data?: any) {
        this.log('debug', message, data);
    }

    info(message: string, data?: any) {
        this.log('info', message, data);
    }

    warn(message: string, data?: any) {
        this.log('warn', message, data);
    }

    error(message: string, data?: any) {
        this.log('error', message, data);
    }

    // Get all logs
    getLogs(): LogEntry[] {
        return [...this.logs];
    }

    // Clear logs
    clear() {
        this.logs = [];
    }

    // Export logs as JSON
    export(): string {
        return JSON.stringify(this.logs, null, 2);
    }
}

// Export singleton instance
export const logger = new Logger();
```

---

### 3. Create Root Layout

Edit `ui/src/routes/+layout.svelte`:

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import '../app.css';
    import { theme } from '$lib/stores/theme';
    import { logger } from '$lib/utils/logger';

    // Initialize theme on mount
    onMount(() => {
        const currentTheme = $theme;
        logger.info('App initialized', { theme: currentTheme });

        if (currentTheme === 'dark') {
            document.documentElement.classList.add('dark');
        }
    });
</script>

<div class="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)]">
    <slot />
</div>
```

---

### 4. Create Header Component

Create `ui/src/lib/components/layout/Header.svelte`:

```svelte
<script lang="ts">
    import { theme } from '$lib/stores/theme';
    import { logger } from '$lib/utils/logger';

    function toggleTheme() {
        theme.toggle();
        logger.info('Theme toggled', { newTheme: $theme });
    }
</script>

<header class="bg-[var(--bg-secondary)] border-b border-[var(--border-color)] px-6 py-4">
    <div class="max-w-7xl mx-auto flex justify-between items-center">
        <!-- App Title -->
        <div class="flex items-center space-x-3">
            <h1 class="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400">
                TF-Engine
            </h1>
            <span class="text-sm text-[var(--text-tertiary)]">Trend Following System</span>
        </div>

        <!-- Right side: Theme toggle + Settings -->
        <div class="flex items-center space-x-4">
            <!-- Theme Toggle Button -->
            <button
                on:click={toggleTheme}
                class="p-2 rounded-lg bg-[var(--bg-primary)] hover:bg-[var(--bg-tertiary)] border border-[var(--border-color)] transition-colors"
                aria-label="Toggle theme"
                title={$theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
            >
                {#if $theme === 'dark'}
                    <!-- Sun icon (switch to light) -->
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                {:else}
                    <!-- Moon icon (switch to dark) -->
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                    </svg>
                {/if}
            </button>

            <!-- Settings Icon (placeholder) -->
            <button
                class="p-2 rounded-lg bg-[var(--bg-primary)] hover:bg-[var(--bg-tertiary)] border border-[var(--border-color)] transition-colors"
                aria-label="Settings"
                title="Settings"
            >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
            </button>
        </div>
    </div>
</header>
```

---

### 5. Create Navigation Component

Create `ui/src/lib/components/layout/Navigation.svelte`:

```svelte
<script lang="ts">
    import { page } from '$app/stores';
    import { logger } from '$lib/utils/logger';

    // Navigation items
    const navItems = [
        { path: '/', label: 'Dashboard', icon: 'home' },
        { path: '/scanner', label: 'Scanner', icon: 'search' },
        { path: '/checklist', label: 'Checklist', icon: 'clipboard' },
        { path: '/sizing', label: 'Position Sizing', icon: 'calculator' },
        { path: '/heat', label: 'Heat Check', icon: 'fire' },
        { path: '/entry', label: 'Trade Entry', icon: 'check' },
        { path: '/calendar', label: 'Calendar', icon: 'calendar' }
    ];

    // Check if route is active
    function isActive(path: string): boolean {
        return $page.url.pathname === path;
    }

    function handleNavigation(path: string, label: string) {
        logger.info('Navigation', { from: $page.url.pathname, to: path, label });
    }

    // Simple icon component (replace with Lucide icons later if desired)
    function getIconPath(icon: string): string {
        const icons: Record<string, string> = {
            home: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6',
            search: 'M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z',
            clipboard: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01',
            calculator: 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z',
            fire: 'M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z',
            check: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
            calendar: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z'
        };
        return icons[icon] || '';
    }
</script>

<nav class="w-64 bg-[var(--bg-secondary)] border-r border-[var(--border-color)] min-h-[calc(100vh-4rem)] p-4">
    <ul class="space-y-2">
        {#each navItems as item (item.path)}
            <li>
                <a
                    href={item.path}
                    on:click={() => handleNavigation(item.path, item.label)}
                    class="flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors {isActive(item.path)
                        ? 'bg-gradient-to-r from-blue-500 to-purple-500 text-white font-semibold'
                        : 'hover:bg-[var(--bg-tertiary)] text-[var(--text-primary)]'}"
                >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getIconPath(item.icon)} />
                    </svg>
                    <span>{item.label}</span>
                </a>
            </li>
        {/each}
    </ul>
</nav>
```

---

### 6. Create Main Layout with Header and Navigation

Edit `ui/src/routes/+layout.svelte`:

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import '../app.css';
    import { theme } from '$lib/stores/theme';
    import { logger } from '$lib/utils/logger';
    import Header from '$lib/components/layout/Header.svelte';
    import Navigation from '$lib/components/layout/Navigation.svelte';

    // Initialize theme on mount
    onMount(() => {
        const currentTheme = $theme;
        logger.info('App initialized', { theme: currentTheme });

        if (currentTheme === 'dark') {
            document.documentElement.classList.add('dark');
        }
    });
</script>

<div class="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)]">
    <Header />

    <div class="flex">
        <Navigation />

        <!-- Main content area -->
        <main class="flex-1 p-6 page-transition">
            <slot />
        </main>
    </div>
</div>
```

---

### 7. Create Placeholder Routes

Create placeholder pages for all main screens:

**Dashboard (`ui/src/routes/+page.svelte`):**

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import { logger } from '$lib/utils/logger';

    onMount(() => {
        logger.info('Dashboard page mounted');
    });
</script>

<div class="max-w-7xl mx-auto">
    <h2 class="text-3xl font-bold mb-6">Dashboard</h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- Portfolio Summary Card -->
        <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
            <h3 class="text-lg font-semibold mb-4">Portfolio Summary</h3>
            <p class="text-[var(--text-secondary)]">Coming in Step 7</p>
        </div>

        <!-- Open Positions Card -->
        <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
            <h3 class="text-lg font-semibold mb-4">Open Positions</h3>
            <p class="text-[var(--text-secondary)]">Coming in Step 7</p>
        </div>

        <!-- Candidates Card -->
        <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
            <h3 class="text-lg font-semibold mb-4">Today's Candidates</h3>
            <p class="text-[var(--text-secondary)]">Coming in Step 7</p>
        </div>
    </div>
</div>
```

**Scanner (`ui/src/routes/scanner/+page.svelte`):**

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import { logger } from '$lib/utils/logger';

    onMount(() => {
        logger.info('Scanner page mounted');
    });
</script>

<div class="max-w-7xl mx-auto">
    <h2 class="text-3xl font-bold mb-6">FINVIZ Scanner</h2>

    <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
        <p class="text-[var(--text-secondary)]">Scanner functionality coming in Step 8</p>
    </div>
</div>
```

**Checklist (`ui/src/routes/checklist/+page.svelte`):**

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import { logger } from '$lib/utils/logger';

    onMount(() => {
        logger.info('Checklist page mounted');
    });
</script>

<div class="max-w-7xl mx-auto">
    <h2 class="text-3xl font-bold mb-6">Checklist</h2>

    <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
        <p class="text-[var(--text-secondary)]">Checklist functionality coming in Phase 2</p>
    </div>
</div>
```

Create similar placeholders for:
- `ui/src/routes/sizing/+page.svelte` - Position Sizing
- `ui/src/routes/heat/+page.svelte` - Heat Check
- `ui/src/routes/entry/+page.svelte` - Trade Entry
- `ui/src/routes/calendar/+page.svelte` - Calendar

---

### 8. Test Navigation

```bash
cd /home/kali/fresh-start-trading-platform/ui

# Run development server
npm run dev

# Open browser to: http://localhost:5173/

# Test:
# 1. Header displays with "TF-Engine" title
# 2. Theme toggle button visible (sun/moon icon)
# 3. Settings icon visible
# 4. Left sidebar with 7 navigation items
# 5. Dashboard is active (highlighted)
# 6. Click theme toggle - smooth transition to dark mode
# 7. Refresh page - theme persists
# 8. Click Scanner - navigates to /scanner
# 9. Check browser console - navigation events logged
# 10. Try all navigation items
```

---

### 9. Build and Test with Go Server

```bash
# Build Svelte UI and sync to Go
cd /home/kali/fresh-start-trading-platform
./scripts/sync-ui-to-go.sh

# Run Go server
cd backend
go run cmd/tf-engine/main.go server --listen 127.0.0.1:8080

# Open browser to: http://localhost:8080/

# Test:
# 1. App loads with header and navigation
# 2. Theme toggle works
# 3. Navigation works
# 4. Theme persists on reload
# 5. All placeholder pages accessible
```

---

## Verification Checklist

Before proceeding to Step 7, verify:

- [ ] Theme store created with localStorage persistence
- [ ] Logger utility created with console output
- [ ] Header component created with theme toggle
- [ ] Theme toggle button works (sun/moon icons)
- [ ] Navigation component created with 7 items
- [ ] Active route highlighted with gradient
- [ ] Root layout combines Header + Navigation + content
- [ ] Placeholder routes created for all 7 screens
- [ ] Navigation between screens works
- [ ] Theme persists on page reload
- [ ] Smooth transitions (0.3s) on theme change
- [ ] Console logs navigation events
- [ ] Build and sync to Go works
- [ ] Go server serves UI correctly
- [ ] No console errors

---

## Expected Outputs

After completing this step, you should have:

1. **Stores:**
   - `src/lib/stores/theme.ts` - Theme management

2. **Utilities:**
   - `src/lib/utils/logger.ts` - Frontend logging

3. **Layout Components:**
   - `src/lib/components/layout/Header.svelte` - App header
   - `src/lib/components/layout/Navigation.svelte` - Sidebar navigation

4. **Root Layout:**
   - `src/routes/+layout.svelte` - Main layout with Header + Nav

5. **Placeholder Pages:**
   - `src/routes/+page.svelte` - Dashboard
   - `src/routes/scanner/+page.svelte` - Scanner
   - `src/routes/checklist/+page.svelte` - Checklist
   - `src/routes/sizing/+page.svelte` - Position Sizing
   - `src/routes/heat/+page.svelte` - Heat Check
   - `src/routes/entry/+page.svelte` - Trade Entry
   - `src/routes/calendar/+page.svelte` - Calendar

6. **Working Features:**
   - Theme toggle (day/night mode)
   - Theme persistence (localStorage)
   - Navigation between all screens
   - Frontend logging to console

---

## Troubleshooting

### Theme toggle doesn't work

**Problem:** Clicking button doesn't change theme
**Solution:**
- Check browser console for errors
- Verify `theme.ts` store logic
- Ensure `document.documentElement` is accessible
- Check CSS variables are defined in `app.css`

### Navigation doesn't highlight active route

**Problem:** Active route not showing gradient
**Solution:**
- Check `$page.url.pathname` in console
- Verify `isActive()` function logic
- Ensure gradient classes exist in Tailwind config

### Theme doesn't persist on reload

**Problem:** Theme resets to light mode
**Solution:**
- Check localStorage in DevTools (Application tab)
- Verify `localStorage.getItem('theme')` returns value
- Ensure `onMount` in layout runs

### Page transitions not smooth

**Problem:** No animations
**Solution:**
- Verify `page-transition` class in `app.css`
- Check CSS animation is defined
- Ensure TailwindCSS is processing the CSS

---

## Time Estimate

- **Theme Store:** 30 minutes
- **Logger Utility:** 30 minutes
- **Header Component:** 1 hour
- **Navigation Component:** 1-2 hours
- **Root Layout:** 30 minutes
- **Placeholder Routes:** 1 hour
- **Testing:** 1-2 hours
- **Integration with Go:** 30 minutes

**Total:** ~6-8 hours (1-2 days)

---

## References

- [SvelteKit Routing](https://kit.svelte.dev/docs/routing)
- [SvelteKit Stores](https://kit.svelte.dev/docs/state-management)
- [Svelte Transitions](https://svelte.dev/docs/svelte-transition)
- [roadmap.md - Step 6 Details](../plans/roadmap.md#step-6-application-layout--navigation)
- [overview-plan.md - Frontend Components](../plans/overview-plan.md#frontend-components-to-build)

---

## Next Step

Proceed to: **[Phase 1 - Step 7: Dashboard Screen](phase1-step7-dashboard.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
**Last Updated:** 2025-10-29
