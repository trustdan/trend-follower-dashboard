# M23: Heat Analysis Command Implementation Plan

## Overview
Implement the `tf-engine heat` command to analyze historical position performance and generate heat scores based on P&L, win rate, and trade count. This will enable traders to identify the best-performing position sizes for each signal type.

**Estimated Time:** ~5 hours
**Dependencies:** M1-M22 (all completed)
**Blocks:** 4 heat analysis tests currently skipped

---

## Phase 1: Core Heat Calculation Engine (~2 hours)

### 1.1 Heat Metrics Calculation
**File:** `cmd/heat.go` (new)
**Time:** 45 minutes

```go
type HeatMetrics struct {
    SignalType    string
    PositionSize  int
    TotalTrades   int
    WinningTrades int
    LosingTrades  int
    TotalPnL      float64
    AvgWin        float64
    AvgLoss       float64
    WinRate       float64
    HeatScore     float64  // Composite score
}

func calculateHeatMetrics(decisions []Decision) map[string]map[int]*HeatMetrics
```

**Key Logic:**
- Group decisions by signal type and position size
- Calculate win/loss stats for each combination
- Compute heat score: `HeatScore = (WinRate * 0.4) + (AvgPnL * 0.4) + (TradeCount * 0.2)`
- Handle edge cases (no trades, all losses, etc.)

### 1.2 Decision History Loading
**File:** `cmd/heat.go`
**Time:** 30 minutes

```go
func loadDecisionHistory(dataDir string) ([]Decision, error) {
    // Read all decision JSON files
    // Parse and aggregate into decision slice
    // Sort by timestamp for chronological analysis
}
```

**Requirements:**
- Read from `~/.tf-engine/data/decisions/`
- Parse JSON files matching pattern `decision_YYYYMMDD_HHMMSS.json`
- Handle missing/corrupted files gracefully
- Support date range filtering (future enhancement)

### 1.3 Heat Score Algorithm
**File:** `pkg/analysis/heat.go` (new)
**Time:** 45 minutes

```go
type HeatScorer interface {
    CalculateScore(metrics *HeatMetrics) float64
    NormalizeScores(scores map[string]float64) map[string]float64
}

type DefaultHeatScorer struct {
    WinRateWeight   float64  // 0.4
    PnLWeight       float64  // 0.4
    VolumeWeight    float64  // 0.2
}
```

**Algorithm:**
1. Normalize each metric to 0-100 scale
2. Apply weights: `score = (winRate * 0.4) + (pnl * 0.4) + (volume * 0.2)`
3. Scale final score to 0-100
4. Handle minimum trade thresholds

---

## Phase 2: CLI Command & Output (~1.5 hours)

### 2.1 Command Structure
**File:** `cmd/heat.go`
**Time:** 30 minutes

```go
var heatCmd = &cobra.Command{
    Use:   "heat",
    Short: "Analyze position size performance (heat analysis)",
    Long: `Analyzes historical trading decisions to calculate "heat scores"
for each signal type and position size combination.

Heat scores range from 0-100 and indicate:
- Win rate performance (40%)
- Average P&L (40%)
- Trade volume/confidence (20%)

Higher scores indicate better historical performance.`,
    RunE: runHeat,
}

// Flags:
// --min-trades int (default 5) - Minimum trades required for analysis
// --top int (default 3) - Show top N position sizes per signal
// --format string (default "table") - Output format: table, json, csv
```

### 2.2 Table Output Formatting
**File:** `cmd/heat.go`
**Time:** 45 minutes

```
Heat Analysis Report
Generated: 2025-10-28 14:30:00
Data Range: 2025-10-01 to 2025-10-28

SIGNAL: BULLISH_ENGULFING
┌──────────┬────────┬──────┬─────────┬──────────┬──────────┬───────────┐
│ Position │ Trades │ W/L  │ Win %   │ Avg P&L  │ Total    │ Heat      │
│ Size     │        │      │         │          │ P&L      │ Score     │
├──────────┼────────┼──────┼─────────┼──────────┼──────────┼───────────┤
│ 500 ⭐   │ 23     │ 18/5 │ 78.3%   │ $127.50  │ $2,932.50│ 🔥 92.5   │
│ 400      │ 31     │ 19/12│ 61.3%   │ $89.25   │ $2,766.75│ 🟠 78.2   │
│ 300      │ 45     │ 24/21│ 53.3%   │ $45.10   │ $2,029.50│ 🟡 65.1   │
└──────────┴────────┴──────┴─────────┴──────────┴──────────┴───────────┘

SIGNAL: HAMMER
┌──────────┬────────┬──────┬─────────┬──────────┬──────────┬───────────┐
│ Position │ Trades │ W/L  │ Win %   │ Avg P&L  │ Total    │ Heat      │
│ Size     │        │      │         │          │ P&L      │ Score     │
├──────────┼────────┼──────┼─────────┼──────────┼──────────┼───────────┤
│ 300 ⭐   │ 18     │ 12/6 │ 66.7%   │ $92.10   │ $1,657.80│ 🔥 85.3   │
│ 200      │ 25     │ 14/11│ 56.0%   │ $54.20   │ $1,355.00│ 🟠 71.8   │
│ 100      │ 34     │ 17/17│ 50.0%   │ $23.50   │ $799.00  │ 🟡 58.9   │
└──────────┴────────┴──────┴─────────┴──────────┴──────────┴───────────┘

📊 Summary:
- Total signals analyzed: 4
- Total trades: 176
- Overall win rate: 61.4%
- Best performer: BULLISH_ENGULFING @ 500 shares (Heat: 92.5)

⭐ = Recommended position size for this signal
```

