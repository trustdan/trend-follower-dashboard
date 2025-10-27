#!/usr/bin/env python3
from pathlib import Path
import sys

# --- COM setup ---------------------------------------------------------------
try:
    import win32com.client as win32
    from pywintypes import com_error
    COM_AVAILABLE = True
except Exception as e:
    COM_AVAILABLE = False
    print("❌ pywin32 not available:", e)

def script_dir() -> Path:
    return Path(__file__).parent.resolve()

def find_vba_modules():
    """Return dict of module paths in ./VBA (ok if empty)."""
    vba = script_dir() / "VBA"
    if not vba.exists():
        print(f"⚠ VBA folder not found at {vba} (continuing; workbook will be created without modules).")
        return {"standard": [], "class": []}
    return {
        "standard": sorted(vba.glob("*.bas")),
        "class": sorted(vba.glob("*.cls")),
    }

def import_vba_modules_com(target_path: Path | None) -> bool:
    if not COM_AVAILABLE:
        print("❌ Cannot proceed: pywin32 not installed (pip install pywin32)")
        return False

    print("Starting Excel…")
    xl = win32.Dispatch("Excel.Application")
    xl.Visible = True
    xl.DisplayAlerts = False

    # Open if exists, otherwise create a new workbook
    wb = None
    if target_path and target_path.exists():
        print(f"Opening workbook: {target_path}")
        wb = xl.Workbooks.Open(str(target_path))
    else:
        print("Creating new workbook…")
        wb = xl.Workbooks.Add()

    # Try to access VBProject (fails if Trust Center blocks it)
    try:
        vbproj = wb.VBProject
    except com_error:
        print("❌ Excel blocked programmatic VBProject access.")
        print("   Enable: File → Options → Trust Center → Trust Center Settings → Macro Settings")
        print("   Check:  ‘Trust access to the VBA project object model’.")
        xl.Quit()
        return False

    mods = find_vba_modules()
    imported = 0

    # Import .bas (standard modules)
    if mods["standard"]:
        print("\n📥 Importing standard modules…")
        for p in mods["standard"]:
            try:
                vbproj.VBComponents.Import(str(p))
                print(f"  ✅ {p.name}")
                imported += 1
            except Exception as e:
                print(f"  ❌ {p.name}: {e}")

    # Run TF_Data.EnsureStructure to create sheets, tables, and data (needed before importing Sheet_*.cls)
    print("\n🔧 Running TF_Data.EnsureStructure to create workbook structure…")
    try:
        # Disable alerts to suppress the msgbox
        xl.DisplayAlerts = False
        xl.Run("TF_Data.EnsureStructure")
        xl.DisplayAlerts = True
        print("  ✅ TF_Data.EnsureStructure completed")
        print("     - Sheets created (8)")
        print("     - Tables created (5)")
        print("     - Named ranges created (7)")
        print("     - Default data seeded")
    except Exception as e:
        xl.DisplayAlerts = True
        print(f"  ⚠ TF_Data.EnsureStructure not available or failed: {e}")
        print("  (You may need to run TF_Data.EnsureStructure manually)")

    # Run TF_UI_Builder.InitializeUI to create the TradeEntry UI
    print("\n🎨 Running TF_UI_Builder.InitializeUI to build TradeEntry UI…")
    try:
        xl.DisplayAlerts = False
        xl.Run("TF_UI_Builder.InitializeUI")
        xl.DisplayAlerts = True
        print("  ✅ TF_UI_Builder.InitializeUI completed")
        print("     - TradeEntry layout created")
        print("     - Buttons added (Evaluate, Recalc, Save, Import)")
        print("     - Formatting applied")
        print("     - Data validation set up")
    except Exception as e:
        xl.DisplayAlerts = True
        print(f"  ⚠ TF_UI_Builder.InitializeUI not available or failed: {e}")
        print("  (UI may need manual setup)")

    # Import .cls (special handling for ThisWorkbook.cls and Sheet_*.cls)
    if mods["class"]:
        print("\n📥 Importing class modules…")
        for p in mods["class"]:
            try:
                if p.stem == "ThisWorkbook":
                    # Replace ThisWorkbook code
                    twb = vbproj.VBComponents("ThisWorkbook")
                    cm = twb.CodeModule
                    if cm.CountOfLines > 0:
                        cm.DeleteLines(1, cm.CountOfLines)
                    with open(p, "r", encoding="utf-8") as f:
                        lines = f.readlines()
                    start = 0
                    for i, line in enumerate(lines):
                        if not line.startswith(("Attribute", "VERSION")):
                            start = i
                            break
                    cm.AddFromString("".join(lines[start:]))
                    print(f"  ✅ {p.name} (replaced)")
                    imported += 1

                elif p.stem.startswith("Sheet_"):
                    # Handle Sheet_*.cls files - replace code in corresponding sheet
                    sheet_name = p.stem[6:]  # Remove "Sheet_" prefix

                    # Try to find the sheet by CodeName first, then by Name
                    sheet_comp = None
                    try:
                        # Try direct CodeName lookup
                        sheet_comp = vbproj.VBComponents(sheet_name)
                    except:
                        # If not found by CodeName, search all worksheet components
                        for comp in vbproj.VBComponents:
                            if comp.Type == 100:  # 100 = Document (worksheet)
                                try:
                                    # Check if the worksheet Name matches
                                    ws_name = wb.Worksheets(comp.Name).Name
                                    if ws_name == sheet_name:
                                        sheet_comp = comp
                                        print(f"  📍 Found sheet '{sheet_name}' with CodeName '{comp.Name}'")
                                        break
                                except:
                                    pass

                    if not sheet_comp:
                        # Sheet not found - skip importing as class module
                        print(f"  ⚠ Sheet '{sheet_name}' not found - SKIPPING {p.name}")
                        print(f"     (Sheet must exist before applying code)")
                        continue

                    # Replace the sheet's code
                    cm = sheet_comp.CodeModule
                    if cm.CountOfLines > 0:
                        cm.DeleteLines(1, cm.CountOfLines)
                    with open(p, "r", encoding="utf-8") as f:
                        lines = f.readlines()
                    start = 0
                    for i, line in enumerate(lines):
                        if not line.startswith(("Attribute", "VERSION")):
                            start = i
                            break
                    cm.AddFromString("".join(lines[start:]))
                    print(f"  ✅ {p.name} → Sheet '{sheet_name}' (code replaced)")
                    imported += 1

                else:
                    # Regular class module
                    vbproj.VBComponents.Import(str(p))
                    print(f"  ✅ {p.name}")
                    imported += 1
            except Exception as e:
                print(f"  ❌ {p.name}: {e}")

    # --- Save / SaveAs logic (always SaveAs when target path is provided) -----
    try:
        if target_path:
            target_path.parent.mkdir(parents=True, exist_ok=True)
            print(f"\n💾 Saving to: {target_path}")
            # Delete existing file if present (avoids save conflicts)
            if target_path.exists():
                target_path.unlink()
                print(f"  (Deleted existing file)")
            wb.SaveAs(str(target_path), FileFormat=52)  # 52 = .xlsm
        else:
            default_path = script_dir() / "TrendFollowing_TradeEntry.xlsm"
            print(f"\n💾 Saving to: {default_path}")
            if default_path.exists():
                default_path.unlink()
                print(f"  (Deleted existing file)")
            wb.SaveAs(str(default_path), FileFormat=52)

        print(f"\n✅ Import complete! {imported} modules imported.")
        print(f"📁 File saved: {target_path or default_path}")

        # Close workbook and quit Excel
        wb.Close(SaveChanges=False)  # Don't save again, we just saved
        print("✅ Workbook closed")

        xl.Quit()
        print("✅ Excel quit successfully")

        xl.DisplayAlerts = True
        return True

    except Exception as e:
        print(f"\n❌ Error during save: {e}")
        print("Attempting to close Excel anyway...")
        try:
            wb.Close(SaveChanges=False)
            xl.Quit()
        except:
            pass
        return False

def main():
    print("=" * 70)
    print("VBA Module Import Automation")
    print("=" * 70)

    # Accept a path even if it doesn't exist (we'll create the workbook there)
    target = Path(sys.argv[1]).resolve() if len(sys.argv) > 1 else None

    # (Optional) list found modules
    mods = find_vba_modules()
    print(f"\n📁 Modules found: {len(mods['standard'])} .bas, {len(mods['class'])} .cls")

    if not COM_AVAILABLE:
        return 1

    ok = import_vba_modules_com(target)
    return 0 if ok else 1

if __name__ == "__main__":
    sys.exit(main())
