# Phase 4 - Step 23: Performance Optimization

**TF = Trend Following** - Systematic trading discipline enforcement system

**Phase:** 4 - Calendar & Polish
**Step:** 23 of 28
**Duration:** 1-2 days
**Dependencies:** Step 22 complete (UI Polish & Refinements complete)

---

## Objectives

Optimize the application for performance using data from logs collected in previous steps. Identify slow API endpoints, sluggish UI components, and long database queries. Implement caching, lazy loading, bundle optimization, and database indexing to ensure the app feels instant.

**Purpose:** Deliver a responsive, fast application that provides immediate feedback and minimal waiting time.

---

## Success Criteria

- [ ] Performance analysis complete (using logged metrics)
- [ ] API responses < 100ms for calculations (excluding network-dependent operations)
- [ ] Database queries < 50ms
- [ ] FINVIZ scan < 5s (network dependent, but optimized parsing)
- [ ] UI interactions feel instant (no perceived lag)
- [ ] Page load < 2 seconds
- [ ] Bundle size optimized (< 500KB for initial load)
- [ ] Lazy loading implemented for heavy components
- [ ] API caching implemented where appropriate
- [ ] Database indexes added for frequent queries
- [ ] Performance benchmarks documented (before/after)
- [ ] No memory leaks detected

---

## Prerequisites

**Completed:**
- Debug panel with performance logging (Step 22)
- All features functional
- Comprehensive logs available for analysis

**Tools:**
- Browser DevTools (Performance tab)
- Lighthouse for scoring
- Bundle analyzer (vite-plugin-analyze)

---

## Implementation Plan

### Part 1: Performance Analysis (2-3 hours)

#### Task 1.1: Review Performance Logs (1 hour)

Extract performance data from logs:

**File:** `scripts/analyze-performance.js` (Node.js script)

```javascript
const fs = require('fs');

// Load logs from exported JSON
const logs = JSON.parse(fs.readFileSync('logs/frontend-logs.json', 'utf8'));

// Analyze API call durations
const apiCalls = logs.filter(log => log.message.includes('API call'));
const slowCalls = apiCalls.filter(log => log.data?.duration > 500); // > 500ms

console.log('=== Slow API Calls (>500ms) ===');
slowCalls.forEach(call => {
    console.log(`${call.data.endpoint}: ${call.data.duration}ms`);
});

// Analyze component render times
const renders = logs.filter(log => log.message.includes('Component rendered'));
const slowRenders = renders.filter(log => log.data?.duration > 100); // > 100ms

console.log('\n=== Slow Component Renders (>100ms) ===');
slowRenders.forEach(render => {
    console.log(`${render.data.component}: ${render.data.duration}ms`);
});

// Analyze database query times (from backend logs)
const dbQueries = logs.filter(log => log.message.includes('DB query'));
const slowQueries = dbQueries.filter(log => log.data?.duration > 50); // > 50ms

console.log('\n=== Slow Database Queries (>50ms) ===');
slowQueries.forEach(query => {
    console.log(`${query.data.query}: ${query.data.duration}ms`);
});
```

Run analysis:
```bash
node scripts/analyze-performance.js > docs/PERFORMANCE_ANALYSIS.md
```

#### Task 1.2: Lighthouse Audit (30 min)

Run Lighthouse on all major screens:

```bash
# Install Lighthouse CLI
npm install -g lighthouse

# Audit each screen
lighthouse http://localhost:8080/dashboard --output=html --output-path=./docs/lighthouse-dashboard.html
lighthouse http://localhost:8080/checklist --output=html --output-path=./docs/lighthouse-checklist.html
lighthouse http://localhost:8080/calendar --output=html --output-path=./docs/lighthouse-calendar.html
```

Target scores:
- Performance: > 90
- Accessibility: > 90
- Best Practices: > 90
- SEO: > 80 (not critical for desktop app)

#### Task 1.3: Bundle Size Analysis (30 min)

**File:** `ui/vite.config.ts`

Add bundle analyzer:

