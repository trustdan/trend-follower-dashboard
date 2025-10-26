#!/usr/bin/env python3
"""
VBA Import Automation Script for Excel Trading Workbook

This script automates the import of VBA modules into an Excel workbook.

Requirements:
    - Python 3.7+
    - openpyxl (install: pip install openpyxl)
    - xlwings (install: pip install xlwings)
    - Windows OS (for COM automation)

Usage:
    python import_to_excel.py [workbook_path]

    If no workbook path provided, creates a new workbook:
        TrendFollowing_TradeEntry.xlsm

Note:
    This script uses Windows COM automation (requires Windows + Excel installed).
    For manual import, see VBA_SETUP_GUIDE.md
"""

import os
import sys
from pathlib import Path

# Try to import required libraries
try:
    import win32com.client
    COM_AVAILABLE = True
except ImportError:
    COM_AVAILABLE = False
    print("⚠ WARNING: win32com not available. Install with: pip install pywin32")

def get_script_dir():
    """Returns the directory where this script is located."""
    return Path(__file__).parent.resolve()

def find_vba_modules():
    """Finds all VBA module files in the VBA/ folder."""
    script_dir = get_script_dir()
    vba_dir = script_dir / "VBA"

    if not vba_dir.exists():
        print(f"❌ Error: VBA folder not found at {vba_dir}")
        return None

    # Find all .bas and .cls files
    modules = {
        'standard': list(vba_dir.glob("*.bas")),
        'class': list(vba_dir.glob("*.cls"))
    }

    return modules

def import_vba_modules_com(workbook_path=None):
    """
    Import VBA modules using Windows COM automation.

    Args:
        workbook_path: Path to existing .xlsm file, or None to create new
    """
    if not COM_AVAILABLE:
        print("❌ Cannot proceed: pywin32 not installed")
        print("   Install with: pip install pywin32")
        return False

    script_dir = get_script_dir()

    # Start Excel
    print("Starting Excel...")
    excel = win32com.client.Dispatch("Excel.Application")
    excel.Visible = True
    excel.DisplayAlerts = False

    # Open or create workbook
    if workbook_path and Path(workbook_path).exists():
        print(f"Opening workbook: {workbook_path}")
        wb = excel.Workbooks.Open(str(Path(workbook_path).resolve()))
    else:
        print("Creating new workbook...")
        wb = excel.Workbooks.Add()
        default_path = script_dir / "TrendFollowing_TradeEntry.xlsm"
        print(f"Will save as: {default_path}")

    # Find VBA modules
    modules = find_vba_modules()
    if not modules:
        excel.Quit()
        return False

    # Import standard modules
    vb_project = wb.VBProject
    imported_count = 0

    print("\n📥 Importing standard modules...")
    for module_path in modules['standard']:
        try:
            vb_project.VBComponents.Import(str(module_path))
            print(f"  ✅ {module_path.name}")
            imported_count += 1
        except Exception as e:
            print(f"  ❌ {module_path.name}: {e}")

    print("\n📥 Importing class modules...")
    for module_path in modules['class']:
        try:
            # Special handling for ThisWorkbook (replace existing)
            if module_path.stem == "ThisWorkbook":
                # Remove existing ThisWorkbook code
                thisworkbook = vb_project.VBComponents("ThisWorkbook")
                code_module = thisworkbook.CodeModule
                if code_module.CountOfLines > 0:
                    code_module.DeleteLines(1, code_module.CountOfLines)

                # Import new code
                with open(module_path, 'r', encoding='utf-8') as f:
                    lines = f.readlines()

                # Skip VBA header lines (Attribute statements)
                code_start = 0
                for i, line in enumerate(lines):
                    if not line.startswith('Attribute') and not line.startswith('VERSION'):
                        code_start = i
                        break

                new_code = ''.join(lines[code_start:])
                code_module.AddFromString(new_code)
                print(f"  ✅ {module_path.name} (replaced)")
                imported_count += 1
            else:
                # Other class modules can be imported normally
                vb_project.VBComponents.Import(str(module_path))
                print(f"  ✅ {module_path.name}")
                imported_count += 1

        except Exception as e:
            print(f"  ❌ {module_path.name}: {e}")

    # Save workbook
    if not workbook_path:
        save_path = script_dir / "TrendFollowing_TradeEntry.xlsm"
        wb.SaveAs(str(save_path), FileFormat=52)  # 52 = xlOpenXMLWorkbookMacroEnabled
        print(f"\n💾 Saved as: {save_path}")
    else:
        wb.Save()
        print(f"\n💾 Saved: {workbook_path}")

    print(f"\n✅ Import complete! {imported_count} modules imported.")
    print("\nNext steps:")
    print("1. Run 'EnsureStructure' macro (Alt+F11 → Ctrl+G → type 'EnsureStructure' → Enter)")
    print("2. Build TradeEntry UI (see VBA_SETUP_GUIDE.md Part 2)")
    print("3. Test workflow")

    # Don't quit Excel - leave it open for user
    excel.DisplayAlerts = True

    return True

