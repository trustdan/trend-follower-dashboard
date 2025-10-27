VERSION 1.0 CLASS
BEGIN
  MultiUse = -1  'True
END
Attribute VB_Name = "ThisWorkbook"
Attribute VB_GlobalNameSpace = False
Attribute VB_Creatable = False
Attribute VB_PredeclaredId = True
Attribute VB_Exposed = True

' ========================================
' ThisWorkbook Class Module
' Workbook-level event handlers
' ========================================

Private Sub Workbook_Open()
    ' Runs when workbook is opened
    On Error Resume Next

    ' Initialize logger first
    Call TF_Logger.InitializeLogger
    Call TF_Logger.WriteLog("========== WORKBOOK OPENED ==========")

    ' Check if setup has been run (Setup sheet exists = setup complete)
    If Not SheetExists("Setup") Then
        ' First time opening - run full setup
        Call TF_Logger.WriteLog("First time setup detected - running RunInitialSetup")

        MsgBox "Welcome to the Trading System!" & vbCrLf & vbCrLf & _
               "This appears to be your first time opening this workbook." & vbCrLf & _
               "Running initial setup now..." & vbCrLf & vbCrLf & _
               "Debug log will be saved to: TradingSystem_Debug.log", vbInformation, "First Time Setup"

        Call Setup.RunInitialSetup

        ' Open user guide after first setup
        Call OpenUserGuide
    Else
        ' Already set up - just activate TradeEntry
        Call TF_Logger.WriteLog("Setup sheet exists - workbook already configured")

        If SheetExists("TradeEntry") Then
            Call TF_Logger.WriteLog("Activating TradeEntry sheet")
            Worksheets("TradeEntry").Activate
            Call TF_UI.BindControls
        Else
            ' TradeEntry missing - show Setup sheet
            Call TF_Logger.WriteLog("TradeEntry sheet missing - showing Setup")
            Worksheets("Setup").Activate
        End If

        ' Optionally open user guide on every start (comment out if annoying)
        ' Call OpenUserGuide
    End If

    Call TF_Logger.WriteLog("Workbook_Open complete")

    On Error GoTo 0
End Sub

Private Function SheetExists(sheetName As String) As Boolean
    ' Helper function to check if sheet exists
    Dim ws As Worksheet
    On Error Resume Next
    Set ws = ThisWorkbook.Worksheets(sheetName)
    SheetExists = Not ws Is Nothing
    On Error GoTo 0
End Function

Private Sub Workbook_BeforeClose(Cancel As Boolean)
    ' Runs before workbook closes
    ' Could add cleanup tasks here if needed
End Sub

Private Sub Workbook_SheetActivate(ByVal Sh As Object)
    ' Runs when any sheet is activated
    If Sh.Name = "TradeEntry" Then
        ' Refresh bindings when TradeEntry is activated
        On Error Resume Next
        Call TF_UI.BindControls
        On Error GoTo 0
    End If
End Sub

Private Sub OpenUserGuide()
    ' Opens the USER_GUIDE.md file in default markdown viewer
    Dim guidePath As String
    Dim result As Long

    On Error Resume Next

    ' Build path to USER_GUIDE.md (in same folder as workbook)
    guidePath = ThisWorkbook.Path & "\USER_GUIDE.md"

    ' Check if file exists
    If Dir(guidePath) <> "" Then
        ' Open with default application (Shell execute)
        result = Shell("cmd /c start """" """ & guidePath & """", vbHide)

        If Err.Number <> 0 Then
            ' Fallback: Try opening with notepad
            result = Shell("notepad.exe """ & guidePath & """", vbNormalFocus)
        End If
    Else
        ' File not found - show message
        MsgBox "USER_GUIDE.md not found." & vbCrLf & vbCrLf & _
               "Expected location:" & vbCrLf & guidePath & vbCrLf & vbCrLf & _
               "You can find documentation in the project folder.", _
               vbInformation, "Guide Not Found"
    End If

    On Error GoTo 0
End Sub
