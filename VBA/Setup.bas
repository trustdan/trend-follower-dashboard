Attribute VB_Name = "Setup"
Option Explicit

Public Sub RunOnce(Optional ByVal wb As Workbook)
    On Error GoTo CleanFail
    Application.ScreenUpdating = False
    Application.EnableEvents = False
    If wb Is Nothing Then Set wb = ThisWorkbook
    EnsureBasicSheets wb

    On Error Resume Next
    Application.Run "TF_UI.InitializeUI"
    Application.Run "TF_UI.ApplyTheme"
    On Error GoTo 0

CleanExit:
    Application.EnableEvents = True
    Application.ScreenUpdating = True
    Exit Sub
CleanFail:
    Resume CleanExit
End Sub

Public Sub EnsureBasicSheets(ByVal wb As Workbook)
    AddSheetIfMissing wb, "TradeEntry", 1
    AddSheetIfMissing wb, "Decisions"
    AddSheetIfMissing wb, "Positions"
    AddSheetIfMissing wb, "Contracts"
    AddSheetIfMissing wb, "Screener"
    AddSheetIfMissing wb, "Config"
    AddSheetIfMissing wb, "Data"

    With wb.Worksheets("TradeEntry")
        If Len(.Range("A1").Value) = 0 Then
            .Range("A1:F1").Value = Array("Ticker", "Date", "Entry Type", "Qty", "Entry Price", "Notes")
            .Rows(1).Font.Bold = True
        End If
    End With
End Sub

Private Sub AddSheetIfMissing(ByVal wb As Workbook, ByVal sName As String, Optional ByVal position As Long = 0)
    Dim ws As Worksheet
    On Error Resume Next
    Set ws = wb.Worksheets(sName)
    On Error GoTo 0
    If ws Is Nothing Then
        If position > 0 And position <= wb.Worksheets.Count Then
            Set ws = wb.Worksheets.Add(Before:=wb.Worksheets(position))
        Else
            Set ws = wb.Worksheets.Add(After:=wb.Worksheets(wb.Worksheets.Count))
        End If
        On Error Resume Next
        ws.Name = sName
        On Error GoTo 0
    End If
End Sub
