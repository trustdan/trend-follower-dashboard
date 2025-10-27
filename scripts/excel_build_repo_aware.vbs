Option Explicit

Dim fso:   Set fso   = CreateObject("Scripting.FileSystemObject")
Dim shell: Set shell = CreateObject("WScript.Shell")

' --- Inputs ---
Dim wbPath, vbaFolder
If WScript.Arguments.Count >= 1 Then
  wbPath = WScript.Arguments(0)
Else
  wbPath = fso.GetAbsolutePathName("TrendFollowing_TradeEntry.xlsm")
End If
If WScript.Arguments.Count >= 2 Then
  vbaFolder = WScript.Arguments(1)
Else
  vbaFolder = "VBA"
End If

' --- Make Excel run in "silent setup" mode for Workbook_Open guards ---
shell.Environment("Process")("XL_SILENT_SETUP") = "1"

' Try to ensure VBOM (non-fatal)
On Error Resume Next
shell.RegWrite "HKCU\Software\Microsoft\Office\16.0\Excel\Security\AccessVBOM", 1, "REG_DWORD"
Err.Clear: On Error GoTo 0

Dim xl: Set xl = CreateObject("Excel.Application")
xl.DisplayAlerts = False
xl.Visible = False

Dim wb
If fso.FileExists(wbPath) Then
  Set wb = xl.Workbooks.Open(wbPath)
Else
  Set wb = xl.Workbooks.Add()
  wb.SaveAs wbPath, 52 ' 52 = xlOpenXMLWorkbookMacroEnabled
End If

Dim vbproj: Set vbproj = wb.VBProject

' --- Helpers ---
Function ReadTextSansAttributes(p)
  Dim s, t, ts, arr, i, line
  Set ts = fso.OpenTextFile(p, 1, False)
  s = ts.ReadAll
  ts.Close
  s = Replace(s, vbCrLf, vbLf)
  s = Replace(s, vbCr, vbLf)
  arr = Split(s, vbLf)
  For i = 0 To UBound(arr)
    line = Trim(arr(i))
    If LCase(Left(line, 7)) = "version" Then
    ElseIf LCase(Left(line, 16)) = "attribute vb_name" Then
    ElseIf LCase(Left(line, 23)) = "attribute vb_globalsnamespace" Then
    ElseIf LCase(Left(line, 20)) = "attribute vb_creatable" Then
    ElseIf LCase(Left(line, 24)) = "attribute vb_predeclaredid" Then
    ElseIf LCase(Left(line, 17)) = "attribute vb_exposed" Then
    Else
      t = t & arr(i) & vbCrLf
    End If
  Next
  ReadTextSansAttributes = t
End Function

Function FindSheetByCodeName(book, codeName)
  Dim ws
  For Each ws In book.Worksheets
    On Error Resume Next
    If StrComp(ws.CodeName, codeName, 0) = 0 Then
      Set FindSheetByCodeName = ws
      Exit Function
    End If
    Err.Clear
    On Error GoTo 0
  Next
  Set FindSheetByCodeName = Nothing
End Function

Function EnsureSheet(book, desiredTabName, desiredCodeName, atPosition)
  Dim ws, comp
  Set ws = FindSheetByCodeName(book, desiredCodeName)
  If ws Is Nothing Then
    If atPosition >= 1 And atPosition <= book.Worksheets.Count Then
      Set ws = book.Worksheets.Add(book.Worksheets(atPosition))
    Else
      Set ws = book.Worksheets.Add()
    End If
    ws.Name = desiredTabName
    Set comp = vbproj.VBComponents(ws.CodeName)
    On Error Resume Next
    comp.Name = desiredCodeName ' renames CodeName
    Err.Clear
    On Error GoTo 0
  End If
  Set EnsureSheet = ws
End Function

Sub RemoveComponentIfExists(vbproj, compName)
  Dim c
  On Error Resume Next
  Set c = vbproj.VBComponents(compName)
  If Not c Is Nothing Then
    If c.Type = 1 Or c.Type = 2 Or c.Type = 3 Then
      vbproj.VBComponents.Remove c
    End If
  End If
  Err.Clear: On Error GoTo 0
End Sub

Sub ReplaceDocFromFile(vbproj, docName, filePath)
  Dim comp, cm, total, code
  On Error Resume Next
  Set comp = vbproj.VBComponents(docName)
  If comp Is Nothing Then Exit Sub
  Set cm = comp.CodeModule
  total = cm.CountOfLines
  If total > 0 Then cm.DeleteLines 1, total
  code = ReadTextSansAttributes(filePath)
  If Len(code) > 0 Then cm.AddFromString code
  Err.Clear: On Error GoTo 0
End Sub

Sub ImportModuleIdempotent(vbproj, filePath)
  Dim base: base = fso.GetBaseName(filePath)
  RemoveComponentIfExists vbproj, base
  vbproj.VBComponents.Import filePath
End Sub

' === Pass 0: ensure any Sheet_<CodeName>.cls implied sheets exist ===
If fso.FolderExists(vbaFolder) Then
  Dim folder, file, base, ext, codeName
  Set folder = fso.GetFolder(vbaFolder)
  For Each file In folder.Files
    ext = LCase(fso.GetExtensionName(file.Path))
    base = fso.GetBaseName(file.Path)
    If ext = "cls" Then
      If LCase(Left(base, 6)) = "sheet_" Then
        codeName = Mid(base, 7)
        If Len(codeName) > 0 Then
          Call EnsureSheet(wb, codeName, codeName, 1)
        End If
      End If
    End If
  Next
End If

' === Pass 1: import/replace modules ===
If fso.FolderExists(vbaFolder) Then
  Dim folder2, file2, base2, ext2, cn
  Set folder2 = fso.GetFolder(vbaFolder)
  For Each file2 In folder2.Files
    ext2 = LCase(fso.GetExtensionName(file2.Path))
    base2 = fso.GetBaseName(file2.Path)

    If ext2 = "bas" Then
      ImportModuleIdempotent vbproj, file2.Path

    ElseIf ext2 = "cls" Then
      If LCase(base2) = "thisworkbook" Then
        ReplaceDocFromFile vbproj, "ThisWorkbook", file2.Path
      ElseIf LCase(Left(base2, 6)) = "sheet_" Then
        cn = Mid(base2, 7)
        If Len(cn) > 0 Then ReplaceDocFromFile vbproj, cn, file2.Path
      Else
        ImportModuleIdempotent vbproj, file2.Path
      End If

    ElseIf ext2 = "frm" Then
      ImportModuleIdempotent vbproj, file2.Path
    End If
  Next
End If

' Attempt known bootstrap macros (ignore missing)
Dim attempts, i
attempts = Array( _
  "'" & wb.Name & "'!Setup.RunOnce", _
  "'" & wb.Name & "'!TF_Data.EnsureStructure", _
  "'" & wb.Name & "'!TF_Data.SetupWorkbook", _
  "'" & wb.Name & "'!TF_Utils.EnsureStructure", _
  "'" & wb.Name & "'!EnsureStructure" _
)

For i = LBound(attempts) To UBound(attempts)
  On Error Resume Next
  xl.Run attempts(i)
  If Err.Number = 0 Then Exit For
  Err.Clear
  On Error GoTo 0
Next

wb.Save
wb.Close False
xl.Quit

WScript.Echo "[OK] Workbook ready: " & wbPath
WScript.Quit 0
