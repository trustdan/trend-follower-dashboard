# Phase 4 - Step 20: Calendar View Implementation

**TF = Trend Following** - Systematic trading discipline enforcement system

**Phase:** 4 - Calendar & Polish
**Step:** 20 of 28
**Duration:** 3-4 days
**Dependencies:** Phase 3 complete (Heat Check & Trade Entry functional)

---

## Objectives

Build the Calendar screen showing a rolling 10-week view (2 weeks back + 8 weeks forward) with sectors as rows and weeks as columns. Each cell represents a sector × week combination and displays tickers active in that bucket during that week.

**Purpose:** Provide an at-a-glance view of sector diversification and help identify basket crowding (too many trades in one sector/week).

---

## Success Criteria

- [ ] Calendar displays 10-week grid (2 weeks historical + 8 weeks forward)
- [ ] Rows represent sector buckets
- [ ] Columns represent weeks (Monday-Sunday date ranges)
- [ ] Each cell shows tickers active in that sector/week
- [ ] Cells are color-coded based on heat level (low=green, medium=yellow, high=red)
- [ ] Current week is highlighted distinctly
- [ ] Tooltips show position details on hover (entry date, risk, expected exit)
- [ ] Calendar is responsive and horizontally scrollable if needed
- [ ] Legend explains color coding
- [ ] Calendar data fetches from `GET /api/calendar?weeks=10`

---

## Prerequisites

**Backend Requirements:**
- Positions table with sector, entry_date, and risk data
- API endpoint `GET /api/calendar?weeks=10` returning positions grouped by sector × week

**Frontend Requirements:**
- Navigation functional (can navigate to Calendar screen)
- Theme system working (day/night mode)
- Component library established (Card, Badge, etc.)
- Tooltip component available

---

## Implementation Plan

### Step 1: Backend API Endpoint (1-2 hours)

**File:** `backend/internal/server/calendar.go`

Create the calendar endpoint that returns positions grouped by sector and week.

```go
package server

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "your-module/internal/storage"
)

// CalendarResponse represents the calendar data structure
type CalendarResponse struct {
    Weeks   []WeekInfo        `json:"weeks"`
    Buckets []string          `json:"buckets"`
    Grid    map[string][]Cell `json:"grid"` // key: "bucket|week_index"
}

type WeekInfo struct {
    Index     int    `json:"index"`
    StartDate string `json:"start_date"` // "2025-10-27"
    EndDate   string `json:"end_date"`   // "2025-11-02"
    IsCurrent bool   `json:"is_current"`
}

type Cell struct {
    Ticker     string  `json:"ticker"`
    EntryDate  string  `json:"entry_date"`
    Risk       float64 `json:"risk"`
    ExitDate   string  `json:"exit_date,omitempty"` // Expected exit (if known)
    HeatLevel  string  `json:"heat_level"` // "low" | "medium" | "high"
}

func (s *Server) GetCalendar(c *gin.Context) {
    weeks := c.DefaultQuery("weeks", "10")
    weeksInt := 10 // parse weeks parameter

    // Calculate week boundaries
    now := time.Now()
    currentWeekStart := getMonday(now)

    // Build week info (2 back + 8 forward)
    weekInfos := make([]WeekInfo, weeksInt)
    for i := -2; i < 8; i++ {
        weekStart := currentWeekStart.AddDate(0, 0, i*7)
        weekEnd := weekStart.AddDate(0, 0, 6)
        weekInfos[i+2] = WeekInfo{
            Index:     i + 2,
            StartDate: weekStart.Format("2006-01-02"),
            EndDate:   weekEnd.Format("2006-01-02"),
            IsCurrent: i == 0,
        }
    }

    // Fetch all open positions
    positions, err := s.db.GetPositions()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Get unique buckets
    buckets := getUniqueBuckets(positions)

    // Build grid: group positions by bucket and week
    grid := make(map[string][]Cell)
    for _, pos := range positions {
        for idx, week := range weekInfos {
            if positionActiveInWeek(pos, week) {
                key := pos.Bucket + "|" + fmt.Sprintf("%d", idx)
                grid[key] = append(grid[key], Cell{
                    Ticker:    pos.Ticker,
                    EntryDate: pos.EntryDate,
                    Risk:      pos.Risk,
                    HeatLevel: calculateHeatLevel(pos.Risk),
                })
            }
        }
    }

    c.JSON(http.StatusOK, CalendarResponse{
        Weeks:   weekInfos,
        Buckets: buckets,
        Grid:    grid,
    })
}

// Helper: Get Monday of the week containing the given date
func getMonday(t time.Time) time.Time {
    offset := int(time.Monday - t.Weekday())
    if offset > 0 {
        offset = -6
    }
    return t.AddDate(0, 0, offset).Truncate(24 * time.Hour)
}

// Helper: Check if position is active during this week
func positionActiveInWeek(pos Position, week WeekInfo) bool {
    entryDate, _ := time.Parse("2006-01-02", pos.EntryDate)
    weekStart, _ := time.Parse("2006-01-02", week.StartDate)
    weekEnd, _ := time.Parse("2006-01-02", week.EndDate)

    // Position is active if entry <= weekEnd and (no exit or exit >= weekStart)
    if entryDate.After(weekEnd) {
        return false
    }
    // For now, assume all open positions are active (no exit date)
    return true
}

// Helper: Calculate heat level based on risk amount
func calculateHeatLevel(risk float64) string {
    // Define thresholds (customize as needed)
    if risk < 500 {
        return "low"
    } else if risk < 1000 {
        return "medium"
    }
    return "high"
}

func getUniqueBuckets(positions []Position) []string {
    bucketMap := make(map[string]bool)
    for _, pos := range positions {
        bucketMap[pos.Bucket] = true
    }
    buckets := make([]string, 0, len(bucketMap))
    for bucket := range bucketMap {
        buckets = append(buckets, bucket)
    }
    sort.Strings(buckets)
    return buckets
}
```

