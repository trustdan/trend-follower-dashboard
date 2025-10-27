Option Explicit

' ===== Logging =====
Dim fso:   Set fso   = CreateObject("Scripting.FileSystemObject")
Dim shell: Set shell = CreateObject("WScript.Shell")
Dim logTS

Function NowISO()
  Dim d: d = Now
  NowISO = Year(d) & "-" & Right("0"&Month(d),2) & "-" & Right("0"&Day(d),2) & _
           " " & Right("0"&Hour(d),2) & ":" & Right("0"&Minute(d),2) & ":" & Right("0"&Second(d),2)
End Function
Sub Log(msg)
  Dim line: line = NowISO() & " " & msg
  ' Write to log file if available
  On Error Resume Next
  If Not logTS Is Nothing Then logTS.WriteLine line
  On Error GoTo 0
  ' Always write to console
  WScript.Echo line
End Sub
Sub LogErr(prefix)
  If Err.Number <> 0 Then
    Log "[ERR] " & prefix & " (code=" & Err.Number & "): " & Err.Description
    Err.Clear
  End If
End Sub

' ===== Inputs =====
Dim wbPath, vbaFolder, logPath
If WScript.Arguments.Count >= 1 Then wbPath    = WScript.Arguments(0) Else wbPath    = fso.GetAbsolutePathName("TrendFollowing_TradeEntry.xlsm")
If WScript.Arguments.Count >= 2 Then vbaFolder = WScript.Arguments(1) Else vbaFolder = "VBA"
If WScript.Arguments.Count >= 3 Then logPath   = WScript.Arguments(2) Else logPath   = fso.BuildPath(".", "logs\build.log")

' Open log (append)
On Error Resume Next
Set logTS = fso.OpenTextFile(logPath, 8, True)
If Err.Number <> 0 Then
  Set logTS = fso.CreateTextFile(logPath, True)
  Err.Clear
End If
On Error GoTo 0
Log "[STEP] Builder start"
Log "[INFO] wbPath=" & wbPath
Log "[INFO] vbaFolder=" & vbaFolder

' ===== Excel startup (hidden, events OFF) =====
Dim xl, wb, vbproj, rc
rc = 0

On Error Resume Next
Set xl = CreateObject("Excel.Application")
Log "[STEP] Excel.Application created"
If Err.Number <> 0 Or xl Is Nothing Then
  LogErr "CreateObject(Excel.Application)"
  rc = 10
End If

If rc = 0 Then
  xl.DisplayAlerts = False
  xl.Visible = False
  xl.EnableEvents = False   ' prevent Workbook_Open prompts while hidden

  If fso.FileExists(wbPath) Then
    Set wb = xl.Workbooks.Open(wbPath)
    Log "[STEP] Opened existing workbook"
  Else
    Set wb = xl.Workbooks.Add()
    Log "[STEP] Added new workbook"
    wb.SaveAs wbPath, 52    ' 52 = .xlsm
    Log "[STEP] Saved as .xlsm"
  End If
  Log "[INFO] Excel.Version=" & xl.Version

  Set vbproj = wb.VBProject
  If vbproj Is Nothing Then
    Log "[ERR] Access VBProject failed. Enable Excel → Options → Trust Center → Trust Center Settings → Macro Settings → 'Trust access to the VBA project object model'."
    rc = 11
  End If
End If
On Error GoTo 0

' ===== Helpers =====
Function ReadTextSansAttributes(p)
  Dim s, t, ts, arr, i, line
  Set ts = fso.OpenTextFile(p, 1, False)
  s = ts.ReadAll: ts.Close
  s = Replace(s, vbCrLf, vbLf)
  s = Replace(s, vbCr, vbLf)
  arr = Split(s, vbLf)
  For i = 0 To UBound(arr)
    line = Trim(arr(i))
    If LCase(Left(line, 7))  = "version" Then
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
    On Error Resume Next
    Set comp = vbproj.VBComponents(ws.CodeName)
    comp.Name = desiredCodeName         ' rename document component -> sets CodeName
    Log "[INFO] Created sheet '" & desiredTabName & "' as CodeName '" & desiredCodeName & "'"
    LogErr "Rename sheet CodeName"
    On Error GoTo 0
  Else
    Log "[INFO] Sheet with CodeName '" & desiredCodeName & "' already exists"
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
      Log "[INFO] Removed existing component: " & compName
    End If
  End If
  Err.Clear
End Sub

