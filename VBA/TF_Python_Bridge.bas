Attribute VB_Name = "TF_Python_Bridge"
Option Explicit

' ========================================
' TF_Python_Bridge Module
' VBA-Python integration for Excel Python (=PY() formulas)
' ========================================

Function IsPythonAvailable() As Boolean
    ' Check if Python in Excel is available (modern Python in Excel)
    ' Returns True if Python can be used

    Call TF_Logger.WriteLogSection("IsPythonAvailable() - Start")

    On Error Resume Next

    Dim testCell As Range
    Dim testFormula As String
    Dim testValue As Variant
    Dim excelVersion As String
    Dim formulaProperty As String

    excelVersion = Application.Version
    Call TF_Logger.WriteLog("Excel Version: " & excelVersion)

    ' Check if Control sheet exists
    On Error Resume Next
    Set testCell = Worksheets("Control").Range("Z1")
    If Err.Number <> 0 Then
        Call TF_Logger.WriteLogError("IsPythonAvailable", Err.Number, "Control sheet not found: " & Err.Description)
        IsPythonAvailable = False
        On Error GoTo 0
        Exit Function
    End If
    On Error GoTo 0

    Call TF_Logger.WriteLog("Test cell: " & testCell.Address & " on sheet: " & testCell.Worksheet.Name)

    ' Clear any existing content
    testCell.ClearContents
    Call TF_Logger.WriteLog("Cleared test cell")

    On Error Resume Next
    Err.Clear

    ' Modern Python in Excel uses =PY() function
    ' Try multiple approaches to detect Python

    ' Approach 1: Try Formula2 property (Office 365 modern)
    testFormula = "=PY(1+1)"
    Call TF_Logger.WriteLog("Attempting Formula2 with: " & testFormula)

    testCell.Formula2 = testFormula

    If Err.Number <> 0 Then
        Call TF_Logger.WriteLog("Formula2 assignment error: [" & Err.Number & "] " & Err.Description)

        ' Approach 2: Try regular Formula property
        Err.Clear
        Call TF_Logger.WriteLog("Trying regular Formula property")
        testCell.Formula = testFormula

        If Err.Number <> 0 Then
            Call TF_Logger.WriteLogError("IsPythonAvailable", Err.Number, "Both Formula2 and Formula failed: " & Err.Description)
            testCell.ClearContents
            IsPythonAvailable = False
            On Error GoTo 0
            Exit Function
        Else
            formulaProperty = "Formula"
            Call TF_Logger.WriteLog("Formula property succeeded")
        End If
    Else
        formulaProperty = "Formula2"
        Call TF_Logger.WriteLog("Formula2 property succeeded")
    End If

    ' Check what formula was actually set
    Dim actualFormula As String
    actualFormula = testCell.Formula
    Call TF_Logger.WriteLog("Cell formula after assignment: " & actualFormula)

    ' Force calculation
    Call TF_Logger.WriteLog("Forcing calculation...")
    Application.Calculate
    DoEvents

    ' Wait for cloud Python execution (modern Python in Excel is cloud-based)
    Call TF_Logger.WriteLog("Waiting 2 seconds for cloud execution...")
    Application.Wait Now + TimeValue("00:00:02")

    ' Check the cell value
    testValue = testCell.Value

    Call TF_Logger.WriteLog("Cell value type: " & TypeName(testValue))
    Call TF_Logger.WriteLog("Cell value: " & CStr(testValue))

    ' Check if it's an error value
    If IsError(testValue) Then
        Call TF_Logger.WriteLog("Cell contains error: " & CStr(testValue))

        ' Try to get error details
        On Error Resume Next
        Dim errType As Long
        errType = CVErr(testValue)
        Call TF_Logger.WriteLog("Error type: " & errType)
        On Error GoTo 0

        IsPythonAvailable = False
    ElseIf IsEmpty(testValue) Then
        Call TF_Logger.WriteLog("Cell is empty - Python may still be calculating")
        IsPythonAvailable = False
    ElseIf testValue = 2 Then
        Call TF_Logger.WriteLog("SUCCESS: Python returned correct value (2)")
        IsPythonAvailable = True
    Else
        Call TF_Logger.WriteLog("UNEXPECTED: Python returned: " & testValue & " (expected 2)")
        IsPythonAvailable = False
    End If

    ' Additional diagnostics
    Call TF_Logger.WriteLog("Cell HasFormula: " & testCell.HasFormula)
    If testCell.HasFormula Then
        Call TF_Logger.WriteLog("Cell Formula: " & testCell.Formula)
    End If

    ' Check if PY function exists
    On Error Resume Next
    Err.Clear
    Dim pyFunctionExists As Boolean
    pyFunctionExists = Application.WorksheetFunction.IsAvailable("PY")
    If Err.Number = 0 Then
        Call TF_Logger.WriteLog("PY function availability check: " & pyFunctionExists)
    Else
        Call TF_Logger.WriteLog("Cannot check PY function availability: " & Err.Description)
    End If

    ' Clean up
    testCell.ClearContents
    Call TF_Logger.WriteLog("Test cell cleared")

    Call TF_Logger.WriteLog("IsPythonAvailable() - Result: " & IsPythonAvailable)

    On Error GoTo 0