```typescript
import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { visualizer } from 'rollup-plugin-visualizer';

export default defineConfig({
    plugins: [
        sveltekit(),
        visualizer({
            filename: './docs/bundle-stats.html',
            open: true,
            gzipSize: true,
            brotliSize: true
        })
    ]
});
```

Build and review:
```bash
cd ui/
npm run build
# Opens bundle-stats.html in browser
```

Identify large dependencies to split or replace.

---

### Part 2: Backend Optimization (3-4 hours)

#### Task 2.1: Database Indexing (1 hour)

**File:** `backend/internal/storage/migrations.go`

Add indexes for frequent queries:

```go
func AddPerformanceIndexes(db *sql.DB) error {
    indexes := []string{
        // Index on positions table
        "CREATE INDEX IF NOT EXISTS idx_positions_ticker ON positions(ticker)",
        "CREATE INDEX IF NOT EXISTS idx_positions_bucket ON positions(bucket)",
        "CREATE INDEX IF NOT EXISTS idx_positions_entry_date ON positions(entry_date)",

        // Index on candidates table
        "CREATE INDEX IF NOT EXISTS idx_candidates_date ON candidates(date)",
        "CREATE INDEX IF NOT EXISTS idx_candidates_ticker ON candidates(ticker)",

        // Index on decisions table
        "CREATE INDEX IF NOT EXISTS idx_decisions_timestamp ON decisions(timestamp DESC)",
        "CREATE INDEX IF NOT EXISTS idx_decisions_ticker ON decisions(ticker)",

        // Index on cooldowns table
        "CREATE INDEX IF NOT EXISTS idx_cooldowns_ticker ON cooldowns(ticker)",
        "CREATE INDEX IF NOT EXISTS idx_cooldowns_until ON cooldowns(until_date)",
    }

    for _, idx := range indexes {
        if _, err := db.Exec(idx); err != nil {
            return fmt.Errorf("failed to create index: %w", err)
        }
    }

    return nil
}
```

Run migration on app startup:

```go
func InitDB(path string) (*DB, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    // ... existing table creation

    // Add indexes
    if err := AddPerformanceIndexes(db); err != nil {
        return nil, err
    }

    return &DB{db: db}, nil
}
```

#### Task 2.2: Query Optimization (1-2 hours)

Optimize calendar query (identified as potentially slow):

**File:** `backend/internal/storage/calendar.go`

Before (N+1 query problem):
```go
func (db *DB) GetCalendarData(weeks int) ([]Position, error) {
    positions, err := db.GetAllPositions()
    for _, pos := range positions {
        // Additional queries for each position
        bucket := db.GetBucketForPosition(pos.ID) // N+1!
    }
}
```

After (single JOIN query):
```go
func (db *DB) GetCalendarData(weeks int) ([]Position, error) {
    query := `
        SELECT p.id, p.ticker, p.entry_date, p.risk, p.bucket, p.entry, p.current_stop
        FROM positions p
        WHERE p.status = 'open'
        ORDER BY p.entry_date DESC
    `
    rows, err := db.db.Query(query)
    // ... parse rows
    return positions, nil
}
```

#### Task 2.3: Response Caching (1 hour)

Cache settings to avoid repeated DB reads:

**File:** `backend/internal/server/cache.go`

```go
package server

import (
    "sync"
    "time"
)

type Cache struct {
    mu    sync.RWMutex
    data  map[string]cacheEntry
}

type cacheEntry struct {
    value      interface{}
    expiration time.Time
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]cacheEntry),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, exists := c.data[key]
    if !exists || time.Now().After(entry.expiration) {
        return nil, false
    }
    return entry.value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = cacheEntry{
        value:      value,
        expiration: time.Now().Add(ttl),
    }
}

func (c *Cache) Clear(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    delete(c.data, key)
}
```

**Usage in server:**

```go
type Server struct {
    db    *storage.DB
    cache *Cache
}

func (s *Server) GetSettings(c *gin.Context) {
    // Try cache first
    if settings, ok := s.cache.Get("settings"); ok {
        c.JSON(http.StatusOK, settings)
        return
    }

    // Cache miss - fetch from DB
    settings, err := s.db.GetSettings()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Cache for 5 minutes
    s.cache.Set("settings", settings, 5*time.Minute)

    c.JSON(http.StatusOK, settings)
}

func (s *Server) UpdateSettings(c *gin.Context) {
    // ... update settings in DB

    // Invalidate cache
    s.cache.Clear("settings")

    c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
```