**Register route:**
```go
router.GET("/api/calendar", s.GetCalendar)
```

**Test:**
```bash
curl http://localhost:8080/api/calendar?weeks=10
```

---

### Step 2: Frontend Calendar Component (2-3 hours)

**File:** `ui/src/routes/calendar/+page.svelte`

```svelte
<script lang="ts">
    import { onMount } from 'svelte';
    import { api } from '$lib/api/client';
    import CalendarCell from '$lib/components/CalendarCell.svelte';
    import CalendarLegend from '$lib/components/CalendarLegend.svelte';
    import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

    let calendarData: any = null;
    let loading = true;
    let error = '';

    onMount(async () => {
        try {
            calendarData = await api.calendar.get(10);
        } catch (err) {
            error = err.message;
        } finally {
            loading = false;
        }
    });
</script>

<div class="calendar-page">
    <h1>Calendar View</h1>
    <p class="subtitle">Rolling 10-week sector diversification</p>

    {#if loading}
        <LoadingSpinner />
    {:else if error}
        <div class="error">{error}</div>
    {:else if calendarData}
        <CalendarLegend />

        <div class="calendar-grid-container">
            <div class="calendar-grid">
                <!-- Header row: week dates -->
                <div class="grid-header">
                    <div class="corner-cell">Sector</div>
                    {#each calendarData.weeks as week}
                        <div class="week-header" class:current={week.is_current}>
                            <div class="week-dates">
                                {formatDateShort(week.start_date)} - {formatDateShort(week.end_date)}
                            </div>
                            {#if week.is_current}
                                <div class="current-badge">Current</div>
                            {/if}
                        </div>
                    {/each}
                </div>

                <!-- Grid rows: one per sector bucket -->
                {#each calendarData.buckets as bucket}
                    <div class="grid-row">
                        <div class="bucket-label">{bucket}</div>
                        {#each calendarData.weeks as week, idx}
                            {@const key = `${bucket}|${idx}`}
                            {@const cells = calendarData.grid[key] || []}
                            <CalendarCell {cells} {bucket} {week} />
                        {/each}
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>

<style>
    .calendar-page {
        padding: var(--space-6);
    }

    h1 {
        font-size: var(--text-3xl);
        font-weight: 700;
        margin-bottom: var(--space-2);
    }

    .subtitle {
        color: var(--text-secondary);
        margin-bottom: var(--space-6);
    }

    .calendar-grid-container {
        overflow-x: auto;
        border-radius: 12px;
        border: 1px solid var(--border-color);
        background: var(--bg-secondary);
    }

    .calendar-grid {
        display: table;
        width: 100%;
        border-collapse: collapse;
    }

    .grid-header,
    .grid-row {
        display: table-row;
    }

    .corner-cell,
    .bucket-label,
    .week-header {
        display: table-cell;
        padding: var(--space-4);
        border: 1px solid var(--border-color);
        text-align: center;
        vertical-align: middle;
    }

    .corner-cell {
        background: var(--bg-tertiary);
        font-weight: 600;
        width: 120px;
    }

    .bucket-label {
        background: var(--bg-tertiary);
        font-weight: 500;
        text-align: left;
        width: 120px;
    }

    .week-header {
        background: var(--bg-primary);
        min-width: 100px;
    }

    .week-header.current {
        background: var(--gradient-blue);
        color: white;
    }

    .week-dates {
        font-size: var(--text-sm);
    }

    .current-badge {
        font-size: var(--text-xs);
        font-weight: 600;
        margin-top: var(--space-1);
    }

    .error {
        color: #DC2626;
        padding: var(--space-4);
        background: rgba(220, 38, 38, 0.1);
        border-radius: 8px;
    }
</style>
```

