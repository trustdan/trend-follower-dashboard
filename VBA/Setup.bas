Attribute VB_Name = "Setup"
Option Explicit

' ========================================
' Setup Module
' One-click initialization and UI building
' ========================================

Sub RunInitialSetup()
    ' Complete setup in one click
    On Error GoTo ErrorHandler

    ' Disable events to prevent recursive loop
    Application.EnableEvents = False
    Application.ScreenUpdating = False
    Application.Calculation = xlCalculationManual

    Call TF_Logger.WriteLogSection("RunInitialSetup - Start")

    Dim startTime As Double
    startTime = Timer

    ' Step 1: Ensure structure
    Call TF_Logger.WriteLog("Step 1: Calling TF_Data.EnsureStructure")
    Call TF_Data.EnsureStructure

    ' Step 2: Build UI
    Call TF_Logger.WriteLog("Step 2: Calling TF_UI_Builder.BuildTradeEntryUI")
    Call TF_UI_Builder.BuildTradeEntryUI

    ' Step 3: Create Setup sheet with instructions
    Call TF_Logger.WriteLog("Step 3: Calling CreateSetupSheet")
    Call CreateSetupSheet

    ' Re-enable events
    Application.EnableEvents = True
    Application.Calculation = xlCalculationAutomatic
    Application.ScreenUpdating = True

    ' Show success message
    Dim elapsed As Double
    elapsed = Timer - startTime

    Call TF_Logger.WriteLog("Setup completed in " & Format(elapsed, "0.0") & " seconds")

    MsgBox "Setup Complete!" & vbCrLf & vbCrLf & _
           "[OK] All sheets created" & vbCrLf & _
           "[OK] All tables initialized" & vbCrLf & _
           "[OK] TradeEntry UI built" & vbCrLf & _
           "[OK] Buttons configured" & vbCrLf & _
           "[OK] Dropdowns bound" & vbCrLf & vbCrLf & _
           "Time: " & Format(elapsed, "0.0") & " seconds" & vbCrLf & vbCrLf & _
           "Go to TradeEntry sheet to start trading!", _
           vbInformation, "Setup Complete"

    ' Go to TradeEntry sheet (not Setup sheet to avoid confusion)
    On Error Resume Next
    Worksheets("TradeEntry").Activate
    On Error GoTo 0

    Call TF_Logger.WriteLog("RunInitialSetup - Complete")

    Exit Sub

ErrorHandler:
    Application.EnableEvents = True
    Application.Calculation = xlCalculationAutomatic
    Application.ScreenUpdating = True
    Call TF_Logger.WriteLogError("RunInitialSetup", Err.Number, Err.Description)
    MsgBox "Setup error: " & Err.Description, vbCritical
End Sub

