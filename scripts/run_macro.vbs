' Usage: cscript //nologo run_macro.vbs "<workbook_path>" "Macro1" "Macro2" ...
Dim xl, wb, path, i
If WScript.Arguments.Count < 2 Then
  WScript.Echo "Usage: run_macro.vbs <workbook_path> <MacroName1> [MacroName2] ..."
  WScript.Quit 1
End If

path = WScript.Arguments(0)
Set xl = CreateObject("Excel.Application")
xl.DisplayAlerts = False
xl.Visible = False
Set wb = xl.Workbooks.Open(path)

For i = 1 To WScript.Arguments.Count - 1
  xl.Run WScript.Arguments(i)
Next

wb.Close True
xl.Quit
