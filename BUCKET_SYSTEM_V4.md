# TF-Engine Bucket System v4 - Complete Implementation

**Date:** October 30, 2025
**Build:** ui/tf-gui-v4.exe
**Status:** ✅ Complete bucket grouping system based on sector correlation

---

## Summary

Implemented comprehensive bucket grouping system based on user specifications:
- **12 individual sectors** (FINVIZ order + ETFs)
- **8 heat tracking buckets** (4 grouped + 4 standalone)
- Auto-mapping from sector to bucket
- Prevents concentration in correlated sectors

---

## The Bucket Strategy

### Design Philosophy

**Purpose:** Prevent over-concentration in CORRELATED sectors that crash together

**Example:**
- With $100,000 account and 1.5% bucket cap = $1,500 max per bucket
- If you have $1,500 in "Tech/Comm" bucket (2 positions), you CANNOT add another Technology OR Communication Services position
- Forces diversification into uncorrelated buckets (Healthcare, Energy, Real Estate, etc.)

### Why These Groupings?

**1. Materials/Industrials** (Basic Materials + Industrials)
- **Correlation:** Both driven by economic cycles and commodity prices
- **Example:** When manufacturing slows, steel (Basic Materials) and factories (Industrials) both suffer
- **Max exposure:** 1.5% of equity

**2. Tech/Comm** (Communication Services + Technology)
- **Correlation:** Both growth/innovation sectors, share investor sentiment
- **Example:** When tech growth story falters, Meta (Comm) and NVDA (Tech) drop together
- **Max exposure:** 1.5% of equity

**3. Financial/Cyclical** (Financial + Consumer Cyclical)
- **Correlation:** Both economically sensitive, benefit from consumer spending and credit growth
- **Example:** When economy booms, banks (Financial) and retailers (Cyclical) thrive together
- **Max exposure:** 1.5% of equity

**4. Defensive/Utilities** (Consumer Defensive + Utilities)
- **Correlation:** Both defensive sectors with stable dividends, low volatility
- **Example:** During market fear, investors flee to safety in both (PG and utilities)
- **Max exposure:** 1.5% of equity

**5. Energy** (standalone)
- **Rationale:** Short burst outperformance, commodity-driven, less correlated to market cycles
- **Example:** Oil spikes are independent of tech growth or defensive plays
- **Max exposure:** 1.5% of equity

**6. Healthcare** (standalone)
- **Rationale:** Less correlated to economic cycles, driven by demographics/innovation
- **Example:** Healthcare performs well in both bull and bear markets
- **Max exposure:** 1.5% of equity

**7. Real Estate** (standalone)
- **Rationale:** Interest rate sensitivity, unique risk profile
- **Example:** REITs move based on rate expectations, not economic growth
- **Max exposure:** 1.5% of equity

**8. ETFs** (standalone)
- **Rationale:** ETFs can span multiple sectors, need separate tracking
- **Example:** XLE (Energy ETF) vs XLK (Tech ETF) have different exposures
- **Max exposure:** 1.5% of equity

---

## Sector Classification (FINVIZ Order)

**All 12 sectors available in dropdown:**

1. **Basic Materials** → Materials/Industrials bucket
2. **Communication Services** → Tech/Comm bucket
3. **Consumer Cyclical** → Financial/Cyclical bucket
4. **Consumer Defensive** → Defensive/Utilities bucket
5. **Energy** → Energy bucket
6. **Financial** → Financial/Cyclical bucket
7. **Healthcare** → Healthcare bucket
8. **Industrials** → Materials/Industrials bucket
9. **Real Estate** → Real Estate bucket
10. **Technology** → Tech/Comm bucket
11. **Utilities** → Defensive/Utilities bucket
12. **ETFs** → ETFs bucket

---

## Auto-Mapping Table

| Sector | Maps To Bucket | Shares Bucket With |
|--------|----------------|-------------------|
| Basic Materials | Materials/Industrials | Industrials |
| Communication Services | Tech/Comm | Technology |
| Consumer Cyclical | Financial/Cyclical | Financial |
| Consumer Defensive | Defensive/Utilities | Utilities |
| Energy | Energy | (standalone) |
| Financial | Financial/Cyclical | Consumer Cyclical |
| Healthcare | Healthcare | (standalone) |
| Industrials | Materials/Industrials | Basic Materials |
| Real Estate | Real Estate | (standalone) |
| Technology | Tech/Comm | Communication Services |
| Utilities | Defensive/Utilities | Consumer Defensive |
| ETFs | ETFs | (standalone) |