End Function

Function CallPythonFinvizScraper(queryString As String) As Variant
    ' Calls Python FINVIZ scraper and returns array of tickers
    ' Returns Empty array if Python not available or error occurs

    If Not IsPythonAvailable() Then
        CallPythonFinvizScraper = Array()
        Exit Function
    End If

    Dim pyCell As Range
    Dim pyFormula As String
    Dim result As Variant
    Dim i As Integer

    Set pyCell = Worksheets("Control").Range("Z2")

    ' Build Python formula for modern Python in Excel
    ' Use multi-line Python code approach
    pyFormula = "=PY(" & _
                "import sys;" & _
                "sys.path.append(r'../Python');" & _
                "from finviz_scraper import fetch_finviz_tickers;" & _
                "fetch_finviz_tickers('" & queryString & "')" & _
                ")"

    On Error Resume Next
    pyCell.Formula2 = pyFormula
    Application.Calculate
    DoEvents

    ' Give Python time to execute
    Application.Wait Now + TimeValue("00:00:10")

    result = pyCell.Value
    pyCell.ClearContents

    If Err.Number <> 0 Or IsError(result) Then
        CallPythonFinvizScraper = Array()
    Else
        CallPythonFinvizScraper = result
    End If

    On Error GoTo 0
End Function

Function CallPythonHeatCheck(addR As Double, bucket As String) As Variant
    ' Calls Python heat calculator and returns validation dict
    ' Returns Nothing if Python not available

    If Not IsPythonAvailable() Then
        Set CallPythonHeatCheck = Nothing
        Exit Function
    End If

    Dim pyCell As Range
    Dim pyFormula As String

    Set pyCell = Worksheets("Control").Range("Z3")

    ' Build Python formula to call heat_calculator for modern Python in Excel
    pyFormula = "=PY(" & _
                "import sys;" & _
                "sys.path.append(r'../Python');" & _
                "import pandas as pd;" & _
                "from heat_calculator import check_heat_caps;" & _
                "positions = xl('Positions[#All]').to_pandas();" & _
                "check_heat_caps(positions, " & addR & ", '" & bucket & "', " & _
                Range("Equity_E").Value & ", " & _
                Range("HeatCap_H_pct").Value & ", " & _
                Range("BucketHeatCap_pct").Value & ")" & _
                ")"

    On Error Resume Next
    pyCell.Formula2 = pyFormula
    Application.Calculate
    DoEvents

    Application.Wait Now + TimeValue("00:00:05")

    CallPythonHeatCheck = pyCell.Value
    pyCell.ClearContents

    On Error GoTo 0
End Function

