Attribute VB_Name = "TF_Python_Bridge"
' ============================================================================
' Module: TF_Python_Bridge
' Purpose: VBA-Python integration layer for Excel
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Function: CallPythonFinvizScraper
' Calls Python FINVIZ scraper and returns array of tickers
' ----------------------------------------------------------------------------
Function CallPythonFinvizScraper(queryString As String) As Variant
    Dim pyFormula As String
    Dim resultCell As Range
    Dim result As Variant
    Dim i As Long
    Dim tickers() As String
    Dim tempSheet As Worksheet

    On Error GoTo ErrorHandler

    ' Use Control sheet for Python calculations (hidden from user)
    Set tempSheet = Sheets("Control")

    ' Build Python formula
    ' Syntax: =PY("module.function", arg1, arg2, ...)
    pyFormula = "=PY(""finviz_scraper.fetch_finviz_tickers"", """ & queryString & """)"

    ' Write formula to hidden cell
    Set resultCell = tempSheet.Range("Z10")
    resultCell.Formula = pyFormula

    ' Wait for calculation
    Application.Calculate
    DoEvents

    ' Read result (should be an array of tickers)
    result = resultCell.Value

    ' Handle different result types
    If IsArray(result) Then
        ' Python returned a list - convert to VBA array
        ReDim tickers(LBound(result) To UBound(result))
        For i = LBound(result) To UBound(result)
            tickers(i) = CStr(result(i))
        Next i
        CallPythonFinvizScraper = tickers

    ElseIf Not IsEmpty(result) And result <> "" Then
        ' Single ticker returned
        ReDim tickers(0 To 0)
        tickers(0) = CStr(result)
        CallPythonFinvizScraper = tickers

    Else
        ' Empty result
        CallPythonFinvizScraper = Array()
    End If

    Exit Function

ErrorHandler:
    MsgBox "Error calling Python scraper: " & Err.Description & vbCrLf & vbCrLf & _
           "Make sure Python in Excel is enabled and finviz_scraper.py is loaded.", _
           vbCritical, "Python Integration Error"
    CallPythonFinvizScraper = Array()
End Function

' ----------------------------------------------------------------------------
' Function: CallPythonHeatCheck
' Calls Python heat calculator and returns validation dict
' ----------------------------------------------------------------------------
Function CallPythonHeatCheck(addR As Double, bucket As String) As Variant
    Dim pyFormula As String
    Dim resultCell As Range
    Dim result As Variant
    Dim tempSheet As Worksheet
    Dim equity As Double, portCapPct As Double, bucketCapPct As Double

    On Error GoTo ErrorHandler

    ' Get settings from named ranges
    equity = Range("Equity_E").Value
    portCapPct = Range("HeatCap_H_pct").Value
    bucketCapPct = Range("BucketHeatCap_pct").Value

    ' Use Control sheet
    Set tempSheet = Sheets("Control")

    ' Build Python formula
    ' Pass Positions table as DataFrame, along with parameters
    pyFormula = "=PY(""heat_calculator.check_heat_caps"", " & _
                "xl(""Positions[#All]""), " & _
                addR & ", " & _
                """" & bucket & """, " & _
                equity & ", " & _
                portCapPct & ", " & _
                bucketCapPct & ")"

    ' Write formula to hidden cell
    Set resultCell = tempSheet.Range("Z20")
    resultCell.Formula = pyFormula

    ' Wait for calculation
    Application.Calculate
    DoEvents

    ' Read result (should be a dict/object)
    result = resultCell.Value

    ' Return the result (Excel will convert Python dict to VBA object)
    CallPythonHeatCheck = result

    Exit Function

ErrorHandler:
    MsgBox "Error calling Python heat calculator: " & Err.Description & vbCrLf & vbCrLf & _
           "Falling back to VBA heat calculations.", _
           vbExclamation, "Python Integration Warning"

    ' Return empty variant to signal fallback to VBA
    CallPythonHeatCheck = Empty
End Function

' ----------------------------------------------------------------------------
' Function: IsPythonAvailable
' Checks if Python in Excel is available and modules are loaded
' ----------------------------------------------------------------------------
Function IsPythonAvailable() As Boolean
    Dim testCell As Range
    Dim testFormula As String
    Dim result As Variant

    On Error Resume Next

    ' Try a simple Python formula
    Set testCell = Sheets("Control").Range("Z1")
    testFormula = "=PY(""1+1"")"

    testCell.Formula = testFormula
    Application.Calculate
    DoEvents

    result = testCell.Value

    ' Check if result is 2 (meaning Python worked)
    If result = 2 Then
        IsPythonAvailable = True
    Else
        IsPythonAvailable = False
    End If

    ' Clean up
    testCell.ClearContents

    On Error GoTo 0
End Function

