Attribute VB_Name = "TF_Utils"
' ============================================================================
' Module: TF_Utils
' Purpose: Helper functions for sheet/table/name management
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Function: SheetExists
' Returns: True if worksheet with given name exists
' ----------------------------------------------------------------------------
Function SheetExists(sheetName As String) As Boolean
    Dim ws As Worksheet
    On Error Resume Next
    Set ws = ThisWorkbook.Worksheets(sheetName)
    SheetExists = Not ws Is Nothing
    On Error GoTo 0
End Function

' ----------------------------------------------------------------------------
' Function: GetOrCreateSheet
' Returns: Worksheet with given name (creates if missing)
' ----------------------------------------------------------------------------
Function GetOrCreateSheet(sheetName As String) As Worksheet
    If SheetExists(sheetName) Then
        Set GetOrCreateSheet = ThisWorkbook.Worksheets(sheetName)
    Else
        Set GetOrCreateSheet = ThisWorkbook.Worksheets.Add(After:=ThisWorkbook.Worksheets(ThisWorkbook.Worksheets.Count))
        GetOrCreateSheet.Name = sheetName
    End If
End Function

' ----------------------------------------------------------------------------
' Function: GetOrCreateTable
' Returns: ListObject (Excel table) with given name and headers
' Creates table if missing, updates headers if exists
' ----------------------------------------------------------------------------
Function GetOrCreateTable(ws As Worksheet, tableName As String, headers As Variant) As ListObject
    Dim tbl As ListObject
    Dim i As Integer
    Dim headerRange As Range

    On Error Resume Next
    Set tbl = ws.ListObjects(tableName)
    On Error GoTo 0

    If tbl Is Nothing Then
        ' Create new table
        ws.Cells.Clear

        ' Write headers
        For i = LBound(headers) To UBound(headers)
            ws.Cells(1, i + 1).Value = headers(i)
        Next i

        ' Create table (single row with headers)
        Set headerRange = ws.Range(ws.Cells(1, 1), ws.Cells(1, UBound(headers) + 1))
        Set tbl = ws.ListObjects.Add(xlSrcRange, headerRange, , xlYes)
        tbl.Name = tableName
        tbl.TableStyle = "TableStyleMedium2"
    Else
        ' Verify headers match (optional - update if needed)
        For i = LBound(headers) To UBound(headers)
            If tbl.HeaderRowRange.Cells(1, i + 1).Value <> headers(i) Then
                tbl.HeaderRowRange.Cells(1, i + 1).Value = headers(i)
            End If
        Next i
    End If

    Set GetOrCreateTable = tbl
End Function

' ----------------------------------------------------------------------------
' Function: EnsureName
' Creates or updates a named range with default value
' ----------------------------------------------------------------------------
Sub EnsureName(rangeName As String, refersTo As String, defaultValue As Variant)
    Dim nm As Name
    Dim rng As Range

    ' Delete existing name if present
    On Error Resume Next
    ThisWorkbook.Names(rangeName).Delete
    On Error GoTo 0

    ' Create new name
    ThisWorkbook.Names.Add Name:=rangeName, RefersTo:="=" & refersTo

    ' Set default value if cell is empty
    Set rng = Range(rangeName)
    If IsEmpty(rng.Value) Or rng.Value = 0 Then
        rng.Value = defaultValue
    End If
End Sub

' ----------------------------------------------------------------------------
' Function: NzD (Null-to-Zero for Doubles)
' Returns: Value or default if null/empty
' ----------------------------------------------------------------------------
Function NzD(v As Variant, defaultVal As Double) As Double
    If IsEmpty(v) Or IsNull(v) Or v = "" Then
        NzD = defaultVal
    Else
        NzD = CDbl(v)
    End If
End Function

' ----------------------------------------------------------------------------
' Function: NzS (Null-to-String)
' Returns: Value or default if null/empty
' ----------------------------------------------------------------------------
Function NzS(v As Variant, defaultVal As String) As String
    If IsEmpty(v) Or IsNull(v) Or v = "" Then
        NzS = defaultVal
    Else
        NzS = CStr(v)
    End If
End Function

' ----------------------------------------------------------------------------
' Function: NormalizeTicker
' Returns: Cleaned ticker symbol (uppercase, trimmed, max 5 chars)
' ----------------------------------------------------------------------------
Function NormalizeTicker(rawTicker As String) As String
    Dim cleaned As String
    cleaned = Trim(UCase(rawTicker))
    cleaned = Replace(cleaned, ".", "-")  ' Convert periods to hyphens

    ' Limit to 5 characters
    If Len(cleaned) > 5 Then
        cleaned = Left(cleaned, 5)
    End If

    NormalizeTicker = cleaned
End Function