---

## How It Works

### Trade Entry Screen

**Dual Dropdown System:**

1. **Sector Dropdown** - Select individual sector classification
   - Shows all 12 FINVIZ sectors + ETFs
   - In FINVIZ order (top to bottom)

2. **Bucket Dropdown** - Shows grouped bucket for heat tracking
   - Auto-populates based on sector selection
   - Can be manually overridden if needed
   - Used for heat cap calculations

### Example Workflow

**Scenario 1: Adding Technology Position**

1. User selects: **Sector:** Technology
2. System auto-selects: **Bucket:** Tech/Comm
3. User enters: **Risk:** $750
4. Click "Check All 5 Gates"
5. System checks: Current Tech/Comm bucket heat + $750 ≤ $1,500 cap
6. **Result:** If Tech/Comm bucket already has $800, new total would be $1,550
7. **Gate 4 FAILS:** "Bucket heat exceeds cap by $50"

**Scenario 2: Diversification Example**

**Current Portfolio:**
- NVDA (Technology) - $750 → Tech/Comm bucket
- META (Communication Services) - $750 → Tech/Comm bucket
- **Tech/Comm bucket:** $1,500 / $1,500 (100% - AT CAP!)

**User tries to add MSFT (Technology):**
- Would add $750 to Tech/Comm bucket
- New total: $2,250 (150% of cap)
- **REJECTED** - exceeds bucket cap by $750

**User must choose different sector:**
- ✅ UNH (Healthcare) → Healthcare bucket (uncorrelated)
- ✅ XLE (Energy) → Energy bucket (uncorrelated)
- ✅ JPM (Financial) → Financial/Cyclical bucket (uncorrelated)
- ✅ WMT (Consumer Defensive) → Defensive/Utilities bucket (uncorrelated)

**Result:** True diversification across uncorrelated sectors!

---

## Calendar View Integration

**10-Week Rolling Calendar**
- Rows: 8 bucket categories
- Columns: 10 weeks (2 back + 8 forward)
- Cells: Count of positions in that bucket/week

**Visual Diversification Check:**
```
Bucket                  | Week1 | Week2 | Week3 | Week4 | ...
------------------------|-------|-------|-------|-------|----
Materials/Industrials   |   1   |   -   |   -   |   2   | ...
Tech/Comm               |   2   |   1   |   -   |   1   | ...
Financial/Cyclical      |   -   |   1   |   1   |   -   | ...
Defensive/Utilities     |   1   |   -   |   1   |   -   | ...
Energy                  |   -   |   -   |   1   |   -   | ...
Healthcare              |   1   |   1   |   -   |   1   | ...
Real Estate             |   -   |   -   |   -   |   1   | ...
ETFs                    |   -   |   1   |   -   |   -   | ...
```

**Goals:**
- No bucket should dominate (avoid 3+ positions in one bucket)
- Spread entries across multiple weeks (avoid clustering)
- Identify gaps (weeks/buckets with no coverage)

---

## Heat Cap Strategy

### Portfolio Level (4% cap)
- **Purpose:** Limit total risk across ALL positions
- **Example:** $100,000 account → $4,000 max total heat
- **Allows:** ~5-6 positions at ~$750 each

### Bucket Level (1.5% cap each)
- **Purpose:** Limit concentration in correlated sectors
- **Example:** $100,000 account → $1,500 max per bucket
- **Allows:** ~2 positions per bucket at ~$750 each

### Enforcement
- **Cannot bypass:** Gate 4 fails if either cap exceeded
- **No exceptions:** Cannot "temporarily" exceed caps
- **Automatic:** System checks before allowing GO decision

### Why Both Caps?

**Portfolio cap alone is insufficient:**
- Could have 5 positions all in Tech/Comm ($3,750 in one correlated bucket)
- When tech crashes, entire portfolio crashes together
- Defeats purpose of diversification

**Bucket caps force true diversification:**
- Max 2 positions in Tech/Comm ($1,500)
- Must spread across Healthcare, Energy, Financials, etc.
- Reduces portfolio volatility
- Survives sector-specific crashes

---

## Benefits

### 1. Prevents Correlated Risk
**Before bucket system:**
- Could have: NVDA (Tech), META (Comm), GOOGL (Comm), MSFT (Tech)
- All crash together when growth story fails
- Total loss: 4 positions × $750 = $3,000 (3% of account)

