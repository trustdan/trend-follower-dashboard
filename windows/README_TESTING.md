# M21 Phase 4 - Automated Integration Testing

**TL;DR:** One command to run all tests

```cmd
run-integration-tests.bat
```

## Quick Start

1. **One-time setup:** Enable VBA project access
   - Excel → File → Options → Trust Center → Trust Center Settings
   - Macro Settings → Check "Trust access to VBA project object model"

2. **Run tests:**
   ```cmd
   cd windows
   run-integration-tests.bat
   ```

3. **Check results:**
   - Excel opens with "Integration Tests" worksheet
   - Log file opens: `logs/integration-tests-YYYYMMDD-HHMMSS.log`

## What Gets Tested

**19 automated tests (1-2 minutes):**
- ✅ Position Sizing: 4 tests (stock, overrides, option delta-ATR, option max-loss)
- ✅ Checklist Evaluation: 5 tests (GREEN, YELLOW, RED banners, persistence)
- ✅ Heat Management: 4 tests (clean state, portfolio cap, bucket cap, edge cases)
- ✅ Save Decision: 6 tests (happy path, Gate 1/2/5 rejections)

**6 manual tests (optional, timing-dependent):**
- Gate 3: 2-minute impulse brake
- Gate 4: Bucket cooldown (24 hours)
- Cumulative heat with open positions

## Files

**Automation:**
- `run-integration-tests.bat` - One-click test runner
- `../excel/vba/TFIntegrationTests.bas` - VBA test module (auto-imported)

**Documentation:**
- `../docs/milestones/M21_PHASE4_AUTOMATED.md` - Complete automation guide
- `../docs/milestones/M21_PHASE4_TEST_SCRIPTS.md` - Manual test procedures
- `../docs/milestones/M21_PHASE4_CHECKLIST.md` - Manual test checklist

**Test Data:**
- `../test-data/phase4-test-data.sql` - Pre-populated test data
- `../test-data/phase4-test-scenarios.csv` - Test scenarios table

## Success Criteria

**All tests PASS:**
- Position sizing calculations correct
- Checklist banners enforce correctly (only GREEN allows save)
- Heat caps enforced (4% portfolio, 1.5% bucket)
- Save decision gates reject invalid trades

## Troubleshooting

**"Trust access to VBA project" error:**
- Enable in Excel options (see Quick Start step 1)

**"TradingPlatform.xlsm not found":**
- Run `setup-all.bat` first

**"Engine not working":**
- Verify: `tf-engine.exe --version`
- Rebuild if needed: `cd .. && go build -o windows/tf-engine.exe ./cmd/tf-engine`

**Tests fail with "Failed to parse JSON":**
- Check `TradingSystem_Debug.log`
- Verify engine output: `tf-engine.exe get-settings`
- Reset: Re-run `setup-all.bat`

## More Information

See: `../docs/milestones/M21_PHASE4_AUTOMATED.md` for complete guide

---

**Created:** 2025-10-27 (M21 Phase 4)
**Method:** Automated via VBA + Batch script
**Duration:** 1-2 minutes (19 tests)
