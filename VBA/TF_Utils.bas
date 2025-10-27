Attribute VB_Name = "TF_Utils"
Option Explicit

' ========================================
' TF_Utils Module
' Utility functions for sheet/table/name management
' ========================================

' Check if a worksheet exists
Function SheetExists(sheetName As String) As Boolean
    Dim ws As Worksheet
    On Error Resume Next
    Set ws = ThisWorkbook.Worksheets(sheetName)
    SheetExists = Not ws Is Nothing
    On Error GoTo 0
End Function

' Get existing sheet or create new one
Function GetOrCreateSheet(sheetName As String) As Worksheet
    Dim ws As Worksheet

    If SheetExists(sheetName) Then
        Set GetOrCreateSheet = ThisWorkbook.Worksheets(sheetName)
    Else
        Set ws = ThisWorkbook.Worksheets.Add(After:=ThisWorkbook.Worksheets(ThisWorkbook.Worksheets.Count))
        ws.Name = sheetName
        Set GetOrCreateSheet = ws
    End If
End Function

' Get existing table or create new one
Function GetOrCreateTable(ws As Worksheet, tableName As String, headers As Variant) As ListObject
    Dim tbl As ListObject
    Dim i As Integer
    Dim rng As Range

    ' Check if table already exists
    On Error Resume Next
    Set tbl = ws.ListObjects(tableName)
    On Error GoTo 0

    If Not tbl Is Nothing Then
        Set GetOrCreateTable = tbl
        Exit Function
    End If

    ' Create new table
    ws.Cells.Clear

    ' Add headers
    For i = LBound(headers) To UBound(headers)
        ws.Cells(1, i + 1).Value = headers(i)
    Next i

    ' Create table
    Set rng = ws.Range(ws.Cells(1, 1), ws.Cells(1, UBound(headers) - LBound(headers) + 1))
    Set tbl = ws.ListObjects.Add(xlSrcRange, rng, , xlYes)
    tbl.Name = tableName
    tbl.TableStyle = "TableStyleMedium2"

    Set GetOrCreateTable = tbl
End Function

' Ensure named range exists with default value
Sub EnsureName(rangeName As String, refersTo As String, defaultValue As Variant)
    Dim nm As Name
    Dim rng As Range

    ' Check if name exists
    On Error Resume Next
    Set nm = ThisWorkbook.Names(rangeName)
    On Error GoTo 0

    If nm Is Nothing Then
        ' Create new name
        ThisWorkbook.Names.Add Name:=rangeName, RefersTo:="=" & refersTo

        ' Set default value
        Set rng = Range(refersTo)
        If IsEmpty(rng.Value) Then
            rng.Value = defaultValue
        End If
    End If
End Sub

' Null-safe double value with default
Function NzD(value As Variant, defaultValue As Double) As Double
    If IsEmpty(value) Or IsNull(value) Or Not IsNumeric(value) Then
        NzD = defaultValue
    Else
        NzD = CDbl(value)
    End If
End Function

' Null-safe string value with default
Function NzS(value As Variant, defaultValue As String) As String
    If IsEmpty(value) Or IsNull(value) Then
        NzS = defaultValue
    Else
        NzS = CStr(value)
    End If
End Function

' Normalize ticker symbol (uppercase, trim, handle variations)
Function NormalizeTicker(rawTicker As String) As String
    Dim cleaned As String

    ' Basic cleanup
    cleaned = Trim(UCase(rawTicker))

    ' Replace dots with dashes (e.g., BRK.B -> BRK-B)
    cleaned = Replace(cleaned, ".", "-")

    ' Remove any other special characters except dashes
    Dim i As Integer
    Dim result As String
    result = ""

    For i = 1 To Len(cleaned)
        Dim char As String
        char = Mid(cleaned, i, 1)

        ' Keep letters, numbers, and dashes only
        If (char >= "A" And char <= "Z") Or (char >= "0" And char <= "9") Or char = "-" Then
            result = result & char
        End If
    Next i

    ' Validate length
    If Len(result) > 0 And Len(result) <= 5 Then
        NormalizeTicker = result
    Else
        NormalizeTicker = ""
    End If
End Function

' Check if a named range exists
Function NameExists(rangeName As String) As Boolean
    Dim nm As Name
    On Error Resume Next
    Set nm = ThisWorkbook.Names(rangeName)
    NameExists = Not nm Is Nothing
    On Error GoTo 0
End Function

' Safe cell value reader
Function GetCellValue(ws As Worksheet, cellAddress As String) As Variant
    On Error Resume Next
    GetCellValue = ws.Range(cellAddress).Value
    If Err.Number <> 0 Then
        GetCellValue = Empty
    End If
    On Error GoTo 0
End Function

' Safe cell value writer
Sub SetCellValue(ws As Worksheet, cellAddress As String, value As Variant)
    On Error Resume Next
    ws.Range(cellAddress).Value = value
    On Error GoTo 0
End Sub
