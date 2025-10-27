Attribute VB_Name = "Python_Run"
Option Explicit

Public Sub FetchScreenedAndRefresh()
    Dim bat As String
    bat = ThisWorkbook.Path & "\scripts\refresh_data.bat"

    Dim sh As Object, rc As Long
    Set sh = CreateObject("WScript.Shell")

    If Dir(bat) = "" Then
        MsgBox "Missing: " & bat, vbCritical: Exit Sub
    End If

    rc = sh.Run(Chr(34) & bat & Chr(34), 1, True)
    If rc <> 0 Then
        MsgBox "Data refresh failed. Exit code " & rc, vbCritical
        Exit Sub
    End If

    ' Refresh the PQ query/table
    ThisWorkbook.Connections("Query - PQ_Screened").Refresh
    MsgBox "Screened tickers updated.", vbInformation
End Sub
