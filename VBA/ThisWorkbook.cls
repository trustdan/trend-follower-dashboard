VERSION 1.0 CLASS
BEGIN
  MultiUse = -1  'True
END
Attribute VB_Name = "ThisWorkbook"
Attribute VB_GlobalNameSpace = False
Attribute VB_Creatable = False
Attribute VB_PredeclaredId = True
Attribute VB_Exposed = True
' ============================================================================
' Class Module: ThisWorkbook
' Purpose: Workbook-level event handlers
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Event: Workbook_Open
' Runs when workbook is opened - ensures structure exists
' ----------------------------------------------------------------------------
Private Sub Workbook_Open()
    On Error Resume Next

    ' Check if structure exists (look for key sheet)
    If Not SheetExists("TradeEntry") Then
        Dim response As VbMsgBoxResult
        response = MsgBox( _
            "This workbook has not been set up yet." & vbCrLf & vbCrLf & _
            "Would you like to run the setup now?" & vbCrLf & _
            "(This will create all required sheets and tables)", _
            vbYesNo + vbQuestion, "First-Time Setup")

        If response = vbYes Then
            Call EnsureStructure
        End If
    Else
        ' Structure exists - just activate TradeEntry
        Sheets("TradeEntry").Activate
    End If

    On Error GoTo 0
End Sub

' ----------------------------------------------------------------------------
' Event: Workbook_BeforeClose
' Optional: Remind to update cooldowns or backup
' ----------------------------------------------------------------------------
Private Sub Workbook_BeforeClose(Cancel As Boolean)
    ' Optional: Add any cleanup or reminder logic here
    ' Example: Prompt to update cooldowns if it's been a week
End Sub