Sub ReplaceDocFromFile(vbproj, docName, filePath)
  Dim comp, cm, total, code
  On Error Resume Next
  Set comp = vbproj.VBComponents(docName)
  If comp Is Nothing Then
    Log "[WARN] Document component '" & docName & "' not found"
    Exit Sub
  End If
  Set cm = comp.CodeModule
  total = cm.CountOfLines
  If total > 0 Then cm.DeleteLines 1, total
  code = ReadTextSansAttributes(filePath)
  If Len(code) > 0 Then cm.AddFromString code
  Log "[INFO] Replaced doc module '" & docName & "' from " & fso.GetFileName(filePath)
  LogErr "ReplaceDocFromFile(" & docName & ")"
End Sub

Sub ImportModuleIdempotent(vbproj, filePath)
  Dim base: base = fso.GetBaseName(filePath)
  RemoveComponentIfExists vbproj, base
  vbproj.VBComponents.Import filePath
  Log "[INFO] Imported: " & fso.GetFileName(filePath)
  LogErr "Import " & base
End Sub

' ===== Pass 0: ensure any Sheet_<CodeName>.cls sheets =====
If rc = 0 And fso.FolderExists(vbaFolder) Then
  Dim folder, file, base, ext, codeName
  Set folder = fso.GetFolder(vbaFolder)
  For Each file In folder.Files
    ext = LCase(fso.GetExtensionName(file.Path))
    base = fso.GetBaseName(file.Path)
    If ext = "cls" And LCase(Left(base, 6)) = "sheet_" Then
      codeName = Mid(base, 7)
      If Len(codeName) > 0 Then Call EnsureSheet(wb, codeName, codeName, 1)
    End If
  Next
Else
  If rc = 0 Then Log "[WARN] VBA folder not found: " & vbaFolder
End If

' ===== Pass 1: import/replace modules =====
If rc = 0 And fso.FolderExists(vbaFolder) Then
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

' List components after import
If rc = 0 Then
  On Error Resume Next
  Dim comp
  Log "[STEP] VBComponents after import:"
  For Each comp In vbproj.VBComponents
    Log "  - " & comp.Name & " (Type=" & comp.Type & ")"
  Next
  On Error GoTo 0
End If

' ===== Run a non-interactive bootstrap (ignore missing) =====
If rc = 0 Then
  On Error Resume Next
  Dim attempts, i
  ' Use simple macro names without workbook prefix
  attempts = Array( _
    "Setup.RunOnce", _
    "TF_Data.EnsureStructure", _
    "TF_Data.SetupWorkbook", _
    "TF_Utils.EnsureStructure", _
    "EnsureStructure" _
  )
  For i = LBound(attempts) To UBound(attempts)
    Err.Clear
    xl.Run attempts(i)
    If Err.Number = 0 Then
      Log "[STEP] Ran bootstrap: " & attempts(i)
      Exit For
    Else
      Log "[INFO] Bootstrap not found/failed: " & attempts(i) & " (code=" & Err.Number & ")"
    End If
  Next
  On Error GoTo 0
End If

' ===== Save & close =====
' ALWAYS run cleanup, even if there were errors
On Error Resume Next

' Save and close workbook
If Not wb Is Nothing Then
  ' Try regular Save first
  wb.Save
  If Err.Number <> 0 Then
    Log "[WARN] wb.Save failed (" & Err.Number & "): " & Err.Description
    Err.Clear

    ' Delete existing file if it exists
    If fso.FileExists(wbPath) Then
      fso.DeleteFile wbPath, True
      Log "[INFO] Deleted old workbook file"
    End If

    ' Try SaveAs instead
    wb.SaveAs wbPath, 52
    If Err.Number = 0 Then
      Log "[STEP] Saved workbook (via SaveAs)"
    Else
      Log "[ERROR] SaveAs also failed (" & Err.Number & "): " & Err.Description
      Err.Clear
    End If
  Else
    Log "[STEP] Saved workbook"
  End If

  ' Close without saving (we already saved above)
  wb.Close False
  If Err.Number = 0 Then
    Log "[STEP] Closed workbook"
  Else
    Log "[WARN] Close failed (" & Err.Number & "): " & Err.Description
    Err.Clear
  End If
End If
Set wb = Nothing

' Quit Excel
If Not xl Is Nothing Then
  xl.DisplayAlerts = False
  xl.Quit
  If Err.Number = 0 Then
    Log "[STEP] Excel quit successfully"
  Else
    Log "[WARN] Excel quit failed: " & Err.Description
    Err.Clear
  End If
End If
Set xl = Nothing

' Wait a moment for Excel to fully close
WScript.Sleep 1000

On Error GoTo 0

Log "[OK] Build completed with rc=" & rc

' Close log file
On Error Resume Next
If Not logTS Is Nothing Then
  logTS.Close
  Set logTS = Nothing
End If
On Error GoTo 0

WScript.Quit rc