**With bucket system:**
- Tech/Comm limited to $1,500 (2 positions)
- Must add: UNH (Healthcare), XLE (Energy), JPM (Financial)
- Tech crash only impacts 2 positions
- Max loss from Tech/Comm: $1,500 (1.5% of account)

### 2. Enforces Market Diversification
- Cannot load up on growth (Tech/Comm only)
- Cannot hide in defensive (Defensive/Utilities only)
- Must hold positions across market spectrum
- Reduces dependency on single market narrative

### 3. Limits Commodity/Cycle Risk
- Materials + Industrials grouped (both commodity/cycle driven)
- Max exposure to economic cycle risk: 1.5%
- Prevents betting too heavily on manufacturing boom

### 4. Separates Independent Sectors
- Healthcare standalone (demographics/innovation driven)
- Energy standalone (commodity/geopolitical driven)
- Real Estate standalone (interest rate driven)
- Each can be independently sized to cap

### 5. ETF Segregation
- ETFs tracked separately
- Prevents double-counting (XLK overlaps with MSFT)
- Allows sector ETF strategies without confusion

---

## Implementation Details

### Files Modified

**1. ui/trade_entry.go (Lines 48-105)**
- Updated sector list to 12 FINVIZ sectors + ETFs
- Updated bucket list to 8 grouped buckets
- Created sectorToBucket mapping with new groupings
- Auto-mapping on sector selection

**2. ui/calendar.go (Lines 26-36)**
- Updated bucket list to match trade entry
- Calendar rows now show 8 buckets
- Consistent naming across app

### Code Structure

**Sector Dropdown:**
```go
sectorOptions := []string{
    "Basic Materials",
    "Communication Services",
    "Consumer Cyclical",
    "Consumer Defensive",
    "Energy",
    "Financial",
    "Healthcare",
    "Industrials",
    "Real Estate",
    "Technology",
    "Utilities",
    "ETFs",
}
```

**Bucket Dropdown:**
```go
bucketOptions := []string{
    "Materials/Industrials",  // Basic Materials + Industrials
    "Tech/Comm",              // Communication Services + Technology
    "Financial/Cyclical",     // Financial + Consumer Cyclical
    "Defensive/Utilities",    // Consumer Defensive + Utilities
    "Energy",                 // Energy (standalone)
    "Healthcare",             // Healthcare (standalone)
    "Real Estate",            // Real Estate (standalone)
    "ETFs",                   // ETFs (standalone)
}
```

**Auto-Mapping:**
```go
sectorToBucket := map[string]string{
    "Basic Materials":        "Materials/Industrials",
    "Industrials":            "Materials/Industrials",
    "Communication Services": "Tech/Comm",
    "Technology":             "Tech/Comm",
    "Financial":              "Financial/Cyclical",
    "Consumer Cyclical":      "Financial/Cyclical",
    "Consumer Defensive":     "Defensive/Utilities",
    "Utilities":              "Defensive/Utilities",
    "Energy":                 "Energy",
    "Healthcare":             "Healthcare",
    "Real Estate":            "Real Estate",
    "ETFs":                   "ETFs",
}
```

---

## Testing Checklist

### Sector-to-Bucket Mapping
- [ ] Navigate to Trade Entry screen
- [ ] Select each sector, verify correct bucket auto-populates:
  - [ ] Basic Materials → Materials/Industrials
  - [ ] Communication Services → Tech/Comm
  - [ ] Consumer Cyclical → Financial/Cyclical
  - [ ] Consumer Defensive → Defensive/Utilities
  - [ ] Energy → Energy
  - [ ] Financial → Financial/Cyclical
  - [ ] Healthcare → Healthcare
  - [ ] Industrials → Materials/Industrials
  - [ ] Real Estate → Real Estate
  - [ ] Technology → Tech/Comm
  - [ ] Utilities → Defensive/Utilities
  - [ ] ETFs → ETFs

### Bucket Grouping Logic
- [ ] Select Technology → verify shares bucket with Communication Services
- [ ] Select Basic Materials → verify shares bucket with Industrials
- [ ] Select Financial → verify shares bucket with Consumer Cyclical
- [ ] Select Consumer Defensive → verify shares bucket with Utilities
- [ ] Verify Healthcare, Energy, Real Estate, ETFs are standalone

