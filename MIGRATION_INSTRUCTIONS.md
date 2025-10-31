# Database Migration Instructions

## Problem
You're seeing this error when creating a trade session:

```
Failed to create session with options: SQL logic error: table trade_sessions has no column named instrument_type (1)
```

This means your database schema is outdated and needs the new options trading columns added in Phases 1-7.

## Solution: Run the Migration

### Quick Fix (Easiest)

**Windows:**
1. Open PowerShell or Command Prompt in the project folder
2. Run: `.\migrate-db.exe`
3. Done! You should see:
   ```
   ✅ Migration completed successfully!
   ```

The migration tool will automatically:
- Add all 27 missing columns to the `trade_sessions` table
- Set appropriate defaults so existing sessions aren't affected
- Enable options trading features

### What Gets Added

The migration adds these column groups:

**Options Metadata:**
- `instrument_type` - 'STOCK' or 'OPTION'
- `options_strategy` - Long Call, Iron Condor, etc.
- `entry_date` - When position entered
- `primary_expiration_date` - Option expiration
- `dte` - Days to expiration at entry
- `roll_threshold_dte` - When to roll (default 21 days)
- `time_exit_mode` - 'None', 'Close', or 'Roll'
- `legs_json` - Multi-leg options structure
- `net_debit` - Total cost
- `max_profit` / `max_loss` - Risk/reward
- `breakeven_lower` / `breakeven_upper` - Breakeven prices
- `underlying_at_entry` - Underlying stock price

**Pyramid Tracking:**
- `max_units` - Maximum position units (default 4)
- `add_step_n` - Add-on interval (default 0.5×N)
- `current_units` - Current units held
- `add_price_1` / `add_price_2` / `add_price_3` - Calculated add-on prices

**Breakout System:**
- `entry_lookback` - Entry breakout period (20 for System-1, 55 for System-2)
- `exit_lookback` - Exit breakout period (default 10)

### Manual Method (Advanced)

If you prefer to run the migration manually:

1. Open SQLite browser or sqlite3 command line
2. Connect to your `trading.db` database
3. Run the SQL from: `backend/migrations/002_add_options_columns.sql`

### After Migration

Once migration is complete:
1. Restart the GUI: `.\tf-gui.exe`
2. Click "Start New Trade"
3. You should now be able to create sessions without errors!

The GUI will now support:
- ✅ Options trading strategies (26 types)
- ✅ Expiration tracking and DTE alerts
- ✅ Pyramid add-on price calculations
- ✅ System-1 vs System-2 breakout selection
- ✅ Multi-leg options structures
- ✅ Trade calendar with options metadata

### Troubleshooting

**Error: "no such table: trade_sessions"**
- Your database is too old. Run `.\tf-engine.exe init` first, then run migration

**Error: "duplicate column name"**
- Columns already exist. Migration already ran successfully. Ignore this error.

**Error: "cannot open database"**
- Make sure `trading.db` exists in the current directory
- Check file permissions (make sure you can write to it)
- Close any other programs accessing the database

### Safe Migration

✅ **Non-destructive** - Adds columns only, doesn't delete data
✅ **Backward compatible** - Existing sessions still work
✅ **Default values** - All new columns have sensible defaults
✅ **Automatic** - Just run the .exe, no manual SQL needed

---

**Created:** 2025-10-30
**Purpose:** Enable options trading enhancement (Phases 1-7)
**Status:** Required for all users upgrading from pre-options version