Sub TestPythonIntegration()
    ' Comprehensive test of Python integration
    Dim result As String

    result = "Python Integration Test" & vbCrLf & String(50, "=") & vbCrLf & vbCrLf

    ' Test 1: Python availability
    result = result & "1. Python Availability: "
    If IsPythonAvailable() Then
        result = result & "[OK] AVAILABLE" & vbCrLf
    Else
        result = result & "[X] NOT AVAILABLE" & vbCrLf
        result = result & vbCrLf & "Python in Excel is not enabled." & vbCrLf
        result = result & "Enable it in: Data -> Python in Excel" & vbCrLf
        result = result & "(Requires Microsoft 365 Insider)"
        MsgBox result, vbInformation, "Python Test Results"
        Exit Sub
    End If

    ' Test 2: FINVIZ scraper
    result = result & vbCrLf & "2. FINVIZ Scraper Test: "
    Dim tickers As Variant
    tickers = CallPythonFinvizScraper("v=211&f=ta_highlow52w_nh&ft=4")

    If IsArray(tickers) And UBound(tickers) >= 0 Then
        result = result & "[OK] SUCCESS" & vbCrLf
        result = result & "   Found " & UBound(tickers) + 1 & " tickers" & vbCrLf

        ' Show first 5 tickers
        Dim i As Integer
        Dim sampleTickers As String
        sampleTickers = ""
        For i = 0 To WorksheetFunction.Min(4, UBound(tickers))
            If i > 0 Then sampleTickers = sampleTickers & ", "
            sampleTickers = sampleTickers & tickers(i)
        Next i
        result = result & "   Sample: " & sampleTickers & "..."
    Else
        result = result & "[X] FAILED" & vbCrLf
        result = result & "   No tickers returned"
    End If

    ' Test 3: Heat calculator (requires some position data)
    result = result & vbCrLf & vbCrLf & "3. Heat Calculator Test: "
    If Worksheets("Positions").ListObjects("tblPositions").ListRows.Count > 0 Then
        Dim heatResult As Variant
        heatResult = CallPythonHeatCheck(75, "Tech/Comm")

        If Not IsEmpty(heatResult) Then
            result = result & "[OK] SUCCESS" & vbCrLf
            result = result & "   Heat calculation completed"
        Else
            result = result & "[X] FAILED"
        End If
    Else
        result = result & "[SKIP] SKIPPED (no positions to test)"
    End If

    result = result & vbCrLf & vbCrLf & String(50, "=")
    result = result & vbCrLf & "All tests completed!"

    MsgBox result, vbInformation, "Python Integration Test Results"
End Sub

Function GetPythonVersion() As String
    ' Returns Python version string
    If Not IsPythonAvailable() Then
        GetPythonVersion = "Python not available"
        Exit Function
    End If

    Dim pyCell As Range
    Set pyCell = Worksheets("Control").Range("Z4")

    On Error Resume Next
    pyCell.Formula2 = "=PY(import sys; sys.version)"
    Application.Calculate
    Application.Wait Now + TimeValue("00:00:02")

    GetPythonVersion = pyCell.Value
    pyCell.ClearContents
    On Error GoTo 0
End Function

Sub InstallPythonPackages()
    ' Helper to show instructions for installing Python packages
    Dim msg As String

    msg = "Python Package Installation" & vbCrLf & String(50, "=") & vbCrLf & vbCrLf
    msg = msg & "Excel Python runs in the cloud and packages are pre-installed." & vbCrLf & vbCrLf
    msg = msg & "Required packages (should already be available):" & vbCrLf
    msg = msg & "  • pandas" & vbCrLf
    msg = msg & "  • numpy" & vbCrLf & vbCrLf
    msg = msg & "For local Python development:" & vbCrLf
    msg = msg & "  1. Open Command Prompt" & vbCrLf
    msg = msg & "  2. Navigate to: " & ThisWorkbook.Path & "\Python" & vbCrLf
    msg = msg & "  3. Run: pip install -r requirements.txt"

    MsgBox msg, vbInformation, "Python Packages"
End Sub