---

### Part 3: Frontend Optimization (3-4 hours)

#### Task 3.1: Lazy Loading Routes (1 hour)

**File:** `ui/src/routes/+layout.svelte`

Implement route-based code splitting:

```svelte
<script lang="ts">
    import { page } from '$app/stores';
    import { onMount } from 'svelte';

    // Lazy load heavy components
    let CalendarComponent;
    let DashboardComponent;

    $: if ($page.url.pathname === '/calendar' && !CalendarComponent) {
        import('./calendar/+page.svelte').then(module => {
            CalendarComponent = module.default;
        });
    }

    $: if ($page.url.pathname === '/dashboard' && !DashboardComponent) {
        import('./dashboard/+page.svelte').then(module => {
            DashboardComponent = module.default;
        });
    }
</script>
```

**Note:** SvelteKit automatically code-splits routes. Verify in build output:

```bash
npm run build
# Check .svelte-kit/output/client/_app/immutable/chunks/ for split bundles
```

#### Task 3.2: Memoization for Expensive Calculations (1 hour)

**File:** `ui/src/lib/utils/memoize.ts`

```typescript
export function memoize<T extends (...args: any[]) => any>(fn: T): T {
    const cache = new Map();

    return ((...args: any[]) => {
        const key = JSON.stringify(args);
        if (cache.has(key)) {
            return cache.get(key);
        }
        const result = fn(...args);
        cache.set(key, result);
        return result;
    }) as T;
}
```

**Usage in Calendar:**

```svelte
<script lang="ts">
    import { memoize } from '$lib/utils/memoize';

    // Memoize expensive heat level calculation
    const calculateHeatLevel = memoize((risk: number): string => {
        if (risk < 500) return 'low';
        if (risk < 1000) return 'medium';
        return 'high';
    });

    $: cells.forEach(cell => {
        cell.heatLevel = calculateHeatLevel(cell.risk);
    });
</script>
```

#### Task 3.3: Virtual Scrolling for Long Lists (1-2 hours)

If candidate or decision lists get very long, implement virtual scrolling:

**File:** `ui/src/lib/components/VirtualList.svelte`

```svelte
<script lang="ts">
    export let items: any[] = [];
    export let itemHeight = 50; // px
    export let visibleCount = 10;

    let scrollTop = 0;
    let containerHeight = itemHeight * visibleCount;

    $: startIndex = Math.floor(scrollTop / itemHeight);
    $: endIndex = Math.min(startIndex + visibleCount + 1, items.length);
    $: visibleItems = items.slice(startIndex, endIndex);
    $: offsetY = startIndex * itemHeight;

    function handleScroll(e: Event) {
        scrollTop = (e.target as HTMLElement).scrollTop;
    }
</script>

<div class="virtual-list-container" style="height: {containerHeight}px" on:scroll={handleScroll}>
    <div class="virtual-list-spacer" style="height: {items.length * itemHeight}px">
        <div class="virtual-list-items" style="transform: translateY({offsetY}px)">
            {#each visibleItems as item, i (startIndex + i)}
                <div class="virtual-list-item" style="height: {itemHeight}px">
                    <slot {item} index={startIndex + i} />
                </div>
            {/each}
        </div>
    </div>
</div>

<style>
    .virtual-list-container {
        overflow-y: auto;
        position: relative;
    }

    .virtual-list-spacer {
        position: relative;
    }

    .virtual-list-items {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
    }
</style>
```

**Usage:**

```svelte
<VirtualList items={candidates} itemHeight={60}>
    <div slot="item" let:item>
        <CandidateRow candidate={item} />
    </div>
</VirtualList>
```

#### Task 3.4: Image/Icon Optimization (30 min)

Replace PNG icons with SVG:

**File:** `ui/src/lib/icons/index.ts`