### 2.3 JSON/CSV Export
**File:** `cmd/heat.go`
**Time:** 15 minutes

```json
{
  "generated_at": "2025-10-28T14:30:00Z",
  "date_range": {
    "start": "2025-10-01",
    "end": "2025-10-28"
  },
  "signals": {
    "BULLISH_ENGULFING": [
      {
        "position_size": 500,
        "trades": 23,
        "wins": 18,
        "losses": 5,
        "win_rate": 0.783,
        "avg_pnl": 127.50,
        "total_pnl": 2932.50,
        "heat_score": 92.5,
        "recommended": true
      }
    ]
  }
}
```

---

## Phase 3: Integration & Testing (~1.5 hours)

### 3.1 Update Main Command
**File:** `cmd/root.go`
**Time:** 10 minutes

```go
func init() {
    rootCmd.AddCommand(heatCmd)
    rootCmd.AddCommand(sizeCmd)
    rootCmd.AddCommand(checklistCmd)
    rootCmd.AddCommand(saveDecisionCmd)
}
```

### 3.2 Automated Tests
**File:** `tests/test_heat.sh`
**Time:** 45 minutes

```bash
#!/bin/bash
# Test cases:
# 1. Heat calculation with sample data
# 2. Minimum trade threshold filtering
# 3. Top N recommendations
# 4. JSON/CSV output validation
# 5. Empty decision history handling
# 6. Single signal type analysis
```

**Test Data:**
- Create `tests/fixtures/decisions/` with sample decision files
- Mix of winning/losing trades
- Multiple signal types and position sizes
- Edge cases (all wins, all losses, exactly min trades)

### 3.3 Documentation
**File:** `README.md` update
**Time:** 15 minutes

Add section:
```markdown
### Heat Analysis

Analyze historical performance to find optimal position sizes:

```bash
# Basic heat analysis
./tf-engine heat

# Require minimum 10 trades per combination
./tf-engine heat --min-trades 10

# Show top 5 position sizes per signal
./tf-engine heat --top 5

# Export to JSON
./tf-engine heat --format json > heat_report.json
```

Heat scores combine:
- **Win Rate (40%)**: Percentage of winning trades
- **Average P&L (40%)**: Mean profit per trade
- **Trade Volume (20%)**: Number of trades (confidence indicator)

Position sizes marked with ⭐ are recommended based on highest heat score.
```

### 3.4 Integration Testing
**File:** Manual testing checklist
**Time:** 20 minutes

- [ ] Run heat with fresh decision history
- [ ] Verify scores match manual calculations
- [ ] Test with 100+ historical decisions
- [ ] Validate recommendations make sense
- [ ] Check all output formats
- [ ] Test edge cases (no data, single trade, etc.)

---

## Phase 4: Success Criteria & Validation (~30 minutes)

### Acceptance Tests
1. **Heat calculation accuracy**
   - Manually verify 3 signal/size combinations
   - Heat scores within 0.1 of expected values

2. **CLI usability**
   - Command completes in < 2 seconds for 1000 decisions
   - Output is readable and aligned
   - Help text is clear

3. **Blocked tests unblocked**
   - `test_heat.sh` passes all 4 test cases
   - No skipped tests remain for heat functionality

4. **Documentation complete**
   - README explains heat analysis
   - Examples show common usage patterns
   - Algorithm is documented

### Definition of Done
- [ ] `tf-engine heat` command runs successfully
- [ ] All 4 heat tests pass
- [ ] Heat scores calculated correctly (verified manually)
- [ ] Table output is formatted and readable
- [ ] JSON export works
- [ ] README updated with heat section
- [ ] No regressions in existing tests

---

## File Structure After M23

```
excel-trading-platform/
├── cmd/
│   ├── heat.go          (NEW - ~300 lines)
│   ├── root.go          (UPDATED - add heat command)
│   ├── size.go          (existing)
│   └── checklist.go     (existing)
├── pkg/
│   ├── analysis/
│   │   └── heat.go      (NEW - heat scoring logic)
│   └── models/
│       └── decision.go  (existing)
├── tests/
│   ├── test_heat.sh     (UPDATED - unskip tests)
│   └── fixtures/
│       └── decisions/   (NEW - test data)
└── docs/
    └── M23_HEAT_COMMAND_PLAN.md (this file)
```

---

## Risk Mitigation

### Risk: Heat scores don't match trader intuition
**Mitigation:**
- Make weights configurable
- Show raw metrics alongside heat score
- Allow filtering by date range

### Risk: Large decision history causes slow performance
**Mitigation:**
- Add progress indicator for > 1000 decisions
- Implement date range filtering
- Consider caching/indexing for future enhancement

### Risk: Not enough historical data for meaningful analysis
**Mitigation:**
- Require minimum trades (default 5)
- Show warning when data is limited
- Provide "insufficient data" message instead of misleading scores

---

## Next Steps After M23
1. Complete M24 (UI implementation)
2. Run full integration testing
3. Windows cross-platform validation
4. User acceptance testing with real trading data
5. Production deployment preparation

---

**Plan Status:** Ready for implementation
**Created:** 2025-10-28
**Dependencies:** All clear (M1-M22 complete)