### Calendar View
- [ ] Navigate to Calendar screen
- [ ] Verify 8 bucket rows display:
  - [ ] Materials/Industrials
  - [ ] Tech/Comm
  - [ ] Financial/Cyclical
  - [ ] Defensive/Utilities
  - [ ] Energy
  - [ ] Healthcare
  - [ ] Real Estate
  - [ ] ETFs

### Manual Override
- [ ] Select Technology (auto-selects Tech/Comm)
- [ ] Manually change bucket to Healthcare
- [ ] Verify bucket dropdown allows override
- [ ] Verify results display shows manual selection

---

## Scenarios

### Scenario 1: Tech Concentration Prevented

**Setup:**
- Account: $100,000
- Bucket cap: 1.5% = $1,500
- Current positions:
  - NVDA (Technology, $750) → Tech/Comm bucket
  - GOOGL (Communication Services, $750) → Tech/Comm bucket
  - **Tech/Comm heat:** $1,500 / $1,500 (100%)

**User Action:** Try to add MSFT (Technology, $750)

**System Response:**
```
Ticker: MSFT
Risk: $750.00
Sector: Technology
Bucket: Tech/Comm

Gate Results:
✗ Gate 4 - Heat Caps FAILED
  Portfolio heat: $2,250 / $4,000 (56%) ✓
  Bucket heat (Tech/Comm): $2,250 / $1,500 (150%) ✗ EXCEEDS by $750

Overall: NO-GO

Suggestion: Reduce position size to $0 (bucket at cap)
          OR close existing Tech/Comm position
          OR choose different sector (Healthcare, Energy, Financial/Cyclical)
```

### Scenario 2: Diversification Encouraged

**User adds positions in sequence:**

1. **NVDA (Technology)** - $750 → Tech/Comm ($750 / $1,500 - 50%)
2. **UNH (Healthcare)** - $750 → Healthcare ($750 / $1,500 - 50%)
3. **XLE (Energy)** - $750 → Energy ($750 / $1,500 - 50%)
4. **JPM (Financial)** - $750 → Financial/Cyclical ($750 / $1,500 - 50%)
5. **FCX (Basic Materials)** - $750 → Materials/Industrials ($750 / $1,500 - 50%)

**Result:**
- Portfolio heat: $3,750 / $4,000 (94%)
- 5 different buckets used
- Maximum correlation protection
- True market diversification

### Scenario 3: ETF Tracking

**Positions:**
- XLK (Technology ETF) - $750 → ETFs bucket
- MSFT (Technology stock) - $750 → Tech/Comm bucket

**Heat Summary:**
- ETFs bucket: $750 / $1,500 (50%)
- Tech/Comm bucket: $750 / $1,500 (50%)
- Total portfolio: $1,500 / $4,000 (38%)

**Note:** ETFs and individual stocks tracked separately, prevents double-counting

---

## Advanced: Adjusting Bucket Groupings

If you need to modify bucket groupings in the future:

**1. Edit ui/trade_entry.go (Lines 69-78):**
```go
bucketOptions := []string{
    // Add/remove/rename buckets here
}
```

**2. Update sectorToBucket mapping (Lines 82-95):**
```go
sectorToBucket := map[string]string{
    "Sector Name": "Bucket Name",
    // Update mappings here
}
```

**3. Edit ui/calendar.go (Lines 27-36):**
```go
buckets := []string{
    // Must match trade_entry.go bucket list
}
```

**4. Rebuild:**
```bash
cd ui
go build -o tf-gui-v4.exe .
```

---

## Conclusion

The v4 bucket system implements your exact specifications:

✅ **12 FINVIZ sectors** (in order, top to bottom) + ETFs
✅ **8 bucket groupings** based on correlation:
  - 4 grouped buckets (2 sectors each)
  - 4 standalone buckets

✅ **Correlation logic:**
  - Materials + Industrials (commodity/cycle)
  - Tech + Comm (growth/innovation)
  - Financial + Cyclical (economic sensitivity)
  - Defensive + Utilities (safety/dividends)
  - Energy, Healthcare, Real Estate, ETFs (standalone)

✅ **Auto-mapping** with manual override capability
✅ **Calendar integration** for visual diversification tracking
✅ **Heat cap enforcement** prevents correlated risk

**This system forces true market diversification while maintaining the flexibility to make individual sector decisions.**

---

**Build:** ui/tf-gui-v4.exe (49MB)
**Ready for testing!**