```typescript
export const icons = {
    chart: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path d="M3 3h18v18H3z" stroke-width="2"/>
    </svg>`,

    check: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path d="M5 12l5 5L20 7" stroke-width="2"/>
    </svg>`,

    // ... more icons
};
```

**Component:**

```svelte
<script lang="ts">
    import { icons } from '$lib/icons';
    export let name: string;
    export let size = 24;
</script>

<span class="icon" style="width: {size}px; height: {size}px">
    {@html icons[name]}
</span>

<style>
    .icon {
        display: inline-block;
    }
    .icon :global(svg) {
        width: 100%;
        height: 100%;
    }
</style>
```

---

### Part 4: Performance Monitoring (1-2 hours)

#### Task 4.1: Performance Overlay (Optional) (1 hour)

**File:** `ui/src/lib/components/PerformanceOverlay.svelte`

```svelte
<script lang="ts">
    import { onMount } from 'svelte';

    let fps = 0;
    let memoryUsage = 0;
    let visible = false;

    let lastTime = performance.now();
    let frameCount = 0;

    function measurePerformance() {
        frameCount++;
        const now = performance.now();
        const delta = now - lastTime;

        if (delta >= 1000) {
            fps = Math.round((frameCount * 1000) / delta);
            frameCount = 0;
            lastTime = now;

            // Memory (if available)
            if (performance.memory) {
                memoryUsage = Math.round(performance.memory.usedJSHeapSize / 1048576); // MB
            }
        }

        requestAnimationFrame(measurePerformance);
    }

    onMount(() => {
        measurePerformance();
    });

    // Toggle with Ctrl+Shift+P
    if (typeof window !== 'undefined') {
        window.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.shiftKey && e.key === 'P') {
                visible = !visible;
            }
        });
    }
</script>

{#if visible && import.meta.env.DEV}
    <div class="perf-overlay">
        <div class="perf-stat">
            <span class="label">FPS:</span>
            <span class="value" class:good={fps >= 55} class:bad={fps < 55}>{fps}</span>
        </div>
        {#if memoryUsage > 0}
            <div class="perf-stat">
                <span class="label">Memory:</span>
                <span class="value">{memoryUsage} MB</span>
            </div>
        {/if}
    </div>
{/if}

<style>
    .perf-overlay {
        position: fixed;
        top: 10px;
        left: 10px;
        background: rgba(0, 0, 0, 0.8);
        color: white;
        padding: 8px 12px;
        border-radius: 6px;
        font-family: monospace;
        font-size: 12px;
        z-index: 9999;
        display: flex;
        gap: 16px;
    }

    .perf-stat {
        display: flex;
        gap: 6px;
    }

    .label {
        opacity: 0.7;
    }

    .value {
        font-weight: bold;
    }

    .value.good {
        color: #10B981;
    }

    .value.bad {
        color: #DC2626;
    }
</style>
```

#### Task 4.2: API Call Timing Middleware (30 min)

**File:** `ui/src/lib/api/client.ts`

Add automatic timing to all API calls:

```typescript
export async function get(endpoint: string) {
    const start = performance.now();

    try {
        const res = await fetch(`${API_BASE}${endpoint}`);
        const duration = performance.now() - start;

        logger.info('API call', {
            endpoint,
            method: 'GET',
            status: res.status,
            duration: Math.round(duration)
        });

        if (!res.ok) throw new Error(await res.text());
        return res.json();
    } catch (err) {
        const duration = performance.now() - start;
        logger.error('API call failed', {
            endpoint,
            method: 'GET',
            duration: Math.round(duration),
            error: err.message
        });
        throw err;
    }
}
```

---

### Part 5: Benchmarking & Documentation (2 hours)

#### Task 5.1: Before/After Benchmarks (1 hour)

**File:** `docs/PERFORMANCE_BENCHMARKS.md`