---

### Step 3: CalendarCell Component (1-2 hours)

**File:** `ui/src/lib/components/CalendarCell.svelte`

```svelte
<script lang="ts">
    export let cells: any[] = [];
    export let bucket: string;
    export let week: any;

    function getHeatClass(level: string): string {
        switch (level) {
            case 'low': return 'heat-low';
            case 'medium': return 'heat-medium';
            case 'high': return 'heat-high';
            default: return '';
        }
    }

    function getTotalHeat(): number {
        return cells.reduce((sum, cell) => sum + cell.risk, 0);
    }
</script>

<div class="calendar-cell" class:empty={cells.length === 0}>
    {#if cells.length > 0}
        <div class="cell-content {getHeatClass(cells[0]?.heat_level || 'low')}">
            {#each cells as cell}
                <div class="ticker-tag" title="Entry: {cell.entry_date}, Risk: ${cell.risk.toFixed(2)}">
                    {cell.ticker}
                </div>
            {/each}
            {#if cells.length > 1}
                <div class="heat-indicator">
                    Heat: ${getTotalHeat().toFixed(0)}
                </div>
            {/if}
        </div>
    {:else}
        <div class="empty-cell">—</div>
    {/if}
</div>

<style>
    .calendar-cell {
        display: table-cell;
        padding: var(--space-2);
        border: 1px solid var(--border-color);
        vertical-align: top;
        min-width: 100px;
        min-height: 60px;
    }

    .cell-content {
        display: flex;
        flex-direction: column;
        gap: var(--space-1);
        padding: var(--space-2);
        border-radius: 6px;
        transition: all 0.2s ease;
    }

    .heat-low {
        background: rgba(16, 185, 129, 0.1); /* Green tint */
        border-left: 3px solid #10B981;
    }

    .heat-medium {
        background: rgba(245, 158, 11, 0.1); /* Yellow tint */
        border-left: 3px solid #F59E0B;
    }

    .heat-high {
        background: rgba(220, 38, 38, 0.1); /* Red tint */
        border-left: 3px solid #DC2626;
    }

    .ticker-tag {
        font-size: var(--text-sm);
        font-weight: 500;
        padding: 2px 6px;
        background: var(--bg-primary);
        border-radius: 4px;
        cursor: pointer;
    }

    .ticker-tag:hover {
        background: var(--bg-tertiary);
    }

    .heat-indicator {
        font-size: var(--text-xs);
        color: var(--text-secondary);
        margin-top: var(--space-1);
    }

    .empty-cell {
        color: var(--text-tertiary);
        text-align: center;
        font-size: var(--text-lg);
    }
</style>
```

---

### Step 4: CalendarLegend Component (30 min)

**File:** `ui/src/lib/components/CalendarLegend.svelte`