def create_manual_import_script():
    """
    Creates a VBA script that can be run inside Excel to import modules.
    This is a fallback for when Python COM automation doesn't work.
    """
    script_dir = get_script_dir()
    vba_dir = script_dir / "VBA"

    modules = find_vba_modules()
    if not modules:
        return

    script_content = [
        "' ============================================================================",
        "' VBA Module Import Script",
        "' Run this in the VBA Immediate Window (Ctrl+G) to import all modules",
        "' ============================================================================",
        "",
        "Sub ImportAllModules()",
        "    Dim modulePath As String",
        f"    Dim vbaFolder As String",
        f"    vbaFolder = ThisWorkbook.Path & \"\\VBA\\\"",
        "",
        "    ' Import standard modules",
    ]

    for module_path in modules['standard']:
        script_content.append(f"    ThisWorkbook.VBProject.VBComponents.Import vbaFolder & \"{module_path.name}\"")

    script_content.extend([
        "",
        "    ' Import class modules (except ThisWorkbook - needs manual handling)",
    ])

    for module_path in modules['class']:
        if module_path.stem != "ThisWorkbook":
            script_content.append(f"    ThisWorkbook.VBProject.VBComponents.Import vbaFolder & \"{module_path.name}\"")

    script_content.extend([
        "",
        "    MsgBox \"Modules imported! Note: ThisWorkbook must be imported manually.\", vbInformation",
        "End Sub",
        "",
        "' ============================================================================",
        "' Manual Import Instructions:",
        "' ============================================================================",
        "' 1. Press Alt+F11 to open VBA Editor",
        "' 2. Press Ctrl+G to open Immediate Window",
        "' 3. Paste this entire script into a new module (Insert → Module)",
        "' 4. Run: ImportAllModules",
        "' 5. Manually import ThisWorkbook.cls:",
        "'    - File → Import File → select VBA/ThisWorkbook.cls",
        "'    - Replace the existing ThisWorkbook module",
        "' ============================================================================",
    ])

    output_path = script_dir / "VBA_IMPORT_SCRIPT.txt"
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write('\n'.join(script_content))

    print(f"✅ Created manual import script: {output_path}")
    print("   Copy contents to VBA Editor if Python automation doesn't work")

def main():
    """Main entry point."""
    print("=" * 70)
    print("VBA Module Import Automation")
    print("=" * 70)

    # Check for workbook path argument
    workbook_path = sys.argv[1] if len(sys.argv) > 1 else None

    if workbook_path and not Path(workbook_path).exists():
        print(f"❌ Error: Workbook not found: {workbook_path}")
        return 1

    # Show what we found
    modules = find_vba_modules()
    if modules:
        print(f"\n📁 Found VBA modules:")
        print(f"   Standard: {len(modules['standard'])} files")
        for m in modules['standard']:
            print(f"     - {m.name}")
        print(f"   Class: {len(modules['class'])} files")
        for m in modules['class']:
            print(f"     - {m.name}")
    else:
        return 1

    # Check if COM is available
    if not COM_AVAILABLE:
        print("\n⚠ Python COM automation not available (pywin32 not installed)")
        print("  Creating manual import script instead...")
        create_manual_import_script()
        print("\n📖 See VBA_SETUP_GUIDE.md for manual import instructions")
        return 1

    # Proceed with automated import
    print("\n🚀 Starting automated import...")
    print("   This will open Excel and import all VBA modules.")

    success = import_vba_modules_com(workbook_path)

    if success:
        # Also create the manual script as backup
        create_manual_import_script()
        return 0
    else:
        return 1

if __name__ == "__main__":
    sys.exit(main())