```markdown
# Performance Benchmarks

## Methodology

All tests performed on:
- OS: Windows 10
- CPU: Intel i7-9700K
- RAM: 16GB
- Browser: Chrome 120

## Results

### API Response Times

| Endpoint | Before | After | Improvement |
|----------|--------|-------|-------------|
| GET /api/settings | 120ms | 45ms | 62% faster |
| GET /api/positions | 280ms | 95ms | 66% faster |
| GET /api/calendar?weeks=10 | 450ms | 180ms | 60% faster |
| POST /api/size/calculate | 85ms | 72ms | 15% faster |
| POST /api/heat/check | 110ms | 65ms | 41% faster |

### Database Query Times

| Query | Before | After | Improvement |
|-------|--------|-------|-------------|
| Get all positions | 95ms | 25ms | 74% faster |
| Get candidates by date | 120ms | 30ms | 75% faster |
| Calendar data fetch | 200ms | 60ms | 70% faster |

### Frontend Metrics

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Initial page load | 2.8s | 1.6s | 43% faster |
| Dashboard render | 320ms | 140ms | 56% faster |
| Calendar render | 650ms | 280ms | 57% faster |
| Checklist interaction | 50ms | 30ms | 40% faster |

### Bundle Size

| Bundle | Before | After | Reduction |
|--------|--------|-------|-----------|
| Initial load | 680KB | 420KB | 38% smaller |
| Dashboard chunk | 180KB | 120KB | 33% smaller |
| Calendar chunk | 220KB | 140KB | 36% smaller |

### Lighthouse Scores

| Screen | Performance | Accessibility | Best Practices |
|--------|-------------|---------------|----------------|
| Dashboard | 92 | 95 | 100 |
| Checklist | 94 | 96 | 100 |
| Calendar | 88 | 94 | 100 |

## Key Optimizations Applied

1. **Database indexing** on frequently queried columns
2. **Response caching** for settings (5-min TTL)
3. **Query optimization** to eliminate N+1 queries
4. **Lazy loading** for heavy route components
5. **Memoization** for expensive calculations
6. **SVG icons** instead of PNG (smaller, scalable)
7. **Bundle splitting** via SvelteKit automatic code splitting
```

---

## Testing Checklist

### Backend Performance
- [ ] All API endpoints respond < 100ms (except FINVIZ scan)
- [ ] Database queries complete < 50ms
- [ ] Indexes created successfully
- [ ] Cache invalidation works correctly
- [ ] No N+1 query problems remain

### Frontend Performance
- [ ] Initial page load < 2s
- [ ] Route transitions < 200ms
- [ ] No layout shifts (CLS = 0)
- [ ] Smooth 60fps animations
- [ ] No memory leaks after 30 min usage
- [ ] Bundle size < 500KB for initial load

### Lighthouse Scores
- [ ] Performance > 90 on all screens
- [ ] Accessibility > 90 on all screens
- [ ] Best Practices > 90 on all screens

### Real-World Testing
- [ ] Test on lower-end hardware (if available)
- [ ] Test with slow network (throttle to 3G)
- [ ] Test with large datasets (100+ candidates, 50+ positions)
- [ ] Monitor memory usage over extended session (1 hour)

---

## Troubleshooting

**Problem:** Database queries still slow despite indexes
**Solution:** Run `ANALYZE` on SQLite database to update query planner statistics

**Problem:** Bundle size still too large
**Solution:** Use bundle analyzer to identify heavy dependencies; consider alternatives

**Problem:** Calendar render is slow with many positions
**Solution:** Implement virtual scrolling or pagination for calendar cells

**Problem:** Memory usage grows over time
**Solution:** Check for event listener leaks; ensure components unmount cleanly

---

## Documentation Requirements

- [ ] Create `docs/PERFORMANCE_BENCHMARKS.md` with before/after metrics
- [ ] Update `docs/PROGRESS.md` with completion status
- [ ] Document caching strategy in API reference
- [ ] Add performance best practices to developer guide

---

## Next Steps

After completing Step 23:
1. **Phase 4 Complete!** Calendar & Polish finished
2. Proceed to **Phase 5: Testing & Packaging**
3. Begin comprehensive testing suite (Step 24)
4. Prepare for Windows deployment

---

**Estimated Completion Time:** 1-2 days
**Phase 4 Progress:** 4 of 4 steps complete ✓
**Overall Progress:** 23 of 28 steps complete (82%)

---

**End of Step 23 - Phase 4 Complete!**