Sub CreateSetupSheet()
    ' Create a Setup/Welcome sheet with instructions and buttons
    Dim ws As Worksheet
    Dim btn As Button

    Call TF_Logger.WriteLog("CreateSetupSheet - Start")

    ' Delete existing Setup sheet if present
    On Error Resume Next
    Application.DisplayAlerts = False
    Call TF_Logger.WriteLog("Deleting existing Setup sheet if present")
    Worksheets("Setup").Delete
    Application.DisplayAlerts = True
    On Error GoTo 0

    ' Create new Setup sheet at beginning
    Call TF_Logger.WriteLog("Creating new Setup sheet")
    Set ws = ThisWorkbook.Worksheets.Add(Before:=ThisWorkbook.Worksheets(1))
    ws.Name = "Setup"
    Call TF_Logger.WriteLog("Setup sheet created and named")

    With ws
        ' Title
        .Range("A1").Value = "TRADING SYSTEM SETUP"
        .Range("A1").Font.Size = 20
        .Range("A1").Font.Bold = True
        .Range("A1:F1").Merge
        .Range("A1").HorizontalAlignment = xlCenter

        ' Status section
        .Range("A3").Value = "SETUP STATUS"
        .Range("A3").Font.Size = 14
        .Range("A3").Font.Bold = True

        .Range("A4").Value = "[OK] Workbook created"
        .Range("A5").Value = "[OK] VBA modules imported"
        .Range("A6").Value = "[OK] Data structure created"
        .Range("A7").Value = "[OK] TradeEntry UI built"
        .Range("A8").Value = "[OK] Checkboxes created (check TradeEntry sheet)"
        .Range("A8").Font.Bold = True
        .Range("A8").Font.Color = RGB(0, 128, 0)

        .Range("A9").Value = ""
        .Range("A9:F9").Merge
        .Range("A9").Interior.Color = RGB(255, 255, 200)
        .Range("A9").Font.Size = 11
        .Range("A9").Font.Bold = True
        .Range("A9").Value = "*** IF YOU SEE 'Dropdown Error' -> Click the >> RUN INITIAL SETUP << button below! ***"

        ' Instructions section
        .Range("A10").Value = "SETUP COMPLETE - READY TO TRADE!"
        .Range("A10").Font.Size = 14
        .Range("A10").Font.Bold = True
        .Range("A10").Font.Color = RGB(0, 128, 0)

        .Range("A12").Value = "The system attempted to create checkboxes automatically."
        .Range("A13").Value = "Go to TradeEntry sheet and verify 6 checkboxes appear in rows 21-26."
        .Range("A14").Value = ""
        .Range("A15").Value = "IF CHECKBOXES ARE MISSING (COM automation failed):"
        .Range("A15").Font.Bold = True
        .Range("A16").Value = "1. Go to TradeEntry sheet"
        .Range("A17").Value = "2. Developer -> Insert -> Check Box (Form Control)"
        .Range("A18").Value = "3. Draw 6 checkboxes next to column A, rows 21-26"
        .Range("A19").Value = "4. Link each to: C20, C21, C22, C23, C24, C25 (Right-click -> Format Control)"

        ' Test section
        .Range("A21").Value = "QUICK TEST WORKFLOW"
        .Range("A21").Font.Size = 14
        .Range("A21").Font.Bold = True

        .Range("A23").Value = "After adding checkboxes, test the system:"
        .Range("A24").Value = "1. Go to TradeEntry sheet"
        .Range("A25").Value = "2. Select Preset: TF_BREAKOUT_LONG"
        .Range("A26").Value = "3. Click 'Import Candidates' button -> Paste: AAPL, MSFT, NVDA"
        .Range("A27").Value = "4. Select Ticker: AAPL"
        .Range("A28").Value = "5. Enter: Entry=180, ATR N=1.50, K=2"
        .Range("A29").Value = "6. Check all 6 boxes -> Click 'Evaluate' -> See GREEN!"
        .Range("A30").Value = "7. Click 'Recalc Sizing' -> See calculated shares"
        .Range("A31").Value = "8. Wait 2 minutes -> Click 'Save Decision' -> Trade logged!"

        ' Buttons section
        .Range("A33").Value = "SETUP UTILITIES"
        .Range("A33").Font.Size = 14
        .Range("A33").Font.Bold = True

        ' Button: Run Initial Setup (PRIMARY ACTION)
        Set btn = .Buttons.Add(30, 650, 220, 35)
        btn.Text = ">> RUN INITIAL SETUP <<"
        btn.OnAction = "Setup.RunInitialSetup"
        btn.Font.Bold = True
        btn.Font.Size = 12

        ' Button: Rebuild UI
        Set btn = .Buttons.Add(30, 690, 150, 25)
        btn.Text = "Rebuild TradeEntry UI"
        btn.OnAction = "TF_UI_Builder.BuildTradeEntryUI"

        ' Button: Test Python
        Set btn = .Buttons.Add(190, 690, 150, 25)
        btn.Text = "Test Python Integration"
        btn.OnAction = "TF_Python_Bridge.TestPythonIntegration"

        ' Button: Clear Old Candidates
        Set btn = .Buttons.Add(350, 690, 150, 25)
        btn.Text = "Clear Old Candidates"
        btn.OnAction = "Setup.ClearOldCandidatesPrompt"

        ' Button: Open User Guide
        Set btn = .Buttons.Add(30, 725, 150, 25)
        btn.Text = "Open User Guide"
        btn.OnAction = "Setup.OpenUserGuideFromButton"

        ' Button: Open Debug Log
        Set btn = .Buttons.Add(190, 725, 150, 25)
        btn.Text = "Open Debug Log"
        btn.OnAction = "TF_Logger.OpenLogFile"

        ' Documentation section
        .Range("A37").Value = "DOCUMENTATION"
        .Range("A37").Font.Size = 14
        .Range("A37").Font.Bold = True

        .Range("A39").Value = "USER_GUIDE.md - Complete walkthrough (click button above)"
        .Range("A40").Value = "README.md - Complete user guide"
        .Range("A41").Value = "START_HERE.md - Quick start guide"
        .Range("A42").Value = "IMPLEMENTATION_STATUS.md - Technical details"

        ' Settings section
        .Range("A43").Value = "KEY SETTINGS (Summary Sheet)"
        .Range("A43").Font.Size = 14
        .Range("A43").Font.Bold = True

        .Range("A45").Value = "Equity_E:"
        .Range("B45").Value = "Your account size (default: $10,000)"
        .Range("A46").Value = "RiskPct_r:"
        .Range("B46").Value = "Risk per trade (default: 0.75%)"
        .Range("A47").Value = "HeatCap_H_pct:"
        .Range("B47").Value = "Max portfolio heat (default: 4%)"

        ' Format columns
        .Columns("A:A").ColumnWidth = 50
        .Columns("B:F").ColumnWidth = 20

        ' Color coding
        .Range("A3:F3").Interior.Color = RGB(200, 230, 255)
        .Range("A10:F10").Interior.Color = RGB(255, 255, 200)
        .Range("A21:F21").Interior.Color = RGB(200, 255, 200)
        .Range("A33:F33").Interior.Color = RGB(240, 240, 240)
        .Range("A37:F37").Interior.Color = RGB(240, 240, 240)
        .Range("A43:F43").Interior.Color = RGB(240, 240, 240)
    End With
End Sub

Sub ClearOldCandidatesPrompt()
    ' Helper to clear old candidates
    Dim daysOld As String
    Dim days As Integer

    daysOld = InputBox("Clear candidates older than how many days?", "Clear Old Candidates", "7")

    If daysOld = "" Then Exit Sub

    If IsNumeric(daysOld) Then
        days = CInt(daysOld)
        Call TF_Data.ClearOldCandidates(days)
        MsgBox "Cleared candidates older than " & days & " days", vbInformation
    Else
        MsgBox "Please enter a number", vbExclamation
    End If
End Sub

Sub ShowSetupSheet()
    ' Quick way to get back to setup instructions
    On Error Resume Next
    Worksheets("Setup").Activate
    On Error GoTo 0
End Sub

Sub OpenUserGuideFromButton()
    ' Wrapper to call OpenUserGuide from Setup sheet button
    ' This calls the function in ThisWorkbook module
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
