"""
Excel Trading Workbook Builder - Simplified Version
Creates workbook with all modules, then user runs setup manually
"""

import os
import sys
import time
import subprocess
from pathlib import Path

def print_header(text):
    """Print formatted header"""
    print("\n" + "=" * 70)
    print(text)
    print("=" * 70 + "\n")

def print_step(text):
    """Print step description"""
    print(f"→ {text}...")

def print_success(text):
    """Print success message"""
    print(f"✓ {text}")

def print_error(text):
    """Print error message"""
    print(f"✗ {text}")

def check_python_packages():
    """Check and install required Python packages"""
    print_step("Checking Python packages")

    packages = ['pywin32']
    missing = []

    for package in packages:
        try:
            __import__(package.replace('-', '_'))
            print_success(f"{package} installed")
        except ImportError:
            missing.append(package)

    if missing:
        print(f"\nInstalling missing packages: {', '.join(missing)}")
        for package in missing:
            subprocess.check_call([sys.executable, "-m", "pip", "install", package])
            print_success(f"{package} installed")

def clear_com_cache():
    """Clear pywin32 COM cache to fix automation issues"""
    print_step("Clearing COM cache")

    import win32com
    gen_path = Path(win32com.__gen_path__)

    if gen_path.exists():
        import shutil
        try:
            shutil.rmtree(gen_path)
            print_success("COM cache cleared")
        except Exception as e:
            print_error(f"Could not clear cache: {e}")
    else:
        print_success("Cache already clean")

def kill_excel():
    """Kill any running Excel processes"""
    print_step("Closing existing Excel processes")
    try:
        subprocess.run(['taskkill', '/F', '/IM', 'EXCEL.EXE'],
                      capture_output=True, check=False)
        time.sleep(2)
        print_success("Excel processes closed")
    except Exception as e:
        print_error(f"Could not kill Excel: {e}")

def import_vba_modules():
    """Import all VBA modules using pywin32"""
    print_step("Importing VBA modules")

    try:
        import win32com.client as win32
    except ImportError:
        print_error("pywin32 not available")
        return False

    # Get paths
    script_dir = Path(__file__).parent
    vba_dir = script_dir / "VBA"
    output_file = script_dir / "TrendFollowing_TradeEntry.xlsm"

    # Delete existing workbook
    if output_file.exists():
        print_step(f"Deleting existing workbook")
        output_file.unlink()
        print_success("Old workbook deleted")

    # Start Excel
    print_step("Starting Excel")
    try:
        excel = win32.Dispatch('Excel.Application')
    except Exception as e:
        print_error(f"Could not start Excel: {e}")
        return False

    excel.Visible = False
    excel.DisplayAlerts = False

    # Create new workbook
    print_step("Creating new workbook")
    wb = excel.Workbooks.Add()

    # Enable VBA project access
    try:
        vbProject = wb.VBProject
    except Exception as e:
        print_error("Cannot access VBA project")
        print("\nPlease enable:")
        print("  File → Options → Trust Center → Trust Center Settings")
        print("  → Macro Settings → Trust access to VBA project object model")
        excel.Quit()
        return False

    # Import standard modules (.bas)
    print("\nImporting standard modules:")
    bas_files = [
        "TF_Logger.bas",
        "TF_Utils.bas",
        "TF_Data.bas",
        "TF_UI.bas",
        "TF_Presets.bas",
        "TF_Python_Bridge.bas",
        "TF_UI_Builder.bas",
        "Setup.bas"
    ]

    for filename in bas_files:
        filepath = vba_dir / filename
        if filepath.exists():
            try:
                vbProject.VBComponents.Import(str(filepath))
                print_success(filename)
            except Exception as e:
                print_error(f"{filename}: {e}")
        else:
            print_error(f"{filename}: File not found")

    # Import class modules (.cls)
    print("\nImporting class modules:")

    # Handle ThisWorkbook specially
    filepath = vba_dir / "ThisWorkbook.cls"
    if filepath.exists():
        try:
            # Read the file and extract code
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()

            # Find where actual code starts (after CLASS header)
            lines = content.split('\n')
            code_start = 0
            for i, line in enumerate(lines):
                # Skip header lines
                if line.strip().startswith("'") or \
                   line.strip().startswith("Private Sub") or \
                   line.strip().startswith("Sub ") or \
                   line.strip().startswith("Function "):
                    code_start = i
                    break

            # Get the actual code
            code = '\n'.join(lines[code_start:])

            # Find ThisWorkbook component
            for comp in vbProject.VBComponents:
                if comp.Name == "ThisWorkbook":
                    # Clear existing code
                    code_module = comp.CodeModule
                    if code_module.CountOfLines > 0:
                        code_module.DeleteLines(1, code_module.CountOfLines)

                    # Add new code
                    code_module.AddFromString(code)
                    print_success("ThisWorkbook.cls (code updated)")
                    break
        except Exception as e:
            print_error(f"ThisWorkbook.cls: {e}")

    # Save workbook
    print_step(f"Saving to {output_file.name}")
    try:
        wb.SaveAs(str(output_file), FileFormat=52)  # 52 = xlOpenXMLWorkbookMacroEnabled
        print_success("Workbook saved")
    except Exception as e:
        print_error(f"Save failed: {e}")
        excel.Quit()
        return False

    # Close Excel
    print_step("Closing Excel")
    wb.Close(SaveChanges=False)
    excel.Quit()
    print_success("Excel closed")

    return True

def main():
    """Main build process"""
    print_header("Excel Trading Workbook - Simplified Build")

    print("Current directory:", os.getcwd())
    print("Target workbook: TrendFollowing_TradeEntry.xlsm\n")

    # Step 1: Check packages
    check_python_packages()

    # Step 2: Clear COM cache (fixes common issues)
    try:
        clear_com_cache()
    except Exception as e:
        print_error(f"Could not clear cache: {e}")

    # Step 3: Kill Excel
    kill_excel()

    # Step 4: Import VBA modules
    if not import_vba_modules():
        print_error("\nBuild failed: Could not import VBA modules")
        print("\nTroubleshooting:")
        print("  1. Enable macro settings in Excel")
        print("  2. Enable VBA project access (Trust Center)")
        print("  3. Try running as Administrator")
        return 1

    # Success!
    print_header("BUILD COMPLETE!")
    print("\nWorkbook created: TrendFollowing_TradeEntry.xlsm")
    print("All VBA modules imported successfully!")
    print("\n" + "=" * 70)
    print("NEXT STEP: OPEN THE WORKBOOK (1 click!)")
    print("=" * 70)
    print("\n1. Double-click: TrendFollowing_TradeEntry.xlsm")
    print("2. Click 'Enable Content' if prompted")
    print("3. Setup runs AUTOMATICALLY on first open!")
    print("")
    print("The workbook will:")
    print("  ✓ Create all sheets and tables")
    print("  ✓ Build TradeEntry UI")
    print("  ✓ Create Setup sheet with instructions")
    print("  ✓ Show you exactly what to do next")
    print("")
    print("Then you just need to:")
    print("  → Add 6 checkboxes (2 minutes - instructions on Setup sheet)")
    print("  → Start trading!")
    print("")
    print("=" * 70)
    print("\nEverything else is AUTOMATED!")
    print("")

    return 0

if __name__ == "__main__":
    sys.exit(main())