' ----------------------------------------------------------------------------
' Function: LoadPythonModules
' Attempts to import required Python modules
' Returns: True if successful, False otherwise
' ----------------------------------------------------------------------------
Function LoadPythonModules() As Boolean
    Dim importCell As Range
    Dim importFormula As String
    Dim result As Variant

    On Error Resume Next

    Set importCell = Sheets("Control").Range("Z2")

    ' Try to import both modules
    importFormula = "=PY(""" & _
                   "import finviz_scraper; " & _
                   "import heat_calculator; " & _
                   "'Modules loaded'" & _
                   """)"

    importCell.Formula = importFormula
    Application.Calculate
    DoEvents

    result = importCell.Value

    If InStr(CStr(result), "loaded") > 0 Then
        LoadPythonModules = True
    Else
        LoadPythonModules = False
    End If

    On Error GoTo 0
End Function

' ----------------------------------------------------------------------------
' Sub: TestPythonIntegration
' Test procedure to verify Python integration is working
' ----------------------------------------------------------------------------
Sub TestPythonIntegration()
    Dim pythonOK As Boolean
    Dim modulesOK As Boolean
    Dim tickers As Variant
    Dim heatResult As Variant
    Dim i As Long

    Application.ScreenUpdating = False

    ' Test 1: Python availability
    pythonOK = IsPythonAvailable()

    If Not pythonOK Then
        MsgBox "Python in Excel is not available." & vbCrLf & vbCrLf & _
               "To enable:" & vbCrLf & _
               "1. Update to Microsoft 365 Insider" & vbCrLf & _
               "2. Data tab → Python → Enable", _
               vbExclamation, "Python Not Available"
        GoTo CleanExit
    End If

    ' Test 2: Module loading
    modulesOK = LoadPythonModules()

    If Not modulesOK Then
        MsgBox "Python modules not found." & vbCrLf & vbCrLf & _
               "Make sure finviz_scraper.py and heat_calculator.py are in:" & vbCrLf & _
               ThisWorkbook.Path & "\Python\", _
               vbExclamation, "Modules Not Loaded"
        GoTo CleanExit
    End If

    ' Test 3: FINVIZ scraper (use a simple query)
    Dim testQuery As String
    testQuery = "v=211&p=d&s=ta_newhigh&f=cap_largeover"

    MsgBox "Testing FINVIZ scraper with sample query..." & vbCrLf & _
           "This may take 5-10 seconds.", vbInformation

    tickers = CallPythonFinvizScraper(testQuery)

    If IsArray(tickers) And UBound(tickers) >= 0 Then
        MsgBox "✅ FINVIZ Scraper Test PASSED" & vbCrLf & vbCrLf & _
               "Found " & (UBound(tickers) + 1) & " tickers" & vbCrLf & _
               "First 5: " & Join(Array(tickers(0), tickers(1), tickers(2), tickers(3), tickers(4)), ", "), _
               vbInformation
    Else
        MsgBox "❌ FINVIZ Scraper Test FAILED" & vbCrLf & _
               "No tickers returned.", vbCritical
    End If

    ' Test 4: Heat calculator
    heatResult = CallPythonHeatCheck(75, "Tech/Comm")

    If Not IsEmpty(heatResult) Then
        MsgBox "✅ Heat Calculator Test PASSED" & vbCrLf & vbCrLf & _
               "Python heat calculations are working.", vbInformation
    Else
        MsgBox "⚠ Heat Calculator Test returned empty" & vbCrLf & _
               "Falling back to VBA calculations.", vbExclamation
    End If

    ' Summary
    MsgBox "Python Integration Tests Complete!" & vbCrLf & vbCrLf & _
           "Python Available: " & IIf(pythonOK, "✅ Yes", "❌ No") & vbCrLf & _
           "Modules Loaded: " & IIf(modulesOK, "✅ Yes", "❌ No") & vbCrLf & _
           "FINVIZ Scraper: " & IIf(IsArray(tickers), "✅ Working", "❌ Failed") & vbCrLf & _
           "Heat Calculator: " & IIf(Not IsEmpty(heatResult), "✅ Working", "⚠ Fallback to VBA"), _
           vbInformation, "Python Integration Status"

CleanExit:
    Application.ScreenUpdating = True
End Sub

' ----------------------------------------------------------------------------
' Function: GetPythonVersion
' Returns Python version string if available
' ----------------------------------------------------------------------------
Function GetPythonVersion() As String
    Dim versionCell As Range
    Dim versionFormula As String
    Dim result As Variant

    On Error Resume Next

    Set versionCell = Sheets("Control").Range("Z3")
    versionFormula = "=PY(""import sys; sys.version"")"

    versionCell.Formula = versionFormula
    Application.Calculate
    DoEvents

    result = versionCell.Value

    If Not IsEmpty(result) And result <> "" Then
        GetPythonVersion = CStr(result)
    Else
        GetPythonVersion = "Not Available"
    End If

    versionCell.ClearContents

    On Error GoTo 0
End Function