```svelte
<div class="legend">
    <h3>Heat Levels</h3>
    <div class="legend-items">
        <div class="legend-item">
            <div class="color-box heat-low"></div>
            <span>Low (&lt; $500)</span>
        </div>
        <div class="legend-item">
            <div class="color-box heat-medium"></div>
            <span>Medium ($500-$1000)</span>
        </div>
        <div class="legend-item">
            <div class="color-box heat-high"></div>
            <span>High (&gt; $1000)</span>
        </div>
    </div>
</div>

<style>
    .legend {
        margin-bottom: var(--space-4);
        padding: var(--space-4);
        background: var(--bg-secondary);
        border-radius: 8px;
        border: 1px solid var(--border-color);
    }

    h3 {
        font-size: var(--text-lg);
        font-weight: 600;
        margin-bottom: var(--space-3);
    }

    .legend-items {
        display: flex;
        gap: var(--space-5);
    }

    .legend-item {
        display: flex;
        align-items: center;
        gap: var(--space-2);
    }

    .color-box {
        width: 24px;
        height: 24px;
        border-radius: 4px;
        border: 1px solid var(--border-color);
    }

    .heat-low {
        background: rgba(16, 185, 129, 0.3);
    }

    .heat-medium {
        background: rgba(245, 158, 11, 0.3);
    }

    .heat-high {
        background: rgba(220, 38, 38, 0.3);
    }
</style>
```

---

### Step 5: API Client Integration (15 min)

**File:** `ui/src/lib/api/client.ts`

Add calendar endpoint:

```typescript
export const api = {
    // ... existing endpoints

    calendar: {
        get: (weeks: number) => get(`/calendar?weeks=${weeks}`)
    }
};
```

---

### Step 6: Navigation Link (10 min)

**File:** `ui/src/lib/components/Navigation.svelte`

Add Calendar link:

```svelte
<nav>
    <!-- ... existing links -->
    <a href="/calendar" class:active={$page.url.pathname === '/calendar'}>
        📅 Calendar
    </a>
</nav>
```

---

## Testing Checklist

### Backend Tests

- [ ] `GET /api/calendar?weeks=10` returns valid JSON
- [ ] Week boundaries are calculated correctly (Monday-Sunday)
- [ ] Current week is marked with `is_current: true`
- [ ] Grid grouping works (positions appear in correct bucket × week cells)
- [ ] Heat levels are calculated correctly
- [ ] Empty weeks/buckets return empty arrays

### Frontend Tests

- [ ] Calendar page loads without errors
- [ ] Grid displays with correct number of rows (buckets) and columns (weeks)
- [ ] Current week is highlighted visually
- [ ] Ticker tags display in correct cells
- [ ] Color coding reflects heat levels (green/yellow/red)
- [ ] Tooltips show position details on hover
- [ ] Calendar is horizontally scrollable if content exceeds viewport
- [ ] Legend displays and explains color coding
- [ ] Theme toggle works (day/night mode)

### Integration Tests

- [ ] Complete workflow: Save GO decision → Navigate to Calendar → Position appears
- [ ] Multiple positions in same bucket/week stack correctly
- [ ] Empty buckets/weeks display "—" placeholder
- [ ] Calendar updates when new positions are added

---

## Troubleshooting

**Problem:** Calendar grid is too wide and breaks layout
**Solution:** Wrap in `.calendar-grid-container` with `overflow-x: auto`

**Problem:** Week boundaries are off by one day
**Solution:** Verify `getMonday()` calculation and timezone handling

**Problem:** Positions don't appear in expected weeks
**Solution:** Check `positionActiveInWeek()` logic; ensure date parsing is correct

**Problem:** Colors don't match heat levels
**Solution:** Review `calculateHeatLevel()` thresholds and CSS class names

---

## Documentation Requirements

- [ ] Update `docs/UI_QUICK_REFERENCE.md` with Calendar screen description
- [ ] Document calendar API endpoint in `docs/HTTP_CLI_PARITY.md`
- [ ] Add screenshots to User Guide (if applicable)
- [ ] Update `docs/PROGRESS.md` with completion status

---

## Next Steps

After completing Step 20:
1. Proceed to **Step 21: TradingView Integration**
2. Test calendar with multiple positions across different sectors
3. Gather feedback on color coding and heat thresholds
4. Consider enhancements (week numbers, month boundaries, export to CSV)

---

**Estimated Completion Time:** 3-4 days
**Phase 4 Progress:** 1 of 4 steps complete
**Overall Progress:** 20 of 28 steps complete (71%)

---

**End of Step 20**
