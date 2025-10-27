Attribute VB_Name = "Sheet1"
Attribute VB_GlobalNameSpace = False
Attribute VB_Creatable = False
Attribute VB_PredeclaredId = True
Attribute VB_Exposed = True
' ============================================================================
' Class Module: TradeEntry Sheet
' Purpose: Event handlers for TradeEntry worksheet
' IMPORTANT: Rename this to match your TradeEntry sheet's CodeName
'            (Right-click sheet tab → View Code → see "(Name)" in Properties)
' ============================================================================

Option Explicit

' ----------------------------------------------------------------------------
' Event: Worksheet_Activate
' Runs when user switches to TradeEntry sheet
' ----------------------------------------------------------------------------
Private Sub Worksheet_Activate()
    On Error Resume Next
    Call BindControls
    On Error GoTo 0
End Sub

' ----------------------------------------------------------------------------
' Event: Worksheet_Change
' Runs when any cell is changed - handles auto-updates
' ----------------------------------------------------------------------------
Private Sub Worksheet_Change(ByVal Target As Range)
    On Error GoTo ErrorHandler

    Application.EnableEvents = False  ' Prevent infinite loops

    ' If Method choice changed (C13), toggle field visibility
    If Not Intersect(Target, Me.Range("C13")) Is Nothing Then
        Call ToggleMethodFields
    End If

    ' If Sector changed (B7), auto-suggest bucket
    If Not Intersect(Target, Me.Range("B7")) Is Nothing Then
        Dim sector As String
        sector = Me.Range("B7").Value

        ' Auto-populate bucket based on sector
        Select Case sector
            Case "Technology"
                Me.Range("B8").Value = "Tech/Comm"
            Case "Consumer"
                Me.Range("B8").Value = "Consumer"
            Case "Financials"
                Me.Range("B8").Value = "Financials"
            Case "Industrials"
                Me.Range("B8").Value = "Industrials"
            Case "Energy"
                Me.Range("B8").Value = "Energy/Materials"
            Case "Healthcare"
                Me.Range("B8").Value = "Defensives/REITs"
        End Select
    End If

    ' If any checklist item changed (C20:C25), clear banner
    If Not Intersect(Target, Me.Range("C20:C25")) Is Nothing Then
        Me.Range("A2").Value = ""
        Me.Range("A2").Interior.ColorIndex = xlNone
        Me.Range("A3").Value = "Click Evaluate to update banner"
    End If

    ' If any input changed (B5:B18), clear banner
    If Not Intersect(Target, Me.Range("B5:B18")) Is Nothing Then
        Me.Range("A2").Value = ""
        Me.Range("A2").Interior.ColorIndex = xlNone
        Me.Range("A3").Value = "Click Evaluate after making changes"
    End If

CleanExit:
    Application.EnableEvents = True
    Exit Sub

ErrorHandler:
    Application.EnableEvents = True
    ' Silently ignore errors in event handlers to avoid disrupting user
End Sub

' ----------------------------------------------------------------------------
' Event: Worksheet_SelectionChange
' Optional: Show tooltips or hints based on selected cell
' ----------------------------------------------------------------------------
Private Sub Worksheet_SelectionChange(ByVal Target As Range)
    ' Optional: Add context-sensitive help here
    ' Example: Display help text in a label when user selects specific cells
End Sub
