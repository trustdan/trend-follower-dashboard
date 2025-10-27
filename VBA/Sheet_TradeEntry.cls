VERSION 1.0 CLASS
BEGIN
  MultiUse = -1  'True
END
Attribute VB_Name = "Sheet_TradeEntry"
Attribute VB_GlobalNameSpace = False
Attribute VB_Creatable = False
Attribute VB_PredeclaredId = True
Attribute VB_Exposed = True

' ========================================
' Sheet_TradeEntry Class Module
' TradeEntry sheet event handlers
' Note: This needs to be copied to the actual TradeEntry sheet's code module
' ========================================

Private Sub Worksheet_Activate()
    ' Runs when TradeEntry sheet is activated
    On Error Resume Next
    Call TF_UI.BindControls
    On Error GoTo 0
End Sub

Private Sub Worksheet_Change(ByVal Target As Range)
    ' Runs when any cell on TradeEntry changes
    On Error Resume Next

    ' If Method changes, toggle fields
    If Not Intersect(Target, Range("C13")) Is Nothing Then
        Call TF_UI.ToggleMethodFields
    End If

    ' If Sector changes, suggest bucket
    If Not Intersect(Target, Range("B7")) Is Nothing Then
        Call AutoSuggestBucket
    End If

    ' If any input changes, clear banner (reminds user to re-evaluate)
    If Not Intersect(Target, Range("B5:B11,B16:B18,C20:C25")) Is Nothing Then
        Range("A2").Value = ""
        Range("A2").Interior.ColorIndex = xlNone
        Range("A3").Value = "Inputs changed - click Evaluate to refresh"
    End If

    On Error GoTo 0
End Sub

Private Sub AutoSuggestBucket()
    ' Auto-suggest bucket based on selected sector
    Dim sector As String
    Dim suggestedBucket As String

    sector = Range("B7").Value

    Select Case sector
        Case "Technology": suggestedBucket = "Tech/Comm"
        Case "Healthcare": suggestedBucket = "Healthcare"
        Case "Financials": suggestedBucket = "Financials"
        Case "Consumer": suggestedBucket = "Consumer"
        Case "Industrials": suggestedBucket = "Industrials"
        Case "Energy": suggestedBucket = "Energy/Materials"
        Case Else: suggestedBucket = ""
    End Select

    If suggestedBucket <> "" And Range("B8").Value = "" Then
        Range("B8").Value = suggestedBucket
    End If
End Sub

Private Sub Worksheet_BeforeDoubleClick(ByVal Target As Range, Cancel As Boolean)
    ' Double-click on ticker to search
    If Not Intersect(Target, Range("B6")) Is Nothing Then
        Cancel = True
        Call SearchTicker(Target.Value)
    End If
End Sub

Private Sub SearchTicker(ticker As String)
    ' Open FINVIZ chart for ticker
    Dim url As String

    If ticker = "" Then
        Exit Sub
    End If

    url = "https://finviz.com/quote.ashx?t=" & ticker

    On Error Resume Next
    CreateObject("WScript.Shell").Run url
    On Error GoTo 0
End Sub
