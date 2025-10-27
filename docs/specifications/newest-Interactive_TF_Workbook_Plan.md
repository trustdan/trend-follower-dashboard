# Interactive Trend-Following Trade Entry Workbook — Build Plan (Excel VBA + Python)

> **Goal**: One **Trade Entry** worksheet that lets you do everything—run the screener checklist, size the position, enforce impulse limits, and save a decision—using dropdowns, checkboxes, buttons, color cues, and automatic **GO / NO-GO** logic. Behind the scenes, structured tables keep a clean log (Decisions/Positions) and compute portfolio/bucket heat.
>
> **Tech Stack**: VBA for UI/automation + Python in Excel (via `xl()`) for web scraping and complex calculations.

---

## Table of Contents

1. [Concept in One Picture](#concept-in-one-picture)
2. [Sheets & Tables](#sheets--tables)
3. [Named Settings & Defaults](#named-settings--defaults)
4. [UI Spec — *Trade Entry* sheet](#ui-spec--trade-entry-sheet)
5. [Cell Reference Map](#cell-reference-map)
6. [Interaction Flow](#interaction-flow)
7. [Checklist Rules & GO/NO-GO Engine](#checklist-rules--gono-go-engine)
8. [Position Sizing Logic](#position-sizing-logic)
9. [Heat Caps, Cooldown & Impulse Brake](#heat-caps-cooldown--impulse-brake)
10. [Conditional Formatting & Data Validation](#conditional-formatting--data-validation)
11. [Buttons & What They Call](#buttons--what-they-call)
12. [Python Integration Strategy](#python-integration-strategy)
13. [Module Layout (VBA pseudo-code)](#module-layout-vba-pseudo-code)
14. [Worksheet Event Hooks (pseudo-code)](#worksheet-event-hooks-pseudo-code)
15. [Data Dictionary](#data-dictionary)
16. [Mini LLM Prompts](#mini-llm-prompts)
17. [Gherkin / Cucumber Acceptance Tests](#gherkin--cucumber-acceptance-tests)
18. [Sample Walkthrough (MSFT, $10k portfolio)](#sample-walkthrough-msft-10k-portfolio)
19. [Implementation Roadmap](#implementation-roadmap)
20. [Definition of Done](#definition-of-done)

---

## Concept in One Picture

**Trade Entry** (one sheet) → user picks/enters: *Preset → Ticker → Sector → Bucket → Entry → N → K → Method → Delta/DTE or MaxLoss* and ticks required checks.  
On **Evaluate**, a rules engine returns **GREEN (GO)**, **YELLOW (Caution)**, or **RED (NO-GO)** with reasons.  
On **Recalc Sizing**, it computes **R $, shares/contracts, stop, add-levels**.  
On **Save Decision**, it hard-gates: **GREEN only**, **among today’s Candidates**, **heat within caps**, **bucket not cooling down**, **2‑minute impulse delay** satisfied → logs to **Decisions** and updates **Positions**.  

---

## Sheets & Tables

Keep the familiar structure but make **Trade Entry** the main UI. Other sheets are storage or dashboards.

- **Trade Entry** (interactive UI) — the only sheet you touch during the day.
- **Presets** — `tblPresets(Name, QueryString)`; saved screener URLs/params.
- **Buckets** — `tblBuckets(Sector, Bucket, BucketHeatCapPct, StopoutsToCooldown, StopoutsWindowBars, CooldownBars, CooldownActive, CooldownEndsOn)`.
- **Candidates** — `tblCandidates(Date, Ticker, Preset, Sector, Bucket)` *(minimal; more fields are set in Trade Entry)*.
- **Decisions** — `tblDecisions(DateTime, Ticker, Preset, Bucket, N_ATR, K, Entry, RiskPct_r, R_dollars, Size_Shares, Size_Contracts, Method, Delta, DTE, InitialStop, Banner, HeatAtEntry, BucketHeatPost, PortHeatPost, Outcome, Notes)`.
- **Positions** — `tblPositions(Ticker, Bucket, OpenDate, UnitsOpen, RperUnit, TotalOpenR, Method, Status, LastAddPrice, NextAddPrice)`.
- **Summary** — global settings (named ranges), small KPIs.
- **Calendar/Review** — optional dashboards (later).

---

## Named Settings & Defaults

Define these on **Summary** and name the cells:

| Name | Meaning | Suggested Default |
|---|---|---|
| `Equity_E` | Account equity used for sizing | 10,000 *(or your real number)* |
| `RiskPct_r` | % of equity risked per unit | 0.0075 *(0.75%)* |
| `StopMultiple_K` | ATR multiple for initial stop | 2 |
| `HeatCap_H_pct` | Portfolio heat cap (% of equity) | 0.04 |
| `BucketHeatCap_pct` | Per-bucket heat cap (% of equity) | 0.015 |
| `AddStepN` | Add-on step in ATR N units | 0.5 |
| `EarningsBufferDays` | Buffer days around earnings | 3 |

---

## UI Spec — *Trade Entry* sheet

**Layout (single screen):**

- **Row 1-2:** Title + *GO/NO-GO banner* (wide merged cell).  
- **Left column: Inputs** (with dropdowns and checkboxes):
  - *Preset* (dropdown from `tblPresets[Name]`)
  - *Ticker* (dropdown from `tblCandidates[Ticker]` filtered to today; allow typing)
  - *Sector* (dropdown: common sectors)
  - *Bucket* (dropdown from `tblBuckets[Bucket]`, auto-suggest from sector)
  - *Entry Price* (number)
  - *ATR N* (number)  
  - *K (Stop Multiple)* (spin button 1–4; default 2)
  - *Method* (option buttons): **Stock**, **Opt-DeltaATR**, **Opt-MaxLoss**
  - *Delta* (if Opt-DeltaATR)
  - *DTE* (if options)
  - *MaxLossPerContract* (if Opt-MaxLoss)
  - Checklist (checkboxes): **FromPreset**, **TrendPass**, **LiquidityPass**, **TVConfirm**, **EarningsOK**, **JournalOK**
- **Right column: Outputs** (formula cells filled by macros):
  - **R ($)**, **StopDist (K×N)**, **Initial Stop**, **Shares**, **Contracts**, **Add1**, **Add2**, **Add3**
  - **Heat Preview**: *Portfolio heat post*, *Bucket heat post* vs caps (horizontal bars).
- **Buttons:** `Evaluate`, `Recalc Sizing`, `Save Decision`, `Start 2‑min Timer`, `Import Candidates`, `Open Preset`.

**Visual cues:**
- Banner background: **GREEN/YELLOW/RED** with reason tooltip.
- Heat preview bars: **green ≤ 70%**, **amber 70–100%**, **red > 100%**.
- Disabled controls if Method doesn't need them (e.g., hide Delta when Stock).

---

## Cell Reference Map

**Trade Entry Sheet Layout** (concrete cell addresses):

| Cell/Range | Content | Type | Notes |
|---|---|---|---|
| **A1:F1** | "TRADE ENTRY SYSTEM" | Label | Merged, bold, 16pt |
| **A2:F2** | Banner Text | Dynamic | GREEN / YELLOW / RED (set by VBA) |
| **A3:F3** | Reason String | Dynamic | List of failed checks/blocks |
| | | | |
| **A5** | "Preset:" | Label | |
| **B5** | Preset dropdown | Input | DV from tblPresets[Name] |
| **A6** | "Ticker:" | Label | |
| **B6** | Ticker dropdown | Input | DV from tblCandidates[Ticker] where Date=TODAY() |
| **A7** | "Sector:" | Label | |
| **B7** | Sector dropdown | Input | DV list: Technology, Healthcare, Financials, etc. |
| **A8** | "Bucket:" | Label | |
| **B8** | Bucket dropdown | Input | DV from tblBuckets[Bucket], auto-suggest from B7 |
| **A9** | "Entry Price:" | Label | |
| **B9** | Entry price | Input | Number, 2 decimal places |
| **A10** | "ATR N:" | Label | |
| **B10** | N value | Input | Number, 3 decimal places |
| **A11** | "K (Stop Multiple):" | Label | |
| **B11** | K value | Input | Spin button 1-4, default 2 |
| | | | |
| **A13** | "Method:" | Label | |
| **B13:B15** | Option buttons | Input | Stock / Opt-DeltaATR / Opt-MaxLoss (linked to C13) |
| **C13** | Method choice | Hidden | 1=Stock, 2=Opt-DeltaATR, 3=Opt-MaxLoss |
| **A16** | "Delta:" | Label | Show if C13=2 |
| **B16** | Delta | Input | Number, 0.01-0.99 |
| **A17** | "DTE:" | Label | Show if C13>1 |
| **B17** | DTE | Input | Integer, 21-90 |
| **A18** | "Max Loss/Contract:" | Label | Show if C13=3 |
| **B18** | MaxLoss | Input | Number, dollar amount |
| | | | |
| **A20:A25** | Checklist labels | Label | FromPreset, TrendPass, LiquidityPass, TVConfirm, EarningsOK, JournalOK |
| **B20:B25** | Checkboxes | Input | Linked to C20:C25 (TRUE/FALSE) |
| **C20:C25** | Checkbox values | Hidden | Used by VBA |
| | | | |
| **E5** | "R ($):" | Label | |
| **F5** | R_dollars | Output | =Equity_E * RiskPct_r |
| **E6** | "Stop Distance:" | Label | |
| **F6** | StopDist | Output | =B11 * B10 |
| **E7** | "Initial Stop:" | Label | |
| **F7** | InitialStop | Output | =B9 - F6 |
| **E8** | "Shares:" | Label | |
| **F8** | Shares | Output | Computed by RecalcSizing (method-dependent) |
| **E9** | "Contracts:" | Label | |
| **F9** | Contracts | Output | Computed by RecalcSizing (method-dependent) |
| **E10** | "Add 1:" | Label | |
| **F10** | Add1 | Output | =B9 + (AddStepN * B10) |
| **E11** | "Add 2:" | Label | |
| **F11** | Add2 | Output | =B9 + (2 * AddStepN * B10) |
| **E12** | "Add 3:" | Label | |
| **F12** | Add3 | Output | =B9 + (3 * AddStepN * B10) |
| | | | |
| **E14** | "Portfolio Heat:" | Label | |
| **F14** | Portfolio heat bar | Output | Data bar CF, computed by Python/VBA |
| **E15** | "Bucket Heat:" | Label | |
| **F15** | Bucket heat bar | Output | Data bar CF, computed by Python/VBA |
| | | | |
| **A28** | "Open Preset" button | Button | Calls OpenPreset() |
| **B28** | "Import Candidates" button | Button | Calls ImportCandidatesPrompt() |
| **A29** | "Evaluate" button | Button | Calls EvaluateChecklist() |
| **B29** | "Recalc Sizing" button | Button | Calls RecalcSizing() |
| **A30** | "Save Decision" button | Button | Calls SaveDecision() |
| **B30** | "Start 2-min Timer" button | Button | Stores Now in Control!A1 |

**Hidden Control Sheet:**
- **Control!A1**: Impulse timer timestamp (set by Evaluate, checked by Save)

---

## Interaction Flow

1. **Import candidates** from a FINVIZ preset → rows added to `tblCandidates` (today’s date).  
2. In **Trade Entry**, select *Preset* then *Ticker* (auto-populate Sector/ Bucket if known).  
3. Fill price/ATR N/K and choose sizing method + parameters.  
4. Tick required checklist boxes.  
5. Click **Evaluate** → shows **Banner** with **GO/NO-GO** and reasons; starts a **2‑minute impulse timer** if GO.  
6. Click **Recalc Sizing** → fills outputs.  
7. Click **Save Decision** → hard-gates; if pass, logs row to `tblDecisions`, updates `tblPositions`, and resets the timer.

---

## Checklist Rules & GO/NO-GO Engine

**Required checks for GO:** `FromPreset`, `TrendPass`, `LiquidityPass`, `TVConfirm`, `EarningsOK`, `JournalOK`.  
**Banner logic:**  
- **GREEN (GO)**: all required checks true.  
- **YELLOW (Caution)**: exactly one missing (show which).  
- **RED (NO-GO)**: two or more missing, or any hard fail (heat caps, cooldown, not a candidate, timer not elapsed).

**Additional blockers (automatic NO-GO):**
- Ticker not in today’s **Candidates**.  
- **Bucket cooldown** active.  
- **Portfolio heat post** > `HeatCap_H_pct × Equity_E`.  
- **Bucket heat post** > `BucketHeatCap_pct × Equity_E`.  
- **2‑minute impulse** not elapsed since **Evaluate**.

**Reason string (for the banner tooltip):** list the failing checks / blocks.

---

## Position Sizing Logic

Let:
- `E` = `Equity_E`  
- `r` = risk % per unit (`RiskPct_r`)  
- `R` = `E × r` (risk dollars)  
- `N` = ATR value (your “N”)  
- `K` = stop multiple  
- `entry` = entry price

**Common:** `StopDist = K × N`, `InitialStop = entry − StopDist`

**Stock:**  
`Shares = floor(R / StopDist)`

**Options – Delta-ATR:**  
`Contracts = floor( R / (K × N × Delta × 100) )`

**Options – MaxLoss:**  
`Contracts = floor( R / (MaxLossPerContract × 100) )` *(use debit of vertical or width minus credit)*

**Add levels:**  
`Add1 = entry + AddStepN × N`, `Add2 = entry + 2×AddStepN × N`, `Add3 = entry + 3×AddStepN × N`

---

## Heat Caps, Cooldown & Impulse Brake

- **Portfolio Heat** = sum of **RperUnit** of open positions in `tblPositions`; cap = `HeatCap_H_pct × E`.
- **Bucket Heat** = sum of open **RperUnit** where `Bucket = selected`; cap = bucket’s `BucketHeatCapPct × E` if set, else `BucketHeatCap_pct × E`.
- **Cooldown**: if a bucket has ≥ `StopoutsToCooldown` **StopOut** outcomes within `StopoutsWindowBars` recent days → set `CooldownActive=TRUE` and `CooldownEndsOn=Today+CooldownBars`.
- **Impulse brake**: on successful **Evaluate** with GREEN, record timestamp `T0`. **Save Decision** only if `Now ≥ T0 + 2 minutes`.

---

## Conditional Formatting & Data Validation

**Data Validation (DV):**
- Dropdowns for *Preset, Ticker, Sector, Bucket, Method*.  
- Ticker DV uses a dynamic named range filtering `tblCandidates` to **today** (can be simplified to “all tickers” if filtering in code).

**Conditional Formatting (CF):**
- Banner cell fill by text: `GREEN`, `YELLOW`, `RED`.
- Heat preview bars: color by % of cap.
- Disable/grey-out blocks when Method = “Stock” (format cell if not applicable).

---

## Buttons & What They Call

| Button | Action (high-level) |
|---|---|
| **Open Preset** | Opens the FINVIZ URL for the selected preset in browser. |
| **Import Candidates** | Prompts user to paste tickers → normalizes → adds to `tblCandidates` with today’s date + selected preset. |
| **Evaluate** | Runs **GO/NO-GO** engine. Writes Banner + reason. If GREEN, starts impulse timer. |
| **Recalc Sizing** | Computes R, stop, shares/contracts, add levels. Fills output cells. |
| **Save Decision** | Validates again (GREEN, caps, cooldown, timer). Appends to **Decisions**. Opens/updates **Positions**. |
| **Start 2‑min Timer** | Manual starter; useful if Evaluate already run. |
| **Refresh Review** *(optional)* | Recomputes KPIs like "% taken non‑GREEN (last 30)". |

---

## Python Integration Strategy

**Why Python in Excel?** Excel now supports Python via the `=PY()` function (Microsoft 365 Insider). Use Python for:
1. **Web scraping FINVIZ** (replaces manual copy/paste)
2. **Complex heat calculations** (vectorized operations)
3. **Data validation** (batch ticker normalization)

**When to use VBA vs Python:**
- **VBA**: UI events, button clicks, worksheet manipulation, banner coloring
- **Python**: HTTP requests, pandas DataFrame operations, numerical computations

---

### Python Module: `finviz_scraper.py`

**Purpose:** Fetch tickers from FINVIZ preset URLs and return as list.

```python
# finviz_scraper.py (stored in Excel Python environment)
import requests
from bs4 import BeautifulSoup
import pandas as pd

def fetch_finviz_tickers(query_string: str) -> list:
    """
    Scrapes FINVIZ screener page and returns list of tickers.

    Args:
        query_string: The FINVIZ query string (from tblPresets[QueryString])

    Returns:
        List of ticker symbols (normalized, uppercase, deduped)
    """
    base_url = "https://finviz.com/screener.ashx"
    url = f"{base_url}?{query_string}"

    headers = {'User-Agent': 'Mozilla/5.0'}
    response = requests.get(url, headers=headers)

    if response.status_code != 200:
        return []

    soup = BeautifulSoup(response.content, 'html.parser')
    table = soup.find('table', {'class': 'table-light'})

    if not table:
        return []

    tickers = []
    for row in table.find_all('tr')[1:]:  # Skip header
        cells = row.find_all('td')
        if cells:
            ticker = cells[1].text.strip()  # Ticker is in 2nd column
            tickers.append(ticker.upper())

    return list(set(tickers))  # Dedupe

def normalize_tickers(raw_list: list) -> list:
    """
    Cleans and normalizes ticker symbols.

    Removes: special chars, duplicates, blanks
    Returns: List of clean uppercase tickers
    """
    normalized = []
    for ticker in raw_list:
        clean = ticker.strip().upper().replace('.', '-')
        if clean and len(clean) <= 5:  # Basic validation
            normalized.append(clean)
    return list(set(normalized))
```

**Excel Integration:**
```python
# In Excel cell (example usage in hidden sheet):
=PY("finviz_scraper.fetch_finviz_tickers", Presets!B2)
# Returns array of tickers which VBA reads and writes to tblCandidates
```

---

### Python Module: `heat_calculator.py`

**Purpose:** Fast vectorized heat calculations using pandas.

```python
# heat_calculator.py
import pandas as pd

def portfolio_heat_after(positions_df: pd.DataFrame, add_r: float) -> float:
    """
    Calculates total portfolio heat + proposed trade.

    Args:
        positions_df: DataFrame with columns [Ticker, Status, TotalOpenR]
        add_r: Additional R dollars from proposed trade

    Returns:
        Total heat in dollars
    """
    open_positions = positions_df[positions_df['Status'] != 'Closed']
    current_heat = open_positions['TotalOpenR'].sum()
    return current_heat + add_r

def bucket_heat_after(positions_df: pd.DataFrame, bucket: str, add_r: float) -> float:
    """
    Calculates bucket-specific heat + proposed trade.

    Args:
        positions_df: DataFrame with columns [Bucket, Status, TotalOpenR]
        bucket: Bucket name to filter by
        add_r: Additional R dollars from proposed trade

    Returns:
        Bucket heat in dollars
    """
    bucket_positions = positions_df[
        (positions_df['Bucket'] == bucket) &
        (positions_df['Status'] != 'Closed')
    ]
    current_heat = bucket_positions['TotalOpenR'].sum()
    return current_heat + add_r

def check_heat_caps(positions_df: pd.DataFrame,
                   add_r: float,
                   bucket: str,
                   equity: float,
                   port_cap_pct: float,
                   bucket_cap_pct: float) -> dict:
    """
    Validates proposed trade against heat caps.

    Returns:
        {
            'portfolio_ok': bool,
            'bucket_ok': bool,
            'portfolio_heat': float,
            'bucket_heat': float,
            'portfolio_cap': float,
            'bucket_cap': float
        }
    """
    port_heat = portfolio_heat_after(positions_df, add_r)
    buck_heat = bucket_heat_after(positions_df, bucket, add_r)

    port_cap = equity * port_cap_pct
    buck_cap = equity * bucket_cap_pct

    return {
        'portfolio_ok': port_heat <= port_cap,
        'bucket_ok': buck_heat <= buck_cap,
        'portfolio_heat': port_heat,
        'bucket_heat': buck_heat,
        'portfolio_cap': port_cap,
        'bucket_cap': buck_cap
    }
```

**Excel Integration:**
```python
# In Excel cell F14 (Portfolio Heat):
=PY("heat_calculator.portfolio_heat_after",
    xl("Positions[#All]"),
    TradeEntry!F5)

# VBA reads this result to validate caps
```

---

### VBA-Python Bridge Pattern

**VBA calls Python, reads result:**

```vba
' Module: TF_Python_Bridge

Function CallPythonHeatCheck(addR As Double, bucket As String) As Variant
    ' Calls Python heat_calculator and returns dict as array
    Dim pyFormula As String
    Dim resultCell As Range

    ' Write Python formula to hidden calc sheet
    Set resultCell = Worksheets("Control").Range("Z1")

    pyFormula = "=PY(""heat_calculator.check_heat_caps"", " & _
                "xl(""Positions[#All]""), " & _
                addR & ", " & _
                """" & bucket & """, " & _
                "Summary!B2, " & _    ' Equity_E
                "Summary!B4, " & _    ' HeatCap_H_pct
                "Summary!B5)"         ' BucketHeatCap_pct

    resultCell.Formula = pyFormula
    DoEvents

    ' Parse returned dict (Excel converts to array)
    CallPythonHeatCheck = resultCell.Value
End Function

Sub TestPythonIntegration()
    ' Example: Call from VBA, get heat check result
    Dim result As Variant
    result = CallPythonHeatCheck(75, "Tech/Comm")

    If result("portfolio_ok") = True And result("bucket_ok") = True Then
        Debug.Print "Heat caps OK"
    Else
        Debug.Print "Heat caps EXCEEDED"
    End If
End Sub
```

---

### Hybrid Architecture Decision Matrix

| Task | Tool | Reason |
|---|---|---|
| Button clicks, form events | VBA | Native Excel automation |
| Banner coloring, CF rules | VBA | Direct worksheet manipulation |
| FINVIZ web scraping | Python | requests/BeautifulSoup superior |
| Heat calculations (batch) | Python | pandas vectorization faster |
| Ticker normalization | Python | Regex/string ops cleaner |
| SaveDecision validation | VBA | Orchestrates Python + writes tables |
| Conditional formatting | VBA | Native Excel API |
| Data table reads/writes | VBA | ListObject operations built-in |

**Rule of thumb:**
- **UI/Events** → VBA
- **Data/Computation** → Python (when complex)
- **Orchestration** → VBA calls Python via `=PY()` and reads results

---

## Module Layout (VBA pseudo-code)

> **Note:** Pseudo-code, not full VBA. Use this as a blueprint when you ask an LLM to generate actual code.

**Module: `TF_Utils`**
```
Function SheetExists(name) As Boolean
Function GetOrCreateSheet(name) As Worksheet
Function GetOrCreateTable(ws, tableName, headers) As ListObject
Function EnsureName(name, refersTo, default)
Function NzD(v, default) As Double
Function NzS(v, default) As String
Function NormalizeTicker(raw) As String
```

**Module: `TF_Data`**
```
Sub EnsureStructure()
    ' Create all sheets if missing
    Call GetOrCreateSheet("TradeEntry")
    Call GetOrCreateSheet("Presets")
    Call GetOrCreateSheet("Buckets")
    Call GetOrCreateSheet("Candidates")
    Call GetOrCreateSheet("Decisions")
    Call GetOrCreateSheet("Positions")
    Call GetOrCreateSheet("Summary")
    Call GetOrCreateSheet("Control")  ' Hidden helper sheet

    ' Create tables with headers
    Call GetOrCreateTable(Sheets("Presets"), "tblPresets", Array("Name", "QueryString"))
    Call GetOrCreateTable(Sheets("Buckets"), "tblBuckets", _
        Array("Sector", "Bucket", "BucketHeatCapPct", "StopoutsToCooldown", _
              "StopoutsWindowBars", "CooldownBars", "CooldownActive", "CooldownEndsOn"))
    Call GetOrCreateTable(Sheets("Candidates"), "tblCandidates", _
        Array("Date", "Ticker", "Preset", "Sector", "Bucket"))
    Call GetOrCreateTable(Sheets("Decisions"), "tblDecisions", _
        Array("DateTime", "Ticker", "Preset", "Bucket", "N_ATR", "K", "Entry", _
              "RiskPct_r", "R_dollars", "Size_Shares", "Size_Contracts", "Method", _
              "Delta", "DTE", "InitialStop", "Banner", "HeatAtEntry", _
              "BucketHeatPost", "PortHeatPost", "Outcome", "Notes"))
    Call GetOrCreateTable(Sheets("Positions"), "tblPositions", _
        Array("Ticker", "Bucket", "OpenDate", "UnitsOpen", "RperUnit", _
              "TotalOpenR", "Method", "Status", "LastAddPrice", "NextAddPrice"))

    ' Define named ranges on Summary sheet
    Call EnsureName("Equity_E", "Summary!B2", 10000)
    Call EnsureName("RiskPct_r", "Summary!B3", 0.0075)
    Call EnsureName("StopMultiple_K", "Summary!B4", 2)
    Call EnsureName("HeatCap_H_pct", "Summary!B5", 0.04)
    Call EnsureName("BucketHeatCap_pct", "Summary!B6", 0.015)
    Call EnsureName("AddStepN", "Summary!B7", 0.5)
    Call EnsureName("EarningsBufferDays", "Summary!B8", 3)

    ' Seed default presets (5 FINVIZ queries)
    Call SeedPresets
    ' Seed default buckets (6 correlation groups)
    Call SeedBuckets
End Sub

Function TodayCandidates() As Range
    ' Returns dynamic range of tickers where Date = Today()
    ' Used for Data Validation dropdown in B6
    Dim tbl As ListObject
    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    ' Filter logic handled by VBA or Python; return range
    Set TodayCandidates = tbl.ListColumns("Ticker").DataBodyRange
End Function

Function PortfolioHeatAfter(addR As Double) As Double
    ' Option 1: Pure VBA (loop through tblPositions)
    Dim tbl As ListObject, row As ListRow
    Dim total As Double
    Set tbl = Sheets("Positions").ListObjects("tblPositions")
    total = 0
    For Each row In tbl.ListRows
        If row.Range.Columns(8).Value <> "Closed" Then  ' Status column
            total = total + row.Range.Columns(6).Value  ' TotalOpenR column
        End If
    Next row
    PortfolioHeatAfter = total + addR

    ' Option 2: Call Python (if available)
    ' PortfolioHeatAfter = CallPythonHeatCheck(addR, "")("portfolio_heat")
End Function

Function BucketHeatAfter(bucket As String, addR As Double) As Double
    ' Sum TotalOpenR where Bucket = bucket AND Status <> "Closed"
    Dim tbl As ListObject, row As ListRow
    Dim total As Double
    Set tbl = Sheets("Positions").ListObjects("tblPositions")
    total = 0
    For Each row In tbl.ListRows
        If row.Range.Columns(2).Value = bucket And _
           row.Range.Columns(8).Value <> "Closed" Then
            total = total + row.Range.Columns(6).Value
        End If
    Next row
    BucketHeatAfter = total + addR
End Function

Function IsBucketInCooldown(bucket As String) As Boolean
    ' Check tblBuckets: CooldownActive = TRUE AND CooldownEndsOn >= Today
    Dim tbl As ListObject, row As ListRow
    Set tbl = Sheets("Buckets").ListObjects("tblBuckets")
    For Each row In tbl.ListRows
        If row.Range.Columns(2).Value = bucket Then  ' Bucket column
            If row.Range.Columns(7).Value = True Then  ' CooldownActive
                If row.Range.Columns(8).Value >= Date Then  ' CooldownEndsOn
                    IsBucketInCooldown = True
                    Exit Function
                End If
            End If
        End If
    Next row
    IsBucketInCooldown = False
End Function

Sub UpdateCooldowns()
    ' Scan last StopoutsWindowBars decisions for each bucket
    ' If >= StopoutsToCooldown StopOuts found, set CooldownActive = TRUE
    ' Set CooldownEndsOn = Today + CooldownBars
    ' Called weekly or after each stop-out

    Dim bucketTbl As ListObject, decisTbl As ListObject
    Dim bucket As String, stopoutCount As Integer
    Dim windowStart As Date

    Set bucketTbl = Sheets("Buckets").ListObjects("tblBuckets")
    Set decisTbl = Sheets("Decisions").ListObjects("tblDecisions")

    For Each bRow In bucketTbl.ListRows
        bucket = bRow.Range.Columns(2).Value
        stopoutThreshold = bRow.Range.Columns(4).Value  ' StopoutsToCooldown
        windowBars = bRow.Range.Columns(5).Value        ' StopoutsWindowBars
        cooldownBars = bRow.Range.Columns(6).Value      ' CooldownBars

        windowStart = Date - windowBars

        ' Count StopOuts in window
        stopoutCount = 0
        For Each dRow In decisTbl.ListRows
            If dRow.Range.Columns(4).Value = bucket And _
               dRow.Range.Columns(1).Value >= windowStart And _
               dRow.Range.Columns(20).Value = "StopOut" Then  ' Outcome column
                stopoutCount = stopoutCount + 1
            End If
        Next dRow

        ' Update cooldown flags
        If stopoutCount >= stopoutThreshold Then
            bRow.Range.Columns(7).Value = True  ' CooldownActive
            bRow.Range.Columns(8).Value = Date + cooldownBars  ' CooldownEndsOn
        Else
            ' Clear cooldown if past end date
            If bRow.Range.Columns(8).Value < Date Then
                bRow.Range.Columns(7).Value = False
            End If
        End If
    Next bRow
End Sub
```

**Module: `TF_UI`**
```
Sub BindControls()
    ' Set up Data Validation dropdowns
    With Sheets("TradeEntry")
        .Range("B5").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblPresets[Name]"
        .Range("B6").Validation.Add Type:=xlValidateList, _
            Formula1:="=TodayCandidates()"  ' Dynamic function
        .Range("B7").Validation.Add Type:=xlValidateList, _
            Formula1:="Technology,Healthcare,Financials,Consumer,Industrials,Energy"
        .Range("B8").Validation.Add Type:=xlValidateList, _
            Formula1:="=tblBuckets[Bucket]"

        ' Link checkboxes to cells (if using Form Controls)
        ' Or link ActiveX checkboxes via .LinkedCell = "C20", etc.
    End With

    ' Set default Method to Stock (option button 1)
    Sheets("TradeEntry").Range("C13").Value = 1

    ' Show/hide method-specific fields
    Call ToggleMethodFields
End Sub

Sub ToggleMethodFields()
    ' Show/hide Delta, DTE, MaxLoss based on Method choice in C13
    Dim methodChoice As Integer
    methodChoice = Sheets("TradeEntry").Range("C13").Value

    With Sheets("TradeEntry")
        If methodChoice = 1 Then  ' Stock
            .Rows("16:18").Hidden = True
        ElseIf methodChoice = 2 Then  ' Opt-DeltaATR
            .Rows("16:17").Hidden = False
            .Rows("18:18").Hidden = True
        ElseIf methodChoice = 3 Then  ' Opt-MaxLoss
            .Rows("17:18").Hidden = False
            .Rows("16:16").Hidden = True
        End If
    End With
End Sub

Sub EvaluateChecklist()
    ' Read checklist values from C20:C25
    Dim checks(1 To 6) As Boolean
    Dim missingCount As Integer, reasons As String
    Dim banner As String, bannerColor As Long

    With Sheets("TradeEntry")
        checks(1) = .Range("C20").Value  ' FromPreset
        checks(2) = .Range("C21").Value  ' TrendPass
        checks(3) = .Range("C22").Value  ' LiquidityPass
        checks(4) = .Range("C23").Value  ' TVConfirm
        checks(5) = .Range("C24").Value  ' EarningsOK
        checks(6) = .Range("C25").Value  ' JournalOK
    End With

    ' Count missing
    missingCount = 0
    reasons = ""
    If Not checks(1) Then: missingCount = missingCount + 1: reasons = reasons & "FromPreset, "
    If Not checks(2) Then: missingCount = missingCount + 1: reasons = reasons & "TrendPass, "
    If Not checks(3) Then: missingCount = missingCount + 1: reasons = reasons & "LiquidityPass, "
    If Not checks(4) Then: missingCount = missingCount + 1: reasons = reasons & "TVConfirm, "
    If Not checks(5) Then: missingCount = missingCount + 1: reasons = reasons & "EarningsOK, "
    If Not checks(6) Then: missingCount = missingCount + 1: reasons = reasons & "JournalOK, "

    ' Determine banner
    If missingCount = 0 Then
        banner = "GREEN - GO"
        bannerColor = RGB(0, 200, 0)
        ' Start impulse timer
        Sheets("Control").Range("A1").Value = Now
    ElseIf missingCount = 1 Then
        banner = "YELLOW - CAUTION"
        bannerColor = RGB(255, 200, 0)
    Else
        banner = "RED - NO-GO"
        bannerColor = RGB(255, 0, 0)
    End If

    ' Write banner and reasons
    With Sheets("TradeEntry")
        .Range("A2").Value = banner
        .Range("A2").Interior.Color = bannerColor
        .Range("A3").Value = "Missing: " & Left(reasons, Len(reasons) - 2)
    End With
End Sub

Sub RecalcSizing()
    ' Read inputs from TradeEntry sheet
    Dim E As Double, r As Double, R As Double
    Dim entry As Double, N As Double, K As Double
    Dim methodChoice As Integer, delta As Double, maxLoss As Double
    Dim shares As Long, contracts As Long
    Dim stopDist As Double, initialStop As Double

    E = Range("Equity_E").Value
    r = Range("RiskPct_r").Value
    R = E * r

    With Sheets("TradeEntry")
        entry = .Range("B9").Value
        N = .Range("B10").Value
        K = .Range("B11").Value
        methodChoice = .Range("C13").Value
        delta = .Range("B16").Value
        maxLoss = .Range("B18").Value
    End With

    ' Common calculations
    stopDist = K * N
    initialStop = entry - stopDist

    ' Method-specific sizing
    Select Case methodChoice
        Case 1  ' Stock
            shares = WorksheetFunction.Floor_Precise(R / stopDist)
            contracts = 0
        Case 2  ' Opt-DeltaATR
            contracts = WorksheetFunction.Floor_Precise(R / (K * N * delta * 100))
            shares = 0
        Case 3  ' Opt-MaxLoss
            contracts = WorksheetFunction.Floor_Precise(R / (maxLoss * 100))
            shares = 0
    End Select

    ' Write outputs
    With Sheets("TradeEntry")
        .Range("F5").Value = R
        .Range("F6").Value = stopDist
        .Range("F7").Value = initialStop
        .Range("F8").Value = shares
        .Range("F9").Value = contracts
        ' Add levels computed via formulas (already set in Cell Reference Map)
    End With
End Sub

Sub SaveDecision()
    ' Hard-gate validation sequence
    Dim banner As String, ticker As String, bucket As String
    Dim addR As Double, portHeat As Double, buckHeat As Double
    Dim timerStart As Date, elapsed As Double

    banner = Sheets("TradeEntry").Range("A2").Value
    ticker = Sheets("TradeEntry").Range("B6").Value
    bucket = Sheets("TradeEntry").Range("B8").Value
    addR = Sheets("TradeEntry").Range("F5").Value

    ' 1. Banner must be GREEN
    If InStr(banner, "GREEN") = 0 Then
        MsgBox "Cannot save: Banner is not GREEN", vbCritical
        Exit Sub
    End If

    ' 2. Ticker must be in today's Candidates
    If Not IsTickerInCandidates(ticker, Date) Then
        MsgBox "Cannot save: Ticker not in today's Candidates", vbCritical
        Exit Sub
    End If

    ' 3. Check impulse timer (2 minutes)
    timerStart = Sheets("Control").Range("A1").Value
    elapsed = (Now - timerStart) * 24 * 60  ' Convert to minutes
    If elapsed < 2 Then
        MsgBox "Cannot save: 2-minute cool-off not elapsed (" & _
               Format(elapsed, "0.0") & " min)", vbCritical
        Exit Sub
    End If

    ' 4. Check bucket cooldown
    If IsBucketInCooldown(bucket) Then
        MsgBox "Cannot save: Bucket " & bucket & " is in cooldown", vbCritical
        Exit Sub
    End If

    ' 5. Check heat caps (call Python or VBA)
    portHeat = PortfolioHeatAfter(addR)
    buckHeat = BucketHeatAfter(bucket, addR)

    If portHeat > Range("HeatCap_H_pct").Value * Range("Equity_E").Value Then
        MsgBox "Cannot save: Portfolio heat would exceed cap", vbCritical
        Exit Sub
    End If

    If buckHeat > Range("BucketHeatCap_pct").Value * Range("Equity_E").Value Then
        MsgBox "Cannot save: Bucket heat would exceed cap", vbCritical
        Exit Sub
    End If

    ' All gates passed - append to Decisions
    Call AppendDecisionRow
    ' Update/open position in Positions table
    Call UpdatePositions
    ' Reset banner
    Sheets("TradeEntry").Range("A2").Value = ""
    Sheets("TradeEntry").Range("A2").Interior.ColorIndex = xlNone

    MsgBox "Decision saved successfully!", vbInformation
End Sub

Function IsTickerInCandidates(ticker As String, tradeDate As Date) As Boolean
    ' Check if ticker exists in tblCandidates with Date = tradeDate
    Dim tbl As ListObject, row As ListRow
    Set tbl = Sheets("Candidates").ListObjects("tblCandidates")
    For Each row In tbl.ListRows
        If row.Range.Columns(2).Value = ticker And _
           row.Range.Columns(1).Value = tradeDate Then
            IsTickerInCandidates = True
            Exit Function
        End If
    Next row
    IsTickerInCandidates = False
End Function

Sub AppendDecisionRow()
    ' Append current trade to tblDecisions with all fields
    ' (Implementation: read all inputs from TradeEntry, add new row to table)
End Sub

Sub UpdatePositions()
    ' Open new position or add to existing in tblPositions
    ' (Implementation: check if ticker exists, update UnitsOpen/RperUnit/NextAddPrice)
End Sub
```

**Module: `TF_Presets`**
```
Sub OpenPreset()
Sub ImportCandidatesPrompt()
    ' paste tickers → normalize → add to tblCandidates
End Sub
```

---

## Worksheet Event Hooks (pseudo-code)

**Trade Entry Sheet Code (`TradeEntry`):**
```
On Worksheet_Activate:
    Call TF_UI.BindControls

On Worksheet_Change(Target):
    If Target intersects inputs (Preset/Sector/Bucket/Method) Then
        ' auto-suggest bucket by sector
        ' show/hide delta/maxloss fields
        ' clear banner when inputs change
    End If
```

**ThisWorkbook:**
```
On Workbook_Open:
    Call TF_Data.EnsureStructure
```

---

## Data Dictionary

### `tblDecisions` (log)
- **DateTime** *(timestamp)*  
- **Ticker** *(text)*  
- **Preset** *(text)*  
- **Bucket** *(text)*  
- **N_ATR, K, Entry, RiskPct_r** *(numbers)*  
- **R_dollars** *(computed when saved)*  
- **Size_Shares, Size_Contracts** *(integers)*  
- **Method, Delta, DTE, InitialStop, Banner**  
- **HeatAtEntry, BucketHeatPost, PortHeatPost** *(numbers in $R)*  
- **Outcome** *(StopOut | ExitRule | Manual | —)*  
- **Notes** *(text)*

### `tblPositions` (open-book)
- **Ticker, Bucket, OpenDate, UnitsOpen, RperUnit, TotalOpenR, Method, Status, LastAddPrice, NextAddPrice**

---

## Mini LLM Prompts

Copy/paste these one by one to an LLM (Claude, GPT-4) to generate implementation-ready code. Context from this document should be provided.

### VBA Prompts

1) **Structure seeding** *(10 min, foundational)*
```
Write a VBA Sub named EnsureStructure that:
- Creates sheets: TradeEntry, Presets, Buckets, Candidates, Decisions, Positions, Summary, Control
- Creates listobjects with these exact headers:
  * tblPresets: Name, QueryString
  * tblBuckets: Sector, Bucket, BucketHeatCapPct, StopoutsToCooldown, StopoutsWindowBars, CooldownBars, CooldownActive, CooldownEndsOn
  * tblCandidates: Date, Ticker, Preset, Sector, Bucket
  * tblDecisions: DateTime, Ticker, Preset, Bucket, N_ATR, K, Entry, RiskPct_r, R_dollars, Size_Shares, Size_Contracts, Method, Delta, DTE, InitialStop, Banner, HeatAtEntry, BucketHeatPost, PortHeatPost, Outcome, Notes
  * tblPositions: Ticker, Bucket, OpenDate, UnitsOpen, RperUnit, TotalOpenR, Method, Status, LastAddPrice, NextAddPrice
- Defines named ranges on Summary: Equity_E (B2, default 10000), RiskPct_r (B3, 0.0075), StopMultiple_K (B4, 2), HeatCap_H_pct (B5, 0.04), BucketHeatCap_pct (B6, 0.015), AddStepN (B7, 0.5), EarningsBufferDays (B8, 3)
- Seeds 5 FINVIZ presets and 6 default buckets (see CLAUDE.md for preset URLs)
Return only VBA code with helper functions GetOrCreateSheet, GetOrCreateTable, EnsureName.
```

2) **UI binding and controls** *(5 min, interactive)*
```
Write a VBA Sub BindControls for sheet TradeEntry that:
- Sets data validation on B5 (Preset dropdown from tblPresets[Name])
- Sets data validation on B6 (Ticker dropdown from tblCandidates[Ticker] filtered to today)
- Sets data validation on B7 (Sector list: Technology,Healthcare,Financials,Consumer,Industrials,Energy)
- Sets data validation on B8 (Bucket dropdown from tblBuckets[Bucket])
- Links 6 checkboxes (B20:B25) to cells C20:C25 for TRUE/FALSE values
- Sets default Method (C13) to 1 (Stock)
- Calls ToggleMethodFields to show/hide rows based on method choice

Also write ToggleMethodFields that:
- Hides rows 16-18 when C13=1 (Stock)
- Shows rows 16-17, hides 18 when C13=2 (Opt-DeltaATR)
- Shows rows 17-18, hides 16 when C13=3 (Opt-MaxLoss)

Return code only.
```

3) **Checklist evaluation with banner** *(10 min, critical)*
```
Write a VBA Sub EvaluateChecklist that:
- Reads 6 Boolean values from cells C20:C25 (FromPreset, TrendPass, LiquidityPass, TVConfirm, EarningsOK, JournalOK)
- Counts how many are FALSE
- Sets banner logic:
  * 0 missing → banner="GREEN - GO", color=RGB(0,200,0), store Now in Control!A1 (impulse timer)
  * 1 missing → banner="YELLOW - CAUTION", color=RGB(255,200,0)
  * 2+ missing → banner="RED - NO-GO", color=RGB(255,0,0)
- Writes banner text to TradeEntry!A2 with background color
- Writes reason string (list of missing items) to TradeEntry!A3

Return code only. Use clear variable names.
```

4) **Position sizing engine** *(10 min, math-heavy)*
```
Write a VBA Sub RecalcSizing that:
- Reads named ranges: Equity_E, RiskPct_r, AddStepN
- Reads inputs from TradeEntry: B9 (Entry), B10 (N), B11 (K), C13 (Method 1/2/3), B16 (Delta), B18 (MaxLoss)
- Computes:
  * R = Equity_E × RiskPct_r
  * StopDist = K × N
  * InitialStop = Entry - StopDist
  * Sizing by method:
    - Stock (1): Shares = floor(R / StopDist)
    - Opt-DeltaATR (2): Contracts = floor(R / (K × N × Delta × 100))
    - Opt-MaxLoss (3): Contracts = floor(R / (MaxLoss × 100))
- Writes outputs to TradeEntry: F5 (R), F6 (StopDist), F7 (InitialStop), F8 (Shares), F9 (Contracts)

Return code only. Use WorksheetFunction.Floor_Precise for rounding.
```

5) **Heat and cooldown checks** *(10 min, data processing)*
```
Write VBA functions in module TF_Data:

Function PortfolioHeatAfter(addR As Double) As Double
  - Loops through tblPositions
  - Sums TotalOpenR (column 6) where Status (column 8) <> "Closed"
  - Returns sum + addR

Function BucketHeatAfter(bucket As String, addR As Double) As Double
  - Loops through tblPositions
  - Sums TotalOpenR where Bucket (column 2) = bucket AND Status <> "Closed"
  - Returns sum + addR

Function IsBucketInCooldown(bucket As String) As Boolean
  - Loops through tblBuckets
  - For matching bucket, checks if CooldownActive (column 7) = TRUE AND CooldownEndsOn (column 8) >= Today
  - Returns TRUE if in cooldown, else FALSE

Return code only.
```

6) **Save decision with hard gates** *(15 min, orchestration)*
```
Write a VBA Sub SaveDecision in module TF_UI that:
1. Reads banner (TradeEntry!A2), ticker (B6), bucket (B8), addR (F5)
2. Validates 5 hard gates (exit with MsgBox if any fail):
   a. Banner must contain "GREEN"
   b. Ticker must exist in tblCandidates with Date = Today (call helper IsTickerInCandidates)
   c. Impulse timer: (Now - Control!A1) >= 2 minutes
   d. NOT IsBucketInCooldown(bucket)
   e. PortfolioHeatAfter(addR) <= HeatCap_H_pct × Equity_E
   f. BucketHeatAfter(bucket, addR) <= BucketHeatCap_pct × Equity_E
3. If all pass:
   - Call AppendDecisionRow (stub: adds row to tblDecisions with all TradeEntry inputs)
   - Call UpdatePositions (stub: opens or adds to position in tblPositions)
   - Clear banner (A2), reset color
   - MsgBox success

Also write helper IsTickerInCandidates(ticker As String, tradeDate As Date) As Boolean

Return code only.
```

7) **FINVIZ integration** *(5 min, web automation)*
```
Write VBA in module TF_Presets:

Sub OpenPreset()
  - Reads selected preset name from TradeEntry!B5
  - Looks up QueryString in tblPresets
  - Opens URL: "https://finviz.com/screener.ashx?" & QueryString in default browser (CreateObject("Shell.Application").ShellExecute)

Sub ImportCandidatesPrompt()
  - Shows InputBox for user to paste tickers (comma or line-separated)
  - Normalizes each ticker (uppercase, trim, max 5 chars)
  - Dedupes
  - Appends to tblCandidates with Date=Today, Preset=TradeEntry!B5
  - MsgBox count imported

Return code only.
```

### Python Prompts

8) **FINVIZ scraper** *(10 min, web scraping)*
```python
# Paste this into Excel's Python environment (Data > Python in Excel)

Write a Python module finviz_scraper.py with function:

def fetch_finviz_tickers(query_string: str) -> list:
    """
    Scrapes FINVIZ screener and returns tickers.
    - Builds URL: https://finviz.com/screener.ashx?{query_string}
    - Uses requests + BeautifulSoup to parse table with class 'table-light'
    - Extracts ticker symbols from 2nd column
    - Returns list of uppercase, deduped tickers
    - Handles HTTP errors gracefully (return empty list)
    """

Use requests, beautifulsoup4, pandas. Return only Python code.
```

9) **Heat calculator** *(5 min, data processing)*
```python
Write Python module heat_calculator.py with function:

def check_heat_caps(positions_df: pd.DataFrame,
                   add_r: float,
                   bucket: str,
                   equity: float,
                   port_cap_pct: float,
                   bucket_cap_pct: float) -> dict:
    """
    Validates trade against heat caps.
    - Filters positions_df to Status != 'Closed'
    - Sums TotalOpenR for portfolio heat
    - Filters to specified bucket for bucket heat
    - Compares to caps
    - Returns dict: {'portfolio_ok': bool, 'bucket_ok': bool, 'portfolio_heat': float, ...}
    """

Use pandas. Return only Python code.
```

### Integration Prompt

10) **VBA-Python bridge** *(5 min, hybrid)*
```
Write VBA function in module TF_Python_Bridge:

Function CallPythonFinvizScraper(queryString As String) As Variant
  - Writes formula to hidden cell Control!Z1: =PY("finviz_scraper.fetch_finviz_tickers", "<queryString>")
  - Waits for calc (DoEvents)
  - Reads result as array
  - Returns array of tickers

Use this in ImportCandidatesPrompt instead of manual paste (if Python available).

Return code only.
```

---

## Gherkin / Cucumber Acceptance Tests

### 1) Banner logic
```
Feature: Checklist GO/NO-GO banner
  Scenario: All checks pass
    Given FromPreset, TrendPass, LiquidityPass, TVConfirm, EarningsOK, JournalOK are TRUE
    When I click Evaluate
    Then the banner shows "GREEN"
    And no reasons are listed

  Scenario: One check missing
    Given all checks except JournalOK are TRUE
    When I click Evaluate
    Then the banner shows "YELLOW"
    And reasons include "JournalOK missing"

  Scenario: Two checks missing
    Given TrendPass and LiquidityPass are FALSE
    When I click Evaluate
    Then the banner shows "RED"
```

### 2) Impulse timer
```
Feature: 2-minute impulse brake
  Scenario: Attempt to save too early
    Given I evaluated and got GREEN at 10:00:00
    And it is now 10:01:00
    When I click Save Decision
    Then I see a message "2-minute cool-off not elapsed"
    And the decision is not saved
```

### 3) Heat caps
```
Feature: Enforce portfolio and bucket heat caps
  Scenario: Portfolio heat exceeded
    Given Current open heat is $380 and HeatCap_H_pct × Equity_E = $400
    When I attempt a new trade requiring $75 R
    Then the system blocks with "Portfolio heat would exceed cap"
```

### 4) Candidate gating
```
Feature: Only trade tickers from today's candidates
  Scenario: Ticker not imported today
    Given "MSFT" is not in today's Candidates
    When I try to Save Decision for "MSFT"
    Then the system blocks with "Ticker not in today's Candidates"
```

### 5) Bucket cooldown
```
Feature: Bucket cooldown halts entries
  Scenario: Tech/Comm in cooldown
    Given "Tech/Comm" bucket has CooldownActive = TRUE until tomorrow
    When I try to Save Decision for a Tech/Comm trade
    Then the system blocks with "Bucket is in cooldown"
```

### 6) Sizing math (Delta-ATR)
```
Feature: Option sizing by Delta-ATR
  Scenario: Delta-ATR sizing yields at least 1 contract
    Given E=10000, r=0.0075, N=1.2, K=2, Delta=0.30
    When I Recalc Sizing
    Then R = $75 and Contracts = floor(75 / (2*1.2*0.3*100)) = 1
```

---

## Sample Walkthrough (MSFT, $10k portfolio)

**Inputs**: Preset=TF_BREAKOUT_LONG, Ticker=MSFT, Sector=Technology → Bucket=Tech/Comm, Entry=420.00, N=1.20, K=2, Method=Opt-DeltaATR, Delta=0.30, DTE=30–45. Checks: all TRUE. Settings: `Equity_E=10000`, `RiskPct_r=0.0075`.

**Evaluate** → Banner=GREEN, start impulse timer.
**Recalc Sizing** → R=$75; StopDist=2.40; InitialStop=417.60; Shares=31; Contracts=1; Add1=420.60; Add2=421.20; Add3=421.80.
**Save Decision** → passes caps (`HeatCap_H_pct=0.04 → $400`, `BucketHeatCap_pct=0.015 → $150`), logs row in Decisions, opens/updates Positions.

---

## Implementation Roadmap

**Phase 1: Foundation (Day 1, ~2 hours)**

1. **Create blank workbook**, save as `TrendFollowing_TradeEntry.xlsm`
2. **Enable Developer tab** (File → Options → Customize Ribbon → Developer)
3. **Open VBA Editor** (Alt+F11)
4. **Create 5 standard modules**:
   - Insert → Module (rename to `TF_Utils`)
   - Repeat for: `TF_Data`, `TF_UI`, `TF_Presets`, `TF_Python_Bridge`
5. **Generate utility code** using Mini LLM Prompt #1 (Structure seeding)
   - Paste into `TF_Data` module
   - Run `EnsureStructure` from Immediate window: `Call EnsureStructure`
   - Verify all 8 sheets created
   - Check Summary sheet has named ranges with default values

**Phase 2: UI Layout (Day 1, ~1 hour)**

6. **Manually build TradeEntry sheet** using Cell Reference Map:
   - Merge A1:F1, add title
   - Add labels in column A (rows 5-25)
   - Format cells for inputs (B column) and outputs (F column)
   - Hide column C (stores hidden values)
7. **Add Form Controls** (Developer tab → Insert → Form Controls):
   - Option buttons for Method (B13:B15), link to C13
   - Checkboxes for checklist (B20:B25), link to C20:C25
   - Spin button for K value (optional enhancement)
8. **Add Command Buttons** (Developer tab → Insert → Button):
   - Place 6 buttons as per Cell Reference Map (rows 28-30)
   - Assign macros later (leave unassigned for now)

**Phase 3: Core Logic (Day 2, ~3 hours)**

9. **Generate VBA code** using Mini LLM Prompts #2-7:
   - Paste each into appropriate module
   - Fix any syntax errors (missing End Sub, typos)
10. **Wire buttons to macros**:
   - Right-click each button → Assign Macro
   - Map to: OpenPreset, ImportCandidatesPrompt, EvaluateChecklist, RecalcSizing, SaveDecision
11. **Test basic flow** without data:
   - Click Evaluate → should show RED banner (no checks ticked)
   - Tick all checkboxes → click Evaluate → should show GREEN
   - Fill dummy values (Entry=100, N=2, K=2) → click Recalc Sizing → outputs populate

**Phase 4: Python Integration (Day 2, ~1 hour) [OPTIONAL]**

12. **Enable Python in Excel** (requires Microsoft 365 Insider):
   - Data tab → Python → Enable
13. **Create Python cells** in Control sheet:
   - Z5: `=PY("import requests, bs4, pandas as pd; 'ready'")`
   - Verify no errors
14. **Generate Python code** using Mini LLM Prompts #8-9:
   - Paste into Excel Python cells or external .py files
15. **Generate VBA-Python bridge** using Prompt #10:
   - Paste into `TF_Python_Bridge` module
16. **Update ImportCandidatesPrompt** to optionally call Python scraper:
   ```vba
   If MsgBox("Use Python scraper?", vbYesNo) = vbYes Then
       tickers = CallPythonFinvizScraper(queryString)
   Else
       ' Fallback to manual paste
   End If
   ```

**Phase 5: Data Validation & Formatting (Day 3, ~1 hour)**

17. **Set Data Validation** (if not already set by BindControls):
   - Select B5 → Data → Data Validation → List → Source: `=tblPresets[Name]`
   - Repeat for B6, B7, B8 per Cell Reference Map
18. **Add Conditional Formatting** for heat bars:
   - Select F14:F15
   - Conditional Formatting → Data Bars → Solid Fill (green/amber/red)
   - Adjust min/max based on cap values
19. **Format output cells**:
   - F5: Currency format ($0.00)
   - F6-F7: Number format (0.00)
   - F8-F9: Integer format (0)
   - F10-F12: Number format (0.00)

**Phase 6: Testing & Refinement (Day 3, ~2 hours)**

20. **Manual test each Gherkin scenario**:
   - Banner logic (all checks, 1 missing, 2 missing)
   - Impulse timer (attempt save before 2 min)
   - Heat caps (simulate high open positions)
   - Candidate gating (try non-imported ticker)
   - Cooldown (manually set bucket cooldown flag)
21. **Seed test data**:
   - Add 2-3 rows to tblCandidates with today's date
   - Add 2-3 rows to tblPositions with Open status
   - Run full workflow: Import → Evaluate → Recalc → Save
22. **Fix bugs**:
   - Add error handlers (`On Error Resume Next` where appropriate)
   - Improve MsgBox messages with clearer guidance
   - Add validation for empty/invalid inputs

**Phase 7: Documentation & Deployment (Day 4, ~1 hour)**

23. **Create Quick Start sheet**:
   - Brief instructions for daily workflow
   - Troubleshooting common errors
   - Link to this plan document
24. **Protect worksheets** (optional):
   - Protect Summary sheet (allow only B2:B8 edits)
   - Lock data tables, allow only TradeEntry edits
25. **Backup and version control**:
   - Save backup copy
   - If using git: export modules as .bas files (File → Export File)
26. **Train on live workflow**:
   - Run through 3-5 real trades using TradingView validation
   - Adjust settings (Equity_E, risk%, caps) based on account size

---

## Definition of Done

**Functional Requirements:**
- ✅ One **Trade Entry** sheet performs the entire trade workflow (no context switching)
- ✅ **GREEN-only** saves enforced with 5 hard gates (banner, candidate, timer, cooldown, heat caps)
- ✅ Sizing works for all 3 methods: **Stock**, **Opt-DeltaATR**, **Opt-MaxLoss**
- ✅ Decisions & Positions tables update automatically on save
- ✅ Heat preview shows real-time portfolio & bucket heat vs caps with color bars
- ✅ 2-minute impulse timer prevents impulsive entries
- ✅ Bucket cooldown activates after N stop-outs in M days

**UI Requirements:**
- ✅ Dropdowns auto-populate from tables (presets, tickers, sectors, buckets)
- ✅ Banner shows GREEN/YELLOW/RED with clear reason strings
- ✅ Method-specific fields (Delta, DTE, MaxLoss) auto-hide when not applicable
- ✅ All 6 buttons work without errors: Open Preset, Import Candidates, Evaluate, Recalc Sizing, Save Decision, Start Timer
- ✅ Checkboxes update banner in real-time (optional: add Worksheet_Change event)

**Data Integrity:**
- ✅ All tables have correct headers per Data Dictionary
- ✅ Named ranges exist on Summary with sensible defaults
- ✅ Control sheet stores impulse timer timestamp
- ✅ No #REF! or #NAME? errors in any formula cells
- ✅ Decisions log captures all 20 fields including heat metrics

**Python Integration (Optional but Recommended):**
- ✅ Python scraper fetches FINVIZ tickers (eliminates manual copy/paste)
- ✅ Heat calculator validates caps using pandas (faster for large position tables)
- ✅ VBA-Python bridge functions work without errors

**Testing:**
- ✅ All 6 Gherkin scenarios pass manual testing
- ✅ Can reproduce MSFT walkthrough example in < 2 minutes
- ✅ Error handlers prevent crashes (all critical subs have On Error handlers)
- ✅ Works with empty tables (no division by zero, null reference errors)

**Documentation:**
- ✅ Quick Start sheet guides new users
- ✅ Cell Reference Map matches actual layout
- ✅ Module code includes comments for key logic blocks

**Performance:**
- ✅ Evaluate button responds in < 1 second
- ✅ Recalc Sizing computes in < 0.5 seconds
- ✅ Save Decision (with heat checks) completes in < 2 seconds
- ✅ No screen flickering (use `Application.ScreenUpdating = False` in loops)
