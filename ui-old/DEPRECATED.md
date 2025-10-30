# ⚠️ DEPRECATED - Old Browser-Based UI

This directory contains the old Svelte browser-based UI that has been **deprecated** in favor of the native Fyne GUI.

## Why Deprecated?

The browser-based approach had several limitations:
- ❌ Missing most backend features (Scanner, full position sizing, etc.)
- ❌ Incomplete implementation (Calendar broken, basic styling)
- ❌ Required browser + HTTP server overhead
- ❌ Less responsive than native
- ❌ Not aligned with project philosophy of discipline enforcement

## Replacement

The **native Fyne GUI** in `ui/` directory replaces this entirely:
- ✅ All backend features exposed
- ✅ Professional native UI
- ✅ Single executable (no browser required)
- ✅ Direct in-process function calls
- ✅ 49MB standalone binary

## Running the New GUI

```bash
cd ui
./tf-gui.exe
```

See `ui/README.md` for full documentation.

---

**This old UI is kept only for historical reference. Do not use for actual trading.**
